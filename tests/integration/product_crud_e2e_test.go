//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/tests"
)

// ============================================================================
// E2E Workflow Test: Create → Get → Update → Get → Delete → Get
// ============================================================================

// TestProductCRUD_E2E_FullWorkflow verifies complete CRUD lifecycle
func TestProductCRUD_E2E_FullWorkflow(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	t.Log("========== E2E WORKFLOW: Create → Get → Update → Get → Delete → Get ==========")

	// ===== STEP 1: Create Product =====
	t.Log("Step 1: Creating product...")

	attributes, _ := structpb.NewStruct(map[string]interface{}{
		"brand": "E2E Brand",
		"model": "E2E-001",
	})

	sku := "E2E-PROD-999"
	barcode := "1234567899999"

	createReq := &pb.CreateProductRequest{
		StorefrontId:  storefrontID,
		Name:          "E2E Test Product",
		Description:   "Product for end-to-end testing",
		Price:         299.99,
		Currency:      "USD",
		CategoryId:    9000,
		Sku:           &sku,
		Barcode:       &barcode,
		StockQuantity: 100,
		IsActive:      true,
		Attributes:    attributes,
	}

	createResp, err := client.CreateProduct(ctx, createReq)
	require.NoError(t, err, "CreateProduct should succeed")
	require.NotNil(t, createResp, "Create response should not be nil")
	require.NotNil(t, createResp.Product, "Created product should be returned")

	productID := createResp.Product.Id
	assert.Greater(t, productID, int64(0), "Product ID should be positive")
	assert.Equal(t, "E2E Test Product", createResp.Product.Name)
	assert.Equal(t, 299.99, createResp.Product.Price)
	assert.Equal(t, int32(100), createResp.Product.StockQuantity)

	t.Logf("  ✓ Product created with ID: %d", productID)

	// ===== STEP 2: Get Product (verify creation) =====
	t.Log("Step 2: Getting created product...")

	getReq := &pb.GetProductRequest{
		ProductId:    productID,
		StorefrontId: &storefrontID,
	}

	getResp1, err := client.GetProduct(ctx, getReq)
	require.NoError(t, err, "GetProduct should succeed")
	require.NotNil(t, getResp1.Product, "Product should be found")
	assert.Equal(t, productID, getResp1.Product.Id)
	assert.Equal(t, "E2E Test Product", getResp1.Product.Name)
	assert.Equal(t, 299.99, getResp1.Product.Price)
	assert.Equal(t, "E2E-PROD-999", *getResp1.Product.Sku)

	t.Log("  ✓ Product retrieved successfully")

	// ===== STEP 3: Update Product =====
	t.Log("Step 3: Updating product...")

	newName := "E2E Test Product (Updated)"
	newPrice := 349.99

	updateReq := &pb.UpdateProductRequest{
		ProductId:    productID,
		StorefrontId: storefrontID,
		Name:         &newName,
		Price:        &newPrice,
	}

	updateResp, err := client.UpdateProduct(ctx, updateReq)
	require.NoError(t, err, "UpdateProduct should succeed")
	require.NotNil(t, updateResp, "Update response should not be nil")
	require.NotNil(t, updateResp.Product, "Updated product should be returned")
	assert.Equal(t, newName, updateResp.Product.Name)
	assert.Equal(t, newPrice, updateResp.Product.Price)

	t.Log("  ✓ Product updated successfully")

	// ===== STEP 4: Get Product (verify update) =====
	t.Log("Step 4: Getting updated product...")

	getResp2, err := client.GetProduct(ctx, getReq)
	require.NoError(t, err, "GetProduct should succeed after update")
	assert.Equal(t, newName, getResp2.Product.Name, "Name should be updated")
	assert.Equal(t, newPrice, getResp2.Product.Price, "Price should be updated")

	t.Log("  ✓ Updated product retrieved successfully")

	// ===== STEP 5: Delete Product =====
	t.Log("Step 5: Deleting product...")

	deleteReq := &pb.DeleteProductRequest{
		ProductId:    productID,
		StorefrontId: storefrontID,
		HardDelete:   true,
	}

	deleteResp, err := client.DeleteProduct(ctx, deleteReq)
	require.NoError(t, err, "DeleteProduct should succeed")
	assert.True(t, deleteResp.Success, "Deletion should be successful")

	t.Log("  ✓ Product deleted successfully")

	// ===== STEP 6: Get Product (should NOT exist) =====
	t.Log("Step 6: Verifying product is deleted...")

	_, err = client.GetProduct(ctx, getReq)
	require.Error(t, err, "GetProduct should fail after deletion")

	t.Log("  ✓ Product not found after deletion (expected)")
	t.Log("========== E2E WORKFLOW COMPLETED SUCCESSFULLY ==========")
}

