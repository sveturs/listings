package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/behavior"
	"backend/internal/logger"
)

// behaviorTrackingRepository реализация репозитория для поведенческих событий
type behaviorTrackingRepository struct {
	pool *pgxpool.Pool
}

// NewBehaviorTrackingRepository создает новый репозиторий
func NewBehaviorTrackingRepository(pool *pgxpool.Pool) BehaviorTrackingRepository {
	return &behaviorTrackingRepository{
		pool: pool,
	}
}

// SaveEvent сохраняет одно событие
func (r *behaviorTrackingRepository) SaveEvent(ctx context.Context, event *behavior.BehaviorEvent) error {
	query := `
		INSERT INTO user_behavior_events (
			event_type, user_id, session_id, search_query, 
			item_id, item_type, position, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	// Конвертируем metadata в JSON
	var metadataJSON []byte
	var err error
	if event.Metadata != nil {
		metadataJSON, err = json.Marshal(event.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	} else {
		metadataJSON = []byte("{}")
	}

	// Выполняем запрос
	err = r.pool.QueryRow(
		ctx, query,
		event.EventType, event.UserID, event.SessionID, event.SearchQuery,
		event.ItemID, event.ItemType, event.Position, metadataJSON,
	).Scan(&event.ID, &event.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}

// SaveEventsBatch сохраняет пакет событий
func (r *behaviorTrackingRepository) SaveEventsBatch(ctx context.Context, events []*behavior.BehaviorEvent) error {
	if len(events) == 0 {
		return nil
	}

	logger.Info().
		Int("batch_size", len(events)).
		Msg("SaveEventsBatch: starting to save behavior events")

	// Используем отдельное соединение для batch операций с retry логикой
	var conn *pgxpool.Conn
	var err error
	
	// Пробуем получить соединение с несколькими попытками
	for i := 0; i < 3; i++ {
		conn, err = r.pool.Acquire(ctx)
		if err == nil {
			break
		}
		if i < 2 {
			// Ждем немного перед следующей попыткой
			time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
		}
	}
	
	if err != nil {
		return fmt.Errorf("failed to acquire connection after retries: %w", err)
	}
	defer conn.Release()

	// Начинаем транзакцию
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Флаг для отслеживания успешности операции
	committed := false
	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
				logger.Error().Err(rollbackErr).Msg("Failed to rollback transaction")
			}
		}
	}()

	// Используем простой INSERT без batch для устранения проблемы conn busy
	query := `
		INSERT INTO user_behavior_events (
			event_type, user_id, session_id, search_query, 
			item_id, item_type, position, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// Сохраняем события по одному в рамках транзакции
	validEvents := 0
	for _, event := range events {
		var metadataJSON []byte
		if event.Metadata != nil {
			metadataJSON, err = json.Marshal(event.Metadata)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to marshal metadata")
				continue
			}
		} else {
			metadataJSON = []byte("{}")
		}

		// Выполняем INSERT
		_, err = tx.Exec(ctx, query,
			event.EventType, event.UserID, event.SessionID, event.SearchQuery,
			event.ItemID, event.ItemType, event.Position, metadataJSON,
		)
		
		if err != nil {
			logger.Error().Err(err).Int("index", validEvents).Msg("Failed to insert event")
			return fmt.Errorf("failed to insert event: %w", err)
		}
		
		validEvents++
	}

	// Если нет валидных событий, откатываем транзакцию
	if validEvents == 0 {
		return fmt.Errorf("no valid events to save")
	}

	logger.Info().
		Int("valid_events", validEvents).
		Msg("SaveEventsBatch: successfully prepared events for insert")

	// Коммитим транзакцию
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	
	logger.Info().
		Int("saved_events", validEvents).
		Msg("SaveEventsBatch: successfully saved behavior events")
		
	return nil
}

