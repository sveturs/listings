package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

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

	if radiusStr := c.Query("radius_km"); radiusStr != "" {
		radius, err := strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidRadius")
		}
		params.RadiusKm = radius
	}

	// Категории
	categories := c.Query("categories")
	if categories != "" {
		params.Categories = strings.Split(categories, ",")
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
	params.Currency = c.Query("currency")
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

// GetClusters получение кластеров
// @Summary Получение кластеров объявлений
// @Description Получение кластеризованных объявлений для отображения на карте
// @Tags GIS
// @Accept json
// @Produce json
// @Param bounds query string true "Границы видимой области карты в формате: north,south,east,west"
// @Param zoom_level query int true "Уровень зума карты (0-20)"
// @Param categories query array false "Категории объявлений"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Param currency query string false "Валюта"
// @Param grid_size query int false "Размер сетки для кластеризации"
// @Success 200 {object} utils.SuccessResponseSwag{data=types.ClusterResponse} "Кластеры и объявления"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректные параметры"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/gis/clusters [get]
func (h *SpatialHandler) GetClusters(c *fiber.Ctx) error {
	// Парсим обязательные параметры
	boundsStr := c.Query("bounds")
	if boundsStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.boundsRequired")
	}

	bounds, err := parseBounds(boundsStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidBounds")
	}

	zoomLevel := c.QueryInt("zoom_level", -1)
	if zoomLevel < 0 || zoomLevel > 20 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "gis.invalidZoomLevel")
	}

	// Создаем параметры кластеризации
	params := types.ClusterParams{
		Bounds:    bounds,
		ZoomLevel: zoomLevel,
		GridSize:  c.QueryInt("grid_size", 0),
	}

	// Категории
	categories := c.Query("categories")
	if categories != "" {
		params.Categories = strings.Split(categories, ",")
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

	params.Currency = c.Query("currency")

	// Получаем кластеры
	response, err := h.service.GetClusters(c.Context(), params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "gis.clusterError")
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

// Helper функции

func parseBounds(boundsStr string) (types.Bounds, error) {
	var bounds types.Bounds
	_, err := fmt.Sscanf(boundsStr, "%f,%f,%f,%f",
		&bounds.North, &bounds.South, &bounds.East, &bounds.West)
	return bounds, err
}

func parsePoint(pointStr string) (types.Point, error) {
	var point types.Point
	_, err := fmt.Sscanf(pointStr, "%f,%f", &point.Lat, &point.Lng)
	return point, err
}
