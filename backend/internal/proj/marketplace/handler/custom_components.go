// Package handler
// backend/internal/proj/marketplace/handler/custom_components.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// CustomComponentHandler handles requests for custom components
type CustomComponentHandler struct {
	storage postgres.CustomComponentStorage
}

// NewCustomComponentHandler creates a new handler for custom components
func NewCustomComponentHandler(storage postgres.CustomComponentStorage) *CustomComponentHandler {
	return &CustomComponentHandler{
		storage: storage,
	}
}

// CreateComponent creates a new custom component
// @Summary Create custom component
// @Description Creates a new custom UI component
// @Tags marketplace-admin-custom-components
// @Accept json
// @Produce json
// @Param component body models.CreateCustomComponentRequest true "Component data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=models.CustomUIComponent}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components [post]
func (h *CustomComponentHandler) CreateComponent(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var req models.CreateCustomComponentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	component := &models.CustomUIComponent{
		Name:          req.Name,
		ComponentType: req.ComponentType,
		Description:   req.Description,
		TemplateCode:  req.TemplateCode,
		Styles:        req.Styles,
		PropsSchema:   req.PropsSchema,
		IsActive:      req.IsActive,
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
	}

	id, err := h.storage.CreateComponent(c.Context(), component)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	component.ID = id
	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, component)
}

// GetComponent gets a component by ID
// @Summary Get custom component
// @Description Gets a custom UI component by ID
// @Tags marketplace-admin-custom-components
// @Produce json
// @Param id path int true "Component ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.CustomUIComponent}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.componentNotFound"
// @Security BearerAuth
// @Router /api/admin/custom-components/{id} [get]
func (h *CustomComponentHandler) GetComponent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	component, err := h.storage.GetComponent(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.componentNotFound")
	}

	return utils.SuccessResponse(c, component)
}

