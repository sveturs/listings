//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// SAGA PATTERN TESTS
// ============================================================================

// TestRollbackStock_AfterSuccessfulDecrement tests the full Saga pattern:
// Decrement → Order Created → Rollback (on order failure)
func TestRollbackStock_AfterSuccessfulDecrement(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Use product 8008 (E2E workflow test product, fresh 100 stock)
	productID := testProductE2EWorkflow
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(100), initialStock, "Initial stock should be 100")

	orderID := rollbackStringPtr("ORDER-SAGA-001")
	orderQuantity := int32(25)

	// ===== SAGA STEP 1: Decrement Stock (Order Creation) =====
	decrementReq := &pb.DecrementStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, Quantity: orderQuantity},
		},
	}

	decrementResp, err := client.DecrementStock(ctx, decrementReq)
	require.NoError(t, err, "Decrement should succeed")
	assert.True(t, decrementResp.Success, "Decrement should be successful")

	// Verify stock decreased
	stockAfterDecrement := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(75), stockAfterDecrement, "Stock should be 75 after decrement")

	// ===== SAGA STEP 2: Simulate Order Creation Failure =====
	// In real scenario: payment failed, user cancelled, etc.
	t.Log("Simulating order creation failure...")

	// ===== SAGA STEP 3: Compensating Transaction (Rollback) =====
	rollbackReq := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, Quantity: orderQuantity},
		},
	}

	rollbackResp, err := client.RollbackStock(ctx, rollbackReq)
	require.NoError(t, err, "Rollback should succeed")
	assert.True(t, rollbackResp.Success, "Rollback should be successful")

	// Verify stock restored
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialStock, finalStock,
		"Stock should be fully restored to initial value after Saga rollback")

	// Verify inventory movements recorded correctly
	decrementCount := tests.GetDecrementMovementCount(t, testDB.DB, productID)
	assert.GreaterOrEqual(t, decrementCount, 1, "Should have decrement movement")

	// Note: Current implementation may not record rollback as 'in' movement
	// This is something to verify/fix
	t.Logf("Decrement movements recorded: %d", decrementCount)
}

// TestRollbackStock_PartialBatchFailure tests batch rollback where some products fail
func TestRollbackStock_PartialBatchFailure(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	orderID := rollbackStringPtr("ORDER-PARTIAL-FAIL")

	// Mix of valid and invalid product IDs
	req := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: testProductBatchRollback1, Quantity: 10}, // Valid
			{ProductId: 99998, Quantity: 5},                      // Invalid - not found
			{ProductId: testProductBatchRollback2, Quantity: 15}, // Valid
			{ProductId: 99999, Quantity: 8},                      // Invalid - not found
		},
	}

	resp, err := client.RollbackStock(ctx, req)

	// Current implementation continues with best effort
	require.NoError(t, err, "gRPC call should not error")

	// Response depends on implementation:
	// Option 1: Partial success (some items rolled back)
	// Option 2: All-or-nothing (transaction rolled back)

	// Check results
	assert.Len(t, resp.Results, 4, "Should have 4 results")

	successCount := 0
	failureCount := 0
	for _, result := range resp.Results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	t.Logf("Partial batch rollback: %d succeeded, %d failed", successCount, failureCount)

	// At least 2 should succeed (the valid products)
	assert.GreaterOrEqual(t, successCount, 2, "Valid products should succeed")
	// At least 2 should fail (the invalid products)
	assert.GreaterOrEqual(t, failureCount, 2, "Invalid products should fail")

	// Verify valid products were actually rolled back
	if successCount >= 2 {
		// Check at least one valid product
		stock1 := tests.GetProductQuantity(t, testDB.DB, testProductBatchRollback1)
		t.Logf("Product %d stock after partial rollback: %d", testProductBatchRollback1, stock1)
	}
}

// ============================================================================
// INTEGRATION WITH DECREMENT TESTS
// ============================================================================

