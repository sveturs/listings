package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/gis/service"
	"backend/internal/proj/gis/types"
	"backend/pkg/utils"
)

// SpatialHandler обработчик пространственных запросов
type SpatialHandler struct {
	service *service.SpatialService
}

// NewSpatialHandler создает новый обработчик
func NewSpatialHandler(service *service.SpatialService) *SpatialHandler {
	return &SpatialHandler{
		service: service,
	}
}

// SearchListings поиск объявлений
// @Summary Пространственный поиск объявлений
// @Description Поиск объявлений с учетом геолокации, фильтров и сортировки
// @Tags GIS
// @Accept json
// @Produce json
// @Param bounds query string false "Границы поиска в формате: north,south,east,west"
// @Param center query string false "Центр поиска в формате: lat,lng"
// @Param radius_km query number false "Радиус поиска в километрах"
// @Param categories query array false "Категории объявлений"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param currency query string false "Валюта"
// @Param q query string false "Текстовый поиск"
// @Param sort_by query string false "Поле сортировки (distance, price, created_at)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Param limit query int false "Количество результатов"
// @Param offset query int false "Смещение"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.SearchResponse} "Результаты поиска"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/search [get]
func (h *SpatialHandler) SearchListings(c *fiber.Ctx) error {
	// Парсим параметры запроса
	params := types.SearchParams{}

	// Границы поиска
	if boundsStr := c.Query("bounds"); boundsStr != "" {
		bounds, err := parseBounds(boundsStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidBounds")
		}
		params.Bounds = &bounds
	}

	// Центр и радиус поиска
	if centerStr := c.Query("center"); centerStr != "" {
		center, err := parsePoint(centerStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidCenter")
		}
		params.Center = &center
	}

	// Поддерживаем также latitude/longitude отдельно
	if lat := c.QueryFloat("latitude", 0); lat != 0 {
		if lng := c.QueryFloat("longitude", 0); lng != 0 {
			params.Center = &types.Point{Lat: lat, Lng: lng}
		}
	}

	if radiusStr := c.Query("radius_km"); radiusStr != "" {
		radius, err := strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidRadius")
		}
		params.RadiusKm = radius
	}

	// Поддерживаем параметр distance в формате "10km"
	if distanceStr := c.Query("distance"); distanceStr != "" {
		if strings.HasSuffix(distanceStr, "km") {
			distanceStr = strings.TrimSuffix(distanceStr, "km")
			if radius, err := strconv.ParseFloat(distanceStr, 64); err == nil {
				params.RadiusKm = radius
			}
		}
	}

	// Категории - поддерживаем оба параметра для совместимости
	categories := c.Query("categories")
	if categories != "" {
		params.Categories = strings.Split(categories, ",")
	}

	// Также поддерживаем category_id для фильтрации по ID категории
	categoryID := c.Query("category_id")
	if categoryID != "" && len(params.Categories) == 0 {
		params.CategoryIDs = []int{}
		if id, err := strconv.Atoi(categoryID); err == nil {
			params.CategoryIDs = append(params.CategoryIDs, id)
		}
	}

	// Фильтры по цене
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidMinPrice")
		}
		params.MinPrice = &minPrice
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidMaxPrice")
		}
		params.MaxPrice = &maxPrice
	}

	// Остальные параметры
	params.SearchQuery = c.Query("q")
	params.SortBy = c.Query("sort_by", "created_at")
	params.SortOrder = c.Query("sort_order", "desc")
	params.Limit = c.QueryInt("limit", 50)
	params.Offset = c.QueryInt("offset", 0)

	// Выполняем поиск
	response, err := h.service.SearchListings(c.Context(), params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.searchError")
	}

	return utils.SuccessResponse(c, response)
}

// GetNearbyListings получение ближайших объявлений
// @Summary Получение ближайших объявлений
// @Description Получение объявлений в радиусе от заданной точки
// @Tags GIS
// @Accept json
// @Produce json
// @Param lat query number true "Широта"
// @Param lng query number true "Долгота"
// @Param radius_km query number false "Радиус поиска в километрах (по умолчанию 5)"
// @Param limit query int false "Количество результатов (по умолчанию 20)"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.SearchResponse} "Ближайшие объявления"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/nearby [get]
func (h *SpatialHandler) GetNearbyListings(c *fiber.Ctx) error {
	// Парсим координаты
	lat := c.QueryFloat("lat", 0)
	lng := c.QueryFloat("lng", 0)

	if lat == 0 || lng == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.coordinatesRequired")
	}

	center := types.Point{Lat: lat, Lng: lng}
	radiusKm := c.QueryFloat("radius_km", 5.0)
	limit := c.QueryInt("limit", 20)

	// Получаем ближайшие объявления
	response, err := h.service.GetNearbyListings(c.Context(), center, radiusKm, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.searchError")
	}

	return utils.SuccessResponse(c, response)
}