// UpdateComponent updates a component
// @Summary Update custom component
// @Description Updates a custom UI component
// @Tags marketplace-admin-custom-components
// @Accept json
// @Produce json
// @Param id path int true "Component ID"
// @Param component body models.UpdateCustomComponentRequest true "Component data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.CustomUIComponent}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidId or marketplace.invalidData"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components/{id} [put]
func (h *CustomComponentHandler) UpdateComponent(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	var req models.UpdateCustomComponentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	updates := map[string]interface{}{
		"updated_by": userID,
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ComponentType != "" {
		updates["component_type"] = req.ComponentType
	}
	if req.TemplateCode != "" {
		updates["template_code"] = req.TemplateCode
	}
	if req.Styles != "" {
		updates["styles"] = req.Styles
	}
	if req.PropsSchema != nil {
		updates["props_schema"] = req.PropsSchema
	}
	updates["is_active"] = req.IsActive

	if err := h.storage.UpdateComponent(c.Context(), id, updates); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	// Получаем обновленный компонент
	component, err := h.storage.GetComponent(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	return utils.SuccessResponse(c, component)
}

// DeleteComponent deletes a component
// @Summary Delete custom component
// @Description Deletes a custom UI component
// @Tags marketplace-admin-custom-components
// @Param id path int true "Component ID"
// @Success 204 "Component deleted"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/admin/custom-components/{id} [delete]
func (h *CustomComponentHandler) DeleteComponent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	if err := h.storage.DeleteComponent(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListComponents returns a list of components
// @Summary List custom components
// @Description Returns a list of all custom UI components
// @Tags marketplace-admin-custom-components
// @Produce json
// @Param component_type query string false "Filter by component type"
// @Param active query bool false "Filter by active status"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.CustomUIComponent}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components [get]
func (h *CustomComponentHandler) ListComponents(c *fiber.Ctx) error {
	logger.Info().Msg("ListComponents called")

	filters := map[string]interface{}{
		"component_type": c.Query("component_type"),
		"active":         c.Query("active"),
	}

	components, err := h.storage.ListComponents(c.Context(), filters)
	if err != nil {
		// Добавляем логирование ошибки
		logger.Error().Err(err).Msg("Error listing components")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	// Логируем успешный результат
	logger.Info().Int("componentsCount", len(components)).Msg("Listed components")

	return utils.SuccessResponse(c, components)
}

// AddComponentUsage adds component usage for a category
// @Summary Add component usage
// @Description Adds component usage for a specific category
// @Tags marketplace-admin-custom-components
// @Accept json
// @Produce json
// @Param usage body models.CreateComponentUsageRequest true "Usage data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=models.CustomUIComponentUsage}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components/usage [post]
func (h *CustomComponentHandler) AddComponentUsage(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var req models.CreateComponentUsageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	usage := &models.CustomUIComponentUsage{
		ComponentID:     req.ComponentID,
		CategoryID:      req.CategoryID,
		UsageContext:    req.UsageContext,
		Placement:       req.Placement,
		Priority:        req.Priority,
		Configuration:   req.Configuration,
		ConditionsLogic: req.ConditionsLogic,
		IsActive:        req.IsActive,
		CreatedBy:       &userID,
		UpdatedBy:       &userID,
	}

	id, err := h.storage.AddComponentUsage(c.Context(), usage)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	usage.ID = id
	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, usage)
}

// GetComponentUsages gets all component usages
// @Summary Get component usages
// @Description Returns a list of all component usages with filtering
// @Tags marketplace-admin-custom-components
// @Produce json
// @Param component_id query int false "Component ID"
// @Param category_id query int false "Category ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.CustomUIComponentUsage}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components/usage [get]
func (h *CustomComponentHandler) GetComponentUsages(c *fiber.Ctx) error {
	var componentID, categoryID *int

	if compID := c.QueryInt("component_id", 0); compID > 0 {
		componentID = &compID
	}

	if catID := c.QueryInt("category_id", 0); catID > 0 {
		categoryID = &catID
	}

	usages, err := h.storage.GetComponentUsages(c.Context(), componentID, categoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	return utils.SuccessResponse(c, usages)
}

// RemoveComponentUsage removes component usage
// @Summary Remove component usage
// @Description Removes component usage for a category
// @Tags marketplace-admin-custom-components
// @Param id path int true "Usage ID"
// @Success 204 "Usage removed"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components/usage/{id} [delete]
func (h *CustomComponentHandler) RemoveComponentUsage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	if err := h.storage.RemoveComponentUsage(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetCategoryComponents returns components for a category
// @Summary Get category components
// @Description Returns custom components for a specific category
// @Tags marketplace-admin-custom-components
// @Produce json
// @Param category_id path int true "Category ID"
// @Param context query string false "Usage context"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.CustomUIComponentUsage}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/categories/{category_id}/components [get]
func (h *CustomComponentHandler) GetCategoryComponents(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	context := c.Query("context")
	components, err := h.storage.GetCategoryComponents(c.Context(), categoryID, context)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	return utils.SuccessResponse(c, components)
}

// CreateTemplate creates a component template
// @Summary Create template
// @Description Creates a template for a custom component
// @Tags marketplace-admin-custom-components
// @Accept json
// @Produce json
// @Param template body models.CreateTemplateRequest true "Template data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=models.ComponentTemplate}
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components/templates [post]
func (h *CustomComponentHandler) CreateTemplate(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var req models.CreateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	template := &models.ComponentTemplate{
		ComponentID:    req.ComponentID,
		Name:           req.Name,
		Description:    req.Description,
		TemplateConfig: req.Variables,
		CreatedBy:      &userID,
	}

	id, err := h.storage.CreateTemplate(c.Context(), template)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	template.ID = id
	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, template)
}

// ListTemplates returns a list of templates
// @Summary List templates
// @Description Returns a list of all component templates
// @Tags marketplace-admin-custom-components
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.ComponentTemplate}
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.internalServerError"
// @Security BearerAuth
// @Router /api/admin/custom-components/templates [get]
func (h *CustomComponentHandler) ListTemplates(c *fiber.Ctx) error {
	logger.Info().Msg("ListTemplates called")

	componentID := c.QueryInt("component_id", 0)
	templates, err := h.storage.ListTemplates(c.Context(), componentID)
	if err != nil {
		// Добавляем логирование ошибки
		logger.Error().Err(err).Msg("Error listing templates")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.internalServerError")
	}

	// Логируем успешный результат
	logger.Info().Int("templatesCount", len(templates)).Msg("Listed templates")

	return utils.SuccessResponse(c, templates)
}
