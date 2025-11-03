//go:build ignore
// +build ignore

// backend/tests/e2e/marketplace_microservice_e2e_test.go
// DEPRECATED: E2E tests for unified architecture that was removed in Phase 7
// These tests are kept for reference but disabled
// To enable: remove //go:build ignore directive and restore unified/service import
package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"

	// "backend/internal/proj/unified/service" // REMOVED: unified architecture deleted in Phase 7

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Skip all e2e tests - unified architecture was removed
const skipUnifiedE2ETests = true

// defaultRoutingContext создаёт дефолтный routing context для тестов
// NOTE: Commented out as unified/service was removed
// type RoutingContext struct {
//     UserID  int
//     IsAdmin bool
// }
// func defaultRoutingContext() *RoutingContext {
// 	return &RoutingContext{
// 		UserID:  100,
// 		IsAdmin: false,
// 	}
// }

// TestE2E_FullFlow_MonolithToMicroservice проверяет полный поток через monolith → microservice
func TestE2E_FullFlow_MonolithToMicroservice(t *testing.T) {
	t.Skip("DEPRECATED: Test disabled - unified architecture was removed in Phase 7")
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Create via monolith and Read via microservice", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

		// Enable microservice
		svc.SetListingsGRPCClient(mockGRPC, true)

		// Create listing
		unified := &models.UnifiedListing{
			UserID:      100,
			SourceType:  service.SourceTypeC2C,
			Title:       "E2E Test Listing",
			Description: "Created via monolith",
			Price:       99.99,
			CategoryID:  1,
		}

		// Act - Create
		id, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Act - Read
		retrieved, err := svc.GetListing(ctx, id, service.SourceTypeC2C, defaultRoutingContext())

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, "E2E Test Listing", retrieved.Title)
		assert.Equal(t, "Created via monolith", retrieved.Description)
		assert.Equal(t, 99.99, retrieved.Price)

		// Verify data landed in microservice (not in monolith DB)
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should create via microservice")
		assert.False(t, mockC2C.createCalled, "Should NOT use local DB")

		listings := mockGRPC.GetAllListings()
		assert.Len(t, listings, 1)
	})

	t.Run("Update flow via microservice", func(t *testing.T) {
		// Arrange
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		// Create initial listing
		created, err := mockGRPC.CreateListing(ctx, &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Original Title",
			Price:      50.00,
			CategoryID: 1,
		})
		require.NoError(t, err)

		// Modify
		created.Title = "Updated Title"
		created.Price = 75.00

		// Act - Update
		err = svc.UpdateListing(ctx, created, defaultRoutingContext())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, mockGRPC.UpdateCallCount)

		// Verify update persisted
		updated, err := svc.GetListing(ctx, int64(created.ID), service.SourceTypeC2C, defaultRoutingContext())
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, 75.00, updated.Price)
	})

	t.Run("Delete flow via microservice", func(t *testing.T) {
		// Arrange
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		// Create listing
		created, err := mockGRPC.CreateListing(ctx, &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "To Delete",
			Price:      10.00,
			CategoryID: 1,
		})
		require.NoError(t, err)

		// Act - Delete
		err = svc.DeleteListing(ctx, int64(created.ID), service.SourceTypeC2C, defaultRoutingContext())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, mockGRPC.DeleteCallCount)

		// Verify deletion
		listings := mockGRPC.GetAllListings()
		assert.Len(t, listings, 0, "Listing should be deleted")
	})
}

// TestE2E_FeatureFlag проверяет переключение feature flag в runtime
func TestE2E_FeatureFlag(t *testing.T) {
	t.Skip("DEPRECATED: Test disabled - unified architecture was removed in Phase 7")
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Toggle microservice on and off", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Feature Flag Test",
			Price:      50.00,
			CategoryID: 1,
		}

		// Act 1 - With microservice OFF (default)
		id1, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		require.NoError(t, err)
		assert.True(t, mockC2C.createCalled, "Should use local DB when microservice disabled")
		assert.Equal(t, 0, mockGRPC.CreateCallCount)

		// Reset flags
		mockC2C.createCalled = false

		// Act 2 - Enable microservice
		svc.SetListingsGRPCClient(mockGRPC, true)
		id2, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		require.NoError(t, err)
		assert.False(t, mockC2C.createCalled, "Should NOT use local DB when microservice enabled")
		assert.Equal(t, 1, mockGRPC.CreateCallCount)
		assert.NotEqual(t, id1, id2, "IDs should be different")

		// Act 3 - Disable microservice again
		svc.SetListingsGRPCClient(nil, false)
		_, err = svc.CreateListing(ctx, unified, defaultRoutingContext())
		require.NoError(t, err)
		assert.True(t, mockC2C.createCalled, "Should fallback to local DB when microservice disabled")
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should not increment microservice calls")
	})
}

