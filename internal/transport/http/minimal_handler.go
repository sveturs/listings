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

	"github.com/sveturs/listings/internal/domain"
)

// MinimalService defines minimal interface for testing
type MinimalService interface {
	CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error)
	GetListing(ctx context.Context, id int64) (*domain.Listing, error)
	ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error)

	// Storefront methods
	GetStorefront(ctx context.Context, storefrontID int64) (*domain.Storefront, error)
	GetStorefrontBySlug(ctx context.Context, slug string) (*domain.Storefront, error)
	ListStorefronts(ctx context.Context, limit, offset int) ([]*domain.Storefront, int64, error)
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

	// Storefronts endpoints
	api.Get("/storefronts", h.ListStorefronts)
	api.Get("/storefronts/:id", h.GetStorefront)
	api.Get("/storefronts/slug/:slug", h.GetStorefrontBySlug)
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

// GetStorefront handles GET storefront by ID
func (h *MinimalHandler) GetStorefront(c *fiber.Ctx) error {
	// Add X-Served-By header for traffic distribution measurement
	c.Set("X-Served-By", "microservice")

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id parameter"})
	}

	storefront, err := h.service.GetStorefront(c.Context(), int64(id))
	if err != nil {
		h.logger.Error().Err(err).Int64("id", int64(id)).Msg("failed to get storefront")
		return c.Status(404).JSON(fiber.Map{"error": "storefront not found"})
	}

	return c.JSON(storefront)
}

// GetStorefrontBySlug handles GET storefront by slug
func (h *MinimalHandler) GetStorefrontBySlug(c *fiber.Ctx) error {
	// Add X-Served-By header for traffic distribution measurement
	c.Set("X-Served-By", "microservice")

	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(fiber.Map{"error": "slug parameter is required"})
	}

	storefront, err := h.service.GetStorefrontBySlug(c.Context(), slug)
	if err != nil {
		h.logger.Error().Err(err).Str("slug", slug).Msg("failed to get storefront by slug")
		return c.Status(404).JSON(fiber.Map{"error": "storefront not found"})
	}

	return c.JSON(storefront)
}

// ListStorefronts handles GET storefronts list
func (h *MinimalHandler) ListStorefronts(c *fiber.Ctx) error {
	// Add X-Served-By header for traffic distribution measurement
	c.Set("X-Served-By", "microservice")

	// Parse pagination parameters
	limit := c.QueryInt("limit", 20)
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	storefronts, total, err := h.service.ListStorefronts(c.Context(), limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Int("limit", limit).Int("offset", offset).Msg("failed to list storefronts")
		return c.Status(500).JSON(fiber.Map{"error": "failed to list storefronts"})
	}

	return c.JSON(fiber.Map{
		"storefronts": storefronts,
		"total":       total,
	})
}

// StartMinimalServer starts minimal HTTP server
func StartMinimalServer(host string, port int, handler *MinimalHandler, healthHandler *HealthHandler, logger zerolog.Logger) (*fiber.App, error) {
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

	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		if err := app.Listen(addr); err != nil {
			logger.Fatal().Err(err).Msg("failed to start HTTP server")
		}
	}()

	logger.Info().Str("addr", addr).Msg("HTTP server started")
	return app, nil
}
