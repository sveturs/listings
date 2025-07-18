package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/gis/types"
)

// PostGISRepository репозиторий для работы с пространственными данными
type PostGISRepository struct {
	db *sqlx.DB
}

// NewPostGISRepository создает новый репозиторий
func NewPostGISRepository(db *sqlx.DB) *PostGISRepository {
	return &PostGISRepository{db: db}
}

// SearchListings поиск объявлений в заданной области (обновлено для unified_geo)
func (r *PostGISRepository) SearchListings(ctx context.Context, params types.SearchParams) ([]types.GeoListing, int64, error) {
	// Проверяем, существует ли таблица unified_geo
	var tableExists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'unified_geo')").Scan(&tableExists)
	if err == nil && tableExists {
		return r.searchUnifiedGeo(ctx, params)
	}

	// Fallback к старому методу
	return r.searchLegacy(ctx, params)
}

// searchUnifiedGeo новый метод поиска с unified geo system
func (r *PostGISRepository) searchUnifiedGeo(ctx context.Context, params types.SearchParams) ([]types.GeoListing, int64, error) {
	query := `
		SELECT
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.id
				WHEN ug.source_type = 'storefront_product' THEN sp.id
			END as id,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.title
				WHEN ug.source_type = 'storefront_product' THEN sp.name
			END as title,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN COALESCE(ml.description, '')
				WHEN ug.source_type = 'storefront_product' THEN COALESCE(sp.description, '')
			END as description,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.price
				WHEN ug.source_type = 'storefront_product' THEN sp.price
			END as price,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN COALESCE(mc1.name, '')
				WHEN ug.source_type = 'storefront_product' THEN COALESCE(mc2.name, '')
			END as category,
			ST_Y(ug.location::geometry) as lat,
			ST_X(ug.location::geometry) as lng,
			COALESCE(ug.formatted_address, '') as address,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.user_id
				WHEN ug.source_type = 'storefront_product' THEN s.user_id
			END as user_id,
			CASE
				WHEN ug.source_type = 'storefront_product' THEN sp.storefront_id
				ELSE NULL
			END as storefront_id,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.status
				WHEN ug.source_type = 'storefront_product' THEN CASE WHEN sp.is_active THEN 'active' ELSE 'inactive' END
			END as status,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.created_at
				WHEN ug.source_type = 'storefront_product' THEN sp.created_at
			END as created_at,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.updated_at
				WHEN ug.source_type = 'storefront_product' THEN sp.updated_at
			END as updated_at,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN ml.views_count
				WHEN ug.source_type = 'storefront_product' THEN sp.view_count
			END as views_count,
			CASE
				WHEN ug.source_type = 'marketplace_listing' THEN COALESCE(rc.average_rating, 0)
				ELSE 0
			END as rating,
			ug.source_type::text as item_type`

	// Добавляем расчет расстояния если есть центр поиска
	if params.Center != nil {
		query += fmt.Sprintf(`,
				ST_Distance(
					ug.location,
					ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geography
				) as distance`, params.Center.Lng, params.Center.Lat)
	}

	query += `
			FROM unified_geo ug
			LEFT JOIN marketplace_listings ml ON ug.source_type = 'marketplace_listing' AND ug.source_id = ml.id
			LEFT JOIN storefront_products sp ON ug.source_type = 'storefront_product' AND ug.source_id = sp.id
			LEFT JOIN storefronts s ON sp.storefront_id = s.id
			LEFT JOIN marketplace_categories mc1 ON ml.category_id = mc1.id
			LEFT JOIN marketplace_categories mc2 ON sp.category_id = mc2.id
			LEFT JOIN rating_cache rc ON rc.entity_type = 'listing' AND rc.entity_id = ml.id
			WHERE (
				(ug.source_type = 'marketplace_listing' AND ml.status = 'active') OR
				(ug.source_type = 'storefront_product' AND sp.is_active = true)
			)`

	var args []interface{}
	argIndex := 1

	// Добавляем фильтры
	if params.Center != nil && params.RadiusKm > 0 {
		query += fmt.Sprintf(` AND ST_DWithin(
			ug.location,
			ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography,
			$%d)`, argIndex, argIndex+1, argIndex+2)
		args = append(args, params.Center.Lng, params.Center.Lat, params.RadiusKm*1000)
		argIndex += 3
	}

	// Подсчет общего количества
	countQuery := "SELECT COUNT(*) " + strings.Replace(query, query[:strings.Index(query, "FROM")], "", 1)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count listings: %w", err)
	}

	// Добавляем сортировку и лимит
	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, params.Offset)

	// Выполняем запрос
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search query: %w", err)
	}
	defer rows.Close()

	var listings []types.GeoListing
	for rows.Next() {
		var listing types.GeoListing
		var lat, lng float64
		var distance sql.NullFloat64
		var storefrontID sql.NullInt64

		scanArgs := []interface{}{
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Category,
			&lat,
			&lng,
			&listing.Address,
			&listing.UserID,
			&storefrontID,
			&listing.Status,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.ViewsCount,
			&listing.Rating,
			&listing.ItemType,
		}

		if params.Center != nil {
			scanArgs = append(scanArgs, &distance)
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan listing: %w", err)
		}

		listing.Location = types.Point{
			Lat: lat,
			Lng: lng,
		}

		if distance.Valid {
			listing.Distance = &distance.Float64
		}

		if storefrontID.Valid {
			storefrontIDInt := int(storefrontID.Int64)
			listing.StorefrontID = &storefrontIDInt
		}

		// Загружаем изображения
		images, err := r.getListingImages(ctx, listing.ID, listing.ItemType)
		if err != nil {
			// Игнорируем ошибки загрузки изображений
			images = []string{}
		}
		listing.Images = images

		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}