// GetSearchMetrics возвращает метрики поиска
func (r *behaviorTrackingRepository) GetSearchMetrics(ctx context.Context, query *behavior.SearchMetricsQuery) ([]*behavior.SearchMetrics, int, error) {
	// Базовый запрос
	baseQuery := `
		FROM search_behavior_metrics
		WHERE 1=1
	`

	// Добавляем условия
	var conditions []string
	var args []interface{}
	argCount := 0

	if query.Query != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("search_query ILIKE $%d", argCount))
		args = append(args, "%"+query.Query+"%")
	}

	// Проверка пересечения периодов
	// Метрика пересекается с запрашиваемым периодом если:
	// - period_start метрики <= period_end запроса
	// - period_end метрики >= period_start запроса
	if !query.PeriodStart.IsZero() && !query.PeriodEnd.IsZero() {
		argCount++
		conditions = append(conditions, fmt.Sprintf("period_start <= $%d", argCount))
		args = append(args, query.PeriodEnd)

		argCount++
		conditions = append(conditions, fmt.Sprintf("period_end >= $%d", argCount))
		args = append(args, query.PeriodStart)
	} else if !query.PeriodStart.IsZero() {
		argCount++
		conditions = append(conditions, fmt.Sprintf("period_end >= $%d", argCount))
		args = append(args, query.PeriodStart)
	} else if !query.PeriodEnd.IsZero() {
		argCount++
		conditions = append(conditions, fmt.Sprintf("period_start <= $%d", argCount))
		args = append(args, query.PeriodEnd)
	}

	// Собираем WHERE условия
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " AND " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// Считаем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var total int

	logger.Info().
		Str("count_query", countQuery).
		Interface("args", args).
		Msg("GetSearchMetrics: executing count query")

	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count metrics: %w", err)
	}

	logger.Info().
		Int("total", total).
		Msg("GetSearchMetrics: count result")

	// Определяем сортировку
	orderBy := "ctr DESC" // по умолчанию
	if query.SortBy != "" {
		orderBy = query.SortBy
		if query.OrderBy != "" {
			orderBy += " " + query.OrderBy
		} else {
			orderBy += " DESC"
		}
	}

	// Основной запрос с пагинацией
	selectQuery := `
		SELECT 
			id, search_query, total_searches, total_clicks, ctr,
			avg_click_position, conversions, conversion_rate,
			period_start, period_end, created_at, updated_at
	` + baseQuery + whereClause + " ORDER BY " + orderBy

	// Добавляем лимит и оффсет
	if query.Limit > 0 {
		argCount++
		selectQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
	}
	if query.Offset > 0 {
		argCount++
		selectQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, query.Offset)
	}

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get search metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*behavior.SearchMetrics
	for rows.Next() {
		var m behavior.SearchMetrics
		err := rows.Scan(
			&m.ID, &m.SearchQuery, &m.TotalSearches, &m.TotalClicks, &m.CTR,
			&m.AvgClickPosition, &m.Conversions, &m.ConversionRate,
			&m.PeriodStart, &m.PeriodEnd, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan search metrics: %w", err)
		}
		metrics = append(metrics, &m)
	}

	return metrics, total, nil
}

