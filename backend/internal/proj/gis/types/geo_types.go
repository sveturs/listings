package types

import (
	"encoding/json"
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
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	Price           float64   `json:"price"`
	Category        string    `json:"category"`
	Location        Point     `json:"location"`
	Address         string    `json:"address,omitempty"`
	Images          []string  `json:"images,omitempty"`
	UserID          int       `json:"user_id"`
	StorefrontID    *int      `json:"storefront_id,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ViewsCount      int       `json:"views_count"`
	Rating          float64   `json:"rating,omitempty"`
	Distance        *float64  `json:"distance,omitempty"` // Расстояние от точки поиска (в метрах)
	ItemType        string    `json:"item_type"`          // marketplace_listing, storefront_product, storefront
	DisplayStrategy string    `json:"display_strategy"`   // individual, storefront_grouped
	PrivacyLevel    string    `json:"privacy_level,omitempty"`
	BlurRadius      int       `json:"blur_radius_meters,omitempty"`
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
	ListingID           int64                `json:"listing_id" db:"listing_id"`
	Location            Point                `json:"location"`
	BlurredLocation     *Point               `json:"blurred_location,omitempty"`
	Geohash             string               `json:"geohash" db:"geohash"`
	IsPrecise           bool                 `json:"is_precise" db:"is_precise"`
	BlurRadius          float64              `json:"blur_radius" db:"blur_radius"`
	AddressComponents   *AddressComponents   `json:"address_components,omitempty"`
	FormattedAddress    string               `json:"formatted_address" db:"formatted_address"`
	GeocodingConfidence float64              `json:"geocoding_confidence" db:"geocoding_confidence"`
	AddressVerified     bool                 `json:"address_verified" db:"address_verified"`
	InputMethod         InputMethod          `json:"input_method" db:"input_method"`
	LocationPrivacy     LocationPrivacyLevel `json:"location_privacy" db:"location_privacy"`
	CreatedAt           time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at" db:"updated_at"`
}

// GeocodeValidateRequest запрос валидации геокодирования
type GeocodeValidateRequest struct {
	Address     string          `json:"address" validate:"required,min=5"`
	Context     *GeocodeContext `json:"context,omitempty"`
	Language    string          `json:"language,omitempty"`
	CountryCode string          `json:"country_code,omitempty"`
}

// GeocodeContext контекст для улучшения геокодирования
type GeocodeContext struct {
	Bounds      *Bounds  `json:"bounds,omitempty"`
	ProximityTo *Point   `json:"proximity_to,omitempty"`
	PlaceTypes  []string `json:"place_types,omitempty"`
}

// GeocodeValidateResponse ответ валидации геокодирования
type GeocodeValidateResponse struct {
	Success           bool                `json:"success"`
	Location          *Point              `json:"location"`
	AddressComponents *AddressComponents  `json:"address_components"`
	FormattedAddress  string              `json:"formatted_address"`
	Confidence        float64             `json:"confidence"`
	Warnings          []string            `json:"warnings,omitempty"`
	Suggestions       []AddressSuggestion `json:"suggestions,omitempty"`
}

// AddressSuggestion предложение адреса
type AddressSuggestion struct {
	ID                string            `json:"id"`
	Text              string            `json:"text"`
	PlaceName         string            `json:"place_name"`
	Location          Point             `json:"location"`
	AddressComponents AddressComponents `json:"address_components"`
	Confidence        float64           `json:"confidence"`
	PlaceTypes        []string          `json:"place_types"`
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
	ID                int64             `json:"id" db:"id"`
	InputAddress      string            `json:"input_address" db:"input_address"`
	NormalizedAddress string            `json:"normalized_address" db:"normalized_address"`
	Location          Point             `json:"location"`
	AddressComponents AddressComponents `json:"address_components"`
	FormattedAddress  string            `json:"formatted_address" db:"formatted_address"`
	Confidence        float64           `json:"confidence" db:"confidence"`
	Provider          string            `json:"provider" db:"provider"`
	Language          string            `json:"language" db:"language"`
	CountryCode       string            `json:"country_code" db:"country_code"`
	CacheHits         int64             `json:"cache_hits" db:"cache_hits"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
	ExpiresAt         time.Time         `json:"expires_at" db:"expires_at"`
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

// RadiusSearchRequest запрос радиусного поиска
type RadiusSearchRequest struct {
	Latitude  float64        `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64        `json:"longitude" validate:"required,min=-180,max=180"`
	Radius    float64        `json:"radius" validate:"required,min=0.1,max=100000"` // в метрах
	Filters   *RadiusFilters `json:"filters,omitempty"`
	Limit     int            `json:"limit,omitempty" validate:"min=1,max=1000"`
	Offset    int            `json:"offset,omitempty" validate:"min=0"`
}

