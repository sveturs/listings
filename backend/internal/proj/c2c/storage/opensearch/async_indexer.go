// backend/internal/proj/c2c/storage/opensearch/async_indexer.go
package opensearch

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// IndexTask представляет задачу индексации
type IndexTask struct {
	ListingID int
	Action    string // "index", "delete"
	Data      *models.MarketplaceListing
	Attempt   int
	CreatedAt time.Time
}

// AsyncIndexer управляет асинхронной индексацией в OpenSearch
type AsyncIndexer struct {
	taskQueue chan IndexTask
	workers   int
	repo      *Repository
	db        *sqlx.DB
	wg        sync.WaitGroup
	shutdown  chan struct{}
	once      sync.Once

	// Prometheus metrics
	queueSize      prometheus.Gauge
	successCounter prometheus.Counter
	failureCounter prometheus.Counter
	retryCounter   prometheus.Counter
	latencyHist    prometheus.Histogram
}

// Prometheus metrics (singleton для регистрации)
var (
	metricsOnce sync.Once
	metrics     struct {
		queueSize      prometheus.Gauge
		successCounter prometheus.Counter
		failureCounter prometheus.Counter
		retryCounter   prometheus.Counter
		latencyHist    prometheus.Histogram
	}
)

func initMetrics() {
	metricsOnce.Do(func() {
		metrics.queueSize = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "opensearch_indexing_queue_size",
			Help: "Current size of the indexing queue",
		})
		metrics.successCounter = promauto.NewCounter(prometheus.CounterOpts{
			Name: "opensearch_indexing_success_total",
			Help: "Total number of successful indexing operations",
		})
		metrics.failureCounter = promauto.NewCounter(prometheus.CounterOpts{
			Name: "opensearch_indexing_failure_total",
			Help: "Total number of failed indexing operations",
		})
		metrics.retryCounter = promauto.NewCounter(prometheus.CounterOpts{
			Name: "opensearch_indexing_retry_total",
			Help: "Total number of retry attempts",
		})
		metrics.latencyHist = promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "opensearch_indexing_latency_seconds",
			Help:    "Indexing operation latency in seconds",
			Buckets: prometheus.DefBuckets,
		})
	})
}

// NewAsyncIndexer создаёт новый асинхронный индексер
func NewAsyncIndexer(repo *Repository, db *sqlx.DB, workers int, queueSize int) *AsyncIndexer {
	if workers <= 0 {
		workers = 5 // default
	}
	if queueSize <= 0 {
		queueSize = 1000 // default
	}

	initMetrics()

	indexer := &AsyncIndexer{
		taskQueue:      make(chan IndexTask, queueSize),
		workers:        workers,
		repo:           repo,
		db:             db,
		shutdown:       make(chan struct{}),
		queueSize:      metrics.queueSize,
		successCounter: metrics.successCounter,
		failureCounter: metrics.failureCounter,
		retryCounter:   metrics.retryCounter,
		latencyHist:    metrics.latencyHist,
	}

	// Запускаем workers
	for i := 0; i < workers; i++ {
		indexer.wg.Add(1)
		go indexer.worker(i)
	}

	logger.Info().
		Int("workers", workers).
		Int("queueSize", queueSize).
		Msg("Async indexer started")

	return indexer
}

// Enqueue добавляет задачу в очередь
func (ai *AsyncIndexer) Enqueue(task IndexTask) error {
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	select {
	case ai.taskQueue <- task:
		ai.queueSize.Set(float64(len(ai.taskQueue)))
		logger.Debug().
			Int("listingID", task.ListingID).
			Str("action", task.Action).
			Msg("Task enqueued")
		return nil
	case <-ai.shutdown:
		return fmt.Errorf("indexer is shutting down")
	default:
		// Очередь переполнена - fallback на синхронную индексацию
		logger.Warn().
			Int("listingID", task.ListingID).
			Msg("Queue full, falling back to sync indexing")
		return ai.executeSyncIndexing(task)
	}
}

