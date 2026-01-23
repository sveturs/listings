//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/tests"
)

// TestCheckStockAvailability_SingleProduct_Sufficient verifies stock check for single product with sufficient stock
// Business requirement: Order validation must confirm adequate inventory before processing
func TestCheckStockAvailability_SingleProduct_Sufficient(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Product 5000 has 100 units in stock
	productID := int64(5000)
	requestedQty := int32(50)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  requestedQty,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.AllAvailable, "All items should be available")
	require.Len(t, resp.Items, 1, "Should return exactly 1 item")

	// Verify item details
	item := resp.Items[0]
	assert.Equal(t, productID, item.ProductId, "Product ID should match")
	assert.Nil(t, item.VariantId, "Variant ID should be nil for product-level check")
	assert.Equal(t, requestedQty, item.RequestedQuantity, "Requested quantity should match")
	assert.Equal(t, int32(100), item.AvailableQuantity, "Available quantity should be 100")
	assert.True(t, item.IsAvailable, "Item should be marked as available")

	// Verify database state unchanged (read-only operation)
	actualStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(100), actualStock, "Stock quantity should remain unchanged")
}

// TestCheckStockAvailability_SingleProduct_Insufficient verifies stock check when quantity exceeds availability
// Business requirement: System must reject orders that exceed available inventory
func TestCheckStockAvailability_SingleProduct_Insufficient(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Product 5001 has only 5 units in stock
	productID := int64(5001)
	requestedQty := int32(50) // Exceeds available stock

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  requestedQty,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should not return error for insufficient stock")
	require.NotNil(t, resp, "Response should not be nil")
	assert.False(t, resp.AllAvailable, "AllAvailable should be false when any item unavailable")
	require.Len(t, resp.Items, 1, "Should return exactly 1 item")

	// Verify item details
	item := resp.Items[0]
	assert.Equal(t, productID, item.ProductId, "Product ID should match")
	assert.Equal(t, requestedQty, item.RequestedQuantity, "Requested quantity should match")
	assert.Equal(t, int32(5), item.AvailableQuantity, "Available quantity should be 5")
	assert.False(t, item.IsAvailable, "Item should be marked as unavailable")

	// Verify database state unchanged
	actualStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(5), actualStock, "Stock quantity should remain unchanged")
}

// TestCheckStockAvailability_SingleProduct_ExactMatch verifies stock check when requested equals available
// Edge case: Ensures exact match is treated as available (>= comparison)
func TestCheckStockAvailability_SingleProduct_ExactMatch(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Product 5001 has exactly 5 units in stock
	productID := int64(5001)
	requestedQty := int32(5) // Exact match

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  requestedQty,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.AllAvailable, "Exact match should be available")
	require.Len(t, resp.Items, 1, "Should return exactly 1 item")

	// Verify item details
	item := resp.Items[0]
	assert.Equal(t, productID, item.ProductId, "Product ID should match")
	assert.Equal(t, requestedQty, item.RequestedQuantity, "Requested quantity should match")
	assert.Equal(t, int32(5), item.AvailableQuantity, "Available quantity should be 5")
	assert.True(t, item.IsAvailable, "Item should be available when requested == available")

	// Verify database state unchanged
	actualStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(5), actualStock, "Stock quantity should remain unchanged")
}

