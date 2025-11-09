// backend/internal/storage/postgres/marketplace_grpc_client.go
package postgres

import (
	"context"
	"fmt"
	"time"

	listingsv1 "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
)

// MarketplaceGRPCClient wraps gRPC client for listings microservice
type MarketplaceGRPCClient struct {
	client listingsv1.ListingsServiceClient
	conn   *grpc.ClientConn
}

// NewMarketplaceGRPCClient creates new gRPC client
func NewMarketplaceGRPCClient(ctx context.Context, address string) (*MarketplaceGRPCClient, error) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to listings service at %s: %w", address, err)
	}

	client := listingsv1.NewListingsServiceClient(conn)

	logger.Info().Str("address", address).Msg("Connected to listings microservice")

	return &MarketplaceGRPCClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes gRPC connection
func (c *MarketplaceGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IndexListing indexes listing in OpenSearch via microservice
func (c *MarketplaceGRPCClient) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	// For now, we will use CreateListing RPC to index
	// In the future, we should add dedicated IndexListing RPC to listings service

	req := &listingsv1.CreateListingRequest{
		UserId:      int64(listing.UserID),
		Title:       listing.Title,
		Description: &listing.Description,
		Price:       listing.Price,
		Currency:    "RSD", // default currency
		CategoryId:  int64(listing.CategoryID),
		Quantity:    1, // default quantity for marketplace listing
	}

	if listing.StorefrontID != nil {
		storefrontID := int64(*listing.StorefrontID)
		req.StorefrontId = &storefrontID
	}

	_, err := c.client.CreateListing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to index listing %d: %w", listing.ID, err)
	}

	logger.Debug().Int("listing_id", listing.ID).Msg("Listing indexed successfully")
	return nil
}

// DeleteListingIndex removes listing from OpenSearch via microservice
func (c *MarketplaceGRPCClient) DeleteListingIndex(ctx context.Context, id string) error {
	// Convert string id to int64
	var listingID int64
	_, err := fmt.Sscanf(id, "%d", &listingID)
	if err != nil {
		return fmt.Errorf("invalid listing ID format: %s", id)
	}

	req := &listingsv1.DeleteListingRequest{
		Id:     listingID,
		UserId: 0, // Admin deletion - no user check
	}

	_, err = c.client.DeleteListing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete listing index %s: %w", id, err)
	}

	logger.Debug().Str("listing_id", id).Msg("Listing index deleted successfully")
	return nil
}

// SearchListings performs search via microservice
func (c *MarketplaceGRPCClient) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	req := &listingsv1.SearchListingsRequest{
		Query:  params.Query,
		Limit:  int32(params.Size),
		Offset: int32((params.Page - 1) * params.Size),
	}

	// Add category filter if present
	if params.CategoryID != nil {
		categoryID := int64(*params.CategoryID)
		req.CategoryId = &categoryID
	}

	// Add price filters if present
	if params.PriceMin != nil {
		req.MinPrice = params.PriceMin
	}
	if params.PriceMax != nil {
		req.MaxPrice = params.PriceMax
	}

	resp, err := c.client.SearchListings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to search listings: %w", err)
	}

	// Convert protobuf listings to domain models
	listings := make([]*models.MarketplaceListing, len(resp.Listings))
	for i, pbListing := range resp.Listings {
		listings[i] = convertProtoToListing(pbListing)
	}

	result := &search.SearchResult{
		Listings:     listings,
		Total:        int(resp.Total),
		Took:         0, // not provided by proto
		Aggregations: make(map[string][]search.Bucket),
		Suggestions:  []string{},
	}

	return result, nil
}

// SuggestListings provides autocomplete via microservice
func (c *MarketplaceGRPCClient) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	// For now, use search with prefix query
	req := &listingsv1.SearchListingsRequest{
		Query:  prefix,
		Limit:  int32(size),
		Offset: 0,
	}

	resp, err := c.client.SearchListings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions: %w", err)
	}

	// Extract titles as suggestions
	suggestions := make([]string, 0, len(resp.Listings))
	for _, listing := range resp.Listings {
		suggestions = append(suggestions, listing.Title)
	}

	return suggestions, nil
}

// CreateListing creates a new listing via microservice
func (c *MarketplaceGRPCClient) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	req := &listingsv1.CreateListingRequest{
		UserId:      int64(listing.UserID),
		Title:       listing.Title,
		Description: &listing.Description,
		Price:       listing.Price,
		Currency:    "RSD", // default currency
		CategoryId:  int64(listing.CategoryID),
		Quantity:    1, // default quantity for marketplace listing
	}

	if listing.StorefrontID != nil {
		storefrontID := int64(*listing.StorefrontID)
		req.StorefrontId = &storefrontID
	}

	resp, err := c.client.CreateListing(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to create listing via gRPC: %w", err)
	}

	logger.Debug().Int64("listing_id", resp.Listing.Id).Msg("Listing created successfully via gRPC")
	return int(resp.Listing.Id), nil
}

