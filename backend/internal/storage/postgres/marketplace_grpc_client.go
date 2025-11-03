// backend/internal/storage/postgres/marketplace_grpc_client.go
package postgres

import (
	"context"
	"fmt"
	"time"

	listingsv1 "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
func NewMarketplaceGRPCClient(address string) (*MarketplaceGRPCClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
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
		fmt.Sscanf(userIDStr, "%d", &userID)
		req.UserId = &userID
	}
	if storefrontIDStr, ok := filters["storefront_id"]; ok {
		var storefrontID int64
		fmt.Sscanf(storefrontIDStr, "%d", &storefrontID)
		req.StorefrontId = &storefrontID
	}
	if categoryIDStr, ok := filters["category_id"]; ok {
		var categoryID int64
		fmt.Sscanf(categoryIDStr, "%d", &categoryID)
		req.CategoryId = &categoryID
	}
	if status, ok := filters["status"]; ok {
		req.Status = &status
	}
	if minPriceStr, ok := filters["min_price"]; ok {
		var minPrice float64
		fmt.Sscanf(minPriceStr, "%f", &minPrice)
		req.MinPrice = &minPrice
	}
	if maxPriceStr, ok := filters["max_price"]; ok {
		var maxPrice float64
		fmt.Sscanf(maxPriceStr, "%f", &maxPrice)
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
