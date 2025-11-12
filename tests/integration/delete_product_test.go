//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/tests"
)

// ============================================================================
// DeleteProduct Happy Path Tests
// ============================================================================

// TestDeleteProduct_Success verifies successful hard delete
func TestDeleteProduct_Success(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	req := &pb.DeleteProductRequest{
		ProductId:    9100, // Delete Test Product 1
		StorefrontId: storefrontID,
		HardDelete:   true,
	}

	resp, err := client.DeleteProduct(ctx, req)

	require.NoError(t, err, "DeleteProduct should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.Success, "Deletion should be successful")

	// Verify product is gone
	getReq := &pb.GetProductRequest{ProductId: 9100}
	_, err = client.GetProduct(ctx, getReq)
	assert.Error(t, err, "Product should not exist after hard delete")
}

// TestDeleteProduct_SoftDelete verifies soft delete sets deleted_at
func TestDeleteProduct_SoftDelete(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	req := &pb.DeleteProductRequest{
		ProductId:    9101, // Delete Test Product 2
		StorefrontId: storefrontID,
		HardDelete:   false, // Soft delete
	}

	resp, err := client.DeleteProduct(ctx, req)

	require.NoError(t, err, "Soft delete should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.True(t, resp.Success, "Soft deletion should be successful")

	// CRITICAL: GetProduct should NOT find soft-deleted product (bug fix validation)
	getReq := &pb.GetProductRequest{ProductId: 9101}
	_, err = client.GetProduct(ctx, getReq)
	require.Error(t, err, "GetProduct should return error for soft-deleted product")

	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be gRPC status")
	assert.Equal(t, codes.NotFound, st.Code(), "Should return NotFound for soft-deleted product")

	// Verify deleted_at is set in database
	var deletedAt *string
	err = testDB.DB.QueryRow("SELECT deleted_at FROM listings WHERE id = $1 AND source_type = 'b2c'", 9101).Scan(&deletedAt)
	require.NoError(t, err, "Should query deleted_at")
	assert.NotNil(t, deletedAt, "deleted_at should be set")
	assert.NotEmpty(t, *deletedAt, "deleted_at should not be empty")
}

// TestDeleteProduct_WithVariants verifies CASCADE delete of variants
func TestDeleteProduct_WithVariants(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	// Product 9102 has 5 variants
	req := &pb.DeleteProductRequest{
		ProductId:    9102,
		StorefrontId: storefrontID,
		HardDelete:   true,
	}

	resp, err := client.DeleteProduct(ctx, req)

	require.NoError(t, err, "Delete with variants should succeed")
	assert.True(t, resp.Success, "Deletion should be successful")
	assert.Equal(t, int32(5), resp.VariantsDeleted, "Should delete 5 variants")

	// Verify variants are also deleted
	var variantCount int
	err = testDB.DB.QueryRow("SELECT COUNT(*) FROM product_variants WHERE listing_id = $1", 9102).Scan(&variantCount)
	require.NoError(t, err, "Should query variant count")
	assert.Equal(t, 0, variantCount, "All variants should be deleted (CASCADE)")
}

// ============================================================================
// Validation Tests
// ============================================================================

// TestDeleteProduct_NonExistent verifies error when product doesn't exist
func TestDeleteProduct_NonExistent(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	req := &pb.DeleteProductRequest{
		ProductId:    99999, // Non-existent product
		StorefrontId: storefrontID,
		HardDelete:   true,
	}

	resp, err := client.DeleteProduct(ctx, req)

	require.Error(t, err, "Should return error for non-existent product")
	assert.Nil(t, resp, "Response should be nil on error")

	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be gRPC status")
	assert.Equal(t, codes.NotFound, st.Code(), "Should return NotFound status code")
}

// TestDeleteProduct_InvalidID verifies validation of product ID
func TestDeleteProduct_InvalidID(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()
	_ = testDB

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

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
			req := &pb.DeleteProductRequest{
				ProductId:    tc.productID,
				StorefrontId: storefrontID,
				HardDelete:   true,
			}

			resp, err := client.DeleteProduct(ctx, req)

			require.Error(t, err, "Should return error for invalid ID")
			assert.Nil(t, resp, "Response should be nil")

			st, ok := status.FromError(err)
			require.True(t, ok, "Error should be gRPC status")
			assert.Equal(t, tc.wantCode, st.Code(), "Should return expected status code")
		})
	}
}

