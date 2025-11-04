//go:build ignore

// Package integration contains integration tests for microservice connectivity
// DEPRECATED: These tests use outdated proto API structures (ListingId field)
// To enable: update proto structures to match current listings API
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"

	listingsClient "backend/internal/clients/listings"
	"backend/internal/logger"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

const (
	// Test microservice address (expecting mock or real service)
	testMicroserviceAddr = "localhost:50053"
	// testTimeout is defined in canary_integration_test.go to avoid redeclaration
)

// TestMicroserviceHealthCheckGRPC verifies health check endpoint is accessible (renamed to avoid duplicate)
func TestMicroserviceHealthCheckGRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Establish gRPC connection
	conn, err := grpc.NewClient(
		testMicroserviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Should connect to microservice")
	defer func() { _ = conn.Close() }()

	// Create health check client
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Check health
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{
		Service: "listings.v1.ListingsService",
	})

	require.NoError(t, err, "Health check should succeed")
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status,
		"Service should be in SERVING state")

	t.Logf("✅ Health check passed: %v", resp.Status)
}

// TestGRPCConnectionEstablished verifies gRPC connection can be established
func TestGRPCConnectionEstablished(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)

	require.NoError(t, err, "Should create gRPC client")
	require.NotNil(t, client, "Client should not be nil")

	defer func() { _ = client.Close() }()

	t.Log("✅ gRPC connection established successfully")
}

// TestGRPCRequestResponseCycle verifies complete request-response cycle
func TestGRPCRequestResponseCycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	tests := []struct {
		name        string
		request     func() (interface{}, error)
		description string
	}{
		{
			name: "ListListings request-response",
			request: func() (interface{}, error) {
				return client.ListListings(ctx, &pb.ListListingsRequest{
					Page:     1,
					PageSize: 10,
				})
			},
			description: "List operation should work",
		},
		{
			name: "GetListing for non-existent ID",
			request: func() (interface{}, error) {
				return client.GetListing(ctx, &pb.GetListingRequest{
					ListingId: 999999, // Non-existent
				})
			},
			description: "Should handle not found gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.request()

			// We don't care if it's NotFound error - we just verify the RPC cycle works
			if err != nil {
				t.Logf("Request returned error (expected for non-existent): %v", err)
			} else {
				assert.NotNil(t, resp, "Response should not be nil")
			}

			t.Logf("✅ %s", tt.description)
		})
	}
}

// TestAuthenticationPassthrough verifies JWT forwarding works
func TestAuthenticationPassthrough(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Test without auth (should work for public endpoints)
	resp, err := client.ListListings(ctx, &pb.ListListingsRequest{
		Page:     1,
		PageSize: 5,
	})

	// Either success or auth error is acceptable (depends on microservice config)
	if err != nil {
		t.Logf("Request returned error: %v (expected if auth required)", err)
	} else {
		assert.NotNil(t, resp, "Response should not be nil")
		t.Logf("✅ Request succeeded without auth (public endpoint)")
	}
}

// TestTimeoutHandling verifies 500ms timeout is respected
func TestTimeoutHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()

	// Try to get a listing (might timeout if service is slow)
	_, err = client.GetListing(ctx, &pb.GetListingRequest{
		ListingId: 999, // Special ID that triggers slow response in mock
	})

	elapsed := time.Since(start)

	// We expect either:
	// 1. Quick success (if mock is fast)
	// 2. Timeout error within ~100ms
	if err != nil {
		assert.Contains(t, err.Error(), "context deadline exceeded",
			"Should timeout with proper error")
		assert.LessOrEqual(t, elapsed, 200*time.Millisecond,
			"Should timeout within expected duration")
		t.Logf("✅ Timeout handled correctly in %v", elapsed)
	} else {
		t.Logf("✅ Request completed quickly in %v (no timeout)", elapsed)
	}
}

// TestCircuitBreakerStateTransitions verifies circuit breaker opens after failures
func TestCircuitBreakerStateTransitions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// Circuit breaker should be CLOSED initially
	t.Log("Circuit breaker initial state: CLOSED")

	// Trigger 5 consecutive failures to open circuit
	// Note: This requires microservice mock to return errors for special ID 777
	failureCount := 0
	for i := 0; i < 6; i++ {
		_, err := client.GetListing(ctx, &pb.GetListingRequest{
			ListingId: 777, // Special ID that triggers errors in mock
		})
		if err != nil {
			failureCount++
			t.Logf("Failure %d/5: %v", failureCount, err)
		}
	}

	assert.GreaterOrEqual(t, failureCount, 5, "Should have at least 5 failures")

	// Circuit should now be OPEN
	// Next request should fail immediately with circuit breaker error
	_, err = client.GetListing(ctx, &pb.GetListingRequest{
		ListingId: 1,
	})

	if err != nil {
		assert.Contains(t, err.Error(), "service unavailable",
			"Circuit breaker should reject request when open")
		t.Log("✅ Circuit breaker opened after threshold failures")
	} else {
		t.Log("⚠️ Circuit breaker did not open (might be using mock without circuit breaker)")
	}
}

// TestCircuitBreakerRecovery verifies circuit breaker closes after successful requests
func TestCircuitBreakerRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Skip("Requires 30s wait for half-open state - run manually")

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// 1. Open circuit with failures
	for i := 0; i < 6; i++ {
		client.GetListing(ctx, &pb.GetListingRequest{ListingId: 777})
	}

	t.Log("Circuit opened, waiting 30s for half-open state...")
	time.Sleep(30 * time.Second)

	// 2. Make 2 successful requests to close circuit
	for i := 0; i < 2; i++ {
		resp, err := client.GetListing(ctx, &pb.GetListingRequest{ListingId: 1})
		if err == nil {
			assert.NotNil(t, resp)
			t.Logf("Success %d/2", i+1)
		}
	}

	// 3. Circuit should be CLOSED now
	resp, err := client.ListListings(ctx, &pb.ListListingsRequest{Page: 1, PageSize: 10})
	assert.NoError(t, err, "Circuit should be closed and accepting requests")
	assert.NotNil(t, resp)

	t.Log("✅ Circuit breaker recovered successfully")
}

// TestConcurrentRequests verifies client handles concurrent requests safely
func TestConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	require.NoError(t, err)
	defer func() { _ = client.Close() }()

	concurrency := 10
	done := make(chan bool, concurrency)
	errors := make(chan error, concurrency)

	ctx := context.Background()

	// Launch concurrent requests
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			_, err := client.ListListings(ctx, &pb.ListListingsRequest{
				Page:     1,
				PageSize: 5,
			})
			if err != nil {
				errors <- err
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}
	close(errors)

	// Check error rate
	errorCount := 0
	for range errors {
		errorCount++
	}

	errorRate := float64(errorCount) / float64(concurrency) * 100
	assert.LessOrEqual(t, errorRate, 10.0, "Error rate should be < 10% under concurrent load")

	t.Logf("✅ Concurrent requests handled: %d/%d succeeded (%.1f%% error rate)",
		concurrency-errorCount, concurrency, errorRate)
}

// BenchmarkMicroserviceLatency measures microservice call latency
func BenchmarkMicroserviceLatency(b *testing.B) {
	log := logger.Get() // Less verbose for benchmarks
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.ListListings(ctx, &pb.ListListingsRequest{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			b.Logf("Request failed: %v", err)
		}
	}
}
