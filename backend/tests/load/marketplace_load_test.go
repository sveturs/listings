//go:build ignore
// +build ignore

// backend/tests/load/marketplace_load_test.go
// DEPRECATED: Load tests for unified architecture that was removed in Phase 7
// To enable: remove //go:build ignore directive and update to use marketplace service
package load

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"backend/internal/domain/search"

	"backend/internal/domain/models"
	"backend/internal/proj/unified/service"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadTestMetrics собирает метрики нагрузочного тестирования
type LoadTestMetrics struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	TotalLatencyMs  int64
	MinLatencyMs    int64
	MaxLatencyMs    int64
	Latencies       []int64
	mu              sync.Mutex
}

func (m *LoadTestMetrics) RecordRequest(latencyMs int64, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	atomic.AddInt64(&m.TotalRequests, 1)
	atomic.AddInt64(&m.TotalLatencyMs, latencyMs)

	if success {
		atomic.AddInt64(&m.SuccessRequests, 1)
	} else {
		atomic.AddInt64(&m.FailedRequests, 1)
	}

	// Update min/max
	if m.MinLatencyMs == 0 || latencyMs < m.MinLatencyMs {
		m.MinLatencyMs = latencyMs
	}
	if latencyMs > m.MaxLatencyMs {
		m.MaxLatencyMs = latencyMs
	}

	m.Latencies = append(m.Latencies, latencyMs)
}

func (m *LoadTestMetrics) CalculatePercentile(p float64) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.Latencies) == 0 {
		return 0
	}

	// Simple percentile calculation (copy and sort)
	sorted := make([]int64, len(m.Latencies))
	copy(sorted, m.Latencies)

	// Bubble sort (good enough for tests)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	idx := int(float64(len(sorted)) * p)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}

	return sorted[idx]
}

func (m *LoadTestMetrics) AvgLatencyMs() int64 {
	if m.TotalRequests == 0 {
		return 0
	}
	return m.TotalLatencyMs / m.TotalRequests
}

func (m *LoadTestMetrics) ErrorRate() float64 {
	if m.TotalRequests == 0 {
		return 0
	}
	return float64(m.FailedRequests) / float64(m.TotalRequests) * 100
}

func (m *LoadTestMetrics) Throughput(durationSec int) float64 {
	if durationSec == 0 {
		return 0
	}
	return float64(m.TotalRequests) / float64(durationSec)
}

// TestLoad_Baseline проверяет baseline (только monolith, без microservice)
// defaultRoutingContext создаёт дефолтный routing context для тестов
func defaultRoutingContext() *service.RoutingContext {
	return &service.RoutingContext{
		UserID:  100,
		IsAdmin: false,
	}
}

func TestLoad_Baseline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	targetRPS := 100  // requests per second
	durationSec := 10 // test duration
	totalRequests := targetRPS * durationSec

	t.Run("Baseline - Monolith only (no microservice)", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		// Microservice disabled (default)

		metrics := &LoadTestMetrics{}
		var wg sync.WaitGroup

		startTime := time.Now()

		// Act - Generate load
		for i := 0; i < totalRequests; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				reqStart := time.Now()

				unified := &models.UnifiedListing{
					UserID:     int(rand.Int31n(1000)),
					SourceType: service.SourceTypeC2C,
					Title:      fmt.Sprintf("Load Test %d", idx),
					Price:      float64(rand.Intn(1000)),
					CategoryID: rand.Intn(10) + 1,
				}

				_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

				latency := time.Since(reqStart).Milliseconds()
				metrics.RecordRequest(latency, err == nil)
			}(i)

			// Rate limiting (spread requests over time)
			time.Sleep(time.Second / time.Duration(targetRPS))
		}

		wg.Wait()
		duration := time.Since(startTime).Seconds()

		// Assert
		assert.Equal(t, int64(totalRequests), metrics.TotalRequests)
		assert.True(t, metrics.ErrorRate() < 0.1, "Error rate should be < 0.1%%")

		p50 := metrics.CalculatePercentile(0.50)
		p95 := metrics.CalculatePercentile(0.95)
		p99 := metrics.CalculatePercentile(0.99)

		t.Logf("Baseline Metrics:")
		t.Logf("  Total Requests: %d", metrics.TotalRequests)
		t.Logf("  Success: %d", metrics.SuccessRequests)
		t.Logf("  Failed: %d", metrics.FailedRequests)
		t.Logf("  Error Rate: %.2f%%", metrics.ErrorRate())
		t.Logf("  Throughput: %.2f req/sec", metrics.Throughput(int(duration)))
		t.Logf("  Latency - Avg: %dms, Min: %dms, Max: %dms", metrics.AvgLatencyMs(), metrics.MinLatencyMs, metrics.MaxLatencyMs)
		t.Logf("  Latency - P50: %dms, P95: %dms, P99: %dms", p50, p95, p99)

		// Performance criteria (mocks are fast, so should be < 10ms)
		assert.Less(t, p99, int64(50), "P99 latency should be < 50ms for baseline")
	})
}

