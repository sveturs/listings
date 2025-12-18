package search

import (
	"testing"
)

// ============================================================================
// PHASE 21.2: Types Validation Tests
// ============================================================================

// TestFacetsRequestValidate tests FacetsRequest validation
func TestFacetsRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *FacetsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request - minimal",
			req: &FacetsRequest{
				Query:    "",
				UseCache: true,
			},
			wantErr: false,
		},
		{
			name: "valid request - with query and category",
			req: &FacetsRequest{
				Query:      "laptop",
				CategoryID: ptrString("cat-uuid-1001"),
				UseCache:   true,
			},
			wantErr: false,
		},
		{
			name: "valid request - with filters",
			req: &FacetsRequest{
				Query: "laptop",
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: ptrFloat64(100),
						Max: ptrFloat64(1000),
					},
				},
				UseCache: true,
			},
			wantErr: false,
		},
		{
			name: "invalid - query too long",
			req: &FacetsRequest{
				Query: generateLongString(501),
			},
			wantErr: true,
			errMsg:  "search query too long (max 500 characters)",
		},
		{
			name: "invalid - invalid filters",
			req: &FacetsRequest{
				Query: "laptop",
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: ptrFloat64(1000),
						Max: ptrFloat64(100), // Min > Max
					},
				},
			},
			wantErr: true,
			errMsg:  "invalid price range: min and max must be >= 0, min must be <= max",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FacetsRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("FacetsRequest.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestSearchFiltersRequestValidate tests SearchFiltersRequest validation
func TestSearchFiltersRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *SearchFiltersRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request - minimal",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
			},
			wantErr: false,
		},
		{
			name: "valid request - with filters and sort",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: ptrFloat64(100),
						Max: ptrFloat64(1000),
					},
				},
				Sort: &SortConfig{
					Field: "price",
					Order: "asc",
				},
			},
			wantErr: false,
		},
		{
			name: "auto-fix limit - too small",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  0,
				Offset: 0,
			},
			wantErr: false, // Auto-fixed to 20
		},
		{
			name: "auto-fix limit - too large",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  200,
				Offset: 0,
			},
			wantErr: false, // Auto-capped to 100
		},
		{
			name: "auto-fix offset - negative",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: -10,
			},
			wantErr: false, // Auto-fixed to 0
		},
		{
			name: "invalid - query too long",
			req: &SearchFiltersRequest{
				Query:  generateLongString(501),
				Limit:  20,
				Offset: 0,
			},
			wantErr: true,
			errMsg:  "search query too long (max 500 characters)",
		},
		{
			name: "invalid - invalid sort field",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Sort: &SortConfig{
					Field: "invalid_field",
					Order: "asc",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchFiltersRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("SearchFiltersRequest.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}

			// Check auto-fixes
			if !tt.wantErr {
				if tt.req.Limit < 1 || tt.req.Limit > 100 {
					t.Errorf("SearchFiltersRequest.Validate() limit not auto-fixed: got %d", tt.req.Limit)
				}
				if tt.req.Offset < 0 {
					t.Errorf("SearchFiltersRequest.Validate() offset not auto-fixed: got %d", tt.req.Offset)
				}
			}
		})
	}
}