// TestCheckStockAvailability_MultipleProducts_AllAvailable verifies batch check with all items available
// Business requirement: Order validation for multi-item carts must verify all items
func TestCheckStockAvailability_MultipleProducts_AllAvailable(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: 5000, // 100 units available
				Quantity:  50,
			},
			{
				ProductId: 5003, // 50 units available
				Quantity:  25,
			},
			{
				ProductId: 5004, // 75 units available
				Quantity:  30,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.AllAvailable, "All items should be available")
	require.Len(t, resp.Items, 3, "Should return all 3 items")

	// Verify all items are marked as available
	for i, item := range resp.Items {
		assert.True(t, item.IsAvailable, "Item %d should be available", i)
		assert.GreaterOrEqual(t, item.AvailableQuantity, item.RequestedQuantity,
			"Item %d: available >= requested", i)
	}

	// Verify specific items
	assert.Equal(t, int64(5000), resp.Items[0].ProductId)
	assert.Equal(t, int32(100), resp.Items[0].AvailableQuantity)

	assert.Equal(t, int64(5003), resp.Items[1].ProductId)
	assert.Equal(t, int32(50), resp.Items[1].AvailableQuantity)

	assert.Equal(t, int64(5004), resp.Items[2].ProductId)
	assert.Equal(t, int32(75), resp.Items[2].AvailableQuantity)

	// Verify database state unchanged for all products
	for _, stockItem := range req.Items {
		actualStock := tests.GetProductQuantity(t, testDB.DB, stockItem.ProductId)
		assert.Greater(t, actualStock, int32(0), "Product %d should have stock", stockItem.ProductId)
	}
}

// TestCheckStockAvailability_MultipleProducts_PartialAvailable verifies batch check with mixed availability
// Business requirement: System must report which items are unavailable for partial order scenarios
func TestCheckStockAvailability_MultipleProducts_PartialAvailable(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()
	_ = testDB // Avoid unused variable error

	ctx := tests.TestContext(t)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: 5000, // 100 units - AVAILABLE
				Quantity:  50,
			},
			{
				ProductId: 5001, // 5 units - UNAVAILABLE (requesting 50)
				Quantity:  50,
			},
			{
				ProductId: 5002, // 0 units - UNAVAILABLE (out of stock)
				Quantity:  10,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed even with partial availability")
	require.NotNil(t, resp, "Response should not be nil")
	assert.False(t, resp.AllAvailable, "AllAvailable should be false when any item unavailable")
	require.Len(t, resp.Items, 3, "Should return all 3 items")

	// Verify item 0: AVAILABLE
	item0 := resp.Items[0]
	assert.Equal(t, int64(5000), item0.ProductId)
	assert.Equal(t, int32(50), item0.RequestedQuantity)
	assert.Equal(t, int32(100), item0.AvailableQuantity)
	assert.True(t, item0.IsAvailable, "First item should be available")

	// Verify item 1: UNAVAILABLE (insufficient stock)
	item1 := resp.Items[1]
	assert.Equal(t, int64(5001), item1.ProductId)
	assert.Equal(t, int32(50), item1.RequestedQuantity)
	assert.Equal(t, int32(5), item1.AvailableQuantity)
	assert.False(t, item1.IsAvailable, "Second item should be unavailable")

	// Verify item 2: UNAVAILABLE (out of stock)
	item2 := resp.Items[2]
	assert.Equal(t, int64(5002), item2.ProductId)
	assert.Equal(t, int32(10), item2.RequestedQuantity)
	assert.Equal(t, int32(0), item2.AvailableQuantity)
	assert.False(t, item2.IsAvailable, "Third item should be unavailable (out of stock)")

	// Count available vs unavailable
	availableCount := 0
	for _, item := range resp.Items {
		if item.IsAvailable {
			availableCount++
		}
	}
	assert.Equal(t, 1, availableCount, "Should have exactly 1 available item")
}

// TestCheckStockAvailability_ProductNotFound verifies handling of non-existent product
// Business requirement: System must gracefully handle invalid product IDs without error
func TestCheckStockAvailability_ProductNotFound(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()
	_ = testDB // Use testDB to avoid unused variable error

	ctx := tests.TestContext(t)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: 99999, // Non-existent product
				Quantity:  10,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions - should NOT return error, but mark as unavailable
	require.NoError(t, err, "CheckStockAvailability should not error for non-existent product")
	require.NotNil(t, resp, "Response should not be nil")
	assert.False(t, resp.AllAvailable, "AllAvailable should be false")
	require.Len(t, resp.Items, 1, "Should return 1 item")

	// Verify item marked as unavailable with zero stock
	item := resp.Items[0]
	assert.Equal(t, int64(99999), item.ProductId)
	assert.Equal(t, int32(10), item.RequestedQuantity)
	assert.Equal(t, int32(0), item.AvailableQuantity, "Non-existent product should show 0 available")
	assert.False(t, item.IsAvailable, "Non-existent product should be unavailable")
}

