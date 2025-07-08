// Package service backend/internal/proj/searchlogs/service/interface.go
package service

import (
	"context"
	"time"

	"backend/internal/proj/searchlogs/types"
)

// ServiceInterface определяет интерфейс сервиса логирования поиска
type ServiceInterface interface {
	// LogSearch асинхронно логирует поисковый запрос
	LogSearch(ctx context.Context, entry *types.SearchLogEntry) error

	// GetSearchStats возвращает статистику поиска за период
	GetSearchStats(ctx context.Context, from, to time.Time) (*types.SearchStats, error)

	// GetPopularSearches возвращает популярные поисковые запросы
	GetPopularSearches(ctx context.Context, limit int, period time.Duration) ([]types.PopularSearch, error)

	// GetUserSearchHistory возвращает историю поиска пользователя
	GetUserSearchHistory(ctx context.Context, userID int, limit int) ([]types.SearchLogEntry, error)
}
