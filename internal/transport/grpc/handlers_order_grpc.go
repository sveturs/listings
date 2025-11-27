package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service"
)

// ============================================================================
// ORDER SERVICE gRPC METHODS
// These methods are added to the Server struct and implement OrderService RPC
// ============================================================================

// ============================================================================
// CART OPERATIONS (6 methods)
// ============================================================================

// AddToCart adds an item to shopping cart
// Creates cart if doesn't exist
// Validates: listing exists, storefront matches, stock available
func (s *Server) AddToCart(ctx context.Context, req *listingspb.AddToCartRequest) (*listingspb.AddToCartResponse, error) {
	s.logger.Debug().
		Interface("user_id", req.UserId).
		Interface("session_id", req.SessionId).
		Int64("storefront_id", req.StorefrontId).
		Int64("listing_id", req.ListingId).
		Int32("quantity", req.Quantity).
		Msg("AddToCart called")

	// Validate input
	if err := validateAddToCartRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Call service layer
	cart, err := s.cartService.AddToCart(ctx, &service.AddToCartRequest{
		UserID:       req.UserId,
		SessionID:    req.SessionId,
		StorefrontID: req.StorefrontId,
		ListingID:    req.ListingId,
		VariantID:    req.VariantId,
		Quantity:     req.Quantity,
	})

	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Cart to proto Cart
	pbCart := domainCartToProtoCart(cart)

	return &listingspb.AddToCartResponse{
		Cart:    pbCart,
		Message: "Item added to cart successfully",
	}, nil
}

// UpdateCartItem updates quantity or variant for a cart item
// Validates: ownership, stock availability
func (s *Server) UpdateCartItem(ctx context.Context, req *listingspb.UpdateCartItemRequest) (*listingspb.UpdateCartItemResponse, error) {
	s.logger.Debug().
		Int64("cart_item_id", req.CartItemId).
		Interface("user_id", req.UserId).
		Interface("quantity", req.Quantity).
		Msg("UpdateCartItem called")

	// Validate input
	if req.CartItemId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "cart_item_id must be greater than 0")
	}

	if req.Quantity != nil && *req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
	}

	// Call service layer with ownership verification built-in
	updatedCart, err := s.cartService.UpdateCartItemByItemID(ctx, req.CartItemId, *req.Quantity, req.UserId, req.SessionId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Find the updated item in the cart
	var updatedItem *domain.CartItem
	for _, item := range updatedCart.Items {
		if item.ID == req.CartItemId {
			updatedItem = item
			break
		}
	}

	// Convert domain.CartItem to proto CartItem
	pbCartItem := domainCartItemToProtoCartItem(updatedItem)

	return &listingspb.UpdateCartItemResponse{
		Item:    pbCartItem,
		Message: "Cart item updated successfully",
	}, nil
}

// RemoveFromCart removes an item from cart
// Validates: ownership
func (s *Server) RemoveFromCart(ctx context.Context, req *listingspb.RemoveFromCartRequest) (*listingspb.RemoveFromCartResponse, error) {
	s.logger.Debug().
		Int64("cart_item_id", req.CartItemId).
		Interface("user_id", req.UserId).
		Msg("RemoveFromCart called")

	// Validate input
	if req.CartItemId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "cart_item_id must be greater than 0")
	}

	// Call service layer with ownership verification built-in
	err := s.cartService.RemoveFromCartByItemID(ctx, req.CartItemId, req.UserId, req.SessionId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	return &listingspb.RemoveFromCartResponse{
		Message: "Item removed from cart successfully",
	}, nil
}

// GetCart retrieves user's cart for a storefront
// Returns: cart with items, prices, stock availability
func (s *Server) GetCart(ctx context.Context, req *listingspb.GetCartRequest) (*listingspb.GetCartResponse, error) {
	s.logger.Debug().
		Interface("user_id", req.UserId).
		Interface("session_id", req.SessionId).
		Int64("storefront_id", req.StorefrontId).
		Msg("GetCart called")

	// Validate input
	if (req.UserId == nil && req.SessionId == nil) || (req.UserId != nil && req.SessionId != nil) {
		return nil, status.Error(codes.InvalidArgument, "must provide either user_id or session_id (not both)")
	}

	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id must be greater than 0")
	}

	// Call service layer
	cart, err := s.cartService.GetCart(ctx, req.UserId, req.SessionId, req.StorefrontId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Cart to proto Cart
	pbCart := domainCartToProtoCart(cart)

	// Calculate cart summary
	summary := calculateCartSummary(cart)

	return &listingspb.GetCartResponse{
		Cart:    pbCart,
		Summary: summary,
	}, nil
}

