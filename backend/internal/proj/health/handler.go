package health

import (
	"database/sql"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	db    *sql.DB
	redis *redis.Client
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// NewHandler creates a new health handler
func NewHandler(db *sql.DB, redis *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

// RegisterRoutes registers health check routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/health/live", h.LiveCheck)
	router.Get("/health/ready", h.ReadyCheck)
	router.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}

// LiveCheck checks if the service is alive
// @Summary Liveness probe
// @Description Check if the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} HealthStatus
// @Router /health/live [get]
func (h *Handler) LiveCheck(c *fiber.Ctx) error {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return c.JSON(status)
}

// ReadyCheck checks if the service is ready to handle requests
// @Summary Readiness probe
// @Description Check if the service is ready to handle requests
// @Tags health
// @Produce json
// @Success 200 {object} HealthStatus
// @Success 503 {object} HealthStatus
// @Router /health/ready [get]
func (h *Handler) ReadyCheck(c *fiber.Ctx) error {
	checks := make(map[string]string)
	isReady := true

	// Check database connection
	if h.db != nil {
		ctx, cancel := c.Context(), func() {}
		defer cancel()

		err := h.db.PingContext(ctx)
		if err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			isReady = false
		} else {
			// Check database stats
			stats := h.db.Stats()
			if stats.OpenConnections > 90 { // Warning threshold
				checks["database"] = "degraded: high connection count"
			} else {
				checks["database"] = "healthy"
			}
		}
	} else {
		checks["database"] = "not configured"
	}

	// Check Redis connection
	if h.redis != nil {
		ctx := c.Context()
		_, err := h.redis.Ping(ctx).Result()
		if err != nil {
			checks["redis"] = "unhealthy: " + err.Error()
			isReady = false
		} else {
			checks["redis"] = "healthy"
		}
	} else {
		checks["redis"] = "not configured"
	}

	// Check disk space (basic check)
	checks["disk"] = "healthy" // TODO: Implement actual disk space check

	// Determine overall status
	statusText := "ok"
	statusCode := fiber.StatusOK
	if !isReady {
		statusText = "unhealthy"
		statusCode = fiber.StatusServiceUnavailable
	}

	status := HealthStatus{
		Status:    statusText,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	}

	return c.Status(statusCode).JSON(status)
}

// MetricsCheck returns Prometheus metrics
// This is handled by the Prometheus middleware adapter
