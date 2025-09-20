package interfaces

import (
	"context"
	"time"
)

// DeliveryProvider - универсальный интерфейс для всех провайдеров доставки
type DeliveryProvider interface {
	// GetCode возвращает уникальный код провайдера
	GetCode() string

	// GetName возвращает название провайдера
	GetName() string

	// IsActive проверяет, активен ли провайдер
	IsActive() bool

	// CalculateRate рассчитывает стоимость доставки
	CalculateRate(ctx context.Context, req *RateRequest) (*RateResponse, error)

	// CreateShipment создает отправление
	CreateShipment(ctx context.Context, req *ShipmentRequest) (*ShipmentResponse, error)

	// TrackShipment отслеживает отправление
	TrackShipment(ctx context.Context, trackingNumber string) (*TrackingResponse, error)

	// CancelShipment отменяет отправление
	CancelShipment(ctx context.Context, shipmentID string) error

	// GetLabel получает этикетку для печати
	GetLabel(ctx context.Context, shipmentID string) (*LabelResponse, error)

	// ValidateAddress проверяет корректность адреса
	ValidateAddress(ctx context.Context, address *Address) (*AddressValidationResponse, error)

	// GetCapabilities возвращает возможности провайдера
	GetCapabilities() *ProviderCapabilities

	// HandleWebhook обрабатывает webhook от провайдера
	HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookResponse, error)
}

