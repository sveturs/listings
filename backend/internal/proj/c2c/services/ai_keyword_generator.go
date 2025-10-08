package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// KeywordGenerationRequest represents request for keyword generation
type KeywordGenerationRequest struct {
	CategoryID   int32  `json:"categoryId" validate:"required"`
	CategoryName string `json:"categoryName" validate:"required"`
	Language     string `json:"language,omitempty"`
	MinKeywords  int    `json:"minKeywords,omitempty"`
}

// KeywordGenerationResult represents generated keywords
type KeywordGenerationResult struct {
	CategoryID       int32                    `json:"categoryId"`
	GeneratedCount   int                      `json:"generatedCount"`
	Keywords         []GeneratedKeyword       `json:"keywords"`
	ProcessingTimeMs int64                    `json:"processingTimeMs"`
	Categories       []CategoryKeywordMapping `json:"categories"`
}

// GeneratedKeyword represents a single generated keyword
type GeneratedKeyword struct {
	Keyword     string  `json:"keyword"`
	Type        string  `json:"type"`        // main, synonym, brand, attribute, context
	Weight      float64 `json:"weight"`      // 0.1-2.0
	Confidence  float64 `json:"confidence"`  // AI confidence in this keyword
	Description string  `json:"description"` // Why this keyword is relevant
}

// CategoryKeywordMapping for bulk operations
type CategoryKeywordMapping struct {
	CategoryID int32              `json:"categoryId"`
	Keywords   []GeneratedKeyword `json:"keywords"`
}

// AIKeywordGenerator generates comprehensive keywords for categories using AI
type AIKeywordGenerator struct {
	logger      *zap.Logger
	redisClient *redis.Client
	httpClient  *HTTPClient
	validator   *AICategoryValidator
}

// NewAIKeywordGenerator creates a new keyword generator
func NewAIKeywordGenerator(logger *zap.Logger, redisClient *redis.Client, validator *AICategoryValidator) *AIKeywordGenerator {
	return &AIKeywordGenerator{
		logger:      logger,
		redisClient: redisClient,
		httpClient:  NewHTTPClient(),
		validator:   validator,
	}
}

// GenerateKeywordsForCategory generates comprehensive keywords for a specific category
func (g *AIKeywordGenerator) GenerateKeywordsForCategory(ctx context.Context, req KeywordGenerationRequest) (*KeywordGenerationResult, error) {
	startTime := time.Now()

	// Check cache first
	cacheKey := fmt.Sprintf("ai_keywords:%d:%s", req.CategoryID, req.Language)
	if cached, err := g.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		g.logger.Info("Keyword generation cache hit", zap.String("cacheKey", cacheKey))
		return cached, nil
	}

	// Set minimum keywords if not specified
	if req.MinKeywords == 0 {
		req.MinKeywords = 50
	}

	// Create keyword generation prompt
	prompt := g.createKeywordPrompt(req)

	// Call AI service
	aiResponse, err := g.callAIKeywordGeneration(ctx, prompt)
	if err != nil {
		g.logger.Error("AI keyword generation failed", zap.Error(err))
		return nil, fmt.Errorf("failed to generate keywords: %w", err)
	}

	// Parse AI response
	result, err := g.parseKeywordResponse(aiResponse, req.CategoryID)
	if err != nil {
		g.logger.Error("Failed to parse AI keyword response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse keywords: %w", err)
	}

	// Add processing time
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()

	// Cache result for 7 days
	g.cacheResult(ctx, cacheKey, result, 7*24*time.Hour)

	g.logger.Info("Keywords generated successfully",
		zap.Int32("categoryId", req.CategoryID),
		zap.String("categoryName", req.CategoryName),
		zap.Int("generatedCount", result.GeneratedCount),
		zap.Int64("processingTime", result.ProcessingTimeMs))

	return result, nil
}

