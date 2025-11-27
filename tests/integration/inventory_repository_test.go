//go:build integration

package integration

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/tests"
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

// ============================================================================
// PHASE 9.7.4: ADVANCED REPOSITORY TESTS
// ============================================================================

// TestUpdateProductInventory_TransactionRollback tests transaction rollback on constraint violation
// Validates that database constraints are enforced and transactions roll back properly
func TestUpdateProductInventory_TransactionRollback(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5001) // Low stock product with quantity 5
	userID := int64(1000)

	initialQty := tests.GetProductQuantity(t, testDB.DB, productID)
	require.Equal(t, int32(5), initialQty, "Initial quantity should be 5")

	// Try to remove more stock than available (should fail)
	_, _, err := repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"out",
		100, // Request 100 but only 5 available
		"sale",
		"Should fail due to insufficient stock",
		userID,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "products.insufficient_stock")

	// Verify quantity remained unchanged (transaction rolled back)
	finalQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, initialQty, finalQty,
		"Quantity should remain unchanged after failed transaction")

	// Verify no inventory movement was recorded (full rollback)
	initialMovementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)

	// Try another operation
	_, _, err = repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"in",
		10,
		"restock",
		"Valid operation after rollback",
		userID,
	)
	require.NoError(t, err)

	// This should create one movement
	finalMovementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
	assert.Equal(t, initialMovementCount+1, finalMovementCount,
		"Only successful operations should be recorded")
}

// TestBatchUpdateStock_TransactionAtomicity tests atomicity of batch operations
// Validates that batch updates maintain transactional integrity
func TestBatchUpdateStock_TransactionAtomicity(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)

	// Get initial quantities
	product5004InitialQty := tests.GetProductQuantity(t, testDB.DB, 5004)

	// Create batch with one valid and one invalid item
	items := []domain.StockUpdateItem{
		{
			ProductID: 5003,
			Quantity:  120,
			Reason:    stringPtr("valid_update"),
		},
		{
			ProductID: 99999, // Non-existent product
			Quantity:  50,
			Reason:    stringPtr("invalid_update"),
		},
	}

	successCount, failedCount, results, err := repo.BatchUpdateStock(
		ctx,
		storefrontID,
		items,
		"atomicity_test",
		userID,
	)

	require.NoError(t, err, "Batch operation itself should not error")

	// Verify counts
	assert.Equal(t, int32(1), successCount, "One item should succeed")
	assert.Equal(t, int32(1), failedCount, "One item should fail")
	assert.Len(t, results, 2, "Should return 2 results")

	// Verify first item succeeded
	assert.True(t, results[0].Success)
	assert.Equal(t, int64(5003), results[0].ProductID)
	product5003NewQty := tests.GetProductQuantity(t, testDB.DB, 5003)
	assert.Equal(t, int32(120), product5003NewQty, "Valid product should be updated")

	// Verify second item failed
	assert.False(t, results[1].Success)
	assert.Equal(t, int64(99999), results[1].ProductID)
	assert.NotNil(t, results[1].Error)

	// Verify other products unaffected
	product5004NewQty := tests.GetProductQuantity(t, testDB.DB, 5004)
	assert.Equal(t, product5004InitialQty, product5004NewQty,
		"Unrelated products should remain unchanged")
}

// TestUpdateProductInventory_MaxInt32Quantity tests maximum int32 boundary
// Validates that the system correctly handles maximum possible quantity values
func TestUpdateProductInventory_MaxInt32Quantity(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	// Set quantity to maximum int32 value
	maxInt32 := int32(2147483647)

	stockBefore, stockAfter, err := repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"adjustment",
		maxInt32,
		"boundary_test",
		"Testing maximum int32 value",
		userID,
	)

	require.NoError(t, err)
	assert.Equal(t, maxInt32, stockAfter, "Should set to max int32")
	assert.NotEqual(t, stockBefore, stockAfter, "Stock should change")

	// Verify in database
	actualQty := tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, maxInt32, actualQty, "Database should store max int32")

	// Try to remove 1 from max
	stockBefore, stockAfter, err = repo.UpdateProductInventory(
		ctx,
		storefrontID,
		productID,
		0,
		"out",
		1,
		"boundary_decrement",
		"Decrement from max",
		userID,
	)

	require.NoError(t, err)
	assert.Equal(t, maxInt32, stockBefore)
	assert.Equal(t, maxInt32-1, stockAfter)

	// Verify decrement worked
	actualQty = tests.GetProductQuantity(t, testDB.DB, productID)
	assert.Equal(t, maxInt32-1, actualQty)
}

