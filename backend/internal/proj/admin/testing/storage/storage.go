// Package storage defines storage interfaces for testing module
// backend/internal/proj/admin/testing/storage/storage.go
package storage

import (
	"context"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// TestStorage defines the interface for test data persistence
type TestStorage interface {
	// TestRuns
	CreateTestRun(ctx context.Context, testRun *domain.TestRun) error
	GetTestRunByID(ctx context.Context, id int64) (*domain.TestRun, error)
	GetTestRunByUUID(ctx context.Context, uuid string) (*domain.TestRun, error)
	UpdateTestRunStatus(ctx context.Context, id int64, status domain.TestRunStatus) error
	UpdateTestRunCompletion(ctx context.Context, id int64, completedAt *time.Time, durationMs int) error
	UpdateTestRunStats(ctx context.Context, id int64, total, passed, failed, skipped int) error
	ListTestRuns(ctx context.Context, limit, offset int) ([]*domain.TestRun, error)
	DeleteTestRun(ctx context.Context, id int64) error

	// TestResults
	CreateTestResult(ctx context.Context, result *domain.TestResult) error
	GetTestResultsByRunID(ctx context.Context, runID int64) ([]*domain.TestResult, error)

	// TestLogs
	CreateTestLog(ctx context.Context, log *domain.TestLog) error
	GetTestLogsByRunID(ctx context.Context, runID int64, limit int) ([]*domain.TestLog, error)

	// Test Data Management
	GetTestDataStats(ctx context.Context) (*domain.TestDataStats, error)
	CleanupTestData(ctx context.Context, types []string) (map[string]int64, error)
}
