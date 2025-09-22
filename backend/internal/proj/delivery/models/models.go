package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// Константы статусов отправления
const (
	ShipmentStatusPending    = "pending"
	ShipmentStatusProcessing = "processing"
	ShipmentStatusShipped    = "shipped"
	ShipmentStatusInTransit  = "in_transit"
	ShipmentStatusDelivered  = "delivered"
	ShipmentStatusCancelled  = "canceled"
	ShipmentStatusFailed     = "failed"
)

// Константы типов зон
const (
	ZoneTypeLocal         = "local"
	ZoneTypeRegional      = "regional"
	ZoneTypeNational      = "national"
	ZoneTypeInternational = "international"
)

// Provider - модель провайдера доставки
type Provider struct {
	ID                int              `json:"id" db:"id"`
	Code              string           `json:"code" db:"code"`
	Name              string           `json:"name" db:"name"`
	LogoURL           *string          `json:"logo_url,omitempty" db:"logo_url"`
	IsActive          bool             `json:"is_active" db:"is_active"`
	SupportsCOD       bool             `json:"supports_cod" db:"supports_cod"`
	SupportsInsurance bool             `json:"supports_insurance" db:"supports_insurance"`
	SupportsTracking  bool             `json:"supports_tracking" db:"supports_tracking"`
	APIConfig         *json.RawMessage `json:"api_config,omitempty" db:"api_config"`
	Capabilities      *json.RawMessage `json:"capabilities,omitempty" db:"capabilities"`
	CreatedAt         time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at" db:"updated_at"`
}

// Shipment - модель отправления
type Shipment struct {
	ID                 int             `json:"id" db:"id"`
	ProviderID         int             `json:"provider_id" db:"provider_id"`
	OrderID            *int            `json:"order_id,omitempty" db:"order_id"`
	ExternalID         *string         `json:"external_id,omitempty" db:"external_id"`
	TrackingNumber     *string         `json:"tracking_number,omitempty" db:"tracking_number"`
	Status             string          `json:"status" db:"status"`
	SenderInfo         json.RawMessage `json:"sender_info" db:"sender_info"`
	RecipientInfo      json.RawMessage `json:"recipient_info" db:"recipient_info"`
	PackageInfo        json.RawMessage `json:"package_info" db:"package_info"`
	DeliveryCost       *float64        `json:"delivery_cost,omitempty" db:"delivery_cost"`
	InsuranceCost      *float64        `json:"insurance_cost,omitempty" db:"insurance_cost"`
	CODAmount          *float64        `json:"cod_amount,omitempty" db:"cod_amount"`
	CostBreakdown      json.RawMessage `json:"cost_breakdown,omitempty" db:"cost_breakdown"`
	PickupDate         *time.Time      `json:"pickup_date,omitempty" db:"pickup_date"`
	EstimatedDelivery  *time.Time      `json:"estimated_delivery,omitempty" db:"estimated_delivery"`
	ActualDeliveryDate *time.Time      `json:"actual_delivery_date,omitempty" db:"actual_delivery_date"`
	ProviderResponse   json.RawMessage `json:"provider_response,omitempty" db:"provider_response"`
	Labels             json.RawMessage `json:"labels,omitempty" db:"labels"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`

	// Связанные данные (не из БД, заполняются при необходимости)
	Provider *Provider       `json:"provider,omitempty" db:"-"`
	Events   []TrackingEvent `json:"events,omitempty" db:"-"`
}

