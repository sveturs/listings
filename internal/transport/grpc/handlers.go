package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	attributespb "github.com/vondi-global/listings/api/proto/attributes/v1"
	categoriespb "github.com/vondi-global/listings/api/proto/categories/v1"
	chatsvcv1 "github.com/vondi-global/listings/api/proto/chat/v1"
	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/metrics"
	minioclient "github.com/vondi-global/listings/internal/repository/minio"
	"github.com/vondi-global/listings/internal/service"
	"github.com/vondi-global/listings/internal/service/listings"
)

// contains checks if substring is in string (case-insensitive helper)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Server implements gRPC ListingsServiceServer, AttributeServiceServer, CategoryServiceServer, OrderServiceServer, AnalyticsServiceServer, and ChatServiceServer
type Server struct {
	listingspb.UnimplementedListingsServiceServer
	attributespb.UnimplementedAttributeServiceServer
	categoriespb.UnimplementedCategoryServiceServer
	listingspb.UnimplementedOrderServiceServer
	listingspb.UnimplementedAnalyticsServiceServer
	chatsvcv1.UnimplementedChatServiceServer
	service                    *listings.Service
	storefrontService          *listings.StorefrontService
	attrService                service.AttributeService
	categoryService            service.CategoryService
	orderService               service.OrderService
	cartService                service.CartService
	chatService                service.ChatService
	analyticsService           service.AnalyticsService
	storefrontAnalyticsService service.StorefrontAnalyticsService
	inventoryService           service.InventoryService
	invitationService          *service.InvitationService
	minioClient                *minioclient.Client
	metrics                    *metrics.Metrics
	logger                     zerolog.Logger
}

// NewServer creates a new gRPC server instance
func NewServer(
	service *listings.Service,
	storefrontService *listings.StorefrontService,
	attrService service.AttributeService,
	categoryService service.CategoryService,
	orderService service.OrderService,
	cartService service.CartService,
	chatService service.ChatService,
	analyticsService service.AnalyticsService,
	storefrontAnalyticsService service.StorefrontAnalyticsService,
	inventoryService service.InventoryService,
	invitationService *service.InvitationService,
	minioClient *minioclient.Client,
	m *metrics.Metrics,
	logger zerolog.Logger,
) *Server {
	return &Server{
		service:                    service,
		storefrontService:          storefrontService,
		attrService:                attrService,
		categoryService:            categoryService,
		orderService:               orderService,
		cartService:                cartService,
		chatService:                chatService,
		analyticsService:           analyticsService,
		storefrontAnalyticsService: storefrontAnalyticsService,
		inventoryService:           inventoryService,
		invitationService:          invitationService,
		minioClient:                minioClient,
		metrics:                    m,
		logger:                     logger.With().Str("component", "grpc_handler").Logger(),
	}
}

// GetListing retrieves a single listing by ID
func (s *Server) GetListing(ctx context.Context, req *listingspb.GetListingRequest) (*listingspb.GetListingResponse, error) {
	// Extract requested language
	requestedLang := ""
	if req.Lang != nil {
		requestedLang = *req.Lang
	}

	s.logger.Debug().
		Int64("listing_id", req.Id).
		Str("requested_lang", requestedLang).
		Msg("GetListing called")

	// Validate request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Get listing from service
	listing, err := s.service.GetListing(ctx, req.Id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to get listing")
		return nil, status.Error(codes.NotFound, fmt.Sprintf("listing not found: %v", err))
	}

	// Check for nil listing (defensive programming)
	if listing == nil {
		s.logger.Error().Int64("listing_id", req.Id).Msg("listing returned nil without error")
		return nil, status.Error(codes.NotFound, "listing not found")
	}

	// Log translation info for debugging
	s.logger.Debug().
		Int64("listing_id", req.Id).
		Str("requested_lang", requestedLang).
		Str("original_lang", listing.OriginalLanguage).
		Bool("has_title_translations", len(listing.TitleTranslations) > 0).
		Bool("has_description_translations", len(listing.DescriptionTranslations) > 0).
		Msg("GetListing translation info")

	// Apply translation if requested and different from original language
	if requestedLang != "" && requestedLang != listing.OriginalLanguage {
		s.logger.Info().
			Int64("listing_id", req.Id).
			Str("requested_lang", requestedLang).
			Str("original_lang", listing.OriginalLanguage).
			Str("original_title", listing.Title).
			Msg("Applying translation to listing")

		ApplyTranslation(listing, requestedLang)

		s.logger.Info().
			Int64("listing_id", req.Id).
			Str("lang", requestedLang).
			Str("translated_title", listing.Title).
			Msg("Translation applied successfully")
	}

	// Convert to proto
	pbListing := DomainToProtoListing(listing)

	return &listingspb.GetListingResponse{
		Listing: pbListing,
	}, nil
}