// GetItemMetrics возвращает метрики товаров
func (r *behaviorTrackingRepository) GetItemMetrics(ctx context.Context, query *behavior.ItemMetricsQuery) ([]*behavior.ItemMetrics, int, error) {
	// Запрос для получения метрик товаров из событий
	baseQuery := `
		WITH item_events AS (
			SELECT 
				item_id,
				item_type,
				event_type,
				position,
				created_at
			FROM user_behavior_events
			WHERE item_id IS NOT NULL
	`

	var conditions []string
	var args []interface{}
	argCount := 0

	if query.ItemType != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("item_type = $%d", argCount))
		args = append(args, query.ItemType)
	}

	if !query.PeriodStart.IsZero() {
		argCount++
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argCount))
		args = append(args, query.PeriodStart)
	}

	if !query.PeriodEnd.IsZero() {
		argCount++
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argCount))
		args = append(args, query.PeriodEnd)
	}

	// Добавляем условия к базовому запросу
	if len(conditions) > 0 {
		baseQuery += " AND " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	baseQuery += `
		)
		SELECT 
			item_id,
			item_type,
			COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END) as views,
			COUNT(CASE WHEN event_type = 'result_clicked' THEN 1 END) as clicks,
			COUNT(CASE WHEN event_type = 'item_purchased' THEN 1 END) as purchases,
			CASE 
				WHEN COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END) > 0 
				THEN COUNT(CASE WHEN event_type = 'result_clicked' THEN 1 END)::float / 
					 COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END)
				ELSE 0 
			END as ctr,
			CASE 
				WHEN COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END) > 0 
				THEN COUNT(CASE WHEN event_type = 'item_purchased' THEN 1 END)::float / 
					 COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END)
				ELSE 0 
			END as conversion_rate,
			COALESCE(AVG(CASE WHEN event_type = 'result_clicked' THEN position END), 0) as avg_position
		FROM item_events
		GROUP BY item_id, item_type
	`

	// Определяем сортировку
	orderBy := "views DESC" // по умолчанию
	if query.SortBy != "" {
		orderBy = query.SortBy
		if query.OrderBy != "" {
			orderBy += " " + query.OrderBy
		} else {
			orderBy += " DESC"
		}
	}

	// Считаем общее количество
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") t"
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count item metrics: %w", err)
	}

	// Основной запрос с пагинацией
	selectQuery := baseQuery + " ORDER BY " + orderBy

	// Добавляем лимит и оффсет
	if query.Limit > 0 {
		argCount++
		selectQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
	}
	if query.Offset > 0 {
		argCount++
		selectQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, query.Offset)
	}

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query item metrics: %w", err)
	}
	defer rows.Close()

	var metrics []*behavior.ItemMetrics
	for rows.Next() {
		var m behavior.ItemMetrics
		err := rows.Scan(
			&m.ItemID, &m.ItemType, &m.Views, &m.Clicks,
			&m.Purchases, &m.CTR, &m.ConversionRate, &m.AvgPosition,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan item metrics: %w", err)
		}

		// Устанавливаем период из запроса
		m.PeriodStart = query.PeriodStart
		m.PeriodEnd = query.PeriodEnd

		metrics = append(metrics, &m)
	}

	return metrics, total, nil
}

// UpdateSearchMetrics обновляет агрегированные метрики поиска
func (r *behaviorTrackingRepository) UpdateSearchMetrics(ctx context.Context, periodStart, periodEnd time.Time) error {
	// Этот запрос агрегирует данные из user_behavior_events и обновляет search_behavior_metrics
	query := `
		INSERT INTO search_behavior_metrics (
			search_query, total_searches, total_clicks, ctr,
			avg_click_position, conversions, conversion_rate,
			period_start, period_end
		)
		SELECT 
			search_query,
			COUNT(DISTINCT CASE WHEN event_type = 'search_performed' THEN session_id END) as total_searches,
			COUNT(CASE WHEN event_type = 'result_clicked' THEN 1 END) as total_clicks,
			CASE 
				WHEN COUNT(DISTINCT CASE WHEN event_type = 'search_performed' THEN session_id END) > 0
				THEN COUNT(CASE WHEN event_type = 'result_clicked' THEN 1 END)::float / 
					 COUNT(DISTINCT CASE WHEN event_type = 'search_performed' THEN session_id END)
				ELSE 0
			END as ctr,
			COALESCE(AVG(CASE WHEN event_type = 'result_clicked' THEN position END), 0) as avg_click_position,
			COUNT(CASE WHEN event_type = 'item_purchased' THEN 1 END) as conversions,
			CASE 
				WHEN COUNT(DISTINCT CASE WHEN event_type = 'search_performed' THEN session_id END) > 0
				THEN COUNT(CASE WHEN event_type = 'item_purchased' THEN 1 END)::float / 
					 COUNT(DISTINCT CASE WHEN event_type = 'search_performed' THEN session_id END)
				ELSE 0
			END as conversion_rate,
			$1 as period_start,
			$2 as period_end
		FROM user_behavior_events
		WHERE search_query IS NOT NULL 
			AND search_query != ''
			AND created_at >= $1 
			AND created_at < $2
		GROUP BY search_query
		ON CONFLICT (search_query, period_start) 
		DO UPDATE SET
			total_searches = EXCLUDED.total_searches,
			total_clicks = EXCLUDED.total_clicks,
			ctr = EXCLUDED.ctr,
			avg_click_position = EXCLUDED.avg_click_position,
			conversions = EXCLUDED.conversions,
			conversion_rate = EXCLUDED.conversion_rate,
			period_end = EXCLUDED.period_end,
			updated_at = NOW()
	`

	_, err := r.pool.Exec(ctx, query, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("failed to update search metrics: %w", err)
	}

	return nil
}

