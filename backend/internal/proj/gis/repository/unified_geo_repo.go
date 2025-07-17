package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"backend/internal/proj/gis/types"
	"github.com/jmoiron/sqlx"
)

// UnifiedGeoRepository репозиторий для работы с unified geo system
type UnifiedGeoRepository struct {
	db *sqlx.DB
}

// NewUnifiedGeoRepository создает новый репозиторий
func NewUnifiedGeoRepository(db *sqlx.DB) *UnifiedGeoRepository {
	return &UnifiedGeoRepository{db: db}
}

// SearchListings выполняет пространственный поиск с unified geo system
func (r *UnifiedGeoRepository) SearchListings(ctx context.Context, params types.SearchParams) ([]types.GeoListing, int64, error) {
	// Базовый запрос из materialized view
	query := `
		SELECT
			mic.id,
			mic.name as title,
			COALESCE(mic.description, '') as description,
			COALESCE(mic.price, 0) as price,
			COALESCE(mic.category_name, '') as category,
			mic.latitude as lat,
			mic.longitude as lng,
			COALESCE(mic.formatted_address, '') as address,
			mic.user_id,
			mic.storefront_id,
			mic.status,
			mic.created_at,
			mic.updated_at,
			mic.views_count,
			mic.rating,
			mic.item_type,
			mic.display_strategy,
			COALESCE(mic.privacy_level::text, 'exact') as privacy_level,
			COALESCE(mic.blur_radius_meters, 0) as blur_radius_meters`

	// Добавляем расчет расстояния если есть центр поиска
	if params.Center != nil {
		query += fmt.Sprintf(`,
				ST_Distance(
					ST_SetSRID(ST_MakePoint(mic.longitude, mic.latitude), 4326)::geography,
					ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geography
				) as distance`, params.Center.Lng, params.Center.Lat)
	}

	query += `
			FROM map_items_cache mic
			WHERE mic.status = 'active' AND mic.latitude IS NOT NULL AND mic.longitude IS NOT NULL`

	var args []interface{}
	argIndex := 1

	// Добавляем фильтры
	if params.Bounds != nil {
		query += fmt.Sprintf(` AND mic.latitude BETWEEN $%d AND $%d AND mic.longitude BETWEEN $%d AND $%d`,
			argIndex, argIndex+1, argIndex+2, argIndex+3)
		args = append(args, params.Bounds.South, params.Bounds.North, params.Bounds.West, params.Bounds.East)
		argIndex += 4
	}

	if params.Center != nil && params.RadiusKm > 0 {
		query += fmt.Sprintf(` AND ST_DWithin(
			ST_SetSRID(ST_MakePoint(mic.longitude, mic.latitude), 4326)::geography,
			ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography,
			$%d)`, argIndex, argIndex+1, argIndex+2)
		args = append(args, params.Center.Lng, params.Center.Lat, params.RadiusKm*1000)
		argIndex += 3
	}

	if len(params.Categories) > 0 {
		placeholders := make([]string, len(params.Categories))
		for i, category := range params.Categories {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, category)
			argIndex++
		}
		query += fmt.Sprintf(` AND mic.category_name IN (%s)`, strings.Join(placeholders, ","))
	}

	if params.MinPrice != nil {
		query += fmt.Sprintf(` AND mic.price >= $%d`, argIndex)
		args = append(args, *params.MinPrice)
		argIndex++
	}

	if params.MaxPrice != nil {
		query += fmt.Sprintf(` AND mic.price <= $%d`, argIndex)
		args = append(args, *params.MaxPrice)
		argIndex++
	}

	if params.UserID != nil {
		query += fmt.Sprintf(` AND mic.user_id = $%d`, argIndex)
		args = append(args, *params.UserID)
		argIndex++
	}

	if params.SearchQuery != "" {
		query += fmt.Sprintf(` AND (mic.name ILIKE $%d OR mic.description ILIKE $%d)`, argIndex, argIndex)
		searchPattern := "%" + params.SearchQuery + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	// Подсчет общего количества
	countQuery := "SELECT COUNT(*) " + strings.Replace(query, query[:strings.Index(query, "FROM")], "", 1)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count listings: %w", err)
	}

	// Добавляем сортировку
	switch params.SortBy {
	case "distance":
		if params.Center != nil {
			query += " ORDER BY distance"
		} else {
			query += " ORDER BY mic.created_at"
		}
	case "price":
		query += " ORDER BY mic.price"
	case "created_at":
		query += " ORDER BY mic.created_at"
	default:
		query += " ORDER BY mic.created_at"
	}

	if params.SortOrder == "desc" {
		query += " DESC"
	} else {
		query += " ASC"
	}

	// Добавляем лимит и оффсет
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
			&listing.DisplayStrategy,
			&listing.PrivacyLevel,
			&listing.BlurRadius,
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

		// Загружаем изображения в зависимости от типа
		images, err := r.getItemImages(ctx, listing.ID, listing.ItemType)
		if err != nil {
			log.Printf("Ошибка загрузки изображений для %s %d: %v", listing.ItemType, listing.ID, err)
		}
		listing.Images = images

		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}