// CreateListing creates a new listing
func (s *Server) CreateListing(ctx context.Context, req *listingspb.CreateListingRequest) (*listingspb.CreateListingResponse, error) {
	s.logger.Debug().Int64("user_id", req.UserId).Str("title", req.Title).Msg("CreateListing called")

	// DEBUG: Log incoming proto fields
	s.logger.Debug().
		Bool("has_condition", req.Condition != nil).
		Bool("has_location", req.Location != nil).
		Bool("has_show_on_map", req.ShowOnMap != nil).
		Int("attributes_count", len(req.Attributes)).
		Msg("DEBUG: incoming proto request fields")

	if req.Location != nil {
		logEvent := s.logger.Debug()
		if req.Location.Country != nil {
			logEvent = logEvent.Str("country", *req.Location.Country)
		}
		if req.Location.City != nil {
			logEvent = logEvent.Str("city", *req.Location.City)
		}
		if req.Location.Latitude != nil {
			logEvent = logEvent.Float64("lat", *req.Location.Latitude)
		}
		if req.Location.Longitude != nil {
			logEvent = logEvent.Float64("lng", *req.Location.Longitude)
		}
		logEvent.Msg("DEBUG: incoming location data")
	}

	// Validate request
	if err := s.validateCreateListingRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain input
	input := ProtoToCreateListingInput(req)

	// DEBUG: Log converted domain input fields
	s.logger.Debug().
		Bool("has_condition", input.Condition != nil).
		Bool("has_location", input.Location != nil).
		Bool("has_show_on_map", input.ShowOnMap != nil).
		Int("attributes_count", len(input.Attributes)).
		Msg("DEBUG: converted domain input fields")

	if input.Location != nil {
		logEvent := s.logger.Debug()
		if input.Location.Country != nil {
			logEvent = logEvent.Str("country", *input.Location.Country)
		}
		if input.Location.City != nil {
			logEvent = logEvent.Str("city", *input.Location.City)
		}
		if input.Location.Latitude != nil {
			logEvent = logEvent.Float64("lat", *input.Location.Latitude)
		}
		if input.Location.Longitude != nil {
			logEvent = logEvent.Float64("lng", *input.Location.Longitude)
		}
		logEvent.Msg("DEBUG: converted location data")
	}

	// Create listing via service
	listing, err := s.service.CreateListing(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create listing")

		// Map errors to appropriate gRPC codes (use contains for wrapped errors)
		errMsg := err.Error()
		if contains(errMsg, "category not found") || contains(errMsg, "fk_listings_category_id") {
			return nil, status.Error(codes.InvalidArgument, "category not found")
		}
		if contains(errMsg, "validation failed") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create listing: %v", err))
	}

	// Convert to proto
	pbListing := DomainToProtoListing(listing)

	s.logger.Info().Int64("listing_id", listing.ID).Msg("listing created successfully")
	return &listingspb.CreateListingResponse{
		Listing: pbListing,
	}, nil
}

