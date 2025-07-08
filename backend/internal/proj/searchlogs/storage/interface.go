// Package storage backend/internal/proj/searchlogs/storage/interface.go
package storage

import (
	"context"
	"time"

	"backend/internal/proj/searchlogs/types"
)

// Interface определяет интерфейс хранилища для логов поиска
type Interface interface {
	// SaveBatch сохраняет батч записей логов
	SaveBatch(ctx context.Context, entries []*types.SearchLogEntry) error

	// GetSearchStats возвращает статистику поиска за период
	GetSearchStats(ctx context.Context, from, to time.Time) (*types.SearchStats, error)

	// GetPopularSearches возвращает популярные поисковые запросы
	GetPopularSearches(ctx context.Context, from, to time.Time, limit int) ([]types.PopularSearch, error)

	// GetUserSearchHistory возвращает историю поиска пользователя
	GetUserSearchHistory(ctx context.Context, userID int, limit int) ([]types.SearchLogEntry, error)
}