// getListingImages загружает изображения в зависимости от типа элемента
func (r *PostGISRepository) getListingImages(ctx context.Context, itemID int, itemType string) ([]string, error) {
	var query string

	switch itemType {
	case "marketplace_listing":
		query = `
			SELECT public_url
			FROM marketplace_images
			WHERE listing_id = $1 AND public_url IS NOT NULL AND public_url != ''
			ORDER BY is_main DESC, created_at
			LIMIT 5`
	case "storefront_product":
		query = `
			SELECT spvi.image_url
			FROM storefront_product_variants spv
			JOIN storefront_product_variant_images spvi ON spv.id = spvi.variant_id
			WHERE spv.product_id = $1 AND spv.is_active = true AND spvi.image_url IS NOT NULL
			ORDER BY spv.is_default DESC, spvi.is_main DESC, spvi.display_order
			LIMIT 5`
	default:
		return []string{}, nil
	}

	var images []string
	err := r.db.SelectContext(ctx, &images, query, itemID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// searchLegacy старый метод поиска (для совместимости)
func (r *PostGISRepository) searchLegacy(ctx context.Context, params types.SearchParams) ([]types.GeoListing, int64, error) {
	var listings []types.GeoListing
	var totalCount int64

	// Базовый запрос
	query := `
		WITH filtered_listings AS (
			SELECT
				ml.id,
				ml.title,
				ml.description,
				ml.price,
				mc.name as category,
				ST_Y(lg.location::geometry) as lat,
				ST_X(lg.location::geometry) as lng,
				ml.location as address,
				ml.user_id,
				ml.status,
				ml.created_at,
				ml.updated_at,
				ml.views_count,
				COALESCE(rc.average_rating, 0) as rating`

	// Добавляем расчет расстояния если есть центр поиска
	if params.Center != nil {
		query += fmt.Sprintf(`,
				ST_Distance(
					lg.location::geography,
					ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geography
				) as distance`, params.Center.Lng, params.Center.Lat)
	}

	query += `
			FROM map_items_cache mic
			WHERE mic.status = 'active' AND mic.latitude IS NOT NULL AND mic.longitude IS NOT NULL`

	// Применяем фильтры
	var conditions []string
	var args []interface{}
	argCount := 1

	// Фильтр по границам
	if params.Bounds != nil {
		conditions = append(conditions, fmt.Sprintf(
			"lg.location && ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)",
			argCount, argCount+1, argCount+2, argCount+3,
		))
		args = append(args, params.Bounds.West, params.Bounds.South, params.Bounds.East, params.Bounds.North)
		argCount += 4
	}

	// Фильтр по радиусу
	if params.Center != nil && params.RadiusKm > 0 {
		conditions = append(conditions, fmt.Sprintf(
			"ST_DWithin(lg.location::geography, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d)",
			argCount, argCount+1, argCount+2,
		))
		args = append(args, params.Center.Lng, params.Center.Lat, params.RadiusKm*1000) // в метрах
		argCount += 3
	}

	// Фильтр по категориям (по названию)
	if len(params.Categories) > 0 {
		placeholders := make([]string, len(params.Categories))
		for i, cat := range params.Categories {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, cat)
			argCount++
		}
		conditions = append(conditions, fmt.Sprintf("mc.name IN (%s)", strings.Join(placeholders, ",")))
	}

	// Фильтр по категориям (по ID)
	if len(params.CategoryIDs) > 0 {
		placeholders := make([]string, len(params.CategoryIDs))
		for i, catID := range params.CategoryIDs {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, catID)
			argCount++
		}
		conditions = append(conditions, fmt.Sprintf("mc.id IN (%s)", strings.Join(placeholders, ",")))
	}

	// Фильтр по цене
	if params.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("ml.price >= $%d", argCount))
		args = append(args, *params.MinPrice)
		argCount++
	}

	if params.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("ml.price <= $%d", argCount))
		args = append(args, *params.MaxPrice)
		argCount++
	}

	// Фильтр по валюте (пропускаем, так как все объявления в RSD)
	// if params.Currency != "" {
	//	conditions = append(conditions, fmt.Sprintf("ml.currency = $%d", argCount))
	//	args = append(args, params.Currency)
	//	argCount++
	// }

	// Фильтр по пользователю
	if params.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("ml.user_id = $%d", argCount))
		args = append(args, *params.UserID)
		argCount++
	}

	// Фильтр по статусу
	if params.Status != "" {
		conditions = append(conditions, fmt.Sprintf("ml.status = $%d", argCount))
		args = append(args, params.Status)
		argCount++
	}

	// Текстовый поиск
	if params.SearchQuery != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(ml.title ILIKE $%d OR ml.description ILIKE $%d)",
			argCount, argCount+1,
		))
		searchPattern := "%" + params.SearchQuery + "%"
		args = append(args, searchPattern, searchPattern)
		_ = argCount + 2 // Note: argCount is used for documentation purposes
	}

	// Применяем условия
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += "\n)"

	// Подсчет общего количества
	countQuery := "SELECT COUNT(*) FROM filtered_listings"

	// Основной запрос с сортировкой и пагинацией
	selectQuery := query + "\nSELECT * FROM filtered_listings"

	// Сортировка
	switch params.SortBy {
	case "distance":
		if params.Center != nil {
			selectQuery += " ORDER BY distance"
		} else {
			selectQuery += " ORDER BY created_at DESC"
		}
	case "price":
		selectQuery += " ORDER BY price"
	case "created_at":
		selectQuery += " ORDER BY created_at"
	default:
		selectQuery += " ORDER BY created_at DESC"
	}

	if params.SortOrder == "desc" && params.SortBy != "" {
		selectQuery += " DESC"
	} else if params.SortOrder == "asc" {
		selectQuery += " ASC"
	}

	// Пагинация
	if params.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT %d", params.Limit)
	} else {
		selectQuery += " LIMIT 100" // Дефолтный лимит
	}

	if params.Offset > 0 {
		selectQuery += fmt.Sprintf(" OFFSET %d", params.Offset)
	}

	// Выполняем запросы
	err := r.db.GetContext(ctx, &totalCount, query+"\n"+countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count listings: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search listings: %w", err)
	}
	defer rows.Close()

	// Сканируем результаты
	for rows.Next() {
		var listing types.GeoListing
		var lat, lng float64
		var distance sql.NullFloat64

		scanArgs := []interface{}{
			&listing.ID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Category,
			&lat,
			&lng,
			&listing.Address,
			&listing.UserID,
			&listing.Status,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.ViewsCount,
			&listing.Rating,
		}

		if params.Center != nil {
			scanArgs = append(scanArgs, &distance)
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan listing: %w", err)
		}

		listing.Location = types.Point{Lat: lat, Lng: lng}
		if distance.Valid {
			listing.Distance = &distance.Float64
		}

		// Загружаем изображения
		listing.Images, _ = r.getListingImages(ctx, listing.ID, "marketplace_listing")

		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}

