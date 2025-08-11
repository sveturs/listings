package models

import (
	"time"
)

// TranslationVersion represents a version of a translation
type TranslationVersion struct {
	ID             int                    `json:"id" db:"id"`
	TranslationID  int                    `json:"translation_id" db:"translation_id"`
	VersionNumber  int                    `json:"version_number" db:"version_number"`
	EntityType     string                 `json:"entity_type" db:"entity_type"`
	EntityID       int                    `json:"entity_id" db:"entity_id"`
	Language       string                 `json:"language" db:"language"`
	FieldName      string                 `json:"field_name" db:"field_name"`
	TranslatedText string                 `json:"translated_text" db:"translated_text"`
	ChangedBy      *int                   `json:"changed_by" db:"changed_by"`
	ChangedAt      time.Time              `json:"changed_at" db:"changed_at"`
	ChangeComment  *string                `json:"change_comment" db:"change_comment"`
	Metadata       map[string]interface{} `json:"metadata" db:"metadata"`
}

// TranslationSyncConflict represents a conflict during synchronization
type TranslationSyncConflict struct {
	ID               int        `json:"id" db:"id"`
	SourceType       string     `json:"source_type" db:"source_type"`
	TargetType       string     `json:"target_type" db:"target_type"`
	EntityIdentifier string     `json:"entity_identifier" db:"entity_identifier"`
	SourceValue      *string    `json:"source_value" db:"source_value"`
	TargetValue      *string    `json:"target_value" db:"target_value"`
	ConflictType     string     `json:"conflict_type" db:"conflict_type"`
	Resolved         bool       `json:"resolved" db:"resolved"`
	ResolvedBy       *int       `json:"resolved_by" db:"resolved_by"`
	ResolvedAt       *time.Time `json:"resolved_at" db:"resolved_at"`
	ResolutionType   *string    `json:"resolution_type" db:"resolution_type"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// TranslationAuditLog represents an audit log entry
type TranslationAuditLog struct {
	ID         int       `json:"id" db:"id"`
	UserID     *int      `json:"user_id" db:"user_id"`
	Action     string    `json:"action" db:"action"`
	EntityType *string   `json:"entity_type" db:"entity_type"`
	EntityID   *int      `json:"entity_id" db:"entity_id"`
	OldValue   *string   `json:"old_value" db:"old_value"`
	NewValue   *string   `json:"new_value" db:"new_value"`
	IPAddress  *string   `json:"ip_address" db:"ip_address"`
	UserAgent  *string   `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// TranslationProvider represents an AI translation provider
type TranslationProvider struct {
	ID           int                    `json:"id" db:"id"`
	Name         string                 `json:"name" db:"name"`
	ProviderType string                 `json:"provider_type" db:"provider_type"`
	APIKey       *string                `json:"-" db:"api_key"` // Hidden from JSON
	Settings     map[string]interface{} `json:"settings" db:"settings"`
	UsageLimit   *int                   `json:"usage_limit" db:"usage_limit"`
	UsageCurrent int                    `json:"usage_current" db:"usage_current"`
	IsActive     bool                   `json:"is_active" db:"is_active"`
	Priority     int                    `json:"priority" db:"priority"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// TranslationQualityMetrics represents quality metrics for a translation
type TranslationQualityMetrics struct {
	ID              int           `json:"id" db:"id"`
	TranslationID   int           `json:"translation_id" db:"translation_id"`
	QualityScore    *float64      `json:"quality_score" db:"quality_score"`
	CharacterCount  *int          `json:"character_count" db:"character_count"`
	WordCount       *int          `json:"word_count" db:"word_count"`
	HasPlaceholders bool          `json:"has_placeholders" db:"has_placeholders"`
	HasHTMLTags     bool          `json:"has_html_tags" db:"has_html_tags"`
	CheckedAt       time.Time     `json:"checked_at" db:"checked_at"`
	CheckedBy       string        `json:"checked_by" db:"checked_by"`
	Issues          []interface{} `json:"issues" db:"issues"`
}

// TranslationTask represents a translation task
type TranslationTask struct {
	ID               int                    `json:"id" db:"id"`
	TaskType         string                 `json:"task_type" db:"task_type"`
	Status           string                 `json:"status" db:"status"`
	SourceLanguage   *string                `json:"source_language" db:"source_language"`
	TargetLanguages  []string               `json:"target_languages" db:"target_languages"`
	EntityReferences []interface{}          `json:"entity_references" db:"entity_references"`
	ProviderID       *int                   `json:"provider_id" db:"provider_id"`
	CreatedBy        *int                   `json:"created_by" db:"created_by"`
	AssignedTo       *int                   `json:"assigned_to" db:"assigned_to"`
	StartedAt        *time.Time             `json:"started_at" db:"started_at"`
	CompletedAt      *time.Time             `json:"completed_at" db:"completed_at"`
	ErrorMessage     *string                `json:"error_message" db:"error_message"`
	Metadata         map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
}

// FrontendTranslation represents a translation from frontend JSON files
type FrontendTranslation struct {
	Module       string              `json:"module"`
	Key          string              `json:"key"`
	Path         string              `json:"path"`         // Full path like "marketplace.filters.priceFrom"
	Translations map[string]string   `json:"translations"` // lang -> text
	Status       TranslationStatus   `json:"status"`
	Metadata     TranslationMetadata `json:"metadata"`
}

// TranslationStatus represents the status of a translation
type TranslationStatus string

const (
	StatusComplete    TranslationStatus = "complete"
	StatusIncomplete  TranslationStatus = "incomplete"
	StatusPlaceholder TranslationStatus = "placeholder"
	StatusMissing     TranslationStatus = "missing"
	StatusOutdated    TranslationStatus = "outdated"
)

// TranslationMetadata contains metadata about a translation
type TranslationMetadata struct {
	Provider       string     `json:"provider,omitempty"`
	TranslatedBy   *int       `json:"translated_by,omitempty"`
	TranslatedAt   *time.Time `json:"translated_at,omitempty"`
	VerifiedBy     *int       `json:"verified_by,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	CharacterCount int        `json:"character_count,omitempty"`
	WordCount      int        `json:"word_count,omitempty"`
	QualityScore   float64    `json:"quality_score,omitempty"`
}

// FrontendModule represents a frontend translation module
type FrontendModule struct {
	Name         string                 `json:"name"`
	Keys         int                    `json:"keys"`
	Complete     int                    `json:"complete"`
	Incomplete   int                    `json:"incomplete"`
	Placeholders int                    `json:"placeholders"`
	Missing      int                    `json:"missing"`
	Languages    map[string]ModuleStats `json:"languages"`
}

// ModuleStats represents statistics for a module per language
type ModuleStats struct {
	Total        int `json:"total"`
	Complete     int `json:"complete"`
	Incomplete   int `json:"incomplete"`
	Placeholders int `json:"placeholders"`
	Missing      int `json:"missing"`
}

// TranslationStatistics represents overall translation statistics
type TranslationStatistics struct {
	TotalTranslations    int                      `json:"total_translations"`
	CompleteTranslations int                      `json:"complete_translations"`
	MissingTranslations  int                      `json:"missing_translations"`
	PlaceholderCount     int                      `json:"placeholder_count"`
	LanguageStats        map[string]LanguageStats `json:"language_stats"`
	ModuleStats          []FrontendModule         `json:"module_stats"`
	RecentChanges        []TranslationAuditLog    `json:"recent_changes"`
}

// LanguageStats represents statistics per language
type LanguageStats struct {
	Total             int     `json:"total"`
	Complete          int     `json:"complete"`
	MachineTranslated int     `json:"machine_translated"`
	Verified          int     `json:"verified"`
	Coverage          float64 `json:"coverage"` // Percentage
}

// SyncStatus represents the status of synchronization
type SyncStatus struct {
	InProgress       bool                   `json:"in_progress"`
	LastSync         *time.Time             `json:"last_sync"`
	Conflicts        int                    `json:"conflicts"`
	PendingConflicts int                    `json:"pending_conflicts"`
	FrontendModified int                    `json:"frontend_modified"`
	DatabaseModified int                    `json:"database_modified"`
	Details          map[string]interface{} `json:"details"`
}

// BatchTranslateRequest represents a batch translation request
type BatchTranslateRequest struct {
	Items           []TranslateItem `json:"items"`
	SourceLanguage  string          `json:"source_language,omitempty"`
	TargetLanguages []string        `json:"target_languages"`
	ProviderID      *int            `json:"provider_id,omitempty"`
	AutoApprove     bool            `json:"auto_approve"`
}

// TranslateItem represents a single item to translate
type TranslateItem struct {
	Module  string `json:"module,omitempty"`
	Key     string `json:"key"`
	Text    string `json:"text"`
	Context string `json:"context,omitempty"`
}

// ValidateTranslationsRequest represents a validation request
type ValidateTranslationsRequest struct {
	Module      string   `json:"module,omitempty"`
	Languages   []string `json:"languages,omitempty"`
	CheckHTML   bool     `json:"check_html"`
	CheckVars   bool     `json:"check_variables"`
	CheckLength bool     `json:"check_length"`
}

// ValidationResult represents validation results
type ValidationResult struct {
	Module string            `json:"module"`
	Key    string            `json:"key"`
	Issues []ValidationIssue `json:"issues"`
}

// ValidationIssue represents a single validation issue
type ValidationIssue struct {
	Type     string `json:"type"` // "missing", "placeholder", "variable_mismatch", "length", "html"
	Language string `json:"language"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // "error", "warning", "info"
}

// TranslationUpdateRequest represents a request to update a translation
type TranslationUpdateRequest struct {
	TranslatedText      string                 `json:"translated_text"`
	IsVerified          *bool                  `json:"is_verified,omitempty"`
	IsMachineTranslated *bool                  `json:"is_machine_translated,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// BatchOperationsRequest represents a batch operations request
type BatchOperationsRequest struct {
	Create []Translation `json:"create,omitempty"`
	Update []Translation `json:"update,omitempty"`
	Delete []int         `json:"delete,omitempty"`
}

// BatchOperationsResult represents the result of batch operations
type BatchOperationsResult struct {
	Created int      `json:"created"`
	Updated int      `json:"updated"`
	Deleted int      `json:"deleted"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

// ImportTranslationsRequest represents a request to import translations
type ImportTranslationsRequest struct {
	Translations      map[string]interface{} `json:"translations"`
	OverwriteExisting bool                   `json:"overwrite_existing"`
}

// ImportResult represents the result of import operation
type ImportResult struct {
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Skipped int `json:"skipped"`
}

// SyncResult represents the result of synchronization
type SyncResult struct {
	Added      int `json:"added"`
	Updated    int `json:"updated"`
	Conflicts  int `json:"conflicts"`
	TotalItems int `json:"total_items"`
}
