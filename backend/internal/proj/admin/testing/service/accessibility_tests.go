// Package service implements accessibility functional tests
// backend/internal/proj/admin/testing/service/accessibility_tests.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"backend/internal/proj/admin/testing/domain"
)

// AccessibilityTests returns list of accessibility tests
var AccessibilityTests = []FunctionalTest{
	{
		Name:        "a11y-wcag-compliance",
		Category:    domain.TestCategoryAccessibility,
		Description: "Test WCAG 2.1 AA compliance using axe-core",
		RunFunc:     testWCAGCompliance,
	},
	{
		Name:        "a11y-keyboard-navigation",
		Category:    domain.TestCategoryAccessibility,
		Description: "Test keyboard navigation on all interactive elements",
		RunFunc:     testKeyboardNavigation,
	},
}

// PlaywrightAxeResult represents the result from Playwright axe-core test
type PlaywrightAxeResult struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error,omitempty"`
	Spec    string                 `json:"spec"`
	Stats   PlaywrightStats        `json:"stats"`
	Results map[string]interface{} `json:"results,omitempty"`
}

// PlaywrightStats represents test execution stats
type PlaywrightStats struct {
	Duration int `json:"duration"`
	Expected int `json:"expected"`
	Failures int `json:"failures"`
	Passes   int `json:"passes"`
	Skipped  int `json:"skipped"`
}

// testWCAGCompliance tests WCAG 2.1 AA compliance using Playwright + axe-core
func testWCAGCompliance(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "a11y-wcag-compliance",
		TestSuite: "accessibility",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Run Playwright test through Go wrapper
	testResult := runPlaywrightAccessibilityTest(ctx, "a11y-wcag-compliance.spec.ts", token)
	if testResult == nil {
		return failTestA11y(result, "Failed to run Playwright accessibility test", fmt.Errorf("test result is nil"))
	}

	if !testResult.Success {
		return failTestA11y(result, fmt.Sprintf("WCAG compliance test failed: %s", testResult.Error), nil)
	}

	if testResult.Stats.Failures > 0 {
		return failTestA11y(result, fmt.Sprintf("Found %d WCAG compliance violations", testResult.Stats.Failures), nil)
	}

	result.CompletedAt = time.Now().UTC()
	return result
}

// testKeyboardNavigation tests keyboard navigation functionality
func testKeyboardNavigation(ctx context.Context, baseURL, token string) *domain.TestResult {
	result := &domain.TestResult{
		TestName:  "a11y-keyboard-navigation",
		TestSuite: "accessibility",
		Status:    domain.TestResultStatusPassed,
		StartedAt: time.Now().UTC(),
	}

	// Run Playwright test through Go wrapper
	testResult := runPlaywrightAccessibilityTest(ctx, "a11y-keyboard-navigation.spec.ts", token)
	if testResult == nil {
		return failTestA11y(result, "Failed to run Playwright keyboard test", fmt.Errorf("test result is nil"))
	}

	if !testResult.Success {
		return failTestA11y(result, fmt.Sprintf("Keyboard navigation test failed: %s", testResult.Error), nil)
	}

	if testResult.Stats.Failures > 0 {
		return failTestA11y(result, fmt.Sprintf("Found %d keyboard navigation issues", testResult.Stats.Failures), nil)
	}

	result.CompletedAt = time.Now().UTC()
	return result
}

// runPlaywrightAccessibilityTest executes a Playwright accessibility test and returns JSON result
func runPlaywrightAccessibilityTest(ctx context.Context, testFile, token string) *PlaywrightAxeResult {
	// Find frontend directory (relative to backend)
	backendDir, _ := os.Getwd()
	frontendDir := filepath.Join(backendDir, "..", "frontend", "svetu")

	// Check if frontend directory exists
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		// Try alternative path
		frontendDir = "/data/hostel-booking-system/frontend/svetu"
		if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
			return &PlaywrightAxeResult{
				Success: false,
				Error:   fmt.Sprintf("Frontend directory not found: %s", frontendDir),
			}
		}
	}

	// Create axe tests directory if it doesn't exist
	axeTestsDir := filepath.Join(frontendDir, "e2e", "axe")
	if err := os.MkdirAll(axeTestsDir, 0o755); err != nil {
		return &PlaywrightAxeResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to create axe tests directory: %v", err),
		}
	}

	// Check if test file exists
	testFilePath := filepath.Join(axeTestsDir, testFile)
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		return &PlaywrightAxeResult{
			Success: false,
			Error:   fmt.Sprintf("Test file not found: %s", testFilePath),
		}
	}

	// Get admin credentials from environment
	adminEmail := os.Getenv("TEST_ADMIN_EMAIL")
	adminPassword := os.Getenv("TEST_ADMIN_PASSWORD")

	if adminEmail == "" {
		adminEmail = "admin@admin.rs"
	}
	if adminPassword == "" {
		adminPassword = "P@$S4@dmi№"
	}

	// Prepare command: npx playwright test e2e/axe/{testFile} --reporter=json
	cmd := exec.CommandContext(ctx, "npx", "playwright", "test", fmt.Sprintf("e2e/axe/%s", testFile), "--reporter=json")
	cmd.Dir = frontendDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TEST_ADMIN_EMAIL=%s", adminEmail),
		fmt.Sprintf("TEST_ADMIN_PASSWORD=%s", adminPassword),
		"CI=true", // Run in CI mode for cleaner output
	)

	// Execute command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Parse JSON output even on error (tests may have failed)
		var result PlaywrightAxeResult
		if jsonErr := json.Unmarshal(output, &result); jsonErr == nil {
			return &result
		}

		// If JSON parsing fails, return error
		return &PlaywrightAxeResult{
			Success: false,
			Error:   fmt.Sprintf("Playwright command failed: %v, output: %s", err, string(output)),
		}
	}

	// Parse JSON output
	var result PlaywrightAxeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return &PlaywrightAxeResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse Playwright JSON output: %v, output: %s", err, string(output)),
		}
	}

	// Mark as success if there are no failures
	result.Success = result.Stats.Failures == 0

	return &result
}

// failTestA11y helper to mark accessibility test as failed
func failTestA11y(result *domain.TestResult, message string, err error) *domain.TestResult {
	result.Status = domain.TestResultStatusFailed
	result.CompletedAt = time.Now().UTC()

	errorMsg := message
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %v", message, err)
	}
	result.ErrorMsg = &errorMsg

	// Stack trace for error
	if err != nil {
		stackTrace := fmt.Sprintf("%+v", err)
		result.StackTrace = &stackTrace
	}

	return result
}
