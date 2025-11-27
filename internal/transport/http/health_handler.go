package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/health"
)

// HealthHandler provides health check HTTP endpoints
type HealthHandler struct {
	checker health.Checker
	logger  zerolog.Logger
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(checker health.Checker, logger zerolog.Logger) *HealthHandler {
	return &HealthHandler{
		checker: checker,
		logger:  logger.With().Str("component", "health_handler").Logger(),
	}
}

// RegisterRoutes registers health check routes on the given Fiber app
func (h *HealthHandler) RegisterRoutes(app *fiber.App) {
	// Kubernetes-compatible health endpoints
	app.Get("/health", h.Health)          // Overall health
	app.Get("/health/live", h.Liveness)   // Liveness probe
	app.Get("/health/ready", h.Readiness) // Readiness probe
	app.Get("/health/startup", h.Startup) // Startup probe
	app.Get("/health/deep", h.DeepHealth) // Deep diagnostics

	// Backwards compatibility with existing /ready endpoint
	app.Get("/ready", h.Readiness)

	h.logger.Info().Msg("health check routes registered")
}

// Health returns overall health status with all dependency checks
// GET /health
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	ctx := c.Context()
	response := h.checker.CheckAll(ctx)

	// Convert response times to milliseconds for JSON
	checks := make(map[string]interface{})
	for name, result := range response.Checks {
		checks[name] = map[string]interface{}{
			"status":           result.Status,
			"response_time_ms": result.ResponseTime.Milliseconds(),
			"details":          result.Details,
		}
	}

	output := map[string]interface{}{
		"status":    response.Status,
		"version":   response.Version,
		"uptime":    response.Uptime,
		"checks":    checks,
		"timestamp": response.Timestamp.Format(time.RFC3339),
	}

	// Set HTTP status based on health status
	statusCode := fiber.StatusOK
	switch response.Status {
	case health.HealthStatusUnhealthy:
		statusCode = fiber.StatusServiceUnavailable
	case health.HealthStatusDegraded:
		statusCode = fiber.StatusOK // Still accepting traffic
	}

	return c.Status(statusCode).JSON(output)
}

// Liveness performs minimal liveness check (K8s liveness probe)
// GET /health/live
func (h *HealthHandler) Liveness(c *fiber.Ctx) error {
	ctx := c.Context()
	isAlive := h.checker.CheckLiveness(ctx)

	if !isAlive {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "unhealthy",
			"timestamp": time.Now().Unix(),
		})
	}

	return c.JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// Readiness performs readiness check (K8s readiness probe)
// GET /health/ready
func (h *HealthHandler) Readiness(c *fiber.Ctx) error {
	ctx := c.Context()
	isReady := h.checker.CheckReadiness(ctx)

	if !isReady {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "not_ready",
			"timestamp": time.Now().Unix(),
		})
	}

	return c.JSON(fiber.Map{
		"status":    "ready",
		"timestamp": time.Now().Unix(),
	})
}

// Startup performs startup probe check (K8s startup probe)
// GET /health/startup
func (h *HealthHandler) Startup(c *fiber.Ctx) error {
	ctx := c.Context()
	isStarted := h.checker.CheckStartup(ctx)

	if !isStarted {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":    "starting",
			"timestamp": time.Now().Unix(),
		})
	}

	return c.JSON(fiber.Map{
		"status":    "started",
		"timestamp": time.Now().Unix(),
	})
}

// DeepHealth performs extended diagnostics check
// GET /health/deep
func (h *HealthHandler) DeepHealth(c *fiber.Ctx) error {
	ctx := c.Context()
	response := h.checker.CheckDeep(ctx)

	// Convert response times to milliseconds for JSON
	checks := make(map[string]interface{})
	for name, result := range response.Checks {
		checks[name] = map[string]interface{}{
			"status":           result.Status,
			"response_time_ms": result.ResponseTime.Milliseconds(),
			"details":          result.Details,
		}
	}

	// Convert pool metrics wait time to milliseconds
	poolMetrics := make(map[string]interface{})
	for name, metrics := range response.Diagnostics.ConnectionPools {
		poolMetrics[name] = map[string]interface{}{
			"active":       metrics.Active,
			"idle":         metrics.Idle,
			"max_open":     metrics.MaxOpen,
			"wait_time_ms": metrics.WaitTime,
		}
	}

	output := map[string]interface{}{
		"status":    response.Status,
		"version":   response.Version,
		"uptime":    response.Uptime,
		"checks":    checks,
		"timestamp": response.Timestamp.Format(time.RFC3339),
		"diagnostics": map[string]interface{}{
			"goroutines":       response.Diagnostics.Goroutines,
			"memory_alloc_mb":  response.Diagnostics.MemoryAllocMB,
			"memory_sys_mb":    response.Diagnostics.MemorySysMB,
			"num_gc":           response.Diagnostics.NumGC,
			"recent_errors":    response.Diagnostics.RecentErrors,
			"connection_pools": poolMetrics,
		},
	}

	// Set HTTP status based on health status
	statusCode := fiber.StatusOK
	if response.Status == health.HealthStatusUnhealthy {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(output)
}
