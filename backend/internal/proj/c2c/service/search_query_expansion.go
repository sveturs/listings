package service

import (
	"context"
	"strings"

	"backend/internal/logger"
)

// ExpandQueryWithSynonyms расширяет поисковый запрос синонимами
func (s *MarketplaceService) ExpandQueryWithSynonyms(ctx context.Context, query string, language string) (string, error) {
	// Используем функцию PostgreSQL для расширения запроса
	expandedQuery, err := s.storage.ExpandSearchQuery(ctx, query, language)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Str("language", language).Msg("Failed to expand query with synonyms")
		// В случае ошибки возвращаем исходный запрос
		return query, nil
	}

	logger.Info().Str("original_query", query).Str("expanded_query", expandedQuery).Str("language", language).Msg("Query expanded with synonyms")
	return expandedQuery, nil
}

// SearchCategoriesFuzzy выполняет нечеткий поиск по категориям
func (s *MarketplaceService) SearchCategoriesFuzzy(ctx context.Context, searchTerm string, language string, similarityThreshold float64) ([]CategorySearchResult, error) {
	if similarityThreshold <= 0 {
		similarityThreshold = 0.3
	}

	results, err := s.storage.SearchCategoriesFuzzy(ctx, searchTerm, language, similarityThreshold)
	if err != nil {
		logger.Error().Err(err).Str("searchTerm", searchTerm).Msg("Failed to search categories with fuzzy matching")
		return nil, err
	}

	// Преобразуем результаты из interface{} в CategorySearchResult
	var categoryResults []CategorySearchResult
	for _, result := range results {
		if resultMap, ok := result.(map[string]interface{}); ok {
			categoryResult := CategorySearchResult{
				CategoryID:      resultMap["category_id"].(int),
				CategorySlug:    resultMap["category_slug"].(string),
				CategoryName:    resultMap["category_name"].(string),
				SimilarityScore: resultMap["similarity_score"].(float64),
			}
			categoryResults = append(categoryResults, categoryResult)
		}
	}

	return categoryResults, nil
}

// CategorySearchResult представляет результат поиска категории
type CategorySearchResult struct {
	CategoryID      int     `json:"category_id"`
	CategorySlug    string  `json:"category_slug"`
	CategoryName    string  `json:"category_name"`
	SimilarityScore float64 `json:"similarity_score"`
}

// AnalyzeSearchQuery анализирует поисковый запрос и выделяет ключевые слова
func AnalyzeSearchQuery(query string) QueryAnalysis {
	query = strings.ToLower(strings.TrimSpace(query))
	words := strings.Fields(query)

	analysis := QueryAnalysis{
		OriginalQuery: query,
		Words:         words,
		WordCount:     len(words),
	}

	// Определяем тип запроса
	switch {
	case len(words) == 0:
		analysis.QueryType = "empty"
	case len(words) == 1:
		analysis.QueryType = "single_word"
	case len(words) == 2:
		analysis.QueryType = "two_words"
		// Для двухсловных запросов часто это марка+модель
		analysis.PossibleMakeModel = true
	default:
		analysis.QueryType = "multi_word"
	}

	// Проверяем на наличие чисел (может быть год, размер и т.д.)
	for _, word := range words {
		if isNumeric(word) {
			analysis.HasNumbers = true
			break
		}
	}

	return analysis
}

// QueryAnalysis представляет анализ поискового запроса
type QueryAnalysis struct {
	OriginalQuery     string   `json:"original_query"`
	Words             []string `json:"words"`
	WordCount         int      `json:"word_count"`
	QueryType         string   `json:"query_type"`
	HasNumbers        bool     `json:"has_numbers"`
	PossibleMakeModel bool     `json:"possible_make_model"`
}

// isNumeric проверяет, является ли строка числом
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
