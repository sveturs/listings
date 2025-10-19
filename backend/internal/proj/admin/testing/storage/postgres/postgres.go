// Package postgres implements PostgreSQL storage for testing module
// backend/internal/proj/admin/testing/storage/postgres/postgres.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

	if err == sql.ErrNoRows {
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

	if err == sql.ErrNoRows {
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
	defer rows.Close()

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
	defer rows.Close()

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
		if pqErr, ok := err.(*pq.Error); ok {
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
	defer rows.Close()

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
