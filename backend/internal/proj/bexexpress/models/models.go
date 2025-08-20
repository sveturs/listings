package models

import (
	"encoding/json"
	"time"
)

// BEXSettings представляет настройки интеграции с BEX Express
type BEXSettings struct {
	ID          int    `json:"id" db:"id"`
	AuthToken   string `json:"-" db:"auth_token"` // Скрываем токен в JSON
	ClientID    string `json:"client_id" db:"client_id"`
	APIEndpoint string `json:"api_endpoint" db:"api_endpoint"`

	// Данные отправителя (магазина)
	SenderClientID   string `json:"sender_client_id" db:"sender_client_id"`
	SenderName       string `json:"sender_name" db:"sender_name"`
	SenderAddress    string `json:"sender_address" db:"sender_address"`
	SenderCity       string `json:"sender_city" db:"sender_city"`
	SenderPostalCode string `json:"sender_postal_code" db:"sender_postal_code"`
	SenderPhone      string `json:"sender_phone" db:"sender_phone"`
	SenderEmail      string `json:"sender_email" db:"sender_email"`

	// Настройки
	Enabled            bool `json:"enabled" db:"enabled"`
	TestMode           bool `json:"test_mode" db:"test_mode"`
	AutoPrintLabels    bool `json:"auto_print_labels" db:"auto_print_labels"`
	AutoTrackShipments bool `json:"auto_track_shipments" db:"auto_track_shipments"`
	UseAddressLookup   bool `json:"use_address_lookup" db:"use_address_lookup"`

	// Уведомления
	NotifyOnPickup         bool `json:"notify_on_pickup" db:"notify_on_pickup"`
	NotifyOnDelivery       bool `json:"notify_on_delivery" db:"notify_on_delivery"`
	NotifyOnFailedDelivery bool `json:"notify_on_failed_delivery" db:"notify_on_failed_delivery"`

	// Статистика
	TotalShipments       int `json:"total_shipments" db:"total_shipments"`
	SuccessfulDeliveries int `json:"successful_deliveries" db:"successful_deliveries"`
	FailedDeliveries     int `json:"failed_deliveries" db:"failed_deliveries"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// BEXMunicipality представляет муниципалитет из справочника BEX
type BEXMunicipality struct {
	ID           int       `json:"id" db:"id"`
	BexID        int       `json:"bex_id" db:"bex_id"`
	Name         string    `json:"name" db:"name"`
	NameCyrillic string    `json:"name_cyrillic" db:"name_cyrillic"`
	Code         string    `json:"code" db:"code"`
	Region       string    `json:"region" db:"region"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BEXPlace представляет населенный пункт из справочника BEX
type BEXPlace struct {
	ID             int       `json:"id" db:"id"`
	BexID          string    `json:"bex_id" db:"bex_id"`
	MunicipalityID int       `json:"municipality_id" db:"municipality_id"`
	Name           string    `json:"name" db:"name"`
	NameCyrillic   string    `json:"name_cyrillic" db:"name_cyrillic"`
	PostalCode     string    `json:"postal_code" db:"postal_code"`
	Latitude       *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude      *float64  `json:"longitude,omitempty" db:"longitude"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// BEXStreet представляет улицу из справочника BEX
type BEXStreet struct {
	ID           int       `json:"id" db:"id"`
	BexID        string    `json:"bex_id" db:"bex_id"`
	PlaceID      int       `json:"place_id" db:"place_id"`
	Name         string    `json:"name" db:"name"`
	NameCyrillic string    `json:"name_cyrillic" db:"name_cyrillic"`
	StreetType   string    `json:"street_type" db:"street_type"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BEXParcelShop представляет пункт выдачи BEX
type BEXParcelShop struct {
	ID           int             `json:"id" db:"id"`
	BexID        int             `json:"bex_id" db:"bex_id"`
	Code         string          `json:"code" db:"code"`
	Name         string          `json:"name" db:"name"`
	Address      string          `json:"address" db:"address"`
	City         string          `json:"city" db:"city"`
	PostalCode   string          `json:"postal_code" db:"postal_code"`
	Phone        string          `json:"phone" db:"phone"`
	WorkingHours json.RawMessage `json:"working_hours" db:"working_hours"`
	Latitude     *float64        `json:"latitude,omitempty" db:"latitude"`
	Longitude    *float64        `json:"longitude,omitempty" db:"longitude"`
	IsActive     bool            `json:"is_active" db:"is_active"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// ShipmentStatus представляет статус отправления BEX
type ShipmentStatus int

const (
	ShipmentStatusNotFound         ShipmentStatus = 1
	ShipmentStatusDeleted          ShipmentStatus = 2
	ShipmentStatusNotSentYet       ShipmentStatus = 3
	ShipmentStatusPickedUp         ShipmentStatus = 4
	ShipmentStatusDelivered        ShipmentStatus = 5
	ShipmentStatusAddressVerify    ShipmentStatus = 6
	ShipmentStatusReturnToSender   ShipmentStatus = 7
	ShipmentStatusReturnedToSender ShipmentStatus = 8
	ShipmentStatusCODPaid          ShipmentStatus = 9
)

// BEXShipment представляет отправление BEX
type BEXShipment struct {
	ID                 int    `json:"id" db:"id"`
	MarketplaceOrderID *int   `json:"marketplace_order_id,omitempty" db:"marketplace_order_id"`
	StorefrontOrderID  *int64 `json:"storefront_order_id,omitempty" db:"storefront_order_id"`

	// BEX идентификаторы
	BexShipmentID  *int    `json:"bex_shipment_id,omitempty" db:"bex_shipment_id"`
	TrackingNumber *string `json:"tracking_number,omitempty" db:"tracking_number"`

	// Отправитель
	SenderName       string `json:"sender_name" db:"sender_name"`
	SenderAddress    string `json:"sender_address" db:"sender_address"`
	SenderCity       string `json:"sender_city" db:"sender_city"`
	SenderPostalCode string `json:"sender_postal_code" db:"sender_postal_code"`
	SenderPhone      string `json:"sender_phone" db:"sender_phone"`
	SenderEmail      string `json:"sender_email" db:"sender_email"`

	// Получатель
	RecipientName       string  `json:"recipient_name" db:"recipient_name"`
	RecipientAddress    string  `json:"recipient_address" db:"recipient_address"`
	RecipientCity       string  `json:"recipient_city" db:"recipient_city"`
	RecipientPostalCode string  `json:"recipient_postal_code" db:"recipient_postal_code"`
	RecipientPhone      string  `json:"recipient_phone" db:"recipient_phone"`
	RecipientEmail      *string `json:"recipient_email,omitempty" db:"recipient_email"`

	// Параметры посылки
	ShipmentType     int     `json:"shipment_type" db:"shipment_type"`
	ShipmentCategory int     `json:"shipment_category" db:"shipment_category"`
	ShipmentContents int     `json:"shipment_contents" db:"shipment_contents"`
	WeightKg         float64 `json:"weight_kg" db:"weight_kg"`
	TotalPackages    int     `json:"total_packages" db:"total_packages"`

	// Услуги
	PayType                  int      `json:"pay_type" db:"pay_type"`
	CODAmount                *float64 `json:"cod_amount,omitempty" db:"cod_amount"`
	InsuranceAmount          *float64 `json:"insurance_amount,omitempty" db:"insurance_amount"`
	PersonalDelivery         bool     `json:"personal_delivery" db:"personal_delivery"`
	ReturnSignedInvoices     bool     `json:"return_signed_invoices" db:"return_signed_invoices"`
	ReturnSignedConfirmation bool     `json:"return_signed_confirmation" db:"return_signed_confirmation"`
	ReturnPackage            bool     `json:"return_package" db:"return_package"`

	// Комментарии
	CommentPublic        *string `json:"comment_public,omitempty" db:"comment_public"`
	CommentPrivate       *string `json:"comment_private,omitempty" db:"comment_private"`
	DeliveryInstructions *string `json:"delivery_instructions,omitempty" db:"delivery_instructions"`

	// Статус
	Status       ShipmentStatus `json:"status" db:"status"`
	StatusText   *string        `json:"status_text,omitempty" db:"status_text"`
	FailedReason *string        `json:"failed_reason,omitempty" db:"failed_reason"`

	// Документы
	LabelBase64 *string `json:"label_base64,omitempty" db:"label_base64"`
	LabelURL    *string `json:"label_url,omitempty" db:"label_url"`

	// Временные метки
	RegisteredAt *time.Time `json:"registered_at,omitempty" db:"registered_at"`
	PickedUpAt   *time.Time `json:"picked_up_at,omitempty" db:"picked_up_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	FailedAt     *time.Time `json:"failed_at,omitempty" db:"failed_at"`
	ReturnedAt   *time.Time `json:"returned_at,omitempty" db:"returned_at"`

	// История статусов
	StatusHistory json.RawMessage `json:"status_history,omitempty" db:"status_history"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// API Request/Response структуры

// CreateShipmentRequest представляет запрос на создание отправления BEX
type CreateShipmentRequest struct {
	OrderID           *int   `json:"order_id,omitempty"`
	StorefrontOrderID *int64 `json:"storefront_order_id,omitempty"`

	RecipientName       string  `json:"recipient_name" validate:"required"`
	RecipientAddress    string  `json:"recipient_address" validate:"required"`
	RecipientCity       string  `json:"recipient_city" validate:"required"`
	RecipientPostalCode string  `json:"recipient_postal_code" validate:"required"`
	RecipientPhone      string  `json:"recipient_phone" validate:"required"`
	RecipientEmail      *string `json:"recipient_email,omitempty"`

	WeightKg         float64 `json:"weight_kg" validate:"required,min=0.1,max=50"`
	TotalPackages    int     `json:"total_packages" validate:"required,min=1,max=99"`
	ShipmentCategory int     `json:"shipment_category" validate:"required"`
	ShipmentContents int     `json:"shipment_contents" validate:"required"`

	CODAmount            *float64 `json:"cod_amount,omitempty"`
	InsuranceAmount      *float64 `json:"insurance_amount,omitempty"`
	PersonalDelivery     bool     `json:"personal_delivery,omitempty"`
	DeliveryInstructions *string  `json:"delivery_instructions,omitempty"`
	Notes                *string  `json:"notes,omitempty"`

	// Уведомления
	PreNotificationMinutes int `json:"pre_notification_minutes,omitempty"`
}

// BEXShipmentTask представляет задачу (pickup/delivery) для API BEX
type BEXShipmentTask struct {
	Type            int    `json:"type"`
	NameType        int    `json:"nameType"`
	Name1           string `json:"name1"`
	Name2           string `json:"name2"`
	TaxID           string `json:"taxId"`
	AddressType     int    `json:"adressType"`
	Municipalities  int    `json:"municipalities"`
	Place           string `json:"place"`
	Street          string `json:"street"`
	HouseNumber     int    `json:"houseNumber"`
	Apartment       string `json:"apartment"`
	ContactPerson   string `json:"contactPerson"`
	Phone           string `json:"phone"`
	Date            string `json:"date"`
	TimeFrom        string `json:"timeFrom"`
	TimeTo          string `json:"timeTo"`
	PreNotification int    `json:"preNotification"`
	Comment         string `json:"comment"`
	ParcelShop      int    `json:"parcelShop"`
}

// BEXShipmentData представляет данные отправления для API BEX
type BEXShipmentData struct {
	ShipmentID               int               `json:"shipmentId"`
	ServiceSpeed             int               `json:"serviceSpeed"`
	ShipmentType             int               `json:"shipmentType"`
	ShipmentCategory         int               `json:"shipmentCategory"`
	ShipmentWeight           float64           `json:"shipmentWeight"`
	TotalPackages            int               `json:"totalPackages"`
	InvoiceAmount            float64           `json:"invoiceAmount"`
	ShipmentContents         int               `json:"shipmentContents"`
	CommentPublic            string            `json:"commentPublic"`
	CommentPrivate           string            `json:"commentPrivate"`
	PersonalDelivery         bool              `json:"personalDelivery"`
	ReturnSignedInvoices     bool              `json:"returnSignedInvoices"`
	ReturnSignedConfirmation bool              `json:"returnSignedConfirmation"`
	ReturnPackage            bool              `json:"returnPackage"`
	PayType                  int               `json:"payType"`
	InsuranceAmount          float64           `json:"insuranceAmount"`
	PayToSender              float64           `json:"payToSender"`
	PayToSenderViaAccount    bool              `json:"payToSenderViaAccount"`
	SendersAccountNumber     string            `json:"sendersAccountNumber"`
	BankTransferComment      string            `json:"bankTransferComment"`
	Tasks                    []BEXShipmentTask `json:"tasks"`
	Reports                  []interface{}     `json:"reports"`
}

// BEXShipmentRequest представляет запрос к API BEX
type BEXShipmentRequest struct {
	ShipmentsList []BEXShipmentData `json:"shipmentslist"`
}

// BEXShipmentResult представляет результат создания отправления
type BEXShipmentResult struct {
	State      bool   `json:"state"`
	ShipmentID int    `json:"shipmentId"`
	Error      string `json:"err"`
}

// BEXShipmentResponse представляет ответ от API BEX
type BEXShipmentResponse struct {
	ShipmentsResultList []BEXShipmentResult `json:"shipmentsResultList"`
	RequestState        bool                `json:"reqstate"`
	RequestError        string              `json:"reqerr"`
}

// BEXLabelRequest представляет запрос на получение этикетки
type BEXLabelRequest struct {
	PageSize     int `json:"pageSize"`
	PagePosition int `json:"pagePosition"`
	ShipmentID   int `json:"shipmentId"`
	ParcelNo     int `json:"parcelNo"`
}

// BEXLabelResponse представляет ответ с этикеткой
type BEXLabelResponse struct {
	State       bool   `json:"state"`
	ShipmentID  int    `json:"shipmentId"`
	ParcelNo    int    `json:"parcelNo"`
	ParcelLabel string `json:"parcelLabel"`
	Error       string `json:"err"`
}

// BEXStatusRequest представляет запрос статуса отправления
type BEXStatusRequest struct {
	ShipmentID int `json:"shipmentid"`
	MType      int `json:"mtype"`
	Lang       int `json:"lang"`
}

// BEXStatusResponse представляет ответ со статусом
type BEXStatusResponse struct {
	Status        int     `json:"status"`
	StatusText    string  `json:"status_text"`
	DateSent      string  `json:"date_sent"`
	DateDelivered string  `json:"date_delivered"`
	CODAmount     float64 `json:"cod_amount"`
	Notes         string  `json:"notes"`
}

// SearchAddressRequest представляет запрос поиска адреса
type SearchAddressRequest struct {
	Query          string `json:"query" validate:"required,min=2"`
	City           string `json:"city,omitempty"`
	MunicipalityID *int   `json:"municipality_id,omitempty"`
	PlaceID        *int   `json:"place_id,omitempty"`
	Limit          int    `json:"limit,omitempty"`
}

// AddressSuggestion представляет предложение адреса
type AddressSuggestion struct {
	PlaceID        string   `json:"place_id"`
	PlaceName      string   `json:"place_name"`
	StreetID       string   `json:"street_id"`
	StreetName     string   `json:"street_name"`
	PostalCode     string   `json:"postal_code"`
	MunicipalityID int      `json:"municipality_id"`
	Municipality   string   `json:"municipality"`
	FullAddress    string   `json:"full_address"`
	Latitude       *float64 `json:"latitude,omitempty"`
	Longitude      *float64 `json:"longitude,omitempty"`
}

// CalculateRateRequest представляет запрос расчета стоимости
type CalculateRateRequest struct {
	RecipientPostalCode string   `json:"recipient_postal_code" validate:"required"`
	WeightKg            float64  `json:"weight_kg" validate:"required,min=0.1,max=50"`
	ShipmentCategory    int      `json:"shipment_category" validate:"required"`
	CODAmount           *float64 `json:"cod_amount,omitempty"`
	InsuranceAmount     *float64 `json:"insurance_amount,omitempty"`
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
