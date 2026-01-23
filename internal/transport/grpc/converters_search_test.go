package grpc

import (
	"reflect"
	"testing"

	searchv1 "github.com/vondi-global/listings/api/proto/search/v1"
	"github.com/vondi-global/listings/internal/service/search"
)

// ============================================================================
// PHASE 21.2: Converters Tests (Proto ↔ Domain Round-Trip)
// ============================================================================

// TestProtoToFacetsRequest_RoundTrip tests proto to domain conversion for FacetsRequest
func TestProtoToFacetsRequest_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		proto *searchv1.GetSearchFacetsRequest
	}{
		{
			name:  "minimal request",
			proto: &searchv1.GetSearchFacetsRequest{},
		},
		{
			name: "with query",
			proto: &searchv1.GetSearchFacetsRequest{
				Query: ptrString("laptop"),
			},
		},
		{
			name: "with category",
			proto: &searchv1.GetSearchFacetsRequest{
				CategoryId: "cat-1001",
			},
		},
		{
			name: "with filters",
			proto: &searchv1.GetSearchFacetsRequest{
				Query:      ptrString("laptop"),
				CategoryId: "cat-1001",
				Filters: &searchv1.Filters{
					Price: &searchv1.PriceRange{
						Min: ptrFloat64(100),
						Max: ptrFloat64(1000),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ProtoToFacetsRequest(tt.proto)

			// Verify domain has expected values
			if tt.proto.Query != nil && domain.Query != *tt.proto.Query {
				t.Errorf("Query mismatch: got %s, want %s", domain.Query, *tt.proto.Query)
			}

			if tt.proto.CategoryId != "" {
				if domain.CategoryID == nil {
					t.Error("CategoryID is nil, expected value")
				} else if *domain.CategoryID != tt.proto.CategoryId {
					t.Errorf("CategoryID mismatch: got %s, want %s", *domain.CategoryID, tt.proto.CategoryId)
				}
			}

			if tt.proto.Filters != nil && domain.Filters == nil {
				t.Error("Filters is nil, expected value")
			}
		})
	}
}

// TestFacetsResponseToProto_RoundTrip tests domain to proto conversion for FacetsResponse
func TestFacetsResponseToProto_RoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		domain *search.FacetsResponse
	}{
		{
			name: "minimal response",
			domain: &search.FacetsResponse{
				TookMs: 50,
				Cached: false,
			},
		},
		{
			name: "with categories",
			domain: &search.FacetsResponse{
				Categories: []search.CategoryFacet{
					{CategoryID: "cat-1001", Count: 100},
					{CategoryID: "cat-1002", Count: 50},
				},
				TookMs: 50,
				Cached: false,
			},
		},
		{
			name: "with price ranges",
			domain: &search.FacetsResponse{
				PriceRanges: []search.PriceRangeFacet{
					{Min: 0, Max: 100, Count: 20},
					{Min: 100, Max: 200, Count: 30},
				},
				TookMs: 50,
				Cached: false,
			},
		},
		{
			name: "with attributes",
			domain: &search.FacetsResponse{
				Attributes: map[string]search.AttributeFacet{
					"brand": {
						Key: "brand",
						Values: []search.AttributeValueCount{
							{Value: "Nike", Count: 50},
							{Value: "Adidas", Count: 30},
						},
					},
				},
				TookMs: 50,
				Cached: false,
			},
		},
		{
			name: "complete response",
			domain: &search.FacetsResponse{
				Categories: []search.CategoryFacet{
					{CategoryID: "cat-1001", Count: 100},
				},
				PriceRanges: []search.PriceRangeFacet{
					{Min: 0, Max: 100, Count: 20},
				},
				Attributes: map[string]search.AttributeFacet{
					"brand": {
						Key: "brand",
						Values: []search.AttributeValueCount{
							{Value: "Nike", Count: 50},
						},
					},
				},
				SourceTypes: []search.Facet{
					{Key: "c2c", Count: 80},
					{Key: "b2c", Count: 20},
				},
				StockStatuses: []search.Facet{
					{Key: "in_stock", Count: 90},
					{Key: "out_of_stock", Count: 10},
				},
				TookMs: 50,
				Cached: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proto := FacetsResponseToProto(tt.domain)

			// Verify proto has expected values
			if proto.TookMs != tt.domain.TookMs {
				t.Errorf("TookMs mismatch: got %d, want %d", proto.TookMs, tt.domain.TookMs)
			}

			if proto.Cached != tt.domain.Cached {
				t.Errorf("Cached mismatch: got %v, want %v", proto.Cached, tt.domain.Cached)
			}

			if len(proto.Categories) != len(tt.domain.Categories) {
				t.Errorf("Categories length mismatch: got %d, want %d", len(proto.Categories), len(tt.domain.Categories))
			}

			if len(proto.PriceRanges) != len(tt.domain.PriceRanges) {
				t.Errorf("PriceRanges length mismatch: got %d, want %d", len(proto.PriceRanges), len(tt.domain.PriceRanges))
			}

			if len(proto.Attributes) != len(tt.domain.Attributes) {
				t.Errorf("Attributes length mismatch: got %d, want %d", len(proto.Attributes), len(tt.domain.Attributes))
			}
		})
	}
}

