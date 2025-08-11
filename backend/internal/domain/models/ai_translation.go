package models

// AITranslateRequest represents a single translation request
type AITranslateRequest struct {
	Text       string `json:"text" validate:"required"`
	SourceLang string `json:"source_lang" validate:"required"`
	TargetLang string `json:"target_lang" validate:"required"`
	Provider   string `json:"provider,omitempty"`
}

// AITranslateResponse represents a translation response
type AITranslateResponse struct {
	Translation string  `json:"translation"`
	Confidence  float64 `json:"confidence"`
	Provider    string  `json:"provider"`
}

// AITranslateBatchRequest represents a batch translation request
type AITranslateBatchRequest struct {
	Items       []TranslationItem `json:"items" validate:"required,min=1"`
	SourceLang  string            `json:"source_lang" validate:"required"`
	TargetLangs []string          `json:"target_langs" validate:"required,min=1"`
	Provider    string            `json:"provider,omitempty"`
}

// TranslationItem represents an item to translate
type TranslationItem struct {
	Key    string `json:"key" validate:"required"`
	Module string `json:"module,omitempty"`
	Text   string `json:"text" validate:"required"`
}

// AITranslateBatchResponse represents batch translation response
type AITranslateBatchResponse struct {
	Results      []AITranslationResult `json:"results"`
	TotalItems   int                   `json:"total_items"`
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
}

// AITranslationResult represents a single translation result
type AITranslationResult struct {
	Key          string            `json:"key"`
	Module       string            `json:"module,omitempty"`
	Translations map[string]string `json:"translations"`
	Provider     string            `json:"provider"`
	Confidence   float64           `json:"confidence,omitempty"`
	Error        string            `json:"error,omitempty"`
}

// TranslateModuleRequest represents a request to translate entire module
type TranslateModuleRequest struct {
	Module      string   `json:"module" validate:"required"`
	SourceLang  string   `json:"source_lang" validate:"required"`
	TargetLangs []string `json:"target_langs" validate:"required,min=1"`
	Provider    string   `json:"provider,omitempty"`
	OnlyMissing bool     `json:"only_missing,omitempty"`
}

// TranslateModuleResponse represents module translation response
type TranslateModuleResponse struct {
	Module         string                `json:"module"`
	Results        []AITranslationResult `json:"results"`
	TotalKeys      int                   `json:"total_keys"`
	TranslatedKeys int                   `json:"translated_keys"`
}
