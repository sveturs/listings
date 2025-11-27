package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// AddProductImage adds a new image to a B2C product
// Includes ownership validation via storefront_id
func (s *Server) AddProductImage(ctx context.Context, req *listingspb.AddProductImageRequest) (*listingspb.ProductImageResponse, error) {
	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int64("storefront_id", req.StorefrontId).
		Str("url", req.Url).
		Msg("AddProductImage called")

	// ============================================================================
	// VALIDATION
	// ============================================================================

	if req.ProductId <= 0 {
		s.logger.Warn().Int64("product_id", req.ProductId).Msg("invalid product_id")
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}

	if req.StorefrontId <= 0 {
		s.logger.Warn().Int64("storefront_id", req.StorefrontId).Msg("invalid storefront_id")
		return nil, status.Error(codes.InvalidArgument, "storefront_id must be positive")
	}

	if req.Url == "" {
		s.logger.Warn().Msg("empty URL")
		return nil, status.Error(codes.InvalidArgument, "image URL is required")
	}

	// ============================================================================
	// AUTHORIZATION: Verify product belongs to storefront
	// ============================================================================

	product, err := s.service.GetProduct(ctx, req.ProductId, nil)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("product not found")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	if product.StorefrontID != req.StorefrontId {
		s.logger.Warn().
			Int64("product_storefront_id", product.StorefrontID).
			Int64("requested_storefront_id", req.StorefrontId).
			Msg("product does not belong to storefront")
		return nil, status.Error(codes.PermissionDenied, "product does not belong to this storefront")
	}

	// ============================================================================
	// CREATE IMAGE
	// ============================================================================

	productID := req.ProductId
	image := &domain.ProductImage{
		ProductID:    &productID,
		URL:          req.Url,
		DisplayOrder: req.DisplayOrder,
		IsPrimary:    req.IsPrimary,
	}

	// Optional fields
	if req.StoragePath != nil {
		image.StoragePath = req.StoragePath
	}
	if req.ThumbnailUrl != nil {
		image.ThumbnailURL = req.ThumbnailUrl
	}
	if req.Width != nil {
		w := *req.Width
		image.Width = &w
	}
	if req.Height != nil {
		h := *req.Height
		image.Height = &h
	}
	if req.FileSize != nil {
		fs := *req.FileSize
		image.FileSize = &fs
	}
	if req.MimeType != nil {
		image.MimeType = req.MimeType
	}

	// Add image to database
	newImage, err := s.service.AddProductImage(ctx, image)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to add product image")
		return nil, status.Error(codes.Internal, "failed to add product image")
	}

	s.logger.Info().
		Int64("image_id", newImage.ID).
		Int64("product_id", req.ProductId).
		Msg("product image added successfully")

	return &listingspb.ProductImageResponse{
		Image: domainToProtoProductImage(newImage),
	}, nil
}

// GetProductImages retrieves all images for a B2C product
func (s *Server) GetProductImages(ctx context.Context, req *listingspb.GetProductImagesRequest) (*listingspb.ProductImagesResponse, error) {
	s.logger.Debug().Int64("product_id", req.ProductId).Msg("GetProductImages called")

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}

	images, err := s.service.GetProductImages(ctx, req.ProductId)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to get product images")
		return nil, status.Error(codes.Internal, "failed to get product images")
	}

	pbImages := make([]*listingspb.ProductImage, len(images))
	for i, img := range images {
		pbImages[i] = domainToProtoProductImage(img)
	}

	s.logger.Debug().Int("count", len(images)).Int64("product_id", req.ProductId).Msg("product images retrieved")
	return &listingspb.ProductImagesResponse{
		Images: pbImages,
	}, nil
}

