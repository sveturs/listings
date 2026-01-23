// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// StorefrontAnalyticsRepository handles storefront analytics queries
type StorefrontAnalyticsRepository interface {
	// GetStorefrontStats retrieves performance analytics for a storefront
	GetStorefrontStats(ctx context.Context, storefrontID int64, period string) (*domain.StorefrontStats, error)

	// GetTopListings retrieves top-performing listings for a storefront
	GetTopListings(ctx context.Context, storefrontID int64, period string, limit int) ([]*domain.TopListingInfo, error)

	// GetStorefrontOwnerID retrieves the owner user_id for a storefront
	GetStorefrontOwnerID(ctx context.Context, storefrontID int64) (int64, error)
}

// storefrontAnalyticsRepository implements StorefrontAnalyticsRepository
type storefrontAnalyticsRepository struct {
	db     *pgxpool.Pool
	logger zerolog.Logger
}

// NewStorefrontAnalyticsRepository creates a new storefront analytics repository
func NewStorefrontAnalyticsRepository(pool *pgxpool.Pool, logger zerolog.Logger) StorefrontAnalyticsRepository {
	return &storefrontAnalyticsRepository{
		db:     pool,
		logger: logger.With().Str("component", "storefront_analytics_repository").Logger(),
	}
}

// ============================================================================
// PUBLIC METHODS
// ============================================================================

// GetStorefrontStats retrieves performance analytics for a storefront
// Uses the materialized view for fast access with optional period filtering
// Target: < 200ms
func (r *storefrontAnalyticsRepository) GetStorefrontStats(ctx context.Context, storefrontID int64, period string) (*domain.StorefrontStats, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Str("period", period).
		Msg("fetching storefront stats")

	// Get the period start time
	periodStart, err := domain.GetPeriodStartTime(period, time.Now())
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	// Query with period filtering
	// If period is "all", we use the materialized view data directly
	// Otherwise, we need to recalculate with period filter
	var stats domain.StorefrontStats

	if period == "all" {
		// Use materialized view directly (fastest)
		query := `
			SELECT
				storefront_id,
				storefront_name,
				owner_id,
				total_sales,
				total_revenue,
				average_order_value,
				active_listings,
				total_listings,
				total_views,
				total_favorites,
				conversion_rate,
				last_updated_at
			FROM analytics_storefront_stats
			WHERE storefront_id = $1
		`

		err = r.db.QueryRow(ctx, query, storefrontID).Scan(
			&stats.StorefrontID,
			&stats.StorefrontName,
			&stats.OwnerID,
			&stats.TotalSales,
			&stats.TotalRevenue,
			&stats.AverageOrderValue,
			&stats.ActiveListings,
			&stats.TotalListings,
			&stats.TotalViews,
			&stats.TotalFavorites,
			&stats.ConversionRate,
			&stats.LastUpdatedAt,
		)
	} else {
		// Calculate stats with period filter (slower but necessary for time-bounded queries)
		query := `
			WITH storefront_info AS (
				SELECT id, name, user_id
				FROM storefronts
				WHERE id = $1 AND deleted_at IS NULL
			),
			storefront_listings AS (
				SELECT id
				FROM listings
				WHERE storefront_id = $1 AND is_deleted = false
			),
			sales_stats AS (
				SELECT
					COUNT(*) as total_sales,
					COALESCE(SUM(total), 0) as total_revenue,
					CASE
						WHEN COUNT(*) > 0 THEN COALESCE(SUM(total), 0) / COUNT(*)
						ELSE 0
					END as average_order_value
				FROM orders
				WHERE storefront_id = $1
					AND status = 'delivered'
					AND created_at >= $2
			),
			listings_stats AS (
				SELECT
					COUNT(*) FILTER (WHERE status = 'active') as active_listings,
					COUNT(*) as total_listings
				FROM listings
				WHERE storefront_id = $1 AND is_deleted = false
			),
			views_stats AS (
				SELECT COALESCE(COUNT(*), 0) as total_views
				FROM analytics_events
				WHERE event_type = 'view'
					AND entity_type = 'listing'
					AND entity_id IN (SELECT id FROM storefront_listings)
					AND created_at >= $2
			),
			favorites_stats AS (
				SELECT COALESCE(SUM(favorites_count), 0) as total_favorites
				FROM listings
				WHERE storefront_id = $1 AND is_deleted = false
			)
			SELECT
				si.id,
				si.name,
				si.user_id,
				ss.total_sales,
				ss.total_revenue,
				ss.average_order_value,
				ls.active_listings::INT,
				ls.total_listings::INT,
				vs.total_views,
				fs.total_favorites,
				CASE
					WHEN vs.total_views > 0 THEN (ss.total_sales::DECIMAL / vs.total_views) * 100
					ELSE 0
				END as conversion_rate,
				NOW() as last_updated_at
			FROM storefront_info si
			CROSS JOIN sales_stats ss
			CROSS JOIN listings_stats ls
			CROSS JOIN views_stats vs
			CROSS JOIN favorites_stats fs
		`

		err = r.db.QueryRow(ctx, query, storefrontID, periodStart).Scan(
			&stats.StorefrontID,
			&stats.StorefrontName,
			&stats.OwnerID,
			&stats.TotalSales,
			&stats.TotalRevenue,
			&stats.AverageOrderValue,
			&stats.ActiveListings,
			&stats.TotalListings,
			&stats.TotalViews,
			&stats.TotalFavorites,
			&stats.ConversionRate,
			&stats.LastUpdatedAt,
		)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("storefront not found: %d", storefrontID)
		}
		r.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to get storefront stats")
		return nil, fmt.Errorf("failed to get storefront stats: %w", err)
	}

	stats.Period = period
	stats.GeneratedAt = time.Now()

	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int64("total_sales", stats.TotalSales).
		Float64("total_revenue", stats.TotalRevenue).
		Msg("storefront stats retrieved")

	return &stats, nil
}