// TestLoad_Microservice_10Percent проверяет с 10% трафика на microservice
func TestLoad_Microservice_10Percent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	targetRPS := 100
	durationSec := 10
	totalRequests := targetRPS * durationSec
	microservicePercent := 10 // 10% traffic to microservice

	t.Run("10% traffic to microservice", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

		metrics := &LoadTestMetrics{}
		microserviceCount := int64(0)
		var wg sync.WaitGroup

		startTime := time.Now()

		// Act
		for i := 0; i < totalRequests; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				// 10% chance to use microservice
				useMicroservice := rand.Intn(100) < microservicePercent
				if useMicroservice {
					svc.SetListingsGRPCClient(mockGRPC, true)
					atomic.AddInt64(&microserviceCount, 1)
				} else {
					svc.SetListingsGRPCClient(nil, false)
				}

				reqStart := time.Now()

				unified := &models.UnifiedListing{
					UserID:     int(rand.Int31n(1000)),
					SourceType: service.SourceTypeC2C,
					Title:      fmt.Sprintf("Load Test %d", idx),
					Price:      float64(rand.Intn(1000)),
					CategoryID: rand.Intn(10) + 1,
				}

				_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

				latency := time.Since(reqStart).Milliseconds()
				metrics.RecordRequest(latency, err == nil)
			}(i)

			time.Sleep(time.Second / time.Duration(targetRPS))
		}

		wg.Wait()
		duration := time.Since(startTime).Seconds()

		// Assert
		assert.Equal(t, int64(totalRequests), metrics.TotalRequests)
		assert.True(t, metrics.ErrorRate() < 0.1, "Error rate should be < 0.1%%")

		p50 := metrics.CalculatePercentile(0.50)
		p95 := metrics.CalculatePercentile(0.95)
		p99 := metrics.CalculatePercentile(0.99)

		t.Logf("10%% Microservice Metrics:")
		t.Logf("  Total Requests: %d", metrics.TotalRequests)
		t.Logf("  Microservice: %d (%.1f%%)", microserviceCount, float64(microserviceCount)/float64(totalRequests)*100)
		t.Logf("  Success: %d", metrics.SuccessRequests)
		t.Logf("  Failed: %d", metrics.FailedRequests)
		t.Logf("  Error Rate: %.2f%%", metrics.ErrorRate())
		t.Logf("  Throughput: %.2f req/sec", metrics.Throughput(int(duration)))
		t.Logf("  Latency - Avg: %dms, Min: %dms, Max: %dms", metrics.AvgLatencyMs(), metrics.MinLatencyMs, metrics.MaxLatencyMs)
		t.Logf("  Latency - P50: %dms, P95: %dms, P99: %dms", p50, p95, p99)

		assert.Less(t, p99, int64(100), "P99 latency should be < 100ms with 10%% microservice")
	})
}

