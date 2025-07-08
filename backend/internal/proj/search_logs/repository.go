package searchlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/domain"

	"github.com/jmoiron/sqlx"
)

// PostgresRepository представляет репозиторий для работы с логами поиска в PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository создает новый экземпляр репозитория
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateSearchLog создает новую запись лога поискового запроса
func (r *PostgresRepository) CreateSearchLog(ctx context.Context, input *domain.SearchLogInput) (*domain.SearchLog, error) {
	query := `
		INSERT INTO search_logs (
			user_id, session_id, query_text, filters, category_id, 
			location, results_count, response_time_ms, page, per_page, 
			sort_by, user_agent, ip_address, referer
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id, created_at`

	var log domain.SearchLog
	var filtersJSON, locationJSON []byte
	var err error

	if input.Filters != nil {
		filtersJSON, err = json.Marshal(input.Filters)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal filters: %w", err)
		}
	} else {
		filtersJSON = nil
	}

	if input.Location != nil {
		locationJSON, err = json.Marshal(input.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal location: %w", err)
		}
	} else {
		locationJSON = nil
	}

	var ipAddress, userAgent, referer *string
	if input.IPAddress != "" {
		ipAddress = &input.IPAddress
	}
	if input.UserAgent != "" {
		userAgent = &input.UserAgent
	}
	if input.Referer != "" {
		referer = &input.Referer
	}

	// Используем интерфейс{} для правильной обработки nil значений в JSONB
	var filtersParam, locationParam interface{}
	if filtersJSON != nil {
		filtersParam = filtersJSON
	}
	if locationJSON != nil {
		locationParam = locationJSON
	}

	err = r.db.QueryRowContext(ctx, query,
		input.UserID, input.SessionID, input.QueryText, filtersParam,
		input.CategoryID, locationParam, input.ResultsCount,
		input.ResponseTimeMs, input.Page, input.PerPage,
		input.SortBy, userAgent, ipAddress, referer,
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create search log: %w", err)
	}

	// Заполняем остальные поля
	log.UserID = input.UserID
	log.SessionID = input.SessionID
	log.QueryText = input.QueryText
	log.Filters = filtersJSON
	log.CategoryID = input.CategoryID
	log.Location = locationJSON
	log.ResultsCount = input.ResultsCount
	log.ResponseTimeMs = input.ResponseTimeMs
	log.Page = input.Page
	log.PerPage = input.PerPage
	log.SortBy = input.SortBy
	log.UserAgent = userAgent
	log.IPAddress = ipAddress
	log.Referer = referer

	return &log, nil
}

