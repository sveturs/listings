package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
)

// GetListingImage retrieves a single image by ID
func (s *Server) GetListingImage(ctx context.Context, req *listingspb.ImageIDRequest) (*listingspb.ImageResponse, error) {
	s.logger.Debug().Int64("image_id", req.ImageId).Msg("GetListingImage called")

	if req.ImageId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "image ID must be greater than 0")
	}

	// Get image from repository (need to add this method)
	image, err := s.service.GetImageByID(ctx, req.ImageId)
	if err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("failed to get image")
		return nil, status.Error(codes.NotFound, fmt.Sprintf("image not found: %v", err))
	}

	return &listingspb.ImageResponse{
		Image: DomainToProtoImage(image),
	}, nil
}

// DeleteListingImage removes an image from a listing
func (s *Server) DeleteListingImage(ctx context.Context, req *listingspb.ImageIDRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("image_id", req.ImageId).Msg("DeleteListingImage called")

	if req.ImageId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "image ID must be greater than 0")
	}

	// Delete image via service
	err := s.service.DeleteImage(ctx, req.ImageId)
	if err != nil {
		s.logger.Error().Err(err).Int64("image_id", req.ImageId).Msg("failed to delete image")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete image: %v", err))
	}

	s.logger.Info().Int64("image_id", req.ImageId).Msg("image deleted successfully")
	return &emptypb.Empty{}, nil
}

// AddListingImage adds a new image to a listing
func (s *Server) AddListingImage(ctx context.Context, req *listingspb.AddImageRequest) (*listingspb.ImageResponse, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Msg("AddListingImage called")

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "image URL is required")
	}

	// Convert proto to domain
	image := ProtoToAddImageInput(req)

	// Add image via service
	newImage, err := s.service.AddImage(ctx, image)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to add image")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to add image: %v", err))
	}

	s.logger.Info().Int64("image_id", newImage.ID).Int64("listing_id", req.ListingId).Msg("image added successfully")
	return &listingspb.ImageResponse{
		Image: DomainToProtoImage(newImage),
	}, nil
}

// GetListingImages retrieves all images for a listing
func (s *Server) GetListingImages(ctx context.Context, req *listingspb.ListingIDRequest) (*listingspb.ImagesResponse, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Msg("GetListingImages called")

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Get images from service
	images, err := s.service.GetImages(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to get images")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get images: %v", err))
	}

	// Convert to proto
	pbImages := make([]*listingspb.ListingImage, len(images))
	for i, img := range images {
		pbImages[i] = DomainToProtoImage(img)
	}

	s.logger.Debug().Int("count", len(images)).Msg("images retrieved")
	return &listingspb.ImagesResponse{
		Images: pbImages,
	}, nil
}

// ReorderListingImages updates display order for multiple images
func (s *Server) ReorderListingImages(ctx context.Context, req *listingspb.ReorderImagesRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Int("count", len(req.ImageOrders)).Msg("ReorderListingImages called")

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	if len(req.ImageOrders) == 0 {
		return nil, status.Error(codes.InvalidArgument, "image orders cannot be empty")
	}

	// Convert proto orders to repository orders
	orders := make([]postgres.ImageOrder, len(req.ImageOrders))
	for i, pbOrder := range req.ImageOrders {
		orders[i] = postgres.ImageOrder{
			ImageID:      pbOrder.ImageId,
			DisplayOrder: pbOrder.DisplayOrder,
		}
	}

	// Reorder images via service
	err := s.service.ReorderImages(ctx, req.ListingId, orders)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to reorder images")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to reorder images: %v", err))
	}

	s.logger.Info().Int64("listing_id", req.ListingId).Int("count", len(orders)).Msg("images reordered successfully")
	return &emptypb.Empty{}, nil
}