// getItemImages загружает изображения в зависимости от типа элемента
func (r *UnifiedGeoRepository) getItemImages(ctx context.Context, itemID int, itemType string) ([]string, error) {
	var query string

	switch itemType {
	case "marketplace_listing":
		query = `
			SELECT public_url
			FROM marketplace_images
			WHERE listing_id = $1 AND public_url IS NOT NULL AND public_url != ''
			ORDER BY is_main DESC, created_at`
	case "storefront_product":
		query = `
			SELECT spvi.image_url
			FROM storefront_product_variants spv
			JOIN storefront_product_variant_images spvi ON spv.id = spvi.variant_id
			WHERE spv.product_id = $1 AND spv.is_active = true AND spvi.image_url IS NOT NULL
			ORDER BY spv.is_default DESC, spvi.is_main DESC, spvi.display_order`
	case "storefront":
		// Для витрин пока нет изображений, возвращаем пустой массив
		return []string{}, nil
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

// RefreshMaterializedView обновляет materialized view
func (r *UnifiedGeoRepository) RefreshMaterializedView(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY map_items_cache")
	return err
}

// GetStorefrontProducts возвращает товары конкретной витрины для группированного отображения
func (r *UnifiedGeoRepository) GetStorefrontProducts(ctx context.Context, storefrontID int, limit int) ([]types.GeoListing, error) {
	query := `
		SELECT
			mic.id,
			mic.name as title,
			COALESCE(mic.description, '') as description,
			COALESCE(mic.price, 0) as price,
			COALESCE(mic.category_name, '') as category,
			mic.latitude as lat,
			mic.longitude as lng,
			COALESCE(mic.formatted_address, '') as address,
			mic.user_id,
			mic.storefront_id,
			mic.status,
			mic.created_at,
			mic.updated_at,
			mic.views_count,
			mic.rating,
			mic.item_type,
			mic.display_strategy,
			COALESCE(mic.privacy_level::text, 'exact') as privacy_level,
			COALESCE(mic.blur_radius_meters, 0) as blur_radius_meters
		FROM map_items_cache mic
		WHERE mic.storefront_id = $1
		  AND mic.item_type = 'storefront_product'
		  AND mic.status = 'active'
		ORDER BY mic.created_at DESC
		LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, storefrontID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront products: %w", err)
	}
	defer rows.Close()

	var products []types.GeoListing
	for rows.Next() {
		var product types.GeoListing
		var lat, lng float64
		var storefrontIDVal sql.NullInt64

		err = rows.Scan(
			&product.ID,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.Category,
			&lat,
			&lng,
			&product.Address,
			&product.UserID,
			&storefrontIDVal,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.ViewsCount,
			&product.Rating,
			&product.ItemType,
			&product.DisplayStrategy,
			&product.PrivacyLevel,
			&product.BlurRadius,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		product.Location = types.Point{
			Lat: lat,
			Lng: lng,
		}

		if storefrontIDVal.Valid {
			storefrontIDInt := int(storefrontIDVal.Int64)
			product.StorefrontID = &storefrontIDInt
		}

		// Загружаем изображения
		images, err := r.getItemImages(ctx, product.ID, product.ItemType)
		if err != nil {
			log.Printf("Ошибка загрузки изображений для товара %d: %v", product.ID, err)
		}
		product.Images = images

		products = append(products, product)
	}

	return products, nil
}
