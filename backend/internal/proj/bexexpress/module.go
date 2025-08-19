package bexexpress

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/proj/bexexpress/handler"
	"backend/internal/proj/bexexpress/service"
)

// Module представляет модуль BEX Express
type Module struct {
	handler *handler.Handler
	service *service.Service
}

// NewModule создает новый модуль BEX Express
func NewModule(db *sql.DB, cfg *config.Config) (*Module, error) {
	// Создаем сервис
	bexService, err := service.NewService(db, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create BEX service: %w", err)
	}

	// Создаем handler
	bexHandler := handler.NewHandler(bexService)

	return &Module{
		handler: bexHandler,
		service: bexService,
	}, nil
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	api := app.Group("/api/v1")
	
	// BEX API требует авторизацию
	api.Use(mw.AuthRequiredJWT)
	
	// Регистрируем маршруты BEX
	m.handler.RegisterRoutes(api)
	
	return nil
}

// GetPrefix возвращает префикс маршрутов модуля
func (m *Module) GetPrefix() string {
	return "/api/v1/bex"
}

// GetService возвращает сервис для использования в других модулях
func (m *Module) GetService() *service.Service {
	return m.service
}

// GetHandler возвращает handler
func (m *Module) GetHandler() *handler.Handler {
	return m.handler
}