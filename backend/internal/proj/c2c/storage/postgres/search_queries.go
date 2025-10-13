// backend/internal/proj/c2c/storage/postgres/search_queries.go
package postgres

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/proj/c2c/service"
)

// GetPopularSearchQueries получает популярные поисковые запросы
func (s *Storage) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]service.SearchQuery, error) {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	sqlQuery := `
		SELECT
			id,
			query,
			normalized_query,
			search_count,
			to_char(last_searched, 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') as last_searched,
			language,
			results_count
		FROM search_queries
		WHERE normalized_query LIKE '%' || $1 || '%'
		ORDER BY search_count DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, sqlQuery, normalizedQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying popular searches: %w", err)
	}
	defer rows.Close()

	var queries []service.SearchQuery
	for rows.Next() {
		var q service.SearchQuery
		if err := rows.Scan(
			&q.ID,
			&q.Query,
			&q.NormalizedQuery,
			&q.SearchCount,
			&q.LastSearched,
			&q.Language,
			&q.ResultsCount,
		); err != nil {
			return nil, fmt.Errorf("error scanning search query: %w", err)
		}
		queries = append(queries, q)
	}

	return queries, nil
}

// SaveSearchQuery сохраняет или обновляет поисковый запрос
func (s *Storage) SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error {
	if normalizedQuery == "" {
		normalizedQuery = strings.ToLower(strings.TrimSpace(query))
	}

	if normalizedQuery == "" {
		return nil // Не сохраняем пустые запросы
	}

	// Используем UPSERT для обновления существующих записей
	sqlQuery := `
		INSERT INTO search_queries (
			query, normalized_query, search_count, last_searched,
			language, results_count
		) VALUES ($1, $2, 1, NOW(), $3, $4)
		ON CONFLICT (normalized_query, language)
		DO UPDATE SET
			query = EXCLUDED.query,
			search_count = search_queries.search_count + 1,
			last_searched = NOW(),
			results_count = EXCLUDED.results_count
	`

	_, err := s.pool.Exec(ctx, sqlQuery, query, normalizedQuery, language, resultsCount)
	if err != nil {
		return fmt.Errorf("error saving search query: %w", err)
	}

	return nil
}
