package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ContactsHandler struct {
	services globalService.ServicesInterface
}

func NewContactsHandler(services globalService.ServicesInterface) *ContactsHandler {
	return &ContactsHandler{
		services: services,
	}
}

// Добавить контакт
func (h *ContactsHandler) AddContact(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.AddContactRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Валидация
	if req.ContactUserID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Contact user ID is required")
	}

	contact, err := h.services.Contacts().AddContact(c.Context(), userID, &req)
	if err != nil {
		// Проверяем специфичные ошибки
		if err.Error() == "contact already exists" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "contact already exists")
		}
		if err.Error() == "cannot add yourself as contact" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "cannot add yourself as contact")
		}
		if err.Error() == "user does not allow contact requests or has blocked you" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "user does not allow contact requests or has blocked you")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, contact)
}

// Обновить статус контакта
func (h *ContactsHandler) UpdateContactStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	contactUserID, err := c.ParamsInt("contact_user_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid contact user ID")
	}

	var req models.UpdateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Валидация статуса
	if req.Status != models.ContactStatusAccepted && req.Status != models.ContactStatusBlocked {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status")
	}

	err = h.services.Contacts().UpdateContactStatus(c.Context(), userID, contactUserID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Contact status updated"})
}

// Получить список контактов
func (h *ContactsHandler) GetContacts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	status := c.Query("status", "")
	page := utils.StringToInt(c.Query("page"), 1)
	limit := utils.StringToInt(c.Query("limit"), 20)

	// Валидация статуса
	if status != "" && status != models.ContactStatusPending &&
		status != models.ContactStatusAccepted && status != models.ContactStatusBlocked {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid status filter")
	}

	contacts, err := h.services.Contacts().GetContacts(c.Context(), userID, status, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching contacts")
	}

	return utils.SuccessResponse(c, contacts)
}

// Удалить контакт
func (h *ContactsHandler) RemoveContact(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	contactUserID, err := c.ParamsInt("contact_user_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid contact user ID")
	}

	err = h.services.Contacts().RemoveContact(c.Context(), userID, contactUserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Contact removed"})
}

// Получить настройки приватности
func (h *ContactsHandler) GetPrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	settings, err := h.services.Contacts().GetPrivacySettings(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching privacy settings")
	}

	return utils.SuccessResponse(c, settings)
}

// Обновить настройки приватности
func (h *ContactsHandler) UpdatePrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.UpdatePrivacySettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	settings, err := h.services.Contacts().UpdatePrivacySettings(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, settings)
}

// Проверить статус контакта
func (h *ContactsHandler) GetContactStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	contactUserIDStr := c.Params("contact_user_id")
	contactUserID, err := strconv.Atoi(contactUserIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid contact user ID")
	}

	// Проверяем, являются ли пользователи контактами
	areContacts, err := h.services.Contacts().AreContacts(c.Context(), userID, contactUserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error checking contact status")
	}

	response := fiber.Map{
		"are_contacts": areContacts,
		"user_id":      userID,
		"contact_id":   contactUserID,
	}

	return utils.SuccessResponse(c, response)
}
