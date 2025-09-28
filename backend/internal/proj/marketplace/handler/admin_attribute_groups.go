// Package handler
// backend/internal/proj/marketplace/handler/admin_attribute_groups.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"
)

// CreateAttributeGroup создает новую группу атрибутов
// @Summary Create attribute group
// @Description Creates a new attribute group for organizing attributes
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param body body backend_internal_proj_marketplace_models.CreateAttributeGroupRequest true "Attribute group data"
// @Success 201 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Group created successfully"

// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidRequest or marketplace.groupNameRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups [post]
func (h *MarketplaceHandler) CreateAttributeGroup(c *fiber.Ctx) error {
	var req models.CreateAttributeGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Валидация
	if req.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.groupNameRequired")
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
		logger.Error().Err(err).Msg("Error creating attribute group")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createGroupError")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, IDMessageResponse{ID: id, Message: "marketplace.groupCreated"})
}

// ListAttributeGroups возвращает список всех групп атрибутов
// @Summary List attribute groups
// @Description Returns all attribute groups
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=AttributeGroupsResponse} "List of attribute groups"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.listGroupsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups [get]
func (h *MarketplaceHandler) ListAttributeGroups(c *fiber.Ctx) error {
	groups, err := h.storage.AttributeGroups.ListAttributeGroups(c.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Error getting attribute groups list")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.listGroupsError")
	}

	return utils.SuccessResponse(c, AttributeGroupsResponse{
		Groups: groups,
	})
}

// GetAttributeGroup получает информацию о группе атрибутов
// @Summary Get attribute group
// @Description Returns information about a specific attribute group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=AttributeGroupResponse} "Attribute group information"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidGroupId"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.groupNotFound"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id} [get]
func (h *MarketplaceHandler) GetAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	group, err := h.storage.AttributeGroups.GetAttributeGroup(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting attribute group")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.groupNotFound")
	}

	return utils.SuccessResponse(c, AttributeGroupResponse{
		Group: group,
	})
}

// UpdateAttributeGroup обновляет группу атрибутов
// @Summary Update attribute group
// @Description Updates an existing attribute group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param body body backend_internal_proj_marketplace_models.UpdateAttributeGroupRequest true "Updated group data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "marketplace.groupUpdated"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidGroupId or marketplace.invalidRequest"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id} [put]
func (h *MarketplaceHandler) UpdateAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	var req models.UpdateAttributeGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
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
		logger.Error().Err(err).Msg("Error updating attribute group")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateGroupError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.groupUpdated",
	})
}

// DeleteAttributeGroup удаляет группу атрибутов
// @Summary Delete attribute group
// @Description Deletes an attribute group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "marketplace.groupDeleted"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidGroupId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id} [delete]
func (h *MarketplaceHandler) DeleteAttributeGroup(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	err = h.storage.AttributeGroups.DeleteAttributeGroup(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("Error deleting attribute group")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteGroupError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.groupDeleted",
	})
}

// GetAttributeGroupWithItems получает группу с её атрибутами
// @Summary Get attribute group with items
// @Description Returns attribute group with all its attributes
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=AttributeGroupWithItemsData} "Group with items"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidGroupId"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.groupNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getGroupItemsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id}/items [get]
func (h *MarketplaceHandler) GetAttributeGroupWithItems(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	group, err := h.storage.AttributeGroups.GetAttributeGroup(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting attribute group")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.groupNotFound")
	}

	items, err := h.storage.AttributeGroups.GetGroupItems(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting group items")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getGroupItemsError")
	}

	return utils.SuccessResponse(c, AttributeGroupWithItemsData{
		Group: group,
		Items: items,
	})
}

