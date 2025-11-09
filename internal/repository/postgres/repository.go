// Package postgres implements PostgreSQL repository layer for listings microservice.
// It provides data access operations including CRUD, search, and indexing queue management.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver registration
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// Repository implements PostgreSQL data access
type Repository struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewRepository creates a new PostgreSQL repository
func NewRepository(db *sqlx.DB, logger zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger.With().Str("component", "postgres_repository").Logger(),
	}
}

// BeginTx starts a new database transaction
func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// GetDB returns the underlying database connection
func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

// InitDB initializes database connection with proper pool settings
func InitDB(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime, connMaxIdleTime time.Duration, logger zerolog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().
		Int("max_open_conns", maxOpenConns).
		Int("max_idle_conns", maxIdleConns).
		Dur("conn_max_lifetime", connMaxLifetime).
		Dur("conn_max_idle_time", connMaxIdleTime).
		Msg("PostgreSQL connection pool initialized")

	return db, nil
}

// validateCreateListingInput validates input before creating a listing
func validateCreateListingInput(input *domain.CreateListingInput) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	if input.UserID <= 0 {
		return fmt.Errorf("user_id must be greater than 0")
	}

	trimmedTitle := strings.TrimSpace(input.Title)
	if trimmedTitle == "" {
		return fmt.Errorf("title is required")
	}

	if len(trimmedTitle) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}

	if input.CategoryID <= 0 {
		return fmt.Errorf("category_id must be greater than 0")
	}

	if input.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if input.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if len(input.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters (ISO 4217)")
	}

	if input.SourceType != "c2c" && input.SourceType != "b2c" {
		return fmt.Errorf("source_type must be either 'c2c' or 'b2c'")
	}

	return nil
}

// CreateListing creates a new listing in the database
func (r *Repository) CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error) {
	// Validate input before database operation
	if err := validateCreateListingInput(input); err != nil {
		r.logger.Warn().Err(err).Msg("invalid create listing input")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	query := `
		INSERT INTO listings (user_id, storefront_id, title, description, price, currency, category_id, quantity, sku, source_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, uuid, user_id, storefront_id, title, description, price, currency, category_id,
		          status, visibility, quantity, sku, source_type, view_count, favorites_count,
		          created_at, updated_at, published_at, deleted_at, is_deleted
	`

	var listing domain.Listing
	err := r.db.QueryRowxContext(
		ctx,
		query,
		input.UserID,
		input.StorefrontID,
		input.Title,
		input.Description,
		input.Price,
		input.Currency,
		input.CategoryID,
		input.Quantity,
		input.SKU,
		input.SourceType,
	).StructScan(&listing)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to create listing")
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	r.logger.Info().Int64("listing_id", listing.ID).Str("title", listing.Title).Msg("listing created")
	return &listing, nil
}

// GetListingByID retrieves a listing by its ID
func (r *Repository) GetListingByID(ctx context.Context, id int64) (*domain.Listing, error) {
	query := `
		SELECT id, uuid, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       created_at, updated_at, published_at, deleted_at, is_deleted
		FROM listings
		WHERE id = $1 AND is_deleted = false
	`

	var listing domain.Listing
	err := r.db.GetContext(ctx, &listing, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return &listing, nil
}

// GetListingByUUID retrieves a listing by its UUID
func (r *Repository) GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error) {
	query := `
		SELECT id, uuid, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       created_at, updated_at, published_at, deleted_at, is_deleted
		FROM listings
		WHERE uuid = $1::uuid AND is_deleted = false
	`

	var listing domain.Listing
	err := r.db.GetContext(ctx, &listing, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		r.logger.Error().Err(err).Str("uuid", uuid).Msg("failed to get listing by UUID")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return &listing, nil
}

// GetListingBySlug retrieves a listing by its slug
func (r *Repository) GetListingBySlug(ctx context.Context, slug string) (*domain.Listing, error) {
	query := `
		SELECT id, uuid, slug, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
		FROM listings
		WHERE slug = $1 AND is_deleted = false
	`

	var listing domain.Listing
	err := r.db.GetContext(ctx, &listing, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		r.logger.Error().Err(err).Str("slug", slug).Msg("failed to get listing by slug")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return &listing, nil
}

// validateUpdateListingInput validates input before updating a listing
func validateUpdateListingInput(input *domain.UpdateListingInput) error {
	if input == nil {
		return fmt.Errorf("input cannot be nil")
	}

	if input.Title != nil {
		trimmedTitle := strings.TrimSpace(*input.Title)
		if trimmedTitle == "" {
			return fmt.Errorf("title cannot be empty")
		}
		if len(trimmedTitle) < 3 {
			return fmt.Errorf("title must be at least 3 characters")
		}
	}

	if input.Price != nil && *input.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if input.Quantity != nil && *input.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if input.Status != nil {
		validStatuses := map[string]bool{
			domain.StatusDraft:    true,
			domain.StatusActive:   true,
			domain.StatusInactive: true,
			domain.StatusSold:     true,
			domain.StatusArchived: true,
		}
		if !validStatuses[*input.Status] {
			return fmt.Errorf("invalid status: %s", *input.Status)
		}
	}

	return nil
}

// UpdateListing updates an existing listing
func (r *Repository) UpdateListing(ctx context.Context, id int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	// Validate input before database operation
	if err := validateUpdateListingInput(input); err != nil {
		r.logger.Warn().Err(err).Int64("listing_id", id).Msg("invalid update listing input")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if input.Title != nil {
		updates = append(updates, fmt.Sprintf("title = $%d", argPos))
		args = append(args, *input.Title)
		argPos++
	}

	if input.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *input.Description)
		argPos++
	}

	if input.Price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *input.Price)
		argPos++
	}

	if input.Quantity != nil {
		updates = append(updates, fmt.Sprintf("quantity = $%d", argPos))
		args = append(args, *input.Quantity)
		argPos++
	}

	if input.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *input.Status)
		argPos++

		// If status is being set to 'active', also set published_at
		if *input.Status == domain.StatusActive {
			updates = append(updates, "published_at = CURRENT_TIMESTAMP")
		}
	}

	if len(updates) == 0 {
		return r.GetListingByID(ctx, id)
	}

	// Add listing ID as last parameter
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE listings
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d AND is_deleted = false
		RETURNING id, uuid, user_id, storefront_id, title, description, price, currency, category_id,
		          status, visibility, quantity, sku, view_count, favorites_count,
		          created_at, updated_at, published_at, deleted_at, is_deleted
	`, strings.Join(updates, ", "), argPos)

	var listing domain.Listing
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&listing)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to update listing")
		return nil, fmt.Errorf("failed to update listing: %w", err)
	}

	r.logger.Info().Int64("listing_id", listing.ID).Msg("listing updated")
	return &listing, nil
}

// DeleteListing soft-deletes a listing
func (r *Repository) DeleteListing(ctx context.Context, id int64) error {
	query := `
		UPDATE listings
		SET is_deleted = true, deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_deleted = false
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to delete listing")
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found or already deleted")
	}

	r.logger.Info().Int64("listing_id", id).Msg("listing deleted")
	return nil
}

