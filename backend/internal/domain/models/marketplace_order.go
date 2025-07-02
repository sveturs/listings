package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

// MarketplaceOrderStatus представляет статус заказа маркетплейса
type MarketplaceOrderStatus string

const (
	MarketplaceOrderStatusPending   MarketplaceOrderStatus = "pending"
	MarketplaceOrderStatusPaid      MarketplaceOrderStatus = "paid"
	MarketplaceOrderStatusShipped   MarketplaceOrderStatus = "shipped"
	MarketplaceOrderStatusDelivered MarketplaceOrderStatus = "delivered"
	MarketplaceOrderStatusCompleted MarketplaceOrderStatus = "completed"
	MarketplaceOrderStatusDisputed  MarketplaceOrderStatus = "disputed"
	MarketplaceOrderStatusCancelled MarketplaceOrderStatus = "cancelled"
	MarketplaceOrderStatusRefunded  MarketplaceOrderStatus = "refunded"
)

// Scan implements sql.Scanner interface
func (s *MarketplaceOrderStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		*s = MarketplaceOrderStatus(v)
	case []byte:
		*s = MarketplaceOrderStatus(v)
	default:
		return errors.New("invalid marketplace order status type")
	}
	return nil
}

// Value implements driver.Valuer interface
func (s MarketplaceOrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// IsValid проверяет валидность статуса
func (s MarketplaceOrderStatus) IsValid() bool {
	switch s {
	case MarketplaceOrderStatusPending, MarketplaceOrderStatusPaid, MarketplaceOrderStatusShipped,
		MarketplaceOrderStatusDelivered, MarketplaceOrderStatusCompleted, MarketplaceOrderStatusDisputed,
		MarketplaceOrderStatusCancelled, MarketplaceOrderStatusRefunded:
		return true
	}
	return false
}

// CanTransitionTo проверяет возможность перехода в новый статус
func (s MarketplaceOrderStatus) CanTransitionTo(newStatus MarketplaceOrderStatus) bool {
	transitions := map[MarketplaceOrderStatus][]MarketplaceOrderStatus{
		MarketplaceOrderStatusPending:   {MarketplaceOrderStatusPaid, MarketplaceOrderStatusCancelled},
		MarketplaceOrderStatusPaid:      {MarketplaceOrderStatusShipped, MarketplaceOrderStatusCancelled, MarketplaceOrderStatusRefunded},
		MarketplaceOrderStatusShipped:   {MarketplaceOrderStatusDelivered, MarketplaceOrderStatusDisputed},
		MarketplaceOrderStatusDelivered: {MarketplaceOrderStatusCompleted, MarketplaceOrderStatusDisputed},
		MarketplaceOrderStatusDisputed:  {MarketplaceOrderStatusCompleted, MarketplaceOrderStatusRefunded},
		// Финальные статусы
		MarketplaceOrderStatusCompleted: {},
		MarketplaceOrderStatusCancelled: {},
		MarketplaceOrderStatusRefunded:  {},
	}

	allowedStatuses, exists := transitions[s]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}
	return false
}

