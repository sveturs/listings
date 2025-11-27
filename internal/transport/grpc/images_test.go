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

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
	minioclient "github.com/vondi-global/listings/internal/repository/minio"
	"github.com/vondi-global/listings/internal/service/listings"
)

// testServer wraps Server for testing with interface-based dependencies
type testServer struct {
	pb.UnimplementedListingsServiceServer
	service     listingsServiceInterface
	minioClient minioClientInterface
	logger      zerolog.Logger
}

// listingsServiceInterface defines methods needed for DeleteListingImage
type listingsServiceInterface interface {
	GetListing(ctx context.Context, id int64) (*domain.Listing, error)
	GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error)
	DeleteImage(ctx context.Context, imageID int64) error
}

// minioClientInterface defines MinIO operations needed for testing
type minioClientInterface interface {
	DeleteImage(ctx context.Context, objectName string) error
}

// DeleteListingImage implements the RPC method for testServer
func (s *testServer) DeleteListingImage(ctx context.Context, req *pb.DeleteListingImageRequest) (*pb.DeleteListingImageResponse, error) {
	// Use the same implementation but with interface-based dependencies
	s.logger.Info().
		Int64("listing_id", req.ListingId).
		Int64("image_id", req.ImageId).
		Int64("user_id", req.UserId).
		Msg("DeleteListingImage called")

	if req.ListingId <= 0 {
		s.logger.Warn().Int64("listing_id", req.ListingId).Msg("invalid listing_id")
		return nil, status.Error(codes.InvalidArgument, "listing_id must be positive")
	}

	if req.ImageId <= 0 {
		s.logger.Warn().Int64("image_id", req.ImageId).Msg("invalid image_id")
		return nil, status.Error(codes.InvalidArgument, "image_id must be positive")
	}

	if req.UserId <= 0 {
		s.logger.Warn().Int64("user_id", req.UserId).Msg("invalid user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

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

	image, err := s.service.GetImageByID(ctx, req.ImageId)
	if err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("image not found")
		return nil, status.Error(codes.NotFound, "image not found")
	}

	if image.ListingID != req.ListingId {
		s.logger.Warn().
			Int64("image_id", req.ImageId).
			Int64("image_listing_id", image.ListingID).
			Int64("requested_listing_id", req.ListingId).
			Msg("image does not belong to this listing")
		return nil, status.Error(codes.InvalidArgument, "image does not belong to this listing")
	}

	if s.minioClient == nil {
		s.logger.Error().Msg("MinIO client not configured")
		return nil, status.Error(codes.Internal, "storage system not available")
	}

	var deletedFromMinio bool
	var minioError error

	if image.StoragePath != nil && *image.StoragePath != "" {
		originalKey := *image.StoragePath

		if err := s.minioClient.DeleteImage(ctx, originalKey); err != nil {
			s.logger.Error().Err(err).Str("key", originalKey).Msg("failed to delete original image from MinIO")
			minioError = err
		} else {
			s.logger.Debug().Str("key", originalKey).Msg("original image deleted from MinIO")
			deletedFromMinio = true
		}

		thumbnailKey := getThumbnailPath(originalKey)
		if err := s.minioClient.DeleteImage(ctx, thumbnailKey); err != nil {
			s.logger.Warn().Err(err).Str("key", thumbnailKey).Msg("failed to delete thumbnail from MinIO (non-critical)")
		} else {
			s.logger.Debug().Str("key", thumbnailKey).Msg("thumbnail deleted from MinIO")
		}
	} else {
		s.logger.Warn().Int64("image_id", req.ImageId).Msg("no storage_path found, skipping MinIO deletion")
	}

	if err := s.service.DeleteImage(ctx, req.ImageId); err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("failed to delete image from database")

		if deletedFromMinio {
			s.logger.Error().
				Int64("image_id", req.ImageId).
				Str("storage_path", func() string {
					if image.StoragePath != nil {
						return *image.StoragePath
					}
					return ""
				}()).
				Msg("ORPHANED FILE IN MINIO: DB deletion failed after MinIO deletion succeeded")
		}

		return nil, status.Error(codes.Internal, "failed to delete image from database")
	}

	s.logger.Debug().Int64("image_id", req.ImageId).Msg("image deleted from database")

	var responseMessage string
	if minioError != nil {
		responseMessage = "Image deleted from database, but storage cleanup failed (files may be orphaned)"
		s.logger.Warn().
			Err(minioError).
			Int64("image_id", req.ImageId).
			Msg("partial success: DB deleted but MinIO cleanup failed")
	} else {
		responseMessage = "Image deleted successfully"
	}

	s.logger.Info().
		Int64("image_id", req.ImageId).
		Int64("listing_id", req.ListingId).
		Bool("minio_deleted", deletedFromMinio).
		Msg("DeleteListingImage completed")

	return &pb.DeleteListingImageResponse{
		Success: true,
		Message: responseMessage,
	}, nil
}

