package search

import (
	"testing"
)

// ============================================================================
// PHASE 21.2: Query Builders Tests
// ============================================================================

// TestBuildFacetsQuery tests facets query building
func TestBuildFacetsQuery(t *testing.T) {
	tests := []struct {
		name        string
		req         *FacetsRequest
		expectSize  int
		expectAggs  bool
		expectQuery bool
	}{
		{
			name: "minimal request - no filters",
			req: &FacetsRequest{
				UseCache: true,
			},
			expectSize:  0,
			expectAggs:  true,
			expectQuery: true, // Always has status filter
		},
		{
			name: "with query filter",
			req: &FacetsRequest{
				Query:    "laptop",
				UseCache: true,
			},
			expectSize:  0,
			expectAggs:  true,
			expectQuery: true,
		},
		{
			name: "with category filter",
			req: &FacetsRequest{
				CategoryID: ptrString("cat-uuid-1001"),
				UseCache:   true,
			},
			expectSize:  0,
			expectAggs:  true,
			expectQuery: true,
		},
		{
			name: "with filters",
			req: &FacetsRequest{
				Query:      "laptop",
				CategoryID: ptrString("cat-uuid-1001"),
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: ptrFloat64(100),
						Max: ptrFloat64(1000),
					},
				},
				UseCache: true,
			},
			expectSize:  0,
			expectAggs:  true,
			expectQuery: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := BuildFacetsQuery(tt.req)

			// Check size (should always be 0 for facets)
			if size, ok := query["size"].(int); ok {
				if size != tt.expectSize {
					t.Errorf("BuildFacetsQuery() size = %v, want %v", size, tt.expectSize)
				}
			}

			// Check aggregations presence
			if _, ok := query["aggs"]; ok != tt.expectAggs {
				t.Errorf("BuildFacetsQuery() has aggs = %v, want %v", ok, tt.expectAggs)
			}

			// Check query presence
			if _, ok := query["query"]; ok != tt.expectQuery {
				t.Errorf("BuildFacetsQuery() has query = %v, want %v", ok, tt.expectQuery)
			}

			// Validate query structure
			if err := ValidateQuery(query); err != nil {
				t.Errorf("BuildFacetsQuery() produced invalid query: %v", err)
			}

			// Pretty print for manual inspection (comment out in CI)
			// prettyPrint(t, query)
		})
	}
}

// TestBuildFilteredSearchQuery tests filtered search query building
func TestBuildFilteredSearchQuery(t *testing.T) {
	tests := []struct {
		name            string
		req             *SearchFiltersRequest
		expectSort      bool
		expectAggs      bool
		expectGeoFilter bool
	}{
		{
			name: "minimal request",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
			},
			expectSort: true,
			expectAggs: false,
		},
		{
			name: "with sort",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Sort: &SortConfig{
					Field: "price",
					Order: "asc",
				},
			},
			expectSort: true,
			expectAggs: false,
		},
		{
			name: "with facets",
			req: &SearchFiltersRequest{
				Query:         "laptop",
				Limit:         20,
				Offset:        0,
				IncludeFacets: true,
			},
			expectSort: true,
			expectAggs: true,
		},
		{
			name: "with all filters",
			req: &SearchFiltersRequest{
				Query:  "laptop",
				Limit:  20,
				Offset: 0,
				Filters: &SearchFilters{
					Price: &PriceRange{
						Min: ptrFloat64(100),
						Max: ptrFloat64(1000),
					},
					Attributes: map[string][]string{
						"brand": {"Nike", "Adidas"},
						"color": {"Red"},
					},
					Location: &LocationFilter{
						Lat:      55.7558,
						Lon:      37.6173,
						RadiusKm: 10,
					},
					SourceType:  ptrString("c2c"),
					StockStatus: ptrString("in_stock"),
				},
			},
			expectSort:      true,
			expectAggs:      false,
			expectGeoFilter: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := BuildFilteredSearchQuery(tt.req)

			// Check sort presence
			if _, ok := query["sort"]; ok != tt.expectSort {
				t.Errorf("BuildFilteredSearchQuery() has sort = %v, want %v", ok, tt.expectSort)
			}

			// Check aggregations presence
			if _, ok := query["aggs"]; ok != tt.expectAggs {
				t.Errorf("BuildFilteredSearchQuery() has aggs = %v, want %v", ok, tt.expectAggs)
			}

			// Check size and from
			if size, ok := query["size"].(int32); ok {
				if size != tt.req.Limit {
					t.Errorf("BuildFilteredSearchQuery() size = %v, want %v", size, tt.req.Limit)
				}
			}

			if from, ok := query["from"].(int32); ok {
				if from != tt.req.Offset {
					t.Errorf("BuildFilteredSearchQuery() from = %v, want %v", from, tt.req.Offset)
				}
			}

			// Validate query structure
			if err := ValidateQuery(query); err != nil {
				t.Errorf("BuildFilteredSearchQuery() produced invalid query: %v", err)
			}

			// Pretty print for manual inspection (comment out in CI)
			// prettyPrint(t, query)
		})
	}
}

