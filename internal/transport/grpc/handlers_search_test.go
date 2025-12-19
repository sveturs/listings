package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	searchv1 "github.com/vondi-global/listings/api/proto/search/v1"
	"github.com/vondi-global/listings/internal/service/search"
)

// ============================================================================
// Mock Search Service
// ============================================================================

type mockSearchService struct {
	searchListingsFunc    func(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error)
	getFacetsFunc         func(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error)
	searchWithFiltersFunc func(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error)
	getSuggestionsFunc    func(ctx context.Context, req *search.SuggestionsRequest) (*search.SuggestionsResponse, error)
	getPopularFunc        func(ctx context.Context, req *search.PopularSearchesRequest) (*search.PopularSearchesResponse, error)
	getSimilarFunc        func(ctx context.Context, listingID int64, limit int32) ([]search.ListingSearchResult, int64, error)
	getTrendingFunc       func(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error)
	getHistoryFunc        func(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error)
}

func (m *mockSearchService) SearchListings(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	if m.searchListingsFunc != nil {
		return m.searchListingsFunc(ctx, req)
	}
	return &search.SearchResponse{}, nil
}

func (m *mockSearchService) GetSearchFacets(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error) {
	if m.getFacetsFunc != nil {
		return m.getFacetsFunc(ctx, req)
	}
	return &search.FacetsResponse{}, nil
}

func (m *mockSearchService) SearchWithFilters(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error) {
	if m.searchWithFiltersFunc != nil {
		return m.searchWithFiltersFunc(ctx, req)
	}
	return &search.SearchFiltersResponse{}, nil
}

func (m *mockSearchService) GetSuggestions(ctx context.Context, req *search.SuggestionsRequest) (*search.SuggestionsResponse, error) {
	if m.getSuggestionsFunc != nil {
		return m.getSuggestionsFunc(ctx, req)
	}
	return &search.SuggestionsResponse{}, nil
}

func (m *mockSearchService) GetPopularSearches(ctx context.Context, req *search.PopularSearchesRequest) (*search.PopularSearchesResponse, error) {
	if m.getPopularFunc != nil {
		return m.getPopularFunc(ctx, req)
	}
	return &search.PopularSearchesResponse{}, nil
}

func (m *mockSearchService) GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]search.ListingSearchResult, int64, error) {
	if m.getSimilarFunc != nil {
		return m.getSimilarFunc(ctx, listingID, limit)
	}
	return []search.ListingSearchResult{}, 0, nil
}

func (m *mockSearchService) GetTrendingSearches(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error) {
	if m.getTrendingFunc != nil {
		return m.getTrendingFunc(ctx, req)
	}
	return &search.TrendingSearchesResponse{}, nil
}

