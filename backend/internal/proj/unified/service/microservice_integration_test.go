// backend/internal/proj/unified/service/microservice_integration_test.go
package service

import (
	"context"
	"testing"

	"backend/internal/domain/models"
	"backend/internal/domain/search"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockC2CRepository - простой mock для C2C репозитория
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
	return &models.MarketplaceListing{ID: id}, nil
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

// mockOpenSearchRepository - простой mock для OpenSearch
type mockOpenSearchRepository struct{}

func (m *mockOpenSearchRepository) SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	// For test purposes just return empty result
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

// TestMarketplaceService_WithMicroservice проверяет что service корректно использует микросервис
func TestMarketplaceService_WithMicroservice(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	t.Run("CreateListing uses microservice when enabled", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()
		service := NewMarketplaceService(
			mockC2C,
			nil, // b2cRepo not needed for this test
			&mockOpenSearchRepository{},
			logger,
		)

		// Enable microservice
		service.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "Test Listing",
			Price:      99.99,
			CategoryID: 1,
		}

		// Act
		id, err := service.CreateListing(context.Background(), unified)

		// Assert
		require.NoError(t, err)
		assert.Greater(t, id, int64(0), "Should return valid ID")
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should call gRPC CreateListing")
		assert.False(t, mockC2C.createCalled, "Should NOT call local C2C repo")

		// Verify listing was created in mock
		listings := mockGRPC.GetAllListings()
		assert.Len(t, listings, 1)
		assert.Equal(t, "Test Listing", listings[0].Title)
	})

	t.Run("CreateListing falls back to local DB when microservice fails", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()
		mockGRPC.ShouldFailCreate = true // Simulate microservice failure

		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)
		service.SetListingsGRPCClient(mockGRPC, true)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "Test Listing",
			Price:      99.99,
			CategoryID: 1,
		}

		// Act
		id, err := service.CreateListing(context.Background(), unified)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should attempt gRPC call")
		assert.True(t, mockC2C.createCalled, "Should fallback to local DB")
		assert.Equal(t, int64(123), id, "Should return local DB ID")
	})

	t.Run("GetListing uses microservice when enabled", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()

		// Pre-populate mock with a listing
		_, _ = mockGRPC.CreateListing(context.Background(), &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "Existing Listing",
			Price:      50.00,
			CategoryID: 1,
		})

		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)
		service.SetListingsGRPCClient(mockGRPC, true)

		// Act
		listing, err := service.GetListing(context.Background(), 1, SourceTypeC2C)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, listing)
		assert.Equal(t, "Existing Listing", listing.Title)
		assert.Equal(t, 1, mockGRPC.GetCallCount, "Should call gRPC GetListing")
		assert.False(t, mockC2C.getCalled, "Should NOT call local C2C repo")
	})

	t.Run("GetListing falls back to local DB when microservice fails", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()
		mockGRPC.ShouldFailGet = true

		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)
		service.SetListingsGRPCClient(mockGRPC, true)

		// Act
		listing, err := service.GetListing(context.Background(), 1, SourceTypeC2C)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, listing)
		assert.Equal(t, 1, mockGRPC.GetCallCount, "Should attempt gRPC call")
		assert.True(t, mockC2C.getCalled, "Should fallback to local DB")
	})

	t.Run("UpdateListing uses microservice when enabled", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()

		// Pre-populate mock
		created, _ := mockGRPC.CreateListing(context.Background(), &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "Old Title",
			Price:      50.00,
			CategoryID: 1,
		})

		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)
		service.SetListingsGRPCClient(mockGRPC, true)

		// Modify listing
		created.Title = "New Title"

		// Act
		err := service.UpdateListing(context.Background(), created)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, mockGRPC.UpdateCallCount, "Should call gRPC UpdateListing")
		assert.False(t, mockC2C.updateCalled, "Should NOT call local C2C repo")

		// Verify update
		updated, _ := mockGRPC.GetListing(context.Background(), int64(created.ID))
		assert.Equal(t, "New Title", updated.Title)
	})

	t.Run("DeleteListing uses microservice when enabled", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()

		// Pre-populate mock
		created, _ := mockGRPC.CreateListing(context.Background(), &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "To Delete",
			Price:      50.00,
			CategoryID: 1,
		})

		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)
		service.SetListingsGRPCClient(mockGRPC, true)

		// Act
		err := service.DeleteListing(context.Background(), int64(created.ID), SourceTypeC2C)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, mockGRPC.GetCallCount, "Should get listing first")
		assert.Equal(t, 1, mockGRPC.DeleteCallCount, "Should call gRPC DeleteListing")
		assert.False(t, mockC2C.deleteCalled, "Should NOT call local C2C repo")

		// Verify deletion
		listings := mockGRPC.GetAllListings()
		assert.Len(t, listings, 0, "Listing should be deleted")
	})
}

// TestMarketplaceService_WithoutMicroservice проверяет что service работает с локальной БД когда микросервис выключен
func TestMarketplaceService_WithoutMicroservice(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	t.Run("CreateListing uses local DB when microservice disabled", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()
		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)

		// Microservice is NOT enabled (default)
		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "Test Listing",
			Price:      99.99,
			CategoryID: 1,
		}

		// Act
		id, err := service.CreateListing(context.Background(), unified)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, int64(123), id)
		assert.Equal(t, 0, mockGRPC.CreateCallCount, "Should NOT call gRPC")
		assert.True(t, mockC2C.createCalled, "Should call local C2C repo")
	})

	t.Run("GetListing uses local DB when microservice disabled", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()
		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)

		// Act
		listing, err := service.GetListing(context.Background(), 1, SourceTypeC2C)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, listing)
		assert.Equal(t, 0, mockGRPC.GetCallCount, "Should NOT call gRPC")
		assert.True(t, mockC2C.getCalled, "Should call local C2C repo")
	})
}

// TestMarketplaceService_FeatureFlagToggle проверяет переключение feature flag
func TestMarketplaceService_FeatureFlagToggle(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	t.Run("Can toggle microservice on and off", func(t *testing.T) {
		// Arrange
		mockC2C := &mockC2CRepository{}
		mockGRPC := NewMockListingsGRPCClient()
		service := NewMarketplaceService(
			mockC2C,
			nil,
			&mockOpenSearchRepository{},
			logger,
		)

		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: SourceTypeC2C,
			Title:      "Test",
			Price:      50.00,
			CategoryID: 1,
		}

		// Initially uses local DB
		_, _ = service.CreateListing(context.Background(), unified)
		assert.True(t, mockC2C.createCalled, "Should use local DB initially")
		assert.Equal(t, 0, mockGRPC.CreateCallCount)

		// Reset
		mockC2C.createCalled = false

		// Enable microservice
		service.SetListingsGRPCClient(mockGRPC, true)
		_, _ = service.CreateListing(context.Background(), unified)
		assert.False(t, mockC2C.createCalled, "Should NOT use local DB when enabled")
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should use microservice")

		// Disable microservice
		service.SetListingsGRPCClient(nil, false)
		_, _ = service.CreateListing(context.Background(), unified)
		assert.True(t, mockC2C.createCalled, "Should use local DB when disabled")
		assert.Equal(t, 1, mockGRPC.CreateCallCount, "Should NOT call microservice again")
	})
}
