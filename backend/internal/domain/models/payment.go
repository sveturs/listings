// backend/internal/domain/models/payment.go

package models

import "time"

type PaymentSession struct {
	ID            string     `json:"id"` // Изменено на string для поддержки external IDs
	UserID        int        `json:"user_id"`
	OrderID       *int       `json:"order_id,omitempty"` // Добавлено для поддержки заказов
	Amount        float64    `json:"amount"`
	Currency      string     `json:"currency"`
	PaymentMethod string     `json:"payment_method"`
	ExternalID    string     `json:"external_id,omitempty"`
	Status        string     `json:"status"`
	ExpiresAt     time.Time  `json:"expires_at"`
	PaymentURL    string     `json:"payment_url"`
	CreatedAt     time.Time  `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}
