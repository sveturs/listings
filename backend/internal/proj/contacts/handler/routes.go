package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует маршруты для модуля contacts
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Группа маршрутов для контактов
	// БЕЗ CSRF - используем BFF proxy архитектуру
	contacts := app.Group("/api/v1/contacts", h.jwtParserMW, authMiddleware.RequireAuthString(), mw.RateLimitByUser(300, time.Minute))

	// Маршруты для работы с контактами
	contacts.Get("/", h.GetContacts)                                // Получить список контактов
	contacts.Post("/", h.AddContact)                                // Добавить контакт
	contacts.Get("/incoming", h.GetIncomingRequests)                // Получить входящие запросы
	contacts.Put("/:contact_user_id/status", h.UpdateContactStatus) // Обновить статус контакта
	contacts.Delete("/:contact_user_id", h.RemoveContact)           // Удалить контакт
	contacts.Get("/:contact_user_id/check", h.GetContactStatus)     // Проверить статус контакта
	contacts.Get("/privacy", h.GetPrivacySettings)                  // Получить настройки приватности
	contacts.Put("/privacy", h.UpdatePrivacySettings)               // Обновить настройки приватности

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/contacts"
}
