package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// ValidationResult represents the result of AI category validation
type ValidationResult struct {
	IsCorrect         bool     `json:"isCorrect"`
	Confidence        float64  `json:"confidence"`
	Reasoning         string   `json:"reasoning"`
	SuggestedCategory string   `json:"suggestedCategory"`
	Keywords          []string `json:"keywords"`
	ProcessingTimeMs  int64    `json:"processingTimeMs"`
}

// ValidationRequest represents the request for category validation
type ValidationRequest struct {
	Title        string `json:"title" validate:"required"`
	Description  string `json:"description"`
	CategoryName string `json:"categoryName" validate:"required"`
	Language     string `json:"language,omitempty"`
}

// AICategoryValidator validates category choices using AI
type AICategoryValidator struct {
	logger      *zap.Logger
	redisClient *redis.Client
	httpClient  *HTTPClient
}

// NewAICategoryValidator creates a new AI category validator
func NewAICategoryValidator(logger *zap.Logger, redisClient *redis.Client) *AICategoryValidator {
	return &AICategoryValidator{
		logger:      logger,
		redisClient: redisClient,
		httpClient:  NewHTTPClient(),
	}
}

// ValidateCategory validates if a product belongs to the suggested category using AI
func (v *AICategoryValidator) ValidateCategory(ctx context.Context, req ValidationRequest) (*ValidationResult, error) {
	startTime := time.Now()

	// Check cache first
	cacheKey := fmt.Sprintf("ai_validation:%s:%s", hashString(req.Title), hashString(req.CategoryName))
	if cached, err := v.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		v.logger.Info("AI validation cache hit", zap.String("cacheKey", cacheKey))
		return cached, nil
	}

	// Create validation prompt
	prompt := v.createValidationPrompt(req)

	// Call AI service (Claude via frontend API)
	aiResponse, err := v.callAIValidation(ctx, prompt)
	if err != nil {
		v.logger.Error("AI validation failed", zap.Error(err))
		return v.createFallbackResult(), nil
	}

	// Parse AI response
	result, err := v.parseValidationResponse(aiResponse)
	if err != nil {
		v.logger.Error("Failed to parse AI validation response", zap.Error(err))
		return v.createFallbackResult(), nil
	}

	// Add processing time
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	// Cache result for 24 hours
	v.cacheResult(ctx, cacheKey, result, 24*time.Hour)

	v.logger.Info("AI validation completed",
		zap.String("title", req.Title),
		zap.String("category", req.CategoryName),
		zap.Bool("isCorrect", result.IsCorrect),
		zap.Float64("confidence", result.Confidence),
		zap.Int64("processingTime", result.ProcessingTimeMs))

	return result, nil
}

// createValidationPrompt creates a prompt for AI validation
func (v *AICategoryValidator) createValidationPrompt(req ValidationRequest) string {
	prompt := fmt.Sprintf(`
Проанализируй, правильно ли товар отнесен к данной категории.

ТОВАР: "%s"
ОПИСАНИЕ: "%s"
ПРЕДЛОЖЕННАЯ КАТЕГОРИЯ: "%s"

Проанализируй и ответь строго в формате JSON:
{
    "isCorrect": true/false,
    "confidence": 0.0-1.0,
    "reasoning": "краткое объяснение почему подходит или не подходит",
    "suggestedCategory": "лучшая категория если текущая неправильная, иначе null",
    "keywords": ["5-10", "ключевых", "слов", "для", "этой", "категории", "товаров"]
}

ВАЖНО:
- isCorrect: true только если категория точно подходит товару
- confidence: насколько уверен в оценке (0.0-1.0)
- suggestedCategory: предложи лучшую категорию только если текущая явно неправильная
- keywords: ключевые слова которые помогут в будущем определять такие товары
- reasoning: краткое объяснение решения

Ответ должен быть валидным JSON без дополнительного текста.
`, req.Title, req.Description, req.CategoryName)

	return prompt
}

// callAIValidation calls the AI service to validate category
func (v *AICategoryValidator) callAIValidation(ctx context.Context, prompt string) (string, error) {
	// Use the existing AI analyze endpoint
	requestBody := map[string]interface{}{
		"prompt":      prompt,
		"maxTokens":   500,
		"temperature": 0.1, // Low temperature for consistent validation
	}

	// Call the AI service (through frontend API route)
	response, err := v.httpClient.PostJSON(ctx, "http://localhost:3001/api/ai/validate", requestBody)
	if err != nil {
		return "", fmt.Errorf("AI validation request failed: %w", err)
	}

	return response, nil
}

// parseValidationResponse parses the AI validation response
func (v *AICategoryValidator) parseValidationResponse(response string) (*ValidationResult, error) {
	var result ValidationResult

	// Try to extract JSON from response
	jsonStart := findJSONStart(response)
	if jsonStart == -1 {
		return nil, fmt.Errorf("no JSON found in AI response")
	}

	jsonStr := response[jsonStart:]
	jsonEnd := findJSONEnd(jsonStr)
	if jsonEnd == -1 {
		return nil, fmt.Errorf("incomplete JSON in AI response")
	}

	jsonStr = jsonStr[:jsonEnd+1]

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate result
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 1 {
		result.Confidence = 1
	}

	return &result, nil
}

// createFallbackResult creates a fallback result when AI validation fails
func (v *AICategoryValidator) createFallbackResult() *ValidationResult {
	return &ValidationResult{
		IsCorrect:         true, // Assume correct to not block user
		Confidence:        0.5,  // Medium confidence
		Reasoning:         "AI validation unavailable, assuming correct",
		SuggestedCategory: "",
		Keywords:          []string{},
		ProcessingTimeMs:  0,
	}
}

// getFromCache retrieves validation result from cache
func (v *AICategoryValidator) getFromCache(ctx context.Context, key string) (*ValidationResult, error) {
	data, err := v.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result ValidationResult
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// cacheResult stores validation result in cache
func (v *AICategoryValidator) cacheResult(ctx context.Context, key string, result *ValidationResult, duration time.Duration) {
	data, err := json.Marshal(result)
	if err != nil {
		v.logger.Error("Failed to marshal validation result for cache", zap.Error(err))
		return
	}

	err = v.redisClient.Set(ctx, key, data, duration).Err()
	if err != nil {
		v.logger.Error("Failed to cache validation result", zap.Error(err))
	}
}

// Helper functions for JSON parsing
func findJSONStart(text string) int {
	for i, char := range text {
		if char == '{' {
			return i
		}
	}
	return -1
}

func findJSONEnd(text string) int {
	braceCount := 0
	for i, char := range text {
		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount == 0 {
				return i
			}
		}
	}
	return -1
}

// hashString creates a simple hash for caching
func hashString(s string) string {
	if len(s) > 50 {
		return s[:50]
	}
	return s
}
