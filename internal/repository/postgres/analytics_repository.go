// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// analyticsRepository implements AnalyticsRepository using PostgreSQL with optimized queries
type analyticsRepository struct {
	db     *pgxpool.Pool
	logger zerolog.Logger
}

// NewAnalyticsRepository creates a new analytics repository
func NewAnalyticsRepository(pool *pgxpool.Pool, logger zerolog.Logger) *analyticsRepository {
	return &analyticsRepository{
		db:     pool,
		logger: logger.With().Str("component", "analytics_repository").Logger(),
	}
}

// ============================================================================
// PUBLIC METHODS
// ============================================================================

// GetOverviewStats retrieves aggregated platform-wide analytics
// Uses CTEs and materialized views for optimal performance
// Target: < 500ms without cache
func (r *analyticsRepository) GetOverviewStats(ctx context.Context, filter *domain.GetOverviewStatsFilter) (*domain.OverviewStats, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	r.logger.Debug().
		Time("start_date", filter.StartDate).
		Time("end_date", filter.EndDate).
		Str("granularity", filter.Granularity).
		Msg("fetching overview stats")

	// Build dynamic WHERE clause for optional filters
	whereConditions := []string{"created_at >= $1 AND created_at <= $2"}
	args := []interface{}{filter.StartDate, filter.EndDate}
	argIdx := 3

	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, *filter.UserID)
		argIdx++
	}

	if filter.StorefrontID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("metadata->>'storefront_id' = $%d::TEXT", argIdx))
		args = append(args, *filter.StorefrontID)
		argIdx++
	}

	if filter.CategoryID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("metadata->>'category_id' = $%d::TEXT", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}

	if filter.SourceType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("metadata->>'source_type' = $%d", argIdx))
		args = append(args, *filter.SourceType)
		argIdx++
	}

	whereClause := whereConditions[0]
	for i := 1; i < len(whereConditions); i++ {
		whereClause += " AND " + whereConditions[i]
	}

	// Main aggregation query with CTEs for performance
	query := fmt.Sprintf(`
		WITH event_metrics AS (
			SELECT
				COUNT(*) FILTER (WHERE event_type = 'view' AND entity_type = 'listing') AS total_views,
				COUNT(*) FILTER (WHERE event_type = 'favorite' AND entity_type = 'listing') AS total_favorites,
				COUNT(*) FILTER (WHERE event_type = 'order_created') AS total_orders,
				COALESCE(SUM((metadata->>'total_amount')::DECIMAL) FILTER (
					WHERE event_type = 'order_completed'
				), 0) AS total_revenue,
				COUNT(DISTINCT entity_id) FILTER (WHERE entity_type = 'listing') AS active_listings,
				COUNT(DISTINCT user_id) FILTER (WHERE user_id IS NOT NULL) AS active_users
			FROM analytics_events
			WHERE %s
		)
		SELECT
			total_views,
			total_favorites,
			total_orders,
			total_revenue,
			active_listings,
			active_users,
			CASE
				WHEN total_views > 0 THEN (total_orders::DECIMAL / total_views) * 100
				ELSE 0
			END AS conversion_rate,
			CASE
				WHEN total_orders > 0 THEN total_revenue / total_orders
				ELSE 0
			END AS average_order_value,
			CASE
				WHEN total_views > 0 THEN (total_favorites::DECIMAL / total_views) * 100
				ELSE 0
			END AS average_favorites_rate
		FROM event_metrics
	`, whereClause)

	var stats domain.OverviewStats
	stats.PeriodStart = filter.StartDate
	stats.PeriodEnd = filter.EndDate

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stats.TotalViews,
		&stats.TotalFavorites,
		&stats.TotalOrders,
		&stats.TotalRevenue,
		&stats.ActiveListings,
		&stats.ActiveUsers,
		&stats.ConversionRate,
		&stats.AverageOrderValue,
		&stats.AverageFavoritesRate,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Return empty stats instead of error
			stats.EnrichWithCalculatedFields()
			return &stats, nil
		}
		r.logger.Error().Err(err).Msg("failed to get overview stats")
		return nil, fmt.Errorf("failed to get overview stats: %w", err)
	}

	// Fetch time series data if requested
	if filter.Limit > 0 {
		timeSeries, err := r.getOverviewTimeSeries(ctx, filter, whereClause, args)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to get time series data")
			// Don't fail the entire request, just log and continue
		} else {
			stats.TimeSeries = timeSeries
		}
	}

	stats.EnrichWithCalculatedFields()

	r.logger.Debug().
		Int64("total_views", stats.TotalViews).
		Int64("total_orders", stats.TotalOrders).
		Float64("total_revenue", stats.TotalRevenue).
		Msg("overview stats fetched successfully")

	return &stats, nil
}

