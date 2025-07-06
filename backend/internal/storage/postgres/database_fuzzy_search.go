package postgres

import (
	"context"
	"fmt"

	"backend/internal/logger"
)

// ExpandSearchQuery расширяет поисковый запрос синонимами используя функцию PostgreSQL
func (db *Database) ExpandSearchQuery(ctx context.Context, query string, language string) (string, error) {
	var expandedQuery string

	sqlQuery := `SELECT expand_search_query($1, $2)`

	err := db.pool.QueryRow(ctx, sqlQuery, query, language).Scan(&expandedQuery)
	if err != nil {
		logger.Error().Err(err).Str("query", query).Str("language", language).Msg("Failed to expand search query")
		return query, fmt.Errorf("failed to expand search query: %w", err)
	}

	logger.Info().Str("original", query).Str("expanded", expandedQuery).Str("language", language).Msg("Query expanded with synonyms")
	return expandedQuery, nil
}

// SearchCategoriesFuzzy выполняет нечеткий поиск по категориям
func (db *Database) SearchCategoriesFuzzy(ctx context.Context, searchTerm string, language string, similarityThreshold float64) ([]interface{}, error) {
	if language == "" {
		language = "ru"
	}

	if similarityThreshold <= 0 {
		similarityThreshold = 0.3
	}

	sqlQuery := `
		SELECT 
			category_id,
			category_slug,
			category_name,
			similarity_score
		FROM search_categories_fuzzy($1, $2, $3)
		ORDER BY similarity_score DESC
		LIMIT 10
	`

	rows, err := db.pool.Query(ctx, sqlQuery, searchTerm, language, similarityThreshold)
	if err != nil {
		logger.Error().Err(err).Str("searchTerm", searchTerm).Msg("Failed to search categories with fuzzy matching")
		return nil, fmt.Errorf("failed to search categories: %w", err)
	}
	defer rows.Close()

	var results []interface{}

	for rows.Next() {
		var categoryID int
		var categorySlug, categoryName string
		var similarityScore float64

		err := rows.Scan(&categoryID, &categorySlug, &categoryName, &similarityScore)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan category search result")
			continue
		}

		result := map[string]interface{}{
			"category_id":      categoryID,
			"category_slug":    categorySlug,
			"category_name":    categoryName,
			"similarity_score": similarityScore,
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		logger.Error().Err(err).Msg("Error iterating category search results")
		return nil, fmt.Errorf("error iterating results: %w", err)
	}

	logger.Info().Int("results_count", len(results)).Str("searchTerm", searchTerm).Msg("Fuzzy category search completed")
	return results, nil
}