// ListListings returns a filtered list of listings
func (r *Repository) ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error) {
	// Build WHERE clause dynamically
	whereConditions := []string{"is_deleted = false"}
	args := []interface{}{}
	argPos := 1

	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argPos))
		args = append(args, *filter.UserID)
		argPos++
	}

	if filter.StorefrontID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("storefront_id = $%d", argPos))
		args = append(args, *filter.StorefrontID)
		argPos++
	}

	if filter.CategoryID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, *filter.CategoryID)
		argPos++
	}

	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.SourceType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("source_type = $%d", argPos))
		args = append(args, *filter.SourceType)
		argPos++
	}

	if filter.MinPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("price >= $%d", argPos))
		args = append(args, *filter.MinPrice)
		argPos++
	}

	if filter.MaxPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("price <= $%d", argPos))
		args = append(args, *filter.MaxPrice)
		argPos++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listings WHERE %s", whereClause)
	var total int32
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to count listings")
		return nil, 0, fmt.Errorf("failed to count listings: %w", err)
	}

	// Get listings with pagination
	args = append(args, filter.Limit, filter.Offset)
	query := fmt.Sprintf(`
		SELECT id, uuid, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       created_at, updated_at, published_at, deleted_at, is_deleted
		FROM listings
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	var listings []*domain.Listing
	err = r.db.SelectContext(ctx, &listings, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to list listings")
		return nil, 0, fmt.Errorf("failed to list listings: %w", err)
	}

	return listings, total, nil
}

// SearchListings performs full-text search using PostgreSQL's built-in search
func (r *Repository) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error) {
	whereConditions := []string{"status = 'active'"}
	args := []interface{}{}
	argPos := 1

	// Add full-text search condition
	if query.Query != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(
			"to_tsvector('english', title || ' ' || COALESCE(description, '')) @@ plainto_tsquery('english', $%d)",
			argPos,
		))
		args = append(args, query.Query)
		argPos++
	}

	if query.CategoryID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, *query.CategoryID)
		argPos++
	}

	if query.MinPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("price >= $%d", argPos))
		args = append(args, *query.MinPrice)
		argPos++
	}

	if query.MaxPrice != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("price <= $%d", argPos))
		args = append(args, *query.MaxPrice)
		argPos++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM listings WHERE %s", whereClause)
	var total int32
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to count search results")
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Get listings with pagination
	args = append(args, query.Limit, query.Offset)
	searchQuery := fmt.Sprintf(`
		SELECT id, gen_random_uuid() as uuid, user_id, storefront_id, title, description, price,
		       'RSD' as currency, category_id, status, 'public' as visibility,
		       1 as quantity, NULL as sku, 'c2c' as source_type, view_count, 0 as favorites_count,
		       created_at, updated_at, NULL as published_at, NULL as deleted_at, false as is_deleted
		FROM listings
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	var listings []*domain.Listing
	err = r.db.SelectContext(ctx, &listings, searchQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to search listings")
		return nil, 0, fmt.Errorf("failed to search listings: %w", err)
	}

	return listings, total, nil
}

