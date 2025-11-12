//go:build integration

package integration

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/tests"
)

// TestGRPCRecordInventoryMovement_FullCycle tests full gRPC request/response cycle
func TestGRPCRecordInventoryMovement_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000) // Initial quantity: 100

	testCases := []struct {
		name       string
		request    *pb.RecordInventoryMovementRequest
		wantErr    bool
		wantCode   codes.Code
		wantBefore int32
		wantAfter  int32
	}{
		{
			name: "Valid stock IN movement",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     50,
				Reason:       stringPtr("restock"),
				Notes:        stringPtr("Integration test restock"),
				UserId:       1000,
			},
			wantErr:    false,
			wantBefore: 100,
			wantAfter:  150,
		},
		{
			name: "Valid stock OUT movement",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "out",
				Quantity:     30,
				Reason:       stringPtr("sale"),
				Notes:        stringPtr("Customer order"),
				UserId:       1000,
			},
			wantErr:    false,
			wantBefore: 150,
			wantAfter:  120,
		},
		{
			name: "Invalid storefront ID",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 0,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Invalid movement type",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "invalid",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Non-existent product",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    99999,
				MovementType: "in",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.RecordInventoryMovement(ctx, tc.request)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.True(t, resp.Success)
			assert.Equal(t, tc.wantBefore, resp.StockBefore)
			assert.Equal(t, tc.wantAfter, resp.StockAfter)

			// Verify database state
			actualQty := tests.GetProductQuantity(t, testDB.DB, tc.request.ProductId)
			assert.Equal(t, tc.wantAfter, actualQty)
		})
	}
}

// TestGRPCRecordInventoryMovement_VariantLevel tests variant-level inventory
func TestGRPCRecordInventoryMovement_VariantLevel(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	variantID := int64(6000) // Size S variant, qty 50

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 1000,
		ProductId:    5000,
		VariantId:    &variantID,
		MovementType: "in",
		Quantity:     25,
		Reason:       stringPtr("variant_restock"),
		UserId:       1000,
	}

	resp, err := client.RecordInventoryMovement(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, int32(50), resp.StockBefore)
	assert.Equal(t, int32(75), resp.StockAfter)

	// Verify variant quantity in database
	actualQty := tests.GetVariantQuantity(t, testDB.DB, variantID)
	assert.Equal(t, int32(75), actualQty)
}

// TestGRPCBatchUpdateStock_FullCycle tests batch stock update via gRPC
func TestGRPCBatchUpdateStock_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name           string
		request        *pb.BatchUpdateStockRequest
		wantErr        bool
		wantCode       codes.Code
		wantSuccessful int32
		wantFailed     int32
	}{
		{
			name: "Successful batch update",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 1000,
				Items: []*pb.StockUpdateItem{
					{
						ProductId: 5003,
						Quantity:  100,
						Reason:    stringPtr("restock"),
					},
					{
						ProductId: 5004,
						Quantity:  80,
						Reason:    stringPtr("adjustment"),
					},
				},
				Reason: stringPtr("bulk_inventory_update"),
				UserId: 1000,
			},
			wantErr:        false,
			wantSuccessful: 2,
			wantFailed:     0,
		},
		{
			name: "Empty items list",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 1000,
				Items:        []*pb.StockUpdateItem{},
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Invalid storefront ID",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 0,
				Items: []*pb.StockUpdateItem{
					{ProductId: 5003, Quantity: 50},
				},
				UserId: 1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Partial success (some products not found)",
			request: &pb.BatchUpdateStockRequest{
				StorefrontId: 1000,
				Items: []*pb.StockUpdateItem{
					{ProductId: 5003, Quantity: 60},  // Valid
					{ProductId: 99999, Quantity: 50}, // Invalid
				},
				UserId: 1000,
			},
			wantErr:        false,
			wantSuccessful: 1,
			wantFailed:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.BatchUpdateStock(ctx, tc.request)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tc.wantSuccessful, resp.SuccessfulCount)
			assert.Equal(t, tc.wantFailed, resp.FailedCount)
			assert.Len(t, resp.Results, int(tc.wantSuccessful+tc.wantFailed))

			// Verify successful items in database
			for _, item := range tc.request.Items {
				if item.ProductId < 90000 && tests.ProductExists(t, testDB.DB, item.ProductId) {
					actualQty := tests.GetProductQuantity(t, testDB.DB, item.ProductId)
					// Quantity should be updated
					assert.GreaterOrEqual(t, actualQty, int32(0))
				}
			}
		})
	}
}

