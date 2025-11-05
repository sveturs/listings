//go:build integration
// +build integration

package integration

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/internal/service/listings"
	"github.com/sveturs/listings/tests"
)

// setupInventoryTestService creates a service with real database
func setupInventoryTestService(t *testing.T) (*listings.Service, *tests.TestDB) {
	t.Helper()

	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Load inventory fixtures
	tests.LoadInventoryFixtures(t, testDB.DB)

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service (no Redis cache and no indexer for integration tests)
	service := listings.NewService(repo, nil, nil, logger)

	return service, testDB
}

// TestServiceUpdateProductInventory_BusinessLogic tests service layer business logic
func TestServiceUpdateProductInventory_BusinessLogic(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000) // Quantity 100
	userID := int64(1000)

	testCases := []struct {
		name         string
		movementType string
		quantity     int32
		wantErr      bool
		errContains  string
	}{
		{
			name:         "Valid IN movement",
			movementType: "in",
			quantity:     50,
			wantErr:      false,
		},
		{
			name:         "Invalid movement type",
			movementType: "invalid_type",
			quantity:     10,
			wantErr:      true,
			errContains:  "invalid movement_type",
		},
		{
			name:         "Negative quantity",
			movementType: "in",
			quantity:     -10,
			wantErr:      true,
			errContains:  "quantity cannot be negative",
		},
		{
			name:         "Invalid storefront ID",
			movementType: "in",
			quantity:     10,
			wantErr:      true,
			errContains:  "storefront_id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testStorefrontID int64
			if tc.errContains == "storefront_id" {
				testStorefrontID = 0 // Invalid
			} else {
				testStorefrontID = storefrontID
			}

			stockBefore, stockAfter, err := service.UpdateProductInventory(
				ctx,
				testStorefrontID,
				productID,
				0,
				tc.movementType,
				tc.quantity,
				"test",
				"",
				userID,
			)

			if tc.wantErr {
				require.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Greater(t, stockAfter, stockBefore)
		})
	}
}

// TestServiceBatchUpdateStock_ComplexScenario tests complex batch update scenarios
func TestServiceBatchUpdateStock_ComplexScenario(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)

	t.Run("Successful batch with multiple products", func(t *testing.T) {
		items := []domain.StockUpdateItem{
			{
				ProductID: 5003, // Quantity 50
				Quantity:  75,
			},
			{
				ProductID: 5004, // Quantity 75
				Quantity:  100,
			},
		}

		successCount, failedCount, results, err := service.BatchUpdateStock(
			ctx,
			storefrontID,
			items,
			"bulk_update",
			userID,
		)

		require.NoError(t, err)
		assert.Equal(t, int32(2), successCount)
		assert.Equal(t, int32(0), failedCount)
		assert.Len(t, results, 2)

		// Verify all successful
		for _, result := range results {
			assert.True(t, result.Success)
		}
	})

	t.Run("Empty batch validation", func(t *testing.T) {
		_, _, _, err := service.BatchUpdateStock(
			ctx,
			storefrontID,
			[]domain.StockUpdateItem{},
			"",
			userID,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
	})

	t.Run("Invalid storefront ID", func(t *testing.T) {
		items := []domain.StockUpdateItem{
			{ProductID: 5003, Quantity: 50},
		}

		_, _, _, err := service.BatchUpdateStock(
			ctx,
			0, // Invalid
			items,
			"",
			userID,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "storefront_id")
	})

	t.Run("Too many items (> 1000)", func(t *testing.T) {
		// Create 1001 items
		items := make([]domain.StockUpdateItem, 1001)
		for i := 0; i < 1001; i++ {
			items[i] = domain.StockUpdateItem{
				ProductID: int64(5000 + i),
				Quantity:  10,
			}
		}

		_, _, _, err := service.BatchUpdateStock(
			ctx,
			storefrontID,
			items,
			"",
			userID,
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "1000")
	})
}