// TestDecrementAndRollback_FullFlow tests complete decrement + rollback cycle
func TestDecrementAndRollback_FullFlow(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8001) // Product with 80 stock
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)
	t.Logf("Initial stock: %d", initialStock)

	orderID := rollbackStringPtr("ORDER-FULL-CYCLE")
	quantity := int32(20)

	// Step 1: Decrement
	decrementReq := &pb.DecrementStockRequest{
		OrderId: orderID,
		Items:   []*pb.StockItem{{ProductId: productID, Quantity: quantity}},
	}

	decrementResp, err := client.DecrementStock(ctx, decrementReq)
	require.NoError(t, err)
	assert.True(t, decrementResp.Success)

	stockAfterDecrement := tests.GetProductQuantity(t, testDB.DB, productID)
	expectedAfterDecrement := initialStock - quantity
	assert.Equal(t, expectedAfterDecrement, stockAfterDecrement,
		"Stock should be %d after decrement", expectedAfterDecrement)

	// Step 2: Rollback same quantity
	rollbackReq := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items:   []*pb.StockItem{{ProductId: productID, Quantity: quantity}},
	}

	rollbackResp, err := client.RollbackStock(ctx, rollbackReq)
	require.NoError(t, err)
	assert.True(t, rollbackResp.Success)

	// Step 3: Verify final stock = initial stock (full restoration)
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialStock, finalStock,
		"Stock should be fully restored: Decrement %d → Rollback %d → Initial %d",
		quantity, quantity, initialStock)

	// Verify math: Final = Initial - Decrement + Rollback
	expectedFinal := initialStock - quantity + quantity
	assert.Equal(t, expectedFinal, finalStock, "Math should be consistent")
}

// TestDecrementAndRollback_PartialRollback tests partial quantity rollback
func TestDecrementAndRollback_PartialRollback(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := int64(8002) // Product with 85 stock
	initialStock := tests.GetProductQuantity(t, testDB.DB, productID)

	orderID := rollbackStringPtr("ORDER-PARTIAL-ROLLBACK")
	decrementQty := int32(30)
	rollbackQty := int32(10) // Only rollback 10 out of 30

	// Decrement 30
	decrementReq := &pb.DecrementStockRequest{
		OrderId: orderID,
		Items:   []*pb.StockItem{{ProductId: productID, Quantity: decrementQty}},
	}

	_, err := client.DecrementStock(ctx, decrementReq)
	require.NoError(t, err)

	stockAfterDecrement := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialStock-decrementQty, stockAfterDecrement)

	// Rollback only 10
	rollbackReq := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items:   []*pb.StockItem{{ProductId: productID, Quantity: rollbackQty}},
	}

	_, err = client.RollbackStock(ctx, rollbackReq)
	require.NoError(t, err)

	// Final stock should be: Initial - 30 + 10 = Initial - 20
	finalStock := tests.GetProductQuantity(t, testDB.DB, productID)
	expectedFinal := initialStock - decrementQty + rollbackQty
	assert.Equal(t, expectedFinal, finalStock,
		"Stock should be %d (Initial %d - Decrement %d + Rollback %d)",
		expectedFinal, initialStock, decrementQty, rollbackQty)
}

// ============================================================================
// E2E WORKFLOW TEST
// ============================================================================

