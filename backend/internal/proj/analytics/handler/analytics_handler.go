package handler

import (
	"encoding/json"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/logger"
	"backend/internal/proj/analytics/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsHandler обработчик для аналитики
type AnalyticsHandler struct {
	service service.AnalyticsService
}

// NewAnalyticsHandler создает новый обработчик
func NewAnalyticsHandler(service service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

// EventRequest запрос на запись события
type EventRequest struct {
	StorefrontID int             `json:"storefront_id" validate:"required"`
	EventType    string          `json:"event_type" validate:"required,oneof=page_view product_view add_to_cart checkout order"`
	EventData    json.RawMessage `json:"event_data,omitempty"`
	SessionID    string          `json:"session_id" validate:"required"`
	UserID       *int            `json:"user_id,omitempty"`
}

// RecordEvent записывает событие аналитики
// @Summary Record analytics event
// @Description Records an analytics event for a storefront
// @Tags analytics
// @Accept json
// @Produce json
// @Param event body EventRequest true "Event data"
// @Success 200 {object} utils.SuccessResponseSwag "Event recorded"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/event [post]
func (h *AnalyticsHandler) RecordEvent(c *fiber.Ctx) error {
	var req EventRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.invalid_request")
	}

	// Валидация
	if req.StorefrontID <= 0 || req.EventType == "" || req.SessionID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.validation_failed")
	}

	// Получаем IP и User-Agent
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")
	referrer := c.Get("Referer")

	// Если пользователь авторизован, добавляем user_id
	if userID, ok := authMiddleware.GetUserID(c); ok && req.UserID == nil {
		req.UserID = &userID
	}

	// Записываем событие
	err := h.service.RecordEvent(c.Context(), &service.EventData{
		StorefrontID: req.StorefrontID,
		EventType:    req.EventType,
		EventData:    req.EventData,
		SessionID:    req.SessionID,
		UserID:       req.UserID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		Referrer:     referrer,
	})
	if err != nil {
		logger.Error().Err(err).Str("event_type", req.EventType).Msg("Failed to record event")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_record")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Event recorded successfully",
	})
}

// GetSearchMetrics возвращает метрики поиска
// @Summary Get search metrics
// @Description Returns search analytics metrics
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param period query string false "Period (day, week, month)" default(week)
// @Success 200 {object} utils.SuccessResponseSwag{data=service.SearchMetrics} "Search metrics"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/metrics/search [get]
func (h *AnalyticsHandler) GetSearchMetrics(c *fiber.Ctx) error {
	// Проверяем права админа - используем тот же подход как в search_admin
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Hardcoded admin users (как в других админских модулях)
	adminUsers := []int{1, 2, 3, 4, 5, 6}
	isAdmin := false
	for _, adminID := range adminUsers {
		if userID == adminID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin.error.access_denied")
	}

	dateFrom := c.Query("date_from", "")
	dateTo := c.Query("date_to", "")
	period := c.Query("period", "week")

	metrics, err := h.service.GetSearchMetrics(c.Context(), dateFrom, dateTo, period)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get search metrics")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_get_metrics")
	}

	return utils.SuccessResponse(c, metrics)
}

// GetItemsPerformance возвращает производительность товаров
// @Summary Get items performance
// @Description Returns performance metrics for items
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]service.ItemPerformance} "Items performance"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/metrics/items [get]
func (h *AnalyticsHandler) GetItemsPerformance(c *fiber.Ctx) error {
	// Проверяем права админа - используем тот же подход как в search_admin
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Hardcoded admin users (как в других админских модулях)
	adminUsers := []int{1, 2, 3, 4, 5, 6}
	isAdmin := false
	for _, adminID := range adminUsers {
		if userID == adminID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin.error.access_denied")
	}

	dateFrom := c.Query("date_from", "")
	dateTo := c.Query("date_to", "")
	limit := c.QueryInt("limit", 20)

	items, err := h.service.GetItemsPerformance(c.Context(), dateFrom, dateTo, limit)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items performance")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_get_items")
	}

	return utils.SuccessResponse(c, items)
}
