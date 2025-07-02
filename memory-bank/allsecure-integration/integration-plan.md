# План интеграции AllSecure с платформой Sve Tu

## Обзор интеграции

### Цель интеграции
Интегрировать AllSecure payment gateway в маркетплейс Sve Tu для обеспечения безопасных платежей между покупателями и продавцами, включая:
- Платежи за товары и услуги
- Marketplace комиссии
- Escrow функциональность
- Выплаты продавцам

### Выбранная стратегия интеграции
**Server-to-Server REST API** + **SecurePay Widget** для максимальной гибкости и безопасности.

## Архитектура интеграции

### Компоненты системы
```
Frontend (React/Next.js)
    │
    ├── SecurePay Widget (SAQ-A compliant)
    │   └── Tokenization в AllSecure Card Vault
    │
    ▼
Backend (Go)
    │
    ├── Payment Service Layer
    │   ├── AllSecure API Client
    │   ├── Payment Processing
    │   └── Webhook Handler
    │
    ├── Existing Services
    │   ├── User Balance Service
    │   ├── Transaction Service  
    │   └── Marketplace Service
    │
    ▼
Database (PostgreSQL)
    │
    ├── payment_gateways
    ├── payment_transactions  
    ├── escrow_payments
    └── merchant_payouts
```

### Интеграционные точки

#### 1. Frontend Integration
- **SecurePay Widget** для сбора карточных данных
- **Payment confirmation flows** в UI
- **Receipt и status pages**

#### 2. Backend Integration  
- **AllSecure API Client** в Go
- **Payment processing workflows**
- **Webhook endpoints** для асинхронных уведомлений
- **Escrow management**

#### 3. Database Extensions
- **Новые таблицы** для платежных данных
- **Расширение существующих** таблиц транзакций

## Детальный план реализации

### Phase 1: Подготовка и настройка (1-2 недели)

#### 1.1 Бизнес подготовка
- [ ] **Связаться с AllSecure** для получения:
  - Demo/Sandbox credentials
  - Коммерческого предложения  
  - Технической консультации
  - Documentation access

- [ ] **Юридическая подготовка**:
  - Подготовка документов для регистрации мерчанта
  - KYC процедуры
  - Банковские реквизиты в Сербии
  - PCI DSS compliance план

#### 1.2 Техническая подготовка
- [ ] **SSL сертификат** для production домена
- [ ] **Webhook endpoint** защищенный URL
- [ ] **Мониторинг системы** для платежей
- [ ] **Backup стратегия** для критичных данных

### Phase 2: Database Schema (3-5 дней)

#### 2.1 Новые таблицы

```sql
-- Конфигурация платежных шлюзов
CREATE TABLE payment_gateways (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, -- 'allsecure'
    is_active BOOLEAN DEFAULT true,
    config JSONB NOT NULL, -- API credentials, endpoints
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Платежные транзакции
CREATE TABLE payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    gateway_id INT REFERENCES payment_gateways(id),
    user_id INT REFERENCES users(id),
    
    -- Ссылка на покупку
    listing_id INT REFERENCES marketplace_listings(id),
    order_reference VARCHAR(255) UNIQUE,
    
    -- AllSecure данные
    gateway_transaction_id VARCHAR(255),
    gateway_reference_id VARCHAR(255),
    
    -- Финансовые данные
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    marketplace_commission DECIMAL(12,2),
    seller_amount DECIMAL(12,2),
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'pending', -- pending, authorized, captured, failed, refunded
    gateway_status VARCHAR(50),
    
    -- Дополнительная информация
    payment_method VARCHAR(50), -- card, bank_transfer, etc.
    customer_email VARCHAR(255),
    description TEXT,
    
    -- Metadata
    gateway_response JSONB,
    error_details JSONB,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    authorized_at TIMESTAMP WITH TIME ZONE,
    captured_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE
);

-- Escrow платежи для marketplace
CREATE TABLE escrow_payments (
    id BIGSERIAL PRIMARY KEY,
    payment_transaction_id BIGINT REFERENCES payment_transactions(id),
    seller_id INT REFERENCES users(id),
    buyer_id INT REFERENCES users(id),
    listing_id INT REFERENCES marketplace_listings(id),
    
    amount DECIMAL(12,2) NOT NULL,
    marketplace_commission DECIMAL(12,2),
    seller_amount DECIMAL(12,2),
    
    status VARCHAR(50) DEFAULT 'held', -- held, released, refunded
    release_date TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Выплаты продавцам
CREATE TABLE merchant_payouts (
    id BIGSERIAL PRIMARY KEY,
    seller_id INT REFERENCES users(id),
    gateway_id INT REFERENCES payment_gateways(id),
    
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    
    -- AllSecure payout данные
    gateway_payout_id VARCHAR(255),
    gateway_reference_id VARCHAR(255),
    
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    
    -- Банковские данные получателя
    bank_account_info JSONB,
    
    -- Metadata
    gateway_response JSONB,
    error_details JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE
);

-- Индексы для производительности
CREATE INDEX idx_payment_transactions_user_id ON payment_transactions(user_id);
CREATE INDEX idx_payment_transactions_listing_id ON payment_transactions(listing_id);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX idx_payment_transactions_gateway_transaction_id ON payment_transactions(gateway_transaction_id);
CREATE INDEX idx_escrow_payments_seller_id ON escrow_payments(seller_id);
CREATE INDEX idx_escrow_payments_status ON escrow_payments(status);
CREATE INDEX idx_merchant_payouts_seller_id ON merchant_payouts(seller_id);
```

