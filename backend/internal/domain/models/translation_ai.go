package models

// AIProvider represents an AI translation provider configuration
type AIProvider struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // openai, anthropic, google, deepl
	APIKey      string  `json:"api_key,omitempty"`
	Endpoint    string  `json:"endpoint,omitempty"`
	Model       string  `json:"model,omitempty"`
	Enabled     bool    `json:"enabled"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
}

// TranslateRequest represents a single translation request
type TranslateRequest struct {
	Provider        string   `json:"provider" validate:"required"`
	Text            string   `json:"text" validate:"required"`
	Key             string   `json:"key" validate:"required"`
	Module          string   `json:"module" validate:"required"`
	SourceLanguage  string   `json:"source_language" validate:"required"`
	TargetLanguages []string `json:"target_languages" validate:"required,min=1"`
	Context         string   `json:"context,omitempty"`
}

// AIBatchTranslateRequest represents a batch translation request
type AIBatchTranslateRequest struct {
	Provider        string   `json:"provider" validate:"required"`
	Modules         []string `json:"modules" validate:"required,min=1"`
	SourceLanguage  string   `json:"source_language" validate:"required"`
	TargetLanguages []string `json:"target_languages" validate:"required,min=1"`
	MissingOnly     bool     `json:"missing_only"`
}

// TranslateResult represents the result of a translation
type TranslateResult struct {
	Key                     string              `json:"key"`
	Module                  string              `json:"module"`
	Translations            map[string]string   `json:"translations"`
	Provider                string              `json:"provider"`
	Confidence              float32             `json:"confidence,omitempty"`
	AlternativeTranslations map[string][]string `json:"alternative_translations,omitempty"`
}

// BatchTranslateResult represents the result of batch translation
type BatchTranslateResult struct {
	Results         []TranslateResult `json:"results"`
	TranslatedCount int               `json:"translated_count"`
	FailedCount     int               `json:"failed_count"`
	Errors          []string          `json:"errors,omitempty"`
}

// ApplyTranslationsRequest represents a request to apply AI translations
type ApplyTranslationsRequest struct {
	Translations []TranslationUpdate `json:"translations" validate:"required,dive"`
}

// TranslationUpdate represents a single translation update
type TranslationUpdate struct {
	Key      string `json:"key" validate:"required"`
	Module   string `json:"module" validate:"required"`
	Language string `json:"language" validate:"required"`
	Value    string `json:"value" validate:"required"`
}
