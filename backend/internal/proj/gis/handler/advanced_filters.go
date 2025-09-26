package handler

import (
	"strconv"
	"time"

	"backend/internal/logger"
	"backend/internal/proj/gis/service"
	"backend/internal/proj/gis/types"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AdvancedFiltersHandler обработчик расширенных геофильтров
type AdvancedFiltersHandler struct {
	isochroneService *service.IsochroneService
	poiService       *service.POIService
	densityService   *service.DensityService
}

// NewAdvancedFiltersHandler создает новый обработчик
func NewAdvancedFiltersHandler(
	isochroneService *service.IsochroneService,
	poiService *service.POIService,
	densityService *service.DensityService,
) *AdvancedFiltersHandler {
	return &AdvancedFiltersHandler{
		isochroneService: isochroneService,
		poiService:       poiService,
		densityService:   densityService,
	}
}

// GetIsochrone получает изохрону для заданных параметров
// @Summary Get isochrone polygon
// @Description Returns isochrone polygon for travel time analysis
// @Tags gis-advanced
// @Accept json
// @Produce json
// @Param filter body types.TravelTimeFilter true "Travel time filter parameters"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.IsohronResponse}
// @Failure 400 {object} utils.ErrorResponseSwag "validation.failed"
// @Failure 500 {object} utils.ErrorResponseSwag "gis.isochroneError"
// @Router /api/v1/gis/advanced/isochrone [post]
func (h *AdvancedFiltersHandler) GetIsochrone(c *fiber.Ctx) error {
	var filter types.TravelTimeFilter
	if err := c.BodyParser(&filter); err != nil {
		logger.Error().Err(err).Msg("Failed to parse travel time filter")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.failed")
	}

	// TODO: Add validation if needed

	// Получаем изохрону
	isochrone, err := h.isochroneService.GetIsochrone(c.Context(), &filter)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get isochrone")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.isochroneError")
	}

	return utils.SuccessResponse(c, isochrone)
}

// SearchPOI поиск точек интереса
// @Summary Search points of interest
// @Description Search for POIs like schools, hospitals, metro stations
// @Tags gis-advanced
// @Accept json
// @Produce json
// @Query poi_type string true "Type of POI" Enums(school,hospital,metro,supermarket,park,bank,pharmacy,bus_stop)
// @Query lat number true "Latitude"
// @Query lng number true "Longitude"
// @Query radius number false "Search radius in km" default(2)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]types.POISearchResult}
// @Failure 400 {object} utils.ErrorResponseSwag "validation.failed"
// @Failure 500 {object} utils.ErrorResponseSwag "gis.poiSearchError"
// @Router /api/v1/gis/advanced/poi/search [get]
func (h *AdvancedFiltersHandler) SearchPOI(c *fiber.Ctx) error {
	poiTypeStr := c.Query("poi_type")
	if poiTypeStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.missingPOIType")
	}

	latStr := c.Query("lat")
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidLatitude")
	}

	lngStr := c.Query("lng")
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidLongitude")
	}

	radiusStr := c.Query("radius", "2.0")
	radius, _ := strconv.ParseFloat(radiusStr, 64)

	// Поиск POI
	results, err := h.poiService.SearchPOI(
		c.Context(),
		types.POIType(poiTypeStr),
		lat,
		lng,
		radius,
	)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to search POI")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.poiSearchError")
	}

	// Кешируем результаты
	for _, poi := range results {
		if err := h.poiService.CachePOI(c.Context(), poi); err != nil {
			logger.Warn().Err(err).Str("poi_id", poi.ID).Msg("Failed to cache POI")
		}
	}

	return utils.SuccessResponse(c, results)
}

// AnalyzeDensity анализ плотности объявлений
// @Summary Analyze listing density
// @Description Analyze listing density in a given area
// @Tags gis-advanced
// @Accept json
// @Produce json
// @Param bbox body types.BoundingBox true "Bounding box for analysis"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]types.DensityAnalysisResult}
// @Failure 400 {object} utils.ErrorResponseSwag "validation.failed"
// @Failure 500 {object} utils.ErrorResponseSwag "gis.densityAnalysisError"
// @Router /api/v1/gis/advanced/density/analyze [post]
func (h *AdvancedFiltersHandler) AnalyzeDensity(c *fiber.Ctx) error {
	var bbox types.BoundingBox
	if err := c.BodyParser(&bbox); err != nil {
		logger.Error().Err(err).Msg("Failed to parse bounding box")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.failed")
	}

	// TODO: Add validation if needed

	// Вычисляем оптимальный размер сетки
	gridSize := h.densityService.CalculateOptimalGridSize(&bbox)

	// Анализируем плотность
	results, err := h.densityService.AnalyzeDensity(c.Context(), &bbox, gridSize)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to analyze density")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.densityAnalysisError")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"grid_size": gridSize,
		"results":   results,
	})
}

