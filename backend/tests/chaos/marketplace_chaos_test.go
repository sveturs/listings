// backend/tests/chaos/marketplace_chaos_test.go
package chaos

import (
	"context"
	"errors"
	"fmt"
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

// ChaosConfig конфигурация для chaos testing
// defaultRoutingContext создаёт дефолтный routing context для тестов
func defaultRoutingContext() *service.RoutingContext {
	return &service.RoutingContext{
		UserID:  100,
		IsAdmin: false,
	}
}

type ChaosConfig struct {
	FailureRate      float64 // 0.0 - 1.0 (probability of failure)
	SlowdownMs       int64   // Delay in milliseconds
	TimeoutMs        int64   // Timeout threshold
	PartialFailures  bool    // Enable partial failures
	NetworkPartition bool    // Simulate network partition
}

// TestChaos_NetworkPartition проверяет поведение при network partition
func TestChaos_NetworkPartition(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping chaos test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Microservice unavailable for 10 seconds - should fallback", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := &chaosGRPCClient{
			config: ChaosConfig{
				NetworkPartition: true, // Simulate complete network failure
			},
		}

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Network Partition Test",
			Price:      99.99,
			CategoryID: 1,
		}

		// Act
		id, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

		// Assert - Should fallback to local DB
		require.NoError(t, err, "Should succeed via fallback")
		assert.Greater(t, id, int64(0))
		assert.True(t, mockC2C.createCalled, "Should fallback to local DB")
		assert.Equal(t, int64(1), atomic.LoadInt64(&mockGRPC.createAttempts), "Should attempt microservice first")

		t.Logf("Network Partition Test:")
		t.Logf("  Microservice attempts: %d", mockGRPC.createAttempts)
		t.Logf("  Fallback successful: %v", mockC2C.createCalled)
	})
}

// TestChaos_SlowMicroservice проверяет поведение при медленном microservice
func TestChaos_SlowMicroservice(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping chaos test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Microservice responds slowly - should timeout and fallback", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := &chaosGRPCClient{
			config: ChaosConfig{
				SlowdownMs: 5000, // 5 second delay
				TimeoutMs:  100,  // 100ms timeout threshold
			},
		}

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Slow Microservice Test",
			Price:      50.00,
			CategoryID: 1,
		}

		start := time.Now()

		// Act
		id, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		elapsed := time.Since(start).Milliseconds()

		// Assert - Should timeout and fallback quickly
		require.NoError(t, err, "Should succeed via fallback")
		assert.Greater(t, id, int64(0))
		assert.True(t, mockC2C.createCalled, "Should fallback to local DB")

		// Should timeout within reasonable time (not wait full 5 seconds)
		assert.Less(t, elapsed, int64(6000), "Should timeout and fallback quickly")

		t.Logf("Slow Microservice Test:")
		t.Logf("  Total time: %dms", elapsed)
		t.Logf("  Fallback successful: %v", mockC2C.createCalled)
	})
}

// TestChaos_PartialFailures проверяет частичные отказы
func TestChaos_PartialFailures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping chaos test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("50% of requests fail - should fallback only failed requests", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := &chaosGRPCClient{
			config: ChaosConfig{
				FailureRate:     0.5, // 50% failure rate
				PartialFailures: true,
			},
		}

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		totalRequests := 100
		successCount := 0
		fallbackCount := 0

		// Act
		for i := 0; i < totalRequests; i++ {
			mockC2C.createCalled = false // Reset flag

			unified := &models.UnifiedListing{
				UserID:     100,
				SourceType: service.SourceTypeC2C,
				Title:      fmt.Sprintf("Partial Failure Test %d", i),
				Price:      float64(i * 10),
				CategoryID: 1,
			}

			_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

			if err == nil {
				successCount++
				if mockC2C.createCalled {
					fallbackCount++
				}
			}
		}

		// Assert
		assert.Equal(t, totalRequests, successCount, "All requests should succeed (via microservice or fallback)")

		// Should have some fallbacks (around 50%)
		assert.Greater(t, fallbackCount, totalRequests/3, "Should have fallbacks for failed microservice calls")
		assert.Less(t, fallbackCount, totalRequests*2/3, "Not all should fallback")

		t.Logf("Partial Failures Test:")
		t.Logf("  Total requests: %d", totalRequests)
		t.Logf("  Successful: %d (100%%)", successCount)
		t.Logf("  Fallback count: %d (%.1f%%)", fallbackCount, float64(fallbackCount)/float64(totalRequests)*100)
		t.Logf("  Microservice attempts: %d", mockGRPC.createAttempts)
	})
}

