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

// ============================================================================
// PHASE 21.2: Advanced Search Service Methods
// ============================================================================

// GetSearchFacets returns aggregations for building filter UI
func (s *Service) GetSearchFacets(ctx context.Context, req *FacetsRequest) (*FacetsResponse, error) {
	start := time.Now()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.logger.Debug().
		Str("query", req.Query).
		Interface("category_id", req.CategoryID).
		Bool("use_cache", req.UseCache).
		Msg("fetching search facets")

	// Generate cache key
	cacheKey := ""
	if req.UseCache && s.cache != nil {
		filters := make(map[string]interface{})
		if req.Filters != nil {
			// Convert SearchFilters to map for cache key
			if req.Filters.Price != nil {
				filters["price"] = req.Filters.Price
			}
			if req.Filters.Attributes != nil {
				filters["attributes"] = req.Filters.Attributes
			}
			if req.Filters.Location != nil {
				filters["location"] = req.Filters.Location
			}
			if req.Filters.SourceType != nil {
				filters["source_type"] = *req.Filters.SourceType
			}
			if req.Filters.StockStatus != nil {
				filters["stock_status"] = *req.Filters.StockStatus
			}
		}
		cacheKey = s.cache.GenerateFacetsKey(req.Query, req.CategoryID, filters)

		// Check cache
		if cached, err := s.cache.GetFacets(ctx, cacheKey); err == nil && cached != nil {
			s.logger.Debug().Msg("facets cache hit")
			return s.convertCachedFacets(cached, true), nil
		}
	}

	// Build OpenSearch query
	query := BuildFacetsQuery(req)

	// Execute search
	result, err := s.searchClient.Search(ctx, query)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("facets query failed")
		return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}

	// Parse aggregations
	facets, err := s.parseAggregations(result)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("failed to parse aggregations")
		return nil, fmt.Errorf("failed to parse aggregations: %w", err)
	}

	facets.TookMs = int32(result.Took)
	facets.Cached = false

	// Cache result (async, non-blocking)
	if req.UseCache && s.cache != nil && cacheKey != "" {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			facetsMap := s.convertFacetsForCache(facets)
			if err := s.cache.SetFacets(cacheCtx, cacheKey, facetsMap); err != nil {
				s.logger.Warn().Err(err).Msg("failed to cache facets")
			}
		}()
	}

	s.logger.Info().
		Dur("duration", time.Since(start)).
		Int32("took_ms", facets.TookMs).
		Int("categories", len(facets.Categories)).
		Int("price_ranges", len(facets.PriceRanges)).
		Int("attributes", len(facets.Attributes)).
		Msg("facets fetched successfully")

	return facets, nil
}