// CreateSearchLogBatch создает несколько записей логов поисковых запросов
func (r *PostgresRepository) CreateSearchLogBatch(ctx context.Context, logs []*domain.SearchLogInput) error {
	if len(logs) == 0 {
		return nil
	}

	// Подготавливаем данные для batch insert
	valueStrings := make([]string, 0, len(logs))
	valueArgs := make([]interface{}, 0, len(logs)*14)

	for i, log := range logs {
		var filtersJSON, locationJSON []byte
		var err error

		if log.Filters != nil {
			filtersJSON, err = json.Marshal(log.Filters)
			if err != nil {
				return fmt.Errorf("failed to marshal filters for log %d: %w", i, err)
			}
		} else {
			filtersJSON = nil
		}

		if log.Location != nil {
			locationJSON, err = json.Marshal(log.Location)
			if err != nil {
				return fmt.Errorf("failed to marshal location for log %d: %w", i, err)
			}
		} else {
			locationJSON = nil
		}

		var ipAddress, userAgent, referer *string
		if log.IPAddress != "" {
			ipAddress = &log.IPAddress
		}
		if log.UserAgent != "" {
			userAgent = &log.UserAgent
		}
		if log.Referer != "" {
			referer = &log.Referer
		}

		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*14+1, i*14+2, i*14+3, i*14+4, i*14+5, i*14+6, i*14+7,
			i*14+8, i*14+9, i*14+10, i*14+11, i*14+12, i*14+13, i*14+14,
		))

		// Используем интерфейс{} для правильной обработки nil значений в JSONB
		var filtersParam, locationParam interface{}
		if filtersJSON != nil {
			filtersParam = filtersJSON
		}
		if locationJSON != nil {
			locationParam = locationJSON
		}

		valueArgs = append(valueArgs,
			log.UserID, log.SessionID, log.QueryText, filtersParam,
			log.CategoryID, locationParam, log.ResultsCount,
			log.ResponseTimeMs, log.Page, log.PerPage,
			log.SortBy, userAgent, ipAddress, referer,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO search_logs (
			user_id, session_id, query_text, filters, category_id, 
			location, results_count, response_time_ms, page, per_page, 
			sort_by, user_agent, ip_address, referer
		) VALUES %s`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to create search log batch: %w", err)
	}

	return nil
}

// CreateClickLog создает запись о клике по результату поиска
func (r *PostgresRepository) CreateClickLog(ctx context.Context, click *domain.SearchResultClick) error {
	query := `
		INSERT INTO search_result_clicks (search_log_id, listing_id, position, clicked_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query,
		click.SearchLogID, click.ListingID, click.Position, click.ClickedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create click log: %w", err)
	}

	return nil
}

// CreateClickLogBatch создает несколько записей о кликах
func (r *PostgresRepository) CreateClickLogBatch(ctx context.Context, clicks []*domain.SearchResultClick) error {
	if len(clicks) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(clicks))
	valueArgs := make([]interface{}, 0, len(clicks)*4)

	for i, click := range clicks {
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4,
		))

		valueArgs = append(valueArgs,
			click.SearchLogID, click.ListingID, click.Position, click.ClickedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO search_result_clicks (search_log_id, listing_id, position, clicked_at)
		VALUES %s`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to create click log batch: %w", err)
	}

	return nil
}

// GetTrendingQueries возвращает популярные поисковые запросы
func (r *PostgresRepository) GetTrendingQueries(ctx context.Context, limit int, categoryID *int, country *string) ([]*domain.SearchTrendingQuery, error) {
	query := `
		SELECT id, query_text, category_id, location_country, trend_score,
			search_count_24h, search_count_7d, search_count_30d,
			first_seen_at, last_seen_at, updated_at
		FROM search_trending_queries
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	if categoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *categoryID)
	}

	if country != nil {
		argCount++
		query += fmt.Sprintf(" AND location_country = $%d", argCount)
		args = append(args, *country)
	}

	argCount++
	query += fmt.Sprintf(" ORDER BY trend_score DESC LIMIT $%d", argCount)
	args = append(args, limit)

	var queries []*domain.SearchTrendingQuery
	err := r.db.SelectContext(ctx, &queries, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending queries: %w", err)
	}

	return queries, nil
}

