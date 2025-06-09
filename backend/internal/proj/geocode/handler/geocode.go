// backend/internal/proj/geocode/handler/geocode.go
package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/proj/geocode/service"
	"backend/pkg/utils"
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
// @Summary Обратное геокодирование
// @Description Получает адрес по географическим координатам
// @Tags geocode
// @Accept json
// @Produce json
// @Param lat query float64 true "Широта"
// @Param lon query float64 true "Долгота"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.GeoLocation} "Адрес по координатам"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidLatitude или validation.invalidLongitude"
// @Failure 500 {object} utils.ErrorResponseSwag "geocode.reverseError"
// @Router /api/v1/geocode/reverse [get]
func (h *GeocodeHandler) ReverseGeocode(c *fiber.Ctx) error {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidLatitude")
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidLongitude")
	}

	// Получаем адрес по координатам
	var address *models.GeoLocation
	address, err = h.geocodeService.ReverseGeocode(c.Context(), lat, lon)
	if err != nil {
		log.Printf("Error with reverse geocoding: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocode.reverseError")
	}

	return utils.SuccessResponse(c, address)
}

// GetCitySuggestions возвращает предложения городов по частичному названию
// @Summary Поиск городов
// @Description Возвращает список городов по частичному названию для автодополнения
// @Tags geocode
// @Accept json
// @Produce json
// @Param q query string true "Поисковый запрос (минимум 2 символа)"
// @Param limit query int false "Максимальное количество результатов" default(10)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.GeoLocation} "Список городов"
// @Failure 500 {object} utils.ErrorResponseSwag "geocode.suggestionsError"
// @Router /api/v1/geocode/cities [get]
func (h *GeocodeHandler) GetCitySuggestions(c *fiber.Ctx) error {
	query := c.Query("q", "")
	limit := c.QueryInt("limit", 10)

	if query == "" || len(query) < 2 {
		return utils.SuccessResponse(c, []interface{}{})
	}

	suggestions, err := h.geocodeService.GetCitySuggestions(c.Context(), query, limit)
	if err != nil {
		log.Printf("Error fetching city suggestions: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocode.suggestionsError")
	}

	return utils.SuccessResponse(c, suggestions)
}