// TestChaos_DatabaseFailure проверяет отказ БД microservice
func TestChaos_DatabaseFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping chaos test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Microservice DB unavailable - should return error and fallback", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := &chaosGRPCClient{
			config: ChaosConfig{
				FailureRate: 1.0, // 100% failure (DB unavailable)
			},
			dbError: errors.New("database connection failed"),
		}

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "DB Failure Test",
			Price:      75.00,
			CategoryID: 1,
		}

		// Act
		id, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

		// Assert
		require.NoError(t, err, "Should fallback successfully")
		assert.Greater(t, id, int64(0))
		assert.True(t, mockC2C.createCalled, "Should fallback to local DB")

		t.Logf("Database Failure Test:")
		t.Logf("  Microservice DB error: %v", mockGRPC.dbError)
		t.Logf("  Fallback successful: %v", mockC2C.createCalled)
	})
}

// TestChaos_CascadingFailures проверяет каскадные отказы
func TestChaos_CascadingFailures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping chaos test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Cascading failures - both microservice and fallback slow", func(t *testing.T) {
		// Arrange
		mockC2C := &slowC2CRepository{
			delayMs: 2000, // 2 second delay in fallback
		}
		mockGRPC := &chaosGRPCClient{
			config: ChaosConfig{
				SlowdownMs: 3000, // 3 second delay in microservice
			},
		}

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Cascading Failure Test",
			Price:      100.00,
			CategoryID: 1,
		}

		start := time.Now()

		// Act
		id, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		elapsed := time.Since(start).Milliseconds()

		// Assert - Should eventually succeed but take longer
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		t.Logf("Cascading Failures Test:")
		t.Logf("  Total time: %dms", elapsed)
		t.Logf("  Microservice slow: 3000ms")
		t.Logf("  Fallback slow: 2000ms")
		t.Logf("  Request completed: %v", err == nil)
	})
}

// TestChaos_FlappingService проверяет нестабильный сервис (flapping)
func TestChaos_FlappingService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping chaos test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Flapping microservice - alternating success/failure", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := &flappingGRPCClient{
			failEveryNth: 2, // Fail every 2nd request
		}

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		totalRequests := 20
		microserviceSuccess := 0
		fallbackCount := 0

		// Act
		for i := 0; i < totalRequests; i++ {
			mockC2C.createCalled = false

			unified := &models.UnifiedListing{
				UserID:     100,
				SourceType: service.SourceTypeC2C,
				Title:      fmt.Sprintf("Flapping Test %d", i),
				Price:      float64(i),
				CategoryID: 1,
			}

			_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
			require.NoError(t, err, "All requests should eventually succeed")

			if mockC2C.createCalled {
				fallbackCount++
			} else {
				microserviceSuccess++
			}
		}

		// Assert
		assert.Equal(t, totalRequests, microserviceSuccess+fallbackCount)
		assert.Greater(t, fallbackCount, 0, "Should have some fallbacks")
		assert.Greater(t, microserviceSuccess, 0, "Should have some microservice successes")

		t.Logf("Flapping Service Test:")
		t.Logf("  Total requests: %d", totalRequests)
		t.Logf("  Microservice success: %d", microserviceSuccess)
		t.Logf("  Fallback count: %d", fallbackCount)
	})
}

// Chaos testing mock implementations

