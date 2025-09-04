package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Config holds validation configuration
type Config struct {
	DatabaseURL    string
	PrometheusURL  string
	APIBaseURL     string
	SlackWebhook   string
	CanaryPercent  int
	ValidationMode string // "pre-deployment", "during-rollout", "post-deployment"
}

// ValidationResult holds the results of validation checks
type ValidationResult struct {
	Timestamp   time.Time              `json:"timestamp"`
	Mode        string                 `json:"mode"`
	Passed      bool                   `json:"passed"`
	Checks      []Check                `json:"checks"`
	Metrics     map[string]interface{} `json:"metrics"`
	Recommendations []string          `json:"recommendations"`
}

// Check represents a single validation check
type Check struct {
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Passed      bool        `json:"passed"`
	Value       interface{} `json:"value"`
	Expected    interface{} `json:"expected"`
	Message     string      `json:"message"`
	Severity    string      `json:"severity"` // "critical", "warning", "info"
}

// Validator performs canary rollout validation
type Validator struct {
	config   Config
	db       *sql.DB
	promAPI  v1.API
	client   *http.Client
	results  []Check
}

// NewValidator creates a new validator instance
func NewValidator(config Config) (*Validator, error) {
	// Database connection
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Prometheus client
	promClient, err := api.NewClient(api.Config{
		Address: config.PrometheusURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}

	return &Validator{
		config:  config,
		db:      db,
		promAPI: v1.NewAPI(promClient),
		client:  &http.Client{Timeout: 10 * time.Second},
		results: []Check{},
	}, nil
}

// ValidatePreDeployment performs pre-deployment validation
func (v *Validator) ValidatePreDeployment() ValidationResult {
	log.Println("Starting pre-deployment validation...")
	v.results = []Check{}

	// 1. Check database schema
	v.checkDatabaseSchema()

	// 2. Check data migration completeness
	v.checkDataMigration()

	// 3. Verify indexes exist
	v.checkIndexes()

	// 4. Test query performance
	v.checkQueryPerformance()

	// 5. Verify backup exists
	v.checkBackupStatus()

	// 6. Check application readiness
	v.checkApplicationReadiness()

	return v.generateResult("pre-deployment")
}

// ValidateDuringRollout performs validation during canary rollout
func (v *Validator) ValidateDuringRollout() ValidationResult {
	log.Printf("Validating canary rollout at %d%%...\n", v.config.CanaryPercent)
	v.results = []Check{}

	// 1. Check error rates
	v.checkErrorRate()

	// 2. Check latency
	v.checkLatency()

	// 3. Check dual-write consistency
	v.checkDualWriteConsistency()

	// 4. Check cache performance
	v.checkCachePerformance()

	// 5. Check resource usage
	v.checkResourceUsage()

	// 6. Verify traffic distribution
	v.checkTrafficDistribution()

	// 7. Check fallback triggers
	v.checkFallbackTriggers()

	return v.generateResult("during-rollout")
}

// ValidatePostDeployment performs post-deployment validation
func (v *Validator) ValidatePostDeployment() ValidationResult {
	log.Println("Starting post-deployment validation...")
	v.results = []Check{}

	// 1. Verify all data migrated
	v.checkFinalDataIntegrity()

	// 2. Check system stability
	v.checkSystemStability()

	// 3. Verify old system can be decommissioned
	v.checkOldSystemUsage()

	// 4. Performance comparison
	v.comparePerformance()

	return v.generateResult("post-deployment")
}

// Database validation checks
func (v *Validator) checkDatabaseSchema() {
	query := `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_name IN ('unified_attributes', 'unified_category_attributes', 'unified_attribute_values')
	`
	
	var count int
	err := v.db.QueryRow(query).Scan(&count)
	
	v.addCheck(Check{
		Name:     "Database Schema",
		Category: "database",
		Passed:   err == nil && count == 3,
		Value:    count,
		Expected: 3,
		Message:  "Unified attributes tables exist",
		Severity: "critical",
	})
}

func (v *Validator) checkDataMigration() {
	// Check if all old attributes are migrated
	query := `
		SELECT 
			(SELECT COUNT(*) FROM category_attributes) as old_count,
			(SELECT COUNT(*) FROM unified_attributes) as new_count,
			(SELECT COUNT(*) FROM unified_attributes WHERE legacy_id IS NOT NULL) as migrated_count
	`
	
	var oldCount, newCount, migratedCount int
	err := v.db.QueryRow(query).Scan(&oldCount, &newCount, &migratedCount)
	
	passed := err == nil && migratedCount >= oldCount
	
	v.addCheck(Check{
		Name:     "Data Migration Completeness",
		Category: "database",
		Passed:   passed,
		Value:    fmt.Sprintf("%d/%d migrated", migratedCount, oldCount),
		Expected: "100% migrated",
		Message:  "All attributes migrated to new system",
		Severity: "critical",
	})
}

func (v *Validator) checkIndexes() {
	query := `
		SELECT COUNT(*) FROM pg_indexes
		WHERE tablename LIKE 'unified_%'
		AND indexname LIKE 'idx_unified_%'
	`
	
	var count int
	err := v.db.QueryRow(query).Scan(&count)
	
	v.addCheck(Check{
		Name:     "Performance Indexes",
		Category: "database",
		Passed:   err == nil && count >= 7,
		Value:    count,
		Expected: ">=7",
		Message:  "Required indexes for performance",
		Severity: "warning",
	})
}

func (v *Validator) checkQueryPerformance() {
	start := time.Now()
	
	query := `
		SELECT ua.*, uca.required, uca.display_order
		FROM unified_attributes ua
		JOIN unified_category_attributes uca ON ua.id = uca.attribute_id
		WHERE uca.category_id = 1
		ORDER BY uca.display_order
		LIMIT 100
	`
	
	rows, err := v.db.Query(query)
	if err == nil {
		defer rows.Close()
	}
	
	elapsed := time.Since(start).Milliseconds()
	
	v.addCheck(Check{
		Name:     "Query Performance",
		Category: "performance",
		Passed:   err == nil && elapsed < 50,
		Value:    fmt.Sprintf("%dms", elapsed),
		Expected: "<50ms",
		Message:  "Attribute query performance",
		Severity: "warning",
	})
}

// Prometheus metrics checks
func (v *Validator) checkErrorRate() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := `sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		v.addCheck(Check{
			Name:     "Error Rate",
			Category: "metrics",
			Passed:   false,
			Message:  fmt.Sprintf("Failed to query metrics: %v", err),
			Severity: "critical",
		})
		return
	}
	
	errorRate := v.extractValue(result)
	
	v.addCheck(Check{
		Name:     "Error Rate",
		Category: "metrics",
		Passed:   errorRate < 0.001,
		Value:    fmt.Sprintf("%.4f%%", errorRate*100),
		Expected: "<0.1%",
		Message:  "HTTP error rate",
		Severity: "critical",
	})
}

func (v *Validator) checkLatency() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		v.addCheck(Check{
			Name:     "P95 Latency",
			Category: "metrics",
			Passed:   false,
			Message:  fmt.Sprintf("Failed to query metrics: %v", err),
			Severity: "critical",
		})
		return
	}
	
	latency := v.extractValue(result) * 1000 // Convert to ms
	
	v.addCheck(Check{
		Name:     "P95 Latency",
		Category: "metrics",
		Passed:   latency < 50,
		Value:    fmt.Sprintf("%.2fms", latency),
		Expected: "<50ms",
		Message:  "95th percentile latency",
		Severity: "critical",
	})
}

func (v *Validator) checkDualWriteConsistency() {
	query := `
		SELECT 
			COUNT(DISTINCT ua.id) as unified_count,
			COUNT(DISTINCT ca.id) as legacy_count
		FROM unified_attributes ua
		FULL OUTER JOIN category_attributes ca ON ua.legacy_id = ca.id
		WHERE ua.created_at > NOW() - INTERVAL '1 hour'
			OR ca.created_at > NOW() - INTERVAL '1 hour'
	`
	
	var unifiedCount, legacyCount int
	err := v.db.QueryRow(query).Scan(&unifiedCount, &legacyCount)
	
	v.addCheck(Check{
		Name:     "Dual-Write Consistency",
		Category: "data",
		Passed:   err == nil && unifiedCount == legacyCount,
		Value:    fmt.Sprintf("Unified: %d, Legacy: %d", unifiedCount, legacyCount),
		Expected: "Equal counts",
		Message:  "Data consistency between old and new systems",
		Severity: "critical",
	})
}

func (v *Validator) checkTrafficDistribution() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := `sum(rate(http_requests_total{version="green"}[5m])) / sum(rate(http_requests_total[5m])) * 100`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		v.addCheck(Check{
			Name:     "Traffic Distribution",
			Category: "deployment",
			Passed:   false,
			Message:  fmt.Sprintf("Failed to query metrics: %v", err),
			Severity: "warning",
		})
		return
	}
	
	actualPercent := v.extractValue(result)
	expectedPercent := float64(v.config.CanaryPercent)
	
	// Allow 5% deviation
	passed := actualPercent >= expectedPercent-5 && actualPercent <= expectedPercent+5
	
	v.addCheck(Check{
		Name:     "Canary Traffic Distribution",
		Category: "deployment",
		Passed:   passed,
		Value:    fmt.Sprintf("%.1f%%", actualPercent),
		Expected: fmt.Sprintf("%d%% ±5%%", v.config.CanaryPercent),
		Message:  "Traffic routed to new version",
		Severity: "warning",
	})
}

// Helper functions
func (v *Validator) addCheck(check Check) {
	v.results = append(v.results, check)
	
	// Log the check
	status := "✓"
	if !check.Passed {
		status = "✗"
	}
	log.Printf("%s %s: %v (expected: %v)\n", status, check.Name, check.Value, check.Expected)
}

func (v *Validator) extractValue(result model.Value) float64 {
	switch result.Type() {
	case model.ValVector:
		vector := result.(model.Vector)
		if len(vector) > 0 {
			return float64(vector[0].Value)
		}
	case model.ValScalar:
		scalar := result.(*model.Scalar)
		return float64(scalar.Value)
	}
	return 0
}

func (v *Validator) generateResult(mode string) ValidationResult {
	passed := true
	criticalFailed := false
	
	for _, check := range v.results {
		if !check.Passed {
			if check.Severity == "critical" {
				criticalFailed = true
			}
			if check.Severity != "info" {
				passed = false
			}
		}
	}
	
	// Generate recommendations
	recommendations := v.generateRecommendations(criticalFailed)
	
	result := ValidationResult{
		Timestamp:       time.Now(),
		Mode:           mode,
		Passed:         passed && !criticalFailed,
		Checks:         v.results,
		Recommendations: recommendations,
		Metrics: map[string]interface{}{
			"total_checks":    len(v.results),
			"passed_checks":   v.countPassed(),
			"critical_failed": criticalFailed,
		},
	}
	
	// Send alert if critical failure
	if criticalFailed {
		v.sendAlert("CRITICAL", fmt.Sprintf("Validation failed for %s", mode))
	}
	
	return result
}

func (v *Validator) generateRecommendations(criticalFailed bool) []string {
	recommendations := []string{}
	
	if criticalFailed {
		recommendations = append(recommendations, "⚠️ CRITICAL: Rollback recommended due to critical failures")
	}
	
	// Analyze results and provide specific recommendations
	for _, check := range v.results {
		if !check.Passed {
			switch check.Name {
			case "Error Rate":
				recommendations = append(recommendations, "Investigate error logs and consider reducing canary traffic")
			case "P95 Latency":
				recommendations = append(recommendations, "Check database query performance and cache hit rates")
			case "Dual-Write Consistency":
				recommendations = append(recommendations, "Verify dual-write mechanism and check for race conditions")
			case "Memory Usage":
				recommendations = append(recommendations, "Consider scaling horizontally or investigating memory leaks")
			}
		}
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ All checks passed - safe to proceed")
	}
	
	return recommendations
}

func (v *Validator) countPassed() int {
	count := 0
	for _, check := range v.results {
		if check.Passed {
			count++
		}
	}
	return count
}

func (v *Validator) sendAlert(severity, message string) {
	if v.config.SlackWebhook == "" {
		return
	}
	
	payload := map[string]string{
		"text": fmt.Sprintf("*Unified Attributes Validation*\n*Severity:* %s\n*Message:* %s", severity, message),
	}
	
	jsonPayload, _ := json.Marshal(payload)
	http.Post(v.config.SlackWebhook, "application/json", bytes.NewBuffer(jsonPayload))
}

func (v *Validator) checkBackupStatus() {
	// Check if recent backup exists
	query := `
		SELECT COUNT(*) FROM pg_stat_activity
		WHERE query LIKE '%pg_dump%'
		AND state = 'active'
	`
	
	var activeBackups int
	v.db.QueryRow(query).Scan(&activeBackups)
	
	// Also check backup files (simplified check)
	backupDir := "/var/backups/unified-attributes"
	entries, err := os.ReadDir(backupDir)
	recentBackup := false
	
	if err == nil {
		for _, entry := range entries {
			info, _ := entry.Info()
			if time.Since(info.ModTime()) < 24*time.Hour {
				recentBackup = true
				break
			}
		}
	}
	
	v.addCheck(Check{
		Name:     "Backup Status",
		Category: "safety",
		Passed:   recentBackup,
		Value:    fmt.Sprintf("Recent backup: %v", recentBackup),
		Expected: "Backup < 24h old",
		Message:  "Recent backup available for rollback",
		Severity: "critical",
	})
}

func (v *Validator) checkApplicationReadiness() {
	// Check if application endpoints are responding
	endpoints := []string{
		"/health/ready",
		"/api/v2/categories/1/attributes",
	}
	
	allReady := true
	for _, endpoint := range endpoints {
		resp, err := v.client.Get(v.config.APIBaseURL + endpoint)
		if err != nil || resp.StatusCode != 200 {
			allReady = false
			break
		}
		resp.Body.Close()
	}
	
	v.addCheck(Check{
		Name:     "Application Readiness",
		Category: "application",
		Passed:   allReady,
		Value:    allReady,
		Expected: true,
		Message:  "All application endpoints ready",
		Severity: "critical",
	})
}

func (v *Validator) checkCachePerformance() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := `sum(rate(cache_hits_total[5m])) / sum(rate(cache_requests_total[5m]))`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		return
	}
	
	hitRate := v.extractValue(result)
	
	v.addCheck(Check{
		Name:     "Cache Hit Rate",
		Category: "performance",
		Passed:   hitRate > 0.8,
		Value:    fmt.Sprintf("%.1f%%", hitRate*100),
		Expected: ">80%",
		Message:  "Cache effectiveness",
		Severity: "warning",
	})
}

func (v *Validator) checkResourceUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Check memory usage
	memQuery := `avg(container_memory_usage_bytes{pod=~"backend-.*"}) / avg(container_spec_memory_limit_bytes{pod=~"backend-.*"})`
	memResult, _, _ := v.promAPI.Query(ctx, memQuery, time.Now())
	memUsage := v.extractValue(memResult)
	
	v.addCheck(Check{
		Name:     "Memory Usage",
		Category: "resources",
		Passed:   memUsage < 0.9,
		Value:    fmt.Sprintf("%.1f%%", memUsage*100),
		Expected: "<90%",
		Message:  "Container memory usage",
		Severity: "warning",
	})
	
	// Check CPU usage
	cpuQuery := `avg(rate(container_cpu_usage_seconds_total{pod=~"backend-.*"}[5m]))`
	cpuResult, _, _ := v.promAPI.Query(ctx, cpuQuery, time.Now())
	cpuUsage := v.extractValue(cpuResult)
	
	v.addCheck(Check{
		Name:     "CPU Usage",
		Category: "resources",
		Passed:   cpuUsage < 0.8,
		Value:    fmt.Sprintf("%.1f%%", cpuUsage*100),
		Expected: "<80%",
		Message:  "Container CPU usage",
		Severity: "warning",
	})
}

func (v *Validator) checkFallbackTriggers() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	query := `sum(rate(fallback_triggered_total[5m]))`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		return
	}
	
	fallbackRate := v.extractValue(result)
	
	v.addCheck(Check{
		Name:     "Fallback Triggers",
		Category: "reliability",
		Passed:   fallbackRate < 0.01,
		Value:    fmt.Sprintf("%.4f/s", fallbackRate),
		Expected: "<0.01/s",
		Message:  "Fallback mechanism triggers",
		Severity: "warning",
	})
}

func (v *Validator) checkFinalDataIntegrity() {
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE ua.id IS NULL) as missing_unified,
			COUNT(*) FILTER (WHERE ca.id IS NULL) as missing_legacy
		FROM category_attributes ca
		FULL OUTER JOIN unified_attributes ua ON ua.legacy_id = ca.id
	`
	
	var missingUnified, missingLegacy int
	err := v.db.QueryRow(query).Scan(&missingUnified, &missingLegacy)
	
	v.addCheck(Check{
		Name:     "Final Data Integrity",
		Category: "data",
		Passed:   err == nil && missingUnified == 0,
		Value:    fmt.Sprintf("Missing: %d", missingUnified),
		Expected: 0,
		Message:  "All data successfully migrated",
		Severity: "critical",
	})
}

