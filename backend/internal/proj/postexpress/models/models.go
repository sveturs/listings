package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PostExpressSettings представляет настройки интеграции с Post Express
type PostExpressSettings struct {
	ID          int    `json:"id" db:"id"`
	APIUsername string `json:"api_username" db:"api_username"`
	APIPassword string `json:"-" db:"api_password"` // Скрываем пароль в JSON
	APIEndpoint string `json:"api_endpoint" db:"api_endpoint"`

	SenderName       string `json:"sender_name" db:"sender_name"`
	SenderAddress    string `json:"sender_address" db:"sender_address"`
	SenderCity       string `json:"sender_city" db:"sender_city"`
	SenderPostalCode string `json:"sender_postal_code" db:"sender_postal_code"`
	SenderPhone      string `json:"sender_phone" db:"sender_phone"`
	SenderEmail      string `json:"sender_email" db:"sender_email"`

	Enabled            bool `json:"enabled" db:"enabled"`
	TestMode           bool `json:"test_mode" db:"test_mode"`
	AutoPrintLabels    bool `json:"auto_print_labels" db:"auto_print_labels"`
	AutoTrackShipments bool `json:"auto_track_shipments" db:"auto_track_shipments"`

	NotifyOnPickup         bool `json:"notify_on_pickup" db:"notify_on_pickup"`
	NotifyOnDelivery       bool `json:"notify_on_delivery" db:"notify_on_delivery"`
	NotifyOnFailedDelivery bool `json:"notify_on_failed_delivery" db:"notify_on_failed_delivery"`

	TotalShipments       int `json:"total_shipments" db:"total_shipments"`
	SuccessfulDeliveries int `json:"successful_deliveries" db:"successful_deliveries"`
	FailedDeliveries     int `json:"failed_deliveries" db:"failed_deliveries"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PostExpressLocation представляет населенный пункт
type PostExpressLocation struct {
	ID            int     `json:"id" db:"id"`
	PostExpressID int     `json:"post_express_id" db:"post_express_id"`
	Name          string  `json:"name" db:"name"`
	NameCyrillic  *string `json:"name_cyrillic,omitempty" db:"name_cyrillic"`
	PostalCode    *string `json:"postal_code,omitempty" db:"postal_code"`
	Municipality  *string `json:"municipality,omitempty" db:"municipality"`

	Latitude  *float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude *float64 `json:"longitude,omitempty" db:"longitude"`

	Region       *string `json:"region,omitempty" db:"region"`
	District     *string `json:"district,omitempty" db:"district"`
	DeliveryZone *string `json:"delivery_zone,omitempty" db:"delivery_zone"`

	IsActive        bool `json:"is_active" db:"is_active"`
	SupportsCOD     bool `json:"supports_cod" db:"supports_cod"`
	SupportsExpress bool `json:"supports_express" db:"supports_express"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PostExpressOffice представляет почтовое отделение
type PostExpressOffice struct {
	ID         int    `json:"id" db:"id"`
	OfficeCode string `json:"office_code" db:"office_code"`
	LocationID *int   `json:"location_id,omitempty" db:"location_id"`

	Name    string  `json:"name" db:"name"`
	Address string  `json:"address" db:"address"`
	Phone   *string `json:"phone,omitempty" db:"phone"`
	Email   *string `json:"email,omitempty" db:"email"`

	WorkingHours json.RawMessage `json:"working_hours,omitempty" db:"working_hours"`

	Latitude  *float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude *float64 `json:"longitude,omitempty" db:"longitude"`

	AcceptsPackages      bool `json:"accepts_packages" db:"accepts_packages"`
	IssuesPackages       bool `json:"issues_packages" db:"issues_packages"`
	HasATM               bool `json:"has_atm" db:"has_atm"`
	HasParking           bool `json:"has_parking" db:"has_parking"`
	WheelchairAccessible bool `json:"wheelchair_accessible" db:"wheelchair_accessible"`

	IsActive        bool       `json:"is_active" db:"is_active"`
	TemporaryClosed bool       `json:"temporary_closed" db:"temporary_closed"`
	ClosedUntil     *time.Time `json:"closed_until,omitempty" db:"closed_until"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PostExpressRate представляет тариф доставки
type PostExpressRate struct {
	ID         int     `json:"id" db:"id"`
	WeightFrom float64 `json:"weight_from" db:"weight_from"`
	WeightTo   float64 `json:"weight_to" db:"weight_to"`
	BasePrice  float64 `json:"base_price" db:"base_price"`

	InsuranceIncludedUpTo float64 `json:"insurance_included_up_to" db:"insurance_included_up_to"`
	InsuranceRatePercent  float64 `json:"insurance_rate_percent" db:"insurance_rate_percent"`
	CODFee                float64 `json:"cod_fee" db:"cod_fee"`

	MaxLengthCm        int `json:"max_length_cm" db:"max_length_cm"`
	MaxWidthCm         int `json:"max_width_cm" db:"max_width_cm"`
	MaxHeightCm        int `json:"max_height_cm" db:"max_height_cm"`
	MaxDimensionsSumCm int `json:"max_dimensions_sum_cm" db:"max_dimensions_sum_cm"`

	DeliveryDaysMin int `json:"delivery_days_min" db:"delivery_days_min"`
	DeliveryDaysMax int `json:"delivery_days_max" db:"delivery_days_max"`

	IsActive       bool `json:"is_active" db:"is_active"`
	IsSpecialOffer bool `json:"is_special_offer" db:"is_special_offer"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ShipmentStatus представляет статус отправления
type ShipmentStatus string

const (
	ShipmentStatusCreated    ShipmentStatus = "created"
	ShipmentStatusRegistered ShipmentStatus = "registered"
	ShipmentStatusPickedUp   ShipmentStatus = "picked_up"
	ShipmentStatusInTransit  ShipmentStatus = "in_transit"
	ShipmentStatusDelivered  ShipmentStatus = "delivered"
	ShipmentStatusFailed     ShipmentStatus = "failed"
	ShipmentStatusReturned   ShipmentStatus = "returned"
)

// PostExpressShipment представляет отправление
type PostExpressShipment struct {
	ID                 int    `json:"id" db:"id"`
	MarketplaceOrderID *int   `json:"marketplace_order_id,omitempty" db:"marketplace_order_id"`
	StorefrontOrderID  *int64 `json:"storefront_order_id,omitempty" db:"storefront_order_id"`

	TrackingNumber *string `json:"tracking_number,omitempty" db:"tracking_number"`
	Barcode        *string `json:"barcode,omitempty" db:"barcode"`
	PostExpressID  *string `json:"post_express_id,omitempty" db:"post_express_id"`

	// Отправитель
	SenderName       string  `json:"sender_name" db:"sender_name"`
	SenderAddress    string  `json:"sender_address" db:"sender_address"`
	SenderCity       string  `json:"sender_city" db:"sender_city"`
	SenderPostalCode string  `json:"sender_postal_code" db:"sender_postal_code"`
	SenderPhone      string  `json:"sender_phone" db:"sender_phone"`
	SenderEmail      *string `json:"sender_email,omitempty" db:"sender_email"`
	SenderLocationID *int    `json:"sender_location_id,omitempty" db:"sender_location_id"`

	// Получатель
	RecipientName       string  `json:"recipient_name" db:"recipient_name"`
	RecipientAddress    string  `json:"recipient_address" db:"recipient_address"`
	RecipientCity       string  `json:"recipient_city" db:"recipient_city"`
	RecipientPostalCode string  `json:"recipient_postal_code" db:"recipient_postal_code"`
	RecipientPhone      string  `json:"recipient_phone" db:"recipient_phone"`
	RecipientEmail      *string `json:"recipient_email,omitempty" db:"recipient_email"`
	RecipientLocationID *int    `json:"recipient_location_id,omitempty" db:"recipient_location_id"`

	// Параметры посылки
	WeightKg      float64  `json:"weight_kg" db:"weight_kg"`
	LengthCm      *int     `json:"length_cm,omitempty" db:"length_cm"`
	WidthCm       *int     `json:"width_cm,omitempty" db:"width_cm"`
	HeightCm      *int     `json:"height_cm,omitempty" db:"height_cm"`
	DeclaredValue *float64 `json:"declared_value,omitempty" db:"declared_value"`

	// Услуги
	ServiceType     string   `json:"service_type" db:"service_type"`
	CODAmount       *float64 `json:"cod_amount,omitempty" db:"cod_amount"`
	CODReference    *string  `json:"cod_reference,omitempty" db:"cod_reference"`
	InsuranceAmount *float64 `json:"insurance_amount,omitempty" db:"insurance_amount"`

	// Расчет стоимости
	BasePrice    float64 `json:"base_price" db:"base_price"`
	InsuranceFee float64 `json:"insurance_fee" db:"insurance_fee"`
	CODFee       float64 `json:"cod_fee" db:"cod_fee"`
	TotalPrice   float64 `json:"total_price" db:"total_price"`

	// Статусы
	Status         ShipmentStatus `json:"status" db:"status"`
	DeliveryStatus *string        `json:"delivery_status,omitempty" db:"delivery_status"`

	// Документы
	LabelURL   *string `json:"label_url,omitempty" db:"label_url"`
	InvoiceURL *string `json:"invoice_url,omitempty" db:"invoice_url"`
	PODURL     *string `json:"pod_url,omitempty" db:"pod_url"`

	// Временные метки
	RegisteredAt *time.Time `json:"registered_at,omitempty" db:"registered_at"`
	PickedUpAt   *time.Time `json:"picked_up_at,omitempty" db:"picked_up_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	FailedAt     *time.Time `json:"failed_at,omitempty" db:"failed_at"`
	ReturnedAt   *time.Time `json:"returned_at,omitempty" db:"returned_at"`

	// История статусов
	StatusHistory json.RawMessage `json:"status_history,omitempty" db:"status_history"`

	// Дополнительная информация
	Notes                *string `json:"notes,omitempty" db:"notes"`
	InternalNotes        *string `json:"internal_notes,omitempty" db:"internal_notes"`
	DeliveryInstructions *string `json:"delivery_instructions,omitempty" db:"delivery_instructions"`
	FailedReason         *string `json:"failed_reason,omitempty" db:"failed_reason"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateShipmentRequest представляет запрос на создание отправления
type CreateShipmentRequest struct {
	OrderID           *int   `json:"order_id,omitempty"`
	StorefrontOrderID *int64 `json:"storefront_order_id,omitempty"`

	RecipientName       string  `json:"recipient_name" validate:"required"`
	RecipientAddress    string  `json:"recipient_address" validate:"required"`
	RecipientCity       string  `json:"recipient_city" validate:"required"`
	RecipientPostalCode string  `json:"recipient_postal_code" validate:"required"`
	RecipientPhone      string  `json:"recipient_phone" validate:"required"`
	RecipientEmail      *string `json:"recipient_email,omitempty"`

	WeightKg      float64  `json:"weight_kg" validate:"required,min=0.1,max=20"`
	LengthCm      *int     `json:"length_cm,omitempty"`
	WidthCm       *int     `json:"width_cm,omitempty"`
	HeightCm      *int     `json:"height_cm,omitempty"`
	DeclaredValue *float64 `json:"declared_value,omitempty"`

	ServiceType          string   `json:"service_type,omitempty"`
	CODAmount            *float64 `json:"cod_amount,omitempty"`
	InsuranceAmount      *float64 `json:"insurance_amount,omitempty"`
	DeliveryInstructions *string  `json:"delivery_instructions,omitempty"`
	Notes                *string  `json:"notes,omitempty"`
}

// CalculateRateRequest представляет запрос на расчет стоимости
type CalculateRateRequest struct {
	SenderPostalCode    string   `json:"sender_postal_code" validate:"required"`
	RecipientPostalCode string   `json:"recipient_postal_code" validate:"required"`
	WeightKg            float64  `json:"weight_kg" validate:"required,min=0.1,max=20"`
	LengthCm            *int     `json:"length_cm,omitempty"`
	WidthCm             *int     `json:"width_cm,omitempty"`
	HeightCm            *int     `json:"height_cm,omitempty"`
	DeclaredValue       *float64 `json:"declared_value,omitempty"`
	CODAmount           *float64 `json:"cod_amount,omitempty"`
	ServiceType         string   `json:"service_type,omitempty"`
}

// CalculateRateResponse представляет ответ с расчетом стоимости
type CalculateRateResponse struct {
	BasePrice        float64 `json:"base_price"`
	InsuranceFee     float64 `json:"insurance_fee"`
	CODFee           float64 `json:"cod_fee"`
	TotalPrice       float64 `json:"total_price"`
	DeliveryDaysMin  int     `json:"delivery_days_min"`
	DeliveryDaysMax  int     `json:"delivery_days_max"`
	ServiceAvailable bool    `json:"service_available"`
}

// TrackingEvent представляет событие отслеживания
type TrackingEvent struct {
	ID               int             `json:"id" db:"id"`
	ShipmentID       int             `json:"shipment_id" db:"shipment_id"`
	EventCode        string          `json:"event_code" db:"event_code"`
	EventDescription string          `json:"event_description" db:"event_description"`
	EventLocation    *string         `json:"event_location,omitempty" db:"event_location"`
	EventTimestamp   time.Time       `json:"event_timestamp" db:"event_timestamp"`
	AdditionalInfo   json.RawMessage `json:"additional_info,omitempty" db:"additional_info"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// WSP API структуры

// TransactionRequest представляет запрос к WSP API
type TransactionRequest struct {
	TransactionType int    `json:"transaction_type"`
	InputData       string `json:"input_data"`
}

// TransactionResponse представляет ответ от WSP API
type TransactionResponse struct {
	Success      bool            `json:"success"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	OutputData   json.RawMessage `json:"output_data,omitempty"`
}

// ClientData представляет данные клиента для WSP API
type ClientData struct {
	Username          string  `json:"Username"`
	Password          string  `json:"Password"`
	Jezik             string  `json:"Jezik"`
	IdTipUredjaja     int     `json:"IdTipUredjaja"`
	VerzijaOS         string  `json:"VerzijaOS"`
	NazivUredjaja     string  `json:"NazivUredjaja"`
	ModelUredjaja     string  `json:"ModelUredjaja"`
	VerzijaAplikacije string  `json:"VerzijaAplikacije"`
	IPAdresa          string  `json:"IPAdresa"`
	Geolokacija       *string `json:"Geolokacija,omitempty"`
}

// TransakcijaIn представляет входящий запрос к WSP API
type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTranskacije int    `json:"IdVrstaTranskacije"`
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn"`
}

// GenerateGUID генерирует новый GUID для транзакции
func GenerateGUID() string {
	return uuid.New().String()
}

// Warehouse модели

// Warehouse представляет склад
type Warehouse struct {
	ID               int             `json:"id" db:"id"`
	Code             string          `json:"code" db:"code"`
	Name             string          `json:"name" db:"name"`
	Type             string          `json:"type" db:"type"`
	Address          string          `json:"address" db:"address"`
	City             string          `json:"city" db:"city"`
	PostalCode       string          `json:"postal_code" db:"postal_code"`
	Country          string          `json:"country" db:"country"`
	Phone            *string         `json:"phone,omitempty" db:"phone"`
	Email            *string         `json:"email,omitempty" db:"email"`
	ManagerName      *string         `json:"manager_name,omitempty" db:"manager_name"`
	ManagerPhone     *string         `json:"manager_phone,omitempty" db:"manager_phone"`
	Latitude         *float64        `json:"latitude,omitempty" db:"latitude"`
	Longitude        *float64        `json:"longitude,omitempty" db:"longitude"`
	WorkingHours     json.RawMessage `json:"working_hours,omitempty" db:"working_hours"`
	TotalAreaM2      *float64        `json:"total_area_m2,omitempty" db:"total_area_m2"`
	StorageAreaM2    *float64        `json:"storage_area_m2,omitempty" db:"storage_area_m2"`
	MaxCapacityM3    *float64        `json:"max_capacity_m3,omitempty" db:"max_capacity_m3"`
	CurrentOccupancy *float64        `json:"current_occupancy_m3,omitempty" db:"current_occupancy_m3"`
	SupportsFBS      bool            `json:"supports_fbs" db:"supports_fbs"`
	SupportsPickup   bool            `json:"supports_pickup" db:"supports_pickup"`
	HasRefrigeration bool            `json:"has_refrigeration" db:"has_refrigeration"`
	HasLoadingDock   bool            `json:"has_loading_dock" db:"has_loading_dock"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// PickupOrderStatus представляет статус заказа на самовывоз
type PickupOrderStatus string

const (
	PickupOrderStatusPending  PickupOrderStatus = "pending"
	PickupOrderStatusReady    PickupOrderStatus = "ready"
	PickupOrderStatusPickedUp PickupOrderStatus = "picked_up"
	PickupOrderStatusExpired  PickupOrderStatus = "expired"
	PickupOrderStatusCanceled PickupOrderStatus = "canceled"
)

// WarehousePickupOrder представляет заказ на самовывоз
type WarehousePickupOrder struct {
	ID                 int               `json:"id" db:"id"`
	WarehouseID        int               `json:"warehouse_id" db:"warehouse_id"`
	MarketplaceOrderID *int              `json:"marketplace_order_id,omitempty" db:"marketplace_order_id"`
	StorefrontOrderID  *int64            `json:"storefront_order_id,omitempty" db:"storefront_order_id"`
	PickupCode         string            `json:"pickup_code" db:"pickup_code"`
	QRCodeURL          *string           `json:"qr_code_url,omitempty" db:"qr_code_url"`
	Status             PickupOrderStatus `json:"status" db:"status"`
	ReadyAt            *time.Time        `json:"ready_at,omitempty" db:"ready_at"`
	PickedUpAt         *time.Time        `json:"picked_up_at,omitempty" db:"picked_up_at"`
	ExpiresAt          *time.Time        `json:"expires_at,omitempty" db:"expires_at"`
	CustomerName       string            `json:"customer_name" db:"customer_name"`
	CustomerPhone      string            `json:"customer_phone" db:"customer_phone"`
	CustomerEmail      *string           `json:"customer_email,omitempty" db:"customer_email"`
	PickupConfirmedBy  *string           `json:"pickup_confirmed_by,omitempty" db:"pickup_confirmed_by"`
	IDDocumentType     *string           `json:"id_document_type,omitempty" db:"id_document_type"`
	IDDocumentNumber   *string           `json:"id_document_number,omitempty" db:"id_document_number"`
	SignatureURL       *string           `json:"signature_url,omitempty" db:"signature_url"`
	NotificationSentAt *time.Time        `json:"notification_sent_at,omitempty" db:"notification_sent_at"`
	ReminderSentAt     *time.Time        `json:"reminder_sent_at,omitempty" db:"reminder_sent_at"`
	Notes              *string           `json:"notes,omitempty" db:"notes"`
	PickupPhotoURL     *string           `json:"pickup_photo_url,omitempty" db:"pickup_photo_url"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at"`
}

// CreatePickupOrderRequest представляет запрос на создание заказа на самовывоз
type CreatePickupOrderRequest struct {
	OrderID           *int    `json:"order_id,omitempty"`
	StorefrontOrderID *int64  `json:"storefront_order_id,omitempty"`
	CustomerName      string  `json:"customer_name" validate:"required"`
	CustomerPhone     string  `json:"customer_phone" validate:"required"`
	CustomerEmail     *string `json:"customer_email,omitempty"`
	Notes             *string `json:"notes,omitempty"`
}
