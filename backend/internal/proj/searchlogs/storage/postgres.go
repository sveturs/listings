// Package storage backend/internal/proj/searchlogs/storage/postgres.go
package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/logger"
	"backend/internal/proj/searchlogs/types"
)

// PostgresStorage реализует хранилище логов поиска в PostgreSQL
type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage создает новое хранилище
func NewPostgresStorage(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		pool: pool,
	}
}

// SaveBatch сохраняет батч записей логов
func (s *PostgresStorage) SaveBatch(ctx context.Context, entries []*types.SearchLogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	// Используем COPY для эффективной вставки
	_, err := s.pool.CopyFrom(
		ctx,
		pgx.Identifier{"search_logs"},
		[]string{
			"query_text", "user_id", "session_id", "results_count", "response_time_ms",
			"filters", "category_id", "price_min", "price_max", "location",
			"language", "device_type", "user_agent", "ip_address", "search_type",
			"has_spell_correct", "clicked_items",
		},
		pgx.CopyFromSlice(len(entries), func(i int) ([]any, error) {
			e := entries[i]

			// Конвертируем фильтры в JSON
			var filtersJSON []byte
			if e.Filters != nil && len(e.Filters) > 0 {
				var err error
				filtersJSON, err = json.Marshal(e.Filters)
				if err != nil {
					return nil, err
				}
			}

			// Конвертируем clicked_items в JSON
			var clickedItemsJSON []byte
			if e.ClickedItems != nil && len(e.ClickedItems) > 0 {
				var err error
				clickedItemsJSON, err = json.Marshal(e.ClickedItems)
				if err != nil {
					return nil, err
				}
			}

			// Конвертируем location в JSON если есть
			var locationJSON []byte
			if e.Location != nil {
				var err error
				locationJSON, err = json.Marshal(e.Location)
				if err != nil {
					return nil, err
				}
			}

			return []any{
				e.Query,
				e.UserID,
				e.SessionID,
				e.ResultCount,
				e.ResponseTimeMS,
				filtersJSON,
				e.CategoryID,
				e.PriceMin,
				e.PriceMax,
				locationJSON,
				e.Language,
				e.DeviceType,
				e.UserAgent,
				e.IP,
				e.SearchType,
				e.HasSpellCorrect,
				clickedItemsJSON,
			}, nil
		}),
	)
	if err != nil {
		logger.Error().Err(err).Int("batch_size", len(entries)).Msg("Failed to save search logs batch")
	}

	return err
}

// GetSearchStats возвращает статистику поиска за период
func (s *PostgresStorage) GetSearchStats(ctx context.Context, from, to time.Time) (*types.SearchStats, error) {
	stats := &types.SearchStats{
		SearchesByHour: make(map[int]int64),
		DeviceStats:    make(map[string]int64),
	}

	// Основная статистика
	query := `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT user_id) as unique_users,
			AVG(response_time_ms) as avg_response_time,
			COUNT(CASE WHEN results_count = 0 THEN 1 END) as zero_result_searches
		FROM search_logs
		WHERE created_at BETWEEN $1 AND $2
	`

	err := s.pool.QueryRow(ctx, query, from, to).Scan(
		&stats.TotalSearches,
		&stats.UniqueUsers,
		&stats.AvgResponseTimeMS,
		&stats.ZeroResultSearches,
	)
	if err != nil {
		return nil, err
	}

	// Топ запросы
	topQueries, err := s.GetPopularSearches(ctx, from, to, 10)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get popular searches")
	} else {
		stats.TopQueries = topQueries
	}

	// Статистика по часам
	hourQuery := `
		SELECT 
			EXTRACT(HOUR FROM created_at) as hour,
			COUNT(*) as count
		FROM search_logs
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY hour
		ORDER BY hour
	`

	rows, err := s.pool.Query(ctx, hourQuery, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get hourly stats")
	} else {
		defer rows.Close()
		for rows.Next() {
			var hour int
			var count int64
			if err := rows.Scan(&hour, &count); err == nil {
				stats.SearchesByHour[hour] = count
			}
		}
	}

	// Статистика по устройствам
	deviceQuery := `
		SELECT 
			device_type,
			COUNT(*) as count
		FROM search_logs
		WHERE created_at BETWEEN $1 AND $2
		GROUP BY device_type
	`

	rows, err = s.pool.Query(ctx, deviceQuery, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get device stats")
	} else {
		defer rows.Close()
		for rows.Next() {
			var deviceType string
			var count int64
			if err := rows.Scan(&deviceType, &count); err == nil {
				stats.DeviceStats[deviceType] = count
			}
		}
	}

	// Топ категории
	categoryQuery := `
		SELECT 
			sl.category_id,
			COALESCE(c.name, sl.category_id) as category_name,
			COUNT(*) as search_count
		FROM search_logs sl
		LEFT JOIN marketplace_categories c ON c.id::text = sl.category_id
		WHERE sl.created_at BETWEEN $1 AND $2 AND sl.category_id IS NOT NULL
		GROUP BY sl.category_id, c.name
		ORDER BY search_count DESC
		LIMIT 10
	`

	rows, err = s.pool.Query(ctx, categoryQuery, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get category stats")
	} else {
		defer rows.Close()
		for rows.Next() {
			var catStat types.CategoryStats
			if err := rows.Scan(&catStat.CategoryID, &catStat.CategoryName, &catStat.SearchCount); err == nil {
				stats.TopCategories = append(stats.TopCategories, catStat)
			}
		}
	}

	return stats, nil
}