// ClearCart removes all items from cart
// Used after order creation or manual clear
func (s *Server) ClearCart(ctx context.Context, req *listingspb.ClearCartRequest) (*emptypb.Empty, error) {
	s.logger.Debug().
		Interface("user_id", req.UserId).
		Interface("session_id", req.SessionId).
		Int64("storefront_id", req.StorefrontId).
		Msg("ClearCart called")

	// Validate input
	if (req.UserId == nil && req.SessionId == nil) || (req.UserId != nil && req.SessionId != nil) {
		return nil, status.Error(codes.InvalidArgument, "must provide either user_id or session_id (not both)")
	}

	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id must be greater than 0")
	}

	// First, get the cart to get its ID
	cart, err := s.cartService.GetCart(ctx, req.UserId, req.SessionId, req.StorefrontId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Clear cart
	if err := s.cartService.ClearCart(ctx, cart.ID); err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	return &emptypb.Empty{}, nil
}

// GetUserCarts retrieves all carts for authenticated user
// Returns: carts across all storefronts
func (s *Server) GetUserCarts(ctx context.Context, req *listingspb.GetUserCartsRequest) (*listingspb.GetUserCartsResponse, error) {
	s.logger.Debug().
		Int64("user_id", req.UserId).
		Msg("GetUserCarts called")

	// Validate input
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be greater than 0")
	}

	// Call service layer
	carts, err := s.cartService.GetUserCarts(ctx, req.UserId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Cart slice to proto Cart slice
	pbCarts := make([]*listingspb.Cart, 0, len(carts))
	for _, cart := range carts {
		pbCarts = append(pbCarts, domainCartToProtoCart(cart))
	}

	return &listingspb.GetUserCartsResponse{
		Carts:      pbCarts,
		TotalCarts: int32(len(pbCarts)),
	}, nil
}

// ============================================================================
// ORDER OPERATIONS (6 methods)
// ============================================================================

// CreateOrder creates a new order from cart OR direct items
// Transaction: validates cart/items, creates order, reserves inventory, deducts stock, clears cart
func (s *Server) CreateOrder(ctx context.Context, req *listingspb.CreateOrderRequest) (*listingspb.CreateOrderResponse, error) {
	s.logger.Info().
		Interface("user_id", req.UserId).
		Interface("cart_id", req.CartId).
		Int("items_count", len(req.Items)).
		Str("payment_method", req.PaymentMethod).
		Msg("CreateOrder called")

	// Validate input
	if err := validateCreateOrderRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto Struct to map[string]interface{}
	shippingAddress := protoStructToMap(req.ShippingAddress)
	billingAddress := protoStructToMap(req.BillingAddress)

	// Convert proto Items to service.OrderItemInput
	var items []service.OrderItemInput
	if len(req.Items) > 0 {
		items = make([]service.OrderItemInput, 0, len(req.Items))
		for _, protoItem := range req.Items {
			item := service.OrderItemInput{
				ProductID: protoItem.ProductId,
				Quantity:  int(protoItem.Quantity),
			}
			if protoItem.VariantId != nil {
				variantID := *protoItem.VariantId
				item.VariantID = &variantID
			}
			items = append(items, item)
		}
	}

	// Determine cart_id (0 if using direct items)
	cartID := int64(0)
	if req.CartId != nil {
		cartID = *req.CartId
	}

	// Call service layer
	order, err := s.orderService.CreateOrder(ctx, &service.CreateOrderRequest{
		CartID:             cartID,
		Items:              items,
		StorefrontID:       req.StorefrontId,
		UserID:             req.UserId,
		ShippingAddress:    shippingAddress,
		BillingAddress:     billingAddress,
		ShippingCost:       0, // TODO: Calculate shipping cost
		DiscountCode:       nil,
		DiscountAmount:     0,
		PaymentMethod:      req.PaymentMethod,
		CustomerNotes:      req.CustomerNotes,
		AcceptPriceChanges: req.AcceptPriceChanges,
		CustomerName:       req.CustomerName,
		CustomerEmail:      req.CustomerEmail,
		CustomerPhone:      req.CustomerPhone,
	})

	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &listingspb.CreateOrderResponse{
		Order:   pbOrder,
		Message: "Order created successfully",
	}, nil
}

