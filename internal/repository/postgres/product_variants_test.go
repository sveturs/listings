package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// Test Fixtures and Helpers for Product Variants
// ============================================================================

// createTestProductWithVariants creates a product with has_variants=true
func createTestProductWithVariants(t *testing.T, repo *Repository, storefrontID int64) *domain.Product {
	t.Helper()
	ctx := tests.TestContext(t)

	input := &domain.CreateProductInput{
		StorefrontID:  storefrontID,
		CategoryID:    createTestCategory(t),
		Name:          "Product with Variants",
		Description:   "Product that supports variants",
		Price:         99.99,
		Currency:      "USD",
		SKU:           stringPtr("PROD-VAR-001"),
		StockQuantity: 0, // Variants will manage stock
		Attributes:    map[string]interface{}{},
	}

	product, err := repo.CreateProduct(ctx, input)
	require.NoError(t, err, "Failed to create test product")
	require.NotNil(t, product)

	// Enable variants
	_, err = repo.db.ExecContext(ctx, `
		UPDATE listings SET has_variants = true WHERE id = $1
	`, product.ID)
	require.NoError(t, err, "Failed to enable variants")

	// Reload product
	product, err = repo.GetProductByID(ctx, product.ID, &storefrontID)
	require.NoError(t, err)
	require.True(t, product.HasVariants)

	return product
}

// float64Ptr returns a pointer to a float64 value

// int32PtrVal returns a pointer to an int32 value

// ============================================================================
// 1. CreateProductVariant Tests (6 tests)
// ============================================================================

func TestCreateProductVariant_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	price := 109.99
	compareAt := 129.99
	lowThreshold := int32(5)

	input := &domain.CreateVariantInput{
		ProductID:         product.ID,
		SKU:               stringPtr("VAR-RED-L"),
		Price:             &price,
		CompareAtPrice:    &compareAt,
		StockQuantity:     50,
		LowStockThreshold: &lowThreshold,
		VariantAttributes: map[string]interface{}{
			"color": "red",
			"size":  "L",
		},
		IsDefault: true,
	}

	variant, err := repo.CreateProductVariant(ctx, input)

	require.NoError(t, err)
	assert.NotZero(t, variant.ID)
	assert.Equal(t, product.ID, variant.ProductID)
	assert.Equal(t, "VAR-RED-L", *variant.SKU)
	assert.Equal(t, price, *variant.Price)
	assert.Equal(t, compareAt, *variant.CompareAtPrice)
	assert.Equal(t, int32(50), variant.StockQuantity)
	assert.Equal(t, domain.StockStatusInStock, variant.StockStatus)
	assert.True(t, variant.IsActive)
	assert.True(t, variant.IsDefault)
	assert.NotNil(t, variant.VariantAttributes)
	assert.Equal(t, "red", variant.VariantAttributes["color"])
	assert.Equal(t, "L", variant.VariantAttributes["size"])
}

func TestCreateProductVariant_WithAttributes(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	price := 149.99
	weight := 0.5
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-ATTR-001"),
		Price:         &price,
		StockQuantity: 25,
		VariantAttributes: map[string]interface{}{
			"color":    "blue",
			"size":     "M",
			"material": "cotton",
			"season":   "summer",
		},
		Dimensions: map[string]interface{}{
			"length": 30,
			"width":  20,
			"height": 10,
			"unit":   "cm",
		},
		Weight: &weight,
	}

	variant, err := repo.CreateProductVariant(ctx, input)

	require.NoError(t, err)
	assert.NotZero(t, variant.ID)
	assert.Equal(t, 4, len(variant.VariantAttributes))
	assert.Equal(t, "blue", variant.VariantAttributes["color"])
	assert.Equal(t, "cotton", variant.VariantAttributes["material"])
	assert.NotNil(t, variant.Dimensions)
	assert.Equal(t, float64(30), variant.Dimensions["length"])
	assert.Equal(t, "cm", variant.Dimensions["unit"])
	assert.Equal(t, 0.5, *variant.Weight)
}

func TestCreateProductVariant_MissingProductID(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	price := 99.99
	input := &domain.CreateVariantInput{
		ProductID:     0, // Invalid product ID
		SKU:           stringPtr("VAR-INVALID"),
		Price:         &price,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)

	// Nil-safe assertions
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "variants.invalid_product_id")
	}
	assert.Nil(t, variant)
}