// TestCheckStockAvailability_VariantNotFound verifies handling of non-existent variant
// Business requirement: Invalid variant IDs should be treated as unavailable
func TestCheckStockAvailability_VariantNotFound(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)
	nonExistentVariantID := int64(99999)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				VariantId: &nonExistentVariantID,
				Quantity:  10,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions - should NOT return error, but mark as unavailable
	require.NoError(t, err, "CheckStockAvailability should not error for non-existent variant")
	require.NotNil(t, resp, "Response should not be nil")
	assert.False(t, resp.AllAvailable, "AllAvailable should be false")
	require.Len(t, resp.Items, 1, "Should return 1 item")

	// Verify item marked as unavailable
	item := resp.Items[0]
	assert.Equal(t, productID, item.ProductId)
	assert.Equal(t, &nonExistentVariantID, item.VariantId)
	assert.Equal(t, int32(0), item.AvailableQuantity, "Non-existent variant should show 0 available")
	assert.False(t, item.IsAvailable, "Non-existent variant should be unavailable")
}

// TestCheckStockAvailability_ZeroStock verifies handling of product with zero stock
// Business requirement: Out-of-stock products must be clearly marked as unavailable
func TestCheckStockAvailability_ZeroStock(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Product 5002 has 0 units in stock
	productID := int64(5002)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				Quantity:  1, // Even requesting 1 unit should fail
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.False(t, resp.AllAvailable, "AllAvailable should be false for zero stock")
	require.Len(t, resp.Items, 1, "Should return 1 item")

	// Verify item details
	item := resp.Items[0]
	assert.Equal(t, productID, item.ProductId)
	assert.Equal(t, int32(1), item.RequestedQuantity)
	assert.Equal(t, int32(0), item.AvailableQuantity, "Zero stock should be reported")
	assert.False(t, item.IsAvailable, "Zero stock item should be unavailable")

	// Verify database confirms zero stock
	actualStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(0), actualStock, "Database should confirm zero stock")
}

// TestCheckStockAvailability_InvalidQuantity verifies validation of invalid quantity values
// Business requirement: API must reject invalid input (negative/zero quantities)
func TestCheckStockAvailability_InvalidQuantity(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name      string
		quantity  int32
		wantCode  codes.Code
		wantError string
	}{
		{
			name:      "Zero quantity",
			quantity:  0,
			wantCode:  codes.InvalidArgument,
			wantError: "quantity must be positive",
		},
		{
			name:      "Negative quantity",
			quantity:  -5,
			wantCode:  codes.InvalidArgument,
			wantError: "quantity must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CheckStockAvailabilityRequest{
				Items: []*pb.StockItem{
					{
						ProductId: 5000,
						Quantity:  tc.quantity,
					},
				},
			}

			resp, err := client.CheckStockAvailability(ctx, req)

			// Assertions
			require.Error(t, err, "Should return error for invalid quantity")
			assert.Nil(t, resp, "Response should be nil on validation error")

			// Verify gRPC status code
			st, ok := status.FromError(err)
			require.True(t, ok, "Error should be gRPC status")
			assert.Equal(t, tc.wantCode, st.Code(), "Should return InvalidArgument status code")
			assert.Contains(t, st.Message(), tc.wantError, "Error message should mention quantity validation")
		})
	}
}