// GetListings retrieves listings with filters via microservice
func (c *MarketplaceGRPCClient) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	req := &listingsv1.ListListingsRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Apply filters
	if userIDStr, ok := filters["user_id"]; ok {
		var userID int64
		_, _ = fmt.Sscanf(userIDStr, "%d", &userID)
		req.UserId = &userID
	}
	if storefrontIDStr, ok := filters["storefront_id"]; ok {
		var storefrontID int64
		_, _ = fmt.Sscanf(storefrontIDStr, "%d", &storefrontID)
		req.StorefrontId = &storefrontID
	}
	if categoryIDStr, ok := filters["category_id"]; ok {
		var categoryID int64
		_, _ = fmt.Sscanf(categoryIDStr, "%d", &categoryID)
		req.CategoryId = &categoryID
	}
	if status, ok := filters["status"]; ok {
		req.Status = &status
	}
	if minPriceStr, ok := filters["min_price"]; ok {
		var minPrice float64
		_, _ = fmt.Sscanf(minPriceStr, "%f", &minPrice)
		req.MinPrice = &minPrice
	}
	if maxPriceStr, ok := filters["max_price"]; ok {
		var maxPrice float64
		_, _ = fmt.Sscanf(maxPriceStr, "%f", &maxPrice)
		req.MaxPrice = &maxPrice
	}

	resp, err := c.client.ListListings(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get listings via gRPC: %w", err)
	}

	// Convert protobuf listings to domain models
	listings := make([]models.MarketplaceListing, len(resp.Listings))
	for i, pbListing := range resp.Listings {
		converted := convertProtoToListing(pbListing)
		if converted != nil {
			listings[i] = *converted
		}
	}

	return listings, int64(resp.Total), nil
}

// GetListingByID retrieves a single listing by ID via microservice
func (c *MarketplaceGRPCClient) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	req := &listingsv1.GetListingRequest{
		Id: int64(id),
	}

	resp, err := c.client.GetListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing by ID %d via gRPC: %w", id, err)
	}

	return convertProtoToListing(resp.Listing), nil
}

// GetListingBySlug retrieves a single listing by slug via microservice
// Note: Currently gRPC service doesn't support slug lookup, so this returns error
func (c *MarketplaceGRPCClient) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	// TODO: Implement slug support in listings microservice
	return nil, fmt.Errorf("GetListingBySlug not yet implemented in gRPC service (slug: %s)", slug)
}

// UpdateListing updates an existing listing via microservice
func (c *MarketplaceGRPCClient) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	req := &listingsv1.UpdateListingRequest{
		Id:     int64(listing.ID),
		UserId: int64(listing.UserID),
	}

	// Set optional fields
	if listing.Title != "" {
		req.Title = &listing.Title
	}
	if listing.Description != "" {
		req.Description = &listing.Description
	}
	if listing.Price > 0 {
		req.Price = &listing.Price
	}
	if listing.Status != "" {
		req.Status = &listing.Status
	}

	// Note: Quantity is not directly in MarketplaceListing, using default
	quantity := int32(1)
	req.Quantity = &quantity

	_, err := c.client.UpdateListing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update listing %d via gRPC: %w", listing.ID, err)
	}

	logger.Debug().Int("listing_id", listing.ID).Msg("Listing updated successfully via gRPC")
	return nil
}

// DeleteListing soft-deletes a listing (user ownership check) via microservice
func (c *MarketplaceGRPCClient) DeleteListing(ctx context.Context, id int, userID int) error {
	req := &listingsv1.DeleteListingRequest{
		Id:     int64(id),
		UserId: int64(userID),
	}

	resp, err := c.client.DeleteListing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete listing %d via gRPC: %w", id, err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to delete listing %d: microservice returned false", id)
	}

	logger.Debug().Int("listing_id", id).Int("user_id", userID).Msg("Listing deleted successfully via gRPC")
	return nil
}

// DeleteListingAdmin hard-deletes a listing (admin, no ownership check) via microservice
func (c *MarketplaceGRPCClient) DeleteListingAdmin(ctx context.Context, id int) error {
	req := &listingsv1.DeleteListingRequest{
		Id:     int64(id),
		UserId: 0, // 0 = admin deletion (no ownership check)
	}

	resp, err := c.client.DeleteListing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to admin-delete listing %d via gRPC: %w", id, err)
	}

	if !resp.Success {
		return fmt.Errorf("failed to admin-delete listing %d: microservice returned false", id)
	}

	logger.Debug().Int("listing_id", id).Msg("Listing admin-deleted successfully via gRPC")
	return nil
}

// GetListingImage retrieves a single image by ID via microservice
func (c *MarketplaceGRPCClient) GetListingImage(ctx context.Context, imageID int64) (*models.MarketplaceImage, error) {
	req := &listingsv1.ImageIDRequest{
		ImageId: imageID,
	}

	resp, err := c.client.GetListingImage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing image %d via gRPC: %w", imageID, err)
	}

	return convertProtoToImage(resp.Image), nil
}

// DeleteListingImage removes an image from a listing via microservice
func (c *MarketplaceGRPCClient) DeleteListingImage(ctx context.Context, imageID int64) error {
	req := &listingsv1.ImageIDRequest{
		ImageId: imageID,
	}

	_, err := c.client.DeleteListingImage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete listing image %d via gRPC: %w", imageID, err)
	}

	logger.Debug().Int64("image_id", imageID).Msg("Listing image deleted successfully via gRPC")
	return nil
}

