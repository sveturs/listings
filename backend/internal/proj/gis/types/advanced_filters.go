package types

import (
	"time"
)

// TransportMode тип транспорта для расчета изохрон
type TransportMode string

const (
	TransportWalking TransportMode = "walking"
	TransportDriving TransportMode = "driving"
	TransportCycling TransportMode = "cycling"
	TransportTransit TransportMode = "transit"
)

// POIType тип точки интереса
type POIType string

const (
	POISchool      POIType = "school"
	POIHospital    POIType = "hospital"
	POIMetro       POIType = "metro"
	POISupermarket POIType = "supermarket"
	POIPark        POIType = "park"
	POIBank        POIType = "bank"
	POIPharmacy    POIType = "pharmacy"
	POIBusStop     POIType = "bus_stop"
)

// AdvancedGeoFilters расширенные геофильтры для поиска
type AdvancedGeoFilters struct {
	// Фильтр по времени пути
	TravelTime *TravelTimeFilter `json:"travel_time,omitempty"`

	// Фильтр по точкам интереса
	POIFilter *POIFilter `json:"poi_filter,omitempty"`

	// Фильтр по плотности объявлений
	DensityFilter *DensityFilter `json:"density_filter,omitempty"`

	// Базовые геофильтры
	BoundingBox *BoundingBox  `json:"bounding_box,omitempty"`
	Radius      *RadiusFilter `json:"radius,omitempty"`
}

// TravelTimeFilter фильтр по времени пути через изохроны
type TravelTimeFilter struct {
	// Центральная точка для расчета изохрон
	CenterLat float64 `json:"center_lat" validate:"required,latitude"`
	CenterLng float64 `json:"center_lng" validate:"required,longitude"`

	// Максимальное время пути в минутах (5-60)
	MaxMinutes int `json:"max_minutes" validate:"required,min=5,max=60"`

	// Тип транспорта
	TransportMode TransportMode `json:"transport_mode" validate:"required,oneof=walking driving cycling transit"`
}

// POIFilter фильтр по близости к точкам интереса
type POIFilter struct {
	// Тип точки интереса
	POIType POIType `json:"poi_type" validate:"required"`

	// Максимальное расстояние в метрах
	MaxDistance int `json:"max_distance" validate:"required,min=100,max=5000"`

	// Минимальное количество POI в радиусе
	MinCount int `json:"min_count,omitempty" validate:"min=0,max=10"`
}

// DensityFilter фильтр по плотности объявлений
type DensityFilter struct {
	// Избегать переполненных районов
	AvoidCrowded bool `json:"avoid_crowded"`

	// Максимальное количество объявлений на км²
	MaxDensity float64 `json:"max_density,omitempty" validate:"min=0"`

	// Минимальное количество объявлений на км²
	MinDensity float64 `json:"min_density,omitempty" validate:"min=0"`
}

// BoundingBox ограничивающий прямоугольник для поиска
type BoundingBox struct {
	MinLat float64 `json:"min_lat" validate:"required,latitude"`
	MinLng float64 `json:"min_lng" validate:"required,longitude"`
	MaxLat float64 `json:"max_lat" validate:"required,latitude"`
	MaxLng float64 `json:"max_lng" validate:"required,longitude"`
}

// RadiusFilter фильтр по радиусу
type RadiusFilter struct {
	CenterLat float64 `json:"center_lat" validate:"required,latitude"`
	CenterLng float64 `json:"center_lng" validate:"required,longitude"`
	RadiusKm  float64 `json:"radius_km" validate:"required,min=0.1,max=50"`
}

// IsohronResponse ответ от API MapBox для изохрон
type IsohronResponse struct {
	Features []struct {
		Geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Contour int `json:"contour"`
		} `json:"properties"`
	} `json:"features"`
}

// POISearchResult результат поиска точек интереса
type POISearchResult struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Type     POIType `json:"type"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Distance float64 `json:"distance"`
}

// DensityAnalysisResult результат анализа плотности
type DensityAnalysisResult struct {
	GridCellID   string  `json:"grid_cell_id"`
	CenterLat    float64 `json:"center_lat"`
	CenterLng    float64 `json:"center_lng"`
	ListingCount int     `json:"listing_count"`
	AreaKm2      float64 `json:"area_km2"`
	Density      float64 `json:"density"`
}

// FilterAnalytics аналитика использования фильтров
type FilterAnalytics struct {
	UserID         string      `json:"user_id,omitempty"`
	SessionID      string      `json:"session_id"`
	FilterType     string      `json:"filter_type"`
	FilterParams   interface{} `json:"filter_params"`
	ResultCount    int         `json:"result_count"`
	ResponseTimeMs int64       `json:"response_time_ms"`
	CreatedAt      time.Time   `json:"created_at"`
}

// ApplyAdvancedFiltersRequest запрос для применения расширенных фильтров
type ApplyAdvancedFiltersRequest struct {
	Filters    AdvancedGeoFilters `json:"filters"`
	ListingIDs []string           `json:"listing_ids"`
}
