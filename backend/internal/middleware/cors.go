package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "strings"
)

func (m *Middleware) CORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins:     strings.Join([]string{m.config.FrontendURL,"http://localhost:3000","http://localhost:3001,http://landhub.rs,https://landhub.rs"}, ","),
        AllowMethods:     "GET,POST,DELETE,PUT,OPTIONS",
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
        AllowCredentials: true,
        ExposeHeaders:    "Content-Length",
    })
}