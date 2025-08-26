package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	paymentsService "backend/internal/proj/payments/service"
)

// MockPaymentService implements payment operations for testing
type MockPaymentService struct{}

// NewMockPaymentService creates a new mock payment service
func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{}
}

// CreatePayment mocks payment creation
func (m *MockPaymentService) CreatePayment(ctx context.Context, req *paymentsService.CreatePaymentRequest) (string, string, error) {
	// Generate mock payment intent ID
	paymentIntentID := "pi_" + uuid.New().String()

	// For testing, redirect to a success page
	redirectURL := fmt.Sprintf("/subscription/payment-mock?payment_intent=%s&amount=%s",
		paymentIntentID,
		req.Amount.String(),
	)

	return paymentIntentID, redirectURL, nil
}

// InitiatePaymentMock initiates mock payment for subscription
func InitiatePaymentMock(ctx context.Context, userID int, planCode string, billingCycle models.BillingCycle, amount decimal.Decimal, returnURL string) (*PaymentInitiationResponse, error) {
	// For free plan
	if amount.IsZero() {
		return &PaymentInitiationResponse{
			PaymentRequired: false,
			Message:         "Free plan selected",
		}, nil
	}

	// Generate mock payment intent
	paymentIntentID := "pi_mock_" + uuid.New().String()

	// For development, redirect to a mock payment page
	redirectURL := fmt.Sprintf("%s?payment_intent=%s&amount=%s&plan=%s&cycle=%s",
		"/subscription/payment-mock",
		paymentIntentID,
		amount.String(),
		planCode,
		billingCycle,
	)

	return &PaymentInitiationResponse{
		PaymentRequired: true,
		PaymentIntentID: paymentIntentID,
		RedirectURL:     redirectURL,
		Amount:          amount,
		Currency:        "EUR",
	}, nil
}
