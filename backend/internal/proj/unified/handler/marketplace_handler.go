// backend/internal/proj/unified/handler/marketplace_handler.go
package handler

import (
	"net/http"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/unified/service"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// MarketplaceHandler унифицированный handler для C2C и B2C listings
type MarketplaceHandler struct {
	service service.MarketplaceServiceInterface
	logger  zerolog.Logger
}

// NewMarketplaceHandler создает новый unified marketplace handler
func NewMarketplaceHandler(service service.MarketplaceServiceInterface) *MarketplaceHandler {
	return &MarketplaceHandler{
		service: service,
		logger:  logger.Get().With().Str("handler", "unified_marketplace").Logger(),
	}
}

// getRoutingContext извлекает routing контекст из fiber context
func (h *MarketplaceHandler) getRoutingContext(c *fiber.Ctx) *service.RoutingContext {
	userID, _ := authMiddleware.GetUserID(c)
	roles, _ := authMiddleware.GetRoles(c)

	// Проверяем наличие admin роли
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	return &service.RoutingContext{
		UserID:  userID,
		IsAdmin: isAdmin,
	}
}

// CreateListing godoc
// @Summary Create unified listing (C2C or B2C)
// @Description Создает новый listing (C2C или B2C в зависимости от source_type)
// @Tags unified-marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param listing body models.UnifiedListing true "Unified listing data"
// @Success 201 {object} map[string]interface{} "id, source_type"
// @Failure 400 {object} map[string]string "validation error"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 500 {object} map[string]string "internal error"
// @Router /api/v1/marketplace/listings [post]
func (h *MarketplaceHandler) CreateListing(c *fiber.Ctx) error {
	// Получаем user_id из JWT
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("CreateListing: user_id not found in context")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "marketplace.unauthorized",
		})
	}

	// Парсинг request body
	var req models.UnifiedListing
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("CreateListing: failed to parse request body")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_request_body",
		})
	}

	// Валидация обязательных полей
	if req.SourceType == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.source_type_required",
		})
	}

	if req.SourceType != service.SourceTypeC2C && req.SourceType != service.SourceTypeB2C {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_source_type",
		})
	}

	if req.Title == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.title_required",
		})
	}

	if req.Price <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_price",
		})
	}

	if req.CategoryID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.category_id_required",
		})
	}

	// B2C требует storefront_id
	if req.SourceType == service.SourceTypeB2C && (req.StorefrontID == nil || *req.StorefrontID == 0) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.storefront_id_required",
		})
	}

	// Устанавливаем user_id из JWT
	req.UserID = userID

	h.logger.Info().
		Str("source_type", req.SourceType).
		Int("user_id", userID).
		Str("title", req.Title).
		Msg("CreateListing: creating unified listing")

	// Получаем routing context
	routingCtx := h.getRoutingContext(c)

	// Создаем через service
	id, err := h.service.CreateListing(c.Context(), &req, routingCtx)
	if err != nil {
		h.logger.Error().Err(err).Msg("CreateListing: service error")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.failed_to_create_listing",
		})
	}

	h.logger.Info().Int64("id", id).Str("source_type", req.SourceType).Msg("CreateListing: success")

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success":     true,
		"id":          id,
		"source_type": req.SourceType,
	})
}

// GetListing godoc
// @Summary Get listing by ID and source type
// @Description Получить listing по ID и типу источника (c2c или b2c)
// @Tags unified-marketplace
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param source_type query string true "Source type: c2c or b2c"
// @Success 200 {object} models.UnifiedListing
// @Failure 400 {object} map[string]string "invalid parameters"
// @Failure 404 {object} map[string]string "not found"
// @Failure 500 {object} map[string]string "internal error"
// @Router /api/v1/marketplace/listings/{id} [get]
func (h *MarketplaceHandler) GetListing(c *fiber.Ctx) error {
	// Парсинг ID
	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Warn().Err(err).Msg("GetListing: invalid ID parameter")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_id",
		})
	}

	// Парсинг source_type
	sourceType := c.Query("source_type", "")
	if sourceType != service.SourceTypeC2C && sourceType != service.SourceTypeB2C {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_source_type",
		})
	}

	h.logger.Info().Int("id", id).Str("source_type", sourceType).Msg("GetListing: request")

	// Получаем routing context
	routingCtx := h.getRoutingContext(c)

	// Получаем через service
	listing, err := h.service.GetListing(c.Context(), int64(id), sourceType, routingCtx)
	if err != nil {
		h.logger.Error().Err(err).Int("id", id).Msg("GetListing: service error")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "marketplace.listing_not_found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    listing,
	})
}

