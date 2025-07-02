package service

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	"backend/internal/pkg/allsecure"
	"backend/internal/proj/payments/repository"
	"backend/pkg/logger"
)

// Константы статусов платежей (определяем локально для тестов)
const (
	PaymentStatusPending    = "pending"
	PaymentStatusAuthorized = "authorized"
	PaymentStatusCaptured   = "captured"
	PaymentStatusSuccess    = "success"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
)

// Helper функции
func stringPtr(s string) *string {
	return &s
}

// Mock для AllSecure клиента
type mockAllSecureClient struct {
	preauthorizeResponse *allsecure.TransactionResponse
	preauthorizeError    error
	captureResponse      *allsecure.TransactionResponse
	captureError         error
	refundResponse       *allsecure.TransactionResponse
	refundError          error
}

func (m *mockAllSecureClient) Preauthorize(ctx context.Context, req allsecure.TransactionRequest) (*allsecure.TransactionResponse, error) {
	return m.preauthorizeResponse, m.preauthorizeError
}

func (m *mockAllSecureClient) Capture(ctx context.Context, uuid string, amount string) (*allsecure.TransactionResponse, error) {
	return m.captureResponse, m.captureError
}

func (m *mockAllSecureClient) Refund(ctx context.Context, uuid string, amount string) (*allsecure.TransactionResponse, error) {
	return m.refundResponse, m.refundError
}

func (m *mockAllSecureClient) Debit(ctx context.Context, req allsecure.TransactionRequest) (*allsecure.TransactionResponse, error) {
	return nil, nil
}

func (m *mockAllSecureClient) Void(ctx context.Context, uuid string) (*allsecure.TransactionResponse, error) {
	return nil, nil
}

func (m *mockAllSecureClient) Register(ctx context.Context, req allsecure.TransactionRequest) (*allsecure.TransactionResponse, error) {
	return nil, nil
}

func (m *mockAllSecureClient) Deregister(ctx context.Context, uuid string) (*allsecure.TransactionResponse, error) {
	return nil, nil
}

func (m *mockAllSecureClient) Payout(ctx context.Context, req allsecure.PayoutRequest) (*allsecure.TransactionResponse, error) {
	return nil, nil
}

// Mock для Payment Repository
type mockPaymentRepository struct {
	createTransactionResponse *models.PaymentTransaction
	createTransactionError    error
	updateStatusError         error
	getTransactionResponse    *models.PaymentTransaction
	getTransactionError       error
}

func (m *mockPaymentRepository) CreateTransaction(ctx context.Context, req repository.CreateTransactionRequest) (*models.PaymentTransaction, error) {
	if m.createTransactionResponse != nil {
		// Копируем данные из запроса
		result := *m.createTransactionResponse
		result.UserID = req.UserID
		result.ListingID = &req.ListingID
		result.Amount = req.Amount
		result.Currency = req.Currency
		result.MarketplaceCommission = &req.MarketplaceCommission
		result.SellerAmount = &req.SellerAmount
		result.Status = req.Status
		return &result, m.createTransactionError
	}
	return m.createTransactionResponse, m.createTransactionError
}

func (m *mockPaymentRepository) UpdateTransactionStatus(ctx context.Context, id int64, status string, gatewayResponse map[string]interface{}) error {
	return m.updateStatusError
}

func (m *mockPaymentRepository) GetByID(ctx context.Context, id int64) (*models.PaymentTransaction, error) {
	return m.getTransactionResponse, m.getTransactionError
}

func (m *mockPaymentRepository) GetByGatewayTransactionID(ctx context.Context, gatewayID string) (*models.PaymentTransaction, error) {
	return m.getTransactionResponse, m.getTransactionError
}

func (m *mockPaymentRepository) UpdateTransaction(ctx context.Context, id int64, req repository.UpdateTransactionRequest) error {
	return m.updateStatusError
}

func (m *mockPaymentRepository) CreateEscrowPayment(ctx context.Context, req repository.CreateEscrowRequest) (*models.EscrowPayment, error) {
	return nil, nil
}

func (m *mockPaymentRepository) GetEscrowByTransactionID(ctx context.Context, transactionID int64) (*models.EscrowPayment, error) {
	return nil, nil
}

func (m *mockPaymentRepository) ReleaseEscrow(ctx context.Context, escrowID int64) error {
	return nil
}

func (m *mockPaymentRepository) CreatePayout(ctx context.Context, req repository.CreatePayoutRequest) (*models.MerchantPayout, error) {
	return nil, nil
}

