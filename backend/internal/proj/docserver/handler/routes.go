// Package handler
// backend/internal/proj/docserver/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта docs
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные маршруты документации - без аутентификации
	docsRoutes := app.Group("/api/v1/docs")

	// Документация доступна всем
	docsRoutes.Get("/files", h.GetDocFiles)
	docsRoutes.Get("/content", h.GetFileContent)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/docs"
}
