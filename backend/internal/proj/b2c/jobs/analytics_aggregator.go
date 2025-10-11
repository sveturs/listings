package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/domain/models"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"
)

// AnalyticsAggregator агрегирует аналитику витрин
type AnalyticsAggregator struct {
	db     *postgres.Database
	logger *logger.Logger
}

// NewAnalyticsAggregator создает новый агрегатор аналитики
func NewAnalyticsAggregator(db *postgres.Database) *AnalyticsAggregator {
	return &AnalyticsAggregator{
		db:     db,
		logger: logger.GetLogger(),
	}
}

// unmarshalToJSONB преобразует json.RawMessage в models.JSONB
func unmarshalToJSONB(raw json.RawMessage) models.JSONB {
	var result models.JSONB
	if err := json.Unmarshal(raw, &result); err != nil {
		return models.JSONB{}
	}
	return result
}

// Run запускает агрегацию аналитики
func (a *AnalyticsAggregator) Run(ctx context.Context) error {
	a.logger.Info("Starting analytics aggregation")

	// Получаем все активные витрины
	b2c_stores, err := a.getActiveStorefronts(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active b2c_stores: %w", err)
	}

	// Агрегируем данные для каждой витрины за вчерашний день
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	for _, storefront := range b2c_stores {
		if err := a.aggregateStorefrontAnalytics(ctx, storefront.ID, yesterday); err != nil {
			a.logger.Error("Failed to aggregate analytics for storefront %d: %v", storefront.ID, err)
			continue
		}
	}

	a.logger.Info("Analytics aggregation completed")
	return nil
}

// aggregateStorefrontAnalytics агрегирует аналитику для одной витрины
func (a *AnalyticsAggregator) aggregateStorefrontAnalytics(ctx context.Context, storefrontID int, date time.Time) error {
	startOfDay := date
	endOfDay := date.Add(24 * time.Hour).Add(-time.Second)

	// Собираем статистику событий
	eventStats, err := a.getEventStats(ctx, storefrontID, startOfDay, endOfDay)
	if err != nil {
		return err
	}

	// Собираем статистику заказов
	orderStats, err := a.getOrderStats(ctx, storefrontID, startOfDay, endOfDay)
	if err != nil {
		return err
	}

	// Собираем уникальных посетителей
	uniqueVisitors, err := a.getUniqueVisitors(ctx, storefrontID, startOfDay, endOfDay)
	if err != nil {
		return err
	}

	// Собираем источники трафика
	trafficSources, err := a.getTrafficSources(ctx, storefrontID, startOfDay, endOfDay)
	if err != nil {
		return err
	}

	// Собираем топ товары
	topProducts, err := a.getTopProducts(ctx, storefrontID, startOfDay, endOfDay, 10)
	if err != nil {
		return err
	}

	// Собираем распределение заказов по городам
	ordersByCity, err := a.getOrdersByCity(ctx, storefrontID, startOfDay, endOfDay)
	if err != nil {
		return err
	}

	// Рассчитываем метрики
	bounceRate := a.calculateBounceRate(ctx, storefrontID, startOfDay, endOfDay)
	avgSessionTime := a.calculateAvgSessionTime(ctx, storefrontID, startOfDay, endOfDay)
	conversionRate := float64(0)
	if uniqueVisitors > 0 {
		conversionRate = float64(orderStats.OrdersCount) / float64(uniqueVisitors) * 100
	}

	// Создаем запись аналитики
	analytics := &models.StorefrontAnalytics{
		StorefrontID:   storefrontID,
		Date:           date,
		PageViews:      eventStats["page_view"],
		UniqueVisitors: uniqueVisitors,
		BounceRate:     bounceRate,
		AvgSessionTime: avgSessionTime,
		OrdersCount:    orderStats.OrdersCount,
		Revenue:        orderStats.Revenue,
		AvgOrderValue:  orderStats.AvgOrderValue,
		ConversionRate: conversionRate,
		ProductViews:   eventStats["product_view"],
		AddToCartCount: eventStats["add_to_cart"],
		CheckoutCount:  eventStats["checkout"],
		TrafficSources: unmarshalToJSONB(trafficSources),
		TopProducts:    unmarshalToJSONB(topProducts),
		OrdersByCity:   unmarshalToJSONB(ordersByCity),
	}

	// Сохраняем в базу
	return a.saveAnalytics(ctx, analytics)
}

// getActiveStorefronts получает активные витрины
func (a *AnalyticsAggregator) getActiveStorefronts(ctx context.Context) ([]*models.Storefront, error) {
	query := `
		SELECT id, name FROM b2c_stores 
		WHERE is_active = true AND deleted_at IS NULL
	`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логируем ошибку закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var b2c_stores []*models.Storefront
	for rows.Next() {
		s := &models.Storefront{}
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		b2c_stores = append(b2c_stores, s)
	}

	return b2c_stores, nil
}

