// Package postgres implements PostgreSQL storage for testing module
// backend/internal/proj/admin/testing/storage/postgres/postgres.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"

	"backend/internal/proj/admin/testing/domain"
)

// Sentinel errors
var (
	ErrTestRunNotFound = errors.New("test run not found")
)

// Data type constants
const (
	DataTypeTestRuns     = "test_runs"
	DataTypePriceHistory = "price_history"
)

// Storage implements TestStorage interface using PostgreSQL
type Storage struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewStorage creates new PostgreSQL storage instance
func NewStorage(db *sqlx.DB, logger zerolog.Logger) *Storage {
	return &Storage{
		db:     db,
		logger: logger.With().Str("component", "testing_storage").Logger(),
	}
}

// CreateTestRun creates new test run record
func (s *Storage) CreateTestRun(ctx context.Context, testRun *domain.TestRun) error {
	query := `
		INSERT INTO test_runs (
			run_uuid, test_suite, status, started_by_user_id,
			started_at, total_tests, passed_tests, failed_tests,
			skipped_tests, metadata, created_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3,
			$4, $5, $6, $7,
			$8, $9, $10
		) RETURNING id, run_uuid, created_at`

	var metadataJSON interface{}
	var err error
	if testRun.Metadata != nil {
		metadataBytes, marshalErr := json.Marshal(testRun.Metadata)
		if marshalErr != nil {
			s.logger.Error().Err(marshalErr).Msg("Failed to marshal metadata")
			return marshalErr
		}
		metadataJSON = metadataBytes
	} else {
		metadataJSON = nil // Pass NULL to database
	}

	err = s.db.QueryRowxContext(
		ctx, query,
		testRun.TestSuite,
		testRun.Status,
		testRun.StartedByUserID,
		testRun.StartedAt,
		testRun.TotalTests,
		testRun.PassedTests,
		testRun.FailedTests,
		testRun.SkippedTests,
		metadataJSON,
		time.Now(),
	).Scan(&testRun.ID, &testRun.RunUUID, &testRun.CreatedAt)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create test run")
		return err
	}

	s.logger.Info().
		Int64("id", testRun.ID).
		Str("uuid", testRun.RunUUID).
		Str("suite", testRun.TestSuite).
		Msg("Test run created")

	return nil
}

// GetTestRunByID retrieves test run by ID
func (s *Storage) GetTestRunByID(ctx context.Context, id int64) (*domain.TestRun, error) {
	query := `
		SELECT
			id, run_uuid, test_suite, status, started_by_user_id,
			started_at, completed_at, duration_ms, total_tests,
			passed_tests, failed_tests, skipped_tests, metadata, created_at
		FROM test_runs
		WHERE id = $1`

	var testRun domain.TestRun
	var metadataJSON []byte

	err := s.db.QueryRowxContext(ctx, query, id).Scan(
		&testRun.ID,
		&testRun.RunUUID,
		&testRun.TestSuite,
		&testRun.Status,
		&testRun.StartedByUserID,
		&testRun.StartedAt,
		&testRun.CompletedAt,
		&testRun.DurationMs,
		&testRun.TotalTests,
		&testRun.PassedTests,
		&testRun.FailedTests,
		&testRun.SkippedTests,
		&metadataJSON,
		&testRun.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTestRunNotFound
	}
	if err != nil {
		s.logger.Error().Err(err).Int64("id", id).Msg("Failed to get test run")
		return nil, err
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &testRun.Metadata); err != nil {
			s.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
		}
	}

	return &testRun, nil
}

// GetTestRunByUUID retrieves test run by UUID
func (s *Storage) GetTestRunByUUID(ctx context.Context, uuid string) (*domain.TestRun, error) {
	query := `
		SELECT
			id, run_uuid, test_suite, status, started_by_user_id,
			started_at, completed_at, duration_ms, total_tests,
			passed_tests, failed_tests, skipped_tests, metadata, created_at
		FROM test_runs
		WHERE run_uuid = $1`

	var testRun domain.TestRun
	var metadataJSON []byte

	err := s.db.QueryRowxContext(ctx, query, uuid).Scan(
		&testRun.ID,
		&testRun.RunUUID,
		&testRun.TestSuite,
		&testRun.Status,
		&testRun.StartedByUserID,
		&testRun.StartedAt,
		&testRun.CompletedAt,
		&testRun.DurationMs,
		&testRun.TotalTests,
		&testRun.PassedTests,
		&testRun.FailedTests,
		&testRun.SkippedTests,
		&metadataJSON,
		&testRun.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTestRunNotFound
	}
	if err != nil {
		s.logger.Error().Err(err).Str("uuid", uuid).Msg("Failed to get test run")
		return nil, err
	}

	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &testRun.Metadata); err != nil {
			s.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
		}
	}

	return &testRun, nil
}