// Helper to create a test server with mocks
func setupTestServerWithMinio() (*testServer, *MockListingsService, *MockMinioClient) {
	mockService := new(MockListingsService)
	mockMinio := new(MockMinioClient)
	logger := zerolog.Nop()

	server := &testServer{
		service:     mockService,
		minioClient: mockMinio,
		logger:      logger,
	}

	return server, mockService, mockMinio
}

// Ensure MockListingsService implements listingsServiceInterface
var _ listingsServiceInterface = (*MockListingsService)(nil)

// Ensure *listings.Service implements listingsServiceInterface
var _ listingsServiceInterface = (*listings.Service)(nil)

// Ensure *minioclient.Client implements minioClientInterface
var _ minioClientInterface = (*minioclient.Client)(nil)

// MockMinioClient is a mock for MinIO operations
type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) DeleteImage(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

// TestDeleteListingImage_Success tests successful deletion
func TestDeleteListingImage_Success(t *testing.T) {
	ctx := context.Background()
	server, mockService, mockMinio := setupTestServerWithMinio()

	listingID := int64(123)
	imageID := int64(456)
	userID := int64(789)
	storagePath := "listings/123/1638360000_abc123.jpg"

	req := &pb.DeleteListingImageRequest{
		ListingId: listingID,
		ImageId:   imageID,
		UserId:    userID,
	}

	// Mock listing ownership check
	mockListing := &domain.Listing{
		ID:     listingID,
		UserID: userID,
	}
	mockService.On("GetListing", ctx, listingID).Return(mockListing, nil)

	// Mock image retrieval
	mockImage := &domain.ListingImage{
		ID:          imageID,
		ListingID:   listingID,
		StoragePath: &storagePath,
	}
	mockService.On("GetImageByID", ctx, imageID).Return(mockImage, nil)

	// Mock MinIO deletions (original + thumbnail)
	mockMinio.On("DeleteImage", ctx, storagePath).Return(nil)
	mockMinio.On("DeleteImage", ctx, "listings/123/1638360000_abc123_thumb.jpg").Return(nil)

	// Mock database deletion
	mockService.On("DeleteImage", ctx, imageID).Return(nil)

	// Execute
	resp, err := server.DeleteListingImage(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, "Image deleted successfully", resp.Message)

	mockService.AssertExpectations(t)
	mockMinio.AssertExpectations(t)
}

// TestDeleteListingImage_InvalidListingID tests validation
func TestDeleteListingImage_InvalidListingID(t *testing.T) {
	ctx := context.Background()
	server, _, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 0, // Invalid
		ImageId:   456,
		UserId:    789,
	}

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "listing_id must be positive")
}

// TestDeleteListingImage_InvalidImageID tests validation
func TestDeleteListingImage_InvalidImageID(t *testing.T) {
	ctx := context.Background()
	server, _, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   0, // Invalid
		UserId:    789,
	}

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "image_id must be positive")
}

// TestDeleteListingImage_InvalidUserID tests validation
func TestDeleteListingImage_InvalidUserID(t *testing.T) {
	ctx := context.Background()
	server, _, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   456,
		UserId:    0, // Invalid
	}

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "user_id must be positive")
}

// TestDeleteListingImage_ListingNotFound tests authorization check
func TestDeleteListingImage_ListingNotFound(t *testing.T) {
	ctx := context.Background()
	server, mockService, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   456,
		UserId:    789,
	}

	mockService.On("GetListing", ctx, req.ListingId).Return(nil, errors.New("not found"))

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "listing not found")

	mockService.AssertExpectations(t)
}

