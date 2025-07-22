package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"backend/internal/domain/models"
)

// ErrNotFound ошибка когда запись не найдена
var ErrNotFound = errors.New("record not found")

// StorefrontRepository интерфейс репозитория витрин
type StorefrontRepository interface {
	// Основные CRUD операции
	Create(ctx context.Context, storefront *models.StorefrontCreateDTO) (*models.Storefront, error)
	GetByID(ctx context.Context, id int) (*models.Storefront, error)
	GetBySlug(ctx context.Context, slug string) (*models.Storefront, error)
	Update(ctx context.Context, id int, updates *models.StorefrontUpdateDTO) error
	Delete(ctx context.Context, id int) error

	// Поиск и фильтрация
	List(ctx context.Context, filter *models.StorefrontFilter) ([]*models.Storefront, int, error)
	GetMapData(ctx context.Context, bounds GeoBounds, filter *models.StorefrontFilter) ([]*models.StorefrontMapData, error)
	GetClusters(ctx context.Context, bounds GeoBounds, zoomLevel int) ([]*models.MapCluster, error)

	// Геолокация
	FindNearby(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*models.Storefront, error)
	GetBusinessesInBuilding(ctx context.Context, lat, lng float64, radiusM float64) ([]*models.StorefrontMapData, error)

	// Управление персоналом
	AddStaff(ctx context.Context, staff *models.StorefrontStaff) error
	UpdateStaff(ctx context.Context, id int, permissions models.JSONB) error
	RemoveStaff(ctx context.Context, storefrontID, userID int) error
	GetStaff(ctx context.Context, storefrontID int) ([]*models.StorefrontStaff, error)

	// Часы работы
	SetWorkingHours(ctx context.Context, hours []*models.StorefrontHours) error
	GetWorkingHours(ctx context.Context, storefrontID int) ([]*models.StorefrontHours, error)
	IsOpenNow(ctx context.Context, storefrontID int) (bool, error)

	// Методы оплаты
	SetPaymentMethods(ctx context.Context, methods []*models.StorefrontPaymentMethod) error
	GetPaymentMethods(ctx context.Context, storefrontID int) ([]*models.StorefrontPaymentMethod, error)

	// Методы доставки
	SetDeliveryOptions(ctx context.Context, options []*models.StorefrontDeliveryOption) error
	GetDeliveryOptions(ctx context.Context, storefrontID int) ([]*models.StorefrontDeliveryOption, error)

	// Аналитика
	RecordView(ctx context.Context, storefrontID int) error
	RecordAnalytics(ctx context.Context, analytics *models.StorefrontAnalytics) error
	GetAnalytics(ctx context.Context, storefrontID int, from, to time.Time) ([]*models.StorefrontAnalytics, error)
	RecordEvent(ctx context.Context, event *StorefrontEvent) error

	// Проверки прав
	IsOwner(ctx context.Context, storefrontID, userID int) (bool, error)
	HasPermission(ctx context.Context, storefrontID, userID int, permission string) (bool, error)
}

// GeoBounds границы географической области
type GeoBounds struct {
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
}

// storefrontRepo реализация репозитория витрин
type storefrontRepo struct {
	db *Database
}

// NewStorefrontRepository создает новый репозиторий витрин
func NewStorefrontRepository(db *Database) StorefrontRepository {
	return &storefrontRepo{db: db}
}

