package domain

import (
	"encoding/json"
	"time"
)

// SearchLog представляет запись лога поискового запроса
type SearchLog struct {
	ID             int64           `db:"id" json:"id"`
	UserID         *int            `db:"user_id" json:"user_id,omitempty"`
	SessionID      string          `db:"session_id" json:"session_id"`
	QueryText      string          `db:"query_text" json:"query_text"`
	Filters        json.RawMessage `db:"filters" json:"filters,omitempty"`
	CategoryID     *int            `db:"category_id" json:"category_id,omitempty"`
	Location       json.RawMessage `db:"location" json:"location,omitempty"`
	ResultsCount   int             `db:"results_count" json:"results_count"`
	ResponseTimeMs int             `db:"response_time_ms" json:"response_time_ms"`
	Page           int             `db:"page" json:"page"`
	PerPage        int             `db:"per_page" json:"per_page"`
	SortBy         *string         `db:"sort_by" json:"sort_by,omitempty"`
	UserAgent      *string         `db:"user_agent" json:"user_agent,omitempty"`
	IPAddress      *string         `db:"ip_address" json:"ip_address,omitempty"`
	Referer        *string         `db:"referer" json:"referer,omitempty"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
}

// SearchAnalytics представляет агрегированную аналитику поисковых запросов
type SearchAnalytics struct {
	ID                  int64     `db:"id" json:"id"`
	Date                time.Time `db:"date" json:"date"`
	Hour                int       `db:"hour" json:"hour"`
	QueryText           string    `db:"query_text" json:"query_text"`
	CategoryID          *int      `db:"category_id" json:"category_id,omitempty"`
	LocationCountry     *string   `db:"location_country" json:"location_country,omitempty"`
	LocationRegion      *string   `db:"location_region" json:"location_region,omitempty"`
	LocationCity        *string   `db:"location_city" json:"location_city,omitempty"`
	SearchCount         int       `db:"search_count" json:"search_count"`
	UniqueUsersCount    int       `db:"unique_users_count" json:"unique_users_count"`
	UniqueSessionsCount int       `db:"unique_sessions_count" json:"unique_sessions_count"`
	AvgResultsCount     *float64  `db:"avg_results_count" json:"avg_results_count,omitempty"`
	AvgResponseTimeMs   *float64  `db:"avg_response_time_ms" json:"avg_response_time_ms,omitempty"`
	ZeroResultsCount    int       `db:"zero_results_count" json:"zero_results_count"`
	ClickThroughRate    *float64  `db:"click_through_rate" json:"click_through_rate,omitempty"`
	ConversionRate      *float64  `db:"conversion_rate" json:"conversion_rate,omitempty"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// SearchResultClick представляет клик по результату поиска
type SearchResultClick struct {
	ID          int64     `db:"id" json:"id"`
	SearchLogID int64     `db:"search_log_id" json:"search_log_id"`
	ListingID   int       `db:"listing_id" json:"listing_id"`
	Position    int       `db:"position" json:"position"`
	ClickedAt   time.Time `db:"clicked_at" json:"clicked_at"`
}

// SearchTrendingQuery представляет популярный поисковый запрос
type SearchTrendingQuery struct {
	ID              int64     `db:"id" json:"id"`
	QueryText       string    `db:"query_text" json:"query_text"`
	CategoryID      *int      `db:"category_id" json:"category_id,omitempty"`
	LocationCountry *string   `db:"location_country" json:"location_country,omitempty"`
	TrendScore      float64   `db:"trend_score" json:"trend_score"`
	SearchCount24h  int       `db:"search_count_24h" json:"search_count_24h"`
	SearchCount7d   int       `db:"search_count_7d" json:"search_count_7d"`
	SearchCount30d  int       `db:"search_count_30d" json:"search_count_30d"`
	FirstSeenAt     time.Time `db:"first_seen_at" json:"first_seen_at"`
	LastSeenAt      time.Time `db:"last_seen_at" json:"last_seen_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// SearchLocation представляет местоположение для поискового запроса
type SearchLocation struct {
	Country string  `json:"country,omitempty"`
	Region  string  `json:"region,omitempty"`
	City    string  `json:"city,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
}

// SearchFilters представляет фильтры поискового запроса
type SearchFilters struct {
	PriceMin   *float64               `json:"price_min,omitempty"`
	PriceMax   *float64               `json:"price_max,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// SearchLogInput представляет входные данные для создания лога поиска
type SearchLogInput struct {
	UserID         *int            `json:"user_id,omitempty"`
	SessionID      string          `json:"session_id"`
	QueryText      string          `json:"query_text"`
	Filters        *SearchFilters  `json:"filters,omitempty"`
	CategoryID     *int            `json:"category_id,omitempty"`
	Location       *SearchLocation `json:"location,omitempty"`
	ResultsCount   int             `json:"results_count"`
	ResponseTimeMs int             `json:"response_time_ms"`
	Page           int             `json:"page"`
	PerPage        int             `json:"per_page"`
	SortBy         *string         `json:"sort_by,omitempty"`
	UserAgent      string          `json:"user_agent,omitempty"`
	IPAddress      string          `json:"ip_address,omitempty"`
	Referer        string          `json:"referer,omitempty"`
}
