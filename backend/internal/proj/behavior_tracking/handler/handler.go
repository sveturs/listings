package handler

import (
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"

	"backend/internal/domain/behavior"
	"backend/internal/logger"
	"backend/internal/proj/behavior_tracking/service"
	"backend/pkg/utils"
)

// BehaviorTrackingHandler обработчик для поведенческих событий
type BehaviorTrackingHandler struct {
	service   service.BehaviorTrackingService
	validator *validator.Validate
	limiter   *rate.Limiter
}

// NewBehaviorTrackingHandler создает новый обработчик
func NewBehaviorTrackingHandler(service service.BehaviorTrackingService) *BehaviorTrackingHandler {
	return &BehaviorTrackingHandler{
		service:   service,
		validator: validator.New(),
		limiter:   rate.NewLimiter(rate.Every(time.Second), 100), // 100 событий в секунду на IP
	}
}

// TrackEvent обрабатывает запрос на отслеживание события или пакета событий
// @Summary Track behavior events
// @Description Records user behavior events for analytics (supports both single event and batch)
// @Tags analytics
// @Accept json
// @Produce json
// @Param event body behavior.TrackEventBatch true "Event batch data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=map[string]interface{}} "Events tracked successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Failure 429 {object} backend_pkg_utils.ErrorResponseSwag "Too many requests"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/track [post]
func (h *BehaviorTrackingHandler) TrackEvent(c *fiber.Ctx) error {
	// Rate limiting по IP
	if !h.limiter.Allow() {
		return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "analytics.error.rate_limit_exceeded")
	}

	// Получаем user_id из контекста (если пользователь авторизован)
	var userID *int
	if uid, ok := c.Locals("user_id").(int); ok {
		userID = &uid
	}

	// Получаем метаданные из контекста запроса
	contextMetadata := map[string]interface{}{
		"ip":         c.IP(),
		"user_agent": c.Get("User-Agent"),
		"referer":    c.Get("Referer"),
		"timestamp":  time.Now().Unix(),
	}

	// Пробуем распарсить как batch
	var batch behavior.TrackEventBatch
	if err := c.BodyParser(&batch); err == nil && len(batch.Events) > 0 {
		// Это batch событий
		logger.Info().
			Str("batch_id", batch.BatchID).
			Int("events_count", len(batch.Events)).
			Msg("Processing event batch")

		// Валидация batch
		if err := h.validator.Struct(&batch); err != nil {
			logger.Error().Err(err).Msg("Validation failed for event batch")
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.validation_failed")
		}

		// Обрабатываем каждое событие
		processedCount := 0
		var lastSessionID string

		for _, event := range batch.Events {
			// Добавляем метаданные к каждому событию
			if event.Metadata == nil {
				event.Metadata = make(map[string]interface{})
			}
			for k, v := range contextMetadata {
				event.Metadata[k] = v
			}

			// Генерируем session_id если не передан
			if event.SessionID == "" && lastSessionID != "" {
				event.SessionID = lastSessionID
			} else if event.SessionID != "" {
				lastSessionID = event.SessionID
			}

			// Отслеживаем событие
			if err := h.service.TrackEvent(c.Context(), userID, &event); err != nil {
				logger.Error().
					Err(err).
					Str("event_type", string(event.EventType)).
					Msg("Failed to track event in batch")
				// Продолжаем обработку остальных событий
				continue
			}
			processedCount++
		}

		return utils.SuccessResponse(c, fiber.Map{
			"message":         "Events batch processed",
			"batch_id":        batch.BatchID,
			"processed_count": processedCount,
			"failed_count":    len(batch.Events) - processedCount,
		})
	}

	// Если не удалось распарсить как batch, пробуем как одиночное событие
	var req behavior.TrackEventRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to parse track event request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.invalid_request")
	}

	// Валидация
	if err := h.validator.Struct(&req); err != nil {
		logger.Error().Err(err).Msg("Validation failed for track event request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.validation_failed")
	}

	// Добавляем метаданные
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	for k, v := range contextMetadata {
		req.Metadata[k] = v
	}

	// Отслеживаем событие
	if err := h.service.TrackEvent(c.Context(), userID, &req); err != nil {
		logger.Error().Err(err).Msg("Failed to track event")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_track")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":    "Event tracked successfully",
		"session_id": req.SessionID,
	})
}

