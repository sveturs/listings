package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// API метрики
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
	)

	// Business метрики для unified attributes
	unifiedAttributesUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "unified_attributes_usage_total",
			Help: "Usage count of unified attributes API",
		},
		[]string{"version", "operation"},
	)

	featureFlagStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "feature_flag_status",
			Help: "Status of feature flags (1=enabled, 0=disabled)",
		},
		[]string{"flag_name"},
	)

	dualWriteOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dual_write_operations_total",
			Help: "Total number of dual write operations",
		},
		[]string{"status"},
	)

	cacheOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "result"},
	)

	// System метрики
	databaseConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_connections_active",
			Help: "Number of active database connections",
		},
	)

	databaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	redisOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "result"},
	)
)

// PrometheusMiddleware collects metrics for all HTTP requests
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip metrics endpoint itself
		if c.Path() == "/metrics" {
			return c.Next()
		}

		// Increment in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		status := strconv.Itoa(c.Response().StatusCode())

		// Get normalized endpoint (route pattern instead of actual path)
		// This prevents high cardinality labels (e.g. /listings/123 -> /listings/:id)
		endpoint := c.Route().Path
		if endpoint == "" {
			// Fallback to actual path if route is not found (e.g. 404 requests)
			endpoint = c.Path()
		}

		// Record metrics with normalized endpoint
		httpRequestsTotal.WithLabelValues(c.Method(), endpoint, status).Inc()
		httpRequestDuration.WithLabelValues(c.Method(), endpoint).Observe(duration)

		return err
	}
}

// RecordUnifiedAttributesUsage records usage of unified attributes API
func RecordUnifiedAttributesUsage(version, operation string) {
	unifiedAttributesUsage.WithLabelValues(version, operation).Inc()
}

// UpdateFeatureFlagStatus updates the status of a feature flag
func UpdateFeatureFlagStatus(flagName string, enabled bool) {
	value := 0.0
	if enabled {
		value = 1.0
	}
	featureFlagStatus.WithLabelValues(flagName).Set(value)
}

// RecordDualWriteOperation records a dual write operation
func RecordDualWriteOperation(success bool) {
	status := "failure"
	if success {
		status = "success"
	}
	dualWriteOperations.WithLabelValues(status).Inc()
}

// RecordCacheOperation records a cache operation
func RecordCacheOperation(operation string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	cacheOperations.WithLabelValues(operation, result).Inc()
}

// UpdateDatabaseConnections updates the number of active database connections
func UpdateDatabaseConnections(count int) {
	databaseConnections.Set(float64(count))
}

// RecordDatabaseQuery records database query duration
func RecordDatabaseQuery(queryType string, duration time.Duration) {
	databaseQueryDuration.WithLabelValues(queryType).Observe(duration.Seconds())
}

// RecordRedisOperation records a Redis operation
func RecordRedisOperation(operation string, success bool) {
	result := "failure"
	if success {
		result = "success"
	}
	redisOperations.WithLabelValues(operation, result).Inc()
}

// InitializeFeatureFlagMetrics initializes metrics for feature flags
func InitializeFeatureFlagMetrics(flags map[string]bool) {
	for name, enabled := range flags {
		UpdateFeatureFlagStatus(name, enabled)
	}
}