// GetTopListings retrieves top-performing listings for a storefront
// Orders by revenue DESC, limited to top N listings
// Target: < 100ms
func (r *storefrontAnalyticsRepository) GetTopListings(ctx context.Context, storefrontID int64, period string, limit int) ([]*domain.TopListingInfo, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Str("period", period).
		Int("limit", limit).
		Msg("fetching top listings")

	// Get the period start time
	periodStart, err := domain.GetPeriodStartTime(period, time.Now())
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	// Query top listings by revenue
	query := `
		WITH listing_stats AS (
			SELECT
				l.id as listing_id,
				l.title,
				COALESCE(SUM(oi.price * oi.quantity), 0) as revenue,
				COUNT(DISTINCT o.id) as order_count,
				l.view_count
			FROM listings l
			LEFT JOIN order_items oi ON oi.listing_id = l.id
			LEFT JOIN orders o ON o.id = oi.order_id
				AND o.status = 'delivered'
				AND ($2::TIMESTAMP IS NULL OR o.created_at >= $2)
			WHERE l.storefront_id = $1
				AND l.is_deleted = false
			GROUP BY l.id, l.title, l.view_count
		)
		SELECT
			listing_id,
			title,
			revenue,
			order_count::INT,
			view_count::INT,
			CASE
				WHEN view_count > 0 THEN (order_count::DECIMAL / view_count) * 100
				ELSE 0
			END as conversion_rate
		FROM listing_stats
		WHERE revenue > 0 OR order_count > 0
		ORDER BY revenue DESC, order_count DESC
		LIMIT $3
	`

	var periodStartPtr *time.Time
	if period != "all" {
		periodStartPtr = &periodStart
	}

	rows, err := r.db.Query(ctx, query, storefrontID, periodStartPtr, limit)
	if err != nil {
		r.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to get top listings")
		return nil, fmt.Errorf("failed to get top listings: %w", err)
	}
	defer rows.Close()

	var topListings []*domain.TopListingInfo
	for rows.Next() {
		var listing domain.TopListingInfo
		if err := rows.Scan(
			&listing.ListingID,
			&listing.Title,
			&listing.Revenue,
			&listing.OrderCount,
			&listing.ViewCount,
			&listing.ConversionRate,
		); err != nil {
			r.logger.Error().Err(err).Msg("failed to scan top listing")
			return nil, fmt.Errorf("failed to scan top listing: %w", err)
		}
		topListings = append(topListings, &listing)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating top listings")
		return nil, fmt.Errorf("error iterating top listings: %w", err)
	}

	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Int("count", len(topListings)).
		Msg("top listings retrieved")

	return topListings, nil
}

// GetStorefrontOwnerID retrieves the owner user_id for a storefront
// Used for authorization checks
// Target: < 10ms (simple indexed query)
func (r *storefrontAnalyticsRepository) GetStorefrontOwnerID(ctx context.Context, storefrontID int64) (int64, error) {
	r.logger.Debug().
		Int64("storefront_id", storefrontID).
		Msg("fetching storefront owner")

	var ownerID int64
	query := `SELECT user_id FROM storefronts WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(ctx, query, storefrontID).Scan(&ownerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("storefront not found: %d", storefrontID)
		}
		r.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to get storefront owner")
		return 0, fmt.Errorf("failed to get storefront owner: %w", err)
	}

	return ownerID, nil
}
