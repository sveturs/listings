// Package load provides load tests for resilience patterns
// backend/tests/load/resilience_load_test.go
//
// Run with: go test -v -bench=. -benchtime=30s ./tests/load/resilience_load_test.go
package load

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	pb "github.com/sveturs/listings/api/proto/listings/v1"

	"backend/internal/clients/listings"
	"backend/internal/logger"
)

const (
	testGRPCURL = "localhost:50051"
)

// BenchmarkLoadWith10PercentTimeouts simulates 1000 RPS with 10% timeouts
func BenchmarkLoadWith10PercentTimeouts(b *testing.B) {
	log := logger.New("load-test")
	client, err := listings.NewClient(testGRPCURL, log)
	require.NoError(b, err, "Failed to create gRPC client")
	defer client.Close()

	ctx := context.Background()
	requestCounter := int32(0)
	successCounter := int32(0)
	timeoutCounter := int32(0)

	// Track initial memory
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			count := atomic.AddInt32(&requestCounter, 1)

			// 10% of requests timeout (ID 999)
			reqID := int64(count % 100)
			if reqID < 10 {
				reqID = 999 // Timeout request
			}

			req := &pb.GetListingRequest{Id: reqID}
			_, err := client.GetListing(ctx, req)

			if err == nil {
				atomic.AddInt32(&successCounter, 1)
			} else {
				atomic.AddInt32(&timeoutCounter, 1)
			}
		}
	})

	b.StopTimer()

	// Collect final memory stats
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// Calculate metrics
	totalRequests := atomic.LoadInt32(&requestCounter)
	totalSuccess := atomic.LoadInt32(&successCounter)
	totalTimeout := atomic.LoadInt32(&timeoutCounter)

	successRate := float64(totalSuccess) / float64(totalRequests) * 100
	timeoutRate := float64(totalTimeout) / float64(totalRequests) * 100
	memIncrease := m2.Alloc - m1.Alloc

	b.ReportMetric(float64(totalRequests)/b.Elapsed().Seconds(), "req/s")
	b.ReportMetric(successRate, "success_%")
	b.ReportMetric(timeoutRate, "timeout_%")
	b.ReportMetric(float64(memIncrease)/(1024*1024), "mem_mb")

	b.Logf("\n"+
		"==== Load Test Results (10%% Timeout) ====\n"+
		"Total Requests: %d\n"+
		"Success: %d (%.1f%%)\n"+
		"Timeout: %d (%.1f%%)\n"+
		"RPS: %.1f\n"+
		"Memory Increase: %.2f MB\n"+
		"Duration: %v",
		totalRequests, totalSuccess, successRate,
		totalTimeout, timeoutRate,
		float64(totalRequests)/b.Elapsed().Seconds(),
		float64(memIncrease)/(1024*1024),
		b.Elapsed())
}

// BenchmarkLoadWithCircuitBreaker simulates 500 RPS with circuit breaker opening/closing
func BenchmarkLoadWithCircuitBreaker(b *testing.B) {
	log := logger.New("load-test")
	client, err := listings.NewClient(testGRPCURL, log)
	require.NoError(b, err, "Failed to create gRPC client")
	defer client.Close()

	ctx := context.Background()
	requestCounter := int32(0)
	successCounter := int32(0)
	circuitOpenCounter := int32(0)

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	b.ResetTimer()

	// Simulate circuit breaker scenario:
	// - First 100 requests: 50% fail (triggers circuit breaker)
	// - Next 50 requests: circuit open (rejected)
	// - Wait 30s: circuit half-open
	// - Next 100 requests: all succeed (circuit closes)
	// - Remaining: normal operations

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			count := atomic.AddInt32(&requestCounter, 1)

			var reqID int64
			if count < 100 {
				// Phase 1: 50% failures to trigger circuit
				if count%2 == 0 {
					reqID = 777 // Failure
				} else {
					reqID = 1 // Success
				}
			} else if count < 150 {
				// Phase 2: Circuit open, requests rejected
				reqID = 1 // Will be rejected by circuit breaker
			} else {
				// Phase 3+: Normal operations
				reqID = int64(count % 100)
			}

			req := &pb.GetListingRequest{Id: reqID}
			_, err := client.GetListing(ctx, req)

			if err == nil {
				atomic.AddInt32(&successCounter, 1)
			} else if count >= 100 && count < 150 {
				// Circuit open rejections
				atomic.AddInt32(&circuitOpenCounter, 1)
			}
		}
	})

	b.StopTimer()

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	totalRequests := atomic.LoadInt32(&requestCounter)
	totalSuccess := atomic.LoadInt32(&successCounter)
	totalCircuitOpen := atomic.LoadInt32(&circuitOpenCounter)

	successRate := float64(totalSuccess) / float64(totalRequests) * 100
	memIncrease := m2.Alloc - m1.Alloc

	b.ReportMetric(float64(totalRequests)/b.Elapsed().Seconds(), "req/s")
	b.ReportMetric(successRate, "success_%")
	b.ReportMetric(float64(memIncrease)/(1024*1024), "mem_mb")

	b.Logf("\n"+
		"==== Circuit Breaker Load Test ====\n"+
		"Total Requests: %d\n"+
		"Success: %d (%.1f%%)\n"+
		"Circuit Open Rejections: %d\n"+
		"RPS: %.1f\n"+
		"Memory Increase: %.2f MB",
		totalRequests, totalSuccess, successRate,
		totalCircuitOpen,
		float64(totalRequests)/b.Elapsed().Seconds(),
		float64(memIncrease)/(1024*1024))
}