// AddListingImage adds a new image to a listing via microservice
func (c *MarketplaceGRPCClient) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (*models.MarketplaceImage, error) {
	req := &listingsv1.AddImageRequest{
		ListingId:    int64(image.ListingID),
		Url:          image.PublicURL,
		DisplayOrder: int32(image.DisplayOrder),
		IsPrimary:    image.IsMain,
	}

	if image.FilePath != "" {
		req.StoragePath = &image.FilePath
	}
	if image.ThumbnailURL != "" {
		req.ThumbnailUrl = &image.ThumbnailURL
	}
	if image.FileSize > 0 {
		fileSize := int64(image.FileSize)
		req.FileSize = &fileSize
	}
	if image.ContentType != "" {
		req.MimeType = &image.ContentType
	}

	resp, err := c.client.AddListingImage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to add listing image via gRPC: %w", err)
	}

	result := convertProtoToImage(resp.Image)
	logger.Debug().Int("listing_id", image.ListingID).Int("image_id", result.ID).Msg("Listing image added successfully via gRPC")
	return result, nil
}

// GetListingImages retrieves all images for a listing via microservice
func (c *MarketplaceGRPCClient) GetListingImages(ctx context.Context, listingID int64) ([]*models.MarketplaceImage, error) {
	req := &listingsv1.ListingIDRequest{
		ListingId: listingID,
	}

	resp, err := c.client.GetListingImages(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing images via gRPC: %w", err)
	}

	images := make([]*models.MarketplaceImage, len(resp.Images))
	for i, pbImg := range resp.Images {
		images[i] = convertProtoToImage(pbImg)
	}

	return images, nil
}

// GetRootCategories retrieves all top-level categories via microservice
func (c *MarketplaceGRPCClient) GetRootCategories(ctx context.Context) ([]*models.MarketplaceCategory, error) {
	resp, err := c.client.GetRootCategories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get root categories via gRPC: %w", err)
	}

	categories := make([]*models.MarketplaceCategory, len(resp.Categories))
	for i, pbCat := range resp.Categories {
		categories[i] = convertProtoToCategory(pbCat)
	}

	return categories, nil
}

// GetAllCategories retrieves all categories via microservice
func (c *MarketplaceGRPCClient) GetAllCategories(ctx context.Context) ([]*models.MarketplaceCategory, error) {
	resp, err := c.client.GetAllCategories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories via gRPC: %w", err)
	}

	categories := make([]*models.MarketplaceCategory, len(resp.Categories))
	for i, pbCat := range resp.Categories {
		categories[i] = convertProtoToCategory(pbCat)
	}

	return categories, nil
}

// GetPopularCategories retrieves popular categories by listing count via microservice
func (c *MarketplaceGRPCClient) GetPopularCategories(ctx context.Context, limit int) ([]*models.MarketplaceCategory, error) {
	req := &listingsv1.PopularCategoriesRequest{
		Limit: int32(limit),
	}

	resp, err := c.client.GetPopularCategories(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular categories via gRPC: %w", err)
	}

	categories := make([]*models.MarketplaceCategory, len(resp.Categories))
	for i, pbCat := range resp.Categories {
		categories[i] = convertProtoToCategory(pbCat)
	}

	return categories, nil
}

// GetCategoryByID retrieves a single category by ID via microservice
func (c *MarketplaceGRPCClient) GetCategoryByID(ctx context.Context, categoryID int64) (*models.MarketplaceCategory, error) {
	req := &listingsv1.CategoryIDRequest{
		CategoryId: categoryID,
	}

	resp, err := c.client.GetCategory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by ID %d via gRPC: %w", categoryID, err)
	}

	return convertProtoToCategory(resp.Category), nil
}

// GetCategoryTree retrieves category hierarchy starting from a node via microservice
func (c *MarketplaceGRPCClient) GetCategoryTree(ctx context.Context, categoryID int64) (*models.CategoryTreeNode, error) {
	req := &listingsv1.CategoryIDRequest{
		CategoryId: categoryID,
	}

	resp, err := c.client.GetCategoryTree(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get category tree for %d via gRPC: %w", categoryID, err)
	}

	return convertProtoToCategoryTree(resp.Tree), nil
}

// GetFavoritedUsers retrieves list of user IDs who favorited a listing via microservice
func (c *MarketplaceGRPCClient) GetFavoritedUsers(ctx context.Context, listingID int64) ([]string, error) {
	req := &listingsv1.ListingIDRequest{
		ListingId: listingID,
	}

	resp, err := c.client.GetFavoritedUsers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorited users via gRPC: %w", err)
	}

	// Convert int64 user IDs to strings
	userIDs := make([]string, len(resp.UserIds))
	for i, uid := range resp.UserIds {
		userIDs[i] = fmt.Sprintf("%d", uid)
	}

	return userIDs, nil
}