// worker обрабатывает задачи из очереди
func (ai *AsyncIndexer) worker(id int) {
	defer ai.wg.Done()

	logger.Info().Int("workerID", id).Msg("Worker started")

	for {
		select {
		case task := <-ai.taskQueue:
			ai.queueSize.Set(float64(len(ai.taskQueue)))
			ai.processTask(task)
		case <-ai.shutdown:
			logger.Info().Int("workerID", id).Msg("Worker shutting down")
			return
		}
	}
}

// processTask обрабатывает одну задачу с retry механизмом
func (ai *AsyncIndexer) processTask(task IndexTask) {
	// Защита от panic (например, если storage или client nil в тестах)
	defer func() {
		if r := recover(); r != nil {
			logger.Error().
				Int("listingID", task.ListingID).
				Str("action", task.Action).
				Interface("panic", r).
				Msg("Panic recovered in processTask")

			// Трактуем panic как ошибку и применяем retry механизм
			err := fmt.Errorf("panic recovered: %v", r)
			ai.handleFailure(task, err)
		}
	}()

	start := time.Now()
	ctx := context.Background()

	logger.Info().
		Int("listingID", task.ListingID).
		Str("action", task.Action).
		Int("attempt", task.Attempt).
		Msg("Processing indexing task")

	var err error
	switch task.Action {
	case "index":
		if task.Data == nil {
			// Если данных нет, загружаем listing из БД
			task.Data, err = ai.fetchListing(ctx, task.ListingID)
			if err != nil {
				logger.Error().
					Err(err).
					Int("listingID", task.ListingID).
					Msg("Failed to fetch listing for indexing")
				ai.handleFailure(task, err)
				return
			}
		}
		err = ai.repo.IndexListing(ctx, task.Data)
	case "delete":
		err = ai.repo.DeleteListing(ctx, fmt.Sprintf("%d", task.ListingID))
	default:
		err = fmt.Errorf("unknown action: %s", task.Action)
	}

	duration := time.Since(start)
	ai.latencyHist.Observe(duration.Seconds())

	if err != nil {
		logger.Error().
			Err(err).
			Int("listingID", task.ListingID).
			Str("action", task.Action).
			Int("attempt", task.Attempt).
			Dur("duration", duration).
			Msg("Indexing task failed")

		ai.handleFailure(task, err)
	} else {
		logger.Info().
			Int("listingID", task.ListingID).
			Str("action", task.Action).
			Dur("duration", duration).
			Msg("Indexing task completed successfully")

		ai.successCounter.Inc()
	}
}

// handleFailure обрабатывает ошибки с retry механизмом
func (ai *AsyncIndexer) handleFailure(task IndexTask, err error) {
	task.Attempt++

	// Retry policy: 3 attempts with exponential backoff
	if task.Attempt < 3 {
		ai.retryCounter.Inc()

		// Exponential backoff: 0s, 1s, 5s
		var delay time.Duration
		switch task.Attempt {
		case 1:
			delay = 0
		case 2:
			delay = 1 * time.Second
		default:
			delay = 5 * time.Second
		}

		logger.Warn().
			Int("listingID", task.ListingID).
			Int("attempt", task.Attempt).
			Dur("delay", delay).
			Msg("Retrying indexing task")

		// Retry after delay
		if delay > 0 {
			time.Sleep(delay)
		}

		select {
		case ai.taskQueue <- task:
			ai.queueSize.Set(float64(len(ai.taskQueue)))
		case <-ai.shutdown:
			// Shutdown в процессе retry - записываем в DLQ
			ai.saveToDLQ(task, err)
		default:
			// Очередь переполнена - сразу в DLQ
			ai.saveToDLQ(task, err)
		}
	} else {
		// Все попытки исчерпаны - записываем в Dead Letter Queue
		ai.failureCounter.Inc()
		ai.saveToDLQ(task, err)
	}
}

