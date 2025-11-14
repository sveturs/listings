package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/service"
)

// ============================================================================
// ORDER SERVICE gRPC METHODS
// These methods are added to the Server struct and implement OrderService RPC
// ============================================================================

// ============================================================================
// CART OPERATIONS (6 methods)
// ============================================================================

// AddToCart adds an item to shopping cart
func (s *Server) AddToCart(ctx context.Context, req *listingspb.AddToCartRequest) (*listingspb.AddToCartResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Int64("listing_id", req.ListingId).
		Int32("quantity", req.Quantity).
		Msg("AddToCart called")

	// Cart operations not yet implemented in OrderService
	// Return Unimplemented for now
	return nil, status.Error(codes.Unimplemented, "cart operations not yet implemented")
}

// UpdateCartItem updates quantity or variant for a cart item
func (s *Server) UpdateCartItem(ctx context.Context, req *listingspb.UpdateCartItemRequest) (*listingspb.UpdateCartItemResponse, error) {
	s.logger.Debug().
		Int64("cart_item_id", req.CartItemId).
		Msg("UpdateCartItem called")

	// Cart operations not yet implemented
	return nil, status.Error(codes.Unimplemented, "cart operations not yet implemented")
}

// RemoveFromCart removes an item from cart
func (s *Server) RemoveFromCart(ctx context.Context, req *listingspb.RemoveFromCartRequest) (*listingspb.RemoveFromCartResponse, error) {
	s.logger.Debug().
		Int64("cart_item_id", req.CartItemId).
		Msg("RemoveFromCart called")

	// Cart operations not yet implemented
	return nil, status.Error(codes.Unimplemented, "cart operations not yet implemented")
}

// GetCart retrieves user's cart for a storefront
func (s *Server) GetCart(ctx context.Context, req *listingspb.GetCartRequest) (*listingspb.GetCartResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Msg("GetCart called")

	// Cart operations not yet implemented
	return nil, status.Error(codes.Unimplemented, "cart operations not yet implemented")
}

// ClearCart removes all items from cart
func (s *Server) ClearCart(ctx context.Context, req *listingspb.ClearCartRequest) (*emptypb.Empty, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Msg("ClearCart called")

	// Cart operations not yet implemented
	return nil, status.Error(codes.Unimplemented, "cart operations not yet implemented")
}

// GetUserCarts retrieves all carts for authenticated user
func (s *Server) GetUserCarts(ctx context.Context, req *listingspb.GetUserCartsRequest) (*listingspb.GetUserCartsResponse, error) {
	s.logger.Debug().
		Int64("user_id", req.UserId).
		Msg("GetUserCarts called")

	// Cart operations not yet implemented
	return nil, status.Error(codes.Unimplemented, "cart operations not yet implemented")
}

// ============================================================================
// ORDER OPERATIONS (6 methods)
// ============================================================================

// CreateOrder creates a new order from cart
func (s *Server) CreateOrder(ctx context.Context, req *listingspb.CreateOrderRequest) (*listingspb.CreateOrderResponse, error) {
	s.logger.Info().
		Int64("storefront_id", req.StorefrontId).
		Int64("cart_id", req.CartId).
		Str("payment_method", req.PaymentMethod).
		Msg("CreateOrder called")

	// 1. Validate request
	if req.GetStorefrontId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}
	if req.GetCartId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "cart_id is required")
	}
	if req.GetPaymentMethod() == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_method is required")
	}
	if req.GetShippingMethod() == "" {
		return nil, status.Error(codes.InvalidArgument, "shipping_method is required")
	}
	if req.ShippingAddress == nil {
		return nil, status.Error(codes.InvalidArgument, "shipping_address is required")
	}

	// 2. Convert proto to domain request
	domainReq := &service.CreateOrderRequest{
		CartID:          req.GetCartId(),
		UserID:          req.UserId,
		ShippingAddress: req.ShippingAddress.AsMap(),
		ShippingCost:    0, // TODO: Calculate shipping cost
		PaymentMethod:   req.GetPaymentMethod(),
	}

	// Set billing address (defaults to shipping if not provided)
	if req.BillingAddress != nil {
		domainReq.BillingAddress = req.BillingAddress.AsMap()
	}

	// Set optional customer notes
	if req.CustomerNotes != nil {
		domainReq.CustomerNotes = req.CustomerNotes
	}

	// 3. Call service layer
	order, err := s.orderService.CreateOrder(ctx, domainReq)
	if err != nil {
		s.logger.Error().Err(err).Int64("cart_id", req.GetCartId()).Msg("failed to create order")
		return nil, mapOrderServiceError(err)
	}

	// 4. Convert domain to proto
	pbOrder := order.ToProto()

	// 5. Build warnings for price changes
	var warnings []string
	// TODO: Extract price change warnings from service response

	s.logger.Info().
		Int64("order_id", order.ID).
		Str("order_number", order.OrderNumber).
		Msg("order created successfully")

	return &listingspb.CreateOrderResponse{
		Order:    pbOrder,
		Message:  "Order created successfully",
		Warnings: warnings,
	}, nil
}

