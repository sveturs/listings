//backend/internal/middleware/middleware.go
package middleware

import (
	"backend/internal/config"


	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service" 
)

type Middleware struct {
	config   *config.Config
    services globalService.ServicesInterface 
}

func NewMiddleware(cfg *config.Config, services globalService.ServicesInterface) *Middleware {
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
