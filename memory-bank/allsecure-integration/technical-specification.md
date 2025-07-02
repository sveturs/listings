# Техническая спецификация интеграции AllSecure

## Обзор интеграции

### Архитектурное решение
**Hybrid Integration**: Комбинация Server-to-Server API для контроля + SecurePay Widget для PCI compliance.

### Основные компоненты
1. **AllSecure API Client** (Go)
2. **Payment Service Layer** (Business Logic)
3. **SecurePay Widget** (Frontend)
4. **Webhook Handler** (Async notifications)
5. **Database Extensions** (Payment data)

## API Client Implementation

### Структура клиента

```go
// internal/pkg/allsecure/client.go
package allsecure

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/shopspring/decimal"
)

type Client struct {
    BaseURL    string
    Username   string
    Password   string
    HTTPClient *http.Client
    Logger     logger.Logger
}

type Config struct {
    BaseURL  string `json:"baseUrl"`
    Username string `json:"username"`
    Password string `json:"password"`
    Timeout  int    `json:"timeout"`
}

func NewClient(config Config) *Client {
    return &Client{
        BaseURL:  config.BaseURL,
        Username: config.Username,
        Password: config.Password,
        HTTPClient: &http.Client{
            Timeout: time.Duration(config.Timeout) * time.Second,
        },
    }
}
```

### Request/Response структуры

```go
// Transaction request
type TransactionRequest struct {
    Amount        string                 `json:"amount"`
    Currency      string                 `json:"currency"`
    Description   string                 `json:"description"`
    MerchantTxID  string                 `json:"merchantTransactionId"`
    Customer      Customer               `json:"customer"`
    CallbackURL   string                 `json:"callbackUrl,omitempty"`
    SuccessURL    string                 `json:"successUrl,omitempty"`
    CancelURL     string                 `json:"cancelUrl,omitempty"`
    ErrorURL      string                 `json:"errorUrl,omitempty"`
    ExtraData     map[string]interface{} `json:"extraData,omitempty"`
}

type Customer struct {
    Identification string `json:"identification"`
    FirstName      string `json:"firstName"`
    LastName       string `json:"lastName"`
    Email          string `json:"email"`
    IPAddress      string `json:"ipAddress,omitempty"`
    Company        string `json:"company,omitempty"`
    BillingAddress Address `json:"billingAddress,omitempty"`
}

type Address struct {
    FirstName   string `json:"firstName,omitempty"`
    LastName    string `json:"lastName,omitempty"`
    Company     string `json:"company,omitempty"`
    Street1     string `json:"street1,omitempty"`
    Street2     string `json:"street2,omitempty"`
    City        string `json:"city,omitempty"`
    Zip         string `json:"zip,omitempty"`
    State       string `json:"state,omitempty"`
    Country     string `json:"country,omitempty"`
    Phone       string `json:"phone,omitempty"`
    Email       string `json:"email,omitempty"`
}

// Transaction response
type TransactionResponse struct {
    Success         bool                   `json:"success"`
    UUID            string                 `json:"uuid"`
    PurchaseID      string                 `json:"purchaseId"`
    ReturnType      string                 `json:"returnType"`
    Status          string                 `json:"status"`
    ErrorMessage    string                 `json:"errorMessage,omitempty"`
    ErrorCode       string                 `json:"errorCode,omitempty"`
    AdapterMessage  string                 `json:"adapterMessage,omitempty"`
    AdapterCode     string                 `json:"adapterCode,omitempty"`
    RedirectURL     string                 `json:"redirectUrl,omitempty"`
    ExtraData       map[string]interface{} `json:"extraData,omitempty"`
}

// Webhook payload
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
```

### API Methods

```go
// Payment operations
func (c *Client) Debit(ctx context.Context, req TransactionRequest) (*TransactionResponse, error) {
    return c.makeRequest(ctx, "POST", "/debit", req)
}

func (c *Client) Preauthorize(ctx context.Context, req TransactionRequest) (*TransactionResponse, error) {
    return c.makeRequest(ctx, "POST", "/preauthorize", req)
}

func (c *Client) Capture(ctx context.Context, uuid string, amount string) (*TransactionResponse, error) {
    req := map[string]interface{}{
        "uuid":   uuid,
        "amount": amount,
    }
    return c.makeRequest(ctx, "POST", "/capture", req)
}

func (c *Client) Void(ctx context.Context, uuid string) (*TransactionResponse, error) {
    req := map[string]interface{}{
        "uuid": uuid,
    }
    return c.makeRequest(ctx, "POST", "/void", req)
}

func (c *Client) Refund(ctx context.Context, uuid string, amount string) (*TransactionResponse, error) {
    req := map[string]interface{}{
        "uuid":   uuid,
        "amount": amount,
    }
    return c.makeRequest(ctx, "POST", "/refund", req)
}

// Token operations for recurring payments
func (c *Client) Register(ctx context.Context, req TransactionRequest) (*TransactionResponse, error) {
    return c.makeRequest(ctx, "POST", "/register", req)
}

func (c *Client) Deregister(ctx context.Context, uuid string) (*TransactionResponse, error) {
    req := map[string]interface{}{
        "uuid": uuid,
    }
    return c.makeRequest(ctx, "POST", "/deregister", req)
}

// Payout operations
func (c *Client) Payout(ctx context.Context, req PayoutRequest) (*TransactionResponse, error) {
    return c.makeRequest(ctx, "POST", "/payout", req)
}
```

