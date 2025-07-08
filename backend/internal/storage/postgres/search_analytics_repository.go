package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// SearchAnalyticsRepository работает с аналитикой поиска из таблицы search_logs
type SearchAnalyticsRepository struct {
	db *sqlx.DB
}

// NewSearchAnalyticsRepository создает новый репозиторий аналитики
func NewSearchAnalyticsRepository(db *sqlx.DB) *SearchAnalyticsRepository {
	return &SearchAnalyticsRepository{
		db: db,
	}
}

// GetSearchAnalytics возвращает аналитику поиска за указанный период
func (r *SearchAnalyticsRepository) GetSearchAnalytics(ctx context.Context, timeRange string) (map[string]interface{}, error) {
	// Определяем временной интервал
	var interval string
	switch timeRange {
	case "24h":
		interval = "1 day"
	case "7d":
		interval = "7 days"
	case "30d":
		interval = "30 days"
	case "90d":
		interval = "90 days"
	default:
		interval = "7 days"
	}

	analytics := make(map[string]interface{})
	analytics["timeRange"] = timeRange

	// Получаем общую статистику
	var stats struct {
		TotalSearches      int64   `db:"total_searches"`
		UniqueQueries      int64   `db:"unique_queries"`
		UniqueUsers        int64   `db:"unique_users"`
		AvgResponseTime    float64 `db:"avg_response_time"`
		ZeroResultSearches int64   `db:"zero_result_searches"`
	}

	statsQuery := `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT query_text) as unique_queries,
			COUNT(DISTINCT user_id) as unique_users,
			AVG(response_time_ms) as avg_response_time,
			COUNT(CASE WHEN results_count = 0 THEN 1 END) as zero_result_searches
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
	`

	err := r.db.GetContext(ctx, &stats, fmt.Sprintf(statsQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get search stats: %w", err)
	}

	analytics["totalSearches"] = stats.TotalSearches
	analytics["uniqueQueries"] = stats.UniqueQueries
	analytics["uniqueUsers"] = stats.UniqueUsers
	analytics["avgResponseTime"] = stats.AvgResponseTime
	analytics["zeroResultSearches"] = stats.ZeroResultSearches

	// Получаем популярные запросы
	popularQuery := `
		SELECT 
			query_text,
			COUNT(*) as search_count,
			AVG(results_count) as avg_results,
			MAX(created_at) as last_searched
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		AND query_text != ''
		GROUP BY query_text
		ORDER BY search_count DESC
		LIMIT 20
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(popularQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get popular searches: %w", err)
	}
	defer rows.Close()

	var popularSearches []map[string]interface{}
	for rows.Next() {
		var query string
		var count int
		var avgResults float64
		var lastSearched time.Time

		if err := rows.Scan(&query, &count, &avgResults, &lastSearched); err != nil {
			continue
		}

		popularSearches = append(popularSearches, map[string]interface{}{
			"query":        query,
			"count":        count,
			"avgResults":   avgResults,
			"lastSearched": lastSearched.Format(time.RFC3339),
		})
	}
	analytics["popularSearches"] = popularSearches

	// Получаем последние поиски
	recentQuery := `
		SELECT 
			sl.query_text,
			sl.results_count,
			sl.response_time_ms,
			sl.device_type,
			sl.created_at,
			u.email as user_email
		FROM search_logs sl
		LEFT JOIN users u ON u.id = sl.user_id
		WHERE sl.created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		ORDER BY sl.created_at DESC
		LIMIT 100
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(recentQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get recent searches: %w", err)
	}
	defer rows.Close()

	var recentSearches []map[string]interface{}
	for rows.Next() {
		var query string
		var resultsCount int
		var responseTime int64
		var deviceType string
		var createdAt time.Time
		var userEmail sql.NullString

		if err := rows.Scan(&query, &resultsCount, &responseTime, &deviceType, &createdAt, &userEmail); err != nil {
			continue
		}

		search := map[string]interface{}{
			"query":        query,
			"resultsCount": resultsCount,
			"responseTime": responseTime,
			"deviceType":   deviceType,
			"createdAt":    createdAt.Format(time.RFC3339),
		}

		if userEmail.Valid {
			search["userEmail"] = userEmail.String
		}

		recentSearches = append(recentSearches, search)
	}
	analytics["recentSearches"] = recentSearches

	// Получаем статистику по устройствам
	deviceQuery := `
		SELECT 
			device_type,
			COUNT(*) as count
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		GROUP BY device_type
		ORDER BY count DESC
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(deviceQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}
	defer rows.Close()

	deviceStats := make(map[string]int)
	for rows.Next() {
		var deviceType string
		var count int
		if err := rows.Scan(&deviceType, &count); err != nil {
			continue
		}
		deviceStats[deviceType] = count
	}
	analytics["deviceStats"] = deviceStats

	// Получаем метрики для админки
	var metrics struct {
		TotalSearches    int64   `db:"total_searches"`
		UniqueQueries    int64   `db:"unique_queries"`
		AvgSearchTime    float64 `db:"avg_search_time"`
		ZeroResultsRate  float64 `db:"zero_results_rate"`
		ClickThroughRate float64 `db:"click_through_rate"`
	}

	metricsQuery := `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT query_text) as unique_queries,
			AVG(response_time_ms) as avg_search_time,
			CASE 
				WHEN COUNT(*) > 0 THEN 
					(COUNT(CASE WHEN results_count = 0 THEN 1 END)::float / COUNT(*)::float) * 100
				ELSE 0 
			END as zero_results_rate,
			0 as click_through_rate -- TODO: вычислять на основе search_result_clicks
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
	`

	err = r.db.GetContext(ctx, &metrics, fmt.Sprintf(metricsQuery, interval))
	if err != nil {
		// Если ошибка, устанавливаем значения по умолчанию
		metrics.TotalSearches = stats.TotalSearches
		metrics.UniqueQueries = stats.UniqueQueries
		metrics.AvgSearchTime = stats.AvgResponseTime
		if stats.TotalSearches > 0 {
			metrics.ZeroResultsRate = float64(stats.ZeroResultSearches) / float64(stats.TotalSearches) * 100
		}
	}

	analytics["metrics"] = map[string]interface{}{
		"totalSearches":    metrics.TotalSearches,
		"uniqueQueries":    metrics.UniqueQueries,
		"avgSearchTime":    metrics.AvgSearchTime,
		"zeroResultsRate":  metrics.ZeroResultsRate,
		"clickThroughRate": metrics.ClickThroughRate,
	}

	// Получаем запросы без результатов
	zeroResultsQuery := `
		SELECT 
			query_text,
			COUNT(*) as count
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		AND results_count = 0
		AND query_text != ''
		GROUP BY query_text
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(zeroResultsQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get zero result queries: %w", err)
	}
	defer rows.Close()

	var zeroResultQueries []map[string]interface{}
	for rows.Next() {
		var query string
		var count int
		if err := rows.Scan(&query, &count); err != nil {
			continue
		}
		zeroResultQueries = append(zeroResultQueries, map[string]interface{}{
			"query":            query,
			"count":            count,
			"avgResultsCount":  0,
			"avgClickPosition": 0,
			"lastSearched":     time.Now().Format(time.RFC3339),
		})
	}
	analytics["zeroResultQueries"] = zeroResultQueries

	// Получаем топ запросы для совместимости с фронтендом
	var topQueries []map[string]interface{}
	for i, ps := range popularSearches {
		if i >= 10 {
			break
		}
		topQueries = append(topQueries, map[string]interface{}{
			"query":            ps["query"],
			"count":            ps["count"],
			"avgResultsCount":  ps["avgResults"],
			"avgClickPosition": 1.0, // TODO: вычислять на основе search_result_clicks
			"lastSearched":     ps["lastSearched"],
		})
	}
	analytics["topQueries"] = topQueries

	return analytics, nil
}