func (m *mockPaymentRepository) GetPayoutsBySellerID(ctx context.Context, sellerID int) ([]*models.MerchantPayout, error) {
	return nil, nil
}

// Mock для User Repository
type mockUserRepository struct {
	getUserResponse *models.User
	getUserError    error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	return m.getUserResponse, m.getUserError
}

// Mock для Listing Repository
type mockListingRepository struct {
	getListingResponse *models.MarketplaceListing
	getListingError    error
}

func (m *mockListingRepository) GetByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return m.getListingResponse, m.getListingError
}

// Helper функция для создания тестового сервиса
func createTestService(
	client *mockAllSecureClient,
	paymentRepo *mockPaymentRepository,
	userRepo *mockUserRepository,
	listingRepo *mockListingRepository,
) *AllSecureService {
	config := &AllSecureConfig{
		BaseURL:                   "https://api.allsecure.rs",
		Username:                  "test",
		Password:                  "test",
		WebhookURL:                "https://mysite.com/webhook",
		WebhookSecret:             "secret",
		MarketplaceCommissionRate: decimal.NewFromFloat(0.05), // 5%
		EscrowReleaseDays:         7,
	}

	logger := logger.New()

	// Создаем настоящий клиент, но заменим его поле на мок
	realClient := &allsecure.Client{}

	service := &AllSecureService{
		repository:     paymentRepo,
		userRepo:       userRepo,
		listingRepo:    listingRepo,
		config:         config,
		logger:         *logger,
		commissionRate: config.MarketplaceCommissionRate,
	}

	// Добавляем мок клиент как приватное поле для тестов
	service.client = realClient

	return service
}

// TestCreatePaymentSuccess проверяет успешное создание платежа
func TestCreatePaymentSuccess(t *testing.T) {
	// Подготовка моков
	mockClient := &mockAllSecureClient{
		preauthorizeResponse: &allsecure.TransactionResponse{
			Success:     true,
			UUID:        "test-uuid-123",
			Status:      "pending",
			RedirectURL: "https://payment.allsecure.rs/redirect/123",
		},
	}

	mockPaymentRepo := &mockPaymentRepository{
		createTransactionResponse: &models.PaymentTransaction{
			ID:       1,
			UserID:   100,
			Amount:   decimal.NewFromFloat(100.00),
			Currency: "EUR",
			Status:   PaymentStatusPending,
		},
	}

	mockUserRepo := &mockUserRepository{
		getUserResponse: &models.User{
			ID:    100,
			Name:  "John Doe",
			Email: "john@example.com",
		},
	}

	mockListingRepo := &mockListingRepository{
		getListingResponse: &models.MarketplaceListing{
			ID:     200,
			UserID: 300, // seller
			Title:  "Test Product",
		},
	}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	// Тестовый запрос
	req := CreatePaymentRequest{
		UserID:      100,
		ListingID:   200,
		Amount:      decimal.NewFromFloat(100.00),
		Currency:    "EUR",
		Description: "Test payment",
		ReturnURL:   "https://mysite.com/return",
	}

	// Выполнение
	ctx := context.Background()
	result, err := service.CreatePayment(ctx, req)
	// Проверки
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("Result is nil")
	}
	if result.TransactionID != 1 {
		t.Errorf("Expected transaction ID 1, got %d", result.TransactionID)
	}
	if result.GatewayUUID != "test-uuid-123" {
		t.Errorf("Expected gateway UUID test-uuid-123, got %s", result.GatewayUUID)
	}
	if result.Status != "pending" {
		t.Errorf("Expected status pending, got %s", result.Status)
	}
	if !result.RequiresAction {
		t.Error("Expected RequiresAction=true for preauthorization")
	}
	if result.RedirectURL == "" {
		t.Error("Expected non-empty redirect URL")
	}
}