### HTTP Client Implementation

```go
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
    
    // Log request
    c.Logger.Debug("AllSecure API Request", 
        "method", method,
        "url", url,
        "body", string(body))
    
    // Make request
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()
    
    // Read response
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // Log response
    c.Logger.Debug("AllSecure API Response",
        "status", resp.StatusCode,
        "body", string(respBody))
    
    // Parse response
    var result TransactionResponse
    if err := json.Unmarshal(respBody, &result); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }
    
    return &result, nil
}

func (c *Client) basicAuth() string {
    credentials := c.Username + ":" + c.Password
    return "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
}
```

## Service Layer Implementation

### Payment Service

```go
// internal/proj/payments/service/payment_service.go
package service

type PaymentService struct {
    allsecureClient   *allsecure.Client
    paymentRepo       *repository.PaymentRepository
    userRepo          *repository.UserRepository
    listingRepo       *repository.MarketplaceRepository
    config            *config.PaymentConfig
    logger            logger.Logger
}

type CreatePaymentRequest struct {
    UserID      int             `json:"user_id" validate:"required"`
    ListingID   int             `json:"listing_id" validate:"required"`
    Amount      decimal.Decimal `json:"amount" validate:"required,gt=0"`
    Currency    string          `json:"currency" validate:"required,len=3"`
    Description string          `json:"description"`
    ReturnURL   string          `json:"return_url" validate:"required,url"`
}

type PaymentResult struct {
    TransactionID   int64  `json:"transaction_id"`
    GatewayUUID     string `json:"gateway_uuid"`
    Status          string `json:"status"`
    RedirectURL     string `json:"redirect_url,omitempty"`
    RequiresAction  bool   `json:"requires_action"`
}

func (s *PaymentService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResult, error) {
    // 1. Validate request
    if err := s.validatePaymentRequest(ctx, req); err != nil {
        return nil, err
    }
    
    // 2. Get user and listing
    user, err := s.userRepo.GetByID(ctx, req.UserID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    listing, err := s.listingRepo.GetByID(ctx, req.ListingID)
    if err != nil {
        return nil, fmt.Errorf("failed to get listing: %w", err)
    }
    
    // 3. Calculate amounts
    marketplaceCommission := s.calculateCommission(req.Amount)
    sellerAmount := req.Amount.Sub(marketplaceCommission)
    
    // 4. Create transaction record
    transaction, err := s.paymentRepo.CreateTransaction(ctx, repository.CreateTransactionRequest{
        UserID:                req.UserID,
        ListingID:             req.ListingID,
        Amount:                req.Amount,
        Currency:              req.Currency,
        MarketplaceCommission: marketplaceCommission,
        SellerAmount:          sellerAmount,
        Description:           req.Description,
        Status:                "pending",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create transaction: %w", err)
    }
    
    // 5. Create AllSecure transaction
    allsecureReq := allsecure.TransactionRequest{
        Amount:       req.Amount.String(),
        Currency:     req.Currency,
        Description:  req.Description,
        MerchantTxID: fmt.Sprintf("SVT-%d", transaction.ID),
        Customer: allsecure.Customer{
            Identification: fmt.Sprintf("user-%d", user.ID),
            FirstName:      user.Name,
            Email:          user.Email,
        },
        CallbackURL: s.config.WebhookURL,
        SuccessURL:  req.ReturnURL + "?status=success",
        CancelURL:   req.ReturnURL + "?status=cancelled",
        ErrorURL:    req.ReturnURL + "?status=error",
    }
    
    response, err := s.allsecureClient.Preauthorize(ctx, allsecureReq)
    if err != nil {
        // Update transaction status
        s.paymentRepo.UpdateTransactionStatus(ctx, transaction.ID, "failed", map[string]interface{}{
            "error": err.Error(),
        })
        return nil, fmt.Errorf("AllSecure request failed: %w", err)
    }
    
    // 6. Update transaction with gateway data
    err = s.paymentRepo.UpdateTransaction(ctx, transaction.ID, repository.UpdateTransactionRequest{
        GatewayTransactionID: response.UUID,
        GatewayReferenceID:   response.PurchaseID,
        Status:               s.mapAllSecureStatus(response.Status),
        GatewayResponse:      response,
    })
    if err != nil {
        s.logger.Error("Failed to update transaction", "error", err, "transactionID", transaction.ID)
    }
    
    return &PaymentResult{
        TransactionID:  transaction.ID,
        GatewayUUID:    response.UUID,
        Status:         s.mapAllSecureStatus(response.Status),
        RedirectURL:    response.RedirectURL,
        RequiresAction: response.ReturnType == "REDIRECT",
    }, nil
}

func (s *PaymentService) ProcessWebhook(ctx context.Context, payload []byte) error {
    var webhook allsecure.WebhookPayload
    if err := json.Unmarshal(payload, &webhook); err != nil {
        return fmt.Errorf("failed to unmarshal webhook: %w", err)
    }
    
    // Find transaction by merchant transaction ID
    merchantTxID := webhook.MerchantTxID
    if !strings.HasPrefix(merchantTxID, "SVT-") {
        return fmt.Errorf("invalid merchant transaction ID format: %s", merchantTxID)
    }
    
    transactionIDStr := strings.TrimPrefix(merchantTxID, "SVT-")
    transactionID, err := strconv.ParseInt(transactionIDStr, 10, 64)
    if err != nil {
        return fmt.Errorf("failed to parse transaction ID: %w", err)
    }
    
    transaction, err := s.paymentRepo.GetByID(ctx, transactionID)
    if err != nil {
        return fmt.Errorf("failed to get transaction: %w", err)
    }
    
    // Update transaction status
    newStatus := s.mapAllSecureStatus(webhook.Status)
    err = s.paymentRepo.UpdateTransactionStatus(ctx, transactionID, newStatus, map[string]interface{}{
        "webhook_data": webhook,
        "updated_at":   time.Now(),
    })
    if err != nil {
        return fmt.Errorf("failed to update transaction status: %w", err)
    }
    
    // Handle status-specific logic
    switch newStatus {
    case "authorized":
        // Create escrow payment
        err = s.createEscrowPayment(ctx, transaction)
        if err != nil {
            s.logger.Error("Failed to create escrow payment", "error", err, "transactionID", transactionID)
        }
        
    case "captured":
        // Release escrow to seller
        err = s.releaseEscrowPayment(ctx, transaction)
        if err != nil {
            s.logger.Error("Failed to release escrow payment", "error", err, "transactionID", transactionID)
        }
        
    case "failed":
        // Handle failed payment
        s.handleFailedPayment(ctx, transaction)
    }
    
    return nil
}

func (s *PaymentService) mapAllSecureStatus(allsecureStatus string) string {
    switch allsecureStatus {
    case "FINISHED":
        return "captured"
    case "PENDING":
        return "authorized"
    case "ERROR":
        return "failed"
    default:
        return "unknown"
    }
}
```

