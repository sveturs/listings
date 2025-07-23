package handler

import (
	"strconv"

	"backend/internal/proj/gis/service"
	"backend/internal/proj/gis/types"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GeocodingHandler обработчик для геокодирования
type GeocodingHandler struct {
	geocodingService *service.GeocodingService
}

// NewGeocodingHandler создает новый обработчик
func NewGeocodingHandler(geocodingService *service.GeocodingService) *GeocodingHandler {
	return &GeocodingHandler{
		geocodingService: geocodingService,
	}
}

// ValidateGeocode валидация геокодирования адреса
// @Summary Валидация геокодирования адреса
// @Description Проверяет и возвращает геокодированные данные для указанного адреса
// @Tags gis
// @Accept json
// @Produce json
// @Param request body types.GeocodeValidateRequest true "Запрос валидации"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.GeocodeValidateResponse} "Результат валидации"
// @Failure 400 {object} utils.ErrorResponseSwag "Ошибка валидации"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/geocode/validate [post]
func (h *GeocodingHandler) ValidateGeocode(c *fiber.Ctx) error {
	var req types.GeocodeValidateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.parseError")
	}

	// Базовая валидация
	if req.Address == "" || len(req.Address) < 5 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.validationError")
	}

	// Устанавливаем язык по умолчанию если не указан
	if req.Language == "" {
		req.Language = "en"
	}

	ctx := c.Context()
	result, err := h.geocodingService.ValidateAndGeocode(ctx, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocoding.serviceError")
	}

	return utils.SuccessResponse(c, result)
}

// SearchAddressSuggestions поиск предложений адресов
// @Summary Поиск предложений адресов
// @Description Возвращает список предложений адресов для автодополнения
// @Tags gis
// @Accept json
// @Produce json
// @Param q query string true "Поисковый запрос" minlength(2)
// @Param limit query int false "Лимит результатов" default(5)
// @Param language query string false "Язык ответа" default("en")
// @Param country_code query string false "Код страны для фильтрации"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]types.AddressSuggestion} "Список предложений"
// @Failure 400 {object} utils.ErrorResponseSwag "Ошибка валидации"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/geocode/suggestions [get]
func (h *GeocodingHandler) SearchAddressSuggestions(c *fiber.Ctx) error {
	query := c.Query("q")
	if len(query) < 2 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.queryTooShort")
	}

	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 10 {
			limit = l
		}
	}

	language := c.Query("language", "en")
	countryCode := c.Query("country_code")

	ctx := c.Context()
	suggestions, err := h.geocodingService.SearchSuggestions(ctx, query, limit, language, countryCode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocoding.suggestionsError")
	}

	return utils.SuccessResponse(c, suggestions)
}

// ReverseGeocode обратное геокодирование
// @Summary Обратное геокодирование
// @Description Получает адрес по координатам
// @Tags gis
// @Accept json
// @Produce json
// @Param lat query number true "Широта" minimum(-90) maximum(90)
// @Param lng query number true "Долгота" minimum(-180) maximum(180)
// @Param language query string false "Язык ответа" default("en")
// @Success 200 {object} utils.SuccessResponseSwag{data=types.AddressSuggestion} "Найденный адрес"
// @Failure 400 {object} utils.ErrorResponseSwag "Ошибка валидации"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/geocode/reverse [get]
func (h *GeocodingHandler) ReverseGeocode(c *fiber.Ctx) error {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	if latStr == "" || lngStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.missingCoordinates")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil || lat < -90 || lat > 90 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.invalidLatitude")
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil || lng < -180 || lng > 180 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.invalidLongitude")
	}

	language := c.Query("language", "en")

	ctx := c.Context()
	point := types.Point{Lat: lat, Lng: lng}

	result, err := h.geocodingService.ReverseGeocode(ctx, point, language)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocoding.reverseError")
	}

	return utils.SuccessResponse(c, result)
}

// GetCacheStats получение статистики кэша геокодирования
// @Summary Статистика кэша геокодирования
// @Description Возвращает статистику использования кэша геокодирования
// @Tags gis
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Статистика кэша"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/geocode/cache/stats [get]
func (h *GeocodingHandler) GetCacheStats(c *fiber.Ctx) error {
	ctx := c.Context()

	stats, err := h.geocodingService.GetCacheStats(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocoding.cacheStatsError")
	}

	return utils.SuccessResponse(c, stats)
}

// CleanupExpiredCache очистка устаревшего кэша
// @Summary Очистка устаревшего кэша
// @Description Удаляет устаревшие записи из кэша геокодирования
// @Tags gis
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]int64} "Количество удаленных записей"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/geocode/cache/cleanup [post]
func (h *GeocodingHandler) CleanupExpiredCache(c *fiber.Ctx) error {
	ctx := c.Context()

	deletedCount, err := h.geocodingService.CleanupExpiredCache(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocoding.cacheCleanupError")
	}

	result := map[string]int64{
		"deleted_count": deletedCount,
	}

	return utils.SuccessResponse(c, result)
}

// MultilingualReverseGeocode многоязычное обратное геокодирование
// @Summary Многоязычное обратное геокодирование
// @Description Получает адреса по координатам на трех языках (сербский, английский, русский)
// @Tags gis
// @Accept json
// @Produce json
// @Param request body types.MultilingualGeocodeRequest true "Координаты"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.MultilingualGeocodeResponse} "Адреса на разных языках"
// @Failure 400 {object} utils.ErrorResponseSwag "Ошибка валидации"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/geocode/multilingual [post]
func (h *GeocodingHandler) MultilingualReverseGeocode(c *fiber.Ctx) error {
	var req types.MultilingualGeocodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.parseError")
	}

	// Валидация координат
	if req.Latitude < -90 || req.Latitude > 90 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.invalidLatitude")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "geocoding.invalidLongitude")
	}

	ctx := c.Context()
	result, err := h.geocodingService.GetMultilingualAddress(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "geocoding.multilingualError")
	}

	return utils.SuccessResponse(c, result)
}
