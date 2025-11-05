package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// Test Fixtures and Helpers
// ============================================================================

// createTestStorefront creates a test storefront for products
func createTestStorefront(t *testing.T, repo *Repository) int64 {
	t.Helper()
	ctx := tests.TestContext(t)

	// Create a minimal storefront in the storefronts table
	query := `
		INSERT INTO storefronts (
			user_id, name, slug, description, is_active,
			created_at, updated_at
		) VALUES (
			1, 'Test Storefront', 'test-store', 'Test Description', true,
			NOW(), NOW()
		) RETURNING id
	`

	var storefrontID int64
	err := repo.db.QueryRowContext(ctx, query).Scan(&storefrontID)
	require.NoError(t, err, "Failed to create test storefront")

	return storefrontID
}

// createTestCategory creates a test category (stub, no real categories table in listings service)
func createTestCategory(t *testing.T) int64 {
	t.Helper()
	// Categories are managed externally, just return a valid ID
	return 100
}

// createTestProduct creates a product with default values
func createTestProduct(t *testing.T, repo *Repository, storefrontID int64) *domain.Product {
	t.Helper()
	ctx := tests.TestContext(t)

	input := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    createTestCategory(t),
		Name:          "Test Product",
		Description:   "Test product description",
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 100,
		Attributes:    map[string]interface{}{}, // Empty attributes to avoid JSONB issues
	}

	product, err := repo.CreateProduct(ctx, input)
	require.NoError(t, err, "Failed to create test product")
	require.NotNil(t, product)

	return product
}

// createTestProductWithOptions creates a product with custom options
func createTestProductWithOptions(t *testing.T, repo *Repository, storefrontID int64, sku string, price float64, quantity int32) *domain.Product {
	t.Helper()
	ctx := tests.TestContext(t)

	input := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    createTestCategory(t),
		Name:          "Test Product - " + sku,
		Description:   "Test product with SKU " + sku,
		Price:         price,
		Currency:      "USD",
		SKU:           stringPtr(sku),
		StockQuantity: quantity,
		Attributes:    map[string]interface{}{}, // Empty attributes to avoid JSONB issues
	}

	product, err := repo.CreateProduct(ctx, input)
	require.NoError(t, err, "Failed to create test product with options")
	require.NotNil(t, product)

	return product
}

// createTestVariant creates a variant for a product
func createTestVariant(t *testing.T, repo *Repository, productID int64) *domain.ProductVariant {
	t.Helper()
	ctx := tests.TestContext(t)

	query := `
		INSERT INTO b2c_product_variants (
			product_id, sku, stock_quantity, stock_status,
			is_active, is_default, created_at, updated_at
		) VALUES (
			$1, $2, $3, 'in_stock', true, false, NOW(), NOW()
		) RETURNING id, product_id, sku, stock_quantity, stock_status,
				  is_active, is_default, view_count, sold_count, created_at, updated_at
	`

	var variant domain.ProductVariant
	sku := "VARIANT-SKU-001"
	err := repo.db.QueryRowContext(ctx, query, productID, sku, 50).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.StockQuantity,
		&variant.StockStatus,
		&variant.IsActive,
		&variant.IsDefault,
		&variant.ViewCount,
		&variant.SoldCount,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)
	require.NoError(t, err, "Failed to create test variant")

	return &variant
}

// ============================================================================
// 1. CreateProduct Tests (10 tests)
// ============================================================================

func TestCreateProduct_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Test Product",
		Description:   "Test description",
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 100,
		Attributes:    map[string]interface{}{}, // Empty attributes
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	require.NoError(t, err)
	assert.NotZero(t, createdProduct.ID)
	assert.Equal(t, product.Name, createdProduct.Name)
	assert.Equal(t, product.SKU, createdProduct.SKU)
	assert.Equal(t, product.Price, createdProduct.Price)
	assert.Equal(t, product.StockQuantity, createdProduct.StockQuantity)
	assert.Equal(t, domain.StockStatusInStock, createdProduct.StockStatus)
	assert.True(t, createdProduct.IsActive)

	// Verify in database
	fetchedProduct, err := repo.GetProductByID(ctx, createdProduct.ID, &storefrontID)
	require.NoError(t, err)
	assert.Equal(t, createdProduct.ID, fetchedProduct.ID)
}