// TestDeleteListingImage_PermissionDenied tests ownership check
func TestDeleteListingImage_PermissionDenied(t *testing.T) {
	ctx := context.Background()
	server, mockService, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   456,
		UserId:    789,
	}

	// Listing belongs to different user
	mockListing := &domain.Listing{
		ID:     req.ListingId,
		UserID: 999, // Different user!
	}
	mockService.On("GetListing", ctx, req.ListingId).Return(mockListing, nil)

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())
	assert.Contains(t, st.Message(), "you do not own this listing")

	mockService.AssertExpectations(t)
}

// TestDeleteListingImage_ImageNotFound tests image retrieval
func TestDeleteListingImage_ImageNotFound(t *testing.T) {
	ctx := context.Background()
	server, mockService, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   456,
		UserId:    789,
	}

	mockListing := &domain.Listing{
		ID:     req.ListingId,
		UserID: req.UserId,
	}
	mockService.On("GetListing", ctx, req.ListingId).Return(mockListing, nil)
	mockService.On("GetImageByID", ctx, req.ImageId).Return(nil, errors.New("not found"))

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "image not found")

	mockService.AssertExpectations(t)
}

// TestDeleteListingImage_ImageBelongsToDifferentListing tests image-listing relationship
func TestDeleteListingImage_ImageBelongsToDifferentListing(t *testing.T) {
	ctx := context.Background()
	server, mockService, _ := setupTestServerWithMinio()

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   456,
		UserId:    789,
	}

	mockListing := &domain.Listing{
		ID:     req.ListingId,
		UserID: req.UserId,
	}
	mockService.On("GetListing", ctx, req.ListingId).Return(mockListing, nil)

	// Image belongs to different listing
	storagePath := "listings/999/image.jpg"
	mockImage := &domain.ListingImage{
		ID:          req.ImageId,
		ListingID:   999, // Different listing!
		StoragePath: &storagePath,
	}
	mockService.On("GetImageByID", ctx, req.ImageId).Return(mockImage, nil)

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "image does not belong to this listing")

	mockService.AssertExpectations(t)
}

// TestDeleteListingImage_MinioClientNotConfigured tests MinIO availability
func TestDeleteListingImage_MinioClientNotConfigured(t *testing.T) {
	ctx := context.Background()
	mockService := new(MockListingsService)
	logger := zerolog.Nop()

	server := &testServer{
		service:     mockService,
		minioClient: nil, // No MinIO client
		logger:      logger,
	}

	req := &pb.DeleteListingImageRequest{
		ListingId: 123,
		ImageId:   456,
		UserId:    789,
	}

	mockListing := &domain.Listing{
		ID:     req.ListingId,
		UserID: req.UserId,
	}
	mockService.On("GetListing", ctx, req.ListingId).Return(mockListing, nil)

	storagePath := "listings/123/image.jpg"
	mockImage := &domain.ListingImage{
		ID:          req.ImageId,
		ListingID:   req.ListingId,
		StoragePath: &storagePath,
	}
	mockService.On("GetImageByID", ctx, req.ImageId).Return(mockImage, nil)

	resp, err := server.DeleteListingImage(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "storage system not available")

	mockService.AssertExpectations(t)
}

// TestDeleteListingImage_MinioFailure_DBSuccess tests compensating transaction
func TestDeleteListingImage_MinioFailure_DBSuccess(t *testing.T) {
	ctx := context.Background()
	server, mockService, mockMinio := setupTestServerWithMinio()

	listingID := int64(123)
	imageID := int64(456)
	userID := int64(789)
	storagePath := "listings/123/1638360000_abc123.jpg"

	req := &pb.DeleteListingImageRequest{
		ListingId: listingID,
		ImageId:   imageID,
		UserId:    userID,
	}

	mockListing := &domain.Listing{
		ID:     listingID,
		UserID: userID,
	}
	mockService.On("GetListing", ctx, listingID).Return(mockListing, nil)

	mockImage := &domain.ListingImage{
		ID:          imageID,
		ListingID:   listingID,
		StoragePath: &storagePath,
	}
	mockService.On("GetImageByID", ctx, imageID).Return(mockImage, nil)

	// MinIO original deletion FAILS
	mockMinio.On("DeleteImage", ctx, storagePath).Return(errors.New("minio error"))
	// Thumbnail deletion is attempted but not critical
	mockMinio.On("DeleteImage", ctx, "listings/123/1638360000_abc123_thumb.jpg").Return(errors.New("minio error"))

	// DB deletion SUCCEEDS
	mockService.On("DeleteImage", ctx, imageID).Return(nil)

	// Execute
	resp, err := server.DeleteListingImage(ctx, req)

	// Assert: Should succeed but with warning message
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Contains(t, resp.Message, "storage cleanup failed")

	mockService.AssertExpectations(t)
	mockMinio.AssertExpectations(t)
}