// TestBuildSuggestionsQuery tests suggestions query building
func TestBuildSuggestionsQuery(t *testing.T) {
	tests := []struct {
		name           string
		req            *SuggestionsRequest
		expectContexts bool
	}{
		{
			name: "minimal request",
			req: &SuggestionsRequest{
				Prefix:   "lap",
				Limit:    10,
				UseCache: true,
			},
			expectContexts: false,
		},
		{
			name: "with category context",
			req: &SuggestionsRequest{
				Prefix:     "lap",
				CategoryID: ptrString("cat-uuid-1001"),
				Limit:      10,
				UseCache:   true,
			},
			expectContexts: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := BuildSuggestionsQuery(tt.req)

			// Check suggest presence
			if suggest, ok := query["suggest"]; !ok {
				t.Error("BuildSuggestionsQuery() missing suggest field")
			} else {
				// Check listing-suggest
				if suggestMap, ok := suggest.(map[string]interface{}); ok {
					if listingSuggest, ok := suggestMap["listing-suggest"]; !ok {
						t.Error("BuildSuggestionsQuery() missing listing-suggest field")
					} else {
						// Check completion
						if lsMap, ok := listingSuggest.(map[string]interface{}); ok {
							if completion, ok := lsMap["completion"]; !ok {
								t.Error("BuildSuggestionsQuery() missing completion field")
							} else {
								// Check contexts
								if compMap, ok := completion.(map[string]interface{}); ok {
									_, hasContexts := compMap["contexts"]
									if hasContexts != tt.expectContexts {
										t.Errorf("BuildSuggestionsQuery() has contexts = %v, want %v", hasContexts, tt.expectContexts)
									}
								}
							}
						}
					}
				}
			}

			// Validate query structure
			if err := ValidateQuery(query); err != nil {
				t.Errorf("BuildSuggestionsQuery() produced invalid query: %v", err)
			}

			// Pretty print for manual inspection (comment out in CI)
			// prettyPrint(t, query)
		})
	}
}

// TestBuildPopularSearchesQuery tests popular searches query building
func TestBuildPopularSearchesQuery(t *testing.T) {
	tests := []struct {
		name             string
		req              *PopularSearchesRequest
		expectTimeRange  string
		expectCategoryID bool
	}{
		{
			name: "24h time range",
			req: &PopularSearchesRequest{
				Limit:     10,
				TimeRange: "24h",
			},
			expectTimeRange:  "now-24h",
			expectCategoryID: false,
		},
		{
			name: "7d time range",
			req: &PopularSearchesRequest{
				Limit:     10,
				TimeRange: "7d",
			},
			expectTimeRange:  "now-7d",
			expectCategoryID: false,
		},
		{
			name: "30d time range",
			req: &PopularSearchesRequest{
				Limit:     10,
				TimeRange: "30d",
			},
			expectTimeRange:  "now-30d",
			expectCategoryID: false,
		},
		{
			name: "with category filter",
			req: &PopularSearchesRequest{
				CategoryID: ptrString("cat-uuid-1001"),
				Limit:      10,
				TimeRange:  "24h",
			},
			expectTimeRange:  "now-24h",
			expectCategoryID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := BuildPopularSearchesQuery(tt.req)

			// Check size (should be 0 for aggregations only)
			if size, ok := query["size"].(int); ok {
				if size != 0 {
					t.Errorf("BuildPopularSearchesQuery() size = %v, want 0", size)
				}
			}

			// Check aggregations presence
			if _, ok := query["aggs"]; !ok {
				t.Error("BuildPopularSearchesQuery() missing aggs field")
			}

			// Validate query structure
			if err := ValidateQuery(query); err != nil {
				t.Errorf("BuildPopularSearchesQuery() produced invalid query: %v", err)
			}

			// Pretty print for manual inspection (comment out in CI)
			// prettyPrint(t, query)
		})
	}
}

