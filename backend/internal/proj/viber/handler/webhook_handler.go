package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"backend/internal/proj/viber/models"
	"backend/internal/proj/viber/service"
)


// WebhookHandler обрабатывает webhook события от Viber
type WebhookHandler struct {
	botService     *service.BotService
	sessionManager *service.SessionManager
	messageHandler *MessageHandler
	authToken      string
}

// NewWebhookHandler создаёт новый обработчик webhook
func NewWebhookHandler(
	botService *service.BotService,
	sessionManager *service.SessionManager,
	messageHandler *MessageHandler,
	authToken string,
) *WebhookHandler {
	return &WebhookHandler{
		botService:     botService,
		sessionManager: sessionManager,
		messageHandler: messageHandler,
		authToken:      authToken,
	}
}

// HandleWebhook обрабатывает webhook запрос от Viber
func (h *WebhookHandler) HandleWebhook(c *fiber.Ctx) error {
	// Проверяем подпись (пропускаем в development режиме)
	skipValidation := h.authToken == "dev" || h.authToken == "development"
	if !skipValidation {
		if err := h.validateSignature(c); err != nil {
			fmt.Printf("Webhook signature validation failed: %v\n", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid signature",
			})
		}
	} else {
		fmt.Printf("Skipping signature validation in development mode\n")
	}

	// Парсим событие
	var event models.WebhookEvent
	if err := json.Unmarshal(c.Body(), &event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	// Логируем событие
	fmt.Printf("Received Viber event: %s\n", event.Event)

	// Обрабатываем событие в зависимости от типа
	switch event.Event {
	case "webhook":
		// Первый callback при установке webhook
		return h.handleWebhookVerification(c)

	case "subscribed":
		// Пользователь подписался на бота
		return h.handleSubscribed(c, &event)

	case "unsubscribed":
		// Пользователь отписался от бота
		return h.handleUnsubscribed(c, &event)

	case "conversation_started":
		// Пользователь начал разговор
		return h.handleConversationStarted(c, &event)

	case "message":
		// Получено сообщение от пользователя
		return h.handleMessage(c, &event)

	case "delivered":
		// Сообщение доставлено
		return h.handleDelivered(c, &event)

	case "seen":
		// Сообщение прочитано
		return h.handleSeen(c, &event)

	case "failed":
		// Ошибка доставки сообщения
		return h.handleFailed(c, &event)

	default:
		fmt.Printf("Unknown event type: %s\n", event.Event)
	}

	return c.SendStatus(fiber.StatusOK)
}

// validateSignature проверяет подпись запроса от Viber
func (h *WebhookHandler) validateSignature(c *fiber.Ctx) error {
	signature := c.Get("X-Viber-Content-Signature")
	if signature == "" {
		fmt.Printf("Missing X-Viber-Content-Signature header from IP: %s\n", c.IP())
		return fmt.Errorf("missing signature header")
	}

	// Проверяем что auth token настроен
	if h.authToken == "" {
		fmt.Printf("Viber auth token not configured\n")
		return fmt.Errorf("auth token not configured")
	}

	body := c.Body()
	if len(body) == 0 {
		fmt.Printf("Empty body in webhook request from IP: %s\n", c.IP())
		return fmt.Errorf("empty request body")
	}

	// Вычисляем HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(h.authToken))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Для отладки (в продакшене убрать)
	fmt.Printf("Received signature: %s\n", signature)
	fmt.Printf("Expected signature: %s\n", expectedSignature)
	fmt.Printf("Body length: %d bytes\n", len(body))

	// Сравниваем подписи безопасным способом
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		fmt.Printf("Signature validation failed from IP: %s\n", c.IP())
		return fmt.Errorf("invalid signature")
	}

	fmt.Printf("Signature validation successful from IP: %s\n", c.IP())
	return nil
}

// handleWebhookVerification обрабатывает первую верификацию webhook
func (h *WebhookHandler) handleWebhookVerification(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": 0,
		"status_message": "ok",
	})
}