// saveToDLQ сохраняет failed task в Dead Letter Queue (PostgreSQL)
func (ai *AsyncIndexer) saveToDLQ(task IndexTask, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dataJSON []byte
	if task.Data != nil {
		dataJSON, _ = json.Marshal(task.Data)
	}

	query := `
		INSERT INTO opensearch_indexing_dlq
		(listing_id, action, data, attempts, last_error, created_at, last_attempt_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (listing_id, action)
		DO UPDATE SET
			attempts = opensearch_indexing_dlq.attempts + 1,
			last_error = $5,
			last_attempt_at = $7,
			data = $3
	`

	_, dbErr := ai.db.ExecContext(
		ctx,
		query,
		task.ListingID,
		task.Action,
		dataJSON,
		task.Attempt,
		err.Error(),
		task.CreatedAt,
		time.Now(),
	)

	if dbErr != nil {
		logger.Error().
			Err(dbErr).
			Int("listingID", task.ListingID).
			Msg("Failed to save task to DLQ")
	} else {
		logger.Warn().
			Int("listingID", task.ListingID).
			Str("action", task.Action).
			Int("attempts", task.Attempt).
			Str("error", err.Error()).
			Msg("Task saved to Dead Letter Queue")
	}
}

// RetryDLQ повторно обрабатывает задачи из DLQ
func (ai *AsyncIndexer) RetryDLQ(ctx context.Context, limit int) error {
	query := `
		SELECT id, listing_id, action, data, attempts, last_error, created_at, last_attempt_at
		FROM opensearch_indexing_dlq
		WHERE attempts < 10  -- Максимум 10 попыток для DLQ
		ORDER BY created_at ASC
		LIMIT $1
	`

	var dlqTasks []struct {
		ID            int             `db:"id"`
		ListingID     int             `db:"listing_id"`
		Action        string          `db:"action"`
		Data          json.RawMessage `db:"data"`
		Attempts      int             `db:"attempts"`
		LastError     string          `db:"last_error"`
		CreatedAt     time.Time       `db:"created_at"`
		LastAttemptAt time.Time       `db:"last_attempt_at"`
	}

	err := ai.db.SelectContext(ctx, &dlqTasks, query, limit)
	if err != nil {
		return fmt.Errorf("failed to fetch DLQ tasks: %w", err)
	}

	logger.Info().
		Int("count", len(dlqTasks)).
		Msg("Retrying tasks from DLQ")

	for _, dt := range dlqTasks {
		task := IndexTask{
			ListingID: dt.ListingID,
			Action:    dt.Action,
			Attempt:   0, // Сбрасываем счётчик для нового цикла retry
			CreatedAt: dt.CreatedAt,
		}

		// Десериализуем data если есть
		if len(dt.Data) > 0 {
			var listing models.MarketplaceListing
			if err := json.Unmarshal(dt.Data, &listing); err == nil {
				task.Data = &listing
			}
		}

		// Пытаемся enqueue задачу
		if err := ai.Enqueue(task); err != nil {
			logger.Error().
				Err(err).
				Int("dlqID", dt.ID).
				Int("listingID", dt.ListingID).
				Msg("Failed to re-enqueue DLQ task")
			continue
		}

		// Удаляем из DLQ после успешного enqueue
		_, err := ai.db.ExecContext(ctx, "DELETE FROM opensearch_indexing_dlq WHERE id = $1", dt.ID)
		if err != nil {
			logger.Error().
				Err(err).
				Int("dlqID", dt.ID).
				Msg("Failed to delete DLQ task after re-enqueue")
		}
	}

	return nil
}

// Shutdown gracefully останавливает indexer
func (ai *AsyncIndexer) Shutdown(timeout time.Duration) error {
	ai.once.Do(func() {
		logger.Info().Msg("Shutting down async indexer...")

		close(ai.shutdown)

		// Ждём завершения всех workers с timeout
		done := make(chan struct{})
		go func() {
			ai.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			logger.Info().Msg("All workers stopped gracefully")
		case <-time.After(timeout):
			logger.Warn().Msg("Shutdown timeout exceeded, forcing stop")
		}

		// Обрабатываем оставшиеся задачи в очереди
		remaining := len(ai.taskQueue)
		if remaining > 0 {
			logger.Warn().
				Int("remaining", remaining).
				Msg("Saving remaining tasks to DLQ")

			for i := 0; i < remaining; i++ {
				select {
				case task := <-ai.taskQueue:
					ai.saveToDLQ(task, fmt.Errorf("indexer shutdown"))
				default:
					break
				}
			}
		}

		close(ai.taskQueue)
	})

	return nil
}