// TestBuildFilterClauses tests filter clause building
func TestBuildFilterClauses(t *testing.T) {
	tests := []struct {
		name          string
		filters       *SearchFilters
		expectClauses int
	}{
		{
			name: "price filter only",
			filters: &SearchFilters{
				Price: &PriceRange{
					Min: ptrFloat64(100),
					Max: ptrFloat64(1000),
				},
			},
			expectClauses: 1,
		},
		{
			name: "multiple attribute filters",
			filters: &SearchFilters{
				Attributes: map[string][]string{
					"brand": {"Nike", "Adidas"},
					"color": {"Red"},
				},
			},
			expectClauses: 2, // One nested query per attribute
		},
		{
			name: "location filter only",
			filters: &SearchFilters{
				Location: &LocationFilter{
					Lat:      55.7558,
					Lon:      37.6173,
					RadiusKm: 10,
				},
			},
			expectClauses: 1,
		},
		{
			name: "source type and stock status",
			filters: &SearchFilters{
				SourceType:  ptrString("c2c"),
				StockStatus: ptrString("in_stock"),
			},
			expectClauses: 2,
		},
		{
			name: "all filters combined",
			filters: &SearchFilters{
				Price: &PriceRange{
					Min: ptrFloat64(100),
					Max: ptrFloat64(1000),
				},
				Attributes: map[string][]string{
					"brand": {"Nike"},
				},
				Location: &LocationFilter{
					Lat:      55.7558,
					Lon:      37.6173,
					RadiusKm: 10,
				},
				SourceType:  ptrString("c2c"),
				StockStatus: ptrString("in_stock"),
			},
			expectClauses: 5, // price + 1 attribute + location + source + stock
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clauses := buildFilterClauses(tt.filters)

			if len(clauses) != tt.expectClauses {
				t.Errorf("buildFilterClauses() returned %d clauses, want %d", len(clauses), tt.expectClauses)
			}

			// Verify each clause has valid structure
			for i, clause := range clauses {
				if len(clause) == 0 {
					t.Errorf("buildFilterClauses() clause %d is empty", i)
				}
			}
		})
	}
}

// TestBuildAggregations tests aggregations building
func TestBuildAggregations(t *testing.T) {
	aggs := buildAggregations()

	expectedAggs := []string{
		"categories",
		"price_ranges",
		"source_types",
		"stock_statuses",
		"attributes",
	}

	for _, aggName := range expectedAggs {
		if _, ok := aggs[aggName]; !ok {
			t.Errorf("buildAggregations() missing aggregation: %s", aggName)
		}
	}

	// Check nested attributes aggregation structure
	if attrs, ok := aggs["attributes"].(map[string]interface{}); ok {
		if _, ok := attrs["nested"]; !ok {
			t.Error("buildAggregations() attributes missing nested field")
		}
		if _, ok := attrs["aggs"]; !ok {
			t.Error("buildAggregations() attributes missing aggs field")
		}
	} else {
		t.Error("buildAggregations() attributes has wrong type")
	}
}

// TestBuildSort tests sort configuration building
func TestBuildSort(t *testing.T) {
	tests := []struct {
		name      string
		sort      *SortConfig
		expectLen int
	}{
		{
			name: "price ascending",
			sort: &SortConfig{
				Field: "price",
				Order: "asc",
			},
			expectLen: 1,
		},
		{
			name: "created_at descending",
			sort: &SortConfig{
				Field: "created_at",
				Order: "desc",
			},
			expectLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortConfig := buildSort(tt.sort)

			if len(sortConfig) != tt.expectLen {
				t.Errorf("buildSort() returned %d sort configs, want %d", len(sortConfig), tt.expectLen)
			}

			// Check structure
			if len(sortConfig) > 0 {
				if _, ok := sortConfig[0][tt.sort.Field]; !ok {
					t.Errorf("buildSort() missing field %s", tt.sort.Field)
				}
			}
		})
	}
}

// TestValidateQuery tests query validation
func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid - with query",
			query: map[string]interface{}{
				"query": map[string]interface{}{
					"match_all": map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid - with suggest",
			query: map[string]interface{}{
				"suggest": map[string]interface{}{
					"my-suggest": map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name: "valid - with aggs",
			query: map[string]interface{}{
				"aggs": map[string]interface{}{
					"my-agg": map[string]interface{}{},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid - nil query",
			query:   nil,
			wantErr: true,
		},
		{
			name:    "invalid - empty query",
			query:   map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper functions are in test_helpers.go
