package allsecure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient проверяет создание нового клиента
func TestNewClient(t *testing.T) {
	config := Config{
		BaseURL:  "https://api.allsecure.rs",
		Username: "test_user",
		Password: "test_pass",
		Timeout:  30,
	}

	client := NewClient(config)

	if client.BaseURL != config.BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", config.BaseURL, client.BaseURL)
	}
	if client.Username != config.Username {
		t.Errorf("Expected Username %s, got %s", config.Username, client.Username)
	}
	if client.Password != config.Password {
		t.Errorf("Expected Password %s, got %s", config.Password, client.Password)
	}
	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.HTTPClient.Timeout)
	}
}

// TestNewClientDefaultTimeout проверяет дефолтный timeout
func TestNewClientDefaultTimeout(t *testing.T) {
	config := Config{
		BaseURL:  "https://api.allsecure.rs",
		Username: "test_user",
		Password: "test_pass",
		Timeout:  0, // No timeout specified
	}

	client := NewClient(config)

	expectedTimeout := 30 * time.Second
	if client.HTTPClient.Timeout != expectedTimeout {
		t.Errorf("Expected default timeout %v, got %v", expectedTimeout, client.HTTPClient.Timeout)
	}
}

// TestBasicAuth проверяет создание Basic Authentication header
func TestBasicAuth(t *testing.T) {
	client := &Client{
		BaseURL:  "https://api.allsecure.rs",
		Username: "testuser",
		Password: "testpass",
	}

	auth := client.basicAuth()
	expected := "Basic dGVzdHVzZXI6dGVzdHBhc3M=" // base64("testuser:testpass")

	if auth != expected {
		t.Errorf("Expected auth header %s, got %s", expected, auth)
	}
}

// TestDebitSuccess проверяет успешный debit запрос
func TestDebitSuccess(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод и путь
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.URL.Path != "/api/v3/debit" {
			t.Errorf("Expected path /api/v3/debit, got %s", r.URL.Path)
		}

		// Проверяем заголовки
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json")
		}
		if r.Header.Get("Authorization") == "" {
			t.Errorf("Expected Authorization header")
		}

		// Отправляем успешный ответ
		response := TransactionResponse{
			Success:    true,
			UUID:       "test-uuid-123",
			PurchaseID: "purchase-123",
			Status:     "success",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Создаем клиента с mock сервером
	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
	}
	client := NewClient(config)

	// Создаем запрос
	req := TransactionRequest{
		Amount:       "100.00",
		Currency:     "EUR",
		Description:  "Test payment",
		MerchantTxID: "test-tx-123",
		Customer: Customer{
			Identification: "user-123",
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@example.com",
		},
	}

	// Выполняем запрос
	ctx := context.Background()
	response, err := client.Debit(ctx, req)
	// Проверяем результат
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("Response is nil")
	}
	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}
	if response.UUID != "test-uuid-123" {
		t.Errorf("Expected UUID test-uuid-123, got %s", response.UUID)
	}
}

// TestDebitAPIError проверяет обработку ошибки API
func TestDebitAPIError(t *testing.T) {
	// Mock HTTP server with error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := TransactionResponse{
			Success:      false,
			ErrorMessage: "Insufficient funds",
			ErrorCode:    "INSUFFICIENT_FUNDS",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
	}
	client := NewClient(config)

	req := TransactionRequest{
		Amount:       "1000.00",
		Currency:     "EUR",
		Description:  "Test payment",
		MerchantTxID: "test-tx-123",
		Customer: Customer{
			Identification: "user-123",
			Email:          "john@example.com",
		},
	}

	ctx := context.Background()
	response, err := client.Debit(ctx, req)

	// Должна быть ошибка
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if response == nil {
		t.Fatal("Response is nil")
	}
	if response.Success {
		t.Errorf("Expected success=false, got %v", response.Success)
	}
	if response.ErrorCode != "INSUFFICIENT_FUNDS" {
		t.Errorf("Expected error code INSUFFICIENT_FUNDS, got %s", response.ErrorCode)
	}
}

