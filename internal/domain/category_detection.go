package domain

import (
	"time"

	"github.com/google/uuid"
)

// DetectionMethod - метод определения категории
type DetectionMethod string

const (
	MethodBrandMatch   DetectionMethod = "brand_match"
	MethodAIClaude     DetectionMethod = "ai_claude"
	MethodKeywordMatch DetectionMethod = "keyword_match"
	MethodSimilarity   DetectionMethod = "similarity"
	MethodFallback     DetectionMethod = "fallback"
)

// CategoryMatch - результат детекции категории
type CategoryMatch struct {
	CategoryID      uuid.UUID       `json:"category_id"`
	CategoryName    string          `json:"category_name"`
	CategorySlug    string          `json:"category_slug"`
	CategoryPath    string          `json:"category_path"`
	ConfidenceScore float64         `json:"confidence_score"`
	DetectionMethod DetectionMethod `json:"detection_method"`
	MatchedKeywords []string        `json:"matched_keywords,omitempty"`
}

// CategoryDetection - результат детекции с tracking информацией
type CategoryDetection struct {
	ID               uuid.UUID        `json:"id"`
	Primary          *CategoryMatch   `json:"primary"`
	Alternatives     []CategoryMatch  `json:"alternatives,omitempty"`
	ProcessingTimeMs int32            `json:"processing_time_ms"`

	// Входные данные (для tracking)
	InputTitle       string `json:"input_title"`
	InputDescription string `json:"input_description,omitempty"`
	InputLanguage    string `json:"input_language"`

	// Подтверждение пользователя (для обучения)
	UserConfirmed  *bool      `json:"user_confirmed,omitempty"`
	UserSelectedID *uuid.UUID `json:"user_selected_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// CategoryHints - подсказки для улучшения детекции
type CategoryHints struct {
	Domain      string   `json:"domain,omitempty"`       // "electronics", "fashion", etc.
	ProductType string   `json:"product_type,omitempty"` // "smartphone", "dress", etc.
	Keywords    []string `json:"keywords,omitempty"`     // дополнительные ключевые слова
}

// DetectFromTextInput - входные данные для детекции по тексту
type DetectFromTextInput struct {
	Title             string         `json:"title"`
	Description       string         `json:"description,omitempty"`
	Language          string         `json:"language"` // "sr" | "en" | "ru"
	Hints             *CategoryHints `json:"hints,omitempty"`
	SuggestedCategory string         `json:"suggested_category,omitempty"` // AI-предложенная категория (например "Furniture")
}

// DetectBatchInput - входные данные для batch детекции
type DetectBatchInput struct {
	Items []DetectFromTextInput `json:"items"`
}

// DetectBatchResult - результат batch детекции
type DetectBatchResult struct {
	Results             []CategoryDetection `json:"results"`
	TotalProcessingTime int32               `json:"total_processing_time_ms"`
}

// NewCategoryDetection создаёт новый CategoryDetection с UUID
func NewCategoryDetection(primary *CategoryMatch, alternatives []CategoryMatch, processingMs int32, input DetectFromTextInput) *CategoryDetection {
	return &CategoryDetection{
		ID:               uuid.New(),
		Primary:          primary,
		Alternatives:     alternatives,
		ProcessingTimeMs: processingMs,
		InputTitle:       input.Title,
		InputDescription: input.Description,
		InputLanguage:    input.Language,
		CreatedAt:        time.Now(),
	}
}

// NewCategoryMatch создаёт CategoryMatch
func NewCategoryMatch(categoryID uuid.UUID, name, slug, path string, score float64, method DetectionMethod, keywords []string) *CategoryMatch {
	return &CategoryMatch{
		CategoryID:      categoryID,
		CategoryName:    name,
		CategorySlug:    slug,
		CategoryPath:    path,
		ConfidenceScore: score,
		DetectionMethod: method,
		MatchedKeywords: keywords,
	}
}
