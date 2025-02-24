// backend/internal/proj/payments/service/stripe.go

package service

import (
	"backend/internal/domain/models"
	"context"
	"fmt"
	"log"
	"time"
	"encoding/json" 
    balanceService "backend/internal/proj/balance/service" 
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripeService struct {
	apiKey        string
	webhookSecret string
	frontendURL   string
    balanceService  balanceService.BalanceServiceInterface
}

func NewStripeService(apiKey, webhookSecret, frontendURL string, balanceService balanceService.BalanceServiceInterface) *StripeService {
    stripe.Key = apiKey
    return &StripeService{
        apiKey:        apiKey,
        webhookSecret: webhookSecret,
        frontendURL:   frontendURL,
        balanceService: balanceService,
    }
}

// Создаем сессию оплаты
func (s *StripeService) CreatePaymentSession(ctx context.Context, userID int, amount float64, currency, method string) (*models.PaymentSession, error) {
	// Конвертируем в минимальные единицы (центы)
	amountInCents := int64(amount * 100)

	// Создаем параметры сессии
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Пополнение баланса"),
					},
					UnitAmount: stripe.Int64(amountInCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s/balance?success=true", s.frontendURL)),
		CancelURL:  stripe.String(fmt.Sprintf("%s/balance?canceled=true", s.frontendURL)),
		ClientReferenceID: stripe.String(fmt.Sprintf("user_%d", userID)),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"user_id": fmt.Sprintf("%d", userID),
				"method":  method,
			},
		},
	}

	// Создаем сессию в Stripe
	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe session: %w", err)
	}

	// Создаем запись о платежной сессии
	paymentSession := &models.PaymentSession{
		UserID:         userID,
		Amount:         amount,
		Currency:       currency,
		PaymentMethod:  method,
		ExternalID:     sess.ID,
		Status:         "pending",
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		PaymentURL:     sess.URL,
	}

	return paymentSession, nil
}

// Обрабатываем вебхук от Stripe
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
    event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
    if err != nil {
        return fmt.Errorf("webhook error: %w", err)
    }

    log.Printf("Received Stripe webhook: %s", event.Type)

    switch event.Type {
    case "checkout.session.completed":
        var checkoutSession stripe.CheckoutSession
        // Используем правильный метод для разбора JSON
        if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
            return fmt.Errorf("error unmarshalling session: %w", err)
        }

        // Получаем ID пользователя из метаданных
        userIDStr := checkoutSession.PaymentIntent.Metadata["user_id"]
        method := checkoutSession.PaymentIntent.Metadata["method"]
        
        var userID int
        _, err = fmt.Sscanf(userIDStr, "%d", &userID)
        if err != nil {
            return fmt.Errorf("invalid user_id: %w", err)
        }

        // Создаем транзакцию пополнения баланса
        amount := float64(checkoutSession.AmountTotal) / 100 // Конвертируем из центов
        _, err = s.balanceService.CreateDeposit(ctx, userID, amount, method)
        if err != nil {
            return fmt.Errorf("failed to create deposit: %w", err)
        }

        log.Printf("Successfully processed payment for user %d: amount=%f, method=%s", 
            userID, amount, method)

		case "payment_intent.payment_failed":
			var paymentIntent stripe.PaymentIntent
 			if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
				return fmt.Errorf("error unmarshalling payment intent: %w", err)
			}

        if paymentIntent.LastPaymentError != nil {
            log.Printf("Payment failed: %s", paymentIntent.LastPaymentError.Error) // Исправлено на Error вместо Message
        } else {
            log.Printf("Payment failed without specific error message")
        }
    }

    return nil
}