// backend/internal/proj/notifications/handler/notifications.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/notifications/service"
	"backend/pkg/utils"
	"log"
	"os"
   "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService service.NotificationServiceInterface
	bot                *tgbotapi.BotAPI
}

func NewNotificationHandler(service service.NotificationServiceInterface) *NotificationHandler {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Printf("Error initializing Telegram bot: %v", err)
	}
	return &NotificationHandler{
		notificationService: service,
		bot:                bot,
	}
 }

// GetNotifications handler
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	notifications, err := h.notificationService.GetUserNotifications(c.Context(), userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении уведомлений")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": notifications, // Изменено с notifications на data
	})
}
func (h *NotificationHandler) ConnectTelegramWebhook() {
    baseURL := "https://landhub.rs/api/v1/notifications/telegram/webhook"
    webhook, err := tgbotapi.NewWebhookWithCert(baseURL, nil)
    if err != nil {
        log.Printf("Error creating webhook: %v", err)
        return
    }

    if _, err := h.bot.Request(webhook); err != nil {
        log.Printf("Error setting webhook: %v", err)
    }
}
func (h *NotificationHandler) HandleTelegramWebhook(c *fiber.Ctx) error {
    update := tgbotapi.Update{}
    if err := c.BodyParser(&update); err != nil {
        return err
    }

    // Обработка входящих сообщений
    if update.Message != nil {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот успешно подключен!")
        _, err := h.bot.Send(msg)
        if err != nil {
            log.Printf("Error sending message: %v", err)
        }
    }

    return c.SendStatus(fiber.StatusOK)
}
// GetSettings handler
func (h *NotificationHandler) GetSettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	settings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении настроек")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": settings, // Добавлено data
	})
}

// UpdateSettings обновляет настройки уведомлений
func (h *NotificationHandler) UpdateSettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var settings models.NotificationSettings
	if err := c.BodyParser(&settings); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	settings.UserID = userID

	err := h.notificationService.UpdateNotificationSettings(c.Context(), &settings)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при обновлении настроек")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Настройки успешно обновлены",
	})
}

// GetTelegramStatus проверяет статус подключения Telegram
func (h *NotificationHandler) GetTelegramStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

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

// ConnectTelegram связывает аккаунт Telegram
func (h *NotificationHandler) ConnectTelegram(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var data struct {
		ChatID   string `json:"chat_id"`
		Username string `json:"username"`
	}

	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	err := h.notificationService.ConnectTelegram(c.Context(), userID, data.ChatID, data.Username)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при подключении Telegram")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Telegram успешно подключен",
	})
}

// MarkAsRead отмечает уведомление как прочитанное
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	notificationID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid notification ID")
	}

	if err = h.notificationService.MarkNotificationAsRead(c.Context(), userID, notificationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при обновлении статуса уведомления")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Уведомление отмечено как прочитанное",
	})
}
