// backend/internal/storage/postgres/db_storefronts.go
package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
)

// Storefront methods

func (db *Database) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	// Генерируем slug если не указан
	slug := dto.Slug
	if slug == "" {
		slug = generateSlug(dto.Name)
	}

	// Если slug пустой после генерации (только кириллица), используем fallback
	if slug == "" {
		slug = fmt.Sprintf("store-%d", time.Now().Unix()%10000)
	}

	// Проверяем уникальность slug и добавляем суффикс если занят
	originalSlug := slug
	for attempt := 0; attempt < 5; attempt++ {
		var existingID int
		err := db.pool.QueryRow(ctx, "SELECT id FROM b2c_stores WHERE slug = $1 LIMIT 1", slug).Scan(&existingID)
		if err != nil {
			// Slug свободен
			break
		}
		// Slug занят, пробуем с суффиксом
		slug = fmt.Sprintf("%s-%d", originalSlug, time.Now().UnixNano()%10000+int64(attempt))
	}

	// Начинаем транзакцию
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // Rollback is safe to call even after commit
	}()

	// Вставляем витрину
	query := `
		INSERT INTO b2c_stores (
			user_id, slug, name, description,
			phone, email, website,
			address, city, postal_code, country,
			latitude, longitude,
			settings, seo_meta, theme,
			is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10, $11,
			$12, $13,
			$14, $15, $16,
			true, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`

	var storefront models.Storefront
	storefront.UserID = userID
	storefront.Slug = slug
	storefront.Name = dto.Name

	// Подготавливаем значения
	description := &dto.Description
	if dto.Description == "" {
		description = nil
	}

	phone := &dto.Phone
	if dto.Phone == "" {
		phone = nil
	}

	email := &dto.Email
	if dto.Email == "" {
		email = nil
	}

	website := &dto.Website
	if dto.Website == "" {
		website = nil
	}

	// Обрабатываем локацию
	var address, city, postalCode, country *string
	var lat, lng *float64

	if dto.Location.FullAddress != "" {
		address = &dto.Location.FullAddress
	}
	if dto.Location.City != "" {
		city = &dto.Location.City
	}
	if dto.Location.PostalCode != "" {
		postalCode = &dto.Location.PostalCode
	}
	if dto.Location.Country != "" {
		country = &dto.Location.Country
	}
	if dto.Location.UserLat != 0 {
		lat = &dto.Location.UserLat
	}
	if dto.Location.UserLng != 0 {
		lng = &dto.Location.UserLng
	}

	// Конвертируем settings и seo_meta в JSON
	settingsJSON := []byte("{}")
	if len(dto.Settings) > 0 {
		if jsonBytes, err := json.Marshal(dto.Settings); err == nil {
			settingsJSON = jsonBytes
		}
	}

	seoMetaJSON := []byte("{}")
	if len(dto.SEOMeta) > 0 {
		if jsonBytes, err := json.Marshal(dto.SEOMeta); err == nil {
			seoMetaJSON = jsonBytes
		}
	}

	themeJSON := []byte(`{"layout": "grid", "primaryColor": "#1976d2"}`)
	if len(dto.Theme) > 0 {
		if jsonBytes, err := json.Marshal(dto.Theme); err == nil {
			themeJSON = jsonBytes
		}
	}

	err = tx.QueryRow(ctx, query,
		userID, slug, dto.Name, description,
		phone, email, website,
		address, city, postalCode, country,
		lat, lng,
		settingsJSON, seoMetaJSON, themeJSON,
	).Scan(&storefront.ID, &storefront.CreatedAt, &storefront.UpdatedAt)
	if err != nil {
		// Проверяем constraint violation (duplicate slug)
		if strings.Contains(err.Error(), "b2c_stores_slug_unique") {
			return nil, fmt.Errorf("storefront with this slug already exists, please try again")
		}
		return nil, fmt.Errorf("failed to insert storefront: %w", err)
	}

	// Копируем остальные поля
	storefront.Description = description
	storefront.Phone = phone
	storefront.Email = email
	storefront.Website = website
	storefront.Address = address
	storefront.City = city
	storefront.PostalCode = postalCode
	storefront.Country = country
	storefront.Latitude = lat
	storefront.Longitude = lng

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &storefront, nil
}

// generateSlug генерирует slug из строки
func generateSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	reg := regexp.MustCompile("[^a-zA-Z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func (db *Database) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	return nil, fmt.Errorf("storefront service temporarily disabled")
}