// GetListingStats retrieves analytics for a specific listing
// Uses optimized queries with window functions for time-series data
// Target: < 300ms without cache
func (r *analyticsRepository) GetListingStats(ctx context.Context, filter *domain.GetListingStatsFilter) (*domain.ListingStats, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Resolve listing ID if UUID provided
	var listingID int64
	if filter.ListingID != nil {
		listingID = *filter.ListingID
	} else if filter.ListingUUID != nil {
		var err error
		listingID, err = r.resolveListingID(ctx, *filter.ListingUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve listing UUID: %w", err)
		}
	}

	r.logger.Debug().
		Int64("listing_id", listingID).
		Msg("fetching listing stats")

	// Build time range filter
	var timeFilter string
	var args []interface{}
	args = append(args, listingID)
	argIdx := 2

	if filter.StartDate != nil && filter.EndDate != nil {
		timeFilter = fmt.Sprintf("AND created_at >= $%d AND created_at <= $%d", argIdx, argIdx+1)
		args = append(args, *filter.StartDate, *filter.EndDate)
		argIdx += 2
	}

	// Main aggregation query
	query := fmt.Sprintf(`
		WITH listing_info AS (
			SELECT id, uuid
			FROM listings
			WHERE id = $1
		),
		event_metrics AS (
			SELECT
				COUNT(*) FILTER (WHERE event_type = 'view') AS views_count,
				COUNT(*) FILTER (WHERE event_type = 'favorite') AS favorites_count,
				COUNT(*) FILTER (WHERE event_type = 'inquiry') AS inquiries_count,
				COUNT(*) FILTER (WHERE event_type = 'order_created') AS orders_count,
				COALESCE(SUM((metadata->>'total_amount')::DECIMAL) FILTER (
					WHERE event_type = 'order_completed'
				), 0) AS total_revenue,
				MIN(created_at) FILTER (WHERE event_type = 'view') AS first_viewed_at,
				MAX(created_at) FILTER (WHERE event_type = 'view') AS last_viewed_at,
				MAX(created_at) FILTER (WHERE event_type = 'favorite') AS last_favorited_at,
				MAX(created_at) FILTER (WHERE event_type = 'order_completed') AS last_ordered_at
			FROM analytics_events
			WHERE entity_type = 'listing' AND entity_id = $1 %s
		)
		SELECT
			li.id,
			li.uuid,
			em.views_count,
			em.favorites_count,
			em.inquiries_count,
			em.orders_count,
			em.total_revenue,
			em.first_viewed_at,
			em.last_viewed_at,
			em.last_favorited_at,
			em.last_ordered_at
		FROM listing_info li
		CROSS JOIN event_metrics em
	`, timeFilter)

	var stats domain.ListingStats
	var uuid string
	var firstViewedAt, lastViewedAt, lastFavoritedAt, lastOrderedAt *time.Time

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stats.ListingID,
		&uuid,
		&stats.ViewsCount,
		&stats.FavoritesCount,
		&stats.InquiriesCount,
		&stats.OrdersCount,
		&stats.TotalRevenue,
		&firstViewedAt,
		&lastViewedAt,
		&lastFavoritedAt,
		&lastOrderedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("listing not found")
		}
		r.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to get listing stats")
		return nil, fmt.Errorf("failed to get listing stats: %w", err)
	}

	stats.UUID = uuid
	stats.FirstViewedAt = firstViewedAt
	stats.LastViewedAt = lastViewedAt
	stats.LastFavoritedAt = lastFavoritedAt
	stats.LastOrderedAt = lastOrderedAt

	// Set period
	if filter.StartDate != nil {
		stats.PeriodStart = *filter.StartDate
	} else if firstViewedAt != nil {
		stats.PeriodStart = *firstViewedAt
	} else {
		stats.PeriodStart = time.Now().AddDate(0, 0, -30) // Default to 30 days ago
	}

	if filter.EndDate != nil {
		stats.PeriodEnd = *filter.EndDate
	} else {
		stats.PeriodEnd = time.Now()
	}

	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()

	// Fetch time series if requested
	if filter.IncludeTimeSeries {
		timeSeries, err := r.getListingTimeSeries(ctx, listingID, filter)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to get listing time series")
			// Don't fail the entire request, just log and continue
		} else {
			stats.TimeSeries = timeSeries
		}
	}

	stats.EnrichWithCalculatedFields()

	r.logger.Debug().
		Int64("listing_id", listingID).
		Int64("views", stats.ViewsCount).
		Int64("orders", stats.OrdersCount).
		Float64("revenue", stats.TotalRevenue).
		Msg("listing stats fetched successfully")

	return &stats, nil
}