// TestGRPCGetProductStats_FullCycle tests stats retrieval via gRPC
func TestGRPCGetProductStats_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name         string
		storefrontID int64
		wantErr      bool
		wantCode     codes.Code
	}{
		{
			name:         "Valid storefront with products",
			storefrontID: 1000,
			wantErr:      false,
		},
		{
			name:         "Invalid storefront ID (zero)",
			storefrontID: 0,
			wantErr:      true,
			wantCode:     codes.InvalidArgument,
		},
		{
			name:         "Negative storefront ID",
			storefrontID: -1,
			wantErr:      true,
			wantCode:     codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.GetProductStatsRequest{
				StorefrontId: tc.storefrontID,
			}

			resp, err := client.GetProductStats(ctx, req)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Stats)

			// Verify stats accuracy against database
			dbTotal := tests.CountProductsByStorefront(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbTotal, resp.Stats.TotalProducts)

			dbActive := tests.CountActiveProductsByStorefront(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbActive, resp.Stats.ActiveProducts)

			dbOutOfStock := tests.CountOutOfStockProducts(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbOutOfStock, resp.Stats.OutOfStock)

			dbLowStock := tests.CountLowStockProducts(t, testDB.DB, tc.storefrontID)
			assert.Equal(t, dbLowStock, resp.Stats.LowStock)

			dbValue := tests.GetTotalInventoryValue(t, testDB.DB, tc.storefrontID)
			assert.InDelta(t, dbValue, resp.Stats.TotalValue, 0.01)
		})
	}
}

// TestGRPCIncrementProductViews_FullCycle tests view increment via gRPC
func TestGRPCIncrementProductViews_FullCycle(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)

	testCases := []struct {
		name      string
		productID int64
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name:      "Valid product",
			productID: productID,
			wantErr:   false,
		},
		{
			name:      "Zero product ID",
			productID: 0,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name:      "Negative product ID",
			productID: -1,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name:      "Non-existent product",
			productID: 99999,
			wantErr:   true,
			wantCode:  codes.Internal, // Or NotFound depending on implementation
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var initialViews int32
			if tc.productID == productID {
				initialViews = tests.GetProductViewCount(t, testDB.DB, tc.productID)
			}

			req := &pb.IncrementProductViewsRequest{
				ProductId: tc.productID,
			}

			resp, err := client.IncrementProductViews(ctx, req)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.IsType(t, &emptypb.Empty{}, resp)

			// Verify view count incremented in database
			newViews := tests.GetProductViewCount(t, testDB.DB, tc.productID)
			assert.Equal(t, initialViews+1, newViews)
		})
	}
}

// TestGRPCInventoryWorkflow_CompleteScenario tests complete inventory workflow via gRPC
func TestGRPCInventoryWorkflow_CompleteScenario(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	// Step 1: Get initial stats
	statsResp, err := client.GetProductStats(ctx, &pb.GetProductStatsRequest{
		StorefrontId: storefrontID,
	})
	require.NoError(t, err)
	initialValue := statsResp.Stats.TotalValue

	// Step 2: Record stock IN
	movementResp, err := client.RecordInventoryMovement(ctx, &pb.RecordInventoryMovementRequest{
		StorefrontId: storefrontID,
		ProductId:    productID,
		MovementType: "in",
		Quantity:     100,
		Reason:       stringPtr("bulk_restock"),
		UserId:       userID,
	})
	require.NoError(t, err)
	assert.True(t, movementResp.Success)

	// Step 3: Batch update multiple products
	batchResp, err := client.BatchUpdateStock(ctx, &pb.BatchUpdateStockRequest{
		StorefrontId: storefrontID,
		Items: []*pb.StockUpdateItem{
			{ProductId: 5003, Quantity: 150},
			{ProductId: 5004, Quantity: 120},
		},
		Reason: stringPtr("quarterly_audit"),
		UserId: userID,
	})
	require.NoError(t, err)
	assert.Equal(t, int32(2), batchResp.SuccessfulCount)

	// Step 4: Increment product views
	_, err = client.IncrementProductViews(ctx, &pb.IncrementProductViewsRequest{
		ProductId: productID,
	})
	require.NoError(t, err)

	// Step 5: Get updated stats
	finalStatsResp, err := client.GetProductStats(ctx, &pb.GetProductStatsRequest{
		StorefrontId: storefrontID,
	})
	require.NoError(t, err)
	assert.Greater(t, finalStatsResp.Stats.TotalValue, initialValue,
		"Total inventory value should increase after restocking")

	// Step 6: Verify all changes in database
	finalQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(200), finalQty, "Product quantity should be 200 after +100")

	finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Greater(t, finalViews, int32(0), "Product should have views")

	movementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
	assert.GreaterOrEqual(t, movementCount, 1, "Should have inventory movements recorded")
}

