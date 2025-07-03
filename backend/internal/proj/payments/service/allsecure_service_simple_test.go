package service

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestCalculateCommissionLogic проверяет расчет комиссии маркетплейса
func TestCalculateCommissionLogic(t *testing.T) {
	// Создаем сервис с тестовой конфигурацией
	config := &AllSecureConfig{
		MarketplaceCommissionRate: decimal.NewFromFloat(0.05), // 5%
	}

	service := &AllSecureService{
		config:         config,
		commissionRate: config.MarketplaceCommissionRate,
	}

	testCases := []struct {
		name               string
		amount             decimal.Decimal
		expectedCommission decimal.Decimal
	}{
		{
			name:               "100 EUR with 5% commission",
			amount:             decimal.NewFromFloat(100.00),
			expectedCommission: decimal.NewFromFloat(5.00),
		},
		{
			name:               "50 EUR with 5% commission",
			amount:             decimal.NewFromFloat(50.00),
			expectedCommission: decimal.NewFromFloat(2.50),
		},
		{
			name:               "1000 EUR with 5% commission",
			amount:             decimal.NewFromFloat(1000.00),
			expectedCommission: decimal.NewFromFloat(50.00),
		},
		{
			name:               "Zero amount",
			amount:             decimal.Zero,
			expectedCommission: decimal.Zero,
		},
		{
			name:               "Small amount",
			amount:             decimal.NewFromFloat(0.01),
			expectedCommission: decimal.NewFromFloat(0.0005),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commission := service.calculateCommission(tc.amount)
			if !commission.Equal(tc.expectedCommission) {
				t.Errorf("Expected commission %s for amount %s, got %s",
					tc.expectedCommission.String(), tc.amount.String(), commission.String())
			}
		})
	}
}

// TestMapAllSecureStatus проверяет маппинг статусов AllSecure
func TestMapAllSecureStatus(t *testing.T) {
	service := &AllSecureService{}

	testCases := []struct {
		allsecureStatus string
		expectedStatus  string
	}{
		{"FINISHED", "captured"},
		{"PENDING", "authorized"},
		{"ERROR", "failed"},
		{"UNKNOWN", "pending"},
		{"", "pending"},
	}

	for _, tc := range testCases {
		t.Run(tc.allsecureStatus, func(t *testing.T) {
			status := service.mapAllSecureStatus(tc.allsecureStatus)
			if status != tc.expectedStatus {
				t.Errorf("Expected status %s for AllSecure status %s, got %s",
					tc.expectedStatus, tc.allsecureStatus, status)
			}
		})
	}
}

// TestValidatePaymentRequestLogic проверяет валидацию запросов
func TestValidatePaymentRequestLogic(t *testing.T) {
	service := &AllSecureService{}

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
			name: "Zero amount should fail",
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
			name: "Negative amount should fail",
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
			name: "Invalid user ID should fail",
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
			name: "Invalid listing ID should fail",
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
			name: "Empty currency should fail",
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
		{
			name: "Empty return URL should fail",
			request: CreatePaymentRequest{
				UserID:      100,
				ListingID:   200,
				Amount:      decimal.NewFromFloat(100.00),
				Currency:    "EUR",
				Description: "Test payment",
				ReturnURL:   "",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.validatePaymentRequest(nil, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("validatePaymentRequest() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// TestCommissionRateEdgeCases проверяет граничные случаи расчета комиссии
func TestCommissionRateEdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		commissionRate decimal.Decimal
		amount         decimal.Decimal
		expected       decimal.Decimal
	}{
		{
			name:           "Zero commission rate",
			commissionRate: decimal.Zero,
			amount:         decimal.NewFromFloat(100.00),
			expected:       decimal.Zero,
		},
		{
			name:           "100% commission rate",
			commissionRate: decimal.NewFromFloat(1.0),
			amount:         decimal.NewFromFloat(100.00),
			expected:       decimal.NewFromFloat(100.00),
		},
		{
			name:           "Very small commission rate",
			commissionRate: decimal.NewFromFloat(0.001), // 0.1%
			amount:         decimal.NewFromFloat(100.00),
			expected:       decimal.NewFromFloat(0.1),
		},
		{
			name:           "Very large amount",
			commissionRate: decimal.NewFromFloat(0.05), // 5%
			amount:         decimal.NewFromFloat(1000000.00),
			expected:       decimal.NewFromFloat(50000.00),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &AllSecureService{
				commissionRate: tc.commissionRate,
			}

			commission := service.calculateCommission(tc.amount)
			if !commission.Equal(tc.expected) {
				t.Errorf("Expected commission %s, got %s",
					tc.expected.String(), commission.String())
			}
		})
	}
}
