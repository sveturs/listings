package search_admin

import (
	"backend/internal/middleware"
	"backend/internal/proj/search_admin/handler"
	"backend/internal/proj/search_admin/service"
	"backend/internal/storage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль управления конфигурацией поиска
type Module struct {
	handler *handler.Handler
	service *service.Service
}

// NewModule создает новый экземпляр модуля
func NewModule(db *postgres.Database, osClient *opensearch.OpenSearchClient, logger *logger.Logger, b2cIndexName string) *Module {
	searchService := service.NewService(db.GetSQLXDB(), osClient, b2cIndexName)
	searchHandler := handler.NewHandler(searchService, logger)

	return &Module{
		handler: searchHandler,
		service: searchService,
	}
}

// SetStorage устанавливает storage для сервиса
func (m *Module) SetStorage(storage storage.Storage) {
	if m.service != nil {
		m.service.SetStorage(storage)
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
