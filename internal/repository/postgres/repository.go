// Package postgres implements PostgreSQL repository layer for listings microservice.
// It provides data access operations including CRUD, search, and indexing queue management.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver registration
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
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

// InitPgxPool initializes pgxpool connection pool for pgx-based repositories
func InitPgxPool(ctx context.Context, dsn string, maxConns, minConns int32, logger zerolog.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Configure connection pool
	config.MaxConns = maxConns
	config.MinConns = minConns

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	// Test connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().
		Int32("max_conns", maxConns).
		Int32("min_conns", minConns).
		Msg("PgxPool connection pool initialized")

	return pool, nil
}

// extractFieldTranslations extracts translations for a specific field from the translations map
func extractFieldTranslations(translations map[string]map[string]string, field string) map[string]string {
	result := make(map[string]string)
	if translations == nil {
		return result
	}

	for lang, fields := range translations {
		if value, ok := fields[field]; ok && value != "" {
			result[lang] = value
		}
	}
	return result
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

	// Generate unique slug from title
	slug, err := r.GenerateUniqueSlug(ctx, input.Title)
	if err != nil {
		r.logger.Error().Err(err).Str("title", input.Title).Msg("failed to generate slug")
		return nil, fmt.Errorf("failed to generate slug: %w", err)
	}

	// Prepare translation JSONB data
	titleTranslations := extractFieldTranslations(input.Translations, "title")
	descriptionTranslations := extractFieldTranslations(input.Translations, "description")
	locationTranslations := extractFieldTranslations(input.Translations, "location")
	cityTranslations := extractFieldTranslations(input.Translations, "city")
	countryTranslations := extractFieldTranslations(input.Translations, "country")

	// Marshal translations to JSON for PostgreSQL JSONB columns
	titleTranslationsJSON, err := json.Marshal(titleTranslations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal title translations: %w", err)
	}
	descriptionTranslationsJSON, err := json.Marshal(descriptionTranslations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal description translations: %w", err)
	}
	locationTranslationsJSON, err := json.Marshal(locationTranslations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal location translations: %w", err)
	}
	cityTranslationsJSON, err := json.Marshal(cityTranslations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal city translations: %w", err)
	}
	countryTranslationsJSON, err := json.Marshal(countryTranslations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal country translations: %w", err)
	}

	// Default original language to "sr" if not provided
	originalLanguage := input.OriginalLanguage
	if originalLanguage == "" {
		originalLanguage = "sr"
	}

	// Default show_on_map to true if not provided
	showOnMap := true
	if input.ShowOnMap != nil {
		showOnMap = *input.ShowOnMap
	}

	// Default location_privacy to "exact" if not provided
	locationPrivacy := "exact"
	if input.LocationPrivacy != nil {
		locationPrivacy = *input.LocationPrivacy
	}

	// Check if location is provided
	hasIndividualLocation := input.Location != nil

	query := `
		INSERT INTO listings (
			user_id, storefront_id, title, description, price, currency, category_id, quantity, sku, source_type, slug,
			title_translations, description_translations, location_translations, city_translations, country_translations, original_language,
			show_on_map, location_privacy, has_individual_location
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id, uuid, slug, user_id, storefront_id, title, description, price, currency, category_id,
		          status, visibility, quantity, sku, source_type, view_count, favorites_count,
		          title_translations, description_translations, location_translations, city_translations, country_translations, original_language,
		          show_on_map, location_privacy, has_individual_location,
		          created_at, updated_at, published_at, deleted_at, is_deleted
	`

	var listing domain.Listing
	var titleTranslationsBytes, descriptionTranslationsBytes, locationTranslationsBytes, cityTranslationsBytes, countryTranslationsBytes []byte

	err = r.db.QueryRowContext(
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
		slug,
		titleTranslationsJSON,
		descriptionTranslationsJSON,
		locationTranslationsJSON,
		cityTranslationsJSON,
		countryTranslationsJSON,
		originalLanguage,
		showOnMap,
		locationPrivacy,
		hasIndividualLocation,
	).Scan(
		&listing.ID,
		&listing.UUID,
		&listing.Slug,
		&listing.UserID,
		&listing.StorefrontID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Currency,
		&listing.CategoryID,
		&listing.Status,
		&listing.Visibility,
		&listing.Quantity,
		&listing.SKU,
		&listing.SourceType,
		&listing.ViewsCount,
		&listing.FavoritesCount,
		&titleTranslationsBytes,
		&descriptionTranslationsBytes,
		&locationTranslationsBytes,
		&cityTranslationsBytes,
		&countryTranslationsBytes,
		&listing.OriginalLanguage,
		&listing.ShowOnMap,
		&listing.LocationPrivacy,
		&listing.HasIndividualLocation,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&listing.PublishedAt,
		&listing.DeletedAt,
		&listing.IsDeleted,
	)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to create listing")
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	// Unmarshal JSONB translations back into maps
	if len(titleTranslationsBytes) > 0 {
		if err := json.Unmarshal(titleTranslationsBytes, &listing.TitleTranslations); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal title translations")
		}
	}
	if len(descriptionTranslationsBytes) > 0 {
		if err := json.Unmarshal(descriptionTranslationsBytes, &listing.DescriptionTranslations); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal description translations")
		}
	}
	if len(locationTranslationsBytes) > 0 {
		if err := json.Unmarshal(locationTranslationsBytes, &listing.LocationTranslations); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal location translations")
		}
	}
	if len(cityTranslationsBytes) > 0 {
		if err := json.Unmarshal(cityTranslationsBytes, &listing.CityTranslations); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal city translations")
		}
	}
	if len(countryTranslationsBytes) > 0 {
		if err := json.Unmarshal(countryTranslationsBytes, &listing.CountryTranslations); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal country translations")
		}
	}

	// DEBUG: Log what we have before creating location/attributes
	r.logger.Debug().
		Bool("has_location", input.Location != nil).
		Bool("has_condition", input.Condition != nil).
		Int("attributes_count", len(input.Attributes)).
		Int64("listing_id", listing.ID).
		Msg("DEBUG: about to create location and attributes")

	// Create location record if provided
	if input.Location != nil {
		logEvent := r.logger.Debug()
		if input.Location.Country != nil {
			logEvent = logEvent.Str("country", *input.Location.Country)
		}
		if input.Location.City != nil {
			logEvent = logEvent.Str("city", *input.Location.City)
		}
		if input.Location.Latitude != nil {
			logEvent = logEvent.Float64("lat", *input.Location.Latitude)
		}
		if input.Location.Longitude != nil {
			logEvent = logEvent.Float64("lng", *input.Location.Longitude)
		}
		logEvent.Msg("DEBUG: creating listing location")
		if err := r.createListingLocation(ctx, listing.ID, input.Location); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to create listing location")
			// Don't fail the entire operation, location is optional
		} else {
			r.logger.Debug().Int64("listing_id", listing.ID).Msg("DEBUG: listing location created successfully")
		}
	} else {
		r.logger.Debug().Int64("listing_id", listing.ID).Msg("DEBUG: input.Location is nil, skipping location creation")
	}

	// Create attributes (including condition) if provided
	if input.Condition != nil && *input.Condition != "" {
		r.logger.Debug().Str("condition", *input.Condition).Msg("DEBUG: creating condition attribute")
		// Add condition as an attribute
		if err := r.createListingAttribute(ctx, listing.ID, "condition", *input.Condition); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to create condition attribute")
		} else {
			r.logger.Debug().Int64("listing_id", listing.ID).Msg("DEBUG: condition attribute created successfully")
		}
	} else {
		r.logger.Debug().Int64("listing_id", listing.ID).Msg("DEBUG: input.Condition is nil or empty, skipping condition attribute")
	}

	// Create other attributes
	r.logger.Debug().Int("count", len(input.Attributes)).Msg("DEBUG: about to create attributes")
	for _, attr := range input.Attributes {
		r.logger.Debug().Str("key", attr.Key).Str("value", attr.Value).Msg("DEBUG: creating attribute")
		if err := r.createListingAttribute(ctx, listing.ID, attr.Key, attr.Value); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", listing.ID).Str("key", attr.Key).Msg("failed to create listing attribute")
		} else {
			r.logger.Debug().Int64("listing_id", listing.ID).Str("key", attr.Key).Msg("DEBUG: attribute created successfully")
		}
	}

	r.logger.Info().Int64("listing_id", listing.ID).Str("title", listing.Title).Str("slug", slug).Msg("listing created")
	return &listing, nil
}

