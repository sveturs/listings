package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/metrics"
)

// Repository defines repository interface for worker
type Repository interface {
	GetPendingIndexingJobs(ctx context.Context, limit int) ([]*domain.IndexingQueueItem, error)
	GetListingByID(ctx context.Context, id int64) (*domain.Listing, error)
	CompleteIndexingJob(ctx context.Context, jobID int64) error
	FailIndexingJob(ctx context.Context, jobID int64, errorMsg string) error
}

// Indexer defines indexing service interface
type Indexer interface {
	IndexListing(ctx context.Context, listing *domain.Listing) error
	UpdateListing(ctx context.Context, listing *domain.Listing) error
	DeleteListing(ctx context.Context, listingID int64) error
}

// Worker handles async background indexing jobs
type Worker struct {
	repo         Repository
	indexer      Indexer
	metrics      *metrics.Metrics
	concurrency  int
	pollInterval time.Duration
	logger       zerolog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewWorker creates a new background worker
func NewWorker(repo Repository, indexer Indexer, metrics *metrics.Metrics, concurrency int, logger zerolog.Logger) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		repo:         repo,
		indexer:      indexer,
		metrics:      metrics,
		concurrency:  concurrency,
		pollInterval: 2 * time.Second,
		logger:       logger.With().Str("component", "indexing_worker").Logger(),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins processing background indexing jobs
func (w *Worker) Start() error {
	w.logger.Info().Int("concurrency", w.concurrency).Msg("starting indexing worker")

	// Start worker goroutines
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)
		go w.workerLoop(i)
	}

	w.logger.Info().Msg("indexing worker started")
	return nil
}

// Stop gracefully shuts down the worker
func (w *Worker) Stop() error {
	w.logger.Info().Msg("stopping indexing worker")

	w.cancel()
	w.wg.Wait()

	w.logger.Info().Msg("indexing worker stopped")
	return nil
}

// workerLoop is the main worker processing loop
func (w *Worker) workerLoop(workerID int) {
	defer w.wg.Done()

	logger := w.logger.With().Int("worker_id", workerID).Logger()
	logger.Debug().Msg("worker started")

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			logger.Debug().Msg("worker shutting down")
			return

		case <-ticker.C:
			w.processBatch(logger)
		}
	}
}

// processBatch fetches and processes a batch of indexing jobs
func (w *Worker) processBatch(logger zerolog.Logger) {
	ctx, cancel := context.WithTimeout(w.ctx, 30*time.Second)
	defer cancel()

	// Fetch pending jobs (limit 10 per worker iteration)
	jobs, err := w.repo.GetPendingIndexingJobs(ctx, 10)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch pending indexing jobs")
		w.metrics.RecordError("worker", "fetch_jobs_failed")
		return
	}

	if len(jobs) == 0 {
		return // No jobs to process
	}

	logger.Debug().Int("count", len(jobs)).Msg("processing indexing jobs")

	// Process each job
	for _, job := range jobs {
		w.processJob(ctx, job, logger)
	}
}

// processJob processes a single indexing job
func (w *Worker) processJob(ctx context.Context, job *domain.IndexingQueueItem, logger zerolog.Logger) {
	start := time.Now()

	logger = logger.With().
		Int64("job_id", job.ID).
		Int64("listing_id", job.ListingID).
		Str("operation", job.Operation).
		Logger()

	logger.Debug().Msg("processing indexing job")

	var err error

	switch job.Operation {
	case domain.IndexOpIndex, domain.IndexOpUpdate:
		err = w.handleIndexJob(ctx, job)
	case domain.IndexOpDelete:
		err = w.handleDeleteJob(ctx, job)
	default:
		err = fmt.Errorf("unknown operation: %s", job.Operation)
	}

	duration := time.Since(start).Seconds()

	if err != nil {
		logger.Error().Err(err).Msg("failed to process indexing job")

		// Mark job as failed
		if failErr := w.repo.FailIndexingJob(ctx, job.ID, err.Error()); failErr != nil {
			logger.Error().Err(failErr).Msg("failed to mark job as failed")
		}

		w.metrics.RecordIndexingJob(job.Operation, "failed", duration)
		w.metrics.RecordError("worker", "job_failed")
	} else {
		logger.Debug().Msg("indexing job completed successfully")

		// Mark job as completed
		if completeErr := w.repo.CompleteIndexingJob(ctx, job.ID); completeErr != nil {
			logger.Error().Err(completeErr).Msg("failed to mark job as completed")
		}

		w.metrics.RecordIndexingJob(job.Operation, "success", duration)
	}
}

// handleIndexJob handles index/update operations
func (w *Worker) handleIndexJob(ctx context.Context, job *domain.IndexingQueueItem) error {
	// Fetch listing from database
	listing, err := w.repo.GetListingByID(ctx, job.ListingID)
	if err != nil {
		return fmt.Errorf("failed to fetch listing: %w", err)
	}

	// Index listing in OpenSearch
	if job.Operation == domain.IndexOpIndex {
		if err := w.indexer.IndexListing(ctx, listing); err != nil {
			return fmt.Errorf("failed to index listing: %w", err)
		}
	} else {
		if err := w.indexer.UpdateListing(ctx, listing); err != nil {
			return fmt.Errorf("failed to update listing: %w", err)
		}
	}

	return nil
}

// handleDeleteJob handles delete operations
func (w *Worker) handleDeleteJob(ctx context.Context, job *domain.IndexingQueueItem) error {
	// Delete listing from OpenSearch
	if err := w.indexer.DeleteListing(ctx, job.ListingID); err != nil {
		return fmt.Errorf("failed to delete listing from index: %w", err)
	}

	return nil
}
