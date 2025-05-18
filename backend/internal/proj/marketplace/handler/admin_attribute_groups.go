package handler

import (
	"backend/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// CreateAttributeGroup создает новую группу атрибутов
func (h *MarketplaceHandler) CreateAttributeGroup(c *fiber.Ctx) error {
	var req models.CreateAttributeGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Валидация
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Название группы обязательно",
		})
	}

	group := &models.AttributeGroup{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Icon:        req.Icon,
		IsActive:    true,
		SortOrder:   req.SortOrder,
	}

	id, err := h.storage.AttributeGroups.CreateAttributeGroup(c.Context(), group)
	if err != nil {
		log.Printf("Ошибка создания группы атрибутов: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка создания группы",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
	})
}

// ListAttributeGroups возвращает список всех групп атрибутов
func (h *MarketplaceHandler) ListAttributeGroups(c *fiber.Ctx) error {
	groups, err := h.storage.AttributeGroups.ListAttributeGroups(c.Context())
	if err != nil {
		log.Printf("Ошибка получения списка групп: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения списка групп",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"groups":  groups,
	})
}

// GetAttributeGroup получает информацию о группе атрибутов
func (h *MarketplaceHandler) GetAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	group, err := h.storage.AttributeGroups.GetAttributeGroup(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка получения группы: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Группа не найдена",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"group":   group,
	})
}

// UpdateAttributeGroup обновляет группу атрибутов
func (h *MarketplaceHandler) UpdateAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	var req models.UpdateAttributeGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	// For boolean and integer fields we always update them
	updates["is_active"] = req.IsActive
	updates["sort_order"] = req.SortOrder

	err = h.storage.AttributeGroups.UpdateAttributeGroup(c.Context(), id, updates)
	if err != nil {
		log.Printf("Ошибка обновления группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка обновления группы",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Группа успешно обновлена",
	})
}

// DeleteAttributeGroup удаляет группу атрибутов
func (h *MarketplaceHandler) DeleteAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	err = h.storage.AttributeGroups.DeleteAttributeGroup(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка удаления группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка удаления группы",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Группа успешно удалена",
	})
}

// GetAttributeGroupWithItems получает группу с её атрибутами
func (h *MarketplaceHandler) GetAttributeGroupWithItems(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	group, err := h.storage.AttributeGroups.GetAttributeGroup(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка получения группы: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Группа не найдена",
		})
	}

	items, err := h.storage.AttributeGroups.GetGroupItems(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка получения атрибутов группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения атрибутов",
		})
	}

	response := struct {
		Group *models.AttributeGroup      `json:"group"`
		Items []*models.AttributeGroupItem `json:"items"`
	}{
		Group: group,
		Items: items,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}

// AddItemToGroup добавляет атрибут в группу
func (h *MarketplaceHandler) AddItemToGroup(c *fiber.Ctx) error {
	groupIDStr := c.Params("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	var req models.AddItemToGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	item := &models.AttributeGroupItem{
		GroupID:            groupID,
		AttributeID:        req.AttributeID,
		Icon:               req.Icon,
		SortOrder:          req.SortOrder,
		CustomDisplayName:  req.CustomDisplayName,
		VisibilityCondition: req.VisibilityCondition,
	}

	id, err := h.storage.AttributeGroups.AddItemToGroup(c.Context(), groupID, item)
	if err != nil {
		log.Printf("Ошибка добавления атрибута в группу: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка добавления атрибута",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
	})
}

// RemoveItemFromGroup удаляет атрибут из группы
func (h *MarketplaceHandler) RemoveItemFromGroup(c *fiber.Ctx) error {
	groupIDStr := c.Params("id")
	attributeIDStr := c.Params("attributeId")
	
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	attributeID, err := strconv.Atoi(attributeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID атрибута",
		})
	}

	err = h.storage.AttributeGroups.RemoveItemFromGroup(c.Context(), groupID, attributeID)
	if err != nil {
		log.Printf("Ошибка удаления атрибута из группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка удаления атрибута",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Атрибут успешно удален из группы",
	})
}

// GetCategoryGroups получает группы атрибутов, привязанные к категории
func (h *MarketplaceHandler) GetCategoryGroups(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID категории",
		})
	}

	groups, err := h.storage.AttributeGroups.GetCategoryGroups(c.Context(), categoryID)
	if err != nil {
		log.Printf("Ошибка получения групп категории: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения групп",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"groups":  groups,
	})
}

// AttachGroupToCategory привязывает группу атрибутов к категории
func (h *MarketplaceHandler) AttachGroupToCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID категории",
		})
	}

	var req models.AttachGroupToCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	attachment := &models.CategoryAttributeGroup{
		CategoryID: categoryID,
		GroupID:    req.GroupID,
		SortOrder:  req.SortOrder,
	}

	id, err := h.storage.AttributeGroups.AttachGroupToCategory(c.Context(), categoryID, attachment)
	if err != nil {
		log.Printf("Ошибка привязки группы к категории: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка привязки группы",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
	})
}

// DetachGroupFromCategory отвязывает группу атрибутов от категории
func (h *MarketplaceHandler) DetachGroupFromCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	groupIDStr := c.Params("groupId")
	
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID категории",
		})
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ID группы",
		})
	}

	err = h.storage.AttributeGroups.DetachGroupFromCategory(c.Context(), categoryID, groupID)
	if err != nil {
		log.Printf("Ошибка отвязки группы от категории: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка отвязки группы",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Группа успешно отвязана от категории",
	})
}