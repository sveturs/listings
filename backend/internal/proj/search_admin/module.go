package search_admin

import (
	"backend/internal/middleware"
	"backend/internal/proj/search_admin/handler"
	"backend/internal/proj/search_admin/service"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль управления конфигурацией поиска
type Module struct {
	handler *handler.Handler
	service *service.Service
}

// NewModule создает новый экземпляр модуля
func NewModule(db *postgres.Database) *Module {
	searchService := service.NewService(db.GetSQLXDB())
	searchHandler := handler.NewHandler(searchService)

	return &Module{
		handler: searchHandler,
		service: searchService,
	}
}

// RegisterRoutes регистрирует все маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	return m.handler.RegisterRoutes(app, middleware)
}

// GetPrefix возвращает префикс модуля для логирования
func (m *Module) GetPrefix() string {
	return m.handler.GetPrefix()
}

// GetService возвращает сервис для использования в других модулях
func (m *Module) GetService() *service.Service {
	return m.service
}
