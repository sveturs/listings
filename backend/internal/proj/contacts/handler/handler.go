// Package handler
// backend/internal/proj/contacts/handler/handler.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
)

type Handler struct {
	services globalService.ServicesInterface
}

func NewHandler(services globalService.ServicesInterface) *Handler {
	return &Handler{
		services: services,
	}
}

// AddContact добавляет нового контакта
// @Summary Добавить контакт
// @Description Добавляет пользователя в список контактов
// @Tags contacts
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.AddContactRequest true "Данные для добавления контакта"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserContact} "Контакт успешно добавлен"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidRequest"
// @Failure 409 {object} utils.ErrorResponseSwag "contacts.alreadyExists"
// @Failure 403 {object} utils.ErrorResponseSwag "contacts.userNotAllowRequests"
// @Failure 500 {object} utils.ErrorResponseSwag "contacts.addError"
// @Security BearerAuth
// @Router /api/v1/contacts [post]
func (h *Handler) AddContact(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.AddContactRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidRequest")
	}

	// Валидация
	if req.ContactUserID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.contactUserIdRequired")
	}

	contact, err := h.services.Contacts().AddContact(c.Context(), userID, &req)
	if err != nil {
		// Проверяем специфичные ошибки
		if err.Error() == "contact already exists" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "contacts.alreadyExists")
		}
		if err.Error() == "cannot add yourself as contact" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "contacts.cannotAddYourself")
		}
		if err.Error() == "user does not allow contact requests or has blocked you" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "contacts.userNotAllowRequests")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "contacts.addError")
	}

	return utils.SuccessResponse(c, contact)
}

// UpdateContactStatus обновляет статус контакта
// @Summary Обновить статус контакта
// @Description Изменяет статус контакта (принять или заблокировать)
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact_user_id path int true "ID контакта"
// @Param request body backend_internal_domain_models.UpdateContactRequest true "Новый статус контакта"
// @Success 200 {object} utils.SuccessResponseSwag{data=ContactStatusUpdateResponse} "Статус обновлен"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidContactUserId или validation.invalidStatus"
// @Failure 500 {object} utils.ErrorResponseSwag "contacts.updateError"
// @Security BearerAuth
// @Router /api/v1/contacts/{contact_user_id}/status [put]
func (h *Handler) UpdateContactStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	contactUserID, err := c.ParamsInt("contact_user_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidContactUserId")
	}

	var req models.UpdateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidRequest")
	}

	// Валидация статуса
	if req.Status != models.ContactStatusAccepted && req.Status != models.ContactStatusBlocked {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidStatus")
	}

	err = h.services.Contacts().UpdateContactStatus(c.Context(), userID, contactUserID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "contacts.updateError")
	}

	return utils.SuccessResponse(c, ContactStatusUpdateResponse{Message: "contacts.statusUpdated"})
}

// GetContacts возвращает список контактов пользователя
// @Summary Получить список контактов
// @Description Возвращает список контактов с фильтрацией по статусу и пагинацией
// @Tags contacts
// @Accept json
// @Produce json
// @Param status query string false "Фильтр по статусу (pending, accepted, blocked)"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество на странице" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.ContactsListResponse} "Список контактов"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidStatusFilter"
// @Failure 500 {object} utils.ErrorResponseSwag "contacts.fetchError"
// @Security BearerAuth
// @Router /api/v1/contacts [get]
func (h *Handler) GetContacts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	status := c.Query("status", "")
	page := utils.StringToInt(c.Query("page"), 1)
	limit := utils.StringToInt(c.Query("limit"), 20)

	// Валидация статуса
	if status != "" && status != models.ContactStatusPending &&
		status != models.ContactStatusAccepted && status != models.ContactStatusBlocked {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidStatusFilter")
	}

	contacts, err := h.services.Contacts().GetContacts(c.Context(), userID, status, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "contacts.fetchError")
	}

	return utils.SuccessResponse(c, contacts)
}