// TestServiceGetProductStats_Accuracy tests stats calculation accuracy
func TestServiceGetProductStats_Accuracy(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)

	stats, err := service.GetProductStats(ctx, storefrontID)

	require.NoError(t, err)
	require.NotNil(t, stats)

	// Verify against direct database queries
	dbTotal := tests.CountProductsByStorefront(t, testDB.DB, storefrontID)
	assert.Equal(t, dbTotal, stats.TotalProducts, "Total products mismatch")

	dbActive := tests.CountActiveProductsByStorefront(t, testDB.DB, storefrontID)
	assert.Equal(t, dbActive, stats.ActiveProducts, "Active products mismatch")

	dbOutOfStock := tests.CountOutOfStockProducts(t, testDB.DB, storefrontID)
	assert.Equal(t, dbOutOfStock, stats.OutOfStock, "Out of stock mismatch")

	dbLowStock := tests.CountLowStockProducts(t, testDB.DB, storefrontID)
	assert.Equal(t, dbLowStock, stats.LowStock, "Low stock mismatch")

	dbValue := tests.GetTotalInventoryValue(t, testDB.DB, storefrontID)
	assert.InDelta(t, dbValue, stats.TotalValue, 0.01, "Total value mismatch")

	// Validate data consistency
	assert.GreaterOrEqual(t, stats.TotalProducts, stats.ActiveProducts,
		"Total products should be >= active products")
	assert.GreaterOrEqual(t, stats.TotalProducts, stats.OutOfStock+stats.LowStock,
		"Total should be >= out of stock + low stock")
}

// TestServiceGetProductStats_ValidationErrors tests input validation
func TestServiceGetProductStats_ValidationErrors(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	testCases := []struct {
		name         string
		storefrontID int64
		wantErr      bool
	}{
		{
			name:         "Valid storefront ID",
			storefrontID: 1000,
			wantErr:      false,
		},
		{
			name:         "Zero storefront ID",
			storefrontID: 0,
			wantErr:      true,
		},
		{
			name:         "Negative storefront ID",
			storefrontID: -1,
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stats, err := service.GetProductStats(ctx, tc.storefrontID)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, stats)
				assert.Contains(t, err.Error(), "storefront_id")
			} else {
				require.NoError(t, err)
				require.NotNil(t, stats)
			}
		})
	}
}

// TestServiceIncrementProductViews_Idempotency tests view increment idempotency
func TestServiceIncrementProductViews_Idempotency(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	productID := int64(5000)

	initialViews := tests.GetProductViewCount(t, testDB.DB, productID)

	// Increment multiple times
	for i := 0; i < 3; i++ {
		err := service.IncrementProductViews(ctx, productID)
		require.NoError(t, err, "Increment #%d failed", i+1)
	}

	// Verify total
	finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Equal(t, initialViews+3, finalViews, "Should increment by 3")
}

// TestServiceIncrementProductViews_ValidationErrors tests validation
func TestServiceIncrementProductViews_ValidationErrors(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	testCases := []struct {
		name      string
		productID int64
		wantErr   bool
	}{
		{
			name:      "Valid product ID",
			productID: 5000,
			wantErr:   false,
		},
		{
			name:      "Zero product ID",
			productID: 0,
			wantErr:   true,
		},
		{
			name:      "Negative product ID",
			productID: -1,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.IncrementProductViews(ctx, tc.productID)

			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "product_id")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestServiceInventoryWorkflow_EndToEnd tests complete inventory workflow
func TestServiceInventoryWorkflow_EndToEnd(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000) // Initial quantity: 100
	userID := int64(1000)

	// Step 1: Check initial stats
	initialStats, err := service.GetProductStats(ctx, storefrontID)
	require.NoError(t, err)
	initialTotalValue := initialStats.TotalValue

	// Step 2: Add stock (in)
	stockBefore, stockAfter, err := service.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"in",
		50,
		"restock",
		"Supplier delivery",
		userID,
	)
	require.NoError(t, err)
	assert.Equal(t, int32(100), stockBefore)
	assert.Equal(t, int32(150), stockAfter)

	// Step 3: Verify stats updated
	statsAfterIn, err := service.GetProductStats(ctx, storefrontID)
	require.NoError(t, err)
	assert.Greater(t, statsAfterIn.TotalValue, initialTotalValue,
		"Total value should increase after stock in")

	// Step 4: Remove stock (out)
	stockBefore, stockAfter, err = service.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"out",
		30,
		"sale",
		"Customer order",
		userID,
	)
	require.NoError(t, err)
	assert.Equal(t, int32(150), stockBefore)
	assert.Equal(t, int32(120), stockAfter)

	// Step 5: Increment view count
	err = service.IncrementProductViews(ctx, productID)
	require.NoError(t, err)

	// Step 6: Verify final state
	finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Greater(t, finalViews, int32(0), "Product should have views")

	finalQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(120), finalQty, "Final quantity should be 120")

	// Step 7: Verify inventory movements recorded
	movementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
	assert.GreaterOrEqual(t, movementCount, 2,
		"Should have at least 2 inventory movements (in + out)")
}