// handleSubscribed обрабатывает подписку пользователя
func (h *WebhookHandler) handleSubscribed(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()

	if event.User != nil {
		// Сохраняем информацию о пользователе
		if err := h.sessionManager.SaveUserInfo(ctx, event.User); err != nil {
			fmt.Printf("Failed to save user info: %v\n", err)
		}

		// Отмечаем как подписанного
		if err := h.sessionManager.SetUserSubscribed(ctx, event.User.ID, true); err != nil {
			fmt.Printf("Failed to set user subscribed: %v\n", err)
		}

		// Отправляем приветственное сообщение
		welcomeMsg := "🎉 Добро пожаловать в SveTu Marketplace!\n\n" +
			"Я помогу вам:\n" +
			"• 🔍 Найти товары\n" +
			"• 📦 Отследить доставку\n" +
			"• 🏪 Управлять витриной\n" +
			"• 💬 Связаться с продавцом\n\n" +
			"Выберите действие из меню ниже 👇"

		if err := h.botService.SendKeyboard(ctx, event.User.ID, welcomeMsg,
			h.botService.GetMainMenuKeyboard()); err != nil {
			fmt.Printf("Failed to send welcome message: %v\n", err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// handleUnsubscribed обрабатывает отписку пользователя
func (h *WebhookHandler) handleUnsubscribed(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()

	if event.User != nil {
		// Отмечаем как отписанного
		if err := h.sessionManager.SetUserSubscribed(ctx, event.User.ID, false); err != nil {
			fmt.Printf("Failed to set user unsubscribed: %v\n", err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// handleConversationStarted обрабатывает начало разговора
func (h *WebhookHandler) handleConversationStarted(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()

	if event.User != nil {
		// Отмечаем начало разговора
		if err := h.sessionManager.SetConversationStarted(ctx, event.User.ID); err != nil {
			fmt.Printf("Failed to set conversation started: %v\n", err)
		}

		// Отправляем приветственное сообщение
		welcomeMsg := "Привет! 👋\n\nЯ бот маркетплейса SveTu.\nЧем могу помочь?"

		// Создаём Rich Media с кнопками быстрого старта
		richMedia := &models.RichMedia{
			Type:                "rich_media",
			ButtonsGroupColumns: 6,
			ButtonsGroupRows:    2,
			Buttons: []models.RichButton{
				{
					Columns:    3,
					Rows:       1,
					ActionType: "reply",
					ActionBody: "search",
					Text:       "🔍 Поиск товаров",
					TextSize:   "medium",
					TextVAlign: "middle",
					TextHAlign: "center",
					BgColor:    "#e3f2fd",
				},
				{
					Columns:    3,
					Rows:       1,
					ActionType: "reply",
					ActionBody: "my_orders",
					Text:       "📦 Мои заказы",
					TextSize:   "medium",
					TextVAlign: "middle",
					TextHAlign: "center",
					BgColor:    "#f3e5f5",
				},
				{
					Columns:    6,
					Rows:       1,
					ActionType: "open-url",
					ActionBody: "https://svetu.rs",
					Text:       "🌐 Открыть сайт",
					TextSize:   "medium",
					TextVAlign: "middle",
					TextHAlign: "center",
					BgColor:    "#e8f5e9",
				},
			},
		}

		if err := h.botService.SendRichMedia(ctx, event.User.ID, richMedia, welcomeMsg); err != nil {
			fmt.Printf("Failed to send conversation started message: %v\n", err)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

// handleMessage обрабатывает входящее сообщение
func (h *WebhookHandler) handleMessage(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()

	if event.Sender == nil || event.Message == nil {
		return c.SendStatus(fiber.StatusOK)
	}

	// Сохраняем информацию о пользователе
	if err := h.sessionManager.SaveUserInfo(ctx, event.Sender); err != nil {
		fmt.Printf("Failed to save sender info: %v\n", err)
	}

	// Создаём или обновляем сессию
	session, err := h.sessionManager.GetActiveSession(ctx, event.Sender.ID)
	if err != nil {
		fmt.Printf("Failed to get session: %v\n", err)
		return c.SendStatus(fiber.StatusOK)
	}

	if session == nil {
		// Создаём новую сессию
		session, err = h.sessionManager.CreateSession(ctx, event.Sender.ID)
		if err != nil {
			fmt.Printf("Failed to create session: %v\n", err)
			return c.SendStatus(fiber.StatusOK)
		}
	} else {
		// Обновляем существующую сессию
		if err := h.sessionManager.UpdateSession(ctx, session.ID); err != nil {
			fmt.Printf("Failed to update session: %v\n", err)
		}
	}

	// Обрабатываем сообщение в зависимости от типа
	switch event.Message.Type {
	case "text":
		return h.handleTextMessage(c, event)
	case "picture":
		return h.handlePictureMessage(c, event)
	case "location":
		return h.handleLocationMessage(c, event)
	default:
		// Неподдерживаемый тип сообщения
		msg := "Извините, я пока не умею обрабатывать такие сообщения. " +
			"Пожалуйста, используйте текстовые сообщения или кнопки меню."
		h.botService.SendTextMessage(ctx, event.Sender.ID, msg)
	}

	return c.SendStatus(fiber.StatusOK)
}

// handleTextMessage обрабатывает текстовое сообщение
func (h *WebhookHandler) handleTextMessage(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()
	text := strings.TrimSpace(strings.ToLower(event.Message.Text))

	// Обрабатываем команды
	switch {
	case text == "search" || strings.Contains(text, "поиск") || strings.Contains(text, "найти"):
		return h.messageHandler.HandleSearch(ctx, event.Sender.ID, text)

	case text == "my_orders" || strings.Contains(text, "заказ"):
		return h.messageHandler.HandleMyOrders(ctx, event.Sender.ID)

	case text == "cart" || strings.Contains(text, "корзин"):
		return h.messageHandler.HandleCart(ctx, event.Sender.ID)

	case text == "storefronts" || strings.Contains(text, "витрин") || strings.Contains(text, "магазин"):
		return h.messageHandler.HandleStorefronts(ctx, event.Sender.ID)

	case text == "help" || strings.Contains(text, "помощь") || text == "/help":
		return h.messageHandler.HandleHelp(ctx, event.Sender.ID)

	case strings.HasPrefix(text, "track_"):
		// Отслеживание доставки
		trackingToken := strings.TrimPrefix(text, "track_")
		return h.messageHandler.HandleTrackDelivery(ctx, event.Sender.ID, trackingToken)

	default:
		// Пытаемся понять намерение пользователя через поиск
		if len(text) > 3 {
			return h.messageHandler.HandleSearch(ctx, event.Sender.ID, text)
		}

		// Не поняли команду
		msg := "Не понял вашу команду. Выберите действие из меню или напишите 'помощь'."
		return h.botService.SendKeyboard(ctx, event.Sender.ID, msg,
			h.botService.GetMainMenuKeyboard())
	}
}

// handlePictureMessage обрабатывает изображение
func (h *WebhookHandler) handlePictureMessage(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()

	msg := "Вы отправили изображение. Если хотите добавить товар с фото, " +
		"пожалуйста, используйте наш сайт: https://svetu.rs/create-listing"

	return h.botService.SendTextMessage(ctx, event.Sender.ID, msg)
}

// handleLocationMessage обрабатывает локацию
func (h *WebhookHandler) handleLocationMessage(c *fiber.Ctx, event *models.WebhookEvent) error {
	ctx := c.Context()

	if event.Message.Location != nil {
		lat := event.Message.Location.Latitude
		lng := event.Message.Location.Longitude

		// Ищем товары рядом
		return h.messageHandler.HandleNearbyProducts(ctx, event.Sender.ID, lat, lng)
	}

	return c.SendStatus(fiber.StatusOK)
}

// handleDelivered обрабатывает подтверждение доставки
func (h *WebhookHandler) handleDelivered(c *fiber.Ctx, event *models.WebhookEvent) error {
	// Обновляем статус сообщения в БД
	fmt.Printf("Message delivered: %s\n", event.MessageToken)
	return c.SendStatus(fiber.StatusOK)
}

// handleSeen обрабатывает прочтение сообщения
func (h *WebhookHandler) handleSeen(c *fiber.Ctx, event *models.WebhookEvent) error {
	// Обновляем статус сообщения в БД
	fmt.Printf("Message seen: %s\n", event.MessageToken)
	return c.SendStatus(fiber.StatusOK)
}

// handleFailed обрабатывает ошибку доставки
func (h *WebhookHandler) handleFailed(c *fiber.Ctx, event *models.WebhookEvent) error {
	// Логируем ошибку
	fmt.Printf("Message failed: %s\n", event.MessageToken)
	return c.SendStatus(fiber.StatusOK)
}