// GetPopularSearchesFromLogs возвращает популярные поисковые запросы из search_logs
func (r *SearchAnalyticsRepository) GetPopularSearchesFromLogs(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			query_text as query,
			COUNT(*) as count,
			AVG(results_count) as avg_results
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '7 days'
		AND query_text != ''
		GROUP BY query_text
		ORDER BY count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular searches: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var query string
		var count int
		var avgResults float64

		if err := rows.Scan(&query, &count, &avgResults); err != nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"query":       query,
			"count":       count,
			"avg_results": avgResults,
		})
	}

	return results, nil
}

// GetSearchStatisticsFromLogs возвращает статистику из search_logs
func (r *SearchAnalyticsRepository) GetSearchStatisticsFromLogs(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			sl.id,
			sl.query_text as query,
			sl.results_count,
			sl.response_time_ms as search_duration_ms,
			sl.user_id,
			sl.filters as search_filters,
			sl.created_at,
			u.email as user_email
		FROM search_logs sl
		LEFT JOIN users u ON u.id = sl.user_id
		ORDER BY sl.created_at DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get search statistics: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var query string
		var resultsCount int
		var searchDuration int64
		var userID sql.NullInt64
		var filtersJSON []byte
		var createdAt time.Time
		var userEmail sql.NullString

		if err := rows.Scan(&id, &query, &resultsCount, &searchDuration, &userID, &filtersJSON, &createdAt, &userEmail); err != nil {
			continue
		}

		stat := map[string]interface{}{
			"id":                 id,
			"query":              query,
			"results_count":      resultsCount,
			"search_duration_ms": searchDuration,
			"created_at":         createdAt.Format(time.RFC3339),
		}

		if userID.Valid {
			stat["user_id"] = userID.Int64
		}

		if userEmail.Valid {
			stat["user_email"] = userEmail.String
		}

		if len(filtersJSON) > 0 {
			var filters map[string]interface{}
			if err := json.Unmarshal(filtersJSON, &filters); err == nil {
				stat["search_filters"] = filters
			}
		}

		results = append(results, stat)
	}

	return results, nil
}

