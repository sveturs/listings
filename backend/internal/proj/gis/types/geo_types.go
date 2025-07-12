package types

import (
	"math"
	"time"
)

// Point представляет географическую точку
type Point struct {
	Lat float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng float64 `json:"lng" validate:"required,min=-180,max=180"`
}

// Bounds представляет географические границы
type Bounds struct {
	North float64 `json:"north" validate:"required,min=-90,max=90"`
	South float64 `json:"south" validate:"required,min=-90,max=90"`
	East  float64 `json:"east" validate:"required,min=-180,max=180"`
	West  float64 `json:"west" validate:"required,min=-180,max=180"`
}

// GeoListing представляет объявление с геоданными
type GeoListing struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Location    Point     `json:"location"`
	Address     string    `json:"address,omitempty"`
	Images      []string  `json:"images,omitempty"`
	UserID      int       `json:"user_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Distance    *float64  `json:"distance,omitempty"` // Расстояние от точки поиска (в метрах)
}

// SearchParams параметры пространственного поиска
type SearchParams struct {
	Bounds      *Bounds  `json:"bounds,omitempty"`
	Center      *Point   `json:"center,omitempty"`    // Центр поиска для радиусного поиска
	RadiusKm    float64  `json:"radius_km,omitempty"` // Радиус поиска в километрах
	Categories  []string `json:"categories,omitempty"`
	CategoryIDs []int    `json:"category_ids,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	SearchQuery string   `json:"q,omitempty"`          // Текстовый поиск
	SortBy      string   `json:"sort_by,omitempty"`    // distance, price, created_at
	SortOrder   string   `json:"sort_order,omitempty"` // asc, desc
	Limit       int      `json:"limit,omitempty"`
	Offset      int      `json:"offset,omitempty"`
	UserID      *int     `json:"user_id,omitempty"` // Фильтр по пользователю
	Status      string   `json:"status,omitempty"`  // Фильтр по статусу
}

// SearchResponse ответ на поисковый запрос
type SearchResponse struct {
	Listings   []GeoListing `json:"listings"`
	TotalCount int64        `json:"total_count"`
	HasMore    bool         `json:"has_more"`
}

// Validate проверяет корректность границ
func (b *Bounds) Validate() error {
	if b.North < b.South {
		return ErrInvalidBounds
	}
	if b.North > 90 || b.South < -90 {
		return ErrInvalidLatitude
	}
	if b.East > 180 || b.West < -180 {
		return ErrInvalidLongitude
	}
	return nil
}

// Contains проверяет, находится ли точка внутри границ
func (b *Bounds) Contains(p Point) bool {
	latOk := p.Lat >= b.South && p.Lat <= b.North

	// Обработка случая пересечения меридиана 180°
	lngOk := false
	if b.West <= b.East {
		lngOk = p.Lng >= b.West && p.Lng <= b.East
	} else {
		lngOk = p.Lng >= b.West || p.Lng <= b.East
	}

	return latOk && lngOk
}

// Expand расширяет границы на заданное расстояние в километрах
func (b *Bounds) Expand(distanceKm float64) {
	// Примерное преобразование км в градусы
	// 1 градус широты ≈ 111 км
	latDelta := distanceKm / 111.0

	// Для долготы зависит от широты
	avgLat := (b.North + b.South) / 2
	lngDelta := distanceKm / (111.0 * cosine(avgLat))

	b.North = min(90, b.North+latDelta)
	b.South = max(-90, b.South-latDelta)
	b.East = min(180, b.East+lngDelta)
	b.West = max(-180, b.West-lngDelta)
}

// Helper функции
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// UpdateLocationRequest запрос на обновление локации
type UpdateLocationRequest struct {
	Lat     float64 `json:"lat" validate:"required,min=-90,max=90"`
	Lng     float64 `json:"lng" validate:"required,min=-180,max=180"`
	Address string  `json:"address" validate:"required,min=5,max=500"`
}

// ========== PHASE 2: Новые типы для умного ввода адресов ==========

// AddressComponents компоненты структурированного адреса
type AddressComponents struct {
	Country     string `json:"country" db:"country"`
	CountryCode string `json:"country_code" db:"country_code"`
	City        string `json:"city" db:"city"`
	District    string `json:"district,omitempty" db:"district"`
	Street      string `json:"street,omitempty" db:"street"`
	HouseNumber string `json:"house_number,omitempty" db:"house_number"`
	PostalCode  string `json:"postal_code,omitempty" db:"postal_code"`
	Formatted   string `json:"formatted" db:"formatted"`
}

// LocationPrivacyLevel уровни приватности локации
type LocationPrivacyLevel string

const (
	PrivacyExact    LocationPrivacyLevel = "exact"    // Точный адрес
	PrivacyStreet   LocationPrivacyLevel = "street"   // Размытие ±100-200м
	PrivacyDistrict LocationPrivacyLevel = "district" // Размытие ±500-1000м
	PrivacyCity     LocationPrivacyLevel = "city"     // Только город
)

// InputMethod способы ввода адреса
type InputMethod string

const (
	InputManual          InputMethod = "manual"
	InputGeocoded        InputMethod = "geocoded"
	InputMapClick        InputMethod = "map_click"
	InputCurrentLocation InputMethod = "current_location"
)

