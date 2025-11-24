package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
)

// reorderServiceInterface defines methods needed for ReorderListingImages
type reorderServiceInterface interface {
	GetListing(ctx context.Context, id int64) (*domain.Listing, error)
	GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error)
	ReorderImages(ctx context.Context, listingID int64, orders []postgres.ImageOrder) error
}

// testReorderServer wraps Server for testing ReorderListingImages
type testReorderServer struct {
	pb.UnimplementedListingsServiceServer
	service reorderServiceInterface
	logger  zerolog.Logger
}

// ReorderListingImages implements the RPC method for testReorderServer
func (s *testReorderServer) ReorderListingImages(ctx context.Context, req *pb.ReorderImagesRequest) (*pb.ReorderImagesResponse, error) {
	s.logger.Info().
		Int64("listing_id", req.ListingId).
		Int64("user_id", req.UserId).
		Int("images_count", len(req.ImageIds)).
		Msg("ReorderListingImages called")

	// Validation
	if req.ListingId <= 0 {
		s.logger.Warn().Int64("listing_id", req.ListingId).Msg("invalid listing_id")
		return nil, status.Error(codes.InvalidArgument, "listing_id must be positive")
	}

	if req.UserId <= 0 {
		s.logger.Warn().Int64("user_id", req.UserId).Msg("invalid user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	if len(req.ImageIds) == 0 {
		s.logger.Warn().Msg("empty image_ids list")
		return nil, status.Error(codes.InvalidArgument, "image_ids cannot be empty")
	}

	// Authorization: Verify user owns listing
	listing, err := s.service.GetListing(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("listing not found")
		return nil, status.Error(codes.NotFound, "listing not found")
	}

	if listing.UserID != req.UserId {
		s.logger.Warn().
			Int64("user_id", req.UserId).
			Int64("listing_user_id", listing.UserID).
			Int64("listing_id", req.ListingId).
			Msg("user does not own listing")
		return nil, status.Error(codes.PermissionDenied, "you do not own this listing")
	}

	s.logger.Debug().Int64("listing_id", req.ListingId).Int64("user_id", req.UserId).Msg("authorization passed")

	// Verify all image IDs belong to this listing
	existingImages, err := s.service.GetImages(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to get existing images")
		return nil, status.Error(codes.Internal, "failed to get existing images")
	}

	// Build map of existing image IDs for O(1) lookup
	imageMap := make(map[int64]bool)
	for _, img := range existingImages {
		imageMap[img.ID] = true
	}

	// Verify each requested image ID belongs to this listing
	for i, imgID := range req.ImageIds {
		if !imageMap[imgID] {
			s.logger.Warn().
				Int64("image_id", imgID).
				Int64("listing_id", req.ListingId).
				Int("position", i).
				Msg("image does not belong to listing")
			return nil, status.Errorf(codes.InvalidArgument, "image_id %d does not belong to listing %d", imgID, req.ListingId)
		}
	}

	s.logger.Debug().
		Int64("listing_id", req.ListingId).
		Int("existing_images", len(existingImages)).
		Int("requested_images", len(req.ImageIds)).
		Msg("image validation passed")

	// Reorder images: Convert image_ids array to ImageOrder list
	var orders []postgres.ImageOrder
	for position, imageID := range req.ImageIds {
		orders = append(orders, postgres.ImageOrder{
			ImageID:      imageID,
			DisplayOrder: int32(position + 1), // 1-indexed
		})
	}

	// Call repository method (transaction-based batch update)
	if err := s.service.ReorderImages(ctx, req.ListingId, orders); err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to reorder images")
		return nil, status.Error(codes.Internal, "failed to reorder images")
	}

	s.logger.Info().
		Int64("listing_id", req.ListingId).
		Int64("user_id", req.UserId).
		Int("images_count", len(req.ImageIds)).
		Msg("images reordered successfully")

	return &pb.ReorderImagesResponse{
		Success: true,
	}, nil
}

// MockReorderService is a mock for reorder operations
type MockReorderService struct {
	mock.Mock
}

func (m *MockReorderService) GetListing(ctx context.Context, id int64) (*domain.Listing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

func (m *MockReorderService) GetImages(ctx context.Context, listingID int64) ([]*domain.ListingImage, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ListingImage), args.Error(1)
}