// ============================================================================
// E2E Workflow Test: Soft Delete
// ============================================================================

// TestProductCRUD_E2E_SoftDeleteWorkflow verifies soft delete workflow
func TestProductCRUD_E2E_SoftDeleteWorkflow(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	t.Log("========== E2E SOFT DELETE WORKFLOW ==========")

	// Create product
	attributes, _ := structpb.NewStruct(map[string]interface{}{
		"test": "soft_delete",
	})

	softSku := "SOFT-DEL-E2E-001"

	createReq := &pb.CreateProductRequest{
		StorefrontId:  storefrontID,
		Name:          "Soft Delete Test Product",
		Description:   "Product for soft delete testing",
		Price:         99.99,
		Currency:      "USD",
		CategoryId:    9000,
		Sku:           &softSku,
		StockQuantity: 50,
		IsActive:      true,
		Attributes:    attributes,
	}

	createResp, err := client.CreateProduct(ctx, createReq)
	require.NoError(t, err, "CreateProduct should succeed")
	productID := createResp.Product.Id

	t.Logf("  ✓ Product created with ID: %d", productID)

	// Verify product exists
	getReq := &pb.GetProductRequest{
		ProductId:    productID,
		StorefrontId: &storefrontID,
	}

	_, err = client.GetProduct(ctx, getReq)
	require.NoError(t, err, "Product should exist before soft delete")

	t.Log("  ✓ Product exists before soft delete")

	// Soft delete
	deleteReq := &pb.DeleteProductRequest{
		ProductId:    productID,
		StorefrontId: storefrontID,
		HardDelete:   false, // Soft delete
	}

	deleteResp, err := client.DeleteProduct(ctx, deleteReq)
	require.NoError(t, err, "Soft delete should succeed")
	assert.True(t, deleteResp.Success, "Soft deletion should be successful")

	t.Log("  ✓ Product soft-deleted successfully")

	// CRITICAL: Verify product is NOT found after soft delete
	_, err = client.GetProduct(ctx, getReq)
	require.Error(t, err, "GetProduct should fail for soft-deleted product (BUG FIX VALIDATION)")

	t.Log("  ✓ Product not found after soft delete (bug fix working!)")

	// Verify database state
	var deletedAt *string
	err = testDB.DB.QueryRow("SELECT deleted_at FROM listings WHERE id = $1 AND source_type = 'b2c'", productID).Scan(&deletedAt)
	require.NoError(t, err, "Should query deleted_at")
	assert.NotNil(t, deletedAt, "deleted_at should be set")
	assert.NotEmpty(t, *deletedAt, "deleted_at should not be empty")

	t.Log("  ✓ Database confirms soft delete (deleted_at is set)")

	// Verify data is NOT physically deleted
	var title string
	err = testDB.DB.QueryRow("SELECT title FROM listings WHERE id = $1 AND source_type = 'b2c'", productID).Scan(&title)
	require.NoError(t, err, "Data should still exist in database")
	assert.Equal(t, "Soft Delete Test Product", title, "Data should be intact")

	t.Log("  ✓ Data is intact (soft delete preserves data)")
	t.Log("========== SOFT DELETE WORKFLOW COMPLETED ==========")
}

// ============================================================================
// E2E Workflow Test: Product with Variants
// ============================================================================

