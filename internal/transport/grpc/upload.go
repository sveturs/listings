package grpc

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
)

const (
	maxImageSize     = 10 * 1024 * 1024 // 10MB per image
	maxTotalSize     = 50 * 1024 * 1024 // 50MB total per upload batch
	maxFiles         = 10               // Maximum files per upload
	thumbnailSize    = 200              // Thumbnail dimensions (200x200px)
	thumbnailQuality = 85               // JPEG quality for thumbnails
	chunkSize        = 1024 * 1024      // 1MB chunks
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadListingImages implements streaming image upload with authorization and thumbnail generation
func (s *Server) UploadListingImages(stream listingspb.ListingsService_UploadListingImagesServer) error {
	ctx := stream.Context()

	s.logger.Info().Msg("UploadListingImages: starting streaming upload")

	var (
		uploadedImages []*listingspb.ListingImage
		uploadErrors   []string
		uploadedCount  int32
		failedCount    int32
		totalSize      int64
		fileCount      int32

		// Current file state
		currentMetadata   *listingspb.UploadImageMetadata
		currentFileBuffer bytes.Buffer
		listingIDChecked  bool
	)

	// Receive streaming chunks
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Process last file if any
			if currentMetadata != nil {
				if img, uploadErr := s.processImageUpload(ctx, currentMetadata, &currentFileBuffer); uploadErr != nil {
					uploadErrors = append(uploadErrors, fmt.Sprintf("%s: %v", currentMetadata.Filename, uploadErr))
					failedCount++
					s.logger.Error().Err(uploadErr).Str("filename", currentMetadata.Filename).Msg("failed to upload image")
				} else {
					uploadedImages = append(uploadedImages, img)
					uploadedCount++
					s.logger.Info().Str("filename", currentMetadata.Filename).Int64("image_id", img.Id).Msg("image uploaded successfully")
				}
			}
			break
		}
		if err != nil {
			s.logger.Error().Err(err).Msg("UploadListingImages: error receiving chunk")
			return status.Errorf(codes.Internal, "failed to receive chunk: %v", err)
		}

		switch data := req.Data.(type) {
		case *listingspb.UploadImageChunkRequest_Metadata:
			// Process previous file if any
			if currentMetadata != nil {
				if img, uploadErr := s.processImageUpload(ctx, currentMetadata, &currentFileBuffer); uploadErr != nil {
					uploadErrors = append(uploadErrors, fmt.Sprintf("%s: %v", currentMetadata.Filename, uploadErr))
					failedCount++
					s.logger.Error().Err(uploadErr).Str("filename", currentMetadata.Filename).Msg("failed to upload image")
				} else {
					uploadedImages = append(uploadedImages, img)
					uploadedCount++
					s.logger.Info().Str("filename", currentMetadata.Filename).Int64("image_id", img.Id).Msg("image uploaded successfully")
				}
			}

			// Start new file
			currentMetadata = data.Metadata
			currentFileBuffer.Reset()
			fileCount++

			// PRE-FLIGHT VALIDATION: Max files
			if fileCount > maxFiles {
				s.logger.Warn().Int32("file_count", fileCount).Msg("too many files")
				return status.Errorf(codes.InvalidArgument, "too many files: maximum %d allowed", maxFiles)
			}

			// PRE-FLIGHT VALIDATION: Validate metadata
			if err := s.validateImageMetadata(currentMetadata); err != nil {
				s.logger.Error().Err(err).Str("filename", currentMetadata.Filename).Msg("invalid metadata")
				return status.Errorf(codes.InvalidArgument, "invalid metadata: %v", err)
			}

			// AUTHORIZATION CHECK: Verify user owns listing (only once per batch)
			if !listingIDChecked {
				listing, err := s.service.GetListing(ctx, currentMetadata.ListingId)
				if err != nil {
					s.logger.Error().Err(err).Int64("listing_id", currentMetadata.ListingId).Msg("listing not found")
					return status.Errorf(codes.NotFound, "listing not found")
				}

				// Check ownership
				if listing.UserID != currentMetadata.UserId {
					s.logger.Warn().
						Int64("user_id", currentMetadata.UserId).
						Int64("listing_user_id", listing.UserID).
						Int64("listing_id", currentMetadata.ListingId).
						Msg("user does not own listing")
					return status.Errorf(codes.PermissionDenied, "you do not own this listing")
				}

				listingIDChecked = true
				s.logger.Debug().Int64("listing_id", currentMetadata.ListingId).Int64("user_id", currentMetadata.UserId).Msg("authorization passed")
			}

			// Track total size
			totalSize += currentMetadata.FileSize
			if totalSize > maxTotalSize {
				s.logger.Warn().Int64("total_size_mb", totalSize/(1024*1024)).Msg("total size exceeds limit")
				return status.Errorf(codes.InvalidArgument, "total size exceeds %dMB limit", maxTotalSize/(1024*1024))
			}

			s.logger.Debug().
				Str("filename", currentMetadata.Filename).
				Int64("file_size", currentMetadata.FileSize).
				Int32("display_order", currentMetadata.DisplayOrder).
				Bool("is_primary", currentMetadata.IsPrimary).
				Msg("receiving file metadata")

		case *listingspb.UploadImageChunkRequest_Chunk:
			if currentMetadata == nil {
				s.logger.Error().Msg("received chunk without metadata")
				return status.Error(codes.InvalidArgument, "received chunk before metadata")
			}

			// Write chunk to buffer
			if _, err := currentFileBuffer.Write(data.Chunk); err != nil {
				s.logger.Error().Err(err).Msg("failed to write chunk to buffer")
				return status.Errorf(codes.Internal, "failed to write chunk: %v", err)
			}

			// Validate size doesn't exceed limit
			if int64(currentFileBuffer.Len()) > maxImageSize {
				s.logger.Warn().Int64("size_mb", int64(currentFileBuffer.Len())/(1024*1024)).Msg("file size exceeds limit")
				return status.Errorf(codes.InvalidArgument, "file size exceeds %dMB limit", maxImageSize/(1024*1024))
			}

			s.logger.Debug().
				Str("filename", currentMetadata.Filename).
				Int("chunk_size", len(data.Chunk)).
				Int("total_received", currentFileBuffer.Len()).
				Msg("received chunk")
		}
	}

	// Return response
	response := &listingspb.UploadImagesResponse{
		Images:        uploadedImages,
		UploadedCount: uploadedCount,
		FailedCount:   failedCount,
		Errors:        uploadErrors,
	}

	s.logger.Info().
		Int32("uploaded", uploadedCount).
		Int32("failed", failedCount).
		Int32("total", fileCount).
		Msg("UploadListingImages: completed")

	return stream.SendAndClose(response)
}

