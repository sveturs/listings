package search

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/opensearch"
)

func TestSearchRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		req       *SearchRequest
		wantLimit int32
		wantErr   error
	}{
		{
			name: "valid request",
			req: &SearchRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
			},
			wantLimit: 20,
			wantErr:   nil,
		},
		{
			name: "limit too small - use default",
			req: &SearchRequest{
				Query:  "laptop",
				Limit:  0,
				Offset: 0,
			},
			wantLimit: 20, // Default
			wantErr:   nil,
		},
		{
			name: "limit too large - cap to max",
			req: &SearchRequest{
				Query:  "laptop",
				Limit:  150,
				Offset: 0,
			},
			wantLimit: 100, // Max
			wantErr:   nil,
		},
		{
			name: "negative offset - reset to 0",
			req: &SearchRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: -10,
			},
			wantLimit: 20,
			wantErr:   nil,
		},
		{
			name: "query too long",
			req: &SearchRequest{
				Query:  string(make([]byte, 501)), // 501 chars
				Limit:  20,
				Offset: 0,
			},
			wantLimit: 20,
			wantErr:   ErrQueryTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLimit, tt.req.Limit)

				// Check offset is never negative after validation
				assert.GreaterOrEqual(t, tt.req.Offset, int32(0))
			}
		})
	}
}

func TestBuildSearchQuery(t *testing.T) {
	// Create a minimal service for testing query builder
	svc := &Service{}

	tests := []struct {
		name     string
		req      *SearchRequest
		validate func(t *testing.T, query map[string]interface{})
	}{
		{
			name: "query with text",
			req: &SearchRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
			},
			validate: func(t *testing.T, query map[string]interface{}) {
				// Check structure
				require.Contains(t, query, "query")
				require.Contains(t, query, "size")
				require.Contains(t, query, "from")
				require.Contains(t, query, "sort")

				// Check size and from
				assert.Equal(t, int32(20), query["size"])
				assert.Equal(t, int32(0), query["from"])

				// Check bool query
				boolQuery := query["query"].(map[string]interface{})
				require.Contains(t, boolQuery, "bool")

				boolPart := boolQuery["bool"].(map[string]interface{})
				require.Contains(t, boolPart, "must")

				mustClauses := boolPart["must"].([]map[string]interface{})
				assert.Len(t, mustClauses, 2) // status + multi_match

				// Find multi_match clause
				var hasMultiMatch bool
				for _, clause := range mustClauses {
					if _, ok := clause["multi_match"]; ok {
						hasMultiMatch = true
						mm := clause["multi_match"].(map[string]interface{})
						assert.Equal(t, "laptop", mm["query"])
						assert.Equal(t, []string{"title^3", "description"}, mm["fields"])
					}
				}
				assert.True(t, hasMultiMatch, "should have multi_match clause")
			},
		},
		{
			name: "query with category filter",
			req: &SearchRequest{
				Query:      "laptop",
				CategoryID: func() *int64 { v := int64(5); return &v }(),
				Limit:      20,
				Offset:     0,
			},
			validate: func(t *testing.T, query map[string]interface{}) {
				boolQuery := query["query"].(map[string]interface{})
				boolPart := boolQuery["bool"].(map[string]interface{})
				mustClauses := boolPart["must"].([]map[string]interface{})

				assert.Len(t, mustClauses, 3) // status + multi_match + category

				// Find category term
				var hasCategoryTerm bool
				for _, clause := range mustClauses {
					if termClause, ok := clause["term"]; ok {
						term := termClause.(map[string]interface{})
						if catID, ok := term["category_id"]; ok {
							assert.Equal(t, int64(5), catID)
							hasCategoryTerm = true
						}
					}
				}
				assert.True(t, hasCategoryTerm, "should have category term")
			},
		},
		{
			name: "query without text - only filters",
			req: &SearchRequest{
				Query:      "",
				CategoryID: func() *int64 { v := int64(5); return &v }(),
				Limit:      50,
				Offset:     10,
			},
			validate: func(t *testing.T, query map[string]interface{}) {
				assert.Equal(t, int32(50), query["size"])
				assert.Equal(t, int32(10), query["from"])

				boolQuery := query["query"].(map[string]interface{})
				boolPart := boolQuery["bool"].(map[string]interface{})
				mustClauses := boolPart["must"].([]map[string]interface{})

				// Should have only status + category (no multi_match for empty query)
				assert.Len(t, mustClauses, 2)

				// No multi_match clause
				for _, clause := range mustClauses {
					_, hasMultiMatch := clause["multi_match"]
					assert.False(t, hasMultiMatch, "should not have multi_match for empty query")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := svc.buildSearchQuery(tt.req)
			tt.validate(t, query)
		})
	}
}

