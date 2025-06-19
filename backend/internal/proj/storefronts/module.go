package storefronts

import (
	"backend/internal/middleware"
	"backend/internal/proj/storefronts/handler"
	"backend/internal/proj/storefronts/routes"
	"backend/internal/proj/storefronts/service"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль витрин
type Module struct {
	handler *handler.StorefrontHandler
	service service.StorefrontService
}

// NewModule создает новый модуль витрин
func NewModule(db *postgres.Database, fileRepo filestorage.FileStorageInterface) *Module {
	// Создаем репозиторий
	repo := postgres.NewStorefrontRepository(db)
	
	// Создаем сервисы
	storefrontService := service.NewStorefrontService(repo, fileRepo)
	
	// Создаем handler
	h := handler.NewStorefrontHandler(storefrontService)
	
	return &Module{
		handler: h,
		service: storefrontService,
	}
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, authMiddleware *middleware.Middleware) error {
	routes.RegisterStorefrontRoutes(app, m.handler, authMiddleware)
	routes.RegisterStorefrontWebhooks(app)
	return nil
}

// GetPrefix возвращает префикс маршрутов модуля
func (m *Module) GetPrefix() string {
	return "/api/v1/storefronts"
}

// GetService возвращает сервис витрин (для использования другими модулями)
func (m *Module) GetService() service.StorefrontService {
	return m.service
}

// MigrateData выполняет миграцию данных из старой структуры
func (m *Module) MigrateData() error {
	// TODO: Implement data migration from user_storefronts to new structure
	// This will be done as a separate task
	return nil
}