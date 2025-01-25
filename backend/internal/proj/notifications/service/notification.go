package service

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/storage"
)

type NotificationService struct {
    storage storage.Storage
}

func NewNotificationService(storage storage.Storage) NotificationServiceInterface {
    return &NotificationService{
        storage: storage,
    }
}

func (s *NotificationService) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
    return s.storage.GetNotificationSettings(ctx, userID)
}

func (s *NotificationService) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error {
    return s.storage.UpdateNotificationSettings(ctx, settings)
}

func (s *NotificationService) ConnectTelegram(ctx context.Context, userID int, chatID string, username string) error {
    return s.storage.SaveTelegramConnection(ctx, userID, chatID, username)
}

func (s *NotificationService) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
    return s.storage.GetTelegramConnection(ctx, userID)
}

func (s *NotificationService) DisconnectTelegram(ctx context.Context, userID int) error {
    return s.storage.DeleteTelegramConnection(ctx, userID)
}

func (s *NotificationService) SavePushSubscription(ctx context.Context, userID int, subscription *models.PushSubscription) error {
    subscription.UserID = userID
    return s.storage.SavePushSubscription(ctx, subscription)
}

func (s *NotificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
    return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
    return s.storage.GetUserNotifications(ctx, userID, limit, offset)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
    return s.storage.MarkNotificationAsRead(ctx, userID, notificationID)
}

func (s *NotificationService) SendChatNotification(ctx context.Context, userID int, message string) error {
    notification := &models.Notification{
        UserID:  userID,
        Type:    models.NotificationTypeNewMessage,
        Title:   "Новое сообщение",
        Message: message,
    }
    return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationService) SendReviewNotification(ctx context.Context, userID int, entityType string, entityID int) error {
    notification := &models.Notification{
        UserID:  userID,
        Type:    models.NotificationTypeNewReview,
        Title:   "Новый отзыв",
        Message: "Получен новый отзыв для " + entityType,
    }
    return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationService) SendListingUpdateNotification(ctx context.Context, userID int, listingID int, updateType string) error {
    notification := &models.Notification{
        UserID:  userID,
        Type:    models.NotificationTypeListingStatus,
        Title:   "Обновление объявления",
        Message: "Статус объявления изменен на: " + updateType,
    }
    return s.storage.CreateNotification(ctx, notification)
}