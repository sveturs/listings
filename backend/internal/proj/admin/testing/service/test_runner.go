// Package service implements test runner service
// backend/internal/proj/admin/testing/service/test_runner.go
package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"backend/internal/proj/admin/testing/domain"
	"backend/internal/proj/admin/testing/storage"
)

// Sentinel errors
var (
	ErrTestRunNotFound = errors.New("test run not found")
)

// AuthTokenProvider interface for getting auth tokens (both real and mock)
type AuthTokenProvider interface {
	GetToken() (string, error)
	ClearToken()
}

// TestRunContext holds context for running test
type TestRunContext struct {
	TestRunID   int64
	Cancel      context.CancelFunc
	CompletedCh chan struct{}
}

// TestRunner manages test execution
type TestRunner struct {
	storage      storage.TestStorage
	testAuthMgr  AuthTokenProvider
	logger       zerolog.Logger
	backendURL   string
	runningTests map[int64]*TestRunContext
	mu           sync.RWMutex
}

// NewTestRunner creates new test runner instance
func NewTestRunner(
	storage storage.TestStorage,
	testAuthMgr AuthTokenProvider,
	backendURL string,
	logger zerolog.Logger,
) *TestRunner {
	return &TestRunner{
		storage:      storage,
		testAuthMgr:  testAuthMgr,
		backendURL:   backendURL,
		logger:       logger.With().Str("component", "test_runner").Logger(),
		runningTests: make(map[int64]*TestRunContext),
	}
}

// GetAvailableTestSuites returns list of available test suites
func (r *TestRunner) GetAvailableTestSuites() []*domain.TestSuite {
	return []*domain.TestSuite{
		{
			Name:        "api-endpoints",
			Description: "Test all API endpoints functionality",
			Category:    domain.TestCategoryAPI,
			TestCount:   len(APIEndpointTests),
			Enabled:     true,
		},
		{
			Name:        "functional-api",
			Description: "Functional API tests (alias for api-endpoints)",
			Category:    domain.TestCategoryAPI,
			TestCount:   len(APIEndpointTests),
			Enabled:     true,
		},
		{
			Name:        "integration",
			Description: "Integration tests for external services (Redis, OpenSearch, PostgreSQL)",
			Category:    domain.TestCategoryIntegration,
			TestCount:   len(IntegrationTests),
			Enabled:     true,
		},
		{
			Name:        "security",
			Description: "Security tests (SQL injection, XSS, file upload, auth, rate limiting, CSRF)",
			Category:    domain.TestCategorySecurity,
			TestCount:   len(SecurityTests),
			Enabled:     true,
		},
		{
			Name:        "performance",
			Description: "Performance tests (response time, concurrent users, database queries, memory)",
			Category:    domain.TestCategoryPerformance,
			TestCount:   len(PerformanceTests),
			Enabled:     true,
		},
		{
			Name:        "data-integrity",
			Description: "Data integrity tests (listing consistency, transaction rollback, orphan cleanup)",
			Category:    domain.TestCategoryDataIntegrity,
			TestCount:   len(DataIntegrityTests),
			Enabled:     true,
		},
		{
			Name:        "e2e",
			Description: "End-to-end tests (create listing, search & contact, admin moderation)",
			Category:    domain.TestCategoryE2E,
			TestCount:   len(E2ETests),
			Enabled:     true,
		},
		{
			Name:        "monitoring",
			Description: "Monitoring and observability tests (health endpoints, metrics, error logging)",
			Category:    domain.TestCategoryMonitoring,
			TestCount:   len(MonitoringTests),
			Enabled:     true,
		},
		{
			Name:        "accessibility",
			Description: "Accessibility tests (WCAG 2.1 AA compliance, keyboard navigation)",
			Category:    domain.TestCategoryAccessibility,
			TestCount:   len(AccessibilityTests),
			Enabled:     true,
		},
		{
			Name:        "all",
			Description: "Run all available tests",
			Category:    domain.TestCategoryAll,
			TestCount:   len(APIEndpointTests) + len(IntegrationTests) + len(SecurityTests) + len(PerformanceTests) + len(DataIntegrityTests) + len(E2ETests) + len(MonitoringTests) + len(AccessibilityTests),
			Enabled:     true,
		},
	}
}

