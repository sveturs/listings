package viber

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"backend/internal/proj/viber/config"
	"backend/internal/proj/viber/handler"
	"backend/internal/proj/viber/service"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"backend/internal/storage/postgres"
)

// ViberHandler обрабатывает запросы Viber Bot
type ViberHandler struct {
	webhookHandler   *handler.WebhookHandler
	botService       *service.BotService
	infobipService   *service.InfobipBotService
	sessionManager   *service.SessionManager
	config           *config.ViberConfig
}

// NewViberHandler создаёт новый обработчик Viber
func NewViberHandler(db *postgres.Database, services globalService.ServicesInterface) *ViberHandler {
	cfg := config.LoadViberConfig()
	sessionManager := service.NewSessionManager(db)

	// Выбираем сервис в зависимости от конфигурации
	var botService *service.BotService
	var infobipService *service.InfobipBotService

	if cfg.UseInfobip {
		infobipService = service.NewInfobipBotService(cfg, db)
	} else {
		botService = service.NewBotService(cfg, db)
	}

	// Создаем полноценный MessageHandler с доступом к сервисам
	messageHandler := handler.NewMessageHandler(
		botService,
		infobipService,
		services,
		services.Marketplace(),
		services.Storefront(),
		cfg.UseInfobip,
	)

	webhookHandler := handler.NewWebhookHandler(
		botService,
		sessionManager,
		messageHandler,
		cfg.AuthToken,
	)

	return &ViberHandler{
		webhookHandler:   webhookHandler,
		botService:       botService,
		infobipService:   infobipService,
		sessionManager:   sessionManager,
		config:           cfg,
	}
}

// HandleViberWebhook обрабатывает webhook от Viber API
// @Summary Handle Viber webhook
// @Description Processes webhook events from Viber (subscriptions, messages, etc.)
// @Tags viber
// @Accept json
// @Produce json
// @Param X-Viber-Content-Signature header string true "Viber signature"
// @Param event body object true "Viber webhook event"
// @Success 200 {object} utils.SuccessResponseSwag "Event processed"
// @Failure 401 {object} utils.ErrorResponseSwag "Invalid signature"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid event"
// @Router /api/viber/webhook [post]
func (h *ViberHandler) HandleViberWebhook(c *fiber.Ctx) error {
	return h.webhookHandler.HandleWebhook(c)
}

// HandleInfobipWebhook обрабатывает webhook от Infobip
// @Summary Handle Infobip webhook
// @Description Processes webhook events from Infobip Viber API
// @Tags viber
// @Accept json
// @Produce json
// @Param webhook body object true "Infobip webhook event"
// @Success 200 {object} utils.SuccessResponseSwag "Event processed"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid event"
// @Router /api/viber/infobip-webhook [post]
func (h *ViberHandler) HandleInfobipWebhook(c *fiber.Ctx) error {
	if h.infobipService == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "viber.error.infobipNotConfigured")
	}

	var webhook map[string]interface{}
	if err := c.BodyParser(&webhook); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	// TODO: Обработка Infobip webhook
	// Здесь нужно реализовать обработку webhook'ов от Infobip

	return utils.SuccessResponse(c, nil)
}

// SendMessage отправляет сообщение пользователю
// @Summary Send message to Viber user
// @Description Sends a text message to a Viber user
// @Tags viber
// @Accept json
// @Produce json
// @Param message body SendMessageRequest true "Message data"
// @Success 200 {object} utils.SuccessResponseSwag "Message sent"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Send failed"
// @Router /api/viber/send [post]
func (h *ViberHandler) SendMessage(c *fiber.Ctx) error {
	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	if req.ViberID == "" || req.Text == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "viber.error.missingFields")
	}

	// Отправляем через подходящий сервис
	var err error
	if h.config.UseInfobip && h.infobipService != nil {
		err = h.infobipService.SendTextMessage(c.Context(), req.ViberID, req.Text)
	} else if h.botService != nil {
		err = h.botService.SendTextMessage(c.Context(), req.ViberID, req.Text)
	} else {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "viber.error.serviceNotAvailable")
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "viber.error.sendFailed")
	}

	return utils.SuccessResponse(c, nil)
}

