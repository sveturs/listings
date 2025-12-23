package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	opensearchRepo "github.com/vondi-global/listings/internal/repository/opensearch"
)

// OpenSearchMonitoringHandler provides OpenSearch monitoring HTTP endpoints
type OpenSearchMonitoringHandler struct {
	searchClient *opensearchRepo.Client
	logger       zerolog.Logger
}

// NewOpenSearchMonitoringHandler creates a new OpenSearch monitoring handler
func NewOpenSearchMonitoringHandler(searchClient *opensearchRepo.Client, logger zerolog.Logger) *OpenSearchMonitoringHandler {
	return &OpenSearchMonitoringHandler{
		searchClient: searchClient,
		logger:       logger.With().Str("component", "opensearch_monitoring_handler").Logger(),
	}
}

// RegisterRoutes registers OpenSearch monitoring routes
func (h *OpenSearchMonitoringHandler) RegisterRoutes(app *fiber.App) {
	if h.searchClient == nil {
		h.logger.Warn().Msg("OpenSearch client not available, monitoring routes disabled")
		return
	}

	// OpenSearch-specific health and metrics endpoints
	app.Get("/health/opensearch", h.HealthCheck)
	app.Get("/metrics/opensearch", h.IndexStats)
	app.Get("/metrics/opensearch/cluster", h.ClusterHealth)

	h.logger.Info().Msg("OpenSearch monitoring routes registered")
}

// HealthCheck performs OpenSearch health check
// GET /health/opensearch
func (h *OpenSearchMonitoringHandler) HealthCheck(c *fiber.Ctx) error {
	if h.searchClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unavailable",
			"error":  "OpenSearch client not initialized",
		})
	}

	ctx := c.Context()
	health, err := h.searchClient.HealthCheckDetailed(ctx)

	if err != nil {
		statusCode := fiber.StatusServiceUnavailable
		if health != nil && health.Status == "degraded" {
			statusCode = fiber.StatusMultipleChoices // 207 Multi-Status
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":         health.Status,
			"cluster_health": health.ClusterHealth,
			"index_exists":   health.IndexExists,
			"docs_count":     health.DocsCount,
			"error":          err.Error(),
			"timestamp":      time.Now().Unix(),
		})
	}

	// Determine HTTP status code
	statusCode := fiber.StatusOK
	if health.Status == "unhealthy" {
		statusCode = fiber.StatusServiceUnavailable
	} else if health.Status == "degraded" {
		statusCode = fiber.StatusMultipleChoices // 207 Multi-Status
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":         health.Status,
		"cluster_health": health.ClusterHealth,
		"index_exists":   health.IndexExists,
		"docs_count":     health.DocsCount,
		"index_size_mb":  health.IndexSizeMB,
		"latency_ms":     health.Latency.Milliseconds(),
		"shards": fiber.Map{
			"total":      health.Shards.Total,
			"successful": health.Shards.Successful,
			"failed":     health.Shards.Failed,
			"primary":    health.Shards.Primary,
			"replica":    health.Shards.Replica,
		},
		"last_check": health.LastCheck.Format(time.RFC3339),
		"timestamp":  time.Now().Unix(),
	})
}

// IndexStats returns detailed index statistics
// GET /metrics/opensearch
func (h *OpenSearchMonitoringHandler) IndexStats(c *fiber.Ctx) error {
	if h.searchClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "OpenSearch client not initialized",
		})
	}

	ctx := c.Context()
	stats, err := h.searchClient.GetIndexStats(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get index stats")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"index":         stats.Index,
		"docs_count":    stats.DocsCount,
		"store_size":    stats.StoreSize,
		"store_size_mb": stats.StoreSizeMB,
		"segment_count": stats.SegmentCount,
		"shards": fiber.Map{
			"total":      stats.Shards.Total,
			"successful": stats.Shards.Successful,
			"failed":     stats.Shards.Failed,
			"primary":    stats.Shards.Primary,
			"replica":    stats.Shards.Replica,
		},
		"refresh_time_ms": stats.RefreshTime.Milliseconds(),
		"flush_time_ms":   stats.FlushTime.Milliseconds(),
		"timestamp":       time.Now().Unix(),
	})
}

// ClusterHealth returns cluster health information
// GET /metrics/opensearch/cluster
func (h *OpenSearchMonitoringHandler) ClusterHealth(c *fiber.Ctx) error {
	if h.searchClient == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "OpenSearch client not initialized",
		})
	}

	// Redirect to detailed health check
	// This endpoint provides cluster-wide health status
	health, err := h.searchClient.HealthCheckDetailed(c.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("cluster health check failed")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"cluster_status": health.ClusterHealth,
		"index_exists":   health.IndexExists,
		"overall_status": health.Status,
		"timestamp":      time.Now().Unix(),
	})
}
