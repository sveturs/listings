package vin

import (
	"backend/internal/middleware"
	"backend/internal/proj/vin/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Module представляет модуль VIN декодера
type Module struct {
	handler *handler.VINHandler
}

// NewModule создает новый модуль VIN
func NewModule(db *sqlx.DB) *Module {
	return &Module{
		handler: handler.NewVINHandler(db),
	}
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	m.handler.RegisterRoutes(app)
	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (m *Module) GetPrefix() string {
	return "vin"
}
