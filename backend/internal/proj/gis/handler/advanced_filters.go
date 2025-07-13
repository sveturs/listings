package handler

import (
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
		return utils.SendError(c, fiber.StatusBadRequest, "validation.failed")
	}

	// Валидация параметров
	if err := utils.ValidateStruct(&filter); err != nil {
		logger.Error().Err(err).Msg("Invalid travel time filter")
		return utils.SendError(c, fiber.StatusBadRequest, "validation.failed")
	}

	// Получаем изохрону
	isochrone, err := h.isochroneService.GetIsochrone(c.Context(), &filter)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get isochrone")
		return utils.SendError(c, fiber.StatusInternalServerError, "gis.isochroneError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, isochrone)
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
		return utils.SendError(c, fiber.StatusBadRequest, "validation.missingPOIType")
	}

	lat, err := c.QueryFloat("lat")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "validation.invalidLatitude")
	}

	lng, err := c.QueryFloat("lng")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "validation.invalidLongitude")
	}

	radius := c.QueryFloat("radius", 2.0)

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
		return utils.SendError(c, fiber.StatusInternalServerError, "gis.poiSearchError")
	}

	// Кешируем результаты
	for _, poi := range results {
		if err := h.poiService.CachePOI(c.Context(), poi); err != nil {
			logger.Warn().Err(err).Str("poi_id", poi.ID).Msg("Failed to cache POI")
		}
	}

	return utils.SendSuccess(c, fiber.StatusOK, results)
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
		return utils.SendError(c, fiber.StatusBadRequest, "validation.failed")
	}

	// Валидация
	if err := utils.ValidateStruct(&bbox); err != nil {
		logger.Error().Err(err).Msg("Invalid bounding box")
		return utils.SendError(c, fiber.StatusBadRequest, "validation.failed")
	}

	// Вычисляем оптимальный размер сетки
	gridSize := h.densityService.CalculateOptimalGridSize(&bbox)

	// Анализируем плотность
	results, err := h.densityService.AnalyzeDensity(c.Context(), &bbox, gridSize)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to analyze density")
		return utils.SendError(c, fiber.StatusInternalServerError, "gis.densityAnalysisError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, map[string]interface{}{
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
	bbox := types.BoundingBox{
		MinLat: c.QueryFloat("min_lat"),
		MinLng: c.QueryFloat("min_lng"),
		MaxLat: c.QueryFloat("max_lat"),
		MaxLng: c.QueryFloat("max_lng"),
	}

	// Валидация
	if err := utils.ValidateStruct(&bbox); err != nil {
		logger.Error().Err(err).Msg("Invalid bounding box")
		return utils.SendError(c, fiber.StatusBadRequest, "validation.failed")
	}

	// Получаем данные для тепловой карты
	points, err := h.densityService.GetDensityHeatmap(c.Context(), &bbox)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get density heatmap")
		return utils.SendError(c, fiber.StatusInternalServerError, "gis.heatmapError")
	}

	return utils.SendSuccess(c, fiber.StatusOK, points)
}

// ApplyAdvancedFilters применяет расширенные фильтры к списку объявлений
// @Summary Apply advanced geo filters
// @Description Apply travel time, POI, and density filters to listings
// @Tags gis-advanced
// @Accept json
// @Produce json
// @Param filters body types.AdvancedGeoFilters true "Advanced geo filters"
// @Param listing_ids body []string true "List of listing IDs to filter"
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
		return utils.SendError(c, fiber.StatusBadRequest, "validation.failed")
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
			return utils.SendError(c, fiber.StatusInternalServerError, "gis.filterError")
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
			return utils.SendError(c, fiber.StatusInternalServerError, "gis.filterError")
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
			return utils.SendError(c, fiber.StatusInternalServerError, "gis.filterError")
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

	return utils.SendSuccess(c, fiber.StatusOK, map[string]interface{}{
		"filtered_ids":    filteredIDs,
		"applied_filters": appliedFilters,
		"input_count":     len(request.ListingIDs),
		"output_count":    len(filteredIDs),
		"response_time":   responseTime,
	})
}
