package service

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/marketplace/services"
)

// AICategoryMapper manages AI-powered category mapping for imports
type AICategoryMapper struct {
	aiDetector      *services.AICategoryDetector
	categoryService *CategoryMappingService
}

// NewAICategoryMapper creates a new AI category mapper
func NewAICategoryMapper(
	aiDetector *services.AICategoryDetector,
	categoryService *CategoryMappingService,
) *AICategoryMapper {
	return &AICategoryMapper{
		aiDetector:      aiDetector,
		categoryService: categoryService,
	}
}

// CategoryMappingSuggestion represents an AI suggestion for category mapping
type CategoryMappingSuggestion struct {
	ExternalCategory      string   `json:"external_category"`
	SuggestedCategoryID   int      `json:"suggested_category_id"`
	SuggestedCategoryPath string   `json:"suggested_category_path"`
	ConfidenceScore       float64  `json:"confidence_score"` // 0.0-1.0
	ReasoningSteps        []string `json:"reasoning_steps"`
	IsManualReviewNeeded  bool     `json:"is_manual_review_needed"`
}

// MapExternalCategory maps external category path to internal marketplace category
//
// Процесс:
// 1. Разбить external category на уровни (OPREMA > MASKE > SAMSUNG)
// 2. Для каждого уровня найти похожую категорию в маркетплейсе
// 3. Использовать AI для финального выбора
// 4. Вернуть с confidence score
func (m *AICategoryMapper) MapExternalCategory(
	ctx context.Context,
	externalCategory string, // "OPREMA ZA MOBILNI > MASKE > SAMSUNG"
	productName string, // Опционально: название товара для лучшего маппинга
	productDescription string, // Опционально: описание товара
) (*CategoryMappingSuggestion, error) {
	// 1. Нормализация пути категории
	levels := m.splitCategoryPath(externalCategory)
	if len(levels) == 0 {
		return nil, fmt.Errorf("empty category path")
	}

	logger.Debug().
		Str("external", externalCategory).
		Interface("levels", levels).
		Str("product", productName).
		Msg("Mapping external category")

	// 2. Использовать AI для определения категории
	// Комбинируем информацию: категория + название + описание
	searchText := strings.Join(levels, " ")
	if productName != "" {
		searchText = fmt.Sprintf("%s %s", searchText, productName)
	}
	if productDescription != "" {
		// Ограничиваем описание (первые 200 символов)
		if len(productDescription) > 200 {
			productDescription = productDescription[:200]
		}
		searchText = fmt.Sprintf("%s %s", searchText, productDescription)
	}

	// 3. AI детекция
	detectionInput := services.AIDetectionInput{
		Title:       searchText,
		Description: productDescription,
		EntityType:  "product",
	}
	detectionResult, err := m.aiDetector.DetectCategory(ctx, detectionInput)
	if err != nil {
		logger.Warn().
			Err(err).
			Str("external_category", externalCategory).
			Msg("AI category detection failed")
		return m.fallbackMapping(externalCategory, levels)
	}

	// 4. Определить нужен ли ручной review
	needsReview := detectionResult.ConfidenceScore < 0.85

	// 5. Подготовить reasoning steps
	reasoning := []string{
		fmt.Sprintf("Analyzed external category: %s", externalCategory),
		fmt.Sprintf("Combined with product context: %s", productName),
		fmt.Sprintf("AI suggested category ID: %d with confidence: %.2f", detectionResult.CategoryID, detectionResult.ConfidenceScore),
	}

	if needsReview {
		reasoning = append(reasoning, "⚠️ Low confidence - manual review recommended")
	}

	suggestion := &CategoryMappingSuggestion{
		ExternalCategory:      externalCategory,
		SuggestedCategoryID:   int(detectionResult.CategoryID),
		SuggestedCategoryPath: detectionResult.CategoryPath,
		ConfidenceScore:       detectionResult.ConfidenceScore,
		ReasoningSteps:        reasoning,
		IsManualReviewNeeded:  needsReview,
	}

	return suggestion, nil
}

