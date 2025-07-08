package searchlogs

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Worker представляет фоновый обработчик задач
type Worker struct {
	service  *Service
	logger   *logrus.Logger
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewWorker создает новый фоновый обработчик для периодических задач
func NewWorker(service *Service, logger *logrus.Logger) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		service:  service,
		logger:   logger,
		interval: 15 * time.Minute, // Интервал по умолчанию
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start запускает фоновый обработчик
func (w *Worker) Start() {
	w.logger.Info("Starting search logs background worker...")

	// Запускаем обработчик аналитики (каждый час)
	w.wg.Add(1)
	go w.runAnalyticsWorker()

	// Запускаем обработчик трендов (каждые 15 минут)
	w.wg.Add(1)
	go w.runTrendsWorker()

	// Запускаем очистку старых логов (раз в день)
	w.wg.Add(1)
	go w.runCleanupWorker()

	w.logger.Info("Search logs background worker started")
}

// Stop останавливает фоновый обработчик
func (w *Worker) Stop() {
	w.logger.Info("Stopping search logs background worker...")
	w.cancel()
	w.wg.Wait()
	w.logger.Info("Search logs background worker stopped")
}

// runAnalyticsWorker периодически обновляет аналитику
func (w *Worker) runAnalyticsWorker() {
	defer w.wg.Done()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Выполняем первый запуск через 5 минут после старта
	firstRunTimer := time.NewTimer(5 * time.Minute)
	defer firstRunTimer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return

		case <-firstRunTimer.C:
			w.updateAnalytics()

		case <-ticker.C:
			w.updateAnalytics()
		}
	}
}

// runTrendsWorker периодически обновляет трендовые запросы
func (w *Worker) runTrendsWorker() {
	defer w.wg.Done()

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// Выполняем первый запуск через 2 минуты после старта
	firstRunTimer := time.NewTimer(2 * time.Minute)
	defer firstRunTimer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return

		case <-firstRunTimer.C:
			w.updateTrends()

		case <-ticker.C:
			w.updateTrends()
		}
	}
}

// runCleanupWorker периодически очищает старые логи
func (w *Worker) runCleanupWorker() {
	defer w.wg.Done()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Выполняем первый запуск через 1 час после старта
	firstRunTimer := time.NewTimer(1 * time.Hour)
	defer firstRunTimer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return

		case <-firstRunTimer.C:
			w.cleanupOldLogs()

		case <-ticker.C:
			w.cleanupOldLogs()
		}
	}
}

// updateAnalytics обновляет агрегированную аналитику
func (w *Worker) updateAnalytics() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	start := time.Now()
	w.logger.Info("Starting analytics update...")

	if err := w.service.repo.UpdateAnalytics(ctx); err != nil {
		w.logger.WithError(err).Error("Failed to update search analytics")
		return
	}

	duration := time.Since(start)
	w.logger.WithField("duration", duration).Info("Analytics update completed")

	// Также обновляем метрики производительности
	w.updatePerformanceMetrics(ctx)
}

// updateTrends обновляет трендовые запросы
func (w *Worker) updateTrends() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	start := time.Now()
	w.logger.Info("Starting trends update...")

	if err := w.service.repo.UpdateTrendingQueries(ctx); err != nil {
		w.logger.WithError(err).Error("Failed to update trending queries")
		return
	}

	duration := time.Since(start)
	w.logger.WithField("duration", duration).Info("Trends update completed")
}

// cleanupOldLogs удаляет старые логи для экономии места
func (w *Worker) cleanupOldLogs() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	start := time.Now()
	w.logger.Info("Starting old logs cleanup...")

	// Определяем retention период (по умолчанию 90 дней для логов, 1 год для аналитики)
	logsRetention := 90 * 24 * time.Hour
	analyticsRetention := 365 * 24 * time.Hour

	// Очищаем старые логи поиска
	deletedLogs, err := w.cleanupSearchLogs(ctx, logsRetention)
	if err != nil {
		w.logger.WithError(err).Error("Failed to cleanup search logs")
	} else {
		w.logger.WithField("deleted_count", deletedLogs).Info("Cleaned up old search logs")
	}

	// Очищаем старые клики
	deletedClicks, err := w.cleanupClickLogs(ctx, logsRetention)
	if err != nil {
		w.logger.WithError(err).Error("Failed to cleanup click logs")
	} else {
		w.logger.WithField("deleted_count", deletedClicks).Info("Cleaned up old click logs")
	}

	// Очищаем старую аналитику
	deletedAnalytics, err := w.cleanupAnalytics(ctx, analyticsRetention)
	if err != nil {
		w.logger.WithError(err).Error("Failed to cleanup analytics")
	} else {
		w.logger.WithField("deleted_count", deletedAnalytics).Info("Cleaned up old analytics")
	}

	duration := time.Since(start)
	w.logger.WithFields(logrus.Fields{
		"duration":          duration,
		"deleted_logs":      deletedLogs,
		"deleted_clicks":    deletedClicks,
		"deleted_analytics": deletedAnalytics,
	}).Info("Old logs cleanup completed")
}