#### 2.2 Расширение существующих таблиц

```sql
-- Добавить поддержку AllSecure в user_transactions
ALTER TABLE user_transactions ADD COLUMN payment_gateway VARCHAR(50);
ALTER TABLE user_transactions ADD COLUMN gateway_transaction_id VARCHAR(255);
ALTER TABLE user_transactions ADD COLUMN gateway_reference VARCHAR(255);

-- Индексы
CREATE INDEX idx_user_transactions_gateway_transaction_id ON user_transactions(gateway_transaction_id);
```

### Phase 3: Backend Implementation (1-2 недели)

#### 3.1 AllSecure API Client

```go
// internal/pkg/allsecure/client.go
package allsecure

type Client struct {
    BaseURL     string
    Username    string
    Password    string
    HTTPClient  *http.Client
}

type Transaction struct {
    Amount      decimal.Decimal `json:"amount"`
    Currency    string         `json:"currency"`
    Description string         `json:"description"`
    Customer    Customer       `json:"customer"`
    ExtraData   map[string]interface{} `json:"extraData,omitempty"`
}

type Customer struct {
    Identification string `json:"identification"`
    FirstName      string `json:"firstName"`
    LastName       string `json:"lastName"`
    Email          string `json:"email"`
}

type TransactionResult struct {
    Success          bool   `json:"success"`
    UUID             string `json:"uuid"`
    PurchaseID       string `json:"purchaseId"`
    ReturnType       string `json:"returnType"`
    Status           string `json:"status"`
    ErrorMessage     string `json:"errorMessage,omitempty"`
    ErrorCode        string `json:"errorCode,omitempty"`
    RedirectURL      string `json:"redirectUrl,omitempty"`
}

// Методы
func (c *Client) Debit(ctx context.Context, req Transaction) (*TransactionResult, error)
func (c *Client) Preauthorize(ctx context.Context, req Transaction) (*TransactionResult, error)
func (c *Client) Capture(ctx context.Context, uuid string, amount decimal.Decimal) (*TransactionResult, error)
func (c *Client) Void(ctx context.Context, uuid string) (*TransactionResult, error)
func (c *Client) Refund(ctx context.Context, uuid string, amount decimal.Decimal) (*TransactionResult, error)
func (c *Client) Payout(ctx context.Context, req PayoutRequest) (*TransactionResult, error)
```

#### 3.2 Payment Service Layer

```go
// internal/proj/payments/service/allsecure_service.go
package service

type AllSecureService struct {
    client     *allsecure.Client
    repository *repository.PaymentRepository
    config     *config.AllSecureConfig
}

// Основные методы
func (s *AllSecureService) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*PaymentResult, error)
func (s *AllSecureService) CapturePayment(ctx context.Context, transactionID int64) error
func (s *AllSecureService) RefundPayment(ctx context.Context, transactionID int64, amount decimal.Decimal) error
func (s *AllSecureService) ProcessWebhook(ctx context.Context, payload []byte) error
func (s *AllSecureService) CreatePayout(ctx context.Context, req PayoutRequest) error
```

#### 3.3 Repository Layer