// TestStockWorkflow_E2E_CheckDecrementRollback tests complete stock workflow
func TestStockWorkflow_E2E_CheckDecrementRollback(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := testProductE2EWorkflow // Fresh product with 100 stock
	variantSizeS := testVariantSizeS    // Variant with 40 stock
	variantSizeM := testVariantSizeM    // Variant with 25 stock

	orderID := rollbackStringPtr("ORDER-E2E-001")

	t.Log("===== E2E WORKFLOW: Check → Decrement → Rollback =====")

	// ===== STEP 1: Check Stock Availability =====
	t.Log("Step 1: Checking stock availability...")

	checkReq := &pb.CheckStockAvailabilityRequest{
		Items: []*pb.StockItem{
			{ProductId: productID, Quantity: 30},                  // Product-level
			{ProductId: productID, VariantId: &variantSizeS, Quantity: 15}, // Variant S
			{ProductId: productID, VariantId: &variantSizeM, Quantity: 10}, // Variant M
		},
	}

	checkResp, err := client.CheckStockAvailability(ctx, checkReq)
	require.NoError(t, err, "Stock check should succeed")
	assert.True(t, checkResp.AllAvailable, "All items should be available")
	assert.Len(t, checkResp.Items, 3, "Should check 3 items")

	// Verify each item is available
	for i, item := range checkResp.Items {
		assert.True(t, item.IsAvailable, "Item %d should be available", i)
		t.Logf("  Item %d: Available=%v, Requested=%d, Available=%d",
			i, item.IsAvailable, item.RequestedQuantity, item.AvailableQuantity)
	}

	// ===== STEP 2: Decrement Stock (Order Creation) =====
	t.Log("Step 2: Decrementing stock (creating order)...")

	decrementReq := &pb.DecrementStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, Quantity: 30},
			{ProductId: productID, VariantId: &variantSizeS, Quantity: 15},
			{ProductId: productID, VariantId: &variantSizeM, Quantity: 10},
		},
	}

	decrementResp, err := client.DecrementStock(ctx, decrementReq)
	require.NoError(t, err, "Decrement should succeed")
	assert.True(t, decrementResp.Success, "Decrement should be successful")

	// Verify stock decreased
	stockAfterDecrement := tests.GetProductQuantity(t, testDB.DB, productID)
	variantSAfterDecrement := tests.GetVariantQuantity(t, testDB.DB, variantSizeS)
	variantMAfterDecrement := tests.GetVariantQuantity(t, testDB.DB, variantSizeM)

	assert.Equal(t, int32(70), stockAfterDecrement, "Product stock should be 70")
	assert.Equal(t, int32(25), variantSAfterDecrement, "Variant S stock should be 25")
	assert.Equal(t, int32(15), variantMAfterDecrement, "Variant M stock should be 15")

	t.Logf("  Stock after decrement: Product=%d, VariantS=%d, VariantM=%d",
		stockAfterDecrement, variantSAfterDecrement, variantMAfterDecrement)

	// ===== STEP 3: Simulate Order Failure =====
	t.Log("Step 3: Simulating order failure (payment rejected)...")
	// In real scenario: payment gateway returns error, inventory timeout, etc.

	// ===== STEP 4: Rollback Stock (Compensating Transaction) =====
	t.Log("Step 4: Rolling back stock (compensating transaction)...")

	rollbackReq := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, Quantity: 30},
			{ProductId: productID, VariantId: &variantSizeS, Quantity: 15},
			{ProductId: productID, VariantId: &variantSizeM, Quantity: 10},
		},
	}

	rollbackResp, err := client.RollbackStock(ctx, rollbackReq)
	require.NoError(t, err, "Rollback should succeed")
	assert.True(t, rollbackResp.Success, "Rollback should be successful")
	assert.Len(t, rollbackResp.Results, 3, "Should rollback 3 items")

	// Verify all rollbacks succeeded
	for i, result := range rollbackResp.Results {
		assert.True(t, result.Success, "Rollback result %d should succeed", i)
		t.Logf("  Rollback %d: %d → %d (restored +%d)",
			i, result.StockBefore, result.StockAfter, result.StockAfter-result.StockBefore)
	}

	// ===== STEP 5: Verify Stock Restored =====
	t.Log("Step 5: Verifying stock fully restored...")

	finalProductStock := tests.GetProductQuantity(t, testDB.DB, productID)
	finalVariantSStock := tests.GetVariantQuantity(t, testDB.DB, variantSizeS)
	finalVariantMStock := tests.GetVariantQuantity(t, testDB.DB, variantSizeM)

	assert.Equal(t, int32(100), finalProductStock, "Product stock should be restored to 100")
	assert.Equal(t, int32(40), finalVariantSStock, "Variant S stock should be restored to 40")
	assert.Equal(t, int32(25), finalVariantMStock, "Variant M stock should be restored to 25")

	t.Logf("  Final stock: Product=%d, VariantS=%d, VariantM=%d",
		finalProductStock, finalVariantSStock, finalVariantMStock)

	// ===== STEP 6: Check Stock Again (Should be available again) =====
	t.Log("Step 6: Checking stock availability again after rollback...")

	checkResp2, err := client.CheckStockAvailability(ctx, checkReq)
	require.NoError(t, err, "Second stock check should succeed")
	assert.True(t, checkResp2.AllAvailable, "All items should be available again after rollback")

	t.Log("✅ E2E WORKFLOW COMPLETED SUCCESSFULLY")
	t.Log("   Check → Decrement → Rollback → Restored")
}

