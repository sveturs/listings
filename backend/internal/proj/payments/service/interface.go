// backend/internal/proj/payments/service/interface.go

package service

import (
	"backend/internal/domain/models"
	"context"
)

type PaymentServiceInterface interface {
	CreatePaymentSession(ctx context.Context, userID int, amount float64, currency, method string) (*models.PaymentSession, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
}
