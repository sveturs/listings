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
	// БЕЗ CSRF - используем BFF proxy архитектуру
	docsRoutes := app.Group("/api/v1/docs", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))

	// Документация доступна только администраторам
	docsRoutes.Get("/files", h.GetDocFiles)
	docsRoutes.Get("/content", h.GetFileContent)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/docs"
}
