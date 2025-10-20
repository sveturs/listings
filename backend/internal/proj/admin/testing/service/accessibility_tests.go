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
	"strings"
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
		// Extract only first few violations from error message (avoid 196+ lines of JSON per violation)
		conciseError := extractConciseViolationsSummary(testResult.Error, testResult.Stats.Failures)
		errorMsg := fmt.Sprintf("WCAG compliance test failed: %s", conciseError)
		return failTestA11y(result, errorMsg, nil)
	}

	if testResult.Stats.Failures > 0 {
		// Extract only first few violations from error message
		conciseError := extractConciseViolationsSummary(testResult.Error, testResult.Stats.Failures)
		errorMsg := fmt.Sprintf("Found %d WCAG compliance violations:\n%s", testResult.Stats.Failures, conciseError)
		return failTestA11y(result, errorMsg, nil)
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
		// Extract only first few violations from error message (avoid 196+ lines of JSON per violation)
		conciseError := extractConciseViolationsSummary(testResult.Error, testResult.Stats.Failures)
		errorMsg := fmt.Sprintf("Keyboard navigation test failed: %s", conciseError)
		return failTestA11y(result, errorMsg, nil)
	}

	if testResult.Stats.Failures > 0 {
		// Extract only first few violations from error message
		conciseError := extractConciseViolationsSummary(testResult.Error, testResult.Stats.Failures)
		errorMsg := fmt.Sprintf("Found %d keyboard navigation issues:\n%s", testResult.Stats.Failures, conciseError)
		return failTestA11y(result, errorMsg, nil)
	}

	result.CompletedAt = time.Now().UTC()
	return result
}

// runPlaywrightAccessibilityTest executes a Playwright accessibility test and returns JSON result
func runPlaywrightAccessibilityTest(ctx context.Context, testFile, token string) *PlaywrightAxeResult {
	// Create context with timeout for Playwright execution
	// Accessibility tests run 6-12 subtests, each can take up to 5 minutes with axe scans
	// Allow 30 minutes total to handle all scenarios and slow network
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	// Find frontend directory using absolute path
	frontendDir := "/data/hostel-booking-system/frontend/svetu"

	// Check if frontend directory exists
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		// Try relative path from current working directory
		cwd, _ := os.Getwd()
		frontendDir = filepath.Join(cwd, "..", "frontend", "svetu")

		if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
			return &PlaywrightAxeResult{
				Success: false,
				Error:   fmt.Sprintf("Frontend directory not found at %s or relative path", frontendDir),
			}
		}
	}

	// Create axe tests directory if it doesn't exist
	axeTestsDir := filepath.Join(frontendDir, "e2e", "axe")
	if err := os.MkdirAll(axeTestsDir, 0o755); err != nil { //nolint:gosec // G301: 0755 permissions required for Playwright test directory
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
		adminPassword = "P@$S4@dmi№" //nolint:gosec // G101: Test password from env or default
	}

	// Prepare command: npx playwright test e2e/axe/{testFile} --reporter=json --timeout=120000
	// Timeout reduced to 120000ms = 2 minutes per test (tests now use domcontentloaded)
	// 3 tests * ~30-40s each = ~2 min total
	cmd := exec.CommandContext(ctx, "npx", "playwright", "test", fmt.Sprintf("e2e/axe/%s", testFile), "--reporter=json", "--timeout=120000") //nolint:gosec // G204: testFile is from internal test registry, not user input
	cmd.Dir = frontendDir
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TEST_ADMIN_EMAIL=%s", adminEmail),
		fmt.Sprintf("TEST_ADMIN_PASSWORD=%s", adminPassword),
		"CI=true",        // Run in CI mode for cleaner output
		"LOG_LEVEL=warn", // Suppress backend INFO logs during tests
	)

	// Execute command and capture output
	output, err := cmd.CombinedOutput()

	// Parse Playwright JSON output (handle both success and failure cases)
	result := parsePlaywrightJSON(output)

	if err != nil && result == nil {
		// Command failed AND JSON parsing failed - extract error lines
		errorLines := extractErrorLines(string(output))
		return &PlaywrightAxeResult{
			Success: false,
			Error:   fmt.Sprintf("Playwright command failed: %v\n%s", err, errorLines),
		}
	}

	// If we got a result (even if command failed), use it
	if result != nil {
		// Mark as success if there are no failures
		result.Success = result.Stats.Failures == 0
		return result
	}

	// Fallback: shouldn't reach here, but handle it
	return &PlaywrightAxeResult{
		Success: false,
		Error:   "Failed to parse Playwright output (no result available)",
	}
}

