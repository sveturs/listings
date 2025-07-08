package service

import (
	"context"
	"time"

	"backend/internal/domain/behavior"
)

// BehaviorTrackingService определяет интерфейс для работы с поведенческими событиями
type BehaviorTrackingService interface {
	// TrackEvent отслеживает поведенческое событие
	TrackEvent(ctx context.Context, userID *int, req *behavior.TrackEventRequest) error

	// GetSearchMetrics возвращает метрики поиска
	GetSearchMetrics(ctx context.Context, query *behavior.SearchMetricsQuery) ([]*behavior.SearchMetrics, int, error)

	// GetItemMetrics возвращает метрики товаров
	GetItemMetrics(ctx context.Context, query *behavior.ItemMetricsQuery) ([]*behavior.ItemMetrics, int, error)

	// UpdateSearchMetrics обновляет агрегированные метрики поиска за период
	UpdateSearchMetrics(ctx context.Context, periodStart, periodEnd time.Time) error

	// GetUserEvents возвращает события пользователя
	GetUserEvents(ctx context.Context, userID int, limit, offset int) ([]*behavior.BehaviorEvent, int, error)

	// GetSessionEvents возвращает события сессии
	GetSessionEvents(ctx context.Context, sessionID string) ([]*behavior.BehaviorEvent, error)
}
