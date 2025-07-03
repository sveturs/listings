// backend/internal/proj/payments/service/interface.go

package service

import (
	"context"

	"backend/internal/domain/models"
)

type PaymentServiceInterface interface {
	CreatePaymentSession(ctx context.Context, userID int, amount float64, currency, method string) (*models.PaymentSession, error)
	CreateOrderPayment(ctx context.Context, orderID int, userID int, amount float64, currency, method string) (*models.PaymentSession, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
	HandleOrderPaymentWebhook(ctx context.Context, payload []byte, signature string) error
}
