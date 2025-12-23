package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository"
)

// searchQueriesRepository implements repository.SearchQueriesRepository
type searchQueriesRepository struct {
	db     *pgxpool.Pool
	logger zerolog.Logger
}

// NewSearchQueriesRepository creates a new search queries repository
func NewSearchQueriesRepository(db *pgxpool.Pool, logger zerolog.Logger) repository.SearchQueriesRepository {
	return &searchQueriesRepository{
		db:     db,
		logger: logger.With().Str("repository", "search_queries").Logger(),
	}
}

// ============================================================================
// CREATE SEARCH QUERY (Async Logging)
// ============================================================================

// CreateSearchQuery logs a search query for analytics
func (r *searchQueriesRepository) CreateSearchQuery(
	ctx context.Context,
	input *domain.CreateSearchQueryInput,
) (*domain.SearchQuery, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	query := `
		INSERT INTO search_queries (
			query_text,
			category_id,
			user_id,
			session_id,
			results_count,
			clicked_listing_id
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, query_text, category_id, user_id, session_id,
		          results_count, clicked_listing_id, created_at
	`

	var searchQuery domain.SearchQuery
	err := r.db.QueryRow(
		ctx,
		query,
		input.QueryText,
		input.CategoryID,
		input.UserID,
		input.SessionID,
		input.ResultsCount,
		input.ClickedListingID,
	).Scan(
		&searchQuery.ID,
		&searchQuery.QueryText,
		&searchQuery.CategoryID,
		&searchQuery.UserID,
		&searchQuery.SessionID,
		&searchQuery.ResultsCount,
		&searchQuery.ClickedListingID,
		&searchQuery.CreatedAt,
	)

	if err != nil {
		r.logger.Error().
			Err(err).
			Str("query_text", input.QueryText).
			Msg("failed to create search query")
		return nil, fmt.Errorf("failed to create search query: %w", err)
	}

	r.logger.Debug().
		Int64("id", searchQuery.ID).
		Str("query_text", searchQuery.QueryText).
		Interface("category_id", searchQuery.CategoryID).
		Int32("results_count", searchQuery.ResultsCount).
		Msg("search query logged")

	return &searchQuery, nil
}

// ============================================================================
// TRENDING QUERIES (Time-based Aggregation)
// ============================================================================

// GetTrendingQueries returns the most popular search queries within a time range
func (r *searchQueriesRepository) GetTrendingQueries(
	ctx context.Context,
	filter *domain.GetTrendingQueriesFilter,
) ([]domain.TrendingSearch, error) {
	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Build query with dynamic WHERE clause
	query := `
		SELECT
			query_text,
			COUNT(*) as search_count,
			MAX(created_at) as last_searched,
			category_id
		FROM search_queries
		WHERE created_at > NOW() - INTERVAL '%d days'
	`
	query = fmt.Sprintf(query, filter.DaysAgo)

	// Add category filter if provided
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = '%s'", *filter.CategoryID)
	}

	// Add results filter
	if !filter.IncludeZeroResults {
		query += fmt.Sprintf(" AND results_count >= %d", filter.MinResultsCount)
	}

	// Group by and order by
	query += `
		GROUP BY query_text, category_id
		ORDER BY search_count DESC, last_searched DESC
		LIMIT $1
	`

	r.logger.Debug().
		Int32("limit", filter.Limit).
		Int32("days_ago", filter.DaysAgo).
		Interface("category_id", filter.CategoryID).
		Msg("fetching trending queries")

	rows, err := r.db.Query(ctx, query, filter.Limit)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("failed to fetch trending queries")
		return nil, fmt.Errorf("failed to fetch trending queries: %w", err)
	}
	defer rows.Close()

	var results []domain.TrendingSearch
	for rows.Next() {
		var trending domain.TrendingSearch
		err := rows.Scan(
			&trending.QueryText,
			&trending.SearchCount,
			&trending.LastSearched,
			&trending.CategoryID,
		)
		if err != nil {
			r.logger.Error().
				Err(err).
				Msg("failed to scan trending query row")
			return nil, fmt.Errorf("failed to scan trending query: %w", err)
		}
		results = append(results, trending)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating trending queries: %w", err)
	}

	r.logger.Info().
		Int("count", len(results)).
		Int32("days_ago", filter.DaysAgo).
		Msg("trending queries fetched")

	return results, nil
}

// ============================================================================
// USER HISTORY (Personal Search History)
// ============================================================================