// GetRootCategories retrieves all top-level categories
func (s *Server) GetRootCategories(ctx context.Context, req *emptypb.Empty) (*listingspb.CategoriesResponse, error) {
	s.logger.Debug().Msg("GetRootCategories called")

	// Get root categories from service
	categories, err := s.service.GetRootCategories(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get root categories")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get root categories: %v", err))
	}

	// Convert to proto
	pbCategories := make([]*listingspb.Category, len(categories))
	for i, cat := range categories {
		pbCategories[i] = DomainToProtoCategory(cat)
	}

	s.logger.Debug().Int("count", len(categories)).Msg("root categories retrieved")
	return &listingspb.CategoriesResponse{
		Categories: pbCategories,
	}, nil
}

// GetAllCategories retrieves all categories in the system
func (s *Server) GetAllCategories(ctx context.Context, req *emptypb.Empty) (*listingspb.CategoriesResponse, error) {
	s.logger.Debug().Msg("GetAllCategories called")

	// Get all categories from service
	categories, err := s.service.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get all categories")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get all categories: %v", err))
	}

	// Convert to proto
	pbCategories := make([]*listingspb.Category, len(categories))
	for i, cat := range categories {
		pbCategories[i] = DomainToProtoCategory(cat)
	}

	s.logger.Debug().Int("count", len(categories)).Msg("all categories retrieved")
	return &listingspb.CategoriesResponse{
		Categories: pbCategories,
	}, nil
}

// GetPopularCategories retrieves most popular categories by listing count
func (s *Server) GetPopularCategories(ctx context.Context, req *listingspb.PopularCategoriesRequest) (*listingspb.CategoriesResponse, error) {
	s.logger.Debug().Int32("limit", req.Limit).Msg("GetPopularCategories called")

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	// Get popular categories from service
	categories, err := s.service.GetPopularCategories(ctx, limit)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get popular categories")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get popular categories: %v", err))
	}

	// Convert to proto
	pbCategories := make([]*listingspb.Category, len(categories))
	for i, cat := range categories {
		pbCategories[i] = DomainToProtoCategory(cat)
	}

	s.logger.Debug().Int("count", len(categories)).Msg("popular categories retrieved")
	return &listingspb.CategoriesResponse{
		Categories: pbCategories,
	}, nil
}

// GetCategory retrieves a single category by ID
func (s *Server) GetCategory(ctx context.Context, req *listingspb.CategoryIDRequest) (*listingspb.CategoryResponse, error) {
	s.logger.Debug().Int64("category_id", req.CategoryId).Msg("GetCategory called")

	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category ID must be greater than 0")
	}

	// Get category from service
	category, err := s.service.GetCategoryByID(ctx, req.CategoryId)
	if err != nil {
		s.logger.Error().Err(err).Int64("category_id", req.CategoryId).Msg("failed to get category")
		return nil, status.Error(codes.NotFound, fmt.Sprintf("category not found: %v", err))
	}

	return &listingspb.CategoryResponse{
		Category: DomainToProtoCategory(category),
	}, nil
}

// GetCategoryTree retrieves category hierarchy starting from a node
func (s *Server) GetCategoryTree(ctx context.Context, req *listingspb.CategoryIDRequest) (*listingspb.CategoryTreeResponse, error) {
	s.logger.Debug().Int64("category_id", req.CategoryId).Msg("GetCategoryTree called")

	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category ID must be greater than 0")
	}

	// Get category tree from service
	tree, err := s.service.GetCategoryTree(ctx, req.CategoryId)
	if err != nil {
		s.logger.Error().Err(err).Int64("category_id", req.CategoryId).Msg("failed to get category tree")
		return nil, status.Error(codes.NotFound, fmt.Sprintf("category tree not found: %v", err))
	}

	return &listingspb.CategoryTreeResponse{
		Tree: DomainToProtoCategoryTree(tree),
	}, nil
}

// GetFavoritedUsers retrieves list of user IDs who favorited a listing
func (s *Server) GetFavoritedUsers(ctx context.Context, req *listingspb.ListingIDRequest) (*listingspb.UserIDsResponse, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Msg("GetFavoritedUsers called")

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Get favorited users from service
	userIDs, err := s.service.GetFavoritedUsers(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to get favorited users")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get favorited users: %v", err))
	}

	s.logger.Debug().Int("count", len(userIDs)).Msg("favorited users retrieved")
	return &listingspb.UserIDsResponse{
		UserIds: userIDs,
	}, nil
}

