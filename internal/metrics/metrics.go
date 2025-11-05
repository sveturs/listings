package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics contains all Prometheus metrics for the listings service
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge

	// gRPC metrics
	GRPCRequestsTotal   *prometheus.CounterVec
	GRPCRequestDuration *prometheus.HistogramVec

	// Business metrics
	ListingsCreated  prometheus.Counter
	ListingsUpdated  prometheus.Counter
	ListingsDeleted  prometheus.Counter
	ListingsSearched prometheus.Counter

	// Database metrics
	DBConnectionsOpen prometheus.Gauge
	DBConnectionsIdle prometheus.Gauge
	DBQueryDuration   *prometheus.HistogramVec

	// Cache metrics
	CacheHits   *prometheus.CounterVec
	CacheMisses *prometheus.CounterVec

	// Indexing queue metrics
	IndexingQueueSize     prometheus.Gauge
	IndexingJobsProcessed *prometheus.CounterVec
	IndexingJobDuration   prometheus.Histogram

	// Error metrics
	ErrorsTotal *prometheus.CounterVec

	// Inventory-specific metrics
	InventoryProductViews          *prometheus.CounterVec
	InventoryProductViewsErrors    prometheus.Counter
	InventoryStockOperations       *prometheus.CounterVec
	InventoryStockLowThreshold     *prometheus.CounterVec
	InventoryMovementsRecorded     *prometheus.CounterVec
	InventoryMovementsErrors       *prometheus.CounterVec
	InventoryStockValue            *prometheus.GaugeVec
	InventoryOutOfStockProducts    prometheus.Gauge

	// gRPC handler metrics (granular)
	GRPCHandlerRequestsActive *prometheus.GaugeVec

	// Rate limiting metrics
	RateLimitHitsTotal     *prometheus.CounterVec
	RateLimitAllowedTotal  *prometheus.CounterVec
	RateLimitRejectedTotal *prometheus.CounterVec

	// Timeout metrics
	TimeoutsTotal     *prometheus.CounterVec
	NearTimeoutsTotal *prometheus.CounterVec
	TimeoutDuration   *prometheus.HistogramVec
}

