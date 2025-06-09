// Package handler
// backend/internal/proj/notifications/handler/responses.go
package handler

import (
	"time"

	"backend/internal/domain/models"
)

// TelegramTokenResponse структура ответа с токеном для Telegram
// @Description Токен для связывания аккаунта с Telegram ботом
type TelegramTokenResponse struct {
	// Сгенерированный токен
	Token string `json:"token" example:"123_abc..."`
	// Время генерации токена
	GeneratedAt time.Time `json:"generated_at" example:"2024-01-15T10:30:00Z"`
}

// TelegramStatusResponse структура ответа о статусе подключения Telegram
// @Description Информация о подключении Telegram к аккаунту
type TelegramStatusResponse struct {
	// Подключен ли Telegram
	Connected bool `json:"connected" example:"true"`
	// Имя пользователя в Telegram (если подключен)
	Username string `json:"username,omitempty" example:"john_doe"`
}

// TelegramConnectResponse структура ответа при подключении Telegram
// @Description Ответ при успешном подключении Telegram
type TelegramConnectResponse struct {
	// Сообщение об успешном подключении
	Message string `json:"message" example:"telegram.connected"`
}

// NotificationSettingsResponse структура ответа с настройками уведомлений
// @Description Настройки уведомлений пользователя
type NotificationSettingsResponse struct {
	// Массив настроек для разных типов уведомлений
	Data []models.NotificationSettings `json:"data"`
}

// NotificationSettingsUpdateResponse структура ответа при обновлении настроек
// @Description Ответ при обновлении настроек уведомлений
type NotificationSettingsUpdateResponse struct {
	// Сообщение об успешном обновлении
	Message string `json:"message" example:"notifications.settingsUpdated"`
	// Обновленные настройки
	Settings []models.NotificationSettings `json:"settings"`
}

// PublicEmailSendResponse структура ответа при отправке публичного email
// @Description Ответ при отправке email с публичной формы
type PublicEmailSendResponse struct {
	// Флаг успешной отправки
	Success bool `json:"success" example:"true"`
	// Сообщение о результате
	Message string `json:"message" example:"email.sentSuccessfully"`
}

// NotificationMarkReadResponse структура ответа при отметке уведомления как прочитанного
// @Description Ответ при изменении статуса уведомления
type NotificationMarkReadResponse struct {
	// Сообщение об успешной отметке
	Message string `json:"message" example:"notifications.markedAsRead"`
}

// TelegramWebhookResponse структура ответа для webhook Telegram
// @Description Ответ на webhook запрос от Telegram
type TelegramWebhookResponse struct {
	// Статус обработки
	Status string `json:"status" example:"OK"`
}
