package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

// SearchListings поиск объявлений в заданной области
func (r *PostGISRepository) SearchListings(ctx context.Context, params types.SearchParams) ([]types.GeoListing, int64, error) {
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
				'RSD' as currency,
				mc.name as category,
				ST_Y(lg.location::geometry) as lat,
				ST_X(lg.location::geometry) as lng,
				ml.location as address,
				ml.user_id,
				ml.status,
				ml.created_at,
				ml.updated_at`

	// Добавляем расчет расстояния если есть центр поиска
	if params.Center != nil {
		query += fmt.Sprintf(`,
				ST_Distance(
					lg.location::geography,
					ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geography
				) as distance`, params.Center.Lng, params.Center.Lat)
	}

	query += `
			FROM marketplace_listings ml
			LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
			LEFT JOIN listings_geo lg ON ml.id = lg.listing_id
			WHERE ml.status = 'active' AND lg.location IS NOT NULL`

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
			&listing.Currency,
			&listing.Category,
			&lat,
			&lng,
			&listing.Address,
			&listing.UserID,
			&listing.Status,
			&listing.CreatedAt,
			&listing.UpdatedAt,
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
		listing.Images, _ = r.getListingImages(ctx, listing.ID)

		listings = append(listings, listing)
	}

	return listings, totalCount, nil
}

// GetClusters получение кластеров для заданной области и уровня зума
func (r *PostGISRepository) GetClusters(ctx context.Context, params types.ClusterParams) ([]types.Cluster, []types.GeoListing, error) {
	// Логируем входные параметры
	log.Printf("[GetClusters] Starting with params: ZoomLevel=%d, Bounds=%+v, Categories=%v, CategoryIDs=%v, MinPrice=%v, MaxPrice=%v",
		params.ZoomLevel, params.Bounds, params.Categories, params.CategoryIDs, params.MinPrice, params.MaxPrice)

	// Определяем размер сетки на основе уровня зума
	gridSize := params.GridSize
	if gridSize == 0 {
		// Автоматический расчет размера сетки
		gridSize = calculateGridSize(params.ZoomLevel)
	}
	log.Printf("[GetClusters] Using gridSize=%d for zoomLevel=%d", gridSize, params.ZoomLevel)

	// Запрос для кластеризации
	clusterQuery := `
		WITH grid AS (
			SELECT 
				width_bucket(ST_X(lg.location::geometry), $1, $3, $5) as x_bucket,
				width_bucket(ST_Y(lg.location::geometry), $2, $4, $5) as y_bucket,
				COUNT(*) as count,
				AVG(ST_X(lg.location::geometry)) as center_lng,
				AVG(ST_Y(lg.location::geometry)) as center_lat,
				MIN(ST_X(lg.location::geometry)) as min_lng,
				MIN(ST_Y(lg.location::geometry)) as min_lat,
				MAX(ST_X(lg.location::geometry)) as max_lng,
				MAX(ST_Y(lg.location::geometry)) as max_lat
			FROM marketplace_listings ml
			LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
			LEFT JOIN listings_geo lg ON ml.id = lg.listing_id
			WHERE ml.status = 'active' AND lg.location IS NOT NULL
				AND lg.location && ST_MakeEnvelope($1, $2, $3, $4, 4326)`

	args := []interface{}{
		params.Bounds.West,
		params.Bounds.South,
		params.Bounds.East,
		params.Bounds.North,
		gridSize,
	}
	argCount := 6

	// Добавляем фильтры по категориям (по названию)
	if len(params.Categories) > 0 {
		placeholders := make([]string, len(params.Categories))
		for i, cat := range params.Categories {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, cat)
			argCount++
		}
		clusterQuery += fmt.Sprintf(" AND mc.name IN (%s)", strings.Join(placeholders, ","))
	}

	// Фильтр по категориям (по ID)
	if len(params.CategoryIDs) > 0 {
		placeholders := make([]string, len(params.CategoryIDs))
		for i, catID := range params.CategoryIDs {
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, catID)
			argCount++
		}
		clusterQuery += fmt.Sprintf(" AND mc.id IN (%s)", strings.Join(placeholders, ","))
	}

	if params.MinPrice != nil {
		clusterQuery += fmt.Sprintf(" AND ml.price >= $%d", argCount)
		args = append(args, *params.MinPrice)
		argCount++
	}

	if params.MaxPrice != nil {
		clusterQuery += fmt.Sprintf(" AND ml.price <= $%d", argCount)
		args = append(args, *params.MaxPrice)
		_ = argCount + 1 // Note: argCount is used for documentation purposes
	}

	// Убираем фильтр по валюте, так как все объявления в RSD
	// if params.Currency != "" {
	//	clusterQuery += fmt.Sprintf(" AND ml.currency = $%d", argCount)
	//	args = append(args, params.Currency)
	//	argCount++
	// }

	clusterQuery += `
			GROUP BY x_bucket, y_bucket
		)
		SELECT * FROM grid`

	// Логируем финальный запрос и аргументы
	log.Printf("[GetClusters] Executing cluster query with %d arguments", len(args))
	log.Printf("[GetClusters] Query: %s", clusterQuery)
	log.Printf("[GetClusters] Args: %v", args)

	// Выполняем запрос кластеров
	rows, err := r.db.QueryContext(ctx, clusterQuery, args...)
	if err != nil {
		log.Printf("[GetClusters] Error executing query: %v", err)
		return nil, nil, fmt.Errorf("failed to get clusters: %w", err)
	}
	defer rows.Close()

	var clusters []types.Cluster
	rowCount := 0
	for rows.Next() {
		rowCount++
		var xBucket, yBucket int
		var count int
		var centerLng, centerLat float64
		var minLng, minLat, maxLng, maxLat float64

		err := rows.Scan(&xBucket, &yBucket, &count, &centerLng, &centerLat,
			&minLng, &minLat, &maxLng, &maxLat)
		if err != nil {
			log.Printf("[GetClusters] Error scanning row %d: %v", rowCount, err)
			return nil, nil, fmt.Errorf("failed to scan cluster: %w", err)
		}

		cluster := types.Cluster{
			Center: types.Point{Lat: centerLat, Lng: centerLng},
			Count:  count,
			Bounds: types.Bounds{
				North: maxLat,
				South: minLat,
				East:  maxLng,
				West:  minLng,
			},
			ZoomExpand: params.ZoomLevel + 2, // Раскрываем на 2 уровня зума больше
		}

		log.Printf("[GetClusters] Row %d: bucket(%d,%d), count=%d, center=(%.6f,%.6f), bounds=(N:%.6f,S:%.6f,E:%.6f,W:%.6f)",
			rowCount, xBucket, yBucket, count, centerLat, centerLng, maxLat, minLat, maxLng, minLng)

		clusters = append(clusters, cluster)
	}

	log.Printf("[GetClusters] Found %d clusters from query", rowCount)

	// На низких уровнях зума (< 12) показываем все как кластеры
	// На высоких уровнях зума (>= 12) показываем одиночные объявления как отдельные маркеры
	var listings []types.GeoListing
	var finalClusters []types.Cluster

	if params.ZoomLevel >= 12 {
		log.Printf("[GetClusters] ZoomLevel >= 12, separating individual listings from clusters")

		// Отделяем одиночные объявления от кластеров
		for _, cluster := range clusters {
			if cluster.Count > 1 {
				finalClusters = append(finalClusters, cluster)
			}
		}

		// Получаем все объявления в области для показа одиночных
		searchParams := types.SearchParams{
			Bounds:      &params.Bounds,
			Categories:  params.Categories,
			CategoryIDs: params.CategoryIDs,
			MinPrice:    params.MinPrice,
			MaxPrice:    params.MaxPrice,
			Currency:    params.Currency,
			Limit:       1000, // Ограничиваем количество
		}

		allListings, _, err := r.SearchListings(ctx, searchParams)
		if err != nil {
			log.Printf("[GetClusters] Error fetching individual listings: %v", err)
			return nil, nil, fmt.Errorf("failed to get individual listings: %w", err)
		}

		// Фильтруем только те объявления, которые не входят в кластеры
		// Для этого нужно проверить, попадает ли объявление в границы какого-либо кластера
		for _, listing := range allListings {
			isInCluster := false
			for _, cluster := range finalClusters {
				if listing.Location.Lat >= cluster.Bounds.South &&
					listing.Location.Lat <= cluster.Bounds.North &&
					listing.Location.Lng >= cluster.Bounds.West &&
					listing.Location.Lng <= cluster.Bounds.East {
					isInCluster = true
					break
				}
			}
			if !isInCluster {
				listings = append(listings, listing)
			}
		}

		log.Printf("[GetClusters] Separated into %d clusters and %d individual listings", len(finalClusters), len(listings))
		clusters = finalClusters
	} else {
		log.Printf("[GetClusters] ZoomLevel < 12, showing all as clusters")
	}

	// Добавим дополнительный запрос для проверки общего количества записей
	countQuery := `
		SELECT COUNT(*) 
		FROM marketplace_listings ml
		LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
		LEFT JOIN listings_geo lg ON ml.id = lg.listing_id
		WHERE ml.status = 'active' AND lg.location IS NOT NULL
			AND lg.location && ST_MakeEnvelope($1, $2, $3, $4, 4326)`

	var totalCount int
	err = r.db.GetContext(ctx, &totalCount, countQuery, params.Bounds.West, params.Bounds.South, params.Bounds.East, params.Bounds.North)
	if err != nil {
		log.Printf("[GetClusters] Error counting total records: %v", err)
	} else {
		log.Printf("[GetClusters] Total active listings with geo data in bounds: %d", totalCount)
	}

	log.Printf("[GetClusters] Returning %d clusters and %d listings", len(clusters), len(listings))
	return clusters, listings, nil
}

// GetListingByID получение объявления по ID с геоданными
func (r *PostGISRepository) GetListingByID(ctx context.Context, id int) (*types.GeoListing, error) {
	query := `
		SELECT 
			ml.id,
			ml.title,
			ml.description,
			ml.price,
			'RSD' as currency,
			mc.name as category,
			ST_Y(lg.location::geometry) as lat,
			ST_X(lg.location::geometry) as lng,
			ml.location as address,
			ml.user_id,
			ml.status,
			ml.created_at,
			ml.updated_at
		FROM marketplace_listings ml
		LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
		LEFT JOIN listings_geo lg ON ml.id = lg.listing_id
		WHERE ml.id = $1 AND lg.location IS NOT NULL`

	var listing types.GeoListing
	var lat, lng float64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Currency,
		&listing.Category,
		&lat,
		&lng,
		&listing.Address,
		&listing.UserID,
		&listing.Status,
		&listing.CreatedAt,
		&listing.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrLocationNotFound
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	listing.Location = types.Point{Lat: lat, Lng: lng}
	listing.Images, _ = r.getListingImages(ctx, listing.ID)

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

// getListingImages получение изображений объявления
func (r *PostGISRepository) getListingImages(ctx context.Context, listingID int) ([]string, error) {
	query := `
		SELECT image_url
		FROM marketplace_images
		WHERE listing_id = $1
		ORDER BY display_order, created_at`

	var images []string
	err := r.db.SelectContext(ctx, &images, query, listingID)
	if err != nil {
		return nil, err
	}

	return images, nil
}

// calculateGridSize вычисление размера сетки для кластеризации
func calculateGridSize(zoomLevel int) int {
	// Чем выше зум, тем мельче сетка (больше buckets = меньше кластеры)
	switch {
	case zoomLevel <= 5:
		return 5 // Очень крупные кластеры
	case zoomLevel <= 8:
		return 10 // Крупные кластеры
	case zoomLevel <= 11:
		return 20 // Средние кластеры
	case zoomLevel <= 14:
		return 40 // Мелкие кластеры
	default:
		return 80 // Очень мелкие кластеры или отдельные точки
	}
}
