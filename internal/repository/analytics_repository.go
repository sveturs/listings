// Package repository defines repository interfaces for the listings microservice.
package repository

import (
	"context"

	"github.com/vondi-global/listings/internal/domain"
)

// AnalyticsRepository defines operations for analytics and metrics tracking
type AnalyticsRepository interface {
	// GetOverviewStats retrieves aggregated platform-wide analytics
	// Returns overview statistics including views, orders, revenue, and engagement metrics
	// Supports time-series data for trend analysis
	// Performance target: < 500ms without cache
	GetOverviewStats(ctx context.Context, filter *domain.GetOverviewStatsFilter) (*domain.OverviewStats, error)

	// GetListingStats retrieves analytics for a specific listing
	// Returns detailed performance metrics including engagement and conversion data
	// Supports time-series data for trend analysis
	// Performance target: < 300ms without cache
	GetListingStats(ctx context.Context, filter *domain.GetListingStatsFilter) (*domain.ListingStats, error)

	// LogEvent records a single analytics event
	// Events are logged asynchronously for high-throughput scenarios
	// Event types: view, favorite, search, order_created, order_completed
	LogEvent(ctx context.Context, eventType, entityType string, entityID int64, userID *int64, sessionID *string, metadata map[string]interface{}) error

	// RefreshMaterializedViews refreshes all analytics materialized views
	// Should be called periodically (e.g., every 15 minutes via cron)
	// Uses CONCURRENTLY to avoid blocking reads
	RefreshMaterializedViews(ctx context.Context) error

	// GetTrendingStats retrieves platform trending analytics
	// Returns trending categories, hot listings, and popular searches
	// Data is pre-calculated in materialized view for optimal performance
	// Performance target: < 100ms (reading from materialized view)
	GetTrendingStats(ctx context.Context) (*domain.TrendingStats, error)
}
