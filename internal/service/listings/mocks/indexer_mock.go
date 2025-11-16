package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/sveturs/listings/internal/domain"
)

// MockIndexingService is a mock implementation of listings.IndexingService interface
type MockIndexingService struct {
	mock.Mock
}

// IndexListing mocks indexing a listing
func (m *MockIndexingService) IndexListing(ctx context.Context, listing *domain.Listing) error {
	args := m.Called(ctx, listing)
	return args.Error(0)
}

// UpdateListing mocks updating a listing in the index
func (m *MockIndexingService) UpdateListing(ctx context.Context, listing *domain.Listing) error {
	args := m.Called(ctx, listing)
	return args.Error(0)
}

// DeleteListing mocks deleting a listing from the index
func (m *MockIndexingService) DeleteListing(ctx context.Context, listingID int64) error {
	args := m.Called(ctx, listingID)
	return args.Error(0)
}

// GetSimilarListings mocks getting similar listings
func (m *MockIndexingService) GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]*domain.Listing, int32, error) {
	args := m.Called(ctx, listingID, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Listing), args.Get(1).(int32), args.Error(2)
}
