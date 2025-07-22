package service

import (
	"context"
	"fmt"
	"math"

	"backend/internal/logger"
	"backend/internal/proj/gis/types"

	"github.com/jmoiron/sqlx"
)

// DensityService сервис для анализа плотности объявлений
type DensityService struct {
	db *sqlx.DB
}

// NewDensityService создает новый сервис анализа плотности
func NewDensityService(db *sqlx.DB) *DensityService {
	return &DensityService{
		db: db,
	}
}

// AnalyzeDensity анализирует плотность объявлений в заданной области
func (s *DensityService) AnalyzeDensity(ctx context.Context, bbox *types.BoundingBox, gridSize float64) ([]types.DensityAnalysisResult, error) {
	// Создаем сетку для анализа
	query := `
		WITH grid AS (
			SELECT 
				row_number() OVER () as id,
				x,
				y,
				ST_MakeEnvelope(
					x, 
					y, 
					x + $5, 
					y + $5,
					4326
				) as cell
			FROM generate_series(
				$1::numeric, 
				$2::numeric, 
				$5::numeric
			) as x,
			generate_series(
				$3::numeric, 
				$4::numeric, 
				$5::numeric
			) as y
		),
		density_data AS (
			SELECT 
				g.id as grid_cell_id,
				ST_X(ST_Centroid(g.cell)) as center_lng,
				ST_Y(ST_Centroid(g.cell)) as center_lat,
				COUNT(l.id) as listing_count,
				ST_Area(g.cell::geography) / 1000000.0 as area_km2
			FROM grid g
			LEFT JOIN marketplace_listings l ON 
				l.location IS NOT NULL 
				AND ST_Within(l.location::geometry, g.cell)
				AND l.status = 'active'
			GROUP BY g.id, g.cell
		)
		SELECT 
			grid_cell_id::text,
			center_lat,
			center_lng,
			listing_count,
			area_km2,
			CASE 
				WHEN area_km2 > 0 THEN listing_count / area_km2
				ELSE 0
			END as density
		FROM density_data
		ORDER BY density DESC
	`

	var results []types.DensityAnalysisResult
	err := s.db.SelectContext(ctx, &results, query,
		bbox.MinLng,
		bbox.MaxLng,
		bbox.MinLat,
		bbox.MaxLat,
		gridSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze density: %w", err)
	}

	return results, nil
}

// FilterListingsByDensity фильтрует объявления по критериям плотности
func (s *DensityService) FilterListingsByDensity(ctx context.Context, filter *types.DensityFilter, listingIDs []string) ([]string, error) {
	// Сначала анализируем плотность для каждого объявления
	query := `
		WITH listing_density AS (
			SELECT 
				l1.id,
				l1.location,
				COUNT(l2.id) as nearby_count,
				ST_Area(
					ST_MakeEnvelope(
						ST_X(l1.location::geometry) - 0.005,
						ST_Y(l1.location::geometry) - 0.005,
						ST_X(l1.location::geometry) + 0.005,
						ST_Y(l1.location::geometry) + 0.005,
						4326
					)::geography
				) / 1000000.0 as area_km2
			FROM marketplace_listings l1
			LEFT JOIN marketplace_listings l2 ON 
				l2.id != l1.id
				AND l2.status = 'active'
				AND l2.location IS NOT NULL
				AND ST_DWithin(
					l1.location::geography, 
					l2.location::geography, 
					500 -- радиус 500м для подсчета плотности
				)
			WHERE 
				l1.id = ANY($1)
				AND l1.location IS NOT NULL
			GROUP BY l1.id, l1.location
		),
		density_calc AS (
			SELECT 
				id,
				nearby_count,
				area_km2,
				CASE 
					WHEN area_km2 > 0 THEN nearby_count / area_km2
					ELSE 0
				END as density
			FROM listing_density
		)
		SELECT id
		FROM density_calc
		WHERE 1=1
	`

	// Добавляем условия фильтрации
	conditions := []string{}

	if filter.AvoidCrowded {
		// Определяем "переполненные" как верхний квартиль плотности
		conditions = append(conditions, "density < (SELECT percentile_cont(0.75) WITHIN GROUP (ORDER BY density) FROM density_calc)")
	}

	if filter.MaxDensity > 0 {
		conditions = append(conditions, fmt.Sprintf("density <= %f", filter.MaxDensity))
	}

	if filter.MinDensity > 0 {
		conditions = append(conditions, fmt.Sprintf("density >= %f", filter.MinDensity))
	}

	for _, condition := range conditions {
		query += " AND " + condition
	}

	var filteredIDs []string
	err := s.db.SelectContext(ctx, &filteredIDs, query, listingIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to filter listings by density: %w", err)
	}

	return filteredIDs, nil
}

// GetDensityHeatmap возвращает данные для тепловой карты плотности
func (s *DensityService) GetDensityHeatmap(ctx context.Context, bbox *types.BoundingBox) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			ST_X(location::geometry) as lng,
			ST_Y(location::geometry) as lat,
			1.0 as weight
		FROM marketplace_listings
		WHERE 
			location IS NOT NULL
			AND status = 'active'
			AND ST_Within(
				location::geometry,
				ST_MakeEnvelope($1, $2, $3, $4, 4326)
			)
	`

	rows, err := s.db.QueryContext(ctx, query,
		bbox.MinLng,
		bbox.MinLat,
		bbox.MaxLng,
		bbox.MaxLat,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get density heatmap: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var points []map[string]interface{}
	for rows.Next() {
		var lat, lng, weight float64
		if err := rows.Scan(&lng, &lat, &weight); err != nil {
			continue
		}

		points = append(points, map[string]interface{}{
			"lat":    lat,
			"lng":    lng,
			"weight": weight,
		})
	}

	return points, nil
}

// CalculateOptimalGridSize вычисляет оптимальный размер сетки для анализа
func (s *DensityService) CalculateOptimalGridSize(bbox *types.BoundingBox) float64 {
	// Вычисляем размеры области
	width := math.Abs(bbox.MaxLng - bbox.MinLng)
	height := math.Abs(bbox.MaxLat - bbox.MinLat)

	// Целевое количество ячеек сетки (10x10)
	targetCells := 100.0

	// Вычисляем размер ячейки
	area := width * height
	cellSize := math.Sqrt(area / targetCells)

	// Ограничиваем размер ячейки
	minSize := 0.001 // ~100м
	maxSize := 0.01  // ~1км

	if cellSize < minSize {
		cellSize = minSize
	} else if cellSize > maxSize {
		cellSize = maxSize
	}

	return cellSize
}
