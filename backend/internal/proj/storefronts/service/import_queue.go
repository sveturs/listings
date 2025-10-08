package service

import (
	"context"
	"fmt"
	"sync"

	"backend/internal/logger"

	"github.com/rs/zerolog"
)

// ImportQueueManager manages the import job queue and worker pool
type ImportQueueManager struct {
	jobQueue    chan *ImportJobTask
	workers     []*ImportWorker
	workerCount int
	quit        chan bool
	wg          sync.WaitGroup
	service     *ImportService
	ctx         context.Context
	cancel      context.CancelFunc
	logger      *zerolog.Logger
	mu          sync.RWMutex
	isRunning   bool
}

// NewImportQueueManager creates a new import queue manager
func NewImportQueueManager(workerCount int, queueSize int, service *ImportService) *ImportQueueManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &ImportQueueManager{
		jobQueue:    make(chan *ImportJobTask, queueSize),
		workers:     make([]*ImportWorker, 0, workerCount),
		workerCount: workerCount,
		quit:        make(chan bool),
		wg:          sync.WaitGroup{},
		service:     service,
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger.Get(),
		isRunning:   false,
	}
}

// Start starts the queue manager and all workers
func (m *ImportQueueManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("queue manager is already running")
	}

	m.logger.Info().
		Int("worker_count", m.workerCount).
		Int("queue_size", cap(m.jobQueue)).
		Msg("Starting import queue manager")

	// Create and start workers
	for i := 0; i < m.workerCount; i++ {
		worker := NewImportWorker(
			i+1,
			m.jobQueue,
			m.quit,
			&m.wg,
			m.service,
			m.ctx,
		)
		worker.Start()
		m.workers = append(m.workers, worker)
	}

	m.isRunning = true

	m.logger.Info().
		Int("worker_count", len(m.workers)).
		Msg("Import queue manager started successfully")

	return nil
}

// Stop stops the queue manager and all workers
func (m *ImportQueueManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("queue manager is not running")
	}

	m.logger.Info().Msg("Stopping import queue manager")

	// Signal all workers to stop
	close(m.quit)

	// Cancel context
	m.cancel()

	// Wait for all workers to finish
	m.wg.Wait()

	// Close job queue
	close(m.jobQueue)

	m.isRunning = false

	m.logger.Info().Msg("Import queue manager stopped")

	return nil
}

// EnqueueJob adds a new import job to the queue
func (m *ImportQueueManager) EnqueueJob(task *ImportJobTask) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.isRunning {
		return fmt.Errorf("queue manager is not running")
	}

	select {
	case m.jobQueue <- task:
		m.logger.Info().
			Int("job_id", task.JobID).
			Int("storefront_id", task.StorefrontID).
			Str("file_type", task.FileType).
			Msg("Job enqueued successfully")
		return nil
	default:
		return fmt.Errorf("job queue is full, please try again later")
	}
}

// GetQueueStats returns statistics about the queue
func (m *ImportQueueManager) GetQueueStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"worker_count":    m.workerCount,
		"queue_size":      cap(m.jobQueue),
		"pending_jobs":    len(m.jobQueue),
		"available_slots": cap(m.jobQueue) - len(m.jobQueue),
		"is_running":      m.isRunning,
	}
}

// IsRunning returns whether the queue manager is running
func (m *ImportQueueManager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isRunning
}

// QueueLength returns the current number of jobs in the queue
func (m *ImportQueueManager) QueueLength() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.jobQueue)
}

// QueueCapacity returns the maximum capacity of the queue
func (m *ImportQueueManager) QueueCapacity() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return cap(m.jobQueue)
}