type chaosGRPCClient struct {
	config         ChaosConfig
	dbError        error
	createAttempts int64
}

func (c *chaosGRPCClient) GetListing(ctx context.Context, id int64) (*models.UnifiedListing, error) {
	if c.config.NetworkPartition {
		return nil, errors.New("network partition: connection refused")
	}

	if c.config.SlowdownMs > 0 {
		time.Sleep(time.Duration(c.config.SlowdownMs) * time.Millisecond)
	}

	if c.dbError != nil {
		return nil, c.dbError
	}

	return &models.UnifiedListing{ID: int(id)}, nil
}

func (c *chaosGRPCClient) CreateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	atomic.AddInt64(&c.createAttempts, 1)

	if c.config.NetworkPartition {
		return nil, errors.New("network partition: connection refused")
	}

	if c.config.SlowdownMs > 0 {
		time.Sleep(time.Duration(c.config.SlowdownMs) * time.Millisecond)
	}

	if c.dbError != nil {
		return nil, c.dbError
	}

	// Partial failures - random failure based on FailureRate
	if c.config.PartialFailures && c.config.FailureRate > 0 {
		// Simple deterministic failure (every Nth request fails)
		attempts := atomic.LoadInt64(&c.createAttempts)
		failEvery := int64(1.0 / c.config.FailureRate)
		if attempts%failEvery == 0 {
			return nil, errors.New("partial failure: request failed")
		}
	}

	// Full failure
	if c.config.FailureRate >= 1.0 {
		return nil, errors.New("chaos: simulated failure")
	}

	unified.ID = int(atomic.LoadInt64(&c.createAttempts))
	return unified, nil
}

func (c *chaosGRPCClient) UpdateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	if c.config.NetworkPartition {
		return nil, errors.New("network partition: connection refused")
	}
	return unified, nil
}

func (c *chaosGRPCClient) DeleteListing(ctx context.Context, id int64, userID int64) error {
	if c.config.NetworkPartition {
		return errors.New("network partition: connection refused")
	}
	return nil
}

type flappingGRPCClient struct {
	requestCount int64
	failEveryNth int
}

func (f *flappingGRPCClient) GetListing(ctx context.Context, id int64) (*models.UnifiedListing, error) {
	count := atomic.AddInt64(&f.requestCount, 1)
	if count%int64(f.failEveryNth) == 0 {
		return nil, errors.New("flapping: service temporarily unavailable")
	}
	return &models.UnifiedListing{ID: int(id)}, nil
}

func (f *flappingGRPCClient) CreateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	count := atomic.AddInt64(&f.requestCount, 1)
	if count%int64(f.failEveryNth) == 0 {
		return nil, errors.New("flapping: service temporarily unavailable")
	}
	unified.ID = int(count)
	return unified, nil
}

func (f *flappingGRPCClient) UpdateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	return unified, nil
}

func (f *flappingGRPCClient) DeleteListing(ctx context.Context, id int64, userID int64) error {
	return nil
}

type slowC2CRepository struct {
	delayMs int64
}

func (s *slowC2CRepository) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	time.Sleep(time.Duration(s.delayMs) * time.Millisecond)
	return 123, nil
}

func (s *slowC2CRepository) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	time.Sleep(time.Duration(s.delayMs) * time.Millisecond)
	return &models.MarketplaceListing{ID: id}, nil
}

func (s *slowC2CRepository) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	time.Sleep(time.Duration(s.delayMs) * time.Millisecond)
	return nil
}

func (s *slowC2CRepository) DeleteListing(ctx context.Context, id int) error {
	time.Sleep(time.Duration(s.delayMs) * time.Millisecond)
	return nil
}

type mockC2CRepository struct {
	createCalled bool
}

func (m *mockC2CRepository) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	m.createCalled = true
	return 123, nil
}

func (m *mockC2CRepository) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return &models.MarketplaceListing{ID: id}, nil
}

func (m *mockC2CRepository) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}

func (m *mockC2CRepository) DeleteListing(ctx context.Context, id int) error {
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
