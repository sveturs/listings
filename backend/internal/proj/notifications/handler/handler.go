// Package handler
// backend/internal/proj/notifications/handler/handler.go
package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/notifications/service"
	"backend/pkg/utils"
)

type Handler struct {
	notificationService service.NotificationServiceInterface
	bot                 *tgbotapi.BotAPI
}

func NewHandler(service service.NotificationServiceInterface) *Handler {
	handler := &Handler{
		notificationService: service,
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err == nil {
		handler.bot = bot
	}

	return handler
}

// GetNotifications возвращает список уведомлений пользователя
// @Summary Получить уведомления
// @Description Возвращает список уведомлений пользователя с пагинацией
// @Tags notifications
// @Accept json
// @Produce json
// @Param limit query int false "Количество записей" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Список уведомлений"
// @Failure 500 {object} utils.ErrorResponseSwag "notifications.getError"
// @Security BearerAuth
// @Router /api/v1/notifications [get]
func (h *Handler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	notifications, err := h.notificationService.GetUserNotifications(c.Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "notifications.getError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": notifications, // Изменено с notifications на data
	})
}

func (h *Handler) ConnectTelegramWebhook() {
	if h.bot == nil {
		return
	}

	baseURL := "https://svetu.rs/api/v1/notifications/telegram/webhook"

	// Создаем конфигурацию вебхука
	webhookConfig := tgbotapi.NewWebhook(baseURL)

	// Устанавливаем вебхук
	_, err := h.bot.SetWebhook(webhookConfig)
	if err != nil {
		return
	}

	// Опционально: проверяем информацию о вебхуке
	info, err := h.bot.GetWebhookInfo()
	if err != nil {
		return
	}

	if info.LastErrorDate != 0 {
		logger.Error().Str("error_message", info.LastErrorMessage).Msg("Telegram callback failed")
	}

}

// HandleTelegramWebhook обрабатывает webhook от Telegram бота
// @Summary Telegram webhook
// @Description Обрабатывает входящие обновления от Telegram бота
// @Tags notifications
// @Accept json
// @Produce json
// @Param update body tgbotapi.Update true "Telegram update"
// @Success 200 {string} string "OK"
// @Router /api/v1/notifications/telegram/webhook [post]
func (h *Handler) HandleTelegramWebhook(c *fiber.Ctx) error {
	logger.Info().Str("body", string(c.Body())).Msg("Received webhook request")
	var update tgbotapi.Update
	if err := c.BodyParser(&update); err != nil {
		return err
	}

	logger.Info().Interface("update", update).Msg("Received update")

	if update.Message != nil {
		if update.Message.IsCommand() {
			command := update.Message.Command()
			args := update.Message.CommandArguments()

			if command == "start" {
				return h.handleStartCommand(c, update.Message, args)
			}
		}
	}
	return nil
}
func (h *Handler) handleStartCommand(c *fiber.Ctx, message *tgbotapi.Message, args string) error {
	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, используйте ссылку для подключения из приложения")
		_, err := h.bot.Send(msg)
		return err
	}

	userID, err := h.validateUserToken(args)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка валидации токена. Пожалуйста, попробуйте снова.")
		_, err := h.bot.Send(msg)
		return err
	}

	err = h.notificationService.ConnectTelegram(
		c.Context(),
		userID,
		fmt.Sprintf("%d", message.Chat.ID),
		message.From.UserName,
	)
	if err != nil {
		return err
	}

	// Установка базовых настроек уведомлений
	settings := []models.NotificationSettings{
		{
			UserID:           userID,
			NotificationType: "new_message",
			TelegramEnabled:  true,
			EmailEnabled:     true, // Добавляем email
		},
		{
			UserID:           userID,
			NotificationType: "new_review",
			TelegramEnabled:  true,
			EmailEnabled:     true, // Добавляем email
		},
		{
			UserID:           userID,
			NotificationType: "review_vote",
			TelegramEnabled:  true,
			EmailEnabled:     true, // Добавляем email
		},
		{
			UserID:           userID,
			NotificationType: "review_response",
			TelegramEnabled:  true,
			EmailEnabled:     true, // Добавляем email
		},
		{
			UserID:           userID,
			NotificationType: "listing_status",
			TelegramEnabled:  true,
			EmailEnabled:     true, // Добавляем email
		},
		{
			UserID:           userID,
			NotificationType: "favorite_price",
			TelegramEnabled:  true,
			EmailEnabled:     true, // Добавляем email
		},
	}

	for _, setting := range settings {
		if err := h.notificationService.UpdateNotificationSettings(c.Context(), &setting); err != nil {
			logger.Error().Err(err).Msg("Failed to update notification settings")
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Уведомления успешно подключены!")
	_, err = h.bot.Send(msg)
	return errors.Wrap(err, "telegram bot send() error")
}

// GetTelegramToken генерирует токен для подключения Telegram
// @Summary Получить токен для Telegram
// @Description Генерирует токен для связывания аккаунта с Telegram ботом
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Токен для подключения"
// @Failure 500 {object} utils.ErrorResponseSwag "notifications.tokenGenerateError"
// @Security BearerAuth
// @Router /api/v1/notifications/telegram/token [get]
func (h *Handler) GetTelegramToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	token, err := h.generateUserToken(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "notifications.tokenGenerateError")
	}

	// Изменяем структуру ответа для соответствия ожиданиям фронтенда
	return utils.SuccessResponse(c, fiber.Map{
		"token":        token, // упрощаем структуру
		"generated_at": time.Now(),
	})
}

