package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type SearchMetrics struct {
	TotalSearches           int64         `json:"total_searches"`
	UniqueSearches          int64         `json:"unique_searches"`
	AverageSearchDurationMs float64       `json:"average_search_duration_ms"`
	TopQueries              []TopQuery    `json:"top_queries"`
	SearchTrends            []SearchTrend `json:"search_trends"`
	ClickMetrics            ClickMetrics  `json:"click_metrics"`
}

type TopQuery struct {
	Query       string  `json:"query"`
	Count       int     `json:"count"`
	CTR         float64 `json:"ctr"`
	AvgPosition float64 `json:"avg_position"`
	AvgResults  float64 `json:"avg_results"`
}

type SearchTrend struct {
	Date          string  `json:"date"`
	SearchesCount int     `json:"searches_count"`
	ClicksCount   int     `json:"clicks_count"`
	CTR           float64 `json:"ctr"`
}

type ClickMetrics struct {
	TotalClicks          int     `json:"total_clicks"`
	AverageClickPosition float64 `json:"average_click_position"`
	CTR                  float64 `json:"ctr"`
	ConversionRate       float64 `json:"conversion_rate"`
}

func main() {
	// Подключаемся к базе данных
	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		dbName = "postgres"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Проверяем подключение
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Определяем период агрегации
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7) // Последние 7 дней по умолчанию

	if len(os.Args) > 1 {
		days := 7
		fmt.Sscanf(os.Args[1], "%d", &days)
		startTime = endTime.AddDate(0, 0, -days)
	}

	fmt.Printf("Aggregating search metrics from %s to %s\n",
		startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))

	// Агрегируем метрики
	metrics, err := aggregateSearchMetrics(ctx, db, startTime, endTime)
	if err != nil {
		log.Fatal("Failed to aggregate metrics:", err)
	}

	// Выводим результаты в JSON
	jsonData, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal metrics:", err)
	}

	fmt.Println(string(jsonData))
}

func aggregateSearchMetrics(ctx context.Context, db *sql.DB, startTime, endTime time.Time) (*SearchMetrics, error) {
	metrics := &SearchMetrics{
		TopQueries:   []TopQuery{},
		SearchTrends: []SearchTrend{},
	}

	// 1. Получаем общие метрики
	err := db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT search_query) as unique_searches,
			AVG(search_duration_ms) as avg_duration_ms
		FROM search_logs
		WHERE created_at BETWEEN $1 AND $2
	`, startTime, endTime).Scan(
		&metrics.TotalSearches,
		&metrics.UniqueSearches,
		&metrics.AverageSearchDurationMs,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get general metrics: %w", err)
	}

	// 2. Получаем топ запросы с метриками
	rows, err := db.QueryContext(ctx, `
		WITH search_stats AS (
			SELECT 
				sl.search_query,
				COUNT(*) as search_count,
				AVG(sl.results_count) as avg_results,
				COUNT(DISTINCT sl.session_id) as unique_sessions
			FROM search_logs sl
			WHERE sl.created_at BETWEEN $1 AND $2
			GROUP BY sl.search_query
		),
		click_stats AS (
			SELECT 
				be.search_query,
				COUNT(*) as click_count,
				AVG(be.position) as avg_position
			FROM behavior_events be
			WHERE be.event_type = 'result_clicked'
				AND be.created_at BETWEEN $1 AND $2
				AND be.search_query IS NOT NULL
			GROUP BY be.search_query
		)
		SELECT 
			ss.search_query,
			ss.search_count,
			COALESCE(cs.click_count, 0) as click_count,
			COALESCE(cs.avg_position, 0) as avg_position,
			ss.avg_results,
			CASE 
				WHEN ss.search_count > 0 
				THEN COALESCE(cs.click_count, 0)::float / ss.search_count::float
				ELSE 0 
			END as ctr
		FROM search_stats ss
		LEFT JOIN click_stats cs ON ss.search_query = cs.search_query
		ORDER BY ss.search_count DESC
		LIMIT 50
	`, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var q TopQuery
		var clickCount int
		err := rows.Scan(
			&q.Query,
			&q.Count,
			&clickCount,
			&q.AvgPosition,
			&q.AvgResults,
			&q.CTR,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top query: %w", err)
		}
		metrics.TopQueries = append(metrics.TopQueries, q)
	}

	// 3. Получаем тренды по дням
	rows, err = db.QueryContext(ctx, `
		WITH daily_searches AS (
			SELECT 
				DATE(created_at) as search_date,
				COUNT(*) as searches_count
			FROM search_logs
			WHERE created_at BETWEEN $1 AND $2
			GROUP BY DATE(created_at)
		),
		daily_clicks AS (
			SELECT 
				DATE(created_at) as click_date,
				COUNT(*) as clicks_count
			FROM behavior_events
			WHERE event_type = 'result_clicked'
				AND created_at BETWEEN $1 AND $2
			GROUP BY DATE(created_at)
		)
		SELECT 
			ds.search_date,
			ds.searches_count,
			COALESCE(dc.clicks_count, 0) as clicks_count,
			CASE 
				WHEN ds.searches_count > 0 
				THEN COALESCE(dc.clicks_count, 0)::float / ds.searches_count::float
				ELSE 0 
			END as ctr
		FROM daily_searches ds
		LEFT JOIN daily_clicks dc ON ds.search_date = dc.click_date
		ORDER BY ds.search_date
	`, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get search trends: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trend SearchTrend
		var date time.Time
		err := rows.Scan(
			&date,
			&trend.SearchesCount,
			&trend.ClicksCount,
			&trend.CTR,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trend: %w", err)
		}
		trend.Date = date.Format("2006-01-02")
		metrics.SearchTrends = append(metrics.SearchTrends, trend)
	}

	// 4. Получаем общие метрики кликов
	err = db.QueryRowContext(ctx, `
		WITH click_data AS (
			SELECT 
				COUNT(*) as total_clicks,
				AVG(position) as avg_position
			FROM behavior_events
			WHERE event_type = 'result_clicked'
				AND created_at BETWEEN $1 AND $2
		),
		conversion_data AS (
			SELECT COUNT(*) as conversions
			FROM behavior_events
			WHERE event_type = 'item_purchased'
				AND created_at BETWEEN $1 AND $2
				AND session_id IN (
					SELECT DISTINCT session_id
					FROM behavior_events
					WHERE event_type = 'search_performed'
						AND created_at BETWEEN $1 AND $2
				)
		)
		SELECT 
			COALESCE(cd.total_clicks, 0),
			COALESCE(cd.avg_position, 0),
			COALESCE(conv.conversions, 0)
		FROM click_data cd
		CROSS JOIN conversion_data conv
	`, startTime, endTime).Scan(
		&metrics.ClickMetrics.TotalClicks,
		&metrics.ClickMetrics.AverageClickPosition,
		&metrics.ClickMetrics.ConversionRate,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get click metrics: %w", err)
	}

	// Вычисляем CTR и conversion rate
	if metrics.TotalSearches > 0 {
		metrics.ClickMetrics.CTR = float64(metrics.ClickMetrics.TotalClicks) / float64(metrics.TotalSearches)
		// ConversionRate временно содержит количество конверсий, преобразуем в rate
		conversions := metrics.ClickMetrics.ConversionRate
		metrics.ClickMetrics.ConversionRate = conversions / float64(metrics.TotalSearches)
	}

	return metrics, nil
}
