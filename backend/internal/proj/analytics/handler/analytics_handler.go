package handler

import (
	"encoding/json"

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
	if userID, ok := c.Locals("user_id").(int); ok && req.UserID == nil {
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
