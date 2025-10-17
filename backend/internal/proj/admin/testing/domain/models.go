// Package domain defines domain models for testing module
// backend/internal/proj/admin/testing/domain/models.go
package domain

import (
	"encoding/json"
	"time"
)

// TestRunStatus represents test run execution status
type TestRunStatus string

const (
	TestRunStatusPending   TestRunStatus = "pending"
	TestRunStatusRunning   TestRunStatus = "running"
	TestRunStatusCompleted TestRunStatus = "completed"
	TestRunStatusFailed    TestRunStatus = "failed"
	TestRunStatusCancelled TestRunStatus = "cancelled"
)

// TestResultStatus represents individual test result status
type TestResultStatus string

const (
	TestResultStatusPassed  TestResultStatus = "passed"
	TestResultStatusFailed  TestResultStatus = "failed"
	TestResultStatusSkipped TestResultStatus = "skipped"
)

// TestCategory represents test category type
type TestCategory string

const (
	TestCategoryAPI         TestCategory = "api"
	TestCategoryDatabase    TestCategory = "database"
	TestCategoryPerformance TestCategory = "performance"
	TestCategoryIntegration TestCategory = "integration"
	TestCategoryE2E         TestCategory = "e2e"
)

// TestRun represents a test suite execution
type TestRun struct {
	ID              int64                  `json:"id" db:"id"`
	RunUUID         string                 `json:"run_uuid" db:"run_uuid"`
	TestSuite       string                 `json:"test_suite" db:"test_suite"`
	Status          TestRunStatus          `json:"status" db:"status"`
	StartedByUserID int                    `json:"started_by_user_id" db:"started_by_user_id"`
	StartedAt       time.Time              `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	DurationMs      *int                   `json:"duration_ms,omitempty" db:"duration_ms"`
	TotalTests      int                    `json:"total_tests" db:"total_tests"`
	PassedTests     int                    `json:"passed_tests" db:"passed_tests"`
	FailedTests     int                    `json:"failed_tests" db:"failed_tests"`
	SkippedTests    int                    `json:"skipped_tests" db:"skipped_tests"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}

// TestResult represents individual test execution result
type TestResult struct {
	ID          int64            `json:"id" db:"id"`
	TestRunID   int64            `json:"test_run_id" db:"test_run_id"`
	TestName    string           `json:"test_name" db:"test_name"`
	TestSuite   string           `json:"test_suite" db:"test_suite"`
	Status      TestResultStatus `json:"status" db:"status"`
	DurationMs  int              `json:"duration_ms" db:"duration_ms"`
	ErrorMsg    *string          `json:"error_msg,omitempty" db:"error_msg"`
	StackTrace  *string          `json:"stack_trace,omitempty" db:"stack_trace"`
	StartedAt   time.Time        `json:"started_at" db:"started_at"`
	CompletedAt time.Time        `json:"completed_at" db:"completed_at"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
}

// TestLog represents test execution log entry
type TestLog struct {
	ID        int64     `json:"id" db:"id"`
	TestRunID int64     `json:"test_run_id" db:"test_run_id"`
	Level     string    `json:"level" db:"level"` // info, warn, error, debug
	Message   string    `json:"message" db:"message"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// TestSuite represents available test suite configuration
type TestSuite struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Category    TestCategory `json:"category"`
	TestCount   int          `json:"test_count"`
	Enabled     bool         `json:"enabled"`
}

// RunTestRequest represents request to run test suite
type RunTestRequest struct {
	TestSuite string `json:"test_suite" binding:"required"`
	Parallel  bool   `json:"parallel"`
}

// RunTestResponse represents response after initiating test run
type RunTestResponse struct {
	TestRunID int64  `json:"test_run_id"`
	RunUUID   string `json:"run_uuid"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// TestRunDetail represents detailed test run information with results
type TestRunDetail struct {
	*TestRun
	Results []*TestResult `json:"results,omitempty"`
	Logs    []*TestLog    `json:"logs,omitempty"`
}

// MarshalJSON custom JSON marshaller for Metadata field
func (t *TestRun) MarshalJSON() ([]byte, error) {
	type Alias TestRun
	return json.Marshal(&struct {
		*Alias
		Metadata json.RawMessage `json:"metadata,omitempty"`
	}{
		Alias:    (*Alias)(t),
		Metadata: mustMarshalMetadata(t.Metadata),
	})
}

func mustMarshalMetadata(m map[string]interface{}) json.RawMessage {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	return b
}