// CreateListingVariants creates multiple variants for a listing via microservice
func (c *MarketplaceGRPCClient) CreateListingVariants(ctx context.Context, listingID int64, variants []*models.MarketplaceListingVariant) error {
	variantInputs := make([]*listingsv1.VariantInput, len(variants))
	for i, v := range variants {
		variantInputs[i] = &listingsv1.VariantInput{
			Sku:        v.SKU,
			IsActive:   true,
			Attributes: make(map[string]string),
		}

		if v.Price != nil {
			variantInputs[i].Price = v.Price
		}
		if v.Stock != nil {
			stock := int32(*v.Stock)
			variantInputs[i].Stock = &stock
		}
		if v.ImageURL != nil && *v.ImageURL != "" {
			variantInputs[i].ImageUrl = v.ImageURL
		}

		// Convert attributes map to proto map
		if len(v.Attributes) > 0 {
			for key, val := range v.Attributes {
				variantInputs[i].Attributes[key] = val
			}
		}
	}

	req := &listingsv1.CreateVariantsRequest{
		ListingId: listingID,
		Variants:  variantInputs,
	}

	_, err := c.client.CreateVariants(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create listing variants via gRPC: %w", err)
	}

	logger.Debug().Int64("listing_id", listingID).Int("count", len(variants)).Msg("Listing variants created successfully via gRPC")
	return nil
}

// GetListingVariants retrieves all variants for a listing via microservice
func (c *MarketplaceGRPCClient) GetListingVariants(ctx context.Context, listingID int64) ([]*models.MarketplaceListingVariant, error) {
	req := &listingsv1.ListingIDRequest{
		ListingId: listingID,
	}

	resp, err := c.client.GetVariants(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing variants via gRPC: %w", err)
	}

	variants := make([]*models.MarketplaceListingVariant, len(resp.Variants))
	for i, pbVariant := range resp.Variants {
		variants[i] = convertProtoToVariant(pbVariant)
	}

	return variants, nil
}

// UpdateListingVariant updates a specific variant via microservice
func (c *MarketplaceGRPCClient) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error {
	req := &listingsv1.UpdateVariantRequest{
		VariantId:  int64(variant.ID),
		Attributes: make(map[string]string),
	}

	if variant.SKU != "" {
		req.Sku = &variant.SKU
	}
	if variant.Price != nil {
		req.Price = variant.Price
	}
	if variant.Stock != nil {
		stock := int32(*variant.Stock)
		req.Stock = &stock
	}
	if variant.ImageURL != nil && *variant.ImageURL != "" {
		req.ImageUrl = variant.ImageURL
	}

	isActive := variant.IsActive
	req.IsActive = &isActive

	// Convert attributes map to proto map
	if len(variant.Attributes) > 0 {
		for key, val := range variant.Attributes {
			req.Attributes[key] = val
		}
	}

	_, err := c.client.UpdateVariant(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update listing variant %d via gRPC: %w", variant.ID, err)
	}

	logger.Debug().Int("variant_id", variant.ID).Msg("Listing variant updated successfully via gRPC")
	return nil
}

// DeleteListingVariant removes a variant from a listing via microservice
func (c *MarketplaceGRPCClient) DeleteListingVariant(ctx context.Context, variantID int64) error {
	req := &listingsv1.VariantIDRequest{
		VariantId: variantID,
	}

	_, err := c.client.DeleteVariant(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete listing variant %d via gRPC: %w", variantID, err)
	}

	logger.Debug().Int64("variant_id", variantID).Msg("Listing variant deleted successfully via gRPC")
	return nil
}

// GetMarketplaceListingsForReindex retrieves listings that need reindexing via microservice
func (c *MarketplaceGRPCClient) GetMarketplaceListingsForReindex(ctx context.Context, limit int) ([]*models.MarketplaceListing, error) {
	req := &listingsv1.ReindexRequest{
		BatchSize: int32(limit),
	}

	resp, err := c.client.GetListingsForReindex(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings for reindex via gRPC: %w", err)
	}

	listings := make([]*models.MarketplaceListing, len(resp.Listings))
	for i, pbListing := range resp.Listings {
		listings[i] = convertProtoToListing(pbListing)
	}

	return listings, nil
}

// ResetMarketplaceListingsReindexFlag resets reindex flags for specified listings via microservice
func (c *MarketplaceGRPCClient) ResetMarketplaceListingsReindexFlag(ctx context.Context, listingIDs []int64) error {
	req := &listingsv1.ResetFlagsRequest{
		ListingIds: listingIDs,
	}

	_, err := c.client.ResetReindexFlags(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to reset reindex flags via gRPC: %w", err)
	}

	logger.Debug().Int("count", len(listingIDs)).Msg("Reindex flags reset successfully via gRPC")
	return nil
}

// SynchronizeDiscountMetadata synchronizes discount information across listings via microservice
func (c *MarketplaceGRPCClient) SynchronizeDiscountMetadata(ctx context.Context) error {
	_, err := c.client.SyncDiscounts(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to synchronize discount metadata via gRPC: %w", err)
	}

	logger.Info().Msg("Discount metadata synchronized successfully via gRPC")
	return nil
}