// fetchListing загружает listing из БД для индексации
func (ai *AsyncIndexer) fetchListing(ctx context.Context, listingID int) (*models.MarketplaceListing, error) {
	// Используем метод из storage через repo.storage
	listing, err := ai.repo.storage.GetListingByID(ctx, listingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing not found: %d", listingID)
		}
		return nil, fmt.Errorf("failed to fetch listing %d: %w", listingID, err)
	}

	// Загружаем связанные данные
	if err := ai.loadListingRelations(ctx, listing); err != nil {
		logger.Warn().
			Err(err).
			Int("listingID", listingID).
			Msg("Failed to load listing relations")
	}

	return listing, nil
}

// loadListingRelations загружает связанные данные для listing
func (ai *AsyncIndexer) loadListingRelations(ctx context.Context, listing *models.MarketplaceListing) error {
	var err error

	// Загружаем переводы
	translations, err := ai.repo.storage.GetTranslationsForEntity(ctx, "listing", listing.ID)
	if err == nil && len(translations) > 0 {
		transMap := make(models.TranslationMap)
		for _, t := range translations {
			if _, ok := transMap[t.Language]; !ok {
				transMap[t.Language] = make(map[string]string)
			}
			transMap[t.Language][t.FieldName] = t.TranslatedText
		}
		listing.Translations = transMap
	}

	// Загружаем атрибуты
	attrs, err := ai.repo.storage.GetListingAttributes(ctx, listing.ID)
	if err == nil && len(attrs) > 0 {
		listing.Attributes = attrs
	}

	// Загружаем изображения
	images, err := ai.repo.storage.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
	if err == nil && len(images) > 0 {
		listing.Images = images
	}

	return nil
}

// executeSyncIndexing выполняет индексацию синхронно (fallback)
func (ai *AsyncIndexer) executeSyncIndexing(task IndexTask) error {
	// Защита от panic в sync fallback (например, если storage или client nil в тестах)
	defer func() {
		if r := recover(); r != nil {
			logger.Error().
				Int("listingID", task.ListingID).
				Str("action", task.Action).
				Interface("panic", r).
				Msg("Panic recovered in executeSyncIndexing")
		}
	}()

	ctx := context.Background()
	start := time.Now()

	logger.Warn().
		Int("listingID", task.ListingID).
		Msg("Executing sync indexing (fallback)")

	var err error
	switch task.Action {
	case "index":
		if task.Data == nil {
			task.Data, err = ai.fetchListing(ctx, task.ListingID)
			if err != nil {
				return fmt.Errorf("failed to fetch listing: %w", err)
			}
		}
		err = ai.repo.IndexListing(ctx, task.Data)
	case "delete":
		err = ai.repo.DeleteListing(ctx, fmt.Sprintf("%d", task.ListingID))
	default:
		return fmt.Errorf("unknown action: %s", task.Action)
	}

	duration := time.Since(start)
	ai.latencyHist.Observe(duration.Seconds())

	if err != nil {
		ai.failureCounter.Inc()
		return err
	}

	ai.successCounter.Inc()
	logger.Info().
		Int("listingID", task.ListingID).
		Dur("duration", duration).
		Msg("Sync indexing completed")

	return nil
}

// GetQueueSize возвращает текущий размер очереди
func (ai *AsyncIndexer) GetQueueSize() int {
	return len(ai.taskQueue)
}

// IsHealthy проверяет здоровье indexer
func (ai *AsyncIndexer) IsHealthy() bool {
	select {
	case <-ai.shutdown:
		return false
	default:
		// Проверяем, не переполнена ли очередь
		queueUsage := float64(len(ai.taskQueue)) / float64(cap(ai.taskQueue))
		return queueUsage < 0.9 // Здоровье OK если очередь заполнена менее чем на 90%
	}
}