// createListingLocation creates a location record for a listing
func (r *Repository) createListingLocation(ctx context.Context, listingID int64, loc *domain.CreateLocationInput) error {
	query := `
		INSERT INTO listing_locations (listing_id, country, city, postal_code, address_line1, address_line2, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query, listingID, loc.Country, loc.City, loc.PostalCode, loc.AddressLine1, loc.AddressLine2, loc.Latitude, loc.Longitude)
	return err
}

// createListingAttribute creates a single attribute record for a listing
func (r *Repository) createListingAttribute(ctx context.Context, listingID int64, key, value string) error {
	query := `
		INSERT INTO listing_attributes (listing_id, attribute_key, attribute_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (listing_id, attribute_key) DO UPDATE SET attribute_value = $3
	`
	_, err := r.db.ExecContext(ctx, query, listingID, key, value)
	return err
}

// getListingLocation retrieves location for a listing
func (r *Repository) getListingLocation(ctx context.Context, listingID int64) (*domain.ListingLocation, error) {
	query := `
		SELECT id, listing_id, country, city, postal_code, address_line1, address_line2, latitude, longitude, created_at, updated_at
		FROM listing_locations
		WHERE listing_id = $1
	`
	var loc domain.ListingLocation
	err := r.db.QueryRowContext(ctx, query, listingID).Scan(
		&loc.ID,
		&loc.ListingID,
		&loc.Country,
		&loc.City,
		&loc.PostalCode,
		&loc.AddressLine1,
		&loc.AddressLine2,
		&loc.Latitude,
		&loc.Longitude,
		&loc.CreatedAt,
		&loc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No location is not an error
		}
		return nil, fmt.Errorf("failed to get listing location: %w", err)
	}
	return &loc, nil
}

// getListingAttributes retrieves all attributes for a listing
func (r *Repository) getListingAttributes(ctx context.Context, listingID int64) ([]*domain.ListingAttribute, error) {
	query := `
		SELECT id, listing_id, attribute_key, attribute_value, created_at
		FROM listing_attributes
		WHERE listing_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*domain.ListingAttribute
	for rows.Next() {
		var attr domain.ListingAttribute
		if err := rows.Scan(&attr.ID, &attr.ListingID, &attr.AttributeKey, &attr.AttributeValue, &attr.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan listing attribute: %w", err)
		}
		attributes = append(attributes, &attr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listing attributes: %w", err)
	}
	return attributes, nil
}

// GetListingByID retrieves a listing by its ID
func (r *Repository) GetListingByID(ctx context.Context, id int64) (*domain.Listing, error) {
	query := `
		SELECT id, uuid, slug, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       title_translations, description_translations, location_translations, city_translations, country_translations, original_language,
		       show_on_map, location_privacy, has_individual_location,
		       expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
		FROM listings
		WHERE id = $1 AND is_deleted = false
	`

	var listing domain.Listing
	var titleTranslationsJSON, descriptionTranslationsJSON, locationTranslationsJSON, cityTranslationsJSON, countryTranslationsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.UUID,
		&listing.Slug,
		&listing.UserID,
		&listing.StorefrontID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Currency,
		&listing.CategoryID,
		&listing.Status,
		&listing.Visibility,
		&listing.Quantity,
		&listing.SKU,
		&listing.SourceType,
		&listing.ViewsCount,
		&listing.FavoritesCount,
		&titleTranslationsJSON,
		&descriptionTranslationsJSON,
		&locationTranslationsJSON,
		&cityTranslationsJSON,
		&countryTranslationsJSON,
		&listing.OriginalLanguage,
		&listing.ShowOnMap,
		&listing.LocationPrivacy,
		&listing.HasIndividualLocation,
		&listing.ExpiresAt,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&listing.PublishedAt,
		&listing.DeletedAt,
		&listing.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// Unmarshal JSONB translations into maps
	if len(titleTranslationsJSON) > 0 {
		if err := json.Unmarshal(titleTranslationsJSON, &listing.TitleTranslations); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to unmarshal title translations")
		}
	}
	if len(descriptionTranslationsJSON) > 0 {
		if err := json.Unmarshal(descriptionTranslationsJSON, &listing.DescriptionTranslations); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to unmarshal description translations")
		}
	}
	if len(locationTranslationsJSON) > 0 {
		if err := json.Unmarshal(locationTranslationsJSON, &listing.LocationTranslations); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to unmarshal location translations")
		}
	}
	if len(cityTranslationsJSON) > 0 {
		if err := json.Unmarshal(cityTranslationsJSON, &listing.CityTranslations); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to unmarshal city translations")
		}
	}
	if len(countryTranslationsJSON) > 0 {
		if err := json.Unmarshal(countryTranslationsJSON, &listing.CountryTranslations); err != nil {
			r.logger.Error().Err(err).Int64("listing_id", id).Msg("failed to unmarshal country translations")
		}
	}

	// Load images for indexing and API responses
	images, err := r.GetImages(ctx, id)
	if err != nil {
		r.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load images for listing")
		// Don't fail the whole request if images fail to load
	} else {
		listing.Images = images
	}

	// Load location
	location, err := r.getListingLocation(ctx, id)
	if err != nil {
		r.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load location for listing")
	} else {
		listing.Location = location
	}

	// Load attributes
	attributes, err := r.getListingAttributes(ctx, id)
	if err != nil {
		r.logger.Warn().Err(err).Int64("listing_id", id).Msg("failed to load attributes for listing")
	} else {
		listing.Attributes = attributes
	}

	return &listing, nil
}

// GetListingByUUID retrieves a listing by its UUID
func (r *Repository) GetListingByUUID(ctx context.Context, uuid string) (*domain.Listing, error) {
	query := `
		SELECT id, uuid, slug, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
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
		RETURNING id, uuid, slug, user_id, storefront_id, title, description, price, currency, category_id,
		          status, visibility, quantity, sku, source_type, view_count, favorites_count,
		          expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
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
		SELECT id, uuid, slug, user_id, storefront_id, title, description, price, currency, category_id,
		       status, visibility, quantity, sku, source_type, view_count, favorites_count,
		       expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
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
		SELECT id, uuid, slug, user_id, storefront_id, title, description, price,
		       currency, category_id, status, visibility,
		       quantity, sku, source_type, view_count, favorites_count,
		       expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
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
		SELECT id, uuid, slug, user_id, storefront_id, title, description, price,
		       currency, category_id, status, visibility,
		       quantity, sku, source_type, view_count, favorites_count,
		       expires_at, created_at, updated_at, published_at, deleted_at, is_deleted
		FROM listings
		WHERE status = 'active' AND is_deleted = false
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
