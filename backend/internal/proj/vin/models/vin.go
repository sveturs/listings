package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// VINDecodeCache представляет кэшированную декодированную информацию VIN
type VINDecodeCache struct {
	ID  int64  `json:"id" db:"id"`
	VIN string `json:"vin" db:"vin"`

	// Основная информация
	Make               *string `json:"make" db:"make"`
	Model              *string `json:"model" db:"model"`
	Year               *int    `json:"year" db:"year"`
	EngineType         *string `json:"engine_type" db:"engine_type"`
	EngineDisplacement *string `json:"engine_displacement" db:"engine_displacement"`
	TransmissionType   *string `json:"transmission_type" db:"transmission_type"`
	Drivetrain         *string `json:"drivetrain" db:"drivetrain"`
	BodyType           *string `json:"body_type" db:"body_type"`
	FuelType           *string `json:"fuel_type" db:"fuel_type"`

	// Технические характеристики
	Doors         *int    `json:"doors" db:"doors"`
	Seats         *int    `json:"seats" db:"seats"`
	ColorExterior *string `json:"color_exterior" db:"color_exterior"`
	ColorInterior *string `json:"color_interior" db:"color_interior"`

	// Информация о производителе
	Manufacturer    *string `json:"manufacturer" db:"manufacturer"`
	CountryOfOrigin *string `json:"country_of_origin" db:"country_of_origin"`
	AssemblyPlant   *string `json:"assembly_plant" db:"assembly_plant"`

	// Дополнительная информация
	VehicleClass       *string `json:"vehicle_class" db:"vehicle_class"`
	VehicleType        *string `json:"vehicle_type" db:"vehicle_type"`
	GrossVehicleWeight *string `json:"gross_vehicle_weight" db:"gross_vehicle_weight"`

	// Результат декодирования
	DecodeStatus string          `json:"decode_status" db:"decode_status"` // success, partial, failed
	ErrorMessage *string         `json:"error_message" db:"error_message"`
	RawResponse  json.RawMessage `json:"raw_response" db:"raw_response"`

	// Метаданные
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// VINCheckHistory представляет историю проверок VIN
type VINCheckHistory struct {
	ID        int64  `json:"id" db:"id"`
	UserID    *int64 `json:"user_id" db:"user_id"`
	VIN       string `json:"vin" db:"vin"`
	ListingID *int64 `json:"listing_id" db:"listing_id"`

	// Результаты проверки
	DecodeSuccess bool   `json:"decode_success" db:"decode_success"`
	DecodeCacheID *int64 `json:"decode_cache_id" db:"decode_cache_id"`

	// Дополнительная информация
	CheckType string  `json:"check_type" db:"check_type"` // manual, auto_fill, verification
	IPAddress *string `json:"ip_address" db:"ip_address"`
	UserAgent *string `json:"user_agent" db:"user_agent"`

	CheckedAt time.Time `json:"checked_at" db:"checked_at"`

	// Связанные данные (для ответов API)
	DecodeCache *VINDecodeCache `json:"decode_cache,omitempty" db:"-"`
}

// VINDecodeRequest представляет запрос на декодирование VIN
type VINDecodeRequest struct {
	VIN       string `json:"vin" validate:"required,len=17"`
	ListingID *int64 `json:"listing_id,omitempty"`
}

// VINDecodeResponse представляет ответ декодирования VIN
type VINDecodeResponse struct {
	Success bool            `json:"success"`
	Data    *VINDecodeCache `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
	Source  string          `json:"source"` // cache, api
}

// VINHistoryRequest представляет запрос истории проверок
type VINHistoryRequest struct {
	UserID *int64 `json:"user_id,omitempty"`
	VIN    string `json:"vin,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

// JSONB type for PostgreSQL jsonb columns
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, j)
}

// BasicVINInfo содержит базовую информацию, декодируемую из VIN локально
type BasicVINInfo struct {
	Year         int    `json:"year"`
	Manufacturer string `json:"manufacturer"`
	Region       string `json:"region"`
	CheckDigit   bool   `json:"check_digit_valid"`
}

// VINValidationError представляет ошибку валидации VIN
type VINValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
