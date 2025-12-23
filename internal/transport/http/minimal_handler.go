package http

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// MinimalService defines minimal interface for testing
type MinimalService interface {
	CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error)
	GetListing(ctx context.Context, id int64) (*domain.Listing, error)
	ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error)
}

// MinimalHandler provides basic HTTP endpoints
type MinimalHandler struct {
	service MinimalService
	logger  zerolog.Logger
}

// NewMinimalHandler creates minimal HTTP handler
func NewMinimalHandler(service MinimalService, logger zerolog.Logger) *MinimalHandler {
	return &MinimalHandler{
		service: service,
		logger:  logger.With().Str("component", "http_handler").Logger(),
	}
}

// SetupRoutes sets up minimal routes
// Note: Health check routes are now handled by HealthHandler
func (h *MinimalHandler) SetupRoutes(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Prometheus metrics endpoint
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	api := app.Group("/api/v1")

	// Listings endpoints
	api.Get("/listings", h.ListListings)
	api.Get("/listings/:id", h.GetListing)
	api.Post("/listings", h.CreateListing)
}

// GetListing handles GET listing
func (h *MinimalHandler) GetListing(c *fiber.Ctx) error {
	// Add X-Served-By header for traffic distribution measurement (Phase 3)
	c.Set("X-Served-By", "microservice")

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id parameter"})
	}
	listing, err := h.service.GetListing(c.Context(), int64(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(listing)
}

// CreateListing handles POST listing
func (h *MinimalHandler) CreateListing(c *fiber.Ctx) error {
	// Add X-Served-By header for traffic distribution measurement (Phase 3)
	c.Set("X-Served-By", "microservice")

	var input domain.CreateListingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	listing, err := h.service.CreateListing(c.Context(), &input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(listing)
}

// ListListings handles GET listings
func (h *MinimalHandler) ListListings(c *fiber.Ctx) error {
	// Add X-Served-By header for traffic distribution measurement (Phase 3)
	c.Set("X-Served-By", "microservice")

	filter := &domain.ListListingsFilter{Limit: 20, Offset: 0}
	listings, total, err := h.service.ListListings(c.Context(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"listings": listings, "total": total})
}

// StartMinimalServer starts minimal HTTP server
func StartMinimalServer(host string, port int, handler *MinimalHandler, healthHandler *HealthHandler, wsHandler *ChatWebSocketHandler, analyticsHandler *AnalyticsHandler, logger zerolog.Logger) (*fiber.App, error) {
	app := fiber.New(fiber.Config{
		AppName:      "Listings Service",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	// Register API routes first
	handler.SetupRoutes(app)

	// Register health check routes BEFORE starting server
	if healthHandler != nil {
		healthHandler.RegisterRoutes(app)
	}

	// Register WebSocket routes BEFORE starting server
	if wsHandler != nil {
		wsHandler.RegisterWebSocketRoute(app)
	}

	// Register analytics routes (Phase 7)
	if analyticsHandler != nil {
		analyticsHandler.RegisterRoutes(app)
		logger.Info().Msg("Search analytics routes registered")
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		if err := app.Listen(addr); err != nil {
			logger.Fatal().Err(err).Msg("failed to start HTTP server")
		}
	}()

	logger.Info().Str("addr", addr).Msg("HTTP server started")
	return app, nil
}
