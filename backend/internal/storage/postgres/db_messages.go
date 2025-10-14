// backend/internal/storage/postgres/db_messages.go
package postgres

import (
	"context"

	"backend/internal/domain/models"
)

func (db *Database) CreateNotification(ctx context.Context, n *models.Notification) error {
	return db.notificationsDB.CreateNotification(ctx, n)
}

func (db *Database) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
	return db.notificationsDB.GetNotificationSettings(ctx, userID)
}

func (db *Database) UpdateNotificationSettings(ctx context.Context, s *models.NotificationSettings) error {
	return db.notificationsDB.UpdateNotificationSettings(ctx, s)
}

func (db *Database) SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error {
	return db.notificationsDB.SaveTelegramConnection(ctx, userID, chatID, username)
}

func (db *Database) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
	return db.notificationsDB.GetTelegramConnection(ctx, userID)
}

func (db *Database) DeleteTelegramConnection(ctx context.Context, userID int) error {
	return db.notificationsDB.DeleteTelegramConnection(ctx, userID)
}

func (db *Database) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
	return db.notificationsDB.GetUserNotifications(ctx, userID, limit, offset)
}

func (db *Database) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
	return db.notificationsDB.MarkNotificationAsRead(ctx, userID, notificationID)
}

func (db *Database) DeleteNotification(ctx context.Context, userID int, notificationID int) error {
	return db.notificationsDB.DeleteNotification(ctx, userID, notificationID)
}
