package handler

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/storefronts/service"
	"backend/internal/proj/storefronts/storage/opensearch"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// StorefrontHandler HTTP handler для витрин
type StorefrontHandler struct {
	service service.StorefrontService
}

// NewStorefrontHandler создает новый handler
func NewStorefrontHandler(service service.StorefrontService) *StorefrontHandler {
	return &StorefrontHandler{
		service: service,
	}
}

// CreateStorefront создает новую витрину
// @Summary Create new storefront
// @Description Creates a new storefront for the authenticated user
// @Tags storefronts
// @Accept json
// @Produce json
// @Param storefront body backend_internal_domain_models.StorefrontCreateDTO true "Storefront data"
// @Success 201 {object} backend_internal_domain_models.Storefront "Created storefront"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Storefront limit reached"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts [post]
func (h *StorefrontHandler) CreateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var dto models.StorefrontCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	storefront, err := h.service.CreateStorefront(c.Context(), userID, &dto)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStorefrontLimitReached):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.limit_reached")
		case errors.Is(err, service.ErrInvalidLocation):
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_location")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.create_failed")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(storefront)
}

// GetStorefront получает витрину по ID
// @Summary Get storefront by ID
// @Description Returns storefront details by ID
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} backend_internal_domain_models.Storefront "Storefront details"
// @Failure 404 {object} utils.ErrorResponseSwag "Storefront not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/{id} [get]
func (h *StorefrontHandler) GetStorefront(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	storefront, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_failed")
	}

	return c.JSON(storefront)
}

// GetStorefrontBySlug получает витрину по slug или ID
// @Summary Get storefront by slug or ID
// @Description Returns storefront details by slug or ID
// @Tags storefronts
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug or ID"
// @Success 200 {object} backend_internal_domain_models.Storefront "Storefront details"
// @Failure 404 {object} utils.ErrorResponseSwag "Storefront not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/slug/{slug} [get]
func (h *StorefrontHandler) GetStorefrontBySlug(c *fiber.Ctx) error {
	slugOrID := c.Params("slug")

	// Пробуем сначала как ID
	if id, err := strconv.Atoi(slugOrID); err == nil {
		storefront, err := h.service.GetByID(c.Context(), id)
		if err == nil {
			return c.JSON(storefront)
		}
		// Если не нашли по ID, пробуем как slug
	}

	// Пробуем как slug
	storefront, err := h.service.GetBySlug(c.Context(), slugOrID)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_failed")
	}

	return c.JSON(storefront)
}

// UpdateStorefront обновляет витрину
// @Summary Update storefront
// @Description Updates storefront details
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param storefront body backend_internal_domain_models.StorefrontUpdateDTO true "Update data"
// @Success 200 {object} map[string]string "Storefront updated"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 404 {object} utils.ErrorResponseSwag "Storefront not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [put]
func (h *StorefrontHandler) UpdateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	var dto models.StorefrontUpdateDTO
	if err := c.BodyParser(&dto); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	err = h.service.Update(c.Context(), userID, id, &dto)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientPermissions):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		case errors.Is(err, service.ErrFeatureNotAvailable):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.feature_not_available")
		case errors.Is(err, postgres.ErrNotFound):
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.update_failed")
		}
	}

	return c.JSON(fiber.Map{
		"message": "Storefront updated successfully",
	})
}

// DeleteStorefront удаляет витрину
// @Summary Delete storefront
// @Description Deletes a storefront. Soft delete by default (marks as inactive). Admins can use ?hard=true for permanent removal
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param hard query bool false "Hard delete (permanent removal, admin only). Default: false (soft delete)"
// @Success 200 {object} map[string]string "Storefront deleted"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Only owner can delete storefront"
// @Failure 404 {object} utils.ErrorResponseSwag "Storefront not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [delete]
func (h *StorefrontHandler) DeleteStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	// Pass admin status and hard delete flag through context
	isAdmin, _ := c.Locals("is_admin").(bool)
	hardDelete := c.Query("hard") == "true" // Query parameter для выбора жесткого удаления

	ctx := context.WithValue(context.Background(), isAdminKey, isAdmin)
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, hardDeleteKey, hardDelete)

	err = h.service.Delete(ctx, userID, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnauthorized):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.only_owner_can_delete")
		case errors.Is(err, postgres.ErrNotFound):
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.delete_failed")
		}
	}

	return c.JSON(fiber.Map{
		"message": "Storefront deleted successfully",
	})
}