// processImageUpload processes a single image: decode, generate thumbnail, upload to MinIO, save to DB
func (s *Server) processImageUpload(ctx context.Context, metadata *listingspb.UploadImageMetadata, fileBuffer *bytes.Buffer) (*listingspb.ListingImage, error) {
	if s.minioClient == nil {
		return nil, fmt.Errorf("MinIO client not configured")
	}

	// Validate file is actually an image by decoding
	img, format, err := image.Decode(bytes.NewReader(fileBuffer.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("invalid image format: %w", err)
	}

	s.logger.Debug().
		Str("filename", metadata.Filename).
		Str("format", format).
		Int("width", img.Bounds().Dx()).
		Int("height", img.Bounds().Dy()).
		Msg("image decoded successfully")

	// Generate thumbnail (200x200px, maintain aspect ratio)
	thumbnailImg := resize.Thumbnail(thumbnailSize, thumbnailSize, img, resize.Lanczos3)
	var thumbnailBuf bytes.Buffer
	if err := jpeg.Encode(&thumbnailBuf, thumbnailImg, &jpeg.Options{Quality: thumbnailQuality}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	s.logger.Debug().
		Str("filename", metadata.Filename).
		Int("thumbnail_size", thumbnailBuf.Len()).
		Msg("thumbnail generated")

	// Generate MinIO object keys
	ext := strings.ToLower(filepath.Ext(metadata.Filename))
	if ext == "" {
		ext = "." + format // Use detected format if no extension
	}

	timestamp := time.Now().UnixNano()
	uniqueID := uuid.New().String()[:8]

	originalKey := fmt.Sprintf("listings/%d/%d_%s%s", metadata.ListingId, timestamp, uniqueID, ext)
	thumbnailKey := fmt.Sprintf("listings/%d/%d_%s_thumb.jpg", metadata.ListingId, timestamp, uniqueID)

	// Upload original image to MinIO
	if err := s.minioClient.UploadImage(ctx, originalKey, bytes.NewReader(fileBuffer.Bytes()), int64(fileBuffer.Len()), metadata.ContentType); err != nil {
		return nil, fmt.Errorf("failed to upload original image: %w", err)
	}

	s.logger.Debug().Str("key", originalKey).Msg("original image uploaded to MinIO")

	// Upload thumbnail to MinIO
	if err := s.minioClient.UploadImage(ctx, thumbnailKey, &thumbnailBuf, int64(thumbnailBuf.Len()), "image/jpeg"); err != nil {
		// Compensating transaction: Delete original image
		if delErr := s.minioClient.DeleteImage(ctx, originalKey); delErr != nil {
			s.logger.Error().Err(delErr).Str("key", originalKey).Msg("failed to cleanup original image after thumbnail upload failure")
		}
		return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	s.logger.Debug().Str("key", thumbnailKey).Msg("thumbnail uploaded to MinIO")

	// Generate public URLs (permanent, no expiry - bucket is public)
	originalURL := s.minioClient.GetPublicURL(originalKey)
	thumbnailURL := s.minioClient.GetPublicURL(thumbnailKey)

	// Save image metadata to database
	width := int32(img.Bounds().Dx())
	height := int32(img.Bounds().Dy())
	fileSize := int64(fileBuffer.Len())
	mimeType := metadata.ContentType

	dbImage := &domain.ListingImage{
		ListingID:    metadata.ListingId,
		URL:          originalURL,
		StoragePath:  &originalKey,
		ThumbnailURL: &thumbnailURL,
		DisplayOrder: metadata.DisplayOrder,
		IsPrimary:    metadata.IsPrimary,
		Width:        &width,
		Height:       &height,
		FileSize:     &fileSize,
		MimeType:     &mimeType,
	}

	savedImage, err := s.service.AddImage(ctx, dbImage)
	if err != nil {
		// Compensating transaction: Delete both images from MinIO
		_ = s.minioClient.DeleteImage(ctx, originalKey)
		_ = s.minioClient.DeleteImage(ctx, thumbnailKey)
		return nil, fmt.Errorf("failed to save image to database: %w", err)
	}

	s.logger.Info().
		Int64("image_id", savedImage.ID).
		Int64("listing_id", metadata.ListingId).
		Str("filename", metadata.Filename).
		Msg("image uploaded and saved successfully")

	// Convert to proto response
	pbImage := &listingspb.ListingImage{
		Id:           savedImage.ID,
		ListingId:    savedImage.ListingID,
		Url:          savedImage.URL,
		ThumbnailUrl: savedImage.ThumbnailURL,
		DisplayOrder: savedImage.DisplayOrder,
		IsPrimary:    savedImage.IsPrimary,
		Width:        savedImage.Width,
		Height:       savedImage.Height,
		FileSize:     savedImage.FileSize,
		MimeType:     savedImage.MimeType,
	}

	return pbImage, nil
}

// validateImageMetadata validates upload metadata
func (s *Server) validateImageMetadata(metadata *listingspb.UploadImageMetadata) error {
	if metadata.ListingId <= 0 {
		return fmt.Errorf("invalid listing_id: must be positive")
	}
	if metadata.UserId <= 0 {
		return fmt.Errorf("invalid user_id: must be positive")
	}
	if metadata.Filename == "" {
		return fmt.Errorf("filename is required")
	}
	if metadata.FileSize <= 0 {
		return fmt.Errorf("file_size must be positive")
	}
	if metadata.FileSize > maxImageSize {
		return fmt.Errorf("file_size exceeds %dMB limit", maxImageSize/(1024*1024))
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(metadata.Filename))
	if ext == "" {
		return fmt.Errorf("filename must have an extension")
	}
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported file extension: %s (allowed: jpg, jpeg, png, gif, webp)", ext)
	}

	// Validate content type
	if metadata.ContentType == "" {
		return fmt.Errorf("content_type is required")
	}
	if !strings.HasPrefix(metadata.ContentType, "image/") {
		return fmt.Errorf("invalid content_type: must be image/*")
	}

	return nil
}