// TestLoad_Microservice_100Percent проверяет со 100% трафика на microservice
func TestLoad_Microservice_100Percent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	targetRPS := 100
	durationSec := 10
	totalRequests := targetRPS * durationSec

	t.Run("100% traffic to microservice", func(t *testing.T) {
		// Arrange
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		metrics := &LoadTestMetrics{}
		var wg sync.WaitGroup

		startTime := time.Now()

		// Act
		for i := 0; i < totalRequests; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				reqStart := time.Now()

				unified := &models.UnifiedListing{
					UserID:     int(rand.Int31n(1000)),
					SourceType: service.SourceTypeC2C,
					Title:      fmt.Sprintf("Load Test %d", idx),
					Price:      float64(rand.Intn(1000)),
					CategoryID: rand.Intn(10) + 1,
				}

				_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

				latency := time.Since(reqStart).Milliseconds()
				metrics.RecordRequest(latency, err == nil)
			}(i)

			time.Sleep(time.Second / time.Duration(targetRPS))
		}

		wg.Wait()
		duration := time.Since(startTime).Seconds()

		// Assert
		assert.Equal(t, int64(totalRequests), metrics.TotalRequests)
		assert.Equal(t, totalRequests, mockGRPC.CreateCallCount, "All requests should go to microservice")
		assert.True(t, metrics.ErrorRate() < 0.1, "Error rate should be < 0.1%%")

		p50 := metrics.CalculatePercentile(0.50)
		p95 := metrics.CalculatePercentile(0.95)
		p99 := metrics.CalculatePercentile(0.99)

		t.Logf("100%% Microservice Metrics:")
		t.Logf("  Total Requests: %d", metrics.TotalRequests)
		t.Logf("  Success: %d", metrics.SuccessRequests)
		t.Logf("  Failed: %d", metrics.FailedRequests)
		t.Logf("  Error Rate: %.2f%%", metrics.ErrorRate())
		t.Logf("  Throughput: %.2f req/sec", metrics.Throughput(int(duration)))
		t.Logf("  Latency - Avg: %dms, Min: %dms, Max: %dms", metrics.AvgLatencyMs(), metrics.MinLatencyMs, metrics.MaxLatencyMs)
		t.Logf("  Latency - P50: %dms, P95: %dms, P99: %dms", p50, p95, p99)

		// With 100% microservice, should still be fast with mocks
		assert.Less(t, p99, int64(100), "P99 latency should be < 100ms with 100%% microservice")
	})
}

// TestLoad_Spike проверяет spike test (резкое увеличение нагрузки)
func TestLoad_Spike(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Spike test - 0 to 200 RPS", func(t *testing.T) {
		// Arrange
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		metrics := &LoadTestMetrics{}
		var wg sync.WaitGroup

		// Phase 1: Warmup (10 RPS for 2 seconds)
		warmupRPS := 10
		warmupDuration := 2
		for i := 0; i < warmupRPS*warmupDuration; i++ {
			wg.Add(1)
			go executeRequest(ctx, svc, i, metrics, &wg)
			time.Sleep(time.Second / time.Duration(warmupRPS))
		}

		// Phase 2: Spike (200 RPS for 5 seconds)
		spikeRPS := 200
		spikeDuration := 5
		startTime := time.Now()
		for i := 0; i < spikeRPS*spikeDuration; i++ {
			wg.Add(1)
			go executeRequest(ctx, svc, i, metrics, &wg)
			time.Sleep(time.Second / time.Duration(spikeRPS))
		}

		wg.Wait()
		duration := time.Since(startTime).Seconds()

		// Assert
		require.Greater(t, metrics.TotalRequests, int64(0))
		assert.True(t, metrics.ErrorRate() < 1.0, "Error rate should be < 1%% during spike")

		p99 := metrics.CalculatePercentile(0.99)

		t.Logf("Spike Test Metrics:")
		t.Logf("  Total Requests: %d", metrics.TotalRequests)
		t.Logf("  Success: %d", metrics.SuccessRequests)
		t.Logf("  Failed: %d", metrics.FailedRequests)
		t.Logf("  Error Rate: %.2f%%", metrics.ErrorRate())
		t.Logf("  Throughput: %.2f req/sec", metrics.Throughput(int(duration)))
		t.Logf("  Latency - P99: %dms", p99)

		// Should handle spike without major degradation
		assert.Less(t, p99, int64(200), "P99 latency should be < 200ms during spike")
	})
}