// ListStorefronts получает список витрин
// @Summary List storefronts
// @Description Returns paginated list of storefronts with filters
// @Tags storefronts
// @Accept json
// @Produce json
// @Param user_id query int false "Filter by user ID"
// @Param is_active query bool false "Filter by active status"
// @Param is_verified query bool false "Filter by verification status"
// @Param city query string false "Filter by city"
// @Param min_rating query float32 false "Minimum rating filter"
// @Param search query string false "Search in name and description"
// @Param lat query float64 false "Latitude for geo search"
// @Param lng query float64 false "Longitude for geo search"
// @Param radius_km query float64 false "Radius in km for geo search"
// @Param sort_by query string false "Sort field (rating, products_count, distance)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Results per page (max 100)"
// @Param offset query int false "Results offset"
// @Success 200 {object} StorefrontsListResponse "List of storefronts"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts [get]
func (h *StorefrontHandler) ListStorefronts(c *fiber.Ctx) error {
	filter := &models.StorefrontFilter{
		Limit:  20,
		Offset: 0,
	}

	// Проверяем, является ли пользователь админом
	isAdmin := false
	if userIDLocal := c.Locals("user_id"); userIDLocal != nil {
		// Проверяем роль пользователя (предполагаем, что есть поле is_admin или role)
		if adminLocal := c.Locals("is_admin"); adminLocal != nil {
			isAdmin = adminLocal.(bool)
		}
	}

	// Set admin flag in filter for repository
	filter.IsAdminRequest = isAdmin

	// Парсим query параметры
	if userID := c.QueryInt("user_id"); userID > 0 {
		filter.UserID = &userID
	}

	// Для админа и владельца (если указан userID) не применяем фильтр по активности по умолчанию
	// Для публичных запросов - применяем
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == boolValueTrue
		filter.IsActive = &active
	} else if !isAdmin && filter.UserID == nil {
		// Для публичных запросов (не админ и не запрос конкретного пользователя) показываем только активные
		active := true
		filter.IsActive = &active
	}

	if isVerified := c.Query("is_verified"); isVerified != "" {
		verified := isVerified == boolValueTrue
		filter.IsVerified = &verified
	}

	if city := c.Query("city"); city != "" {
		filter.City = &city
	}

	if minRating, err := strconv.ParseFloat(c.Query("min_rating"), 64); err == nil {
		filter.MinRating = &minRating
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	// Геофильтр
	if lat, err := strconv.ParseFloat(c.Query("lat"), 64); err == nil {
		filter.Latitude = &lat
	}
	if lng, err := strconv.ParseFloat(c.Query("lng"), 64); err == nil {
		filter.Longitude = &lng
	}
	if radius, err := strconv.ParseFloat(c.Query("radius_km"), 64); err == nil {
		filter.RadiusKm = &radius
	}

	// Сортировка
	filter.SortBy = c.Query("sort_by", "created_at")
	filter.SortOrder = c.Query("sort_order", "desc")

	// Пагинация
	if limit := c.QueryInt("limit", 20); limit > 0 && limit <= 100 {
		filter.Limit = limit
	}
	filter.Offset = c.QueryInt("offset", 0)

	storefronts, total, err := h.service.Search(c.Context(), filter)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list storefronts")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.list_failed")
	}

	return c.JSON(StorefrontsListResponse{
		Storefronts: storefronts,
		Total:       total,
		Limit:       filter.Limit,
		Offset:      filter.Offset,
	})
}

// GetMyStorefronts получает витрины текущего пользователя
// @Summary Get my storefronts
// @Description Returns list of storefronts owned by current user
// @Tags storefronts
// @Accept json
// @Produce json
// @Success 200 {object} []backend_internal_domain_models.Storefront "List of user's storefronts"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/my [get]
func (h *StorefrontHandler) GetMyStorefronts(c *fiber.Ctx) error {
	logger.Info().Msg("GetMyStorefronts called")

	userID := c.Locals("user_id").(int)
	logger.Info().Int("userID", userID).Msg("Getting storefronts for user")

	storefronts, err := h.service.ListUserStorefronts(c.Context(), userID)
	if err != nil {
		// Логируем конкретную ошибку
		logger.Error().Err(err).Int("userID", userID).Msg("Failed to get user storefronts")

		// Обрабатываем специфичные ошибки
		if errors.Is(err, service.ErrRepositoryNotInitialized) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.repository_not_initialized")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_user_storefronts_failed")
	}

	return c.JSON(storefronts)
}