// EnqueueIndexing adds a listing to the indexing queue
func (r *Repository) EnqueueIndexing(ctx context.Context, listingID int64, operation string) error {
	query := `
		INSERT INTO indexing_queue (listing_id, operation, status, retry_count, max_retries)
		VALUES ($1, $2, 'pending', 0, 3)
		ON CONFLICT (listing_id) WHERE status = 'pending'
		DO UPDATE SET operation = EXCLUDED.operation, updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, query, listingID, operation)
	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", listingID).Str("operation", operation).Msg("failed to enqueue indexing")
		return fmt.Errorf("failed to enqueue indexing: %w", err)
	}

	r.logger.Debug().Int64("listing_id", listingID).Str("operation", operation).Msg("indexing enqueued")
	return nil
}

// GetPendingIndexingJobs retrieves pending indexing jobs from the queue
func (r *Repository) GetPendingIndexingJobs(ctx context.Context, limit int) ([]*domain.IndexingQueueItem, error) {
	query := `
		UPDATE indexing_queue
		SET status = 'processing', updated_at = CURRENT_TIMESTAMP
		WHERE id IN (
			SELECT id FROM indexing_queue
			WHERE status = 'pending' AND retry_count < max_retries
			ORDER BY created_at ASC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, listing_id, operation, status, retry_count, max_retries,
		          error_message, created_at, updated_at, processed_at
	`

	var jobs []*domain.IndexingQueueItem
	err := r.db.SelectContext(ctx, &jobs, query, limit)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get pending indexing jobs")
		return nil, fmt.Errorf("failed to get pending indexing jobs: %w", err)
	}

	return jobs, nil
}

// CompleteIndexingJob marks an indexing job as completed
func (r *Repository) CompleteIndexingJob(ctx context.Context, jobID int64) error {
	query := `
		UPDATE indexing_queue
		SET status = 'completed', processed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, jobID)
	if err != nil {
		r.logger.Error().Err(err).Int64("job_id", jobID).Msg("failed to complete indexing job")
		return fmt.Errorf("failed to complete indexing job: %w", err)
	}

	return nil
}

// FailIndexingJob marks an indexing job as failed with error message
func (r *Repository) FailIndexingJob(ctx context.Context, jobID int64, errorMsg string) error {
	query := `
		UPDATE indexing_queue
		SET status = 'failed', retry_count = retry_count + 1, error_message = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, jobID, errorMsg)
	if err != nil {
		r.logger.Error().Err(err).Int64("job_id", jobID).Msg("failed to mark indexing job as failed")
		return fmt.Errorf("failed to fail indexing job: %w", err)
	}

	return nil
}

// HealthCheck performs a database health check
func (r *Repository) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var result int
	err := r.db.GetContext(ctx, &result, "SELECT 1")
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetConnectionStats returns connection pool statistics
func (r *Repository) GetConnectionStats() sql.DBStats {
	return r.db.Stats()
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

// WithTransaction executes a function within a database transaction (accepts sqlx.Tx)
func (r *Repository) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				r.logger.Error().Err(rbErr).Msg("failed to rollback transaction on panic")
			}
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			r.logger.Error().Err(rbErr).Msg("failed to rollback transaction")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetListingsForReindex retrieves listings that need reindexing
func (r *Repository) GetListingsForReindex(ctx context.Context, limit int) ([]*domain.Listing, error) {
	query := `
		SELECT id, gen_random_uuid() as uuid, user_id, storefront_id, title, description, price,
		       'RSD' as currency, category_id, status, 'public' as visibility,
		       1 as quantity, NULL as sku, view_count, 0 as favorites_count,
		       created_at, updated_at, NULL as published_at, NULL as deleted_at, false as is_deleted
		FROM listings
		WHERE status = 'active'
		ORDER BY id ASC
		LIMIT $1
	`

	var listings []*domain.Listing
	err := r.db.SelectContext(ctx, &listings, query, limit)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get listings for reindex")
		return nil, fmt.Errorf("failed to get listings for reindex: %w", err)
	}

	return listings, nil
}

// ResetReindexFlags resets reindex flags for specified listings (stub implementation)
func (r *Repository) ResetReindexFlags(ctx context.Context, listingIDs []int64) error {
	if len(listingIDs) == 0 {
		return nil
	}

	// Currently listings table doesn't have needs_reindex flag
	// This is a placeholder for future implementation when we add the flag
	r.logger.Info().Int("count", len(listingIDs)).Msg("reset reindex flags (noop - flag not implemented)")
	return nil
}

// SyncDiscounts synchronizes discount information across listings (stub implementation)
func (r *Repository) SyncDiscounts(ctx context.Context) error {
	// Placeholder for discount sync logic
	// This will be implemented when discount system is added
	r.logger.Info().Msg("sync discounts (noop - discounts not implemented)")
	return nil
}