// Create создает новую витрину
func (r *storefrontRepo) Create(ctx context.Context, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Создаем витрину
	var storefront models.Storefront
	// Преобразуем название страны в код
	countryCode := getCountryCode(dto.Location.Country)

	err = tx.QueryRow(ctx, `
		INSERT INTO storefronts (
			user_id, slug, name, description,
			logo_url, banner_url, theme,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude,
			settings, seo_meta,
			is_active, subscription_plan, commission_rate
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18,
			$19, $20, $21
		)
		RETURNING id, created_at, updated_at
	`,
		dto.UserID, dto.Slug, dto.Name, dto.Description,
		"", "", dto.Theme,
		dto.Phone, dto.Email, dto.Website,
		dto.Location.FullAddress, dto.Location.City, dto.Location.PostalCode, countryCode,
		dto.Location.UserLat, dto.Location.UserLng,
		dto.Settings, dto.SEOMeta,
		false, models.SubscriptionPlanStarter, 3.00,
	).Scan(&storefront.ID, &storefront.CreatedAt, &storefront.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create storefront: %w", err)
	}

	// Создаем владельца в staff
	_, err = tx.Exec(ctx, `
		INSERT INTO storefront_staff (storefront_id, user_id, role, permissions)
		VALUES ($1, $2, $3, $4)
	`, storefront.ID, dto.UserID, models.StaffRoleOwner, getOwnerPermissions())
	if err != nil {
		return nil, fmt.Errorf("failed to add owner to staff: %w", err)
	}

	// Устанавливаем дефолтные часы работы (9:00-18:00 пн-пт)
	for day := 1; day <= 5; day++ {
		_, err = tx.Exec(ctx, `
			INSERT INTO storefront_hours (storefront_id, day_of_week, open_time, close_time)
			VALUES ($1, $2, $3, $4)
		`, storefront.ID, day, "09:00", "18:00")
		if err != nil {
			return nil, fmt.Errorf("failed to set working hours: %w", err)
		}
	}

	// Добавляем базовые методы оплаты
	paymentMethods := []struct {
		method models.PaymentMethodType
		fee    float64
	}{
		{models.PaymentMethodCash, 0},
		{models.PaymentMethodCOD, 2.5},
		{models.PaymentMethodCard, 2.0},
	}

	for _, pm := range paymentMethods {
		_, err = tx.Exec(ctx, `
			INSERT INTO storefront_payment_methods (
				storefront_id, method_type, is_enabled, transaction_fee
			) VALUES ($1, $2, $3, $4)
		`, storefront.ID, pm.method, true, pm.fee)
		if err != nil {
			return nil, fmt.Errorf("failed to add payment method %s: %w", pm.method, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Загружаем полную информацию о созданной витрине
	return r.GetByID(ctx, storefront.ID)
}

// GetByID получает витрину по ID
func (r *storefrontRepo) GetByID(ctx context.Context, id int) (*models.Storefront, error) {
	var s models.Storefront
	var theme, settings, seoMeta, aiConfig json.RawMessage

	err := r.db.pool.QueryRow(ctx, `
		SELECT 
			id, user_id, slug, name, description,
			logo_url, banner_url, COALESCE(theme, '{}')::jsonb,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude,
			COALESCE(settings, '{}')::jsonb, COALESCE(seo_meta, '{}')::jsonb,
			is_active, is_verified, verification_date,
			rating, reviews_count, products_count, sales_count, views_count,
			subscription_plan, subscription_expires_at, commission_rate,
			ai_agent_enabled, COALESCE(ai_agent_config, '{}')::jsonb, live_shopping_enabled, group_buying_enabled,
			created_at, updated_at
		FROM storefronts
		WHERE id = $1
	`, id).Scan(
		&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description,
		&s.LogoURL, &s.BannerURL, &theme,
		&s.Phone, &s.Email, &s.Website,
		&s.Address, &s.City, &s.PostalCode, &s.Country, &s.Latitude, &s.Longitude,
		&settings, &seoMeta,
		&s.IsActive, &s.IsVerified, &s.VerificationDate,
		&s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount, &s.ViewsCount,
		&s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
		&s.AIAgentEnabled, &aiConfig, &s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}

	// Конвертируем json.RawMessage в JSONB
	if theme != nil {
		json.Unmarshal(theme, &s.Theme)
	}
	if settings != nil {
		json.Unmarshal(settings, &s.Settings)
	}
	if seoMeta != nil {
		json.Unmarshal(seoMeta, &s.SEOMeta)
	}
	if aiConfig != nil {
		json.Unmarshal(aiConfig, &s.AIAgentConfig)
	}

	return &s, nil
}

// GetBySlug получает витрину по slug
func (r *storefrontRepo) GetBySlug(ctx context.Context, slug string) (*models.Storefront, error) {
	var id int
	err := r.db.pool.QueryRow(ctx, "SELECT id FROM storefronts WHERE slug = $1", slug).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

// Update обновляет витрину
func (r *storefrontRepo) Update(ctx context.Context, id int, dto *models.StorefrontUpdateDTO) error {
	var setClauses []string
	var args []interface{}
	argCount := 1

	// Динамически строим UPDATE запрос
	if dto.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *dto.Name)
		argCount++
	}
	if dto.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *dto.Description)
		argCount++
	}
	if dto.Theme != nil {
		setClauses = append(setClauses, fmt.Sprintf("theme = $%d", argCount))
		args = append(args, dto.Theme)
		argCount++
	}
	if dto.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argCount))
		args = append(args, *dto.Phone)
		argCount++
	}
	if dto.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argCount))
		args = append(args, *dto.Email)
		argCount++
	}
	if dto.Website != nil {
		setClauses = append(setClauses, fmt.Sprintf("website = $%d", argCount))
		args = append(args, *dto.Website)
		argCount++
	}
	if dto.LogoURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("logo_url = $%d", argCount))
		args = append(args, *dto.LogoURL)
		argCount++
	}
	if dto.BannerURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("banner_url = $%d", argCount))
		args = append(args, *dto.BannerURL)
		argCount++
	}
	if dto.Location != nil {
		setClauses = append(setClauses, fmt.Sprintf("address = $%d", argCount))
		args = append(args, dto.Location.FullAddress)
		argCount++

		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argCount))
		args = append(args, dto.Location.City)
		argCount++

		setClauses = append(setClauses, fmt.Sprintf("postal_code = $%d", argCount))
		args = append(args, dto.Location.PostalCode)
		argCount++

		setClauses = append(setClauses, fmt.Sprintf("latitude = $%d", argCount))
		args = append(args, dto.Location.BuildingLat)
		argCount++

		setClauses = append(setClauses, fmt.Sprintf("longitude = $%d", argCount))
		args = append(args, dto.Location.BuildingLng)
		argCount++
	}
	if dto.Settings != nil {
		setClauses = append(setClauses, fmt.Sprintf("settings = $%d", argCount))
		args = append(args, dto.Settings)
		argCount++
	}
	if dto.SEOMeta != nil {
		setClauses = append(setClauses, fmt.Sprintf("seo_meta = $%d", argCount))
		args = append(args, dto.SEOMeta)
		argCount++
	}
	if dto.AIAgentEnabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("ai_agent_enabled = $%d", argCount))
		args = append(args, *dto.AIAgentEnabled)
		argCount++
	}
	if dto.LiveShoppingEnabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("live_shopping_enabled = $%d", argCount))
		args = append(args, *dto.LiveShoppingEnabled)
		argCount++
	}
	if dto.GroupBuyingEnabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("group_buying_enabled = $%d", argCount))
		args = append(args, *dto.GroupBuyingEnabled)
		argCount++
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE storefronts SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		argCount,
	)

	result, err := r.db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update storefront: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete удаляет витрину (soft delete)
