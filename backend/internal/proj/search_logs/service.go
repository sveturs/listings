package searchlogs

import (
	"context"
	"sync"
	"time"

	"backend/internal/domain"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// Service представляет сервис для работы с логами поисковых запросов
type Service struct {
	repo         Repository
	logger       *logrus.Logger
	logQueue     chan *domain.SearchLogInput
	clickQueue   chan *domain.SearchResultClick
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	batchSize    int
	flushTimeout time.Duration
	cron         *cron.Cron
}

// Repository определяет интерфейс для работы с хранилищем логов
type Repository interface {
	CreateSearchLog(ctx context.Context, log *domain.SearchLogInput) (*domain.SearchLog, error)
	CreateSearchLogBatch(ctx context.Context, logs []*domain.SearchLogInput) error
	CreateClickLog(ctx context.Context, click *domain.SearchResultClick) error
	CreateClickLogBatch(ctx context.Context, clicks []*domain.SearchResultClick) error
	GetTrendingQueries(ctx context.Context, limit int, categoryID *int, country *string) ([]*domain.SearchTrendingQuery, error)
	UpdateAnalytics(ctx context.Context) error
	UpdateTrendingQueries(ctx context.Context) error
}

// NewService создает новый экземпляр сервиса логирования поисковых запросов
func NewService(repo Repository, logger *logrus.Logger) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Service{
		repo:         repo,
		logger:       logger,
		logQueue:     make(chan *domain.SearchLogInput, 10000),
		clickQueue:   make(chan *domain.SearchResultClick, 10000),
		ctx:          ctx,
		cancel:       cancel,
		batchSize:    100,
		flushTimeout: 5 * time.Second,
		cron:         cron.New(),
	}

	// Запускаем фоновые воркеры
	s.startWorkers()

	// Настраиваем периодические задачи
	s.setupCronJobs()

	return s
}

// LogSearch асинхронно логирует поисковый запрос
func (s *Service) LogSearch(input *domain.SearchLogInput) {
	select {
	case s.logQueue <- input:
		// Успешно добавлено в очередь
	default:
		// Очередь переполнена, логируем ошибку
		s.logger.WithFields(logrus.Fields{
			"query":   input.QueryText,
			"user_id": input.UserID,
		}).Warn("Search log queue is full, dropping log entry")
	}
}

// LogClick асинхронно логирует клик по результату поиска
func (s *Service) LogClick(searchLogID int64, listingID int, position int) {
	click := &domain.SearchResultClick{
		SearchLogID: searchLogID,
		ListingID:   listingID,
		Position:    position,
		ClickedAt:   time.Now(),
	}

	select {
	case s.clickQueue <- click:
		// Успешно добавлено в очередь
	default:
		// Очередь переполнена, логируем ошибку
		s.logger.WithFields(logrus.Fields{
			"search_log_id": searchLogID,
			"listing_id":    listingID,
		}).Warn("Click log queue is full, dropping click entry")
	}
}

// GetTrendingQueries возвращает популярные поисковые запросы
func (s *Service) GetTrendingQueries(ctx context.Context, limit int, categoryID *int, country *string) ([]*domain.SearchTrendingQuery, error) {
	return s.repo.GetTrendingQueries(ctx, limit, categoryID, country)
}

// startWorkers запускает фоновые воркеры для обработки логов
func (s *Service) startWorkers() {
	// Воркер для обработки логов поисковых запросов
	s.wg.Add(1)
	go s.searchLogWorker()

	// Воркер для обработки кликов
	s.wg.Add(1)
	go s.clickLogWorker()
}

// searchLogWorker обрабатывает очередь логов поисковых запросов
func (s *Service) searchLogWorker() {
	defer s.wg.Done()

	batch := make([]*domain.SearchLogInput, 0, s.batchSize)
	ticker := time.NewTicker(s.flushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			// Сохраняем оставшиеся логи перед выходом
			if len(batch) > 0 {
				s.flushSearchLogs(batch)
			}
			return

		case log := <-s.logQueue:
			batch = append(batch, log)
			if len(batch) >= s.batchSize {
				s.flushSearchLogs(batch)
				batch = make([]*domain.SearchLogInput, 0, s.batchSize)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				s.flushSearchLogs(batch)
				batch = make([]*domain.SearchLogInput, 0, s.batchSize)
			}
		}
	}
}

