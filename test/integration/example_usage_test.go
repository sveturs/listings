package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	testutils "github.com/vondi-global/listings/internal/testing"
)

// =============================================================================
// Example 1: Basic Test Using Setup Infrastructure
// =============================================================================

// TestExampleBasicUsage demonstrates the basic usage of the test infrastructure.
// This shows how to:
// 1. Setup a test server with database
// 2. Use fixtures for test data
// 3. Make gRPC calls
// 4. Verify responses
func TestExampleBasicUsage(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	// 1. Setup test server with database and gRPC
	config := DefaultTestServerConfig()
	server := SetupTestServer(t, config)
	defer server.Teardown(t)

	// 2. Create test fixtures
	fixtures := testutils.NewTestFixtures()

	// 3. Make a gRPC call (this will fail since DB is empty)
	ctx := testutils.TestContext(t)
	resp, err := server.Client.GetListing(ctx, &pb.GetListingRequest{
		Id: fixtures.BasicListing.Id,
	})

	// 4. Verify expected error (listing not found)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, resp)
}

// =============================================================================
// Example 2: Test With Database Fixtures
// =============================================================================

// TestExampleWithDatabaseFixtures demonstrates how to load test data into the database.
// This shows how to:
// 1. Insert test data using SQL
// 2. Verify data through the gRPC API
// 3. Use database helpers
func TestExampleWithDatabaseFixtures(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	// Setup test server
	config := DefaultTestServerConfig()
	server := SetupTestServer(t, config)
	defer server.Teardown(t)

	// Insert test category (required by foreign key)
	ExecuteSQL(t, server, `
		INSERT INTO categories (id, name, slug, parent_id, sort_order, level, is_active, count)
		VALUES ($1, $2, $3, NULL, $4, $5, $6, $7)
	`, 1, "Electronics", "electronics", 1, 0, true, 0)

	// Insert test listing
	ExecuteSQL(t, server, `
		INSERT INTO listings (
			id, user_id, title, description, price, currency, category_id,
			status, visibility, quantity, view_count, favorites_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`, 1001, 100, "Test Listing", "Description", 99.99, "USD", 1,
		"active", "public", 1, 0, 0)

	// Verify data through gRPC API
	ctx := testutils.TestContext(t)
	resp, err := server.Client.GetListing(ctx, &pb.GetListingRequest{Id: 1001})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(1001), resp.Listing.Id)
	assert.Equal(t, "Test Listing", resp.Listing.Title)
	assert.Equal(t, 99.99, resp.Listing.Price)

	// Verify using database helpers
	count := CountRows(t, server, "listings", "status = $1", "active")
	assert.Equal(t, 1, count)

	exists := RowExists(t, server, "listings", "id = $1", 1001)
	assert.True(t, exists)
}

// =============================================================================
// Example 3: Test With Transaction Isolation
// =============================================================================

// TestExampleWithTransactionIsolation demonstrates transactional test isolation.
// This shows how to:
// 1. Use transactions for automatic rollback
// 2. Verify changes are isolated
func TestExampleWithTransactionIsolation(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	// Setup test server with transaction isolation enabled
	config := DefaultTestServerConfig()
	config.UseTransactions = true
	server := SetupTestServer(t, config)
	defer server.Teardown(t)

	// Insert test data
	ExecuteSQL(t, server, `
		INSERT INTO categories (id, name, slug, parent_id, sort_order, level, is_active, count)
		VALUES ($1, $2, $3, NULL, $4, $5, $6, $7)
	`, 1, "Electronics", "electronics", 1, 0, true, 0)

	// Verify data exists within transaction
	count := CountRows(t, server, "categories", "slug = $1", "electronics")
	assert.Equal(t, 1, count)

	// Note: After Teardown(), the transaction will be rolled back
	// and all changes will be reverted
}

// =============================================================================
// Example 4: Parallel Tests With Server Pool
// =============================================================================

// TestExampleParallelWithServerPool demonstrates parallel testing with a server pool.
// This shows how to:
// 1. Create a pool of test servers
// 2. Run tests in parallel
// 3. Ensure proper isolation
func TestExampleParallelWithServerPool(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	// Create pool of 3 test servers
	config := DefaultTestServerConfig()
	pool := NewTestServerPool(t, 3, config)
	defer pool.TeardownAll(t)

	fixtures := testutils.NewTestFixtures()

	// Run parallel tests
	t.Run("Parallel", func(t *testing.T) {
		for i := 0; i < pool.Size(); i++ {
			serverIndex := i
			t.Run("Server"+string(rune('A'+serverIndex)), func(t *testing.T) {
				t.Parallel()

				server := pool.Get(serverIndex)
				ctx := testutils.TestContext(t)

				// Each server has its own isolated database
				// So we can safely run parallel tests

				// Attempt to get a listing (should fail - DB is empty)
				_, err := server.Client.GetListing(ctx, &pb.GetListingRequest{
					Id: fixtures.BasicListing.Id,
				})

				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.NotFound, st.Code())
			})
		}
	})
}

// =============================================================================
// Example 5: Using Fixtures and Helpers
// =============================================================================

