package postgres

import (
	"context"
	"time"

	"backend/internal/domain/behavior"
)

// BehaviorTrackingRepository определяет интерфейс для работы с поведенческими событиями
type BehaviorTrackingRepository interface {
	// SaveEvent сохраняет одно событие
	SaveEvent(ctx context.Context, event *behavior.BehaviorEvent) error

	// SaveEventsBatch сохраняет пакет событий
	SaveEventsBatch(ctx context.Context, events []*behavior.BehaviorEvent) error

	// GetSearchMetrics возвращает метрики поиска
	GetSearchMetrics(ctx context.Context, query *behavior.SearchMetricsQuery) ([]*behavior.SearchMetrics, int, error)

	// GetItemMetrics возвращает метрики товаров
	GetItemMetrics(ctx context.Context, query *behavior.ItemMetricsQuery) ([]*behavior.ItemMetrics, int, error)

	// UpdateSearchMetrics обновляет агрегированные метрики поиска
	UpdateSearchMetrics(ctx context.Context, periodStart, periodEnd time.Time) error

	// GetEventsBySession возвращает события по session_id
	GetEventsBySession(ctx context.Context, sessionID string) ([]*behavior.BehaviorEvent, error)

	// GetEventsByUser возвращает события по user_id
	GetEventsByUser(ctx context.Context, userID int, limit, offset int) ([]*behavior.BehaviorEvent, int, error)

	// GetAggregatedSearchMetrics возвращает агрегированные метрики поиска
	GetAggregatedSearchMetrics(ctx context.Context, periodStart, periodEnd time.Time) (*behavior.AggregatedSearchMetrics, error)

	// GetTopSearchQueries возвращает топ поисковых запросов с полной статистикой
	GetTopSearchQueries(ctx context.Context, periodStart, periodEnd time.Time, limit int) ([]behavior.TopSearchQuery, error)
}
