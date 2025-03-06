// backend/internal/proj/geocode/handler/geocode.go
package handler

import (
    "backend/pkg/utils"
    "backend/internal/proj/geocode/service"
    "github.com/gofiber/fiber/v2"
    "log"
    "strconv"
)

type GeocodeHandler struct {
    geocodeService service.GeocodeServiceInterface
}

func NewGeocodeHandler(geocodeService service.GeocodeServiceInterface) *GeocodeHandler {
    return &GeocodeHandler{
        geocodeService: geocodeService,
    }
}

// ReverseGeocode получает адрес по координатам
func (h *GeocodeHandler) ReverseGeocode(c *fiber.Ctx) error {
    latStr := c.Query("lat")
    lonStr := c.Query("lon")

    lat, err := strconv.ParseFloat(latStr, 64)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid latitude")
    }

    lon, err := strconv.ParseFloat(lonStr, 64)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid longitude")
    }

    // Получаем адрес по координатам
    address, err := h.geocodeService.ReverseGeocode(c.Context(), lat, lon)
    if err != nil {
        log.Printf("Error with reverse geocoding: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error geocoding coordinates")
    }

    return utils.SuccessResponse(c, address)
}

// GetCitySuggestions возвращает предложения городов по частичному названию
func (h *GeocodeHandler) GetCitySuggestions(c *fiber.Ctx) error {
    query := c.Query("q", "")
    limit := c.QueryInt("limit", 10)

    if query == "" || len(query) < 2 {
        return utils.SuccessResponse(c, []interface{}{})
    }

    suggestions, err := h.geocodeService.GetCitySuggestions(c.Context(), query, limit)
    if err != nil {
        log.Printf("Error fetching city suggestions: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching suggestions")
    }

    return utils.SuccessResponse(c, suggestions)
}