// GetEventsBySession возвращает события по session_id
func (r *behaviorTrackingRepository) GetEventsBySession(ctx context.Context, sessionID string) ([]*behavior.BehaviorEvent, error) {
	query := `
		SELECT 
			id, event_type, user_id, session_id, search_query,
			item_id, item_type, position, metadata, created_at
		FROM user_behavior_events
		WHERE session_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query events by session: %w", err)
	}
	defer rows.Close()

	var events []*behavior.BehaviorEvent
	for rows.Next() {
		var event behavior.BehaviorEvent
		var metadataJSON []byte
		var itemID, searchQuery, itemType sql.NullString
		var position sql.NullInt32

		err := rows.Scan(
			&event.ID, &event.EventType, &event.UserID, &event.SessionID,
			&searchQuery, &itemID, &itemType, &position, &metadataJSON, &event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		// Обрабатываем nullable поля
		if searchQuery.Valid {
			event.SearchQuery = searchQuery.String
		}
		if itemID.Valid {
			event.ItemID = itemID.String
		}
		if itemType.Valid {
			event.ItemType = behavior.ItemType(itemType.String)
		}
		if position.Valid {
			pos := int(position.Int32)
			event.Position = &pos
		}

		// Парсим metadata
		if len(metadataJSON) > 0 {
			var metadata map[string]interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			} else {
				event.Metadata = metadata
			}
		}

		events = append(events, &event)
	}

	return events, nil
}

// GetEventsByUser возвращает события по user_id
func (r *behaviorTrackingRepository) GetEventsByUser(ctx context.Context, userID int, limit, offset int) ([]*behavior.BehaviorEvent, int, error) {
	// Считаем общее количество
	countQuery := "SELECT COUNT(*) FROM user_behavior_events WHERE user_id = $1"
	var total int
	err := r.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user events: %w", err)
	}

	// Основной запрос
	query := `
		SELECT 
			id, event_type, user_id, session_id, search_query,
			item_id, item_type, position, metadata, created_at
		FROM user_behavior_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query events by user: %w", err)
	}
	defer rows.Close()

	var events []*behavior.BehaviorEvent
	for rows.Next() {
		var event behavior.BehaviorEvent
		var metadataJSON []byte
		var itemID, searchQuery, itemType sql.NullString
		var position sql.NullInt32

		err := rows.Scan(
			&event.ID, &event.EventType, &event.UserID, &event.SessionID,
			&searchQuery, &itemID, &itemType, &position, &metadataJSON, &event.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}

		// Обрабатываем nullable поля
		if searchQuery.Valid {
			event.SearchQuery = searchQuery.String
		}
		if itemID.Valid {
			event.ItemID = itemID.String
		}
		if itemType.Valid {
			event.ItemType = behavior.ItemType(itemType.String)
		}
		if position.Valid {
			pos := int(position.Int32)
			event.Position = &pos
		}

		// Парсим metadata
		if len(metadataJSON) > 0 {
			var metadata map[string]interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			} else {
				event.Metadata = metadata
			}
		}

		events = append(events, &event)
	}

	return events, total, nil
}