// UpdateAnalytics обновляет агрегированную аналитику
func (r *PostgresRepository) UpdateAnalytics(ctx context.Context) error {
	// Агрегируем данные за последний час
	query := `
		INSERT INTO search_analytics (
			date, hour, query_text, category_id, location_country, 
			location_region, location_city, search_count, unique_users_count, 
			unique_sessions_count, avg_results_count, avg_response_time_ms, 
			zero_results_count
		)
		SELECT 
			DATE(created_at) as date,
			EXTRACT(HOUR FROM created_at)::INTEGER as hour,
			query_text,
			category_id,
			location->>'country' as location_country,
			location->>'region' as location_region,
			location->>'city' as location_city,
			COUNT(*) as search_count,
			COUNT(DISTINCT user_id) as unique_users_count,
			COUNT(DISTINCT session_id) as unique_sessions_count,
			AVG(results_count) as avg_results_count,
			AVG(response_time_ms) as avg_response_time_ms,
			COUNT(*) FILTER (WHERE results_count = 0) as zero_results_count
		FROM search_logs
		WHERE created_at >= NOW() - INTERVAL '1 hour'
		GROUP BY date, hour, query_text, category_id, location_country, location_region, location_city
		ON CONFLICT (date, hour, query_text, category_id, location_country, location_region, location_city)
		DO UPDATE SET
			search_count = search_analytics.search_count + EXCLUDED.search_count,
			unique_users_count = EXCLUDED.unique_users_count,
			unique_sessions_count = EXCLUDED.unique_sessions_count,
			avg_results_count = EXCLUDED.avg_results_count,
			avg_response_time_ms = EXCLUDED.avg_response_time_ms,
			zero_results_count = search_analytics.zero_results_count + EXCLUDED.zero_results_count,
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update analytics: %w", err)
	}

	// Обновляем click-through rate
	updateCTRQuery := `
		UPDATE search_analytics sa
		SET click_through_rate = (
			SELECT COUNT(DISTINCT src.search_log_id)::NUMERIC * 100.0 / sa.search_count
			FROM search_logs sl
			INNER JOIN search_result_clicks src ON sl.id = src.search_log_id
			WHERE sl.query_text = sa.query_text
				AND DATE(sl.created_at) = sa.date
				AND EXTRACT(HOUR FROM sl.created_at) = sa.hour
				AND (sl.category_id = sa.category_id OR (sl.category_id IS NULL AND sa.category_id IS NULL))
				AND (sl.location->>'country' = sa.location_country OR (sl.location->>'country' IS NULL AND sa.location_country IS NULL))
		)
		WHERE sa.updated_at >= NOW() - INTERVAL '1 hour'`

	_, err = r.db.ExecContext(ctx, updateCTRQuery)
	if err != nil {
		return fmt.Errorf("failed to update CTR: %w", err)
	}

	return nil
}

// UpdateTrendingQueries обновляет популярные запросы
func (r *PostgresRepository) UpdateTrendingQueries(ctx context.Context) error {
	// Вычисляем трендовые запросы на основе частоты и новизны
	query := `
		INSERT INTO search_trending_queries (
			query_text, category_id, location_country,
			search_count_24h, search_count_7d, search_count_30d,
			trend_score, first_seen_at, last_seen_at
		)
		SELECT 
			query_text,
			category_id,
			location->>'country' as location_country,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as search_count_24h,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as search_count_7d,
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as search_count_30d,
			-- Формула тренда: вес последних 24 часов больше
			(COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') * 10.0 +
			 COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') * 2.0 +
			 COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') * 1.0) / 
			GREATEST(EXTRACT(EPOCH FROM (NOW() - MIN(created_at))) / 86400.0, 1) as trend_score,
			MIN(created_at) as first_seen_at,
			MAX(created_at) as last_seen_at
		FROM search_logs
		WHERE created_at >= NOW() - INTERVAL '30 days'
			AND query_text != ''
		GROUP BY query_text, category_id, location_country
		HAVING COUNT(*) >= 5 -- Минимум 5 поисков для попадания в тренды
		ON CONFLICT (query_text, category_id, location_country)
		DO UPDATE SET
			search_count_24h = EXCLUDED.search_count_24h,
			search_count_7d = EXCLUDED.search_count_7d,
			search_count_30d = EXCLUDED.search_count_30d,
			trend_score = EXCLUDED.trend_score,
			last_seen_at = EXCLUDED.last_seen_at,
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to update trending queries: %w", err)
	}

	// Удаляем старые тренды
	deleteQuery := `
		DELETE FROM search_trending_queries
		WHERE last_seen_at < NOW() - INTERVAL '7 days'
			OR (trend_score < 1 AND last_seen_at < NOW() - INTERVAL '1 day')`

	_, err = r.db.ExecContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to delete old trends: %w", err)
	}

	return nil
}