// TestDeleteListingImage_MinioSuccess_DBFailure tests orphaned files scenario
func TestDeleteListingImage_MinioSuccess_DBFailure(t *testing.T) {
	ctx := context.Background()
	server, mockService, mockMinio := setupTestServerWithMinio()

	listingID := int64(123)
	imageID := int64(456)
	userID := int64(789)
	storagePath := "listings/123/1638360000_abc123.jpg"

	req := &pb.DeleteListingImageRequest{
		ListingId: listingID,
		ImageId:   imageID,
		UserId:    userID,
	}

	mockListing := &domain.Listing{
		ID:     listingID,
		UserID: userID,
	}
	mockService.On("GetListing", ctx, listingID).Return(mockListing, nil)

	mockImage := &domain.ListingImage{
		ID:          imageID,
		ListingID:   listingID,
		StoragePath: &storagePath,
	}
	mockService.On("GetImageByID", ctx, imageID).Return(mockImage, nil)

	// MinIO deletions SUCCEED
	mockMinio.On("DeleteImage", ctx, storagePath).Return(nil)
	mockMinio.On("DeleteImage", ctx, "listings/123/1638360000_abc123_thumb.jpg").Return(nil)

	// DB deletion FAILS (orphaned files in MinIO)
	mockService.On("DeleteImage", ctx, imageID).Return(errors.New("database error"))

	// Execute
	resp, err := server.DeleteListingImage(ctx, req)

	// Assert: Should fail with Internal error
	assert.Nil(t, resp)
	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to delete image from database")

	mockService.AssertExpectations(t)
	mockMinio.AssertExpectations(t)
}

// TestDeleteListingImage_NoStoragePath tests image without storage path
func TestDeleteListingImage_NoStoragePath(t *testing.T) {
	ctx := context.Background()
	server, mockService, _ := setupTestServerWithMinio()

	listingID := int64(123)
	imageID := int64(456)
	userID := int64(789)

	req := &pb.DeleteListingImageRequest{
		ListingId: listingID,
		ImageId:   imageID,
		UserId:    userID,
	}

	mockListing := &domain.Listing{
		ID:     listingID,
		UserID: userID,
	}
	mockService.On("GetListing", ctx, listingID).Return(mockListing, nil)

	// Image with no storage path
	mockImage := &domain.ListingImage{
		ID:          imageID,
		ListingID:   listingID,
		StoragePath: nil, // No storage path
	}
	mockService.On("GetImageByID", ctx, imageID).Return(mockImage, nil)

	// DB deletion should still work
	mockService.On("DeleteImage", ctx, imageID).Return(nil)

	// Execute
	resp, err := server.DeleteListingImage(ctx, req)

	// Assert: Should succeed (MinIO deletion skipped)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}

// TestGetThumbnailPath tests thumbnail path generation
func TestGetThumbnailPath(t *testing.T) {
	tests := []struct {
		name         string
		originalPath string
		expectedPath string
	}{
		{
			name:         "jpg extension",
			originalPath: "listings/123/1638360000_abc123.jpg",
			expectedPath: "listings/123/1638360000_abc123_thumb.jpg",
		},
		{
			name:         "jpeg extension",
			originalPath: "listings/456/1638360000_def456.jpeg",
			expectedPath: "listings/456/1638360000_def456_thumb.jpg",
		},
		{
			name:         "png extension",
			originalPath: "listings/789/1638360000_ghi789.png",
			expectedPath: "listings/789/1638360000_ghi789_thumb.jpg",
		},
		{
			name:         "webp extension",
			originalPath: "listings/101/1638360000_jkl101.webp",
			expectedPath: "listings/101/1638360000_jkl101_thumb.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getThumbnailPath(tt.originalPath)
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}