// TestDeleteProduct_AlreadyDeleted verifies idempotency (soft delete twice)
func TestDeleteProduct_AlreadyDeleted(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	// Product 9103 is already soft-deleted in fixtures
	req := &pb.DeleteProductRequest{
		ProductId:    9103,
		StorefrontId: storefrontID,
		HardDelete:   false,
	}

	resp, err := client.DeleteProduct(ctx, req)

	// Behavior can be either: error OR success with idempotency
	// Implementation dependent - document actual behavior
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok, "Error should be gRPC status")
		assert.Equal(t, codes.NotFound, st.Code(), "Already deleted product returns NotFound")
	} else {
		assert.True(t, resp.Success, "Idempotent delete should succeed")
	}
}

// ============================================================================
// BulkDeleteProducts Tests
// ============================================================================

// TestBulkDeleteProducts_Success verifies batch deletion
func TestBulkDeleteProducts_Success(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	req := &pb.BulkDeleteProductsRequest{
		StorefrontId: storefrontID,
		ProductIds: []int64{
			9104, 9105, 9106, 9107, 9108, // 5 products from fixtures
		},
		HardDelete: true,
	}

	resp, err := client.BulkDeleteProducts(ctx, req)

	require.NoError(t, err, "BulkDeleteProducts should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.Equal(t, int32(5), resp.SuccessfulCount, "Should delete 5 products")
	assert.Equal(t, int32(0), resp.FailedCount, "Should have no failures")
	assert.Empty(t, resp.Errors, "Should have no errors")

	// Verify all products are deleted
	for _, id := range req.ProductIds {
		var count int
		err := testDB.DB.QueryRow("SELECT COUNT(*) FROM listings WHERE id = $1 AND source_type = 'b2c'", id).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count, "Product %d should be deleted", id)
	}
}

// TestBulkDeleteProducts_EmptyBatch verifies validation of empty list
func TestBulkDeleteProducts_EmptyBatch(t *testing.T) {
	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()
	_ = testDB

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	req := &pb.BulkDeleteProductsRequest{
		StorefrontId: storefrontID,
		ProductIds:   []int64{}, // Empty list
		HardDelete:   true,
	}

	resp, err := client.BulkDeleteProducts(ctx, req)

	// Should either error OR return success with 0 count
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok, "Error should be gRPC status")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "Empty batch should return InvalidArgument")
	} else {
		assert.Equal(t, int32(0), resp.SuccessfulCount, "Empty batch has 0 successes")
	}
}

// TestBulkDeleteProducts_LargeBatch verifies performance with 100 products
func TestBulkDeleteProducts_LargeBatch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large batch test in short mode")
	}

	client, testDB, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	tests.LoadTestFixtures(t, testDB.DB, "../fixtures/get_delete_product_fixtures.sql")

	ctx := tests.TestContext(t)
	storefrontID := int64(9001)

	// Delete 100 products (9150-9249 from fixtures)
	productIDs := make([]int64, 100)
	for i := 0; i < 100; i++ {
		productIDs[i] = int64(9150 + i)
	}

	req := &pb.BulkDeleteProductsRequest{
		StorefrontId: storefrontID,
		ProductIds:   productIDs,
		HardDelete:   true,
	}

	resp, err := client.BulkDeleteProducts(ctx, req)

	require.NoError(t, err, "Large batch delete should succeed")
	assert.Equal(t, int32(100), resp.SuccessfulCount, "Should delete all 100 products")
	assert.Equal(t, int32(0), resp.FailedCount, "Should have no failures")

	// Verify products are deleted
	var remainingCount int
	err = testDB.DB.QueryRow("SELECT COUNT(*) FROM listings WHERE id BETWEEN 9150 AND 9249 AND source_type = 'b2c'").Scan(&remainingCount)
	require.NoError(t, err)
	assert.Equal(t, 0, remainingCount, "All 100 products should be deleted")
}