// GetAggregatedSearchMetrics возвращает агрегированные метрики поиска
func (r *behaviorTrackingRepository) GetAggregatedSearchMetrics(ctx context.Context, periodStart, periodEnd time.Time) (*behavior.AggregatedSearchMetrics, error) {
	metrics := &behavior.AggregatedSearchMetrics{
		SearchTrends: []behavior.SearchTrend{},
		ClickMetrics: behavior.AggregatedClickMetrics{},
	}

	// 1. Получаем общие метрики из user_behavior_events
	err := r.pool.QueryRow(ctx, `
		SELECT 
			COUNT(*) as total_searches,
			COUNT(DISTINCT search_query) as unique_searches,
			COALESCE(AVG((metadata->>'duration_ms')::float), 0) as avg_duration_ms
		FROM user_behavior_events
		WHERE event_type = 'search_performed'
			AND created_at BETWEEN $1 AND $2
	`, periodStart, periodEnd).Scan(
		&metrics.TotalSearches,
		&metrics.UniqueSearches,
		&metrics.AverageSearchDurationMs,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get general metrics: %w", err)
	}

	// 2. Получаем тренды по дням
	rows, err := r.pool.Query(ctx, `
		WITH daily_searches AS (
			SELECT 
				DATE(created_at) as search_date,
				COUNT(*) as searches_count
			FROM user_behavior_events
			WHERE event_type = 'search_performed'
				AND created_at BETWEEN $1 AND $2
			GROUP BY DATE(created_at)
		),
		daily_clicks AS (
			SELECT 
				DATE(created_at) as click_date,
				COUNT(*) as clicks_count
			FROM user_behavior_events
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
	`, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get search trends: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trend behavior.SearchTrend
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

	// 3. Получаем общие метрики кликов
	var totalClicks int
	var avgPosition float64
	var conversions int

	err = r.pool.QueryRow(ctx, `
		WITH click_data AS (
			SELECT 
				COUNT(*) as total_clicks,
				COALESCE(AVG(position), 0) as avg_position
			FROM user_behavior_events
			WHERE event_type = 'result_clicked'
				AND created_at BETWEEN $1 AND $2
		),
		conversion_data AS (
			SELECT COUNT(*) as conversions
			FROM user_behavior_events
			WHERE event_type = 'item_purchased'
				AND created_at BETWEEN $1 AND $2
				AND session_id IN (
					SELECT DISTINCT session_id
					FROM user_behavior_events
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
	`, periodStart, periodEnd).Scan(
		&totalClicks,
		&avgPosition,
		&conversions,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get click metrics: %w", err)
	}

	// Вычисляем CTR и conversion rate
	metrics.ClickMetrics.TotalClicks = totalClicks
	metrics.ClickMetrics.AverageClickPosition = avgPosition

	if metrics.TotalSearches > 0 {
		metrics.ClickMetrics.CTR = float64(totalClicks) / float64(metrics.TotalSearches)
		metrics.ClickMetrics.ConversionRate = float64(conversions) / float64(metrics.TotalSearches)
	}

	return metrics, nil
}

// GetTopSearchQueries возвращает топ поисковых запросов с полной статистикой
func (r *behaviorTrackingRepository) GetTopSearchQueries(ctx context.Context, periodStart, periodEnd time.Time, limit int) ([]behavior.TopSearchQuery, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		WITH search_stats AS (
			SELECT 
				search_query,
				COUNT(*) as search_count,
				COALESCE(AVG((metadata->>'results_count')::int), 0) as avg_results
			FROM user_behavior_events
			WHERE event_type = 'search_performed'
				AND search_query IS NOT NULL
				AND created_at BETWEEN $1 AND $2
			GROUP BY search_query
		),
		click_stats AS (
			SELECT 
				search_query,
				COUNT(*) as click_count,
				AVG(position) as avg_position
			FROM user_behavior_events
			WHERE event_type = 'result_clicked'
				AND created_at BETWEEN $1 AND $2
				AND search_query IS NOT NULL
			GROUP BY search_query
		)
		SELECT 
			ss.search_query,
			ss.search_count,
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
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, periodStart, periodEnd, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries: %w", err)
	}
	defer rows.Close()

	var queries []behavior.TopSearchQuery
	for rows.Next() {
		var q behavior.TopSearchQuery
		err := rows.Scan(
			&q.Query,
			&q.Count,
			&q.AvgPosition,
			&q.AvgResults,
			&q.CTR,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top query: %w", err)
		}
		queries = append(queries, q)
	}

	return queries, nil
}
