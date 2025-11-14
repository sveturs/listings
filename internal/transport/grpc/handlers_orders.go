package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	ordersspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/service"
)

// OrderServiceServer implements gRPC OrderServiceServer interface
// This struct handles all cart and order operations defined in orders.proto
type OrderServiceServer struct {
	ordersspb.UnimplementedOrderServiceServer
	cartService      service.CartService
	orderService     service.OrderService
	inventoryService service.InventoryService
	logger           zerolog.Logger
}

// NewOrderServiceServer creates a new order service gRPC handler
func NewOrderServiceServer(
	cartService service.CartService,
	orderService service.OrderService,
	inventoryService service.InventoryService,
	logger zerolog.Logger,
) *OrderServiceServer {
	return &OrderServiceServer{
		cartService:      cartService,
		orderService:     orderService,
		inventoryService: inventoryService,
		logger:           logger.With().Str("component", "grpc_orders_handler").Logger(),
	}
}

// ============================================================================
// CART OPERATIONS (6 methods)
// ============================================================================

// AddToCart adds an item to shopping cart
// Creates cart if doesn't exist
// Validates: listing exists, storefront matches, stock available
func (s *OrderServiceServer) AddToCart(ctx context.Context, req *ordersspb.AddToCartRequest) (*ordersspb.AddToCartResponse, error) {
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

	return &ordersspb.AddToCartResponse{
		Cart:    pbCart,
		Message: "Item added to cart successfully",
	}, nil
}

// UpdateCartItem updates quantity or variant for a cart item
// Validates: ownership, stock availability
func (s *OrderServiceServer) UpdateCartItem(ctx context.Context, req *ordersspb.UpdateCartItemRequest) (*ordersspb.UpdateCartItemResponse, error) {
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

	// TODO: Get cart_id from cart_item_id (need to query CartRepository)
	// For now, this is a placeholder implementation
	// In production, we need to:
	// 1. Get cart_item to find cart_id
	// 2. Verify ownership (user_id or session_id matches)
	// 3. Call UpdateCartItem service method

	// Placeholder: return error for now
	return nil, status.Error(codes.Unimplemented, "UpdateCartItem not fully implemented yet")
}

// RemoveFromCart removes an item from cart
// Validates: ownership
func (s *OrderServiceServer) RemoveFromCart(ctx context.Context, req *ordersspb.RemoveFromCartRequest) (*ordersspb.RemoveFromCartResponse, error) {
	s.logger.Debug().
		Int64("cart_item_id", req.CartItemId).
		Interface("user_id", req.UserId).
		Msg("RemoveFromCart called")

	// Validate input
	if req.CartItemId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "cart_item_id must be greater than 0")
	}

	// TODO: Similar to UpdateCartItem, need to get cart_id from cart_item_id
	// Placeholder implementation
	return nil, status.Error(codes.Unimplemented, "RemoveFromCart not fully implemented yet")
}

// GetCart retrieves user's cart for a storefront
// Returns: cart with items, prices, stock availability
func (s *OrderServiceServer) GetCart(ctx context.Context, req *ordersspb.GetCartRequest) (*ordersspb.GetCartResponse, error) {
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

	return &ordersspb.GetCartResponse{
		Cart:    pbCart,
		Summary: summary,
	}, nil
}