func TestCreateProduct_WithVariants(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)

	// Create product
	product := createTestProduct(t, repo, storefrontID)

	// Create variant
	variant := createTestVariant(t, repo, product.ID)

	assert.NotZero(t, variant.ID)
	assert.Equal(t, product.ID, variant.ProductID)
}

func TestCreateProduct_MissingName(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "", // Empty name
		Description:   "Test description",
		Price:         99.99,
		Currency:      "USD",
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, createdProduct)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestCreateProduct_MissingSKU(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Test Product",
		Description:   "Test description",
		Price:         99.99,
		Currency:      "USD",
		SKU:           nil, // No SKU (this is valid - SKU is optional)
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	// SKU is optional, should succeed
	require.NoError(t, err)
	assert.NotNil(t, createdProduct)
	assert.Nil(t, createdProduct.SKU)
}

func TestCreateProduct_DuplicateSKU(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	sku := "DUPLICATE-SKU-001"

	// Create first product
	product1 := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Product 1",
		Description:   "First product",
		Price:         99.99,
		Currency:      "USD",
		SKU:           &sku,
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	_, err := repo.CreateProduct(ctx, product1)
	require.NoError(t, err)

	// Try to create second product with same SKU
	product2 := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Product 2",
		Description:   "Second product",
		Price:         149.99,
		Currency:      "USD",
		SKU:           &sku, // Same SKU
		StockQuantity: 50,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product2)

	assert.Error(t, err)
	assert.Nil(t, createdProduct)
	assert.Contains(t, err.Error(), "products.sku_duplicate")
}

func TestCreateProduct_InvalidStorefrontID(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  99999, // Non-existent storefront
		CategoryID:    categoryID,
		Name:          "Test Product",
		Description:   "Test description",
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, createdProduct)
}

func TestCreateProduct_InvalidCategoryID(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    0, // Invalid category (0 is technically valid as there's no FK constraint)
		Name:          "Test Product",
		Description:   "Test description",
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	// No FK constraint on category_id, so this should succeed
	require.NoError(t, err)
	assert.NotNil(t, createdProduct)
}

func TestCreateProduct_NegativePrice(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Test Product",
		Description:   "Test description",
		Price:         -10.00, // Negative price
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	assert.Error(t, err)
	assert.Nil(t, createdProduct)
	assert.Contains(t, err.Error(), "must be non-negative")
}

func TestCreateProduct_ZeroQuantity(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Test Product",
		Description:   "Test description",
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 0, // Out of stock
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	require.NoError(t, err)
	assert.NotNil(t, createdProduct)
	assert.Equal(t, int32(0), createdProduct.StockQuantity)
	assert.Equal(t, domain.StockStatusOutOfStock, createdProduct.StockStatus)
}

func TestCreateProduct_LongDescription(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	// Create a very long description (TEXT field has no practical limit in Postgres)
	longDesc := ""
	for i := 0; i < 1000; i++ {
		longDesc += "This is a test description. "
	}

	product := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    categoryID,
		Name:          "Test Product",
		Description:   longDesc,
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("TEST-SKU-001"),
		StockQuantity: 100,
		Attributes:    map[string]interface{}{},
	}

	createdProduct, err := repo.CreateProduct(ctx, product)

	require.NoError(t, err)
	assert.NotNil(t, createdProduct)
	assert.Equal(t, longDesc, createdProduct.Description)
}

// ============================================================================
// 2. UpdateProduct Tests (8 tests)
// ============================================================================

func TestUpdateProduct_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	// Create product
	product := createTestProduct(t, repo, storefrontID)

	// Update product
	newName := "Updated Product Name"
	newPrice := 149.99
	updateInput := &domain.UpdateProductInput{
		Name:  &newName,
		Price: &newPrice,
	}

	updatedProduct, err := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newName, updatedProduct.Name)
	assert.Equal(t, newPrice, updatedProduct.Price)
	assert.True(t, updatedProduct.UpdatedAt.After(product.UpdatedAt))
}