// GetPopularSearches возвращает популярные поисковые запросы
func (s *PostgresStorage) GetPopularSearches(ctx context.Context, from, to time.Time, limit int) ([]types.PopularSearch, error) {
	query := `
		SELECT 
			query_text,
			COUNT(*) as count,
			AVG(results_count) as avg_results,
			COALESCE(SUM(CASE WHEN jsonb_array_length(clicked_items::jsonb) > 0 THEN 1 ELSE 0 END)::float / COUNT(*), 0) as click_rate
		FROM search_logs
		WHERE created_at BETWEEN $1 AND $2
			AND query_text != ''
		GROUP BY query_text
		HAVING COUNT(*) > 1
		ORDER BY count DESC
		LIMIT $3
	`

	rows, err := s.pool.Query(ctx, query, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var searches []types.PopularSearch
	for rows.Next() {
		var search types.PopularSearch
		err := rows.Scan(&search.Query, &search.Count, &search.AvgResults, &search.ClickRate)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan popular search")
			continue
		}
		searches = append(searches, search)
	}

	return searches, nil
}

// GetUserSearchHistory возвращает историю поиска пользователя
func (s *PostgresStorage) GetUserSearchHistory(ctx context.Context, userID int, limit int) ([]types.SearchLogEntry, error) {
	query := `
		SELECT 
			query_text, user_id, session_id, results_count, response_time_ms,
			filters, category_id, price_min, price_max, location,
			language, device_type, user_agent, ip_address, search_type,
			has_spell_correct, clicked_items, created_at
		FROM search_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []types.SearchLogEntry
	for rows.Next() {
		var e types.SearchLogEntry
		var filtersJSON, clickedItemsJSON []byte

		var createdAt time.Time
		err := rows.Scan(
			&e.Query, &e.UserID, &e.SessionID, &e.ResultCount, &e.ResponseTimeMS,
			&filtersJSON, &e.CategoryID, &e.PriceMin, &e.PriceMax, &e.Location,
			&e.Language, &e.DeviceType, &e.UserAgent, &e.IP, &e.SearchType,
			&e.HasSpellCorrect, &clickedItemsJSON, &createdAt,
		)
		e.Timestamp = createdAt
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan search history entry")
			continue
		}

		// Парсим JSON поля
		if filtersJSON != nil {
			json.Unmarshal(filtersJSON, &e.Filters)
		}
		if clickedItemsJSON != nil {
			json.Unmarshal(clickedItemsJSON, &e.ClickedItems)
		}

		entries = append(entries, e)
	}

	return entries, nil
}
