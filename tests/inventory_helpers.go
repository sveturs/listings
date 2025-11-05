package tests

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// LoadInventoryFixtures loads inventory test fixtures
func LoadInventoryFixtures(t *testing.T, db *sql.DB) {
	t.Helper()
	// Load B2C inventory fixtures (uses b2c_products tables)
	LoadTestFixtures(t, db, "../fixtures/b2c_inventory_fixtures.sql")
}

// LoadRollbackStockFixtures loads rollback stock test fixtures
func LoadRollbackStockFixtures(t *testing.T, db *sql.DB) {
	t.Helper()
	// Load rollback-specific fixtures (includes decrement history)
	LoadTestFixtures(t, db, "../fixtures/rollback_stock_fixtures.sql")
}

// CleanupInventoryTestData removes inventory test data
func CleanupInventoryTestData(t *testing.T, db *sql.DB) {
	t.Helper()

	tables := []string{
		"b2c_inventory_movements",
		"b2c_product_variants",
		"b2c_products",
		"storefronts",
	}

	for _, table := range tables {
		// Delete test data with IDs >= 5000 (our test range)
		// Also include rollback test data (IDs >= 8000)
		var query string
		switch table {
		case "b2c_inventory_movements":
			query = "DELETE FROM b2c_inventory_movements WHERE id >= 7000"
		case "b2c_product_variants":
			query = "DELETE FROM b2c_product_variants WHERE id >= 6000"
		case "b2c_products":
			query = "DELETE FROM b2c_products WHERE id >= 5000"
		case "storefronts":
			query = "DELETE FROM storefronts WHERE id >= 1000"
		default:
			query = "DELETE FROM " + table + " WHERE id >= 5000"
		}

		_, err := db.Exec(query)
		require.NoError(t, err, "Could not cleanup table: %s", table)
	}

	// Clean test users (if users table exists)
	_, err := db.Exec("DELETE FROM users WHERE id >= 1000 AND id < 2000")
	if err != nil {
		// Users table might not exist in listings microservice, ignore error
		t.Logf("Could not cleanup test users (table might not exist): %v", err)
	}

	// Clean test categories (if categories table exists)
	_, err = db.Exec("DELETE FROM categories WHERE id >= 2000 AND id < 3000")
	if err != nil {
		// Categories table might not exist in listings microservice, ignore error
		t.Logf("Could not cleanup test categories (table might not exist): %v", err)
	}
}

// GetProductQuantity returns current product quantity from database
func GetProductQuantity(t *testing.T, db *sql.DB, productID int64) int32 {
	t.Helper()

	var quantity int32
	err := db.QueryRow("SELECT stock_quantity FROM b2c_products WHERE id = $1", productID).Scan(&quantity)
	require.NoError(t, err, "Could not get product quantity")

	return quantity
}

// GetVariantQuantity returns current variant quantity from database
func GetVariantQuantity(t *testing.T, db *sql.DB, variantID int64) int32 {
	t.Helper()

	var quantity int32
	err := db.QueryRow("SELECT stock_quantity FROM b2c_product_variants WHERE id = $1", variantID).Scan(&quantity)
	require.NoError(t, err, "Could not get variant quantity")

	return quantity
}

// GetInventoryMovementCount returns count of inventory movements for a product
func GetInventoryMovementCount(t *testing.T, db *sql.DB, productID int64) int {
	t.Helper()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM b2c_inventory_movements WHERE storefront_product_id = $1", productID).Scan(&count)
	require.NoError(t, err, "Could not count inventory movements")

	return count
}

// GetProductViewCount returns view count for a product
func GetProductViewCount(t *testing.T, db *sql.DB, productID int64) int32 {
	t.Helper()

	var views int32
	err := db.QueryRow("SELECT COALESCE(view_count, 0) FROM b2c_products WHERE id = $1", productID).Scan(&views)
	if err == sql.ErrNoRows {
		return 0
	}
	require.NoError(t, err, "Could not get product view count")

	return views
}