// TrackingEvent - модель события отслеживания
type TrackingEvent struct {
	ID          int             `json:"id" db:"id"`
	ShipmentID  int             `json:"shipment_id" db:"shipment_id"`
	ProviderID  int             `json:"provider_id" db:"provider_id"`
	EventTime   time.Time       `json:"event_time" db:"event_time"`
	Status      string          `json:"status" db:"status"`
	Location    *string         `json:"location,omitempty" db:"location"`
	Description *string         `json:"description,omitempty" db:"description"`
	RawData     json.RawMessage `json:"raw_data,omitempty" db:"raw_data"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
}

// CategoryDefaults - дефолтные атрибуты доставки для категории
type CategoryDefaults struct {
	ID                   int       `json:"id" db:"id"`
	CategoryID           int       `json:"category_id" db:"category_id"`
	DefaultWeightKg      *float64  `json:"default_weight_kg,omitempty" db:"default_weight_kg"`
	DefaultLengthCm      *float64  `json:"default_length_cm,omitempty" db:"default_length_cm"`
	DefaultWidthCm       *float64  `json:"default_width_cm,omitempty" db:"default_width_cm"`
	DefaultHeightCm      *float64  `json:"default_height_cm,omitempty" db:"default_height_cm"`
	DefaultPackagingType *string   `json:"default_packaging_type,omitempty" db:"default_packaging_type"`
	IsTypicallyFragile   bool      `json:"is_typically_fragile" db:"is_typically_fragile"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// PricingRule - правило расчета стоимости
type PricingRule struct {
	ID                       int             `json:"id" db:"id"`
	ProviderID               int             `json:"provider_id" db:"provider_id"`
	RuleType                 string          `json:"rule_type" db:"rule_type"`
	WeightRanges             json.RawMessage `json:"weight_ranges,omitempty" db:"weight_ranges"`
	VolumeRanges             json.RawMessage `json:"volume_ranges,omitempty" db:"volume_ranges"`
	ZoneMultipliers          json.RawMessage `json:"zone_multipliers,omitempty" db:"zone_multipliers"`
	FragileSurcharge         float64         `json:"fragile_surcharge" db:"fragile_surcharge"`
	OversizedSurcharge       float64         `json:"oversized_surcharge" db:"oversized_surcharge"`
	SpecialHandlingSurcharge float64         `json:"special_handling_surcharge" db:"special_handling_surcharge"`
	MinPrice                 *float64        `json:"min_price,omitempty" db:"min_price"`
	MaxPrice                 *float64        `json:"max_price,omitempty" db:"max_price"`
	CustomFormula            *string         `json:"custom_formula,omitempty" db:"custom_formula"`
	Priority                 int             `json:"priority" db:"priority"`
	IsActive                 bool            `json:"is_active" db:"is_active"`
	CreatedAt                time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time       `json:"updated_at" db:"updated_at"`
}

// Zone - зона доставки
type Zone struct {
	ID          int            `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Type        string         `json:"type" db:"type"`
	Countries   pq.StringArray `json:"countries,omitempty" db:"countries"`
	Regions     pq.StringArray `json:"regions,omitempty" db:"regions"`
	Cities      pq.StringArray `json:"cities,omitempty" db:"cities"`
	PostalCodes pq.StringArray `json:"postal_codes,omitempty" db:"postal_codes"`
	// Boundary и CenterPoint пропускаем для простоты, добавим при необходимости GIS функций
	RadiusKm  *float64  `json:"radius_km,omitempty" db:"radius_km"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// DeliveryAttributes - атрибуты доставки (для JSONB полей)
type DeliveryAttributes struct {
	WeightKg                float64     `json:"weight_kg"`
	Dimensions              *Dimensions `json:"dimensions,omitempty"`
	VolumeM3                float64     `json:"volume_m3,omitempty"`
	IsFragile               bool        `json:"is_fragile"`
	RequiresSpecialHandling bool        `json:"requires_special_handling"`
	Stackable               bool        `json:"stackable"`
	MaxStackWeightKg        float64     `json:"max_stack_weight_kg,omitempty"`
	PackagingType           string      `json:"packaging_type"`
	HazmatClass             *string     `json:"hazmat_class,omitempty"`
}

// Dimensions - габариты
type Dimensions struct {
	LengthCm float64 `json:"length_cm"`
	WidthCm  float64 `json:"width_cm"`
	HeightCm float64 `json:"height_cm"`
}

// CalculateVolume - рассчитывает объем в кубических метрах
func (d *Dimensions) CalculateVolume() float64 {
	if d == nil {
		return 0
	}
	return (d.LengthCm * d.WidthCm * d.HeightCm) / 1000000 // см³ -> м³
}

// CalculateVolumetricWeight - рассчитывает объемный вес
func (d *Dimensions) CalculateVolumetricWeight(divisor float64) float64 {
	if d == nil || divisor == 0 {
		return 0
	}
	return (d.LengthCm * d.WidthCm * d.HeightCm) / divisor
}

// Address - адрес доставки
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

// WeightRange - диапазон веса для расчета стоимости
type WeightRange struct {
	From       float64 `json:"from"`
	To         float64 `json:"to"`
	BasePrice  float64 `json:"base_price,omitempty"`
	PricePerKg float64 `json:"price_per_kg,omitempty"`
}

// VolumeRange - диапазон объема для расчета стоимости
type VolumeRange struct {
	From       float64 `json:"from"`
	To         float64 `json:"to"`
	PricePerM3 float64 `json:"price_per_m3"`
}

// ZoneMultiplier - множитель для зоны доставки
type ZoneMultiplier struct {
	Zone       string  `json:"zone"`
	Multiplier float64 `json:"multiplier"`
}

// DeliveryService - дополнительная услуга доставки
type DeliveryService struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	IsIncluded  bool    `json:"is_included"`
	IsAvailable bool    `json:"is_available"`
}

// CostBreakdown - детализация стоимости
type CostBreakdown struct {
	BasePrice           float64 `json:"base_price"`
	WeightSurcharge     float64 `json:"weight_surcharge,omitempty"`
	OversizeSurcharge   float64 `json:"oversize_surcharge,omitempty"`
	FragileSurcharge    float64 `json:"fragile_surcharge,omitempty"`
	InsuranceFee        float64 `json:"insurance_fee,omitempty"`
	CODFee              float64 `json:"cod_fee,omitempty"`
	FuelSurcharge       float64 `json:"fuel_surcharge,omitempty"`
	RemoteAreaSurcharge float64 `json:"remote_area_surcharge,omitempty"`
	Discount            float64 `json:"discount,omitempty"`
	Tax                 float64 `json:"tax,omitempty"`
	Total               float64 `json:"total"`
}

// CalculateTotal - рассчитывает итоговую стоимость
func (c *CostBreakdown) CalculateTotal() {
	c.Total = c.BasePrice +
		c.WeightSurcharge +
		c.OversizeSurcharge +
		c.FragileSurcharge +
		c.InsuranceFee +
		c.CODFee +
		c.FuelSurcharge +
		c.RemoteAreaSurcharge +
		c.Tax -
		c.Discount
}

// Value - реализация driver.Valuer для DeliveryAttributes
func (da DeliveryAttributes) Value() (driver.Value, error) {
	return json.Marshal(da)
}

// Scan - реализация sql.Scanner для DeliveryAttributes
func (da *DeliveryAttributes) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, da)
	case string:
		return json.Unmarshal([]byte(v), da)
	default:
		return json.Unmarshal(value.([]byte), da)
	}
}

// Дополнительные статусы отправления (к уже определенным выше)
const (
	ShipmentStatusConfirmed         = "confirmed"
	ShipmentStatusPickedUp          = "picked_up"
	ShipmentStatusOutForDelivery    = "out_for_delivery"
	ShipmentStatusDeliveryAttempted = "delivery_attempted"
	ShipmentStatusReturning         = "returning"
	ShipmentStatusReturned          = "returned"
	ShipmentStatusLost              = "lost"
	ShipmentStatusDamaged           = "damaged"
)

// RuleType - типы правил расчета стоимости
const (
	RuleTypeWeightBased = "weight_based"
	RuleTypeVolumeBased = "volume_based"
	RuleTypeZoneBased   = "zone_based"
	RuleTypeCombined    = "combined"
)
