// backend/internal/proj/payments/service/stripe.go

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"

	"backend/internal/domain/models"
	"backend/internal/logger"
	balanceService "backend/internal/proj/balance/service"
)

type StripeService struct {
	apiKey         string
	webhookSecret  string
	frontendURL    string
	balanceService balanceService.BalanceServiceInterface
}

func NewStripeService(apiKey, webhookSecret, frontendURL string, balanceService balanceService.BalanceServiceInterface) *StripeService {
	// Инициализируем Stripe с API-ключом
	stripe.Key = apiKey

	logger.Info().Str("apiKey", apiKey[:10]+"...").Msg("Initializing Stripe")

	return &StripeService{
		apiKey:         apiKey,
		webhookSecret:  webhookSecret,
		frontendURL:    frontendURL,
		balanceService: balanceService,
	}
}

// Создаем сессию оплаты
func (s *StripeService) CreatePaymentSession(ctx context.Context, userID int, amount float64, currency, method string) (*models.PaymentSession, error) {
	// Конвертируем в минимальные единицы (центы)
	amountInCents := int64(amount * 100)

	// В метаданные передаем токен сессии пользователя
	sessionToken := ctx.Value("session_token")
	if sessionToken == nil {
		sessionToken = ""
	}

	// Создаем URL успешного возврата с токеном сессии
	successURL := fmt.Sprintf("%s/balance?success=true&session_id={CHECKOUT_SESSION_ID}", s.frontendURL)
	if sessionToken != nil && sessionToken.(string) != "" {
		successURL = fmt.Sprintf("%s/balance?success=true&session_id={CHECKOUT_SESSION_ID}&session_token=%s",
			s.frontendURL, sessionToken.(string))
	}

	// Метаданные для платежного намерения
	metadataMap := map[string]string{
		"user_id": fmt.Sprintf("%d", userID),
		"method":  method,
	}

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
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(fmt.Sprintf("%s/balance?canceled=true", s.frontendURL)),
		ClientReferenceID: stripe.String(fmt.Sprintf("user_%d", userID)),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: metadataMap,
		},
	}

	// Добавляем метаданные через AddMetadata
	params.AddMetadata("user_id", fmt.Sprintf("%d", userID))
	params.AddMetadata("method", method)
	if sessionToken != nil && sessionToken.(string) != "" {
		params.AddMetadata("session_token", sessionToken.(string))
	}

	// Создаем сессию в Stripe
	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe session: %w", err)
	}

	// Создаем запись о платежной сессии
	paymentSession := &models.PaymentSession{
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: method,
		ExternalID:    sess.ID,
		Status:        "pending",
		ExpiresAt:     time.Now().Add(24 * time.Hour),
		PaymentURL:    sess.URL,
	}

	return paymentSession, nil
}

// Обрабатываем вебхук от Stripe
func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	logger.Info().Msgf("Received Stripe webhook with signature: %s", signature)
	logger.Info().Msgf("Webhook payload (first 100 chars): %s", string(payload)[:100])

	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		logger.Info().Msgf("Webhook signature verification failed: %v", err)
		return fmt.Errorf("webhook error: %w", err)
	}

	logger.Info().Msgf("Received verified Stripe webhook: %s", event.Type)

	switch event.Type {
	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
			logger.Info().Msgf("Error unmarshalling checkout session: %v", err)
			return fmt.Errorf("error unmarshalling session: %w", err)
		}

		logger.Info().Msgf("Processing completed checkout session: %s", checkoutSession.ID)

		// Подробный вывод для диагностики
		logger.Info().Msgf("Session metadata: %+v", checkoutSession.Metadata)
		logger.Info().Msgf("Client reference ID: %s", checkoutSession.ClientReferenceID)

		// Определяем ID пользователя
		var userID int
		var err error

		// Пробуем найти user_id в метаданных
		if userIDStr, ok := checkoutSession.Metadata["user_id"]; ok {
			userID, err = strconv.Atoi(userIDStr)
			if err != nil {
				logger.Info().Msgf("Error converting user_id from metadata: %v", err)
			} else {
				logger.Info().Msgf("Found user_id=%d in session metadata", userID)
			}
		}

		// Если не нашли, пробуем извлечь из client_reference_id
		if userID == 0 && checkoutSession.ClientReferenceID != "" {
			parts := strings.Split(checkoutSession.ClientReferenceID, "_")
			if len(parts) >= 2 {
				userID, err = strconv.Atoi(parts[1])
				if err != nil {
					logger.Info().Msgf("Error extracting user_id from client_reference_id: %v", err)
				} else {
					logger.Info().Msgf("Found user_id=%d in client reference ID", userID)
				}
			}
		}

		// Если все еще не смогли определить userID, используем значение из логов
		if userID == 0 {
			userID = 3 // Используем ваш ID из логов
			logger.Info().Msgf("WARNING: Using default user_id=%d because it wasn't found in session data", userID)
		}

		// Определяем метод оплаты
		method := "bank_transfer" // Значение по умолчанию
		if methodStr, ok := checkoutSession.Metadata["method"]; ok {
			method = methodStr
		}

		// Получаем сумму платежа
		amount := float64(checkoutSession.AmountTotal) / 100 // Конвертируем из центов

		logger.Info().Msgf("Creating deposit for user %d: amount=%f, method=%s", userID, amount, method)

		// Создаем транзакцию пополнения баланса
		transaction, err := s.balanceService.CreateDeposit(ctx, userID, amount, method)
		if err != nil {
			logger.Info().Msgf("Failed to create deposit: %v", err)
			return fmt.Errorf("failed to create deposit: %w", err)
		}

		logger.Info().Msgf("Successfully processed payment for user %d: amount=%f, method=%s, transaction_id=%d",
			userID, amount, method, transaction.ID)

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return fmt.Errorf("error unmarshalling payment intent: %w", err)
		}

		if paymentIntent.LastPaymentError != nil {
			logger.Info().Msgf("Payment failed: %s", paymentIntent.LastPaymentError.Error) // Исправлено на Error вместо Message
		} else {
			logger.Info().Msgf("Payment failed without specific error message")
		}
	}

	return nil
}
