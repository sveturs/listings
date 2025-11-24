package search

import (
	"fmt"
	"time"
)

// SearchRequest represents domain search parameters
type SearchRequest struct {
	Query      string // Search query text
	CategoryID *int64 // Optional category filter
	Limit      int32  // Results per page (1-100)
	Offset     int32  // Pagination offset
	UseCache   bool   // Whether to use cache
}

// Validate validates search request parameters
func (r *SearchRequest) Validate() error {
	// Limit must be between 1 and 100
	if r.Limit < 1 {
		r.Limit = 20 // Default
	}
	if r.Limit > 100 {
		r.Limit = 100 // Max
	}

	// Offset must be >= 0
	if r.Offset < 0 {
		r.Offset = 0
	}

	// Query max length check (optional, but good practice)
	if len(r.Query) > 500 {
		return ErrQueryTooLong
	}

	return nil
}

// SearchResponse represents search result
type SearchResponse struct {
	Listings []ListingSearchResult `json:"listings"`
	Total    int64                 `json:"total"`
	TookMs   int32                 `json:"took_ms"`
	Cached   bool                  `json:"cached"`
}

// ListingSearchResult represents a single listing in search results
type ListingSearchResult struct {
	ID           int64                `json:"id"`
	UUID         string               `json:"uuid"`
	Title        string               `json:"title"`
	Description  *string              `json:"description,omitempty"`
	Price        float64              `json:"price"`
	Currency     string               `json:"currency"`
	CategoryID   int64                `json:"category_id"`
	Status       string               `json:"status"`
	Images       []ListingImageResult `json:"images,omitempty"`
	CreatedAt    string               `json:"created_at"`
	UserID       int64                `json:"user_id"`
	StorefrontID *int64               `json:"storefront_id,omitempty"`
	Quantity     int32                `json:"quantity"`
	SKU          *string              `json:"sku,omitempty"`
	SourceType   string               `json:"source_type"`
	StockStatus  string               `json:"stock_status"`
}

