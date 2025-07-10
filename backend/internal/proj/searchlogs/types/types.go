package types

import "time"

// SearchLogEntry представляет запись логов поиска
type SearchLogEntry struct {
	Query           string                 `json:"query"`
	UserID          *int                   `json:"user_id"`
	SessionID       string                 `json:"session_id"`
	ResultCount     int                    `json:"result_count"`
	ResponseTimeMS  int64                  `json:"response_time_ms"`
	Filters         map[string]interface{} `json:"filters"`
	CategoryID      *int                   `json:"category_id"`
	PriceMin        *float64               `json:"price_min"`
	PriceMax        *float64               `json:"price_max"`
	Location        *string                `json:"location"`
	Language        string                 `json:"language"`
	DeviceType      string                 `json:"device_type"`
	UserAgent       string                 `json:"user_agent"`
	IP              string                 `json:"ip"`
	SearchType      string                 `json:"search_type"`
	HasSpellCorrect bool                   `json:"has_spell_correct"`
	ClickedItems    []int                  `json:"clicked_items"`
	Timestamp       time.Time              `json:"timestamp"`
}