func (db *Database) IncrementListingViewCount(ctx context.Context, id int, userIdentifier string) error {
	return nil // Silent no-op
}

// B2C Product methods

func (db *Database) GetB2CProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error) {
	return []models.MarketplaceImage{}, nil
}

// Favorites

func (db *Database) AddToFavorites(ctx context.Context, userID, listingID int) error {
	return fmt.Errorf("favorites service temporarily disabled")
}

func (db *Database) RemoveFromFavorites(ctx context.Context, userID, listingID int) error {
	return fmt.Errorf("favorites service temporarily disabled")
}

func (db *Database) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return []models.MarketplaceListing{}, nil
}

func (db *Database) AddStorefrontToFavorites(ctx context.Context, userID, productID int) error {
	return fmt.Errorf("storefront favorites temporarily disabled")
}

func (db *Database) RemoveStorefrontFromFavorites(ctx context.Context, userID, productID int) error {
	return fmt.Errorf("storefront favorites temporarily disabled")
}

func (db *Database) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return []models.MarketplaceListing{}, nil
}

// DeleteStorefront удаляет витрину по ID (hard delete)
func (db *Database) DeleteStorefront(ctx context.Context, id int) error {
	query := `DELETE FROM b2c_stores WHERE id = $1`

	result, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete storefront: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("storefront not found")
	}

	return nil
}
func (db *Database) DeleteStorefrontIndex(ctx context.Context, id int) error { return nil }
func (db *Database) GetStorefrontOwnerByProductID(ctx context.Context, productID int) (int, error) {
	return 0, fmt.Errorf("disabled")
}

func (db *Database) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	query := `
		SELECT
			id, user_id, slug, name, description,
			logo_url, banner_url, theme,
			phone, email, website,
			address, city, postal_code, country,
			latitude, longitude,
			settings, seo_meta,
			is_active, is_verified, verification_date,
			rating, reviews_count, products_count, sales_count, views_count,
			subscription_plan, subscription_expires_at, commission_rate,
			ai_agent_enabled, ai_agent_config,
			live_shopping_enabled, group_buying_enabled,
			created_at, updated_at,
			formatted_address, geo_strategy, default_privacy_level, address_verified,
			subscription_id, is_subscription_active, followers_count
		FROM b2c_stores
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user storefronts: %w", err)
	}
	defer rows.Close()

	var storefronts []models.Storefront
	for rows.Next() {
		var s models.Storefront
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description,
			&s.LogoURL, &s.BannerURL, &s.Theme,
			&s.Phone, &s.Email, &s.Website,
			&s.Address, &s.City, &s.PostalCode, &s.Country,
			&s.Latitude, &s.Longitude,
			&s.Settings, &s.SEOMeta,
			&s.IsActive, &s.IsVerified, &s.VerificationDate,
			&s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount, &s.ViewsCount,
			&s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
			&s.AIAgentEnabled, &s.AIAgentConfig,
			&s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
			&s.CreatedAt, &s.UpdatedAt,
			&s.FormattedAddress, &s.GeoStrategy, &s.DefaultPrivacyLevel, &s.AddressVerified,
			&s.SubscriptionID, &s.IsSubscriptionActive, &s.FollowersCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan storefront: %w", err)
		}
		storefronts = append(storefronts, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating storefronts: %w", err)
	}

	return storefronts, nil
}
func (db *Database) IncrementViewsCount(ctx context.Context, id int) error               { return nil }
func (db *Database) IndexStorefront(ctx context.Context, store *models.Storefront) error { return nil }
func (db *Database) GetStorefrontBySlug(ctx context.Context, slug string) (*models.Storefront, error) {
	return nil, fmt.Errorf("disabled")
}

func (db *Database) UpdateStorefront(ctx context.Context, store *models.Storefront) error {
	return fmt.Errorf("disabled")
}
func (db *Database) ReindexAllStorefronts(ctx context.Context) error { return nil }
func (db *Database) SearchStorefronts(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	return &search.SearchResult{Listings: []*models.MarketplaceListing{}}, nil
}

func (db *Database) SuggestStorefronts(ctx context.Context, prefix string, size int) ([]string, error) {
	return []string{}, nil
}
func (db *Database) PrepareSearchIndex(ctx context.Context) error { return nil }
func (db *Database) Storefront() interface{}                      { return db.storefrontRepo }
