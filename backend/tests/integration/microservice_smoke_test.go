// Package integration
// backend/tests/integration/microservice_smoke_test.go
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	listingsv1 "github.com/sveturs/listings/api/proto/listings/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	microserviceAddr = "localhost:50053"
	healthTimeout    = 5 * time.Second
	requestTimeout   = 2 * time.Second
)

// TestMicroserviceHealthCheck проверяет что микросервис жив и отвечает
func TestMicroserviceHealthCheck(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()

	conn, err := grpc.NewClient(microserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	client := listingsv1.NewListingsServiceClient(conn)

	// Simple health check - list with limit 1
	resp, err := client.ListListings(ctx, &listingsv1.ListListingsRequest{
		Limit:  1,
		Offset: 0,
	})

	require.NoError(t, err, "Health check should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.GreaterOrEqual(t, resp.Total, int32(0), "Total should be >= 0")

	t.Logf("✅ Microservice health check passed: %d total listings", resp.Total)
}

// TestMicroserviceConnectivity проверяет gRPC connectivity
func TestMicroserviceConnectivity(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()

	conn, err := grpc.NewClient(microserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	client := listingsv1.NewListingsServiceClient(conn)

	// Test ListListings
	resp, err := client.ListListings(ctx, &listingsv1.ListListingsRequest{
		Limit:  10,
		Offset: 0,
	})

	require.NoError(t, err, "ListListings should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	assert.GreaterOrEqual(t, resp.Total, int32(0), "Total should be >= 0")

	t.Logf("✅ gRPC connectivity test passed: got %d listings (total: %d)", len(resp.Listings), resp.Total)
}

// TestMicroserviceResponseTime проверяет время отклика микросервиса
func TestMicroserviceResponseTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()

	conn, err := grpc.NewClient(microserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	client := listingsv1.NewListingsServiceClient(conn)

	start := time.Now()

	// Test ListListings with timeout
	requestCtx, requestCancel := context.WithTimeout(ctx, requestTimeout)
	defer requestCancel()

	resp, err := client.ListListings(requestCtx, &listingsv1.ListListingsRequest{
		Limit:  10,
		Offset: 0,
	})

	duration := time.Since(start)

	require.NoError(t, err, "ListListings should succeed")
	require.NotNil(t, resp, "Response should not be nil")

	// Response should be fast (< 500ms for local microservice)
	assert.Less(t, duration.Milliseconds(), int64(500), "Response time should be < 500ms")

	t.Logf("✅ Response time test passed: %dms", duration.Milliseconds())
}

// TestMicroserviceTimeout проверяет что timeout срабатывает корректно
func TestMicroserviceTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()

	conn, err := grpc.NewClient(microserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	client := listingsv1.NewListingsServiceClient(conn)

	// Create very short timeout context (1ms - should timeout)
	shortCtx, shortCancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer shortCancel()

	// Wait a bit to ensure timeout triggers
	time.Sleep(2 * time.Millisecond)

	_, err = client.ListListings(shortCtx, &listingsv1.ListListingsRequest{
		Limit:  10,
		Offset: 0,
	})

	// Should get context deadline exceeded error
	assert.Error(t, err, "Should get timeout error")

	t.Logf("✅ Timeout test passed: got expected error: %v", err)
}

// TestMicroserviceGetListing проверяет получение конкретного listing
func TestMicroserviceGetListing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()

	conn, err := grpc.NewClient(microserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	client := listingsv1.NewListingsServiceClient(conn)

	// First, get list to find an existing ID
	listResp, err := client.ListListings(ctx, &listingsv1.ListListingsRequest{
		Limit:  1,
		Offset: 0,
	})
	require.NoError(t, err, "ListListings should succeed")

	if len(listResp.Listings) == 0 {
		t.Skip("No listings in database, skipping GetListing test")
	}

	// Get first listing
	listingID := listResp.Listings[0].Id

	getResp, err := client.GetListing(ctx, &listingsv1.GetListingRequest{
		Id: listingID,
	})

	require.NoError(t, err, "GetListing should succeed")
	require.NotNil(t, getResp, "Response should not be nil")
	require.NotNil(t, getResp.Listing, "Listing should not be nil")
	assert.Equal(t, listingID, getResp.Listing.Id, "Listing ID should match")

	t.Logf("✅ GetListing test passed: got listing %d (%s)", getResp.Listing.Id, getResp.Listing.Title)
}

// TestMicroserviceSearchListings проверяет поиск listing
func TestMicroserviceSearchListings(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), healthTimeout)
	defer cancel()

	conn, err := grpc.NewClient(microserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	client := listingsv1.NewListingsServiceClient(conn)

	// Test search with simple query
	searchResp, err := client.SearchListings(ctx, &listingsv1.SearchListingsRequest{
		Query:  "test",
		Limit:  10,
		Offset: 0,
	})

	require.NoError(t, err, "SearchListings should succeed")
	require.NotNil(t, searchResp, "Response should not be nil")
	assert.GreaterOrEqual(t, searchResp.Total, int32(0), "Total should be >= 0")

	t.Logf("✅ SearchListings test passed: found %d results", searchResp.Total)
}