// RadiusFilters фильтры для радиусного поиска
type RadiusFilters struct {
	Categories  []string `json:"categories,omitempty"`
	CategoryIDs []int    `json:"category_ids,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	SearchQuery string   `json:"q,omitempty"`          // Текстовый поиск
	UserID      *int     `json:"user_id,omitempty"`    // Фильтр по пользователю
	Status      string   `json:"status,omitempty"`     // Фильтр по статусу
	SortBy      string   `json:"sort_by,omitempty"`    // distance, price, created_at
	SortOrder   string   `json:"sort_order,omitempty"` // asc, desc
}

// RadiusSearchResponse ответ на радиусный поиск
type RadiusSearchResponse struct {
	Listings     []GeoListing `json:"listings"`
	TotalCount   int64        `json:"total_count"`
	HasMore      bool         `json:"has_more"`
	SearchRadius float64      `json:"search_radius"` // радиус в метрах
	SearchCenter Point        `json:"search_center"` // центр поиска
}

// Validate валидация запроса радиусного поиска
func (r *RadiusSearchRequest) Validate() error {
	if r.Latitude < -90 || r.Latitude > 90 {
		return ErrInvalidLatitude
	}
	if r.Longitude < -180 || r.Longitude > 180 {
		return ErrInvalidLongitude
	}
	if r.Radius <= 0 || r.Radius > 100000 {
		return ErrInvalidRadius
	}
	if r.Limit < 0 {
		r.Limit = 50 // дефолтное значение
	}
	if r.Limit > 1000 {
		r.Limit = 1000 // максимальное значение
	}
	if r.Offset < 0 {
		r.Offset = 0
	}
	return nil
}

// ToSearchParams преобразует в стандартные параметры поиска
func (r *RadiusSearchRequest) ToSearchParams() SearchParams {
	params := SearchParams{
		Center:    &Point{Lat: r.Latitude, Lng: r.Longitude},
		RadiusKm:  r.Radius / 1000.0, // конвертируем метры в километры
		Limit:     r.Limit,
		Offset:    r.Offset,
		SortBy:    "distance", // по умолчанию сортируем по расстоянию
		SortOrder: "asc",
	}

	if r.Filters != nil {
		params.Categories = r.Filters.Categories
		params.CategoryIDs = r.Filters.CategoryIDs
		params.MinPrice = r.Filters.MinPrice
		params.MaxPrice = r.Filters.MaxPrice
		params.SearchQuery = r.Filters.SearchQuery
		params.UserID = r.Filters.UserID
		params.Status = r.Filters.Status

		if r.Filters.SortBy != "" {
			params.SortBy = r.Filters.SortBy
		}
		if r.Filters.SortOrder != "" {
			params.SortOrder = r.Filters.SortOrder
		}
	}

	return params
}

func cosine(degrees float64) float64 {
	radians := degrees * math.Pi / 180
	return math.Cos(radians)
}

// ========== PHASE 2.5: Типы для городов и контекстно-зависимого поиска ==========

// City представляет город с районами
type City struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	CountryCode  string    `json:"country_code" db:"country_code"`
	CenterPoint  *Point    `json:"center_point,omitempty" db:"-"`
	Boundary     *Polygon  `json:"boundary,omitempty" db:"-"`
	Population   *int      `json:"population,omitempty" db:"population"`
	AreaKm2      *float64  `json:"area_km2,omitempty" db:"area_km2"`
	PostalCodes  []string  `json:"postal_codes,omitempty" db:"postal_codes"`
	HasDistricts bool      `json:"has_districts" db:"has_districts"`
	Priority     int       `json:"priority" db:"priority"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CitySearchParams параметры поиска городов
type CitySearchParams struct {
	Bounds       *Bounds `json:"bounds,omitempty"`
	CountryCode  string  `json:"country_code,omitempty"`
	HasDistricts *bool   `json:"has_districts,omitempty"`
	SearchQuery  string  `json:"q,omitempty"`
	Limit        int     `json:"limit,omitempty"`
	Offset       int     `json:"offset,omitempty"`
}

// VisibleCitiesRequest запрос для определения видимых городов
type VisibleCitiesRequest struct {
	Bounds *Bounds `json:"bounds" validate:"required"`
	Center *Point  `json:"center" validate:"required"`
}

// VisibleCitiesResponse ответ с видимыми городами
type VisibleCitiesResponse struct {
	VisibleCities []CityWithDistance `json:"visible_cities"`
	ClosestCity   *CityWithDistance  `json:"closest_city,omitempty"`
}

// CityWithDistance город с расстоянием до центра карты
type CityWithDistance struct {
	City     City    `json:"city"`
	Distance float64 `json:"distance"` // расстояние в метрах
}

// DistrictBoundaryResponse ответ с границами района
type DistrictBoundaryResponse struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	CityID   *string         `json:"city_id,omitempty"`
	Boundary json.RawMessage `json:"boundary"` // GeoJSON polygon as JSON object
}

// Distance вычисляет расстояние между двумя точками в метрах
func (p Point) Distance(other Point) float64 {
	const R = 6371000 // радиус Земли в метрах

	lat1Rad := p.Lat * math.Pi / 180
	lat2Rad := other.Lat * math.Pi / 180
	deltaLatRad := (other.Lat - p.Lat) * math.Pi / 180
	deltaLngRad := (other.Lng - p.Lng) * math.Pi / 180

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLngRad/2)*math.Sin(deltaLngRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// CalculateDistance вычисляет расстояние до города
func (c City) CalculateDistance(center Point) float64 {
	return c.CenterPoint.Distance(center)
}

// Validate валидация параметров поиска видимых городов
func (v *VisibleCitiesRequest) Validate() error {
	if v.Bounds == nil {
		return ErrMissingBounds
	}
	if v.Center == nil {
		return ErrMissingCenter
	}
	return v.Bounds.Validate()
}