// convertProtoToListing converts protobuf Listing to domain MarketplaceListing
func convertProtoToListing(pb *listingsv1.Listing) *models.MarketplaceListing {
	listing := &models.MarketplaceListing{
		ID:         int(pb.Id),
		UserID:     int(pb.UserId),
		CategoryID: int(pb.CategoryId),
		Title:      pb.Title,
		Price:      pb.Price,
		Status:     pb.Status,
		ViewsCount: int(pb.ViewsCount),
		IsFavorite: false, // will be set later if needed
		ShowOnMap:  false, // default value
		Variants:   []models.MarketplaceListingVariant{},
		Attributes: []models.ListingAttributeValue{},
	}

	// Handle optional fields
	if pb.Description != nil {
		listing.Description = *pb.Description
	}

	if pb.StorefrontId != nil {
		storefrontID := int(*pb.StorefrontId)
		listing.StorefrontID = &storefrontID
	}

	// Parse timestamps
	if pb.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, pb.CreatedAt)
		if err == nil {
			listing.CreatedAt = createdAt
		}
	}

	if pb.UpdatedAt != "" {
		updatedAt, err := time.Parse(time.RFC3339, pb.UpdatedAt)
		if err == nil {
			listing.UpdatedAt = updatedAt
		}
	}

	// Convert images
	if len(pb.Images) > 0 {
		listing.Images = make([]models.MarketplaceImage, len(pb.Images))
		for i, pbImg := range pb.Images {
			listing.Images[i] = models.MarketplaceImage{
				ID:           int(pbImg.Id),
				ListingID:    int(pbImg.ListingId),
				PublicURL:    pbImg.Url,
				IsMain:       pbImg.IsPrimary,
				DisplayOrder: int(pbImg.DisplayOrder),
			}

			if pbImg.ThumbnailUrl != nil {
				listing.Images[i].ThumbnailURL = *pbImg.ThumbnailUrl
			}

			// Parse image timestamps
			if pbImg.CreatedAt != "" {
				imgCreatedAt, err := time.Parse(time.RFC3339, pbImg.CreatedAt)
				if err == nil {
					listing.Images[i].CreatedAt = imgCreatedAt
				}
			}
		}
	}

	// Convert location if present
	if pb.Location != nil {
		loc := pb.Location
		if loc.Country != nil {
			listing.Country = *loc.Country
		}
		if loc.City != nil {
			listing.City = *loc.City
		}
		if loc.Latitude != nil {
			lat := *loc.Latitude
			listing.Latitude = &lat
		}
		if loc.Longitude != nil {
			lon := *loc.Longitude
			listing.Longitude = &lon
		}
	}

	// Convert attributes
	if len(pb.Attributes) > 0 {
		listing.Attributes = make([]models.ListingAttributeValue, len(pb.Attributes))
		for i, pbAttr := range pb.Attributes {
			textValue := pbAttr.AttributeValue
			listing.Attributes[i] = models.ListingAttributeValue{
				AttributeName: pbAttr.AttributeKey,
				TextValue:     &textValue,
			}
		}
	}

	return listing
}

// convertProtoToImage converts protobuf ListingImage to domain MarketplaceImage
func convertProtoToImage(pb *listingsv1.ListingImage) *models.MarketplaceImage {
	image := &models.MarketplaceImage{
		ID:           int(pb.Id),
		ListingID:    int(pb.ListingId),
		PublicURL:    pb.Url,
		IsMain:       pb.IsPrimary,
		DisplayOrder: int(pb.DisplayOrder),
	}

	if pb.StoragePath != nil {
		image.FilePath = *pb.StoragePath
	}
	if pb.ThumbnailUrl != nil {
		image.ThumbnailURL = *pb.ThumbnailUrl
	}
	if pb.FileSize != nil {
		image.FileSize = int(*pb.FileSize)
	}
	if pb.MimeType != nil {
		image.ContentType = *pb.MimeType
	}

	// Parse timestamps
	if pb.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, pb.CreatedAt)
		if err == nil {
			image.CreatedAt = createdAt
		}
	}

	return image
}

// convertProtoToCategory converts protobuf Category to domain MarketplaceCategory
func convertProtoToCategory(pb *listingsv1.Category) *models.MarketplaceCategory {
	cat := &models.MarketplaceCategory{
		ID:           int(pb.Id),
		Name:         pb.Name,
		Slug:         pb.Slug,
		IsActive:     true, // default
		SortOrder:    int(pb.SortOrder),
		Level:        int(pb.Level),
		ListingCount: int(pb.ListingCount),
		HasCustomUI:  pb.HasCustomUi,
	}

	if pb.ParentId != nil {
		parentID := int(*pb.ParentId)
		cat.ParentID = &parentID
	}
	if pb.Icon != nil {
		cat.Icon = pb.Icon
	}
	if pb.Description != nil {
		cat.Description = pb.Description
	}
	if pb.CustomUiComponent != nil {
		cat.CustomUIComponent = pb.CustomUiComponent
	}

	// Parse timestamps
	if pb.CreatedAt != "" {
		createdAt, err := time.Parse(time.RFC3339, pb.CreatedAt)
		if err == nil {
			cat.CreatedAt = createdAt
		}
	}

	return cat
}

