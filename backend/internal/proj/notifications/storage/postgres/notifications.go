// backend/internal/proj/notifications/storage/postgres/notifications.go
package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewNotificationStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT user_id, notification_type, telegram_enabled, email_enabled, created_at, updated_at
        FROM notification_settings
        WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying notification settings: %w", err)
	}
	defer rows.Close()

	var settings []models.NotificationSettings
	for rows.Next() {
		var setting models.NotificationSettings
		err := rows.Scan(
			&setting.UserID,
			&setting.NotificationType,
			&setting.TelegramEnabled,
			&setting.EmailEnabled,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning notification setting: %w", err)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func (s *Storage) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO notification_settings (
            user_id,
            notification_type,
            telegram_enabled,
            email_enabled
        ) VALUES (
            $1,
            $2,
            $3,
            $4
        )
        ON CONFLICT (user_id, notification_type)
        DO UPDATE SET
            telegram_enabled = EXCLUDED.telegram_enabled,
            email_enabled = EXCLUDED.email_enabled,
            updated_at = CURRENT_TIMESTAMP
    `,
		settings.UserID,
		settings.NotificationType,
		settings.TelegramEnabled,
		settings.EmailEnabled,
	)
	if err != nil {
		return fmt.Errorf("error updating notification settings: %w", err)
	}

	return nil
}

func (s *Storage) SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error {
	result, err := s.pool.Exec(ctx, `
        INSERT INTO user_telegram_connections (user_id, telegram_chat_id, telegram_username)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id)
        DO UPDATE SET
            telegram_chat_id = $2,
            telegram_username = $3,
            connected_at = CURRENT_TIMESTAMP
        RETURNING user_id
    `, userID, chatID, username)
	if err != nil {
		return fmt.Errorf("error saving telegram connection: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected when saving telegram connection")
	}

	return nil
}

func (s *Storage) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
	connection := &models.TelegramConnection{}
	err := s.pool.QueryRow(ctx, `
        SELECT user_id, telegram_chat_id, telegram_username, connected_at
        FROM user_telegram_connections
        WHERE user_id = $1
    `, userID).Scan(
		&connection.UserID,
		&connection.TelegramChatID,
		&connection.TelegramUsername,
		&connection.ConnectedAt,
	)
	if err != nil {
		return nil, err
	}
	return connection, nil
}

func (s *Storage) DeleteTelegramConnection(ctx context.Context, userID int) error {
	_, err := s.pool.Exec(ctx, `
        DELETE FROM user_telegram_connections
        WHERE user_id = $1
    `, userID)
	return err
}

func (s *Storage) CreateNotification(ctx context.Context, notification *models.Notification) error {
	_, err := s.pool.Exec(ctx, `
        INSERT INTO notifications (
            user_id, type, title, message, data, is_read, delivered_to
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Data,
		notification.IsRead,
		notification.DeliveredTo,
	)
	return err
}

func (s *Storage) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
	rows, err := s.pool.Query(ctx, `
        SELECT id, user_id, type, title, message, data, is_read, delivered_to, created_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.Data,
			&n.IsRead,
			&n.DeliveredTo,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (s *Storage) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
	result, err := s.pool.Exec(ctx, `
        UPDATE notifications
        SET is_read = true
        WHERE id = $1 AND user_id = $2
    `, notificationID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("notification not found or access denied")
	}
	return nil
}

func (s *Storage) DeleteNotification(ctx context.Context, userID int, notificationID int) error {
	_, err := s.pool.Exec(ctx, `
        DELETE FROM notifications
        WHERE id = $1 AND user_id = $2
    `, notificationID, userID)
	return err
}
