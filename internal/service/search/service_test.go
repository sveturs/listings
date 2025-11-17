package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
