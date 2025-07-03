// backend/internal/proj/marketplace/handler/map.go
package handler

import (
	"strconv"

	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetListingsInBounds returns listings within specified map bounds
// @Summary Get listings in bounds
// @Description Returns all listings within the specified geographical bounds
// @Tags marketplace-map
// @Accept json
// @Produce json
// @Param ne_lat query number true "Northeast latitude"
// @Param ne_lng query number true "Northeast longitude"
// @Param sw_lat query number true "Southwest latitude"
// @Param sw_lng query number true "Southwest longitude"
// @Param zoom query int false "Map zoom level" default(10)
// @Param categories query string false "Comma-separated category IDs"
// @Param condition query string false "Item condition filter"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {object} MapBoundsResponse "Listings within bounds"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidBounds"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.mapError"
// @Router /api/v1/marketplace/map/bounds [get]
func (h *MarketplaceHandler) GetListingsInBounds(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.missingBounds")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
	}

	zoom, err := strconv.Atoi(zoomStr)
	if err != nil {
		zoom = 10
	}

	// Получаем фильтры
	categoryIDs := c.Query("categories", "")
	condition := c.Query("condition", "")
	minPrice := c.Query("min_price", "")
	maxPrice := c.Query("max_price", "")

	// Парсим цены
	var minPriceFloat, maxPriceFloat *float64
	if minPrice != "" {
		if parsed, err := strconv.ParseFloat(minPrice, 64); err == nil {
			minPriceFloat = &parsed
		}
	}
	if maxPrice != "" {
		if parsed, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			maxPriceFloat = &parsed
		}
	}

	// Получаем объявления в указанных границах
	listings, err := h.service.GetListingsInBounds(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
	}

	response := MapBoundsResponse{
		Success: true,
		Data: MapBoundsData{
			Listings: listings,
			Bounds: MapBounds{
				NE: Coordinates{Lat: neLat64, Lng: neLng64},
				SW: Coordinates{Lat: swLat64, Lng: swLng64},
			},
			Zoom:  zoom,
			Count: len(listings),
		},
	}
	return c.JSON(response)
}

// GetMapClusters returns clustered data for map view
// @Summary Get map clusters
// @Description Returns clustered listings data for efficient map rendering
// @Tags marketplace-map
// @Accept json
// @Produce json
// @Param ne_lat query number true "Northeast latitude"
// @Param ne_lng query number true "Northeast longitude"
// @Param sw_lat query number true "Southwest latitude"
// @Param sw_lng query number true "Southwest longitude"
// @Param zoom query int false "Map zoom level" default(10)
// @Param categories query string false "Comma-separated category IDs"
// @Param condition query string false "Item condition filter"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {object} MapClustersResponse "Map clusters or markers data"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidBounds"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.mapError"
// @Router /api/v1/marketplace/map/clusters [get]
func (h *MarketplaceHandler) GetMapClusters(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.missingBounds")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
	}

	zoom, err := strconv.Atoi(zoomStr)
	if err != nil {
		zoom = 10
	}

	// Получаем фильтры
	categoryIDs := c.Query("categories", "")
	condition := c.Query("condition", "")
	minPrice := c.Query("min_price", "")
	maxPrice := c.Query("max_price", "")

	// Парсим цены
	var minPriceFloat, maxPriceFloat *float64
	if minPrice != "" {
		if parsed, err := strconv.ParseFloat(minPrice, 64); err == nil {
			minPriceFloat = &parsed
		}
	}
	if maxPrice != "" {
		if parsed, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			maxPriceFloat = &parsed
		}
	}

	// Для больших zoom уровней возвращаем отдельные маркеры
	if zoom >= 15 {
		listings, err := h.service.GetListingsInBounds(c.Context(),
			neLat64, neLng64, swLat64, swLng64, zoom,
			categoryIDs, condition, minPriceFloat, maxPriceFloat)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
		}

		response := MapClustersData{
			Type:  "markers",
			Data:  listings,
			Zoom:  zoom,
			Count: len(listings),
		}
		return utils.SuccessResponse(c, response)
	}

	// Для меньших zoom уровней возвращаем кластеры
	clusters, err := h.service.GetMapClusters(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
	}

	response := MapClustersResponse{
		Success: true,
		Data: MapClustersData{
			Type:  "clusters",
			Data:  clusters,
			Zoom:  zoom,
			Count: len(clusters),
		},
	}
	return c.JSON(response)
}
