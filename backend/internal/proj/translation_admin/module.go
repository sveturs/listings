package translation_admin

import (
	"context"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"backend/internal/middleware"
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
func NewModule(ctx context.Context, db *sqlx.DB, logger zerolog.Logger, frontendPath string) *Module {
	// Create repository
	repo := NewRepository(db, logger)

	// Create service with proper frontend path
	messagesPath := filepath.Join(frontendPath, "frontend", "svetu")
	service := NewService(logger, messagesPath, repo, repo)

	// Create handlers
	handler := NewHandler(service, logger)
	aiHandler := NewAITranslationHandler(logger, service)

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
	ai.Post("/translate", m.aiHandler.TranslateText)
	ai.Post("/translate-batch", m.aiHandler.TranslateBatch)
	ai.Post("/translate-module", m.aiHandler.TranslateModule)

	m.logger.Info().Msg("Translation admin routes registered")

	return nil
}

// GetPrefix returns the module prefix for logging
func (m *Module) GetPrefix() string {
	return "translation_admin"
}
