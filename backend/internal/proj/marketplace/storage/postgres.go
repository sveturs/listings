// TEMPORARY: Will be moved to microservice
package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/domain/models"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type postgresMarketplaceStorage struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

func NewPostgresMarketplaceStorage(db *sqlx.DB, logger zerolog.Logger) MarketplaceStorage {
	return &postgresMarketplaceStorage{
		db:     db,
		logger: logger.With().Str("storage", "marketplace_postgres").Logger(),
	}
}

// GetCategories получает список всех активных категорий
func (s *postgresMarketplaceStorage) GetCategories(ctx context.Context, lang string) ([]models.MarketplaceCategory, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, has_custom_ui, custom_ui_component,
			sort_order, level, COALESCE(count, 0) as count
		FROM c2c_categories
		WHERE is_active = true
		ORDER BY sort_order ASC, name ASC
	`

	var categories []models.MarketplaceCategory
	if err := s.db.SelectContext(ctx, &categories, query); err != nil {
		s.logger.Error().Err(err).Msg("Failed to get categories")
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetPopularCategories получает популярные категории (с наибольшим count)
func (s *postgresMarketplaceStorage) GetPopularCategories(ctx context.Context, lang string, limit int) ([]models.MarketplaceCategory, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, has_custom_ui, custom_ui_component,
			sort_order, level, COALESCE(count, 0) as count
		FROM c2c_categories
		WHERE is_active = true AND parent_id IS NULL
		ORDER BY count DESC, sort_order ASC
		LIMIT $1
	`

	var categories []models.MarketplaceCategory
	if err := s.db.SelectContext(ctx, &categories, query, limit); err != nil {
		s.logger.Error().Err(err).Int("limit", limit).Msg("Failed to get popular categories")
		return nil, fmt.Errorf("failed to get popular categories: %w", err)
	}

	return categories, nil
}

// GetCategoryByID получает категорию по ID
func (s *postgresMarketplaceStorage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, has_custom_ui, custom_ui_component,
			sort_order, level, COALESCE(count, 0) as count
		FROM c2c_categories
		WHERE id = $1 AND is_active = true
	`

	var category models.MarketplaceCategory
	if err := s.db.GetContext(ctx, &category, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %d", id)
		}
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to get category by ID")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// GetCategoryBySlug получает категорию по slug
func (s *postgresMarketplaceStorage) GetCategoryBySlug(ctx context.Context, slug string) (*models.MarketplaceCategory, error) {
	query := `
		SELECT
			id, name, slug, parent_id, icon, description,
			is_active, created_at, has_custom_ui, custom_ui_component,
			sort_order, level, COALESCE(count, 0) as count
		FROM c2c_categories
		WHERE slug = $1 AND is_active = true
	`

	var category models.MarketplaceCategory
	if err := s.db.GetContext(ctx, &category, query, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("category not found: %s", slug)
		}
		s.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get category by slug")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// GetUserFavorites получает список ID избранных объявлений пользователя
func (s *postgresMarketplaceStorage) GetUserFavorites(ctx context.Context, userID int) ([]int, error) {
	query := `
		SELECT listing_id
		FROM c2c_favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var listingIDs []int
	if err := s.db.SelectContext(ctx, &listingIDs, query, userID); err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get user favorites")
		return nil, fmt.Errorf("failed to get user favorites: %w", err)
	}

	return listingIDs, nil
}

// AddToFavorites добавляет объявление в избранное
func (s *postgresMarketplaceStorage) AddToFavorites(ctx context.Context, userID, listingID int) error {
	query := `
		INSERT INTO c2c_favorites (user_id, listing_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, listing_id) DO NOTHING
	`

	if _, err := s.db.ExecContext(ctx, query, userID, listingID); err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Int("listing_id", listingID).Msg("Failed to add to favorites")
		return fmt.Errorf("failed to add to favorites: %w", err)
	}

	return nil
}

// RemoveFromFavorites удаляет объявление из избранного
func (s *postgresMarketplaceStorage) RemoveFromFavorites(ctx context.Context, userID, listingID int) error {
	query := `
		DELETE FROM c2c_favorites
		WHERE user_id = $1 AND listing_id = $2
	`

	if _, err := s.db.ExecContext(ctx, query, userID, listingID); err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Int("listing_id", listingID).Msg("Failed to remove from favorites")
		return fmt.Errorf("failed to remove from favorites: %w", err)
	}

	return nil
}