// UpdateTestRunStatus updates test run status
func (s *Storage) UpdateTestRunStatus(ctx context.Context, id int64, status domain.TestRunStatus) error {
	query := `UPDATE test_runs SET status = $1 WHERE id = $2`

	_, err := s.db.ExecContext(ctx, query, status, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", id).Str("status", string(status)).Msg("Failed to update test run status")
		return err
	}

	return nil
}

// UpdateTestRunCompletion updates test run completion time and duration
func (s *Storage) UpdateTestRunCompletion(ctx context.Context, id int64, completedAt *time.Time, durationMs int) error {
	query := `UPDATE test_runs SET completed_at = $1, duration_ms = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, completedAt, durationMs, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", id).Msg("Failed to update test run completion")
		return err
	}

	return nil
}

// UpdateTestRunStats updates test run statistics
func (s *Storage) UpdateTestRunStats(ctx context.Context, id int64, total, passed, failed, skipped int) error {
	query := `
		UPDATE test_runs
		SET total_tests = $1, passed_tests = $2, failed_tests = $3, skipped_tests = $4
		WHERE id = $5`

	_, err := s.db.ExecContext(ctx, query, total, passed, failed, skipped, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", id).Msg("Failed to update test run stats")
		return err
	}

	return nil
}

// ListTestRuns retrieves list of test runs with pagination
func (s *Storage) ListTestRuns(ctx context.Context, limit, offset int) ([]*domain.TestRun, error) {
	query := `
		SELECT
			id, run_uuid, test_suite, status, started_by_user_id,
			started_at, completed_at, duration_ms, total_tests,
			passed_tests, failed_tests, skipped_tests, metadata, created_at
		FROM test_runs
		ORDER BY started_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := s.db.QueryxContext(ctx, query, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list test runs")
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var testRuns []*domain.TestRun
	for rows.Next() {
		var testRun domain.TestRun
		var metadataJSON []byte

		err := rows.Scan(
			&testRun.ID,
			&testRun.RunUUID,
			&testRun.TestSuite,
			&testRun.Status,
			&testRun.StartedByUserID,
			&testRun.StartedAt,
			&testRun.CompletedAt,
			&testRun.DurationMs,
			&testRun.TotalTests,
			&testRun.PassedTests,
			&testRun.FailedTests,
			&testRun.SkippedTests,
			&metadataJSON,
			&testRun.CreatedAt,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to scan test run")
			continue
		}

		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &testRun.Metadata); err != nil {
				s.logger.Error().Err(err).Msg("Failed to unmarshal metadata")
			}
		}

		testRuns = append(testRuns, &testRun)
	}

	return testRuns, nil
}