// SearchMetricsResponse представляет ответ с агрегированными метриками поиска
type SearchMetricsResponse struct {
	TotalSearches           int64                 `json:"total_searches"`
	UniqueSearches          int64                 `json:"unique_searches"`
	AverageSearchDurationMs float64               `json:"average_search_duration_ms"`
	TopQueries              []TopQueryResponse    `json:"top_queries"`
	SearchTrends            []SearchTrendResponse `json:"search_trends"`
	ClickMetrics            ClickMetricsResponse  `json:"click_metrics"`
}

type TopQueryResponse struct {
	Query       string  `json:"query"`
	Count       int     `json:"count"`
	CTR         float64 `json:"ctr"`
	AvgPosition float64 `json:"avg_position"`
	AvgResults  float64 `json:"avg_results"`
}

type SearchTrendResponse struct {
	Date          string  `json:"date"`
	SearchesCount int     `json:"searches_count"`
	ClicksCount   int     `json:"clicks_count"`
	CTR           float64 `json:"ctr"`
}

type ClickMetricsResponse struct {
	TotalClicks          int     `json:"total_clicks"`
	AverageClickPosition float64 `json:"average_click_position"`
	CTR                  float64 `json:"ctr"`
	ConversionRate       float64 `json:"conversion_rate"`
}

// GetUserEvents возвращает события пользователя
// @Summary Get user events
// @Description Returns behavior events for a specific user
// @Tags analytics
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param limit query int false "Limit results (default: 20, max: 100)"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]behavior.BehaviorEvent} "User events"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Forbidden"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/users/{user_id}/events [get]
func (h *BehaviorTrackingHandler) GetUserEvents(c *fiber.Ctx) error {
	// Парсим user_id из пути
	userID, err := strconv.Atoi(c.Params("user_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.invalid_user_id")
	}

	// Проверяем права доступа (пользователь может смотреть только свои события или админ)
	currentUserID, ok := c.Locals("user_id").(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "errors.unauthorized")
	}

	isAdmin := authMiddleware.IsAdmin(c)
	if currentUserID != userID && !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "errors.forbidden")
	}

	// Парсим параметры пагинации
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Получаем события
	events, total, err := h.service.GetUserEvents(c.Context(), userID, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get user events")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_get_events")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetSessionEvents возвращает события сессии
// @Summary Get session events
// @Description Returns behavior events for a specific session
// @Tags analytics
// @Accept json
// @Produce json
// @Param session_id path string true "Session ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]behavior.BehaviorEvent} "Session events"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/sessions/{session_id}/events [get]
func (h *BehaviorTrackingHandler) GetSessionEvents(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")
	if sessionID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.invalid_session_id")
	}

	// Получаем события
	events, err := h.service.GetSessionEvents(c.Context(), sessionID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get session events")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_get_events")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"events": events,
		"total":  len(events),
	})
}

// UpdateMetrics принудительно обновляет агрегированные метрики
// @Summary Update aggregated metrics
// @Description Forces update of aggregated search metrics for a specific period
// @Tags analytics
// @Accept json
// @Produce json
// @Param period_start query string true "Period start date (RFC3339)"
// @Param period_end query string true "Period end date (RFC3339)"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=map[string]string} "Metrics updated successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Admin access required"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/analytics/metrics/update [post]
// @Security BearerAuth
func (h *BehaviorTrackingHandler) UpdateMetrics(c *fiber.Ctx) error {
	// Только админы могут обновлять метрики
	isAdmin := authMiddleware.IsAdmin(c)
	if !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "errors.admin_required")
	}

	// Парсим даты
	periodStartStr := c.Query("period_start")
	periodEndStr := c.Query("period_end")

	if periodStartStr == "" || periodEndStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.period_required")
	}

	periodStart, err := time.Parse(time.RFC3339, periodStartStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.invalid_period_start")
	}

	periodEnd, err := time.Parse(time.RFC3339, periodEndStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "analytics.error.invalid_period_end")
	}

	// Обновляем метрики
	if err := h.service.UpdateSearchMetrics(c.Context(), periodStart, periodEnd); err != nil {
		logger.Error().Err(err).Msg("Failed to update search metrics")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "analytics.error.failed_to_update_metrics")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":      "Metrics updated successfully",
		"period_start": periodStart.Format(time.RFC3339),
		"period_end":   periodEnd.Format(time.RFC3339),
	})
}