// cleanupSearchLogs удаляет старые логи поиска
func (w *Worker) cleanupSearchLogs(ctx context.Context, retention time.Duration) (int64, error) {
	// Реализация будет добавлена в репозиторий
	// Здесь заглушка для демонстрации
	return 0, nil
}

// cleanupClickLogs удаляет старые логи кликов
func (w *Worker) cleanupClickLogs(ctx context.Context, retention time.Duration) (int64, error) {
	// Реализация будет добавлена в репозиторий
	// Здесь заглушка для демонстрации
	return 0, nil
}

// cleanupAnalytics удаляет старую аналитику
func (w *Worker) cleanupAnalytics(ctx context.Context, retention time.Duration) (int64, error) {
	// Реализация будет добавлена в репозиторий
	// Здесь заглушка для демонстрации
	return 0, nil
}

// updatePerformanceMetrics обновляет метрики производительности поиска
func (w *Worker) updatePerformanceMetrics(ctx context.Context) {
	// Вычисляем средние показатели производительности за последние 24 часа
	metrics, err := w.calculatePerformanceMetrics(ctx)
	if err != nil {
		w.logger.WithError(err).Error("Failed to calculate performance metrics")
		return
	}

	// Логируем важные метрики
	if metrics != nil {
		w.logger.WithFields(logrus.Fields{
			"avg_response_time_ms": metrics.AvgResponseTime,
			"p95_response_time_ms": metrics.P95ResponseTime,
			"p99_response_time_ms": metrics.P99ResponseTime,
			"zero_results_rate":    metrics.ZeroResultsRate,
			"error_rate":           metrics.ErrorRate,
		}).Info("Search performance metrics updated")

		// Проверяем пороги и отправляем алерты при необходимости
		w.checkPerformanceThresholds(metrics)
	}
}

// PerformanceMetrics представляет метрики производительности поиска
type PerformanceMetrics struct {
	AvgResponseTime float64
	P95ResponseTime float64
	P99ResponseTime float64
	ZeroResultsRate float64
	ErrorRate       float64
	TotalSearches   int64
}

// calculatePerformanceMetrics вычисляет метрики производительности
func (w *Worker) calculatePerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error) {
	// Реализация будет добавлена в репозиторий
	// Здесь заглушка для демонстрации
	return &PerformanceMetrics{
		AvgResponseTime: 150.5,
		P95ResponseTime: 450.0,
		P99ResponseTime: 850.0,
		ZeroResultsRate: 0.05,
		ErrorRate:       0.001,
		TotalSearches:   10000,
	}, nil
}

// checkPerformanceThresholds проверяет пороговые значения производительности
func (w *Worker) checkPerformanceThresholds(metrics *PerformanceMetrics) {
	// Пороговые значения
	const (
		maxAvgResponseTime = 500.0  // мс
		maxP99ResponseTime = 2000.0 // мс
		maxZeroResultsRate = 0.15   // 15%
		maxErrorRate       = 0.01   // 1%
	)

	// Проверяем среднее время ответа
	if metrics.AvgResponseTime > maxAvgResponseTime {
		w.logger.WithFields(logrus.Fields{
			"avg_response_time": metrics.AvgResponseTime,
			"threshold":         maxAvgResponseTime,
		}).Warn("Average search response time exceeds threshold")
	}

	// Проверяем P99 время ответа
	if metrics.P99ResponseTime > maxP99ResponseTime {
		w.logger.WithFields(logrus.Fields{
			"p99_response_time": metrics.P99ResponseTime,
			"threshold":         maxP99ResponseTime,
		}).Warn("P99 search response time exceeds threshold")
	}

	// Проверяем процент нулевых результатов
	if metrics.ZeroResultsRate > maxZeroResultsRate {
		w.logger.WithFields(logrus.Fields{
			"zero_results_rate": metrics.ZeroResultsRate,
			"threshold":         maxZeroResultsRate,
		}).Warn("Zero results rate exceeds threshold")
	}

	// Проверяем процент ошибок
	if metrics.ErrorRate > maxErrorRate {
		w.logger.WithFields(logrus.Fields{
			"error_rate": metrics.ErrorRate,
			"threshold":  maxErrorRate,
		}).Error("Search error rate exceeds threshold")
	}
}