// convertProtoToCategoryTree converts protobuf CategoryTreeNode to domain CategoryTreeNode
func convertProtoToCategoryTree(pb *listingsv1.CategoryTreeNode) *models.CategoryTreeNode {
	node := &models.CategoryTreeNode{
		ID:            int(pb.Id),
		Name:          pb.Name,
		Slug:          pb.Slug,
		Level:         int(pb.Level),
		Path:          pb.Path,
		ListingCount:  int(pb.ListingCount),
		ChildrenCount: int(pb.ChildrenCount),
		HasCustomUI:   pb.HasCustomUi,
		Children:      []models.CategoryTreeNode{},
		Translations:  pb.Translations,
	}

	if pb.Icon != nil {
		node.Icon = *pb.Icon
	}
	if pb.ParentId != nil {
		parentID := int(*pb.ParentId)
		node.ParentID = &parentID
	}
	if pb.CustomUiComponent != nil {
		node.CustomUIComponent = *pb.CustomUiComponent
	}

	// Recursively convert children
	if len(pb.Children) > 0 {
		node.Children = make([]models.CategoryTreeNode, len(pb.Children))
		for i, child := range pb.Children {
			converted := convertProtoToCategoryTree(child)
			if converted != nil {
				node.Children[i] = *converted
			}
		}
	}

	return node
}

// convertProtoToVariant converts protobuf ListingVariant to domain MarketplaceListingVariant
func convertProtoToVariant(pb *listingsv1.ListingVariant) *models.MarketplaceListingVariant {
	variant := &models.MarketplaceListingVariant{
		ID:         int(pb.Id),
		ListingID:  int(pb.ListingId),
		SKU:        pb.Sku,
		IsActive:   true, // default
		Attributes: make(map[string]string),
	}

	if pb.Price != nil {
		variant.Price = pb.Price
	}
	if pb.Stock != nil {
		stock := int(*pb.Stock)
		variant.Stock = &stock
	}
	if pb.ImageUrl != nil {
		variant.ImageURL = pb.ImageUrl
	}

	// Convert attributes map
	if len(pb.Attributes) > 0 {
		variant.Attributes = pb.Attributes
	}

	return variant
}

// ============================================================================
// Product CRUD Operations (Phase 9.5.2)
// ============================================================================

// CreateProduct creates a new B2C product via microservice
func (c *MarketplaceGRPCClient) CreateProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	// Apply timeout for single operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Convert request to protobuf
	pbReq := &listingsv1.CreateProductRequest{
		StorefrontId:          int64(storefrontID),
		Name:                  req.Name,
		Description:           req.Description,
		Price:                 req.Price,
		Currency:              req.Currency,
		CategoryId:            int64(req.CategoryID),
		StockQuantity:         int32(req.StockQuantity),
		IsActive:              req.IsActive,
		HasIndividualLocation: req.HasIndividualLocation != nil && *req.HasIndividualLocation,
		ShowOnMap:             req.ShowOnMap == nil || *req.ShowOnMap,
		HasVariants:           req.HasVariants,
	}

	// Optional fields
	if req.SKU != nil {
		pbReq.Sku = req.SKU
	}
	if req.Barcode != nil {
		pbReq.Barcode = req.Barcode
	}
	if req.IndividualAddress != nil {
		pbReq.IndividualAddress = req.IndividualAddress
	}
	if req.IndividualLatitude != nil {
		pbReq.IndividualLatitude = req.IndividualLatitude
	}
	if req.IndividualLongitude != nil {
		pbReq.IndividualLongitude = req.IndividualLongitude
	}
	if req.LocationPrivacy != nil {
		pbReq.LocationPrivacy = req.LocationPrivacy
	}

	// Call gRPC
	resp, err := c.client.CreateProduct(ctx, pbReq)
	if err != nil {
		logger.Error().Err(err).
			Int("storefront_id", storefrontID).
			Str("name", req.Name).
			Msg("Failed to create product via gRPC")
		return nil, fmt.Errorf("failed to create product via gRPC: %w", err)
	}

	// Convert proto to model
	product := convertProtoToProduct(resp.Product)
	logger.Debug().
		Int64("product_id", resp.Product.Id).
		Int("storefront_id", storefrontID).
		Msg("Product created successfully via gRPC")

	return product, nil
}

// UpdateProduct updates an existing product via microservice
func (c *MarketplaceGRPCClient) UpdateProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pbReq := &listingsv1.UpdateProductRequest{
		ProductId:    int64(productID),
		StorefrontId: int64(storefrontID),
	}

	// Set optional fields
	if req.Name != nil {
		pbReq.Name = req.Name
	}
	if req.Description != nil {
		pbReq.Description = req.Description
	}
	if req.Price != nil {
		pbReq.Price = req.Price
	}
	if req.CategoryID != nil {
		categoryID := int64(*req.CategoryID)
		pbReq.CategoryId = &categoryID
	}
	if req.SKU != nil {
		pbReq.Sku = req.SKU
	}
	if req.Barcode != nil {
		pbReq.Barcode = req.Barcode
	}
	if req.IsActive != nil {
		pbReq.IsActive = req.IsActive
	}
	if req.HasIndividualLocation != nil {
		pbReq.HasIndividualLocation = req.HasIndividualLocation
	}
	if req.IndividualAddress != nil {
		pbReq.IndividualAddress = req.IndividualAddress
	}
	if req.IndividualLatitude != nil {
		pbReq.IndividualLatitude = req.IndividualLatitude
	}
	if req.IndividualLongitude != nil {
		pbReq.IndividualLongitude = req.IndividualLongitude
	}
	if req.LocationPrivacy != nil {
		pbReq.LocationPrivacy = req.LocationPrivacy
	}
	if req.ShowOnMap != nil {
		pbReq.ShowOnMap = req.ShowOnMap
	}

	_, err := c.client.UpdateProduct(ctx, pbReq)
	if err != nil {
		logger.Error().Err(err).
			Int("product_id", productID).
			Int("storefront_id", storefrontID).
			Msg("Failed to update product via gRPC")
		return fmt.Errorf("failed to update product via gRPC: %w", err)
	}

	logger.Debug().
		Int("product_id", productID).
		Int("storefront_id", storefrontID).
		Msg("Product updated successfully via gRPC")

	return nil
}