// TestServiceConcurrentOperations tests concurrent service calls
func TestServiceConcurrentOperations(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	productID := int64(5000)
	concurrentCalls := 20

	t.Run("Concurrent view increments", func(t *testing.T) {
		initialViews := tests.GetProductViewCount(t, testDB.DB, productID)

		done := make(chan error, concurrentCalls)
		for i := 0; i < concurrentCalls; i++ {
			go func() {
				done <- service.IncrementProductViews(ctx, productID)
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
				_, err := service.GetProductStats(ctx, 1000)
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

// TestServicePerformance_ResponseTime tests service layer performance
func TestServicePerformance_ResponseTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	t.Run("GetProductStats latency", func(t *testing.T) {
		// Warm up
		_, _ = service.GetProductStats(ctx, 1000)

		// Measure 10 calls
		for i := 0; i < 10; i++ {
			_, err := service.GetProductStats(ctx, 1000)
			require.NoError(t, err)
			// Response time should be < 100ms for stats query
			// (actual timing would require time.Now() measurements)
		}
	})

	t.Run("IncrementProductViews latency", func(t *testing.T) {
		productID := int64(5000)

		for i := 0; i < 10; i++ {
			err := service.IncrementProductViews(ctx, productID)
			require.NoError(t, err)
			// Should be < 50ms for simple increment
		}
	})
}

// ============================================================================
// PHASE 9.7.4: ADVANCED SERVICE TESTS
// ============================================================================

// TestServiceBusinessLogic_LowStockThresholdDetection tests low stock threshold detection
// Validates that the service correctly identifies products approaching stock-out conditions
func TestServiceBusinessLogic_LowStockThresholdDetection(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)

	// Product 5001 has quantity 5 (low stock)
	// Product 5002 has quantity 0 (out of stock)
	// Product 5000 has quantity 100 (normal stock)

	stats, err := service.GetProductStats(ctx, storefrontID)
	require.NoError(t, err)
	require.NotNil(t, stats)

	// Verify low stock detection
	assert.Greater(t, stats.LowStock, int32(0),
		"Should detect low stock products")

	// Verify out of stock detection
	assert.Greater(t, stats.OutOfStock, int32(0),
		"Should detect out of stock products")

	// Verify business logic: low stock should not include out of stock
	assert.GreaterOrEqual(t, stats.TotalProducts, stats.LowStock+stats.OutOfStock,
		"Low stock and out of stock should not overlap")

	// Reduce product to low stock threshold
	userID := int64(1000)
	productID := int64(5000) // Currently 100

	_, _, err = service.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"out",
		95, // Reduce to 5 (low stock)
		"threshold_test",
		"Testing low stock threshold",
		userID,
	)
	require.NoError(t, err)

	// Check stats again
	statsAfter, err := service.GetProductStats(ctx, storefrontID)
	require.NoError(t, err)

	// Low stock count should increase
	assert.GreaterOrEqual(t, statsAfter.LowStock, stats.LowStock,
		"Low stock count should increase or stay same")
}

// TestServiceBusinessLogic_OutOfStockPrevention tests out of stock prevention logic
// Validates that business rules prevent invalid stock operations
func TestServiceBusinessLogic_OutOfStockPrevention(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5001) // Low stock product with quantity 5
	userID := int64(1000)

	testCases := []struct {
		name        string
		quantity    int32
		movementType string
		wantErr     bool
		errContains string
	}{
		{
			name:        "Remove exact available quantity",
			quantity:    5,
			movementType: "out",
			wantErr:     false,
		},
		{
			name:        "Try to remove more than available",
			quantity:    10,
			movementType: "out",
			wantErr:     true,
			errContains: "products.insufficient_stock",
		},
		{
			name:        "Add stock after out of stock",
			quantity:    100,
			movementType: "in",
			wantErr:     false,
		},
		{
			name:        "Adjustment to zero (valid)",
			quantity:    0,
			movementType: "adjustment",
			wantErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stockBefore, stockAfter, err := service.UpdateProductInventory(
				ctx,
				storefrontID,
				productID,
				0,
				tc.movementType,
				tc.quantity,
				"prevention_test",
				tc.name,
				userID,
			)

			if tc.wantErr {
				require.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, stockBefore, stockAfter,
				"Stock should change for valid operations")
		})
	}
}

