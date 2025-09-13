package handler

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"backend/internal/proj/gis/constants"
	"backend/pkg/utils"
)

// ClusterPoint представляет точку кластера на карте
type ClusterPoint struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Count int     `json:"count"`
	IDs   []int   `json:"ids,omitempty"`
}

// ClusterHandler обработчик для кластеризации
type ClusterHandler struct {
	db *sqlx.DB
}

// NewClusterHandler создает новый обработчик кластеров
func NewClusterHandler(db *sqlx.DB) *ClusterHandler {
	return &ClusterHandler{
		db: db,
	}
}

// GetClusters получение кластеров объявлений
// @Summary Получить кластеры объявлений
// @Description Возвращает кластеризованные точки для отображения на карте
// @Tags gis
// @Accept json
// @Produce json
// @Param zoom query integer true "Уровень зума карты (1-20)"
// @Param bounds query string true "Границы видимой области (south,west,north,east)"
// @Param category_id query integer false "ID категории для фильтрации"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]ClusterPoint} "Массив кластеров"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/clusters [get]
func (h *ClusterHandler) GetClusters(c *fiber.Ctx) error {
	// Получаем параметры запроса
	zoom := c.QueryInt("zoom", 10)
	bounds := c.Query("bounds")

	if bounds == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.boundsRequired")
	}

	// Парсим границы
	south, west, north, east, err := parseBounds(bounds)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidBounds")
	}

	// Определяем размер сетки в зависимости от зума
	var gridSize float64
	switch {
	case zoom < 8:
		gridSize = 0.5 // 50км сетка
	case zoom < 12:
		gridSize = 0.1 // 10км сетка
	case zoom < 15:
		gridSize = 0.01 // 1км сетка
	default:
		gridSize = 0.001 // 100м сетка
	}

	// Строим запрос с учетом фильтров
	query := `
		WITH filtered_listings AS (
			SELECT
				lg.listing_id,
				lg.location,
				lg.blurred_location,
				lg.is_precise,
				ml.price,
				ml.category_id
			FROM listings_geo lg
			JOIN marketplace_listings ml ON lg.listing_id = ml.id
			WHERE lg.location && ST_MakeEnvelope($1, $2, $3, $4, 4326)
				AND ml.status = 'active'
	`

	args := []interface{}{west, south, east, north}
	argCount := 4

	// Добавляем фильтр по категории
	if categoryID := c.QueryInt("category_id", 0); categoryID > 0 {
		argCount++
		query += ` AND ml.category_id = $` + strconv.Itoa(argCount)
		args = append(args, categoryID)
	}

	// Добавляем фильтр по цене
	if minPrice := c.QueryFloat("min_price", 0); minPrice > 0 {
		argCount++
		query += ` AND ml.price >= $` + strconv.Itoa(argCount)
		args = append(args, minPrice)
	}

	if maxPrice := c.QueryFloat("max_price", 0); maxPrice > 0 {
		argCount++
		query += ` AND ml.price <= $` + strconv.Itoa(argCount)
		args = append(args, maxPrice)
	}

	// Завершаем CTE и выполняем кластеризацию
	argCount++
	query += `
		),
		clusters AS (
			SELECT
				round(ST_X(ST_Centroid(
					CASE
						WHEN is_precise THEN location::geometry
						ELSE blurred_location::geometry
					END
				))::numeric / $` + strconv.Itoa(argCount) + `) * $` + strconv.Itoa(argCount) + ` as cluster_lng,
				round(ST_Y(ST_Centroid(
					CASE
						WHEN is_precise THEN location::geometry
						ELSE blurred_location::geometry
					END
				))::numeric / $` + strconv.Itoa(argCount) + `) * $` + strconv.Itoa(argCount) + ` as cluster_lat,
				COUNT(*) as count,
				array_agg(listing_id) as ids
			FROM filtered_listings
			GROUP BY cluster_lng, cluster_lat
			HAVING COUNT(*) > 0
		)
		SELECT
			cluster_lat as lat,
			cluster_lng as lng,
			count,
			CASE
				WHEN $` + strconv.Itoa(argCount+1) + ` > 15 THEN ids
				ELSE NULL
			END as ids
		FROM clusters
		ORDER BY count DESC
		LIMIT $` + strconv.Itoa(argCount+2)

	args = append(args, gridSize, zoom, constants.MAX_LIMIT)

	// Выполняем запрос
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.clusterQueryError")
	}
	defer rows.Close()

	// Собираем результаты
	clusters := []ClusterPoint{}
	for rows.Next() {
		var cluster ClusterPoint
		var idsJSON []byte

		err := rows.Scan(&cluster.Lat, &cluster.Lng, &cluster.Count, &idsJSON)
		if err != nil {
			continue
		}

		// Парсим IDs если они есть (только на больших зумах)
		if idsJSON != nil && len(idsJSON) > 0 {
			json.Unmarshal(idsJSON, &cluster.IDs)
		}

		clusters = append(clusters, cluster)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"clusters": clusters,
		"total":    len(clusters),
		"zoom":     zoom,
		"gridSize": gridSize,
	})
}

