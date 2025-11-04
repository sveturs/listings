package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestForeignKeyConstraints runs integration tests for FK constraints
// These tests verify CASCADE DELETE and RESTRICT behavior
func TestForeignKeyConstraints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database connection
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	t.Run("CASCADE_DELETE_c2c_images", testCascadeDeleteC2CImages(db))
	t.Run("CASCADE_DELETE_c2c_attributes", testCascadeDeleteC2CAttributes(db))
	t.Run("CASCADE_DELETE_c2c_favorites", testCascadeDeleteC2CFavorites(db))
	t.Run("CASCADE_DELETE_b2c_product_images", testCascadeDeleteB2CImages(db))
	t.Run("CASCADE_DELETE_b2c_product_variants", testCascadeDeleteB2CVariants(db))
	t.Run("RESTRICT_category_with_listings", testRestrictCategoryWithListings(db))
	t.Run("RESTRICT_storefront_with_products", testRestrictStorefrontWithProducts(db))
	t.Run("MULTI_CASCADE_layers", testMultiCascadeLayers(db))
}

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	require.NoError(t, err, "Failed to connect to test database")

	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	return db
}

// testCascadeDeleteC2CImages verifies that deleting a listing cascades to images
func testCascadeDeleteC2CImages(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var userID, categoryID int
		err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		require.NoError(t, err, "No users found in database")

		err = tx.QueryRow("SELECT id FROM c2c_categories LIMIT 1").Scan(&categoryID)
		require.NoError(t, err, "No categories found in database")

		// Create test listing
		listingID := 888888
		_, err = tx.Exec(`
			INSERT INTO c2c_listings (id, user_id, category_id, title, description, price, status)
			VALUES ($1, $2, $3, 'Test Listing CASCADE', 'Test', 100, 'active')
		`, listingID, userID, categoryID)
		require.NoError(t, err, "Failed to create test listing")

		// Create test images
		_, err = tx.Exec(`
			INSERT INTO c2c_images (listing_id, image_url, sort_order)
			VALUES ($1, 'http://test.com/img1.jpg', 1),
			       ($1, 'http://test.com/img2.jpg', 2)
		`, listingID)
		require.NoError(t, err, "Failed to create test images")

		// Verify images exist
		var imageCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM c2c_images WHERE listing_id = $1", listingID).Scan(&imageCount)
		require.NoError(t, err)
		assert.Equal(t, 2, imageCount, "Expected 2 images before deletion")

		// Delete listing (should CASCADE to images)
		_, err = tx.Exec("DELETE FROM c2c_listings WHERE id = $1", listingID)
		require.NoError(t, err, "Failed to delete listing")

		// Verify images were CASCADE deleted
		err = tx.QueryRow("SELECT COUNT(*) FROM c2c_images WHERE listing_id = $1", listingID).Scan(&imageCount)
		require.NoError(t, err)
		assert.Equal(t, 0, imageCount, "Images should be CASCADE deleted with listing")
	}
}

// testCascadeDeleteC2CAttributes verifies cascade delete for attributes
func testCascadeDeleteC2CAttributes(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var userID, categoryID, attributeID int
		err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		require.NoError(t, err)

		err = tx.QueryRow("SELECT id FROM c2c_categories LIMIT 1").Scan(&categoryID)
		require.NoError(t, err)

		err = tx.QueryRow("SELECT id FROM c2c_attributes_meta LIMIT 1").Scan(&attributeID)
		require.NoError(t, err, "No attributes found")

		// Create test listing
		listingID := 888887
		_, err = tx.Exec(`
			INSERT INTO c2c_listings (id, user_id, category_id, title, description, price, status)
			VALUES ($1, $2, $3, 'Test Attrs CASCADE', 'Test', 100, 'active')
		`, listingID, userID, categoryID)
		require.NoError(t, err)

		// Create test attributes
		_, err = tx.Exec(`
			INSERT INTO c2c_attributes (listing_id, attribute_id, value)
			VALUES ($1, $2, 'Test Value 1'),
			       ($1, $2, 'Test Value 2')
		`, listingID, attributeID)
		require.NoError(t, err)

		// Verify attributes exist
		var attrCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM c2c_attributes WHERE listing_id = $1", listingID).Scan(&attrCount)
		require.NoError(t, err)
		assert.Equal(t, 2, attrCount, "Expected 2 attributes before deletion")

		// Delete listing
		_, err = tx.Exec("DELETE FROM c2c_listings WHERE id = $1", listingID)
		require.NoError(t, err)

		// Verify CASCADE delete
		err = tx.QueryRow("SELECT COUNT(*) FROM c2c_attributes WHERE listing_id = $1", listingID).Scan(&attrCount)
		require.NoError(t, err)
		assert.Equal(t, 0, attrCount, "Attributes should be CASCADE deleted")
	}
}