## Database Models

### Go Models

```go
// internal/domain/models/payment.go
package models

type PaymentGateway struct {
    ID        int                    `json:"id" db:"id"`
    Name      string                 `json:"name" db:"name"`
    IsActive  bool                   `json:"is_active" db:"is_active"`
    Config    map[string]interface{} `json:"config" db:"config"`
    CreatedAt time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

type PaymentTransaction struct {
    ID                    int64                  `json:"id" db:"id"`
    GatewayID             int                    `json:"gateway_id" db:"gateway_id"`
    UserID                int                    `json:"user_id" db:"user_id"`
    ListingID             *int                   `json:"listing_id,omitempty" db:"listing_id"`
    OrderReference        string                 `json:"order_reference" db:"order_reference"`
    GatewayTransactionID  *string                `json:"gateway_transaction_id,omitempty" db:"gateway_transaction_id"`
    GatewayReferenceID    *string                `json:"gateway_reference_id,omitempty" db:"gateway_reference_id"`
    Amount                decimal.Decimal        `json:"amount" db:"amount"`
    Currency              string                 `json:"currency" db:"currency"`
    MarketplaceCommission *decimal.Decimal       `json:"marketplace_commission,omitempty" db:"marketplace_commission"`
    SellerAmount          *decimal.Decimal       `json:"seller_amount,omitempty" db:"seller_amount"`
    Status                string                 `json:"status" db:"status"`
    GatewayStatus         *string                `json:"gateway_status,omitempty" db:"gateway_status"`
    PaymentMethod         *string                `json:"payment_method,omitempty" db:"payment_method"`
    CustomerEmail         *string                `json:"customer_email,omitempty" db:"customer_email"`
    Description           *string                `json:"description,omitempty" db:"description"`
    GatewayResponse       map[string]interface{} `json:"gateway_response,omitempty" db:"gateway_response"`
    ErrorDetails          map[string]interface{} `json:"error_details,omitempty" db:"error_details"`
    CreatedAt             time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
    AuthorizedAt          *time.Time             `json:"authorized_at,omitempty" db:"authorized_at"`
    CapturedAt            *time.Time             `json:"captured_at,omitempty" db:"captured_at"`
    FailedAt              *time.Time             `json:"failed_at,omitempty" db:"failed_at"`
    
    // Relations
    Gateway *PaymentGateway     `json:"gateway,omitempty"`
    User    *User               `json:"user,omitempty"`
    Listing *MarketplaceListing `json:"listing,omitempty"`
}

type EscrowPayment struct {
    ID                   int64           `json:"id" db:"id"`
    PaymentTransactionID int64           `json:"payment_transaction_id" db:"payment_transaction_id"`
    SellerID             int             `json:"seller_id" db:"seller_id"`
    BuyerID              int             `json:"buyer_id" db:"buyer_id"`
    ListingID            int             `json:"listing_id" db:"listing_id"`
    Amount               decimal.Decimal `json:"amount" db:"amount"`
    MarketplaceCommission decimal.Decimal `json:"marketplace_commission" db:"marketplace_commission"`
    SellerAmount         decimal.Decimal `json:"seller_amount" db:"seller_amount"`
    Status               string          `json:"status" db:"status"`
    ReleaseDate          *time.Time      `json:"release_date,omitempty" db:"release_date"`
    CreatedAt            time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
    
    // Relations
    PaymentTransaction *PaymentTransaction `json:"payment_transaction,omitempty"`
    Seller             *User               `json:"seller,omitempty"`
    Buyer              *User               `json:"buyer,omitempty"`
    Listing            *MarketplaceListing `json:"listing,omitempty"`
}
```

