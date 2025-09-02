// backend/internal/proj/marketplace/service/map.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"backend/internal/domain/models"
)

// GetListingsInBounds возвращает маркеры объявлений в указанных границах
func (s *MarketplaceService) GetListingsInBounds(ctx context.Context, neLat, neLng, swLat, swLng float64, zoom int, categoryIDs, condition string, minPrice, maxPrice *float64, attributesFilter string) ([]models.MapMarker, error) {
	// Строим SQL запрос для получения объявлений в границах
	query := `
		SELECT DISTINCT
			ml.id,
			ml.latitude,
			ml.longitude,
			ml.title,
			ml.price,
			COALESCE(ml.condition, '') as condition,
			ml.category_id,
			ml.user_id,
			COALESCE(ml.address_city, '') as city,
			COALESCE(ml.address_country, '') as country,
			ml.created_at,
			ml.views_count,
			COALESCE(rc.average_rating, 0) as rating,
			COALESCE(mi.file_path, '') as main_image
		FROM marketplace_listings ml
		LEFT JOIN marketplace_images mi ON ml.id = mi.listing_id AND mi.is_main = true
		LEFT JOIN rating_cache rc ON rc.entity_type = 'listing' AND rc.entity_id = ml.id
		WHERE ml.status = 'active'
		AND ml.show_on_map = true
		AND ml.latitude IS NOT NULL
		AND ml.longitude IS NOT NULL
		AND ml.latitude BETWEEN $1 AND $2
		AND ml.longitude BETWEEN $3 AND $4
	`

	args := []interface{}{swLat, neLat, swLng, neLng}
	argCounter := 5

	// Добавляем фильтры категорий
	if categoryIDs != "" && categoryIDs != "0" {
		categoryIDList := strings.Split(categoryIDs, ",")
		if len(categoryIDList) > 0 {
			placeholders := make([]string, len(categoryIDList))
			for i, catID := range categoryIDList {
				placeholders[i] = fmt.Sprintf("$%d", argCounter)
				args = append(args, catID)
				argCounter++
			}
			query += fmt.Sprintf(" AND ml.category_id IN (%s)", strings.Join(placeholders, ","))
		}
	}

	// Добавляем фильтр состояния
	if condition != "" && condition != "any" {
		query += fmt.Sprintf(" AND ml.condition = $%d", argCounter)
		args = append(args, condition)
		argCounter++
	}

	// Добавляем фильтры цены
	if minPrice != nil {
		query += fmt.Sprintf(" AND ml.price >= $%d", argCounter)
		args = append(args, *minPrice)
		argCounter++
	}
	if maxPrice != nil {
		query += fmt.Sprintf(" AND ml.price <= $%d", argCounter)
		args = append(args, *maxPrice)
		argCounter++
	}

	// Добавляем фильтры по атрибутам
	if attributesFilter != "" {
		// Парсим JSON с фильтрами атрибутов
		var filters map[string]interface{}
		if err := json.Unmarshal([]byte(attributesFilter), &filters); err == nil {
			for attrID, value := range filters {
				// Добавляем JOIN для каждого атрибута
				query += fmt.Sprintf(`
					AND EXISTS (
						SELECT 1 FROM listing_attribute_values lav
						WHERE lav.listing_id = ml.id
						AND lav.attribute_id = $%d
				`, argCounter)
				args = append(args, attrID)
				argCounter++

				// В зависимости от типа значения добавляем соответствующее условие
				switch v := value.(type) {
				case string:
					query += fmt.Sprintf(" AND lav.text_value = $%d", argCounter)
					args = append(args, v)
					argCounter++
				case float64:
					query += fmt.Sprintf(" AND lav.numeric_value = $%d", argCounter)
					args = append(args, v)
					argCounter++
				case bool:
					query += fmt.Sprintf(" AND lav.boolean_value = $%d", argCounter)
					args = append(args, v)
					argCounter++
				case map[string]interface{}:
					// Для range фильтров
					if min, ok := v["min"].(float64); ok {
						query += fmt.Sprintf(" AND lav.numeric_value >= $%d", argCounter)
						args = append(args, min)
						argCounter++
					}
					if max, ok := v["max"].(float64); ok {
						query += fmt.Sprintf(" AND lav.numeric_value <= $%d", argCounter)
						args = append(args, max)
						argCounter++
					}
				case []interface{}:
					// Для multiselect
					if len(v) > 0 {
						values := make([]string, len(v))
						for i, item := range v {
							if str, ok := item.(string); ok {
								values[i] = str
							}
						}
						if len(values) > 0 {
							query += fmt.Sprintf(" AND lav.json_value ?| ARRAY[%s]", strings.Join(values, ","))
						}
					}
				}

				query += ")"
			}
		}
	}

	// Для больших масштабов ограничиваем количество результатов
	limit := s.calculateLimit(zoom)
	query += fmt.Sprintf(" ORDER BY ml.created_at DESC LIMIT %d", limit)

	rows, err := s.storage.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query listings in bounds: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var markers []models.MapMarker
	for rows.Next() {
		var marker models.MapMarker
		var createdAt interface{}

		err := rows.Scan(
			&marker.ID,
			&marker.Latitude,
			&marker.Longitude,
			&marker.Title,
			&marker.Price,
			&marker.Condition,
			&marker.CategoryID,
			&marker.UserID,
			&marker.City,
			&marker.Country,
			&createdAt,
			&marker.ViewsCount,
			&marker.Rating,
			&marker.MainImage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan marker: %w", err)
		}

		// Форматируем дату
		if createdAt != nil {
			marker.CreatedAt = fmt.Sprintf("%v", createdAt)
		}

		markers = append(markers, marker)
	}

	return markers, nil
}