// MarketplaceOrder представляет заказ в маркетплейсе
type MarketplaceOrder struct {
	ID                   int64                  `db:"id" json:"id"`
	BuyerID              int64                  `db:"buyer_id" json:"buyer_id"`
	SellerID             int64                  `db:"seller_id" json:"seller_id"`
	ListingID            int64                  `db:"listing_id" json:"listing_id"`
	ItemPrice            float64                `db:"item_price" json:"item_price"`
	PlatformFeeRate      float64                `db:"platform_fee_rate" json:"platform_fee_rate"`
	PlatformFeeAmount    float64                `db:"platform_fee_amount" json:"platform_fee_amount"`
	SellerPayoutAmount   float64                `db:"seller_payout_amount" json:"seller_payout_amount"`
	PaymentTransactionID *int64                 `db:"payment_transaction_id" json:"payment_transaction_id,omitempty"`
	Status               MarketplaceOrderStatus `db:"status" json:"status"`
	ProtectionPeriodDays int                    `db:"protection_period_days" json:"protection_period_days"`
	ProtectionExpiresAt  *time.Time             `db:"protection_expires_at" json:"protection_expires_at,omitempty"`
	ShippingMethod       *string                `db:"shipping_method" json:"shipping_method,omitempty"`
	TrackingNumber       *string                `db:"tracking_number" json:"tracking_number,omitempty"`
	ShippedAt            *time.Time             `db:"shipped_at" json:"shipped_at,omitempty"`
	DeliveredAt          *time.Time             `db:"delivered_at" json:"delivered_at,omitempty"`
	CreatedAt            time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time              `db:"updated_at" json:"updated_at"`

	// Связанные данные (заполняются при необходимости)
	Buyer              *User                `db:"-" json:"buyer,omitempty"`
	Seller             *User                `db:"-" json:"seller,omitempty"`
	Listing            *MarketplaceListing  `db:"-" json:"listing,omitempty"`
	PaymentTransaction *PaymentTransaction  `db:"-" json:"payment_transaction,omitempty"`
	StatusHistory      []OrderStatusHistory `db:"-" json:"status_history,omitempty"`
	Messages           []*OrderMessage      `db:"-" json:"messages,omitempty"`
}

// CalculateFees рассчитывает комиссии
func (o *MarketplaceOrder) CalculateFees() {
	o.PlatformFeeAmount = o.ItemPrice * (o.PlatformFeeRate / 100)
	o.SellerPayoutAmount = o.ItemPrice - o.PlatformFeeAmount
}

// SetProtectionExpiry устанавливает дату окончания защитного периода
func (o *MarketplaceOrder) SetProtectionExpiry() {
	if o.DeliveredAt != nil {
		expiresAt := o.DeliveredAt.Add(time.Duration(o.ProtectionPeriodDays) * 24 * time.Hour)
		o.ProtectionExpiresAt = &expiresAt
	}
}

// IsProtectionActive проверяет активен ли защитный период
func (o *MarketplaceOrder) IsProtectionActive() bool {
	if o.ProtectionExpiresAt == nil {
		return false
	}
	return time.Now().Before(*o.ProtectionExpiresAt)
}

// CanBeCaptured проверяет можно ли захватить платеж
func (o *MarketplaceOrder) CanBeCaptured() bool {
	// Можно захватить если:
	// 1. Статус completed
	// 2. Или доставлен и истек защитный период
	if o.Status == MarketplaceOrderStatusCompleted {
		return true
	}
	if o.Status == MarketplaceOrderStatusDelivered && !o.IsProtectionActive() {
		return true
	}
	return false
}

// OrderStatusHistory представляет историю изменения статуса заказа
type OrderStatusHistory struct {
	ID        int64     `db:"id" json:"id"`
	OrderID   int64     `db:"order_id" json:"order_id"`
	OldStatus *string   `db:"old_status" json:"old_status,omitempty"`
	NewStatus string    `db:"new_status" json:"new_status"`
	Reason    *string   `db:"reason" json:"reason,omitempty"`
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// OrderMessageType тип сообщения в заказе
type OrderMessageType string

const (
	OrderMessageTypeText           OrderMessageType = "text"
	OrderMessageTypeShippingUpdate OrderMessageType = "shipping_update"
	OrderMessageTypeDisputeOpened  OrderMessageType = "dispute_opened"
	OrderMessageTypeDisputeMessage OrderMessageType = "dispute_message"
	OrderMessageTypeSystem         OrderMessageType = "system"
)

// OrderMessage представляет сообщение в заказе
type OrderMessage struct {
	ID          int64                  `db:"id" json:"id"`
	OrderID     int64                  `db:"order_id" json:"order_id"`
	SenderID    int64                  `db:"sender_id" json:"sender_id"`
	MessageType OrderMessageType       `db:"message_type" json:"message_type"`
	Content     string                 `db:"content" json:"content"`
	Metadata    map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt   time.Time              `db:"created_at" json:"created_at"`

	// Связанные данные
	Sender *User `db:"-" json:"sender,omitempty"`
}