// getEventStats получает статистику событий
func (a *AnalyticsAggregator) getEventStats(ctx context.Context, storefrontID int, from, to time.Time) (map[string]int, error) {
	query := `
		SELECT event_type, COUNT(*) 
		FROM b2c_events
		WHERE storefront_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY event_type
	`

	rows, err := a.db.QueryContext(ctx, query, storefrontID, from, to)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логируем ошибку закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	stats := make(map[string]int)
	for rows.Next() {
		var eventType string
		var count int
		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, err
		}
		stats[eventType] = count
	}

	return stats, nil
}

// OrderStats статистика заказов
type OrderStats struct {
	OrdersCount   int
	Revenue       float64
	AvgOrderValue float64
}

// getOrderStats получает статистику заказов
func (a *AnalyticsAggregator) getOrderStats(ctx context.Context, storefrontID int, from, to time.Time) (*OrderStats, error) {
	query := `
		SELECT 
			COUNT(*) as orders_count,
			COALESCE(SUM(total_amount), 0) as revenue,
			COALESCE(AVG(total_amount), 0) as avg_order_value
		FROM c2c_orders
		WHERE storefront_id = $1 
			AND created_at >= $2 
			AND created_at <= $3
			AND status NOT IN ('cancelled', 'failed')
	`

	stats := &OrderStats{}
	err := a.db.QueryRowContext(ctx, query, storefrontID, from, to).Scan(
		&stats.OrdersCount,
		&stats.Revenue,
		&stats.AvgOrderValue,
	)

	return stats, err
}

// getUniqueVisitors получает количество уникальных посетителей
func (a *AnalyticsAggregator) getUniqueVisitors(ctx context.Context, storefrontID int, from, to time.Time) (int, error) {
	query := `
		SELECT COUNT(DISTINCT session_id)
		FROM b2c_events
		WHERE storefront_id = $1 
			AND created_at >= $2 
			AND created_at <= $3
			AND event_type = 'page_view'
	`

	var count int
	err := a.db.QueryRowContext(ctx, query, storefrontID, from, to).Scan(&count)
	return count, err
}

// getTrafficSources получает источники трафика
func (a *AnalyticsAggregator) getTrafficSources(ctx context.Context, storefrontID int, from, to time.Time) (json.RawMessage, error) {
	query := `
		SELECT 
			json_build_object(
				'direct', COUNT(*) FILTER (WHERE referrer = '' OR referrer IS NULL),
				'google', COUNT(*) FILTER (WHERE referrer LIKE '%google%'),
				'facebook', COUNT(*) FILTER (WHERE referrer LIKE '%facebook%'),
				'instagram', COUNT(*) FILTER (WHERE referrer LIKE '%instagram%'),
				'other', COUNT(*) FILTER (WHERE referrer NOT LIKE '%google%' 
					AND referrer NOT LIKE '%facebook%' 
					AND referrer NOT LIKE '%instagram%' 
					AND referrer != '' 
					AND referrer IS NOT NULL)
			)
		FROM b2c_events
		WHERE storefront_id = $1 
			AND created_at >= $2 
			AND created_at <= $3
			AND event_type = 'page_view'
	`

	var result json.RawMessage
	err := a.db.QueryRowContext(ctx, query, storefrontID, from, to).Scan(&result)
	return result, err
}

// getTopProducts получает топ товаров
func (a *AnalyticsAggregator) getTopProducts(ctx context.Context, storefrontID int, from, to time.Time, limit int) (json.RawMessage, error) {
	query := `
		SELECT json_agg(product_data)
		FROM (
			SELECT json_build_object(
				'product_id', ml.id,
				'title', ml.title,
				'views', COUNT(*),
				'revenue', COALESCE(SUM(oi.price * oi.quantity), 0)
			) as product_data
			FROM b2c_events se
			JOIN c2c_listings ml ON ml.id = (se.event_data->>'product_id')::int
			LEFT JOIN c2c_order_items oi ON oi.listing_id = ml.id
			WHERE se.storefront_id = $1 
				AND se.created_at >= $2 
				AND se.created_at <= $3
				AND se.event_type = 'product_view'
			GROUP BY ml.id, ml.title
			ORDER BY COUNT(*) DESC
			LIMIT $4
		) t
	`

	var result json.RawMessage
	err := a.db.QueryRowContext(ctx, query, storefrontID, from, to, limit).Scan(&result)
	if err != nil {
		// Если нет данных, логируем и возвращаем пустой массив
		a.logger.Info("No top products data found, returning empty array: %v", err)
		return json.RawMessage("[]"), nil
	}
	return result, err
}