// GetOrder retrieves a single order by ID
// Validates: ownership (user can only see own orders)
func (s *Server) GetOrder(ctx context.Context, req *listingspb.GetOrderRequest) (*listingspb.GetOrderResponse, error) {
	s.logger.Debug().
		Int64("order_id", req.OrderId).
		Interface("user_id", req.UserId).
		Msg("GetOrder called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}

	// Call service layer
	order, err := s.orderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Verify ownership (if user_id provided)
	if req.UserId != nil && order.UserID != nil && *req.UserId != *order.UserID {
		s.logger.Warn().
			Int64("order_id", req.OrderId).
			Int64("requesting_user_id", *req.UserId).
			Int64("order_user_id", *order.UserID).
			Msg("unauthorized access attempt")
		return nil, status.Error(codes.PermissionDenied, "you don't have permission to view this order")
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &listingspb.GetOrderResponse{
		Order: pbOrder,
	}, nil
}

// ListOrders retrieves orders with filters and pagination
// Admin: can see all orders
// User: can only see own orders
func (s *Server) ListOrders(ctx context.Context, req *listingspb.ListOrdersRequest) (*listingspb.ListOrdersResponse, error) {
	s.logger.Debug().
		Interface("user_id", req.UserId).
		Interface("storefront_id", req.StorefrontId).
		Int32("page", req.Page).
		Int32("page_size", req.PageSize).
		Msg("ListOrders called")

	// Validate pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Call service layer
	orders, totalCount, err := s.orderService.ListOrders(ctx, &service.ListOrdersRequest{
		UserID:       req.UserId,
		StorefrontID: req.StorefrontId,
		Status:       domainOrderStatusFromProto(req.Status),
		Limit:        int(pageSize),
		Offset:       int(offset),
	})

	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order slice to proto Order slice
	pbOrders := make([]*listingspb.Order, 0, len(orders))
	for _, order := range orders {
		pbOrders = append(pbOrders, domainOrderToProtoOrder(order))
	}

	return &listingspb.ListOrdersResponse{
		Orders:     pbOrders,
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// CancelOrder cancels an order and releases inventory
// Validates: order status (only pending/confirmed allowed)
// Actions: update status, release reservations, restore stock, publish event
func (s *Server) CancelOrder(ctx context.Context, req *listingspb.CancelOrderRequest) (*listingspb.CancelOrderResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Interface("user_id", req.UserId).
		Bool("refund", req.Refund).
		Msg("CancelOrder called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}

	// For user cancellations, require user_id
	if req.UserId == nil {
		return nil, status.Error(codes.InvalidArgument, "user_id is required for order cancellation")
	}

	// Prepare reason
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	// Call service layer
	order, err := s.orderService.CancelOrder(ctx, req.OrderId, *req.UserId, reason)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &listingspb.CancelOrderResponse{
		Order:           pbOrder,
		Message:         "Order cancelled successfully",
		RefundInitiated: req.Refund,
	}, nil
}

// UpdateOrderStatus updates order status (admin only)
// Validates: status transition rules (state machine)
// On status = shipped: requires tracking_number
// On status = delivered: sets escrow_release_date
func (s *Server) UpdateOrderStatus(ctx context.Context, req *listingspb.UpdateOrderStatusRequest) (*listingspb.UpdateOrderStatusResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Str("new_status", req.NewStatus.String()).
		Msg("UpdateOrderStatus called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}

	// Validate status
	if req.NewStatus == listingspb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "new_status must be specified")
	}

	// Convert proto status to domain status
	domainStatus := domainOrderStatusFromProto(&req.NewStatus)
	if domainStatus == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order status")
	}

	// Call service layer
	order, err := s.orderService.UpdateOrderStatus(ctx, req.OrderId, *domainStatus)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &listingspb.UpdateOrderStatusResponse{
		Order:   pbOrder,
		Message: "Order status updated successfully",
	}, nil
}

// GetOrderStats retrieves order statistics (admin)
// Returns: aggregated stats, status breakdown, daily trends
func (s *Server) GetOrderStats(ctx context.Context, req *listingspb.GetOrderStatsRequest) (*listingspb.GetOrderStatsResponse, error) {
	s.logger.Debug().
		Interface("storefront_id", req.StorefrontId).
		Msg("GetOrderStats called")

	// Call service layer
	stats, err := s.orderService.GetOrderStats(ctx, nil, req.StorefrontId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain stats to proto stats
	pbStats := &listingspb.OrderStatsSummary{
		TotalOrders:     int32(stats.TotalOrders),
		PendingOrders:   int32(stats.PendingOrders),
		ConfirmedOrders: int32(stats.ConfirmedOrders),
		TotalRevenue:    stats.TotalRevenue,
	}

	return &listingspb.GetOrderStatsResponse{
		Stats: pbStats,
	}, nil
}

// ============================================================================
// HELPER FUNCTIONS - Validation
// ============================================================================

// validateAddToCartRequest validates AddToCart request
func validateAddToCartRequest(req *listingspb.AddToCartRequest) error {
	if (req.UserId == nil && req.SessionId == nil) || (req.UserId != nil && req.SessionId != nil) {
		return errors.New("must provide either user_id or session_id (not both)")
	}

	if req.StorefrontId <= 0 {
		return errors.New("storefront_id must be greater than 0")
	}

	if req.ListingId <= 0 {
		return errors.New("listing_id must be greater than 0")
	}

	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	return nil
}

// validateCreateOrderRequest validates CreateOrder request
func validateCreateOrderRequest(req *listingspb.CreateOrderRequest) error {
	// Validate that either cart_id or items is provided
	hasCartID := req.CartId != nil && *req.CartId > 0
	hasItems := len(req.Items) > 0

	if !hasCartID && !hasItems {
		return errors.New("either cart_id or items must be provided")
	}

	// Validate items if using direct checkout
	if hasItems {
		for i, item := range req.Items {
			if item.ProductId <= 0 {
				return fmt.Errorf("items[%d].product_id must be greater than 0", i)
			}
			if item.Quantity <= 0 {
				return fmt.Errorf("items[%d].quantity must be greater than 0", i)
			}
		}
	}

	if req.StorefrontId <= 0 {
		return errors.New("storefront_id must be greater than 0")
	}

	if req.ShippingAddress == nil {
		return errors.New("shipping_address is required")
	}

	if req.ShippingMethod == "" {
		return errors.New("shipping_method is required")
	}

	if req.PaymentMethod == "" {
		return errors.New("payment_method is required")
	}

	// Validate payment method is one of allowed values
	if !isValidPaymentMethod(req.PaymentMethod) {
		return fmt.Errorf("invalid payment_method: %s (allowed: card, cash_on_delivery, cod, cash, bank_transfer)", req.PaymentMethod)
	}

	return nil
}

// isValidPaymentMethod checks if payment method is valid
func isValidPaymentMethod(method string) bool {
	switch method {
	case "card", "cash_on_delivery", "cod", "cash", "bank_transfer", "pouzecem", "pouzećem":
		return true
	default:
		return false
	}
}

// ============================================================================
// HELPER FUNCTIONS - Mappers
// ============================================================================

// domainCartToProtoCart converts domain.Cart to proto Cart
func domainCartToProtoCart(cart *domain.Cart) *listingspb.Cart {
	if cart == nil {
		return nil
	}

	pbCart := &listingspb.Cart{
		Id:           cart.ID,
		UserId:       cart.UserID,
		SessionId:    cart.SessionID,
		StorefrontId: cart.StorefrontID,
		CreatedAt:    timestamppb.New(cart.CreatedAt),
		UpdatedAt:    timestamppb.New(cart.UpdatedAt),
	}

	// Convert cart items
	if len(cart.Items) > 0 {
		pbCart.Items = make([]*listingspb.CartItem, 0, len(cart.Items))
		for _, item := range cart.Items {
			pbCart.Items = append(pbCart.Items, domainCartItemToProtoCartItem(item))
		}
	}

	return pbCart
}

// domainCartItemToProtoCartItem converts domain.CartItem to proto CartItem
// Uses domain's ToProto method to include embedded data (ListingName, ListingImage, etc.)
func domainCartItemToProtoCartItem(item *domain.CartItem) *listingspb.CartItem {
	if item == nil {
		return nil
	}

	// Use domain's ToProto method which includes embedded data fields
	return item.ToProto()
}

// domainOrderToProtoOrder converts domain.Order to proto Order
func domainOrderToProtoOrder(order *domain.Order) *listingspb.Order {
	if order == nil {
		return nil
	}

	pbOrder := &listingspb.Order{
		Id:           order.ID,
		OrderNumber:  order.OrderNumber,
		UserId:       order.UserID,
		StorefrontId: order.StorefrontID,
		Status:       protoOrderStatusFromDomain(order.Status),
		Financials: &listingspb.OrderFinancials{
			Subtotal:     order.Subtotal,
			Tax:          order.Tax,
			ShippingCost: order.Shipping,
			Discount:     order.Discount,
			Total:        order.Total,
			Commission:   order.Commission,
			SellerAmount: order.SellerAmount,
			Currency:     order.Currency,
		},
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   protoPaymentStatusFromDomain(order.PaymentStatus),
		ShippingAddress: mapToProtoStruct(order.ShippingAddress),
		BillingAddress:  mapToProtoStruct(order.BillingAddress),
		EscrowDays:      order.EscrowDays,
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}

	// Optional fields
	if order.PaymentTransactionID != nil {
		pbOrder.PaymentTransactionId = order.PaymentTransactionID
	}

	if order.PaymentCompletedAt != nil {
		pbOrder.PaymentCompletedAt = timestamppb.New(*order.PaymentCompletedAt)
	}

	if order.ShippingMethod != nil {
		pbOrder.ShippingMethod = order.ShippingMethod
	}

	if order.TrackingNumber != nil {
		pbOrder.TrackingNumber = order.TrackingNumber
	}

	if order.EscrowReleaseDate != nil {
		pbOrder.EscrowReleaseDate = timestamppb.New(*order.EscrowReleaseDate)
	}

	if order.CustomerNotes != nil {
		pbOrder.CustomerNotes = order.CustomerNotes
	}

	if order.AdminNotes != nil {
		pbOrder.AdminNotes = order.AdminNotes
	}

	if order.ConfirmedAt != nil {
		pbOrder.ConfirmedAt = timestamppb.New(*order.ConfirmedAt)
	}

	if order.ShippedAt != nil {
		pbOrder.ShippedAt = timestamppb.New(*order.ShippedAt)
	}

	if order.DeliveredAt != nil {
		pbOrder.DeliveredAt = timestamppb.New(*order.DeliveredAt)
	}

	if order.CancelledAt != nil {
		pbOrder.CancelledAt = timestamppb.New(*order.CancelledAt)
	}

	// Shipment workflow fields
	if order.AcceptedAt != nil {
		pbOrder.AcceptedAt = timestamppb.New(*order.AcceptedAt)
	}

	if order.SellerNotes != nil {
		pbOrder.SellerNotes = order.SellerNotes
	}

	if order.LabelURL != nil {
		pbOrder.LabelUrl = order.LabelURL
	}

	if order.ShipmentID != nil {
		pbOrder.ShipmentId = order.ShipmentID
	}

	if order.ShippingProvider != nil {
		pbOrder.ShippingProvider = order.ShippingProvider
	}

	// Storefront name (for frontend display)
	if order.StorefrontName != nil {
		pbOrder.StorefrontName = order.StorefrontName
	}

	// Customer info (for seller's sales page display)
	if order.CustomerName != nil {
		pbOrder.CustomerName = order.CustomerName
	}
	if order.CustomerEmail != nil {
		pbOrder.CustomerEmail = order.CustomerEmail
	}
	if order.CustomerPhone != nil {
		pbOrder.CustomerPhone = order.CustomerPhone
	}

	// Convert order items
	if len(order.Items) > 0 {
		pbOrder.Items = make([]*listingspb.OrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			pbOrder.Items = append(pbOrder.Items, domainOrderItemToProtoOrderItem(item))
		}
	}

	return pbOrder
}

// domainOrderItemToProtoOrderItem converts domain.OrderItem to proto OrderItem
func domainOrderItemToProtoOrderItem(item *domain.OrderItem) *listingspb.OrderItem {
	if item == nil {
		return nil
	}

	pbItem := &listingspb.OrderItem{
		Id:          item.ID,
		OrderId:     item.OrderID,
		ListingId:   item.ListingID,
		VariantId:   item.VariantID,
		ListingName: item.ListingName,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		Subtotal:    item.Subtotal,
		Discount:    item.Discount,
		Total:       item.Total,
		CreatedAt:   timestamppb.New(item.CreatedAt),
	}

	// Add image URL if present
	if item.ImageURL != nil {
		pbItem.ImageUrl = item.ImageURL
	}

	return pbItem
}

// calculateCartSummary calculates cart summary (totals, warnings)
func calculateCartSummary(cart *domain.Cart) *listingspb.CartSummary {
	if cart == nil || len(cart.Items) == 0 {
		return &listingspb.CartSummary{
			TotalItems: 0,
			Subtotal:   0,
		}
	}

	var totalItems int32
	var subtotal float64

	for _, item := range cart.Items {
		totalItems += item.Quantity
		subtotal += float64(item.Quantity) * item.PriceSnapshot
	}

	return &listingspb.CartSummary{
		TotalItems:        totalItems,
		Subtotal:          subtotal,
		EstimatedTotal:    subtotal, // TODO: Add tax and shipping calculation
		EstimatedTax:      0,
		EstimatedShipping: 0,
	}
}

// protoOrderStatusFromDomain converts domain.OrderStatus to proto OrderStatus
func protoOrderStatusFromDomain(status domain.OrderStatus) listingspb.OrderStatus {
	switch status {
	case domain.OrderStatusPending:
		return listingspb.OrderStatus_ORDER_STATUS_PENDING
	case domain.OrderStatusConfirmed:
		return listingspb.OrderStatus_ORDER_STATUS_CONFIRMED
	case domain.OrderStatusAccepted:
		return listingspb.OrderStatus_ORDER_STATUS_ACCEPTED
	case domain.OrderStatusProcessing:
		return listingspb.OrderStatus_ORDER_STATUS_PROCESSING
	case domain.OrderStatusShipped:
		return listingspb.OrderStatus_ORDER_STATUS_SHIPPED
	case domain.OrderStatusDelivered:
		return listingspb.OrderStatus_ORDER_STATUS_DELIVERED
	case domain.OrderStatusCancelled:
		return listingspb.OrderStatus_ORDER_STATUS_CANCELLED
	case domain.OrderStatusRefunded:
		return listingspb.OrderStatus_ORDER_STATUS_REFUNDED
	case domain.OrderStatusFailed:
		return listingspb.OrderStatus_ORDER_STATUS_FAILED
	default:
		return listingspb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

// domainOrderStatusFromProto converts proto OrderStatus to domain.OrderStatus
func domainOrderStatusFromProto(status *listingspb.OrderStatus) *domain.OrderStatus {
	if status == nil {
		return nil
	}

	var domainStatus domain.OrderStatus
	switch *status {
	case listingspb.OrderStatus_ORDER_STATUS_PENDING:
		domainStatus = domain.OrderStatusPending
	case listingspb.OrderStatus_ORDER_STATUS_CONFIRMED:
		domainStatus = domain.OrderStatusConfirmed
	case listingspb.OrderStatus_ORDER_STATUS_ACCEPTED:
		domainStatus = domain.OrderStatusAccepted
	case listingspb.OrderStatus_ORDER_STATUS_PROCESSING:
		domainStatus = domain.OrderStatusProcessing
	case listingspb.OrderStatus_ORDER_STATUS_SHIPPED:
		domainStatus = domain.OrderStatusShipped
	case listingspb.OrderStatus_ORDER_STATUS_DELIVERED:
		domainStatus = domain.OrderStatusDelivered
	case listingspb.OrderStatus_ORDER_STATUS_CANCELLED:
		domainStatus = domain.OrderStatusCancelled
	case listingspb.OrderStatus_ORDER_STATUS_REFUNDED:
		domainStatus = domain.OrderStatusRefunded
	case listingspb.OrderStatus_ORDER_STATUS_FAILED:
		domainStatus = domain.OrderStatusFailed
	default:
		return nil
	}

	return &domainStatus
}

// protoPaymentStatusFromDomain converts domain.PaymentStatus to proto PaymentStatus
func protoPaymentStatusFromDomain(status domain.PaymentStatus) listingspb.PaymentStatus {
	switch status {
	case domain.PaymentStatusPending:
		return listingspb.PaymentStatus_PAYMENT_STATUS_PENDING
	case domain.PaymentStatusProcessing:
		return listingspb.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case domain.PaymentStatusCompleted:
		return listingspb.PaymentStatus_PAYMENT_STATUS_COMPLETED
	case domain.PaymentStatusCODPending:
		return listingspb.PaymentStatus_PAYMENT_STATUS_COD_PENDING
	case domain.PaymentStatusFailed:
		return listingspb.PaymentStatus_PAYMENT_STATUS_FAILED
	case domain.PaymentStatusRefunded:
		return listingspb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return listingspb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

// protoStructToMap converts proto Struct to map[string]interface{}
func protoStructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// mapToProtoStruct converts map[string]interface{} to proto Struct
func mapToProtoStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}

	s, err := structpb.NewStruct(m)
	if err != nil {
		// Log error but return nil instead of panicking
		return nil
	}

	return s
}

// ============================================================================
// HELPER FUNCTIONS - Error Mapping
// ============================================================================

// ============================================================================
// SHIPMENT OPERATIONS (4 methods)
// ============================================================================

// AcceptOrder handles seller accepting an order
// Flow: confirmed -> accepted
// Validates: order status, seller is storefront owner
func (s *Server) AcceptOrder(ctx context.Context, req *listingspb.AcceptOrderRequest) (*listingspb.AcceptOrderResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Int64("seller_id", req.SellerId).
		Msg("AcceptOrder called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}
	if req.SellerId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "seller_id must be greater than 0")
	}

	// Get seller notes (optional)
	var sellerNotes string
	if req.SellerNotes != nil {
		sellerNotes = *req.SellerNotes
	}

	// Call service layer
	order, err := s.orderService.AcceptOrder(ctx, req.OrderId, req.SellerId, sellerNotes)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &listingspb.AcceptOrderResponse{
		Order:   pbOrder,
		Message: "Order accepted successfully",
	}, nil
}

// CreateOrderShipment creates shipment via Delivery Service
// Flow: accepted -> processing
// Validates: order status, seller is storefront owner, package info
func (s *Server) CreateOrderShipment(ctx context.Context, req *listingspb.CreateOrderShipmentRequest) (*listingspb.CreateOrderShipmentResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Int64("seller_id", req.SellerId).
		Str("provider_code", req.ProviderCode).
		Msg("CreateOrderShipment called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}
	if req.SellerId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "seller_id must be greater than 0")
	}
	if req.ProviderCode == "" {
		return nil, status.Error(codes.InvalidArgument, "provider_code is required")
	}

	// Validate package info
	if req.PackageInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "package_info is required")
	}
	if req.PackageInfo.WeightKg <= 0 {
		return nil, status.Error(codes.InvalidArgument, "package weight must be greater than 0")
	}

	// Build service request
	shipmentReq := &service.CreateShipmentRequest{
		OrderID:      req.OrderId,
		SellerID:     req.SellerId,
		ProviderCode: req.ProviderCode,
		PackageInfo: service.PackageInfo{
			WeightKg:    req.PackageInfo.WeightKg,
			LengthCm:    req.PackageInfo.LengthCm,
			WidthCm:     req.PackageInfo.WidthCm,
			HeightCm:    req.PackageInfo.HeightCm,
			IsFragile:   req.PackageInfo.IsFragile,
			Description: req.PackageInfo.Description,
		},
		UseCOD:         req.UseCod,
		CODAmount:      req.CodAmount,
		UseInsurance:   req.UseInsurance,
		InsuranceValue: req.InsuranceValue,
	}

	// Call service layer
	result, err := s.orderService.CreateOrderShipment(ctx, shipmentReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(result.Order)

	// Build shipment info
	pbShipment := &listingspb.ShipmentInfo{
		ShipmentId:     result.ShipmentID,
		TrackingNumber: result.TrackingNumber,
		Provider:       result.Provider,
		Status:         result.Status,
		DeliveryCost:   result.DeliveryCost,
	}
	if result.EstimatedDelivery != "" {
		pbShipment.EstimatedDelivery = &result.EstimatedDelivery
	}

	return &listingspb.CreateOrderShipmentResponse{
		Order:    pbOrder,
		Shipment: pbShipment,
		LabelUrl: result.LabelURL,
		Message:  "Shipment created successfully",
	}, nil
}

// MarkOrderShipped marks order as shipped
// Flow: processing -> shipped
// Validates: order status, tracking number exists
func (s *Server) MarkOrderShipped(ctx context.Context, req *listingspb.MarkOrderShippedRequest) (*listingspb.MarkOrderShippedResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Int64("seller_id", req.SellerId).
		Msg("MarkOrderShipped called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}
	if req.SellerId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "seller_id must be greater than 0")
	}

	// Get seller notes (optional)
	var sellerNotes string
	if req.SellerNotes != nil {
		sellerNotes = *req.SellerNotes
	}

	// Call service layer
	order, err := s.orderService.MarkOrderShipped(ctx, req.OrderId, req.SellerId, sellerNotes)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &listingspb.MarkOrderShippedResponse{
		Order:   pbOrder,
		Message: "Order marked as shipped successfully",
	}, nil
}

// GetOrderTracking gets tracking info from Delivery Service
// Returns: tracking events timeline
func (s *Server) GetOrderTracking(ctx context.Context, req *listingspb.GetOrderTrackingRequest) (*listingspb.GetOrderTrackingResponse, error) {
	s.logger.Debug().
		Int64("order_id", req.OrderId).
		Int64("user_id", req.UserId).
		Msg("GetOrderTracking called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be greater than 0")
	}

	// Call service layer
	tracking, err := s.orderService.GetOrderTracking(ctx, req.OrderId, req.UserId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Build response
	resp := &listingspb.GetOrderTrackingResponse{
		TrackingNumber: tracking.TrackingNumber,
		Provider:       tracking.Provider,
		Status:         tracking.Status,
	}

	if tracking.EstimatedDelivery != "" {
		resp.EstimatedDelivery = &tracking.EstimatedDelivery
	}

	// Convert tracking events
	if len(tracking.Events) > 0 {
		resp.Events = make([]*listingspb.TrackingEvent, 0, len(tracking.Events))
		for _, event := range tracking.Events {
			pbEvent := &listingspb.TrackingEvent{
				Status:      event.Status,
				Location:    event.Location,
				Description: event.Description,
				Timestamp:   timestamppb.New(event.Timestamp),
			}
			resp.Events = append(resp.Events, pbEvent)
		}
	}

	return resp, nil
}

// ============================================================================
// HELPER FUNCTIONS - Error Mapping
// ============================================================================

// mapServiceErrorToGRPC maps service layer errors to gRPC status codes
func mapServiceErrorToGRPC(err error, logger zerolog.Logger) error {
	if err == nil {
		return nil
	}

	// Log the original error
	logger.Error().Err(err).Msg("service error occurred")

	// Check for specific error types
	if service.IsNotFoundError(err) {
		return status.Error(codes.NotFound, err.Error())
	}

	if service.IsValidationError(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if service.IsConflictError(err) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	// Check for unauthorized errors
	if errors.Is(err, service.ErrUnauthorized) {
		return status.Error(codes.PermissionDenied, "unauthorized")
	}

	// Check for chat-specific errors
	if errors.Is(err, service.ErrNotParticipant) {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	if errors.Is(err, service.ErrChatBlocked) {
		return status.Error(codes.FailedPrecondition, err.Error())
	}

	// Check for attachment validation errors (custom error types)
	var attachmentTooLargeErr *service.ErrAttachmentTooLarge
	if errors.As(err, &attachmentTooLargeErr) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	var invalidFileTypeErr *service.ErrInvalidFileType
	if errors.As(err, &invalidFileTypeErr) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// Default to internal error
	return status.Error(codes.Internal, "internal server error")
}
