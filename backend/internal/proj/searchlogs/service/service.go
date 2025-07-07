// Package service backend/internal/proj/searchlogs/service/service.go
package service

import (
	"context"
	"time"

	"backend/internal/logger"
	"backend/internal/proj/searchlogs/storage"
	"backend/internal/proj/searchlogs/types"
)

// Service реализует сервис логирования поиска
type Service struct {
	storage storage.Interface
	// Канал для асинхронной обработки логов
	logChan chan *types.SearchLogEntry
	// Размер буфера канала
	bufferSize int
}

// NewService создает новый экземпляр сервиса
func NewService(storage storage.Interface) *Service {
	s := &Service{
		storage:    storage,
		bufferSize: 1000, // Буфер на 1000 записей
		logChan:    make(chan *types.SearchLogEntry, 1000),
	}

	// Запускаем горутину для обработки логов
	go s.processLogs()

	return s
}

// LogSearch асинхронно логирует поисковый запрос
func (s *Service) LogSearch(ctx context.Context, entry *types.SearchLogEntry) error {
	// Устанавливаем текущее время, если не задано
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	logger.Info().
		Str("query", entry.Query).
		Int("results", entry.ResultCount).
		Msg("LogSearch called")

	// Отправляем в канал неблокирующим способом
	select {
	case s.logChan <- entry:
		// Успешно отправлено
		logger.Info().Msg("Search log entry sent to channel")
		return nil
	default:
		// Канал переполнен, логируем предупреждение
		logger.Warn().Msg("Search log channel is full, dropping log entry")
		return nil // Не возвращаем ошибку, чтобы не блокировать основной запрос
	}
}

// processLogs обрабатывает логи из канала
func (s *Service) processLogs() {
	logger.Info().Msg("Search log processor started")

	// Буфер для батчевой вставки
	batch := make([]*types.SearchLogEntry, 0, 100)
	ticker := time.NewTicker(5 * time.Second) // Флашим каждые 5 секунд

	for {
		select {
		case entry := <-s.logChan:
			logger.Info().
				Str("query", entry.Query).
				Int("batch_size", len(batch)+1).
				Msg("Received log entry from channel")

			batch = append(batch, entry)

			// Если накопилось достаточно записей, сохраняем
			if len(batch) >= 100 {
				s.saveBatch(batch)
				batch = batch[:0] // Очищаем слайс
			}

		case <-ticker.C:
			// Периодически сохраняем накопленные записи
			if len(batch) > 0 {
				logger.Info().Int("batch_size", len(batch)).Msg("Flushing batch on timer")
				s.saveBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// saveBatch сохраняет батч записей в хранилище
func (s *Service) saveBatch(batch []*types.SearchLogEntry) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.storage.SaveBatch(ctx, batch); err != nil {
		logger.Error().Err(err).Int("batch_size", len(batch)).Msg("Failed to save search log batch")
	} else {
		logger.Debug().Int("batch_size", len(batch)).Msg("Search log batch saved successfully")
	}
}

// GetSearchStats возвращает статистику поиска за период
func (s *Service) GetSearchStats(ctx context.Context, from, to time.Time) (*types.SearchStats, error) {
	return s.storage.GetSearchStats(ctx, from, to)
}

// GetPopularSearches возвращает популярные поисковые запросы
func (s *Service) GetPopularSearches(ctx context.Context, limit int, period time.Duration) ([]types.PopularSearch, error) {
	from := time.Now().Add(-period)
	to := time.Now()
	return s.storage.GetPopularSearches(ctx, from, to, limit)
}

// GetUserSearchHistory возвращает историю поиска пользователя
func (s *Service) GetUserSearchHistory(ctx context.Context, userID int, limit int) ([]types.SearchLogEntry, error) {
	return s.storage.GetUserSearchHistory(ctx, userID, limit)
}

// Close закрывает сервис и сохраняет оставшиеся логи
func (s *Service) Close() error {
	close(s.logChan)

	// Сохраняем оставшиеся записи
	batch := make([]*types.SearchLogEntry, 0, 100)
	for entry := range s.logChan {
		batch = append(batch, entry)
		if len(batch) >= 100 {
			s.saveBatch(batch)
			batch = batch[:0]
		}
	}

	// Сохраняем последний батч
	if len(batch) > 0 {
		s.saveBatch(batch)
	}

	return nil
}
