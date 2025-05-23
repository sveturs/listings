// backend/internal/proj/marketplace/handler/map.go
package handler

import (
	"backend/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetListingsInBounds возвращает объявления в указанных границах карты
func (h *MarketplaceHandler) GetListingsInBounds(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing bounding box parameters")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lat")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lng")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lat")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lng")
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get listings")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"listings": listings,
		"bounds": map[string]interface{}{
			"ne": map[string]float64{"lat": neLat64, "lng": neLng64},
			"sw": map[string]float64{"lat": swLat64, "lng": swLng64},
		},
		"zoom":  zoom,
		"count": len(listings),
	})
}

// GetMapClusters возвращает кластеризованные данные для карты
func (h *MarketplaceHandler) GetMapClusters(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing bounding box parameters")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lat")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lng")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lat")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lng")
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
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get listings")
		}

		return utils.SuccessResponse(c, map[string]interface{}{
			"type":  "markers",
			"data":  listings,
			"zoom":  zoom,
			"count": len(listings),
		})
	}

	// Для меньших zoom уровней возвращаем кластеры
	clusters, err := h.service.GetMapClusters(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get clusters")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"type":  "clusters",
		"data":  clusters,
		"zoom":  zoom,
		"count": len(clusters),
	})
}