// GetListingByID получение объявления по ID с геоданными
func (r *PostGISRepository) GetListingByID(ctx context.Context, id int) (*types.GeoListing, error) {
	query := `
		SELECT
			ml.id,
			ml.title,
			ml.description,
			ml.price,
			mc.name as category,
			ST_Y(lg.location::geometry) as lat,
			ST_X(lg.location::geometry) as lng,
			ml.location as address,
			ml.user_id,
			ml.status,
			ml.created_at,
			ml.updated_at,
			ml.views_count,
			COALESCE(rc.average_rating, 0) as rating
		FROM marketplace_listings ml
		LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
		LEFT JOIN listings_geo lg ON ml.id = lg.listing_id
		LEFT JOIN rating_cache rc ON rc.entity_type = 'listing' AND rc.entity_id = ml.id
		WHERE ml.id = $1 AND lg.location IS NOT NULL`

	var listing types.GeoListing
	var lat, lng float64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Category,
		&lat,
		&lng,
		&listing.Address,
		&listing.UserID,
		&listing.Status,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&listing.ViewsCount,
		&listing.Rating,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrLocationNotFound
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	listing.Location = types.Point{Lat: lat, Lng: lng}
	listing.Images, _ = r.getListingImages(ctx, listing.ID, "marketplace_listing")

	return &listing, nil
}

// UpdateListingLocation обновление геолокации объявления
func (r *PostGISRepository) UpdateListingLocation(ctx context.Context, id int, location types.Point, address string) error {
	// Обновляем запись в listings_geo
	query := `
		INSERT INTO listings_geo (listing_id, location, geohash, is_precise, blur_radius)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), ST_GeoHash(ST_SetSRID(ST_MakePoint($2, $3), 4326), 8), true, 0)
		ON CONFLICT (listing_id) DO UPDATE SET
			location = ST_SetSRID(ST_MakePoint($2, $3), 4326),
			geohash = ST_GeoHash(ST_SetSRID(ST_MakePoint($2, $3), 4326), 8),
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query, id, location.Lng, location.Lat)
	if err != nil {
		return fmt.Errorf("failed to update listing location: %w", err)
	}

	// Обновляем адрес в основной таблице
	updateAddressQuery := `
		UPDATE marketplace_listings
		SET location = $1, updated_at = NOW()
		WHERE id = $2`

	_, err = r.db.ExecContext(ctx, updateAddressQuery, address, id)
	if err != nil {
		return fmt.Errorf("failed to update listing address: %w", err)
	}

	return nil
}