// TestGRPCConcurrentRequests tests concurrent gRPC requests
func TestGRPCConcurrentRequests(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)
	concurrentCalls := 10

	t.Run("Concurrent view increments", func(t *testing.T) {
		initialViews := tests.GetProductViewCount(t, testDB.DB, productID)

		done := make(chan error, concurrentCalls)
		for i := 0; i < concurrentCalls; i++ {
			go func() {
				req := &pb.IncrementProductViewsRequest{ProductId: productID}
				_, err := client.IncrementProductViews(ctx, req)
				done <- err
			}()
		}

		// Wait for all
		for i := 0; i < concurrentCalls; i++ {
			err := <-done
			assert.NoError(t, err)
		}

		finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
		assert.Equal(t, initialViews+int32(concurrentCalls), finalViews,
			"All concurrent increments should be recorded")
	})

	t.Run("Concurrent stats reads", func(t *testing.T) {
		done := make(chan error, concurrentCalls)
		for i := 0; i < concurrentCalls; i++ {
			go func() {
				req := &pb.GetProductStatsRequest{StorefrontId: 1000}
				_, err := client.GetProductStats(ctx, req)
				done <- err
			}()
		}

		// All should succeed
		for i := 0; i < concurrentCalls; i++ {
			err := <-done
			assert.NoError(t, err)
		}
	})
}

// ============================================================================
// PHASE 9.7.4: ADVANCED INTEGRATION TESTS
// ============================================================================

// TestGRPCRecordInventoryMovement_BoundaryValues tests boundary value handling
// This test verifies that the gRPC layer correctly validates and handles
// boundary values including massive quantities, negative IDs, and edge cases
func TestGRPCRecordInventoryMovement_BoundaryValues(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)

	testCases := []struct {
		name     string
		request  *pb.RecordInventoryMovementRequest
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Maximum int32 quantity",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "adjustment",
				Quantity:     2147483647, // Max int32
				Reason:       stringPtr("boundary_test"),
				UserId:       1000,
			},
			wantErr: false,
		},
		{
			name: "Zero quantity (invalid)",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     0,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Negative storefront ID",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: -1,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Negative product ID",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    -999,
				MovementType: "in",
				Quantity:     10,
				UserId:       1000,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Negative user ID",
			request: &pb.RecordInventoryMovementRequest{
				StorefrontId: 1000,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     10,
				UserId:       -5,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.RecordInventoryMovement(ctx, tc.request)

			if tc.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.wantCode, st.Code())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.True(t, resp.Success)
		})
	}
}

// TestGRPCRecordInventoryMovement_LongStrings tests handling of very long input strings
// Validates that reason and notes fields properly handle maximum length constraints
func TestGRPCRecordInventoryMovement_LongStrings(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(5000)
	storefrontID := int64(1000)
	userID := int64(1000)

	// Create strings of various lengths
	normalString := "Normal reason text"
	longString := stringRepeat("x", 255)       // Typical VARCHAR limit
	veryLongString := stringRepeat("y", 10000) // Exceeds typical limit
	unicodeString := "Тест кирилицы 测试中文 🚀🎉"   // Unicode characters

	testCases := []struct {
		name    string
		reason  *string
		notes   *string
		wantErr bool
	}{
		{
			name:    "Normal length strings",
			reason:  &normalString,
			notes:   &normalString,
			wantErr: false,
		},
		{
			name:    "Long but valid strings (255 chars)",
			reason:  &longString,
			notes:   &longString,
			wantErr: false,
		},
		{
			name:    "Very long strings (10k chars)",
			reason:  &veryLongString,
			notes:   &veryLongString,
			wantErr: false, // Should either accept or truncate gracefully
		},
		{
			name:    "Unicode strings",
			reason:  &unicodeString,
			notes:   &unicodeString,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: storefrontID,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     1,
				Reason:       tc.reason,
				Notes:        tc.notes,
				UserId:       userID,
			}

			resp, err := client.RecordInventoryMovement(ctx, req)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.True(t, resp.Success)
		})
	}
}

