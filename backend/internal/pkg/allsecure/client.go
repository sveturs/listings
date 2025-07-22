package allsecure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client представляет AllSecure API клиент
type Client struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// Config содержит конфигурацию для AllSecure клиента
type Config struct {
	BaseURL  string `json:"baseUrl"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
}

// NewClient создает новый AllSecure клиент
func NewClient(config Config) *Client {
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		BaseURL:  config.BaseURL,
		Username: config.Username,
		Password: config.Password,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// TransactionRequest представляет запрос на создание транзакции
type TransactionRequest struct {
	Amount       string                 `json:"amount"`
	Currency     string                 `json:"currency"`
	Description  string                 `json:"description"`
	MerchantTxID string                 `json:"merchantTransactionId"`
	Customer     Customer               `json:"customer"`
	CallbackURL  string                 `json:"callbackUrl,omitempty"`
	SuccessURL   string                 `json:"successUrl,omitempty"`
	CancelURL    string                 `json:"cancelUrl,omitempty"`
	ErrorURL     string                 `json:"errorUrl,omitempty"`
	ExtraData    map[string]interface{} `json:"extraData,omitempty"`
}

// Customer представляет информацию о клиенте
type Customer struct {
	Identification string  `json:"identification"`
	FirstName      string  `json:"firstName"`
	LastName       string  `json:"lastName"`
	Email          string  `json:"email"`
	IPAddress      string  `json:"ipAddress,omitempty"`
	Company        string  `json:"company,omitempty"`
	BillingAddress Address `json:"billingAddress,omitempty"`
}

// Address представляет адрес
type Address struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Company   string `json:"company,omitempty"`
	Street1   string `json:"street1,omitempty"`
	Street2   string `json:"street2,omitempty"`
	City      string `json:"city,omitempty"`
	Zip       string `json:"zip,omitempty"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
}

// TransactionResponse представляет ответ от AllSecure API
type TransactionResponse struct {
	Success        bool                   `json:"success"`
	UUID           string                 `json:"uuid"`
	PurchaseID     string                 `json:"purchaseId"`
	ReturnType     string                 `json:"returnType"`
	Status         string                 `json:"status"`
	ErrorMessage   string                 `json:"errorMessage,omitempty"`
	ErrorCode      string                 `json:"errorCode,omitempty"`
	AdapterMessage string                 `json:"adapterMessage,omitempty"`
	AdapterCode    string                 `json:"adapterCode,omitempty"`
	RedirectURL    string                 `json:"redirectUrl,omitempty"`
	ExtraData      map[string]interface{} `json:"extraData,omitempty"`
}

// WebhookPayload представляет payload webhook'а от AllSecure
type WebhookPayload struct {
	UUID           string                 `json:"uuid"`
	MerchantTxID   string                 `json:"merchantTransactionId"`
	Status         string                 `json:"status"`
	ErrorMessage   string                 `json:"errorMessage,omitempty"`
	ErrorCode      string                 `json:"errorCode,omitempty"`
	AdapterMessage string                 `json:"adapterMessage,omitempty"`
	AdapterCode    string                 `json:"adapterCode,omitempty"`
	ExtraData      map[string]interface{} `json:"extraData,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
}

// PayoutRequest представляет запрос на выплату
type PayoutRequest struct {
	Amount      string                 `json:"amount"`
	Currency    string                 `json:"currency"`
	Description string                 `json:"description"`
	Customer    Customer               `json:"customer"`
	ExtraData   map[string]interface{} `json:"extraData,omitempty"`
}

// Debit выполняет прямое списание средств
func (c *Client) Debit(ctx context.Context, req TransactionRequest) (*TransactionResponse, error) {
	return c.makeRequest(ctx, "POST", "/debit", req)
}

// Preauthorize резервирует средства
func (c *Client) Preauthorize(ctx context.Context, req TransactionRequest) (*TransactionResponse, error) {
	return c.makeRequest(ctx, "POST", "/preauthorize", req)
}

// Capture списывает зарезервированные средства
func (c *Client) Capture(ctx context.Context, uuid string, amount string) (*TransactionResponse, error) {
	req := map[string]interface{}{
		"uuid":   uuid,
		"amount": amount,
	}
	return c.makeRequest(ctx, "POST", "/capture", req)
}

// Void отменяет транзакцию
func (c *Client) Void(ctx context.Context, uuid string) (*TransactionResponse, error) {
	req := map[string]interface{}{
		"uuid": uuid,
	}
	return c.makeRequest(ctx, "POST", "/void", req)
}

// Refund возвращает средства
func (c *Client) Refund(ctx context.Context, uuid string, amount string) (*TransactionResponse, error) {
	req := map[string]interface{}{
		"uuid":   uuid,
		"amount": amount,
	}
	return c.makeRequest(ctx, "POST", "/refund", req)
}

// Register регистрирует карту для будущих платежей
func (c *Client) Register(ctx context.Context, req TransactionRequest) (*TransactionResponse, error) {
	return c.makeRequest(ctx, "POST", "/register", req)
}

// Deregister удаляет зарегистрированную карту
func (c *Client) Deregister(ctx context.Context, uuid string) (*TransactionResponse, error) {
	req := map[string]interface{}{
		"uuid": uuid,
	}
	return c.makeRequest(ctx, "POST", "/deregister", req)
}

// Payout выполняет выплату средств
func (c *Client) Payout(ctx context.Context, req PayoutRequest) (*TransactionResponse, error) {
	return c.makeRequest(ctx, "POST", "/payout", req)
}

// makeRequest выполняет HTTP запрос к AllSecure API
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) (*TransactionResponse, error) {
	url := c.BaseURL + "/api/v3" + endpoint

	// Marshal payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.basicAuth())

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log error but don't return it since we're in defer
		}
	}()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var result TransactionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API errors
	if !result.Success {
		return &result, fmt.Errorf("AllSecure API error: %s (code: %s)", result.ErrorMessage, result.ErrorCode)
	}

	return &result, nil
}

// basicAuth создает Basic Authentication заголовок
func (c *Client) basicAuth() string {
	credentials := c.Username + ":" + c.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}
