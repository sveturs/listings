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

// SearchHandler implements SearchService gRPC service
type SearchHandler struct {
	searchv1.UnimplementedSearchServiceServer
	service *search.Service
	logger  zerolog.Logger
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *search.Service, logger zerolog.Logger) *SearchHandler {
	return &SearchHandler{
		service: service,
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