// GetHeatmap получение данных для тепловой карты
// @Summary Получить данные для тепловой карты
// @Description Возвращает точки с весами для построения тепловой карты
// @Tags gis
// @Accept json
// @Produce json
// @Param bounds query string true "Границы видимой области (south,west,north,east)"
// @Param metric query string false "Метрика для веса (price, views, density)" default(density)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]map[string]interface{}} "Данные для тепловой карты"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/heatmap [get]
func (h *ClusterHandler) GetHeatmap(c *fiber.Ctx) error {
	bounds := c.Query("bounds")
	if bounds == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.boundsRequired")
	}

	// Парсим границы
	south, west, north, east, err := parseBounds(bounds)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidBounds")
	}

	metric := c.Query("metric", "density")

	// Определяем запрос в зависимости от метрики
	var query string
	switch metric {
	case "price":
		query = `
			SELECT
				ST_X(lg.location::geometry) as lng,
				ST_Y(lg.location::geometry) as lat,
				ml.price as weight
			FROM listings_geo lg
			JOIN marketplace_listings ml ON lg.listing_id = ml.id
			WHERE lg.location && ST_MakeEnvelope($1, $2, $3, $4, 4326)
				AND ml.status = 'active'
				AND ml.price IS NOT NULL
			LIMIT 1000
		`
	case "views":
		query = `
			SELECT
				ST_X(lg.location::geometry) as lng,
				ST_Y(lg.location::geometry) as lat,
				COALESCE(ml.views_count, 1) as weight
			FROM listings_geo lg
			JOIN marketplace_listings ml ON lg.listing_id = ml.id
			WHERE lg.location && ST_MakeEnvelope($1, $2, $3, $4, 4326)
				AND ml.status = 'active'
			LIMIT 1000
		`
	default: // density
		query = `
			SELECT
				ST_X(lg.location::geometry) as lng,
				ST_Y(lg.location::geometry) as lat,
				1.0 as weight
			FROM listings_geo lg
			JOIN marketplace_listings ml ON lg.listing_id = ml.id
			WHERE lg.location && ST_MakeEnvelope($1, $2, $3, $4, 4326)
				AND ml.status = 'active'
			LIMIT 1000
		`
	}

	// Выполняем запрос
	rows, err := h.db.Query(query, west, south, east, north)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.heatmapQueryError")
	}
	defer rows.Close()

	// Собираем результаты
	heatmapData := []map[string]interface{}{}
	for rows.Next() {
		var lat, lng, weight float64
		err := rows.Scan(&lng, &lat, &weight)
		if err != nil {
			continue
		}

		heatmapData = append(heatmapData, map[string]interface{}{
			"lat":    lat,
			"lng":    lng,
			"weight": weight,
		})
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data":   heatmapData,
		"metric": metric,
		"count":  len(heatmapData),
	})
}

// parseBounds парсит строку границ в формате "south,west,north,east"
func parseBounds(boundsStr string) (south, west, north, east float64, err error) {
	parts := strings.Split(boundsStr, ",")
	if len(parts) != 4 {
		return 0, 0, 0, 0, fiber.ErrBadRequest
	}

	south, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	west, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	north, err = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	east, err = strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return south, west, north, east, nil
}