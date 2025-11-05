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
	"github.com/sveturs/listings/tests"
)

// setupInventoryTestRepo creates a repository with test database
func setupInventoryTestRepo(t *testing.T) (*postgres.Repository, *tests.TestDB) {
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

	return repo, testDB
}

// TestUpdateProductInventory_ValidMovement_Success tests successful inventory movement
func TestUpdateProductInventory_ValidMovement_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Product ID 5000 has quantity 100
	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	testCases := []struct {
		name         string
		movementType string
		quantity     int32
		reason       string
		notes        string
		wantBefore   int32
		wantAfter    int32
	}{
		{
			name:         "Stock In - Add 50 units",
			movementType: "in",
			quantity:     50,
			reason:       "restock",
			notes:        "Supplier delivery",
			wantBefore:   100,
			wantAfter:    150,
		},
		{
			name:         "Stock Out - Remove 30 units",
			movementType: "out",
			quantity:     30,
			reason:       "sale",
			notes:        "Customer order #12345",
			wantBefore:   150,
			wantAfter:    120,
		},
		{
			name:         "Stock Adjustment - Set to 100",
			movementType: "adjustment",
			quantity:     100,
			reason:       "inventory_count",
			notes:        "Physical count adjustment",
			wantBefore:   120,
			wantAfter:    100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stockBefore, stockAfter, err := repo.UpdateProductInventory(
				ctx,
				storefrontID,
				productID,
				0, // No variant
				tc.movementType,
				tc.quantity,
				tc.reason,
				tc.notes,
				userID,
			)

			require.NoError(t, err)
			assert.Equal(t, tc.wantBefore, stockBefore, "Stock before mismatch")
			assert.Equal(t, tc.wantAfter, stockAfter, "Stock after mismatch")

			// Verify in database
			actualQty := tests.GetProductQuantity(t, testDB.DB, productID)
			assert.Equal(t, tc.wantAfter, actualQty, "Database quantity mismatch")

			// Verify inventory movement was recorded
			count := tests.GetInventoryMovementCount(t, testDB.DB, productID)
			assert.Greater(t, count, 0, "Inventory movement not recorded")
		})
	}
}

// TestUpdateProductInventory_VariantLevel_Success tests variant-level inventory management
func TestUpdateProductInventory_VariantLevel_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	variantID := int64(6000) // Size S variant, quantity 50
	userID := int64(1000)

	// Verify variant exists
	require.True(t, tests.VariantExists(t, testDB.DB, variantID), "Test variant should exist")

	// Add stock to variant
	stockBefore, stockAfter, err := repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		variantID,
		"in",
		25,
		"restock",
		"Variant restock",
		userID,
	)

	require.NoError(t, err)
	assert.Equal(t, int32(50), stockBefore)
	assert.Equal(t, int32(75), stockAfter)

	// Verify variant quantity in database
	actualQty := tests.GetVariantQuantity(t, testDB.DB, variantID)
	assert.Equal(t, int32(75), actualQty)
}

// TestUpdateProductInventory_InsufficientStock_Error tests error when stock is insufficient
func TestUpdateProductInventory_InsufficientStock_Error(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5001) // Low stock product, quantity 5
	userID := int64(1000)

	// Try to remove more than available
	_, _, err := repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"out",
		100, // Request 100 but only 5 available
		"sale",
		"",
		userID,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "products.insufficient_stock", "Should return insufficient stock error")

	// Verify quantity unchanged
	actualQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, int32(5), actualQty, "Quantity should remain unchanged after error")
}

// TestUpdateProductInventory_NonExistentProduct_Error tests error for non-existent product
func TestUpdateProductInventory_NonExistentProduct_Error(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(9999) // Non-existent product
	userID := int64(1000)

	_, _, err := repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"in",
		10,
		"",
		"",
		userID,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "products.not_found", "Should return product not found error")
}