```go
// internal/proj/payments/repository/payment_repository.go
package repository

type PaymentRepository struct {
    db *sql.DB
}

func (r *PaymentRepository) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*PaymentTransaction, error)
func (r *PaymentRepository) UpdateTransactionStatus(ctx context.Context, id int64, status string, gatewayResponse map[string]interface{}) error
func (r *PaymentRepository) GetTransactionByGatewayID(ctx context.Context, gatewayID string) (*PaymentTransaction, error)
func (r *PaymentRepository) CreateEscrowPayment(ctx context.Context, req CreateEscrowRequest) (*EscrowPayment, error)
func (r *PaymentRepository) ReleaseEscrow(ctx context.Context, escrowID int64) error
```

#### 3.4 Webhook Handler

```go
// internal/proj/payments/handler/webhook_handler.go
package handler

// @Summary AllSecure Webhook Handler
// @Description Обработка уведомлений от AllSecure
// @Tags payments
// @Accept json
// @Param payload body AllSecureWebhookPayload true "Webhook payload"
// @Success 200 {object} utils.SuccessResponseSwag "OK"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request"
// @Router /webhooks/allsecure [post]
func (h *WebhookHandler) HandleAllSecureWebhook(c *fiber.Ctx) error {
    // Проверка подписи
    signature := c.Get("X-Signature")
    if !h.verifySignature(c.Body(), signature) {
        return utils.ErrorResponse(c, 400, "invalid signature")
    }
    
    // Обработка webhook
    err := h.service.ProcessWebhook(c.Context(), c.Body())
    if err != nil {
        logger.Error("Failed to process AllSecure webhook", "error", err)
        return utils.ErrorResponse(c, 500, "webhook.processError")
    }
    
    return utils.SuccessResponse(c, "OK")
}
```

### Phase 4: Frontend Implementation (1 неделя)

#### 4.1 SecurePay Widget Integration

```typescript
// src/components/payments/AllSecureWidget.tsx
import React, { useEffect, useState } from 'react';

interface AllSecureWidgetProps {
  amount: number;
  currency: string;
  onSuccess: (result: PaymentResult) => void;
  onError: (error: PaymentError) => void;
}

export const AllSecureWidget: React.FC<AllSecureWidgetProps> = ({
  amount,
  currency,
  onSuccess,
  onError
}) => {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Загрузка AllSecure SecurePay script
    const script = document.createElement('script');
    script.src = 'https://securepay.allsecure.rs/js/widget.js';
    script.onload = () => {
      initializeWidget();
    };
    document.head.appendChild(script);

    return () => {
      document.head.removeChild(script);
    };
  }, []);

  const initializeWidget = () => {
    (window as any).SecurePay.init({
      containerSelector: '#allsecure-widget',
      amount: amount,
      currency: currency,
      onComplete: (result: any) => {
        onSuccess(result);
      },
      onError: (error: any) => {
        onError(error);
      }
    });
    setIsLoading(false);
  };

  return (
    <div className="allsecure-payment-widget">
      {isLoading && (
        <div className="loading">
          <span className="loading loading-spinner loading-lg"></span>
          Загрузка платежной формы...
        </div>
      )}
      <div id="allsecure-widget"></div>
    </div>
  );
};
```

#### 4.2 Payment Flow Components

```typescript
// src/components/payments/PaymentFlow.tsx
export const PaymentFlow: React.FC<PaymentFlowProps> = ({ 
  listing, 
  onComplete 
}) => {
  const [currentStep, setCurrentStep] = useState<PaymentStep>('summary');
  const [paymentData, setPaymentData] = useState<PaymentData | null>(null);

  const handlePaymentSuccess = async (result: PaymentResult) => {
    try {
      // Подтверждение платежа на backend
      const response = await apiClient.confirmPayment({
        transactionId: result.transactionId,
        listingId: listing.id
      });
      
      setCurrentStep('success');
      onComplete(response.data);
    } catch (error) {
      setCurrentStep('error');
    }
  };

  return (
    <div className="payment-flow">
      {currentStep === 'summary' && (
        <PaymentSummary 
          listing={listing}
          onProceed={() => setCurrentStep('payment')}
        />
      )}
      
      {currentStep === 'payment' && (
        <AllSecureWidget
          amount={listing.price}
          currency={listing.currency}
          onSuccess={handlePaymentSuccess}
          onError={() => setCurrentStep('error')}
        />
      )}
      
      {currentStep === 'success' && (
        <PaymentSuccess paymentData={paymentData} />
      )}
      
      {currentStep === 'error' && (
        <PaymentError onRetry={() => setCurrentStep('payment')} />
      )}
    </div>
  );
};
```

