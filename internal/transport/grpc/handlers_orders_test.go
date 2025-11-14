package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	ordersspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	grpcTransport "github.com/sveturs/listings/internal/transport/grpc"
	"github.com/sveturs/listings/internal/testing"
)

// TestAddToCart_Success tests successful add to cart operation
func TestAddToCart_Success(t *testing.T) {
	// Setup test environment
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	// Create order service handler
	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	// Create a test storefront and product
	storefrontID := int64(1)
	listingID := int64(100)
	userID := int64(42)

	// Add to cart request
	req := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
		ListingId:    listingID,
		Quantity:     2,
	}

	// Call AddToCart
	resp, err := orderHandler.AddToCart(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Cart)
	assert.Equal(t, storefrontID, resp.Cart.StorefrontId)
	assert.Equal(t, "Item added to cart successfully", resp.Message)
	assert.Greater(t, len(resp.Cart.Items), 0)
}

// TestAddToCart_InvalidInput tests validation errors
func TestAddToCart_InvalidInput(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *ordersspb.AddToCartRequest
		wantErr codes.Code
	}{
		{
			name: "no user_id or session_id",
			req: &ordersspb.AddToCartRequest{
				StorefrontId: 1,
				ListingId:    100,
				Quantity:     1,
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "both user_id and session_id",
			req: &ordersspb.AddToCartRequest{
				UserId:       ptrInt64(1),
				SessionId:    ptrString("session123"),
				StorefrontId: 1,
				ListingId:    100,
				Quantity:     1,
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "invalid storefront_id",
			req: &ordersspb.AddToCartRequest{
				UserId:       ptrInt64(1),
				StorefrontId: 0,
				ListingId:    100,
				Quantity:     1,
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "invalid listing_id",
			req: &ordersspb.AddToCartRequest{
				UserId:       ptrInt64(1),
				StorefrontId: 1,
				ListingId:    0,
				Quantity:     1,
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "invalid quantity",
			req: &ordersspb.AddToCartRequest{
				UserId:       ptrInt64(1),
				StorefrontId: 1,
				ListingId:    100,
				Quantity:     0,
			},
			wantErr: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := orderHandler.AddToCart(ctx, tt.req)

			assert.Nil(t, resp)
			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tt.wantErr, st.Code())
		})
	}
}

// TestGetCart_Success tests successful get cart operation
func TestGetCart_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// First, add item to cart
	addReq := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
		ListingId:    100,
		Quantity:     3,
	}
	_, err := orderHandler.AddToCart(ctx, addReq)
	require.NoError(t, err)

	// Now get cart
	getReq := &ordersspb.GetCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
	}

	resp, err := orderHandler.GetCart(ctx, getReq)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Cart)
	assert.NotNil(t, resp.Summary)
	assert.Equal(t, storefrontID, resp.Cart.StorefrontId)
	assert.Greater(t, len(resp.Cart.Items), 0)
	assert.Equal(t, int32(3), resp.Summary.TotalItems)
	assert.Greater(t, resp.Summary.Subtotal, 0.0)
}

// TestClearCart_Success tests successful clear cart operation
func TestClearCart_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Add item to cart
	addReq := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
		ListingId:    100,
		Quantity:     2,
	}
	_, err := orderHandler.AddToCart(ctx, addReq)
	require.NoError(t, err)

	// Clear cart
	clearReq := &ordersspb.ClearCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
	}

	_, err = orderHandler.ClearCart(ctx, clearReq)
	require.NoError(t, err)

	// Verify cart is empty
	getReq := &ordersspb.GetCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
	}

	resp, err := orderHandler.GetCart(ctx, getReq)
	require.NoError(t, err)
	assert.Equal(t, 0, len(resp.Cart.Items))
}

// TestGetUserCarts_Success tests get all user carts
func TestGetUserCarts_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	userID := int64(42)

	// Add items to multiple storefronts
	addReq1 := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: 1,
		ListingId:    100,
		Quantity:     1,
	}
	_, err := orderHandler.AddToCart(ctx, addReq1)
	require.NoError(t, err)

	addReq2 := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: 2,
		ListingId:    200,
		Quantity:     2,
	}
	_, err = orderHandler.AddToCart(ctx, addReq2)
	require.NoError(t, err)

	// Get all user carts
	req := &ordersspb.GetUserCartsRequest{
		UserId: userID,
	}

	resp, err := orderHandler.GetUserCarts(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int32(2), resp.TotalCarts)
	assert.Len(t, resp.Carts, 2)
}