// RunTestSuite initiates test suite execution
// If testName is provided, only that specific test will be run
func (r *TestRunner) RunTestSuite(
	ctx context.Context,
	suite string,
	testName string,
	userID int,
	parallel bool,
) (*domain.TestRun, error) {
	// Create test run record
	testRun := &domain.TestRun{
		TestSuite:       suite,
		Status:          domain.TestRunStatusPending,
		StartedByUserID: userID,
		StartedAt:       time.Now().UTC(),
		TotalTests:      0,
		PassedTests:     0,
		FailedTests:     0,
		SkippedTests:    0,
	}

	if err := r.storage.CreateTestRun(ctx, testRun); err != nil {
		r.logger.Error().Err(err).Msg("Failed to create test run")
		return nil, fmt.Errorf("failed to create test run: %w", err)
	}

	r.logger.Info().
		Int64("run_id", testRun.ID).
		Str("suite", suite).
		Str("test_name", testName).
		Int("user_id", userID).
		Msg("Test run created")

	// Start test execution in background
	go r.executeTestSuite(ctx, testRun, testName, parallel)

	return testRun, nil
}

// executeTestSuite executes test suite in background
func (r *TestRunner) executeTestSuite(parentCtx context.Context, testRun *domain.TestRun, testName string, parallel bool) {
	// Use background context instead of parentCtx to avoid HTTP request timeout cancellation
	// Tests can run for up to 30 minutes independently of the HTTP request that initiated them
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Register running test
	r.mu.Lock()
	completedCh := make(chan struct{})
	r.runningTests[testRun.ID] = &TestRunContext{
		TestRunID:   testRun.ID,
		Cancel:      cancel,
		CompletedCh: completedCh,
	}
	r.mu.Unlock()

	// Cleanup on completion
	defer func() {
		r.mu.Lock()
		delete(r.runningTests, testRun.ID)
		r.mu.Unlock()
		close(completedCh)
	}()

	// Update status to running
	if err := r.storage.UpdateTestRunStatus(ctx, testRun.ID, domain.TestRunStatusRunning); err != nil {
		r.logger.Error().Err(err).Int64("run_id", testRun.ID).Msg("Failed to update test run status")
		return
	}

	r.logger.Info().Int64("run_id", testRun.ID).Msg("Test run started")

	// Get admin token for tests
	token, err := r.testAuthMgr.GetToken()
	if err != nil {
		r.logger.Error().Err(err).Int64("run_id", testRun.ID).Msg("Failed to get admin token")
		r.logTestError(ctx, testRun.ID, "Failed to get admin token: "+err.Error())
		r.completeTestRun(ctx, testRun.ID, domain.TestRunStatusFailed, 0, 0, 0, 0)
		return
	}

	// Get tests for suite
	tests := r.getTestsForSuite(testRun.TestSuite, testName)
	if len(tests) == 0 {
		r.logger.Warn().
			Str("suite", testRun.TestSuite).
			Str("test_name", testName).
			Msg("No tests found for suite/test")
		r.completeTestRun(ctx, testRun.ID, domain.TestRunStatusCompleted, 0, 0, 0, 0)
		return
	}

	// Execute tests
	var results []*domain.TestResult
	if parallel {
		results = r.executeTestsParallel(ctx, testRun.ID, tests, token)
	} else {
		results = r.executeTestsSequential(ctx, testRun.ID, tests, token)
	}

	// Calculate stats
	totalTests := len(results)
	passedTests := 0
	failedTests := 0
	skippedTests := 0

	for _, result := range results {
		switch result.Status {
		case domain.TestResultStatusPassed:
			passedTests++
		case domain.TestResultStatusFailed:
			failedTests++
		case domain.TestResultStatusSkipped:
			skippedTests++
		}
	}

	// Determine final status
	finalStatus := domain.TestRunStatusCompleted
	if failedTests > 0 {
		finalStatus = domain.TestRunStatusFailed
	}

	// Complete test run
	r.completeTestRun(ctx, testRun.ID, finalStatus, totalTests, passedTests, failedTests, skippedTests)

	r.logger.Info().
		Int64("run_id", testRun.ID).
		Str("status", string(finalStatus)).
		Int("total", totalTests).
		Int("passed", passedTests).
		Int("failed", failedTests).
		Msg("Test run completed")
}