### Phase 5: API Endpoints (3-5 дней)

#### 5.1 Payment API

```go
// API endpoints для платежей
// POST /api/v1/payments/create
// POST /api/v1/payments/{id}/capture  
// POST /api/v1/payments/{id}/refund
// GET /api/v1/payments/{id}/status
// POST /api/v1/webhooks/allsecure

// Swagger документация
// @Summary Create Payment
// @Description Создает новый платеж через AllSecure
// @Tags payments
// @Accept json
// @Produce json
// @Param request body CreatePaymentRequest true "Payment details"
// @Success 200 {object} utils.SuccessResponseSwag{data=CreatePaymentResponse}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Router /payments/create [post]
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error
```

### Phase 6: Testing (1 неделя)

#### 6.1 Unit Tests
- AllSecure API Client тесты
- Payment Service тесты  
- Repository тесты
- Webhook handling тесты

#### 6.2 Integration Tests
- Полный payment flow
- Escrow functionality
- Payout процессы
- Error handling

#### 6.3 End-to-End Tests
- Frontend payment widget
- Complete purchase flow
- Refund scenarios
- Webhook delivery

### Phase 7: Security и Compliance (3-5 дней)

#### 7.1 PCI DSS Compliance
- [ ] SAQ-A questionnaire заполнение
- [ ] Security scan отчеты
- [ ] SSL/TLS конфигурация
- [ ] Data encryption at rest

#### 7.2 Security Measures
- [ ] Webhook signature verification
- [ ] API credentials secure storage
- [ ] Audit logging
- [ ] Rate limiting

### Phase 8: Production Deployment (2-3 дня)

#### 8.1 Environment Setup
- [ ] Production credentials от AllSecure
- [ ] SSL certificates
- [ ] Monitoring dashboards
- [ ] Alerting rules

#### 8.2 Go-Live Checklist
- [ ] Database migrations
- [ ] Configuration deployment
- [ ] Feature flags activation
- [ ] Payment testing с real cards
- [ ] Monitoring verification

## Risk Management

### Технические риски
- **API изменения**: AllSecure может изменить API
- **Network issues**: Проблемы с connectivity
- **Performance**: Медленные ответы от gateway

### Mitigation
- **Retry logic** с exponential backoff
- **Circuit breaker** pattern
- **Health checks** и monitoring
- **Fallback mechanism** на другие payment methods

### Business риски
- **Compliance требования**: PCI DSS, local regulations
- **Chargebacks**: Споры по платежам
- **Fraud**: Мошеннические транзакции

### Mitigation  
- **KYC процедуры**
- **Fraud detection** rules
- **Transaction monitoring**
- **Clear refund policies**

## Success Metrics

### Technical KPIs
- **Payment success rate**: >95%
- **API response time**: <2 seconds
- **Uptime**: >99.5%
- **Error rate**: <1%

### Business KPIs
- **Conversion rate**: Increase from current
- **Transaction volume**: Month-over-month growth
- **Customer satisfaction**: Payment experience rating
- **Chargeback rate**: <0.5%

## Timeline Summary

| Phase | Duration | Dependencies |
|-------|----------|--------------|
| 1. Подготовка | 1-2 недели | AllSecure approval |
| 2. Database | 3-5 дней | - |
| 3. Backend | 1-2 недели | Phase 2 |
| 4. Frontend | 1 неделя | Phase 3 |
| 5. API | 3-5 дней | Phase 3 |
| 6. Testing | 1 неделя | Phase 4,5 |
| 7. Security | 3-5 дней | - |
| 8. Deployment | 2-3 дня | All phases |

**Общая длительность**: 6-8 недель

## Следующие шаги

1. **Немедленно**: Связаться с AllSecure для получения demo credentials
2. **На этой неделе**: Начать Phase 1 (подготовка)
3. **Через неделю**: Приступить к техническому design
4. **В течение месяца**: Завершить core implementation

Этот план обеспечивает поэтапную интеграцию с минимальными рисками и максимальной гибкостью.