func (r *storefrontRepo) Delete(ctx context.Context, id int) error {
	result, err := r.db.pool.Exec(ctx,
		"UPDATE storefronts SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete storefront: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// List возвращает список витрин с фильтрацией
func (r *storefrontRepo) List(ctx context.Context, filter *models.StorefrontFilter) ([]*models.Storefront, int, error) {
	// Строим динамический запрос
	whereConditions := []string{"1=1"}
	countWhereConditions := []string{"1=1"}
	args := []interface{}{}
	argCount := 1

	if filter.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("user_id = $%d", argCount))
		countWhereConditions = append(countWhereConditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.IsActive != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argCount))
		countWhereConditions = append(countWhereConditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filter.IsActive)
		argCount++
	}

	if filter.IsVerified != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_verified = $%d", argCount))
		countWhereConditions = append(countWhereConditions, fmt.Sprintf("is_verified = $%d", argCount))
		args = append(args, *filter.IsVerified)
		argCount++
	}

	if filter.City != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("LOWER(city) = LOWER($%d)", argCount))
		countWhereConditions = append(countWhereConditions, fmt.Sprintf("LOWER(city) = LOWER($%d)", argCount))
		args = append(args, *filter.City)
		argCount++
	}

	if filter.MinRating != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("rating >= $%d", argCount))
		countWhereConditions = append(countWhereConditions, fmt.Sprintf("rating >= $%d", argCount))
		args = append(args, *filter.MinRating)
		argCount++
	}

	if filter.Search != nil && *filter.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argCount, argCount))
		countWhereConditions = append(countWhereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argCount, argCount))
		searchTerm := "%" + *filter.Search + "%"
		args = append(args, searchTerm)
		argCount++
	}

	// Геофильтр
	if filter.Latitude != nil && filter.Longitude != nil && filter.RadiusKm != nil {
		whereConditions = append(whereConditions, fmt.Sprintf(`
			earth_distance(
				ll_to_earth(latitude, longitude),
				ll_to_earth($%d, $%d)
			) <= $%d * 1000
		`, argCount, argCount+1, argCount+2))

		countWhereConditions = append(countWhereConditions, fmt.Sprintf(`
			earth_distance(
				ll_to_earth(latitude, longitude),
				ll_to_earth($%d, $%d)
			) <= $%d * 1000
		`, argCount, argCount+1, argCount+2))

		args = append(args, *filter.Latitude, *filter.Longitude, *filter.RadiusKm)
		argCount += 3
	}

	// Подсчет общего количества
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM storefronts WHERE %s", strings.Join(countWhereConditions, " AND "))
	var totalCount int
	err := r.db.pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count storefronts: %w", err)
	}

	// Сортировка
	orderBy := "created_at DESC"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "rating":
			orderBy = "rating"
		case "products_count":
			orderBy = "products_count"
		case "distance":
			if filter.Latitude != nil && filter.Longitude != nil {
				orderBy = fmt.Sprintf("earth_distance(ll_to_earth(latitude, longitude), ll_to_earth(%f, %f))",
					*filter.Latitude, *filter.Longitude)
			}
		}

		if filter.SortOrder == "desc" {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}
	}

	// Основной запрос
	query := fmt.Sprintf(`
		SELECT 
			id, user_id, slug, name, description,
			logo_url, banner_url, theme,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude,
			settings, seo_meta,
			is_active, is_verified, verification_date,
			rating, reviews_count, products_count, sales_count, views_count,
			subscription_plan, subscription_expires_at, commission_rate,
			ai_agent_enabled, ai_agent_config, live_shopping_enabled, group_buying_enabled,
			created_at, updated_at
		FROM storefronts
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, strings.Join(whereConditions, " AND "), orderBy, argCount, argCount+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list storefronts: %w", err)
	}
	defer rows.Close()

	var storefronts []*models.Storefront
	for rows.Next() {
		s := &models.Storefront{}
		var theme, settings, seoMeta, aiConfig []byte

		err := rows.Scan(
			&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description,
			&s.LogoURL, &s.BannerURL, &theme,
			&s.Phone, &s.Email, &s.Website,
			&s.Address, &s.City, &s.PostalCode, &s.Country, &s.Latitude, &s.Longitude,
			&settings, &seoMeta,
			&s.IsActive, &s.IsVerified, &s.VerificationDate,
			&s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount, &s.ViewsCount,
			&s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
			&s.AIAgentEnabled, &aiConfig, &s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan storefront: %w", err)
		}

		// Парсим JSONB поля
		if theme != nil {
			json.Unmarshal(theme, &s.Theme)
		}
		if settings != nil {
			json.Unmarshal(settings, &s.Settings)
		}
		if seoMeta != nil {
			json.Unmarshal(seoMeta, &s.SEOMeta)
		}
		if aiConfig != nil {
			json.Unmarshal(aiConfig, &s.AIAgentConfig)
		}

		storefronts = append(storefronts, s)
	}

	return storefronts, totalCount, nil
}

// FindNearby находит витрины в радиусе
func (r *storefrontRepo) FindNearby(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*models.Storefront, error) {
	filter := &models.StorefrontFilter{
		Latitude:  &lat,
		Longitude: &lng,
		RadiusKm:  &radiusKm,
		IsActive:  boolPtr(true),
		Limit:     limit,
		Offset:    0,
		SortBy:    "distance",
		SortOrder: "asc",
	}

	storefronts, _, err := r.List(ctx, filter)
	return storefronts, err
}

// GetMapData получает данные для отображения на карте
func (r *storefrontRepo) GetMapData(ctx context.Context, bounds GeoBounds, filter *models.StorefrontFilter) ([]*models.StorefrontMapData, error) {
	query := `
		SELECT 
			s.id, s.slug, s.name, s.latitude, s.longitude, s.rating, s.logo_url,
			s.address, s.phone, s.products_count,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM storefront_hours 
					WHERE storefront_id = s.id 
					AND day_of_week = EXTRACT(DOW FROM CURRENT_TIMESTAMP)
					AND NOT is_closed
					AND CURRENT_TIME BETWEEN open_time AND close_time
				) THEN true 
				ELSE false 
			END as working_now,
			CASE WHEN EXISTS (
				SELECT 1 FROM storefront_payment_methods 
				WHERE storefront_id = s.id AND method_type = 'cod' AND is_enabled
			) THEN true ELSE false END as supports_cod,
			CASE WHEN EXISTS (
				SELECT 1 FROM storefront_delivery_options 
				WHERE storefront_id = s.id AND is_enabled AND provider != 'self_pickup'
			) THEN true ELSE false END as has_delivery,
			CASE WHEN EXISTS (
				SELECT 1 FROM storefront_delivery_options 
				WHERE storefront_id = s.id AND is_enabled AND provider = 'self_pickup'
			) THEN true ELSE false END as has_self_pickup,
			CASE WHEN EXISTS (
				SELECT 1 FROM storefront_payment_methods 
				WHERE storefront_id = s.id AND method_type = 'card' AND is_enabled
			) THEN true ELSE false END as accepts_cards
		FROM storefronts s
		WHERE s.is_active = true
		AND s.latitude BETWEEN $1 AND $2
		AND s.longitude BETWEEN $3 AND $4
	`

	args := []interface{}{bounds.MinLat, bounds.MaxLat, bounds.MinLng, bounds.MaxLng}

	// Добавляем дополнительные фильтры
	if filter != nil {
		if filter.MinRating != nil {
			query += " AND s.rating >= $5"
			args = append(args, *filter.MinRating)
		}
		// ... другие фильтры
	}

	rows, err := r.db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get map data: %w", err)
	}
	defer rows.Close()

	var results []*models.StorefrontMapData
	for rows.Next() {
		var data models.StorefrontMapData
		err := rows.Scan(
			&data.ID, &data.Slug, &data.Name, &data.Latitude, &data.Longitude,
			&data.Rating, &data.LogoURL, &data.Address, &data.Phone, &data.ProductsCount,
			&data.WorkingNow, &data.SupportsCOD, &data.HasDelivery, &data.HasSelfPickup, &data.AcceptsCards,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan map data: %w", err)
		}
		results = append(results, &data)
	}

	return results, nil
}

// GetBusinessesInBuilding получает все бизнесы в здании
func (r *storefrontRepo) GetBusinessesInBuilding(ctx context.Context, lat, lng float64, radiusM float64) ([]*models.StorefrontMapData, error) {
	bounds := GeoBounds{
		MinLat: lat - (radiusM / 111000.0), // приблизительно
		MaxLat: lat + (radiusM / 111000.0),
		MinLng: lng - (radiusM / (111000.0 * math.Cos(lat*math.Pi/180))),
		MaxLng: lng + (radiusM / (111000.0 * math.Cos(lat*math.Pi/180))),
	}

	return r.GetMapData(ctx, bounds, nil)
}

// IsOwner проверяет является ли пользователь владельцем
func (r *storefrontRepo) IsOwner(ctx context.Context, storefrontID, userID int) (bool, error) {
	var exists bool
	err := r.db.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM storefronts 
			WHERE id = $1 AND user_id = $2
		)
	`, storefrontID, userID).Scan(&exists)

	return exists, err
}