// UpdateListing updates an existing listing
func (s *Server) UpdateListing(ctx context.Context, req *listingspb.UpdateListingRequest) (*listingspb.UpdateListingResponse, error) {
	s.logger.Debug().Int64("listing_id", req.Id).Int64("user_id", req.UserId).Msg("UpdateListing called")

	// Validate request
	if err := s.validateUpdateListingRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain input
	input := ProtoToUpdateListingInput(req)

	// Update listing via service (with ownership check)
	listing, err := s.service.UpdateListing(ctx, req.Id, req.UserId, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to update listing")

		// Map errors to appropriate gRPC codes (use contains for wrapped errors)
		errMsg := err.Error()
		if errMsg == "unauthorized: user does not own this listing" {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if contains(errMsg, "listing not found") || contains(errMsg, "sql: no rows in result set") {
			return nil, status.Error(codes.NotFound, "listing not found")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update listing: %v", err))
	}

	// Convert to proto
	pbListing := DomainToProtoListing(listing)

	s.logger.Info().Int64("listing_id", listing.ID).Msg("listing updated successfully")
	return &listingspb.UpdateListingResponse{
		Listing: pbListing,
	}, nil
}

// DeleteListing soft-deletes a listing
func (s *Server) DeleteListing(ctx context.Context, req *listingspb.DeleteListingRequest) (*listingspb.DeleteListingResponse, error) {
	s.logger.Debug().
		Int64("listing_id", req.Id).
		Int64("user_id", req.UserId).
		Bool("is_admin", req.IsAdmin).
		Msg("DeleteListing called")

	// Validate request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	// If admin is requesting deletion, bypass ownership check
	if req.IsAdmin {
		s.logger.Info().
			Int64("listing_id", req.Id).
			Int64("admin_user_id", req.UserId).
			Msg("Admin bypass: deleting listing without ownership check")

		// Get listing first to verify it exists
		listing, err := s.service.GetListing(ctx, req.Id)
		if err != nil {
			s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to get listing for admin deletion")
			return nil, status.Error(codes.NotFound, "listing not found")
		}

		// Delete listing via service using listing owner's ID
		err = s.service.DeleteListing(ctx, req.Id, listing.UserID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Int64("listing_id", req.Id).
				Int64("admin_user_id", req.UserId).
				Int64("listing_owner_id", listing.UserID).
				Msg("failed to delete listing as admin")
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete listing: %v", err))
		}

		s.logger.Info().
			Int64("listing_id", req.Id).
			Int64("admin_user_id", req.UserId).
			Int64("listing_owner_id", listing.UserID).
			Msg("Listing deleted successfully by admin")

		return &listingspb.DeleteListingResponse{
			Success: true,
		}, nil
	}

	// Delete listing via service (with ownership check)
	err := s.service.DeleteListing(ctx, req.Id, req.UserId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to delete listing")

		// Map errors to appropriate gRPC codes (use contains for wrapped errors)
		errMsg := err.Error()
		if errMsg == "unauthorized: user does not own this listing" {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if contains(errMsg, "listing not found") || contains(errMsg, "sql: no rows in result set") {
			return nil, status.Error(codes.NotFound, "listing not found")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete listing: %v", err))
	}

	s.logger.Info().Int64("listing_id", req.Id).Msg("listing deleted successfully")
	return &listingspb.DeleteListingResponse{
		Success: true,
	}, nil
}

// SearchListings performs full-text search on listings
func (s *Server) SearchListings(ctx context.Context, req *listingspb.SearchListingsRequest) (*listingspb.SearchListingsResponse, error) {
	s.logger.Debug().Str("query", req.Query).Msg("SearchListings called")

	// Validate request
	if err := s.validateSearchListingsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain query
	query := ProtoToSearchListingsQuery(req)

	// Search listings via service
	listings, total, err := s.service.SearchListings(ctx, query)
	if err != nil {
		s.logger.Error().Err(err).Str("query", req.Query).Msg("failed to search listings")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to search listings: %v", err))
	}

	// Convert to proto
	pbListings := make([]*listingspb.Listing, len(listings))
	for i, listing := range listings {
		pbListings[i] = DomainToProtoListing(listing)
	}

	s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("search completed")
	return &listingspb.SearchListingsResponse{
		Listings: pbListings,
		Total:    total,
	}, nil
}

// GetSimilarListings retrieves similar listings based on category and price
func (s *Server) GetSimilarListings(ctx context.Context, req *listingspb.GetSimilarListingsRequest) (*listingspb.GetSimilarListingsResponse, error) {
	s.logger.Debug().Int64("listing_id", req.ListingId).Msg("GetSimilarListings called")

	// Validate request
	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing_id must be positive")
	}

	// Default limit
	limit := int32(10)
	if req.Limit != nil {
		limit = *req.Limit
	}

	// Validate limit range
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}

	// Get similar listings from service
	listings, total, err := s.service.GetSimilarListings(ctx, req.ListingId, limit)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.ListingId).Msg("failed to get similar listings")
		// Graceful degradation: return empty array on error
		return &listingspb.GetSimilarListingsResponse{
			Listings: []*listingspb.Listing{},
			Total:    0,
		}, nil
	}

	// Convert to proto
	pbListings := make([]*listingspb.Listing, len(listings))
	for i, listing := range listings {
		pbListings[i] = DomainToProtoListing(listing)
	}

	s.logger.Debug().
		Int64("listing_id", req.ListingId).
		Int("count", len(listings)).
		Int32("total", total).
		Msg("similar listings retrieved")

	return &listingspb.GetSimilarListingsResponse{
		Listings: pbListings,
		Total:    total,
	}, nil
}

