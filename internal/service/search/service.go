package search

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/cache"
	"github.com/sveturs/listings/internal/opensearch"
)

// Service provides search functionality for listings
type Service struct {
	searchClient *opensearch.SearchClient
	cache        *cache.SearchCache
	logger       zerolog.Logger
}

// NewService creates a new search service
func NewService(
	searchClient *opensearch.SearchClient,
	cache *cache.SearchCache,
	logger zerolog.Logger,
) *Service {
	return &Service{
		searchClient: searchClient,
		cache:        cache,
		logger:       logger.With().Str("service", "search").Logger(),
	}
}

// SearchListings searches for listings based on query and filters
func (s *Service) SearchListings(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	start := time.Now()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.logger.Debug().
		Str("query", req.Query).
		Interface("category_id", req.CategoryID).
		Int32("limit", req.Limit).
		Int32("offset", req.Offset).
		Bool("use_cache", req.UseCache).
		Msg("searching listings")

	// Try cache first (if enabled)
	if req.UseCache && s.cache != nil {
		cacheReq := &cache.SearchRequest{
			Query:      req.Query,
			CategoryID: req.CategoryID,
			Limit:      req.Limit,
			Offset:     req.Offset,
		}

		if cached, err := s.cache.Get(ctx, cacheReq); err == nil {
			// Convert cached result to response
			response := &SearchResponse{
				Listings: s.convertCachedListings(cached.Listings),
				Total:    cached.Total,
				TookMs:   cached.TookMs,
				Cached:   true,
			}

			s.logger.Debug().
				Dur("duration", time.Since(start)).
				Int64("total", response.Total).
				Msg("returned cached search results")

			return response, nil
		}
		// Cache miss is fine, continue to OpenSearch
	}

	// Build OpenSearch query
	query := s.buildSearchQuery(req)

	// Execute search
	searchResp, err := s.searchClient.Search(ctx, query)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("OpenSearch query failed")
		return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}

	// Parse results
	listings := s.parseSearchResults(searchResp)

	response := &SearchResponse{
		Listings: listings,
		Total:    searchResp.Hits.Total.Value,
		TookMs:   int32(searchResp.Took),
		Cached:   false,
	}

	// Cache result (async, non-blocking)
	if req.UseCache && s.cache != nil {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			cacheReq := &cache.SearchRequest{
				Query:      req.Query,
				CategoryID: req.CategoryID,
				Limit:      req.Limit,
				Offset:     req.Offset,
			}

			cacheResult := &cache.SearchResult{
				Listings: s.convertListingsForCache(listings),
				Total:    response.Total,
				TookMs:   response.TookMs,
			}

			if err := s.cache.Set(cacheCtx, cacheReq, cacheResult); err != nil {
				s.logger.Warn().Err(err).Msg("failed to cache search results")
			}
		}()
	}

	s.logger.Info().
		Dur("duration", time.Since(start)).
		Int64("total", response.Total).
		Int32("took_ms", response.TookMs).
		Int("results", len(response.Listings)).
		Msg("search completed")

	return response, nil
}

// buildSearchQuery constructs OpenSearch query from request
func (s *Service) buildSearchQuery(req *SearchRequest) map[string]interface{} {
	// Build bool query
	mustClauses := []map[string]interface{}{
		// Filter by active status
		{
			"term": map[string]interface{}{
				"status": "active",
			},
		},
	}

	// Add text search if query provided
	if req.Query != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"title^3", "description"},
				"type":   "best_fields",
			},
		})
	}

	// Add category filter if provided
	if req.CategoryID != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": *req.CategoryID,
			},
		})
	}

	// Build final query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
		"size": req.Limit,
		"from": req.Offset,
		"sort": []map[string]interface{}{
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return query
}

// parseSearchResults converts OpenSearch hits to domain listings
func (s *Service) parseSearchResults(resp *opensearch.SearchResponse) []ListingSearchResult {
	listings := make([]ListingSearchResult, 0, len(resp.Hits.Hits))

	for _, hit := range resp.Hits.Hits {
		listing := s.parseListingFromHit(hit.Source)
		listings = append(listings, listing)
	}

	return listings
}