// TestPreauthorizeSuccess проверяет успешную preauthorization
func TestPreauthorizeSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/preauthorize" {
			t.Errorf("Expected path /api/v3/preauthorize, got %s", r.URL.Path)
		}

		response := TransactionResponse{
			Success:     true,
			UUID:        "preauth-uuid-123",
			Status:      "pending",
			RedirectURL: "https://payment.allsecure.rs/redirect/123",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
	}
	client := NewClient(config)

	req := TransactionRequest{
		Amount:       "50.00",
		Currency:     "EUR",
		Description:  "Test preauth",
		MerchantTxID: "preauth-tx-123",
		Customer: Customer{
			Identification: "user-456",
			Email:          "jane@example.com",
		},
		SuccessURL: "https://mysite.com/success",
		CancelURL:  "https://mysite.com/cancel",
	}

	ctx := context.Background()
	response, err := client.Preauthorize(ctx, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.UUID != "preauth-uuid-123" {
		t.Errorf("Expected UUID preauth-uuid-123, got %s", response.UUID)
	}
	if response.RedirectURL == "" {
		t.Error("Expected non-empty redirect URL")
	}
}

// TestCaptureSuccess проверяет успешный capture
func TestCaptureSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/capture" {
			t.Errorf("Expected path /api/v3/capture, got %s", r.URL.Path)
		}

		// Проверяем тело запроса
		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Logf("Failed to decode request body: %v", err)
		}

		if reqBody["uuid"] != "preauth-uuid-123" {
			t.Errorf("Expected uuid preauth-uuid-123, got %v", reqBody["uuid"])
		}
		if reqBody["amount"] != "50.00" {
			t.Errorf("Expected amount 50.00, got %v", reqBody["amount"])
		}

		response := TransactionResponse{
			Success: true,
			UUID:    "preauth-uuid-123",
			Status:  "success",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
	}
	client := NewClient(config)

	ctx := context.Background()
	response, err := client.Capture(ctx, "preauth-uuid-123", "50.00")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.Status != "success" {
		t.Errorf("Expected status success, got %s", response.Status)
	}
}

// TestRefundSuccess проверяет успешный refund
func TestRefundSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/refund" {
			t.Errorf("Expected path /api/v3/refund, got %s", r.URL.Path)
		}

		response := TransactionResponse{
			Success: true,
			UUID:    "refund-uuid-123",
			Status:  "success",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Logf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
	}
	client := NewClient(config)

	ctx := context.Background()
	response, err := client.Refund(ctx, "payment-uuid-123", "25.00")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if response.UUID != "refund-uuid-123" {
		t.Errorf("Expected UUID refund-uuid-123, got %s", response.UUID)
	}
}

// TestHTTPError проверяет обработку HTTP ошибок
func TestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Internal Server Error")); err != nil {
			t.Logf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
	}
	client := NewClient(config)

	req := TransactionRequest{
		Amount:       "100.00",
		Currency:     "EUR",
		Description:  "Test payment",
		MerchantTxID: "test-tx-123",
		Customer: Customer{
			Identification: "user-123",
			Email:          "test@example.com",
		},
	}

	ctx := context.Background()
	_, err := client.Debit(ctx, req)

	if err == nil {
		t.Error("Expected error for HTTP 500, got nil")
	}
}

// TestContextTimeout проверяет обработку timeout
func TestContextTimeout(t *testing.T) {
	// Создаем медленный сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Задержка больше timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "test",
		Timeout:  1, // 1 секунда timeout
	}
	client := NewClient(config)

	req := TransactionRequest{
		Amount:       "100.00",
		Currency:     "EUR",
		Description:  "Test payment",
		MerchantTxID: "test-tx-123",
		Customer: Customer{
			Identification: "user-123",
			Email:          "test@example.com",
		},
	}

	// Контекст с очень коротким timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.Debit(ctx, req)

	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