// executeTestsSequential executes tests one by one
func (r *TestRunner) executeTestsSequential(
	ctx context.Context,
	runID int64,
	tests []FunctionalTest,
	token string,
) []*domain.TestResult {
	results := make([]*domain.TestResult, 0, len(tests))

	for _, test := range tests {
		r.logTestInfo(ctx, runID, fmt.Sprintf("Running test: %s", test.Name))

		result := test.RunFunc(ctx, r.backendURL, token)
		result.TestRunID = runID

		// Save result
		if err := r.storage.CreateTestResult(ctx, result); err != nil {
			r.logger.Error().Err(err).Str("test", test.Name).Msg("Failed to save test result")
		}

		results = append(results, result)

		// Log result
		if result.Status == domain.TestResultStatusPassed {
			r.logTestInfo(ctx, runID, fmt.Sprintf("✓ %s passed (%dms)", test.Name, result.DurationMs))
		} else {
			msg := fmt.Sprintf("✗ %s failed (%dms)", test.Name, result.DurationMs)
			if result.ErrorMsg != nil {
				msg += ": " + *result.ErrorMsg
			}
			r.logTestError(ctx, runID, msg)
		}
	}

	return results
}

// executeTestsParallel executes tests in parallel
func (r *TestRunner) executeTestsParallel(
	ctx context.Context,
	runID int64,
	tests []FunctionalTest,
	token string,
) []*domain.TestResult {
	results := make([]*domain.TestResult, len(tests))
	var wg sync.WaitGroup

	for i, test := range tests {
		wg.Add(1)
		go func(idx int, t FunctionalTest) {
			defer wg.Done()

			r.logTestInfo(ctx, runID, fmt.Sprintf("Running test: %s", t.Name))

			result := t.RunFunc(ctx, r.backendURL, token)
			result.TestRunID = runID

			// Save result
			if err := r.storage.CreateTestResult(ctx, result); err != nil {
				r.logger.Error().Err(err).Str("test", t.Name).Msg("Failed to save test result")
			}

			results[idx] = result

			// Log result
			if result.Status == domain.TestResultStatusPassed {
				r.logTestInfo(ctx, runID, fmt.Sprintf("✓ %s passed (%dms)", t.Name, result.DurationMs))
			} else {
				msg := fmt.Sprintf("✗ %s failed (%dms)", t.Name, result.DurationMs)
				if result.ErrorMsg != nil {
					msg += ": " + *result.ErrorMsg
				}
				r.logTestError(ctx, runID, msg)
			}
		}(i, test)
	}

	wg.Wait()
	return results
}

// getTestsForSuite returns tests for specified suite
func (r *TestRunner) getTestsForSuite(suite string, testName string) []FunctionalTest {
	var tests []FunctionalTest

	switch suite {
	case "api-endpoints", "functional-api":
		tests = APIEndpointTests
	case "integration":
		tests = IntegrationTests
	case "security":
		tests = SecurityTests
	case "performance":
		tests = PerformanceTests
	case "data-integrity":
		tests = DataIntegrityTests
	case "e2e":
		tests = E2ETests
	case "monitoring":
		tests = MonitoringTests
	case "accessibility":
		tests = AccessibilityTests
	case "all":
		// Combine all test suites
		tests = append(tests, APIEndpointTests...)
		tests = append(tests, IntegrationTests...)
		tests = append(tests, SecurityTests...)
		tests = append(tests, PerformanceTests...)
		tests = append(tests, DataIntegrityTests...)
		tests = append(tests, E2ETests...)
		tests = append(tests, MonitoringTests...)
		tests = append(tests, AccessibilityTests...)
	default:
		return nil
	}

	// If specific test name is provided, filter to only that test
	if testName != "" {
		for _, test := range tests {
			if test.Name == testName {
				return []FunctionalTest{test}
			}
		}
		// Test not found
		return nil
	}

	return tests
}

// completeTestRun marks test run as completed and updates stats
func (r *TestRunner) completeTestRun(
	ctx context.Context,
	runID int64,
	status domain.TestRunStatus,
	total, passed, failed, skipped int,
) {
	now := time.Now().UTC()

	// Get test run to calculate duration
	testRun, err := r.storage.GetTestRunByID(ctx, runID)
	if err != nil {
		r.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to get test run")
		return
	}

	durationMs := int(now.Sub(testRun.StartedAt).Milliseconds())

	// Update status
	if err := r.storage.UpdateTestRunStatus(ctx, runID, status); err != nil {
		r.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to update status")
	}

	// Update completion
	if err := r.storage.UpdateTestRunCompletion(ctx, runID, &now, durationMs); err != nil {
		r.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to update completion")
	}

	// Update stats
	if err := r.storage.UpdateTestRunStats(ctx, runID, total, passed, failed, skipped); err != nil {
		r.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to update stats")
	}
}