// DeleteProductImage removes an image from a B2C product
// Includes authorization and optional MinIO cleanup
func (s *Server) DeleteProductImage(ctx context.Context, req *listingspb.DeleteProductImageRequest) (*listingspb.DeleteProductImageResponse, error) {
	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int64("image_id", req.ImageId).
		Int64("storefront_id", req.StorefrontId).
		Msg("DeleteProductImage called")

	// ============================================================================
	// VALIDATION
	// ============================================================================

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}

	if req.ImageId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "image_id must be positive")
	}

	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id must be positive")
	}

	// ============================================================================
	// AUTHORIZATION: Verify product belongs to storefront
	// ============================================================================

	product, err := s.service.GetProduct(ctx, req.ProductId, nil)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("product not found")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	if product.StorefrontID != req.StorefrontId {
		s.logger.Warn().
			Int64("product_storefront_id", product.StorefrontID).
			Int64("requested_storefront_id", req.StorefrontId).
			Msg("product does not belong to storefront")
		return nil, status.Error(codes.PermissionDenied, "product does not belong to this storefront")
	}

	// ============================================================================
	// GET IMAGE: Verify it belongs to this product
	// ============================================================================

	image, err := s.service.GetProductImageByID(ctx, req.ImageId)
	if err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("image not found")
		return nil, status.Error(codes.NotFound, "image not found")
	}

	if image.ProductID == nil || *image.ProductID != req.ProductId {
		s.logger.Warn().
			Int64("image_id", req.ImageId).
			Int64("requested_product_id", req.ProductId).
			Msg("image does not belong to this product")
		return nil, status.Error(codes.InvalidArgument, "image does not belong to this product")
	}

	// ============================================================================
	// MINIO DELETION (optional, best-effort)
	// ============================================================================

	if s.minioClient != nil && image.StoragePath != nil && *image.StoragePath != "" {
		originalKey := *image.StoragePath

		if err := s.minioClient.DeleteImage(ctx, originalKey); err != nil {
			s.logger.Warn().Err(err).Str("key", originalKey).Msg("failed to delete product image from MinIO (non-critical)")
		} else {
			s.logger.Debug().Str("key", originalKey).Msg("product image deleted from MinIO")
		}

		// Try to delete thumbnail
		thumbnailKey := getThumbnailPath(originalKey)
		if err := s.minioClient.DeleteImage(ctx, thumbnailKey); err != nil {
			s.logger.Warn().Err(err).Str("key", thumbnailKey).Msg("failed to delete product thumbnail from MinIO (non-critical)")
		}
	}

	// ============================================================================
	// DATABASE DELETION
	// ============================================================================

	if err := s.service.DeleteProductImage(ctx, req.ImageId); err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("failed to delete product image from database")
		return nil, status.Error(codes.Internal, "failed to delete product image")
	}

	s.logger.Info().
		Int64("image_id", req.ImageId).
		Int64("product_id", req.ProductId).
		Msg("DeleteProductImage completed")

	return &listingspb.DeleteProductImageResponse{
		Success: true,
		Message: "Product image deleted successfully",
	}, nil
}

// ReorderProductImages updates display order for product images
func (s *Server) ReorderProductImages(ctx context.Context, req *listingspb.ReorderProductImagesRequest) (*listingspb.ReorderProductImagesResponse, error) {
	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int64("storefront_id", req.StorefrontId).
		Int("images_count", len(req.ImageIds)).
		Msg("ReorderProductImages called")

	// ============================================================================
	// VALIDATION
	// ============================================================================

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id must be positive")
	}

	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id must be positive")
	}

	if len(req.ImageIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image_ids cannot be empty")
	}

	// ============================================================================
	// AUTHORIZATION: Verify product belongs to storefront
	// ============================================================================

	product, err := s.service.GetProduct(ctx, req.ProductId, nil)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("product not found")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	if product.StorefrontID != req.StorefrontId {
		s.logger.Warn().
			Int64("product_storefront_id", product.StorefrontID).
			Int64("requested_storefront_id", req.StorefrontId).
			Msg("product does not belong to storefront")
		return nil, status.Error(codes.PermissionDenied, "product does not belong to this storefront")
	}

	// ============================================================================
	// VERIFY ALL IMAGE IDs BELONG TO THIS PRODUCT
	// ============================================================================

	existingImages, err := s.service.GetProductImages(ctx, req.ProductId)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to get existing images")
		return nil, status.Error(codes.Internal, "failed to get existing images")
	}

	existingIDs := make(map[int64]bool)
	for _, img := range existingImages {
		existingIDs[img.ID] = true
	}

	for _, imageID := range req.ImageIds {
		if !existingIDs[imageID] {
			s.logger.Warn().Int64("image_id", imageID).Int64("product_id", req.ProductId).Msg("image does not belong to product")
			return nil, status.Error(codes.InvalidArgument, "one or more images do not belong to this product")
		}
	}

	// ============================================================================
	// REORDER IMAGES
	// ============================================================================

	orders := make([]postgres.ProductImageOrder, len(req.ImageIds))
	for i, imageID := range req.ImageIds {
		orders[i] = postgres.ProductImageOrder{
			ImageID:      imageID,
			DisplayOrder: int32(i),
		}
	}

	if err := s.service.ReorderProductImages(ctx, req.ProductId, orders); err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to reorder product images")
		return nil, status.Error(codes.Internal, "failed to reorder product images")
	}

	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int("count", len(req.ImageIds)).
		Msg("ReorderProductImages completed")

	return &listingspb.ReorderProductImagesResponse{
		Success: true,
	}, nil
}

// domainToProtoProductImage converts domain.ProductImage to proto ProductImage
func domainToProtoProductImage(img *domain.ProductImage) *listingspb.ProductImage {
	if img == nil {
		return nil
	}

	pb := &listingspb.ProductImage{
		Id:           img.ID,
		Url:          img.URL,
		DisplayOrder: img.DisplayOrder,
		IsPrimary:    img.IsPrimary,
		CreatedAt:    img.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    img.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if img.ProductID != nil {
		pb.ProductId = *img.ProductID
	}
	if img.StoragePath != nil {
		pb.StoragePath = img.StoragePath
	}
	if img.ThumbnailURL != nil {
		pb.ThumbnailUrl = img.ThumbnailURL
	}
	if img.Width != nil {
		pb.Width = img.Width
	}
	if img.Height != nil {
		pb.Height = img.Height
	}
	if img.FileSize != nil {
		pb.FileSize = img.FileSize
	}
	if img.MimeType != nil {
		pb.MimeType = img.MimeType
	}

	return pb
}
