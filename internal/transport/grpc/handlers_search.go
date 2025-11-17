package grpc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	searchv1 "github.com/sveturs/listings/api/proto/search/v1"
	"github.com/sveturs/listings/internal/service/search"
)

// SearchServiceInterface defines the interface for search service
type SearchServiceInterface interface {
	SearchListings(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error)
	GetSearchFacets(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error)
	SearchWithFilters(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error)
	GetSuggestions(ctx context.Context, req *search.SuggestionsRequest) (*search.SuggestionsResponse, error)
	GetPopularSearches(ctx context.Context, req *search.PopularSearchesRequest) (*search.PopularSearchesResponse, error)
	GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]search.ListingSearchResult, int64, error)
}

// SearchHandler implements SearchService gRPC service
type SearchHandler struct {
	searchv1.UnimplementedSearchServiceServer
	service SearchServiceInterface
	logger  zerolog.Logger
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *search.Service, logger zerolog.Logger) *SearchHandler {
	return &SearchHandler{
		service: SearchServiceInterface(service),
		logger:  logger.With().Str("handler", "search").Logger(),
	}
}

// SearchListings implements SearchService.SearchListings RPC
func (h *SearchHandler) SearchListings(
	ctx context.Context,
	req *searchv1.SearchListingsRequest,
) (*searchv1.SearchListingsResponse, error) {
	// Log request
	h.logger.Info().
		Str("query", req.Query).
		Interface("category_id", req.CategoryId).
		Int32("limit", req.Limit).
		Int32("offset", req.Offset).
		Bool("use_cache", req.UseCache).
		Msg("SearchListings RPC called")

	// Validate request
	if err := h.validateSearchRequest(req); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("invalid search request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain
	domainReq := h.protoToDomainRequest(req)

	// Execute search
	result, err := h.service.SearchListings(ctx, domainReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("search failed")
		return nil, status.Error(codes.Internal, "search failed")
	}

	// Convert domain to proto
	protoResp := h.domainToProtoResponse(result)

	h.logger.Info().
		Int64("total", protoResp.Total).
		Int32("took_ms", protoResp.TookMs).
		Bool("cached", protoResp.Cached).
		Int("results", len(protoResp.Listings)).
		Msg("SearchListings completed")

	return protoResp, nil
}

// ============================================================================
// PHASE 21.2: Advanced Search Handlers
// ============================================================================

// GetSearchFacets returns aggregations for building filter UI
func (h *SearchHandler) GetSearchFacets(
	ctx context.Context,
	req *searchv1.GetSearchFacetsRequest,
) (*searchv1.GetSearchFacetsResponse, error) {
	// Log request
	h.logger.Info().
		Interface("query", req.Query).
		Interface("category_id", req.CategoryId).
		Msg("GetSearchFacets RPC called")

	// Convert proto to domain
	domainReq := ProtoToFacetsRequest(req)

	// Validate (service layer also validates, but double-check)
	if err := domainReq.Validate(); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("invalid facets request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Call service
	result, err := h.service.GetSearchFacets(ctx, domainReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("facets service failed")

		// Map service errors to gRPC status codes
		if containsError(err, "invalid") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to get search facets")
	}

	// Convert domain to proto
	resp := FacetsResponseToProto(result)

	h.logger.Info().
		Int32("took_ms", resp.TookMs).
		Bool("cached", resp.Cached).
		Int("categories", len(resp.Categories)).
		Int("price_ranges", len(resp.PriceRanges)).
		Int("attributes", len(resp.Attributes)).
		Msg("GetSearchFacets completed")

	return resp, nil
}

// SearchWithFilters performs enhanced search with multiple filters
func (h *SearchHandler) SearchWithFilters(
	ctx context.Context,
	req *searchv1.SearchWithFiltersRequest,
) (*searchv1.SearchWithFiltersResponse, error) {
	// Log request
	h.logger.Info().
		Str("query", req.Query).
		Interface("category_id", req.CategoryId).
		Int32("limit", req.Limit).
		Int32("offset", req.Offset).
		Bool("use_cache", req.UseCache).
		Bool("include_facets", req.IncludeFacets).
		Msg("SearchWithFilters RPC called")

	// Convert proto to domain
	domainReq := ProtoToSearchFiltersRequest(req)

	// Validate
	if err := domainReq.Validate(); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("invalid filtered search request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Call service
	result, err := h.service.SearchWithFilters(ctx, domainReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("filtered search service failed")

		// Map service errors to gRPC status codes
		if containsError(err, "invalid") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to search with filters")
	}

	// Convert domain to proto
	resp := SearchFiltersResponseToProto(result)

	h.logger.Info().
		Int("count", len(resp.Listings)).
		Int64("total", resp.Total).
		Int32("took_ms", resp.TookMs).
		Bool("cached", resp.Cached).
		Bool("has_facets", resp.Facets != nil).
		Msg("SearchWithFilters completed")

	return resp, nil
}

// GetSuggestions provides autocomplete suggestions
func (h *SearchHandler) GetSuggestions(
	ctx context.Context,
	req *searchv1.GetSuggestionsRequest,
) (*searchv1.GetSuggestionsResponse, error) {
	// Log request
	h.logger.Info().
		Str("prefix", req.Prefix).
		Interface("category_id", req.CategoryId).
		Int32("limit", req.Limit).
		Msg("GetSuggestions RPC called")

	// Early validation (prefix length check)
	if len(req.Prefix) < 2 {
		h.logger.Warn().
			Str("prefix", req.Prefix).
			Msg("prefix too short")
		return nil, status.Error(codes.InvalidArgument, "prefix must be at least 2 characters")
	}

	// Convert proto to domain
	domainReq := ProtoToSuggestionsRequest(req)

	// Validate
	if err := domainReq.Validate(); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("invalid suggestions request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Call service
	result, err := h.service.GetSuggestions(ctx, domainReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("suggestions service failed")

		// Map service errors to gRPC status codes
		if containsError(err, "invalid") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to get suggestions")
	}

	// Convert domain to proto
	resp := SuggestionsResponseToProto(result)

	h.logger.Info().
		Str("prefix", req.Prefix).
		Int("count", len(resp.Suggestions)).
		Int32("took_ms", resp.TookMs).
		Bool("cached", resp.Cached).
		Msg("GetSuggestions completed")

	return resp, nil
}

// GetPopularSearches returns trending search queries
func (h *SearchHandler) GetPopularSearches(
	ctx context.Context,
	req *searchv1.GetPopularSearchesRequest,
) (*searchv1.GetPopularSearchesResponse, error) {
	// Log request
	h.logger.Info().
		Interface("category_id", req.CategoryId).
		Interface("time_range", req.TimeRange).
		Int32("limit", req.Limit).
		Msg("GetPopularSearches RPC called")

	// Convert proto to domain
	domainReq := ProtoToPopularSearchesRequest(req)

	// Validate
	if err := domainReq.Validate(); err != nil {
		h.logger.Warn().
			Err(err).
			Msg("invalid popular searches request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Call service
	result, err := h.service.GetPopularSearches(ctx, domainReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("popular searches service failed")

		// Map service errors to gRPC status codes
		if containsError(err, "invalid") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to get popular searches")
	}

	// Convert domain to proto
	resp := PopularSearchesResponseToProto(result)

	h.logger.Info().
		Int("count", len(resp.Searches)).
		Int32("took_ms", resp.TookMs).
		Msg("GetPopularSearches completed")

	return resp, nil
}

// containsError checks if error message contains substring (case-insensitive)
func containsError(err error, substr string) bool {
	if err == nil || substr == "" {
		return false
	}
	errMsg := err.Error()
	return contains(errMsg, substr)
}

// validateSearchRequest validates search request parameters
func (h *SearchHandler) validateSearchRequest(req *searchv1.SearchListingsRequest) error {
	// Validate limit (1-100)
	if req.Limit < 0 {
		return fmt.Errorf("limit must be >= 0")
	}
	if req.Limit > 100 {
		return fmt.Errorf("limit must be <= 100")
	}

	// Validate offset (>= 0)
	if req.Offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}

	// Validate query length
	if len(req.Query) > 500 {
		return fmt.Errorf("query too long (max 500 characters)")
	}

	return nil
}

// protoToDomainRequest converts proto request to domain request
func (h *SearchHandler) protoToDomainRequest(req *searchv1.SearchListingsRequest) *search.SearchRequest {
	domainReq := &search.SearchRequest{
		Query:    req.Query,
		Limit:    req.Limit,
		Offset:   req.Offset,
		UseCache: req.UseCache,
	}

	// Handle optional category_id
	if req.CategoryId != nil {
		categoryID := *req.CategoryId
		domainReq.CategoryID = &categoryID
	}

	// Set defaults
	if domainReq.Limit == 0 {
		domainReq.Limit = 20
	}

	return domainReq
}

// domainToProtoResponse converts domain response to proto response
func (h *SearchHandler) domainToProtoResponse(result *search.SearchResponse) *searchv1.SearchListingsResponse {
	protoListings := make([]*searchv1.Listing, 0, len(result.Listings))

	for _, listing := range result.Listings {
		protoListing := &searchv1.Listing{
			Id:          listing.ID,
			Uuid:        listing.UUID,
			Title:       listing.Title,
			Price:       listing.Price,
			Currency:    listing.Currency,
			CategoryId:  listing.CategoryID,
			Status:      listing.Status,
			CreatedAt:   listing.CreatedAt,
			UserId:      listing.UserID,
			Quantity:    listing.Quantity,
			SourceType:  listing.SourceType,
			StockStatus: listing.StockStatus,
		}

		// Add optional fields
		if listing.Description != nil {
			protoListing.Description = listing.Description
		}
		if listing.StorefrontID != nil {
			protoListing.StorefrontId = listing.StorefrontID
		}
		if listing.SKU != nil {
			protoListing.Sku = listing.SKU
		}

		// Add images
		if len(listing.Images) > 0 {
			protoImages := make([]*searchv1.ListingImage, 0, len(listing.Images))
			for _, img := range listing.Images {
				protoImages = append(protoImages, &searchv1.ListingImage{
					Id:           img.ID,
					Url:          img.URL,
					IsPrimary:    img.IsPrimary,
					DisplayOrder: img.DisplayOrder,
				})
			}
			protoListing.Images = protoImages
		}

		protoListings = append(protoListings, protoListing)
	}

	return &searchv1.SearchListingsResponse{
		Listings: protoListings,
		Total:    result.Total,
		TookMs:   result.TookMs,
		Cached:   result.Cached,
	}
}
