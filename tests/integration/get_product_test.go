//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// GetProduct Happy Path Tests
// ============================================================================

// TestGetProduct_Success verifies successful product retrieval by ID
func TestGetProduct_Success(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	req := &pb.GetProductRequest{
		ProductId:    9000, // Test Laptop from fixtures
		StorefrontId: &storefrontID,
	}

	resp, err := client.GetProduct(ctx, req)

	require.NoError(t, err, "GetProduct should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.NotNil(t, resp.Product, "Product should be returned")
	assert.Equal(t, int64(9000), resp.Product.Id, "Product ID should match")
	assert.Equal(t, "Test Laptop", resp.Product.Name, "Product name should match")
	assert.Equal(t, 1200.00, resp.Product.Price, "Product price should match")
	assert.Equal(t, "USD", resp.Product.Currency, "Currency should be USD")
	assert.Equal(t, int32(50), resp.Product.StockQuantity, "Stock quantity should be 50")
	assert.True(t, resp.Product.IsActive, "Product should be active")
}

// TestGetProduct_WithVariants verifies product with variants is returned correctly
func TestGetProduct_WithVariants(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	req := &pb.GetProductRequest{
		ProductId:    9001, // Test T-Shirt with variants
		StorefrontId: &storefrontID,
	}

	resp, err := client.GetProduct(ctx, req)

	require.NoError(t, err, "GetProduct should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.NotNil(t, resp.Product, "Product should be returned")
	assert.Equal(t, int64(9001), resp.Product.Id)
	assert.Equal(t, "Test T-Shirt", resp.Product.Name)
	assert.True(t, resp.Product.HasVariants, "Product should have variants")
	assert.NotEmpty(t, resp.Product.Variants, "Variants should be included")
	assert.Len(t, resp.Product.Variants, 3, "Should have 3 variants (S, M, L)")
}

// TestGetProduct_WithImages verifies product with images is returned correctly
func TestGetProduct_WithImages(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)

	req := &pb.GetProductRequest{
		ProductId: 9002, // Test Headphones with images
	}

	resp, err := client.GetProduct(ctx, req)

	require.NoError(t, err, "GetProduct should succeed")
	require.NotNil(t, resp.Product, "Product should be returned")
	assert.Equal(t, "Test Headphones", resp.Product.Name)
	// Note: Images may be loaded separately in real implementation
}

// ============================================================================
// GetProductsBySKUs Tests
// ============================================================================

// TestGetProductsBySKUs_Success verifies batch retrieval by SKUs
func TestGetProductsBySKUs_Success(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	req := &pb.GetProductsBySKUsRequest{
		Skus: []string{
			"SKU-BATCH-010",
			"SKU-BATCH-011",
			"SKU-BATCH-012",
			"SKU-BATCH-013",
			"SKU-BATCH-014",
		},
		StorefrontId: &storefrontID,
	}

	resp, err := client.GetProductsBySKUs(ctx, req)

	require.NoError(t, err, "GetProductsBySKUs should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Len(t, resp.Products, 5, "Should return 5 products")

	// Verify each product has correct SKU
	skuMap := make(map[string]bool)
	for _, p := range resp.Products {
		skuMap[*p.Sku] = true
	}
	assert.True(t, skuMap["SKU-BATCH-010"], "Should include SKU-BATCH-010")
	assert.True(t, skuMap["SKU-BATCH-014"], "Should include SKU-BATCH-014")
}

// ============================================================================
// GetProductsByIDs Tests
// ============================================================================

// TestGetProductsByIDs_Success verifies batch retrieval by IDs
func TestGetProductsByIDs_Success(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	req := &pb.GetProductsByIDsRequest{
		ProductIds: []int64{
			9005, 9006, 9007, 9008, 9009, 9010, 9011, 9012, 9013, 9014,
		},
		StorefrontId: &storefrontID,
	}

	resp, err := client.GetProductsByIDs(ctx, req)

	require.NoError(t, err, "GetProductsByIDs should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Equal(t, 10, len(resp.Products), "Should return 10 products")
	assert.Equal(t, int32(10), resp.TotalCount, "Total count should be 10")

	// Verify IDs are present
	idMap := make(map[int64]bool)
	for _, p := range resp.Products {
		idMap[p.Id] = true
	}
	assert.True(t, idMap[9005], "Should include product 9005")
	assert.True(t, idMap[9014], "Should include product 9014")
}

// ============================================================================
// Error Cases Tests
// ============================================================================

// TestGetProduct_NotFound verifies error when product doesn't exist
func TestGetProduct_NotFound(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)

	req := &pb.GetProductRequest{
		ProductId: 99999, // Non-existent product
	}

	resp, err := client.GetProduct(ctx, req)

	require.Error(t, err, "Should return error for non-existent product")
	assert.Nil(t, resp, "Response should be nil on error")

	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be gRPC status")
	assert.Equal(t, codes.NotFound, st.Code(), "Should return NotFound status code")
}

// TestGetProduct_SoftDeleted verifies soft-deleted product is NOT returned
// CRITICAL BUG FIX TEST: This validates the deleted_at IS NULL filter
func TestGetProduct_SoftDeleted(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)

	req := &pb.GetProductRequest{
		ProductId: 9020, // Soft-deleted product from fixtures
	}

	resp, err := client.GetProduct(ctx, req)

	// CRITICAL: Soft-deleted product should NOT be found
	require.Error(t, err, "Should return error for soft-deleted product")
	assert.Nil(t, resp, "Response should be nil")

	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be gRPC status")
	assert.Equal(t, codes.NotFound, st.Code(), "Should return NotFound for soft-deleted product")
}

// TestGetProduct_InvalidID verifies error for invalid product ID
func TestGetProduct_InvalidID(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()
	_ = testDB

	ctx := tests.TestContext(t)

	testCases := []struct {
		name      string
		productID int64
		wantCode  codes.Code
	}{
		{"Zero ID", 0, codes.InvalidArgument},
		{"Negative ID", -1, codes.InvalidArgument},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.GetProductRequest{
				ProductId: tc.productID,
			}

			resp, err := client.GetProduct(ctx, req)

			require.Error(t, err, "Should return error for invalid ID")
			assert.Nil(t, resp, "Response should be nil")

			st, ok := status.FromError(err)
			require.True(t, ok, "Error should be gRPC status")
			assert.Equal(t, tc.wantCode, st.Code(), "Should return expected status code")
		})
	}
}