func (v *Validator) checkSystemStability() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Check error rate over last hour
	query := `sum(increase(http_requests_total{status=~"5.."}[1h]))`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		return
	}
	
	errors := v.extractValue(result)
	
	v.addCheck(Check{
		Name:     "System Stability (1h)",
		Category: "stability",
		Passed:   errors < 100,
		Value:    fmt.Sprintf("%.0f errors", errors),
		Expected: "<100 errors/hour",
		Message:  "System stability over past hour",
		Severity: "warning",
	})
}

func (v *Validator) checkOldSystemUsage() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Check if old endpoints are still being called
	query := `sum(rate(http_requests_total{endpoint=~"/api/v1/.*attributes.*"}[5m]))`
	result, _, err := v.promAPI.Query(ctx, query, time.Now())
	
	if err != nil {
		return
	}
	
	oldAPIUsage := v.extractValue(result)
	
	v.addCheck(Check{
		Name:     "Old API Usage",
		Category: "migration",
		Passed:   oldAPIUsage < 0.1,
		Value:    fmt.Sprintf("%.2f req/s", oldAPIUsage),
		Expected: "<0.1 req/s",
		Message:  "Old API endpoints usage",
		Severity: "info",
	})
}

func (v *Validator) comparePerformance() {
	// Compare performance between old and new system
	// This is a simplified comparison - in reality would be more complex
	
	oldQuery := `
		SELECT AVG(execution_time) FROM (
			SELECT extract(epoch from (clock_timestamp() - statement_timestamp())) * 1000 as execution_time
			FROM category_attributes ca
			JOIN marketplace_category_attributes mca ON ca.id = mca.attribute_id
			WHERE mca.category_id = 1
			LIMIT 100
		) t
	`
	
	newQuery := `
		SELECT AVG(execution_time) FROM (
			SELECT extract(epoch from (clock_timestamp() - statement_timestamp())) * 1000 as execution_time
			FROM unified_attributes ua
			JOIN unified_category_attributes uca ON ua.id = uca.attribute_id
			WHERE uca.category_id = 1
			LIMIT 100
		) t
	`
	
	var oldTime, newTime float64
	v.db.QueryRow(oldQuery).Scan(&oldTime)
	v.db.QueryRow(newQuery).Scan(&newTime)
	
	improvement := ((oldTime - newTime) / oldTime) * 100
	
	v.addCheck(Check{
		Name:     "Performance Improvement",
		Category: "performance",
		Passed:   newTime < oldTime,
		Value:    fmt.Sprintf("%.1f%% faster", improvement),
		Expected: "Faster than old system",
		Message:  "Performance comparison",
		Severity: "info",
	})
}