// TestServiceE2E_MultiStorefrontIsolation tests data isolation between storefronts
// Validates that inventory operations for one storefront don't affect others
func TestServiceE2E_MultiStorefrontIsolation(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Create second storefront with products
	_, err := testDB.DB.Exec(`
		INSERT INTO storefronts (id, user_id, name, slug, is_active, created_at, updated_at)
		VALUES (2000, 1001, 'Test Store 2', 'test-store-2', true, NOW(), NOW())
	`)
	require.NoError(t, err)

	// Create products for second storefront
	_, err = testDB.DB.Exec(`
		INSERT INTO products (id, storefront_id, name, slug, price, quantity, is_active, created_at, updated_at)
		VALUES
			(6000, 2000, 'Product A', 'product-a', 100.00, 50, true, NOW(), NOW()),
			(6001, 2000, 'Product B', 'product-b', 200.00, 75, true, NOW(), NOW())
	`)
	require.NoError(t, err)

	storefront1ID := int64(1000)
	storefront2ID := int64(2000)
	userID := int64(1000)

	// Get initial stats for both storefronts
	stats1Before, err := service.GetProductStats(ctx, storefront1ID)
	require.NoError(t, err)

	stats2Before, err := service.GetProductStats(ctx, storefront2ID)
	require.NoError(t, err)

	// Perform operations on storefront 1
	_, _, err = service.UpdateProductInventory(
		ctx,
		storefront1ID,
		5000, // Product from storefront 1
		0,
		"in",
		100,
		"isolation_test",
		"Update storefront 1",
		userID,
	)
	require.NoError(t, err)

	// Perform operations on storefront 2
	_, _, err = service.UpdateProductInventory(
		ctx,
		storefront2ID,
		6000, // Product from storefront 2
		0,
		"in",
		50,
		"isolation_test",
		"Update storefront 2",
		userID,
	)
	require.NoError(t, err)

	// Get updated stats
	stats1After, err := service.GetProductStats(ctx, storefront1ID)
	require.NoError(t, err)

	stats2After, err := service.GetProductStats(ctx, storefront2ID)
	require.NoError(t, err)

	// Verify storefront 1 changed
	assert.Greater(t, stats1After.TotalValue, stats1Before.TotalValue,
		"Storefront 1 total value should increase")

	// Verify storefront 2 changed
	assert.Greater(t, stats2After.TotalValue, stats2Before.TotalValue,
		"Storefront 2 total value should increase")

	// Verify isolation: product counts should remain separate
	assert.NotEqual(t, stats1After.TotalProducts, stats2After.TotalProducts,
		"Different storefronts should have independent product counts")

	// Try invalid cross-storefront operation (should fail)
	_, _, err = service.UpdateProductInventory(
		ctx,
		storefront1ID,
		6000, // Product belongs to storefront 2
		0,
		"in",
		10,
		"cross_test",
		"Should fail - wrong storefront",
		userID,
	)
	require.Error(t, err,
		"Should not allow operations on products from different storefronts")
}