func TestParseListingFromHit(t *testing.T) {
	svc := &Service{}

	source := map[string]interface{}{
		"id":           float64(123),
		"uuid":         "test-uuid-123",
		"title":        "Test Laptop",
		"description":  "Great laptop for work",
		"price":        999.99,
		"currency":     "EUR",
		"category_id":  float64(5),
		"status":       "active",
		"created_at":   "2025-11-17T10:00:00Z",
		"user_id":      float64(456),
		"quantity":     float64(10),
		"source_type":  "b2c",
		"stock_status": "in_stock",
		"images": []interface{}{
			map[string]interface{}{
				"id":            float64(1),
				"url":           "https://example.com/image1.jpg",
				"is_primary":    true,
				"display_order": float64(1),
			},
		},
	}

	listing := svc.parseListingFromHit(source)

	assert.Equal(t, int64(123), listing.ID)
	assert.Equal(t, "test-uuid-123", listing.UUID)
	assert.Equal(t, "Test Laptop", listing.Title)
	assert.NotNil(t, listing.Description)
	assert.Equal(t, "Great laptop for work", *listing.Description)
	assert.Equal(t, 999.99, listing.Price)
	assert.Equal(t, "EUR", listing.Currency)
	assert.Equal(t, int64(5), listing.CategoryID)
	assert.Equal(t, "active", listing.Status)
	assert.Equal(t, "2025-11-17T10:00:00Z", listing.CreatedAt)
	assert.Equal(t, int64(456), listing.UserID)
	assert.Equal(t, int32(10), listing.Quantity)
	assert.Equal(t, "b2c", listing.SourceType)
	assert.Equal(t, "in_stock", listing.StockStatus)

	require.Len(t, listing.Images, 1)
	assert.Equal(t, int64(1), listing.Images[0].ID)
	assert.Equal(t, "https://example.com/image1.jpg", listing.Images[0].URL)
	assert.True(t, listing.Images[0].IsPrimary)
	assert.Equal(t, int32(1), listing.Images[0].DisplayOrder)
}

func TestParseListingFromHit_OptionalFields(t *testing.T) {
	svc := &Service{}

	source := map[string]interface{}{
		"id":           float64(123),
		"uuid":         "test-uuid-123",
		"title":        "Test Laptop",
		"price":        999.99,
		"currency":     "EUR",
		"category_id":  float64(5),
		"status":       "active",
		"created_at":   "2025-11-17T10:00:00Z",
		"user_id":      float64(456),
		"quantity":     float64(10),
		"source_type":  "c2c",
		"stock_status": "in_stock",
		// No description, storefront_id, sku, images
	}

	listing := svc.parseListingFromHit(source)

	assert.Nil(t, listing.Description)
	assert.Nil(t, listing.StorefrontID)
	assert.Nil(t, listing.SKU)
	assert.Empty(t, listing.Images)
}

// ============================================================================
// PHASE 21.2: Tests for Advanced Search Methods
// ============================================================================

// Note: GetPopularSearches now uses real PostgreSQL queries via repository.
// Tests moved to integration tests or handlers_search_test.go with mocked repository.