// BatchMapCategories maps multiple categories at once
// Returns map: external_category -> suggestion
func (m *AICategoryMapper) BatchMapCategories(
	ctx context.Context,
	externalCategories []string,
	products map[string]models.ImportProductRequest, // external_category -> sample product
) (map[string]*CategoryMappingSuggestion, error) {
	results := make(map[string]*CategoryMappingSuggestion)

	for _, extCat := range externalCategories {
		// Получаем пример товара для лучшего маппинга
		var productName, productDesc string
		if product, ok := products[extCat]; ok {
			productName = product.Name
			productDesc = product.Description
		}

		suggestion, err := m.MapExternalCategory(ctx, extCat, productName, productDesc)
		if err != nil {
			logger.Warn().
				Str("external_category", extCat).
				Err(err).
				Msg("Failed to map category")
			continue
		}

		results[extCat] = suggestion
	}

	logger.Info().
		Int("total", len(externalCategories)).
		Int("mapped", len(results)).
		Msg("Batch category mapping completed")

	return results, nil
}

// AnalyzeMappingQuality анализирует качество маппинга и группирует по confidence
type MappingQuality struct {
	HighConfidence   []*CategoryMappingSuggestion `json:"high_confidence"`   // >= 0.90
	MediumConfidence []*CategoryMappingSuggestion `json:"medium_confidence"` // 0.70-0.90
	LowConfidence    []*CategoryMappingSuggestion `json:"low_confidence"`    // < 0.70
	TotalMapped      int                          `json:"total_mapped"`
	HighPercent      float64                      `json:"high_percent"`
	MediumPercent    float64                      `json:"medium_percent"`
	LowPercent       float64                      `json:"low_percent"`
}

func (m *AICategoryMapper) AnalyzeMappingQuality(
	suggestions map[string]*CategoryMappingSuggestion,
) *MappingQuality {
	quality := &MappingQuality{
		HighConfidence:   []*CategoryMappingSuggestion{},
		MediumConfidence: []*CategoryMappingSuggestion{},
		LowConfidence:    []*CategoryMappingSuggestion{},
	}

	for _, suggestion := range suggestions {
		switch {
		case suggestion.ConfidenceScore >= 0.90:
			quality.HighConfidence = append(quality.HighConfidence, suggestion)
		case suggestion.ConfidenceScore >= 0.70:
			quality.MediumConfidence = append(quality.MediumConfidence, suggestion)
		default:
			quality.LowConfidence = append(quality.LowConfidence, suggestion)
		}
	}

	total := len(suggestions)
	quality.TotalMapped = total

	if total > 0 {
		quality.HighPercent = float64(len(quality.HighConfidence)) / float64(total) * 100
		quality.MediumPercent = float64(len(quality.MediumConfidence)) / float64(total) * 100
		quality.LowPercent = float64(len(quality.LowConfidence)) / float64(total) * 100
	}

	return quality
}

// splitCategoryPath разбивает путь категории на уровни
// Поддерживает разделители: >, /, |, -
func (m *AICategoryMapper) splitCategoryPath(categoryPath string) []string {
	// Нормализуем разделители
	normalized := categoryPath
	for _, sep := range []string{" > ", " / ", " | ", " - "} {
		normalized = strings.ReplaceAll(normalized, sep, ">")
	}

	// Разделяем
	levels := strings.Split(normalized, ">")

	// Очистка
	result := make([]string, 0, len(levels))
	for _, level := range levels {
		trimmed := strings.TrimSpace(level)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// fallbackMapping возвращает fallback маппинг при ошибке AI
func (m *AICategoryMapper) fallbackMapping(externalCategory string, levels []string) (*CategoryMappingSuggestion, error) {
	// Используем дефолтную категорию (Elektronika)
	const defaultCategoryID = 1001
	const defaultCategoryPath = "Elektronika"

	return &CategoryMappingSuggestion{
		ExternalCategory:      externalCategory,
		SuggestedCategoryID:   defaultCategoryID,
		SuggestedCategoryPath: defaultCategoryPath,
		ConfidenceScore:       0.1, // Очень низкая confidence
		ReasoningSteps: []string{
			"AI detection failed",
			fmt.Sprintf("Fallback to default category: %s", defaultCategoryPath),
			"⚠️ Manual review required",
		},
		IsManualReviewNeeded: true,
	}, nil
}