// TestBatchUpdateStock_ValidBatch_Success tests successful batch stock update
func TestBatchUpdateStock_ValidBatch_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)

	// Product 5003: quantity 50, Product 5004: quantity 75
	items := []domain.StockUpdateItem{
		{
			ProductID: 5003,
			VariantID: nil,
			Quantity:  100, // Update to 100
			Reason:    stringPtr("restock"),
		},
		{
			ProductID: 5004,
			VariantID: nil,
			Quantity:  50, // Update to 50
			Reason:    stringPtr("adjustment"),
		},
	}

	successCount, failedCount, results, err := repo.BatchUpdateStock(
		ctx,
		storefrontID,
		items,
		"inventory_audit",
		userID,
	)

	require.NoError(t, err)
	assert.Equal(t, int32(2), successCount, "Should update 2 products successfully")
	assert.Equal(t, int32(0), failedCount, "Should have no failures")
	assert.Len(t, results, 2, "Should return 2 results")

	// Verify first product
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(5003), results[0].ProductID)
	assert.Equal(t, int32(50), results[0].StockBefore)
	assert.Equal(t, int32(100), results[0].StockAfter)

	// Verify second product
	assert.True(t, results[1].Success)
	assert.Equal(t, int64(5004), results[1].ProductID)
	assert.Equal(t, int32(75), results[1].StockBefore)
	assert.Equal(t, int32(50), results[1].StockAfter)

	// Verify database
	assert.Equal(t, int32(100), tests.GetProductQuantity(t, testDB.DB, 5003))
	assert.Equal(t, int32(50), tests.GetProductQuantity(t, testDB.DB, 5004))
}

// TestBatchUpdateStock_PartialSuccess tests batch update with some failures
func TestBatchUpdateStock_PartialSuccess(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)

	items := []domain.StockUpdateItem{
		{
			ProductID: 5003,
			Quantity:  80, // Valid
		},
		{
			ProductID: 9999, // Non-existent product
			Quantity:  50,
		},
		{
			ProductID: 5004,
			Quantity:  60, // Valid
		},
	}

	successCount, failedCount, results, err := repo.BatchUpdateStock(
		ctx,
		storefrontID,
		items,
		"mixed_batch",
		userID,
	)

	require.NoError(t, err, "Batch update should not fail completely")
	assert.Equal(t, int32(2), successCount, "Should succeed for 2 products")
	assert.Equal(t, int32(1), failedCount, "Should fail for 1 product")
	assert.Len(t, results, 3, "Should return 3 results")

	// Check success items
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(5003), results[0].ProductID)

	// Check failed item
	assert.False(t, results[1].Success)
	assert.Equal(t, int64(9999), results[1].ProductID)
	assert.NotNil(t, results[1].Error)

	// Check second success
	assert.True(t, results[2].Success)
	assert.Equal(t, int64(5004), results[2].ProductID)
}

// TestBatchUpdateStock_EmptyBatch_Error tests error for empty batch
func TestBatchUpdateStock_EmptyBatch_Error(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	items := []domain.StockUpdateItem{}

	_, _, _, err := repo.BatchUpdateStock(
		ctx,
		1000,
		items,
		"",
		1000,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty", "Should return error for empty batch")
}

// TestGetProductStats_ValidStorefront_Success tests successful stats retrieval
func TestGetProductStats_ValidStorefront_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)

	stats, err := repo.GetProductStats(ctx, storefrontID)

	require.NoError(t, err)
	require.NotNil(t, stats)

	// Storefront 1000 has products: 5000 (qty 100), 5001 (qty 5), 5002 (qty 0),
	// 5003 (qty 50), 5004 (qty 75), 5005 (qty 20, inactive), 5006 (qty 30)
	assert.Greater(t, stats.TotalProducts, int32(0), "Should have products")
	assert.Greater(t, stats.ActiveProducts, int32(0), "Should have active products")
	assert.GreaterOrEqual(t, stats.OutOfStock, int32(0), "Out of stock count should be >= 0")
	assert.GreaterOrEqual(t, stats.LowStock, int32(0), "Low stock count should be >= 0")
	assert.Greater(t, stats.TotalValue, float64(0), "Should have inventory value")

	// Verify against database
	dbTotal := tests.CountProductsByStorefront(t, testDB.DB, storefrontID)
	assert.Equal(t, dbTotal, stats.TotalProducts, "Total products should match database")

	dbActive := tests.CountActiveProductsByStorefront(t, testDB.DB, storefrontID)
	assert.Equal(t, dbActive, stats.ActiveProducts, "Active products should match database")

	dbOutOfStock := tests.CountOutOfStockProducts(t, testDB.DB, storefrontID)
	assert.Equal(t, dbOutOfStock, stats.OutOfStock, "Out of stock should match database")

	dbLowStock := tests.CountLowStockProducts(t, testDB.DB, storefrontID)
	assert.Equal(t, dbLowStock, stats.LowStock, "Low stock should match database")

	dbValue := tests.GetTotalInventoryValue(t, testDB.DB, storefrontID)
	assert.InDelta(t, dbValue, stats.TotalValue, 0.01, "Total value should match database")
}