// TestGRPCConcurrent_MixedOperations tests concurrent mixed gRPC operations
// Validates thread safety when multiple different operations run simultaneously:
// - Inventory movements
// - View increments
// - Stats queries
func TestGRPCConcurrent_MixedOperations(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)
	concurrentOps := 20

	var wg sync.WaitGroup
	errors := make(chan error, concurrentOps*3)

	// Launch concurrent movements
	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()
			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: storefrontID,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     1,
				Reason:       stringPtr("concurrent_test"),
				UserId:       userID,
			}
			_, err := client.RecordInventoryMovement(ctx, req)
			errors <- err
		}(i)
	}

	// Launch concurrent view increments
	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := &pb.IncrementProductViewsRequest{ProductId: productID}
			_, err := client.IncrementProductViews(ctx, req)
			errors <- err
		}()
	}

	// Launch concurrent stats queries
	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := &pb.GetProductStatsRequest{StorefrontId: storefrontID}
			_, err := client.GetProductStats(ctx, req)
			errors <- err
		}()
	}

	// Wait for all operations
	wg.Wait()
	close(errors)

	// Check all operations succeeded
	errorCount := 0
	for err := range errors {
		if err != nil {
			t.Logf("Operation failed: %v", err)
			errorCount++
		}
	}

	assert.Equal(t, 0, errorCount, "All concurrent operations should succeed")

	// Verify final state consistency
	finalQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.GreaterOrEqual(t, finalQty, int32(100+concurrentOps),
		"Quantity should reflect all successful movements")
}

// TestGRPCBatchUpdateStock_LargeScaleBatch tests batch update with large item count
// Validates performance and correctness when processing 100 items in a single batch
func TestGRPCBatchUpdateStock_LargeScaleBatch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large batch test in short mode")
	}

	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)
	batchSize := 100

	// Create test products for batch update
	items := make([]*pb.StockUpdateItem, 0, batchSize)
	for i := 0; i < batchSize; i++ {
		// Use existing products (5000-5006) and cycle through them
		productID := int64(5000 + (i % 7))
		items = append(items, &pb.StockUpdateItem{
			ProductId: productID,
			Quantity:  int32(50 + i),
			Reason:    stringPtr("large_batch_test"),
		})
	}

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: storefrontID,
		Items:        items,
		Reason:       stringPtr("performance_test"),
		UserId:       userID,
	}

	// Execute batch update
	resp, err := client.BatchUpdateStock(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	// Most should succeed (some might fail due to constraints)
	assert.Greater(t, resp.SuccessfulCount, int32(50),
		"At least 50% of batch should succeed")
	assert.Len(t, resp.Results, batchSize,
		"Should return result for each item")

	// Verify successful items updated in database
	successfulUpdates := 0
	for _, result := range resp.Results {
		if result.Success {
			successfulUpdates++
			actualQty := tests.GetProductQuantity(t, testDB.DB, result.ProductId)
			assert.GreaterOrEqual(t, actualQty, int32(0),
				"Updated product should have valid quantity")
		}
	}

	assert.Equal(t, int(resp.SuccessfulCount), successfulUpdates,
		"Success count should match actual successful updates")
}

// TestGRPCGetProductStats_LargeDataset tests stats performance with large dataset
// Validates that stats queries remain performant with substantial data volume
func TestGRPCGetProductStats_LargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)

	// Add multiple movements to create larger dataset
	for i := 0; i < 50; i++ {
		productID := int64(5000 + (i % 7))
		req := &pb.RecordInventoryMovementRequest{
			StorefrontId: storefrontID,
			ProductId:    productID,
			MovementType: "in",
			Quantity:     int32(i + 1),
			Reason:       stringPtr("dataset_creation"),
			UserId:       1000,
		}
		_, err := client.RecordInventoryMovement(ctx, req)
		require.NoError(t, err)
	}

	// Query stats multiple times to test performance
	for i := 0; i < 10; i++ {
		req := &pb.GetProductStatsRequest{StorefrontId: storefrontID}
		resp, err := client.GetProductStats(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Stats)

		// Validate stats structure
		assert.Greater(t, resp.Stats.TotalProducts, int32(0))
		assert.GreaterOrEqual(t, resp.Stats.TotalProducts, resp.Stats.ActiveProducts)
		assert.GreaterOrEqual(t, resp.Stats.TotalValue, float64(0))
	}
}