// DeleteProduct removes a product via microservice
func (c *MarketplaceGRPCClient) DeleteProduct(ctx context.Context, storefrontID, productID int, hardDelete bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &listingsv1.DeleteProductRequest{
		ProductId:    int64(productID),
		StorefrontId: int64(storefrontID),
		HardDelete:   hardDelete,
	}

	_, err := c.client.DeleteProduct(ctx, req)
	if err != nil {
		logger.Error().Err(err).
			Int("product_id", productID).
			Int("storefront_id", storefrontID).
			Bool("hard_delete", hardDelete).
			Msg("Failed to delete product via gRPC")
		return fmt.Errorf("failed to delete product via gRPC: %w", err)
	}

	logger.Debug().
		Int("product_id", productID).
		Int("storefront_id", storefrontID).
		Bool("hard_delete", hardDelete).
		Msg("Product deleted successfully via gRPC")

	return nil
}

// BulkUpdateProducts updates multiple products via microservice
func (c *MarketplaceGRPCClient) BulkUpdateProducts(ctx context.Context, storefrontID int, updates []models.BulkUpdateItem) ([]int, []error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Convert to protobuf
	pbUpdates := make([]*listingsv1.ProductUpdateInput, len(updates))
	for i, update := range updates {
		pbUpdates[i] = &listingsv1.ProductUpdateInput{
			ProductId: int64(update.ProductID),
		}

		// Set fields from UpdateProductRequest
		if update.Updates.Name != nil {
			pbUpdates[i].Name = update.Updates.Name
		}
		if update.Updates.Description != nil {
			pbUpdates[i].Description = update.Updates.Description
		}
		if update.Updates.Price != nil {
			pbUpdates[i].Price = update.Updates.Price
		}
		if update.Updates.CategoryID != nil {
			categoryID := int64(*update.Updates.CategoryID)
			pbUpdates[i].CategoryId = &categoryID
		}
		if update.Updates.SKU != nil {
			pbUpdates[i].Sku = update.Updates.SKU
		}
		if update.Updates.Barcode != nil {
			pbUpdates[i].Barcode = update.Updates.Barcode
		}
		if update.Updates.IsActive != nil {
			pbUpdates[i].IsActive = update.Updates.IsActive
		}
	}

	req := &listingsv1.BulkUpdateProductsRequest{
		StorefrontId: int64(storefrontID),
		Updates:      pbUpdates,
	}

	resp, err := c.client.BulkUpdateProducts(ctx, req)
	if err != nil {
		logger.Error().Err(err).
			Int("storefront_id", storefrontID).
			Int("count", len(updates)).
			Msg("Failed to bulk update products via gRPC")
		return nil, []error{fmt.Errorf("failed to bulk update products via gRPC: %w", err)}
	}

	// Extract successful IDs
	successIDs := make([]int, len(resp.Products))
	for i, product := range resp.Products {
		successIDs[i] = int(product.Id)
	}

	// Convert errors
	var errors []error
	for _, protoErr := range resp.Errors {
		errors = append(errors, fmt.Errorf("product %d: %s", protoErr.ProductId, protoErr.ErrorMessage))
	}

	logger.Debug().
		Int("storefront_id", storefrontID).
		Int("successful", int(resp.SuccessfulCount)).
		Int("failed", int(resp.FailedCount)).
		Msg("Bulk update products completed via gRPC")

	return successIDs, errors
}

// RecordInventoryMovement records inventory movement via microservice
// Note: This is a stub - listings microservice doesn't have dedicated inventory endpoint yet
// We handle inventory through UpdateProduct with stock_quantity field
func (c *MarketplaceGRPCClient) RecordInventoryMovement(ctx context.Context, productID int, variantID *int, quantity int, reason, notes string, userID int) error {
	// For now, this is a no-op - inventory is handled through stock quantity updates
	logger.Warn().
		Int("product_id", productID).
		Interface("variant_id", variantID).
		Int("quantity", quantity).
		Str("reason", reason).
		Msg("RecordInventoryMovement called but not implemented in listings microservice - using local DB")

	return fmt.Errorf("RecordInventoryMovement not implemented in listings microservice")
}

