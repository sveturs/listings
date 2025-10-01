// Package handler
// backend/internal/proj/docserver/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта docs
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Маршруты документации - требуют авторизации и прав администратора
	docsRoutes := app.Group("/api/v1/docs", mw.JWTParser(), authMiddleware.RequireAuth())

	// Документация доступна только администраторам
	docsRoutes.Get("/files", mw.RequireAdmin(), h.GetDocFiles)
	docsRoutes.Get("/content", mw.RequireAdmin(), h.GetFileContent)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/docs"
}
