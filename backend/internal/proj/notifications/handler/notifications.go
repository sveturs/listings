// backend/internal/proj/notifications/handler/notifications.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/notifications/service"
	"backend/pkg/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	    "net/smtp"
	"fmt"
	"log"
	"os"
	"time"

	//	"errors"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService service.NotificationServiceInterface
	bot                 *tgbotapi.BotAPI
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

	baseURL := "https://svetu.rs/api/v1/notifications/telegram/webhook"

	// Создаем конфигурацию вебхука
	webhookConfig := tgbotapi.NewWebhook(baseURL)

	// Устанавливаем вебхук
	_, err := h.bot.SetWebhook(webhookConfig)
	if err != nil {
		//log.Printf("Error setting webhook: %v", err)
		return
	}

	// Опционально: проверяем информацию о вебхуке
	info, err := h.bot.GetWebhookInfo()
	if err != nil {
		//log.Printf("Error getting webhook info: %v", err)
		return
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	//log.Printf("Successfully set webhook for Telegram bot")
}

// In NotificationHandler
func (h *NotificationHandler) HandleTelegramWebhook(c *fiber.Ctx) error {
	var update tgbotapi.Update
	if err := c.BodyParser(&update); err != nil {
		//log.Printf("Error parsing update: %v", err)
		return err
	}

	log.Printf("Received update: %+v", update)

	if update.Message != nil {
		if update.Message.IsCommand() {
			command := update.Message.Command()
			args := update.Message.CommandArguments()
			//log.Printf("Command %s with args: '%s'", command, args)

			if command == "start" {
				return h.handleStartCommand(c, update.Message, args)
			}
		}
	}
	return nil
}
func (h *NotificationHandler) handleStartCommand(c *fiber.Ctx, message *tgbotapi.Message, args string) error {
	// Добавляем логирование для отладки
	//log.Printf("Handling start command with args: %s", args)

	if args == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Пожалуйста, используйте ссылку для подключения из приложения")
		_, err := h.bot.Send(msg)
		return err
	}

	userID, err := h.validateUserToken(args)
	if err != nil {
		//log.Printf("Token validation error: %v", err)
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
			// log.Printf("Error setting notification: %v", err)
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Уведомления успешно подключены!")
	_, err = h.bot.Send(msg)
	return err
}

func (h *NotificationHandler) GetTelegramToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	token, err := h.generateUserToken(userID)
	if err != nil {
		// log.Printf("Error generating token: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	// log.Printf("Generated telegram token for user %d: %s", userID, token)

	// Изменяем структуру ответа для соответствия ожиданиям фронтенда
	return utils.SuccessResponse(c, fiber.Map{
		"token":        token, // упрощаем структуру
		"generated_at": time.Now(),
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

    // Создаем базовые настройки, если их нет
    baseSettings := []models.NotificationSettings{
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


	settings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
	if err != nil || len(settings) == 0 {
		// Если настроек нет, создаем базовые
		for _, setting := range baseSettings {
			if err := h.notificationService.UpdateNotificationSettings(c.Context(), &setting); err != nil {
				log.Printf("Error creating base settings: %v", err)
			}
		}
		settings = baseSettings
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": settings,
	})
}

// SendPublicEmail обрабатывает запросы на отправку писем с форм обратной связи
func (h *NotificationHandler) SendPublicEmail(c *fiber.Ctx) error {
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
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid data format")
    }
    
    log.Printf("Получен запрос на отправку email от %s (%s): %s", data.Name, data.Email, data.Source)

    // Валидация данных
    if data.Name == "" || data.Email == "" || data.Message == "" {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "All fields are required")
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
        log.Printf("Ошибка соединения с SMTP: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to connect to mail server")
    }
    defer conn.Close()
    
    // Устанавливаем отправителя и получателя
    if err = conn.Mail("info@svetu.rs"); err != nil {
        log.Printf("Ошибка установки отправителя: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set sender")
    }
    
    if err = conn.Rcpt(to); err != nil {
        log.Printf("Ошибка установки получателя: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to set recipient")
    }
    
    // Отправляем данные
    wc, err := conn.Data()
    if err != nil {
        log.Printf("Ошибка при получении writer'а: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get data writer")
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
        log.Printf("Ошибка записи данных: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to write email data")
    }
    
    err = wc.Close()
    if err != nil {
        log.Printf("Ошибка завершения отправки: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to close connection")
    }
    
    // Закрываем соединение
    err = conn.Quit()
    if err != nil {
        log.Printf("Ошибка закрытия соединения: %v", err)
        // Не возвращаем ошибку, т.к. письмо уже должно быть отправлено
    }
    
    log.Printf("Email успешно отправлен на %s", to)
    return utils.SuccessResponse(c, fiber.Map{
        "success": true,
        "message": "Email sent successfully",
    })
}
func (h *NotificationHandler) UpdateSettings(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    // Получаем и логируем необработанное тело запроса
    body := string(c.Body())
    log.Printf("Raw request body: %s", body)
    
    // Парсим запрос, не используя прямую привязку к структуре настроек
    var request struct {
        NotificationType string `json:"notification_type"`
        TelegramEnabled  *bool  `json:"telegram_enabled"`
        EmailEnabled     *bool  `json:"email_enabled"`
    }
    
    if err := c.BodyParser(&request); err != nil {
        log.Printf("Error parsing request: %v", err)
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid data format")
    }
    
    log.Printf("Parsed request: %+v", request)
    
    // Получаем текущие настройки пользователя
    currentSettings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
    if err != nil {
        log.Printf("Error getting current settings: %v", err)
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
        log.Printf("Setting telegram_enabled to: %v", settings.TelegramEnabled)
    }
    
    if request.EmailEnabled != nil {
        settings.EmailEnabled = *request.EmailEnabled
        log.Printf("Setting email_enabled to: %v", settings.EmailEnabled)
    }
    
    log.Printf("Final settings to save: %+v", settings)
    
    err = h.notificationService.UpdateNotificationSettings(c.Context(), &settings)
    if err != nil {
        log.Printf("Error updating settings: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update settings")
    }
    
    // Получаем обновленные настройки для ответа
    updatedSettings, err := h.notificationService.GetNotificationSettings(c.Context(), userID)
    if err != nil {
        log.Printf("Error getting updated settings: %v", err)
        updatedSettings = []models.NotificationSettings{settings}
    }
    
    // Возвращаем обновленные настройки в ответе
    return utils.SuccessResponse(c, fiber.Map{
        "message": "Settings updated",
        "settings": updatedSettings,
    })
}

// GetTelegramStatus проверяет статус подключения Telegram
func (h *NotificationHandler) GetTelegramStatus(c *fiber.Ctx) error {
    // Добавить безопасное получение userID и проверку на авторизацию
    var userID int
    if uidVal := c.Locals("user_id"); uidVal != nil {
        if uid, ok := uidVal.(int); ok {
            userID = uid
        } else {
            // Если user_id есть, но неверного типа, возвращаем ошибку авторизации
            return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
        }
    } else {
        // Для неавторизованных запросов просто возвращаем отсутствие подключения
        return utils.SuccessResponse(c, fiber.Map{
            "connected": false,
        })
    }

    connection, err := h.notificationService.GetTelegramConnection(c.Context(), userID)
    // log.Printf("Checking Telegram status for user %d: connection=%v, err=%v", userID, connection, err)

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