// SearchWithFilters performs enhanced search with multiple filters
func (s *Service) SearchWithFilters(ctx context.Context, req *SearchFiltersRequest) (*SearchFiltersResponse, error) {
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
		Bool("include_facets", req.IncludeFacets).
		Msg("searching with filters")

	// Generate cache key
	cacheKey := ""
	if req.UseCache && s.cache != nil {
		filters := make(map[string]interface{})
		if req.Filters != nil {
			if req.Filters.Price != nil {
				filters["price"] = req.Filters.Price
			}
			if req.Filters.Attributes != nil {
				filters["attributes"] = req.Filters.Attributes
			}
			if req.Filters.Location != nil {
				filters["location"] = req.Filters.Location
			}
			if req.Filters.SourceType != nil {
				filters["source_type"] = *req.Filters.SourceType
			}
			if req.Filters.StockStatus != nil {
				filters["stock_status"] = *req.Filters.StockStatus
			}
		}

		sort := make(map[string]string)
		if req.Sort != nil {
			sort["field"] = req.Sort.Field
			sort["order"] = req.Sort.Order
		}

		cacheKey = s.cache.GenerateFilteredKey(req.Query, req.CategoryID, filters, sort, req.Limit, req.Offset)

		// Check cache
		if cached, err := s.cache.GetFiltered(ctx, cacheKey); err == nil && cached != nil {
			s.logger.Debug().Msg("filtered search cache hit")
			return s.convertCachedFilteredSearch(cached, true), nil
		}
	}

	// Build OpenSearch query
	query := BuildFilteredSearchQuery(req)

	// Execute search
	result, err := s.searchClient.Search(ctx, query)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("filtered search query failed")
		return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}

	// Parse results
	listings := s.parseSearchResults(result)

	response := &SearchFiltersResponse{
		Listings: listings,
		Total:    result.Hits.Total.Value,
		TookMs:   int32(result.Took),
		Cached:   false,
	}

	// Parse facets if requested
	if req.IncludeFacets {
		facets, err := s.parseAggregations(result)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Msg("failed to parse facets (continuing without facets)")
		} else {
			facets.TookMs = int32(result.Took)
			facets.Cached = false
			response.Facets = facets
		}
	}

	// Cache result (async, non-blocking)
	if req.UseCache && s.cache != nil && cacheKey != "" {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			filteredMap := s.convertFilteredSearchForCache(response)
			if err := s.cache.SetFiltered(cacheCtx, cacheKey, filteredMap); err != nil {
				s.logger.Warn().Err(err).Msg("failed to cache filtered search")
			}
		}()
	}

	s.logger.Info().
		Dur("duration", time.Since(start)).
		Int64("total", response.Total).
		Int32("took_ms", response.TookMs).
		Int("results", len(response.Listings)).
		Msg("filtered search completed")

	return response, nil
}

// GetSuggestions provides autocomplete suggestions
func (s *Service) GetSuggestions(ctx context.Context, req *SuggestionsRequest) (*SuggestionsResponse, error) {
	start := time.Now()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.logger.Debug().
		Str("prefix", req.Prefix).
		Interface("category_id", req.CategoryID).
		Int32("limit", req.Limit).
		Bool("use_cache", req.UseCache).
		Msg("fetching suggestions")

	// Generate cache key
	cacheKey := ""
	if req.UseCache && s.cache != nil {
		cacheKey = s.cache.GenerateSuggestionsKey(req.Prefix, req.CategoryID)

		// Check cache
		if cached, err := s.cache.GetSuggestions(ctx, cacheKey); err == nil && cached != nil {
			s.logger.Debug().Msg("suggestions cache hit")
			return s.convertCachedSuggestions(cached, true), nil
		}
	}

	// Build completion suggester query
	query := BuildSuggestionsQuery(req)

	// Execute search
	result, err := s.searchClient.Search(ctx, query)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("suggestions query failed")
		return nil, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}

	// Parse suggestions
	suggestions := s.parseSuggestions(result)

	response := &SuggestionsResponse{
		Suggestions: suggestions,
		TookMs:      int32(result.Took),
		Cached:      false,
	}

	// Cache result (async, non-blocking) - 1 hour TTL for stable data
	if req.UseCache && s.cache != nil && cacheKey != "" {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			suggestionsMap := s.convertSuggestionsForCache(response)
			if err := s.cache.SetSuggestions(cacheCtx, cacheKey, suggestionsMap); err != nil {
				s.logger.Warn().Err(err).Msg("failed to cache suggestions")
			}
		}()
	}

	s.logger.Info().
		Dur("duration", time.Since(start)).
		Int32("took_ms", response.TookMs).
		Int("suggestions", len(response.Suggestions)).
		Msg("suggestions fetched successfully")

	return response, nil
}