// EnhancedListingGeo расширенная модель геоданных объявления
type EnhancedListingGeo struct {
	ID                  int64                `json:"id" db:"id"`
	ListingID          int64                `json:"listing_id" db:"listing_id"`
	Location           Point                `json:"location"`
	BlurredLocation    *Point               `json:"blurred_location,omitempty"`
	Geohash            string               `json:"geohash" db:"geohash"`
	IsPrecise          bool                 `json:"is_precise" db:"is_precise"`
	BlurRadius         float64              `json:"blur_radius" db:"blur_radius"`
	AddressComponents  *AddressComponents   `json:"address_components,omitempty"`
	FormattedAddress   string               `json:"formatted_address" db:"formatted_address"`
	GeocodingConfidence float64             `json:"geocoding_confidence" db:"geocoding_confidence"`
	AddressVerified    bool                 `json:"address_verified" db:"address_verified"`
	InputMethod        InputMethod          `json:"input_method" db:"input_method"`
	LocationPrivacy    LocationPrivacyLevel `json:"location_privacy" db:"location_privacy"`
	CreatedAt          time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at" db:"updated_at"`
}

// GeocodeValidateRequest запрос валидации геокодирования
type GeocodeValidateRequest struct {
	Address     string           `json:"address" validate:"required,min=5"`
	Context     *GeocodeContext  `json:"context,omitempty"`
	Language    string           `json:"language,omitempty"`
	CountryCode string           `json:"country_code,omitempty"`
}

// GeocodeContext контекст для улучшения геокодирования
type GeocodeContext struct {
	Bounds      *Bounds  `json:"bounds,omitempty"`
	ProximityTo *Point   `json:"proximity_to,omitempty"`
	PlaceTypes  []string `json:"place_types,omitempty"`
}

// GeocodeValidateResponse ответ валидации геокодирования
type GeocodeValidateResponse struct {
	Success           bool                 `json:"success"`
	Location          *Point               `json:"location"`
	AddressComponents *AddressComponents   `json:"address_components"`
	FormattedAddress  string               `json:"formatted_address"`
	Confidence        float64              `json:"confidence"`
	Warnings          []string             `json:"warnings,omitempty"`
	Suggestions       []AddressSuggestion  `json:"suggestions,omitempty"`
}

// AddressSuggestion предложение адреса
type AddressSuggestion struct {
	ID                string             `json:"id"`
	Text              string             `json:"text"`
	PlaceName         string             `json:"place_name"`
	Location          Point              `json:"location"`
	AddressComponents AddressComponents  `json:"address_components"`
	Confidence        float64            `json:"confidence"`
	PlaceTypes        []string           `json:"place_types"`
}

// UpdateAddressRequest запрос обновления адреса
type UpdateAddressRequest struct {
	Address           string               `json:"address" validate:"required,min=5"`
	Location          Point                `json:"location" validate:"required"`
	AddressComponents *AddressComponents   `json:"address_components,omitempty"`
	LocationPrivacy   LocationPrivacyLevel `json:"location_privacy" validate:"required"`
	InputMethod       InputMethod          `json:"input_method" validate:"required"`
	Verified          bool                 `json:"verified"`
}

// AddressChangeLog лог изменений адреса
type AddressChangeLog struct {
	ID               int64     `json:"id" db:"id"`
	ListingID        int64     `json:"listing_id" db:"listing_id"`
	UserID           int64     `json:"user_id" db:"user_id"`
	OldAddress       string    `json:"old_address" db:"old_address"`
	NewAddress       string    `json:"new_address" db:"new_address"`
	OldLocation      *Point    `json:"old_location,omitempty"`
	NewLocation      *Point    `json:"new_location,omitempty"`
	ChangeReason     string    `json:"change_reason" db:"change_reason"`
	ConfidenceBefore float64   `json:"confidence_before" db:"confidence_before"`
	ConfidenceAfter  float64   `json:"confidence_after" db:"confidence_after"`
	IPAddress        string    `json:"ip_address" db:"ip_address"`
	UserAgent        string    `json:"user_agent" db:"user_agent"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// GeocodingCacheEntry кэш геокодирования
type GeocodingCacheEntry struct {
	ID                int64              `json:"id" db:"id"`
	InputAddress      string             `json:"input_address" db:"input_address"`
	NormalizedAddress string             `json:"normalized_address" db:"normalized_address"`
	Location          Point              `json:"location"`
	AddressComponents AddressComponents  `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address" db:"formatted_address"`
	Confidence        float64            `json:"confidence" db:"confidence"`
	Provider          string             `json:"provider" db:"provider"`
	Language          string             `json:"language" db:"language"`
	CountryCode       string             `json:"country_code" db:"country_code"`
	CacheHits         int64              `json:"cache_hits" db:"cache_hits"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" db:"updated_at"`
	ExpiresAt         time.Time          `json:"expires_at" db:"expires_at"`
}

// CalculateBlurRadius вычисляет радиус размытия по уровню приватности
func (p LocationPrivacyLevel) CalculateBlurRadius() float64 {
	switch p {
	case PrivacyExact:
		return 0
	case PrivacyStreet:
		return 150 // ±150м
	case PrivacyDistrict:
		return 750 // ±750м
	case PrivacyCity:
		return 5000 // ±5км
	default:
		return 0
	}
}

// GetZoomLevel возвращает подходящий уровень зума для карты
func (p LocationPrivacyLevel) GetZoomLevel() float64 {
	switch p {
	case PrivacyExact:
		return 16
	case PrivacyStreet:
		return 15
	case PrivacyDistrict:
		return 14
	case PrivacyCity:
		return 12
	default:
		return 14
	}
}

// IsValid проверяет корректность уровня приватности
func (p LocationPrivacyLevel) IsValid() bool {
	switch p {
	case PrivacyExact, PrivacyStreet, PrivacyDistrict, PrivacyCity:
		return true
	default:
		return false
	}
}

// IsValid проверяет корректность метода ввода
func (i InputMethod) IsValid() bool {
	switch i {
	case InputManual, InputGeocoded, InputMapClick, InputCurrentLocation:
		return true
	default:
		return false
	}
}

func cosine(degrees float64) float64 {
	radians := degrees * math.Pi / 180
	return math.Cos(radians)
}