// testCascadeDeleteC2CFavorites verifies cascade delete for favorites
func testCascadeDeleteC2CFavorites(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var userID, categoryID int
		err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		require.NoError(t, err)

		err = tx.QueryRow("SELECT id FROM c2c_categories LIMIT 1").Scan(&categoryID)
		require.NoError(t, err)

		// Create test listing
		listingID := 888886
		_, err = tx.Exec(`
			INSERT INTO c2c_listings (id, user_id, category_id, title, description, price, status)
			VALUES ($1, $2, $3, 'Test Favorites CASCADE', 'Test', 100, 'active')
		`, listingID, userID, categoryID)
		require.NoError(t, err)

		// Create favorite
		_, err = tx.Exec(`
			INSERT INTO c2c_favorites (user_id, listing_id)
			VALUES ($1, $2)
		`, userID, listingID)
		require.NoError(t, err)

		// Verify favorite exists
		var favCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM c2c_favorites WHERE listing_id = $1", listingID).Scan(&favCount)
		require.NoError(t, err)
		assert.Equal(t, 1, favCount, "Expected 1 favorite before deletion")

		// Delete listing
		_, err = tx.Exec("DELETE FROM c2c_listings WHERE id = $1", listingID)
		require.NoError(t, err)

		// Verify CASCADE delete
		err = tx.QueryRow("SELECT COUNT(*) FROM c2c_favorites WHERE listing_id = $1", listingID).Scan(&favCount)
		require.NoError(t, err)
		assert.Equal(t, 0, favCount, "Favorites should be CASCADE deleted")
	}
}

// testCascadeDeleteB2CImages verifies B2C product images cascade delete
func testCascadeDeleteB2CImages(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var storefrontID, categoryID int
		err = tx.QueryRow("SELECT id FROM storefronts LIMIT 1").Scan(&storefrontID)
		if err != nil {
			t.Skip("No storefronts found, skipping B2C test")
		}

		err = tx.QueryRow("SELECT id FROM b2c_categories LIMIT 1").Scan(&categoryID)
		if err != nil {
			t.Skip("No B2C categories found, skipping test")
		}

		// Create test product
		productID := 888885
		_, err = tx.Exec(`
			INSERT INTO b2c_products (id, storefront_id, category_id, name, description, price, status)
			VALUES ($1, $2, $3, 'Test B2C Product', 'Test', 100, 'active')
		`, productID, storefrontID, categoryID)
		require.NoError(t, err)

		// Create product images
		_, err = tx.Exec(`
			INSERT INTO b2c_product_images (product_id, image_url, sort_order)
			VALUES ($1, 'http://test.com/product1.jpg', 1),
			       ($1, 'http://test.com/product2.jpg', 2)
		`, productID)
		require.NoError(t, err)

		// Verify images exist
		var imageCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM b2c_product_images WHERE product_id = $1", productID).Scan(&imageCount)
		require.NoError(t, err)
		assert.Equal(t, 2, imageCount, "Expected 2 product images")

		// Delete product
		_, err = tx.Exec("DELETE FROM b2c_products WHERE id = $1", productID)
		require.NoError(t, err)

		// Verify CASCADE delete
		err = tx.QueryRow("SELECT COUNT(*) FROM b2c_product_images WHERE product_id = $1", productID).Scan(&imageCount)
		require.NoError(t, err)
		assert.Equal(t, 0, imageCount, "Product images should be CASCADE deleted")
	}
}