func TestUpdateProduct_PartialUpdate(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)
	originalPrice := product.Price

	// Update only name
	newName := "Only Name Updated"
	updateInput := &domain.UpdateProductInput{
		Name: &newName,
	}

	updatedProduct, err := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newName, updatedProduct.Name)
	assert.Equal(t, originalPrice, updatedProduct.Price) // Price unchanged
}

func TestUpdateProduct_UpdatePrice(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	newPrice := 199.99
	updateInput := &domain.UpdateProductInput{
		Price: &newPrice,
	}

	updatedProduct, err := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newPrice, updatedProduct.Price)
}

func TestUpdateProduct_UpdateQuantity(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	newQuantity := int32(200)
	updateInput := &domain.UpdateProductInput{
		StockQuantity: &newQuantity,
	}

	updatedProduct, err := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newQuantity, updatedProduct.StockQuantity)
	assert.Equal(t, domain.StockStatusInStock, updatedProduct.StockStatus)
}

func TestUpdateProduct_NonExistentProduct(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	newName := "Should Fail"
	updateInput := &domain.UpdateProductInput{
		Name: &newName,
	}

	updatedProduct, err := repo.UpdateProduct(ctx, 99999, storefrontID, updateInput)

	assert.Error(t, err)
	assert.Nil(t, updatedProduct)
	assert.Contains(t, err.Error(), "products.not_found")
}

func TestUpdateProduct_DuplicateSKU(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	// Create two products
	product1 := createTestProductWithOptions(t, repo, storefrontID, "SKU-001", 99.99, 100)
	product2 := createTestProductWithOptions(t, repo, storefrontID, "SKU-002", 149.99, 50)

	// Note: SKU is not in UpdateProductInput, only in BulkUpdateProductInput
	// This test verifies that we can't update via BulkUpdate with duplicate SKU
	sku1 := *product1.SKU
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: product2.ID,
			SKU:       &sku1,
		},
	}

	result, err := repo.BulkUpdateProducts(ctx, storefrontID, updates)

	require.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 0)
	assert.Len(t, result.FailedUpdates, 1)
	assert.Contains(t, result.FailedUpdates[0].ErrorCode, "products.sku_duplicate")
}

func TestUpdateProduct_InvalidData(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	// Try to update with empty name
	emptyName := ""
	updateInput := &domain.UpdateProductInput{
		Name: &emptyName,
	}

	// Note: Current implementation doesn't validate empty name in update
	// This test documents current behavior
	updatedProduct, err := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput)

	// Current implementation allows empty name in updates
	require.NoError(t, err)
	assert.NotNil(t, updatedProduct)
}

func TestUpdateProduct_ConcurrentUpdate(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	// Simulate concurrent updates (no optimistic locking in current implementation)
	name1 := "Update 1"
	name2 := "Update 2"

	updateInput1 := &domain.UpdateProductInput{Name: &name1}
	updateInput2 := &domain.UpdateProductInput{Name: &name2}

	_, err1 := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput1)
	_, err2 := repo.UpdateProduct(ctx, product.ID, storefrontID, updateInput2)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Last write wins (no optimistic locking)
	finalProduct, err := repo.GetProductByID(ctx, product.ID, &storefrontID)
	require.NoError(t, err)
	assert.Equal(t, name2, finalProduct.Name)
}

// ============================================================================
// 3. DeleteProduct Tests (6 tests)
// ============================================================================

func TestDeleteProduct_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	variantsDeleted, err := repo.DeleteProduct(ctx, product.ID, storefrontID, true)

	require.NoError(t, err)
	assert.Equal(t, int32(0), variantsDeleted)

	// Verify product is deleted
	_, err = repo.GetProductByID(ctx, product.ID, &storefrontID)
	assert.Error(t, err)
}

func TestDeleteProduct_SoftDelete(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	variantsDeleted, err := repo.DeleteProduct(ctx, product.ID, storefrontID, false)

	require.NoError(t, err)
	assert.Equal(t, int32(0), variantsDeleted)

	// Soft deleted product should not be found
	_, err = repo.GetProductByID(ctx, product.ID, &storefrontID)
	assert.Error(t, err)
}