// GetPopularSearches returns trending search queries
// NOTE: This is a mock implementation until analytics tracking is built
func (s *Service) GetPopularSearches(ctx context.Context, req *PopularSearchesRequest) (*PopularSearchesResponse, error) {
	start := time.Now()

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	s.logger.Debug().
		Interface("category_id", req.CategoryID).
		Int32("limit", req.Limit).
		Str("time_range", req.TimeRange).
		Msg("fetching popular searches")

	// Generate cache key (15 min TTL - trending data)
	cacheKey := ""
	if s.cache != nil {
		cacheKey = s.cache.GeneratePopularKey(req.CategoryID, req.TimeRange)

		// Check cache
		if cached, err := s.cache.GetPopular(ctx, cacheKey); err == nil && cached != nil {
			s.logger.Debug().Msg("popular searches cache hit")
			return s.convertCachedPopularSearches(cached), nil
		}
	}

	// TODO: Query PostgreSQL search_queries table when implemented
	// For now, return mock trending searches based on category
	searches := s.getMockPopularSearches(req.CategoryID, req.Limit)

	response := &PopularSearchesResponse{
		Searches: searches,
		TookMs:   int32(time.Since(start).Milliseconds()),
	}

	// Cache results
	if s.cache != nil && cacheKey != "" {
		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			popularMap := s.convertPopularSearchesForCache(response)
			if err := s.cache.SetPopular(cacheCtx, cacheKey, popularMap); err != nil {
				s.logger.Warn().Err(err).Msg("failed to cache popular searches")
			}
		}()
	}

	s.logger.Info().
		Dur("duration", time.Since(start)).
		Int32("took_ms", response.TookMs).
		Int("searches", len(response.Searches)).
		Msg("popular searches fetched (mock data)")

	return response, nil
}

// GetSimilarListings finds similar listings using More Like This query (ENHANCED)
func (s *Service) GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]ListingSearchResult, int64, error) {
	start := time.Now()

	s.logger.Debug().
		Int64("listing_id", listingID).
		Int32("limit", limit).
		Msg("fetching similar listings")

	// Fetch original listing to get category
	original, err := s.getListingByID(ctx, listingID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("listing_id", listingID).
			Msg("failed to fetch original listing")
		return nil, 0, fmt.Errorf("failed to fetch original listing: %w", err)
	}

	// Build More Like This query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"more_like_this": map[string]interface{}{
							"fields": []string{"title^3", "description", "attributes_searchable_text"},
							"like": []map[string]interface{}{
								{"_id": fmt.Sprintf("%d", listingID)},
							},
							"min_term_freq":   1,
							"min_doc_freq":    2,
							"max_query_terms": 12,
						},
					},
				},
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"status": "active"}},
					{"term": map[string]interface{}{"category_id": original.CategoryID}},
					{"bool": map[string]interface{}{
						"must_not": []map[string]interface{}{
							{"term": map[string]interface{}{"id": listingID}},
						},
					}},
				},
			},
		},
		"size": limit,
	}

	// Execute search
	result, err := s.searchClient.Search(ctx, query)
	if err != nil {
		s.logger.Error().
			Err(err).
			Int64("listing_id", listingID).
			Msg("similar listings query failed")
		return nil, 0, fmt.Errorf("%w: %v", ErrSearchFailed, err)
	}

	// Parse results
	listings := s.parseSearchResults(result)

	s.logger.Info().
		Dur("duration", time.Since(start)).
		Int64("listing_id", listingID).
		Int64("total", result.Hits.Total.Value).
		Int("results", len(listings)).
		Msg("similar listings fetched successfully")

	return listings, result.Hits.Total.Value, nil
}

// ============================================================================
// Helper Methods for Phase 21.2
// ============================================================================