func TestConvertFacetsForCache_RoundTrip(t *testing.T) {
	svc := &Service{}

	// Create original facets
	original := &FacetsResponse{
		Categories: []CategoryFacet{
			{CategoryID: 1001, Count: 150},
			{CategoryID: 1002, Count: 75},
		},
		PriceRanges: []PriceRangeFacet{
			{Min: 0, Max: 100, Count: 50},
			{Min: 100, Max: 200, Count: 30},
		},
		Attributes: map[string]AttributeFacet{
			"color": {
				Key: "color",
				Values: []AttributeValueCount{
					{Value: "red", Count: 20},
					{Value: "blue", Count: 15},
				},
			},
		},
		SourceTypes: []Facet{
			{Key: "b2c", Count: 100},
			{Key: "c2c", Count: 50},
		},
		StockStatuses: []Facet{
			{Key: "in_stock", Count: 120},
			{Key: "out_of_stock", Count: 30},
		},
		TookMs: 42,
		Cached: false,
	}

	// Convert to cache format
	cached := svc.convertFacetsForCache(original)

	// Simulate JSON marshaling/unmarshaling (what Redis does)
	jsonData, err := json.Marshal(cached)
	require.NoError(t, err)

	var cachedAfterJSON map[string]interface{}
	err = json.Unmarshal(jsonData, &cachedAfterJSON)
	require.NoError(t, err)

	// Convert back
	restored := svc.convertCachedFacets(cachedAfterJSON, true)

	// Verify round-trip
	assert.Equal(t, len(original.Categories), len(restored.Categories))
	assert.Equal(t, len(original.PriceRanges), len(restored.PriceRanges))
	assert.Equal(t, len(original.Attributes), len(restored.Attributes))
	assert.Equal(t, len(original.SourceTypes), len(restored.SourceTypes))
	assert.Equal(t, len(original.StockStatuses), len(restored.StockStatuses))
	assert.Equal(t, original.TookMs, restored.TookMs)
	assert.True(t, restored.Cached)

	// Verify categories
	for i, cat := range original.Categories {
		assert.Equal(t, cat.CategoryID, restored.Categories[i].CategoryID)
		assert.Equal(t, cat.Count, restored.Categories[i].Count)
	}

	// Verify attributes
	for key, attr := range original.Attributes {
		restoredAttr, ok := restored.Attributes[key]
		assert.True(t, ok)
		assert.Equal(t, attr.Key, restoredAttr.Key)
		assert.Equal(t, len(attr.Values), len(restoredAttr.Values))
	}
}

func TestConvertSuggestionsForCache_RoundTrip(t *testing.T) {
	svc := &Service{}

	listingID := int64(123)
	original := &SuggestionsResponse{
		Suggestions: []Suggestion{
			{Text: "laptop", Score: 10.5, ListingID: &listingID},
			{Text: "laptop pro", Score: 8.3, ListingID: nil},
		},
		TookMs: 15,
		Cached: false,
	}

	// Convert to cache format
	cached := svc.convertSuggestionsForCache(original)

	// Simulate JSON marshaling/unmarshaling (what Redis does)
	jsonData, err := json.Marshal(cached)
	require.NoError(t, err)

	var cachedAfterJSON map[string]interface{}
	err = json.Unmarshal(jsonData, &cachedAfterJSON)
	require.NoError(t, err)

	// Convert back
	restored := svc.convertCachedSuggestions(cachedAfterJSON, true)

	// Verify round-trip
	assert.Equal(t, len(original.Suggestions), len(restored.Suggestions))
	assert.Equal(t, original.TookMs, restored.TookMs)
	assert.True(t, restored.Cached)

	// Verify suggestions
	for i, sugg := range original.Suggestions {
		assert.Equal(t, sugg.Text, restored.Suggestions[i].Text)
		assert.Equal(t, sugg.Score, restored.Suggestions[i].Score)
		if sugg.ListingID != nil {
			require.NotNil(t, restored.Suggestions[i].ListingID)
			assert.Equal(t, *sugg.ListingID, *restored.Suggestions[i].ListingID)
		} else {
			assert.Nil(t, restored.Suggestions[i].ListingID)
		}
	}
}

func TestConvertPopularSearchesForCache_RoundTrip(t *testing.T) {
	svc := &Service{}

	original := &PopularSearchesResponse{
		Searches: []PopularSearch{
			{Query: "iphone", SearchCount: 1203, TrendScore: 22.5},
			{Query: "laptop", SearchCount: 654, TrendScore: -1.2},
		},
		TookMs: 8,
	}

	// Convert to cache format
	cached := svc.convertPopularSearchesForCache(original)

	// Simulate JSON marshaling/unmarshaling (what Redis does)
	jsonData, err := json.Marshal(cached)
	require.NoError(t, err)

	var cachedAfterJSON map[string]interface{}
	err = json.Unmarshal(jsonData, &cachedAfterJSON)
	require.NoError(t, err)

	// Convert back
	restored := svc.convertCachedPopularSearches(cachedAfterJSON)

	// Verify round-trip
	assert.Equal(t, len(original.Searches), len(restored.Searches))
	assert.Equal(t, original.TookMs, restored.TookMs)

	// Verify searches
	for i, search := range original.Searches {
		assert.Equal(t, search.Query, restored.Searches[i].Query)
		assert.Equal(t, search.SearchCount, restored.Searches[i].SearchCount)
		assert.Equal(t, search.TrendScore, restored.Searches[i].TrendScore)
	}
}

