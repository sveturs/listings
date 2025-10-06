package service

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/logger"
)

// CategoryInsight представляет insight о категории клиента
type CategoryInsight struct {
	ExternalCategory     string               `json:"external_category"`
	ProductCount         int                  `json:"product_count"`
	Importance           float64              `json:"importance"` // 0-1, основано на количестве товаров
	IsUnique             bool                 `json:"is_unique"`  // Нет аналога в нашей системе
	SuggestedNewCategory *NewCategoryProposal `json:"suggested_new_category,omitempty"`
}

// NewCategoryProposal предложение новой категории
type NewCategoryProposal struct {
	ParentCategoryID  int               `json:"parent_category_id"`
	Name              string            `json:"name"`
	NameTranslations  map[string]string `json:"name_translations"` // ru, en, sr
	Description       string            `json:"description"`
	Reasoning         string            `json:"reasoning"`
	ExpectedProducts  int               `json:"expected_products"`
	SimilarCategories []int             `json:"similar_categories"` // Связанные категории
	Tags              []string          `json:"tags"`
}

// CategoryAnalysisResult результат анализа категорий клиента
type CategoryAnalysisResult struct {
	TotalCategories       int                   `json:"total_categories"`
	UniqueCategoriesFound int                   `json:"unique_categories_found"`
	Insights              []CategoryInsight     `json:"insights"`
	NewCategoryProposals  []NewCategoryProposal `json:"new_category_proposals"`
}

// AICategoryAnalyzer анализирует категории клиента и предлагает новые
type AICategoryAnalyzer struct {
	mapper          *AICategoryMapper
	categoryService *CategoryMappingService
}

// NewAICategoryAnalyzer creates new category analyzer
func NewAICategoryAnalyzer(
	mapper *AICategoryMapper,
	categoryService *CategoryMappingService,
) *AICategoryAnalyzer {
	return &AICategoryAnalyzer{
		mapper:          mapper,
		categoryService: categoryService,
	}
}

// ClientCategoryInfo информация о категории клиента
type ClientCategoryInfo struct {
	Path           string
	ProductCount   int
	SampleProducts []string // Примеры названий товаров
}

// AnalyzeClientCategories анализирует категории клиента и находит уникальные/важные
//
// Процесс:
// 1. Мапит все категории через AI
// 2. Находит категории без хорошего маппинга (низкая confidence)
// 3. Для важных категорий (много товаров) предлагает создание новых
func (a *AICategoryAnalyzer) AnalyzeClientCategories(
	ctx context.Context,
	clientCategories []ClientCategoryInfo,
) (*CategoryAnalysisResult, error) {
	result := &CategoryAnalysisResult{
		TotalCategories:      len(clientCategories),
		Insights:             []CategoryInsight{},
		NewCategoryProposals: []NewCategoryProposal{},
	}

	logger.Info().
		Int("total", len(clientCategories)).
		Msg("Analyzing client categories")

	// 1. Мапим все категории
	for _, category := range clientCategories {
		// Получаем sample product для лучшего маппинга
		var sampleProduct string
		if len(category.SampleProducts) > 0 {
			sampleProduct = category.SampleProducts[0]
		}

		suggestion, err := a.mapper.MapExternalCategory(ctx, category.Path, sampleProduct, "")
		if err != nil {
			logger.Warn().
				Str("category", category.Path).
				Err(err).
				Msg("Failed to map category during analysis")
			continue
		}

		// 2. Определяем важность категории
		importance := a.calculateImportance(category)

		// 3. Проверяем уникальность (низкая confidence = возможно уникальная)
		isUnique := suggestion.ConfidenceScore < 0.5

		insight := CategoryInsight{
			ExternalCategory: category.Path,
			ProductCount:     category.ProductCount,
			Importance:       importance,
			IsUnique:         isUnique,
		}

		// 4. Для важных и уникальных категорий - предлагаем создание новой
		if isUnique && category.ProductCount >= 20 {
			proposal := a.generateNewCategoryProposal(category, suggestion)
			insight.SuggestedNewCategory = &proposal
			result.NewCategoryProposals = append(result.NewCategoryProposals, proposal)
			result.UniqueCategoriesFound++
		}

		result.Insights = append(result.Insights, insight)
	}

	logger.Info().
		Int("total", result.TotalCategories).
		Int("unique", result.UniqueCategoriesFound).
		Int("proposals", len(result.NewCategoryProposals)).
		Msg("Category analysis completed")

	return result, nil
}

// calculateImportance вычисляет важность категории (0-1)
// Основано на количестве товаров
func (a *AICategoryAnalyzer) calculateImportance(category ClientCategoryInfo) float64 {
	// Простая формула: важность растет с количеством товаров
	// 0-10 товаров: 0.1-0.3 (низкая)
	// 10-50 товаров: 0.3-0.6 (средняя)
	// 50-100 товаров: 0.6-0.8 (высокая)
	// 100+ товаров: 0.8-1.0 (очень высокая)

	switch {
	case category.ProductCount >= 100:
		return 0.9
	case category.ProductCount >= 50:
		return 0.7
	case category.ProductCount >= 20:
		return 0.5
	case category.ProductCount >= 10:
		return 0.3
	default:
		return 0.1
	}
}

// generateNewCategoryProposal генерирует предложение новой категории
func (a *AICategoryAnalyzer) generateNewCategoryProposal(
	category ClientCategoryInfo,
	existingSuggestion *CategoryMappingSuggestion,
) NewCategoryProposal {
	// Извлекаем последний уровень категории как название
	levels := a.mapper.splitCategoryPath(category.Path)
	name := levels[len(levels)-1]

	// Parent category - используем suggested или дефолтную
	parentID := existingSuggestion.SuggestedCategoryID
	if parentID == 1001 {
		// Если это дефолтная категория, используем её как parent
		parentID = 1001
	}

	// Генерируем описание
	description := fmt.Sprintf("Категория импортирована из внешнего источника: %s. Содержит %d товаров.",
		category.Path,
		category.ProductCount,
	)

	// Reasoning
	reasoning := fmt.Sprintf("Обнаружена значимая категория '%s' с %d товарами, которая не имеет хорошего аналога в системе (confidence %.2f). Рекомендуется создание новой категории для лучшей организации.",
		category.Path,
		category.ProductCount,
		existingSuggestion.ConfidenceScore,
	)

	// Извлекаем возможные теги из пути категории
	tags := extractTags(category.Path)

	return NewCategoryProposal{
		ParentCategoryID: parentID,
		Name:             name,
		NameTranslations: map[string]string{
			"sr": name, // Оригинальное название (сербское)
			"en": name, // TODO: AI перевод
			"ru": name, // TODO: AI перевод
		},
		Description:       description,
		Reasoning:         reasoning,
		ExpectedProducts:  category.ProductCount,
		SimilarCategories: []int{parentID},
		Tags:              tags,
	}
}

// extractTags извлекает потенциальные теги из пути категории
func extractTags(categoryPath string) []string {
	// Простое извлечение: берем все уровни как теги в lowercase
	tags := []string{}
	parts := strings.Split(categoryPath, ">")

	for _, part := range parts {
		trimmed := strings.TrimSpace(strings.ToLower(part))
		if trimmed != "" && len(trimmed) > 2 {
			tags = append(tags, trimmed)
		}
	}

	return tags
}