## Frontend Integration

### SecurePay Widget

```typescript
// src/lib/allsecure/widget.ts
declare global {
  interface Window {
    SecurePay: {
      init: (config: SecurePayConfig) => void;
      destroy: () => void;
    };
  }
}

interface SecurePayConfig {
  containerSelector: string;
  amount: number;
  currency: string;
  language?: string;
  theme?: 'light' | 'dark';
  onComplete: (result: PaymentResult) => void;
  onError: (error: PaymentError) => void;
  onCancel?: () => void;
}

interface PaymentResult {
  success: boolean;
  transactionId: string;
  amount: number;
  currency: string;
}

interface PaymentError {
  code: string;
  message: string;
  details?: any;
}

export class AllSecureWidget {
  private isLoaded = false;
  private container: HTMLElement | null = null;

  async load(): Promise<void> {
    if (this.isLoaded) return;

    return new Promise((resolve, reject) => {
      const script = document.createElement('script');
      script.src = process.env.NEXT_PUBLIC_ALLSECURE_WIDGET_URL || 
                   'https://securepay.allsecure.rs/js/widget.js';
      script.onload = () => {
        this.isLoaded = true;
        resolve();
      };
      script.onerror = () => {
        reject(new Error('Failed to load AllSecure widget'));
      };
      document.head.appendChild(script);
    });
  }

  init(config: SecurePayConfig): void {
    if (!this.isLoaded) {
      throw new Error('Widget not loaded. Call load() first.');
    }

    this.container = document.querySelector(config.containerSelector);
    if (!this.container) {
      throw new Error(`Container not found: ${config.containerSelector}`);
    }

    window.SecurePay.init(config);
  }

  destroy(): void {
    if (window.SecurePay && window.SecurePay.destroy) {
      window.SecurePay.destroy();
    }
    if (this.container) {
      this.container.innerHTML = '';
    }
  }
}
```

### React Hook