// GetOrder retrieves a single order by ID
func (s *Server) GetOrder(ctx context.Context, req *listingspb.GetOrderRequest) (*listingspb.GetOrderResponse, error) {
	s.logger.Debug().
		Int64("order_id", req.OrderId).
		Msg("GetOrder called")

	// 1. Validate request
	if req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	// 2. Call service layer
	order, err := s.orderService.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		s.logger.Error().Err(err).Int64("order_id", req.GetOrderId()).Msg("failed to get order")
		return nil, mapOrderServiceError(err)
	}

	// 3. Verify ownership if user_id provided
	if req.UserId != nil && order.UserID != nil && *order.UserID != *req.UserId {
		s.logger.Warn().
			Int64("order_id", req.GetOrderId()).
			Int64("requested_by", *req.UserId).
			Int64("order_owner", *order.UserID).
			Msg("unauthorized order access attempt")
		return nil, status.Error(codes.PermissionDenied, "cannot access other users' orders")
	}

	// 4. Convert domain to proto
	pbOrder := order.ToProto()

	return &listingspb.GetOrderResponse{
		Order: pbOrder,
	}, nil
}

// ListOrders retrieves orders with filters and pagination
func (s *Server) ListOrders(ctx context.Context, req *listingspb.ListOrdersRequest) (*listingspb.ListOrdersResponse, error) {
	s.logger.Debug().
		Int32("page", req.Page).
		Int32("page_size", req.PageSize).
		Msg("ListOrders called")

	// 1. Set defaults
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

	// 2. Build domain filter
	filter := &service.ListOrdersRequest{
		Limit:  int(pageSize),
		Offset: int((page - 1) * pageSize),
	}

	// 3. Apply filters
	if req.UserId != nil {
		userID := *req.UserId
		filter.UserID = &userID
	}
	if req.StorefrontId != nil {
		storefrontID := *req.StorefrontId
		filter.StorefrontID = &storefrontID
	}
	if req.Status != nil {
		statusDomain := domain.OrderStatusFromProto(*req.Status)
		filter.Status = &statusDomain
	}

	// 4. Call service layer
	orders, total, err := s.orderService.ListOrders(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list orders")
		return nil, mapOrderServiceError(err)
	}

	// 5. Convert domain to proto
	pbOrders := make([]*listingspb.Order, len(orders))
	for i, order := range orders {
		pbOrders[i] = order.ToProto()
	}

	s.logger.Info().
		Int("count", len(orders)).
		Int64("total", total).
		Msg("orders listed successfully")

	return &listingspb.ListOrdersResponse{
		Orders:     pbOrders,
		TotalCount: int32(total),
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// CancelOrder cancels an order and releases inventory
func (s *Server) CancelOrder(ctx context.Context, req *listingspb.CancelOrderRequest) (*listingspb.CancelOrderResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Bool("refund", req.Refund).
		Msg("CancelOrder called")

	// 1. Validate request
	if req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	// 2. Determine user_id for cancellation
	var userID int64
	if req.UserId != nil {
		userID = *req.UserId
	} else {
		// If no user_id in request, get from order
		order, err := s.orderService.GetOrder(ctx, req.GetOrderId())
		if err != nil {
			return nil, mapOrderServiceError(err)
		}
		if order.UserID == nil {
			return nil, status.Error(codes.Internal, "order has no owner")
		}
		userID = *order.UserID
	}

	// 3. Call service layer
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	order, err := s.orderService.CancelOrder(ctx, req.GetOrderId(), userID, reason)
	if err != nil {
		s.logger.Error().Err(err).Int64("order_id", req.GetOrderId()).Msg("failed to cancel order")
		return nil, mapOrderServiceError(err)
	}

	// 4. Convert domain to proto
	pbOrder := order.ToProto()

	// 5. Determine if refund was initiated
	refundInitiated := req.Refund && order.PaymentStatus == domain.PaymentStatusCompleted

	s.logger.Info().
		Int64("order_id", order.ID).
		Bool("refund_initiated", refundInitiated).
		Msg("order cancelled successfully")

	return &listingspb.CancelOrderResponse{
		Order:           pbOrder,
		Message:         "Order cancelled successfully",
		RefundInitiated: refundInitiated,
	}, nil
}

// UpdateOrderStatus updates order status (admin only)
func (s *Server) UpdateOrderStatus(ctx context.Context, req *listingspb.UpdateOrderStatusRequest) (*listingspb.UpdateOrderStatusResponse, error) {
	s.logger.Info().
		Int64("order_id", req.OrderId).
		Str("new_status", req.NewStatus.String()).
		Msg("UpdateOrderStatus called")

	// 1. Validate request
	if req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.NewStatus == listingspb.OrderStatus_ORDER_STATUS_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "new_status is required")
	}

	// 2. Validate tracking_number if status is SHIPPED
	if req.NewStatus == listingspb.OrderStatus_ORDER_STATUS_SHIPPED {
		if req.TrackingNumber == nil || *req.TrackingNumber == "" {
			return nil, status.Error(codes.InvalidArgument, "tracking_number is required when status is SHIPPED")
		}
	}

	// 3. Convert proto status to domain status
	newStatus := domain.OrderStatusFromProto(req.NewStatus)

	// 4. Call service layer
	order, err := s.orderService.UpdateOrderStatus(ctx, req.GetOrderId(), newStatus)
	if err != nil {
		s.logger.Error().Err(err).Int64("order_id", req.GetOrderId()).Msg("failed to update order status")
		return nil, mapOrderServiceError(err)
	}

	// 5. Convert domain to proto
	pbOrder := order.ToProto()

	// 6. Build warnings for invalid state transitions
	var warnings []string
	// Warnings are already handled by service layer errors

	s.logger.Info().
		Int64("order_id", order.ID).
		Str("new_status", string(order.Status)).
		Msg("order status updated successfully")

	return &listingspb.UpdateOrderStatusResponse{
		Order:    pbOrder,
		Message:  "Order status updated successfully",
		Warnings: warnings,
	}, nil
}