// TestCreateOrder_Success tests successful order creation
func TestCreateOrder_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Add item to cart
	addReq := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
		ListingId:    100,
		Quantity:     2,
	}
	addResp, err := orderHandler.AddToCart(ctx, addReq)
	require.NoError(t, err)

	// Create shipping address
	shippingAddress, err := structpb.NewStruct(map[string]interface{}{
		"street":      "123 Main St",
		"city":        "Belgrade",
		"postal_code": "11000",
		"country":     "RS",
	})
	require.NoError(t, err)

	// Create order
	createReq := &ordersspb.CreateOrderRequest{
		UserId:          &userID,
		StorefrontId:    storefrontID,
		CartId:          addResp.Cart.Id,
		ShippingAddress: shippingAddress,
		ShippingMethod:  "standard",
		PaymentMethod:   "card",
	}

	resp, err := orderHandler.CreateOrder(ctx, createReq)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Order)
	assert.NotEmpty(t, resp.Order.OrderNumber)
	assert.Equal(t, ordersspb.OrderStatus_ORDER_STATUS_PENDING, resp.Order.Status)
	assert.Equal(t, ordersspb.PaymentStatus_PAYMENT_STATUS_PENDING, resp.Order.PaymentStatus)
	assert.Greater(t, len(resp.Order.Items), 0)
}

// TestCreateOrder_EmptyCart tests order creation with empty cart
func TestCreateOrder_EmptyCart(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Create empty cart first
	addReq := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
		ListingId:    100,
		Quantity:     1,
	}
	addResp, err := orderHandler.AddToCart(ctx, addReq)
	require.NoError(t, err)

	// Clear cart
	clearReq := &ordersspb.ClearCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
	}
	_, err = orderHandler.ClearCart(ctx, clearReq)
	require.NoError(t, err)

	// Try to create order from empty cart
	shippingAddress, err := structpb.NewStruct(map[string]interface{}{
		"street": "123 Main St",
		"city":   "Belgrade",
	})
	require.NoError(t, err)

	createReq := &ordersspb.CreateOrderRequest{
		UserId:          &userID,
		StorefrontId:    storefrontID,
		CartId:          addResp.Cart.Id,
		ShippingAddress: shippingAddress,
		ShippingMethod:  "standard",
		PaymentMethod:   "card",
	}

	resp, err := orderHandler.CreateOrder(ctx, createReq)

	// Should fail with FailedPrecondition (cart is empty)
	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	// Could be InvalidArgument or FailedPrecondition depending on implementation
	assert.Contains(t, []codes.Code{codes.InvalidArgument, codes.FailedPrecondition}, st.Code())
}

// TestGetOrder_Success tests successful get order
func TestGetOrder_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Create order first
	order := createTestOrder(t, env, orderHandler, storefrontID, userID)

	// Get order
	req := &ordersspb.GetOrderRequest{
		OrderId: order.Id,
		UserId:  &userID,
	}

	resp, err := orderHandler.GetOrder(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Order)
	assert.Equal(t, order.Id, resp.Order.Id)
	assert.Equal(t, order.OrderNumber, resp.Order.OrderNumber)
}

// TestGetOrder_Unauthorized tests unauthorized access to order
func TestGetOrder_Unauthorized(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)
	otherUserID := int64(99)

	// Create order as user 42
	order := createTestOrder(t, env, orderHandler, storefrontID, userID)

	// Try to get order as user 99
	req := &ordersspb.GetOrderRequest{
		OrderId: order.Id,
		UserId:  &otherUserID,
	}

	resp, err := orderHandler.GetOrder(ctx, req)

	// Should fail with PermissionDenied
	assert.Nil(t, resp)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())
}

