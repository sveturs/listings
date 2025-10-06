package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"

	"github.com/rs/zerolog"
)

// ImportWorker represents a worker that processes import jobs
type ImportWorker struct {
	id       int
	jobQueue <-chan *ImportJobTask
	quit     chan bool
	wg       *sync.WaitGroup
	service  *ImportService
	ctx      context.Context
	logger   *zerolog.Logger
}

// ImportJobTask represents a task to be processed by a worker
type ImportJobTask struct {
	JobID        int
	UserID       int
	StorefrontID int
	FileData     []byte
	FileType     string
	UpdateMode   string
	FileName     *string
	FileURL      *string
}

// NewImportWorker creates a new import worker
func NewImportWorker(
	id int,
	jobQueue <-chan *ImportJobTask,
	quit chan bool,
	wg *sync.WaitGroup,
	service *ImportService,
	ctx context.Context,
) *ImportWorker {
	return &ImportWorker{
		id:       id,
		jobQueue: jobQueue,
		quit:     quit,
		wg:       wg,
		service:  service,
		ctx:      ctx,
		logger:   logger.Get(),
	}
}

// Start starts the worker
func (w *ImportWorker) Start() {
	w.logger.Info().
		Int("worker_id", w.id).
		Msg("Import worker started")

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		for {
			select {
			case task := <-w.jobQueue:
				if task != nil {
					w.processJob(task)
				}
			case <-w.quit:
				w.logger.Info().
					Int("worker_id", w.id).
					Msg("Import worker stopped")
				return
			case <-w.ctx.Done():
				w.logger.Info().
					Int("worker_id", w.id).
					Msg("Import worker stopped due to context cancellation")
				return
			}
		}
	}()
}

// processJob processes a single import job
func (w *ImportWorker) processJob(task *ImportJobTask) {
	w.logger.Info().
		Int("worker_id", w.id).
		Int("job_id", task.JobID).
		Int("storefront_id", task.StorefrontID).
		Msg("Processing import job")

	startTime := time.Now()

	// Get job from database
	job, err := w.service.jobsRepo.GetByID(w.ctx, task.JobID)
	if err != nil {
		w.logger.Error().
			Err(err).
			Int("job_id", task.JobID).
			Msg("Failed to get job from database")
		return
	}

	// Check if job was cancelled
	if job.Status == "cancelled" {
		w.logger.Info().
			Int("job_id", task.JobID).
			Msg("Job was cancelled, skipping")
		return
	}

	// Update job status to processing
	job.Status = "processing"
	job.StartedAt = &startTime
	if err := w.service.jobsRepo.Update(w.ctx, job); err != nil {
		w.logger.Error().
			Err(err).
			Int("job_id", task.JobID).
			Msg("Failed to update job status to processing")
		return
	}

	// Process file based on type
	var products []models.ImportProductRequest
	var validationErrors []models.ImportValidationError

	switch task.FileType {
	case "xml":
		products, validationErrors, err = w.service.processXMLData(task.FileData, task.StorefrontID)
	case "csv":
		products, validationErrors, err = w.service.processCSVData(task.FileData, task.StorefrontID)
	case "zip":
		products, validationErrors, err = w.service.processZIPData(task.FileData, task.StorefrontID)
	default:
		err = fmt.Errorf("unsupported file type: %s", task.FileType)
	}

	if err != nil {
		w.logger.Error().
			Err(err).
			Int("job_id", task.JobID).
			Str("file_type", task.FileType).
			Msg("Failed to process file")

		job.Status = "failed"
		errorMsg := err.Error()
		job.ErrorMessage = &errorMsg
		completedTime := time.Now()
		job.CompletedAt = &completedTime
		_ = w.service.jobsRepo.Update(w.ctx, job)
		return
	}

	// Update job with totals
	job.TotalRecords = len(products) + len(validationErrors)
	job.FailedRecords = len(validationErrors)

	// Save validation errors to import_errors table
	for i, valErr := range validationErrors {
		importError := &models.ImportError{
			JobID:        job.ID,
			LineNumber:   i + 1,
			FieldName:    valErr.Field,
			ErrorMessage: valErr.Message,
			RawData:      fmt.Sprintf("%v", valErr.Value),
		}
		_ = w.service.jobsRepo.AddError(w.ctx, importError)
	}

	// Import products
	w.logger.Info().
		Int("job_id", task.JobID).
		Int("total_products", len(products)).
		Msg("Starting product import")

	successCount, importErrors := w.service.importProducts(w.ctx, job.ID, products, task.StorefrontID, task.UpdateMode)

	job.ProcessedRecords = len(products)
	job.SuccessfulRecords = successCount
	job.FailedRecords += len(importErrors)

	// Complete job
	completedTime := time.Now()
	job.CompletedAt = &completedTime
	job.Status = "completed"

	if len(importErrors) > 0 {
		errorMsg := fmt.Sprintf("Import completed with %d errors", len(importErrors))
		job.ErrorMessage = &errorMsg
	}

	// Update job in database
	if err := w.service.jobsRepo.Update(w.ctx, job); err != nil {
		w.logger.Error().
			Err(err).
			Int("job_id", task.JobID).
			Msg("Failed to update job after completion")
		return
	}

	duration := time.Since(startTime)
	w.logger.Info().
		Int("worker_id", w.id).
		Int("job_id", task.JobID).
		Int("successful_records", successCount).
		Int("failed_records", job.FailedRecords).
		Dur("duration", duration).
		Msg("Import job completed")
}

// Stop stops the worker
func (w *ImportWorker) Stop() {
	w.logger.Info().
		Int("worker_id", w.id).
		Msg("Stopping import worker")
	w.quit <- true
}