// TestProtoToSearchFiltersRequest_RoundTrip tests proto to domain conversion for SearchFiltersRequest
func TestProtoToSearchFiltersRequest_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		proto *searchv1.SearchWithFiltersRequest
	}{
		{
			name: "minimal request",
			proto: &searchv1.SearchWithFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
			},
		},
		{
			name: "with sort",
			proto: &searchv1.SearchWithFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Sort: &searchv1.SortConfig{
					Field: "price",
					Order: "asc",
				},
			},
		},
		{
			name: "with filters",
			proto: &searchv1.SearchWithFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Filters: &searchv1.Filters{
					Price: &searchv1.PriceRange{
						Min: ptrFloat64(100),
						Max: ptrFloat64(1000),
					},
					Attributes: map[string]*searchv1.AttributeValues{
						"brand": {Values: []string{"Nike", "Adidas"}},
					},
					SourceType:  ptrString("c2c"),
					StockStatus: ptrString("in_stock"),
				},
			},
		},
		{
			name: "with facets included",
			proto: &searchv1.SearchWithFiltersRequest{
				Query:         "laptop",
				Limit:         20,
				Offset:        0,
				IncludeFacets: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ProtoToSearchFiltersRequest(tt.proto)

			if domain.Query != tt.proto.Query {
				t.Errorf("Query mismatch: got %s, want %s", domain.Query, tt.proto.Query)
			}

			if domain.Limit != tt.proto.Limit && tt.proto.Limit != 0 {
				t.Errorf("Limit mismatch: got %d, want %d", domain.Limit, tt.proto.Limit)
			}

			if domain.Offset != tt.proto.Offset {
				t.Errorf("Offset mismatch: got %d, want %d", domain.Offset, tt.proto.Offset)
			}

			if domain.IncludeFacets != tt.proto.IncludeFacets {
				t.Errorf("IncludeFacets mismatch: got %v, want %v", domain.IncludeFacets, tt.proto.IncludeFacets)
			}

			if tt.proto.Sort != nil {
				if domain.Sort == nil {
					t.Error("Sort is nil, expected value")
				} else if domain.Sort.Field != tt.proto.Sort.Field {
					t.Errorf("Sort.Field mismatch: got %s, want %s", domain.Sort.Field, tt.proto.Sort.Field)
				}
			}

			if tt.proto.Filters != nil && domain.Filters == nil {
				t.Error("Filters is nil, expected value")
			}
		})
	}
}