// GetMapData получает данные для карты
// @Summary Get storefronts for map
// @Description Returns storefronts within map bounds with minimal data for performance
// @Tags storefronts,map
// @Accept json
// @Produce json
// @Param min_lat query float64 true "Minimum latitude"
// @Param max_lat query float64 true "Maximum latitude"
// @Param min_lng query float64 true "Minimum longitude"
// @Param max_lng query float64 true "Maximum longitude"
// @Param min_rating query float32 false "Minimum rating filter"
// @Success 200 {object} []backend_internal_domain_models.StorefrontMapData "Map markers data"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid bounds"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/map [get]
func (h *StorefrontHandler) GetMapData(c *fiber.Ctx) error {
	// Парсим границы карты
	minLat, err1 := strconv.ParseFloat(c.Query("min_lat"), 64)
	maxLat, err2 := strconv.ParseFloat(c.Query("max_lat"), 64)
	minLng, err3 := strconv.ParseFloat(c.Query("min_lng"), 64)
	maxLng, err4 := strconv.ParseFloat(c.Query("max_lng"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_map_bounds")
	}

	bounds := postgres.GeoBounds{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}

	// Дополнительные фильтры
	filter := &models.StorefrontFilter{}
	if minRating, err := strconv.ParseFloat(c.Query("min_rating"), 64); err == nil {
		filter.MinRating = &minRating
	}

	mapData, err := h.service.GetMapData(c.Context(), bounds, filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_map_data_failed")
	}

	return c.JSON(mapData)
}

// SearchOpenSearch выполняет поиск витрин через OpenSearch
// @Summary Search storefronts using OpenSearch
// @Description Performs advanced search of storefronts using OpenSearch engine
// @Tags storefronts,search
// @Accept json
// @Produce json
// @Param q query string false "Search query"
// @Param city query string false "City filter"
// @Param lat query float64 false "Latitude for geo search"
// @Param lng query float64 false "Longitude for geo search"
// @Param radius_km query int false "Search radius in kilometers"
// @Param min_rating query float64 false "Minimum rating"
// @Param is_verified query bool false "Only verified storefronts"
// @Param is_open_now query bool false "Only currently open storefronts"
// @Param payment_methods query string false "Payment methods (comma separated)"
// @Param has_delivery query bool false "Has delivery option"
// @Param has_self_pickup query bool false "Has self pickup option"
// @Param sort_by query string false "Sort field (rating, distance, products_count, created_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Number of results (max 100)" default(20)
// @Param offset query int false "Results offset" default(0)
// @Success 200 {object} opensearch.StorefrontSearchResult "Search results"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/search [get]
func (h *StorefrontHandler) SearchOpenSearch(c *fiber.Ctx) error {
	params := &opensearch.StorefrontSearchParams{
		Query: c.Query("q"),
		City:  c.Query("city"),
	}

	// Геолокация
	if lat, err := strconv.ParseFloat(c.Query("lat"), 64); err == nil {
		params.Latitude = lat
	}
	if lng, err := strconv.ParseFloat(c.Query("lng"), 64); err == nil {
		params.Longitude = lng
	}
	if radius, err := strconv.Atoi(c.Query("radius_km")); err == nil {
		params.RadiusKm = radius
	}

	// Фильтры
	if minRating, err := strconv.ParseFloat(c.Query("min_rating"), 64); err == nil {
		params.MinRating = minRating
	}
	if isVerified := c.Query("is_verified"); isVerified != "" {
		verified := isVerified == boolValueTrue
		params.IsVerified = &verified
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == boolValueTrue
		params.IsActive = &active
	}
	if isOpenNow := c.Query("is_open_now"); isOpenNow != "" {
		openNow := isOpenNow == boolValueTrue
		params.IsOpenNow = &openNow
	}
	if hasDelivery := c.Query("has_delivery"); hasDelivery != "" {
		delivery := hasDelivery == boolValueTrue
		params.HasDelivery = &delivery
	}
	if hasSelfPickup := c.Query("has_self_pickup"); hasSelfPickup != "" {
		selfPickup := hasSelfPickup == boolValueTrue
		params.HasSelfPickup = &selfPickup
	}

	// Методы оплаты
	if paymentMethods := c.Query("payment_methods"); paymentMethods != "" {
		params.PaymentMethods = strings.Split(paymentMethods, ",")
	}

	// Сортировка
	params.SortBy = c.Query("sort_by", "rating")
	params.SortOrder = c.Query("sort_order", "desc")

	// Пагинация
	params.Limit = c.QueryInt("limit", 20)
	if params.Limit > 100 {
		params.Limit = 100
	}
	params.Offset = c.QueryInt("offset", 0)

	result, err := h.service.SearchOpenSearch(c.Context(), params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.search_failed")
	}

	return c.JSON(result)
}

// GetNearbyStorefronts получает ближайшие витрины
// @Summary Get nearby storefronts
// @Description Returns storefronts within specified radius from coordinates
// @Tags storefronts,map
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lng query float64 true "Longitude"
// @Param radius_km query float64 false "Radius in kilometers (default 5)"
// @Param limit query int false "Maximum results (default 20, max 100)"
// @Success 200 {object} []backend_internal_domain_models.Storefront "Nearby storefronts"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid coordinates"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/nearby [get]
func (h *StorefrontHandler) GetNearbyStorefronts(c *fiber.Ctx) error {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)

	if err1 != nil || err2 != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_coordinates")
	}

	radiusKm := c.QueryFloat("radius_km", 5.0)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}

	storefronts, err := h.service.GetNearby(c.Context(), lat, lng, radiusKm, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_nearby_failed")
	}

	return c.JSON(storefronts)
}