// ClearCart removes all items from cart
// Used after order creation or manual clear
func (s *OrderServiceServer) ClearCart(ctx context.Context, req *ordersspb.ClearCartRequest) (*emptypb.Empty, error) {
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
func (s *OrderServiceServer) GetUserCarts(ctx context.Context, req *ordersspb.GetUserCartsRequest) (*ordersspb.GetUserCartsResponse, error) {
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
	pbCarts := make([]*ordersspb.Cart, 0, len(carts))
	for _, cart := range carts {
		pbCarts = append(pbCarts, domainCartToProtoCart(cart))
	}

	return &ordersspb.GetUserCartsResponse{
		Carts:      pbCarts,
		TotalCarts: int32(len(pbCarts)),
	}, nil
}

// ============================================================================
// ORDER OPERATIONS (6 methods)
// ============================================================================

// CreateOrder creates a new order from cart
// Transaction: validates cart, creates order, reserves inventory, deducts stock, clears cart
func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *ordersspb.CreateOrderRequest) (*ordersspb.CreateOrderResponse, error) {
	s.logger.Info().
		Interface("user_id", req.UserId).
		Int64("cart_id", req.CartId).
		Str("payment_method", req.PaymentMethod).
		Msg("CreateOrder called")

	// Validate input
	if err := validateCreateOrderRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto Struct to map[string]interface{}
	shippingAddress := protoStructToMap(req.ShippingAddress)
	billingAddress := protoStructToMap(req.BillingAddress)

	// Call service layer
	order, err := s.orderService.CreateOrder(ctx, &service.CreateOrderRequest{
		CartID:          req.CartId,
		UserID:          req.UserId,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		ShippingCost:    0, // TODO: Calculate shipping cost
		DiscountCode:    nil,
		DiscountAmount:  0,
		PaymentMethod:   req.PaymentMethod,
		CustomerNotes:   req.CustomerNotes,
	})

	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Order to proto Order
	pbOrder := domainOrderToProtoOrder(order)

	return &ordersspb.CreateOrderResponse{
		Order:   pbOrder,
		Message: "Order created successfully",
	}, nil
}

// GetOrder retrieves a single order by ID
// Validates: ownership (user can only see own orders)
func (s *OrderServiceServer) GetOrder(ctx context.Context, req *ordersspb.GetOrderRequest) (*ordersspb.GetOrderResponse, error) {
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

	return &ordersspb.GetOrderResponse{
		Order: pbOrder,
	}, nil
}

// ListOrders retrieves orders with filters and pagination
// Admin: can see all orders
// User: can only see own orders
func (s *OrderServiceServer) ListOrders(ctx context.Context, req *ordersspb.ListOrdersRequest) (*ordersspb.ListOrdersResponse, error) {
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
	pbOrders := make([]*ordersspb.Order, 0, len(orders))
	for _, order := range orders {
		pbOrders = append(pbOrders, domainOrderToProtoOrder(order))
	}

	return &ordersspb.ListOrdersResponse{
		Orders:     pbOrders,
		TotalCount: int32(totalCount),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// CancelOrder cancels an order and releases inventory
// Validates: order status (only pending/confirmed allowed)
// Actions: update status, release reservations, restore stock, publish event
func (s *OrderServiceServer) CancelOrder(ctx context.Context, req *ordersspb.CancelOrderRequest) (*ordersspb.CancelOrderResponse, error) {
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

	return &ordersspb.CancelOrderResponse{
		Order:           pbOrder,
		Message:         "Order cancelled successfully",
		RefundInitiated: req.Refund,
	}, nil
}

// UpdateOrderStatus updates order status (admin only)
// Validates: status transition rules (state machine)
// On status = shipped: requires tracking_number
// On status = delivered: sets escrow_release_date
func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, req *ordersspb.UpdateOrderStatusRequest) (*ordersspb.UpdateOrderStatusResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Str("new_status", req.NewStatus.String()).
		Msg("UpdateOrderStatus called")

	// Validate input
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be greater than 0")
	}

	// Validate status
	if req.NewStatus == ordersspb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
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

	return &ordersspb.UpdateOrderStatusResponse{
		Order:   pbOrder,
		Message: "Order status updated successfully",
	}, nil
}