// TestCheckStockAvailability_EmptyRequest verifies handling of empty items list
// Edge case: Empty request should be rejected
func TestCheckStockAvailability_EmptyRequest(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{}, // Empty list
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.Error(t, err, "Empty request should return error")
	assert.Nil(t, resp, "Response should be nil")

	// Verify gRPC status code
	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be gRPC status")
	assert.Equal(t, codes.InvalidArgument, st.Code(), "Should return InvalidArgument status code")
	assert.Contains(t, st.Message(), "no items provided", "Error should mention empty items")
}

// TestCheckStockAvailability_InvalidProductID verifies validation of invalid product IDs
// Business requirement: API must reject malformed input
func TestCheckStockAvailability_InvalidProductID(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name      string
		productID int64
		wantCode  codes.Code
	}{
		{
			name:      "Zero product ID",
			productID: 0,
			wantCode:  codes.InvalidArgument,
		},
		{
			name:      "Negative product ID",
			productID: -1,
			wantCode:  codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CheckStockAvailabilityRequest{
				Items: []*pb.StockItem{
					{
						ProductId: tc.productID,
						Quantity:  10,
					},
				},
			}

			resp, err := client.CheckStockAvailability(ctx, req)

			// Assertions
			require.Error(t, err, "Should return error for invalid product ID")
			assert.Nil(t, resp, "Response should be nil")

			// Verify gRPC status code
			st, ok := status.FromError(err)
			require.True(t, ok, "Error should be gRPC status")
			assert.Equal(t, tc.wantCode, st.Code(), "Should return InvalidArgument status code")
		})
	}
}

// TestCheckStockAvailability_VariantLevel verifies stock check for product variants
// Business requirement: System must support variant-level inventory checks
func TestCheckStockAvailability_VariantLevel(t *testing.T) {
	t.Skip("Skipping: variant fixtures not loaded (product 5000 / variant 6000 do not exist)")

	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)
	variantID := int64(6000) // Size S variant with 50 units

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: productID,
				VariantId: &variantID,
				Quantity:  30,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed for variant")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.AllAvailable, "Variant should be available")
	require.Len(t, resp.Items, 1, "Should return 1 item")

	// Verify variant-specific details
	item := resp.Items[0]
	assert.Equal(t, productID, item.ProductId)
	assert.Equal(t, &variantID, item.VariantId, "Variant ID should match")
	assert.Equal(t, int32(30), item.RequestedQuantity)
	assert.Equal(t, int32(50), item.AvailableQuantity, "Variant has 50 units")
	assert.True(t, item.IsAvailable)

	// Verify database state unchanged
	actualStock := tests.GetVariantQuantity(t, testDB.DB, variantID)
	assert.Equal(t, int32(50), actualStock, "Variant stock should remain unchanged")
}

// TestCheckStockAvailability_MixedProductsAndVariants verifies batch check with both products and variants
// Business requirement: Order validation must support mixed cart items
func TestCheckStockAvailability_MixedProductsAndVariants(t *testing.T) {
	t.Skip("Skipping: variant fixtures not loaded (variant 6000 does not exist)")

	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	variantID := int64(6000)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: 5003, // Product without variant
				Quantity:  25,
			},
			{
				ProductId: 5000, // Product with variant specified
				VariantId: &variantID,
				Quantity:  30,
			},
		},
	}

	resp, err := client.CheckStockAvailability(ctx, req)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.AllAvailable, "All items should be available")
	require.Len(t, resp.Items, 2, "Should return 2 items")

	// Verify product-level check
	assert.Equal(t, int64(5003), resp.Items[0].ProductId)
	assert.Nil(t, resp.Items[0].VariantId, "First item should have no variant")
	assert.True(t, resp.Items[0].IsAvailable)

	// Verify variant-level check
	assert.Equal(t, int64(5000), resp.Items[1].ProductId)
	assert.Equal(t, &variantID, resp.Items[1].VariantId, "Second item should have variant")
	assert.True(t, resp.Items[1].IsAvailable)
}

