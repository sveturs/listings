package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// PaymentGateway представляет платежный шлюз
type PaymentGateway struct {
	ID        int                    `json:"id" db:"id"`
	Name      string                 `json:"name" db:"name"`
	IsActive  bool                   `json:"is_active" db:"is_active"`
	Config    map[string]interface{} `json:"config" db:"config"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// PaymentSource представляет тип источника платежа
type PaymentSource string

const (
	PaymentSourceMarketplaceListing PaymentSource = "marketplace_listing"
	PaymentSourceStorefrontOrder    PaymentSource = "storefront_order"
)

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

// EscrowPayment представляет эскроу платеж
type EscrowPayment struct {
	ID                    int64           `json:"id" db:"id"`
	PaymentTransactionID  int64           `json:"payment_transaction_id" db:"payment_transaction_id"`
	SellerID              int             `json:"seller_id" db:"seller_id"`
	BuyerID               int             `json:"buyer_id" db:"buyer_id"`
	ListingID             int             `json:"listing_id" db:"listing_id"`
	Amount                decimal.Decimal `json:"amount" db:"amount"`
	MarketplaceCommission decimal.Decimal `json:"marketplace_commission" db:"marketplace_commission"`
	SellerAmount          decimal.Decimal `json:"seller_amount" db:"seller_amount"`
	Status                string          `json:"status" db:"status"`
	ReleaseDate           *time.Time      `json:"release_date,omitempty" db:"release_date"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`

	// Relations
	PaymentTransaction *PaymentTransaction `json:"payment_transaction,omitempty"`
	Seller             *User               `json:"seller,omitempty"`
	Buyer              *User               `json:"buyer,omitempty"`
	Listing            *MarketplaceListing `json:"listing,omitempty"`
}

// MerchantPayout представляет выплату продавцу
type MerchantPayout struct {
	ID                 int64                  `json:"id" db:"id"`
	SellerID           int                    `json:"seller_id" db:"seller_id"`
	GatewayID          int                    `json:"gateway_id" db:"gateway_id"`
	Amount             decimal.Decimal        `json:"amount" db:"amount"`
	Currency           string                 `json:"currency" db:"currency"`
	GatewayPayoutID    *string                `json:"gateway_payout_id,omitempty" db:"gateway_payout_id"`
	GatewayReferenceID *string                `json:"gateway_reference_id,omitempty" db:"gateway_reference_id"`
	Status             string                 `json:"status" db:"status"`
	BankAccountInfo    map[string]interface{} `json:"bank_account_info,omitempty" db:"bank_account_info"`
	GatewayResponse    map[string]interface{} `json:"gateway_response,omitempty" db:"gateway_response"`
	ErrorDetails       map[string]interface{} `json:"error_details,omitempty" db:"error_details"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
	ProcessedAt        *time.Time             `json:"processed_at,omitempty" db:"processed_at"`

	// Relations
	Seller  *User           `json:"seller,omitempty"`
	Gateway *PaymentGateway `json:"gateway,omitempty"`
}

// PaymentStatuses определяет возможные статусы платежей
const (
	PaymentStatusPending    = "pending"
	PaymentStatusAuthorized = "authorized"
	PaymentStatusCaptured   = "captured"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusVoided     = "voided"
)

// EscrowStatuses определяет возможные статусы эскроу
const (
	EscrowStatusHeld     = "held"
	EscrowStatusReleased = "released"
	EscrowStatusRefunded = "refunded"
)

// PayoutStatuses определяет возможные статусы выплат
const (
	PayoutStatusPending    = "pending"
	PayoutStatusProcessing = "processing"
	PayoutStatusCompleted  = "completed"
	PayoutStatusFailed     = "failed"
)