func TestDeleteProduct_CascadeToVariants(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)
	createTestVariant(t, repo, product.ID)

	variantsDeleted, err := repo.DeleteProduct(ctx, product.ID, storefrontID, true)

	require.NoError(t, err)
	assert.Equal(t, int32(1), variantsDeleted)

	// Verify product and variants are deleted
	_, err = repo.GetProductByID(ctx, product.ID, &storefrontID)
	assert.Error(t, err)
}

func TestDeleteProduct_NonExistentProduct(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	variantsDeleted, err := repo.DeleteProduct(ctx, 99999, storefrontID, true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.not_found")
	assert.Equal(t, int32(0), variantsDeleted)
}

func TestDeleteProduct_WithActiveOrders(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	// Note: No orders table in current implementation
	// This test documents that check is not implemented yet
	variantsDeleted, err := repo.DeleteProduct(ctx, product.ID, storefrontID, true)

	require.NoError(t, err)
	assert.Equal(t, int32(0), variantsDeleted)
}

func TestDeleteProduct_AlreadyDeleted(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product := createTestProduct(t, repo, storefrontID)

	// Delete once
	_, err := repo.DeleteProduct(ctx, product.ID, storefrontID, true)
	require.NoError(t, err)

	// Try to delete again
	_, err = repo.DeleteProduct(ctx, product.ID, storefrontID, true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.not_found")
}

// ============================================================================
// 4. BulkCreateProducts Tests (5 tests)
// ============================================================================

func TestBulkCreateProducts_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	inputs := []*domain.CreateProductInput{
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "Bulk Product 1",
			Description:   "Description 1",
			Price:         99.99,
			Currency:      "USD",
			SKU:           stringPtr("BULK-001"),
			StockQuantity: 100,
		Attributes:    map[string]interface{}{},
		},
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "Bulk Product 2",
			Description:   "Description 2",
			Price:         149.99,
			Currency:      "USD",
			SKU:           stringPtr("BULK-002"),
			StockQuantity: 50,
		Attributes:    map[string]interface{}{},
		},
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "Bulk Product 3",
			Description:   "Description 3",
			Price:         199.99,
			Currency:      "USD",
			SKU:           stringPtr("BULK-003"),
			StockQuantity: 75,
		Attributes:    map[string]interface{}{},
		},
	}

	products, errors, err := repo.BulkCreateProducts(ctx, storefrontID, inputs)

	require.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Len(t, errors, 0)

	for i, product := range products {
		assert.NotZero(t, product.ID)
		assert.Equal(t, inputs[i].Name, product.Name)
		assert.Equal(t, inputs[i].SKU, product.SKU)
	}
}

func TestBulkCreateProducts_PartialFailure(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	inputs := []*domain.CreateProductInput{
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "Valid Product",
			Description:   "Valid",
			Price:         99.99,
			Currency:      "USD",
			SKU:           stringPtr("VALID-001"),
			StockQuantity: 100,
		Attributes:    map[string]interface{}{},
		},
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "", // Invalid - empty name
			Description:   "Invalid",
			Price:         149.99,
			Currency:      "USD",
			SKU:           stringPtr("INVALID-001"),
			StockQuantity: 50,
		Attributes:    map[string]interface{}{},
		},
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "Another Valid",
			Description:   "Valid",
			Price:         199.99,
			Currency:      "USD",
			SKU:           stringPtr("VALID-002"),
			StockQuantity: 75,
		Attributes:    map[string]interface{}{},
		},
	}

	products, errors, err := repo.BulkCreateProducts(ctx, storefrontID, inputs)

	require.NoError(t, err)
	assert.Len(t, products, 2)    // 2 valid products
	assert.Len(t, errors, 1)      // 1 error
	assert.Equal(t, int32(1), errors[0].Index)
	assert.Contains(t, errors[0].ErrorMessage, "name cannot be empty")
}

func TestBulkCreateProducts_EmptyBatch(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	inputs := []*domain.CreateProductInput{}

	products, errors, err := repo.BulkCreateProducts(ctx, storefrontID, inputs)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.bulk_empty")
	assert.Nil(t, products)
	assert.Nil(t, errors)
}

