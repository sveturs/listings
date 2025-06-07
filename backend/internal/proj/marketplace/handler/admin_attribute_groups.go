package handler

import (
	"backend/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// CreateAttributeGroup создает новую группу атрибутов
// @Summary Create attribute group
// @Description Creates a new attribute group for organizing attributes
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param body body models.CreateAttributeGroupRequest true "Attribute group data"
// @Success 201 {object} object{success=bool,id=int} "Group created successfully"
// @Failure 400 {object} object{error=string} "marketplace.invalidRequest or marketplace.groupNameRequired"
// @Failure 500 {object} object{error=string} "marketplace.createGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups [post]
func (h *MarketplaceHandler) CreateAttributeGroup(c *fiber.Ctx) error {
	var req models.CreateAttributeGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidRequest",
		})
	}

	// Валидация
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.groupNameRequired",
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
			"error": "marketplace.createGroupError",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
	})
}

// ListAttributeGroups возвращает список всех групп атрибутов
// @Summary List attribute groups
// @Description Returns all attribute groups
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,groups=[]models.AttributeGroup} "List of attribute groups"
// @Failure 500 {object} object{error=string} "marketplace.listGroupsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups [get]
func (h *MarketplaceHandler) ListAttributeGroups(c *fiber.Ctx) error {
	groups, err := h.storage.AttributeGroups.ListAttributeGroups(c.Context())
	if err != nil {
		log.Printf("Ошибка получения списка групп: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.listGroupsError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"groups":  groups,
	})
}

// GetAttributeGroup получает информацию о группе атрибутов
// @Summary Get attribute group
// @Description Returns information about a specific attribute group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} object{success=bool,group=models.AttributeGroup} "Attribute group information"
// @Failure 400 {object} object{error=string} "marketplace.invalidGroupId"
// @Failure 404 {object} object{error=string} "marketplace.groupNotFound"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id} [get]
func (h *MarketplaceHandler) GetAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	group, err := h.storage.AttributeGroups.GetAttributeGroup(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка получения группы: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "marketplace.groupNotFound",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"group":   group,
	})
}

// UpdateAttributeGroup обновляет группу атрибутов
// @Summary Update attribute group
// @Description Updates an existing attribute group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param body body models.UpdateAttributeGroupRequest true "Updated group data"
// @Success 200 {object} object{success=bool,message=string} "marketplace.groupUpdated"
// @Failure 400 {object} object{error=string} "marketplace.invalidGroupId or marketplace.invalidRequest"
// @Failure 500 {object} object{error=string} "marketplace.updateGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id} [put]
func (h *MarketplaceHandler) UpdateAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	var req models.UpdateAttributeGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidRequest",
		})
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.Icon != nil {
		updates["icon"] = req.Icon
	}
	// For boolean and integer fields we always update them
	updates["is_active"] = req.IsActive
	updates["sort_order"] = req.SortOrder

	err = h.storage.AttributeGroups.UpdateAttributeGroup(c.Context(), id, updates)
	if err != nil {
		log.Printf("Ошибка обновления группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.updateGroupError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.groupUpdated",
	})
}

// DeleteAttributeGroup удаляет группу атрибутов
// @Summary Delete attribute group
// @Description Deletes an attribute group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} object{success=bool,message=string} "marketplace.groupDeleted"
// @Failure 400 {object} object{error=string} "marketplace.invalidGroupId"
// @Failure 500 {object} object{error=string} "marketplace.deleteGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id} [delete]
func (h *MarketplaceHandler) DeleteAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	err = h.storage.AttributeGroups.DeleteAttributeGroup(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка удаления группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.deleteGroupError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.groupDeleted",
	})
}

