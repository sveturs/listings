package repository

import (
	"context"

	"github.com/sveturs/listings/internal/domain"
)

// SearchQueriesRepository defines the interface for search query analytics data access
type SearchQueriesRepository interface {
	// CreateSearchQuery logs a search query for analytics
	// This is called asynchronously after every search operation
	// Returns the created SearchQuery with ID
	CreateSearchQuery(ctx context.Context, input *domain.CreateSearchQueryInput) (*domain.SearchQuery, error)

	// GetTrendingQueries returns the most popular search queries within a time range
	// Used for "Trending Searches" feature
	// Performance target: < 500ms for 1M rows (with proper indexes)
	GetTrendingQueries(ctx context.Context, filter *domain.GetTrendingQueriesFilter) ([]domain.TrendingSearch, error)

	// GetUserHistory returns search history for a specific user or session
	// Used for "Recent Searches" feature in search UI
	// Performance target: < 50ms (indexed by user_id/session_id)
	GetUserHistory(ctx context.Context, filter *domain.GetUserHistoryFilter) ([]domain.SearchQuery, error)

	// GetPopularQueries returns all-time popular searches (no time filter)
	// Used for "Popular Searches" suggestions
	// Performance target: < 200ms (cached, refreshed hourly)
	GetPopularQueries(ctx context.Context, filter *domain.GetPopularQueriesFilter) ([]domain.TrendingSearch, error)

	// UpdateClickedListing records that a user clicked on a listing from search results
	// Used for CTR (Click-Through Rate) analysis
	// Called when user clicks a listing in search results
	UpdateClickedListing(ctx context.Context, searchQueryID int64, listingID int64) error

	// GetCTRAnalysis returns click-through rate statistics for queries
	// Used for admin analytics dashboard
	// Performance target: < 300ms
	GetCTRAnalysis(ctx context.Context, filter *domain.GetCTRAnalysisFilter) ([]domain.SearchQueryCTR, error)

	// CleanupOldQueries deletes search queries older than retention period
	// Called by periodic cleanup job (e.g., daily cron)
	// Default retention: 90 days
	CleanupOldQueries(ctx context.Context, daysToKeep int32) (int64, error)
}
