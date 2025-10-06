package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/domain/models"
)

// ImportJobFilter represents filter options for listing import jobs
type ImportJobFilter struct {
	Status    *string
	FileType  *string
	UserID    *int
	Limit     int
	Offset    int
	SortBy    string // created_at, updated_at, status
	SortOrder string // asc, desc
}

// ImportJobsRepositoryInterface определяет интерфейс для работы с import jobs
type ImportJobsRepositoryInterface interface {
	Create(ctx context.Context, job *models.ImportJob) error
	GetByID(ctx context.Context, id int) (*models.ImportJob, error)
	GetByStorefront(ctx context.Context, storefrontID int, filter *ImportJobFilter) ([]*models.ImportJob, int, error)
	Update(ctx context.Context, job *models.ImportJob) error
	UpdateStatus(ctx context.Context, id int, status string) error
	AddError(ctx context.Context, err *models.ImportError) error
	GetErrors(ctx context.Context, jobID int) ([]*models.ImportError, error)
}

// importJobsRepository реализует интерфейс для работы с import jobs
type importJobsRepository struct {
	pool *pgxpool.Pool
}

// NewImportJobsRepository создает новый репозиторий import jobs
func NewImportJobsRepository(pool *pgxpool.Pool) ImportJobsRepositoryInterface {
	return &importJobsRepository{pool: pool}
}