// GetListingLocation получение геоданных объявления
// @Summary Получение геоданных объявления
// @Description Получение координат и адреса объявления
// @Tags GIS
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.GeoListing} "Геоданные объявления"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный ID"
// @Failure 404 {object} utils.ErrorResponseSwag "Объявление не найдено"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/listings/{id}/location [get]
func (h *SpatialHandler) GetListingLocation(c *fiber.Ctx) error {
	// Парсим ID
	idStr := c.Params("id")
	listingID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidListingId")
	}

	// Получаем геоданные
	listing, err := h.service.GetListingLocation(c.Context(), listingID)
	if err != nil {
		if err == types.ErrLocationNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "gis.listingNotFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.getLocationError")
	}

	return utils.SuccessResponse(c, listing)
}

// UpdateListingLocation обновление геолокации объявления
// @Summary Обновление геолокации объявления
// @Description Обновление координат и адреса объявления
// @Tags GIS
// @Accept json
// @Produce json
// @Param id path string true "ID объявления"
// @Param body body types.UpdateLocationRequest true "Новые координаты и адрес"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "Успешное обновление"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные данные"
// @Failure 403 {object} utils.ErrorResponseSwag "Нет прав на изменение"
// @Failure 404 {object} utils.ErrorResponseSwag "Объявление не найдено"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/listings/{id}/location [put]
// @Security BearerAuth
func (h *SpatialHandler) UpdateListingLocation(c *fiber.Ctx) error {
	// Парсим ID
	idStr := c.Params("id")
	listingID, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidListingId")
	}

	// Парсим тело запроса
	var req types.UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidRequest")
	}

	// Базовая валидация выполнена через теги validate в структуре

	// TODO: Проверка прав доступа
	// userID := c.Locals("userID").(uuid.UUID)
	// Проверить, что пользователь является владельцем объявления

	// Обновляем локацию
	location := types.Point{Lat: req.Lat, Lng: req.Lng}
	err = h.service.UpdateListingLocation(c.Context(), listingID, location, req.Address)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.updateLocationError")
	}

	return utils.SuccessResponse(c, "Location updated successfully")
}

// UpdateListingAddress обновление адреса объявления (Phase 2)
// @Summary Обновление адреса объявления с валидацией
// @Description Обновляет геолокацию и адрес объявления с полной валидацией и логированием
// @Tags gis
// @Accept json
// @Produce json
// @Param id path int true "ID объявления"
// @Param request body types.UpdateAddressRequest true "Данные для обновления"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.EnhancedListingGeo} "Обновленные геоданные"
// @Failure 400 {object} utils.ErrorResponseSwag "Ошибка валидации"
// @Failure 403 {object} utils.ErrorResponseSwag "Нет прав доступа"
// @Failure 404 {object} utils.ErrorResponseSwag "Объявление не найдено"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка"
// @Router /api/v1/gis/listings/{id}/address [put]
// @Security BearerAuth
func (h *SpatialHandler) UpdateListingAddress(c *fiber.Ctx) error {
	// Получаем ID объявления
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidListingID")
	}

	// Парсим запрос
	var req types.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.parseError")
	}

	// Базовая валидация
	if req.Address == "" || len(req.Address) < 5 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.validationError")
	}

	// Валидируем значения enum
	if !req.LocationPrivacy.IsValid() {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidPrivacyLevel")
	}

	if !req.InputMethod.IsValid() {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidInputMethod")
	}

	// Получаем пользователя из контекста (middleware должен установить это)
	// Пока используем захардкоженное значение для тестирования
	userID := int64(1) // TODO: получать из middleware

	ctx := c.Context()

	// Получаем IP и User-Agent для логирования
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	// Обновляем адрес через сервис
	updatedGeo, err := h.service.UpdateListingAddress(
		ctx,
		int64(listingID),
		int64(userID),
		req,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// Разные типы ошибок
		switch err {
		case types.ErrListingNotFound:
			return utils.ErrorResponse(c, fiber.StatusNotFound, "gis.listingNotFound")
		case types.ErrAccessDenied:
			return utils.ErrorResponse(c, fiber.StatusForbidden, "gis.accessDenied")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.updateError")
		}
	}

	return utils.SuccessResponse(c, updatedGeo)
}