// getOrdersByCity получает распределение заказов по городам
func (a *AnalyticsAggregator) getOrdersByCity(ctx context.Context, storefrontID int, from, to time.Time) (json.RawMessage, error) {
	query := `
		SELECT json_object_agg(city, order_count)
		FROM (
			SELECT 
				COALESCE(delivery_city, 'Unknown') as city,
				COUNT(*) as order_count
			FROM c2c_orders
			WHERE storefront_id = $1 
				AND created_at >= $2 
				AND created_at <= $3
				AND status NOT IN ('cancelled', 'failed')
			GROUP BY delivery_city
		) t
	`

	var result json.RawMessage
	err := a.db.QueryRowContext(ctx, query, storefrontID, from, to).Scan(&result)
	if err != nil {
		// Если нет данных, логируем и возвращаем пустой объект
		a.logger.Info("No orders by city data found, returning empty object: %v", err)
		return json.RawMessage("{}"), nil
	}
	return result, err
}

// calculateBounceRate рассчитывает показатель отказов
func (a *AnalyticsAggregator) calculateBounceRate(ctx context.Context, storefrontID int, from, to time.Time) float64 {
	query := `
		WITH session_events AS (
			SELECT 
				session_id,
				COUNT(*) as event_count
			FROM b2c_events
			WHERE storefront_id = $1 
				AND created_at >= $2 
				AND created_at <= $3
			GROUP BY session_id
		)
		SELECT 
			CASE 
				WHEN COUNT(*) = 0 THEN 0
				ELSE (COUNT(*) FILTER (WHERE event_count = 1))::float / COUNT(*)::float * 100
			END as bounce_rate
		FROM session_events
	`

	var bounceRate float64
	_ = a.db.QueryRowContext(ctx, query, storefrontID, from, to).Scan(&bounceRate)
	return bounceRate
}

// calculateAvgSessionTime рассчитывает среднее время сессии
func (a *AnalyticsAggregator) calculateAvgSessionTime(ctx context.Context, storefrontID int, from, to time.Time) int {
	query := `
		WITH session_times AS (
			SELECT 
				session_id,
				EXTRACT(EPOCH FROM (MAX(created_at) - MIN(created_at))) as duration
			FROM b2c_events
			WHERE storefront_id = $1 
				AND created_at >= $2 
				AND created_at <= $3
			GROUP BY session_id
			HAVING COUNT(*) > 1
		)
		SELECT COALESCE(AVG(duration)::int, 0)
		FROM session_times
	`

	var avgTime int
	_ = a.db.QueryRowContext(ctx, query, storefrontID, from, to).Scan(&avgTime)
	return avgTime
}

// saveAnalytics сохраняет аналитику в базу
func (a *AnalyticsAggregator) saveAnalytics(ctx context.Context, analytics *models.StorefrontAnalytics) error {
	query := `
		INSERT INTO b2c_analytics (
			storefront_id, date, page_views, unique_visitors, bounce_rate, avg_session_time,
			orders_count, revenue, avg_order_value, conversion_rate,
			payment_methods_usage, product_views, add_to_cart_count, checkout_count,
			traffic_sources, top_products, top_categories, orders_by_city
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
		ON CONFLICT (storefront_id, date) DO UPDATE SET
			page_views = EXCLUDED.page_views,
			unique_visitors = EXCLUDED.unique_visitors,
			bounce_rate = EXCLUDED.bounce_rate,
			avg_session_time = EXCLUDED.avg_session_time,
			orders_count = EXCLUDED.orders_count,
			revenue = EXCLUDED.revenue,
			avg_order_value = EXCLUDED.avg_order_value,
			conversion_rate = EXCLUDED.conversion_rate,
			product_views = EXCLUDED.product_views,
			add_to_cart_count = EXCLUDED.add_to_cart_count,
			checkout_count = EXCLUDED.checkout_count,
			traffic_sources = EXCLUDED.traffic_sources,
			top_products = EXCLUDED.top_products,
			orders_by_city = EXCLUDED.orders_by_city,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := a.db.ExecContext(ctx, query,
		analytics.StorefrontID, analytics.Date, analytics.PageViews, analytics.UniqueVisitors,
		analytics.BounceRate, analytics.AvgSessionTime, analytics.OrdersCount, analytics.Revenue,
		analytics.AvgOrderValue, analytics.ConversionRate, analytics.PaymentMethodsUsage,
		analytics.ProductViews, analytics.AddToCartCount, analytics.CheckoutCount,
		analytics.TrafficSources, analytics.TopProducts, analytics.TopCategories,
		analytics.OrdersByCity,
	)

	return err
}