// Create создает новую запись import job
func (r *importJobsRepository) Create(ctx context.Context, job *models.ImportJob) error {
	query := `
		INSERT INTO import_jobs (
			storefront_id, user_id, file_name, file_type, file_url,
			status, total_records, processed_records, successful_records, failed_records,
			error_message, started_at, completed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(
		ctx, query,
		job.StorefrontID,
		job.UserID,
		job.FileName,
		job.FileType,
		job.FileURL,
		job.Status,
		job.TotalRecords,
		job.ProcessedRecords,
		job.SuccessfulRecords,
		job.FailedRecords,
		job.ErrorMessage,
		job.StartedAt,
		job.CompletedAt,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create import job: %w", err)
	}

	return nil
}

// GetByID получает import job по ID
func (r *importJobsRepository) GetByID(ctx context.Context, id int) (*models.ImportJob, error) {
	query := `
		SELECT
			id, storefront_id, user_id, file_name, file_type, file_url,
			status, total_records, processed_records, successful_records, failed_records,
			error_message, started_at, completed_at, created_at, updated_at
		FROM import_jobs
		WHERE id = $1`

	var job models.ImportJob
	var fileURL, errorMessage sql.NullString
	var startedAt, completedAt sql.NullTime

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.StorefrontID,
		&job.UserID,
		&job.FileName,
		&job.FileType,
		&fileURL,
		&job.Status,
		&job.TotalRecords,
		&job.ProcessedRecords,
		&job.SuccessfulRecords,
		&job.FailedRecords,
		&errorMessage,
		&startedAt,
		&completedAt,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("import job not found")
		}
		return nil, fmt.Errorf("failed to get import job: %w", err)
	}

	// Обработка NULL значений
	if fileURL.Valid {
		job.FileURL = &fileURL.String
	}
	if errorMessage.Valid {
		job.ErrorMessage = &errorMessage.String
	}
	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}

	return &job, nil
}

// GetByStorefront получает список import jobs для витрины с фильтрацией
func (r *importJobsRepository) GetByStorefront(ctx context.Context, storefrontID int, filter *ImportJobFilter) ([]*models.ImportJob, int, error) {
	// Базовый запрос
	baseQuery := `
		FROM import_jobs
		WHERE storefront_id = $1`

	// Счетчик параметров
	paramCount := 1
	args := []interface{}{storefrontID}
	conditions := []string{}

	// Применяем фильтры
	if filter != nil {
		if filter.Status != nil {
			paramCount++
			conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
			args = append(args, *filter.Status)
		}

		if filter.FileType != nil {
			paramCount++
			conditions = append(conditions, fmt.Sprintf("file_type = $%d", paramCount))
			args = append(args, *filter.FileType)
		}

		if filter.UserID != nil {
			paramCount++
			conditions = append(conditions, fmt.Sprintf("user_id = $%d", paramCount))
			args = append(args, *filter.UserID)
		}
	}

	// Добавляем дополнительные условия к базовому запросу
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Получаем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count import jobs: %w", err)
	}

	// Формируем запрос на получение данных
	selectQuery := `
		SELECT
			id, storefront_id, user_id, file_name, file_type, file_url,
			status, total_records, processed_records, successful_records, failed_records,
			error_message, started_at, completed_at, created_at, updated_at
		` + baseQuery

	// Сортировка
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter != nil {
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = strings.ToUpper(filter.SortOrder)
		}
	}
	selectQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Лимит и offset
	if filter != nil && filter.Limit > 0 {
		paramCount++
		selectQuery += fmt.Sprintf(" LIMIT $%d", paramCount)
		args = append(args, filter.Limit)

		if filter.Offset > 0 {
			paramCount++
			selectQuery += fmt.Sprintf(" OFFSET $%d", paramCount)
			args = append(args, filter.Offset)
		}
	}

	// Выполняем запрос
	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query import jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*models.ImportJob
	for rows.Next() {
		var job models.ImportJob
		var fileURL, errorMessage sql.NullString
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(
			&job.ID,
			&job.StorefrontID,
			&job.UserID,
			&job.FileName,
			&job.FileType,
			&fileURL,
			&job.Status,
			&job.TotalRecords,
			&job.ProcessedRecords,
			&job.SuccessfulRecords,
			&job.FailedRecords,
			&errorMessage,
			&startedAt,
			&completedAt,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan import job: %w", err)
		}

		// Обработка NULL значений
		if fileURL.Valid {
			job.FileURL = &fileURL.String
		}
		if errorMessage.Valid {
			job.ErrorMessage = &errorMessage.String
		}
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}

		jobs = append(jobs, &job)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating import jobs: %w", err)
	}

	return jobs, total, nil
}

// Update обновляет существующий import job
func (r *importJobsRepository) Update(ctx context.Context, job *models.ImportJob) error {
	query := `
		UPDATE import_jobs
		SET
			status = $2,
			total_records = $3,
			processed_records = $4,
			successful_records = $5,
			failed_records = $6,
			error_message = $7,
			started_at = $8,
			completed_at = $9
		WHERE id = $1`

	result, err := r.pool.Exec(
		ctx, query,
		job.ID,
		job.Status,
		job.TotalRecords,
		job.ProcessedRecords,
		job.SuccessfulRecords,
		job.FailedRecords,
		job.ErrorMessage,
		job.StartedAt,
		job.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update import job: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("import job not found")
	}

	return nil
}

// UpdateStatus обновляет только статус import job
func (r *importJobsRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE import_jobs
		SET status = $2
		WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update import job status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("import job not found")
	}

	return nil
}

// AddError добавляет ошибку импорта
func (r *importJobsRepository) AddError(ctx context.Context, importError *models.ImportError) error {
	query := `
		INSERT INTO import_errors (
			job_id, line_number, field_name, error_message, raw_data
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id, created_at`

	err := r.pool.QueryRow(
		ctx, query,
		importError.JobID,
		importError.LineNumber,
		importError.FieldName,
		importError.ErrorMessage,
		importError.RawData,
	).Scan(&importError.ID, &importError.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add import error: %w", err)
	}

	return nil
}

// GetErrors получает список ошибок для import job
func (r *importJobsRepository) GetErrors(ctx context.Context, jobID int) ([]*models.ImportError, error) {
	query := `
		SELECT
			id, job_id, line_number, field_name, error_message, raw_data, created_at
		FROM import_errors
		WHERE job_id = $1
		ORDER BY line_number ASC, created_at ASC`

	rows, err := r.pool.Query(ctx, query, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to query import errors: %w", err)
	}
	defer rows.Close()

	var errors []*models.ImportError
	for rows.Next() {
		var importError models.ImportError
		err := rows.Scan(
			&importError.ID,
			&importError.JobID,
			&importError.LineNumber,
			&importError.FieldName,
			&importError.ErrorMessage,
			&importError.RawData,
			&importError.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan import error: %w", err)
		}

		errors = append(errors, &importError)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating import errors: %w", err)
	}

	return errors, nil
}