// Main function
func main() {
	config := Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		PrometheusURL: os.Getenv("PROMETHEUS_URL"),
		APIBaseURL:    os.Getenv("API_BASE_URL"),
		SlackWebhook:  os.Getenv("SLACK_WEBHOOK"),
		CanaryPercent: 0,
	}
	
	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: validator <mode> [canary_percent]")
		fmt.Println("Modes: pre-deployment, during-rollout, post-deployment")
		os.Exit(1)
	}
	
	mode := os.Args[1]
	if mode == "during-rollout" && len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &config.CanaryPercent)
	}
	
	// Create validator
	validator, err := NewValidator(config)
	if err != nil {
		log.Fatal(err)
	}
	defer validator.db.Close()
	
	// Run validation based on mode
	var result ValidationResult
	
	switch mode {
	case "pre-deployment":
		result = validator.ValidatePreDeployment()
	case "during-rollout":
		result = validator.ValidateDuringRollout()
	case "post-deployment":
		result = validator.ValidatePostDeployment()
	default:
		log.Fatal("Invalid mode. Use: pre-deployment, during-rollout, or post-deployment")
	}
	
	// Output results
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))
	
	// Save to file
	outputFile := fmt.Sprintf("/tmp/validation-result-%s-%d.json", mode, time.Now().Unix())
	os.WriteFile(outputFile, jsonResult, 0644)
	log.Printf("Results saved to %s\n", outputFile)
	
	// Exit with appropriate code
	if !result.Passed {
		os.Exit(1)
	}
}