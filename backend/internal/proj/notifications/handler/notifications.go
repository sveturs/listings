// backend/internal/proj/notifications/handler/notifications.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/notifications/service"
	"backend/pkg/utils"
	"log"
	"os"
	"fmt"
	"time"
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
    if h.bot == nil {
        log.Printf("Bot not initialized")
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Bot not initialized")
    }

    var update tgbotapi.Update
    if err := c.BodyParser(&update); err != nil {
        log.Printf("Error parsing update: %v", err)
        return err
    }

    if update.Message != nil && update.Message.IsCommand() {
        log.Printf("Received command: %s with args: %s", 
            update.Message.Command(), 
            update.Message.CommandArguments())
            
        if update.Message.Command() == "start" {
            token := update.Message.CommandArguments()
            userID, err := h.validateUserToken(token)
            if err != nil {
                log.Printf("Invalid token: %v", err)
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Недействительный токен")
                h.bot.Send(msg)
                return nil
            }
            
            err = h.notificationService.ConnectTelegram(
                c.Context(),
                userID,
                fmt.Sprintf("%d", update.Message.Chat.ID),
                update.Message.From.UserName,
            )
            if err != nil {
                log.Printf("Error connecting telegram: %v", err) 
                msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка подключения")
                h.bot.Send(msg)
                return err
            }
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот успешно подключен!")
            h.bot.Send(msg)
        }
    }
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
        "token": token,  // Убедитесь, что токен передается именно в этом поле
    })
}

// Генерация токена для привязки Telegram
func (h *NotificationHandler) generateUserToken(userID int) (string, error) {
    data := fmt.Sprintf("%d:%d", userID, time.Now().Unix())
    hash := hmac.New(sha256.New, []byte(os.Getenv("TELEGRAM_BOT_SECRET")))
    hash.Write([]byte(data))
    token := base64.URLEncoding.EncodeToString(hash.Sum(nil))
    return fmt.Sprintf("%d.%s", userID, token), nil
}

// Проверка токена
func (h *NotificationHandler) validateUserToken(token string) (int, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return 0, fmt.Errorf("invalid token format: expected 3 parts")
    }

    userID, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, fmt.Errorf("invalid user ID in token")
    }

    timestamp, err := strconv.ParseInt(parts[1], 10, 64)
    if err != nil {
        return 0, fmt.Errorf("invalid timestamp in token")
    }

    // Проверяем срок действия токена (например, 5 минут)
    if time.Now().Unix()-timestamp > 300 {
        return 0, fmt.Errorf("token expired")
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
