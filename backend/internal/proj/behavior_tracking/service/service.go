package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"backend/internal/domain/behavior"
	"backend/internal/logger"
	"backend/internal/proj/behavior_tracking/storage/postgres"
)

// behaviorTrackingService реализация сервиса поведенческих событий
type behaviorTrackingService struct {
	repo      postgres.BehaviorTrackingRepository
	validator *validator.Validate

	// Для батчевой обработки
	eventBuffer []*behavior.BehaviorEvent
	bufferMutex sync.Mutex
	bufferSize  int
	flushTicker *time.Ticker
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewBehaviorTrackingService создает новый сервис
func NewBehaviorTrackingService(ctx context.Context, repo postgres.BehaviorTrackingRepository) BehaviorTrackingService {
	ctx, cancel := context.WithCancel(ctx)

	s := &behaviorTrackingService{
		repo:        repo,
		validator:   validator.New(),
		eventBuffer: make([]*behavior.BehaviorEvent, 0, 100),
		bufferSize:  100,
		flushTicker: time.NewTicker(5 * time.Second), // Флашим буфер каждые 5 секунд
		ctx:         ctx,
		cancel:      cancel,
	}

	// Запускаем горутину для периодического флаша буфера
	go s.flushWorker() //nolint:contextcheck // использует внутренний контекст s.ctx

	return s
}

// flushWorker периодически сохраняет накопленные события
func (s *behaviorTrackingService) flushWorker() {
	for {
		select {
		case <-s.flushTicker.C:
			if err := s.flushBuffer(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Failed to flush event buffer")
			}
		case <-s.ctx.Done():
			// Сохраняем оставшиеся события перед завершением
			if err := s.flushBuffer(context.Background()); err != nil {
				logger.Error().Err(err).Msg("Failed to flush event buffer on shutdown")
			}
			return
		}
	}
}

// flushBuffer сохраняет накопленные события в БД
func (s *behaviorTrackingService) flushBuffer(ctx context.Context) error {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	if len(s.eventBuffer) == 0 {
		return nil
	}

	// Копируем буфер для сохранения
	events := make([]*behavior.BehaviorEvent, len(s.eventBuffer))
	copy(events, s.eventBuffer)

	// Очищаем буфер
	s.eventBuffer = s.eventBuffer[:0]

	// Сохраняем события пакетом с увеличенным таймаутом
	flushCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	logger.Info().
		Int("events_to_save", len(events)).
		Msg("flushBuffer: attempting to save behavior events to database")

	if err := s.repo.SaveEventsBatch(flushCtx, events); err != nil {
		// При ошибке conn busy просто логируем и не возвращаем события в буфер
		// чтобы избежать накопления событий
		if strings.Contains(err.Error(), "conn busy") {
			logger.Warn().
				Int("count", len(events)).
				Err(err).
				Msg("Dropping events due to connection busy")
			return err
		}

		// Для других ошибок возвращаем события обратно в буфер
		s.eventBuffer = append(s.eventBuffer, events...)
		logger.Error().
			Int("count", len(events)).
			Err(err).
			Msg("Failed to save events batch, returned to buffer")
		return fmt.Errorf("failed to save events batch: %w", err)
	}

	logger.Info().Int("count", len(events)).Msg("Successfully flushed behavior events to database")
	return nil
}

// TrackEvent отслеживает поведенческое событие
func (s *behaviorTrackingService) TrackEvent(ctx context.Context, userID *int, req *behavior.TrackEventRequest) error {
	// Валидация запроса
	if err := s.validator.Struct(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Генерируем session_id если не передан
	if req.SessionID == "" {
		req.SessionID = uuid.New().String()
	}

	// Создаем событие
	event := &behavior.BehaviorEvent{
		EventType:   req.EventType,
		UserID:      userID,
		SessionID:   req.SessionID,
		SearchQuery: req.SearchQuery,
		ItemID:      req.ItemID,
		ItemType:    req.ItemType,
		Position:    req.Position,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
	}

	// Добавляем в буфер
	s.bufferMutex.Lock()
	s.eventBuffer = append(s.eventBuffer, event)
	currentBufferSize := len(s.eventBuffer)
	shouldFlush := currentBufferSize >= s.bufferSize
	s.bufferMutex.Unlock()

	logger.Debug().
		Str("event_type", string(event.EventType)).
		Int("buffer_size", currentBufferSize).
		Int("max_buffer_size", s.bufferSize).
		Bool("should_flush", shouldFlush).
		Msg("Added behavior event to buffer")

	// Если буфер заполнен, сохраняем
	if shouldFlush {
		logger.Info().Msg("Buffer is full, triggering flush")
		go func() {
			if err := s.flushBuffer(ctx); err != nil {
				logger.Error().Err(err).Msg("Failed to flush full buffer")
			}
		}()
	}

	return nil
}

// GetSearchMetrics возвращает метрики поиска
func (s *behaviorTrackingService) GetSearchMetrics(ctx context.Context, query *behavior.SearchMetricsQuery) ([]*behavior.SearchMetrics, int, error) {
	// Валидация параметров запроса
	if err := s.validator.Struct(query); err != nil {
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	// Устанавливаем значения по умолчанию
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	if query.OrderBy == "" {
		query.OrderBy = "desc"
	}

	// Если период не указан, берем последние 7 дней
	if query.PeriodStart.IsZero() {
		query.PeriodStart = time.Now().AddDate(0, 0, -7)
	}
	if query.PeriodEnd.IsZero() {
		query.PeriodEnd = time.Now()
	}

	logger.Info().
		Time("period_start", query.PeriodStart).
		Time("period_end", query.PeriodEnd).
		Int("limit", query.Limit).
		Str("sort_by", query.SortBy).
		Msg("GetSearchMetrics: query parameters")

	return s.repo.GetSearchMetrics(ctx, query)
}

// GetItemMetrics возвращает метрики товаров
func (s *behaviorTrackingService) GetItemMetrics(ctx context.Context, query *behavior.ItemMetricsQuery) ([]*behavior.ItemMetrics, int, error) {
	// Валидация параметров запроса
	if err := s.validator.Struct(query); err != nil {
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	// Устанавливаем значения по умолчанию
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	if query.OrderBy == "" {
		query.OrderBy = "desc"
	}

	// Если период не указан, берем последние 7 дней
	if query.PeriodStart.IsZero() {
		query.PeriodStart = time.Now().AddDate(0, 0, -7)
	}
	if query.PeriodEnd.IsZero() {
		query.PeriodEnd = time.Now()
	}

	return s.repo.GetItemMetrics(ctx, query)
}

// UpdateSearchMetrics обновляет агрегированные метрики поиска за период
func (s *behaviorTrackingService) UpdateSearchMetrics(ctx context.Context, periodStart, periodEnd time.Time) error {
	// Форсируем сохранение буфера перед обновлением метрик
	if err := s.flushBuffer(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to flush buffer before updating metrics")
	}

	return s.repo.UpdateSearchMetrics(ctx, periodStart, periodEnd)
}

// GetUserEvents возвращает события пользователя
func (s *behaviorTrackingService) GetUserEvents(ctx context.Context, userID int, limit, offset int) ([]*behavior.BehaviorEvent, int, error) {
	// Валидация параметров
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.GetEventsByUser(ctx, userID, limit, offset)
}

// GetSessionEvents возвращает события сессии
func (s *behaviorTrackingService) GetSessionEvents(ctx context.Context, sessionID string) ([]*behavior.BehaviorEvent, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session_id is required")
	}

	return s.repo.GetEventsBySession(ctx, sessionID)
}

// GetAggregatedSearchMetrics возвращает агрегированные метрики поиска
func (s *behaviorTrackingService) GetAggregatedSearchMetrics(ctx context.Context, query *behavior.SearchMetricsQuery) (*behavior.AggregatedSearchMetrics, error) {
	// Устанавливаем значения по умолчанию
	if query.PeriodStart.IsZero() {
		query.PeriodStart = time.Now().AddDate(0, 0, -7)
	}
	if query.PeriodEnd.IsZero() {
		query.PeriodEnd = time.Now()
	}

	// Получаем агрегированные метрики из репозитория
	return s.repo.GetAggregatedSearchMetrics(ctx, query.PeriodStart, query.PeriodEnd)
}

// GetTopSearchQueries возвращает топ поисковых запросов с полной статистикой
func (s *behaviorTrackingService) GetTopSearchQueries(ctx context.Context, query *behavior.SearchMetricsQuery) ([]behavior.TopSearchQuery, error) {
	// Устанавливаем значения по умолчанию
	if query.PeriodStart.IsZero() {
		query.PeriodStart = time.Now().AddDate(0, 0, -7)
	}
	if query.PeriodEnd.IsZero() {
		query.PeriodEnd = time.Now()
	}
	if query.Limit <= 0 {
		query.Limit = 50
	}

	// Получаем топ запросы из репозитория
	return s.repo.GetTopSearchQueries(ctx, query.PeriodStart, query.PeriodEnd, query.Limit)
}

// Close завершает работу сервиса
func (s *behaviorTrackingService) Close() error {
	// Останавливаем флаш воркер
	s.cancel()

	// Даем время на завершение
	time.Sleep(100 * time.Millisecond)

	// Финальный флаш буфера
	return s.flushBuffer(context.Background())
}
