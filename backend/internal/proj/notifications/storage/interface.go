// backend/internal/proj/notifications/storage/interface.go
package storage

import (
    "context"
    "backend/internal/domain/models"
)

type NotificationRepository interface {
    // Методы для настроек уведомлений
    GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error)
    UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error
    
    // Методы для Telegram
    SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error
    GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error)
    DeleteTelegramConnection(ctx context.Context, userID int) error
    
    // Методы для уведомлений
    CreateNotification(ctx context.Context, notification *models.Notification) error
    GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error)
    MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error
    DeleteNotification(ctx context.Context, userID int, notificationID int) error
}