// clickLogWorker обрабатывает очередь кликов
func (s *Service) clickLogWorker() {
	defer s.wg.Done()

	batch := make([]*domain.SearchResultClick, 0, s.batchSize)
	ticker := time.NewTicker(s.flushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			// Сохраняем оставшиеся клики перед выходом
			if len(batch) > 0 {
				s.flushClickLogs(batch)
			}
			return

		case click := <-s.clickQueue:
			batch = append(batch, click)
			if len(batch) >= s.batchSize {
				s.flushClickLogs(batch)
				batch = make([]*domain.SearchResultClick, 0, s.batchSize)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				s.flushClickLogs(batch)
				batch = make([]*domain.SearchResultClick, 0, s.batchSize)
			}
		}
	}
}

// flushSearchLogs сохраняет батч логов поисковых запросов
func (s *Service) flushSearchLogs(batch []*domain.SearchLogInput) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.repo.CreateSearchLogBatch(ctx, batch); err != nil {
		s.logger.WithError(err).WithField("batch_size", len(batch)).
			Error("Failed to save search log batch")

		// В случае ошибки пытаемся сохранить по одному
		for _, log := range batch {
			if _, err := s.repo.CreateSearchLog(ctx, log); err != nil {
				s.logger.WithError(err).WithField("query", log.QueryText).
					Error("Failed to save individual search log")
			}
		}
	}
}

// flushClickLogs сохраняет батч кликов
func (s *Service) flushClickLogs(batch []*domain.SearchResultClick) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.repo.CreateClickLogBatch(ctx, batch); err != nil {
		s.logger.WithError(err).WithField("batch_size", len(batch)).
			Error("Failed to save click log batch")

		// В случае ошибки пытаемся сохранить по одному
		for _, click := range batch {
			if err := s.repo.CreateClickLog(ctx, click); err != nil {
				s.logger.WithError(err).WithField("listing_id", click.ListingID).
					Error("Failed to save individual click log")
			}
		}
	}
}

// setupCronJobs настраивает периодические задачи
func (s *Service) setupCronJobs() {
	// Обновление аналитики каждый час
	_, err := s.cron.AddFunc("0 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := s.repo.UpdateAnalytics(ctx); err != nil {
			s.logger.WithError(err).Error("Failed to update search analytics")
		}
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to add analytics cron job")
	}

	// Обновление трендовых запросов каждые 15 минут
	_, err = s.cron.AddFunc("*/15 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if err := s.repo.UpdateTrendingQueries(ctx); err != nil {
			s.logger.WithError(err).Error("Failed to update trending queries")
		}
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to add trending queries cron job")
	}

	s.cron.Start()
}

// Stop останавливает сервис и сохраняет все незаписанные логи
func (s *Service) Stop() {
	s.logger.Info("Stopping search log service...")

	// Останавливаем cron
	s.cron.Stop()

	// Отменяем контекст для остановки воркеров
	s.cancel()

	// Ждем завершения всех воркеров
	s.wg.Wait()

	// Закрываем каналы
	close(s.logQueue)
	close(s.clickQueue)

	s.logger.Info("Search log service stopped")
}

// PrepareSearchLogInput подготавливает входные данные для логирования из параметров запроса
func PrepareSearchLogInput(userID *int, sessionID, queryText, userAgent, ipAddress, referer string,
	filters map[string]interface{}, categoryID *int, location *domain.SearchLocation,
	resultsCount int, responseTimeMs int, page int, perPage int, sortBy *string,
) *domain.SearchLogInput {
	input := &domain.SearchLogInput{
		UserID:         userID,
		SessionID:      sessionID,
		QueryText:      queryText,
		CategoryID:     categoryID,
		ResultsCount:   resultsCount,
		ResponseTimeMs: responseTimeMs,
		Page:           page,
		PerPage:        perPage,
		SortBy:         sortBy,
		UserAgent:      userAgent,
		IPAddress:      ipAddress,
		Referer:        referer,
	}

	// Конвертируем фильтры в SearchFilters
	if len(filters) > 0 {
		searchFilters := &domain.SearchFilters{
			Attributes: make(map[string]interface{}),
		}

		for key, value := range filters {
			switch key {
			case "price_min":
				if v, ok := value.(float64); ok {
					searchFilters.PriceMin = &v
				}
			case "price_max":
				if v, ok := value.(float64); ok {
					searchFilters.PriceMax = &v
				}
			case "tags":
				if tags, ok := value.([]string); ok {
					searchFilters.Tags = tags
				}
			default:
				searchFilters.Attributes[key] = value
			}
		}

		input.Filters = searchFilters
	}

	if location != nil {
		input.Location = location
	}

	return input
}
