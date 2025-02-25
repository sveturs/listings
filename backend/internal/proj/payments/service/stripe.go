// backend/internal/proj/payments/service/stripe.go

package service

import (
	"backend/internal/domain/models"
	"context"
	"fmt"
	"log"
	    "strconv"
     "strings"
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
    // Инициализируем Stripe с API-ключом
    stripe.Key = apiKey
    
    log.Printf("Initializing Stripe with key: %s...", apiKey[:10]) 
    
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
    successURL := fmt.Sprintf("%s/balance?success=true&session_id={CHECKOUT_SESSION_ID}", s.frontendURL)
    
    // Получаем токен сессии из контекста
    sessionToken := ctx.Value("session_token")
    if sessionToken == nil {
        sessionToken = ""
    }
    
    // Метаданные платежа
    metadata := map[string]string{
        "user_id": fmt.Sprintf("%d", userID),
        "method":  method,
    }
    
    // Если есть токен сессии, добавляем его в метаданные
    if sessionToken != nil && sessionToken.(string) != "" {
        metadata["session_token"] = sessionToken.(string)
        // Добавляем токен сессии в URL
        successURL = fmt.Sprintf("%s/balance?success=true&session_id={CHECKOUT_SESSION_ID}&session_token=%s", 
                               s.frontendURL, sessionToken.(string))
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
        Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
        SuccessURL: stripe.String(successURL),
        CancelURL:  stripe.String(fmt.Sprintf("%s/balance?canceled=true", s.frontendURL)),
        ClientReferenceID: stripe.String(fmt.Sprintf("user_%d", userID)),
        PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
            Metadata: metadata,
        },
        Metadata: metadata,
    }
	
	// Добавляем метаданные через параметр
	if sessionToken != nil && sessionToken.(string) != "" {
		params.AddMetadata("user_id", fmt.Sprintf("%d", userID))
		params.AddMetadata("method", method)
		params.AddMetadata("session_token", sessionToken.(string))
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
		if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
			return fmt.Errorf("error unmarshalling session: %w", err)
		}
	
		log.Printf("Processing completed checkout session: %s", checkoutSession.ID)
		
		// Подробный вывод для отладки
		log.Printf("Session metadata: %+v", checkoutSession.Metadata)
		log.Printf("Client reference ID: %s", checkoutSession.ClientReferenceID)
		log.Printf("Payment intent: %+v", checkoutSession.PaymentIntent)
		
		var userID int
		var method string
		
		// Извлекаем userID из метаданных
		if userIDStr, ok := checkoutSession.Metadata["user_id"]; ok {
			if id, err := strconv.Atoi(userIDStr); err == nil {
				userID = id
				log.Printf("Found user_id=%d in session metadata", userID)
			}
		}
		
		// Если не нашли в метаданных, извлекаем из ClientReferenceID
		if userID == 0 && checkoutSession.ClientReferenceID != "" {
			parts := strings.Split(checkoutSession.ClientReferenceID, "_")
			if len(parts) >= 2 {
				userIDStr := parts[1]
				if id, err := strconv.Atoi(userIDStr); err == nil {
					userID = id
					log.Printf("Found user_id=%d in client reference ID", userID)
				}
			}
		}
		
		// Проверяем payment intent data для дополнительных метаданных
		if checkoutSession.PaymentIntent != nil && checkoutSession.PaymentIntent.Metadata != nil {
			if userIDStr, ok := checkoutSession.PaymentIntent.Metadata["user_id"]; ok {
				if id, err := strconv.Atoi(userIDStr); err == nil {
					userID = id
					log.Printf("Found user_id=%d in payment intent metadata", userID)
				}
			}
		}
		
		// Если всё ещё не нашли, проверяем payment intent ID
		if userID == 0 && checkoutSession.PaymentIntent != nil {
			// Получаем платежное намерение по ID
			pi, err := payment.PaymentIntents.Get(checkoutSession.PaymentIntent.ID, nil)
			if err == nil && pi.Metadata != nil {
				if userIDStr, ok := pi.Metadata["user_id"]; ok {
					if id, err := strconv.Atoi(userIDStr); err == nil {
						userID = id
						log.Printf("Found user_id=%d in payment intent object", userID)
					}
				}
			}
		}
		
		// Если не смогли найти ID пользователя, логируем ошибку
		if userID == 0 {
			log.Printf("ERROR: Could not determine user_id from Stripe session data")
			return fmt.Errorf("failed to determine user ID from session data")
		}
		
		// Устанавливаем метод оплаты
		method = "bank_transfer" // Значение по умолчанию
		if methodStr, ok := checkoutSession.Metadata["method"]; ok {
			method = methodStr
		}			
        
        // Создаем транзакцию пополнения баланса
        amount := float64(checkoutSession.AmountTotal) / 100 // Конвертируем из центов
        transaction, err := s.balanceService.CreateDeposit(ctx, userID, amount, method)
        if err != nil {
            return fmt.Errorf("failed to create deposit: %w", err)
        }

        log.Printf("Successfully processed payment for user %d: amount=%f, method=%s, transaction_id=%d", 
            userID, amount, method, transaction.ID)


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