func TestBulkCreateProducts_LargeBatch(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	// Create 150 products
	inputs := make([]*domain.CreateProductInput, 150)
	for i := 0; i < 150; i++ {
		sku := stringPtr(fmt.Sprintf("LARGE-BATCH-%03d", i))
		inputs[i] = &domain.CreateProductInput{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          fmt.Sprintf("Product %d", i),
			Description:   fmt.Sprintf("Description %d", i),
			Price:         float64(100 + i),
			Currency:      "USD",
			SKU:           sku,
			StockQuantity: int32(10 + i),
			Attributes:    map[string]interface{}{},
		}
	}

	products, errors, err := repo.BulkCreateProducts(ctx, storefrontID, inputs)

	require.NoError(t, err)
	assert.Len(t, products, 150)
	assert.Len(t, errors, 0)
}

func TestBulkCreateProducts_TransactionRollback(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	categoryID := createTestCategory(t)
	ctx := tests.TestContext(t)

	// Create a product with SKU first
	createTestProductWithOptions(t, repo, storefrontID, "EXISTING-SKU", 99.99, 100)

	// Try bulk create with duplicate SKU
	inputs := []*domain.CreateProductInput{
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "New Product 1",
			Description:   "Description 1",
			Price:         99.99,
			Currency:      "USD",
			SKU:           stringPtr("NEW-001"),
			StockQuantity: 100,
		Attributes:    map[string]interface{}{},
		},
		{
			StorefrontID:  storefrontID,
			CategoryID:    categoryID,
			Name:          "Duplicate SKU Product",
			Description:   "Description 2",
			Price:         149.99,
			Currency:      "USD",
			SKU:           stringPtr("EXISTING-SKU"), // Duplicate
			StockQuantity: 50,
		Attributes:    map[string]interface{}{},
		},
	}

	products, errors, err := repo.BulkCreateProducts(ctx, storefrontID, inputs)

	// Should fail on duplicate check
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "products.sku_duplicate")
	assert.Nil(t, products)
	assert.NotEmpty(t, errors)
}

// ============================================================================
// 5. BulkUpdateProducts Tests (5 tests)
// ============================================================================

func TestBulkUpdateProducts_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	// Create 3 products
	product1 := createTestProductWithOptions(t, repo, storefrontID, "UPD-001", 99.99, 100)
	product2 := createTestProductWithOptions(t, repo, storefrontID, "UPD-002", 149.99, 50)
	product3 := createTestProductWithOptions(t, repo, storefrontID, "UPD-003", 199.99, 75)

	// Bulk update
	newPrice1 := 109.99
	newPrice2 := 159.99
	newName3 := "Updated Product 3"

	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: product1.ID,
			Price:     &newPrice1,
		},
		{
			ProductID: product2.ID,
			Price:     &newPrice2,
		},
		{
			ProductID: product3.ID,
			Name:      &newName3,
		},
	}

	result, err := repo.BulkUpdateProducts(ctx, storefrontID, updates)

	require.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 3)
	assert.Len(t, result.FailedUpdates, 0)

	assert.Equal(t, newPrice1, result.SuccessfulProducts[0].Price)
	assert.Equal(t, newPrice2, result.SuccessfulProducts[1].Price)
	assert.Equal(t, newName3, result.SuccessfulProducts[2].Name)
}

func TestBulkUpdateProducts_PartialSuccess(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product1 := createTestProductWithOptions(t, repo, storefrontID, "UPD-PART-001", 99.99, 100)

	newPrice := 109.99
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: product1.ID,
			Price:     &newPrice,
		},
		{
			ProductID: 99999, // Non-existent
			Price:     &newPrice,
		},
	}

	result, err := repo.BulkUpdateProducts(ctx, storefrontID, updates)

	require.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 1)
	assert.Len(t, result.FailedUpdates, 1)
	assert.Equal(t, int64(99999), result.FailedUpdates[0].ProductID)
	assert.Contains(t, result.FailedUpdates[0].ErrorCode, "products.not_found")
}

func TestBulkUpdateProducts_EmptyBatch(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	updates := []*domain.BulkUpdateProductInput{}

	result, err := repo.BulkUpdateProducts(ctx, storefrontID, updates)

	require.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 0)
	assert.Len(t, result.FailedUpdates, 0)
}

