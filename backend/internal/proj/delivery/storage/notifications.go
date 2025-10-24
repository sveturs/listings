package storage

import (
	"context"
	"encoding/json"
	"time"
)

// DeliveryNotification - уведомление о доставке
type DeliveryNotification struct {
	ID           int             `db:"id"`
	ShipmentID   int             `db:"shipment_id"`
	UserID       int             `db:"user_id"`
	Channel      string          `db:"channel"` // email, sms, viber, telegram, push
	Status       string          `db:"status"`  // pending, sent, failed
	Template     string          `db:"template"`
	Data         json.RawMessage `db:"data"`
	SentAt       *time.Time      `db:"sent_at"`
	ErrorMessage *string         `db:"error_message"`
	CreatedAt    time.Time       `db:"created_at"`
}

// SaveNotification сохраняет уведомление в БД
func (s *Storage) SaveNotification(ctx context.Context, notification *DeliveryNotification) error {
	query := `
		INSERT INTO delivery_notifications (
			shipment_id, channel, status, template, data, sent_at, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := s.db.GetContext(ctx, notification, query,
		notification.ShipmentID,
		notification.Channel,
		notification.Status,
		notification.Template,
		notification.Data,
		notification.SentAt,
		notification.ErrorMessage,
	)

	return err
}

// GetNotificationHistory получает историю уведомлений для отправления
func (s *Storage) GetNotificationHistory(ctx context.Context, shipmentID int) ([]*DeliveryNotification, error) {
	var notifications []*DeliveryNotification

	query := `
		SELECT id, shipment_id, channel, status, template, data, sent_at, error_message, created_at
		FROM delivery_notifications
		WHERE shipment_id = $1
		ORDER BY created_at DESC`

	err := s.db.SelectContext(ctx, &notifications, query, shipmentID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetNotificationByID получает уведомление по ID
func (s *Storage) GetNotificationByID(ctx context.Context, notificationID int) (*DeliveryNotification, error) {
	var notification DeliveryNotification

	query := `SELECT * FROM delivery_notifications WHERE id = $1`

	err := s.db.GetContext(ctx, &notification, query, notificationID)
	if err != nil {
		return nil, err
	}

	return &notification, nil
}

// UpdateNotificationStatus обновляет статус уведомления
func (s *Storage) UpdateNotificationStatus(ctx context.Context, notificationID int, status string, errorMsg *string) error {
	query := `
		UPDATE delivery_notifications
		SET status = $1, error_message = $2
		WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, status, errorMsg, notificationID)
	return err
}