// BatchUpdateStock updates stock for multiple products via microservice
// Note: Currently not implemented in proto - will use BulkUpdateProducts instead
func (c *MarketplaceGRPCClient) BatchUpdateStock(ctx context.Context, updates []struct {
	ProductID int
	VariantID *int
	Quantity  int
}) ([]int, []error) {
	logger.Warn().
		Int("count", len(updates)).
		Msg("BatchUpdateStock called but not implemented - falling back to local DB")

	return nil, []error{fmt.Errorf("BatchUpdateStock not implemented in listings microservice")}
}

// IncrementProductViews increments view count for a product via microservice
// Note: Currently not implemented in proto
func (c *MarketplaceGRPCClient) IncrementProductViews(ctx context.Context, productID int) error {
	logger.Debug().
		Int("product_id", productID).
		Msg("IncrementProductViews called but not implemented - using local DB")

	return fmt.Errorf("IncrementProductViews not implemented in listings microservice")
}

// GetProductStats retrieves product statistics via microservice
// Note: Currently not implemented in proto
func (c *MarketplaceGRPCClient) GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error) {
	logger.Debug().
		Int("storefront_id", storefrontID).
		Msg("GetProductStats called but not implemented - using local DB")

	return nil, fmt.Errorf("GetProductStats not implemented in listings microservice")
}

// ============================================================================
// Helper Functions - Product Conversion
// ============================================================================

// convertProtoToProduct converts protobuf Product to domain StorefrontProduct
func convertProtoToProduct(pb *listingsv1.Product) *models.StorefrontProduct {
	product := &models.StorefrontProduct{
		ID:                     int(pb.Id),
		StorefrontID:           int(pb.StorefrontId),
		Name:                   pb.Name,
		Description:            pb.Description,
		Price:                  pb.Price,
		Currency:               pb.Currency,
		CategoryID:             int(pb.CategoryId),
		StockQuantity:          int(pb.StockQuantity),
		StockStatus:            pb.StockStatus,
		IsActive:               pb.IsActive,
		ViewCount:              int(pb.ViewCount),
		SoldCount:              int(pb.SoldCount),
		HasIndividualLocation:  pb.HasIndividualLocation,
		ShowOnMap:              pb.ShowOnMap,
		HasVariants:            pb.HasVariants,
		Images:                 []models.StorefrontProductImage{},
		Variants:               []models.StorefrontProductVariant{},
	}

	// Optional fields
	if pb.Sku != nil {
		product.SKU = pb.Sku
	}
	if pb.Barcode != nil {
		product.Barcode = pb.Barcode
	}
	if pb.IndividualAddress != nil {
		product.IndividualAddress = pb.IndividualAddress
	}
	if pb.IndividualLatitude != nil {
		product.IndividualLatitude = pb.IndividualLatitude
	}
	if pb.IndividualLongitude != nil {
		product.IndividualLongitude = pb.IndividualLongitude
	}
	if pb.LocationPrivacy != nil {
		product.LocationPrivacy = pb.LocationPrivacy
	}

	// Parse timestamps
	if pb.CreatedAt != nil {
		product.CreatedAt = pb.CreatedAt.AsTime()
	}
	if pb.UpdatedAt != nil {
		product.UpdatedAt = pb.UpdatedAt.AsTime()
	}

	// Convert variants if present
	if len(pb.Variants) > 0 {
		product.Variants = make([]models.StorefrontProductVariant, len(pb.Variants))
		for i, pbVariant := range pb.Variants {
			product.Variants[i] = *convertProtoToProductVariant(pbVariant)
		}
	}

	return product
}

// convertProtoToProductVariant converts protobuf ProductVariant to domain StorefrontProductVariant
func convertProtoToProductVariant(pb *listingsv1.ProductVariant) *models.StorefrontProductVariant {
	variant := &models.StorefrontProductVariant{
		ID:            int(pb.Id),
		ProductID:     int(pb.ProductId),
		StockQuantity: int(pb.StockQuantity),
		StockStatus:   pb.StockStatus,
		IsActive:      pb.IsActive,
		IsDefault:     pb.IsDefault,
		ViewCount:     int(pb.ViewCount),
		SoldCount:     int(pb.SoldCount),
	}

	// Optional fields
	if pb.Sku != nil {
		variant.SKU = pb.Sku
	}
	if pb.Barcode != nil {
		variant.Barcode = pb.Barcode
	}
	if pb.Price != nil {
		price := *pb.Price
		variant.Price = &price
	}
	if pb.CompareAtPrice != nil {
		compareAtPrice := *pb.CompareAtPrice
		variant.CompareAtPrice = &compareAtPrice
	}
	if pb.CostPrice != nil {
		costPrice := *pb.CostPrice
		variant.CostPrice = &costPrice
	}
	if pb.LowStockThreshold != nil {
		threshold := int(*pb.LowStockThreshold)
		variant.LowStockThreshold = &threshold
	}
	if pb.Weight != nil {
		weight := *pb.Weight
		variant.Weight = &weight
	}

	// Parse timestamps
	if pb.CreatedAt != nil {
		variant.CreatedAt = pb.CreatedAt.AsTime()
	}
	if pb.UpdatedAt != nil {
		variant.UpdatedAt = pb.UpdatedAt.AsTime()
	}

	return variant
}