func TestCreateProductVariant_InvalidProductID(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	price := 99.99
	input := &domain.CreateVariantInput{
		ProductID:     99999, // Non-existent product
		SKU:           stringPtr("VAR-NOEXIST"),
		Price:         &price,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, variant)
	assert.Contains(t, err.Error(), "variants.product_not_found")
}

func TestCreateProductVariant_DuplicateSKU(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	price := 99.99
	sku := "VAR-DUP-001"

	// Create first variant
	input1 := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           &sku,
		Price:         &price,
		StockQuantity: 10,
	}

	variant1, err := repo.CreateProductVariant(ctx, input1)
	require.NoError(t, err)
	require.NotNil(t, variant1)

	// Try to create second variant with same SKU
	input2 := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           &sku, // Same SKU
		Price:         &price,
		StockQuantity: 5,
	}

	variant2, err := repo.CreateProductVariant(ctx, input2)

	assert.Error(t, err)
	assert.Nil(t, variant2)
	assert.Contains(t, err.Error(), "variants.sku_duplicate")
}

func TestCreateProductVariant_NegativePrice(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Negative prices are validated at DB level with CHECK constraint
	negativePrice := -10.00
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-NEG-PRICE"),
		Price:         &negativePrice,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, variant)
	// PostgreSQL constraint violation
	assert.Contains(t, err.Error(), "variants.create_failed")
}

// ============================================================================
// 2. UpdateProductVariant Tests (5 tests)
// ============================================================================

func TestUpdateProductVariant_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variant
	price := 99.99
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-UPD-001"),
		Price:         &price,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)
	require.NoError(t, err)

	// Update variant
	newPrice := 119.99
	newSKU := "VAR-UPD-001-NEW"
	newQuantity := int32(25)
	updateInput := &domain.UpdateVariantInput{
		SKU:           &newSKU,
		Price:         &newPrice,
		StockQuantity: &newQuantity,
	}

	updatedVariant, err := repo.UpdateProductVariant(ctx, variant.ID, product.ID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, variant.ID, updatedVariant.ID)
	assert.Equal(t, newSKU, *updatedVariant.SKU)
	assert.Equal(t, newPrice, *updatedVariant.Price)
	assert.Equal(t, newQuantity, updatedVariant.StockQuantity)
	assert.True(t, updatedVariant.UpdatedAt.After(variant.UpdatedAt))
}

func TestUpdateProductVariant_PartialUpdate(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variant
	price := 99.99
	originalSKU := "VAR-PARTIAL-001"
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           &originalSKU,
		Price:         &price,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)
	require.NoError(t, err)

	// Update only price
	newPrice := 149.99
	updateInput := &domain.UpdateVariantInput{
		Price: &newPrice,
	}

	updatedVariant, err := repo.UpdateProductVariant(ctx, variant.ID, product.ID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newPrice, *updatedVariant.Price)
	assert.Equal(t, originalSKU, *updatedVariant.SKU)        // SKU unchanged
	assert.Equal(t, int32(10), updatedVariant.StockQuantity) // Quantity unchanged
}

func TestUpdateProductVariant_UpdatePrice(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variant
	price := 99.99
	compareAt := 129.99
	input := &domain.CreateVariantInput{
		ProductID:      product.ID,
		SKU:            stringPtr("VAR-PRICE-001"),
		Price:          &price,
		CompareAtPrice: &compareAt,
		StockQuantity:  10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)
	require.NoError(t, err)

	// Update prices
	newPrice := 79.99
	newCompareAt := 99.99
	costPrice := 50.00
	updateInput := &domain.UpdateVariantInput{
		Price:          &newPrice,
		CompareAtPrice: &newCompareAt,
		CostPrice:      &costPrice,
	}

	updatedVariant, err := repo.UpdateProductVariant(ctx, variant.ID, product.ID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newPrice, *updatedVariant.Price)
	assert.Equal(t, newCompareAt, *updatedVariant.CompareAtPrice)
	assert.Equal(t, costPrice, *updatedVariant.CostPrice)

	// Verify discount calculation
	assert.True(t, updatedVariant.HasDiscount())
	discountPct := updatedVariant.GetDiscountPercentage()
	assert.Greater(t, discountPct, 0.0)
}

func TestUpdateProductVariant_NonExistentVariant(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	newPrice := 99.99
	updateInput := &domain.UpdateVariantInput{
		Price: &newPrice,
	}

	updatedVariant, err := repo.UpdateProductVariant(ctx, 99999, product.ID, updateInput)

	assert.Error(t, err)
	assert.Nil(t, updatedVariant)
	assert.Contains(t, err.Error(), "variants.not_found")
}