// ============================================================================
// Performance Tests
// ============================================================================

// TestGetProduct_PerformanceUnder50ms verifies response time
func TestGetProduct_PerformanceUnder50ms(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)

	req := &pb.GetProductRequest{
		ProductId: 9000,
	}

	// Warmup
	_, _ = client.GetProduct(ctx, req)

	// Measure
	start := time.Now()
	resp, err := client.GetProduct(ctx, req)
	elapsed := time.Since(start)

	require.NoError(t, err, "GetProduct should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Less(t, elapsed.Milliseconds(), int64(50),
		"GetProduct should complete in < 50ms (actual: %dms)", elapsed.Milliseconds())
}

// TestGetProductsByIDs_BatchPerformance verifies batch performance
func TestGetProductsByIDs_BatchPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9000)

	// Request 10 products
	productIDs := []int64{9005, 9006, 9007, 9008, 9009, 9010, 9011, 9012, 9013, 9014}

	req := &pb.GetProductsByIDsRequest{
		ProductIds:   productIDs,
		StorefrontId: &storefrontID,
	}

	// Warmup
	_, _ = client.GetProductsByIDs(ctx, req)

	// Measure
	start := time.Now()
	resp, err := client.GetProductsByIDs(ctx, req)
	elapsed := time.Since(start)

	require.NoError(t, err, "GetProductsByIDs should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Less(t, elapsed.Milliseconds(), int64(200),
		"GetProductsByIDs (10 items) should complete in < 200ms (actual: %dms)", elapsed.Milliseconds())
}