// TestProductCRUD_E2E_WithVariantsWorkflow verifies product+variants lifecycle
func TestProductCRUD_E2E_WithVariantsWorkflow(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	t.Log("========== E2E WORKFLOW: Product with Variants ==========")

	// ===== Create Product with has_variants=true =====
	t.Log("Step 1: Creating product with variants enabled...")

	attributes, _ := structpb.NewStruct(map[string]interface{}{
		"type": "clothing",
	})

	variantSku := "VAR-E2E-001"

	createReq := &pb.CreateProductRequest{
		StorefrontId:  storefrontID,
		Name:          "E2E Variant Test Product",
		Description:   "Product with multiple variants",
		Price:         79.99,
		Currency:      "USD",
		CategoryId:    9001,
		Sku:           &variantSku,
		StockQuantity: 0, // Variants will have stock
		IsActive:      true,
		Attributes:    attributes,
		HasVariants:   true,
	}

	createResp, err := client.CreateProduct(ctx, createReq)
	require.NoError(t, err, "CreateProduct should succeed")
	productID := createResp.Product.Id

	t.Logf("  ✓ Product created with ID: %d", productID)

	// ===== Create Variants =====
	t.Log("Step 2: Creating variants...")

	variantAttrsS, _ := structpb.NewStruct(map[string]interface{}{
		"size": "S",
	})
	variantAttrsM, _ := structpb.NewStruct(map[string]interface{}{
		"size": "M",
	})
	variantAttrsL, _ := structpb.NewStruct(map[string]interface{}{
		"size": "L",
	})

	varSkuS := "VAR-E2E-001-S"
	varSkuM := "VAR-E2E-001-M"
	varSkuL := "VAR-E2E-001-L"
	varPrice := 79.99

	variants := []*pb.ProductVariantInput{
		{
			Sku:               &varSkuS,
			Price:             &varPrice,
			StockQuantity:     30,
			VariantAttributes: variantAttrsS,
			IsActive:          true,
			IsDefault:         false,
		},
		{
			Sku:               &varSkuM,
			Price:             &varPrice,
			StockQuantity:     40,
			VariantAttributes: variantAttrsM,
			IsActive:          true,
			IsDefault:         true, // Default variant
		},
		{
			Sku:               &varSkuL,
			Price:             &varPrice,
			StockQuantity:     20,
			VariantAttributes: variantAttrsL,
			IsActive:          true,
			IsDefault:         false,
		},
	}

	bulkCreateReq := &pb.BulkCreateProductVariantsRequest{
		ProductId: productID,
		Variants:  variants,
	}

	bulkCreateResp, err := client.BulkCreateProductVariants(ctx, bulkCreateReq)
	require.NoError(t, err, "BulkCreateProductVariants should succeed")
	assert.Equal(t, int32(3), bulkCreateResp.SuccessfulCount, "Should create 3 variants")

	t.Log("  ✓ 3 variants created successfully")

	// ===== Get Product (verify variants) =====
	t.Log("Step 3: Getting product with variants...")

	getReq := &pb.GetProductRequest{
		ProductId:    productID,
		StorefrontId: &storefrontID,
	}

	getResp, err := client.GetProduct(ctx, getReq)
	require.NoError(t, err, "GetProduct should succeed")
	assert.True(t, getResp.Product.HasVariants, "Product should have variants")
	assert.Len(t, getResp.Product.Variants, 3, "Should have 3 variants")

	t.Log("  ✓ Product retrieved with 3 variants")

	// ===== Delete Product (cascade to variants) =====
	t.Log("Step 4: Deleting product (cascade to variants)...")

	deleteReq := &pb.DeleteProductRequest{
		ProductId:    productID,
		StorefrontId: storefrontID,
		HardDelete:   true,
	}

	deleteResp, err := client.DeleteProduct(ctx, deleteReq)
	require.NoError(t, err, "DeleteProduct should succeed")
	assert.True(t, deleteResp.Success, "Deletion should be successful")
	assert.Equal(t, int32(3), deleteResp.VariantsDeleted, "Should cascade delete 3 variants")

	t.Log("  ✓ Product and 3 variants deleted (cascade)")

	// ===== Verify all deleted =====
	var variantCount int
	err = testDB.DB.QueryRow("SELECT COUNT(*) FROM product_variants WHERE listing_id = $1", productID).Scan(&variantCount)
	require.NoError(t, err)
	assert.Equal(t, 0, variantCount, "All variants should be deleted")

	t.Log("  ✓ Verified: All variants deleted from database")
	t.Log("========== VARIANT WORKFLOW COMPLETED ==========")
}