// RadiusSearch радиусный поиск объявлений
// @Summary Радиусный поиск объявлений
// @Description Поиск объявлений в заданном радиусе от точки с детальными фильтрами
// @Tags GIS
// @Accept json
// @Produce json
// @Param latitude query number true "Широта центра поиска"
// @Param longitude query number true "Долгота центра поиска"
// @Param radius query number true "Радиус поиска в метрах"
// @Param categories query array false "Категории объявлений"
// @Param category_ids query array false "ID категорий"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param currency query string false "Валюта"
// @Param q query string false "Текстовый поиск"
// @Param sort_by query string false "Поле сортировки (distance, price, created_at)"
// @Param sort_order query string false "Порядок сортировки (asc, desc)"
// @Param limit query int false "Количество результатов (по умолчанию 50, максимум 1000)"
// @Param offset query int false "Смещение"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.RadiusSearchResponse} "Результаты радиусного поиска"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/search/radius [get]
func (h *SpatialHandler) RadiusSearch(c *fiber.Ctx) error {
	// ВРЕМЕННАЯ БЛОКИРОВКА: Проверяем Referer для блокировки радиусного поиска на странице районов
	referer := c.Get("Referer")
	userAgent := c.Get("User-Agent")

	log.Info().
		Str("referer", referer).
		Str("user_agent", userAgent).
		Str("path", c.Path()).
		Str("query", string(c.Request().URI().QueryString())).
		Msg("🔍 BACKEND: Radius search request details")

	if strings.Contains(referer, "/districts") {
		log.Info().Str("referer", referer).Msg("🚫 BACKEND: Blocking radius search from districts page")

		// Возвращаем пустой результат
		return utils.SuccessResponse(c, types.RadiusSearchResponse{
			Listings:     []types.GeoListing{},
			TotalCount:   0,
			HasMore:      false,
			SearchRadius: 5000,
			SearchCenter: types.Point{
				Lat: c.QueryFloat("latitude", 0),
				Lng: c.QueryFloat("longitude", 0),
			},
		})
	}

	// Получаем основные параметры
	lat := c.QueryFloat("latitude", 0)
	lng := c.QueryFloat("longitude", 0)

	// Проверяем, есть ли заголовок, указывающий на комбинированный поиск
	isCombinedSearch := c.Get("X-Combined-Search") == "true"

	if isCombinedSearch {
		log.Info().
			Float64("lat", lat).
			Float64("lng", lng).
			Msg("✅ BACKEND: Combined district+radius search detected")
	}

	// Создаем запрос радиусного поиска
	var req types.RadiusSearchRequest

	// Парсим основные параметры
	req.Latitude = c.QueryFloat("latitude", 0)
	req.Longitude = c.QueryFloat("longitude", 0)
	req.Radius = c.QueryFloat("radius", 0)
	req.Limit = c.QueryInt("limit", 50)
	req.Offset = c.QueryInt("offset", 0)

	// Валидация основных параметров
	if err := req.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidSearchParams")
	}

	// Парсим фильтры если есть дополнительные параметры
	if c.Query("categories") != "" || c.Query("min_price") != "" || c.Query("max_price") != "" ||
		c.Query("currency") != "" || c.Query("q") != "" || c.Query("sort_by") != "" || c.Query("sort_order") != "" {

		req.Filters = &types.RadiusFilters{}

		// Категории - поддерживаем передачу как строк ID
		if categories := c.Query("categories"); categories != "" {
			// Пытаемся парсить как ID
			categoryStrings := strings.Split(categories, ",")
			categoryIDs := make([]int, 0, len(categoryStrings))
			for _, catStr := range categoryStrings {
				if id, err := strconv.Atoi(strings.TrimSpace(catStr)); err == nil {
					categoryIDs = append(categoryIDs, id)
				}
			}

			if len(categoryIDs) > 0 {
				req.Filters.CategoryIDs = categoryIDs
				log.Info().
					Str("categories_raw", categories).
					Ints("category_ids_parsed", categoryIDs).
					Msg("🔍 BACKEND: Parsed categories as IDs")
			} else {
				// Если не удалось распарсить как ID, сохраняем как строки
				req.Filters.Categories = categoryStrings
				log.Info().
					Str("categories_raw", categories).
					Strs("categories_parsed", req.Filters.Categories).
					Msg("🔍 BACKEND: Parsing categories as strings")
			}
		}

		// ID категорий
		if categoryIDs := c.Query("category_ids"); categoryIDs != "" {
			idStrings := strings.Split(categoryIDs, ",")
			req.Filters.CategoryIDs = make([]int, 0, len(idStrings))
			for _, idStr := range idStrings {
				if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
					req.Filters.CategoryIDs = append(req.Filters.CategoryIDs, id)
				}
			}
		}

		// Фильтры по цене
		if minPriceStr := c.Query("min_price"); minPriceStr != "" {
			if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
				req.Filters.MinPrice = &minPrice
			}
		}

		if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
			if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
				req.Filters.MaxPrice = &maxPrice
			}
		}

		// Остальные фильтры
		req.Filters.SearchQuery = c.Query("q")
		req.Filters.SortBy = c.Query("sort_by", "distance")
		req.Filters.SortOrder = c.Query("sort_order", "asc")

		// Фильтр по пользователю (опционально)
		if userIDStr := c.Query("user_id"); userIDStr != "" {
			if userID, err := strconv.Atoi(userIDStr); err == nil {
				req.Filters.UserID = &userID
			}
		}

		// Фильтр по статусу
		req.Filters.Status = c.Query("status")
	}

	// Выполняем радиусный поиск
	response, err := h.service.SearchByRadius(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.radiusSearchError")
	}

	return utils.SuccessResponse(c, response)
}

// Helper функции

func parseBounds(boundsStr string) (types.Bounds, error) {
	var bounds types.Bounds
	// Ожидаем формат: south,west,north,east (как в Leaflet/OpenStreetMap)
	_, err := fmt.Sscanf(boundsStr, "%f,%f,%f,%f",
		&bounds.South, &bounds.West, &bounds.North, &bounds.East)
	return bounds, err
}

func parsePoint(pointStr string) (types.Point, error) {
	var point types.Point
	_, err := fmt.Sscanf(pointStr, "%f,%f", &point.Lat, &point.Lng)
	return point, err
}