// DeleteTestRun deletes test run and related data
func (s *Storage) DeleteTestRun(ctx context.Context, id int64) error {
	query := `DELETE FROM test_runs WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", id).Msg("Failed to delete test run")
		return err
	}

	s.logger.Info().Int64("id", id).Msg("Test run deleted")
	return nil
}

// CreateTestResult creates new test result record
func (s *Storage) CreateTestResult(ctx context.Context, result *domain.TestResult) error {
	query := `
		INSERT INTO test_results (
			test_run_id, test_name, test_suite, status,
			duration_ms, error_msg, stack_trace,
			started_at, completed_at, created_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10
		) RETURNING id, created_at`

	err := s.db.QueryRowxContext(
		ctx, query,
		result.TestRunID,
		result.TestName,
		result.TestSuite,
		result.Status,
		result.DurationMs,
		result.ErrorMsg,
		result.StackTrace,
		result.StartedAt,
		result.CompletedAt,
		time.Now(),
	).Scan(&result.ID, &result.CreatedAt)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create test result")
		return err
	}

	return nil
}

// GetTestResultsByRunID retrieves all test results for a run
func (s *Storage) GetTestResultsByRunID(ctx context.Context, runID int64) ([]*domain.TestResult, error) {
	query := `
		SELECT
			id, test_run_id, test_name, test_suite, status,
			duration_ms, error_msg, stack_trace,
			started_at, completed_at, created_at
		FROM test_results
		WHERE test_run_id = $1
		ORDER BY started_at ASC`

	rows, err := s.db.QueryxContext(ctx, query, runID)
	if err != nil {
		s.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to get test results")
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []*domain.TestResult
	for rows.Next() {
		var result domain.TestResult
		err := rows.Scan(
			&result.ID,
			&result.TestRunID,
			&result.TestName,
			&result.TestSuite,
			&result.Status,
			&result.DurationMs,
			&result.ErrorMsg,
			&result.StackTrace,
			&result.StartedAt,
			&result.CompletedAt,
			&result.CreatedAt,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to scan test result")
			continue
		}
		results = append(results, &result)
	}

	return results, nil
}

// CreateTestLog creates new test log entry
func (s *Storage) CreateTestLog(ctx context.Context, log *domain.TestLog) error {
	query := `
		INSERT INTO test_logs (
			test_run_id, level, message, timestamp, created_at
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, created_at`

	err := s.db.QueryRowxContext(
		ctx, query,
		log.TestRunID,
		log.Level,
		log.Message,
		log.Timestamp,
		time.Now(),
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		// Don't log this error to avoid infinite loop
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			s.logger.Debug().Str("code", string(pqErr.Code)).Msg("Database error creating test log")
		}
		return err
	}

	return nil
}

// GetTestLogsByRunID retrieves test logs for a run
func (s *Storage) GetTestLogsByRunID(ctx context.Context, runID int64, limit int) ([]*domain.TestLog, error) {
	query := `
		SELECT
			id, test_run_id, level, message, timestamp, created_at
		FROM test_logs
		WHERE test_run_id = $1
		ORDER BY timestamp DESC
		LIMIT $2`

	rows, err := s.db.QueryxContext(ctx, query, runID, limit)
	if err != nil {
		s.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to get test logs")
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var logs []*domain.TestLog
	for rows.Next() {
		var log domain.TestLog
		err := rows.Scan(
			&log.ID,
			&log.TestRunID,
			&log.Level,
			&log.Message,
			&log.Timestamp,
			&log.CreatedAt,
		)
		if err != nil {
			s.logger.Error().Err(err).Msg("Failed to scan test log")
			continue
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

// GetTestDataStats retrieves statistics about test data in database
func (s *Storage) GetTestDataStats(ctx context.Context) (*domain.TestDataStats, error) {
	stats := &domain.TestDataStats{}
	stats.CollectedAt = time.Now()

	// Get test_runs stats
	err := s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('test_runs')),
			pg_total_relation_size('test_runs')
		FROM test_runs
	`).Scan(&stats.TestRuns.Count, &stats.TestRuns.SizeMB, &stats.TestRuns.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get test_runs stats: %w", err)
	}

	// Get test_logs stats
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('test_logs')),
			pg_total_relation_size('test_logs')
		FROM test_logs
	`).Scan(&stats.TestLogs.Count, &stats.TestLogs.SizeMB, &stats.TestLogs.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get test_logs stats: %w", err)
	}

	// Get test_results stats
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('test_results')),
			pg_total_relation_size('test_results')
		FROM test_results
	`).Scan(&stats.TestResults.Count, &stats.TestResults.SizeMB, &stats.TestResults.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get test_results stats: %w", err)
	}

	// Get user_behavior_events stats
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('user_behavior_events')),
			pg_total_relation_size('user_behavior_events')
		FROM user_behavior_events
	`).Scan(&stats.BehaviorEvents.Count, &stats.BehaviorEvents.SizeMB, &stats.BehaviorEvents.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_behavior_events stats: %w", err)
	}

	// Get behavior events by type
	rows, err := s.db.QueryContext(ctx, `
		SELECT event_type, COUNT(*)
		FROM user_behavior_events
		GROUP BY event_type
		ORDER BY COUNT(*) DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get behavior events by type: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			s.logger.Error().Err(closeErr).Msg("Failed to close rows")
		}
	}()

	stats.BehaviorEvents.ByType = make(map[string]int64)
	for rows.Next() {
		var eventType string
		var count int64
		if err := rows.Scan(&eventType, &count); err != nil {
			continue
		}
		stats.BehaviorEvents.ByType[eventType] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate behavior events: %w", err)
	}

	// Get category_detection_feedback stats
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('category_detection_feedback')),
			pg_total_relation_size('category_detection_feedback')
		FROM category_detection_feedback
	`).Scan(&stats.CategoryFeedback.Count, &stats.CategoryFeedback.SizeMB, &stats.CategoryFeedback.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get category_detection_feedback stats: %w", err)
	}

	// Get price_history stats
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('price_history')),
			pg_total_relation_size('price_history')
		FROM price_history
	`).Scan(&stats.PriceHistory.Count, &stats.PriceHistory.SizeMB, &stats.PriceHistory.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get price_history stats: %w", err)
	}

	// Get ai_category_decisions stats
	err = s.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*),
			pg_size_pretty(pg_total_relation_size('ai_category_decisions')),
			pg_total_relation_size('ai_category_decisions')
		FROM ai_category_decisions
	`).Scan(&stats.AIDecisions.Count, &stats.AIDecisions.SizeMB, &stats.AIDecisions.SizeBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get ai_category_decisions stats: %w", err)
	}

	// Calculate total size
	stats.TotalSizeBytes = stats.TestRuns.SizeBytes +
		stats.TestLogs.SizeBytes +
		stats.TestResults.SizeBytes +
		stats.BehaviorEvents.SizeBytes +
		stats.CategoryFeedback.SizeBytes +
		stats.PriceHistory.SizeBytes +
		stats.AIDecisions.SizeBytes

	// Get database size
	err = s.db.QueryRowContext(ctx, `
		SELECT pg_size_pretty(pg_database_size(current_database())),
		       pg_size_pretty($1::bigint)
	`, stats.TotalSizeBytes).Scan(&stats.DatabaseSizeMB, &stats.TotalSizeMB)
	if err != nil {
		return nil, fmt.Errorf("failed to get database size: %w", err)
	}

	return stats, nil
}

// CleanupTestData removes test data from database
func (s *Storage) CleanupTestData(ctx context.Context, types []string) (map[string]int64, error) {
	deletedCount := make(map[string]int64)

	// If no types specified, clean all
	if len(types) == 0 {
		types = []string{DataTypeTestRuns, "logs", "results", "behavior", "feedback", DataTypePriceHistory, "ai_decisions"}
	}

	for _, dataType := range types {
		var err error

		switch dataType {
		case DataTypeTestRuns:
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE test_runs CASCADE")
		case "logs":
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE test_logs CASCADE")
		case "results":
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE test_results CASCADE")
		case "behavior":
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE user_behavior_events CASCADE")
		case "feedback":
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE category_detection_feedback CASCADE")
		case DataTypePriceHistory:
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE price_history CASCADE")
		case "ai_decisions":
			_, err = s.db.ExecContext(ctx, "TRUNCATE TABLE ai_category_decisions CASCADE")
		default:
			s.logger.Warn().Str("type", dataType).Msg("Unknown data type for cleanup")
			continue
		}

		if err != nil {
			s.logger.Error().Err(err).Str("type", dataType).Msg("Failed to cleanup data")
			return nil, fmt.Errorf("failed to cleanup %s: %w", dataType, err)
		}

		// For TRUNCATE, RowsAffected always returns 0, so we'll just mark as cleaned
		deletedCount[dataType] = 1
	}

	// Run VACUUM to reclaim space
	for _, dataType := range types {
		var tableName string
		switch dataType {
		case DataTypeTestRuns:
			tableName = DataTypeTestRuns
		case "logs":
			tableName = "test_logs"
		case "results":
			tableName = "test_results"
		case "behavior":
			tableName = "user_behavior_events"
		case "feedback":
			tableName = "category_detection_feedback"
		case DataTypePriceHistory:
			tableName = DataTypePriceHistory
		case "ai_decisions":
			tableName = "ai_category_decisions"
		default:
			continue
		}

		// VACUUM FULL must be run outside transaction
		// We'll just log it for now, the DBA can run it manually
		s.logger.Info().Str("table", tableName).Msg("Table cleaned, consider running VACUUM FULL")
	}

	return deletedCount, nil
}