// testCascadeDeleteB2CVariants verifies product variants cascade delete
func testCascadeDeleteB2CVariants(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var storefrontID, categoryID int
		err = tx.QueryRow("SELECT id FROM storefronts LIMIT 1").Scan(&storefrontID)
		if err != nil {
			t.Skip("No storefronts found")
		}

		err = tx.QueryRow("SELECT id FROM b2c_categories LIMIT 1").Scan(&categoryID)
		if err != nil {
			t.Skip("No B2C categories found")
		}

		// Create test product
		productID := 888884
		_, err = tx.Exec(`
			INSERT INTO b2c_products (id, storefront_id, category_id, name, description, price, status)
			VALUES ($1, $2, $3, 'Test Product Variants', 'Test', 100, 'active')
		`, productID, storefrontID, categoryID)
		require.NoError(t, err)

		// Create product variants
		_, err = tx.Exec(`
			INSERT INTO b2c_product_variants (product_id, sku, name, price, stock)
			VALUES ($1, 'TEST-SKU-001', 'Variant 1', 150, 10),
			       ($1, 'TEST-SKU-002', 'Variant 2', 200, 5)
		`, productID)
		require.NoError(t, err)

		// Verify variants exist
		var variantCount int
		err = tx.QueryRow("SELECT COUNT(*) FROM b2c_product_variants WHERE product_id = $1", productID).Scan(&variantCount)
		require.NoError(t, err)
		assert.Equal(t, 2, variantCount, "Expected 2 variants")

		// Delete product
		_, err = tx.Exec("DELETE FROM b2c_products WHERE id = $1", productID)
		require.NoError(t, err)

		// Verify CASCADE delete
		err = tx.QueryRow("SELECT COUNT(*) FROM b2c_product_variants WHERE product_id = $1", productID).Scan(&variantCount)
		require.NoError(t, err)
		assert.Equal(t, 0, variantCount, "Variants should be CASCADE deleted")
	}
}

// testRestrictCategoryWithListings verifies RESTRICT constraint on categories
func testRestrictCategoryWithListings(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var userID, categoryID int
		err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		require.NoError(t, err)

		err = tx.QueryRow("SELECT id FROM c2c_categories LIMIT 1").Scan(&categoryID)
		require.NoError(t, err)

		// Create test listing
		listingID := 888883
		_, err = tx.Exec(`
			INSERT INTO c2c_listings (id, user_id, category_id, title, description, price, status)
			VALUES ($1, $2, $3, 'Test RESTRICT', 'Test', 100, 'active')
		`, listingID, userID, categoryID)
		require.NoError(t, err)

		// Try to delete category (should FAIL with FK violation)
		_, err = tx.Exec("DELETE FROM c2c_categories WHERE id = $1", categoryID)
		assert.Error(t, err, "Should not be able to delete category with listings")
		assert.Contains(t, err.Error(), "foreign key", "Should be FK constraint violation")
	}
}

// testRestrictStorefrontWithProducts verifies RESTRICT on storefront deletion
func testRestrictStorefrontWithProducts(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var userID int
		err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		require.NoError(t, err)

		var categoryID int
		err = tx.QueryRow("SELECT id FROM b2c_categories LIMIT 1").Scan(&categoryID)
		if err != nil {
			t.Skip("No B2C categories found")
		}

		// Create test storefront
		var storefrontID int
		err = tx.QueryRow(`
			INSERT INTO storefronts (user_id, name, description, status)
			VALUES ($1, 'Test Storefront RESTRICT', 'Test', 'active')
			RETURNING id
		`, userID).Scan(&storefrontID)
		require.NoError(t, err)

		// Create product for storefront
		productID := 888882
		_, err = tx.Exec(`
			INSERT INTO b2c_products (id, storefront_id, category_id, name, description, price, status)
			VALUES ($1, $2, $3, 'Test Product', 'Test', 100, 'active')
		`, productID, storefrontID, categoryID)
		require.NoError(t, err)

		// Try to delete storefront (should FAIL)
		_, err = tx.Exec("DELETE FROM storefronts WHERE id = $1", storefrontID)
		assert.Error(t, err, "Should not be able to delete storefront with products")
		assert.Contains(t, err.Error(), "foreign key", "Should be FK constraint violation")
	}
}