// BenchmarkLoadWithMixedSuccessFailure simulates 2000 RPS with mixed results
func BenchmarkLoadWithMixedSuccessFailure(b *testing.B) {
	log := logger.New("load-test")
	client, err := listings.NewClient(testGRPCURL, log)
	require.NoError(b, err, "Failed to create gRPC client")
	defer client.Close()

	ctx := context.Background()
	requestCounter := int32(0)
	successCounter := int32(0)
	failureCounter := int32(0)
	timeoutCounter := int32(0)

	// Latency tracking
	var totalLatencyNs int64
	latencies := make([]time.Duration, 0, b.N)
	var latencyMutex sync.Mutex

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			count := atomic.AddInt32(&requestCounter, 1)

			// Mixed load:
			// - 70% success
			// - 20% timeout
			// - 10% error
			var reqID int64
			mod := count % 100
			if mod < 70 {
				reqID = 1 // Success
			} else if mod < 90 {
				reqID = 999 // Timeout
			} else {
				reqID = 777 // Error
			}

			req := &pb.GetListingRequest{Id: reqID}

			start := time.Now()
			_, err := client.GetListing(ctx, req)
			latency := time.Since(start)

			// Record latency
			atomic.AddInt64(&totalLatencyNs, int64(latency))
			latencyMutex.Lock()
			latencies = append(latencies, latency)
			latencyMutex.Unlock()

			if err == nil {
				atomic.AddInt32(&successCounter, 1)
			} else if containsTimeout(err) {
				atomic.AddInt32(&timeoutCounter, 1)
			} else {
				atomic.AddInt32(&failureCounter, 1)
			}
		}
	})

	b.StopTimer()

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// Calculate metrics
	totalRequests := atomic.LoadInt32(&requestCounter)
	totalSuccess := atomic.LoadInt32(&successCounter)
	totalFailure := atomic.LoadInt32(&failureCounter)
	totalTimeout := atomic.LoadInt32(&timeoutCounter)

	avgLatency := time.Duration(atomic.LoadInt64(&totalLatencyNs) / int64(totalRequests))
	p99Latency := calculateP99(latencies)

	successRate := float64(totalSuccess) / float64(totalRequests) * 100
	memIncrease := m2.Alloc - m1.Alloc

	b.ReportMetric(float64(totalRequests)/b.Elapsed().Seconds(), "req/s")
	b.ReportMetric(successRate, "success_%")
	b.ReportMetric(float64(avgLatency.Milliseconds()), "avg_latency_ms")
	b.ReportMetric(float64(p99Latency.Milliseconds()), "p99_latency_ms")
	b.ReportMetric(float64(memIncrease)/(1024*1024), "mem_mb")

	b.Logf("\n"+
		"==== Mixed Load Test ====\n"+
		"Total Requests: %d\n"+
		"Success: %d (%.1f%%)\n"+
		"Timeout: %d\n"+
		"Failure: %d\n"+
		"RPS: %.1f\n"+
		"Avg Latency: %v\n"+
		"P99 Latency: %v\n"+
		"Memory Increase: %.2f MB",
		totalRequests, totalSuccess, successRate,
		totalTimeout, totalFailure,
		float64(totalRequests)/b.Elapsed().Seconds(),
		avgLatency, p99Latency,
		float64(memIncrease)/(1024*1024))

	// Assert P99 < 200ms
	require.Less(b, p99Latency, 200*time.Millisecond,
		"P99 latency exceeded 200ms threshold")
}