// GetUserHistory returns search history for a specific user or session
func (r *searchQueriesRepository) GetUserHistory(
	ctx context.Context,
	filter *domain.GetUserHistoryFilter,
) ([]domain.SearchQuery, error) {
	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Build query with dynamic WHERE clause
	query := `
		SELECT
			id, query_text, category_id, user_id, session_id,
			results_count, clicked_listing_id, created_at
		FROM search_queries
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	// Add user_id or session_id filter
	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	} else if filter.SessionID != nil {
		query += fmt.Sprintf(" AND session_id = $%d", argIndex)
		args = append(args, *filter.SessionID)
		argIndex++
	}

	// Add category filter if provided
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	// Order by most recent first
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argIndex)
	args = append(args, filter.Limit)

	r.logger.Debug().
		Interface("user_id", filter.UserID).
		Interface("session_id", filter.SessionID).
		Int32("limit", filter.Limit).
		Msg("fetching user search history")

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("failed to fetch user history")
		return nil, fmt.Errorf("failed to fetch user history: %w", err)
	}
	defer rows.Close()

	var results []domain.SearchQuery
	for rows.Next() {
		var sq domain.SearchQuery
		err := rows.Scan(
			&sq.ID,
			&sq.QueryText,
			&sq.CategoryID,
			&sq.UserID,
			&sq.SessionID,
			&sq.ResultsCount,
			&sq.ClickedListingID,
			&sq.CreatedAt,
		)
		if err != nil {
			r.logger.Error().
				Err(err).
				Msg("failed to scan user history row")
			return nil, fmt.Errorf("failed to scan user history: %w", err)
		}
		results = append(results, sq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user history: %w", err)
	}

	r.logger.Info().
		Int("count", len(results)).
		Msg("user history fetched")

	return results, nil
}

// ============================================================================
// POPULAR QUERIES (All-Time)
// ============================================================================

// GetPopularQueries returns all-time popular searches (no time filter)
func (r *searchQueriesRepository) GetPopularQueries(
	ctx context.Context,
	filter *domain.GetPopularQueriesFilter,
) ([]domain.TrendingSearch, error) {
	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Build query with dynamic WHERE clause
	query := `
		SELECT
			query_text,
			COUNT(*) as search_count,
			MAX(created_at) as last_searched,
			category_id
		FROM search_queries
		WHERE results_count > 0
	`

	// Add category filter if provided
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = '%s'", *filter.CategoryID)
	}

	// Add min search count filter
	if filter.MinSearchCount > 0 {
		query += " GROUP BY query_text, category_id"
		query += fmt.Sprintf(" HAVING COUNT(*) >= %d", filter.MinSearchCount)
		query += " ORDER BY search_count DESC, last_searched DESC"
	} else {
		query += " GROUP BY query_text, category_id"
		query += " ORDER BY search_count DESC, last_searched DESC"
	}

	query += " LIMIT $1"

	r.logger.Debug().
		Int32("limit", filter.Limit).
		Interface("category_id", filter.CategoryID).
		Int64("min_search_count", filter.MinSearchCount).
		Msg("fetching popular queries")

	rows, err := r.db.Query(ctx, query, filter.Limit)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("failed to fetch popular queries")
		return nil, fmt.Errorf("failed to fetch popular queries: %w", err)
	}
	defer rows.Close()

	var results []domain.TrendingSearch
	for rows.Next() {
		var trending domain.TrendingSearch
		err := rows.Scan(
			&trending.QueryText,
			&trending.SearchCount,
			&trending.LastSearched,
			&trending.CategoryID,
		)
		if err != nil {
			r.logger.Error().
				Err(err).
				Msg("failed to scan popular query row")
			return nil, fmt.Errorf("failed to scan popular query: %w", err)
		}
		results = append(results, trending)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating popular queries: %w", err)
	}

	r.logger.Info().
		Int("count", len(results)).
		Msg("popular queries fetched")

	return results, nil
}

// ============================================================================
// UPDATE CLICKED LISTING (CTR Tracking)
// ============================================================================

// UpdateClickedListing records that a user clicked on a listing from search results
func (r *searchQueriesRepository) UpdateClickedListing(
	ctx context.Context,
	searchQueryID int64,
	listingID int64,
) error {
	query := `
		UPDATE search_queries
		SET clicked_listing_id = $1
		WHERE id = $2
		  AND clicked_listing_id IS NULL
	`

	result, err := r.db.Exec(ctx, query, listingID, searchQueryID)
	if err != nil {
		r.logger.Error().
			Err(err).
			Int64("search_query_id", searchQueryID).
			Int64("listing_id", listingID).
			Msg("failed to update clicked listing")
		return fmt.Errorf("failed to update clicked listing: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		r.logger.Warn().
			Int64("search_query_id", searchQueryID).
			Int64("listing_id", listingID).
			Msg("search query not found or already has clicked_listing_id")
		return fmt.Errorf("search query not found or already updated")
	}

	r.logger.Debug().
		Int64("search_query_id", searchQueryID).
		Int64("listing_id", listingID).
		Msg("clicked listing updated")

	return nil
}

// ============================================================================
// CTR ANALYSIS
// ============================================================================

// GetCTRAnalysis returns click-through rate statistics for queries
func (r *searchQueriesRepository) GetCTRAnalysis(
	ctx context.Context,
	filter *domain.GetCTRAnalysisFilter,
) ([]domain.SearchQueryCTR, error) {
	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	// Build query
	query := `
		SELECT
			query_text,
			COUNT(*) as total_searches,
			COUNT(clicked_listing_id) as total_clicks,
			ROUND(100.0 * COUNT(clicked_listing_id) / NULLIF(COUNT(*), 0), 2) as ctr_percent,
			category_id
		FROM search_queries
		WHERE created_at > NOW() - INTERVAL '%d days'
	`
	query = fmt.Sprintf(query, filter.DaysAgo)

	// Add query_text filter if provided
	if filter.QueryText != nil && *filter.QueryText != "" {
		query += fmt.Sprintf(" AND query_text = '%s'", pgx.Identifier{*filter.QueryText}.Sanitize())
	}

	// Add category filter if provided
	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = '%s'", *filter.CategoryID)
	}

	// Group and order
	query += `
		GROUP BY query_text, category_id
		ORDER BY total_searches DESC, ctr_percent DESC
		LIMIT $1
	`

	r.logger.Debug().
		Interface("query_text", filter.QueryText).
		Interface("category_id", filter.CategoryID).
		Int32("days_ago", filter.DaysAgo).
		Int32("limit", filter.Limit).
		Msg("fetching CTR analysis")

	rows, err := r.db.Query(ctx, query, filter.Limit)
	if err != nil {
		r.logger.Error().
			Err(err).
			Msg("failed to fetch CTR analysis")
		return nil, fmt.Errorf("failed to fetch CTR analysis: %w", err)
	}
	defer rows.Close()

	var results []domain.SearchQueryCTR
	for rows.Next() {
		var ctr domain.SearchQueryCTR
		err := rows.Scan(
			&ctr.QueryText,
			&ctr.TotalSearches,
			&ctr.TotalClicks,
			&ctr.CTRPercent,
			&ctr.CategoryID,
		)
		if err != nil {
			r.logger.Error().
				Err(err).
				Msg("failed to scan CTR analysis row")
			return nil, fmt.Errorf("failed to scan CTR analysis: %w", err)
		}
		results = append(results, ctr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating CTR analysis: %w", err)
	}

	r.logger.Info().
		Int("count", len(results)).
		Msg("CTR analysis fetched")

	return results, nil
}

// ============================================================================
// CLEANUP OLD QUERIES (Retention Policy)
// ============================================================================

// CleanupOldQueries deletes search queries older than retention period
func (r *searchQueriesRepository) CleanupOldQueries(
	ctx context.Context,
	daysToKeep int32,
) (int64, error) {
	if daysToKeep < 1 {
		return 0, fmt.Errorf("daysToKeep must be at least 1")
	}

	query := `
		DELETE FROM search_queries
		WHERE created_at < NOW() - INTERVAL '%d days'
	`
	query = fmt.Sprintf(query, daysToKeep)

	r.logger.Info().
		Int32("days_to_keep", daysToKeep).
		Msg("cleaning up old search queries")

	start := time.Now()
	result, err := r.db.Exec(ctx, query)
	if err != nil {
		r.logger.Error().
			Err(err).
			Int32("days_to_keep", daysToKeep).
			Msg("failed to cleanup old queries")
		return 0, fmt.Errorf("failed to cleanup old queries: %w", err)
	}

	rowsDeleted := result.RowsAffected()

	r.logger.Info().
		Int64("rows_deleted", rowsDeleted).
		Dur("duration", time.Since(start)).
		Int32("days_to_keep", daysToKeep).
		Msg("old search queries cleaned up")

	return rowsDeleted, nil
}