// logTestInfo logs info message
func (r *TestRunner) logTestInfo(ctx context.Context, runID int64, message string) {
	log := &domain.TestLog{
		TestRunID: runID,
		Level:     "info",
		Message:   message,
		Timestamp: time.Now(),
	}
	_ = r.storage.CreateTestLog(ctx, log)
}

// logTestError logs error message
func (r *TestRunner) logTestError(ctx context.Context, runID int64, message string) {
	log := &domain.TestLog{
		TestRunID: runID,
		Level:     "error",
		Message:   message,
		Timestamp: time.Now(),
	}
	_ = r.storage.CreateTestLog(ctx, log)
}

// GetTestRunStatus returns current test run status
func (r *TestRunner) GetTestRunStatus(ctx context.Context, runID int64) (*domain.TestRun, error) {
	return r.storage.GetTestRunByID(ctx, runID)
}

// GetTestRunDetail returns detailed test run information
func (r *TestRunner) GetTestRunDetail(ctx context.Context, runID int64) (*domain.TestRunDetail, error) {
	testRun, err := r.storage.GetTestRunByID(ctx, runID)
	if err != nil {
		return nil, err
	}

	// Get results
	results, err := r.storage.GetTestResultsByRunID(ctx, runID)
	if err != nil {
		r.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to get test results")
		results = nil
	}

	// Get logs
	logs, err := r.storage.GetTestLogsByRunID(ctx, runID, 1000)
	if err != nil {
		r.logger.Error().Err(err).Int64("run_id", runID).Msg("Failed to get test logs")
		logs = nil
	}

	return &domain.TestRunDetail{
		TestRun: testRun,
		Results: results,
		Logs:    logs,
	}, nil
}

// ListTestRuns returns paginated list of test runs
func (r *TestRunner) ListTestRuns(ctx context.Context, limit, offset int) ([]*domain.TestRun, error) {
	return r.storage.ListTestRuns(ctx, limit, offset)
}

// CancelTestRun cancels running test
func (r *TestRunner) CancelTestRun(runID int64) error {
	r.mu.RLock()
	testCtx, exists := r.runningTests[runID]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("test run %d is not running", runID)
	}

	testCtx.Cancel()

	// Wait for completion with timeout
	select {
	case <-testCtx.CompletedCh:
		r.logger.Info().Int64("run_id", runID).Msg("Test run canceled")
	case <-time.After(10 * time.Second):
		r.logger.Warn().Int64("run_id", runID).Msg("Test run cancellation timeout")
	}

	return nil
}

// GetAllAvailableTests returns list of all available tests across all suites
func (r *TestRunner) GetAllAvailableTests() []domain.AvailableTest {
	var allTests []domain.AvailableTest

	// Helper function to convert FunctionalTest to AvailableTest
	convertTest := func(test FunctionalTest, suite string) domain.AvailableTest {
		return domain.AvailableTest{
			ID:          test.Name,
			Name:        test.Name,
			Description: test.Description,
			Category:    test.Category,
			Suite:       suite,
		}
	}

	// Add API tests
	for _, test := range APIEndpointTests {
		allTests = append(allTests, convertTest(test, "functional-api"))
	}

	// Add Integration tests
	for _, test := range IntegrationTests {
		allTests = append(allTests, convertTest(test, "integration"))
	}

	// Add Security tests
	for _, test := range SecurityTests {
		allTests = append(allTests, convertTest(test, "security"))
	}

	// Add Performance tests
	for _, test := range PerformanceTests {
		allTests = append(allTests, convertTest(test, "performance"))
	}

	// Add Data Integrity tests
	for _, test := range DataIntegrityTests {
		allTests = append(allTests, convertTest(test, "data-integrity"))
	}

	// Add E2E tests
	for _, test := range E2ETests {
		allTests = append(allTests, convertTest(test, "e2e"))
	}

	// Add Monitoring tests
	for _, test := range MonitoringTests {
		allTests = append(allTests, convertTest(test, "monitoring"))
	}

	// Add Accessibility tests
	for _, test := range AccessibilityTests {
		allTests = append(allTests, convertTest(test, "accessibility"))
	}

	return allTests
}
