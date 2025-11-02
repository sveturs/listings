package viber

import (
	"context"
	"fmt"
	"strings"

	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/viber/config"
	"backend/internal/proj/viber/handler"
	"backend/internal/proj/viber/infobip"
	"backend/internal/proj/viber/service"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ViberHandler обрабатывает запросы Viber Bot
type ViberHandler struct {
	webhookHandler *handler.WebhookHandler
	messageHandler *handler.MessageHandler
	botService     *service.BotService
	infobipService *service.InfobipBotService
	sessionManager *service.SessionManager
	config         *config.ViberConfig
	db             *postgres.Database
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
		cfg.UseInfobip,
		cfg,
	)

	webhookHandler := handler.NewWebhookHandler(
		botService,
		sessionManager,
		messageHandler,
		cfg.AuthToken,
		cfg,
	)

	return &ViberHandler{
		webhookHandler: webhookHandler,
		messageHandler: messageHandler,
		botService:     botService,
		infobipService: infobipService,
		sessionManager: sessionManager,
		config:         cfg,
		db:             db,
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

	// Infobip отправляет массив results
	type InfobipWebhookPayload struct {
		Results []infobip.ViberWebhook `json:"results"`
	}

	var payload InfobipWebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "common.error.invalidJSON")
	}

	ctx := c.Context()

	// Обрабатываем каждое событие в массиве
	for _, webhook := range payload.Results {
		// Обрабатываем через сервис (сохранение в БД)
		if err := h.infobipService.ProcessWebhook(ctx, &webhook); err != nil {
			// Логируем ошибку, но продолжаем обработку
			fmt.Printf("Failed to process webhook: %v\n", err)
		}

		// ВАЖНО: Если это входящее сообщение, обрабатываем через MessageHandler
		if webhook.InboundContent != nil && webhook.InboundContent.Text != "" {
			viberID := webhook.From
			text := webhook.InboundContent.Text

			fmt.Printf("Received inbound message from %s: %s\n", viberID, text)

			// Обрабатываем команды
			text = strings.TrimSpace(strings.ToLower(text))

			switch {
			case text == "help" || strings.Contains(text, "помощь") || text == "/help" || text == "start":
				_ = h.messageHandler.HandleHelp(ctx, viberID)

			case text == "search" || strings.Contains(text, "поиск") || strings.Contains(text, "найти"):
				_ = h.messageHandler.HandleSearch(ctx, viberID, text)

			case text == "my_orders" || strings.Contains(text, "заказ"):
				_ = h.messageHandler.HandleMyOrders(ctx, viberID)

			case text == "cart" || strings.Contains(text, "корзин"):
				_ = h.messageHandler.HandleCart(ctx, viberID)

			case text == "storefronts" || strings.Contains(text, "витрин") || strings.Contains(text, "магазин"):
				_ = h.messageHandler.HandleStorefronts(ctx, viberID)

			case strings.HasPrefix(text, "track_"):
				trackingToken := strings.TrimPrefix(text, "track_")
				_ = h.messageHandler.HandleTrackDelivery(ctx, viberID, trackingToken)

			default:
				// Пытаемся понять намерение через поиск
				if len(text) > 3 {
					_ = h.messageHandler.HandleSearch(ctx, viberID, text)
				} else {
					// Приветствие или первое сообщение
					welcomeMsg := "👋 Привет! Я бот SveTu Marketplace.\n\n" +
						"Я помогу вам:\n" +
						"🔍 Найти товары - просто напишите что ищете\n" +
						"🏪 Показать витрины - напишите 'витрины'\n" +
						"📦 Отследить доставку - напишите 'track_НОМЕР'\n" +
						"❓ Помощь - напишите 'помощь'\n\n" +
						"Что вас интересует?"
					_ = h.infobipService.SendTextMessage(ctx, viberID, welcomeMsg)
				}
			}
		}
	}

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
	//nolint:gocritic // if-else chain is appropriate here for service selection
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

	// Получаем информацию о доставке
	deliveryInfo, err := h.getDeliveryInfoByID(c.Context(), req.DeliveryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "delivery.error.notFound")
	}

	// Отправляем через подходящий сервис
	//nolint:gocritic // if-else chain is appropriate here for service selection
	if h.config.UseInfobip && h.infobipService != nil {
		err = h.infobipService.SendTrackingNotification(c.Context(), req.ViberID, deliveryInfo)
	} else if h.botService != nil {
		err = h.botService.SendTrackingNotification(c.Context(), req.ViberID, deliveryInfo)
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

// getDeliveryInfoByID получает информацию о доставке из БД
func (h *ViberHandler) getDeliveryInfoByID(ctx context.Context, deliveryID int) (*service.DeliveryInfo, error) {
	query := `
		SELECT
			d.id,
			d.order_id,
			d.tracking_token,
			d.delivery_address,
			d.delivery_latitude,
			d.delivery_longitude,
			d.estimated_delivery_time,
			COALESCE(clh.latitude, d.pickup_latitude) as courier_latitude,
			COALESCE(clh.longitude, d.pickup_longitude) as courier_longitude,
			COALESCE(u.full_name, 'Courier') as courier_name
		FROM deliveries d
		LEFT JOIN users u ON u.id = d.courier_id
		LEFT JOIN LATERAL (
			SELECT latitude, longitude
			FROM courier_location_history
			WHERE courier_id = d.courier_id
			ORDER BY created_at DESC
			LIMIT 1
		) clh ON true
		WHERE d.id = $1
	`

	var info service.DeliveryInfo
	var courierLatitude, courierLongitude, deliveryLatitude, deliveryLongitude *float64

	err := h.db.QueryRowContext(ctx, query, deliveryID).Scan(
		&info.ID,
		&info.OrderID,
		&info.TrackingToken,
		&info.DeliveryAddress,
		&deliveryLatitude,
		&deliveryLongitude,
		&info.EstimatedTime,
		&courierLatitude,
		&courierLongitude,
		&info.CourierName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery info: %w", err)
	}

	// Конвертируем nullable поля
	if courierLatitude != nil {
		info.CourierLatitude = *courierLatitude
	}
	if courierLongitude != nil {
		info.CourierLongitude = *courierLongitude
	}
	if deliveryLatitude != nil {
		info.DeliveryLatitude = *deliveryLatitude
	}
	if deliveryLongitude != nil {
		info.DeliveryLongitude = *deliveryLongitude
	}

	return &info, nil
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