// AddItemToGroup добавляет атрибут в группу
// @Summary Add attribute to group
// @Description Adds an attribute to a group
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param body body backend_internal_proj_marketplace_models.AddItemToGroupRequest true "Attribute data"
// @Success 201 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Item added successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidGroupId or marketplace.invalidRequest"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.addItemError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id}/items [post]
func (h *MarketplaceHandler) AddItemToGroup(c *fiber.Ctx) error {
	groupIDStr := c.Params("id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	var req models.AddItemToGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	item := &models.AttributeGroupItem{
		GroupID:             groupID,
		AttributeID:         req.AttributeID,
		Icon:                req.Icon,
		SortOrder:           req.SortOrder,
		CustomDisplayName:   req.CustomDisplayName,
		VisibilityCondition: req.VisibilityCondition,
	}

	id, err := h.storage.AttributeGroups.AddItemToGroup(c.Context(), groupID, item)
	if err != nil {
		logger.Error().Err(err).Msg("Error adding item to group")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.addItemError")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, IDMessageResponse{
		ID:      id,
		Message: "marketplace.itemAdded",
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
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "marketplace.itemRemoved"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidGroupId or marketplace.invalidAttributeId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.removeItemError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/attribute-groups/{id}/items/{attributeId} [delete]
func (h *MarketplaceHandler) RemoveItemFromGroup(c *fiber.Ctx) error {
	groupIDStr := c.Params("id")
	attributeIDStr := c.Params("attributeId")

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	attributeID, err := strconv.Atoi(attributeIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	err = h.storage.AttributeGroups.RemoveItemFromGroup(c.Context(), groupID, attributeID)
	if err != nil {
		logger.Error().Err(err).Msg("Error removing item from group")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.removeItemError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.itemRemoved",
	})
}

// GetCategoryGroups получает группы атрибутов, привязанные к категории
// @Summary Get category attribute groups
// @Description Returns attribute groups attached to a category
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=AttributeGroupsResponse} "Category groups"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getCategoryGroupsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/categories/{id}/groups [get]
func (h *MarketplaceHandler) GetCategoryGroups(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	categoryGroups, err := h.storage.AttributeGroups.GetCategoryGroups(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting category groups")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getCategoryGroupsError")
	}

	// Извлекаем только группы из CategoryAttributeGroup
	groups := make([]*models.AttributeGroup, 0, len(categoryGroups))
	for _, cg := range categoryGroups {
		if cg.Group != nil {
			groups = append(groups, cg.Group)
		}
	}

	return utils.SuccessResponse(c, AttributeGroupsResponse{
		Groups: groups,
	})
}

// AttachGroupToCategory привязывает группу атрибутов к категории
// @Summary Attach group to category
// @Description Attaches an attribute group to a category
// @Tags marketplace-admin-attribute-groups
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body backend_internal_proj_marketplace_models.AttachGroupToCategoryRequest true "Group attachment data"
// @Success 201 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Group attached successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidRequest"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.attachGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/categories/{id}/groups [post]
func (h *MarketplaceHandler) AttachGroupToCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	var req models.AttachGroupToCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	attachment := &models.CategoryAttributeGroup{
		CategoryID: categoryID,
		GroupID:    req.GroupID,
		SortOrder:  req.SortOrder,
	}

	id, err := h.storage.AttributeGroups.AttachGroupToCategory(c.Context(), categoryID, attachment)
	if err != nil {
		logger.Error().Err(err).Msg("Error attaching group to category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.attachGroupError")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, IDMessageResponse{
		ID:      id,
		Message: "marketplace.groupAttached",
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
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "marketplace.groupDetached"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidGroupId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.detachGroupError"
// @Security BearerAuth
// @Router /api/v1/marketplace/admin/categories/{id}/groups/{groupId} [delete]
func (h *MarketplaceHandler) DetachGroupFromCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("id")
	groupIDStr := c.Params("groupId")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	err = h.storage.AttributeGroups.DetachGroupFromCategory(c.Context(), categoryID, groupID)
	if err != nil {
		logger.Error().Err(err).Msg("Error detaching group from category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.detachGroupError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.groupDetached",
	})
}
