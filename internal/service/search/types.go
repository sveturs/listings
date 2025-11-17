package search

// SearchRequest represents domain search parameters
type SearchRequest struct {
	Query      string  // Search query text
	CategoryID *int64  // Optional category filter
	Limit      int32   // Results per page (1-100)
	Offset     int32   // Pagination offset
	UseCache   bool    // Whether to use cache
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
	ID            int64                `json:"id"`
	UUID          string               `json:"uuid"`
	Title         string               `json:"title"`
	Description   *string              `json:"description,omitempty"`
	Price         float64              `json:"price"`
	Currency      string               `json:"currency"`
	CategoryID    int64                `json:"category_id"`
	Status        string               `json:"status"`
	Images        []ListingImageResult `json:"images,omitempty"`
	CreatedAt     string               `json:"created_at"`
	UserID        int64                `json:"user_id"`
	StorefrontID  *int64               `json:"storefront_id,omitempty"`
	Quantity      int32                `json:"quantity"`
	SKU           *string              `json:"sku,omitempty"`
	SourceType    string               `json:"source_type"`
	StockStatus   string               `json:"stock_status"`
}

// ListingImageResult represents an image in search results
type ListingImageResult struct {
	ID           int64  `json:"id"`
	URL          string `json:"url"`
	IsPrimary    bool   `json:"is_primary"`
	DisplayOrder int32  `json:"display_order"`
}