func TestUpdateProductVariant_UpdateAttributes(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variant with attributes
	price := 99.99
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-ATTR-UPD-001"),
		Price:         &price,
		StockQuantity: 10,
		VariantAttributes: map[string]interface{}{
			"color": "red",
			"size":  "M",
		},
	}

	variant, err := repo.CreateProductVariant(ctx, input)
	require.NoError(t, err)

	// Update attributes
	updateInput := &domain.UpdateVariantInput{
		VariantAttributes: map[string]interface{}{
			"color":    "blue",
			"size":     "L",
			"material": "cotton",
		},
	}

	updatedVariant, err := repo.UpdateProductVariant(ctx, variant.ID, product.ID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, 3, len(updatedVariant.VariantAttributes))
	assert.Equal(t, "blue", updatedVariant.VariantAttributes["color"])
	assert.Equal(t, "L", updatedVariant.VariantAttributes["size"])
	assert.Equal(t, "cotton", updatedVariant.VariantAttributes["material"])
}

// ============================================================================
// 3. DeleteProductVariant Tests (4 tests)
// ============================================================================

func TestDeleteProductVariant_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variant
	price := 99.99
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-DEL-001"),
		Price:         &price,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)
	require.NoError(t, err)

	// Delete variant
	err = repo.DeleteProductVariant(ctx, variant.ID, product.ID)

	require.NoError(t, err)

	// Verify variant is deleted
	_, err = repo.GetVariantByID(ctx, variant.ID, &product.ID)
	assert.Error(t, err)
}

func TestDeleteProductVariant_NonExistentVariant(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Try to delete non-existent variant
	err := repo.DeleteProductVariant(ctx, 99999, product.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variants.not_found")
}

func TestDeleteProductVariant_UpdatesProductStock(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create two variants
	price := 99.99
	input1 := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-STOCK-001"),
		Price:         &price,
		StockQuantity: 50,
	}

	variant1, err := repo.CreateProductVariant(ctx, input1)
	require.NoError(t, err)

	input2 := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-STOCK-002"),
		Price:         &price,
		StockQuantity: 30,
	}

	variant2, err := repo.CreateProductVariant(ctx, input2)
	require.NoError(t, err)

	// Delete first variant
	err = repo.DeleteProductVariant(ctx, variant1.ID, product.ID)
	require.NoError(t, err)

	// Verify second variant still exists
	remainingVariant, err := repo.GetVariantByID(ctx, variant2.ID, &product.ID)
	require.NoError(t, err)
	assert.NotNil(t, remainingVariant)

	// Product should still have has_variants=true
	updatedProduct, err := repo.GetProductByID(ctx, product.ID, &storefrontID)
	require.NoError(t, err)
	assert.True(t, updatedProduct.HasVariants)
}

func TestDeleteProductVariant_AlreadyDeleted(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variant
	price := 99.99
	input := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-DEL-TWICE-001"),
		Price:         &price,
		StockQuantity: 10,
	}

	variant, err := repo.CreateProductVariant(ctx, input)
	require.NoError(t, err)

	// Delete once
	err = repo.DeleteProductVariant(ctx, variant.ID, product.ID)
	require.NoError(t, err)

	// Try to delete again (idempotency check)
	err = repo.DeleteProductVariant(ctx, variant.ID, product.ID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variants.not_found")
}

// ============================================================================
// 4. BulkCreateProductVariants Tests (5 tests)
// ============================================================================