func TestConvertFilteredSearchForCache_RoundTrip(t *testing.T) {
	svc := &Service{}

	desc := "Test description"
	original := &SearchFiltersResponse{
		Listings: []ListingSearchResult{
			{
				ID:          123,
				UUID:        "test-uuid",
				Title:       "Test Laptop",
				Description: &desc,
				Price:       999.99,
				Currency:    "EUR",
				CategoryID:  1001,
			},
		},
		Total:  1,
		TookMs: 25,
		Cached: false,
		Facets: &FacetsResponse{
			Categories: []CategoryFacet{
				{CategoryID: 1001, Count: 1},
			},
			PriceRanges:   []PriceRangeFacet{},
			Attributes:    make(map[string]AttributeFacet),
			SourceTypes:   []Facet{},
			StockStatuses: []Facet{},
			TookMs:        25,
			Cached:        false,
		},
	}

	// Convert to cache format
	cached := svc.convertFilteredSearchForCache(original)

	// Simulate JSON marshaling/unmarshaling (what Redis does)
	jsonData, err := json.Marshal(cached)
	require.NoError(t, err)

	var cachedAfterJSON map[string]interface{}
	err = json.Unmarshal(jsonData, &cachedAfterJSON)
	require.NoError(t, err)

	// Convert back
	restored := svc.convertCachedFilteredSearch(cachedAfterJSON, true)

	// Verify round-trip
	assert.Equal(t, original.Total, restored.Total)
	assert.Equal(t, original.TookMs, restored.TookMs)
	assert.True(t, restored.Cached)
	assert.Equal(t, len(original.Listings), len(restored.Listings))

	// Verify listing
	if len(restored.Listings) > 0 {
		assert.Equal(t, original.Listings[0].ID, restored.Listings[0].ID)
		assert.Equal(t, original.Listings[0].UUID, restored.Listings[0].UUID)
		assert.Equal(t, original.Listings[0].Title, restored.Listings[0].Title)
	}

	// Verify facets
	assert.NotNil(t, restored.Facets)
	assert.Equal(t, len(original.Facets.Categories), len(restored.Facets.Categories))
}

func TestParseAggregations_EmptyResult(t *testing.T) {
	svc := &Service{}

	// Test with empty result (current implementation returns empty facets)
	result := &opensearch.SearchResponse{}
	facets, err := svc.parseAggregations(result)

	assert.NoError(t, err)
	assert.NotNil(t, facets)
	assert.Empty(t, facets.Categories)
	assert.Empty(t, facets.PriceRanges)
	assert.Empty(t, facets.Attributes)
	assert.Empty(t, facets.SourceTypes)
	assert.Empty(t, facets.StockStatuses)
}

func TestParseSuggestions_EmptyResult(t *testing.T) {
	svc := &Service{}

	// Test with empty result (current implementation returns empty suggestions)
	result := &opensearch.SearchResponse{}
	suggestions := svc.parseSuggestions(result)

	assert.NotNil(t, suggestions)
	assert.Empty(t, suggestions)
}

// ============================================================================
// Integration-style Tests (with mock dependencies)
// ============================================================================

// ============================================================================
// Validation Tests for Phase 21.2 Request Types
// ============================================================================

func TestFacetsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *FacetsRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: &FacetsRequest{
				Query:      "laptop",
				CategoryID: func() *int64 { v := int64(1001); return &v }(),
				UseCache:   true,
			},
			wantErr: nil,
		},
		{
			name: "query too long",
			req: &FacetsRequest{
				Query: string(make([]byte, 501)),
			},
			wantErr: ErrQueryTooLong,
		},
		{
			name: "valid with filters",
			req: &FacetsRequest{
				Query: "laptop",
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: func() *float64 { v := 100.0; return &v }(),
						Max: func() *float64 { v := 500.0; return &v }(),
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearchFiltersRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		req       *SearchFiltersRequest
		wantLimit int32
		wantErr   error
	}{
		{
			name: "valid request",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
			},
			wantLimit: 20,
			wantErr:   nil,
		},
		{
			name: "limit too small - use default",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  0,
				Offset: 0,
			},
			wantLimit: 20,
			wantErr:   nil,
		},
		{
			name: "limit too large - cap to max",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  150,
				Offset: 0,
			},
			wantLimit: 100,
			wantErr:   nil,
		},
		{
			name: "valid with filters and sort",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: func() *float64 { v := 100.0; return &v }(),
						Max: func() *float64 { v := 500.0; return &v }(),
					},
				},
				Sort: &SortConfig{
					Field: "price",
					Order: "asc",
				},
			},
			wantLimit: 20,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLimit, tt.req.Limit)
			}
		})
	}
}

func TestSuggestionsRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		req       *SuggestionsRequest
		wantLimit int32
		wantErr   error
	}{
		{
			name: "valid request",
			req: &SuggestionsRequest{
				Prefix: "lap",
				Limit:  10,
			},
			wantLimit: 10,
			wantErr:   nil,
		},
		{
			name: "prefix too short",
			req: &SuggestionsRequest{
				Prefix: "l",
				Limit:  10,
			},
			wantErr: ErrPrefixTooShort,
		},
		{
			name: "prefix too long",
			req: &SuggestionsRequest{
				Prefix: string(make([]byte, 101)),
				Limit:  10,
			},
			wantErr: ErrPrefixTooLong,
		},
		{
			name: "limit too large - cap to max",
			req: &SuggestionsRequest{
				Prefix: "laptop",
				Limit:  50,
			},
			wantLimit: 20,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLimit, tt.req.Limit)
			}
		})
	}
}

func TestPopularSearchesRequest_Validate(t *testing.T) {
	tests := []struct {
		name          string
		req           *PopularSearchesRequest
		wantLimit     int32
		wantTimeRange string
		wantErr       error
	}{
		{
			name: "valid request",
			req: &PopularSearchesRequest{
				Limit:     10,
				TimeRange: "24h",
			},
			wantLimit:     10,
			wantTimeRange: "24h",
			wantErr:       nil,
		},
		{
			name: "default time range",
			req: &PopularSearchesRequest{
				Limit:     10,
				TimeRange: "",
			},
			wantLimit:     10,
			wantTimeRange: "24h",
			wantErr:       nil,
		},
		{
			name: "invalid time range",
			req: &PopularSearchesRequest{
				Limit:     10,
				TimeRange: "48h",
			},
			wantErr: ErrInvalidTimeRange,
		},
		{
			name: "limit too large - cap to max",
			req: &PopularSearchesRequest{
				Limit:     50,
				TimeRange: "7d",
			},
			wantLimit:     20,
			wantTimeRange: "7d",
			wantErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLimit, tt.req.Limit)
				assert.Equal(t, tt.wantTimeRange, tt.req.TimeRange)
			}
		})
	}
}

// ============================================================================
// Parser Tests - parseAggregations()
// ============================================================================