// GetBusinessesInBuilding получает все бизнесы в здании
// @Summary Get businesses in building
// @Description Returns all storefronts located in the same building
// @Tags storefronts,map
// @Accept json
// @Produce json
// @Param lat query float64 true "Building latitude"
// @Param lng query float64 true "Building longitude"
// @Success 200 {object} []backend_internal_domain_models.StorefrontMapData "Businesses in building"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid coordinates"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/building [get]
func (h *StorefrontHandler) GetBusinessesInBuilding(c *fiber.Ctx) error {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)

	if err1 != nil || err2 != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_coordinates")
	}

	businesses, err := h.service.GetBusinessesInBuilding(c.Context(), lat, lng)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_businesses_in_building_failed")
	}

	return c.JSON(businesses)
}

// UpdateWorkingHours обновляет часы работы
// @Summary Update working hours
// @Description Updates storefront working hours
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param hours body []backend_internal_domain_models.StorefrontHours true "Working hours"
// @Success 200 {object} map[string]string "Hours updated"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/hours [put]
func (h *StorefrontHandler) UpdateWorkingHours(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	var hours []*models.StorefrontHours
	if err := c.BodyParser(&hours); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	err = h.service.UpdateWorkingHours(c.Context(), userID, id, hours)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.update_working_hours_failed")
	}

	return c.JSON(fiber.Map{
		"message": "Working hours updated successfully",
	})
}

// UpdatePaymentMethods обновляет методы оплаты
// @Summary Update payment methods
// @Description Updates storefront payment methods
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param methods body []backend_internal_domain_models.StorefrontPaymentMethod true "Payment methods"
// @Success 200 {object} map[string]string "Payment methods updated"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/payment-methods [put]
func (h *StorefrontHandler) UpdatePaymentMethods(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	var methods []*models.StorefrontPaymentMethod
	if err := c.BodyParser(&methods); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	err = h.service.UpdatePaymentMethods(c.Context(), userID, id, methods)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.update_payment_methods_failed")
	}

	return c.JSON(fiber.Map{
		"message": "Payment methods updated successfully",
	})
}

// UpdateDeliveryOptions обновляет опции доставки
// @Summary Update delivery options
// @Description Updates storefront delivery options
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param options body []backend_internal_domain_models.StorefrontDeliveryOption true "Delivery options"
// @Success 200 {object} map[string]string "Delivery options updated"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/delivery-options [put]
func (h *StorefrontHandler) UpdateDeliveryOptions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	var options []*models.StorefrontDeliveryOption
	if err := c.BodyParser(&options); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	err = h.service.UpdateDeliveryOptions(c.Context(), userID, id, options)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.update_delivery_options_failed")
	}

	return c.JSON(fiber.Map{
		"message": "Delivery options updated successfully",
	})
}