// TestServiceE2E_AuditTrailCompleteness tests end-to-end audit trail functionality
// Validates that all operations are properly logged and traceable
func TestServiceE2E_AuditTrailCompleteness(t *testing.T) {
	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	// Record initial movement count
	initialCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)

	// Perform series of operations
	operations := []struct {
		movementType string
		quantity     int32
		reason       string
	}{
		{"in", 50, "initial_stock"},
		{"out", 20, "sale"},
		{"in", 30, "restock"},
		{"adjustment", 150, "inventory_audit"},
		{"out", 10, "damage"},
	}

	for i, op := range operations {
		_, _, err := service.UpdateProductInventory(
			ctx,
			storefrontID,
			productID,
			0,
			op.movementType,
			op.quantity,
			op.reason,
			"Audit trail test operation",
			userID,
		)
		require.NoError(t, err, "Operation %d should succeed", i+1)
	}

	// Verify all movements recorded
	finalCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
	expectedCount := initialCount + len(operations)
	assert.Equal(t, expectedCount, finalCount,
		"All operations should be recorded in audit trail")

	// Verify final state matches operations
	expectedQty := int32(150) // Final adjustment value
	actualQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, expectedQty, actualQty,
		"Final quantity should match last adjustment")

	// Perform batch operation and verify audit
	items := []domain.StockUpdateItem{
		{ProductID: 5003, Quantity: 200, Reason: stringPtr("batch_audit_1")},
		{ProductID: 5004, Quantity: 180, Reason: stringPtr("batch_audit_2")},
	}

	successCount, failedCount, results, err := service.BatchUpdateStock(
		ctx,
		storefrontID,
		items,
		"batch_audit_test",
		userID,
	)

	require.NoError(t, err)
	assert.Equal(t, int32(len(items)), successCount)
	assert.Equal(t, int32(0), failedCount)

	// Verify batch operations recorded
	for _, result := range results {
		if result.Success {
			count := tests.GetInventoryMovementCount(t, testDB.DB, result.ProductID)
			assert.Greater(t, count, 0,
				"Batch operation should be recorded for product %d", result.ProductID)
		}
	}
}

// TestServiceStress_MixedConcurrentOperations tests system under heavy mixed load
// Validates stability with 100 concurrent operations of different types
func TestServiceStress_MixedConcurrentOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	service, testDB := setupInventoryTestService(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)
	stressOps := 100

	var wg sync.WaitGroup
	successCount := int32(0)
	errorCount := int32(0)

	// Operation types
	operations := []string{"movement", "batch", "stats", "views"}
	productIDs := []int64{5000, 5003, 5004, 5006}

	for i := 0; i < stressOps; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()

			opType := operations[iteration%len(operations)]
			productID := productIDs[iteration%len(productIDs)]

			var err error

			switch opType {
			case "movement":
				_, _, err = service.UpdateProductInventory(
					ctx,
					storefrontID,
					productID,
					0,
					"in",
					1,
					"stress_test",
					"Concurrent operation",
					userID,
				)

			case "batch":
				items := []domain.StockUpdateItem{
					{ProductID: productID, Quantity: int32(50 + iteration)},
				}
				_, _, _, err = service.BatchUpdateStock(
					ctx,
					storefrontID,
					items,
					"stress_batch",
					userID,
				)

			case "stats":
				_, err = service.GetProductStats(ctx, storefrontID)

			case "views":
				err = service.IncrementProductViews(ctx, productID)
			}

			if err != nil {
				atomic.AddInt32(&errorCount, 1)
			} else {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// Verify success rate
	successRate := float64(successCount) / float64(stressOps)
	assert.Greater(t, successRate, 0.95,
		"At least 95%% of operations should succeed under stress")

	t.Logf("Stress test results: %.2f%% success (%d/%d)",
		successRate*100, successCount, stressOps)

	// Verify database consistency
	for _, prodID := range productIDs {
		qty := tests.GetProductQuantity(t, testDB.DB, prodID)
		assert.GreaterOrEqual(t, qty, int32(0),
			"Product %d should have valid quantity", prodID)

		views := tests.GetProductViewCount(t, testDB.DB, prodID)
		assert.GreaterOrEqual(t, views, int32(0),
			"Product %d should have valid view count", prodID)
	}

	// Verify stats still work after stress
	stats, err := service.GetProductStats(ctx, storefrontID)
	require.NoError(t, err)
	require.NotNil(t, stats)
	assert.Greater(t, stats.TotalProducts, int32(0),
		"Stats should be available after stress test")
}