// GetMapClusters возвращает кластеризованные данные для карты
func (s *MarketplaceService) GetMapClusters(ctx context.Context, neLat, neLng, swLat, swLng float64, zoom int, categoryIDs, condition string, minPrice, maxPrice *float64, attributesFilter string) ([]models.MapCluster, error) {
	// Для серверной кластеризации используем простую сетку
	gridSize := s.calculateGridSize(zoom)

	query := `
		SELECT
			FLOOR(ml.latitude / $5) * $5 as cluster_lat,
			FLOOR(ml.longitude / $6) * $6 as cluster_lng,
			COUNT(*) as count,
			AVG(ml.price) as avg_price,
			MIN(ml.latitude) as min_lat,
			MAX(ml.latitude) as max_lat,
			MIN(ml.longitude) as min_lng,
			MAX(ml.longitude) as max_lng,
			array_agg(DISTINCT ml.category_id) as categories
		FROM marketplace_listings ml
		WHERE ml.status = 'active'
		AND ml.show_on_map = true
		AND ml.latitude IS NOT NULL
		AND ml.longitude IS NOT NULL
		AND ml.latitude BETWEEN $1 AND $2
		AND ml.longitude BETWEEN $3 AND $4
	`

	args := []interface{}{swLat, neLat, swLng, neLng, gridSize, gridSize}
	argCounter := 7

	// Добавляем фильтры (аналогично GetListingsInBounds)
	if categoryIDs != "" && categoryIDs != "0" {
		categoryIDList := strings.Split(categoryIDs, ",")
		if len(categoryIDList) > 0 {
			placeholders := make([]string, len(categoryIDList))
			for i, catID := range categoryIDList {
				placeholders[i] = fmt.Sprintf("$%d", argCounter)
				args = append(args, catID)
				argCounter++
			}
			query += fmt.Sprintf(" AND ml.category_id IN (%s)", strings.Join(placeholders, ","))
		}
	}

	if condition != "" && condition != "any" {
		query += fmt.Sprintf(" AND ml.condition = $%d", argCounter)
		args = append(args, condition)
		argCounter++
	}

	if minPrice != nil {
		query += fmt.Sprintf(" AND ml.price >= $%d", argCounter)
		args = append(args, *minPrice)
		argCounter++
	}
	if maxPrice != nil {
		query += fmt.Sprintf(" AND ml.price <= $%d", argCounter)
		args = append(args, *maxPrice)
		argCounter++
	}

	// Добавляем фильтры по атрибутам
	if attributesFilter != "" {
		// Парсим JSON с фильтрами атрибутов
		var filters map[string]interface{}
		if err := json.Unmarshal([]byte(attributesFilter), &filters); err == nil {
			for attrID, value := range filters {
				// Добавляем JOIN для каждого атрибута
				query += fmt.Sprintf(`
					AND EXISTS (
						SELECT 1 FROM listing_attribute_values lav
						WHERE lav.listing_id = ml.id
						AND lav.attribute_id = $%d
				`, argCounter)
				args = append(args, attrID)
				argCounter++

				// В зависимости от типа значения добавляем соответствующее условие
				switch v := value.(type) {
				case string:
					query += fmt.Sprintf(" AND lav.text_value = $%d", argCounter)
					args = append(args, v)
					argCounter++
				case float64:
					query += fmt.Sprintf(" AND lav.numeric_value = $%d", argCounter)
					args = append(args, v)
					argCounter++
				case bool:
					query += fmt.Sprintf(" AND lav.boolean_value = $%d", argCounter)
					args = append(args, v)
					argCounter++
				case map[string]interface{}:
					// Для range фильтров
					if min, ok := v["min"].(float64); ok {
						query += fmt.Sprintf(" AND lav.numeric_value >= $%d", argCounter)
						args = append(args, min)
						argCounter++
					}
					if max, ok := v["max"].(float64); ok {
						query += fmt.Sprintf(" AND lav.numeric_value <= $%d", argCounter)
						args = append(args, max)
						argCounter++
					}
				case []interface{}:
					// Для multiselect
					if len(v) > 0 {
						values := make([]string, len(v))
						for i, item := range v {
							if str, ok := item.(string); ok {
								values[i] = str
							}
						}
						if len(values) > 0 {
							query += fmt.Sprintf(" AND lav.json_value ?| ARRAY[%s]", strings.Join(values, ","))
						}
					}
				}

				query += ")"
			}
		}
	}

	query += " GROUP BY cluster_lat, cluster_lng HAVING COUNT(*) > 1"

	rows, err := s.storage.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clusters: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var clusters []models.MapCluster
	for rows.Next() {
		var cluster models.MapCluster
		var categoriesStr string
		var minLat, maxLat, minLng, maxLng float64

		err := rows.Scan(
			&cluster.Latitude,
			&cluster.Longitude,
			&cluster.Count,
			&cluster.AvgPrice,
			&minLat,
			&maxLat,
			&minLng,
			&maxLng,
			&categoriesStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cluster: %w", err)
		}

		// Создаем уникальный ID для кластера
		cluster.ID = fmt.Sprintf("%f_%f", cluster.Latitude, cluster.Longitude)

		// Устанавливаем границы кластера
		cluster.Bounds = models.MapBounds{
			NorthEast: models.MapPoint{Latitude: maxLat, Longitude: maxLng},
			SouthWest: models.MapPoint{Latitude: minLat, Longitude: minLng},
		}

		// Парсим категории
		if categoriesStr != "" {
			categoriesStr = strings.Trim(categoriesStr, "{}")
			if categoriesStr != "" {
				catStrings := strings.Split(categoriesStr, ",")
				for _, catStr := range catStrings {
					if catID, err := strconv.Atoi(strings.TrimSpace(catStr)); err == nil {
						cluster.Categories = append(cluster.Categories, catID)
					}
				}
			}
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

// calculateLimit возвращает лимит результатов в зависимости от zoom
func (s *MarketplaceService) calculateLimit(zoom int) int {
	switch {
	case zoom >= 17:
		return 50 // Очень близко - показываем до 50 объявлений
	case zoom >= 15:
		return 100 // Близко - показываем до 100 объявлений
	case zoom >= 13:
		return 200 // Средне - показываем до 200 объявлений
	case zoom >= 10:
		return 500 // Далеко - показываем до 500 объявлений
	default:
		return 1000 // Очень далеко - показываем до 1000 объявлений
	}
}

// calculateGridSize возвращает размер сетки для кластеризации в зависимости от zoom
func (s *MarketplaceService) calculateGridSize(zoom int) float64 {
	// Чем больше zoom, тем меньше размер сетки (больше деталей)
	baseGridSize := 0.1 // Базовый размер сетки в градусах

	switch {
	case zoom >= 15:
		return baseGridSize / 16 // 0.00625 градуса
	case zoom >= 13:
		return baseGridSize / 8 // 0.0125 градуса
	case zoom >= 11:
		return baseGridSize / 4 // 0.025 градуса
	case zoom >= 9:
		return baseGridSize / 2 // 0.05 градуса
	default:
		return baseGridSize // 0.1 градуса
	}
}