func (m *mockSearchService) GetSearchHistory(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error) {
	if m.getHistoryFunc != nil {
		return m.getHistoryFunc(ctx, req)
	}
	return &search.SearchHistoryResponse{}, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

func newTestSearchHandler(mockSvc *mockSearchService) *SearchHandler {
	logger := zerolog.Nop() // Silent logger for tests
	return &SearchHandler{
		service: mockSvc,
		logger:  logger,
	}
}

func int64PtrTest(v int64) *int64 {
	return &v
}

func stringPtrTest(v string) *string {
	return &v
}

// ============================================================================
// Test GetSearchFacets Handler
// ============================================================================

func TestSearchHandler_GetSearchFacets_Success(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getFacetsFunc: func(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error) {
			assert.Equal(t, "", req.Query)
			assert.Nil(t, req.CategoryID)

			return &search.FacetsResponse{
				Categories: []search.CategoryFacet{
					{CategoryID: "cat-1001", Count: 25},
					{CategoryID: "cat-1002", Count: 18},
				},
				PriceRanges: []search.PriceRangeFacet{
					{Min: 0, Max: 100, Count: 10},
					{Min: 100, Max: 500, Count: 15},
				},
				Attributes:    make(map[string]search.AttributeFacet),
				SourceTypes:   []search.Facet{{Key: "b2c", Count: 30}},
				StockStatuses: []search.Facet{{Key: "in_stock", Count: 40}},
				TookMs:        25,
				Cached:        false,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetSearchFacetsRequest{}

	// Act
	resp, err := handler.GetSearchFacets(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(25), resp.TookMs)
	assert.False(t, resp.Cached)
	assert.Len(t, resp.Categories, 2)
	assert.Len(t, resp.PriceRanges, 2)
	assert.Len(t, resp.SourceTypes, 1)
	assert.Len(t, resp.StockStatuses, 1)
}

func TestSearchHandler_GetSearchFacets_WithFilters(t *testing.T) {
	// Arrange
	categoryID := "cat-1301"
	query := "laptop"
	mockSvc := &mockSearchService{
		getFacetsFunc: func(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error) {
			assert.Equal(t, query, req.Query)
			assert.Equal(t, categoryID, req.CategoryID)

			return &search.FacetsResponse{
				Categories:    []search.CategoryFacet{{CategoryID: categoryID, Count: 50}},
				PriceRanges:   []search.PriceRangeFacet{},
				Attributes:    make(map[string]search.AttributeFacet),
				SourceTypes:   []search.Facet{},
				StockStatuses: []search.Facet{},
				TookMs:        15,
				Cached:        true,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetSearchFacetsRequest{
		Query:      &query,
		CategoryId: categoryID,
	}

	// Act
	resp, err := handler.GetSearchFacets(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Cached)
	assert.Len(t, resp.Categories, 1)
}

func TestSearchHandler_GetSearchFacets_ValidationError(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getFacetsFunc: func(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error) {
			return nil, errors.New("invalid request: query too long")
		},
	}
	handler := newTestSearchHandler(mockSvc)

	longQuery := stringPtr(string(make([]byte, 501))) // 501 chars
	req := &searchv1.GetSearchFacetsRequest{
		Query: longQuery,
	}

	// Act
	resp, err := handler.GetSearchFacets(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSearchHandler_GetSearchFacets_ServiceError(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getFacetsFunc: func(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error) {
			return nil, errors.New("opensearch connection failed")
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetSearchFacetsRequest{}

	// Act
	resp, err := handler.GetSearchFacets(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

// ============================================================================
// Test SearchWithFilters Handler
// ============================================================================

func TestSearchHandler_SearchWithFilters_Success(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		searchWithFiltersFunc: func(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error) {
			assert.Equal(t, "laptop", req.Query)
			assert.Equal(t, int32(20), req.Limit)
			assert.Equal(t, int32(0), req.Offset)

			return &search.SearchFiltersResponse{
				Listings: []search.ListingSearchResult{
					{
						ID:          281,
						UUID:        "test-uuid-1",
						Title:       "Gaming Laptop",
						Price:       1500.00,
						Currency:    "EUR",
						CategoryID:  "cat-1001",
						Status:      "active",
						SourceType:  "b2c",
						StockStatus: "in_stock",
					},
				},
				Total:  1,
				TookMs: 35,
				Cached: false,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.SearchWithFiltersRequest{
		Query:  "laptop",
		Limit:  20,
		Offset: 0,
	}

	// Act
	resp, err := handler.SearchWithFilters(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Listings, 1)
	assert.Equal(t, "Gaming Laptop", resp.Listings[0].Title)
	assert.False(t, resp.Cached)
}

func TestSearchHandler_SearchWithFilters_WithFacets(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		searchWithFiltersFunc: func(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error) {
			assert.True(t, req.IncludeFacets)

			return &search.SearchFiltersResponse{
				Listings: []search.ListingSearchResult{},
				Total:    0,
				TookMs:   20,
				Cached:   false,
				Facets: &search.FacetsResponse{
					Categories: []search.CategoryFacet{
						{CategoryID: "cat-1001", Count: 10},
					},
					PriceRanges:   []search.PriceRangeFacet{},
					Attributes:    make(map[string]search.AttributeFacet),
					SourceTypes:   []search.Facet{},
					StockStatuses: []search.Facet{},
					TookMs:        20,
					Cached:        false,
				},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.SearchWithFiltersRequest{
		Query:         "test",
		Limit:         20,
		IncludeFacets: true,
	}

	// Act
	resp, err := handler.SearchWithFilters(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Facets)
	assert.Len(t, resp.Facets.Categories, 1)
}

func TestSearchHandler_SearchWithFilters_ValidationError(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		searchWithFiltersFunc: func(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error) {
			return nil, errors.New("invalid limit: must be between 1 and 100")
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.SearchWithFiltersRequest{
		Query:  "test",
		Limit:  150, // Invalid: too large
		Offset: 0,
	}

	// Act
	resp, err := handler.SearchWithFilters(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSearchHandler_SearchWithFilters_EmptyResults(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		searchWithFiltersFunc: func(ctx context.Context, req *search.SearchFiltersRequest) (*search.SearchFiltersResponse, error) {
			return &search.SearchFiltersResponse{
				Listings: []search.ListingSearchResult{},
				Total:    0,
				TookMs:   10,
				Cached:   false,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.SearchWithFiltersRequest{
		Query: "nonexistent",
		Limit: 20,
	}

	// Act
	resp, err := handler.SearchWithFilters(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Listings, 0)
}

// ============================================================================
// Test GetSuggestions Handler
// ============================================================================

func TestSearchHandler_GetSuggestions_Success(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getSuggestionsFunc: func(ctx context.Context, req *search.SuggestionsRequest) (*search.SuggestionsResponse, error) {
			assert.Equal(t, "lap", req.Prefix)
			assert.Equal(t, int32(10), req.Limit)

			return &search.SuggestionsResponse{
				Suggestions: []search.Suggestion{
					{Text: "laptop", Score: 0.95, ListingID: nil},
					{Text: "laptop bag", Score: 0.78, ListingID: nil},
					{Text: "laptop stand", Score: 0.65, ListingID: nil},
				},
				TookMs: 5,
				Cached: false,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetSuggestionsRequest{
		Prefix: "lap",
		Limit:  10,
	}

	// Act
	resp, err := handler.GetSuggestions(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Suggestions, 3)
	assert.Equal(t, "laptop", resp.Suggestions[0].Text)
	assert.InDelta(t, 0.95, resp.Suggestions[0].Score, 0.01)
}

func TestSearchHandler_GetSuggestions_PrefixTooShort(t *testing.T) {
	// Arrange
	handler := newTestSearchHandler(&mockSearchService{})

	req := &searchv1.GetSuggestionsRequest{
		Prefix: "a", // Too short (< 2 chars)
		Limit:  10,
	}

	// Act
	resp, err := handler.GetSuggestions(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "at least 2 characters")
}

func TestSearchHandler_GetSuggestions_ValidationError(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getSuggestionsFunc: func(ctx context.Context, req *search.SuggestionsRequest) (*search.SuggestionsResponse, error) {
			return nil, errors.New("invalid limit")
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetSuggestionsRequest{
		Prefix: "test",
		Limit:  -1, // Invalid
	}

	// Act
	resp, err := handler.GetSuggestions(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSearchHandler_GetSuggestions_NoResults(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getSuggestionsFunc: func(ctx context.Context, req *search.SuggestionsRequest) (*search.SuggestionsResponse, error) {
			return &search.SuggestionsResponse{
				Suggestions: []search.Suggestion{},
				TookMs:      3,
				Cached:      false,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetSuggestionsRequest{
		Prefix: "xyz",
		Limit:  10,
	}

	// Act
	resp, err := handler.GetSuggestions(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Suggestions, 0)
}

// ============================================================================
// Test GetPopularSearches Handler
// ============================================================================

func TestSearchHandler_GetPopularSearches_Success(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getPopularFunc: func(ctx context.Context, req *search.PopularSearchesRequest) (*search.PopularSearchesResponse, error) {
			assert.Equal(t, int32(10), req.Limit)
			assert.Equal(t, "24h", req.TimeRange)

			return &search.PopularSearchesResponse{
				Searches: []search.PopularSearch{
					{Query: "iphone", SearchCount: 1203, TrendScore: 22.5},
					{Query: "laptop", SearchCount: 891, TrendScore: 15.3},
					{Query: "headphones", SearchCount: 654, TrendScore: -3.2},
				},
				TookMs: 8,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetPopularSearchesRequest{
		Limit: 10,
	}

	// Act
	resp, err := handler.GetPopularSearches(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Searches, 3)
	assert.Equal(t, "iphone", resp.Searches[0].Query)
	assert.Equal(t, int64(1203), resp.Searches[0].SearchCount)
	assert.InDelta(t, 22.5, resp.Searches[0].TrendScore, 0.01)
}

func TestSearchHandler_GetPopularSearches_WithCategoryFilter(t *testing.T) {
	// Arrange
	categoryID := "cat-1301"
	mockSvc := &mockSearchService{
		getPopularFunc: func(ctx context.Context, req *search.PopularSearchesRequest) (*search.PopularSearchesResponse, error) {
			assert.Equal(t, categoryID, req.CategoryID)

			return &search.PopularSearchesResponse{
				Searches: []search.PopularSearch{
					{Query: "lamborghini", SearchCount: 542, TrendScore: 15.3},
				},
				TookMs: 5,
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	req := &searchv1.GetPopularSearchesRequest{
		CategoryId: categoryID,
		Limit:      10,
	}

	// Act
	resp, err := handler.GetPopularSearches(context.Background(), req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Searches, 1)
	assert.Equal(t, "lamborghini", resp.Searches[0].Query)
}

func TestSearchHandler_GetPopularSearches_ValidationError(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getPopularFunc: func(ctx context.Context, req *search.PopularSearchesRequest) (*search.PopularSearchesResponse, error) {
			return nil, errors.New("invalid time_range")
		},
	}
	handler := newTestSearchHandler(mockSvc)

	invalidTimeRange := "invalid"
	req := &searchv1.GetPopularSearchesRequest{
		TimeRange: &invalidTimeRange,
		Limit:     10,
	}

	// Act
	resp, err := handler.GetPopularSearches(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// ============================================================================
// Test Error Mapping
// ============================================================================

func TestSearchHandler_ErrorMapping(t *testing.T) {
	tests := []struct {
		name         string
		serviceError error
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name:         "Invalid request error maps to InvalidArgument",
			serviceError: errors.New("invalid limit: must be positive"),
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "invalid",
		},
		{
			name:         "Generic error maps to Internal",
			serviceError: errors.New("database connection failed"),
			expectedCode: codes.Internal,
			expectedMsg:  "failed",
		},
		{
			name:         "Search failed error maps to Internal",
			serviceError: errors.New("opensearch timeout"),
			expectedCode: codes.Internal,
			expectedMsg:  "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockSvc := &mockSearchService{
				getFacetsFunc: func(ctx context.Context, req *search.FacetsRequest) (*search.FacetsResponse, error) {
					return nil, tt.serviceError
				},
			}
			handler := newTestSearchHandler(mockSvc)

			req := &searchv1.GetSearchFacetsRequest{}

			// Act
			resp, err := handler.GetSearchFacets(context.Background(), req)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, resp)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Contains(t, st.Message(), tt.expectedMsg)
		})
	}
}

// ============================================================================
// Test Helper Functions
// ============================================================================

func TestContainsErrorFunction(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		substr   string
		expected bool
	}{
		{"Exact match", errors.New("invalid"), "invalid", true},
		{"Case insensitive", errors.New("Invalid Request"), "invalid", true},
		{"Substring", errors.New("This is an invalid request"), "invalid", true},
		{"Not found", errors.New("database error"), "invalid", false},
		{"Nil error", nil, "invalid", false},
		{"Empty substring", errors.New("invalid"), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsError(tt.err, tt.substr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// PHASE 28: Test GetTrendingSearches Handler
// ============================================================================

func TestSearchHandler_GetTrendingSearches_Success(t *testing.T) {
	// Arrange
	now := testTimeNow()
	mockSvc := &mockSearchService{
		getTrendingFunc: func(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error) {
			return &search.TrendingSearchesResponse{
				Searches: []search.TrendingSearchResult{
					{QueryText: "iphone", SearchCount: 542, LastSearched: now},
					{QueryText: "laptop", SearchCount: 321, LastSearched: now},
				},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 10,
		Days:  7,
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Searches, 2)
	assert.Equal(t, "iphone", resp.Searches[0].QueryText)
	assert.Equal(t, int32(542), resp.Searches[0].SearchCount)
	assert.NotNil(t, resp.Searches[0].LastSearched)
	assert.Equal(t, "laptop", resp.Searches[1].QueryText)
}

func TestSearchHandler_GetTrendingSearches_WithCategoryFilter(t *testing.T) {
	// Arrange
	now := testTimeNow()
	categoryID := "cat-1001"
	mockSvc := &mockSearchService{
		getTrendingFunc: func(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error) {
			assert.NotNil(t, req.CategoryID)
			assert.Equal(t, categoryID, req.CategoryID)
			return &search.TrendingSearchesResponse{
				Searches: []search.TrendingSearchResult{
					{QueryText: "macbook", SearchCount: 150, LastSearched: now},
				},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		CategoryId: categoryID,
		Limit:      10,
		Days:       7,
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Searches, 1)
	assert.Equal(t, "macbook", resp.Searches[0].QueryText)
}

func TestSearchHandler_GetTrendingSearches_InvalidLimit_TooLow(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 0, // Invalid
		Days:  7,
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "limit must be at least 1")
}

func TestSearchHandler_GetTrendingSearches_InvalidLimit_TooHigh(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 51, // Invalid (max 50)
		Days:  7,
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "limit must not exceed 50")
}

func TestSearchHandler_GetTrendingSearches_InvalidDays_TooLow(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 10,
		Days:  0, // Invalid
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "days must be at least 1")
}

func TestSearchHandler_GetTrendingSearches_InvalidDays_TooHigh(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 10,
		Days:  31, // Invalid (max 30)
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), "days must not exceed 30")
}

func TestSearchHandler_GetTrendingSearches_EmptyResult(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getTrendingFunc: func(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error) {
			return &search.TrendingSearchesResponse{
				Searches: []search.TrendingSearchResult{},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 10,
		Days:  7,
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Searches)
}

func TestSearchHandler_GetTrendingSearches_ServiceError(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getTrendingFunc: func(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 10,
		Days:  7,
	}

	// Act
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "failed to get trending searches")
}

func TestSearchHandler_GetTrendingSearches_DefaultParams(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{
		getTrendingFunc: func(ctx context.Context, req *search.TrendingSearchesRequest) (*search.TrendingSearchesResponse, error) {
			// Validate defaults are applied
			assert.Equal(t, int32(10), req.Limit)
			assert.Equal(t, int32(7), req.Days)
			return &search.TrendingSearchesResponse{
				Searches: []search.TrendingSearchResult{},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	// Request with zeros (should get defaults)
	protoReq := &searchv1.GetTrendingSearchesRequest{
		Limit: 0,
		Days:  0,
	}

	// Act
	// Note: Validation will reject Limit=0 and Days=0
	// So we need to provide valid values
	protoReq.Limit = 10
	protoReq.Days = 7
	resp, err := handler.GetTrendingSearches(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// Helper function for tests
func testTimeNow() time.Time {
	return time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
}

// ============================================================================
// PHASE 28: GetSearchHistory Tests
// ============================================================================

func TestGetSearchHistory_Success_UserID(t *testing.T) {
	// Arrange
	userID := int64(123)
	mockSvc := &mockSearchService{
		getHistoryFunc: func(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error) {
			assert.NotNil(t, req.UserID)
			assert.Equal(t, userID, *req.UserID)
			assert.Nil(t, req.SessionID)
			assert.Equal(t, int32(50), req.Limit)

			return &search.SearchHistoryResponse{
				Entries: []search.SearchHistoryEntry{
					{
						QueryText:    "iphone 15",
						CategoryID:   stringPtrTest("cat-1001"),
						ResultsCount: 42,
						SearchedAt:   testTimeNow(),
					},
					{
						QueryText:    "macbook pro",
						CategoryID:   stringPtrTest("cat-1001"),
						ResultsCount: 15,
						SearchedAt:   testTimeNow().Add(-1 * time.Hour),
					},
				},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		UserId: &userID,
		Limit:  50,
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Entries, 2)

	// Verify first entry
	assert.Equal(t, "iphone 15", resp.Entries[0].QueryText)
	assert.NotNil(t, resp.Entries[0].CategoryId)
	assert.Equal(t, "cat-1001", resp.Entries[0].CategoryId)
	assert.Equal(t, int32(42), resp.Entries[0].ResultsCount)
	assert.NotNil(t, resp.Entries[0].SearchedAt)

	// Verify second entry
	assert.Equal(t, "macbook pro", resp.Entries[1].QueryText)
}

func TestGetSearchHistory_Success_SessionID(t *testing.T) {
	// Arrange
	sessionID := "550e8400-e29b-41d4-a716-446655440000"
	mockSvc := &mockSearchService{
		getHistoryFunc: func(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error) {
			assert.Nil(t, req.UserID)
			assert.NotNil(t, req.SessionID)
			assert.Equal(t, sessionID, *req.SessionID)
			assert.Equal(t, int32(20), req.Limit)

			return &search.SearchHistoryResponse{
				Entries: []search.SearchHistoryEntry{
					{
						QueryText:        "tesla model 3",
						CategoryID:       stringPtrTest("cat-1301"),
						ResultsCount:     8,
						ClickedListingID: int64PtrTest(999),
						SearchedAt:       testTimeNow(),
					},
				},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		SessionId: &sessionID,
		Limit:     20,
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Entries, 1)

	// Verify entry with clicked listing
	assert.Equal(t, "tesla model 3", resp.Entries[0].QueryText)
	assert.NotNil(t, resp.Entries[0].ClickedListingId)
	assert.Equal(t, int64(999), *resp.Entries[0].ClickedListingId)
}

func TestGetSearchHistory_MissingBothIDs(t *testing.T) {
	// Arrange
	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		Limit: 50,
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify error code
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "either user_id or session_id must be provided")
}

func TestGetSearchHistory_BothIDsProvided(t *testing.T) {
	// Arrange
	userID := int64(123)
	sessionID := "550e8400-e29b-41d4-a716-446655440000"

	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		UserId:    &userID,
		SessionId: &sessionID, // Both provided - XOR violation
		Limit:     50,
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify error code
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "only one of user_id or session_id must be provided")
}

func TestGetSearchHistory_InvalidLimit(t *testing.T) {
	// Arrange
	userID := int64(123)

	mockSvc := &mockSearchService{}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		UserId: &userID,
		Limit:  150, // Exceeds max of 100
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify error code
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "limit")
}

func TestGetSearchHistory_EmptyHistory(t *testing.T) {
	// Arrange
	userID := int64(999)
	mockSvc := &mockSearchService{
		getHistoryFunc: func(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error) {
			return &search.SearchHistoryResponse{
				Entries: []search.SearchHistoryEntry{}, // Empty history
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		UserId: &userID,
		Limit:  50,
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Entries)
}

func TestGetSearchHistory_ServiceError(t *testing.T) {
	// Arrange
	userID := int64(123)
	serviceErr := errors.New("database connection failed")

	mockSvc := &mockSearchService{
		getHistoryFunc: func(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error) {
			return nil, serviceErr
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		UserId: &userID,
		Limit:  50,
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)

	// Verify error code
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to get search history")
}

func TestGetSearchHistory_DefaultLimit(t *testing.T) {
	// Arrange
	userID := int64(123)
	mockSvc := &mockSearchService{
		getHistoryFunc: func(ctx context.Context, req *search.SearchHistoryRequest) (*search.SearchHistoryResponse, error) {
			// Should use default limit of 50
			assert.Equal(t, int32(50), req.Limit)

			return &search.SearchHistoryResponse{
				Entries: []search.SearchHistoryEntry{},
			}, nil
		},
	}
	handler := newTestSearchHandler(mockSvc)

	protoReq := &searchv1.GetSearchHistoryRequest{
		UserId: &userID,
		Limit:  0, // Will be set to default 50
	}

	// Act
	resp, err := handler.GetSearchHistory(context.Background(), protoReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