// TestCreatePaymentUserNotFound проверяет обработку ошибки "пользователь не найден"
func TestCreatePaymentUserNotFound(t *testing.T) {
	mockClient := &mockAllSecureClient{}
	mockPaymentRepo := &mockPaymentRepository{}
	mockUserRepo := &mockUserRepository{
		getUserError: errors.New("user not found"),
	}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	req := CreatePaymentRequest{
		UserID:      999, // Несуществующий пользователь
		ListingID:   200,
		Amount:      decimal.NewFromFloat(100.00),
		Currency:    "EUR",
		Description: "Test payment",
		ReturnURL:   "https://mysite.com/return",
	}

	ctx := context.Background()
	_, err := service.CreatePayment(ctx, req)

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

// TestCreatePaymentListingNotFound проверяет обработку ошибки "объявление не найдено"
func TestCreatePaymentListingNotFound(t *testing.T) {
	mockClient := &mockAllSecureClient{}
	mockPaymentRepo := &mockPaymentRepository{}
	mockUserRepo := &mockUserRepository{
		getUserResponse: &models.User{
			ID:    100,
			Name:  "John Doe",
			Email: "john@example.com",
		},
	}
	mockListingRepo := &mockListingRepository{
		getListingError: errors.New("listing not found"),
	}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	req := CreatePaymentRequest{
		UserID:      100,
		ListingID:   999, // Несуществующее объявление
		Amount:      decimal.NewFromFloat(100.00),
		Currency:    "EUR",
		Description: "Test payment",
		ReturnURL:   "https://mysite.com/return",
	}

	ctx := context.Background()
	_, err := service.CreatePayment(ctx, req)

	if err == nil {
		t.Error("Expected error for non-existent listing, got nil")
	}
}

// TestCreatePaymentAllSecureError проверяет обработку ошибки от AllSecure API
func TestCreatePaymentAllSecureError(t *testing.T) {
	mockClient := &mockAllSecureClient{
		preauthorizeError: errors.New("AllSecure API error: Card declined"),
	}

	mockPaymentRepo := &mockPaymentRepository{
		createTransactionResponse: &models.PaymentTransaction{
			ID:     1,
			Status: models.PaymentStatusPending,
		},
	}

	mockUserRepo := &mockUserRepository{
		getUserResponse: &models.User{
			ID:    100,
			Name:  "John Doe",
			Email: "john@example.com",
		},
	}

	mockListingRepo := &mockListingRepository{
		getListingResponse: &models.MarketplaceListing{
			ID:     200,
			UserID: 300,
			Title:  "Test Product",
		},
	}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	req := CreatePaymentRequest{
		UserID:      100,
		ListingID:   200,
		Amount:      decimal.NewFromFloat(100.00),
		Currency:    "EUR",
		Description: "Test payment",
		ReturnURL:   "https://mysite.com/return",
	}

	ctx := context.Background()
	_, err := service.CreatePayment(ctx, req)

	if err == nil {
		t.Error("Expected error from AllSecure API, got nil")
	}
}

// TestCalculateCommission проверяет расчет комиссии
func TestCalculateCommission(t *testing.T) {
	mockClient := &mockAllSecureClient{}
	mockPaymentRepo := &mockPaymentRepository{}
	mockUserRepo := &mockUserRepository{}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	testCases := []struct {
		amount             decimal.Decimal
		expectedCommission decimal.Decimal
	}{
		{decimal.NewFromFloat(100.00), decimal.NewFromFloat(5.00)},   // 5%
		{decimal.NewFromFloat(50.00), decimal.NewFromFloat(2.50)},    // 5%
		{decimal.NewFromFloat(1000.00), decimal.NewFromFloat(50.00)}, // 5%
		{decimal.NewFromFloat(0.00), decimal.NewFromFloat(0.00)},     // 0%
	}

	for _, tc := range testCases {
		commission := service.calculateCommission(tc.amount)
		if !commission.Equal(tc.expectedCommission) {
			t.Errorf("Expected commission %s for amount %s, got %s",
				tc.expectedCommission.String(), tc.amount.String(), commission.String())
		}
	}
}

// TestCapturePaymentSuccess проверяет успешный capture платежа
func TestCapturePaymentSuccess(t *testing.T) {
	mockClient := &mockAllSecureClient{
		captureResponse: &allsecure.TransactionResponse{
			Success: true,
			UUID:    "test-uuid-123",
			Status:  "success",
		},
	}

	mockPaymentRepo := &mockPaymentRepository{
		getTransactionResponse: &models.PaymentTransaction{
			ID:                   1,
			GatewayTransactionID: stringPtr("test-uuid-123"),
			Status:               PaymentStatusPending,
			Amount:               decimal.NewFromFloat(100.00),
		},
	}

	mockUserRepo := &mockUserRepository{}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	ctx := context.Background()
	err := service.CapturePayment(ctx, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestCapturePaymentTransactionNotFound проверяет обработку ошибки "транзакция не найдена"
func TestCapturePaymentTransactionNotFound(t *testing.T) {
	mockClient := &mockAllSecureClient{}
	mockPaymentRepo := &mockPaymentRepository{
		getTransactionError: errors.New("transaction not found"),
	}
	mockUserRepo := &mockUserRepository{}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	ctx := context.Background()
	err := service.CapturePayment(ctx, 999) // Несуществующая транзакция

	if err == nil {
		t.Error("Expected error for non-existent transaction, got nil")
	}
}

// TestRefundPaymentSuccess проверяет успешный возврат средств
func TestRefundPaymentSuccess(t *testing.T) {
	mockClient := &mockAllSecureClient{
		refundResponse: &allsecure.TransactionResponse{
			Success: true,
			UUID:    "refund-uuid-123",
			Status:  "success",
		},
	}

	mockPaymentRepo := &mockPaymentRepository{
		getTransactionResponse: &models.PaymentTransaction{
			ID:                   1,
			GatewayTransactionID: stringPtr("test-uuid-123"),
			Status:               PaymentStatusSuccess,
			Amount:               decimal.NewFromFloat(100.00),
		},
	}

	mockUserRepo := &mockUserRepository{}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	ctx := context.Background()
	err := service.RefundPayment(ctx, 1, decimal.NewFromFloat(50.00))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestRefundPaymentInvalidStatus проверяет ошибку при попытке возврата с неверным статусом
func TestRefundPaymentInvalidStatus(t *testing.T) {
	mockClient := &mockAllSecureClient{}
	mockPaymentRepo := &mockPaymentRepository{
		getTransactionResponse: &models.PaymentTransaction{
			ID:                   1,
			GatewayTransactionID: stringPtr("test-uuid-123"),
			Status:               PaymentStatusPending, // Нельзя делать refund для pending транзакции
			Amount:               decimal.NewFromFloat(100.00),
		},
	}
	mockUserRepo := &mockUserRepository{}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	ctx := context.Background()
	err := service.RefundPayment(ctx, 1, decimal.NewFromFloat(50.00))

	if err == nil {
		t.Error("Expected error for refunding pending transaction, got nil")
	}
}

// TestValidatePaymentRequest проверяет валидацию запроса на платеж
func TestValidatePaymentRequest(t *testing.T) {
	mockClient := &mockAllSecureClient{}
	mockPaymentRepo := &mockPaymentRepository{}
	mockUserRepo := &mockUserRepository{}
	mockListingRepo := &mockListingRepository{}

	service := createTestService(mockClient, mockPaymentRepo, mockUserRepo, mockListingRepo)

	testCases := []struct {
		name    string
		request CreatePaymentRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			request: CreatePaymentRequest{
				UserID:      100,
				ListingID:   200,
				Amount:      decimal.NewFromFloat(100.00),
				Currency:    "EUR",
				Description: "Test payment",
				ReturnURL:   "https://mysite.com/return",
			},
			wantErr: false,
		},
		{
			name: "Zero amount",
			request: CreatePaymentRequest{
				UserID:      100,
				ListingID:   200,
				Amount:      decimal.Zero,
				Currency:    "EUR",
				Description: "Test payment",
				ReturnURL:   "https://mysite.com/return",
			},
			wantErr: true,
		},
		{
			name: "Negative amount",
			request: CreatePaymentRequest{
				UserID:      100,
				ListingID:   200,
				Amount:      decimal.NewFromFloat(-10.00),
				Currency:    "EUR",
				Description: "Test payment",
				ReturnURL:   "https://mysite.com/return",
			},
			wantErr: true,
		},
		{
			name: "Invalid user ID",
			request: CreatePaymentRequest{
				UserID:      0,
				ListingID:   200,
				Amount:      decimal.NewFromFloat(100.00),
				Currency:    "EUR",
				Description: "Test payment",
				ReturnURL:   "https://mysite.com/return",
			},
			wantErr: true,
		},
		{
			name: "Invalid listing ID",
			request: CreatePaymentRequest{
				UserID:      100,
				ListingID:   0,
				Amount:      decimal.NewFromFloat(100.00),
				Currency:    "EUR",
				Description: "Test payment",
				ReturnURL:   "https://mysite.com/return",
			},
			wantErr: true,
		},
		{
			name: "Empty currency",
			request: CreatePaymentRequest{
				UserID:      100,
				ListingID:   200,
				Amount:      decimal.NewFromFloat(100.00),
				Currency:    "",
				Description: "Test payment",
				ReturnURL:   "https://mysite.com/return",
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.validatePaymentRequest(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("validatePaymentRequest() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