// Генерация токена для привязки Telegram
func (h *Handler) generateUserToken(userID int) (string, error) {
	data := fmt.Sprintf("%d", userID)
	secret := []byte(os.Getenv("TELEGRAM_BOT_TOKEN"))
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%d_%s", userID, signature), nil
}

// Проверка токена
func (h *Handler) validateUserToken(token string) (int, error) {
	parts := strings.Split(token, "_")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid token format")
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid user ID")
	}

	expectedToken, err := h.generateUserToken(userID)
	if err != nil {
		return 0, err
	}

	if token != expectedToken {
		return 0, fmt.Errorf("invalid signature")
	}

	return userID, nil
}

// GetSettings возвращает настройки уведомлений
// @Summary Получить настройки уведомлений
// @Description Возвращает настройки уведомлений пользователя для всех типов
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string][]models.NotificationSettings} "Настройки уведомлений"
// @Security BearerAuth
// @Router /api/v1/notifications/settings [get]
func (h *Handler) GetSettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	// Создаем базовые настройки, если их нет
	baseSettings := []models.NotificationSettings{
		{
			UserID:           userID,
			NotificationType: "new_message",
			TelegramEnabled:  true,
			EmailEnabled:     true,
		},
		{
			UserID:           userID,
			NotificationType: "new_review",
			TelegramEnabled:  true,
			EmailEnabled:     true,
		},
		{
			UserID:           userID,
			NotificationType: "review_vote",
			TelegramEnabled:  true,
			EmailEnabled:     true,
		},
		{
			UserID:           userID,
			NotificationType: "review_response",
			TelegramEnabled:  true,
			EmailEnabled:     true,
		},
		{
			UserID:           userID,
			NotificationType: "listing_status",
			TelegramEnabled:  true,
			EmailEnabled:     true,
		},
		{
			UserID:           userID,
			NotificationType: "favorite_price",
			TelegramEnabled:  true,
			EmailEnabled:     true,
		},
	}

	settings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
	if err != nil || len(settings) == 0 {
		// Если настроек нет, создаем базовые
		for _, setting := range baseSettings {
			if err := h.notificationService.UpdateNotificationSettings(c.Context(), &setting); err != nil {
				logger.Error().Err(err).Msg("Error creating base settings")
			}
		}
		settings = baseSettings
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": settings,
	})
}