// TestSearchFiltersValidate tests SearchFilters validation
func TestSearchFiltersValidate(t *testing.T) {
	tests := []struct {
		name    string
		filters *SearchFilters
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid - price range",
			filters: &SearchFilters{
				Price: &PriceRange{
					Min: ptrFloat64(100),
					Max: ptrFloat64(1000),
				},
			},
			wantErr: false,
		},
		{
			name: "valid - location",
			filters: &SearchFilters{
				Location: &LocationFilter{
					Lat:      55.7558,
					Lon:      37.6173,
					RadiusKm: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "valid - source type c2c",
			filters: &SearchFilters{
				SourceType: ptrString("c2c"),
			},
			wantErr: false,
		},
		{
			name: "valid - stock status in_stock",
			filters: &SearchFilters{
				StockStatus: ptrString("in_stock"),
			},
			wantErr: false,
		},
		{
			name: "valid - attributes",
			filters: &SearchFilters{
				Attributes: map[string][]string{
					"brand": {"Nike", "Adidas"},
					"color": {"Red"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid - price range min > max",
			filters: &SearchFilters{
				Price: &PriceRange{
					Min: ptrFloat64(1000),
					Max: ptrFloat64(100),
				},
			},
			wantErr: true,
			errMsg:  "invalid price range: min and max must be >= 0, min must be <= max",
		},
		{
			name: "invalid - negative price",
			filters: &SearchFilters{
				Price: &PriceRange{
					Min: ptrFloat64(-100),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid - invalid source type",
			filters: &SearchFilters{
				SourceType: ptrString("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid source type: must be 'c2c' or 'b2c'",
		},
		{
			name: "invalid - invalid stock status",
			filters: &SearchFilters{
				StockStatus: ptrString("invalid"),
			},
			wantErr: true,
			errMsg:  "invalid stock status: must be 'in_stock', 'out_of_stock', or 'low_stock'",
		},
		{
			name: "invalid - too many attributes",
			filters: &SearchFilters{
				Attributes: map[string][]string{
					"attr1":  {"value"},
					"attr2":  {"value"},
					"attr3":  {"value"},
					"attr4":  {"value"},
					"attr5":  {"value"},
					"attr6":  {"value"},
					"attr7":  {"value"},
					"attr8":  {"value"},
					"attr9":  {"value"},
					"attr10": {"value"},
					"attr11": {"value"}, // 11 > 10 max
				},
			},
			wantErr: true,
			errMsg:  "too many attribute filters: maximum 10 allowed",
		},
		{
			name: "invalid - empty attribute key",
			filters: &SearchFilters{
				Attributes: map[string][]string{
					"": {"value"},
				},
			},
			wantErr: true,
			errMsg:  "attribute key cannot be empty",
		},
		{
			name: "invalid - empty attribute values",
			filters: &SearchFilters{
				Attributes: map[string][]string{
					"brand": {},
				},
			},
			wantErr: true,
			errMsg:  "attribute values list cannot be empty",
		},
		{
			name: "invalid - too many attribute values",
			filters: &SearchFilters{
				Attributes: map[string][]string{
					"brand": generateStringSlice(21), // 21 > 20 max
				},
			},
			wantErr: true,
			errMsg:  "too many values for attribute: maximum 20 allowed per attribute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filters.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchFilters.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("SearchFilters.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestPriceRangeValidate tests PriceRange validation
func TestPriceRangeValidate(t *testing.T) {
	tests := []struct {
		name    string
		price   *PriceRange
		wantErr bool
	}{
		{
			name:    "valid - both bounds",
			price:   &PriceRange{Min: ptrFloat64(100), Max: ptrFloat64(1000)},
			wantErr: false,
		},
		{
			name:    "valid - only min",
			price:   &PriceRange{Min: ptrFloat64(100)},
			wantErr: false,
		},
		{
			name:    "valid - only max",
			price:   &PriceRange{Max: ptrFloat64(1000)},
			wantErr: false,
		},
		{
			name:    "valid - zero min",
			price:   &PriceRange{Min: ptrFloat64(0), Max: ptrFloat64(1000)},
			wantErr: false,
		},
		{
			name:    "invalid - negative min",
			price:   &PriceRange{Min: ptrFloat64(-100)},
			wantErr: true,
		},
		{
			name:    "invalid - negative max",
			price:   &PriceRange{Max: ptrFloat64(-100)},
			wantErr: true,
		},
		{
			name:    "invalid - min > max",
			price:   &PriceRange{Min: ptrFloat64(1000), Max: ptrFloat64(100)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.price.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PriceRange.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestLocationFilterValidate tests LocationFilter validation
func TestLocationFilterValidate(t *testing.T) {
	tests := []struct {
		name     string
		location *LocationFilter
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid - Moscow",
			location: &LocationFilter{Lat: 55.7558, Lon: 37.6173, RadiusKm: 10},
			wantErr:  false,
		},
		{
			name:     "valid - max latitude",
			location: &LocationFilter{Lat: 90, Lon: 0, RadiusKm: 10},
			wantErr:  false,
		},
		{
			name:     "valid - min latitude",
			location: &LocationFilter{Lat: -90, Lon: 0, RadiusKm: 10},
			wantErr:  false,
		},
		{
			name:     "valid - max longitude",
			location: &LocationFilter{Lat: 0, Lon: 180, RadiusKm: 10},
			wantErr:  false,
		},
		{
			name:     "valid - min longitude",
			location: &LocationFilter{Lat: 0, Lon: -180, RadiusKm: 10},
			wantErr:  false,
		},
		{
			name:     "invalid - latitude too high",
			location: &LocationFilter{Lat: 91, Lon: 0, RadiusKm: 10},
			wantErr:  true,
			errMsg:   "invalid latitude: must be between -90 and 90",
		},
		{
			name:     "invalid - latitude too low",
			location: &LocationFilter{Lat: -91, Lon: 0, RadiusKm: 10},
			wantErr:  true,
			errMsg:   "invalid latitude: must be between -90 and 90",
		},
		{
			name:     "invalid - longitude too high",
			location: &LocationFilter{Lat: 0, Lon: 181, RadiusKm: 10},
			wantErr:  true,
			errMsg:   "invalid longitude: must be between -180 and 180",
		},
		{
			name:     "invalid - longitude too low",
			location: &LocationFilter{Lat: 0, Lon: -181, RadiusKm: 10},
			wantErr:  true,
			errMsg:   "invalid longitude: must be between -180 and 180",
		},
		{
			name:     "invalid - radius zero",
			location: &LocationFilter{Lat: 0, Lon: 0, RadiusKm: 0},
			wantErr:  true,
			errMsg:   "invalid radius: must be between 0 and 1000 km",
		},
		{
			name:     "invalid - radius negative",
			location: &LocationFilter{Lat: 0, Lon: 0, RadiusKm: -10},
			wantErr:  true,
			errMsg:   "invalid radius: must be between 0 and 1000 km",
		},
		{
			name:     "invalid - radius too large",
			location: &LocationFilter{Lat: 0, Lon: 0, RadiusKm: 1001},
			wantErr:  true,
			errMsg:   "invalid radius: must be between 0 and 1000 km",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.location.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LocationFilter.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("LocationFilter.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestSortConfigValidate tests SortConfig validation
func TestSortConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		sort    *SortConfig
		wantErr bool
	}{
		{
			name:    "valid - price asc",
			sort:    &SortConfig{Field: "price", Order: "asc"},
			wantErr: false,
		},
		{
			name:    "valid - created_at desc",
			sort:    &SortConfig{Field: "created_at", Order: "desc"},
			wantErr: false,
		},
		{
			name:    "valid - relevance desc",
			sort:    &SortConfig{Field: "relevance", Order: "desc"},
			wantErr: false,
		},
		{
			name:    "invalid - unknown field",
			sort:    &SortConfig{Field: "unknown", Order: "asc"},
			wantErr: true,
		},
		{
			name:    "invalid - unknown order",
			sort:    &SortConfig{Field: "price", Order: "unknown"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sort.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SortConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSuggestionsRequestValidate tests SuggestionsRequest validation
func TestSuggestionsRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *SuggestionsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - minimal",
			req:     &SuggestionsRequest{Prefix: "lap", Limit: 10},
			wantErr: false,
		},
		{
			name:    "valid - with category",
			req:     &SuggestionsRequest{Prefix: "lap", CategoryID: ptrString("cat-uuid-1001"), Limit: 10},
			wantErr: false,
		},
		{
			name:    "auto-fix limit - zero",
			req:     &SuggestionsRequest{Prefix: "lap", Limit: 0},
			wantErr: false,
		},
		{
			name:    "auto-cap limit - too large",
			req:     &SuggestionsRequest{Prefix: "lap", Limit: 50},
			wantErr: false,
		},
		{
			name:    "invalid - prefix too short",
			req:     &SuggestionsRequest{Prefix: "l", Limit: 10},
			wantErr: true,
			errMsg:  "prefix too short: minimum 2 characters required",
		},
		{
			name:    "invalid - prefix too long",
			req:     &SuggestionsRequest{Prefix: generateLongString(101), Limit: 10},
			wantErr: true,
			errMsg:  "prefix too long: maximum 100 characters allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SuggestionsRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("SuggestionsRequest.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}

			// Check auto-fixes
			if !tt.wantErr {
				if tt.req.Limit < 1 || tt.req.Limit > 20 {
					t.Errorf("SuggestionsRequest.Validate() limit not auto-fixed: got %d", tt.req.Limit)
				}
			}
		})
	}
}

// TestPopularSearchesRequestValidate tests PopularSearchesRequest validation
func TestPopularSearchesRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *PopularSearchesRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid - 24h",
			req:     &PopularSearchesRequest{Limit: 10, TimeRange: "24h"},
			wantErr: false,
		},
		{
			name:    "valid - 7d",
			req:     &PopularSearchesRequest{Limit: 10, TimeRange: "7d"},
			wantErr: false,
		},
		{
			name:    "valid - 30d",
			req:     &PopularSearchesRequest{Limit: 10, TimeRange: "30d"},
			wantErr: false,
		},
		{
			name:    "auto-default time range",
			req:     &PopularSearchesRequest{Limit: 10, TimeRange: ""},
			wantErr: false,
		},
		{
			name:    "invalid - unknown time range",
			req:     &PopularSearchesRequest{Limit: 10, TimeRange: "invalid"},
			wantErr: true,
			errMsg:  "invalid time range: must be '24h', '7d', or '30d'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PopularSearchesRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg {
					t.Errorf("PopularSearchesRequest.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}

			// Check auto-defaults
			if !tt.wantErr {
				if tt.req.TimeRange == "" {
					t.Errorf("PopularSearchesRequest.Validate() time range not defaulted")
				}
			}
		})
	}
}

// Helper functions are in test_helpers.go
