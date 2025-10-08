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

	// Check if job was canceled
	if job.Status == statusCanceled {
		w.logger.Info().
			Int("job_id", task.JobID).
			Msg("Job was canceled, skipping")
		return
	}

	// Update job status to processing
	job.Status = statusProcessing
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

	// Import products with variant detection
	w.logger.Info().
		Int("job_id", task.JobID).
		Int("total_products", len(products)).
		Msg("Starting product import with variant detection")

	// Шаг 1: Группировка товаров в варианты
	variantGroups := w.service.groupAndDetectVariants(products)

	w.logger.Info().
		Int("job_id", task.JobID).
		Int("total_groups", len(variantGroups)).
		Msg("Variant detection completed")

	// Шаг 2: Импорт variant groups
	successCount := 0
	var importErrors []error

	// Создаем map для быстрой проверки какие товары вошли в группы
	productsInGroups := make(map[string]bool)

	for _, group := range variantGroups {
		// Проверяем что группа валидна (минимум 2 варианта, confidence > 0.5)
		if len(group.Variants) < 2 || group.Confidence < 0.5 {
			// Не группируем - импортируем как отдельные товары
			continue
		}

		w.logger.Info().
			Int("job_id", task.JobID).
			Str("base_name", group.BaseName).
			Int("variants_count", len(group.Variants)).
			Float64("confidence", group.Confidence).
			Msg("Importing variant group")

		// Импортируем группу вариантов
		if err := w.service.importVariantGroup(w.ctx, group, task.StorefrontID); err != nil {
			w.logger.Error().
				Err(err).
				Int("job_id", task.JobID).
				Str("base_name", group.BaseName).
				Msg("Failed to import variant group")
			importErrors = append(importErrors, err)
		} else {
			// Помечаем все товары группы как успешно импортированные
			for _, variant := range group.Variants {
				productsInGroups[variant.SKU] = true
				successCount++
			}
		}
	}

	// Шаг 3: Импорт оставшихся товаров (не вошедших в группы)
	remainingProducts := make([]models.ImportProductRequest, 0)
	for _, product := range products {
		if !productsInGroups[product.SKU] {
			remainingProducts = append(remainingProducts, product)
		}
	}

	if len(remainingProducts) > 0 {
		w.logger.Info().
			Int("job_id", task.JobID).
			Int("remaining_products", len(remainingProducts)).
			Msg("Importing remaining products (not grouped)")

		remainingSuccess, remainingErrors := w.service.importProducts(w.ctx, job.ID, remainingProducts, task.StorefrontID, task.UpdateMode)
		successCount += remainingSuccess
		importErrors = append(importErrors, remainingErrors...)
	}

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

	// Инкрементальная индексация товаров с needs_reindex=true
	if successCount > 0 {
		w.logger.Info().
			Int("job_id", task.JobID).
			Int("storefront_id", task.StorefrontID).
			Msg("Starting incremental indexing of imported products...")

		// Индексируем marketplace listings (созданные триггером из storefront_products)
		if err := w.service.IndexMarketplaceListings(w.ctx, 100); err != nil {
			w.logger.Error().
				Err(err).
				Int("job_id", task.JobID).
				Msg("Failed to index marketplace listings (non-fatal)")
			// Не прерываем выполнение - товары импортированы, индексация может быть выполнена позже
		} else {
			w.logger.Info().
				Int("job_id", task.JobID).
				Msg("Successfully indexed marketplace listings")
		}
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
