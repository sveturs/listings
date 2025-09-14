// backend/internal/middleware/cors.go
package middleware

import (
	"backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (m *Middleware) CORS() fiber.Handler {
	// Используем стандартный CORS middleware от Fiber
	// с конкретными origins для работы с credentials
	// Добавляем VPN адрес из конфигурации
	allowedOrigins := "https://svetu.rs,https://www.svetu.rs,https://preprod.svetu.rs,https://dev.svetu.rs,https://auth.svetu.net,http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:3003,http://100.88.44.15:3001"

	// Если есть FRONTEND_URL в конфигурации, добавляем его тоже
	if m.config != nil && m.config.FrontendURL != "" && m.config.FrontendURL != "http://100.88.44.15:3001" {
		allowedOrigins += "," + m.config.FrontendURL
	}

	logger.Info().
		Str("allowed_origins", allowedOrigins).
		Msg("CORS middleware configured with origins")

	corsHandler := cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token",
		ExposeHeaders:    "Content-Length,Set-Cookie",
		AllowCredentials: true, // Разрешаем credentials
		MaxAge:           300,
	})

	// Обёртка для дополнительного логирования
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin != "" && c.Method() == "OPTIONS" {
			logger.Debug().
				Str("origin", origin).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Msg("CORS preflight request")
		}
		return corsHandler(c)
	}
}