// TestProtoToSearchFilters tests Filters conversion
func TestProtoToSearchFilters(t *testing.T) {
	tests := []struct {
		name  string
		proto *searchv1.Filters
	}{
		{
			name: "price range only",
			proto: &searchv1.Filters{
				Price: &searchv1.PriceRange{
					Min: ptrFloat64(100),
					Max: ptrFloat64(1000),
				},
			},
		},
		{
			name: "attributes only",
			proto: &searchv1.Filters{
				Attributes: map[string]*searchv1.AttributeValues{
					"brand": {Values: []string{"Nike", "Adidas"}},
					"color": {Values: []string{"Red"}},
				},
			},
		},
		{
			name: "location only",
			proto: &searchv1.Filters{
				Location: &searchv1.LocationFilter{
					Lat:      55.7558,
					Lon:      37.6173,
					RadiusKm: 10,
				},
			},
		},
		{
			name: "all filters",
			proto: &searchv1.Filters{
				Price: &searchv1.PriceRange{
					Min: ptrFloat64(100),
					Max: ptrFloat64(1000),
				},
				Attributes: map[string]*searchv1.AttributeValues{
					"brand": {Values: []string{"Nike"}},
				},
				Location: &searchv1.LocationFilter{
					Lat:      55.7558,
					Lon:      37.6173,
					RadiusKm: 10,
				},
				SourceType:  ptrString("c2c"),
				StockStatus: ptrString("in_stock"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ProtoToSearchFilters(tt.proto)

			if tt.proto.Price != nil {
				if domain.Price == nil {
					t.Error("Price is nil, expected value")
				} else {
					if tt.proto.Price.Min != nil && *domain.Price.Min != *tt.proto.Price.Min {
						t.Errorf("Price.Min mismatch: got %f, want %f", *domain.Price.Min, *tt.proto.Price.Min)
					}
					if tt.proto.Price.Max != nil && *domain.Price.Max != *tt.proto.Price.Max {
						t.Errorf("Price.Max mismatch: got %f, want %f", *domain.Price.Max, *tt.proto.Price.Max)
					}
				}
			}

			if len(tt.proto.Attributes) != len(domain.Attributes) {
				t.Errorf("Attributes length mismatch: got %d, want %d", len(domain.Attributes), len(tt.proto.Attributes))
			}

			if tt.proto.Location != nil {
				if domain.Location == nil {
					t.Error("Location is nil, expected value")
				} else {
					if domain.Location.Lat != tt.proto.Location.Lat {
						t.Errorf("Location.Lat mismatch: got %f, want %f", domain.Location.Lat, tt.proto.Location.Lat)
					}
				}
			}

			if tt.proto.SourceType != nil {
				if domain.SourceType == nil {
					t.Error("SourceType is nil, expected value")
				} else if *domain.SourceType != *tt.proto.SourceType {
					t.Errorf("SourceType mismatch: got %s, want %s", *domain.SourceType, *tt.proto.SourceType)
				}
			}
		})
	}
}

// TestProtoToSuggestionsRequest_RoundTrip tests suggestions request conversion
func TestProtoToSuggestionsRequest_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		proto *searchv1.GetSuggestionsRequest
	}{
		{
			name: "minimal request",
			proto: &searchv1.GetSuggestionsRequest{
				Prefix: "lap",
				Limit:  10,
			},
		},
		{
			name: "with category",
			proto: &searchv1.GetSuggestionsRequest{
				Prefix:     "lap",
				CategoryId: "cat-1001",
				Limit:      10,
			},
		},
		{
			name: "zero limit (should default)",
			proto: &searchv1.GetSuggestionsRequest{
				Prefix: "lap",
				Limit:  0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ProtoToSuggestionsRequest(tt.proto)

			if domain.Prefix != tt.proto.Prefix {
				t.Errorf("Prefix mismatch: got %s, want %s", domain.Prefix, tt.proto.Prefix)
			}

			if tt.proto.Limit > 0 && domain.Limit != tt.proto.Limit {
				t.Errorf("Limit mismatch: got %d, want %d", domain.Limit, tt.proto.Limit)
			}

			if tt.proto.CategoryId != "" {
				if domain.CategoryID == nil {
					t.Error("CategoryID is nil, expected value")
				} else if *domain.CategoryID != tt.proto.CategoryId {
					t.Errorf("CategoryID mismatch: got %s, want %s", *domain.CategoryID, tt.proto.CategoryId)
				}
			}
		})
	}
}