// testMultiCascadeLayers verifies multiple CASCADE layers work together
func testMultiCascadeLayers(db *sql.DB) func(*testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)
		defer func() { _ = tx.Rollback() }()

		// Get test data
		var userID, categoryID, attributeID int
		err = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
		require.NoError(t, err)

		err = tx.QueryRow("SELECT id FROM c2c_categories LIMIT 1").Scan(&categoryID)
		require.NoError(t, err)

		err = tx.QueryRow("SELECT id FROM c2c_attributes_meta LIMIT 1").Scan(&attributeID)
		require.NoError(t, err)

		// Create test listing
		listingID := 888881
		_, err = tx.Exec(`
			INSERT INTO c2c_listings (id, user_id, category_id, title, description, price, status)
			VALUES ($1, $2, $3, 'Test Multi CASCADE', 'Test', 100, 'active')
		`, listingID, userID, categoryID)
		require.NoError(t, err)

		// Create multiple child records
		_, err = tx.Exec(`
			INSERT INTO c2c_images (listing_id, image_url, sort_order)
			VALUES ($1, 'http://test.com/multi1.jpg', 1),
			       ($1, 'http://test.com/multi2.jpg', 2),
			       ($1, 'http://test.com/multi3.jpg', 3)
		`, listingID)
		require.NoError(t, err)

		_, err = tx.Exec(`
			INSERT INTO c2c_attributes (listing_id, attribute_id, value)
			VALUES ($1, $2, 'Multi Value 1'),
			       ($1, $2, 'Multi Value 2')
		`, listingID, attributeID)
		require.NoError(t, err)

		_, err = tx.Exec(`
			INSERT INTO c2c_favorites (user_id, listing_id)
			VALUES ($1, $2)
		`, userID, listingID)
		require.NoError(t, err)

		// Verify all created
		var imageCount, attrCount, favCount int
		_ = tx.QueryRow("SELECT COUNT(*) FROM c2c_images WHERE listing_id = $1", listingID).Scan(&imageCount)
		_ = tx.QueryRow("SELECT COUNT(*) FROM c2c_attributes WHERE listing_id = $1", listingID).Scan(&attrCount)
		_ = tx.QueryRow("SELECT COUNT(*) FROM c2c_favorites WHERE listing_id = $1", listingID).Scan(&favCount)

		assert.Equal(t, 3, imageCount, "Expected 3 images")
		assert.Equal(t, 2, attrCount, "Expected 2 attributes")
		assert.Equal(t, 1, favCount, "Expected 1 favorite")

		// Delete listing (should CASCADE to all children)
		_, err = tx.Exec("DELETE FROM c2c_listings WHERE id = $1", listingID)
		require.NoError(t, err)

		// Verify ALL children CASCADE deleted
		_ = tx.QueryRow("SELECT COUNT(*) FROM c2c_images WHERE listing_id = $1", listingID).Scan(&imageCount)
		_ = tx.QueryRow("SELECT COUNT(*) FROM c2c_attributes WHERE listing_id = $1", listingID).Scan(&attrCount)
		_ = tx.QueryRow("SELECT COUNT(*) FROM c2c_favorites WHERE listing_id = $1", listingID).Scan(&favCount)

		assert.Equal(t, 0, imageCount, "All images should be CASCADE deleted")
		assert.Equal(t, 0, attrCount, "All attributes should be CASCADE deleted")
		assert.Equal(t, 0, favCount, "All favorites should be CASCADE deleted")
	}
}