// GetOrderStats retrieves order statistics (admin)
func (s *Server) GetOrderStats(ctx context.Context, req *listingspb.GetOrderStatsRequest) (*listingspb.GetOrderStatsResponse, error) {
	s.logger.Debug().
		Msg("GetOrderStats called")

	// 1. Build filters
	var storefrontIDPtr *int64
	if req.StorefrontId != nil {
		storefrontID := *req.StorefrontId
		storefrontIDPtr = &storefrontID
	}

	// 2. Call service layer
	stats, err := s.orderService.GetOrderStats(ctx, nil, storefrontIDPtr)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get order stats")
		return nil, mapOrderServiceError(err)
	}

	// 3. Convert domain to proto
	pbStats := &listingspb.OrderStatsSummary{
		TotalOrders:     int32(stats.TotalOrders),
		PendingOrders:   int32(stats.PendingOrders),
		ConfirmedOrders: int32(stats.ConfirmedOrders),
		// CompletedOrders field doesn't exist in proto, using total
		CancelledOrders: int32(stats.CancelledOrders),
		TotalRevenue:    stats.TotalRevenue,
	}

	s.logger.Info().
		Int64("total_orders", stats.TotalOrders).
		Float64("total_revenue", stats.TotalRevenue).
		Msg("order stats retrieved successfully")

	return &listingspb.GetOrderStatsResponse{
		Stats: pbStats,
		// StatusBreakdown: nil, // TODO: Implement status breakdown
		// DailyStats:      nil, // TODO: Implement daily stats
	}, nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// extractAuthFromContext extracts user_id from gRPC metadata
// Returns (userID *int64, sessionID *string)
func extractAuthFromContext(ctx context.Context) (*int64, *string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil
	}

	// Try to extract user_id from authorization header (JWT)
	if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
		// authHeader format: "Bearer <JWT>"
		// JWT should contain user_id claim
		// TODO: Parse JWT and extract user_id
	}

	// Try to extract session_id from x-session-id header
	if sessionHeaders := md.Get("x-session-id"); len(sessionHeaders) > 0 {
		sessionID := sessionHeaders[0]
		return nil, &sessionID
	}

	return nil, nil
}

// mapOrderServiceError converts service layer errors to gRPC status codes
func mapOrderServiceError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Order-specific errors
	if contains(errMsg, "order not found") {
		return status.Error(codes.NotFound, "order not found")
	}

	if contains(errMsg, "cart is empty") || contains(errMsg, "cart empty") {
		return status.Error(codes.FailedPrecondition, "cart is empty")
	}

	if contains(errMsg, "unauthorized") {
		return status.Error(codes.PermissionDenied, "unauthorized")
	}

	// Insufficient stock errors
	if contains(errMsg, "insufficient stock") {
		return status.Error(codes.FailedPrecondition, errMsg)
	}

	// Price change errors
	if contains(errMsg, "price changed") {
		return status.Error(codes.FailedPrecondition, errMsg)
	}

	// Order cannot be cancelled errors
	if contains(errMsg, "cannot cancel") {
		return status.Error(codes.FailedPrecondition, errMsg)
	}

	// Invalid status transition errors
	if contains(errMsg, "cannot update status") || contains(errMsg, "invalid status transition") {
		return status.Error(codes.FailedPrecondition, errMsg)
	}

	// Validation errors
	if contains(errMsg, "validation") || contains(errMsg, "invalid") {
		return status.Error(codes.InvalidArgument, errMsg)
	}

	// Not found errors
	if contains(errMsg, "not found") {
		return status.Error(codes.NotFound, errMsg)
	}

	// Already exists errors
	if contains(errMsg, "already exists") || contains(errMsg, "duplicate") {
		return status.Error(codes.AlreadyExists, errMsg)
	}

	// Default to Internal error
	return status.Error(codes.Internal, fmt.Sprintf("internal error: %v", err))
}