// CountProductsByStorefront returns count of products in storefront
func CountProductsByStorefront(t *testing.T, db *sql.DB, storefrontID int64) int32 {
	t.Helper()

	var count int32
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM b2c_products
		WHERE storefront_id = $1 AND deleted_at IS NULL
	`, storefrontID).Scan(&count)
	require.NoError(t, err, "Could not count products")

	return count
}

// CountActiveProductsByStorefront returns count of active products in storefront
func CountActiveProductsByStorefront(t *testing.T, db *sql.DB, storefrontID int64) int32 {
	t.Helper()

	var count int32
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM b2c_products
		WHERE storefront_id = $1 AND is_active = true AND deleted_at IS NULL
	`, storefrontID).Scan(&count)
	require.NoError(t, err, "Could not count active products")

	return count
}

// CountOutOfStockProducts returns count of out-of-stock products in storefront
func CountOutOfStockProducts(t *testing.T, db *sql.DB, storefrontID int64) int32 {
	t.Helper()

	var count int32
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM b2c_products
		WHERE storefront_id = $1 AND stock_status = 'out_of_stock' AND deleted_at IS NULL
	`, storefrontID).Scan(&count)
	require.NoError(t, err, "Could not count out-of-stock products")

	return count
}

// CountLowStockProducts returns count of low-stock products (stock_status = 'low_stock') in storefront
func CountLowStockProducts(t *testing.T, db *sql.DB, storefrontID int64) int32 {
	t.Helper()

	var count int32
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM b2c_products
		WHERE storefront_id = $1 AND stock_status = 'low_stock' AND deleted_at IS NULL
	`, storefrontID).Scan(&count)
	require.NoError(t, err, "Could not count low-stock products")

	return count
}

// GetTotalInventoryValue returns total value of inventory in storefront
func GetTotalInventoryValue(t *testing.T, db *sql.DB, storefrontID int64) float64 {
	t.Helper()

	var total float64
	err := db.QueryRow(`
		SELECT COALESCE(SUM(price * stock_quantity), 0)
		FROM b2c_products
		WHERE storefront_id = $1 AND deleted_at IS NULL
	`, storefrontID).Scan(&total)
	require.NoError(t, err, "Could not calculate total inventory value")

	return total
}

// ProductExists checks if product exists in database
func ProductExists(t *testing.T, db *sql.DB, productID int64) bool {
	t.Helper()

	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM b2c_products WHERE id = $1 AND deleted_at IS NULL)
	`, productID).Scan(&exists)
	require.NoError(t, err, "Could not check product existence")

	return exists
}

// VariantExists checks if variant exists in database
func VariantExists(t *testing.T, db *sql.DB, variantID int64) bool {
	t.Helper()

	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM b2c_product_variants WHERE id = $1)
	`, variantID).Scan(&exists)
	require.NoError(t, err, "Could not check variant existence")

	return exists
}

// GetRollbackMovementCount returns count of rollback ('in') movements for a product/order
func GetRollbackMovementCount(t *testing.T, db *sql.DB, productID int64, notes string) int {
	t.Helper()

	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM b2c_inventory_movements
		WHERE storefront_product_id = $1 AND type = 'in' AND notes LIKE $2
	`, productID, "%"+notes+"%").Scan(&count)
	require.NoError(t, err, "Could not count rollback movements")

	return count
}

// GetDecrementMovementCount returns count of decrement ('out') movements for a product
func GetDecrementMovementCount(t *testing.T, db *sql.DB, productID int64) int {
	t.Helper()

	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM b2c_inventory_movements
		WHERE storefront_product_id = $1 AND type = 'out'
	`, productID).Scan(&count)
	require.NoError(t, err, "Could not count decrement movements")

	return count
}

// GetLatestStockMovement returns the most recent inventory movement for a product
func GetLatestStockMovement(t *testing.T, db *sql.DB, productID int64) (movementType string, quantity int32, notes string) {
	t.Helper()

	err := db.QueryRow(`
		SELECT type, quantity, COALESCE(notes, '')
		FROM b2c_inventory_movements
		WHERE storefront_product_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, productID).Scan(&movementType, &quantity, &notes)

	if err == sql.ErrNoRows {
		return "", 0, ""
	}
	require.NoError(t, err, "Could not get latest stock movement")

	return movementType, quantity, notes
}
