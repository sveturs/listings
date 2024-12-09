//backend/internal/middleware/middleware.go
package middleware

import (
	"backend/internal/config"
	"backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	config   *config.Config
	services *services.Services
}

func NewMiddleware(cfg *config.Config, services *services.Services) *Middleware {
	return &Middleware{
		config:   cfg,
		services: services,
	}
}

func (m *Middleware) Setup(app *fiber.App) {
	app.Use(m.Logger())
	app.Use(m.CORS())
}

// ErrorHandler обрабатывает все ошибки приложения
func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
