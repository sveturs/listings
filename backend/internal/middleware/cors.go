// backend/internal/middleware/cors.go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (m *Middleware) CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: strings.Join([]string{
			m.config.FrontendURL,
			"http://localhost:3000",
			"http://localhost:3001",
			"http://SveTu.rs",
			"https://SveTu.rs",
		}, ","),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
		AllowCredentials: true,
		MaxAge:           300, // 5 минут в секундах
	})
}