// LogEvent records a single analytics event
// Uses prepared statement for high-throughput logging
func (r *analyticsRepository) LogEvent(ctx context.Context, eventType, entityType string, entityID int64, userID *int64, sessionID *string, metadata map[string]interface{}) error {
	// Convert metadata to JSONB
	var metadataJSON []byte
	var err error
	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	} else {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO analytics_events (
			event_type, entity_type, entity_id,
			user_id, session_id, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING uuid
	`

	var eventUUID string
	err = r.db.QueryRow(ctx, query,
		eventType, entityType, entityID,
		userID, sessionID, metadataJSON,
	).Scan(&eventUUID)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("event_type", eventType).
			Str("entity_type", entityType).
			Int64("entity_id", entityID).
			Msg("failed to log analytics event")
		return fmt.Errorf("failed to log analytics event: %w", err)
	}

	r.logger.Debug().
		Str("event_uuid", eventUUID).
		Str("event_type", eventType).
		Str("entity_type", entityType).
		Int64("entity_id", entityID).
		Msg("analytics event logged")

	return nil
}

// RefreshMaterializedViews refreshes all analytics materialized views concurrently
// Should be called periodically (e.g., every 15 minutes via cron)
func (r *analyticsRepository) RefreshMaterializedViews(ctx context.Context) error {
	r.logger.Info().Msg("refreshing materialized views")

	// Refresh both views concurrently
	query := `
		SELECT refresh_analytics_views();
	`

	_, err := r.db.Exec(ctx, query)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to refresh materialized views")
		return fmt.Errorf("failed to refresh materialized views: %w", err)
	}

	r.logger.Info().Msg("materialized views refreshed successfully")
	return nil
}

// ============================================================================
// PRIVATE HELPER METHODS
// ============================================================================

// getOverviewTimeSeries retrieves time-series data for overview stats
func (r *analyticsRepository) getOverviewTimeSeries(ctx context.Context, filter *domain.GetOverviewStatsFilter, whereClause string, baseArgs []interface{}) ([]*domain.TimeSeriesDataPoint, error) {
	// Determine time bucket based on granularity
	var timeBucket string
	switch filter.Granularity {
	case "hourly":
		timeBucket = "1 hour"
	case "daily":
		timeBucket = "1 day"
	default:
		timeBucket = "1 day"
	}

	// Build time series query with window functions
	query := fmt.Sprintf(`
		SELECT
			time_bucket,
			SUM(views) AS views,
			SUM(favorites) AS favorites,
			SUM(orders) AS orders,
			SUM(revenue) AS revenue,
			COUNT(DISTINCT user_id) AS active_users
		FROM (
			SELECT
				date_trunc('%s', created_at) AS time_bucket,
				COUNT(*) FILTER (WHERE event_type = 'view' AND entity_type = 'listing') AS views,
				COUNT(*) FILTER (WHERE event_type = 'favorite' AND entity_type = 'listing') AS favorites,
				COUNT(*) FILTER (WHERE event_type = 'order_created') AS orders,
				COALESCE(SUM((metadata->>'total_amount')::DECIMAL) FILTER (
					WHERE event_type = 'order_completed'
				), 0) AS revenue,
				user_id
			FROM analytics_events
			WHERE %s
			GROUP BY time_bucket, user_id
		) AS bucketed_events
		GROUP BY time_bucket
		ORDER BY time_bucket DESC
		LIMIT $%d OFFSET $%d
	`, timeBucket, whereClause, len(baseArgs)+1, len(baseArgs)+2)

	args := append(baseArgs, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query time series: %w", err)
	}
	defer rows.Close()

	var points []*domain.TimeSeriesDataPoint
	for rows.Next() {
		var point domain.TimeSeriesDataPoint
		err := rows.Scan(
			&point.Timestamp,
			&point.Views,
			&point.Favorites,
			&point.Orders,
			&point.Revenue,
			&point.ActiveUsers,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan time series point: %w", err)
		}

		point.ConversionRate = domain.CalculateConversionRate(point.Views, point.Orders)
		points = append(points, &point)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating time series rows: %w", err)
	}

	return points, nil
}

// getListingTimeSeries retrieves time-series data for a specific listing
func (r *analyticsRepository) getListingTimeSeries(ctx context.Context, listingID int64, filter *domain.GetListingStatsFilter) ([]*domain.ListingTimeSeriesPoint, error) {
	// Determine time bucket based on granularity
	var timeBucket string
	switch filter.Granularity {
	case "hourly":
		timeBucket = "1 hour"
	case "daily":
		timeBucket = "1 day"
	default:
		timeBucket = "1 day"
	}

	// Build time range filter
	var timeFilter string
	var args []interface{}
	args = append(args, listingID)
	argIdx := 2

	if filter.StartDate != nil && filter.EndDate != nil {
		timeFilter = fmt.Sprintf("AND created_at >= $%d AND created_at <= $%d", argIdx, argIdx+1)
		args = append(args, *filter.StartDate, *filter.EndDate)
		argIdx += 2
	}

	query := fmt.Sprintf(`
		SELECT
			date_trunc('%s', created_at) AS time_bucket,
			COUNT(*) FILTER (WHERE event_type = 'view') AS views,
			COUNT(*) FILTER (WHERE event_type = 'favorite') AS favorites,
			COUNT(*) FILTER (WHERE event_type = 'order_created') AS orders,
			COALESCE(SUM((metadata->>'total_amount')::DECIMAL) FILTER (
				WHERE event_type = 'order_completed'
			), 0) AS revenue
		FROM analytics_events
		WHERE entity_type = 'listing' AND entity_id = $1 %s
		GROUP BY time_bucket
		ORDER BY time_bucket DESC
		LIMIT $%d
	`, timeBucket, timeFilter, argIdx)

	args = append(args, filter.TimeSeriesLimit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query listing time series: %w", err)
	}
	defer rows.Close()

	var points []*domain.ListingTimeSeriesPoint
	for rows.Next() {
		var point domain.ListingTimeSeriesPoint
		point.ListingID = listingID

		err := rows.Scan(
			&point.Timestamp,
			&point.Views,
			&point.Favorites,
			&point.Orders,
			&point.Revenue,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing time series point: %w", err)
		}

		point.ConversionRate = domain.CalculateConversionRate(point.Views, point.Orders)
		points = append(points, &point)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listing time series rows: %w", err)
	}

	return points, nil
}

// resolveListingID resolves a listing UUID to its ID
func (r *analyticsRepository) resolveListingID(ctx context.Context, uuid string) (int64, error) {
	query := `SELECT id FROM listings WHERE uuid = $1`

	var listingID int64
	err := r.db.QueryRow(ctx, query, uuid).Scan(&listingID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("listing not found with UUID: %s", uuid)
		}
		return 0, fmt.Errorf("failed to resolve listing UUID: %w", err)
	}

	return listingID, nil
}

// GetTrendingStats retrieves platform trending analytics from materialized view
// This method reads pre-calculated data from analytics_trending_cache for optimal performance
// Target: < 100ms (reading from materialized view)
func (r *analyticsRepository) GetTrendingStats(ctx context.Context) (*domain.TrendingStats, error) {
	r.logger.Debug().Msg("fetching trending stats from materialized view")

	// Read from materialized view
	query := `SELECT trending_data FROM analytics_trending_cache LIMIT 1`

	var trendingDataJSON []byte
	err := r.db.QueryRow(ctx, query).Scan(&trendingDataJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn().Msg("trending cache is empty, returning empty stats")
			return &domain.TrendingStats{
				TrendingCategories: []*domain.TrendingCategory{},
				HotListings:        []*domain.HotListing{},
				PopularSearches:    []*domain.PopularSearch{},
				GeneratedAt:        time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to query trending stats: %w", err)
	}

	// Parse JSON data
	var rawData struct {
		TrendingCategories []struct {
			CategoryID    int64   `json:"category_id"`
			CategoryName  string  `json:"category_name"`
			OrderCount30d int32   `json:"order_count_30d"`
			OrderCount7d  int32   `json:"order_count_7d"`
			GrowthRate    float64 `json:"growth_rate"`
			TrendScore    float64 `json:"trend_score"`
		} `json:"trending_categories"`
		HotListings []struct {
			ListingID       int64   `json:"listing_id"`
			Title           string  `json:"title"`
			Orders24h       int64   `json:"orders_24h"`
			Orders7d        int64   `json:"orders_7d"`
			OrdersGrowth    float64 `json:"orders_growth"`
			QuantitySold24h int64   `json:"quantity_sold_24h"`
			Price           float64 `json:"price"`
		} `json:"hot_listings"`
		PopularSearches []struct {
			Query       string `json:"query"`
			SearchCount int64  `json:"search_count"`
		} `json:"popular_searches"`
		GeneratedAt time.Time `json:"generated_at"`
	}

	if err := json.Unmarshal(trendingDataJSON, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal trending data: %w", err)
	}

	// Convert to domain models
	stats := &domain.TrendingStats{
		TrendingCategories: make([]*domain.TrendingCategory, 0, len(rawData.TrendingCategories)),
		HotListings:        make([]*domain.HotListing, 0, len(rawData.HotListings)),
		PopularSearches:    make([]*domain.PopularSearch, 0, len(rawData.PopularSearches)),
		GeneratedAt:        rawData.GeneratedAt,
	}

	// Convert trending categories
	for _, tc := range rawData.TrendingCategories {
		stats.TrendingCategories = append(stats.TrendingCategories, &domain.TrendingCategory{
			CategoryID:    tc.CategoryID,
			CategoryName:  tc.CategoryName,
			OrderCount30d: tc.OrderCount30d,
			OrderCount7d:  tc.OrderCount7d,
			GrowthRate:    tc.GrowthRate,
			TrendScore:    tc.TrendScore,
		})
	}

	// Convert hot listings
	for _, hl := range rawData.HotListings {
		stats.HotListings = append(stats.HotListings, &domain.HotListing{
			ListingID:       hl.ListingID,
			Title:           hl.Title,
			Orders24h:       hl.Orders24h,
			Orders7d:        hl.Orders7d,
			OrdersGrowth:    hl.OrdersGrowth,
			QuantitySold24h: hl.QuantitySold24h,
			Price:           hl.Price,
		})
	}

	// Convert popular searches
	for _, ps := range rawData.PopularSearches {
		stats.PopularSearches = append(stats.PopularSearches, &domain.PopularSearch{
			Query:       ps.Query,
			SearchCount: ps.SearchCount,
		})
	}

	r.logger.Debug().
		Int("trending_categories", len(stats.TrendingCategories)).
		Int("hot_listings", len(stats.HotListings)).
		Int("popular_searches", len(stats.PopularSearches)).
		Time("generated_at", stats.GeneratedAt).
		Msg("successfully fetched trending stats")

	return stats, nil
}