// RecordView записывает просмотр витрины
// @Summary Record storefront view
// @Description Records a view for analytics
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} map[string]string "View recorded"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/{id}/view [post]
func (h *StorefrontHandler) RecordView(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	err = h.service.RecordView(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.record_view_failed")
	}

	return c.JSON(fiber.Map{
		"message": "View recorded",
	})
}

// GetAnalytics получает аналитику витрины
// @Summary Get storefront analytics
// @Description Returns analytics data for specified period
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param from query string true "Start date (YYYY-MM-DD)"
// @Param to query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} []backend_internal_domain_models.StorefrontAnalytics "Analytics data"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/analytics [get]
func (h *StorefrontHandler) GetAnalytics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	// Парсим ISO формат даты (RFC3339) который приходит с frontend
	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		// Попробуем короткий формат
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_from_date")
		}
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		// Попробуем короткий формат
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_to_date")
		}
	}

	analytics, err := h.service.GetAnalytics(c.Context(), userID, id, from, to)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_analytics_failed")
	}

	return c.JSON(analytics)
}

// Response types

type StorefrontsListResponse struct {
	Storefronts []*models.Storefront `json:"storefronts"`
	Total       int                  `json:"total"`
	Limit       int                  `json:"limit"`
	Offset      int                  `json:"offset"`
}

// UploadLogo загружает логотип витрины
// @Summary Upload storefront logo
// @Description Uploads a logo image for the storefront
// @Tags storefronts
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Storefront ID"
// @Param logo formData file true "Logo file"
// @Success 200 {object} map[string]string "Logo URL"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/logo [post]
func (h *StorefrontHandler) UploadLogo(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	file, err := c.FormFile("logo")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.failed_to_get_logo_file")
	}

	// Читаем файл
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_open_file")
	}
	defer func() {
		if err := src.Close(); err != nil {
			// Логирование ошибки закрытия file
			_ = err // Explicitly ignore error
		}
	}()

	// Читаем данные
	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_read_file")
	}

	url, err := h.service.UploadLogo(c.Context(), userID, storefrontID, data, file.Filename)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_upload_logo")
	}

	return c.JSON(fiber.Map{
		"logo_url": url,
	})
}

// UploadBanner загружает баннер витрины
// @Summary Upload storefront banner
// @Description Uploads a banner image for the storefront
// @Tags storefronts
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Storefront ID"
// @Param banner formData file true "Banner file"
// @Success 200 {object} map[string]string "Banner URL"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/banner [post]
func (h *StorefrontHandler) UploadBanner(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	file, err := c.FormFile("banner")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.failed_to_get_banner_file")
	}

	// Читаем файл
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_open_file")
	}
	defer func() {
		if err := src.Close(); err != nil {
			// Логирование ошибки закрытия file
			_ = err // Explicitly ignore error
		}
	}()

	// Читаем данные
	data := make([]byte, file.Size)
	if _, readErr := src.Read(data); readErr != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_read_file")
	}

	url, err := h.service.UploadBanner(c.Context(), userID, storefrontID, data, file.Filename)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_upload_banner")
	}

	return c.JSON(fiber.Map{
		"banner_url": url,
	})
}

// GetStorefrontAnalytics возвращает аналитику витрины
// @Summary Get storefront analytics
// @Description Returns analytics data for a storefront
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param from query string false "Start date (RFC3339 format)"
// @Param to query string false "End date (RFC3339 format)"
// @Success 200 {object} backend_internal_domain_models.StorefrontAnalytics "Analytics data"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 404 {object} utils.ErrorResponseSwag "Storefront not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/analytics [get]
func (h *StorefrontHandler) GetStorefrontAnalytics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	// Парсим даты из query параметров
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var from, to time.Time
	now := time.Now()

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_from_date")
		}
	} else {
		// По умолчанию - последние 30 дней
		from = now.AddDate(0, 0, -30)
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_to_date")
		}
	} else {
		to = now
	}

	analytics, err := h.service.GetAnalytics(c.Context(), userID, storefrontID, from, to)
	if err != nil {
		if errors.Is(err, service.ErrStorefrontNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
		}
		if errors.Is(err, service.ErrInsufficientPermissions) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		}
		logger.Error().Err(err).Int("storefront_id", storefrontID).Msg("Failed to get storefront analytics")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.failed_to_get_analytics")
	}

	return c.JSON(analytics)
}