// GenerateKeywordsForAllCategories generates keywords for all categories that need them
func (g *AIKeywordGenerator) GenerateKeywordsForAllCategories(ctx context.Context, categories []Category) (*KeywordGenerationResult, error) {
	startTime := time.Now()

	var allResults []CategoryKeywordMapping
	totalGenerated := 0

	for _, category := range categories {
		req := KeywordGenerationRequest{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			Language:     "ru",
			MinKeywords:  50,
		}

		result, err := g.GenerateKeywordsForCategory(ctx, req)
		if err != nil {
			g.logger.Error("Failed to generate keywords for category",
				zap.Int32("categoryId", category.ID),
				zap.String("categoryName", category.Name),
				zap.Error(err))
			continue
		}

		allResults = append(allResults, CategoryKeywordMapping{
			CategoryID: category.ID,
			Keywords:   result.Keywords,
		})
		totalGenerated += result.GeneratedCount

		// Small delay to avoid overwhelming AI service
		time.Sleep(2 * time.Second)
	}

	return &KeywordGenerationResult{
		GeneratedCount:   totalGenerated,
		Keywords:         []GeneratedKeyword{}, // Individual keywords not included in bulk
		Categories:       allResults,
		ProcessingTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// createKeywordPrompt creates AI prompt for keyword generation
func (g *AIKeywordGenerator) createKeywordPrompt(req KeywordGenerationRequest) string {
	prompt := fmt.Sprintf(`
Создай максимально полный список ключевых слов для категории товаров "%s".

ТРЕБОВАНИЯ:
- Минимум %d ключевых слов
- Включи ВСЕ возможные варианты названий товаров этой категории
- Добавь популярные бренды (Bosch, Makita, DeWalt, Apple, Samsung и др.)
- Включи синонимы и альтернативные названия
- Добавь технические термины и профессиональный сленг
- Включи характеристики и атрибуты товаров
- Добавь сезонные и региональные названия

ТИПЫ КЛЮЧЕВЫХ СЛОВ:
1. "main" - основные названия товаров (вес 2.0)
2. "synonym" - синонимы и альтернативы (вес 1.5)
3. "brand" - названия брендов (вес 1.8)
4. "attribute" - характеристики и свойства (вес 1.2)
5. "context" - контекстные слова (вес 1.0)

Ответь строго в формате JSON:
{
    "keywords": [
        {
            "keyword": "название",
            "type": "main|synonym|brand|attribute|context",
            "weight": 1.0-2.0,
            "confidence": 0.0-1.0,
            "description": "почему это слово подходит для категории"
        }
    ],
    "totalCount": 123,
    "categoryAnalysis": "краткий анализ категории и стратегии подбора слов"
}

ПРИМЕРЫ для строительных инструментов:
- main: "болгарка", "углошлифовальная машина", "дрель", "перфоратор"
- synonym: "УШМ", "бормашина", "дрель-шуруповерт"
- brand: "Bosch", "Makita", "DeWalt", "Hitachi", "Hilti"
- attribute: "аккумуляторная", "сетевая", "ударная", "безударная"
- context: "строительство", "ремонт", "профессиональный", "бытовой"

КАТЕГОРИЯ: "%s"
ЯЗЫК: %s

Ответ должен быть валидным JSON без дополнительного текста.
`, req.CategoryName, req.MinKeywords, req.CategoryName, getLanguageName(req.Language))

	return prompt
}

// callAIKeywordGeneration calls AI service for keyword generation
func (g *AIKeywordGenerator) callAIKeywordGeneration(ctx context.Context, prompt string) (string, error) {
	requestBody := map[string]interface{}{
		"prompt":      prompt,
		"maxTokens":   3000, // More tokens for comprehensive keyword lists
		"temperature": 0.3,  // Some creativity but still focused
	}

	// Call AI service (through frontend API route)
	response, err := g.httpClient.PostJSON(ctx, "http://localhost:3001/api/ai/validate", requestBody)
	if err != nil {
		return "", fmt.Errorf("AI keyword generation request failed: %w", err)
	}

	return response, nil
}

// parseKeywordResponse parses AI response into structured keywords
func (g *AIKeywordGenerator) parseKeywordResponse(response string, categoryID int32) (*KeywordGenerationResult, error) {
	// Find JSON in response
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

	// Parse JSON structure
	var parsedResponse struct {
		Keywords         []GeneratedKeyword `json:"keywords"`
		TotalCount       int                `json:"totalCount"`
		CategoryAnalysis string             `json:"categoryAnalysis"`
	}

	err := json.Unmarshal([]byte(jsonStr), &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate and normalize keywords
	validKeywords := make([]GeneratedKeyword, 0)
	for _, kw := range parsedResponse.Keywords {
		if g.validateKeyword(kw) {
			// Normalize values
			kw.Weight = g.normalizeWeight(kw.Weight, kw.Type)
			kw.Confidence = g.normalizeConfidence(kw.Confidence)
			kw.Keyword = strings.ToLower(strings.TrimSpace(kw.Keyword))

			validKeywords = append(validKeywords, kw)
		}
	}

	return &KeywordGenerationResult{
		CategoryID:     categoryID,
		GeneratedCount: len(validKeywords),
		Keywords:       validKeywords,
	}, nil
}

// validateKeyword validates a generated keyword
func (g *AIKeywordGenerator) validateKeyword(kw GeneratedKeyword) bool {
	if len(strings.TrimSpace(kw.Keyword)) < 2 {
		return false
	}

	validTypes := map[string]bool{
		"main": true, "synonym": true, "brand": true,
		"attribute": true, "context": true,
	}

	return validTypes[kw.Type]
}

// normalizeWeight normalizes keyword weight based on type
func (g *AIKeywordGenerator) normalizeWeight(weight float64, keywordType string) float64 {
	if weight < 0.1 {
		weight = 0.1
	}
	if weight > 2.0 {
		weight = 2.0
	}

	// Default weights by type
	defaults := map[string]float64{
		"main":      2.0,
		"synonym":   1.5,
		"brand":     1.8,
		"attribute": 1.2,
		"context":   1.0,
	}

	if weight == 0 {
		if defaultWeight, exists := defaults[keywordType]; exists {
			return defaultWeight
		}
		return 1.0
	}

	return weight
}

// normalizeConfidence normalizes confidence value
func (g *AIKeywordGenerator) normalizeConfidence(confidence float64) float64 {
	if confidence < 0 {
		return 0
	}
	if confidence > 1 {
		return 1
	}
	if confidence == 0 {
		return 0.8 // Default confidence
	}
	return confidence
}

// getFromCache retrieves keyword result from cache
func (g *AIKeywordGenerator) getFromCache(ctx context.Context, key string) (*KeywordGenerationResult, error) {
	data, err := g.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result KeywordGenerationResult
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// cacheResult stores keyword result in cache
func (g *AIKeywordGenerator) cacheResult(ctx context.Context, key string, result *KeywordGenerationResult, duration time.Duration) {
	data, err := json.Marshal(result)
	if err != nil {
		g.logger.Error("Failed to marshal keyword result for cache", zap.Error(err))
		return
	}

	err = g.redisClient.Set(ctx, key, data, duration).Err()
	if err != nil {
		g.logger.Error("Failed to cache keyword result", zap.Error(err))
	}
}

// getLanguageName returns language name for prompts
func getLanguageName(lang string) string {
	switch lang {
	case "ru":
		return "Русский"
	case "en":
		return "English"
	case "sr":
		return "Српски"
	default:
		return "Русский"
	}
}

// Category represents a marketplace category
type Category struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
