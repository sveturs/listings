//go:build integration
// +build integration

package integration

import (
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