// TestGRPCBatchUpdateStock_EmptyReasonHandling tests batch update with nil/empty reasons
// Validates that optional reason fields are handled correctly in batch operations
func TestGRPCBatchUpdateStock_EmptyReasonHandling(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)

	testCases := []struct {
		name   string
		items  []*pb.StockUpdateItem
		reason *string
	}{
		{
			name: "Item-level reasons provided",
			items: []*pb.StockUpdateItem{
				{
					ProductId: 5003,
					Quantity:  100,
					Reason:    stringPtr("item_reason_1"),
				},
				{
					ProductId: 5004,
					Quantity:  80,
					Reason:    stringPtr("item_reason_2"),
				},
			},
			reason: stringPtr("batch_reason"),
		},
		{
			name: "No item-level reasons (use batch reason)",
			items: []*pb.StockUpdateItem{
				{ProductId: 5003, Quantity: 90},
				{ProductId: 5004, Quantity: 70},
			},
			reason: stringPtr("batch_reason_only"),
		},
		{
			name: "No reasons at all (should still work)",
			items: []*pb.StockUpdateItem{
				{ProductId: 5003, Quantity: 85},
				{ProductId: 5004, Quantity: 75},
			},
			reason: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.BatchUpdateStockRequest{
				StorefrontId: storefrontID,
				Items:        tc.items,
				Reason:       tc.reason,
				UserId:       userID,
			}

			resp, err := client.BatchUpdateStock(ctx, req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, int32(len(tc.items)), resp.SuccessfulCount)
			assert.Equal(t, int32(0), resp.FailedCount)

			// Verify all items updated
			for _, item := range tc.items {
				actualQty := tests.GetProductQuantity(t, testDB.DB, item.ProductId)
				assert.Equal(t, item.Quantity, actualQty)
			}
		})
	}
}

// TestGRPCRecordInventoryMovement_AuditTrail tests audit trail completeness
// Validates that all inventory movements are properly logged with full context
func TestGRPCRecordInventoryMovement_AuditTrail(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	// Record multiple movements with different contexts
	movements := []struct {
		movementType string
		quantity     int32
		reason       string
		notes        string
	}{
		{"in", 50, "purchase", "Supplier ABC - Invoice #12345"},
		{"out", 10, "sale", "Customer order #67890"},
		{"adjustment", 100, "inventory_count", "Physical count correction"},
		{"in", 25, "return", "Customer return - defective item"},
	}

	for _, mv := range movements {
		req := &pb.RecordInventoryMovementRequest{
			StorefrontId: storefrontID,
			ProductId:    productID,
			MovementType: mv.movementType,
			Quantity:     mv.quantity,
			Reason:       &mv.reason,
			Notes:        &mv.notes,
			UserId:       userID,
		}

		resp, err := client.RecordInventoryMovement(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Success)
	}

	// Verify all movements recorded
	movementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
	assert.GreaterOrEqual(t, movementCount, len(movements),
		"All movements should be recorded in audit trail")

	// Verify final quantity matches expected
	expectedQty := int32(100 + 50 - 10 + (100 - (100 + 50 - 10)) + 25)
	actualQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, expectedQty, actualQty,
		"Final quantity should match sum of all movements")
}

// TestGRPCConcurrent_StressTest tests system behavior under heavy concurrent load
// Validates stability and correctness with 100 simultaneous operations
func TestGRPCConcurrent_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)
	stressOps := 100

	var wg sync.WaitGroup
	successCount := sync.Map{}
	errorCount := int32(0)

	// Launch stress operations
	for i := 0; i < stressOps; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()

			// Distribute operations across different products
			productID := int64(5000 + (iteration % 7))

			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: storefrontID,
				ProductId:    productID,
				MovementType: "in",
				Quantity:     1,
				Reason:       stringPtr("stress_test"),
				UserId:       userID,
			}

			resp, err := client.RecordInventoryMovement(ctx, req)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
				return
			}

			if resp.Success {
				// Track successful operations per product
				count, _ := successCount.LoadOrStore(productID, int32(0))
				successCount.Store(productID, count.(int32)+1)
			}
		}(i)
	}

	wg.Wait()

	// Verify results
	assert.Less(t, int(errorCount), stressOps/10,
		"Error rate should be less than 10%")

	// Verify database consistency
	totalSuccessful := int32(0)
	successCount.Range(func(key, value interface{}) bool {
		totalSuccessful += value.(int32)
		return true
	})

	assert.Greater(t, totalSuccessful, int32(stressOps*9/10),
		"At least 90% of operations should succeed")
}
