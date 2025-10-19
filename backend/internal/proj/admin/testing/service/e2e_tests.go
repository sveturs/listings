// Package service implements E2E tests through Playwright
// backend/internal/proj/admin/testing/service/e2e_tests.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// E2E test definitions - these tests run Playwright scripts
var E2ETests = []FunctionalTest{
	{
		Name:        "e2e-user-journey-create-listing",
		Description: "Full flow: login → create listing → upload images → publish",
		Category:    domain.TestCategoryE2E,
		RunFunc:     testE2EUserJourneyCreateListing,
	},
	{
		Name:        "e2e-user-journey-search-contact",
		Description: "Search → view listing → contact seller",
		Category:    domain.TestCategoryE2E,
		RunFunc:     testE2EUserJourneySearchContact,
	},
	{
		Name:        "e2e-admin-moderation",
		Description: "Admin reviews and approves/rejects listing",
		Category:    domain.TestCategoryE2E,
		RunFunc:     testE2EAdminModeration,
	},
}

// PlaywrightTestResult represents Playwright test result
type PlaywrightTestResult struct {
	Status   string `json:"status"`   // passed, failed, timedOut, skipped
	Duration int    `json:"duration"` // milliseconds
	Error    string `json:"error,omitempty"`
}

// PlaywrightSuiteResult represents entire suite result
type PlaywrightSuiteResult struct {
	Stats struct {
		Expected   int `json:"expected"`
		Unexpected int `json:"unexpected"`
		Skipped    int `json:"skipped"`
	} `json:"stats"`
	Suites []struct {
		Tests []struct {
			Title   string `json:"title"`
			Results []struct {
				Status   string `json:"status"`
				Duration int    `json:"duration"`
				Error    struct {
					Message string `json:"message"`
				} `json:"error"`
			} `json:"results"`
		} `json:"tests"`
	} `json:"suites"`
}

// runPlaywrightTest executes a Playwright test and returns result
func runPlaywrightTest(ctx context.Context, testFile string, testName string, backendURL string) *domain.TestResult {
	startTime := time.Now()

	// Find frontend directory
	frontendDir := filepath.Join(getProjectRoot(), "frontend", "svetu")

	// Prepare Playwright command
	// Run specific test file with grep pattern for test name
	args := []string{
		"playwright",
		"test",
		filepath.Join("e2e", testFile),
		"--reporter=json",
		"--project=chromium",
	}

	// Set environment variables
	env := []string{
		"NEXT_PUBLIC_FRONTEND_URL=http://localhost:3001",
		fmt.Sprintf("TEST_ADMIN_EMAIL=%s", getEnvOrDefault("TEST_ADMIN_EMAIL", "admin@admin.rs")),
		fmt.Sprintf("TEST_ADMIN_PASSWORD=%s", getEnvOrDefault("TEST_ADMIN_PASSWORD", "P@$S4@dmi№")),
		"CI=true", // Don't start dev server, assume it's already running
	}

	// Execute Playwright test
	cmd := exec.CommandContext(ctx, "npx", args...) //nolint:gosec // G204: args constructed from internal test files, not user input
	cmd.Dir = frontendDir
	cmd.Env = append(os.Environ(), env...)

	output, err := cmd.CombinedOutput()
	duration := int(time.Since(startTime).Milliseconds())

	// Parse Playwright JSON output
	result := &domain.TestResult{
		TestName:    testName,
		TestSuite:   "e2e",
		StartedAt:   startTime,
		CompletedAt: time.Now(),
		DurationMs:  duration,
	}

	if err != nil {
		// Test failed
		errorMsg := fmt.Sprintf("Playwright test failed: %v\nOutput: %s", err, string(output))
		result.Status = domain.TestResultStatusFailed
		result.ErrorMsg = &errorMsg

		// Try to extract error from output
		if len(output) > 0 {
			stackTrace := string(output)
			result.StackTrace = &stackTrace
		}

		return result
	}

	// Try to parse JSON output
	var suiteResult PlaywrightSuiteResult
	if jsonErr := json.Unmarshal(output, &suiteResult); jsonErr == nil {
		// Successfully parsed JSON
		switch {
		case suiteResult.Stats.Unexpected > 0:
			// Tests failed
			errorMsg := fmt.Sprintf("%d tests failed", suiteResult.Stats.Unexpected)

			// Extract first error message
			for _, suite := range suiteResult.Suites {
				for _, test := range suite.Tests {
					for _, testResult := range test.Results {
						if testResult.Status == "failed" && testResult.Error.Message != "" {
							errorMsg = testResult.Error.Message
							break
						}
					}
				}
			}

			result.Status = domain.TestResultStatusFailed
			result.ErrorMsg = &errorMsg
		case suiteResult.Stats.Expected > 0:
			// All tests passed
			result.Status = domain.TestResultStatusPassed
		default:
			// Skipped
			result.Status = domain.TestResultStatusSkipped
		}
	} else {
		// Could not parse JSON, but command succeeded
		// Check output for success indicators
		outputStr := string(output)
		if strings.Contains(outputStr, "passed") || strings.Contains(outputStr, "✓") {
			result.Status = domain.TestResultStatusPassed
		} else {
			result.Status = domain.TestResultStatusFailed
			errorMsg := fmt.Sprintf("Could not parse test result: %s", outputStr)
			result.ErrorMsg = &errorMsg
		}
	}

	return result
}

// testE2EUserJourneyCreateListing tests full listing creation flow
func testE2EUserJourneyCreateListing(ctx context.Context, backendURL string, token string) *domain.TestResult {
	return runPlaywrightTest(
		ctx,
		"user-journey-create-listing.spec.ts",
		"e2e-user-journey-create-listing",
		backendURL,
	)
}

// testE2EUserJourneySearchContact tests search and contact flow
func testE2EUserJourneySearchContact(ctx context.Context, backendURL string, token string) *domain.TestResult {
	return runPlaywrightTest(
		ctx,
		"user-journey-search-contact.spec.ts",
		"e2e-user-journey-search-contact",
		backendURL,
	)
}

// testE2EAdminModeration tests admin moderation workflow
func testE2EAdminModeration(ctx context.Context, backendURL string, token string) *domain.TestResult {
	return runPlaywrightTest(
		ctx,
		"admin-moderation-flow.spec.ts",
		"e2e-admin-moderation",
		backendURL,
	)
}

// getProjectRoot returns project root directory
func getProjectRoot() string {
	// Try to get from environment
	if root := getEnvOrDefault("PROJECT_ROOT", ""); root != "" {
		return root
	}

	// Try absolute path first (production deployment)
	absolutePath := "/data/hostel-booking-system"
	if _, err := os.Stat(absolutePath); err == nil {
		return absolutePath
	}

	// Fallback: try relative path from current working directory
	cwd, err := os.Getwd()
	if err == nil {
		// If we're in backend directory, go up one level
		if strings.Contains(cwd, "/backend") {
			return filepath.Join(cwd, "..")
		}
		// If we're already at project root
		if strings.Contains(cwd, "hostel-booking-system") && !strings.Contains(cwd, "/backend") {
			return cwd
		}
	}

	// Last resort: assume we're in backend/internal/proj/admin/testing/service
	return filepath.Join("..", "..", "..", "..", "..", "..")
}

// getEnvOrDefault gets environment variable or returns default
func getEnvOrDefault(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
