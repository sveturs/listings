package models

import (
	"encoding/json"
	"time"
)

// Storefront представляет витрину магазина
type Storefront struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	LogoPath            string    `json:"logo_path"`
	Slug                string    `json:"slug"`
	Status              string    `json:"status"` // active, inactive, suspended
	CreationTransactionID *int     `json:"creation_transaction_id,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ImportSource представляет источник данных для импорта
type ImportSource struct {
    ID               int             `json:"id"`
    StorefrontID     int             `json:"storefront_id"`
    Type             string          `json:"type"` // csv, xml, json
    URL              string          `json:"url"`
    AuthData         json.RawMessage `json:"auth_data,omitempty"`
    Schedule         string          `json:"schedule,omitempty"`
    Mapping          json.RawMessage `json:"mapping,omitempty"`
    LastImportAt     *time.Time      `json:"last_import_at,omitempty"`
    LastImportStatus string          `json:"last_import_status,omitempty"` // Изменяем на обычную строку, но обрабатываем NULL при сканировании
    LastImportLog    string          `json:"last_import_log,omitempty"`    // Изменяем на обычную строку, но обрабатываем NULL при сканировании
    CreatedAt        time.Time       `json:"created_at"`
    UpdatedAt        time.Time       `json:"updated_at"`
}

// ImportHistory представляет запись об импорте
type ImportHistory struct {
	ID            int        `json:"id"`
	SourceID      int        `json:"source_id"`
	Status        string     `json:"status"` // success, failed, partial
	ItemsTotal    int        `json:"items_total"`
	ItemsImported int        `json:"items_imported"`
	ItemsFailed   int        `json:"items_failed"`
	Log           string     `json:"log,omitempty"`
	StartedAt     time.Time  `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at,omitempty"`
}

// StorefrontCreate содержит данные для создания витрины
type StorefrontCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Logo        []byte `json:"-"`
}

// ImportSourceCreate содержит данные для создания источника импорта
type ImportSourceCreate struct {
	StorefrontID int             `json:"storefront_id" validate:"required"`
	Type         string          `json:"type" validate:"required,oneof=csv xml json"`
	URL          string          `json:"url"`
	AuthData     json.RawMessage `json:"auth_data,omitempty"`
	Schedule     string          `json:"schedule,omitempty"`
	Mapping      json.RawMessage `json:"mapping,omitempty"`
}