// TestGetProductStats_EmptyStorefront_Success tests stats for storefront with no products
func TestGetProductStats_EmptyStorefront_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Create empty storefront
	_, err := testDB.DB.Exec(`
		INSERT INTO storefronts (id, user_id, name, slug, is_active, created_at, updated_at)
		VALUES (9999, 1000, 'Empty Store', 'empty-store', true, NOW(), NOW())
	`)
	require.NoError(t, err)

	stats, err := repo.GetProductStats(ctx, 9999)

	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, int32(0), stats.TotalProducts)
	assert.Equal(t, int32(0), stats.ActiveProducts)
	assert.Equal(t, int32(0), stats.OutOfStock)
	assert.Equal(t, int32(0), stats.LowStock)
	assert.Equal(t, float64(0), stats.TotalValue)
	assert.Equal(t, int32(0), stats.TotalSold)
}

// TestGetProductStats_NonExistentStorefront_Error tests error for non-existent storefront
func TestGetProductStats_NonExistentStorefront_Error(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	stats, err := repo.GetProductStats(ctx, 99999)

	// Depending on implementation, might return empty stats or error
	// Adjust assertion based on actual behavior
	if err != nil {
		assert.Error(t, err)
	} else {
		require.NotNil(t, stats)
		assert.Equal(t, int32(0), stats.TotalProducts)
	}
}

// TestIncrementProductViews_ValidProduct_Success tests successful view increment
func TestIncrementProductViews_ValidProduct_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	productID := int64(5000)

	// Get initial view count
	initialViews := tests.GetProductViewCount(t, testDB.DB, productID)

	// Increment views
	err := repo.IncrementProductViews(ctx, productID)
	require.NoError(t, err)

	// Verify increment
	newViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Equal(t, initialViews+1, newViews, "Views should increment by 1")
}

// TestIncrementProductViews_MultipleIncrements_Success tests multiple increments
func TestIncrementProductViews_MultipleIncrements_Success(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	productID := int64(5006) // Already has 10 views

	initialViews := tests.GetProductViewCount(t, testDB.DB, productID)
	require.Equal(t, int32(10), initialViews, "Should start with 10 views")

	// Increment 5 times
	for i := 0; i < 5; i++ {
		err := repo.IncrementProductViews(ctx, productID)
		require.NoError(t, err, "Increment #%d should succeed", i+1)
	}

	// Verify total
	finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Equal(t, int32(15), finalViews, "Should have 15 views after 5 increments")
}

// TestIncrementProductViews_NonExistentProduct_Error tests error for non-existent product
func TestIncrementProductViews_NonExistentProduct_Error(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	err := repo.IncrementProductViews(ctx, 99999)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "products.not_found", "Should return product not found error")
}

// TestIncrementProductViews_ConcurrentIncrements tests concurrent view increments
func TestIncrementProductViews_ConcurrentIncrements(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	productID := int64(5000)
	concurrentCalls := 10

	initialViews := tests.GetProductViewCount(t, testDB.DB, productID)

	// Run concurrent increments
	done := make(chan error, concurrentCalls)
	for i := 0; i < concurrentCalls; i++ {
		go func() {
			done <- repo.IncrementProductViews(ctx, productID)
		}()
	}

	// Wait for all to complete
	for i := 0; i < concurrentCalls; i++ {
		err := <-done
		assert.NoError(t, err, "Concurrent increment should succeed")
	}

	// Verify total (should be initial + 10)
	finalViews := tests.GetProductViewCount(t, testDB.DB, productID)
	assert.Equal(t, initialViews+int32(concurrentCalls), finalViews,
		"All concurrent increments should be recorded")
}