// GetIncomingRequests возвращает список входящих запросов в контакты
// @Summary Получить входящие запросы в контакты
// @Description Возвращает список входящих запросов в контакты со статусом pending
// @Tags contacts
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество на странице" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.ContactsListResponse} "Список входящих запросов"
// @Failure 500 {object} utils.ErrorResponseSwag "contacts.fetchError"
// @Security BearerAuth
// @Router /api/v1/contacts/incoming [get]
func (h *Handler) GetIncomingRequests(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	page := utils.StringToInt(c.Query("page"), 1)
	limit := utils.StringToInt(c.Query("limit"), 20)

	requests, err := h.services.Contacts().GetIncomingContactRequests(c.Context(), userID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "contacts.fetchError")
	}

	return utils.SuccessResponse(c, requests)
}

// RemoveContact удаляет контакт из списка
// @Summary Удалить контакт
// @Description Удаляет пользователя из списка контактов
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact_user_id path int true "ID контакта для удаления"
// @Success 200 {object} utils.SuccessResponseSwag{data=ContactRemoveResponse} "Контакт удален"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidContactUserId"
// @Failure 500 {object} utils.ErrorResponseSwag "contacts.removeError"
// @Security BearerAuth
// @Router /api/v1/contacts/{contact_user_id} [delete]
func (h *Handler) RemoveContact(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	contactUserID, err := c.ParamsInt("contact_user_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidContactUserId")
	}

	err = h.services.Contacts().RemoveContact(c.Context(), userID, contactUserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "contacts.removeError")
	}

	return utils.SuccessResponse(c, ContactRemoveResponse{Message: "contacts.removed"})
}

// GetPrivacySettings возвращает настройки приватности пользователя
// @Summary Получить настройки приватности
// @Description Возвращает текущие настройки приватности пользователя
// @Tags contacts
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserPrivacySettings} "Настройки приватности"
// @Failure 500 {object} utils.ErrorResponseSwag "privacy.fetchError"
// @Security BearerAuth
// @Router /api/v1/contacts/privacy [get]
func (h *Handler) GetPrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	settings, err := h.services.Contacts().GetPrivacySettings(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "privacy.fetchError")
	}

	return utils.SuccessResponse(c, settings)
}

// UpdatePrivacySettings обновляет настройки приватности
// @Summary Обновить настройки приватности
// @Description Изменяет настройки приватности пользователя
// @Tags contacts
// @Accept json
// @Produce json
// @Param request body backend_internal_domain_models.UpdatePrivacySettingsRequest true "Новые настройки приватности"
// @Success 200 {object} utils.SuccessResponseSwag{data=backend_internal_domain_models.UserPrivacySettings} "Обновленные настройки"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidRequest"
// @Failure 500 {object} utils.ErrorResponseSwag "privacy.updateError"
// @Security BearerAuth
// @Router /api/v1/contacts/privacy [put]
func (h *Handler) UpdatePrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.UpdatePrivacySettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidRequest")
	}

	settings, err := h.services.Contacts().UpdatePrivacySettings(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "privacy.updateError")
	}

	return utils.SuccessResponse(c, settings)
}

// GetContactStatus проверяет статус контакта между пользователями
// @Summary Проверить статус контакта
// @Description Проверяет, являются ли два пользователя контактами
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact_user_id path int true "ID пользователя для проверки"
// @Success 200 {object} utils.SuccessResponseSwag{data=ContactStatusCheckResponse} "Статус контакта"
// @Failure 400 {object} utils.ErrorResponseSwag "validation.invalidContactUserId"
// @Failure 500 {object} utils.ErrorResponseSwag "contacts.checkError"
// @Security BearerAuth
// @Router /api/v1/contacts/{contact_user_id}/check [get]
func (h *Handler) GetContactStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	contactUserIDStr := c.Params("contact_user_id")
	contactUserID, err := strconv.Atoi(contactUserIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.invalidContactUserId")
	}

	// Проверяем, являются ли пользователи контактами
	areContacts, err := h.services.Contacts().AreContacts(c.Context(), userID, contactUserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "contacts.checkError")
	}

	response := ContactStatusCheckResponse{
		AreContacts: areContacts,
		UserID:      userID,
		ContactID:   contactUserID,
	}

	return utils.SuccessResponse(c, response)
}
