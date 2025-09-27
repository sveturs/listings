package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// SavedSearch представляет сохраненный поиск пользователя
type SavedSearch struct {
	ID              int         `json:"id" db:"id"`
	UserID          int         `json:"user_id" db:"user_id"`
	Name            string      `json:"name" db:"name"`
	Filters         FiltersJSON `json:"filters" db:"filters"`
	SearchType      string      `json:"search_type" db:"search_type"`
	NotifyEnabled   bool        `json:"notify_enabled" db:"notify_enabled"`
	NotifyFrequency string      `json:"notify_frequency" db:"notify_frequency"`
	LastNotifiedAt  *time.Time  `json:"last_notified_at,omitempty" db:"last_notified_at"`
	ResultsCount    int         `json:"results_count" db:"results_count"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}

// FiltersJSON для хранения JSON фильтров в базе данных
type FiltersJSON map[string]interface{}

// Scan реализует интерфейс sql.Scanner для чтения из базы данных
func (f *FiltersJSON) Scan(value interface{}) error {
	if value == nil {
		*f = make(map[string]interface{})
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, f)
	case string:
		return json.Unmarshal([]byte(v), f)
	default:
		*f = make(map[string]interface{})
		return nil
	}
}

// Value реализует интерфейс driver.Valuer для записи в базу данных
func (f FiltersJSON) Value() (driver.Value, error) {
	if f == nil {
		return "{}", nil
	}
	return json.Marshal(f)
}

// SavedSearchNotification представляет уведомление о новых результатах поиска
type SavedSearchNotification struct {
	ID               int        `json:"id" db:"id"`
	SavedSearchID    int        `json:"saved_search_id" db:"saved_search_id"`
	NewListingsCount int        `json:"new_listings_count" db:"new_listings_count"`
	NotificationSent bool       `json:"notification_sent" db:"notification_sent"`
	SentAt           *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	ErrorMessage     *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// UserCarViewHistory представляет историю просмотров автомобилей
type UserCarViewHistory struct {
	ID                  int       `json:"id" db:"id"`
	UserID              *int      `json:"user_id,omitempty" db:"user_id"`
	ListingID           int       `json:"listing_id" db:"listing_id"`
	SessionID           *string   `json:"session_id,omitempty" db:"session_id"`
	ViewedAt            time.Time `json:"viewed_at" db:"viewed_at"`
	ViewDurationSeconds *int      `json:"view_duration_seconds,omitempty" db:"view_duration_seconds"`
	Referrer            *string   `json:"referrer,omitempty" db:"referrer"`
	DeviceType          *string   `json:"device_type,omitempty" db:"device_type"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}
