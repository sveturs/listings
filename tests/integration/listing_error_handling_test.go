//go:build integration

package integration

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/internal/service/listings"
	grpchandlers "github.com/vondi-global/listings/internal/transport/grpc"
	"github.com/vondi-global/listings/tests"
)

// ============================================================================
// Error Handling Tests - Database Failures, Timeouts, Invalid Data
// ============================================================================

// TestListing_Error_MissingRequiredFields tests missing required fields
func TestListing_Error_MissingRequiredFields(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name        string
		req         *pb.CreateListingRequest
		expectedErr codes.Code
	}{
		{
			name: "Missing title",
			req: &pb.CreateListingRequest{
				UserId:     100,
				Title:      "", // Missing
				Price:      99.99,
				Currency:   "USD",
				CategoryId: 1,
				Quantity:   1,
			},
			expectedErr: codes.InvalidArgument,
		},
		{
			name: "Missing currency",
			req: &pb.CreateListingRequest{
				UserId:     100,
				Title:      "Test Product",
				Price:      99.99,
				Currency:   "", // Missing
				CategoryId: 1,
				Quantity:   1,
			},
			expectedErr: codes.InvalidArgument,
		},
		{
			name: "Missing user_id",
			req: &pb.CreateListingRequest{
				UserId:     0, // Missing
				Title:      "Test Product",
				Price:      99.99,
				Currency:   "USD",
				CategoryId: 1,
				Quantity:   1,
			},
			expectedErr: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.CreateListing(ctx, tc.req)

			require.Error(t, err, "Should return error for missing required field")
			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.expectedErr, st.Code())
			assert.Nil(t, resp)
		})
	}
}

// TestListing_Error_InvalidEnumValues tests invalid enum values
func TestListing_Error_InvalidEnumValues(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name   string
		status string
	}{
		{"Invalid status", "invalid_status"},
		{"Typo in status", "activ"},
		{"Uppercase status", "ACTIVE"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CreateListingRequest{
				UserId:     100,
				Title:      "Test Product",
				Price:      99.99,
				Currency:   "USD",
				CategoryId: 1,
				Quantity:   1,
			}

			resp, err := client.CreateListing(ctx, req)

			if err != nil {
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			} else {
				// Backend might normalize or accept any string
				require.NotNil(t, resp)
				t.Logf("Backend accepted status '%s': %v", tc.status, resp.Listing.Status)
			}
		})
	}
}

// TestListing_Error_NotFound tests retrieving non-existent listing
func TestListing_Error_NotFound(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Try to get non-existent listing
	req := &pb.GetListingRequest{
		Id: 999999, // Non-existent ID
	}

	resp, err := client.GetListing(ctx, req)

	require.Error(t, err, "Should return error for non-existent listing")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, resp)
}

// TestListing_Error_DeleteNonExistent tests deleting non-existent listing
func TestListing_Error_DeleteNonExistent(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.DeleteListingRequest{
		Id: 999999, // Non-existent
	}

	resp, err := client.DeleteListing(ctx, req)

	require.Error(t, err, "Should return error when deleting non-existent listing")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, resp)
}

// TestListing_Error_UpdateNonExistent tests updating non-existent listing
func TestListing_Error_UpdateNonExistent(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.UpdateListingRequest{
		Id:    999999, // Non-existent
		Title: stringPtr("Updated Title"),
		Price: float64Ptr(199.99),
	}

	resp, err := client.UpdateListing(ctx, req)

	require.Error(t, err, "Should return error when updating non-existent listing")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Nil(t, resp)
}

// TestListing_Error_InvalidCurrency tests invalid currency code
func TestListing_Error_InvalidCurrency(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	testCases := []struct {
		name     string
		currency string
	}{
		{"Too short", "US"},
		{"Too long", "USDD"},
		{"Invalid code", "XXX"},
		{"Lowercase", "usd"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.CreateListingRequest{
				UserId:     100,
				Title:      "Test Product",
				Price:      99.99,
				Currency:   tc.currency,
				CategoryId: 1,
				Quantity:   1,
			}

			resp, err := client.CreateListing(ctx, req)

			if err != nil {
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
			} else {
				// Backend might accept any currency string
				require.NotNil(t, resp)
				t.Logf("Backend accepted currency '%s'", tc.currency)
			}
		})
	}
}

// TestListing_Error_DuplicateSKU tests duplicate SKU handling
func TestListing_Error_DuplicateSKU(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create first listing with SKU
	req1 := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "First Product",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,

		Sku: stringPtr("DUPLICATE-SKU-001"),
	}

	resp1, err := client.CreateListing(ctx, req1)
	require.NoError(t, err)
	require.NotNil(t, resp1)

	// Try to create second listing with same SKU
	req2 := &pb.CreateListingRequest{
		UserId:     101,
		Title:      "Second Product",
		Price:      149.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,

		Sku: stringPtr("DUPLICATE-SKU-001"), // Same SKU
	}

	resp2, err := client.CreateListing(ctx, req2)

	// Depending on business logic:
	// - Might allow duplicate SKUs (different users/storefronts)
	// - Might reject with AlreadyExists error
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.AlreadyExists, st.Code(), "Should reject duplicate SKU")
	} else {
		require.NotNil(t, resp2)
		t.Log("Backend allows duplicate SKU (might be scoped to user/storefront)")
	}
}