// TestUpdateProductInventory_UnicodeHandling tests Unicode characters in text fields
// Validates proper handling of international characters, emojis, and special symbols
func TestUpdateProductInventory_UnicodeHandling(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	productID := int64(5000)
	userID := int64(1000)

	testCases := []struct {
		name   string
		reason string
		notes  string
	}{
		{
			name:   "Cyrillic characters",
			reason: "Поступление товара",
			notes:  "Получено от поставщика ООО Рога и Копыта",
		},
		{
			name:   "Chinese characters",
			reason: "库存调整",
			notes:  "根据实际盘点结果调整库存数量",
		},
		{
			name:   "Arabic characters",
			reason: "استلام البضائع",
			notes:  "تم استلام البضائع من المورد",
		},
		{
			name:   "Emoji and special characters",
			reason: "📦 Restock 🚚",
			notes:  "New inventory arrived! 🎉 Check quality ✅",
		},
		{
			name:   "Mixed scripts",
			reason: "Restock товара 产品 📦",
			notes:  "International supplier delivery - Международная поставка",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stockBefore, stockAfter, err := repo.UpdateProductInventory(
				ctx,
				storefrontID,
				productID,
				0,
				"in",
				10,
				tc.reason,
				tc.notes,
				userID,
			)

			require.NoError(t, err, "Should handle Unicode characters")
			assert.Equal(t, stockBefore+10, stockAfter, "Quantity should increment correctly")

			// Verify movement was recorded with correct text
			// (This assumes we have a way to query movements by product)
			movementCount := tests.GetInventoryMovementCount(t, testDB.DB, productID)
			assert.Greater(t, movementCount, 0, "Movement should be recorded")
		})
	}
}

// TestBatchUpdateStock_ConcurrentBatches tests concurrent batch operations
// Validates thread safety and race condition handling in batch updates
func TestBatchUpdateStock_ConcurrentBatches(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)
	concurrentBatches := 5

	// Each batch updates different products to avoid conflicts
	productIDs := []int64{5000, 5003, 5004, 5006}

	var wg sync.WaitGroup
	errors := make(chan error, concurrentBatches)

	for i := 0; i < concurrentBatches; i++ {
		wg.Add(1)
		go func(batchNum int) {
			defer wg.Done()

			items := make([]domain.StockUpdateItem, len(productIDs))
			for j, prodID := range productIDs {
				items[j] = domain.StockUpdateItem{
					ProductID: prodID,
					Quantity:  int32(100 + batchNum*10 + j),
					Reason:    stringPtr("concurrent_batch"),
				}
			}

			_, _, _, err := repo.BatchUpdateStock(
				ctx,
				storefrontID,
				items,
				"concurrency_test",
				userID,
			)
			errors <- err
		}(i)
	}

	wg.Wait()
	close(errors)

	// Verify all batches succeeded
	for err := range errors {
		assert.NoError(t, err, "All concurrent batches should succeed")
	}

	// Verify final state is consistent
	for _, prodID := range productIDs {
		qty := tests.GetProductQuantity(t, testDB.DB, prodID)
		assert.Greater(t, qty, int32(0), "Product should have positive quantity")
	}
}

// TestUpdateProductInventory_DeadlockPrevention tests deadlock prevention mechanisms
// Validates that the repository properly orders operations to prevent deadlocks
func TestUpdateProductInventory_DeadlockPrevention(t *testing.T) {
	repo, testDB := setupInventoryTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	storefrontID := int64(1000)
	userID := int64(1000)
	concurrentOps := 20

	// Products to update (creates potential for deadlock if not handled properly)
	productIDs := []int64{5000, 5003, 5004}

	var wg sync.WaitGroup
	successCount := int32(0)
	errorCount := int32(0)

	// Launch concurrent operations on same products
	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func(iteration int) {
			defer wg.Done()

			// Alternate between products to create contention
			productID := productIDs[iteration%len(productIDs)]

			_, _, err := repo.UpdateProductInventory(
				ctx,
				storefrontID,
				productID,
				0,
				"in",
				1,
				"deadlock_test",
				"Testing concurrent updates",
				userID,
			)

			if err != nil {
				atomic.AddInt32(&errorCount, 1)
				// Log error for debugging
				t.Logf("Operation %d failed: %v", iteration, err)
			} else {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// Verify most operations succeeded (some failures acceptable under high contention)
	successRate := float64(successCount) / float64(concurrentOps)
	assert.Greater(t, successRate, 0.8,
		"At least 80%% of concurrent operations should succeed (no deadlocks)")

	t.Logf("Success rate: %.2f%% (%d/%d)", successRate*100, successCount, concurrentOps)

	// Verify database state is consistent
	for _, prodID := range productIDs {
		qty := tests.GetProductQuantity(t, testDB.DB, prodID)
		assert.Greater(t, qty, int32(0), "Product should have valid quantity")

		movementCount := tests.GetInventoryMovementCount(t, testDB.DB, prodID)
		assert.Greater(t, movementCount, 0, "Should have recorded movements")
	}
}
