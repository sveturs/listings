package grpc

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
)

// DeleteListingImage implements TRUE MICROSERVICE pattern for image deletion
// Includes: Authorization, MinIO cleanup (original + thumbnail), DB deletion, compensating transactions
func (s *Server) DeleteListingImage(ctx context.Context, req *listingspb.DeleteListingImageRequest) (*listingspb.DeleteListingImageResponse, error) {
	s.logger.Info().
		Int64("listing_id", req.ListingId).
		Int64("image_id", req.ImageId).
		Int64("user_id", req.UserId).
		Msg("DeleteListingImage called")

	// ============================================================================
	// VALIDATION
	// ============================================================================

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

	// ============================================================================
	// AUTHORIZATION: Verify user owns listing
	// ============================================================================

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

	// ============================================================================
	// GET IMAGE RECORD: Need storage_path for MinIO deletion
	// ============================================================================

	image, err := s.service.GetImageByID(ctx, req.ImageId)
	if err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("image not found")
		return nil, status.Error(codes.NotFound, "image not found")
	}

	// Verify image belongs to this listing
	if image.ListingID != req.ListingId {
		s.logger.Warn().
			Int64("image_id", req.ImageId).
			Int64("image_listing_id", image.ListingID).
			Int64("requested_listing_id", req.ListingId).
			Msg("image does not belong to this listing")
		return nil, status.Error(codes.InvalidArgument, "image does not belong to this listing")
	}

	// ============================================================================
	// MINIO DELETION: Delete original + thumbnail
	// ============================================================================

	if s.minioClient == nil {
		s.logger.Error().Msg("MinIO client not configured")
		return nil, status.Error(codes.Internal, "storage system not available")
	}

	var deletedFromMinio bool
	var minioError error

	// Delete original image if storage_path exists
	if image.StoragePath != nil && *image.StoragePath != "" {
		originalKey := *image.StoragePath

		if err := s.minioClient.DeleteImage(ctx, originalKey); err != nil {
			s.logger.Error().Err(err).Str("key", originalKey).Msg("failed to delete original image from MinIO")
			minioError = fmt.Errorf("failed to delete original image: %w", err)
		} else {
			s.logger.Debug().Str("key", originalKey).Msg("original image deleted from MinIO")
			deletedFromMinio = true
		}

		// Generate thumbnail path and delete (Phase 24 pattern: "listings/<id>/<timestamp>_<uuid>_thumb.jpg")
		thumbnailKey := getThumbnailPath(originalKey)

		if err := s.minioClient.DeleteImage(ctx, thumbnailKey); err != nil {
			// Non-critical: thumbnail might not exist or deletion failed
			s.logger.Warn().Err(err).Str("key", thumbnailKey).Msg("failed to delete thumbnail from MinIO (non-critical)")
		} else {
			s.logger.Debug().Str("key", thumbnailKey).Msg("thumbnail deleted from MinIO")
		}
	} else {
		s.logger.Warn().Int64("image_id", req.ImageId).Msg("no storage_path found, skipping MinIO deletion")
	}

	// ============================================================================
	// DATABASE DELETION
	// ============================================================================

	if err := s.service.DeleteImage(ctx, req.ImageId); err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("failed to delete image from database")

		// COMPENSATING TRANSACTION: MinIO deletion succeeded but DB failed
		// In this case, we accept orphaned files in MinIO (better than inconsistent state)
		// Production systems should have cleanup job to remove orphaned MinIO files
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

	// ============================================================================
	// SUCCESS RESPONSE
	// ============================================================================

	var responseMessage string
	if minioError != nil {
		// DB deletion succeeded but MinIO failed (inconsistent state)
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

	return &listingspb.DeleteListingImageResponse{
		Success: true,
		Message: responseMessage,
	}, nil
}

// getThumbnailPath generates thumbnail path from original image path
// Input:  "listings/123/1638360000_abc123.jpg"
// Output: "listings/123/1638360000_abc123_thumb.jpg"
func getThumbnailPath(originalPath string) string {
	ext := filepath.Ext(originalPath)
	baseWithoutExt := strings.TrimSuffix(originalPath, ext)
	return baseWithoutExt + "_thumb.jpg" // Thumbnails are always JPEG (Phase 24 pattern)
}