// GetSearchAnalyticsWithPagination возвращает аналитику поиска с поддержкой пагинации
func (r *SearchAnalyticsRepository) GetSearchAnalyticsWithPagination(ctx context.Context, timeRange string, offsetTop, offsetZero, limit int) (map[string]interface{}, error) {
	// Определяем временной интервал
	var interval string
	switch timeRange {
	case "24h":
		interval = "1 day"
	case "7d":
		interval = "7 days"
	case "30d":
		interval = "30 days"
	case "90d":
		interval = "90 days"
	default:
		interval = "7 days"
	}

	analytics := make(map[string]interface{})
	analytics["timeRange"] = timeRange

	// Получаем общую статистику (как и раньше)
	var stats struct {
		TotalSearches      int64   `db:"total_searches"`
		UniqueQueries      int64   `db:"unique_queries"`
		UniqueUsers        int64   `db:"unique_users"`
		AvgResponseTime    float64 `db:"avg_response_time"`
		ZeroResultSearches int64   `db:"zero_result_searches"`
	}

	statsQuery := `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT query_text) as unique_queries,
			COUNT(DISTINCT user_id) as unique_users,
			AVG(response_time_ms) as avg_response_time,
			COUNT(CASE WHEN results_count = 0 THEN 1 END) as zero_result_searches
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
	`

	err := r.db.GetContext(ctx, &stats, fmt.Sprintf(statsQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get search stats: %w", err)
	}

	analytics["totalSearches"] = stats.TotalSearches
	analytics["uniqueQueries"] = stats.UniqueQueries
	analytics["uniqueUsers"] = stats.UniqueUsers
	analytics["avgResponseTime"] = stats.AvgResponseTime
	analytics["zeroResultSearches"] = stats.ZeroResultSearches

	// Получаем общее количество популярных запросов
	countQuery := `
		SELECT COUNT(DISTINCT query_text) 
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		AND query_text != ''
	`
	var totalTopQueries int
	err = r.db.GetContext(ctx, &totalTopQueries, fmt.Sprintf(countQuery, interval))
	if err != nil {
		totalTopQueries = 0
	}
	analytics["totalTopQueries"] = totalTopQueries

	// Получаем популярные запросы с пагинацией
	popularQuery := `
		SELECT 
			query_text,
			COUNT(*) as search_count,
			AVG(results_count) as avg_results,
			MAX(created_at) as last_searched
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		AND query_text != ''
		GROUP BY query_text
		ORDER BY search_count DESC
		LIMIT %d OFFSET %d
	`

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(popularQuery, interval, limit, offsetTop))
	if err != nil {
		return nil, fmt.Errorf("failed to get popular searches: %w", err)
	}
	defer rows.Close()

	var topQueries []map[string]interface{}
	for rows.Next() {
		var query string
		var count int
		var avgResults float64
		var lastSearched time.Time

		if err := rows.Scan(&query, &count, &avgResults, &lastSearched); err != nil {
			continue
		}

		topQueries = append(topQueries, map[string]interface{}{
			"query":            query,
			"count":            count,
			"avgResultsCount":  avgResults,
			"avgClickPosition": 1.0, // TODO: вычислять на основе search_result_clicks
			"lastSearched":     lastSearched.Format(time.RFC3339),
		})
	}
	analytics["topQueries"] = topQueries

	// Получаем общее количество запросов без результатов
	zeroCountQuery := `
		SELECT COUNT(DISTINCT query_text)
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		AND results_count = 0
		AND query_text != ''
	`
	var totalZeroQueries int
	err = r.db.GetContext(ctx, &totalZeroQueries, fmt.Sprintf(zeroCountQuery, interval))
	if err != nil {
		totalZeroQueries = 0
	}
	analytics["totalZeroQueries"] = totalZeroQueries

	// Получаем запросы без результатов с пагинацией
	zeroResultsQuery := `
		SELECT 
			query_text,
			COUNT(*) as count,
			MAX(created_at) as last_searched
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		AND results_count = 0
		AND query_text != ''
		GROUP BY query_text
		ORDER BY count DESC
		LIMIT %d OFFSET %d
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(zeroResultsQuery, interval, limit, offsetZero))
	if err != nil {
		return nil, fmt.Errorf("failed to get zero result queries: %w", err)
	}
	defer rows.Close()

	var zeroResultQueries []map[string]interface{}
	for rows.Next() {
		var query string
		var count int
		var lastSearched time.Time
		if err := rows.Scan(&query, &count, &lastSearched); err != nil {
			continue
		}
		zeroResultQueries = append(zeroResultQueries, map[string]interface{}{
			"query":            query,
			"count":            count,
			"avgResultsCount":  0,
			"avgClickPosition": 0,
			"lastSearched":     lastSearched.Format(time.RFC3339),
		})
	}
	analytics["zeroResultQueries"] = zeroResultQueries

	// Получаем метрики (как и раньше)
	var metrics struct {
		TotalSearches    int64   `db:"total_searches"`
		UniqueQueries    int64   `db:"unique_queries"`
		AvgSearchTime    float64 `db:"avg_search_time"`
		ZeroResultsRate  float64 `db:"zero_results_rate"`
		ClickThroughRate float64 `db:"click_through_rate"`
	}

	metricsQuery := `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT query_text) as unique_queries,
			AVG(response_time_ms) as avg_search_time,
			CASE 
				WHEN COUNT(*) > 0 THEN 
					(COUNT(CASE WHEN results_count = 0 THEN 1 END)::float / COUNT(*)::float)
				ELSE 0 
			END as zero_results_rate,
			0 as click_through_rate -- TODO: вычислять на основе search_result_clicks
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
	`

	err = r.db.GetContext(ctx, &metrics, fmt.Sprintf(metricsQuery, interval))
	if err != nil {
		// Если ошибка, устанавливаем значения по умолчанию
		metrics.TotalSearches = stats.TotalSearches
		metrics.UniqueQueries = stats.UniqueQueries
		metrics.AvgSearchTime = stats.AvgResponseTime
		if stats.TotalSearches > 0 {
			metrics.ZeroResultsRate = float64(stats.ZeroResultSearches) / float64(stats.TotalSearches)
		}
	}

	analytics["metrics"] = map[string]interface{}{
		"totalSearches":    metrics.TotalSearches,
		"uniqueQueries":    metrics.UniqueQueries,
		"avgSearchTime":    metrics.AvgSearchTime,
		"zeroResultsRate":  metrics.ZeroResultsRate,
		"clickThroughRate": metrics.ClickThroughRate,
	}

	// Получаем статистику по устройствам (без пагинации)
	deviceQuery := `
		SELECT 
			device_type,
			COUNT(*) as count
		FROM search_logs
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '%s'
		GROUP BY device_type
		ORDER BY count DESC
	`

	rows, err = r.db.QueryContext(ctx, fmt.Sprintf(deviceQuery, interval))
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}
	defer rows.Close()

	deviceStats := make(map[string]int)
	for rows.Next() {
		var deviceType string
		var count int
		if err := rows.Scan(&deviceType, &count); err != nil {
			continue
		}
		deviceStats[deviceType] = count
	}
	analytics["deviceStats"] = deviceStats

	return analytics, nil
}