// parseAggregations extracts facets from OpenSearch aggregations
func (s *Service) parseAggregations(result *opensearch.SearchResponse) (*FacetsResponse, error) {
	facets := &FacetsResponse{
		Categories:    []CategoryFacet{},
		PriceRanges:   []PriceRangeFacet{},
		Attributes:    make(map[string]AttributeFacet),
		SourceTypes:   []Facet{},
		StockStatuses: []Facet{},
	}

	if result.Aggregations == nil || len(result.Aggregations) == 0 {
		s.logger.Debug().Msg("no aggregations in response")
		return facets, nil
	}

	// Parse categories aggregation
	if categoriesAgg, ok := result.Aggregations["categories"].(map[string]interface{}); ok {
		if buckets, ok := categoriesAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					categoryID, _ := b["key"].(float64)
					docCount, _ := b["doc_count"].(float64)
					facets.Categories = append(facets.Categories, CategoryFacet{
						CategoryID: int64(categoryID),
						Count:      int64(docCount),
					})
				}
			}
		}
	}

	// Parse price_ranges aggregation (histogram)
	if priceRangesAgg, ok := result.Aggregations["price_ranges"].(map[string]interface{}); ok {
		if buckets, ok := priceRangesAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					key, _ := b["key"].(float64)
					docCount, _ := b["doc_count"].(float64)
					if docCount > 0 { // Only include non-empty buckets
						facets.PriceRanges = append(facets.PriceRanges, PriceRangeFacet{
							Min:   key,
							Max:   key + 100, // Interval is 100 (from buildAggregations)
							Count: int64(docCount),
						})
					}
				}
			}
		}
	}

	// Parse source_types aggregation
	if sourceTypesAgg, ok := result.Aggregations["source_types"].(map[string]interface{}); ok {
		if buckets, ok := sourceTypesAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					key, _ := b["key"].(string)
					docCount, _ := b["doc_count"].(float64)
					facets.SourceTypes = append(facets.SourceTypes, Facet{
						Key:   key,
						Count: int64(docCount),
					})
				}
			}
		}
	}

	// Parse stock_statuses aggregation
	if stockStatusesAgg, ok := result.Aggregations["stock_statuses"].(map[string]interface{}); ok {
		if buckets, ok := stockStatusesAgg["buckets"].([]interface{}); ok {
			for _, bucket := range buckets {
				if b, ok := bucket.(map[string]interface{}); ok {
					key, _ := b["key"].(string)
					docCount, _ := b["doc_count"].(float64)
					facets.StockStatuses = append(facets.StockStatuses, Facet{
						Key:   key,
						Count: int64(docCount),
					})
				}
			}
		}
	}

	// Parse attributes aggregation (nested)
	if attributesAgg, ok := result.Aggregations["attributes"].(map[string]interface{}); ok {
		if attributeKeys, ok := attributesAgg["attribute_keys"].(map[string]interface{}); ok {
			if buckets, ok := attributeKeys["buckets"].([]interface{}); ok {
				for _, bucket := range buckets {
					if b, ok := bucket.(map[string]interface{}); ok {
						attributeCode, _ := b["key"].(string)
						if attributeCode == "" {
							continue
						}

						// Parse attribute values sub-aggregation
						values := []AttributeValueCount{}
						if attributeValues, ok := b["attribute_values"].(map[string]interface{}); ok {
							if valueBuckets, ok := attributeValues["buckets"].([]interface{}); ok {
								for _, valueBucket := range valueBuckets {
									if vb, ok := valueBucket.(map[string]interface{}); ok {
										value, _ := vb["key"].(string)
										docCount, _ := vb["doc_count"].(float64)
										if value != "" {
											values = append(values, AttributeValueCount{
												Value: value,
												Count: int64(docCount),
											})
										}
									}
								}
							}
						}

						if len(values) > 0 {
							facets.Attributes[attributeCode] = AttributeFacet{
								Key:    attributeCode,
								Values: values,
							}
						}
					}
				}
			}
		}
	}

	s.logger.Debug().
		Int("categories", len(facets.Categories)).
		Int("price_ranges", len(facets.PriceRanges)).
		Int("attributes", len(facets.Attributes)).
		Int("source_types", len(facets.SourceTypes)).
		Int("stock_statuses", len(facets.StockStatuses)).
		Msg("aggregations parsed successfully")

	return facets, nil
}

