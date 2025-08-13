package translation_admin

import (
	"context"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"backend/internal/middleware"
	"backend/internal/proj/translation_admin/ratelimit"
)

// Module represents the translation admin module
type Module struct {
	handler   *Handler
	aiHandler *AITranslationHandler
	service   *Service
	repo      *Repository
	logger    zerolog.Logger
}

// NewModule creates a new translation admin module
func NewModule(ctx context.Context, db *sqlx.DB, logger zerolog.Logger, frontendPath string, redisClient *redis.Client, translationService interface{}) *Module {
	// Create repository
	repo := NewRepository(db, logger)

	// Create service with proper frontend path and Redis
	messagesPath := filepath.Join(frontendPath, "frontend", "svetu")
	service := NewService(ctx, logger, messagesPath, repo, repo, redisClient, db.DB, translationService)

	// Create rate limiter for AI translations
	rateLimiter := ratelimit.NewMultiProviderRateLimiter(redisClient, ratelimit.DefaultConfig())

	// Create handlers
	handler := NewHandler(service, logger)
	aiHandler := NewAITranslationHandler(logger, service, rateLimiter)

	// Warm up cache if Redis is available
	if redisClient != nil {
		go func(ctx context.Context) {
			if err := service.WarmUpCache(ctx); err != nil {
				logger.Error().Err(err).Msg("Failed to warm up translation cache")
			}
		}(ctx)
	}

	return &Module{
		handler:   handler,
		aiHandler: aiHandler,
		service:   service,
		repo:      repo,
		logger:    logger,
	}
}

// RegisterRoutes registers the module routes
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	// Public test endpoint for translation testing
	app.Post("/api/v1/test/translate", m.aiHandler.TranslateText)

	// Admin-only endpoints for translation management
	admin := app.Group("/api/v1/admin/translations",
		middleware.AuthRequiredJWT,
		middleware.AdminRequired,
	)

	// Register all translation admin routes
	m.handler.RegisterRoutes(admin)

	// Register AI translation routes
	ai := admin.Group("/ai")
	ai.Get("/providers", m.aiHandler.GetAvailableProviders)
	ai.Get("/rate-limit-status", m.aiHandler.GetRateLimitStatus)
	ai.Post("/translate", m.aiHandler.TranslateText)
	ai.Post("/translate-batch", m.aiHandler.TranslateBatch)
	ai.Post("/translate-module", m.aiHandler.TranslateModule)

	// Cost tracking endpoints
	ai.Get("/costs", m.aiHandler.GetCostsSummary)
	ai.Get("/costs/alerts", m.aiHandler.GetCostAlerts)
	ai.Get("/costs/:provider", m.aiHandler.GetProviderCostDetails)
	ai.Post("/costs/:provider/reset", m.aiHandler.ResetProviderCosts)

	m.logger.Info().Msg("Translation admin routes registered")

	return nil
}

// GetPrefix returns the module prefix for logging
func (m *Module) GetPrefix() string {
	return "translation_admin"
}
