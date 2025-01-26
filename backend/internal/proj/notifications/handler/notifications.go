// backend/internal/proj/notifications/handler/notifications.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/notifications/service"
	"backend/pkg/utils"
	"log"
	"os"
	"fmt"
//	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
//	"errors"
	"strconv"
	"strings"

   "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService service.NotificationServiceInterface
	bot                *tgbotapi.BotAPI
}

func NewNotificationHandler(service service.NotificationServiceInterface) *NotificationHandler {
    handler := &NotificationHandler{
        notificationService: service,
    }
    
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
    if err == nil {
        handler.bot = bot
    }
    
    return handler
}

 func (h *NotificationHandler) SubscribePush(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    var subscription models.PushSubscription
    if err := c.BodyParser(&subscription); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid subscription data")
    }
    
    err := h.notificationService.SavePushSubscription(c.Context(), userID, &subscription)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save subscription")
    }
    
    return utils.SuccessResponse(c, fiber.Map{
        "message": "Successfully subscribed to push notifications",
    })
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
	if h.bot == nil {
        return
    }
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


// In NotificationHandler
func (h *NotificationHandler) HandleTelegramWebhook(c *fiber.Ctx) error {
    var update tgbotapi.Update
    if err := c.BodyParser(&update); err != nil {
        log.Printf("Error parsing update: %v", err)
        return err
    }

    log.Printf("Received update: %+v", update)

    if update.Message != nil {
        if update.Message.IsCommand() {
            command := update.Message.Command()
            args := update.Message.CommandArguments()
            log.Printf("Command %s with args: '%s'", command, args)
            
            if command == "start" {
                return h.handleStartCommand(c, update.Message, args)
            }
        }
    }
    return nil
}
func (h *NotificationHandler) handleStartCommand(c *fiber.Ctx, message *tgbotapi.Message, token string) error {
    userID, err := h.validateUserToken(token)
    if err != nil {
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

    // Установим базовые настройки уведомлений
    settings := []models.NotificationSettings{
        {UserID: userID, NotificationType: "new_message", TelegramEnabled: true},
        {UserID: userID, NotificationType: "new_review", TelegramEnabled: true},
        {UserID: userID, NotificationType: "review_vote", TelegramEnabled: true},
        {UserID: userID, NotificationType: "review_response", TelegramEnabled: true},
        {UserID: userID, NotificationType: "listing_status", TelegramEnabled: true},
        {UserID: userID, NotificationType: "favorite_price", TelegramEnabled: true},
    }

    for _, setting := range settings {
        if err := h.notificationService.UpdateNotificationSettings(c.Context(), &setting); err != nil {
            log.Printf("Error setting notification: %v", err)
        }
    }

    msg := tgbotapi.NewMessage(message.Chat.ID, "Уведомления успешно подключены!")
    h.bot.Send(msg)
    return nil
}


func (h *NotificationHandler) GetTelegramToken(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    token, err := h.generateUserToken(userID)
    if err != nil {
        log.Printf("Error generating token: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
    }
    
     log.Printf("Generated telegram token for user %d: %s", userID, token)
    
	return utils.SuccessResponse(c, fiber.Map{
		"data": fiber.Map{
			"token": token,
		},
	})
}

// Генерация токена для привязки Telegram
func (h *NotificationHandler) generateUserToken(userID int) (string, error) {
    data := fmt.Sprintf("%d", userID)
    secret := []byte(os.Getenv("TELEGRAM_BOT_TOKEN"))
    hash := hmac.New(sha256.New, secret)
    hash.Write([]byte(data))
    signature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
    return fmt.Sprintf("%d_%s", userID, signature), nil
}



// Проверка токена
func (h *NotificationHandler) validateUserToken(token string) (int, error) {
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
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid data format")
    }

    settings.UserID = userID
    err := h.notificationService.UpdateNotificationSettings(c.Context(), &settings)
    if err != nil {
        log.Printf("Error updating settings: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update settings")
    }

    return utils.SuccessResponse(c, fiber.Map{"message": "Settings updated"})
}

// GetTelegramStatus проверяет статус подключения Telegram
func (h *NotificationHandler) GetTelegramStatus(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    connection, err := h.notificationService.GetTelegramConnection(c.Context(), userID)
    log.Printf("Checking Telegram status for user %d: connection=%v, err=%v", userID, connection, err)
    
    if err != nil {
        return utils.SuccessResponse(c, fiber.Map{
            "connected": false,
        })
    }

    return utils.SuccessResponse(c, fiber.Map{
        "connected": true,
        "username": connection.TelegramUsername,
    })
}

func (h *NotificationHandler) SendTestNotification(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    connection, err := h.notificationService.GetTelegramConnection(c.Context(), userID)
    if err != nil {
        log.Printf("Error getting telegram connection: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Telegram not connected")
    }

    msg := tgbotapi.NewMessage(connection.TelegramChatID, "Тестовое уведомление")
    _, err = h.bot.Send(msg)
    if err != nil {
        log.Printf("Error sending test notification: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error sending notification")
    }

    return utils.SuccessResponse(c, fiber.Map{"message": "Test notification sent"})
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