// TestCheckStockAvailability_PerformanceUnder100ms verifies response time for single item check
// Performance requirement: Stock availability checks must complete quickly for order validation
func TestCheckStockAvailability_PerformanceUnder100ms(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{
				ProductId: 5000,
				Quantity:  50,
			},
		},
	}

	// Warmup call
	_, _ = client.CheckStockAvailability(ctx, req)

	// Measure actual call
	start := time.Now()
	resp, err := client.CheckStockAvailability(ctx, req)
	elapsed := time.Since(start)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Less(t, elapsed.Milliseconds(), int64(100),
		"Single item check should complete in < 100ms (actual: %dms)", elapsed.Milliseconds())
}

// TestCheckStockAvailability_BatchPerformance verifies response time for batch check
// Performance requirement: Batch checks (10 items) should complete efficiently
func TestCheckStockAvailability_BatchPerformance(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create request with 10 items
	items := make([]*pb.StockItem, 10)
	for i := 0; i < 10; i++ {
		items[i] = &pb.StockItem{
			ProductId: 5000, // Reusing same product for consistency
			Quantity:  int32(10 + i),
		}
	}

	req := &pb.CheckStockAvailabilityRequest{
		Items: items,
	}

	// Warmup call
	_, _ = client.CheckStockAvailability(ctx, req)

	// Measure actual call
	start := time.Now()
	resp, err := client.CheckStockAvailability(ctx, req)
	elapsed := time.Since(start)

	// Assertions
	require.NoError(t, err, "CheckStockAvailability should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Less(t, elapsed.Milliseconds(), int64(200),
		"Batch check (10 items) should complete in < 200ms (actual: %dms)", elapsed.Milliseconds())
	assert.Len(t, resp.Items, 10, "Should return all 10 items")
}

// TestCheckStockAvailability_ConcurrentRequests verifies thread safety
// Performance requirement: Concurrent stock checks must not interfere with each other
func TestCheckStockAvailability_ConcurrentRequests(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	_ = testDB // Avoid unused variable error
	defer cleanup()

	ctx := tests.TestContext(t)

	concurrentCalls := 20
	done := make(chan error, concurrentCalls)

	// Launch concurrent requests
	for i := 0; i < concurrentCalls; i++ {
		go func(callNum int) {
			req := &pb.CheckStockAvailabilityRequest{
				Items: []*pb.StockItem{
					{
						ProductId: 5000,
						Quantity:  int32(10 + callNum),
					},
				},
			}

			resp, err := client.CheckStockAvailability(ctx, req)
			if err != nil {
				done <- err
				return
			}

			if resp == nil || len(resp.Items) != 1 {
				done <- assert.AnError
				return
			}

			done <- nil
		}(i)
	}

	// Wait for all goroutines
	successCount := 0
	for i := 0; i < concurrentCalls; i++ {
		err := <-done
		if err == nil {
			successCount++
		}
	}

	assert.Equal(t, concurrentCalls, successCount,
		"All %d concurrent calls should succeed", concurrentCalls)
}

// TestCheckStockAvailability_ReadOnlyOperation verifies that check doesn't modify data
// Business requirement: Stock availability checks must be read-only operations
func TestCheckStockAvailability_ReadOnlyOperation(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)

	// Get initial stock
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)

	// Perform multiple checks
	for i := 0; i < 5; i++ {
		req := &pb.CheckStockAvailabilityRequest{
			Items: []*pb.StockItem{
				{
					ProductId: productID,
					Quantity:  int32(50),
				},
			},
		}

		resp, err := client.CheckStockAvailability(ctx, req)
		require.NoError(t, err, "Check %d should succeed", i+1)
		require.NotNil(t, resp, "Check %d response should not be nil", i+1)
	}

	// Verify stock unchanged after multiple checks
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialStock, finalStock,
		"Stock quantity should remain unchanged after multiple availability checks")
}