func TestParseAggregations_FullResponse(t *testing.T) {
	svc := &Service{}

	// Create a realistic OpenSearch aggregations response
	result := &opensearch.SearchResponse{
		Aggregations: map[string]interface{}{
			"categories": map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{
						"key":       float64(1001),
						"doc_count": float64(125),
					},
					map[string]interface{}{
						"key":       float64(1002),
						"doc_count": float64(87),
					},
					map[string]interface{}{
						"key":       float64(1003),
						"doc_count": float64(43),
					},
				},
			},
			"price_ranges": map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{
						"key":       float64(0),
						"doc_count": float64(50),
					},
					map[string]interface{}{
						"key":       float64(100),
						"doc_count": float64(75),
					},
					map[string]interface{}{
						"key":       float64(200),
						"doc_count": float64(0), // Should be filtered out
					},
					map[string]interface{}{
						"key":       float64(300),
						"doc_count": float64(25),
					},
				},
			},
			"source_types": map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{
						"key":       "c2c",
						"doc_count": float64(150),
					},
					map[string]interface{}{
						"key":       "b2c",
						"doc_count": float64(105),
					},
				},
			},
			"stock_statuses": map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{
						"key":       "in_stock",
						"doc_count": float64(200),
					},
					map[string]interface{}{
						"key":       "out_of_stock",
						"doc_count": float64(55),
					},
				},
			},
			"attributes": map[string]interface{}{
				"attribute_keys": map[string]interface{}{
					"buckets": []interface{}{
						map[string]interface{}{
							"key": "brand",
							"attribute_values": map[string]interface{}{
								"buckets": []interface{}{
									map[string]interface{}{
										"key":       "Apple",
										"doc_count": float64(80),
									},
									map[string]interface{}{
										"key":       "Samsung",
										"doc_count": float64(60),
									},
								},
							},
						},
						map[string]interface{}{
							"key": "color",
							"attribute_values": map[string]interface{}{
								"buckets": []interface{}{
									map[string]interface{}{
										"key":       "Black",
										"doc_count": float64(45),
									},
									map[string]interface{}{
										"key":       "White",
										"doc_count": float64(35),
									},
									map[string]interface{}{
										"key":       "Silver",
										"doc_count": float64(20),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	facets, err := svc.parseAggregations(result)

	require.NoError(t, err)
	require.NotNil(t, facets)

	// Verify Categories
	assert.Len(t, facets.Categories, 3)
	assert.Equal(t, int64(1001), facets.Categories[0].CategoryID)
	assert.Equal(t, int64(125), facets.Categories[0].Count)
	assert.Equal(t, int64(1002), facets.Categories[1].CategoryID)
	assert.Equal(t, int64(87), facets.Categories[1].Count)

	// Verify PriceRanges (zero doc_count filtered out)
	assert.Len(t, facets.PriceRanges, 3)
	assert.Equal(t, float64(0), facets.PriceRanges[0].Min)
	assert.Equal(t, float64(100), facets.PriceRanges[0].Max)
	assert.Equal(t, int64(50), facets.PriceRanges[0].Count)

	assert.Equal(t, float64(100), facets.PriceRanges[1].Min)
	assert.Equal(t, float64(200), facets.PriceRanges[1].Max)
	assert.Equal(t, int64(75), facets.PriceRanges[1].Count)

	// Verify SourceTypes
	assert.Len(t, facets.SourceTypes, 2)
	assert.Equal(t, "c2c", facets.SourceTypes[0].Key)
	assert.Equal(t, int64(150), facets.SourceTypes[0].Count)
	assert.Equal(t, "b2c", facets.SourceTypes[1].Key)
	assert.Equal(t, int64(105), facets.SourceTypes[1].Count)

	// Verify StockStatuses
	assert.Len(t, facets.StockStatuses, 2)
	assert.Equal(t, "in_stock", facets.StockStatuses[0].Key)
	assert.Equal(t, int64(200), facets.StockStatuses[0].Count)

	// Verify Attributes
	assert.Len(t, facets.Attributes, 2)

	brandFacet, hasBrand := facets.Attributes["brand"]
	require.True(t, hasBrand)
	assert.Len(t, brandFacet.Values, 2)
	assert.Equal(t, "Apple", brandFacet.Values[0].Value)
	assert.Equal(t, int64(80), brandFacet.Values[0].Count)

	colorFacet, hasColor := facets.Attributes["color"]
	require.True(t, hasColor)
	assert.Len(t, colorFacet.Values, 3)
	assert.Equal(t, "Black", colorFacet.Values[0].Value)
	assert.Equal(t, int64(45), colorFacet.Values[0].Count)
}

func TestParseAggregations_EmptyAggregations(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name   string
		result *opensearch.SearchResponse
	}{
		{
			name: "nil aggregations",
			result: &opensearch.SearchResponse{
				Aggregations: nil,
			},
		},
		{
			name: "empty aggregations map",
			result: &opensearch.SearchResponse{
				Aggregations: map[string]interface{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facets, err := svc.parseAggregations(tt.result)

			require.NoError(t, err)
			require.NotNil(t, facets)

			// All facet types should be empty but initialized
			assert.Empty(t, facets.Categories)
			assert.Empty(t, facets.PriceRanges)
			assert.Empty(t, facets.SourceTypes)
			assert.Empty(t, facets.StockStatuses)
			assert.NotNil(t, facets.Attributes)
			assert.Empty(t, facets.Attributes)
		})
	}
}

func TestParseAggregations_InvalidTypeAssertions(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name   string
		result *opensearch.SearchResponse
	}{
		{
			name: "categories - invalid bucket structure",
			result: &opensearch.SearchResponse{
				Aggregations: map[string]interface{}{
					"categories": map[string]interface{}{
						"buckets": []interface{}{
							map[string]interface{}{
								"key":       "invalid_string", // Should be float64
								"doc_count": float64(10),
							},
						},
					},
				},
			},
		},
		{
			name: "price_ranges - invalid key type",
			result: &opensearch.SearchResponse{
				Aggregations: map[string]interface{}{
					"price_ranges": map[string]interface{}{
						"buckets": []interface{}{
							map[string]interface{}{
								"key":       "invalid", // Should be float64
								"doc_count": float64(10),
							},
						},
					},
				},
			},
		},
		{
			name: "source_types - invalid key type",
			result: &opensearch.SearchResponse{
				Aggregations: map[string]interface{}{
					"source_types": map[string]interface{}{
						"buckets": []interface{}{
							map[string]interface{}{
								"key":       123, // Should be string
								"doc_count": float64(10),
							},
						},
					},
				},
			},
		},
		{
			name: "attributes - missing attribute_values",
			result: &opensearch.SearchResponse{
				Aggregations: map[string]interface{}{
					"attributes": map[string]interface{}{
						"attribute_keys": map[string]interface{}{
							"buckets": []interface{}{
								map[string]interface{}{
									"key": "brand",
									// Missing attribute_values
								},
							},
						},
					},
				},
			},
		},
		{
			name: "buckets - not an array",
			result: &opensearch.SearchResponse{
				Aggregations: map[string]interface{}{
					"categories": map[string]interface{}{
						"buckets": "invalid", // Should be []interface{}
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic, should gracefully skip invalid entries
			facets, err := svc.parseAggregations(tt.result)

			require.NoError(t, err)
			require.NotNil(t, facets)

			// Should return empty facets for invalid data
			assert.NotNil(t, facets.Categories)
			assert.NotNil(t, facets.PriceRanges)
			assert.NotNil(t, facets.SourceTypes)
			assert.NotNil(t, facets.StockStatuses)
			assert.NotNil(t, facets.Attributes)
		})
	}
}

func TestParseAggregations_PartialData(t *testing.T) {
	svc := &Service{}

	// Only some aggregations present
	result := &opensearch.SearchResponse{
		Aggregations: map[string]interface{}{
			"categories": map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{
						"key":       float64(1001),
						"doc_count": float64(50),
					},
				},
			},
			// Other aggregations missing
		},
	}

	facets, err := svc.parseAggregations(result)

	require.NoError(t, err)
	require.NotNil(t, facets)

	// Categories should be populated
	assert.Len(t, facets.Categories, 1)
	assert.Equal(t, int64(1001), facets.Categories[0].CategoryID)

	// Other facets should be empty but initialized
	assert.Empty(t, facets.PriceRanges)
	assert.Empty(t, facets.SourceTypes)
	assert.Empty(t, facets.StockStatuses)
	assert.NotNil(t, facets.Attributes)
	assert.Empty(t, facets.Attributes)
}

// ============================================================================
// Parser Tests - parseSuggestions()
// ============================================================================

func TestParseSuggestions_FullResponse(t *testing.T) {
	svc := &Service{}

	// Create a realistic completion suggester response
	result := &opensearch.SearchResponse{
		Suggest: map[string]interface{}{
			"listing-suggest": []interface{}{
				map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{
							"text":   "iPhone 14 Pro",
							"_score": float64(95.5),
							"_source": map[string]interface{}{
								"id": float64(12345),
							},
						},
						map[string]interface{}{
							"text":   "iPhone 13",
							"_score": float64(88.3),
							"_source": map[string]interface{}{
								"id": float64(12346),
							},
						},
						map[string]interface{}{
							"text":   "iPhone 12 Mini",
							"_score": float64(75.0),
							"_source": map[string]interface{}{
								"id": float64(12347),
							},
						},
					},
				},
			},
		},
	}

	suggestions := svc.parseSuggestions(result)

	require.Len(t, suggestions, 3)

	// Verify first suggestion
	assert.Equal(t, "iPhone 14 Pro", suggestions[0].Text)
	assert.Equal(t, float64(95.5), suggestions[0].Score)
	require.NotNil(t, suggestions[0].ListingID)
	assert.Equal(t, int64(12345), *suggestions[0].ListingID)

	// Verify second suggestion
	assert.Equal(t, "iPhone 13", suggestions[1].Text)
	assert.Equal(t, float64(88.3), suggestions[1].Score)
	require.NotNil(t, suggestions[1].ListingID)
	assert.Equal(t, int64(12346), *suggestions[1].ListingID)

	// Verify third suggestion
	assert.Equal(t, "iPhone 12 Mini", suggestions[2].Text)
	assert.Equal(t, float64(75.0), suggestions[2].Score)
	require.NotNil(t, suggestions[2].ListingID)
	assert.Equal(t, int64(12347), *suggestions[2].ListingID)
}

func TestParseSuggestions_EmptySuggestions(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name   string
		result *opensearch.SearchResponse
	}{
		{
			name: "nil suggest",
			result: &opensearch.SearchResponse{
				Suggest: nil,
			},
		},
		{
			name: "empty suggest map",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{},
			},
		},
		{
			name: "empty listing-suggest array",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": []interface{}{},
				},
			},
		},
		{
			name: "empty options array",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": []interface{}{
						map[string]interface{}{
							"options": []interface{}{},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestions := svc.parseSuggestions(tt.result)

			require.NotNil(t, suggestions)
			assert.Empty(t, suggestions)
		})
	}
}