// Address - структура адреса
type Address struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Email       string `json:"email,omitempty"`
	Street      string `json:"street"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	CompanyName string `json:"company_name,omitempty"`
	Note        string `json:"note,omitempty"`
}

// Package - параметры посылки
type Package struct {
	Weight       float64       `json:"weight_kg"`          // Вес в килограммах
	Dimensions   *Dimensions   `json:"dimensions"`         // Габариты
	IsFragile    bool          `json:"is_fragile"`         // Хрупкий товар
	Value        float64       `json:"value,omitempty"`    // Объявленная ценность
	Description  string        `json:"description"`        // Описание содержимого
	Attributes   *DeliveryAttrs `json:"attributes,omitempty"` // Дополнительные атрибуты
}

// Dimensions - габариты посылки
type Dimensions struct {
	Length float64 `json:"length_cm"` // Длина в сантиметрах
	Width  float64 `json:"width_cm"`  // Ширина в сантиметрах
	Height float64 `json:"height_cm"` // Высота в сантиметрах
}

// DeliveryAttrs - атрибуты доставки товара
type DeliveryAttrs struct {
	WeightKg              float64     `json:"weight_kg"`
	Dimensions            *Dimensions `json:"dimensions"`
	VolumeM3              float64     `json:"volume_m3,omitempty"`
	IsFragile             bool        `json:"is_fragile"`
	RequiresSpecialHandling bool      `json:"requires_special_handling"`
	Stackable             bool        `json:"stackable"`
	MaxStackWeightKg      float64     `json:"max_stack_weight_kg,omitempty"`
	PackagingType         string      `json:"packaging_type"` // box, envelope, pallet, custom
	HazmatClass           *string     `json:"hazmat_class,omitempty"`
}

// RateRequest - запрос расчета стоимости
type RateRequest struct {
	FromAddress   *Address        `json:"from_address"`
	ToAddress     *Address        `json:"to_address"`
	Packages      []Package       `json:"packages"`
	DeliveryType  string          `json:"delivery_type,omitempty"` // standard, express, same_day
	InsuranceValue float64        `json:"insurance_value,omitempty"`
	CODAmount     float64         `json:"cod_amount,omitempty"`
	Services      []string        `json:"services,omitempty"` // дополнительные услуги
}

// RateResponse - ответ с расчетом стоимости
type RateResponse struct {
	ProviderCode    string           `json:"provider_code"`
	ProviderName    string           `json:"provider_name"`
	DeliveryOptions []DeliveryOption `json:"delivery_options"`
	Currency        string           `json:"currency"`
	ValidUntil      time.Time        `json:"valid_until"`
}

// DeliveryOption - вариант доставки
type DeliveryOption struct {
	Type            string          `json:"type"` // standard, express, same_day
	Name            string          `json:"name"`
	TotalCost       float64         `json:"total_cost"`
	CostBreakdown   *CostBreakdown  `json:"cost_breakdown,omitempty"`
	EstimatedDays   int             `json:"estimated_days"`
	EstimatedDate   *time.Time      `json:"estimated_date,omitempty"`
	Services        []string        `json:"services,omitempty"`
}

// CostBreakdown - детализация стоимости
type CostBreakdown struct {
	BasePrice          float64 `json:"base_price"`
	WeightSurcharge    float64 `json:"weight_surcharge,omitempty"`
	OversizeSurcharge  float64 `json:"oversize_surcharge,omitempty"`
	FragileSurcharge   float64 `json:"fragile_surcharge,omitempty"`
	InsuranceFee       float64 `json:"insurance_fee,omitempty"`
	CODFee             float64 `json:"cod_fee,omitempty"`
	FuelSurcharge      float64 `json:"fuel_surcharge,omitempty"`
	RemoteAreaSurcharge float64 `json:"remote_area_surcharge,omitempty"`
	Discount           float64 `json:"discount,omitempty"`
	Tax                float64 `json:"tax,omitempty"`
}

// ShipmentRequest - запрос создания отправления
type ShipmentRequest struct {
	OrderID        int             `json:"order_id"`
	FromAddress    *Address        `json:"from_address"`
	ToAddress      *Address        `json:"to_address"`
	Packages       []Package       `json:"packages"`
	DeliveryType   string          `json:"delivery_type"`
	PickupDate     *time.Time      `json:"pickup_date,omitempty"`
	InsuranceValue float64         `json:"insurance_value,omitempty"`
	CODAmount      float64         `json:"cod_amount,omitempty"`
	Services       []string        `json:"services,omitempty"`
	Reference      string          `json:"reference,omitempty"` // внутренний номер заказа
	Notes          string          `json:"notes,omitempty"`
}

// ShipmentResponse - ответ создания отправления
type ShipmentResponse struct {
	ShipmentID      string          `json:"shipment_id"`
	TrackingNumber  string          `json:"tracking_number"`
	ExternalID      string          `json:"external_id,omitempty"` // ID в системе провайдера
	Status          string          `json:"status"`
	TotalCost       float64         `json:"total_cost"`
	CostBreakdown   *CostBreakdown  `json:"cost_breakdown,omitempty"`
	EstimatedDate   *time.Time      `json:"estimated_date,omitempty"`
	Labels          []LabelInfo     `json:"labels,omitempty"`
	PickupInfo      *PickupInfo     `json:"pickup_info,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// LabelInfo - информация об этикетке
type LabelInfo struct {
	Type   string `json:"type"` // shipping, return, customs
	Format string `json:"format"` // pdf, zpl, png
	URL    string `json:"url,omitempty"`
	Data   []byte `json:"data,omitempty"`
}

// PickupInfo - информация о заборе груза
type PickupInfo struct {
	PickupDate     time.Time `json:"pickup_date"`
	PickupWindow   string    `json:"pickup_window,omitempty"` // например, "09:00-18:00"
	PickupLocation string    `json:"pickup_location,omitempty"`
	Instructions   string    `json:"instructions,omitempty"`
}

// TrackingResponse - ответ отслеживания
type TrackingResponse struct {
	TrackingNumber string           `json:"tracking_number"`
	Status         string           `json:"status"`
	StatusText     string           `json:"status_text"`
	CurrentLocation string          `json:"current_location,omitempty"`
	EstimatedDate  *time.Time       `json:"estimated_date,omitempty"`
	DeliveredDate  *time.Time       `json:"delivered_date,omitempty"`
	Events         []TrackingEvent  `json:"events"`
	ProofOfDelivery *ProofOfDelivery `json:"proof_of_delivery,omitempty"`
}