// parsePlaywrightJSON safely parses Playwright JSON output, handling UTF-8 and different formats
func parsePlaywrightJSON(output []byte) *PlaywrightAxeResult {
	if len(output) == 0 {
		return nil
	}

	// Playwright JSON reporter outputs to stdout with possible non-JSON prefix (logs)
	// Find the start of JSON (first { character)
	jsonStart := -1
	for i := 0; i < len(output); i++ {
		if output[i] == '{' {
			jsonStart = i
			break
		}
	}

	if jsonStart == -1 {
		// No JSON found in output
		return nil
	}

	// Extract JSON portion only
	jsonData := output[jsonStart:]

	// Parse Playwright test reporter JSON format
	// Format: { "config": {...}, "suites": [...], "errors": [...], ... }
	var playwrightReport struct {
		Config struct {
			Projects []struct {
				Name string `json:"name"`
			} `json:"projects"`
		} `json:"config"`
		Suites []struct {
			Suites []struct {
				Specs []struct {
					Title string `json:"title"`
					OK    bool   `json:"ok"`
					Tests []struct {
						Results []struct {
							Status   string        `json:"status"`
							Duration float64       `json:"duration"` // Duration can be float (e.g., 8441.503ms)
							Errors   []interface{} `json:"errors"`
						} `json:"results"`
					} `json:"tests"`
				} `json:"specs"`
			} `json:"suites"`
		} `json:"suites"`
		Stats struct {
			StartTime      string  `json:"startTime"`
			Duration       float64 `json:"duration"` // Duration can be float (e.g., 8441.503ms)
			Expected       int     `json:"expected"`
			Skipped        int     `json:"skipped"`
			Unexpected     int     `json:"unexpected"`
			Flaky          int     `json:"flaky"`
			InterruptedInt int     `json:"interrupted"`
		} `json:"stats"`
		Errors []interface{} `json:"errors"`
	}

	if err := json.Unmarshal(jsonData, &playwrightReport); err != nil {
		// Failed to parse as Playwright report JSON
		return nil
	}

	// Convert Playwright format to our PlaywrightAxeResult format
	result := &PlaywrightAxeResult{
		Success: playwrightReport.Stats.Unexpected == 0,
		Stats: PlaywrightStats{
			Duration: int(playwrightReport.Stats.Duration), // Convert float64 to int
			Expected: playwrightReport.Stats.Expected,
			Failures: playwrightReport.Stats.Unexpected,
			Passes:   playwrightReport.Stats.Expected,
			Skipped:  playwrightReport.Stats.Skipped,
		},
	}

	// If there are failures, extract error messages
	if playwrightReport.Stats.Unexpected > 0 {
		var errorMessages []string
		for _, suite := range playwrightReport.Suites {
			for _, innerSuite := range suite.Suites {
				for _, spec := range innerSuite.Specs {
					if !spec.OK {
						for _, test := range spec.Tests {
							for _, testResult := range test.Results {
								if testResult.Status != "passed" && len(testResult.Errors) > 0 {
									// Extract error message from first error
									if errMap, ok := testResult.Errors[0].(map[string]interface{}); ok {
										if msg, exists := errMap["message"]; exists {
											errorMessages = append(errorMessages, fmt.Sprintf("%s: %v", spec.Title, msg))
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if len(errorMessages) > 0 {
			result.Error = strings.Join(errorMessages, "\n\n")
		}
	}

	return result
}

// extractConciseViolationsSummary извлекает краткое резюме violations из Playwright error message
// Пропускает полный JSON объект violations (196+ строк) и оставляет только formatViolations() вывод
func extractConciseViolationsSummary(errorMessage string, failuresCount int) string {
	if errorMessage == "" {
		return "No error details available"
	}

	// Playwright error.message содержит:
	// 1. formatViolations() вывод: "Found X violations:\n1. id: desc\n   Impact:..."
	// 2. Затем идёт diff с JSON: "- Expected  -   1\n+ Received  + 196\n+ Array [\n+   Object {..."
	//
	// Нам нужно взять ТОЛЬКО часть 1 и отбросить часть 2

	// Ищем границу между formatViolations() и JSON diff
	// JSON diff начинается со строк "- Expected" или "+ Array [" или "+ Object {"
	endOfFormattedSection := strings.Index(errorMessage, "\\u001b[32m- Expected")
	if endOfFormattedSection == -1 {
		endOfFormattedSection = strings.Index(errorMessage, "- Expected")
	}
	if endOfFormattedSection == -1 {
		endOfFormattedSection = strings.Index(errorMessage, "+ Array [")
	}
	if endOfFormattedSection == -1 {
		endOfFormattedSection = strings.Index(errorMessage, "+ Received")
	}

	// Если нашли границу - обрезаем
	cleanMessage := errorMessage
	if endOfFormattedSection > 0 {
		cleanMessage = errorMessage[:endOfFormattedSection]
	}

	// Теперь парсим cleanMessage и извлекаем только нужные строки
	lines := strings.Split(cleanMessage, "\n")
	var conciseLines []string
	maxLines := 20 // Максимум 20 строк (вместо 196+!)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Пропускаем пустые строки
		if trimmed == "" {
			continue
		}

		// Пропускаем служебные строки JSON reporter
		if strings.HasPrefix(trimmed, "\"message\":") ||
			strings.HasPrefix(trimmed, "\"stack\":") ||
			strings.HasPrefix(trimmed, "\"location\":") ||
			strings.HasPrefix(trimmed, "\"snippet\":") {
			continue
		}

		// Извлекаем чистый текст из экранированных строк вида \"Error: Found...\"
		// Убираем \" в начале и конце
		if strings.HasPrefix(trimmed, "\"Error:") {
			trimmed = strings.Trim(trimmed, "\",")
		}

		// Заменяем \\n на реальные переносы
		if strings.Contains(trimmed, "\\n") {
			// Разбиваем по \\n и добавляем как отдельные строки
			sublines := strings.Split(trimmed, "\\n")
			for _, subline := range sublines {
				sublineTrimmed := strings.TrimSpace(subline)
				if sublineTrimmed != "" && !strings.HasPrefix(sublineTrimmed, "\\u001b") {
					// Обрезаем Example если слишком длинный
					if strings.Contains(sublineTrimmed, "Example:") && len(sublineTrimmed) > 200 {
						conciseLines = append(conciseLines, sublineTrimmed[:200]+"...")
					} else {
						conciseLines = append(conciseLines, sublineTrimmed)
					}

					if len(conciseLines) >= maxLines {
						break
					}
				}
			}
		} else if !strings.HasPrefix(trimmed, "\\u001b") {
			// Обычная строка
			conciseLines = append(conciseLines, trimmed)
		}

		if len(conciseLines) >= maxLines {
			break
		}
	}

	if len(conciseLines) == 0 {
		// Fallback - берём первые 150 символов
		if len(errorMessage) > 150 {
			return errorMessage[:150] + "... (error details truncated)"
		}
		return errorMessage
	}

	summary := strings.Join(conciseLines, "\n")
	if failuresCount > 1 {
		summary += fmt.Sprintf("\n\n... (%d more violations, details truncated to save space)", failuresCount-1)
	}

	return summary
}

// extractErrorLines извлекает только строки с ошибками из вывода Playwright
// Избегаем сохранения всего JSON config и debug информации
func extractErrorLines(output string) string {
	lines := strings.Split(output, "\n")
	var errorLines []string
	var inErrorBlock bool

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip JSON config noise (config, globalConfig, etc.)
		if strings.Contains(trimmed, `"config":`) ||
			strings.Contains(trimmed, `"globalConfig":`) ||
			strings.Contains(trimmed, `"rootDir":`) ||
			strings.Contains(trimmed, `"testMatch":`) {
			continue
		}

		// Include actual error messages
		if strings.Contains(line, "Error:") ||
			strings.Contains(line, "FAIL") ||
			strings.Contains(line, "✕") ||
			strings.Contains(line, "violations") ||
			strings.Contains(line, "accessibility") {
			errorLines = append(errorLines, trimmed)
			inErrorBlock = true
			continue
		}

		// Include lines after error marker (but limit depth)
		if inErrorBlock && len(errorLines) < 20 {
			// Stop if we hit next section or JSON
			if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
				inErrorBlock = false
				continue
			}
			errorLines = append(errorLines, trimmed)
		}
	}

	// Limit output to 30 lines max
	if len(errorLines) > 30 {
		errorLines = errorLines[:30]
		errorLines = append(errorLines, "... (output truncated)")
	}

	if len(errorLines) == 0 {
		return "No specific error details found in output"
	}

	return strings.Join(errorLines, "\n")
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
