package types

import (
	"time"
)

// SearchLogEntry представляет запись лога поиска
type SearchLogEntry struct {
	ID              int64                  `json:"id"`
	UserID          *int                   `json:"user_id,omitempty"`
	SessionID       string                 `json:"session_id"`
	Query           string                 `json:"query"`      // Для совместимости с postgres.go
	QueryText       string                 `json:"query_text"` // Для совместимости с существующим кодом
	Filters         map[string]interface{} `json:"filters,omitempty"`
	CategoryID      *int                   `json:"category_id,omitempty"`
	Location        *LocationInfo          `json:"location,omitempty"`
	ResultCount     int                    `json:"result_count"`  // Для совместимости с postgres.go
	ResultsCount    int                    `json:"results_count"` // Для совместимости с существующим кодом
	ResponseTime    int                    `json:"response_time_ms"`
	ResponseTimeMS  int64                  `json:"response_time_ms_int64"` // Для совместимости с postgres.go
	Page            int                    `json:"page"`
	ItemsPerPage    int                    `json:"items_per_page"`
	SortBy          string                 `json:"sort_by,omitempty"`
	SearchType      string                 `json:"search_type"`
	ClickedItems    []int                  `json:"clicked_items,omitempty"`
	PurchasedItem   *int                   `json:"purchased_item,omitempty"`
	UserAgent       string                 `json:"user_agent,omitempty"`
	ClientIP        string                 `json:"client_ip,omitempty"`
	IP              string                 `json:"ip"` // Для совместимости с postgres.go
	Referrer        string                 `json:"referrer,omitempty"`
	DeviceType      string                 `json:"device_type,omitempty"`
	Language        string                 `json:"language"`
	PriceMin        *float64               `json:"price_min,omitempty"` // Для совместимости с postgres.go
	PriceMax        *float64               `json:"price_max,omitempty"` // Для совместимости с postgres.go
	HasSpellCorrect bool                   `json:"has_spell_correct"`   // Для совместимости с postgres.go
	CreatedAt       time.Time              `json:"created_at"`
	Timestamp       time.Time              `json:"timestamp"` // Для совместимости с postgres.go
}

// LocationInfo информация о местоположении
type LocationInfo struct {
	Country string  `json:"country,omitempty"`
	Region  string  `json:"region,omitempty"`
	City    string  `json:"city,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
}

// SearchStats статистика поиска
type SearchStats struct {
	TotalSearches       int64            `json:"total_searches"`
	UniqueUsers         int64            `json:"unique_users"`
	AvgResponseTime     float64          `json:"avg_response_time_ms"`
	AvgResponseTimeMS   float64          `json:"avg_response_time_ms_alt"` // Для совместимости с postgres.go
	SearchesWithResults int64            `json:"searches_with_results"`
	ConversionRate      float64          `json:"conversion_rate"`
	TopCategories       []CategoryStats  `json:"top_categories"`
	PeakHours           []HourStats      `json:"peak_hours"`
	SearchesByHour      map[int]int64    `json:"searches_by_hour"`     // Для совместимости с postgres.go
	DeviceStats         map[string]int64 `json:"device_stats"`         // Для совместимости с postgres.go
	ZeroResultSearches  int64            `json:"zero_result_searches"` // Для совместимости с postgres.go
	TopQueries          []PopularSearch  `json:"top_queries"`          // Для совместимости с postgres.go
}

// CategoryStats статистика по категории
type CategoryStats struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	SearchCount  int64   `json:"search_count"`
	Percentage   float64 `json:"percentage"`
}

// HourStats статистика по часам
type HourStats struct {
	Hour        int   `json:"hour"`
	SearchCount int64 `json:"search_count"`
}

// PopularSearch популярный поисковый запрос
type PopularSearch struct {
	Query        string    `json:"query"`
	Count        int64     `json:"count"`
	LastSearched time.Time `json:"last_searched"`
	AvgResults   float64   `json:"avg_results"` // Для совместимости с postgres.go
	ClickRate    float64   `json:"click_rate"`  // Для совместимости с postgres.go
}
