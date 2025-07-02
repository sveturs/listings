// backend/internal/proj/notifications/service/interface.go
package service

import (
	"context"

	"backend/internal/domain/models"
)

type NotificationServiceInterface interface {
	GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error
	ConnectTelegram(ctx context.Context, userID int, chatID string, username string) error
	GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error)
	DisconnectTelegram(ctx context.Context, userID int) error
	CreateNotification(ctx context.Context, notification *models.Notification) error
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error
	SendNotification(ctx context.Context, userID int, notificationType string, message string, listingID int) error
}
