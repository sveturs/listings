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
