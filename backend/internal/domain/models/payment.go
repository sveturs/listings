// backend/internal/domain/models/payment.go

package models

import "time"

type PaymentSession struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	ExternalID    string    `json:"external_id"`
	Status        string    `json:"status"`
	ExpiresAt     time.Time `json:"expires_at"`
	PaymentURL    string    `json:"payment_url"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
}