// parseSuggestions extracts suggestions from completion suggester response
func (s *Service) parseSuggestions(result *opensearch.SearchResponse) []Suggestion {
	suggestions := []Suggestion{}

	if result.Suggest == nil || len(result.Suggest) == 0 {
		s.logger.Debug().Msg("no suggestions in response")
		return suggestions
	}

	// The suggest field structure from BuildSuggestionsQuery is:
	// { "suggest": { "listing-suggest": [...] } }
	if listingSuggest, ok := result.Suggest["listing-suggest"].([]interface{}); ok {
		for _, suggestionGroup := range listingSuggest {
			if group, ok := suggestionGroup.(map[string]interface{}); ok {
				// Each suggestion group has an "options" array
				if options, ok := group["options"].([]interface{}); ok {
					for _, option := range options {
						if opt, ok := option.(map[string]interface{}); ok {
							text, _ := opt["text"].(string)
							score, _ := opt["_score"].(float64)

							if text == "" {
								continue
							}

							suggestion := Suggestion{
								Text:  text,
								Score: score,
							}

							// Extract listing_id from _source if available
							if source, ok := opt["_source"].(map[string]interface{}); ok {
								if id, ok := source["id"].(float64); ok {
									listingID := int64(id)
									suggestion.ListingID = &listingID
								}
							}

							suggestions = append(suggestions, suggestion)
						}
					}
				}
			}
		}
	}

	s.logger.Debug().
		Int("count", len(suggestions)).
		Msg("suggestions parsed successfully")

	return suggestions
}

// getListingByID fetches a listing by ID to extract category
func (s *Service) getListingByID(ctx context.Context, listingID int64) (*ListingSearchResult, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": listingID,
			},
		},
		"size": 1,
	}

	result, err := s.searchClient.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	listings := s.parseSearchResults(result)
	if len(listings) == 0 {
		return nil, fmt.Errorf("listing not found: %d", listingID)
	}

	return &listings[0], nil
}

// getMockPopularSearches returns mock trending data
// TODO: Replace with real analytics query
func (s *Service) getMockPopularSearches(categoryID *int64, limit int32) []PopularSearch {
	// Mock data - will be replaced with real analytics
	allSearches := map[int64][]PopularSearch{
		1301: { // Cars category
			{Query: "lamborghini", SearchCount: 542, TrendScore: 15.3},
			{Query: "mercedes", SearchCount: 387, TrendScore: 8.2},
			{Query: "bmw", SearchCount: 312, TrendScore: -3.1},
			{Query: "audi", SearchCount: 298, TrendScore: 12.5},
			{Query: "tesla", SearchCount: 276, TrendScore: 25.8},
		},
		1001: { // Electronics
			{Query: "iphone", SearchCount: 1203, TrendScore: 22.5},
			{Query: "samsung", SearchCount: 891, TrendScore: 5.7},
			{Query: "laptop", SearchCount: 654, TrendScore: -1.2},
			{Query: "airpods", SearchCount: 543, TrendScore: 18.3},
			{Query: "macbook", SearchCount: 432, TrendScore: 7.9},
		},
	}

	var searches []PopularSearch

	if categoryID != nil {
		if categorySearches, ok := allSearches[*categoryID]; ok {
			searches = categorySearches
		}
	}

	// Default: combine top searches from all categories
	if len(searches) == 0 {
		searches = []PopularSearch{
			{Query: "iphone", SearchCount: 1203, TrendScore: 22.5},
			{Query: "tesla", SearchCount: 276, TrendScore: 25.8},
			{Query: "lamborghini", SearchCount: 542, TrendScore: 15.3},
			{Query: "airpods", SearchCount: 543, TrendScore: 18.3},
			{Query: "samsung", SearchCount: 891, TrendScore: 5.7},
		}
	}

	// Apply limit
	if int32(len(searches)) > limit {
		searches = searches[:limit]
	}

	return searches
}

// ============================================================================
// Cache Conversion Methods
// ============================================================================

