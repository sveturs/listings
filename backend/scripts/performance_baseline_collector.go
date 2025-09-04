package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

// MetricType represents the type of metric being collected
type MetricType string

const (
	MetricResponseTime MetricType = "response_time"
	MetricThroughput   MetricType = "throughput"
	MetricErrorRate    MetricType = "error_rate"
	MetricCacheHit     MetricType = "cache_hit"
	MetricMemoryUsage  MetricType = "memory_usage"
	MetricCPUUsage     MetricType = "cpu_usage"
)

// PerformanceBaseline represents baseline performance metrics
type PerformanceBaseline struct {
	Timestamp   time.Time                `json:"timestamp"`
	Environment string                   `json:"environment"`
	Version     string                   `json:"version"`
	Metrics     map[string]MetricStats   `json:"metrics"`
	Endpoints   map[string]EndpointStats `json:"endpoints"`
	Anomalies   []Anomaly                `json:"anomalies,omitempty"`
	HealthScore float64                  `json:"health_score"`
}

// MetricStats contains statistical information about a metric
type MetricStats struct {
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	P95    float64 `json:"p95"`
	P99    float64 `json:"p99"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	StdDev float64 `json:"std_dev"`
	Count  int     `json:"count"`
}

// EndpointStats contains performance stats for a specific endpoint
type EndpointStats struct {
	URL          string      `json:"url"`
	Method       string      `json:"method"`
	ResponseTime MetricStats `json:"response_time"`
	SuccessRate  float64     `json:"success_rate"`
	ErrorCount   int         `json:"error_count"`
	SampleSize   int         `json:"sample_size"`
}

// Anomaly represents a detected performance anomaly
type Anomaly struct {
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Expected    float64   `json:"expected"`
	Deviation   float64   `json:"deviation"`
	Severity    string    `json:"severity"`
	DetectedAt  time.Time `json:"detected_at"`
	Description string    `json:"description"`
}

// BaselineCollector collects and analyzes performance metrics
type BaselineCollector struct {
	config           *Config
	httpClient       *http.Client
	prometheusURL    string
	baseline         *PerformanceBaseline
	previousBaseline *PerformanceBaseline
	mu               sync.RWMutex
	alertChan        chan Anomaly
}

// Config holds collector configuration
type Config struct {
	APIBaseURL         string        `json:"api_base_url"`
	PrometheusURL      string        `json:"prometheus_url"`
	CollectionInterval time.Duration `json:"collection_interval"`
	SampleSize         int           `json:"sample_size"`
	WarmupRequests     int           `json:"warmup_requests"`
	BaselineFile       string        `json:"baseline_file"`
	ReportDir          string        `json:"report_dir"`
	AlertThresholds    AlertConfig   `json:"alert_thresholds"`
}

// AlertConfig defines thresholds for alerts
type AlertConfig struct {
	ResponseTimeP95MS float64 `json:"response_time_p95_ms"`
	ErrorRatePercent  float64 `json:"error_rate_percent"`
	CacheHitPercent   float64 `json:"cache_hit_percent"`
	MemoryUsageMB     float64 `json:"memory_usage_mb"`
	CPUUsagePercent   float64 `json:"cpu_usage_percent"`
}

// NewCollector creates a new baseline collector
func NewCollector(config *Config) *BaselineCollector {
	return &BaselineCollector{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		prometheusURL: config.PrometheusURL,
		alertChan:     make(chan Anomaly, 100),
	}
}

// CollectEndpointMetrics collects metrics for a specific endpoint
func (c *BaselineCollector) CollectEndpointMetrics(endpoint string, method string) EndpointStats {
	stats := EndpointStats{
		URL:        endpoint,
		Method:     method,
		SampleSize: c.config.SampleSize,
	}

	// Warmup requests
	for i := 0; i < c.config.WarmupRequests; i++ {
		req, _ := http.NewRequest(method, endpoint, nil)
		c.httpClient.Do(req)
		time.Sleep(100 * time.Millisecond)
	}

	// Collect samples
	var responseTimes []float64
	var successCount int

	for i := 0; i < c.config.SampleSize; i++ {
		start := time.Now()
		req, err := http.NewRequest(method, endpoint, nil)
		if err != nil {
			stats.ErrorCount++
			continue
		}

		resp, err := c.httpClient.Do(req)
		elapsed := time.Since(start).Seconds() * 1000 // Convert to milliseconds

		if err != nil {
			stats.ErrorCount++
		} else {
			responseTimes = append(responseTimes, elapsed)
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				successCount++
			} else {
				stats.ErrorCount++
			}
			resp.Body.Close()
		}

		time.Sleep(50 * time.Millisecond) // Small delay between requests
	}

	// Calculate statistics
	if len(responseTimes) > 0 {
		stats.ResponseTime = calculateStats(responseTimes)
		stats.SuccessRate = float64(successCount) / float64(c.config.SampleSize) * 100
	}

	return stats
}

// calculateStats calculates statistical metrics from a slice of values
func calculateStats(values []float64) MetricStats {
	if len(values) == 0 {
		return MetricStats{}
	}

	sort.Float64s(values)

	stats := MetricStats{
		Count: len(values),
		Min:   values[0],
		Max:   values[len(values)-1],
	}

	// Calculate mean
	var sum float64
	for _, v := range values {
		sum += v
	}
	stats.Mean = sum / float64(len(values))

	// Calculate median
	if len(values)%2 == 0 {
		stats.Median = (values[len(values)/2-1] + values[len(values)/2]) / 2
	} else {
		stats.Median = values[len(values)/2]
	}

	// Calculate percentiles
	stats.P95 = percentile(values, 95)
	stats.P99 = percentile(values, 99)

	// Calculate standard deviation
	var variance float64
	for _, v := range values {
		variance += math.Pow(v-stats.Mean, 2)
	}
	stats.StdDev = math.Sqrt(variance / float64(len(values)))

	return stats
}

// percentile calculates the nth percentile of a sorted slice
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}

	index := (p / 100) * float64(len(sorted)-1)
	lower := math.Floor(index)
	upper := math.Ceil(index)

	if lower == upper {
		return sorted[int(index)]
	}

	// Linear interpolation
	return sorted[int(lower)]*(upper-index) + sorted[int(upper)]*(index-lower)
}

// CollectPrometheusMetrics queries Prometheus for system metrics
func (c *BaselineCollector) CollectPrometheusMetrics() (map[string]MetricStats, error) {
	metrics := make(map[string]MetricStats)

	// Define queries
	queries := map[string]string{
		"error_rate":      `rate(http_requests_total{status=~"5.."}[5m])`,
		"request_rate":    `rate(http_requests_total[5m])`,
		"cache_hit_rate":  `rate(unified_attributes_cache_hits_total[5m])/(rate(unified_attributes_cache_hits_total[5m])+rate(unified_attributes_cache_misses_total[5m]))`,
		"memory_usage_mb": `container_memory_usage_bytes/1024/1024`,
		"cpu_usage":       `rate(container_cpu_usage_seconds_total[5m])*100`,
		"db_connections":  `pg_stat_database_numbackends`,
		"db_query_time":   `pg_stat_statements_mean_seconds*1000`,
	}

	for name, query := range queries {
		value, err := c.queryPrometheus(query)
		if err != nil {
			log.Printf("Error querying %s: %v", name, err)
			continue
		}

		// For single value metrics, create a simple stats object
		metrics[name] = MetricStats{
			Mean:   value,
			Median: value,
			P95:    value,
			P99:    value,
			Min:    value,
			Max:    value,
			Count:  1,
		}
	}

	return metrics, nil
}

// queryPrometheus executes a PromQL query
func (c *BaselineCollector) queryPrometheus(query string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/query?query=%s", c.prometheusURL, query)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Parse Prometheus response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	// Extract value (simplified - assumes single result)
	if data, ok := result["data"].(map[string]interface{}); ok {
		if results, ok := data["result"].([]interface{}); ok && len(results) > 0 {
			if res, ok := results[0].(map[string]interface{}); ok {
				if value, ok := res["value"].([]interface{}); ok && len(value) > 1 {
					if v, ok := value[1].(string); ok {
						var val float64
						fmt.Sscanf(v, "%f", &val)
						return val, nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("no data found")
}

// DetectAnomalies compares current metrics with baseline
func (c *BaselineCollector) DetectAnomalies() []Anomaly {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var anomalies []Anomaly

	if c.previousBaseline == nil {
		return anomalies
	}

	// Check response time
	for endpoint, stats := range c.baseline.Endpoints {
		if prevStats, exists := c.previousBaseline.Endpoints[endpoint]; exists {
			deviation := (stats.ResponseTime.P95 - prevStats.ResponseTime.P95) / prevStats.ResponseTime.P95 * 100

			if math.Abs(deviation) > 20 { // 20% deviation threshold
				anomalies = append(anomalies, Anomaly{
					Metric:      fmt.Sprintf("%s_response_time_p95", endpoint),
					Value:       stats.ResponseTime.P95,
					Expected:    prevStats.ResponseTime.P95,
					Deviation:   deviation,
					Severity:    getSeverity(deviation),
					DetectedAt:  time.Now(),
					Description: fmt.Sprintf("Response time deviation for %s: %.2f%%", endpoint, deviation),
				})
			}
		}
	}

	// Check error rate
	if errorRate, exists := c.baseline.Metrics["error_rate"]; exists {
		if errorRate.Mean > c.config.AlertThresholds.ErrorRatePercent {
			anomalies = append(anomalies, Anomaly{
				Metric:      "error_rate",
				Value:       errorRate.Mean,
				Expected:    c.config.AlertThresholds.ErrorRatePercent,
				Deviation:   errorRate.Mean - c.config.AlertThresholds.ErrorRatePercent,
				Severity:    "critical",
				DetectedAt:  time.Now(),
				Description: fmt.Sprintf("Error rate %.2f%% exceeds threshold %.2f%%", errorRate.Mean, c.config.AlertThresholds.ErrorRatePercent),
			})
		}
	}

	return anomalies
}

// getSeverity determines anomaly severity based on deviation
func getSeverity(deviation float64) string {
	absDeviation := math.Abs(deviation)
	switch {
	case absDeviation > 50:
		return "critical"
	case absDeviation > 30:
		return "warning"
	default:
		return "info"
	}
}

// CalculateHealthScore calculates overall system health score (0-100)
func (c *BaselineCollector) CalculateHealthScore() float64 {
	score := 100.0

	// Deduct points for anomalies
	for _, anomaly := range c.baseline.Anomalies {
		switch anomaly.Severity {
		case "critical":
			score -= 20
		case "warning":
			score -= 10
		case "info":
			score -= 5
		}
	}

	// Deduct points for poor metrics
	if errorRate, exists := c.baseline.Metrics["error_rate"]; exists && errorRate.Mean > 0.1 {
		score -= 15
	}

	if cacheHit, exists := c.baseline.Metrics["cache_hit_rate"]; exists && cacheHit.Mean < 70 {
		score -= 10
	}

	// Ensure score stays within bounds
	if score < 0 {
		score = 0
	}

	return score
}

// CollectBaseline performs a full baseline collection
func (c *BaselineCollector) CollectBaseline() error {
	log.Println("Starting baseline collection...")

	baseline := &PerformanceBaseline{
		Timestamp:   time.Now(),
		Environment: os.Getenv("ENVIRONMENT"),
		Version:     os.Getenv("VERSION"),
		Metrics:     make(map[string]MetricStats),
		Endpoints:   make(map[string]EndpointStats),
	}

	// Define endpoints to test
	endpoints := []struct {
		URL    string
		Method string
	}{
		{c.config.APIBaseURL + "/api/v1/unified-attributes", "GET"},
		{c.config.APIBaseURL + "/api/v1/marketplace/search", "GET"},
		{c.config.APIBaseURL + "/api/v1/categories", "GET"},
		{c.config.APIBaseURL + "/health", "GET"},
	}

	// Collect endpoint metrics
	var wg sync.WaitGroup
	for _, ep := range endpoints {
		wg.Add(1)
		go func(endpoint, method string) {
			defer wg.Done()
			stats := c.CollectEndpointMetrics(endpoint, method)
			c.mu.Lock()
			baseline.Endpoints[endpoint] = stats
			c.mu.Unlock()
		}(ep.URL, ep.Method)
	}
	wg.Wait()

	// Collect system metrics from Prometheus
	systemMetrics, err := c.CollectPrometheusMetrics()
	if err != nil {
		log.Printf("Error collecting Prometheus metrics: %v", err)
	} else {
		baseline.Metrics = systemMetrics
	}

	// Update baseline
	c.mu.Lock()
	c.previousBaseline = c.baseline
	c.baseline = baseline
	c.mu.Unlock()

	// Detect anomalies
	baseline.Anomalies = c.DetectAnomalies()

	// Calculate health score
	baseline.HealthScore = c.CalculateHealthScore()

	// Save baseline
	if err := c.SaveBaseline(baseline); err != nil {
		return fmt.Errorf("error saving baseline: %w", err)
	}

	log.Printf("Baseline collected successfully. Health score: %.2f", baseline.HealthScore)
	return nil
}

// SaveBaseline saves baseline to file
func (c *BaselineCollector) SaveBaseline(baseline *PerformanceBaseline) error {
	// Create report directory if it doesn't exist
	if err := os.MkdirAll(c.config.ReportDir, 0o755); err != nil {
		return err
	}

	// Save current baseline
	filename := fmt.Sprintf("%s/baseline_%s.json", c.config.ReportDir, time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(baseline); err != nil {
		return err
	}

	// Also save as latest baseline
	latestFile := fmt.Sprintf("%s/latest_baseline.json", c.config.ReportDir)
	if err := copyFile(filename, latestFile); err != nil {
		return err
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// MonitorContinuously runs continuous monitoring
func (c *BaselineCollector) MonitorContinuously(ctx context.Context) {
	ticker := time.NewTicker(c.config.CollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Monitoring stopped")
			return
		case <-ticker.C:
			if err := c.CollectBaseline(); err != nil {
				log.Printf("Error collecting baseline: %v", err)
			}

			// Process alerts
			c.ProcessAlerts()
		}
	}
}

// ProcessAlerts handles detected anomalies
func (c *BaselineCollector) ProcessAlerts() {
	for _, anomaly := range c.baseline.Anomalies {
		select {
		case c.alertChan <- anomaly:
			log.Printf("Alert: %s - %s", anomaly.Severity, anomaly.Description)
		default:
			// Alert channel full, log to file
			log.Printf("Alert channel full, dropping alert: %s", anomaly.Description)
		}
	}
}

// GenerateReport generates a markdown report
func (c *BaselineCollector) GenerateReport() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.baseline == nil {
		return fmt.Errorf("no baseline data available")
	}

	filename := fmt.Sprintf("%s/performance_report_%s.md", c.config.ReportDir, time.Now().Format("20060102"))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "# Performance Baseline Report\n\n")
	fmt.Fprintf(file, "**Generated:** %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "**Environment:** %s\n", c.baseline.Environment)
	fmt.Fprintf(file, "**Version:** %s\n", c.baseline.Version)
	fmt.Fprintf(file, "**Health Score:** %.2f/100\n\n", c.baseline.HealthScore)

	fmt.Fprintf(file, "## Endpoint Performance\n\n")
	fmt.Fprintf(file, "| Endpoint | P95 (ms) | P99 (ms) | Success Rate | Errors |\n")
	fmt.Fprintf(file, "|----------|----------|----------|--------------|--------|\n")

	for url, stats := range c.baseline.Endpoints {
		fmt.Fprintf(file, "| %s | %.2f | %.2f | %.2f%% | %d |\n",
			url, stats.ResponseTime.P95, stats.ResponseTime.P99, stats.SuccessRate, stats.ErrorCount)
	}

	fmt.Fprintf(file, "\n## System Metrics\n\n")
	fmt.Fprintf(file, "| Metric | Value |\n")
	fmt.Fprintf(file, "|--------|-------|\n")

	for name, stats := range c.baseline.Metrics {
		fmt.Fprintf(file, "| %s | %.2f |\n", name, stats.Mean)
	}

	if len(c.baseline.Anomalies) > 0 {
		fmt.Fprintf(file, "\n## Detected Anomalies\n\n")
		for _, anomaly := range c.baseline.Anomalies {
			fmt.Fprintf(file, "- **%s** [%s]: %s (deviation: %.2f%%)\n",
				anomaly.Metric, anomaly.Severity, anomaly.Description, anomaly.Deviation)
		}
	}

	log.Printf("Report generated: %s", filename)
	return nil
}

func main() {
	// Load configuration
	config := &Config{
		APIBaseURL:         os.Getenv("API_BASE_URL"),
		PrometheusURL:      os.Getenv("PROMETHEUS_URL"),
		CollectionInterval: 5 * time.Minute,
		SampleSize:         100,
		WarmupRequests:     5,
		BaselineFile:       "/var/reports/unified-attributes/baseline.json",
		ReportDir:          "/var/reports/unified-attributes",
		AlertThresholds: AlertConfig{
			ResponseTimeP95MS: 100,
			ErrorRatePercent:  0.1,
			CacheHitPercent:   70,
			MemoryUsageMB:     512,
			CPUUsagePercent:   80,
		},
	}

	// Set defaults if not provided
	if config.APIBaseURL == "" {
		config.APIBaseURL = "https://api.svetu.rs"
	}
	if config.PrometheusURL == "" {
		config.PrometheusURL = "https://prometheus.svetu.rs"
	}

	// Create collector
	collector := NewCollector(config)

	// Perform initial collection
	if err := collector.CollectBaseline(); err != nil {
		log.Fatalf("Error collecting initial baseline: %v", err)
	}

	// Generate initial report
	if err := collector.GenerateReport(); err != nil {
		log.Printf("Error generating report: %v", err)
	}

	// Start continuous monitoring
	ctx := context.Background()
	collector.MonitorContinuously(ctx)
}