// TestFKConstraintMetadata verifies FK constraints exist in database schema
func TestFKConstraintMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	query := `
		SELECT
			COUNT(*) as total_fks,
			SUM(CASE WHEN rc.delete_rule = 'CASCADE' THEN 1 ELSE 0 END) as cascade_count,
			SUM(CASE WHEN rc.delete_rule IN ('RESTRICT', 'NO ACTION') THEN 1 ELSE 0 END) as restrict_count
		FROM information_schema.table_constraints tc
		JOIN information_schema.referential_constraints rc
			ON rc.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		AND tc.table_schema = 'public'
		AND tc.table_name IN (
			'c2c_listings', 'c2c_images', 'c2c_attributes', 'c2c_favorites',
			'b2c_products', 'b2c_product_images', 'b2c_product_variants',
			'storefronts', 'c2c_categories', 'b2c_categories'
		)
	`

	var totalFKs, cascadeCount, restrictCount int
	err := db.QueryRow(query).Scan(&totalFKs, &cascadeCount, &restrictCount)
	require.NoError(t, err)

	t.Logf("FK Constraints Found:")
	t.Logf("  Total: %d", totalFKs)
	t.Logf("  CASCADE: %d", cascadeCount)
	t.Logf("  RESTRICT/NO ACTION: %d", restrictCount)

	// Assert we have FKs defined
	assert.Greater(t, totalFKs, 0, "Should have FK constraints defined")
	assert.Greater(t, cascadeCount, 0, "Should have CASCADE constraints")
}

// BenchmarkCascadeDelete benchmarks CASCADE DELETE performance
func BenchmarkCascadeDelete(b *testing.B) {
	db := setupTestDB(&testing.T{})
	defer func() { _ = db.Close() }()

	ctx := context.Background()

	b.Run("Single_Listing_With_Children", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tx, _ := db.BeginTx(ctx, nil)

			// Create listing with children
			var userID, categoryID int
			_ = tx.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
			_ = tx.QueryRow("SELECT id FROM c2c_categories LIMIT 1").Scan(&categoryID)

			listingID := 777000 + i
			_, _ = tx.Exec(`
				INSERT INTO c2c_listings (id, user_id, category_id, title, description, price, status)
				VALUES ($1, $2, $3, 'Benchmark', 'Test', 100, 'active')
			`, listingID, userID, categoryID)

			_, _ = tx.Exec(`INSERT INTO c2c_images (listing_id, image_url, sort_order) VALUES ($1, 'test.jpg', 1)`, listingID)

			// Measure CASCADE DELETE
			b.StartTimer()
			_, _ = tx.Exec("DELETE FROM c2c_listings WHERE id = $1", listingID)
			b.StopTimer()

			_ = tx.Rollback()
		}
	})
}

// Helper function to print FK constraint details (for debugging)
func PrintFKConstraints(t *testing.T, db *sql.DB) {
	query := `
		SELECT
			tc.table_name,
			tc.constraint_name,
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			rc.delete_rule
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage ccu
			ON ccu.constraint_name = tc.constraint_name
		JOIN information_schema.referential_constraints rc
			ON rc.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		AND tc.table_schema = 'public'
		ORDER BY tc.table_name, tc.constraint_name
	`

	rows, err := db.Query(query)
	require.NoError(t, err)
	defer func() { _ = rows.Close() }()

	t.Log("FK Constraints:")
	for rows.Next() {
		var table, constraint, column, foreignTable, foreignColumn, deleteRule string
		_ = rows.Scan(&table, &constraint, &column, &foreignTable, &foreignColumn, &deleteRule)
		t.Logf("  %s.%s -> %s.%s (ON DELETE %s)",
			table, column, foreignTable, foreignColumn, deleteRule)
	}
	_ = rows.Err()
}