// TestSuggestionsResponseToProto_RoundTrip tests suggestions response conversion
func TestSuggestionsResponseToProto_RoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		domain *search.SuggestionsResponse
	}{
		{
			name: "empty suggestions",
			domain: &search.SuggestionsResponse{
				Suggestions: []search.Suggestion{},
				TookMs:      10,
				Cached:      false,
			},
		},
		{
			name: "with suggestions",
			domain: &search.SuggestionsResponse{
				Suggestions: []search.Suggestion{
					{Text: "laptop", Score: 10.5},
					{Text: "laptop bag", Score: 8.2, ListingID: ptrInt64(123)},
				},
				TookMs: 10,
				Cached: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proto := SuggestionsResponseToProto(tt.domain)

			if proto.TookMs != tt.domain.TookMs {
				t.Errorf("TookMs mismatch: got %d, want %d", proto.TookMs, tt.domain.TookMs)
			}

			if proto.Cached != tt.domain.Cached {
				t.Errorf("Cached mismatch: got %v, want %v", proto.Cached, tt.domain.Cached)
			}

			if len(proto.Suggestions) != len(tt.domain.Suggestions) {
				t.Errorf("Suggestions length mismatch: got %d, want %d", len(proto.Suggestions), len(tt.domain.Suggestions))
			}

			for i, domainSug := range tt.domain.Suggestions {
				if i < len(proto.Suggestions) {
					if proto.Suggestions[i].Text != domainSug.Text {
						t.Errorf("Suggestion[%d].Text mismatch: got %s, want %s", i, proto.Suggestions[i].Text, domainSug.Text)
					}
				}
			}
		})
	}
}

// TestProtoToPopularSearchesRequest_RoundTrip tests popular searches request conversion
func TestProtoToPopularSearchesRequest_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		proto *searchv1.GetPopularSearchesRequest
	}{
		{
			name: "minimal request",
			proto: &searchv1.GetPopularSearchesRequest{
				Limit: 10,
			},
		},
		{
			name: "with time range",
			proto: &searchv1.GetPopularSearchesRequest{
				Limit:     10,
				TimeRange: ptrString("7d"),
			},
		},
		{
			name: "with category",
			proto: &searchv1.GetPopularSearchesRequest{
				CategoryId: "cat-1001",
				Limit:      10,
				TimeRange:  ptrString("24h"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := ProtoToPopularSearchesRequest(tt.proto)

			if domain.Limit != tt.proto.Limit && tt.proto.Limit != 0 {
				t.Errorf("Limit mismatch: got %d, want %d", domain.Limit, tt.proto.Limit)
			}

			if tt.proto.TimeRange != nil && domain.TimeRange != *tt.proto.TimeRange {
				t.Errorf("TimeRange mismatch: got %s, want %s", domain.TimeRange, *tt.proto.TimeRange)
			}

			if tt.proto.CategoryId != "" {
				if domain.CategoryID == nil {
					t.Error("CategoryID is nil, expected value")
				} else if *domain.CategoryID != tt.proto.CategoryId {
					t.Errorf("CategoryID mismatch: got %s, want %s", *domain.CategoryID, tt.proto.CategoryId)
				}
			}
		})
	}
}

// TestPopularSearchesResponseToProto_RoundTrip tests popular searches response conversion
func TestPopularSearchesResponseToProto_RoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		domain *search.PopularSearchesResponse
	}{
		{
			name: "empty searches",
			domain: &search.PopularSearchesResponse{
				Searches: []search.PopularSearch{},
				TookMs:   10,
			},
		},
		{
			name: "with searches",
			domain: &search.PopularSearchesResponse{
				Searches: []search.PopularSearch{
					{Query: "laptop", SearchCount: 1000, TrendScore: 15.5},
					{Query: "phone", SearchCount: 800, TrendScore: -2.3},
				},
				TookMs: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proto := PopularSearchesResponseToProto(tt.domain)

			if proto.TookMs != tt.domain.TookMs {
				t.Errorf("TookMs mismatch: got %d, want %d", proto.TookMs, tt.domain.TookMs)
			}

			if len(proto.Searches) != len(tt.domain.Searches) {
				t.Errorf("Searches length mismatch: got %d, want %d", len(proto.Searches), len(tt.domain.Searches))
			}

			for i, domainSearch := range tt.domain.Searches {
				if i < len(proto.Searches) {
					if proto.Searches[i].Query != domainSearch.Query {
						t.Errorf("Search[%d].Query mismatch: got %s, want %s", i, proto.Searches[i].Query, domainSearch.Query)
					}
				}
			}
		})
	}
}