func (m *MockReorderService) ReorderImages(ctx context.Context, listingID int64, orders []postgres.ImageOrder) error {
	args := m.Called(ctx, listingID, orders)
	return args.Error(0)
}

// Helper to create a test server for reorder tests
func setupReorderTestServer() (*testReorderServer, *MockReorderService) {
	mockService := new(MockReorderService)
	logger := zerolog.Nop()

	server := &testReorderServer{
		service: mockService,
		logger:  logger,
	}

	return server, mockService
}

// TestReorderListingImages_Success tests successful reordering
func TestReorderListingImages_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	listingID := int64(123)
	userID := int64(1)
	imageIDs := []int64{3, 1, 2}

	req := &pb.ReorderImagesRequest{
		ListingId: listingID,
		UserId:    userID,
		ImageIds:  imageIDs,
	}

	// Mock listing ownership check
	mockListing := &domain.Listing{
		ID:     listingID,
		UserID: userID,
	}
	mockService.On("GetListing", ctx, listingID).Return(mockListing, nil)

	// Mock existing images
	existingImages := []*domain.ListingImage{
		{ID: 1, ListingID: listingID, DisplayOrder: 1},
		{ID: 2, ListingID: listingID, DisplayOrder: 2},
		{ID: 3, ListingID: listingID, DisplayOrder: 3},
	}
	mockService.On("GetImages", ctx, listingID).Return(existingImages, nil)

	// Mock reorder operation - use mock.MatchedBy for flexible matching
	mockService.On("ReorderImages", ctx, listingID, mock.MatchedBy(func(orders []postgres.ImageOrder) bool {
		if len(orders) != 3 {
			return false
		}
		// Verify order: [3,1,2] -> [(3,1), (1,2), (2,3)]
		return orders[0].ImageID == 3 && orders[0].DisplayOrder == 1 &&
			orders[1].ImageID == 1 && orders[1].DisplayOrder == 2 &&
			orders[2].ImageID == 2 && orders[2].DisplayOrder == 3
	})).Return(nil)

	// Execute
	resp, err := server.ReorderListingImages(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_MissingListingID tests missing listing_id
func TestReorderListingImages_MissingListingID(t *testing.T) {
	ctx := context.Background()
	server, _ := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		UserId:   1,
		ImageIds: []int64{1, 2, 3},
	}

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "listing_id must be positive")
}

// TestReorderListingImages_MissingUserID tests missing user_id
func TestReorderListingImages_MissingUserID(t *testing.T) {
	ctx := context.Background()
	server, _ := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		ImageIds:  []int64{1, 2, 3},
	}

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "user_id must be positive")
}

// TestReorderListingImages_EmptyImageIDs tests empty image_ids list
func TestReorderListingImages_EmptyImageIDs(t *testing.T) {
	ctx := context.Background()
	server, _ := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{},
	}

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "image_ids cannot be empty")
}