// NewMetrics creates and registers all Prometheus metrics
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5},
			},
			[]string{"method", "path"},
		),
		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "http_requests_in_flight",
				Help:      "Current number of HTTP requests being processed",
			},
		),

		// gRPC metrics
		GRPCRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "grpc_requests_total",
				Help:      "Total number of gRPC requests",
			},
			[]string{"method", "status"},
		),
		GRPCRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "grpc_request_duration_seconds",
				Help:      "gRPC request latency in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
			},
			[]string{"method"},
		),

		// Business metrics
		ListingsCreated: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "listings_created_total",
				Help:      "Total number of listings created",
			},
		),
		ListingsUpdated: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "listings_updated_total",
				Help:      "Total number of listings updated",
			},
		),
		ListingsDeleted: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "listings_deleted_total",
				Help:      "Total number of listings deleted",
			},
		),
		ListingsSearched: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "listings_searched_total",
				Help:      "Total number of search queries executed",
			},
		),

		// Database metrics
		DBConnectionsOpen: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_open",
				Help:      "Current number of open database connections",
			},
		),
		DBConnectionsIdle: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connections_idle",
				Help:      "Current number of idle database connections",
			},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_query_duration_seconds",
				Help:      "Database query execution time in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
			},
			[]string{"operation"},
		),

		// Cache metrics
		CacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_hits_total",
				Help:      "Total number of cache hits",
			},
			[]string{"cache_type"},
		),
		CacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_misses_total",
				Help:      "Total number of cache misses",
			},
			[]string{"cache_type"},
		),

		// Indexing queue metrics
		IndexingQueueSize: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "indexing_queue_size",
				Help:      "Current size of the indexing queue",
			},
		),
		IndexingJobsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "indexing_jobs_processed_total",
				Help:      "Total number of indexing jobs processed",
			},
			[]string{"operation", "status"},
		),
		IndexingJobDuration: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "indexing_job_duration_seconds",
				Help:      "Indexing job processing time in seconds",
				Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
			},
		),

		// Error metrics
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "errors_total",
				Help:      "Total number of errors",
			},
			[]string{"component", "error_type"},
		),

		// Inventory-specific metrics
		InventoryProductViews: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "inventory_product_views_total",
				Help:      "Total number of product view increments",
			},
			[]string{"product_id"},
		),
		InventoryProductViewsErrors: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "inventory_product_views_errors_total",
				Help:      "Total number of product view increment errors",
			},
		),
		InventoryStockOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "inventory_stock_operations_total",
				Help:      "Total number of stock operations (update/batch)",
			},
			[]string{"operation", "status"},
		),
		InventoryStockLowThreshold: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "inventory_stock_low_threshold_reached_total",
				Help:      "Number of times stock fell below low threshold",
			},
			[]string{"product_id", "storefront_id"},
		),
		InventoryMovementsRecorded: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "inventory_movements_recorded_total",
				Help:      "Total number of inventory movements recorded",
			},
			[]string{"movement_type"},
		),
		InventoryMovementsErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "inventory_movements_errors_total",
				Help:      "Total number of inventory movement recording errors",
			},
			[]string{"reason"},
		),
		InventoryStockValue: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "inventory_stock_value",
				Help:      "Current stock value for products",
			},
			[]string{"storefront_id", "product_id"},
		),
		InventoryOutOfStockProducts: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "inventory_out_of_stock_products",
				Help:      "Current number of out-of-stock products",
			},
		),

		// gRPC handler metrics (granular)
		GRPCHandlerRequestsActive: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "grpc_handler_requests_active",
				Help:      "Current number of active gRPC requests",
			},
			[]string{"method"},
		),

		// Rate limiting metrics
		RateLimitHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rate_limit_hits_total",
				Help:      "Total number of rate limit evaluations",
			},
			[]string{"method", "identifier_type"},
		),
		RateLimitAllowedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rate_limit_allowed_total",
				Help:      "Total number of allowed requests (under rate limit)",
			},
			[]string{"method", "identifier_type"},
		),
		RateLimitRejectedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rate_limit_rejected_total",
				Help:      "Total number of rejected requests (rate limit exceeded)",
			},
			[]string{"method", "identifier_type"},
		),

		// Timeout metrics
		TimeoutsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "timeouts_total",
				Help:      "Total number of timed out requests",
			},
			[]string{"method"},
		),
		NearTimeoutsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "near_timeouts_total",
				Help:      "Total number of requests that approached timeout threshold (>80%)",
			},
			[]string{"method"},
		),
		TimeoutDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "timeout_duration_seconds",
				Help:      "Duration when timeout occurred",
				Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 15, 20, 30},
			},
			[]string{"method"},
		),
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, path, status string, duration float64) {
	m.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// RecordGRPCRequest records gRPC request metrics
func (m *Metrics) RecordGRPCRequest(method, status string, duration float64) {
	m.GRPCRequestsTotal.WithLabelValues(method, status).Inc()
	m.GRPCRequestDuration.WithLabelValues(method).Observe(duration)
}

// RecordCacheHit records a cache hit
func (m *Metrics) RecordCacheHit(cacheType string) {
	m.CacheHits.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss records a cache miss
func (m *Metrics) RecordCacheMiss(cacheType string) {
	m.CacheMisses.WithLabelValues(cacheType).Inc()
}

// RecordError records an error
func (m *Metrics) RecordError(component, errorType string) {
	m.ErrorsTotal.WithLabelValues(component, errorType).Inc()
}

// RecordIndexingJob records indexing job metrics
func (m *Metrics) RecordIndexingJob(operation, status string, duration float64) {
	m.IndexingJobsProcessed.WithLabelValues(operation, status).Inc()
	m.IndexingJobDuration.Observe(duration)
}

// UpdateDBConnectionStats updates database connection pool metrics
func (m *Metrics) UpdateDBConnectionStats(open, idle int) {
	m.DBConnectionsOpen.Set(float64(open))
	m.DBConnectionsIdle.Set(float64(idle))
}

// UpdateIndexingQueueSize updates indexing queue size metric
func (m *Metrics) UpdateIndexingQueueSize(size int) {
	m.IndexingQueueSize.Set(float64(size))
}

// RecordInventoryProductView records a product view increment
func (m *Metrics) RecordInventoryProductView(productID string) {
	m.InventoryProductViews.WithLabelValues(productID).Inc()
}

// RecordInventoryProductViewError records a product view increment error
func (m *Metrics) RecordInventoryProductViewError() {
	m.InventoryProductViewsErrors.Inc()
}

// RecordInventoryStockOperation records a stock operation (update/batch)
func (m *Metrics) RecordInventoryStockOperation(operation, status string) {
	m.InventoryStockOperations.WithLabelValues(operation, status).Inc()
}

// RecordInventoryLowStockThreshold records when stock falls below threshold
func (m *Metrics) RecordInventoryLowStockThreshold(productID, storefrontID string) {
	m.InventoryStockLowThreshold.WithLabelValues(productID, storefrontID).Inc()
}

// RecordInventoryMovement records an inventory movement
func (m *Metrics) RecordInventoryMovement(movementType string) {
	m.InventoryMovementsRecorded.WithLabelValues(movementType).Inc()
}

// RecordInventoryMovementError records an inventory movement error
func (m *Metrics) RecordInventoryMovementError(reason string) {
	m.InventoryMovementsErrors.WithLabelValues(reason).Inc()
}

// UpdateInventoryStockValue updates the stock value gauge for a product
func (m *Metrics) UpdateInventoryStockValue(storefrontID, productID string, value float64) {
	m.InventoryStockValue.WithLabelValues(storefrontID, productID).Set(value)
}

// UpdateInventoryOutOfStock updates the count of out-of-stock products
func (m *Metrics) UpdateInventoryOutOfStock(count int) {
	m.InventoryOutOfStockProducts.Set(float64(count))
}

// RecordRateLimitEvaluation records a rate limit evaluation
func (m *Metrics) RecordRateLimitEvaluation(method, identifierType string, allowed bool) {
	m.RateLimitHitsTotal.WithLabelValues(method, identifierType).Inc()
	if allowed {
		m.RateLimitAllowedTotal.WithLabelValues(method, identifierType).Inc()
	} else {
		m.RateLimitRejectedTotal.WithLabelValues(method, identifierType).Inc()
	}
}