// ListingImageResult represents an image in search results
type ListingImageResult struct {
	ID           int64  `json:"id"`
	URL          string `json:"url"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int32  `json:"display_order"`
}

// ============================================================================
// PHASE 21.2: Advanced Search Types
// ============================================================================

// FacetsRequest - request for GetSearchFacets
type FacetsRequest struct {
	Query      string         // Optional pre-filter
	CategoryID *int64         // Optional pre-filter
	Filters    *SearchFilters // Optional pre-filter (price, attributes)
	UseCache   bool           // Whether to use cache
}

// Validate validates facets request parameters
func (r *FacetsRequest) Validate() error {
	// Query max length check
	if len(r.Query) > 500 {
		return ErrQueryTooLong
	}

	// Validate filters if provided
	if r.Filters != nil {
		if err := r.Filters.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// FacetsResponse - aggregation results
type FacetsResponse struct {
	Categories    []CategoryFacet           `json:"categories"`
	PriceRanges   []PriceRangeFacet         `json:"price_ranges"`
	Attributes    map[string]AttributeFacet `json:"attributes"`
	SourceTypes   []Facet                   `json:"source_types"`
	StockStatuses []Facet                   `json:"stock_statuses"`
	TookMs        int32                     `json:"took_ms"`
	Cached        bool                      `json:"cached"`
}

// CategoryFacet represents category distribution
type CategoryFacet struct {
	CategoryID int64 `json:"category_id"`
	Count      int64 `json:"count"`
}

// PriceRangeFacet represents price histogram bucket
type PriceRangeFacet struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

// AttributeFacet contains all values for a specific attribute
type AttributeFacet struct {
	Key    string                `json:"key"`
	Values []AttributeValueCount `json:"values"`
}

// AttributeValueCount represents count for specific attribute value
type AttributeValueCount struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// Facet represents generic key-count pair
type Facet struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// SearchFiltersRequest - enhanced search with filters
type SearchFiltersRequest struct {
	Query         string         // Search query text
	CategoryID    *int64         // Optional category filter
	Limit         int32          // Results per page (1-100)
	Offset        int32          // Pagination offset
	Filters       *SearchFilters // Advanced filters
	Sort          *SortConfig    // Sort configuration
	UseCache      bool           // Whether to use cache
	IncludeFacets bool           // Return facets with results
}

// Validate validates search filters request parameters
func (r *SearchFiltersRequest) Validate() error {
	// Limit validation
	if r.Limit < 1 {
		r.Limit = 20 // Default
	}
	if r.Limit > 100 {
		r.Limit = 100 // Max
	}

	// Offset validation
	if r.Offset < 0 {
		r.Offset = 0
	}

	// Query max length check
	if len(r.Query) > 500 {
		return ErrQueryTooLong
	}

	// Validate filters if provided
	if r.Filters != nil {
		if err := r.Filters.Validate(); err != nil {
			return err
		}
	}

	// Validate sort if provided
	if r.Sort != nil {
		if err := r.Sort.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// SearchFilters - all filter options
type SearchFilters struct {
	Price       *PriceRange         // Price range filter
	Attributes  map[string][]string // Attribute filters (key -> values)
	Location    *LocationFilter     // Geo location filter
	SourceType  *string             // "c2c" | "b2c"
	StockStatus *string             // "in_stock" | "out_of_stock" | "low_stock"
}

// Validate validates search filters
func (f *SearchFilters) Validate() error {
	// Validate price range
	if f.Price != nil {
		if err := f.Price.Validate(); err != nil {
			return err
		}
	}

	// Validate location filter
	if f.Location != nil {
		if err := f.Location.Validate(); err != nil {
			return err
		}
	}

	// Validate source type
	if f.SourceType != nil {
		validSourceTypes := map[string]bool{"c2c": true, "b2c": true}
		if !validSourceTypes[*f.SourceType] {
			return ErrInvalidSourceType
		}
	}

	// Validate stock status
	if f.StockStatus != nil {
		validStockStatuses := map[string]bool{
			"in_stock":     true,
			"out_of_stock": true,
			"low_stock":    true,
		}
		if !validStockStatuses[*f.StockStatus] {
			return ErrInvalidStockStatus
		}
	}

	// Validate attributes (max 10 attribute filters)
	if len(f.Attributes) > 10 {
		return ErrTooManyAttributeFilters
	}

	// Validate each attribute filter (max 20 values per attribute)
	for key, values := range f.Attributes {
		if len(key) == 0 {
			return ErrEmptyAttributeKey
		}
		if len(values) == 0 {
			return ErrEmptyAttributeValues
		}
		if len(values) > 20 {
			return ErrTooManyAttributeValues
		}
	}

	return nil
}

// PriceRange filter
type PriceRange struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

// Validate validates price range
func (p *PriceRange) Validate() error {
	if p.Min != nil && *p.Min < 0 {
		return ErrInvalidPriceRange
	}
	if p.Max != nil && *p.Max < 0 {
		return ErrInvalidPriceRange
	}
	if p.Min != nil && p.Max != nil && *p.Min > *p.Max {
		return ErrInvalidPriceRange
	}
	return nil
}

// LocationFilter for geo distance filtering
type LocationFilter struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	RadiusKm float64 `json:"radius_km"`
}

// Validate validates location filter
func (l *LocationFilter) Validate() error {
	if l.Lat < -90 || l.Lat > 90 {
		return ErrInvalidLatitude
	}
	if l.Lon < -180 || l.Lon > 180 {
		return ErrInvalidLongitude
	}
	if l.RadiusKm <= 0 || l.RadiusKm > 1000 {
		return ErrInvalidRadius
	}
	return nil
}

// SortConfig defines sorting parameters
type SortConfig struct {
	Field string `json:"field"` // "price" | "created_at" | "relevance" | "views_count" | "favorites_count"
	Order string `json:"order"` // "asc" | "desc"
}

// Validate validates sort configuration
func (s *SortConfig) Validate() error {
	// Validate sort field
	validFields := map[string]bool{
		"price":           true,
		"created_at":      true,
		"relevance":       true,
		"views_count":     true,
		"favorites_count": true,
	}
	if !validFields[s.Field] {
		return ErrInvalidSortField
	}

	// Validate sort order
	validOrders := map[string]bool{"asc": true, "desc": true}
	if !validOrders[s.Order] {
		return ErrInvalidSortOrder
	}

	return nil
}

// SearchFiltersResponse - search results + optional facets
type SearchFiltersResponse struct {
	Listings []ListingSearchResult `json:"listings"`
	Total    int64                 `json:"total"`
	TookMs   int32                 `json:"took_ms"`
	Cached   bool                  `json:"cached"`
	Facets   *FacetsResponse       `json:"facets,omitempty"` // If IncludeFacets=true
}

// SuggestionsRequest - autocomplete request
type SuggestionsRequest struct {
	Prefix     string // Search prefix (min 2 chars)
	CategoryID *int64 // Optional category filter
	Limit      int32  // Max suggestions (1-20)
	UseCache   bool   // Whether to use cache
}

// Validate validates suggestions request parameters
func (r *SuggestionsRequest) Validate() error {
	// Prefix min length
	if len(r.Prefix) < 2 {
		return ErrPrefixTooShort
	}

	// Prefix max length
	if len(r.Prefix) > 100 {
		return ErrPrefixTooLong
	}

	// Limit validation
	if r.Limit < 1 {
		r.Limit = 10 // Default
	}
	if r.Limit > 20 {
		r.Limit = 20 // Max
	}

	return nil
}

// SuggestionsResponse - autocomplete results
type SuggestionsResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
	TookMs      int32        `json:"took_ms"`
	Cached      bool         `json:"cached"`
}

// Suggestion represents a single autocomplete suggestion
type Suggestion struct {
	Text      string  `json:"text"`
	Score     float64 `json:"score"`
	ListingID *int64  `json:"listing_id,omitempty"` // Optional
}

// PopularSearchesRequest - trending queries request
type PopularSearchesRequest struct {
	CategoryID *int64 // Optional category filter
	Limit      int32  // Max results (1-20)
	TimeRange  string // "24h" | "7d" | "30d"
}

// Validate validates popular searches request parameters
func (r *PopularSearchesRequest) Validate() error {
	// Limit validation
	if r.Limit < 1 {
		r.Limit = 10 // Default
	}
	if r.Limit > 20 {
		r.Limit = 20 // Max
	}

	// Time range validation
	validTimeRanges := map[string]bool{
		"24h": true,
		"7d":  true,
		"30d": true,
	}
	if r.TimeRange == "" {
		r.TimeRange = "24h" // Default
	}
	if !validTimeRanges[r.TimeRange] {
		return ErrInvalidTimeRange
	}

	return nil
}

// PopularSearchesResponse - trending queries
type PopularSearchesResponse struct {
	Searches []PopularSearch `json:"searches"`
	TookMs   int32           `json:"took_ms"`
}

// PopularSearch represents a trending search query
type PopularSearch struct {
	Query       string  `json:"query"`
	SearchCount int64   `json:"search_count"`
	TrendScore  float64 `json:"trend_score"` // +/- change percentage
}

// ============================================================================
// PHASE 28: Search Analytics - Trending Searches
// ============================================================================

// TrendingSearchesRequest - real trending queries from analytics
type TrendingSearchesRequest struct {
	CategoryID *int64 // Optional category filter
	Limit      int32  // Max results (1-50)
	Days       int32  // Period in days (1-30)
}

// Validate validates trending searches request parameters
func (r *TrendingSearchesRequest) Validate() error {
	// Limit validation
	if r.Limit < 1 {
		r.Limit = 10 // Default
	}
	if r.Limit > 50 {
		r.Limit = 50 // Max
	}

	// Days validation
	if r.Days < 1 {
		r.Days = 7 // Default
	}
	if r.Days > 30 {
		r.Days = 30 // Max
	}

	return nil
}

// TrendingSearchesResponse - trending queries from analytics
type TrendingSearchesResponse struct {
	Searches []TrendingSearchResult `json:"searches"`
}

// TrendingSearchResult represents a single trending search query
type TrendingSearchResult struct {
	QueryText    string    `json:"query_text"`
	SearchCount  int32     `json:"search_count"`
	LastSearched time.Time `json:"last_searched"`
}

// ============================================================================
// PHASE 28: Search Analytics - Personal Search History
// ============================================================================

// SearchHistoryRequest represents a request for user's search history
type SearchHistoryRequest struct {
	UserID    *int64  `json:"user_id,omitempty"`
	SessionID *string `json:"session_id,omitempty"`
	Limit     int32   `json:"limit"`
}

// Validate validates search history request parameters
func (r *SearchHistoryRequest) Validate() error {
	// Must have either user_id or session_id (XOR)
	if (r.UserID == nil && r.SessionID == nil) || (r.UserID != nil && r.SessionID != nil) {
		return fmt.Errorf("exactly one of user_id or session_id must be provided")
	}

	// Limit validation [1-100]
	if r.Limit < 1 {
		r.Limit = 50 // Default
	}
	if r.Limit > 100 {
		r.Limit = 100 // Max
	}

	return nil
}

// SearchHistoryResponse represents user's search history
type SearchHistoryResponse struct {
	Entries []SearchHistoryEntry `json:"entries"`
}

// SearchHistoryEntry represents a single search from user's history
type SearchHistoryEntry struct {
	QueryText        string    `json:"query_text"`
	CategoryID       *int64    `json:"category_id,omitempty"`
	ResultsCount     int32     `json:"results_count"`
	ClickedListingID *int64    `json:"clicked_listing_id,omitempty"`
	SearchedAt       time.Time `json:"searched_at"`
}