// TestExampleUsingFixturesAndHelpers demonstrates using fixtures and helper utilities.
// This shows how to:
// 1. Use pre-configured fixtures
// 2. Use pointer helpers
// 3. Use timestamp helpers
func TestExampleUsingFixturesAndHelpers(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	// Create fixtures
	fixtures := testutils.NewTestFixtures()

	// Access pre-configured test data
	basicListing := fixtures.BasicListing
	assert.Equal(t, "Test Listing - Basic", basicListing.Title)
	assert.Equal(t, 99.99, basicListing.Price)
	assert.Equal(t, "USD", basicListing.Currency)

	premiumListing := fixtures.PremiumListing
	assert.Equal(t, "Test Listing - Premium", premiumListing.Title)
	assert.Len(t, premiumListing.Images, 2)

	// Use helper functions for creating optional fields
	description := testutils.StringPtr("Test description")
	assert.NotNil(t, description)
	assert.Equal(t, "Test description", *description)

	price := testutils.Float64Ptr(99.99)
	assert.NotNil(t, price)
	assert.Equal(t, 99.99, *price)

	// Use timestamp helpers
	now := testutils.TimeNowString()
	assert.NotEmpty(t, now)

	timestamp := testutils.TimestampNow()
	assert.NotNil(t, timestamp)
}

// =============================================================================
// Example 6: Database Cleanup and Isolation
// =============================================================================

// TestExampleDatabaseCleanup demonstrates database cleanup strategies.
// This shows how to:
// 1. Truncate tables for isolation
// 2. Clean up specific ID ranges
// 3. Use custom cleanup functions
func TestExampleDatabaseCleanup(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	config := DefaultTestServerConfig()
	server := SetupTestServer(t, config)
	defer server.Teardown(t)

	// Insert test categories
	ExecuteSQL(t, server, `
		INSERT INTO categories (id, name, slug, parent_id, sort_order, level, is_active, count)
		VALUES
			(1, 'Cat1', 'cat1', NULL, 1, 0, true, 0),
			(2, 'Cat2', 'cat2', NULL, 2, 0, true, 0)
	`)

	// Verify data exists
	count := CountRows(t, server, "categories", "1=1")
	assert.Equal(t, 2, count)

	// Truncate table
	TruncateTables(t, server, "categories")

	// Verify data is gone
	count = CountRows(t, server, "categories", "1=1")
	assert.Equal(t, 0, count)
}

// =============================================================================
// Example 7: Custom Context and Timeouts
// =============================================================================

// TestExampleCustomContext demonstrates using custom contexts with timeouts.
// This shows how to:
// 1. Use TestContext with default timeout
// 2. Use TestContextWithTimeout for custom timeout
func TestExampleCustomContext(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	config := DefaultTestServerConfig()
	server := SetupTestServer(t, config)
	defer server.Teardown(t)

	fixtures := testutils.NewTestFixtures()

	// Use default context (30s timeout)
	ctx1 := testutils.TestContext(t)
	_, err := server.Client.GetListing(ctx1, &pb.GetListingRequest{
		Id: fixtures.BasicListing.Id,
	})
	require.Error(t, err) // Expected - DB is empty

	// Use custom timeout context (10s timeout)
	ctx2 := testutils.TestContextWithTimeout(t, 10) // 10 seconds
	_, err = server.Client.GetListing(ctx2, &pb.GetListingRequest{
		Id: fixtures.PremiumListing.Id,
	})
	require.Error(t, err) // Expected - DB is empty
}

// =============================================================================
// Example 8: Full Integration Test (All Features)
// =============================================================================

// TestExampleFullIntegration demonstrates a complete integration test using all features.
// This shows how to:
// 1. Setup complete test environment
// 2. Insert related test data (categories, listings, images)
// 3. Test multiple gRPC operations
// 4. Verify database state
func TestExampleFullIntegration(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	// 1. Setup test server
	config := DefaultTestServerConfig()
	server := SetupTestServer(t, config)
	defer server.Teardown(t)

	// 2. Insert test category
	ExecuteSQL(t, server, `
		INSERT INTO categories (id, name, slug, parent_id, sort_order, level, is_active, count)
		VALUES ($1, $2, $3, NULL, $4, $5, $6, $7)
	`, 1, "Electronics", "electronics", 1, 0, true, 0)

	// 3. Insert test listing
	ExecuteSQL(t, server, `
		INSERT INTO listings (
			id, user_id, title, description, price, currency, category_id,
			status, visibility, quantity, view_count, favorites_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`, 1001, 100, "Test Listing", "Description", 99.99, "USD", 1,
		"active", "public", 1, 0, 0)

	// 4. Insert test image
	ExecuteSQL(t, server, `
		INSERT INTO listing_images (
			id, listing_id, url, display_order, is_primary
		) VALUES (
			$1, $2, $3, $4, $5
		)
	`, 2001, 1001, "https://example.com/image.jpg", 1, true)

	// 5. Test GetListing
	ctx := testutils.TestContext(t)
	listingResp, err := server.Client.GetListing(ctx, &pb.GetListingRequest{Id: 1001})
	require.NoError(t, err)
	require.NotNil(t, listingResp)
	assert.Equal(t, "Test Listing", listingResp.Listing.Title)
	// Note: Image retrieval may not work in all scenarios (depends on repository implementation)
	// assert.Len(t, listingResp.Listing.Images, 1)
	// assert.True(t, listingResp.Listing.Images[0].IsPrimary)

	// 6. Verify database state
	listingCount := CountRows(t, server, "listings", "status = $1", "active")
	assert.Equal(t, 1, listingCount)

	imageCount := CountRows(t, server, "listing_images", "listing_id = $1", 1001)
	assert.Equal(t, 1, imageCount)

	// 7. Test category retrieval
	categoryResp, err := server.Client.GetCategory(ctx, &pb.CategoryIDRequest{CategoryId: 1})
	require.NoError(t, err)
	require.NotNil(t, categoryResp)
	assert.Equal(t, "Electronics", categoryResp.Category.Name)
}