// AddToFavorites adds a listing to user's favorites
func (s *Server) AddToFavorites(ctx context.Context, req *listingspb.AddToFavoritesRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("AddToFavorites called")

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Add to favorites via service
	err := s.service.AddToFavorites(ctx, req.UserId, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("failed to add to favorites")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to add to favorites: %v", err))
	}

	s.logger.Info().Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("added to favorites successfully")
	return &emptypb.Empty{}, nil
}

// RemoveFromFavorites removes a listing from user's favorites
func (s *Server) RemoveFromFavorites(ctx context.Context, req *listingspb.RemoveFromFavoritesRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("RemoveFromFavorites called")

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Remove from favorites via service
	err := s.service.RemoveFromFavorites(ctx, req.UserId, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("failed to remove from favorites")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to remove from favorites: %v", err))
	}

	s.logger.Info().Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("removed from favorites successfully")
	return &emptypb.Empty{}, nil
}

// GetUserFavorites retrieves list of listing IDs favorited by a user
func (s *Server) GetUserFavorites(ctx context.Context, req *listingspb.GetUserFavoritesRequest) (*listingspb.GetUserFavoritesResponse, error) {
	s.logger.Debug().Int64("user_id", req.UserId).Msg("GetUserFavorites called")

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	// Get user favorites from service
	listingIDs, err := s.service.GetUserFavorites(ctx, req.UserId)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", req.UserId).Msg("failed to get user favorites")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user favorites: %v", err))
	}

	s.logger.Debug().Int("count", len(listingIDs)).Msg("user favorites retrieved")
	return &listingspb.GetUserFavoritesResponse{
		ListingIds: listingIDs,
		Total:      int32(len(listingIDs)),
	}, nil
}

// IsFavorite checks if a listing is in user's favorites
func (s *Server) IsFavorite(ctx context.Context, req *listingspb.IsFavoriteRequest) (*listingspb.IsFavoriteResponse, error) {
	s.logger.Debug().Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("IsFavorite called")

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Check favorite status via service
	isFavorite, err := s.service.IsFavorite(ctx, req.UserId, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", req.UserId).Int64("listing_id", req.ListingId).Msg("failed to check favorite status")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to check favorite status: %v", err))
	}

	s.logger.Debug().Bool("is_favorite", isFavorite).Msg("favorite status checked")
	return &listingspb.IsFavoriteResponse{
		IsFavorite: isFavorite,
	}, nil
}

// CreateVariants creates multiple variants for a listing
func (s *Server) CreateVariants(ctx context.Context, req *listingspb.CreateVariantsRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Int("count", len(req.Variants)).Msg("CreateVariants called")

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	if len(req.Variants) == 0 {
		return &emptypb.Empty{}, nil // No variants to create
	}

	// Convert proto variants to domain
	variants := make([]*domain.ListingVariant, len(req.Variants))
	for i, v := range req.Variants {
		variants[i] = ProtoToVariantInput(v, req.ListingId)
	}

	// Create variants via service
	err := s.service.CreateVariants(ctx, variants)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to create variants")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create variants: %v", err))
	}

	s.logger.Info().Int64("listing_id", req.ListingId).Int("count", len(variants)).Msg("variants created successfully")
	return &emptypb.Empty{}, nil
}

// GetVariants retrieves all variants for a listing
func (s *Server) GetVariants(ctx context.Context, req *listingspb.ListingIDRequest) (*listingspb.VariantsResponse, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Msg("GetVariants called")

	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Get variants from service
	variants, err := s.service.GetVariants(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to get variants")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get variants: %v", err))
	}

	// Convert to proto
	pbVariants := make([]*listingspb.ListingVariant, len(variants))
	for i, v := range variants {
		pbVariants[i] = DomainToProtoVariant(v)
	}

	s.logger.Debug().Int("count", len(variants)).Msg("variants retrieved")
	return &listingspb.VariantsResponse{
		Variants: pbVariants,
	}, nil
}

