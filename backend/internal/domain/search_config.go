package domain

import (
	"time"
)

// SearchWeight представляет веса для полей поиска
type SearchWeight struct {
	ID          int64     `json:"id" db:"id"`
	FieldName   string    `json:"field_name" db:"field_name"`
	Weight      float64   `json:"weight" db:"weight"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SearchSynonym представляет синонимы для поиска
type SearchSynonym struct {
	ID        int64     `json:"id" db:"id"`
	Term      string    `json:"term" db:"term"`
	Synonyms  []string  `json:"synonyms" db:"synonyms"`
	Language  string    `json:"language" db:"language"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TransliterationRule представляет правило транслитерации
type TransliterationRule struct {
	ID          int64     `json:"id" db:"id"`
	FromScript  string    `json:"from_script" db:"from_script"`
	ToScript    string    `json:"to_script" db:"to_script"`
	FromPattern string    `json:"from_pattern" db:"from_pattern"`
	ToPattern   string    `json:"to_pattern" db:"to_pattern"`
	Priority    int       `json:"priority" db:"priority"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SearchStatistics представляет статистику поиска
type SearchStatistics struct {
	ID             int64     `json:"id" db:"id"`
	Query          string    `json:"query" db:"query"`
	ResultsCount   int       `json:"results_count" db:"results_count"`
	SearchDuration int64     `json:"search_duration_ms" db:"search_duration_ms"`
	UserID         *int64    `json:"user_id,omitempty" db:"user_id"`
	SearchFilters  string    `json:"search_filters,omitempty" db:"search_filters"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// SearchConfig представляет общую конфигурацию поиска
type SearchConfig struct {
	ID                     int64     `json:"id" db:"id"`
	MinSearchLength        int       `json:"min_search_length" db:"min_search_length"`
	MaxSuggestions         int       `json:"max_suggestions" db:"max_suggestions"`
	FuzzyEnabled           bool      `json:"fuzzy_enabled" db:"fuzzy_enabled"`
	FuzzyMaxEdits          int       `json:"fuzzy_max_edits" db:"fuzzy_max_edits"`
	SynonymsEnabled        bool      `json:"synonyms_enabled" db:"synonyms_enabled"`
	TransliterationEnabled bool      `json:"transliteration_enabled" db:"transliteration_enabled"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}