// GetDensityHeatmap получение данных для тепловой карты
// @Summary Get density heatmap data
// @Description Get point data for density heatmap visualization
// @Tags gis-advanced
// @Accept json
// @Produce json
// @Query min_lat number true "Minimum latitude"
// @Query min_lng number true "Minimum longitude"
// @Query max_lat number true "Maximum latitude"
// @Query max_lng number true "Maximum longitude"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]map[string]interface{}}
// @Failure 400 {object} utils.ErrorResponseSwag "validation.failed"
// @Failure 500 {object} utils.ErrorResponseSwag "gis.heatmapError"
// @Router /api/v1/gis/advanced/density/heatmap [get]
func (h *AdvancedFiltersHandler) GetDensityHeatmap(c *fiber.Ctx) error {
	minLatStr := c.Query("min_lat")
	minLat, _ := strconv.ParseFloat(minLatStr, 64)
	minLngStr := c.Query("min_lng")
	minLng, _ := strconv.ParseFloat(minLngStr, 64)
	maxLatStr := c.Query("max_lat")
	maxLat, _ := strconv.ParseFloat(maxLatStr, 64)
	maxLngStr := c.Query("max_lng")
	maxLng, _ := strconv.ParseFloat(maxLngStr, 64)

	bbox := types.BoundingBox{
		MinLat: minLat,
		MinLng: minLng,
		MaxLat: maxLat,
		MaxLng: maxLng,
	}

	// TODO: Add validation if needed

	// Получаем данные для тепловой карты
	points, err := h.densityService.GetDensityHeatmap(c.Context(), &bbox)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get density heatmap")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.heatmapError")
	}

	return utils.SuccessResponse(c, points)
}

// ApplyAdvancedFilters применяет расширенные фильтры к списку объявлений
// @Summary Apply advanced geo filters
// @Description Apply travel time, POI, and density filters to listings
// @Tags gis-advanced
// @Accept json
// @Produce json
// @Param request body types.ApplyAdvancedFiltersRequest true "Advanced filters request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}}
// @Failure 400 {object} utils.ErrorResponseSwag "validation.failed"
// @Failure 500 {object} utils.ErrorResponseSwag "gis.filterError"
// @Router /api/v1/gis/advanced/apply-filters [post]
func (h *AdvancedFiltersHandler) ApplyAdvancedFilters(c *fiber.Ctx) error {
	var request struct {
		Filters    types.AdvancedGeoFilters `json:"filters"`
		ListingIDs []string                 `json:"listing_ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.failed")
	}

	startTime := time.Now()
	filteredIDs := request.ListingIDs
	appliedFilters := []string{}

	// Применяем фильтр по времени пути
	if request.Filters.TravelTime != nil {
		ids, err := h.isochroneService.FilterListingsByIsochrone(
			c.Context(),
			request.Filters.TravelTime,
			filteredIDs,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to apply travel time filter")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.filterError")
		}
		filteredIDs = ids
		appliedFilters = append(appliedFilters, "travel_time")
	}

	// Применяем фильтр по POI
	if request.Filters.POIFilter != nil {
		ids, err := h.poiService.FilterListingsByPOI(
			c.Context(),
			request.Filters.POIFilter,
			filteredIDs,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to apply POI filter")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.filterError")
		}
		filteredIDs = ids
		appliedFilters = append(appliedFilters, "poi")
	}

	// Применяем фильтр по плотности
	if request.Filters.DensityFilter != nil {
		ids, err := h.densityService.FilterListingsByDensity(
			c.Context(),
			request.Filters.DensityFilter,
			filteredIDs,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to apply density filter")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.filterError")
		}
		filteredIDs = ids
		appliedFilters = append(appliedFilters, "density")
	}

	responseTime := time.Since(startTime).Milliseconds()

	// Логируем аналитику
	analytics := types.FilterAnalytics{
		SessionID:      c.Get("X-Session-ID", ""),
		FilterType:     "advanced_geo",
		FilterParams:   request.Filters,
		ResultCount:    len(filteredIDs),
		ResponseTimeMs: responseTime,
		CreatedAt:      time.Now(),
	}

	// TODO: Сохранить аналитику в базу данных
	logger.Info().
		Str("session_id", analytics.SessionID).
		Strs("applied_filters", appliedFilters).
		Int("input_count", len(request.ListingIDs)).
		Int("output_count", len(filteredIDs)).
		Int64("response_time_ms", responseTime).
		Msg("Advanced geo filters applied")

	return utils.SuccessResponse(c, map[string]interface{}{
		"filtered_ids":    filteredIDs,
		"applied_filters": appliedFilters,
		"input_count":     len(request.ListingIDs),
		"output_count":    len(filteredIDs),
		"response_time":   responseTime,
	})
}