// UpdateVariant updates a specific variant
func (s *Server) UpdateVariant(ctx context.Context, req *listingspb.UpdateVariantRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("variant_id", req.VariantId).Msg("UpdateVariant called")

	if req.VariantId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "variant ID must be greater than 0")
	}

	// Build update - first get existing variant
	variant, err := s.service.GetVariantByID(ctx, req.VariantId)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", req.VariantId).Msg("variant not found")
		return nil, status.Error(codes.NotFound, fmt.Sprintf("variant not found: %v", err))
	}

	// Update fields
	if req.Sku != nil {
		variant.SKU = *req.Sku
	}
	if req.Price != nil {
		variant.Price = req.Price
	}
	if req.Stock != nil {
		variant.Stock = req.Stock
	}
	if req.ImageUrl != nil {
		variant.ImageURL = req.ImageUrl
	}
	if req.IsActive != nil {
		variant.IsActive = *req.IsActive
	}
	if len(req.Attributes) > 0 {
		variant.Attributes = req.Attributes
	}

	// Update via service
	err = s.service.UpdateVariant(ctx, variant)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", req.VariantId).Msg("failed to update variant")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update variant: %v", err))
	}

	s.logger.Info().Int64("variant_id", req.VariantId).Msg("variant updated successfully")
	return &emptypb.Empty{}, nil
}

// DeleteVariant removes a variant from a listing
func (s *Server) DeleteVariant(ctx context.Context, req *listingspb.VariantIDRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int64("variant_id", req.VariantId).Msg("DeleteVariant called")

	if req.VariantId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "variant ID must be greater than 0")
	}

	// Delete variant via service
	err := s.service.DeleteVariant(ctx, req.VariantId)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", req.VariantId).Msg("failed to delete variant")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete variant: %v", err))
	}

	s.logger.Info().Int64("variant_id", req.VariantId).Msg("variant deleted successfully")
	return &emptypb.Empty{}, nil
}

// GetListingsForReindex retrieves listings that need reindexing
func (s *Server) GetListingsForReindex(ctx context.Context, req *listingspb.ReindexRequest) (*listingspb.ListingsResponse, error) {
	s.logger.Debug().Int32("batch_size", req.BatchSize).Msg("GetListingsForReindex called")

	batchSize := int(req.BatchSize)
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}
	if batchSize > 1000 {
		batchSize = 1000 // Max batch size
	}

	// Get listings for reindex from service
	listings, err := s.service.GetListingsForReindex(ctx, batchSize)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get listings for reindex")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get listings for reindex: %v", err))
	}

	// Convert to proto
	pbListings := make([]*listingspb.Listing, len(listings))
	for i, listing := range listings {
		pbListings[i] = DomainToProtoListing(listing)
	}

	s.logger.Debug().Int("count", len(listings)).Msg("listings for reindex retrieved")
	return &listingspb.ListingsResponse{
		Listings: pbListings,
		Total:    int32(len(listings)),
	}, nil
}

// ResetReindexFlags resets reindex flags for specified listings
func (s *Server) ResetReindexFlags(ctx context.Context, req *listingspb.ResetFlagsRequest) (*emptypb.Empty, error) {
	s.logger.Debug().Int("count", len(req.ListingIds)).Msg("ResetReindexFlags called")

	if len(req.ListingIds) == 0 {
		return &emptypb.Empty{}, nil // Nothing to reset
	}

	// Reset reindex flags via service
	err := s.service.ResetReindexFlags(ctx, req.ListingIds)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to reset reindex flags")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to reset reindex flags: %v", err))
	}

	s.logger.Info().Int("count", len(req.ListingIds)).Msg("reindex flags reset successfully")
	return &emptypb.Empty{}, nil
}

// SyncDiscounts synchronizes discount information across listings
func (s *Server) SyncDiscounts(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	s.logger.Debug().Msg("SyncDiscounts called")

	// Sync discounts via service
	err := s.service.SyncDiscounts(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to sync discounts")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to sync discounts: %v", err))
	}

	s.logger.Info().Msg("discounts synced successfully")
	return &emptypb.Empty{}, nil
}