// TestMemoryStabilityUnderLoad verifies no memory leaks during sustained load
func TestMemoryStabilityUnderLoad(t *testing.T) {
	log := logger.New("memory-test")
	client, err := listings.NewClient(testGRPCURL, log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	ctx := context.Background()

	// Baseline measurement
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Run sustained load for 10 seconds
	const duration = 10 * time.Second
	const targetRPS = 100

	done := make(chan struct{})
	var wg sync.WaitGroup

	// Request generator
	ticker := time.NewTicker(time.Second / targetRPS)
	defer ticker.Stop()

	go func() {
		time.Sleep(duration)
		close(done)
	}()

	requestCount := 0

loop:
	for {
		select {
		case <-done:
			break loop
		case <-ticker.C:
			wg.Add(1)
			go func() {
				defer wg.Done()
				req := &pb.GetListingRequest{Id: 1}
				client.GetListing(ctx, req)
			}()
			requestCount++
		}
	}

	wg.Wait()

	// Final measurement
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// Calculate memory increase
	memIncrease := int64(m2.Alloc) - int64(m1.Alloc)
	memIncreasePercent := float64(memIncrease) / float64(m1.Alloc) * 100

	t.Logf("Memory Stability Test Results:")
	t.Logf("  Requests: %d", requestCount)
	t.Logf("  Initial Memory: %.2f MB", float64(m1.Alloc)/(1024*1024))
	t.Logf("  Final Memory: %.2f MB", float64(m2.Alloc)/(1024*1024))
	t.Logf("  Increase: %.2f MB (%.1f%%)", float64(memIncrease)/(1024*1024), memIncreasePercent)

	// Memory increase should be < 10%
	require.Less(t, memIncreasePercent, 10.0,
		"Memory increased more than 10%% - possible leak")

	t.Logf("✅ Memory stable under load: %.1f%% increase", memIncreasePercent)
}

// TestNoGoroutineLeaksUnderLoad verifies goroutine cleanup
func TestNoGoroutineLeaksUnderLoad(t *testing.T) {
	log := logger.New("goroutine-test")
	client, err := listings.NewClient(testGRPCURL, log)
	require.NoError(t, err, "Failed to create gRPC client")
	defer client.Close()

	ctx := context.Background()

	// Baseline
	baseline := runtime.NumGoroutine()

	// Generate load
	const numRequests = 1000
	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req := &pb.GetListingRequest{Id: int64(id % 100)}
			client.GetListing(ctx, req)
		}(i)
	}

	wg.Wait()

	// Wait for cleanup
	time.Sleep(2 * time.Second)
	runtime.GC()

	// Final count
	final := runtime.NumGoroutine()
	leaked := final - baseline

	t.Logf("Goroutine Leak Test Results:")
	t.Logf("  Baseline: %d", baseline)
	t.Logf("  Final: %d", final)
	t.Logf("  Leaked: %d", leaked)

	// Allow some variance (max 10 goroutines)
	require.Less(t, leaked, 10,
		"Too many goroutines leaked")

	t.Logf("✅ No goroutine leaks: %d leaked", leaked)
}

// Helper functions

func containsTimeout(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return containsAny(errStr, []string{"timeout", "deadline exceeded"})
}

func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if len(s) >= len(substr) && contains(s, substr) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func calculateP99(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	// Sort latencies
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sortDurations(sorted)

	// Calculate 99th percentile index
	index := int(float64(len(sorted)) * 0.99)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

func sortDurations(durations []time.Duration) {
	// Simple bubble sort (good enough for testing)
	n := len(durations)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if durations[j] > durations[j+1] {
				durations[j], durations[j+1] = durations[j+1], durations[j]
			}
		}
	}
}