// ListListings returns a paginated list of listings
func (s *Server) ListListings(ctx context.Context, req *listingspb.ListListingsRequest) (*listingspb.ListListingsResponse, error) {
	s.logger.Debug().Int32("limit", req.Limit).Int32("offset", req.Offset).Msg("ListListings called")

	// Validate request
	if err := s.validateListListingsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain filter
	filter := ProtoToListListingsFilter(req)

	// List listings via service
	listings, total, err := s.service.ListListings(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list listings")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list listings: %v", err))
	}

	// Convert to proto
	pbListings := make([]*listingspb.Listing, len(listings))
	for i, listing := range listings {
		pbListings[i] = DomainToProtoListing(listing)
	}

	s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("listings retrieved")
	return &listingspb.ListListingsResponse{
		Listings: pbListings,
		Total:    total,
	}, nil
}

// Validation helpers

func (s *Server) validateCreateListingRequest(req *listingspb.CreateListingRequest) error {
	if req.UserId <= 0 {
		return fmt.Errorf("user_id must be greater than 0")
	}

	if req.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(req.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}

	if len(req.Title) > 255 {
		return fmt.Errorf("title must not exceed 255 characters")
	}

	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters (ISO 4217)")
	}

	if req.CategoryId <= 0 {
		return fmt.Errorf("category_id must be greater than 0")
	}

	if req.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	return nil
}

func (s *Server) validateUpdateListingRequest(req *listingspb.UpdateListingRequest) error {
	if req.Id <= 0 {
		return fmt.Errorf("listing ID must be greater than 0")
	}

	if req.UserId <= 0 {
		return fmt.Errorf("user ID must be greater than 0")
	}

	// At least one field must be set
	if req.Title == nil && req.Description == nil && req.Price == nil && req.Quantity == nil && req.Status == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Validate individual fields if present
	if req.Title != nil {
		title := *req.Title
		if title == "" {
			return fmt.Errorf("title cannot be empty")
		}
		if len(title) < 3 {
			return fmt.Errorf("title must be at least 3 characters")
		}
		if len(title) > 255 {
			return fmt.Errorf("title must not exceed 255 characters")
		}
	}

	if req.Price != nil && *req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if req.Quantity != nil && *req.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if req.Status != nil {
		validStatuses := map[string]bool{
			"draft":    true,
			"active":   true,
			"inactive": true,
			"sold":     true,
			"archived": true,
		}
		status := *req.Status
		if !validStatuses[status] {
			return fmt.Errorf("invalid status: %s", status)
		}
	}

	return nil
}

func (s *Server) validateSearchListingsRequest(req *listingspb.SearchListingsRequest) error {
	// Query is optional - if provided, validate it
	if req.Query != "" && len(req.Query) < 2 {
		return fmt.Errorf("search query must be at least 2 characters")
	}

	if req.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if req.Limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}

	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MinPrice > *req.MaxPrice {
			return fmt.Errorf("min_price cannot be greater than max_price")
		}
	}

	return nil
}

func (s *Server) validateListListingsRequest(req *listingspb.ListListingsRequest) error {
	if req.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if req.Limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}

	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MinPrice > *req.MaxPrice {
			return fmt.Errorf("min_price cannot be greater than max_price")
		}
	}

	return nil
}

// ============================================================================
// ANALYTICS SERVICE HANDLERS
// ============================================================================