// GetOrderStats retrieves order statistics (admin)
// Returns: aggregated stats, status breakdown, daily trends
func (s *OrderServiceServer) GetOrderStats(ctx context.Context, req *ordersspb.GetOrderStatsRequest) (*ordersspb.GetOrderStatsResponse, error) {
	s.logger.Debug().
		Interface("storefront_id", req.StorefrontId).
		Msg("GetOrderStats called")

	// Call service layer
	stats, err := s.orderService.GetOrderStats(ctx, nil, req.StorefrontId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain stats to proto stats
	pbStats := &ordersspb.OrderStatsSummary{
		TotalOrders:     int32(stats.TotalOrders),
		PendingOrders:   int32(stats.PendingOrders),
		ConfirmedOrders: int32(stats.ConfirmedOrders),
		TotalRevenue:    stats.TotalRevenue,
	}

	return &ordersspb.GetOrderStatsResponse{
		Stats: pbStats,
	}, nil
}

// ============================================================================
// HELPER FUNCTIONS - Validation
// ============================================================================

// validateAddToCartRequest validates AddToCart request
func validateAddToCartRequest(req *ordersspb.AddToCartRequest) error {
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
func validateCreateOrderRequest(req *ordersspb.CreateOrderRequest) error {
	if req.CartId <= 0 {
		return errors.New("cart_id must be greater than 0")
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

	return nil
}

// ============================================================================
// HELPER FUNCTIONS - Mappers
// ============================================================================

// domainCartToProtoCart converts domain.Cart to proto Cart
func domainCartToProtoCart(cart *domain.Cart) *ordersspb.Cart {
	if cart == nil {
		return nil
	}

	pbCart := &ordersspb.Cart{
		Id:           cart.ID,
		UserId:       cart.UserID,
		SessionId:    cart.SessionID,
		StorefrontId: cart.StorefrontID,
		CreatedAt:    timestamppb.New(cart.CreatedAt),
		UpdatedAt:    timestamppb.New(cart.UpdatedAt),
	}

	// Convert cart items
	if len(cart.Items) > 0 {
		pbCart.Items = make([]*ordersspb.CartItem, 0, len(cart.Items))
		for _, item := range cart.Items {
			pbCart.Items = append(pbCart.Items, domainCartItemToProtoCartItem(item))
		}
	}

	return pbCart
}

// domainCartItemToProtoCartItem converts domain.CartItem to proto CartItem
func domainCartItemToProtoCartItem(item *domain.CartItem) *ordersspb.CartItem {
	if item == nil {
		return nil
	}

	return &ordersspb.CartItem{
		Id:            item.ID,
		CartId:        item.CartID,
		ListingId:     item.ListingID,
		VariantId:     item.VariantID,
		Quantity:      item.Quantity,
		PriceSnapshot: item.PriceSnapshot,
		CreatedAt:     timestamppb.New(item.CreatedAt),
		UpdatedAt:     timestamppb.New(item.UpdatedAt),
	}
}

// domainOrderToProtoOrder converts domain.Order to proto Order
func domainOrderToProtoOrder(order *domain.Order) *ordersspb.Order {
	if order == nil {
		return nil
	}

	pbOrder := &ordersspb.Order{
		Id:           order.ID,
		OrderNumber:  order.OrderNumber,
		UserId:       order.UserID,
		StorefrontId: order.StorefrontID,
		Status:       protoOrderStatusFromDomain(order.Status),
		Financials: &ordersspb.OrderFinancials{
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

	// Convert order items
	if len(order.Items) > 0 {
		pbOrder.Items = make([]*ordersspb.OrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			pbOrder.Items = append(pbOrder.Items, domainOrderItemToProtoOrderItem(item))
		}
	}

	return pbOrder
}

// domainOrderItemToProtoOrderItem converts domain.OrderItem to proto OrderItem
func domainOrderItemToProtoOrderItem(item *domain.OrderItem) *ordersspb.OrderItem {
	if item == nil {
		return nil
	}

	return &ordersspb.OrderItem{
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
}

// calculateCartSummary calculates cart summary (totals, warnings)
func calculateCartSummary(cart *domain.Cart) *ordersspb.CartSummary {
	if cart == nil || len(cart.Items) == 0 {
		return &ordersspb.CartSummary{
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

	return &ordersspb.CartSummary{
		TotalItems:      totalItems,
		Subtotal:        subtotal,
		EstimatedTotal:  subtotal, // TODO: Add tax and shipping calculation
		EstimatedTax:    0,
		EstimatedShipping: 0,
	}
}

// protoOrderStatusFromDomain converts domain.OrderStatus to proto OrderStatus
func protoOrderStatusFromDomain(status domain.OrderStatus) ordersspb.OrderStatus {
	switch status {
	case domain.OrderStatusPending:
		return ordersspb.OrderStatus_ORDER_STATUS_PENDING
	case domain.OrderStatusConfirmed:
		return ordersspb.OrderStatus_ORDER_STATUS_CONFIRMED
	case domain.OrderStatusProcessing:
		return ordersspb.OrderStatus_ORDER_STATUS_PROCESSING
	case domain.OrderStatusShipped:
		return ordersspb.OrderStatus_ORDER_STATUS_SHIPPED
	case domain.OrderStatusDelivered:
		return ordersspb.OrderStatus_ORDER_STATUS_DELIVERED
	case domain.OrderStatusCancelled:
		return ordersspb.OrderStatus_ORDER_STATUS_CANCELLED
	case domain.OrderStatusRefunded:
		return ordersspb.OrderStatus_ORDER_STATUS_REFUNDED
	case domain.OrderStatusFailed:
		return ordersspb.OrderStatus_ORDER_STATUS_FAILED
	default:
		return ordersspb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

// domainOrderStatusFromProto converts proto OrderStatus to domain.OrderStatus
func domainOrderStatusFromProto(status *ordersspb.OrderStatus) *domain.OrderStatus {
	if status == nil {
		return nil
	}

	var domainStatus domain.OrderStatus
	switch *status {
	case ordersspb.OrderStatus_ORDER_STATUS_PENDING:
		domainStatus = domain.OrderStatusPending
	case ordersspb.OrderStatus_ORDER_STATUS_CONFIRMED:
		domainStatus = domain.OrderStatusConfirmed
	case ordersspb.OrderStatus_ORDER_STATUS_PROCESSING:
		domainStatus = domain.OrderStatusProcessing
	case ordersspb.OrderStatus_ORDER_STATUS_SHIPPED:
		domainStatus = domain.OrderStatusShipped
	case ordersspb.OrderStatus_ORDER_STATUS_DELIVERED:
		domainStatus = domain.OrderStatusDelivered
	case ordersspb.OrderStatus_ORDER_STATUS_CANCELLED:
		domainStatus = domain.OrderStatusCancelled
	case ordersspb.OrderStatus_ORDER_STATUS_REFUNDED:
		domainStatus = domain.OrderStatusRefunded
	case ordersspb.OrderStatus_ORDER_STATUS_FAILED:
		domainStatus = domain.OrderStatusFailed
	default:
		return nil
	}

	return &domainStatus
}

// protoPaymentStatusFromDomain converts domain.PaymentStatus to proto PaymentStatus
func protoPaymentStatusFromDomain(status domain.PaymentStatus) ordersspb.PaymentStatus {
	switch status {
	case domain.PaymentStatusPending:
		return ordersspb.PaymentStatus_PAYMENT_STATUS_PENDING
	case domain.PaymentStatusProcessing:
		return ordersspb.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case domain.PaymentStatusCompleted:
		return ordersspb.PaymentStatus_PAYMENT_STATUS_COMPLETED
	case domain.PaymentStatusFailed:
		return ordersspb.PaymentStatus_PAYMENT_STATUS_FAILED
	case domain.PaymentStatusRefunded:
		return ordersspb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return ordersspb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
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

	// Default to internal error
	return status.Error(codes.Internal, "internal server error")
}
