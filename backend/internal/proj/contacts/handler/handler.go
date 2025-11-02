// backend/internal/proj/contacts/handler/handler.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/proj/global/service"
	"backend/pkg/utils"
)

// Handler представляет обработчик для контактов пользователей
type Handler struct {
	services service.ServicesInterface
}

// NewHandler создает новый экземпляр Handler
func NewHandler(services service.ServicesInterface) *Handler {
	return &Handler{
		services: services,
	}
}

// AddContact добавляет новый контакт - TODO: temporarily disabled
func (h *Handler) AddContact(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// UpdateContactStatus обновляет статус контакта - TODO: temporarily disabled  
func (h *Handler) UpdateContactStatus(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetContacts возвращает список контактов - TODO: temporarily disabled
func (h *Handler) GetContacts(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetIncomingContactRequests возвращает входящие запросы - TODO: temporarily disabled
func (h *Handler) GetIncomingContactRequests(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// RemoveContact удаляет контакт - TODO: temporarily disabled
func (h *Handler) RemoveContact(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetPrivacySettings возвращает настройки приватности - TODO: temporarily disabled
func (h *Handler) GetPrivacySettings(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// UpdatePrivacySettings обновляет настройки приватности - TODO: temporarily disabled
func (h *Handler) UpdatePrivacySettings(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// CheckContactStatus проверяет статус контакта - TODO: temporarily disabled
func (h *Handler) CheckContactStatus(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}