// UpdateListing godoc
// @Summary Update listing
// @Description Обновить существующий listing (C2C или B2C)
// @Tags unified-marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Listing ID"
// @Param listing body models.UnifiedListing true "Updated listing data"
// @Success 200 {object} map[string]bool "success"
// @Failure 400 {object} map[string]string "validation error"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "forbidden"
// @Failure 404 {object} map[string]string "not found"
// @Failure 500 {object} map[string]string "internal error"
// @Router /api/v1/marketplace/listings/{id} [put]
func (h *MarketplaceHandler) UpdateListing(c *fiber.Ctx) error {
	// Получаем user_id из JWT
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("UpdateListing: user_id not found in context")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "marketplace.unauthorized",
		})
	}

	// Парсинг ID
	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Warn().Err(err).Msg("UpdateListing: invalid ID parameter")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_id",
		})
	}

	// Парсинг request body
	var req models.UnifiedListing
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn().Err(err).Msg("UpdateListing: failed to parse request body")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_request_body",
		})
	}

	// Валидация source_type
	if req.SourceType == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.source_type_required",
		})
	}

	if req.SourceType != service.SourceTypeC2C && req.SourceType != service.SourceTypeB2C {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_source_type",
		})
	}

	// Получаем routing context
	routingCtx := h.getRoutingContext(c)

	// Получаем существующий listing для проверки владельца
	existing, err := h.service.GetListing(c.Context(), int64(id), req.SourceType, routingCtx)
	if err != nil {
		h.logger.Error().Err(err).Int("id", id).Msg("UpdateListing: failed to get existing listing")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "marketplace.listing_not_found",
		})
	}

	// Проверка прав доступа (только владелец может обновлять)
	// TODO: добавить проверку роли admin
	if existing.UserID != userID {
		h.logger.Warn().
			Int("listing_user_id", existing.UserID).
			Int("request_user_id", userID).
			Msg("UpdateListing: permission denied")
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "marketplace.permission_denied",
		})
	}

	// Устанавливаем ID и user_id
	req.ID = id
	req.UserID = userID

	h.logger.Info().
		Int("id", id).
		Str("source_type", req.SourceType).
		Int("user_id", userID).
		Msg("UpdateListing: updating listing")

	// Обновляем через service (используем тот же routingCtx)
	if err := h.service.UpdateListing(c.Context(), &req, routingCtx); err != nil {
		h.logger.Error().Err(err).Int("id", id).Msg("UpdateListing: service error")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.failed_to_update_listing",
		})
	}

	h.logger.Info().Int("id", id).Msg("UpdateListing: success")

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// DeleteListing godoc
// @Summary Delete listing
// @Description Удалить listing (C2C или B2C)
// @Tags unified-marketplace
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Listing ID"
// @Param source_type query string true "Source type: c2c or b2c"
// @Success 200 {object} map[string]bool "success"
// @Failure 400 {object} map[string]string "invalid parameters"
// @Failure 401 {object} map[string]string "unauthorized"
// @Failure 403 {object} map[string]string "forbidden"
// @Failure 404 {object} map[string]string "not found"
// @Failure 500 {object} map[string]string "internal error"
// @Router /api/v1/marketplace/listings/{id} [delete]
func (h *MarketplaceHandler) DeleteListing(c *fiber.Ctx) error {
	// Получаем user_id из JWT
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("DeleteListing: user_id not found in context")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "marketplace.unauthorized",
		})
	}

	// Парсинг ID
	id, err := c.ParamsInt("id")
	if err != nil {
		h.logger.Warn().Err(err).Msg("DeleteListing: invalid ID parameter")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_id",
		})
	}

	// Парсинг source_type
	sourceType := c.Query("source_type", "")
	if sourceType != service.SourceTypeC2C && sourceType != service.SourceTypeB2C {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalid_source_type",
		})
	}

	// Получаем routing context
	routingCtx := h.getRoutingContext(c)

	// Получаем существующий listing для проверки владельца
	existing, err := h.service.GetListing(c.Context(), int64(id), sourceType, routingCtx)
	if err != nil {
		h.logger.Error().Err(err).Int("id", id).Msg("DeleteListing: failed to get existing listing")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "marketplace.listing_not_found",
		})
	}

	// Проверка прав доступа (только владелец может удалять)
	// TODO: добавить проверку роли admin
	if existing.UserID != userID {
		h.logger.Warn().
			Int("listing_user_id", existing.UserID).
			Int("request_user_id", userID).
			Msg("DeleteListing: permission denied")
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "marketplace.permission_denied",
		})
	}

	h.logger.Info().
		Int("id", id).
		Str("source_type", sourceType).
		Int("user_id", userID).
		Msg("DeleteListing: deleting listing")

	// Удаляем через service (используем тот же routingCtx)
	if err := h.service.DeleteListing(c.Context(), int64(id), sourceType, routingCtx); err != nil {
		h.logger.Error().Err(err).Int("id", id).Msg("DeleteListing: service error")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.failed_to_delete_listing",
		})
	}

	h.logger.Info().Int("id", id).Msg("DeleteListing: success")

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// SearchListings godoc
// @Summary Search unified listings
// @Description Поиск по C2C и B2C listings через OpenSearch
// @Tags unified-marketplace
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Param category_id query int false "Category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param condition query string false "Condition: new, used, refurbished"
// @Param source_type query string false "Source type: c2c, b2c, all" default(all)
// @Param storefront_id query int false "Storefront ID (for B2C)"
// @Param user_id query int false "User ID (for C2C)"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param sort_by query string false "Sort field: price, created_at, views" default(created_at)
// @Param sort_order query string false "Sort order: asc, desc" default(desc)
// @Success 200 {object} map[string]interface{} "data, total, limit, offset"
// @Failure 400 {object} map[string]string "validation error"
// @Failure 500 {object} map[string]string "internal error"
// @Router /api/v1/marketplace/search [get]
func (h *MarketplaceHandler) SearchListings(c *fiber.Ctx) error {
	// Парсинг query parameters
	params := &service.SearchParams{
		Query:        c.Query("query", ""),
		CategoryID:   c.QueryInt("category_id", 0),
		MinPrice:     parseFloat(c.Query("min_price", "0")),
		MaxPrice:     parseFloat(c.Query("max_price", "0")),
		Condition:    c.Query("condition", ""),
		SourceType:   c.Query("source_type", "all"),
		StorefrontID: c.QueryInt("storefront_id", 0),
		UserID:       c.QueryInt("user_id", 0),
		Limit:        c.QueryInt("limit", 20),
		Offset:       c.QueryInt("offset", 0),
		SortBy:       c.Query("sort_by", "created_at"),
		SortOrder:    c.Query("sort_order", "desc"),
	}

	// Валидация
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	if params.SourceType != "all" && params.SourceType != service.SourceTypeC2C && params.SourceType != service.SourceTypeB2C {
		params.SourceType = "all"
	}

	h.logger.Info().
		Str("query", params.Query).
		Int("category_id", params.CategoryID).
		Str("source_type", params.SourceType).
		Int("limit", params.Limit).
		Int("offset", params.Offset).
		Msg("SearchListings: request")

	// Получаем routing context
	routingCtx := h.getRoutingContext(c)

	// Поиск через service
	listings, total, err := h.service.SearchListings(c.Context(), params, routingCtx)
	if err != nil {
		h.logger.Error().Err(err).Msg("SearchListings: service error")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.search_failed",
		})
	}

	h.logger.Info().
		Int("count", len(listings)).
		Int64("total", total).
		Msg("SearchListings: success")

	return c.JSON(fiber.Map{
		"success": true,
		"data":    listings,
		"total":   total,
		"limit":   params.Limit,
		"offset":  params.Offset,
	})
}

// parseFloat парсит float64 из string
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// GetPrefix возвращает префикс проекта для логирования
func (h *MarketplaceHandler) GetPrefix() string {
	return "/api/v1/marketplace"
}
