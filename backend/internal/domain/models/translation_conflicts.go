package models

import "time"

// TranslationConflict represents a conflict between frontend and database translations
type TranslationConflict struct {
	ID                   int        `json:"id" db:"id"`
	Key                  string     `json:"key" db:"key"`
	Module               string     `json:"module" db:"module"`
	Language             string     `json:"language" db:"language"`
	FrontendValue        *string    `json:"frontend_value" db:"frontend_value"`
	DatabaseValue        *string    `json:"database_value" db:"database_value"`
	LastModifiedFrontend *time.Time `json:"last_modified_frontend" db:"last_modified_frontend"`
	LastModifiedDatabase *time.Time `json:"last_modified_database" db:"last_modified_database"`
	ConflictType         string     `json:"conflict_type" db:"conflict_type"` // value_mismatch, missing_in_frontend, missing_in_database
	Resolved             bool       `json:"resolved" db:"resolved"`
	Resolution           *string    `json:"resolution" db:"resolution"` // use_frontend, use_database, use_custom
	CustomValue          *string    `json:"custom_value" db:"custom_value"`
	ResolvedAt           *time.Time `json:"resolved_at" db:"resolved_at"`
	ResolvedBy           *int       `json:"resolved_by" db:"resolved_by"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// ConflictResolution represents a single conflict resolution request
type ConflictResolution struct {
	ConflictID  int     `json:"conflict_id" validate:"required"`
	Resolution  string  `json:"resolution" validate:"required,oneof=use_frontend use_database use_custom"`
	CustomValue *string `json:"custom_value,omitempty"`
}

// ConflictResolutionBatch represents a batch of conflict resolutions
type ConflictResolutionBatch struct {
	Resolutions []ConflictResolution `json:"resolutions" validate:"required,dive"`
}

// ConflictResolutionResult represents the result of conflict resolution
type ConflictResolutionResult struct {
	TotalProcessed int      `json:"total_processed"`
	Resolved       int      `json:"resolved"`
	Failed         int      `json:"failed"`
	Errors         []string `json:"errors,omitempty"`
}

// ConflictsFilter represents filters for getting conflicts
type ConflictsFilter struct {
	Module   *string `json:"module,omitempty"`
	Language *string `json:"language,omitempty"`
	Type     *string `json:"type,omitempty"`
	Resolved *bool   `json:"resolved,omitempty"`
	Limit    int     `json:"limit,omitempty"`
	Offset   int     `json:"offset,omitempty"`
}

// ConflictsList represents a list of conflicts with metadata
type ConflictsList struct {
	Conflicts     []TranslationConflict `json:"conflicts"`
	Total         int                   `json:"total"`
	TotalResolved int                   `json:"total_resolved"`
	TotalPending  int                   `json:"total_pending"`
}
