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
			uca.id,
			uca.category_id,
			uca.attribute_id,
			uca.is_enabled,
			uca.is_required,
			uca.sort_order,
			uca.category_specific_options,
			ua.name as attribute_name,
			ua.slug as attribute_slug,
			ua.type as attribute_type,
			ua.is_filterable,
			ua.is_searchable,
			ua.options,
			ua.validation_rules,
			ua.is_multi_select,
			ua.is_custom_value_allowed,
			ua.is_dynamic
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