// TestLoad_Endurance проверяет endurance test (длительная нагрузка)
func TestLoad_Endurance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	targetRPS := 50
	durationSec := 30 // 30 seconds (в реальности было бы 10 минут)
	totalRequests := targetRPS * durationSec

	t.Run("Endurance test - 30 seconds at 50 RPS", func(t *testing.T) {
		// Arrange
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		metrics := &LoadTestMetrics{}
		var wg sync.WaitGroup

		startTime := time.Now()

		// Act
		for i := 0; i < totalRequests; i++ {
			wg.Add(1)
			go executeRequest(ctx, svc, i, metrics, &wg)
			time.Sleep(time.Second / time.Duration(targetRPS))
		}

		wg.Wait()
		duration := time.Since(startTime).Seconds()

		// Assert
		assert.Equal(t, int64(totalRequests), metrics.TotalRequests)
		assert.True(t, metrics.ErrorRate() < 0.1, "Error rate should be < 0.1%% for endurance")

		p95 := metrics.CalculatePercentile(0.95)
		p99 := metrics.CalculatePercentile(0.99)

		t.Logf("Endurance Test Metrics:")
		t.Logf("  Total Requests: %d", metrics.TotalRequests)
		t.Logf("  Duration: %.1f seconds", duration)
		t.Logf("  Success: %d", metrics.SuccessRequests)
		t.Logf("  Failed: %d", metrics.FailedRequests)
		t.Logf("  Error Rate: %.2f%%", metrics.ErrorRate())
		t.Logf("  Throughput: %.2f req/sec", metrics.Throughput(int(duration)))
		t.Logf("  Latency - P95: %dms, P99: %dms", p95, p99)

		// No memory leaks or degradation
		assert.Less(t, p99, int64(100), "P99 should remain stable during endurance test")
	})
}

// Helper functions

func executeRequest(ctx context.Context, svc *service.MarketplaceService, idx int, metrics *LoadTestMetrics, wg *sync.WaitGroup) {
	defer wg.Done()

	reqStart := time.Now()

	unified := &models.UnifiedListing{
		UserID:     int(rand.Int31n(1000)),
		SourceType: service.SourceTypeC2C,
		Title:      fmt.Sprintf("Load Test %d", idx),
		Price:      float64(rand.Intn(1000)),
		CategoryID: rand.Intn(10) + 1,
	}

	_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

	latency := time.Since(reqStart).Milliseconds()
	metrics.RecordRequest(latency, err == nil)
}

// Mock repositories

type mockC2CRepository struct{}

func (m *mockC2CRepository) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	// Simulate small delay
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
	return rand.Intn(10000), nil
}

func (m *mockC2CRepository) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3)))
	return &models.MarketplaceListing{ID: id}, nil
}

func (m *mockC2CRepository) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
	return nil
}

func (m *mockC2CRepository) DeleteListing(ctx context.Context, id int) error {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3)))
	return nil
}

type mockOpenSearchRepository struct{}

func (m *mockOpenSearchRepository) SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	return &search.ServiceResult{Items: []*models.MarketplaceListing{}, Total: 0, Page: 0, Size: 10}, nil
}

func (m *mockOpenSearchRepository) Index(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}

func (m *mockOpenSearchRepository) Delete(ctx context.Context, listingID int) error {
	return nil
}