// TestListingSearchResultToProto tests listing conversion
func TestListingSearchResultToProto(t *testing.T) {
	listing := &search.ListingSearchResult{
		ID:          123,
		UUID:        "test-uuid",
		Title:       "Test Listing",
		Description: ptrString("Test description"),
		Price:       99.99,
		Currency:    "EUR",
		CategoryID:  "cat-1001",
		Status:      "active",
		Images: []search.ListingImageResult{
			{ID: 1, URL: "http://example.com/img1.jpg", IsPrimary: true, DisplayOrder: 0},
			{ID: 2, URL: "http://example.com/img2.jpg", IsPrimary: false, DisplayOrder: 1},
		},
		CreatedAt:    "2025-11-17T10:00:00Z",
		UserID:       456,
		StorefrontID: ptrInt64(789),
		Quantity:     10,
		SKU:          ptrString("SKU-123"),
		SourceType:   "c2c",
		StockStatus:  "in_stock",
	}

	proto := ListingSearchResultToProto(listing)

	if proto.Id != listing.ID {
		t.Errorf("ID mismatch: got %d, want %d", proto.Id, listing.ID)
	}

	if proto.Title != listing.Title {
		t.Errorf("Title mismatch: got %s, want %s", proto.Title, listing.Title)
	}

	if proto.Description == nil || *proto.Description != *listing.Description {
		t.Error("Description mismatch")
	}

	if proto.Price != listing.Price {
		t.Errorf("Price mismatch: got %f, want %f", proto.Price, listing.Price)
	}

	if len(proto.Images) != len(listing.Images) {
		t.Errorf("Images length mismatch: got %d, want %d", len(proto.Images), len(listing.Images))
	}
}

// TestSearchFiltersResponseToProto_Complete tests complete filtered search response conversion
func TestSearchFiltersResponseToProto_Complete(t *testing.T) {
	domain := &search.SearchFiltersResponse{
		Listings: []search.ListingSearchResult{
			{
				ID:          123,
				UUID:        "test-uuid",
				Title:       "Test Listing",
				Price:       99.99,
				Currency:    "EUR",
				CategoryID:  "cat-1001",
				Status:      "active",
				CreatedAt:   "2025-11-17T10:00:00Z",
				UserID:      456,
				Quantity:    10,
				SourceType:  "c2c",
				StockStatus: "in_stock",
			},
		},
		Total:  100,
		TookMs: 50,
		Cached: true,
		Facets: &search.FacetsResponse{
			Categories: []search.CategoryFacet{
				{CategoryID: "cat-1001", Count: 100},
			},
			TookMs: 50,
			Cached: true,
		},
	}

	proto := SearchFiltersResponseToProto(domain)

	if proto.Total != domain.Total {
		t.Errorf("Total mismatch: got %d, want %d", proto.Total, domain.Total)
	}

	if proto.TookMs != domain.TookMs {
		t.Errorf("TookMs mismatch: got %d, want %d", proto.TookMs, domain.TookMs)
	}

	if proto.Cached != domain.Cached {
		t.Errorf("Cached mismatch: got %v, want %v", proto.Cached, domain.Cached)
	}

	if len(proto.Listings) != len(domain.Listings) {
		t.Errorf("Listings length mismatch: got %d, want %d", len(proto.Listings), len(domain.Listings))
	}

	if proto.Facets == nil {
		t.Error("Facets is nil, expected value")
	}
}

// TestProtoToSortConfig tests sort config conversion
func TestProtoToSortConfig(t *testing.T) {
	tests := []struct {
		name  string
		proto *searchv1.SortConfig
		want  *search.SortConfig
	}{
		{
			name: "price ascending",
			proto: &searchv1.SortConfig{
				Field: "price",
				Order: "asc",
			},
			want: &search.SortConfig{
				Field: "price",
				Order: "asc",
			},
		},
		{
			name:  "empty (should default)",
			proto: &searchv1.SortConfig{},
			want: &search.SortConfig{
				Field: "relevance",
				Order: "desc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProtoToSortConfig(tt.proto)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProtoToSortConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrString(v string) *string {
	return &v
}