// GetOverviewStats retrieves platform-wide analytics statistics (admin only)
func (s *Server) GetOverviewStats(
	ctx context.Context,
	req *listingspb.GetOverviewStatsRequest,
) (*listingspb.GetOverviewStatsResponse, error) {
	// Log request
	s.logger.Info().
		Time("date_from", req.DateFrom.AsTime()).
		Time("date_to", req.DateTo.AsTime()).
		Interface("period", req.Period).
		Interface("storefront_id", req.StorefrontId).
		Interface("category_id", req.CategoryId).
		Interface("listing_type", req.ListingType).
		Msg("GetOverviewStats RPC called")

	// Extract auth from gRPC metadata
	userID, isAdmin, err := s.extractAuthFromMetadata(ctx)
	if err != nil {
		s.logger.Warn().
			Err(err).
			Msg("failed to extract auth from metadata")
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Validate request
	if err := s.validateOverviewStatsRequest(req); err != nil {
		s.logger.Warn().
			Err(err).
			Int64("user_id", userID).
			Msg("invalid overview stats request")
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}

	// Call service
	response, err := s.analyticsService.GetOverviewStats(ctx, req, userID, isAdmin)
	if err != nil {
		// Check error type to determine appropriate gRPC status code
		if contains(err.Error(), "admin") || contains(err.Error(), "unauthorized") {
			s.logger.Warn().
				Err(err).
				Int64("user_id", userID).
				Bool("is_admin", isAdmin).
				Msg("permission denied for overview stats")
			return nil, status.Error(codes.PermissionDenied, "admin access required")
		}

		if contains(err.Error(), "invalid") {
			s.logger.Warn().
				Err(err).
				Int64("user_id", userID).
				Msg("service validation error for overview stats")
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid input: %v", err))
		}

		s.logger.Error().
			Err(err).
			Int64("user_id", userID).
			Time("date_from", req.DateFrom.AsTime()).
			Time("date_to", req.DateTo.AsTime()).
			Msg("service error getting overview stats")
		return nil, status.Error(codes.Internal, "failed to retrieve analytics")
	}

	// Log success
	s.logger.Info().
		Int64("user_id", userID).
		Time("date_from", req.DateFrom.AsTime()).
		Time("date_to", req.DateTo.AsTime()).
		Int32("total_views", response.Engagement.TotalViews).
		Int32("total_orders", response.Orders.TotalOrders).
		Msg("GetOverviewStats completed successfully")

	return response, nil
}

// GetListingStats retrieves analytics for a specific listing (owner or admin)
func (s *Server) GetListingStats(
	ctx context.Context,
	req *listingspb.GetListingStatsRequest,
) (*listingspb.GetListingStatsResponse, error) {
	// Log request
	s.logger.Info().
		Int64("listing_id", req.GetListingId()).
		Int64("product_id", req.GetProductId()).
		Time("date_from", req.DateFrom.AsTime()).
		Time("date_to", req.DateTo.AsTime()).
		Interface("period", req.Period).
		Interface("include_variants", req.IncludeVariants).
		Interface("include_geo", req.IncludeGeo).
		Msg("GetListingStats RPC called")

	// Extract auth from gRPC metadata
	userID, isAdmin, err := s.extractAuthFromMetadata(ctx)
	if err != nil {
		s.logger.Warn().
			Err(err).
			Msg("failed to extract auth from metadata")
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Validate request
	if err := s.validateListingStatsRequest(req); err != nil {
		s.logger.Warn().
			Err(err).
			Int64("user_id", userID).
			Int64("listing_id", req.GetListingId()).
			Msg("invalid listing stats request")
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid request: %v", err))
	}

	// Call service
	response, err := s.analyticsService.GetListingStats(ctx, req, userID, isAdmin)
	if err != nil {
		// Check error type to determine appropriate gRPC status code
		if contains(err.Error(), "permission") || contains(err.Error(), "unauthorized") {
			s.logger.Warn().
				Err(err).
				Int64("user_id", userID).
				Bool("is_admin", isAdmin).
				Int64("listing_id", req.GetListingId()).
				Msg("permission denied for listing stats")
			return nil, status.Error(codes.PermissionDenied, "access denied to this listing")
		}

		if contains(err.Error(), "invalid") || contains(err.Error(), "required") {
			s.logger.Warn().
				Err(err).
				Int64("user_id", userID).
				Msg("service validation error for listing stats")
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid input: %v", err))
		}

		s.logger.Error().
			Err(err).
			Int64("user_id", userID).
			Int64("listing_id", req.GetListingId()).
			Time("date_from", req.DateFrom.AsTime()).
			Time("date_to", req.DateTo.AsTime()).
			Msg("service error getting listing stats")
		return nil, status.Error(codes.Internal, "failed to retrieve listing analytics")
	}

	// Log success
	s.logger.Info().
		Int64("user_id", userID).
		Int64("listing_id", response.ListingId).
		Str("listing_name", response.ListingName).
		Int32("total_views", response.TotalViews).
		Int32("total_sales", response.TotalSales).
		Float64("total_revenue", response.TotalRevenue).
		Msg("GetListingStats completed successfully")

	return response, nil
}

// GetStorefrontStats retrieves analytics for a storefront (owner or admin)
func (s *Server) GetStorefrontStats(
	ctx context.Context,
	req *listingspb.GetStorefrontStatsRequest,
) (*listingspb.GetStorefrontStatsResponse, error) {
	// Log request
	s.logger.Info().
		Int64("storefront_id", req.StorefrontId).
		Str("period", req.GetPeriod()).
		Int64("user_id", req.UserId).
		Strs("roles", req.Roles).
		Msg("GetStorefrontStats RPC called")

	// Validate request basics
	if req == nil {
		s.logger.Warn().Msg("nil request received")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	if req.StorefrontId <= 0 {
		s.logger.Warn().Int64("storefront_id", req.StorefrontId).Msg("invalid storefront_id")
		return nil, status.Error(codes.InvalidArgument, "storefront_id must be positive")
	}

	if req.UserId <= 0 {
		s.logger.Warn().Int64("user_id", req.UserId).Msg("invalid user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Call service (authorization happens inside service layer)
	response, err := s.storefrontAnalyticsService.GetStorefrontStats(ctx, req)
	if err != nil {
		// Map service errors to gRPC status codes
		errMsg := err.Error()

		if contains(errMsg, "permission") || contains(errMsg, "unauthorized") || contains(errMsg, "owner") {
			s.logger.Warn().
				Err(err).
				Int64("user_id", req.UserId).
				Int64("storefront_id", req.StorefrontId).
				Msg("permission denied for storefront stats")
			return nil, status.Error(codes.PermissionDenied, "you don't have permission to view this storefront's analytics")
		}

		if contains(errMsg, "invalid") || contains(errMsg, "required") {
			s.logger.Warn().
				Err(err).
				Int64("storefront_id", req.StorefrontId).
				Msg("validation error for storefront stats")
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid input: %v", err))
		}

		if contains(errMsg, "not found") {
			s.logger.Warn().
				Err(err).
				Int64("storefront_id", req.StorefrontId).
				Msg("storefront not found")
			return nil, status.Error(codes.NotFound, fmt.Sprintf("storefront not found: %d", req.StorefrontId))
		}

		s.logger.Error().
			Err(err).
			Int64("user_id", req.UserId).
			Int64("storefront_id", req.StorefrontId).
			Str("period", req.GetPeriod()).
			Msg("service error getting storefront stats")
		return nil, status.Error(codes.Internal, "failed to retrieve storefront analytics")
	}

	// Log success
	s.logger.Info().
		Int64("user_id", req.UserId).
		Int64("storefront_id", req.StorefrontId).
		Str("storefront_name", response.StorefrontName).
		Int64("total_sales", response.TotalSales).
		Float64("total_revenue", response.TotalRevenue).
		Int32("active_listings", response.ActiveListings).
		Msg("GetStorefrontStats completed successfully")

	return response, nil
}

// GetTrendingStats retrieves platform trending analytics (admin only)
func (s *Server) GetTrendingStats(
	ctx context.Context,
	req *listingspb.GetTrendingStatsRequest,
) (*listingspb.GetTrendingStatsResponse, error) {
	// Log request
	s.logger.Info().Msg("GetTrendingStats RPC called")

	// No validation needed - request is empty (admin-only, no parameters)

	// Call service
	response, err := s.analyticsService.GetTrendingStats(ctx, req)
	if err != nil {
		// Map service errors to gRPC status codes
		errMsg := err.Error()

		if contains(errMsg, "permission") || contains(errMsg, "unauthorized") || contains(errMsg, "admin") {
			s.logger.Warn().
				Err(err).
				Msg("permission denied for trending stats")
			return nil, status.Error(codes.PermissionDenied, "admin access required")
		}

		if contains(errMsg, "invalid") || contains(errMsg, "required") {
			s.logger.Warn().
				Err(err).
				Msg("validation error for trending stats")
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid input: %v", err))
		}

		s.logger.Error().
			Err(err).
			Msg("service error getting trending stats")
		return nil, status.Error(codes.Internal, "failed to retrieve trending analytics")
	}

	// Log success
	s.logger.Info().
		Int("trending_categories", len(response.TrendingCategories)).
		Int("hot_listings", len(response.HotListings)).
		Int("popular_searches", len(response.PopularSearches)).
		Msg("GetTrendingStats completed successfully")

	return response, nil
}

// validateOverviewStatsRequest validates GetOverviewStatsRequest
func (s *Server) validateOverviewStatsRequest(req *listingspb.GetOverviewStatsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate date_from and date_to are present
	if req.DateFrom == nil {
		return fmt.Errorf("date_from is required")
	}
	if req.DateTo == nil {
		return fmt.Errorf("date_to is required")
	}

	// Validate date range is valid
	dateFrom := req.DateFrom.AsTime()
	dateTo := req.DateTo.AsTime()

	if dateFrom.After(dateTo) {
		return fmt.Errorf("date_from must be before or equal to date_to")
	}

	// Validate listing_type if provided
	if req.ListingType != nil && *req.ListingType != "" {
		listingType := *req.ListingType
		if listingType != "b2c" && listingType != "c2c" {
			return fmt.Errorf("listing_type must be either 'b2c' or 'c2c'")
		}
	}

	return nil
}

// validateListingStatsRequest validates GetListingStatsRequest
func (s *Server) validateListingStatsRequest(req *listingspb.GetListingStatsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Must have listing_id or product_id
	if req.GetListingId() == 0 && req.GetProductId() == 0 {
		return fmt.Errorf("either listing_id or product_id must be provided")
	}

	// Validate date_from and date_to are present
	if req.DateFrom == nil {
		return fmt.Errorf("date_from is required")
	}
	if req.DateTo == nil {
		return fmt.Errorf("date_to is required")
	}

	// Validate date range is valid
	dateFrom := req.DateFrom.AsTime()
	dateTo := req.DateTo.AsTime()

	if dateFrom.After(dateTo) {
		return fmt.Errorf("date_from must be before or equal to date_to")
	}

	return nil
}

// extractAuthFromMetadata extracts user_id and is_admin from gRPC metadata
func (s *Server) extractAuthFromMetadata(ctx context.Context) (int64, bool, error) {
	// Need to import metadata first
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Debug().Msg("no metadata found in gRPC context")
		return 0, false, fmt.Errorf("no authentication metadata")
	}

	// Extract user_id from metadata
	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		s.logger.Debug().Msg("user_id not found in metadata")
		return 0, false, fmt.Errorf("user_id not found in metadata")
	}

	// Parse user_id
	var userID int64
	_, err := fmt.Sscanf(userIDs[0], "%d", &userID)
	if err != nil {
		s.logger.Warn().
			Err(err).
			Str("user_id_str", userIDs[0]).
			Msg("failed to parse user_id from metadata")
		return 0, false, fmt.Errorf("invalid user_id format")
	}

	// Extract is_admin from metadata (optional)
	isAdmin := false
	adminValues := md.Get("is_admin")
	if len(adminValues) > 0 {
		isAdmin = adminValues[0] == "true" || adminValues[0] == "1"
	}

	// Also check roles in metadata
	roles := md.Get("roles")
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	s.logger.Debug().
		Int64("user_id", userID).
		Bool("is_admin", isAdmin).
		Msg("auth extracted from metadata")

	return userID, isAdmin, nil
}