// TestE2E_Fallback проверяет fallback на monolith при отказе microservice
func TestE2E_Fallback(t *testing.T) {
	t.Skip("DEPRECATED: Test disabled - unified architecture was removed in Phase 7")
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Fallback to monolith when microservice fails", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		mockGRPC.ShouldFailCreate = true // Simulate microservice failure

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Fallback Test",
			Price:      25.00,
			CategoryID: 1,
		}

		// Act
		id, err := svc.CreateListing(ctx, unified, defaultRoutingContext())

		// Assert - Should succeed via fallback
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should attempt microservice first")
		assert.True(t, mockC2C.createCalled, "Should fallback to local DB")
	})

	t.Run("Fallback for Get operation", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		mockGRPC.ShouldFailGet = true

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		// Act
		listing, err := svc.GetListing(ctx, 123, service.SourceTypeC2C, defaultRoutingContext())

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, listing)
		assert.Equal(t, 1, mockGRPC.GetCallCount, "Should attempt microservice first")
		assert.True(t, mockC2C.getCalled, "Should fallback to local DB")
	})

	t.Run("Fallback for Update operation", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		mockGRPC.ShouldFailUpdate = true

		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			ID:         100,
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Updated",
			Price:      99.99,
			CategoryID: 1,
		}

		// Act
		err := svc.UpdateListing(ctx, unified, defaultRoutingContext())

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, mockGRPC.UpdateCallCount, "Should attempt microservice first")
		assert.True(t, mockC2C.updateCalled, "Should fallback to local DB")
	})
}

// TestE2E_DataConsistency проверяет согласованность данных
func TestE2E_DataConsistency(t *testing.T) {
	t.Skip("DEPRECATED: Test disabled - unified architecture was removed in Phase 7")
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Data consistency across microservice and monolith", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

		// Create same listing in both systems
		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      "Consistency Test",
			Price:      50.00,
			CategoryID: 1,
		}

		// Create in local DB (microservice disabled)
		localID, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		require.NoError(t, err)

		// Create in microservice
		svc.SetListingsGRPCClient(mockGRPC, true)
		microID, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
		require.NoError(t, err)

		// Verify both have same data (different IDs is OK)
		assert.True(t, mockC2C.createCalled)
		assert.Equal(t, 1, mockGRPC.CreateCallCount)

		// Both should return same content
		svc.SetListingsGRPCClient(nil, false)
		localListing, _ := svc.GetListing(ctx, localID, service.SourceTypeC2C, defaultRoutingContext())

		svc.SetListingsGRPCClient(mockGRPC, true)
		microListing, _ := svc.GetListing(ctx, microID, service.SourceTypeC2C, defaultRoutingContext())

		assert.Equal(t, localListing.Title, microListing.Title)
		assert.Equal(t, localListing.Price, microListing.Price)
		assert.Equal(t, localListing.CategoryID, microListing.CategoryID)
	})
}

// TestE2E_ConcurrentOperations проверяет конкурентные операции
func TestE2E_ConcurrentOperations(t *testing.T) {
	t.Skip("DEPRECATED: Test disabled - unified architecture was removed in Phase 7")
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	t.Run("Concurrent creates via microservice", func(t *testing.T) {
		// Arrange
		mockGRPC := service.NewMockListingsGRPCClient()
		svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
		svc.SetListingsGRPCClient(mockGRPC, true)

		concurrency := 10
		done := make(chan bool, concurrency)
		errors := make(chan error, concurrency)

		// Act - Create multiple listings concurrently
		for i := 0; i < concurrency; i++ {
			go func(idx int) {
				unified := &models.UnifiedListing{
					UserID:     100,
					SourceType: service.SourceTypeC2C,
					Title:      fmt.Sprintf("Concurrent Test %d", idx),
					Price:      float64(idx * 10),
					CategoryID: 1,
				}

				_, err := svc.CreateListing(ctx, unified, defaultRoutingContext())
				if err != nil {
					errors <- err
				}
				done <- true
			}(i)
		}

		// Wait for completion
		for i := 0; i < concurrency; i++ {
			select {
			case <-done:
				// OK
			case err := <-errors:
				t.Errorf("Concurrent create failed: %v", err)
			case <-time.After(5 * time.Second):
				t.Fatal("Timeout waiting for concurrent creates")
			}
		}

		// Assert
		listings := mockGRPC.GetAllListings()
		assert.Len(t, listings, concurrency, "Should create all listings")
		assert.Equal(t, concurrency, mockGRPC.CreateCallCount)
	})
}

// Mock repositories for E2E tests

type mockC2CRepository struct {
	createCalled bool
	getCalled    bool
	updateCalled bool
	deleteCalled bool
	shouldFail   bool
}

func (m *mockC2CRepository) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	m.createCalled = true
	if m.shouldFail {
		return 0, assert.AnError
	}
	return 123, nil
}

func (m *mockC2CRepository) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	m.getCalled = true
	if m.shouldFail {
		return nil, assert.AnError
	}
	return &models.MarketplaceListing{ID: id, Title: "Mock Listing"}, nil
}

func (m *mockC2CRepository) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	m.updateCalled = true
	if m.shouldFail {
		return assert.AnError
	}
	return nil
}

func (m *mockC2CRepository) DeleteListing(ctx context.Context, id int) error {
	m.deleteCalled = true
	if m.shouldFail {
		return assert.AnError
	}
	return nil
}

type mockOpenSearchRepository struct{}

func (m *mockOpenSearchRepository) SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	return &search.ServiceResult{
		Items: []*models.MarketplaceListing{},
		Total: 0,
		Page:  0,
		Size:  10,
	}, nil
}

func (m *mockOpenSearchRepository) Index(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}

func (m *mockOpenSearchRepository) Delete(ctx context.Context, listingID int) error {
	return nil
}