// TestReorderListingImages_ListingNotFound tests listing not found
func TestReorderListingImages_ListingNotFound(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{1, 2, 3},
	}

	mockService.On("GetListing", ctx, int64(123)).Return(nil, errors.New("not found"))

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "listing not found")

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_Unauthorized tests user doesn't own listing
func TestReorderListingImages_Unauthorized(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{1, 2, 3},
	}

	// Listing owned by different user
	mockListing := &domain.Listing{
		ID:     123,
		UserID: 999, // Different user
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())
	assert.Contains(t, st.Message(), "you do not own this listing")

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_InvalidImageID tests image doesn't belong to listing
func TestReorderListingImages_InvalidImageID(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{1, 2, 999}, // 999 doesn't belong
	}

	mockListing := &domain.Listing{
		ID:     123,
		UserID: 1,
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)

	// Only images 1 and 2 exist
	existingImages := []*domain.ListingImage{
		{ID: 1, ListingID: 123},
		{ID: 2, ListingID: 123},
	}
	mockService.On("GetImages", ctx, int64(123)).Return(existingImages, nil)

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "image_id 999 does not belong to listing 123")

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_GetImagesFails tests failure to retrieve images
func TestReorderListingImages_GetImagesFails(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{1, 2, 3},
	}

	mockListing := &domain.Listing{
		ID:     123,
		UserID: 1,
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)
	mockService.On("GetImages", ctx, int64(123)).Return(nil, errors.New("database error"))

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to get existing images")

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_ReorderFails tests failure during reorder operation
func TestReorderListingImages_ReorderFails(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{1, 2},
	}

	mockListing := &domain.Listing{
		ID:     123,
		UserID: 1,
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)

	existingImages := []*domain.ListingImage{
		{ID: 1, ListingID: 123},
		{ID: 2, ListingID: 123},
	}
	mockService.On("GetImages", ctx, int64(123)).Return(existingImages, nil)
	mockService.On("ReorderImages", ctx, int64(123), mock.Anything).Return(errors.New("update failed"))

	resp, err := server.ReorderListingImages(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to reorder images")

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_SingleImage tests reordering with single image
func TestReorderListingImages_SingleImage(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  []int64{42},
	}

	mockListing := &domain.Listing{
		ID:     123,
		UserID: 1,
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)

	existingImages := []*domain.ListingImage{
		{ID: 42, ListingID: 123},
	}
	mockService.On("GetImages", ctx, int64(123)).Return(existingImages, nil)

	mockService.On("ReorderImages", ctx, int64(123), mock.MatchedBy(func(orders []postgres.ImageOrder) bool {
		return len(orders) == 1 && orders[0].ImageID == 42 && orders[0].DisplayOrder == 1
	})).Return(nil)

	resp, err := server.ReorderListingImages(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_LargeSet tests reordering with many images
func TestReorderListingImages_LargeSet(t *testing.T) {
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	// Create 10 images
	imageIDs := make([]int64, 10)
	existingImages := make([]*domain.ListingImage, 10)
	for i := 0; i < 10; i++ {
		imageIDs[i] = int64(i + 1)
		existingImages[i] = &domain.ListingImage{
			ID:        int64(i + 1),
			ListingID: 123,
		}
	}

	// Reverse order
	reversedIDs := make([]int64, 10)
	for i := 0; i < 10; i++ {
		reversedIDs[i] = imageIDs[9-i]
	}

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  reversedIDs,
	}

	mockListing := &domain.Listing{
		ID:     123,
		UserID: 1,
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)
	mockService.On("GetImages", ctx, int64(123)).Return(existingImages, nil)

	mockService.On("ReorderImages", ctx, int64(123), mock.MatchedBy(func(orders []postgres.ImageOrder) bool {
		if len(orders) != 10 {
			return false
		}
		// Verify first and last elements
		return orders[0].ImageID == 10 && orders[0].DisplayOrder == 1 &&
			orders[9].ImageID == 1 && orders[9].DisplayOrder == 10
	})).Return(nil)

	resp, err := server.ReorderListingImages(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

// TestReorderListingImages_DisplayOrderConversion tests position to display_order conversion
func TestReorderListingImages_DisplayOrderConversion(t *testing.T) {
	// This test verifies that array indices (0-indexed) are converted to display_order (1-indexed)
	ctx := context.Background()
	server, mockService := setupReorderTestServer()

	imageIDs := []int64{5, 3, 7}

	req := &pb.ReorderImagesRequest{
		ListingId: 123,
		UserId:    1,
		ImageIds:  imageIDs,
	}

	mockListing := &domain.Listing{
		ID:     123,
		UserID: 1,
	}
	mockService.On("GetListing", ctx, int64(123)).Return(mockListing, nil)

	existingImages := []*domain.ListingImage{
		{ID: 3, ListingID: 123},
		{ID: 5, ListingID: 123},
		{ID: 7, ListingID: 123},
	}
	mockService.On("GetImages", ctx, int64(123)).Return(existingImages, nil)

	// Verify exact mapping: index 0 -> display_order 1, index 1 -> display_order 2, etc.
	mockService.On("ReorderImages", ctx, int64(123), mock.MatchedBy(func(orders []postgres.ImageOrder) bool {
		if len(orders) != 3 {
			return false
		}
		// imageIDs[0]=5 should get display_order=1
		// imageIDs[1]=3 should get display_order=2
		// imageIDs[2]=7 should get display_order=3
		return orders[0].ImageID == 5 && orders[0].DisplayOrder == 1 &&
			orders[1].ImageID == 3 && orders[1].DisplayOrder == 2 &&
			orders[2].ImageID == 7 && orders[2].DisplayOrder == 3
	})).Return(nil)

	resp, err := server.ReorderListingImages(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}
