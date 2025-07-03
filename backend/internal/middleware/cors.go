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
			"https://svetu.rs",
			"https://www.svetu.rs",
			"http://localhost:3000",       // для бэкенда
			"http://localhost:3001",       // для фронтенда
			"http://localhost:3002",       // для фронтенда
			"http://hostel_frontend:3001", // для Docker
			"http://backend:3000",         // для Docker
			"http://100.88.44.15:3000",    // IP адрес бэкенда
			"http://100.88.44.15:3001",    // IP адрес фронтенда
		}, ","),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token",
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
		AllowCredentials: true,
		MaxAge:           300,
	})
}
