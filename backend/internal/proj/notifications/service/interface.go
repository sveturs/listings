// backend/internal/proj/notifications/service/interface.go
package service

import (
    "context"
    "backend/internal/domain/models"
)

type NotificationServiceInterface interface {
    // Методы настроек
    GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error)
    UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error
    
    // Методы для каналов уведомлений
    ConnectTelegram(ctx context.Context, userID int, chatID string, username string) error
	GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error)
    DisconnectTelegram(ctx context.Context, userID int) error
    SavePushSubscription(ctx context.Context, userID int, subscription *models.PushSubscription) error
    
    // Методы для уведомлений
    CreateNotification(ctx context.Context, notification *models.Notification) error
    GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error)
    MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error
    
    // Отправка уведомлений
    SendChatNotification(ctx context.Context, userID int, message string) error
    SendReviewNotification(ctx context.Context, userID int, entityType string, entityID int) error
    SendListingUpdateNotification(ctx context.Context, userID int, listingID int, updateType string) error
}