// SendPublicEmail отправляет email с публичной формы обратной связи
// @Summary Отправить email с формы обратной связи
// @Description Отправляет email сообщение с публичной формы контактов
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body map[string]string true "Данные формы" example({"name":"Имя","email":"email@example.com","message":"Сообщение","source":"website"})
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Email отправлен"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidDataFormat или validation.allFieldsRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "email.sendError"
// @Router /api/v1/notifications/email/public [post]
func (h *Handler) SendPublicEmail(c *fiber.Ctx) error {
	// Настройка CORS
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Content-Type")

	// Получаем данные из запроса
	var data struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
		Source  string `json:"source"` // Добавим поле для определения источника запроса
	}

	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidDataFormat")
	}

	logger.Info().Str("name", data.Name).Str("email", data.Email).Str("source", data.Source).Msg("Получен запрос на отправку email")

	// Валидация данных
	if data.Name == "" || data.Email == "" || data.Message == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.allFieldsRequired")
	}

	// Определяем получателя в зависимости от источника
	to := "info@svetu.rs" // По умолчанию
	if data.Source == "klimagrad" {
		to = "klimagrad@svetu.rs"
	}

	// Формируем заголовок
	subject := "Сообщение с сайта"
	if data.Source == "klimagrad" {
		subject = "Сообщение с сайта KlimaGrad"
	}

	// Формируем текст сообщения
	message := fmt.Sprintf("Имя: %s\nEmail: %s\n\nСообщение:\n%s",
		data.Name, data.Email, data.Message)

	// Используем ручное соединение без TLS вместо smtp.SendMail
	conn, err := smtp.Dial("mailserver:25")
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка соединения с SMTP")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "email.connectError")
	}
	defer conn.Close()

	// Устанавливаем отправителя и получателя
	if err = conn.Mail("info@svetu.rs"); err != nil {
		logger.Error().Err(err).Msg("Ошибка установки отправителя")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "email.senderError")
	}

	if err = conn.Rcpt(to); err != nil {
		logger.Error().Err(err).Msg("Ошибка установки получателя")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "email.recipientError")
	}

	// Отправляем данные
	wc, err := conn.Data()
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка при получении writer'а")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "email.writerError")
	}

	// Заголовки письма
	headers := "From: info@svetu.rs\r\n"
	headers += "Reply-To: " + data.Email + "\r\n"
	headers += "To: " + to + "\r\n"
	headers += "Subject: " + subject + "\r\n"
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/plain; charset=UTF-8\r\n\r\n"

	_, err = fmt.Fprintf(wc, headers+message)
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка записи данных")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "email.writeError")
	}

	err = wc.Close()
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка завершения отправки")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "email.closeError")
	}

	// Закрываем соединение
	err = conn.Quit()
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка закрытия соединения")
		// Не возвращаем ошибку, т.к. письмо уже должно быть отправлено
	}

	logger.Info().Str("recipient", to).Msg("Email успешно отправлен")
	return utils.SuccessResponse(c, fiber.Map{
		"success": true,
		"message": "email.sentSuccessfully",
	})
}

// UpdateSettings обновляет настройки уведомлений
// @Summary Обновить настройки уведомлений
// @Description Обновляет настройки уведомлений для конкретного типа
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Настройки" example({"notification_type":"new_message","telegram_enabled":true,"email_enabled":false})
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Настройки обновлены"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidDataFormat"
// @Failure 500 {object} utils.ErrorResponseSwag "notifications.updateError"
// @Security BearerAuth
// @Router /api/v1/notifications/settings [put]
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	body := string(c.Body())
	logger.Debug().Str("body", body).Msg("Raw request body")

	// Парсим запрос, не используя прямую привязку к структуре настроек
	var request struct {
		NotificationType string `json:"notification_type"`
		TelegramEnabled  *bool  `json:"telegram_enabled"`
		EmailEnabled     *bool  `json:"email_enabled"`
	}

	if err := c.BodyParser(&request); err != nil {
		logger.Error().Err(err).Msg("Error parsing request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidDataFormat")
	}

	logger.Debug().Interface("request", request).Msg("Parsed request")

	// Получаем текущие настройки пользователя
	currentSettings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting current settings")
		// Если ошибка, создаем новую настройку с дефолтными значениями
		currentSettings = []models.NotificationSettings{}
	}

	// Ищем настройку для этого типа уведомления
	var settings models.NotificationSettings
	settings.UserID = userID
	settings.NotificationType = request.NotificationType

	// По умолчанию настройка активна для обоих каналов, если создается новая
	settings.TelegramEnabled = true
	settings.EmailEnabled = true

	// Проверяем, есть ли уже такая настройка
	for _, s := range currentSettings {
		if s.NotificationType == request.NotificationType {
			// Если нашли, используем текущие настройки как базовые
			settings = s
			break
		}
	}

	// Обновляем только те настройки, которые пришли в запросе
	if request.TelegramEnabled != nil {
		settings.TelegramEnabled = *request.TelegramEnabled
		logger.Debug().Bool("telegram_enabled", settings.TelegramEnabled).Msg("Setting telegram_enabled")
	}

	if request.EmailEnabled != nil {
		settings.EmailEnabled = *request.EmailEnabled
		logger.Debug().Bool("email_enabled", settings.EmailEnabled).Msg("Setting email_enabled")
	}

	logger.Debug().Interface("settings", settings).Msg("Final settings to save")

	err = h.notificationService.UpdateNotificationSettings(c.Context(), &settings)
	if err != nil {
		logger.Error().Err(err).Msg("Error updating settings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "notifications.updateError")
	}

	// Получаем обновленные настройки для ответа
	updatedSettings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting updated settings")
		updatedSettings = []models.NotificationSettings{settings}
	}

	// Возвращаем обновленные настройки в ответе
	return utils.SuccessResponse(c, fiber.Map{
		"message":  "notifications.settingsUpdated",
		"settings": updatedSettings,
	})
}