// convertCachedFacets converts cached facets map to FacetsResponse
func (s *Service) convertCachedFacets(cached map[string]interface{}, isCached bool) *FacetsResponse {
	facets := &FacetsResponse{
		Categories:    []CategoryFacet{},
		PriceRanges:   []PriceRangeFacet{},
		Attributes:    make(map[string]AttributeFacet),
		SourceTypes:   []Facet{},
		StockStatuses: []Facet{},
		Cached:        isCached,
	}

	if tookMs, ok := cached["took_ms"].(float64); ok {
		facets.TookMs = int32(tookMs)
	}

	// Parse categories
	if categories, ok := cached["categories"].([]interface{}); ok {
		for _, cat := range categories {
			if catMap, ok := cat.(map[string]interface{}); ok {
				facet := CategoryFacet{}
				if id, ok := catMap["category_id"].(float64); ok {
					facet.CategoryID = int64(id)
				}
				if count, ok := catMap["count"].(float64); ok {
					facet.Count = int64(count)
				}
				facets.Categories = append(facets.Categories, facet)
			}
		}
	}

	// Parse price ranges
	if priceRanges, ok := cached["price_ranges"].([]interface{}); ok {
		for _, pr := range priceRanges {
			if prMap, ok := pr.(map[string]interface{}); ok {
				facet := PriceRangeFacet{}
				if min, ok := prMap["min"].(float64); ok {
					facet.Min = min
				}
				if max, ok := prMap["max"].(float64); ok {
					facet.Max = max
				}
				if count, ok := prMap["count"].(float64); ok {
					facet.Count = int64(count)
				}
				facets.PriceRanges = append(facets.PriceRanges, facet)
			}
		}
	}

	// Parse attributes
	if attributes, ok := cached["attributes"].(map[string]interface{}); ok {
		for key, attr := range attributes {
			if attrMap, ok := attr.(map[string]interface{}); ok {
				facet := AttributeFacet{Key: key}
				if values, ok := attrMap["values"].([]interface{}); ok {
					for _, val := range values {
						if valMap, ok := val.(map[string]interface{}); ok {
							vc := AttributeValueCount{}
							if value, ok := valMap["value"].(string); ok {
								vc.Value = value
							}
							if count, ok := valMap["count"].(float64); ok {
								vc.Count = int64(count)
							}
							facet.Values = append(facet.Values, vc)
						}
					}
				}
				facets.Attributes[key] = facet
			}
		}
	}

	// Parse source types
	if sourceTypes, ok := cached["source_types"].([]interface{}); ok {
		for _, st := range sourceTypes {
			if stMap, ok := st.(map[string]interface{}); ok {
				facet := Facet{}
				if key, ok := stMap["key"].(string); ok {
					facet.Key = key
				}
				if count, ok := stMap["count"].(float64); ok {
					facet.Count = int64(count)
				}
				facets.SourceTypes = append(facets.SourceTypes, facet)
			}
		}
	}

	// Parse stock statuses
	if stockStatuses, ok := cached["stock_statuses"].([]interface{}); ok {
		for _, ss := range stockStatuses {
			if ssMap, ok := ss.(map[string]interface{}); ok {
				facet := Facet{}
				if key, ok := ssMap["key"].(string); ok {
					facet.Key = key
				}
				if count, ok := ssMap["count"].(float64); ok {
					facet.Count = int64(count)
				}
				facets.StockStatuses = append(facets.StockStatuses, facet)
			}
		}
	}

	return facets
}

// convertFacetsForCache converts FacetsResponse to cacheable map
func (s *Service) convertFacetsForCache(facets *FacetsResponse) map[string]interface{} {
	return map[string]interface{}{
		"categories":     facets.Categories,
		"price_ranges":   facets.PriceRanges,
		"attributes":     facets.Attributes,
		"source_types":   facets.SourceTypes,
		"stock_statuses": facets.StockStatuses,
		"took_ms":        facets.TookMs,
	}
}