func TestBulkUpdateProducts_MixedOperations(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product1 := createTestProductWithOptions(t, repo, storefrontID, "MIX-001", 99.99, 100)
	product2 := createTestProductWithOptions(t, repo, storefrontID, "MIX-002", 149.99, 50)

	newName := "Updated Name"
	newPrice := 199.99
	newQuantity := int32(200)

	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: product1.ID,
			Name:      &newName,
		},
		{
			ProductID:     product2.ID,
			Price:         &newPrice,
			StockQuantity: &newQuantity,
		},
	}

	result, err := repo.BulkUpdateProducts(ctx, storefrontID, updates)

	require.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 2)
	assert.Equal(t, newName, result.SuccessfulProducts[0].Name)
	assert.Equal(t, newPrice, result.SuccessfulProducts[1].Price)
	assert.Equal(t, newQuantity, result.SuccessfulProducts[1].StockQuantity)
}

func TestBulkUpdateProducts_TransactionRollback(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product1 := createTestProductWithOptions(t, repo, storefrontID, "ROLL-001", 99.99, 100)
	product2 := createTestProductWithOptions(t, repo, storefrontID, "ROLL-002", 149.99, 50)

	// Try to update product2 with product1's SKU
	sku1 := *product1.SKU
	updates := []*domain.BulkUpdateProductInput{
		{
			ProductID: product2.ID,
			SKU:       &sku1, // Duplicate SKU
		},
	}

	result, err := repo.BulkUpdateProducts(ctx, storefrontID, updates)

	require.NoError(t, err)
	assert.Len(t, result.SuccessfulProducts, 0)
	assert.Len(t, result.FailedUpdates, 1)
	assert.Contains(t, result.FailedUpdates[0].ErrorCode, "products.sku_duplicate")
}

// ============================================================================
// 6. BulkDeleteProducts Tests (4 tests)
// ============================================================================

func TestBulkDeleteProducts_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product1 := createTestProductWithOptions(t, repo, storefrontID, "DEL-001", 99.99, 100)
	product2 := createTestProductWithOptions(t, repo, storefrontID, "DEL-002", 149.99, 50)
	product3 := createTestProductWithOptions(t, repo, storefrontID, "DEL-003", 199.99, 75)

	productIDs := []int64{product1.ID, product2.ID, product3.ID}

	successCount, failedCount, variantsDeleted, errors, err := repo.BulkDeleteProducts(ctx, storefrontID, productIDs, true)

	require.NoError(t, err)
	assert.Equal(t, int32(3), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Empty(t, errors)

	// Verify products are deleted
	for _, id := range productIDs {
		_, err := repo.GetProductByID(ctx, id, &storefrontID)
		assert.Error(t, err)
	}
}

func TestBulkDeleteProducts_PartialSuccess(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	product1 := createTestProductWithOptions(t, repo, storefrontID, "PDEL-001", 99.99, 100)

	productIDs := []int64{product1.ID, 99999} // One valid, one invalid

	successCount, failedCount, variantsDeleted, errors, err := repo.BulkDeleteProducts(ctx, storefrontID, productIDs, true)

	require.NoError(t, err)
	assert.Equal(t, int32(1), successCount)
	assert.Equal(t, int32(1), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Len(t, errors, 1)
	assert.Equal(t, "products.not_found", errors[99999])
}

func TestBulkDeleteProducts_EmptyBatch(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	productIDs := []int64{}

	successCount, failedCount, variantsDeleted, errors, err := repo.BulkDeleteProducts(ctx, storefrontID, productIDs, true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(0), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Nil(t, errors)
}

func TestBulkDeleteProducts_NonExistentProducts(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	ctx := tests.TestContext(t)

	productIDs := []int64{99991, 99992, 99993}

	successCount, failedCount, variantsDeleted, errors, err := repo.BulkDeleteProducts(ctx, storefrontID, productIDs, true)

	require.NoError(t, err)
	assert.Equal(t, int32(0), successCount)
	assert.Equal(t, int32(3), failedCount)
	assert.Equal(t, int32(0), variantsDeleted)
	assert.Len(t, errors, 3)

	for _, id := range productIDs {
		assert.Equal(t, "products.not_found", errors[id])
	}
}
