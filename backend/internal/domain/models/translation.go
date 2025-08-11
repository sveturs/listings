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
}
