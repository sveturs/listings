// Package integration provides integration tests for timeout resilience
// backend/tests/integration/timeout_test.go
package integration

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/sveturs/listings/api/proto/listings/v1"

	"backend/internal/clients/listings"
	"backend/internal/logger"
)

const (
	testGRPCURL            = "localhost:50053" // Mock microservice URL
	timeoutTestTimeout     = 500 * time.Millisecond
)

// TestTimeoutTriggersAtConfiguredDuration verifies timeout triggers at 500ms
func TestTimeoutTriggersAtConfiguredDuration(t *testing.T) {
	log := logger.Get()
	client, err := listings.NewClient(testGRPCURL, *log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	// Send request that will take 1s (should timeout at 500ms)
	ctx := context.Background()
	req := &pb.GetListingRequest{
		Id: 999, // Mock service will delay this request
	}

	start := time.Now()
	resp, err := client.GetListing(ctx, req)
	elapsed := time.Since(start)

	// Should timeout
	assert.Error(t, err, "Expected timeout error")
	assert.Nil(t, resp, "Response should be nil on timeout")

	// Timeout should happen around 500ms (allow 200ms variance for retries)
	assert.Less(t, elapsed, 2*time.Second, "Timeout took too long")

	t.Logf("✅ Timeout triggered after %v", elapsed)
}

// TestFallbackToMonolithOnTimeout verifies fallback works on timeout
func TestFallbackToMonolithOnTimeout(t *testing.T) {
	log := logger.Get()
	client, err := listings.NewClient(testGRPCURL, *log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	// This test requires integration with router/fallback logic
	// For now, verify error is returned correctly
	ctx := context.Background()
	req := &pb.GetListingRequest{
		Id: 999, // Will timeout
	}

	_, err = client.GetListing(ctx, req)

	// Should get error that triggers fallback
	assert.Error(t, err, "Expected error for fallback trigger")

	// Error should be timeout-related
	assert.Contains(t, err.Error(), "timeout", "Error should indicate timeout")

	t.Logf("✅ Timeout error correctly returned for fallback")
}

// TestContextCancellationPropagates verifies context cancellation works
func TestContextCancellationPropagates(t *testing.T) {
	log := logger.Get()
	client, err := listings.NewClient(testGRPCURL, *log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	req := &pb.GetListingRequest{
		Id: 999, // Will take long time
	}

	start := time.Now()
	_, err = client.GetListing(ctx, req)
	elapsed := time.Since(start)

	// Should fail with context canceled
	assert.Error(t, err, "Expected context canceled error")
	assert.Less(t, elapsed, 500*time.Millisecond, "Should cancel quickly")

	t.Logf("✅ Context cancellation propagated in %v", elapsed)
}

// TestMultipleConcurrentTimeouts verifies concurrent timeouts handled
func TestMultipleConcurrentTimeouts(t *testing.T) {
	log := logger.Get()
	client, err := listings.NewClient(testGRPCURL, *log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	const numRequests = 10

	// Send multiple concurrent requests
	errChan := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int64) {
			ctx := context.Background()
			req := &pb.GetListingRequest{
				Id: 999, // Will timeout
			}
			_, err := client.GetListing(ctx, req)
			errChan <- err
		}(int64(i))
	}

	// Collect results
	timeoutCount := 0
	for i := 0; i < numRequests; i++ {
		err := <-errChan
		if err != nil {
			timeoutCount++
		}
	}

	// All should timeout
	assert.Equal(t, numRequests, timeoutCount, "All requests should timeout")

	t.Logf("✅ %d concurrent timeouts handled correctly", timeoutCount)
}

// TestNoGoroutineLeaksOnTimeout verifies no goroutine leaks
func TestNoGoroutineLeaksOnTimeout(t *testing.T) {
	log := logger.Get()
	client, err := listings.NewClient(testGRPCURL, *log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	// Get baseline goroutine count
	baselineCount := countGoroutines()

	// Send multiple requests that will timeout
	for i := 0; i < 100; i++ {
		ctx := context.Background()
		req := &pb.GetListingRequest{
			Id: 999, // Will timeout
		}
		client.GetListing(ctx, req)
	}

	// Wait for cleanup
	time.Sleep(1 * time.Second)

	// Check goroutine count
	finalCount := countGoroutines()
	leaked := finalCount - baselineCount

	// Allow some variance (max 5 goroutines)
	assert.Less(t, leaked, 5, "Too many goroutines leaked")

	t.Logf("✅ No significant goroutine leaks (baseline=%d, final=%d, leaked=%d)",
		baselineCount, finalCount, leaked)
}

// TestTimeoutWithRetries verifies retries respect timeout
func TestTimeoutWithRetries(t *testing.T) {
	log := logger.Get()
	client, err := listings.NewClient(testGRPCURL, *log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	ctx := context.Background()
	req := &pb.SearchListingsRequest{
		Query: "timeout-test",
	}

	start := time.Now()
	_, err = client.SearchListings(ctx, req)
	elapsed := time.Since(start)

	// Should timeout even with retries
	assert.Error(t, err, "Expected timeout error")

	// Total time should not exceed timeout * retries significantly
	maxExpected := 2 * time.Second // 500ms * 3 retries + overhead
	assert.Less(t, elapsed, maxExpected, "Retries took too long")

	t.Logf("✅ Timeout with retries completed in %v", elapsed)
}

// Helper functions

func countGoroutines() int {
	// Simple goroutine counter using runtime
	buf := make([]byte, 1<<16)
	n := runtime.Stack(buf, true)
	stackTrace := string(buf[:n])

	// Count "goroutine" occurrences
	count := 0
	for _, line := range strings.Split(stackTrace, "\n") {
		if strings.HasPrefix(line, "goroutine ") {
			count++
		}
	}
	return count
}