// SendTrackingNotification отправляет уведомление о трекинге
// @Summary Send tracking notification
// @Description Sends a rich media message with tracking information
// @Tags viber
// @Accept json
// @Produce json
// @Param notification body SendTrackingRequest true "Tracking data"
// @Success 200 {object} utils.SuccessResponseSwag "Notification sent"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Send failed"
// @Router /api/viber/send-tracking [post]
func (h *ViberHandler) SendTrackingNotification(c *fiber.Ctx) error {
	var req SendTrackingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	if req.ViberID == "" || req.DeliveryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "viber.error.missingFields")
	}

	// TODO: Получаем информацию о доставке
	// deliveryInfo := getDeliveryInfo(req.DeliveryID)

	// Отправляем через подходящий сервис
	var err error
	if h.config.UseInfobip && h.infobipService != nil {
		// err = h.infobipService.SendTrackingNotification(c.Context(), req.ViberID, deliveryInfo)
		err = h.infobipService.SendTextMessage(c.Context(), req.ViberID, "🚚 Ваш заказ в пути! Отследить: https://svetu.rs/track/"+strconv.Itoa(req.DeliveryID))
	} else if h.botService != nil {
		// err = h.botService.SendTrackingNotification(c.Context(), req.ViberID, deliveryInfo)
		err = h.botService.SendTextMessage(c.Context(), req.ViberID, "🚚 Ваш заказ в пути! Отследить: https://svetu.rs/track/"+strconv.Itoa(req.DeliveryID))
	} else {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "viber.error.serviceNotAvailable")
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "viber.error.sendFailed")
	}

	return utils.SuccessResponse(c, nil)
}

// GetSessionStats получает статистику сессий
// @Summary Get session statistics
// @Description Returns statistics about Viber bot sessions
// @Tags viber
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Session statistics"
// @Failure 500 {object} utils.ErrorResponseSwag "Failed to get statistics"
// @Router /api/viber/stats [get]
func (h *ViberHandler) GetSessionStats(c *fiber.Ctx) error {
	stats, err := h.sessionManager.GetSessionStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "viber.error.statsFailed")
	}

	return utils.SuccessResponse(c, stats)
}

// EstimateMessageCost оценивает стоимость сообщения
// @Summary Estimate message cost
// @Description Estimates the cost of sending a message to a Viber user
// @Tags viber
// @Accept json
// @Produce json
// @Param estimate body EstimateCostRequest true "Cost estimation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]float64} "Estimated cost"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Estimation failed"
// @Router /api/viber/estimate-cost [post]
func (h *ViberHandler) EstimateMessageCost(c *fiber.Ctx) error {
	var req EstimateCostRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	if req.ViberID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "viber.error.missingViberID")
	}

	var cost float64
	var err error

	if h.config.UseInfobip && h.infobipService != nil {
		cost, err = h.infobipService.EstimateMessageCost(c.Context(), req.ViberID, req.IsRichMedia)
	} else {
		// Прямой Viber API бесплатный
		cost = 0
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "viber.error.estimationFailed")
	}

	result := map[string]float64{
		"estimated_cost": cost,
		"currency":       0, // EUR
	}

	return utils.SuccessResponse(c, result)
}

// DTOs для запросов

// SendMessageRequest запрос на отправку сообщения
type SendMessageRequest struct {
	ViberID string `json:"viber_id" validate:"required"`
	Text    string `json:"text" validate:"required,min=1,max=7000"`
}

// SendTrackingRequest запрос на отправку уведомления о трекинге
type SendTrackingRequest struct {
	ViberID    string `json:"viber_id" validate:"required"`
	DeliveryID int    `json:"delivery_id" validate:"required,min=1"`
}

// EstimateCostRequest запрос на оценку стоимости
type EstimateCostRequest struct {
	ViberID     string `json:"viber_id" validate:"required"`
	IsRichMedia bool   `json:"is_rich_media"`
}