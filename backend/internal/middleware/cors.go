// backend/internal/middleware/cors.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (m *Middleware) CORS() fiber.Handler {
	// Используем стандартный CORS middleware от Fiber
	// с конкретными origins для работы с credentials
	return cors.New(cors.Config{
		AllowOrigins:     "https://svetu.rs,https://www.svetu.rs,https://dev.svetu.rs,http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:3003",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true, // Разрешаем credentials
		MaxAge:           300,
	})
}