// HasPermission проверяет права пользователя
func (r *storefrontRepo) HasPermission(ctx context.Context, storefrontID, userID int, permission string) (bool, error) {
	// Владелец имеет все права
	isOwner, err := r.IsOwner(ctx, storefrontID, userID)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}

	// Проверяем права сотрудника
	var permissions models.JSONB
	err = r.db.pool.QueryRow(ctx, `
		SELECT permissions FROM storefront_staff 
		WHERE storefront_id = $1 AND user_id = $2
	`, storefrontID, userID).Scan(&permissions)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// Проверяем наличие конкретного разрешения
	if val, ok := permissions[permission]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal, nil
		}
	}

	return false, nil
}

// Helper функции

func generateSlug(name string) string {
	// Простая генерация slug - в продакшене нужно использовать библиотеку
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Добавляем timestamp для уникальности
	return fmt.Sprintf("%s-%d", slug, time.Now().Unix())
}

func getOwnerPermissions() models.JSONB {
	return models.JSONB{
		"can_add_products":     true,
		"can_edit_products":    true,
		"can_delete_products":  true,
		"can_view_orders":      true,
		"can_process_orders":   true,
		"can_refund_orders":    true,
		"can_edit_storefront":  true,
		"can_manage_staff":     true,
		"can_view_analytics":   true,
		"can_manage_payments":  true,
		"can_reply_to_reviews": true,
		"can_send_messages":    true,
	}
}

func boolPtr(b bool) *bool {
	return &b
}

// getCountryCode преобразует название страны в ISO код
func getCountryCode(countryName string) string {
	// Маппинг названий стран на ISO коды
	countryMap := map[string]string{
		"Сербия":         "RS",
		"Serbia":         "RS",
		"Россия":         "RU",
		"Russia":         "RU",
		"США":            "US",
		"USA":            "US",
		"United States":  "US",
		"Германия":       "DE",
		"Germany":        "DE",
		"Франция":        "FR",
		"France":         "FR",
		"Италия":         "IT",
		"Italy":          "IT",
		"Испания":        "ES",
		"Spain":          "ES",
		"Великобритания": "GB",
		"United Kingdom": "GB",
		"UK":             "GB",
	}

	// Пробуем найти код по названию
	if code, ok := countryMap[countryName]; ok {
		return code
	}

	// Если уже передан код из 2 символов, возвращаем его
	if len(countryName) == 2 {
		return strings.ToUpper(countryName)
	}

	// По умолчанию возвращаем код Сербии
	return "RS"
}
