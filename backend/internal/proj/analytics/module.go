package analytics

import (
	"backend/internal/middleware"
	"backend/internal/proj/analytics/handler"
	"backend/internal/proj/analytics/routes"
	"backend/internal/proj/analytics/service"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

// Module модуль аналитики
type Module struct {
	handler *handler.AnalyticsHandler
}

// NewModule создает новый модуль аналитики
func NewModule(db *postgres.Database, osClient *opensearch.OpenSearchClient) *Module {
	storefrontRepo := postgres.NewStorefrontRepository(db)
	analyticsService := service.NewAnalyticsService(storefrontRepo, osClient, db)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	return &Module{
		handler: analyticsHandler,
	}
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	routes.RegisterAnalyticsRoutes(app, m.handler, middleware)
	return nil
}

// GetPrefix возвращает префикс маршрутов модуля
func (m *Module) GetPrefix() string {
	return "/api/v1/analytics"
}