// GetTelegramStatus проверяет статус подключения Telegram
// @Summary Статус подключения Telegram
// @Description Проверяет, подключен ли Telegram к аккаунту пользователя
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Статус подключения"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.unauthorized"
// @Router /api/v1/notifications/telegram/status [get]
func (h *Handler) GetTelegramStatus(c *fiber.Ctx) error {
	// Добавить безопасное получение userID и проверку на авторизацию
	var userID int
	if uidVal := c.Locals("user_id"); uidVal != nil {
		if uid, ok := uidVal.(int); ok {
			userID = uid
		} else {
			// Если user_id есть, но неверного типа, возвращаем ошибку авторизации
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
		}
	} else {
		// Для неавторизованных запросов просто возвращаем отсутствие подключения
		return utils.SuccessResponse(c, fiber.Map{
			"connected": false,
		})
	}

	connection, err := h.notificationService.GetTelegramConnection(c.Context(), userID)

	if err != nil {
		return utils.SuccessResponse(c, fiber.Map{
			"connected": false,
		})
	}

	return utils.SuccessResponse(c, fiber.Map{
		"connected": true,
		"username":  connection.TelegramUsername,
	})
}

// ConnectTelegram связывает аккаунт с Telegram
// @Summary Подключить Telegram
// @Description Связывает аккаунт пользователя с Telegram для получения уведомлений
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body map[string]string true "Данные для подключения" example({"chat_id":"123456789","username":"username"})
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]string} "Telegram подключен"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidDataFormat"
// @Failure 500 {object} utils.ErrorResponseSwag "telegram.connectionError"
// @Security BearerAuth
// @Router /api/v1/notifications/telegram/connect [post]
func (h *Handler) ConnectTelegram(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var data struct {
		ChatID   string `json:"chat_id"`
		Username string `json:"username"`
	}

	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidDataFormat")
	}

	err := h.notificationService.ConnectTelegram(c.Context(), userID, data.ChatID, data.Username)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "telegram.connectionError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "telegram.connected",
	})
}

// MarkAsRead отмечает уведомление как прочитанное
// @Summary Отметить как прочитанное
// @Description Отмечает уведомление как прочитанное
// @Tags notifications
// @Accept json
// @Produce json
// @Param id path int true "ID уведомления"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]string} "Уведомление отмечено как прочитанное"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidNotificationId"
// @Failure 500 {object} utils.ErrorResponseSwag "notifications.updateError"
// @Security BearerAuth
// @Router /api/v1/notifications/{id}/read [put]
func (h *Handler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	notificationID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidNotificationId")
	}

	if err = h.notificationService.MarkNotificationAsRead(c.Context(), userID, notificationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "notifications.updateError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "notifications.markedAsRead",
	})
}