func TestParseSuggestions_MissingListingID(t *testing.T) {
	svc := &Service{}

	// Suggestion without _source or id field
	result := &opensearch.SearchResponse{
		Suggest: map[string]interface{}{
			"listing-suggest": []interface{}{
				map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{
							"text":   "Laptop",
							"_score": float64(85.0),
							// No _source field
						},
						map[string]interface{}{
							"text":   "Tablet",
							"_score": float64(80.0),
							"_source": map[string]interface{}{
								// No id field in _source
								"title": "Some Tablet",
							},
						},
					},
				},
			},
		},
	}

	suggestions := svc.parseSuggestions(result)

	require.Len(t, suggestions, 2)

	// First suggestion - no _source at all
	assert.Equal(t, "Laptop", suggestions[0].Text)
	assert.Equal(t, float64(85.0), suggestions[0].Score)
	assert.Nil(t, suggestions[0].ListingID)

	// Second suggestion - _source without id
	assert.Equal(t, "Tablet", suggestions[1].Text)
	assert.Equal(t, float64(80.0), suggestions[1].Score)
	assert.Nil(t, suggestions[1].ListingID)
}

func TestParseSuggestions_InvalidTypeAssertions(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name   string
		result *opensearch.SearchResponse
	}{
		{
			name: "text is not string",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": []interface{}{
						map[string]interface{}{
							"options": []interface{}{
								map[string]interface{}{
									"text":   123, // Invalid type
									"_score": float64(80.0),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "score is not float64",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": []interface{}{
						map[string]interface{}{
							"options": []interface{}{
								map[string]interface{}{
									"text":   "Laptop",
									"_score": "invalid", // Invalid type
								},
							},
						},
					},
				},
			},
		},
		{
			name: "id is not float64",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": []interface{}{
						map[string]interface{}{
							"options": []interface{}{
								map[string]interface{}{
									"text":   "Laptop",
									"_score": float64(80.0),
									"_source": map[string]interface{}{
										"id": "invalid", // Invalid type
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "listing-suggest not an array",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": "invalid", // Should be []interface{}
				},
			},
		},
		{
			name: "options not an array",
			result: &opensearch.SearchResponse{
				Suggest: map[string]interface{}{
					"listing-suggest": []interface{}{
						map[string]interface{}{
							"options": "invalid", // Should be []interface{}
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic, should gracefully skip invalid entries
			suggestions := svc.parseSuggestions(tt.result)

			require.NotNil(t, suggestions)
			// Invalid entries should be filtered out
		})
	}
}

func TestParseSuggestions_EmptyText(t *testing.T) {
	svc := &Service{}

	// Suggestion with empty text (should be filtered out)
	result := &opensearch.SearchResponse{
		Suggest: map[string]interface{}{
			"listing-suggest": []interface{}{
				map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{
							"text":   "", // Empty text
							"_score": float64(80.0),
						},
						map[string]interface{}{
							"text":   "Valid Suggestion",
							"_score": float64(75.0),
						},
					},
				},
			},
		},
	}

	suggestions := svc.parseSuggestions(result)

	// Empty text should be filtered out
	require.Len(t, suggestions, 1)
	assert.Equal(t, "Valid Suggestion", suggestions[0].Text)
	assert.Equal(t, float64(75.0), suggestions[0].Score)
}

func TestParseSuggestions_MultipleGroups(t *testing.T) {
	svc := &Service{}

	// Multiple suggestion groups (from multiple input texts)
	result := &opensearch.SearchResponse{
		Suggest: map[string]interface{}{
			"listing-suggest": []interface{}{
				map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{
							"text":   "Laptop Pro",
							"_score": float64(90.0),
						},
					},
				},
				map[string]interface{}{
					"options": []interface{}{
						map[string]interface{}{
							"text":   "Laptop Air",
							"_score": float64(85.0),
						},
					},
				},
			},
		},
	}

	suggestions := svc.parseSuggestions(result)

	// Should aggregate all suggestions from all groups
	require.Len(t, suggestions, 2)
	assert.Equal(t, "Laptop Pro", suggestions[0].Text)
	assert.Equal(t, "Laptop Air", suggestions[1].Text)
}