// parseListingFromHit parses a single listing from OpenSearch hit source
func (s *Service) parseListingFromHit(source map[string]interface{}) ListingSearchResult {
	listing := ListingSearchResult{}

	// Parse required fields
	if id, ok := source["id"].(float64); ok {
		listing.ID = int64(id)
	}
	if uuid, ok := source["uuid"].(string); ok {
		listing.UUID = uuid
	}
	if title, ok := source["title"].(string); ok {
		listing.Title = title
	}
	if price, ok := source["price"].(float64); ok {
		listing.Price = price
	}
	if currency, ok := source["currency"].(string); ok {
		listing.Currency = currency
	}
	if categoryID, ok := source["category_id"].(float64); ok {
		listing.CategoryID = int64(categoryID)
	}
	if status, ok := source["status"].(string); ok {
		listing.Status = status
	}
	if createdAt, ok := source["created_at"].(string); ok {
		listing.CreatedAt = createdAt
	}
	if userID, ok := source["user_id"].(float64); ok {
		listing.UserID = int64(userID)
	}
	if quantity, ok := source["quantity"].(float64); ok {
		listing.Quantity = int32(quantity)
	}
	if sourceType, ok := source["source_type"].(string); ok {
		listing.SourceType = sourceType
	}
	if stockStatus, ok := source["stock_status"].(string); ok {
		listing.StockStatus = stockStatus
	}

	// Parse optional fields
	if desc, ok := source["description"].(string); ok && desc != "" {
		listing.Description = &desc
	}
	if storefrontID, ok := source["storefront_id"].(float64); ok {
		id := int64(storefrontID)
		listing.StorefrontID = &id
	}
	if sku, ok := source["sku"].(string); ok && sku != "" {
		listing.SKU = &sku
	}

	// Parse images
	if imagesData, ok := source["images"].([]interface{}); ok {
		listing.Images = s.parseImages(imagesData)
	}

	return listing
}

// parseImages parses images from OpenSearch source
func (s *Service) parseImages(imagesData []interface{}) []ListingImageResult {
	images := make([]ListingImageResult, 0, len(imagesData))

	for _, imgData := range imagesData {
		if imgMap, ok := imgData.(map[string]interface{}); ok {
			img := ListingImageResult{}

			if id, ok := imgMap["id"].(float64); ok {
				img.ID = int64(id)
			}
			if url, ok := imgMap["url"].(string); ok {
				img.URL = url
			}
			if isPrimary, ok := imgMap["is_primary"].(bool); ok {
				img.IsPrimary = isPrimary
			}
			if displayOrder, ok := imgMap["display_order"].(float64); ok {
				img.DisplayOrder = int32(displayOrder)
			}

			images = append(images, img)
		}
	}

	return images
}

// convertCachedListings converts cached listings to domain type
func (s *Service) convertCachedListings(cached []map[string]interface{}) []ListingSearchResult {
	listings := make([]ListingSearchResult, 0, len(cached))

	for _, item := range cached {
		listing := s.parseListingFromHit(item)
		listings = append(listings, listing)
	}

	return listings
}

// convertListingsForCache converts listings to cacheable format
func (s *Service) convertListingsForCache(listings []ListingSearchResult) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(listings))

	for _, listing := range listings {
		item := map[string]interface{}{
			"id":           listing.ID,
			"uuid":         listing.UUID,
			"title":        listing.Title,
			"price":        listing.Price,
			"currency":     listing.Currency,
			"category_id":  listing.CategoryID,
			"status":       listing.Status,
			"created_at":   listing.CreatedAt,
			"user_id":      listing.UserID,
			"quantity":     listing.Quantity,
			"source_type":  listing.SourceType,
			"stock_status": listing.StockStatus,
		}

		if listing.Description != nil {
			item["description"] = *listing.Description
		}
		if listing.StorefrontID != nil {
			item["storefront_id"] = *listing.StorefrontID
		}
		if listing.SKU != nil {
			item["sku"] = *listing.SKU
		}
		if len(listing.Images) > 0 {
			item["images"] = listing.Images
		}

		result = append(result, item)
	}

	return result
}