// convertCachedFilteredSearch converts cached filtered search map to response
func (s *Service) convertCachedFilteredSearch(cached map[string]interface{}, isCached bool) *SearchFiltersResponse {
	response := &SearchFiltersResponse{
		Listings: []ListingSearchResult{},
		Cached:   isCached,
	}

	if total, ok := cached["total"].(float64); ok {
		response.Total = int64(total)
	}

	if tookMs, ok := cached["took_ms"].(float64); ok {
		response.TookMs = int32(tookMs)
	}

	// Parse listings
	if listings, ok := cached["listings"].([]interface{}); ok {
		for _, listingData := range listings {
			if listingMap, ok := listingData.(map[string]interface{}); ok {
				listing := s.parseListingFromHit(listingMap)
				response.Listings = append(response.Listings, listing)
			}
		}
	}

	// Parse facets if present
	if facetsData, ok := cached["facets"].(map[string]interface{}); ok {
		response.Facets = s.convertCachedFacets(facetsData, isCached)
	}

	return response
}

// convertFilteredSearchForCache converts SearchFiltersResponse to cacheable map
func (s *Service) convertFilteredSearchForCache(response *SearchFiltersResponse) map[string]interface{} {
	result := map[string]interface{}{
		"listings": s.convertListingsForCache(response.Listings),
		"total":    response.Total,
		"took_ms":  response.TookMs,
	}

	if response.Facets != nil {
		result["facets"] = s.convertFacetsForCache(response.Facets)
	}

	return result
}

// convertCachedSuggestions converts cached suggestions map to response
func (s *Service) convertCachedSuggestions(cached map[string]interface{}, isCached bool) *SuggestionsResponse {
	response := &SuggestionsResponse{
		Suggestions: []Suggestion{},
		Cached:      isCached,
	}

	if tookMs, ok := cached["took_ms"].(float64); ok {
		response.TookMs = int32(tookMs)
	}

	if suggestions, ok := cached["suggestions"].([]interface{}); ok {
		for _, suggData := range suggestions {
			if suggMap, ok := suggData.(map[string]interface{}); ok {
				sugg := Suggestion{}
				if text, ok := suggMap["text"].(string); ok {
					sugg.Text = text
				}
				if score, ok := suggMap["score"].(float64); ok {
					sugg.Score = score
				}
				if listingID, ok := suggMap["listing_id"].(float64); ok {
					id := int64(listingID)
					sugg.ListingID = &id
				}
				response.Suggestions = append(response.Suggestions, sugg)
			}
		}
	}

	return response
}

// convertSuggestionsForCache converts SuggestionsResponse to cacheable map
func (s *Service) convertSuggestionsForCache(response *SuggestionsResponse) map[string]interface{} {
	return map[string]interface{}{
		"suggestions": response.Suggestions,
		"took_ms":     response.TookMs,
	}
}

// convertCachedPopularSearches converts cached popular searches map to response
func (s *Service) convertCachedPopularSearches(cached map[string]interface{}) *PopularSearchesResponse {
	response := &PopularSearchesResponse{
		Searches: []PopularSearch{},
	}

	if tookMs, ok := cached["took_ms"].(float64); ok {
		response.TookMs = int32(tookMs)
	}

	if searches, ok := cached["searches"].([]interface{}); ok {
		for _, searchData := range searches {
			if searchMap, ok := searchData.(map[string]interface{}); ok {
				search := PopularSearch{}
				if query, ok := searchMap["query"].(string); ok {
					search.Query = query
				}
				if count, ok := searchMap["search_count"].(float64); ok {
					search.SearchCount = int64(count)
				}
				if trend, ok := searchMap["trend_score"].(float64); ok {
					search.TrendScore = trend
				}
				response.Searches = append(response.Searches, search)
			}
		}
	}

	return response
}

// convertPopularSearchesForCache converts PopularSearchesResponse to cacheable map
func (s *Service) convertPopularSearchesForCache(response *PopularSearchesResponse) map[string]interface{} {
	return map[string]interface{}{
		"searches": response.Searches,
		"took_ms":  response.TookMs,
	}
}