// IsFavorite проверяет, находится ли объявление в избранном
func (s *postgresMarketplaceStorage) IsFavorite(ctx context.Context, userID, listingID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM c2c_favorites
			WHERE user_id = $1 AND listing_id = $2
		)
	`

	var exists bool
	if err := s.db.GetContext(ctx, &exists, query, userID, listingID); err != nil {
		s.logger.Error().Err(err).Int("user_id", userID).Int("listing_id", listingID).Msg("Failed to check favorite")
		return false, fmt.Errorf("failed to check favorite: %w", err)
	}

	return exists, nil
}

// GetCategoryAttributes получает атрибуты категории
func (s *postgresMarketplaceStorage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	query := `
		SELECT
			ua.id,
			ua.name,
			ua.display_name,
			ua.attribute_type,
			ua.icon,
			ua.options,
			ua.validation_rules,
			ua.is_searchable,
			ua.is_filterable,
			uca.is_required,
			ua.show_in_card,
			false as show_in_list,
			uca.sort_order,
			ua.created_at,
			ua.is_variant_compatible,
			ua.affects_stock
		FROM unified_category_attributes uca
		JOIN unified_attributes ua ON uca.attribute_id = ua.id
		WHERE uca.category_id = $1 AND uca.is_enabled = true
		ORDER BY uca.sort_order ASC, ua.name ASC
	`

	var attributes []models.CategoryAttribute
	if err := s.db.SelectContext(ctx, &attributes, query, categoryID); err != nil {
		s.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to get category attributes")
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}

	return attributes, nil
}

// GetVariantAttributes получает вариативные атрибуты категории
func (s *postgresMarketplaceStorage) GetVariantAttributes(ctx context.Context, categorySlug string) ([]models.CategoryVariantAttribute, error) {
	query := `
		SELECT
			cva.id,
			cva.category_id,
			cva.attribute_id,
			cva.attribute_name,
			cva.is_required,
			cva.is_filterable,
			cva.sort_order
		FROM category_variant_attributes cva
		JOIN c2c_categories c ON cva.category_id = c.id
		WHERE c.slug = $1 AND c.is_active = true
		ORDER BY cva.sort_order ASC, cva.attribute_name ASC
	`

	var attributes []models.CategoryVariantAttribute
	if err := s.db.SelectContext(ctx, &attributes, query, categorySlug); err != nil {
		s.logger.Error().Err(err).Str("category_slug", categorySlug).Msg("Failed to get variant attributes")
		return nil, fmt.Errorf("failed to get variant attributes: %w", err)
	}

	return attributes, nil
}

// GetStorefronts получает список витрин с фильтрацией
func (s *postgresMarketplaceStorage) GetStorefronts(ctx context.Context, filters StorefrontFilters) ([]models.Storefront, int, error) {
	// Базовый запрос
	baseQuery := `
		FROM b2c_stores
		WHERE 1=1
	`

	// Добавляем фильтр is_active
	if filters.IsActive != nil {
		baseQuery += fmt.Sprintf(" AND is_active = %t", *filters.IsActive)
	}

	// Подсчитываем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := s.db.GetContext(ctx, &total, countQuery); err != nil {
		s.logger.Error().Err(err).Msg("Failed to count storefronts")
		return nil, 0, fmt.Errorf("failed to count storefronts: %w", err)
	}

	// Если результатов нет, возвращаем пустой список
	if total == 0 {
		return []models.Storefront{}, 0, nil
	}

	// Запрос данных с сортировкой и пагинацией
	selectQuery := `
		SELECT
			id, user_id, slug, name, description, logo_url, banner_url,
			theme, phone, email, website, address, city, postal_code, country,
			latitude, longitude, settings, seo_meta, is_active, is_verified,
			verification_date, rating, reviews_count, products_count, sales_count,
			views_count, subscription_plan, subscription_expires_at, commission_rate,
			ai_agent_enabled, ai_agent_config, live_shopping_enabled, group_buying_enabled,
			created_at, updated_at, formatted_address, geo_strategy, default_privacy_level,
			address_verified, subscription_id, is_subscription_active, followers_count
	` + baseQuery

	// Добавляем сортировку
	sortColumn := "products_count"
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "products_count", "rating", "created_at", "views_count":
			sortColumn = filters.SortBy
		}
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	selectQuery += fmt.Sprintf(" ORDER BY %s %s, id DESC", sortColumn, sortOrder)

	// Добавляем пагинацию
	if filters.Limit > 0 {
		offset := 0
		if filters.Page > 0 {
			offset = (filters.Page - 1) * filters.Limit
		}
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filters.Limit, offset)
	}

	var storefronts []models.Storefront
	if err := s.db.SelectContext(ctx, &storefronts, selectQuery); err != nil {
		s.logger.Error().Err(err).Msg("Failed to get storefronts")
		return nil, 0, fmt.Errorf("failed to get storefronts: %w", err)
	}

	return storefronts, total, nil
}

// GetStorefrontByID получает витрину по ID
func (s *postgresMarketplaceStorage) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	query := `
		SELECT
			id, user_id, slug, name, description, logo_url, banner_url,
			theme, phone, email, website, address, city, postal_code, country,
			latitude, longitude, settings, seo_meta, is_active, is_verified,
			verification_date, rating, reviews_count, products_count, sales_count,
			views_count, subscription_plan, subscription_expires_at, commission_rate,
			ai_agent_enabled, ai_agent_config, live_shopping_enabled, group_buying_enabled,
			created_at, updated_at, formatted_address, geo_strategy, default_privacy_level,
			address_verified, subscription_id, is_subscription_active, followers_count
		FROM b2c_stores
		WHERE id = $1
	`

	var storefront models.Storefront
	if err := s.db.GetContext(ctx, &storefront, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("storefront not found: %d", id)
		}
		s.logger.Error().Err(err).Int("id", id).Msg("Failed to get storefront by ID")
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}

	return &storefront, nil
}

// CreateListing создает новое объявление
// TEMPORARY: Direct DB insert until microservice fully integrated
func (s *postgresMarketplaceStorage) CreateListing(
	ctx context.Context,
	userID int,
	categoryID int,
	title string,
	description *string,
	price float64,
	currency string,
	quantity int32,
	sku *string,
	storefrontID *int,
) (*models.MarketplaceListing, error) {
	query := `
		INSERT INTO c2c_listings (
			user_id, category_id, title, description, price,
			status, created_at, updated_at, storefront_id
		)
		VALUES ($1, $2, $3, $4, $5, 'draft', NOW(), NOW(), $6)
		RETURNING id, user_id, category_id, title, description, price,
				  status, created_at, updated_at, storefront_id, views_count
	`

	var listing models.MarketplaceListing
	var dbDescription sql.NullString
	var dbStorefrontID sql.NullInt32

	err := s.db.QueryRowContext(ctx, query,
		userID, categoryID, title, description, price, storefrontID,
	).Scan(
		&listing.ID,
		&listing.UserID,
		&listing.CategoryID,
		&listing.Title,
		&dbDescription,
		&listing.Price,
		&listing.Status,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&dbStorefrontID,
		&listing.ViewsCount,
	)
	if err != nil {
		s.logger.Error().Err(err).
			Int("user_id", userID).
			Int("category_id", categoryID).
			Msg("Failed to create listing")
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	// Convert nullable fields
	if dbDescription.Valid {
		listing.Description = dbDescription.String
	}
	if dbStorefrontID.Valid {
		storefrontIDValue := int(dbStorefrontID.Int32)
		listing.StorefrontID = &storefrontIDValue
	}

	s.logger.Info().
		Int("listing_id", listing.ID).
		Int("user_id", userID).
		Str("title", title).
		Msg("Listing created successfully")

	return &listing, nil
}

// GetListing получает объявление по ID
// TEMPORARY: Direct DB query until microservice fully integrated
func (s *postgresMarketplaceStorage) GetListing(ctx context.Context, listingID int) (*models.MarketplaceListing, error) {
	query := `
		SELECT
			id, user_id, category_id, title, description, price,
			status, created_at, updated_at, storefront_id, views_count,
			condition
		FROM c2c_listings
		WHERE id = $1
	`

	var listing models.MarketplaceListing
	var dbDescription sql.NullString
	var dbStorefrontID sql.NullInt32
	var dbCondition sql.NullString

	err := s.db.QueryRowContext(ctx, query, listingID).Scan(
		&listing.ID,
		&listing.UserID,
		&listing.CategoryID,
		&listing.Title,
		&dbDescription,
		&listing.Price,
		&listing.Status,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&dbStorefrontID,
		&listing.ViewsCount,
		&dbCondition,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("listing not found: %d", listingID)
		}
		s.logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to get listing")
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	// Convert nullable fields
	if dbDescription.Valid {
		listing.Description = dbDescription.String
	}
	if dbStorefrontID.Valid {
		storefrontIDValue := int(dbStorefrontID.Int32)
		listing.StorefrontID = &storefrontIDValue
	}
	if dbCondition.Valid {
		listing.Condition = dbCondition.String
	}

	// Load images for the listing
	imagesQuery := `
		SELECT
			id, listing_id, file_path, file_name, file_size,
			content_type, is_main, storage_type,
			COALESCE(storage_bucket, '') as storage_bucket,
			COALESCE(public_url, '') as public_url
		FROM c2c_images
		WHERE listing_id = $1
		ORDER BY is_main DESC, id ASC
	`

	var images []models.MarketplaceImage
	err = s.db.SelectContext(ctx, &images, imagesQuery, listingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error().Err(err).Int("listing_id", listingID).Msg("Failed to load images for listing")
		// Don't fail the whole request if images fail to load
	}
	listing.Images = images

	return &listing, nil
}

// CreateStorefront создает новую витрину (проксирует к Database)
func (s *postgresMarketplaceStorage) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	// Этот метод должен быть реализован через Database, который имеет доступ к pool
	// Пока возвращаем ошибку, так как нужен доступ к pgx pool, а у нас только sqlx
	return nil, fmt.Errorf("CreateStorefront should be called through Database.CreateStorefront, not through marketplace storage")
}
