package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// OrderStatus представляет статус заказа
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // ожидает оплаты
	OrderStatusConfirmed  OrderStatus = "confirmed"  // оплачен, подтвержден
	OrderStatusProcessing OrderStatus = "processing" // в обработке
	OrderStatusShipped    OrderStatus = "shipped"    // отправлен
	OrderStatusDelivered  OrderStatus = "delivered"  // доставлен
	OrderStatusCancelled  OrderStatus = "canceled"   // отменен
	OrderStatusRefunded   OrderStatus = "refunded"   // возвращен
)

// ReservationStatus представляет статус резервирования товара
type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "active"    // активное резервирование
	ReservationStatusCommitted ReservationStatus = "committed" // подтверждено (списано)
	ReservationStatusReleased  ReservationStatus = "released"  // освобождено
	ReservationStatusExpired   ReservationStatus = "expired"   // истекло
)

// ShoppingCart представляет корзину покупок
type ShoppingCart struct {
	ID           int64              `json:"id"`
	UserID       *int               `json:"user_id,omitempty"`
	StorefrontID int                `json:"storefront_id"`
	SessionID    *string            `json:"session_id,omitempty"`
	Items        []ShoppingCartItem `json:"items,omitempty"`
	Storefront   *Storefront        `json:"storefront,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// ShoppingCartItem представляет позицию в корзине
type ShoppingCartItem struct {
	ID           int64                     `json:"id"`
	CartID       int64                     `json:"cart_id"`
	ProductID    int64                     `json:"product_id"`
	VariantID    *int64                    `json:"variant_id,omitempty"`
	Quantity     int                       `json:"quantity"`
	PricePerUnit decimal.Decimal           `json:"price_per_unit"`
	TotalPrice   decimal.Decimal           `json:"total_price"`
	Product      *StorefrontProduct        `json:"product,omitempty"`
	Variant      *StorefrontProductVariant `json:"variant,omitempty"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

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
	PickupAddress    JSONB   `json:"pickup_address,omitempty" db:"pickup_address"` // Адрес забора товара у продавца
	ShippingMethod   *string `json:"shipping_method,omitempty" db:"shipping_method"`
	ShippingProvider *string `json:"shipping_provider,omitempty" db:"shipping_provider"`
	TrackingNumber   *string `json:"tracking_number,omitempty" db:"tracking_number"`
	ShipmentID       *int64  `json:"shipment_id,omitempty" db:"shipment_id"`
	DeliveryProvider *string `json:"delivery_provider,omitempty" db:"delivery_provider"` // Provider code для delivery microservice (только для UI)

	// Заметки и метаданные
	Notes         *string                `json:"notes,omitempty" db:"notes"`
	CustomerNotes *string                `json:"customer_notes,omitempty" db:"customer_notes"`
	SellerNotes   *string                `json:"seller_notes,omitempty" db:"seller_notes"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" db:"metadata"`

	// Связанные данные
	Items              []StorefrontOrderItem `json:"items,omitempty"`
	Storefront         *Storefront           `json:"storefront,omitempty"`
	Customer           *User                 `json:"customer,omitempty"`
	Seller             *User                 `json:"seller,omitempty"`
	PaymentTransaction *PaymentTransaction   `json:"payment_transaction,omitempty"`

	// Временные метки
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	ShippedAt   *time.Time `json:"shipped_at,omitempty" db:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	CancelledAt *time.Time `json:"canceled_at,omitempty" db:"canceled_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// StorefrontOrderItem представляет позицию заказа
type StorefrontOrderItem struct {
	ID        int64  `json:"id"`
	OrderID   int64  `json:"order_id"`
	ProductID int64  `json:"product_id"`
	VariantID *int64 `json:"variant_id,omitempty"`

	// Snapshot данных на момент заказа
	ProductName string  `json:"product_name"`
	ProductSKU  *string `json:"product_sku,omitempty"`
	VariantName *string `json:"variant_name,omitempty"`

	Quantity     int             `json:"quantity"`
	PricePerUnit decimal.Decimal `json:"price_per_unit"`
	TotalPrice   decimal.Decimal `json:"total_price"`

	// Snapshot атрибутов
	ProductAttributes JSONB `json:"product_attributes,omitempty"`

	// Связанные данные
	Product *StorefrontProduct        `json:"product,omitempty"`
	Variant *StorefrontProductVariant `json:"variant,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// InventoryReservation представляет резервирование товара
type InventoryReservation struct {
	ID        int64             `json:"id" db:"id"`
	ProductID int64             `json:"product_id" db:"product_id"`
	VariantID *int64            `json:"variant_id,omitempty" db:"variant_id"`
	OrderID   int64             `json:"order_id" db:"order_id"`
	Quantity  int               `json:"quantity" db:"quantity"`
	Status    ReservationStatus `json:"status" db:"status"`
	ExpiresAt time.Time         `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`

	// Deprecated - для обратной совместимости
	ReservedQuantity int        `json:"reserved_quantity,omitempty"`
	ReleasedAt       *time.Time `json:"released_at,omitempty"`

	// Связанные данные
	Product *StorefrontProduct        `json:"product,omitempty"`
	Variant *StorefrontProductVariant `json:"variant,omitempty"`
	Order   *StorefrontOrder          `json:"order,omitempty"`
}

// Расширенные структуры для API запросов/ответов

// ShippingAddress представляет адрес доставки
type ShippingAddress struct {
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email,omitempty"`
	Street      string `json:"street"`
	HouseNumber string `json:"house_number"`
	Apartment   string `json:"apartment,omitempty"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	Notes       string `json:"notes,omitempty"`
}

// CreateOrderRequest представляет запрос на создание заказа
type CreateOrderRequest struct {
	StorefrontID    int                `json:"storefront_id" validate:"required"`
	CartID          *int64             `json:"cart_id,omitempty"`
	Items           []OrderItemRequest `json:"items,omitempty"` // альтернатива cart_id
	ShippingAddress ShippingAddress    `json:"shipping_address" validate:"required"`
	BillingAddress  ShippingAddress    `json:"billing_address" validate:"required"`
	ShippingMethod  string             `json:"shipping_method" validate:"required"`
	PaymentMethod   string             `json:"payment_method" validate:"required"`
	CustomerNotes   string             `json:"customer_notes,omitempty"`
}

// OrderItemRequest представляет позицию в запросе создания заказа
type OrderItemRequest struct {
	ProductID int64  `json:"product_id" validate:"required"`
	VariantID *int64 `json:"variant_id,omitempty"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// UpdateOrderStatusRequest представляет запрос на обновление статуса заказа
type UpdateOrderStatusRequest struct {
	Status         OrderStatus `json:"status" validate:"required"`
	TrackingNumber *string     `json:"tracking_number,omitempty"`
	SellerNotes    *string     `json:"seller_notes,omitempty"`
}

// AddToCartRequest представляет запрос на добавление товара в корзину
type AddToCartRequest struct {
	ProductID int64  `json:"product_id" validate:"required"`
	VariantID *int64 `json:"variant_id,omitempty"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// UpdateCartItemRequest представляет запрос на обновление позиции в корзине
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=0"` // 0 = удалить из корзины
}

// CancelOrderRequest представляет запрос на отмену заказа
type CancelOrderRequest struct {
	Reason string `json:"reason,omitempty"`
}

// OrderFilter представляет фильтры для поиска заказов
type OrderFilter struct {
	StorefrontID   *int             `json:"storefront_id,omitempty"`
	UserID         *int             `json:"user_id,omitempty"` // Алиас для CustomerID
	CustomerID     *int             `json:"customer_id,omitempty"`
	SellerID       *int             `json:"seller_id,omitempty"`
	Status         *OrderStatus     `json:"status,omitempty"`
	PaymentStatus  *string          `json:"payment_status,omitempty"`
	DateFrom       *time.Time       `json:"date_from,omitempty"`
	DateTo         *time.Time       `json:"date_to,omitempty"`
	MinAmount      *decimal.Decimal `json:"min_amount,omitempty"`
	MaxAmount      *decimal.Decimal `json:"max_amount,omitempty"`
	OrderNumber    *string          `json:"order_number,omitempty"`
	TrackingNumber *string          `json:"tracking_number,omitempty"`

	// Пагинация
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	// Сортировка
	SortBy    string `json:"sort_by,omitempty"`    // created_at, total_amount, status
	SortOrder string `json:"sort_order,omitempty"` // asc, desc
}

// OrderSummary представляет краткую информацию о заказе
type OrderSummary struct {
	ID             int64           `json:"id"`
	OrderNumber    string          `json:"order_number"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	Currency       string          `json:"currency"`
	Status         OrderStatus     `json:"status"`
	ItemsCount     int             `json:"items_count"`
	CustomerName   string          `json:"customer_name"`
	StorefrontName string          `json:"storefront_name"`
	CreatedAt      time.Time       `json:"created_at"`
}

// Методы для работы с заказами

// CalculateTotals рассчитывает итоговые суммы заказа
func (o *StorefrontOrder) CalculateTotals() {
	o.SubtotalAmount = decimal.Zero
	for _, item := range o.Items {
		o.SubtotalAmount = o.SubtotalAmount.Add(item.TotalPrice)
	}
	o.TotalAmount = o.SubtotalAmount.Add(o.ShippingAmount).Add(o.TaxAmount)
	o.SellerAmount = o.TotalAmount.Sub(o.CommissionAmount)
}

// CanBeCancelled проверяет можно ли отменить заказ
func (o *StorefrontOrder) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

// CanBeRefunded проверяет можно ли вернуть заказ
func (o *StorefrontOrder) CanBeRefunded() bool {
	return o.Status == OrderStatusConfirmed || o.Status == OrderStatusProcessing ||
		o.Status == OrderStatusShipped || o.Status == OrderStatusDelivered
}

// IsEscrowExpired проверяет истек ли срок escrow
func (o *StorefrontOrder) IsEscrowExpired() bool {
	return o.EscrowReleaseDate != nil && time.Now().After(*o.EscrowReleaseDate)
}

// GetTotalQuantity возвращает общее количество товаров в корзине
func (c *ShoppingCart) GetTotalQuantity() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// GetTotalAmount возвращает общую сумму корзины
func (c *ShoppingCart) GetTotalAmount() decimal.Decimal {
	total := decimal.Zero
	for _, item := range c.Items {
		total = total.Add(item.TotalPrice)
	}
	return total
}

// UpdateTotalPrice обновляет общую стоимость позиции корзины
func (i *ShoppingCartItem) UpdateTotalPrice() {
	i.TotalPrice = i.PricePerUnit.Mul(decimal.NewFromInt(int64(i.Quantity)))
}

// IsExpired проверяет истекло ли резервирование
func (r *InventoryReservation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// StockUpdate представляет обновление остатков товара
type StockUpdate struct {
	ProductID      int64  `json:"product_id"`
	VariantID      *int64 `json:"variant_id,omitempty"`
	QuantityChange int    `json:"quantity_change"`
}

// LowStockItem представляет товар с низким остатком
type LowStockItem struct {
	ProductID         int64   `json:"product_id"`
	VariantID         *int64  `json:"variant_id,omitempty"`
	ProductName       string  `json:"product_name"`
	VariantName       *string `json:"variant_name,omitempty"`
	Quantity          int     `json:"quantity"`
	AvailableQuantity int     `json:"available_quantity"`
	LowStockThreshold int     `json:"low_stock_threshold"`
}

// InventoryMovementDTO представляет движение товаров для API
type InventoryMovementDTO struct {
	ID          int64     `json:"id"`
	ProductID   int64     `json:"product_id"`
	VariantID   *int64    `json:"variant_id,omitempty"`
	Type        string    `json:"type"`
	Quantity    int       `json:"quantity"`
	Reference   *string   `json:"reference,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ProductName *string   `json:"product_name,omitempty"`
}

// InventoryStock представляет текущий запас товара
type InventoryStock struct {
	ProductID         int64     `json:"product_id" db:"product_id"`
	VariantID         *int64    `json:"variant_id,omitempty" db:"variant_id"`
	Quantity          int       `json:"quantity" db:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity" db:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity" db:"available_quantity"`
	LowStockThreshold int       `json:"low_stock_threshold" db:"low_stock_threshold"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// InventoryMovement представляет движение товара в базе данных
type InventoryMovement struct {
	ID            int64                  `json:"id" db:"id"`
	ProductID     int64                  `json:"product_id" db:"product_id"`
	VariantID     *int64                 `json:"variant_id,omitempty" db:"variant_id"`
	Type          string                 `json:"type" db:"type"`
	Quantity      int                    `json:"quantity" db:"quantity"`
	ReferenceType *string                `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID   *int64                 `json:"reference_id,omitempty" db:"reference_id"`
	Notes         *string                `json:"notes,omitempty" db:"notes"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}
