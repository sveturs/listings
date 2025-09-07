package config

import (
	"backend/internal/config"
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль конфигурации
type Module struct {
	Handler *Handler // Экспортируем для доступа извне
}

// NewModule создает новый экземпляр модуля конфигурации
func NewModule(cfg *config.Config) *Module {
	return &Module{
		Handler: NewHandler(cfg),
	}
}

// GetPrefix возвращает префикс для API конфигурации
func (m *Module) GetPrefix() string {
	return m.Handler.GetPrefix()
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	return m.Handler.RegisterRoutes(app, mw)
}
