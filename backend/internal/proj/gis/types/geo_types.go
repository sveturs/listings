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

func cosine(degrees float64) float64 {
	radians := degrees * math.Pi / 180
	return math.Cos(radians)
}