// TestListOrders_Success tests successful list orders
func TestListOrders_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Create multiple orders
	order1 := createTestOrder(t, env, orderHandler, storefrontID, userID)
	order2 := createTestOrder(t, env, orderHandler, storefrontID, userID)

	// List orders
	req := &ordersspb.ListOrdersRequest{
		UserId:   &userID,
		Page:     1,
		PageSize: 10,
	}

	resp, err := orderHandler.ListOrders(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.GreaterOrEqual(t, len(resp.Orders), 2)
	assert.GreaterOrEqual(t, resp.TotalCount, int32(2))

	// Verify orders are in response
	orderIDs := make([]int64, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		orderIDs = append(orderIDs, order.Id)
	}
	assert.Contains(t, orderIDs, order1.Id)
	assert.Contains(t, orderIDs, order2.Id)
}

// TestCancelOrder_Success tests successful order cancellation
func TestCancelOrder_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Create order
	order := createTestOrder(t, env, orderHandler, storefrontID, userID)

	// Cancel order
	reason := "Changed my mind"
	req := &ordersspb.CancelOrderRequest{
		OrderId: order.Id,
		UserId:  &userID,
		Reason:  &reason,
		Refund:  true,
	}

	resp, err := orderHandler.CancelOrder(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Order)
	assert.Equal(t, ordersspb.OrderStatus_ORDER_STATUS_CANCELLED, resp.Order.Status)
	assert.True(t, resp.RefundInitiated)
}

// TestUpdateOrderStatus_Success tests successful status update (admin operation)
func TestUpdateOrderStatus_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)
	userID := int64(42)

	// Create order (status: PENDING)
	order := createTestOrder(t, env, orderHandler, storefrontID, userID)

	// Update status to CONFIRMED (simulating payment confirmation)
	req := &ordersspb.UpdateOrderStatusRequest{
		OrderId:   order.Id,
		NewStatus: ordersspb.OrderStatus_ORDER_STATUS_CONFIRMED,
	}

	resp, err := orderHandler.UpdateOrderStatus(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Order)
	assert.Equal(t, ordersspb.OrderStatus_ORDER_STATUS_CONFIRMED, resp.Order.Status)
}

// TestGetOrderStats_Success tests order statistics retrieval
func TestGetOrderStats_Success(t *testing.T) {
	env := testing.NewTestEnvironment(t)
	defer env.Cleanup()

	orderHandler := grpcTransport.NewOrderServiceServer(
		env.CartService,
		env.OrderService,
		env.InventoryService,
		env.Logger,
	)

	ctx := context.Background()

	storefrontID := int64(1)

	// Get stats
	req := &ordersspb.GetOrderStatsRequest{
		StorefrontId: &storefrontID,
	}

	resp, err := orderHandler.GetOrderStats(ctx, req)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Stats)
	// Stats should be non-negative
	assert.GreaterOrEqual(t, resp.Stats.TotalOrders, int32(0))
	assert.GreaterOrEqual(t, resp.Stats.TotalRevenue, 0.0)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// createTestOrder creates a test order and returns it
func createTestOrder(
	t *testing.T,
	env *testing.TestEnvironment,
	handler *grpcTransport.OrderServiceServer,
	storefrontID int64,
	userID int64,
) *ordersspb.Order {
	ctx := context.Background()

	// Add item to cart
	addReq := &ordersspb.AddToCartRequest{
		UserId:       &userID,
		StorefrontId: storefrontID,
		ListingId:    100,
		Quantity:     1,
	}
	addResp, err := handler.AddToCart(ctx, addReq)
	require.NoError(t, err)

	// Create shipping address
	shippingAddress, err := structpb.NewStruct(map[string]interface{}{
		"street":      "123 Main St",
		"city":        "Belgrade",
		"postal_code": "11000",
		"country":     "RS",
	})
	require.NoError(t, err)

	// Create order
	createReq := &ordersspb.CreateOrderRequest{
		UserId:          &userID,
		StorefrontId:    storefrontID,
		CartId:          addResp.Cart.Id,
		ShippingAddress: shippingAddress,
		ShippingMethod:  "standard",
		PaymentMethod:   "card",
	}

	resp, err := handler.CreateOrder(ctx, createReq)
	require.NoError(t, err)

	return resp.Order
}

// Helper functions for pointer types
func ptrInt64(v int64) *int64 {
	return &v
}

func ptrString(v string) *string {
	return &v
}
