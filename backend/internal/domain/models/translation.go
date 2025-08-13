package models

import (
	"time"
)

type Translation struct {
	ID                  int                    `json:"id"`
	EntityType          string                 `json:"entity_type"`
	EntityID            int                    `json:"entity_id"`
	Language            string                 `json:"language"`
	FieldName           string                 `json:"field_name"`
	TranslatedText      string                 `json:"translated_text"`
	IsMachineTranslated bool                   `json:"is_machine_translated"`
	IsVerified          bool                   `json:"is_verified"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	UpdatedBy           int                    `json:"updated_by,omitempty"`
	Version             int                    `json:"version,omitempty"`
	CurrentVersion      int                    `json:"current_version,omitempty"`
	LastModifiedBy      *int64                 `json:"last_modified_by,omitempty"`
	LastModifiedAt      *time.Time             `json:"last_modified_at,omitempty"`
}

// TranslationVersion представляет версию перевода
type TranslationVersion struct {
	ID             int64                  `json:"id"`
	TranslationID  int                    `json:"translation_id"`
	EntityType     string                 `json:"entity_type"`
	EntityID       int                    `json:"entity_id"`
	FieldName      string                 `json:"field_name"`
	Language       string                 `json:"language"`
	TranslatedText string                 `json:"translated_text"`
	PreviousText   *string                `json:"previous_text,omitempty"`
	Version        int                    `json:"version"`
	ChangeType     string                 `json:"change_type"`
	ChangedBy      *int64                 `json:"changed_by,omitempty"`
	ChangedAt      time.Time              `json:"changed_at"`
	ChangeReason   *string                `json:"change_reason,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// TranslationSnapshot представляет снимок всех переводов
type TranslationSnapshot struct {
	ID                int64                  `json:"id"`
	Name              string                 `json:"name"`
	Description       *string                `json:"description,omitempty"`
	CreatedBy         *int64                 `json:"created_by,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	TranslationsCount int                    `json:"translations_count"`
	SnapshotData      map[string]interface{} `json:"snapshot_data"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// TranslationVersionFilter фильтры для поиска версий
type TranslationVersionFilter struct {
	TranslationID *int       `json:"translation_id,omitempty"`
	EntityType    *string    `json:"entity_type,omitempty"`
	EntityID      *int       `json:"entity_id,omitempty"`
	Language      *string    `json:"language,omitempty"`
	FieldName     *string    `json:"field_name,omitempty"`
	ChangedBy     *int64     `json:"changed_by,omitempty"`
	ChangeType    *string    `json:"change_type,omitempty"`
	DateFrom      *time.Time `json:"date_from,omitempty"`
	DateTo        *time.Time `json:"date_to,omitempty"`
	Limit         int        `json:"limit,omitempty"`
	Offset        int        `json:"offset,omitempty"`
}
