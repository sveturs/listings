// backend/internal/proj/payments/service/mock_service.go

package service

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"backend/internal/domain/models"
)

// MockPaymentService представляет mock-реализацию для разработки
type MockPaymentService struct {
	frontendURL string
}

// NewMockPaymentService создает новый mock payment service
func NewMockPaymentService(frontendURL string) PaymentServiceInterface {
	return &MockPaymentService{
		frontendURL: frontendURL,
	}
}

// extractLocaleFromURL извлекает локаль из URL
func (m *MockPaymentService) extractLocaleFromURL(returnURL string) string {
	// Ищем паттерн /{locale}/ в URL после домена
	// Поддерживаем как двухбуквенные (en, ru) так и региональные (en-US, ru-RU) локали
	re := regexp.MustCompile(`/([a-z]{2}(?:-[A-Z]{2})?)/`)
	matches := re.FindStringSubmatch(returnURL)
	if len(matches) > 1 {
		return matches[1]
	}
	// По умолчанию возвращаем 'en'
	return "en"
}

// CreatePaymentSession создает mock платежную сессию
func (m *MockPaymentService) CreatePaymentSession(ctx context.Context, userID int, amount float64, currency, method string) (*models.PaymentSession, error) {
	log.Printf("MockPaymentService: Creating payment session for user %d, amount %f %s, method %s", userID, amount, currency, method)

	// Генерируем mock данные
	sessionID := fmt.Sprintf("mock_session_%d_%d", userID, time.Now().Unix())

	// Извлекаем локаль из return_url в контексте
	locale := "en" // дефолтное значение
	if returnURL, ok := ctx.Value("return_url").(string); ok && returnURL != "" {
		locale = m.extractLocaleFromURL(returnURL)
		log.Printf("MockPaymentService: Extracted locale '%s' from return_url: %s", locale, returnURL)
	}

	// Генерируем URL с правильной локалью
	paymentURL := fmt.Sprintf("%s/%s/payment/mock?session_id=%s&amount=%f&currency=%s", m.frontendURL, locale, sessionID, amount, currency)

	session := &models.PaymentSession{
		ID:            sessionID,
		UserID:        userID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: method,
		Status:        "pending",
		PaymentURL:    paymentURL,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(30 * time.Minute), // 30 минут на оплату
	}

	log.Printf("MockPaymentService: Created payment session: %+v", session)
	return session, nil
}

// CreateOrderPayment создает mock платежную сессию для заказа
func (m *MockPaymentService) CreateOrderPayment(ctx context.Context, orderID int, userID int, amount float64, currency, method string) (*models.PaymentSession, error) {
	log.Printf("MockPaymentService: Creating order payment for order %d, user %d, amount %f %s, method %s", orderID, userID, amount, currency, method)

	// Генерируем mock данные для заказа
	sessionID := fmt.Sprintf("mock_order_session_%d_%d_%d", orderID, userID, time.Now().Unix())

	// Извлекаем локаль из return_url в контексте (для заказов также может понадобиться)
	locale := "en" // дефолтное значение
	if returnURL, ok := ctx.Value("return_url").(string); ok && returnURL != "" {
		locale = m.extractLocaleFromURL(returnURL)
		log.Printf("MockPaymentService: Extracted locale '%s' from return_url for order: %s", locale, returnURL)
	}

	paymentURL := fmt.Sprintf("%s/%s/payment/mock?session_id=%s&amount=%f&currency=%s&order_id=%d", m.frontendURL, locale, sessionID, amount, currency, orderID)

	session := &models.PaymentSession{
		ID:            sessionID,
		UserID:        userID,
		OrderID:       &orderID,
		Amount:        amount,
		Currency:      currency,
		PaymentMethod: method,
		Status:        "pending",
		PaymentURL:    paymentURL,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(30 * time.Minute),
	}

	log.Printf("MockPaymentService: Created order payment session: %+v", session)
	return session, nil
}

// HandleWebhook обрабатывает mock webhook для баланса
func (m *MockPaymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	log.Printf("MockPaymentService: Handling balance webhook, payload size: %d", len(payload))
	// В реальной реализации здесь будет обработка webhook от AllSecure
	// Для mock просто логируем
	return nil
}

// HandleOrderPaymentWebhook обрабатывает mock webhook для заказов
func (m *MockPaymentService) HandleOrderPaymentWebhook(ctx context.Context, payload []byte, signature string) error {
	log.Printf("MockPaymentService: Handling order payment webhook, payload size: %d", len(payload))
	// В реальной реализации здесь будет обработка webhook от AllSecure для заказов
	// Для mock просто логируем
	return nil
}
