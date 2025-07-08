package interfaces

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RouteRegistrar интерфейс для регистрации маршрутов проектов
type RouteRegistrar interface {
	// RegisterRoutes регистрирует все маршруты проекта
	RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error

	// GetPrefix возвращает префикс проекта для логирования
	GetPrefix() string
}
