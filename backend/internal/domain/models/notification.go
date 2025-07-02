// backend/internal/domain/models/notification.go
package models

import (
	"encoding/json"
	"time"
)

type NotificationSettings struct {
	UserID           int       `json:"user_id"`
	NotificationType string    `json:"notification_type"`
	TelegramEnabled  bool      `json:"telegram_enabled"`
	EmailEnabled     bool      `json:"email_enabled"` // Добавляем поле для email
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Notification struct {
	ID          int             `json:"id"`
	UserID      int             `json:"user_id"`
	Type        string          `json:"type"`
	Title       string          `json:"title"`
	Message     string          `json:"message"`
	ListingID   int             `json:"listing_id,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
	IsRead      bool            `json:"is_read"`
	DeliveredTo json.RawMessage `json:"delivered_to"`
	CreatedAt   time.Time       `json:"created_at"`
}

type TelegramConnection struct {
	UserID           int       `json:"user_id"`
	TelegramChatID   string    `json:"telegram_chat_id"`
	TelegramUsername string    `json:"telegram_username"`
	ConnectedAt      time.Time `json:"connected_at"`
}

// Константы типов уведомлений