```typescript
// src/hooks/useAllSecurePayment.ts
import { useState, useCallback } from 'react';
import { AllSecureWidget } from '@/lib/allsecure/widget';
import { apiClient } from '@/lib/api-client';

interface UseAllSecurePaymentProps {
  onSuccess?: (result: PaymentResult) => void;
  onError?: (error: PaymentError) => void;
}

export const useAllSecurePayment = ({
  onSuccess,
  onError
}: UseAllSecurePaymentProps = {}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [widget] = useState(() => new AllSecureWidget());

  const initializePayment = useCallback(async (
    listingId: number,
    amount: number,
    currency: string = 'RSD'
  ) => {
    try {
      setIsLoading(true);
      setError(null);

      // Create payment on backend
      const response = await apiClient.post('/payments/create', {
        listing_id: listingId,
        amount: amount.toString(),
        currency,
        return_url: window.location.href
      });

      const { transaction_id, gateway_uuid, redirect_url, requires_action } = response.data;

      if (requires_action && redirect_url) {
        // Redirect to AllSecure hosted page
        window.location.href = redirect_url;
      } else {
        // Use widget for direct payment
        await widget.load();
        
        widget.init({
          containerSelector: '#allsecure-widget',
          amount,
          currency,
          onComplete: (result) => {
            onSuccess?.(result);
          },
          onError: (error) => {
            setError(error.message);
            onError?.(error);
          }
        });
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Payment initialization failed';
      setError(message);
      onError?.({ code: 'INIT_ERROR', message });
    } finally {
      setIsLoading(false);
    }
  }, [widget, onSuccess, onError]);

  const cleanup = useCallback(() => {
    widget.destroy();
  }, [widget]);

  return {
    initializePayment,
    cleanup,
    isLoading,
    error
  };
};
```

## Configuration

### Environment Variables

```bash
# AllSecure Configuration
ALLSECURE_BASE_URL=https://asxgw.com
ALLSECURE_USERNAME=your_username
ALLSECURE_PASSWORD=your_password
ALLSECURE_WEBHOOK_URL=https://svetu.rs/api/v1/webhooks/allsecure
ALLSECURE_WEBHOOK_SECRET=your_webhook_secret

# Frontend
NEXT_PUBLIC_ALLSECURE_WIDGET_URL=https://securepay.allsecure.rs/js/widget.js
```

### Go Configuration

```go
// internal/config/payment.go
type PaymentConfig struct {
    AllSecure AllSecureConfig `mapstructure:"allsecure"`
}

type AllSecureConfig struct {
    BaseURL       string `mapstructure:"base_url"`
    Username      string `mapstructure:"username"`
    Password      string `mapstructure:"password"`
    WebhookURL    string `mapstructure:"webhook_url"`
    WebhookSecret string `mapstructure:"webhook_secret"`
    Timeout       int    `mapstructure:"timeout"`
}
```

## Error Handling

### Error Types

```go
type PaymentError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

const (
    ErrInvalidAmount        = "INVALID_AMOUNT"
    ErrInsufficientFunds   = "INSUFFICIENT_FUNDS"
    ErrPaymentDeclined     = "PAYMENT_DECLINED"
    ErrGatewayError        = "GATEWAY_ERROR"
    ErrTransactionNotFound = "TRANSACTION_NOT_FOUND"
    ErrInvalidSignature    = "INVALID_SIGNATURE"
)
```

### Retry Logic

```go
func (c *Client) makeRequestWithRetry(ctx context.Context, method, endpoint string, payload interface{}) (*TransactionResponse, error) {
    var lastErr error
    
    for attempt := 0; attempt < 3; attempt++ {
        response, err := c.makeRequest(ctx, method, endpoint, payload)
        if err == nil {
            return response, nil
        }
        
        lastErr = err
        
        // Don't retry on client errors
        if isClientError(err) {
            break
        }
        
        // Exponential backoff
        backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(backoff):
            continue
        }
    }
    
    return nil, lastErr
}
```

## Security Implementation

### Webhook Signature Verification

```go
func (h *WebhookHandler) verifySignature(payload []byte, signature string) bool {
    expectedSignature := h.calculateSignature(payload)
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (h *WebhookHandler) calculateSignature(payload []byte) string {
    mac := hmac.New(sha256.New, []byte(h.webhookSecret))
    mac.Write(payload)
    return hex.EncodeToString(mac.Sum(nil))
}
```

### Request Logging

```go
func (c *Client) logRequest(req *http.Request, body []byte) {
    c.Logger.Info("AllSecure API Request",
        "method", req.Method,
        "url", req.URL.String(),
        "headers", sanitizeHeaders(req.Header),
        "body_size", len(body))
}

func sanitizeHeaders(headers http.Header) map[string]string {
    sanitized := make(map[string]string)
    for key, values := range headers {
        if strings.ToLower(key) == "authorization" {
            sanitized[key] = "[REDACTED]"
        } else {
            sanitized[key] = strings.Join(values, ", ")
        }
    }
    return sanitized
}
```

Эта техническая спецификация предоставляет полную roadmap для имплементации интеграции с AllSecure, включая все необходимые компоненты, структуры данных и error handling.