// TrackingEvent - событие отслеживания
type TrackingEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location,omitempty"`
	Details     string    `json:"details,omitempty"`
}

// ProofOfDelivery - подтверждение доставки
type ProofOfDelivery struct {
	RecipientName string    `json:"recipient_name"`
	SignatureURL  string    `json:"signature_url,omitempty"`
	PhotoURL      string    `json:"photo_url,omitempty"`
	DeliveredAt   time.Time `json:"delivered_at"`
	Notes         string    `json:"notes,omitempty"`
}

// LabelResponse - ответ с этикеткой
type LabelResponse struct {
	Labels []LabelInfo `json:"labels"`
}

// AddressValidationResponse - ответ валидации адреса
type AddressValidationResponse struct {
	IsValid           bool      `json:"is_valid"`
	SuggestedAddress  *Address  `json:"suggested_address,omitempty"`
	ValidationErrors  []string  `json:"validation_errors,omitempty"`
	DeliveryAvailable bool      `json:"delivery_available"`
	Zone              string    `json:"zone,omitempty"` // local, regional, national, international
}

// ProviderCapabilities - возможности провайдера
type ProviderCapabilities struct {
	MaxWeightKg      float64  `json:"max_weight_kg"`
	MaxVolumeM3      float64  `json:"max_volume_m3"`
	MaxLengthCm      float64  `json:"max_length_cm"`
	DeliveryZones    []string `json:"delivery_zones"`
	DeliveryTypes    []string `json:"delivery_types"`
	SupportsCOD      bool     `json:"supports_cod"`
	SupportsInsurance bool    `json:"supports_insurance"`
	SupportsTracking bool     `json:"supports_tracking"`
	SupportsPickup   bool     `json:"supports_pickup"`
	SupportsReturn   bool     `json:"supports_return"`
	Services         []string `json:"services"` // список дополнительных услуг
}

// DeliveryStatus - универсальные статусы доставки
const (
	StatusPending          = "pending"           // ожидает обработки
	StatusConfirmed        = "confirmed"         // подтверждено
	StatusPickedUp         = "picked_up"         // забрано курьером
	StatusInTransit        = "in_transit"        // в пути
	StatusOutForDelivery   = "out_for_delivery"  // передано на доставку
	StatusDelivered        = "delivered"         // доставлено
	StatusDeliveryAttempted = "delivery_attempted" // попытка доставки
	StatusReturning        = "returning"         // возвращается
	StatusReturned         = "returned"          // возвращено
	StatusCancelled        = "cancelled"         // отменено
	StatusLost             = "lost"              // потеряно
	StatusDamaged          = "damaged"           // повреждено
)

// DeliveryType - типы доставки
const (
	DeliveryTypeStandard = "standard" // стандартная
	DeliveryTypeExpress  = "express"  // экспресс
	DeliveryTypeSameDay  = "same_day" // в тот же день
	DeliveryTypeNextDay  = "next_day" // на следующий день
	DeliveryTypeEconomy  = "economy"  // эконом
)

// PackagingType - типы упаковки
const (
	PackagingTypeBox      = "box"      // коробка
	PackagingTypeEnvelope = "envelope" // конверт
	PackagingTypePallet   = "pallet"   // паллета
	PackagingTypeCustom   = "custom"   // нестандартная
)

// WebhookResponse - ответ обработки webhook
type WebhookResponse struct {
	TrackingNumber string           `json:"tracking_number"`
	Status         string           `json:"status"`
	StatusDetails  string           `json:"status_details,omitempty"`
	Location       string           `json:"location,omitempty"`
	Timestamp      time.Time        `json:"timestamp"`
	DeliveryDetails *ProofOfDelivery `json:"delivery_details,omitempty"`
	Events         []TrackingEvent  `json:"events,omitempty"`
	Processed      bool             `json:"processed"`
	Error          string           `json:"error,omitempty"`
}