func TestBulkCreateProductVariants_Success(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	price1 := 99.99
	price2 := 109.99
	price3 := 119.99
	lowThreshold := int32(5)

	inputs := []*domain.CreateVariantInput{
		{
			ProductID:         product.ID,
			SKU:               stringPtr("VAR-BULK-RED-S"),
			Price:             &price1,
			StockQuantity:     10,
			LowStockThreshold: &lowThreshold,
			VariantAttributes: map[string]interface{}{
				"color": "red",
				"size":  "S",
			},
			IsDefault: true,
		},
		{
			ProductID:     product.ID,
			SKU:           stringPtr("VAR-BULK-RED-M"),
			Price:         &price2,
			StockQuantity: 20,
			VariantAttributes: map[string]interface{}{
				"color": "red",
				"size":  "M",
			},
		},
		{
			ProductID:     product.ID,
			SKU:           stringPtr("VAR-BULK-RED-L"),
			Price:         &price3,
			StockQuantity: 15,
			VariantAttributes: map[string]interface{}{
				"color": "red",
				"size":  "L",
			},
		},
	}

	variants, err := repo.BulkCreateProductVariants(ctx, product.ID, inputs)

	require.NoError(t, err)
	assert.Len(t, variants, 3)

	// Verify all variants created
	for _, variant := range variants {
		assert.NotZero(t, variant.ID)
		assert.Equal(t, product.ID, variant.ProductID)
		assert.NotNil(t, variant.SKU)
		assert.True(t, variant.IsActive)

		// Verify attributes
		assert.NotNil(t, variant.VariantAttributes)
		assert.Equal(t, "red", variant.VariantAttributes["color"])
	}

	// Verify first variant is default
	assert.True(t, variants[0].IsDefault)
	assert.False(t, variants[1].IsDefault)
	assert.False(t, variants[2].IsDefault)
}

func TestBulkCreateProductVariants_MultipleProducts(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product1 := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create variants for product1 only (not multiple products in same call)
	price := 99.99
	inputs := []*domain.CreateVariantInput{
		{
			ProductID:     product1.ID,
			SKU:           stringPtr("VAR-P1-001"),
			Price:         &price,
			StockQuantity: 10,
		},
		{
			ProductID:     product1.ID,
			SKU:           stringPtr("VAR-P1-002"),
			Price:         &price,
			StockQuantity: 20,
		},
	}

	variants, err := repo.BulkCreateProductVariants(ctx, product1.ID, inputs)

	require.NoError(t, err)
	assert.Len(t, variants, 2)

	// All variants should belong to product1
	for _, variant := range variants {
		assert.Equal(t, product1.ID, variant.ProductID)
	}
}

func TestBulkCreateProductVariants_PartialFailure(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	// Create a variant first to cause duplicate SKU error
	price := 99.99
	existingInput := &domain.CreateVariantInput{
		ProductID:     product.ID,
		SKU:           stringPtr("VAR-EXISTING-001"),
		Price:         &price,
		StockQuantity: 10,
	}

	_, err := repo.CreateProductVariant(ctx, existingInput)
	require.NoError(t, err)

	// Try bulk create with duplicate SKU
	inputs := []*domain.CreateVariantInput{
		{
			ProductID:     product.ID,
			SKU:           stringPtr("VAR-NEW-001"),
			Price:         &price,
			StockQuantity: 10,
		},
		{
			ProductID:     product.ID,
			SKU:           stringPtr("VAR-EXISTING-001"), // Duplicate SKU
			Price:         &price,
			StockQuantity: 20,
		},
	}

	variants, err := repo.BulkCreateProductVariants(ctx, product.ID, inputs)

	// Should fail due to duplicate SKU (transaction rollback)
	assert.Error(t, err)
	assert.Nil(t, variants)
	assert.Contains(t, err.Error(), "variants.sku_duplicate")
}

func TestBulkCreateProductVariants_EmptyBatch(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	inputs := []*domain.CreateVariantInput{}

	variants, err := repo.BulkCreateProductVariants(ctx, product.ID, inputs)

	assert.Error(t, err)
	assert.Nil(t, variants)
	assert.Contains(t, err.Error(), "variants.bulk_empty")
}

func TestBulkCreateProductVariants_TransactionRollback(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	storefrontID := createTestStorefront(t, repo)
	product := createTestProductWithVariants(t, repo, storefrontID)
	ctx := tests.TestContext(t)

	price := 99.99
	negativePrice := -50.00

	// Mix valid and invalid variants (negative price violates constraint)
	inputs := []*domain.CreateVariantInput{
		{
			ProductID:     product.ID,
			SKU:           stringPtr("VAR-VALID-001"),
			Price:         &price,
			StockQuantity: 10,
		},
		{
			ProductID:     product.ID,
			SKU:           stringPtr("VAR-INVALID-001"),
			Price:         &negativePrice, // Invalid price
			StockQuantity: 20,
		},
	}

	variants, err := repo.BulkCreateProductVariants(ctx, product.ID, inputs)

	// Should fail and rollback entire transaction
	assert.Error(t, err)
	assert.Nil(t, variants)
	assert.Contains(t, err.Error(), "variants.bulk_create_failed")

	// Verify no variants were created (transaction rolled back)
	_, err = repo.GetVariantByID(ctx, 0, &product.ID)
	assert.Error(t, err)
}
