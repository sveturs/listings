# Определения структур Storefront, StorefrontOrder и PaymentTransaction

## Структура Storefront
Файл: `/data/hostel-booking-system/backend/internal/domain/models/storefront.go`

```go
// Storefront представляет структуру витрины
type Storefront struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Брендинг
	LogoURL   string `json:"logo_url,omitempty"`
	BannerURL string `json:"banner_url,omitempty"`
	Theme     JSONB  `json:"theme"`

	// Контактная информация
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Website string `json:"website,omitempty"`

	// Локация
	Address    string   `json:"address,omitempty"`
	City       string   `json:"city,omitempty"`
	PostalCode string   `json:"postal_code,omitempty"`
	Country    string   `json:"country"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`

	// Настройки бизнеса
	Settings JSONB `json:"settings"`
	SEOMeta  JSONB `json:"seo_meta"`

	// Статус и статистика
	IsActive         bool       `json:"is_active"`
	IsVerified       bool       `json:"is_verified"`
	VerificationDate *time.Time `json:"verification_date,omitempty"`
	Rating           float64    `json:"rating"`
	ReviewsCount     int        `json:"reviews_count"`
	ProductsCount    int        `json:"products_count"`
	SalesCount       int        `json:"sales_count"`
	ViewsCount       int        `json:"views_count"`

	// Подписка (монетизация)
	SubscriptionPlan      SubscriptionPlan `json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time       `json:"subscription_expires_at,omitempty"`
	CommissionRate        float64          `json:"commission_rate"`

	// AI и killer features
	AIAgentEnabled      bool  `json:"ai_agent_enabled"`
	AIAgentConfig       JSONB `json:"ai_agent_config"`
	LiveShoppingEnabled bool  `json:"live_shopping_enabled"`
	GroupBuyingEnabled  bool  `json:"group_buying_enabled"`

	// Временные метки
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

## Структура StorefrontOrder
Файл: `/data/hostel-booking-system/backend/internal/domain/models/storefront_order.go`