// GetAttributeGroupWithItems получает группу с её атрибутами
// @Summary Get attribute group with items
// @Description Returns attribute group with all its attributes
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} object{success=bool,data=object{group=models.AttributeGroup,items=[]models.AttributeGroupItem}} "Group with items"
// @Failure 400 {object} object{error=string} "marketplace.invalidGroupId"
// @Failure 404 {object} object{error=string} "marketplace.groupNotFound"
// @Failure 500 {object} object{error=string} "marketplace.getGroupItemsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id}/items [get]
func (h *MarketplaceHandler) GetAttributeGroupWithItems(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	group, err := h.storage.AttributeGroups.GetAttributeGroup(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка получения группы: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "marketplace.groupNotFound",
		})
	}

	items, err := h.storage.AttributeGroups.GetGroupItems(c.Context(), id)
	if err != nil {
		log.Printf("Ошибка получения атрибутов группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.getGroupItemsError",
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
// @Summary Add attribute to group
// @Description Adds an attribute to a group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param body body models.AddItemToGroupRequest true "Attribute data"
// @Success 201 {object} object{success=bool,id=int} "Item added successfully"
// @Failure 400 {object} object{error=string} "marketplace.invalidGroupId or marketplace.invalidRequest"
// @Failure 500 {object} object{error=string} "marketplace.addItemError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id}/items [post]
func (h *MarketplaceHandler) AddItemToGroup(c *fiber.Ctx) error {
	groupIDStr := c.Params("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	var req models.AddItemToGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidRequest",
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
			"error": "marketplace.addItemError",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
	})
}

// RemoveItemFromGroup удаляет атрибут из группы
// @Summary Remove attribute from group
// @Description Removes an attribute from a group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param attributeId path int true "Attribute ID"
// @Success 200 {object} object{success=bool,message=string} "marketplace.itemRemoved"
// @Failure 400 {object} object{error=string} "marketplace.invalidGroupId or marketplace.invalidAttributeId"
// @Failure 500 {object} object{error=string} "marketplace.removeItemError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id}/items/{attributeId} [delete]
func (h *MarketplaceHandler) RemoveItemFromGroup(c *fiber.Ctx) error {
	groupIDStr := c.Params("id")
	attributeIDStr := c.Params("attributeId")
	
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	attributeID, err := strconv.Atoi(attributeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidAttributeId",
		})
	}

	err = h.storage.AttributeGroups.RemoveItemFromGroup(c.Context(), groupID, attributeID)
	if err != nil {
		log.Printf("Ошибка удаления атрибута из группы: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.removeItemError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.itemRemoved",
	})
}

// GetCategoryGroups получает группы атрибутов, привязанные к категории
// @Summary Get category attribute groups
// @Description Returns attribute groups attached to a category
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} object{success=bool,groups=[]models.AttributeGroup} "Category groups"
// @Failure 400 {object} object{error=string} "marketplace.invalidCategoryId"
// @Failure 500 {object} object{error=string} "marketplace.getCategoryGroupsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/categories/{id}/groups [get]
func (h *MarketplaceHandler) GetCategoryGroups(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidCategoryId",
		})
	}

	groups, err := h.storage.AttributeGroups.GetCategoryGroups(c.Context(), categoryID)
	if err != nil {
		log.Printf("Ошибка получения групп категории: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.getCategoryGroupsError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"groups":  groups,
	})
}

// AttachGroupToCategory привязывает группу атрибутов к категории
// @Summary Attach group to category
// @Description Attaches an attribute group to a category
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body models.AttachGroupToCategoryRequest true "Group attachment data"
// @Success 201 {object} object{success=bool,id=int} "Group attached successfully"
// @Failure 400 {object} object{error=string} "marketplace.invalidCategoryId or marketplace.invalidRequest"
// @Failure 500 {object} object{error=string} "marketplace.attachGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/categories/{id}/groups [post]
func (h *MarketplaceHandler) AttachGroupToCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidCategoryId",
		})
	}

	var req models.AttachGroupToCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidRequest",
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
			"error": "marketplace.attachGroupError",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"id":      id,
	})
}

// DetachGroupFromCategory отвязывает группу атрибутов от категории
// @Summary Detach group from category
// @Description Detaches an attribute group from a category
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param groupId path int true "Group ID"
// @Success 200 {object} object{success=bool,message=string} "marketplace.groupDetached"
// @Failure 400 {object} object{error=string} "marketplace.invalidCategoryId or marketplace.invalidGroupId"
// @Failure 500 {object} object{error=string} "marketplace.detachGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/categories/{id}/groups/{groupId} [delete]
func (h *MarketplaceHandler) DetachGroupFromCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	groupIDStr := c.Params("groupId")
	
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidCategoryId",
		})
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "marketplace.invalidGroupId",
		})
	}

	err = h.storage.AttributeGroups.DetachGroupFromCategory(c.Context(), categoryID, groupID)
	if err != nil {
		log.Printf("Ошибка отвязки группы от категории: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "marketplace.detachGroupError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.groupDetached",
	})
}