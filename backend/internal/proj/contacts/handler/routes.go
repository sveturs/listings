package handler

import (
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes регистрирует маршруты для модуля contacts
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Группа маршрутов для контактов
	contacts := app.Group("/api/v1/contacts", mw.AuthRequiredJWT)

	// Маршруты для работы с контактами
	contacts.Get("/", h.GetContacts)                             // Получить список контактов
	contacts.Post("/", h.AddContact)                             // Добавить контакт
	contacts.Put("/:contact_user_id", h.UpdateContactStatus)     // Обновить статус контакта
	contacts.Delete("/:contact_user_id", h.RemoveContact)        // Удалить контакт
	contacts.Get("/status/:contact_user_id", h.GetContactStatus) // Проверить статус контакта
	contacts.Get("/privacy", h.GetPrivacySettings)               // Получить настройки приватности
	contacts.Put("/privacy", h.UpdatePrivacySettings)            // Обновить настройки приватности

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "contacts"
}