```go
// StorefrontOrder представляет заказ в витрине
type StorefrontOrder struct {
	ID           int64  `json:"id" db:"id"`
	OrderNumber  string `json:"order_number" db:"order_number"`
	StorefrontID int    `json:"storefront_id" db:"storefront_id"`
	UserID       int    `json:"user_id" db:"customer_id"`     // Для совместимости с API
	CustomerID   int    `json:"customer_id" db:"customer_id"` // Реальное поле в БД

	// Финансовые данные
	Subtotal         decimal.Decimal `json:"subtotal" db:"subtotal"` // Алиас для SubtotalAmount
	SubtotalAmount   decimal.Decimal `json:"subtotal_amount" db:"subtotal"`
	Tax              decimal.Decimal `json:"tax" db:"tax"` // Алиас для TaxAmount
	TaxAmount        decimal.Decimal `json:"tax_amount" db:"tax"`
	Shipping         decimal.Decimal `json:"shipping" db:"shipping"` // Алиас для ShippingAmount
	ShippingAmount   decimal.Decimal `json:"shipping_amount" db:"shipping"`
	Discount         decimal.Decimal `json:"discount" db:"discount"`
	Total            decimal.Decimal `json:"total" db:"total"` // Алиас для TotalAmount
	TotalAmount      decimal.Decimal `json:"total_amount" db:"total"`
	CommissionAmount decimal.Decimal `json:"commission_amount" db:"commission_amount"`
	SellerAmount     decimal.Decimal `json:"seller_amount" db:"seller_amount"`
	Currency         string          `json:"currency" db:"currency"`

	// Платежные данные
	PaymentMethod        string  `json:"payment_method" db:"payment_method"`
	PaymentStatus        string  `json:"payment_status" db:"payment_status"`
	PaymentTransactionID *string `json:"payment_transaction_id,omitempty" db:"payment_transaction_id"`

	// Статус и escrow
	Status            OrderStatus `json:"status" db:"status"`
	EscrowReleaseDate *time.Time  `json:"escrow_release_date,omitempty" db:"escrow_release_date"`
	EscrowDays        int         `json:"escrow_days" db:"escrow_days"`

	// Доставка и адреса
	ShippingAddress  JSONB   `json:"shipping_address,omitempty" db:"shipping_address"`
	BillingAddress   JSONB   `json:"billing_address,omitempty" db:"billing_address"`
	ShippingMethod   *string `json:"shipping_method,omitempty" db:"shipping_method"`
	ShippingProvider *string `json:"shipping_provider,omitempty" db:"shipping_provider"`
	TrackingNumber   *string `json:"tracking_number,omitempty" db:"tracking_number"`

	// Заметки и метаданные
	Notes         *string                `json:"notes,omitempty" db:"notes"`
	CustomerNotes *string                `json:"customer_notes,omitempty" db:"customer_notes"`
	SellerNotes   *string                `json:"seller_notes,omitempty" db:"seller_notes"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" db:"metadata"`

	// Связанные данные
	Items              []StorefrontOrderItem `json:"items,omitempty"`
	Storefront         *Storefront           `json:"storefront,omitempty"`
	Customer           *User                 `json:"customer,omitempty"`
	PaymentTransaction *PaymentTransaction   `json:"payment_transaction,omitempty"`

	// Временные метки
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty" db:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
```

## Структура PaymentTransaction
Файл: `/data/hostel-booking-system/backend/internal/domain/models/payment_gateway.go`

```go
// PaymentTransaction представляет платежную транзакцию
type PaymentTransaction struct {
	ID        int64 `json:"id" db:"id"`
	GatewayID int   `json:"gateway_id" db:"gateway_id"`
	UserID    int   `json:"user_id" db:"user_id"`

	// Legacy поле для обратной совместимости
	ListingID *int `json:"listing_id,omitempty" db:"listing_id"`

	// Новые поля для унифицированной системы
	SourceType   PaymentSource `json:"source_type" db:"source_type"`
	SourceID     *int64        `json:"source_id,omitempty" db:"source_id"`
	StorefrontID *int          `json:"storefront_id,omitempty" db:"storefront_id"`

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

	// Delayed capture fields
	CaptureMode        string     `json:"capture_mode" db:"capture_mode"` // 'auto' или 'manual'
	AutoCaptureAt      *time.Time `json:"auto_capture_at,omitempty" db:"auto_capture_at"`
	CaptureDeadlineAt  *time.Time `json:"capture_deadline_at,omitempty" db:"capture_deadline_at"`
	CaptureAttemptedAt *time.Time `json:"capture_attempted_at,omitempty" db:"capture_attempted_at"`
	CaptureAttempts    int        `json:"capture_attempts" db:"capture_attempts"`

	// Relations
	Gateway *PaymentGateway     `json:"gateway,omitempty"`
	User    *User               `json:"user,omitempty"`
	Listing *MarketplaceListing `json:"listing,omitempty"`
}
```

## Связанные типы и константы

### Тип JSONB
```go
// JSONB тип для работы с JSONB полями PostgreSQL
type JSONB map[string]interface{}
```

### Типы подписки
```go
type SubscriptionPlan string

const (
	SubscriptionPlanStarter      SubscriptionPlan = "starter"
	SubscriptionPlanProfessional SubscriptionPlan = "professional"
	SubscriptionPlanBusiness     SubscriptionPlan = "business"
	SubscriptionPlanEnterprise   SubscriptionPlan = "enterprise"
)
```

### Статусы заказов
```go
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // ожидает оплаты
	OrderStatusConfirmed  OrderStatus = "confirmed"  // оплачен, подтвержден
	OrderStatusProcessing OrderStatus = "processing" // в обработке
	OrderStatusShipped    OrderStatus = "shipped"    // отправлен
	OrderStatusDelivered  OrderStatus = "delivered"  // доставлен
	OrderStatusCancelled  OrderStatus = "cancelled"  // отменен
	OrderStatusRefunded   OrderStatus = "refunded"   // возвращен
)
```

### Типы источников платежей
```go
type PaymentSource string

const (
	PaymentSourceMarketplaceListing PaymentSource = "marketplace_listing"
	PaymentSourceStorefrontOrder    PaymentSource = "storefront_order"
)
```

### Статусы платежей
```go
const (
	PaymentStatusPending    = "pending"
	PaymentStatusAuthorized = "authorized"
	PaymentStatusCaptured   = "captured"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusVoided     = "voided"
)
```