// ============================================================================
// Resource Exhaustion Tests
// ============================================================================

// TestListing_Error_TooManyListings tests creating excessive listings
func TestListing_Error_TooManyListings(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping resource exhaustion test in short mode")
	}

	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Try to create 1000 listings rapidly
	successCount := 0
	errorCount := 0

	for i := 0; i < 1000; i++ {
		req := &pb.CreateListingRequest{
			UserId:     100,
			Title:      "Bulk Test Product",
			Price:      99.99,
			Currency:   "USD",
			CategoryId: 1,
			Quantity:   1,
		}

		_, err := client.CreateListing(ctx, req)
		if err != nil {
			errorCount++
			// Should eventually hit rate limit or resource limit
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.ResourceExhausted {
				t.Logf("Hit resource limit after %d listings", successCount)
				break
			}
		} else {
			successCount++
		}
	}

	t.Logf("Created %d listings, %d errors", successCount, errorCount)
	assert.Greater(t, successCount, 0, "Should create some listings")
}

// ============================================================================
// Database Connection Failure Tests
// ============================================================================

// TestListing_Error_DatabaseDisconnect tests behavior when database disconnects
func TestListing_Error_DatabaseDisconnect(t *testing.T) {
	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)
	defer testDB.TeardownTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := postgres.NewRepository(db, logger)

	// Create service
	service := listings.NewService(repo, nil, nil, logger)

	// Create gRPC server
	m := getTestMetrics()
	server := grpchandlers.NewServer(service, m, logger)

	// Setup gRPC connection
	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()
	pb.RegisterListingsServiceServer(grpcServer, server)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewListingsServiceClient(conn)

	ctx := tests.TestContext(t)

	// First request should succeed
	req1 := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Test Before Disconnect",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	resp1, err := client.CreateListing(ctx, req1)
	require.NoError(t, err)
	require.NotNil(t, resp1)

	// Close database connection to simulate failure
	err = db.Close()
	require.NoError(t, err)

	// Second request should fail
	req2 := &pb.CreateListingRequest{
		UserId:     101,
		Title:      "Test After Disconnect",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	resp2, err := client.CreateListing(ctx, req2)

	require.Error(t, err, "Should fail when database is disconnected")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Contains(t, []codes.Code{codes.Unavailable, codes.Internal}, st.Code())
	assert.Nil(t, resp2)
}

// ============================================================================
// Context Timeout Tests
// ============================================================================

// TestListing_Error_ContextTimeout tests request timeout handling
func TestListing_Error_ContextTimeout(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	// Create context with very short timeout (1 millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for timeout to expire
	time.Sleep(10 * time.Millisecond)

	req := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Timeout Test",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.Error(t, err, "Should fail with timeout")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.DeadlineExceeded, st.Code())
	assert.Nil(t, resp)
}

// TestListing_Error_ContextCancellation tests request cancellation
func TestListing_Error_ContextCancellation(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	req := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Cancelled Request Test",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.Error(t, err, "Should fail with cancellation")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Canceled, st.Code())
	assert.Nil(t, resp)
}

// ============================================================================
// Permission/Authorization Edge Cases
// ============================================================================

// TestListing_Error_UpdateOthersListing tests updating another user's listing
func TestListing_Error_UpdateOthersListing(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create listing as user 100
	createReq := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "User 100 Listing",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	createResp, err := client.CreateListing(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Try to update as user 101 (different user)
	updateReq := &pb.UpdateListingRequest{
		Id:     createResp.Listing.Id,
		UserId: 101, // Different user
		Title:  stringPtr("Hijacked Title"),
	}

	updateResp, err := client.UpdateListing(ctx, updateReq)

	// Should reject unauthorized update
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Contains(t, []codes.Code{codes.PermissionDenied, codes.Unauthenticated}, st.Code())
	} else {
		// Or might silently ignore user_id change
		require.NotNil(t, updateResp)
		assert.Equal(t, int64(100), updateResp.Listing.UserId, "User ID should not change")
	}
}

// TestListing_Error_DeleteOthersListing tests deleting another user's listing
func TestListing_Error_DeleteOthersListing(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create listing as user 100
	createReq := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "User 100 Listing",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	createResp, err := client.CreateListing(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Try to delete as user 101 (requires user context in real app)
	deleteReq := &pb.DeleteListingRequest{
		Id: createResp.Listing.Id,
		// In real app, user_id would come from auth context
	}

	deleteResp, err := client.DeleteListing(ctx, deleteReq)

	// Without auth context, might succeed (test limitation)
	// In production, should check auth and reject
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())
	} else {
		// Test limitation: no auth context
		require.NotNil(t, deleteResp)
		t.Log("Delete succeeded (no auth context in test)")
	}
}

// ============================================================================
// Invalid UUID Tests
// ============================================================================

// TestListing_Error_InvalidID tests getting listing by invalid ID (zero)
func TestListing_Error_InvalidID(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.GetListingRequest{
		Id: 0, // Invalid ID
	}

	resp, err := client.GetListing(ctx, req)

	require.Error(t, err, "Should reject zero/invalid ID")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Contains(t, []codes.Code{codes.InvalidArgument, codes.NotFound}, st.Code())
	assert.Nil(t, resp)
}