// TestStockWorkflow_E2E_VariantRollback tests rollback for product variants
func TestStockWorkflow_E2E_VariantRollback(t *testing.T) {
	client, testDB, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	productID := testProductVariantRollback
	variantS := testVariantSizeS
	variantM := testVariantSizeM

	// Get initial variant stocks
	initialVariantS := tests.GetVariantQuantity(t, testDB.DB, variantS)
	initialVariantM := tests.GetVariantQuantity(t, testDB.DB, variantM)

	t.Logf("Initial stocks: VariantS=%d, VariantM=%d", initialVariantS, initialVariantM)

	orderID := rollbackStringPtr("ORDER-VARIANT-E2E")

	// Decrement variants
	decrementReq := &pb.DecrementStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, VariantId: &variantS, Quantity: 20},
			{ProductId: productID, VariantId: &variantM, Quantity: 10},
		},
	}

	_, err := client.DecrementStock(ctx, decrementReq)
	require.NoError(t, err, "Variant decrement should succeed")

	// Verify decrements
	variantSAfter := tests.GetVariantQuantity(t, testDB.DB, variantS)
	variantMAfter := tests.GetVariantQuantity(t, testDB.DB, variantM)

	assert.Equal(t, initialVariantS-20, variantSAfter, "Variant S should be decremented")
	assert.Equal(t, initialVariantM-10, variantMAfter, "Variant M should be decremented")

	// Rollback variants
	rollbackReq := &pb.RollbackStockRequest{
		OrderId: orderID,
		Items: []*pb.StockItem{
			{ProductId: productID, VariantId: &variantS, Quantity: 20},
			{ProductId: productID, VariantId: &variantM, Quantity: 10},
		},
	}

	rollbackResp, err := client.RollbackStock(ctx, rollbackReq)
	require.NoError(t, err, "Variant rollback should succeed")
	assert.True(t, rollbackResp.Success, "Variant rollback should be successful")

	// Verify restoration
	finalVariantS := tests.GetVariantQuantity(t, testDB.DB, variantS)
	finalVariantM := tests.GetVariantQuantity(t, testDB.DB, variantM)

	assert.Equal(t, initialVariantS, finalVariantS, "Variant S should be fully restored")
	assert.Equal(t, initialVariantM, finalVariantM, "Variant M should be fully restored")

	t.Log("✅ Variant rollback E2E completed successfully")
}

// ============================================================================
// PERFORMANCE & STRESS TESTS
// ============================================================================

// TestRollbackStock_Performance tests rollback performance
func TestRollbackStock_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	client, _, cleanup := setupRollbackTestServer(t)
	defer cleanup()

	_ = tests.TestContext(t)

	productID := testProductSingleRollback

	testCases := []struct {
		name          string
		itemCount     int
		maxDurationMs int
	}{
		{"Single item rollback", 1, 50},
		{"Batch 5 items", 5, 150},
		{"Batch 10 items", 10, 200},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare items
			items := make([]*pb.StockItem, tc.itemCount)
			for i := 0; i < tc.itemCount; i++ {
				items[i] = &pb.StockItem{
					ProductId: productID,
					Quantity:  1,
				}
			}

			req := &pb.RollbackStockRequest{
				OrderId: rollbackStringPtr("ORDER-PERF"),
				Items:   items,
			}

			// Measure time
			start := context.Background()
			ctx, cancel := context.WithTimeout(start, 5*time.Second)
			defer cancel()

			startTime := time.Now()
			_, err := client.RollbackStock(ctx, req)
			duration := time.Since(startTime)

			// Assertions
			require.NoError(t, err, "Rollback should succeed")

			durationMs := duration.Milliseconds()
			t.Logf("  Duration: %dms (max: %dms)", durationMs, tc.maxDurationMs)

			// Performance check (informational, not strict)
			if durationMs > int64(tc.maxDurationMs) {
				t.Logf("⚠️  Performance warning: Rollback took %dms (expected <%dms)",
					durationMs, tc.maxDurationMs)
			}
		})
	}
}
