package postexpress

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config содержит конфигурацию для Post Express API
type Config struct {
	APIURL        string
	Username      string
	Password      string
	Brand         string
	Warehouse     string
	Timeout       time.Duration
	RetryAttempts int
	IsProduction  bool

	// COD (Otkupnina) Configuration
	BankAccount  string // Банковский счёт для перевода откупнины
	PaymentCode  string // Шифра плаћања (обычно "189")
	PaymentModel string // Модель платежа (обычно "97")
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	apiURL := os.Getenv("POST_EXPRESS_API_URL")
	if apiURL == "" {
		apiURL = "https://wsp-test.posta.rs/api" // default to test environment
	}

	username := os.Getenv("POST_EXPRESS_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("POST_EXPRESS_USERNAME environment variable is required")
	}

	password := os.Getenv("POST_EXPRESS_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("POST_EXPRESS_PASSWORD environment variable is required")
	}

	brand := os.Getenv("POST_EXPRESS_BRAND")
	if brand == "" {
		brand = "SVETU" // default brand
	}

	warehouse := os.Getenv("POST_EXPRESS_WAREHOUSE")
	if warehouse == "" {
		warehouse = brand // default warehouse same as brand
	}

	timeoutSeconds := 30
	if timeoutStr := os.Getenv("POST_EXPRESS_TIMEOUT_SECONDS"); timeoutStr != "" {
		if parsed, err := strconv.Atoi(timeoutStr); err == nil {
			timeoutSeconds = parsed
		}
	}

	retryAttempts := 3
	if retryStr := os.Getenv("POST_EXPRESS_RETRY_ATTEMPTS"); retryStr != "" {
		if parsed, err := strconv.Atoi(retryStr); err == nil {
			retryAttempts = parsed
		}
	}

	isProduction := apiURL == "https://wsp.posta.rs/api"

	// COD Configuration
	bankAccount := os.Getenv("POST_EXPRESS_BANK_ACCOUNT")
	if bankAccount == "" {
		// Тестовый банковский счёт для COD (замените на реальный в production!)
		bankAccount = "160-12345678-90" // default test bank account
	}
	paymentCode := os.Getenv("POST_EXPRESS_PAYMENT_CODE")
	if paymentCode == "" {
		paymentCode = "189" // default payment code
	}
	paymentModel := os.Getenv("POST_EXPRESS_PAYMENT_MODEL")
	if paymentModel == "" {
		paymentModel = "97" // default payment model
	}

	return &Config{
		APIURL:        apiURL,
		Username:      username,
		Password:      password,
		Brand:         brand,
		Warehouse:     warehouse,
		Timeout:       time.Duration(timeoutSeconds) * time.Second,
		RetryAttempts: retryAttempts,
		IsProduction:  isProduction,
		BankAccount:   bankAccount,
		PaymentCode:   paymentCode,
		PaymentModel:  paymentModel,
	}, nil
}
