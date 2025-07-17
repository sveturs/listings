// backend/internal/domain/models/map.go
package models

// MapMarker представляет маркер на карте
type MapMarker struct {
	ID         int     `json:"id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	Condition  string  `json:"condition"`
	CategoryID int     `json:"category_id"`
	MainImage  string  `json:"main_image,omitempty"`
	UserID     int     `json:"user_id"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	CreatedAt  string  `json:"created_at"`
	ViewsCount int     `json:"views_count"`
	Rating     float64 `json:"rating,omitempty"`
	IsFavorite bool    `json:"is_favorite,omitempty"`
}

// MapCluster представляет кластер маркеров на карте
type MapCluster struct {
	ID         string    `json:"id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Count      int       `json:"count"`
	Bounds     MapBounds `json:"bounds"`
	AvgPrice   float64   `json:"avg_price,omitempty"`
	Categories []int     `json:"categories,omitempty"`
}

// MapBounds представляет границы области на карте
type MapBounds struct {
	NorthEast MapPoint `json:"northeast"`
	SouthWest MapPoint `json:"southwest"`
}

// MapPoint представляет точку на карте
type MapPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// MapFilters представляет фильтры для поиска на карте
type MapFilters struct {
	CategoryIDs []int    `json:"category_ids,omitempty"`
	Condition   string   `json:"condition,omitempty"`
	MinPrice    *float64 `json:"min_price,omitempty"`
	MaxPrice    *float64 `json:"max_price,omitempty"`
	Query       string   `json:"query,omitempty"`
	Radius      *float64 `json:"radius,omitempty"` // В километрах
	CenterLat   *float64 `json:"center_lat,omitempty"`
	CenterLng   *float64 `json:"center_lng,omitempty"`
}

// MapSearchRequest представляет запрос поиска на карте
type MapSearchRequest struct {
	Bounds  MapBounds  `json:"bounds"`
	Zoom    int        `json:"zoom"`
	Filters MapFilters `json:"filters"`
	Limit   int        `json:"limit,omitempty"`
}

// MapSearchResponse представляет ответ поиска на карте
type MapSearchResponse struct {
	Type     string       `json:"type"` // "markers" или "clusters"
	Markers  []MapMarker  `json:"markers,omitempty"`
	Clusters []MapCluster `json:"clusters,omitempty"`
	Total    int          `json:"total"`
	Bounds   MapBounds    `json:"bounds"`
	Zoom     int          `json:"zoom"`
}
