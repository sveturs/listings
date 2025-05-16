package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	postgres "backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// CustomComponentHandler обрабатывает запросы для работы с кастомными компонентами
type CustomComponentHandler struct {
	storage postgres.CustomComponentStorage
}

// NewCustomComponentHandler создает новый обработчик для кастомных компонентов
func NewCustomComponentHandler(storage postgres.CustomComponentStorage) *CustomComponentHandler {
	return &CustomComponentHandler{
		storage: storage,
	}
}

// CreateComponent создает новый кастомный компонент
// @Summary Создание кастомного компонента
// @Description Создает новый кастомный UI компонент
// @Tags CustomComponents
// @Accept json
// @Produce json
// @Param component body models.CreateCustomComponentRequest true "Данные компонента"
// @Success 201 {object} models.CustomUIComponent
// @Failure 400 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components [post]
func (h *CustomComponentHandler) CreateComponent(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var req models.CreateCustomComponentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	component := &models.CustomUIComponent{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		ComponentType: req.ComponentType,
		ComponentCode: req.ComponentCode,
		Configuration: req.Configuration,
		Dependencies:  req.Dependencies,
		IsActive:      req.IsActive,
		CreatedBy:     &userID,
		UpdatedBy:     &userID,
	}

	id, err := h.storage.CreateComponent(c.Context(), component)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	component.ID = id
	return c.Status(fiber.StatusCreated).JSON(component)
}

// GetComponent получает компонент по ID
// @Summary Получение компонента
// @Description Получает кастомный UI компонент по ID
// @Tags CustomComponents
// @Produce json
// @Param id path int true "ID компонента"
// @Success 200 {object} models.CustomUIComponent
// @Failure 404 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components/{id} [get]
func (h *CustomComponentHandler) GetComponent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID")
	}

	component, err := h.storage.GetComponent(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Компонент не найден")
	}

	return c.JSON(component)
}

// UpdateComponent обновляет компонент
// @Summary Обновление компонента
// @Description Обновляет кастомный UI компонент
// @Tags CustomComponents
// @Accept json
// @Produce json
// @Param id path int true "ID компонента"
// @Param component body models.UpdateCustomComponentRequest true "Данные компонента"
// @Success 200 {object} models.CustomUIComponent
// @Failure 400 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components/{id} [put]
func (h *CustomComponentHandler) UpdateComponent(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID")
	}

	var req models.UpdateCustomComponentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	updates := map[string]interface{}{
		"updated_by": userID,
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ComponentType != "" {
		updates["component_type"] = req.ComponentType
	}
	if req.ComponentCode != "" {
		updates["component_code"] = req.ComponentCode
	}
	if req.Configuration != nil {
		updates["configuration"] = req.Configuration
	}
	if req.Dependencies != nil {
		updates["dependencies"] = req.Dependencies
	}
	updates["is_active"] = req.IsActive

	if err := h.storage.UpdateComponent(c.Context(), id, updates); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Получаем обновленный компонент
	component, err := h.storage.GetComponent(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(component)
}

// DeleteComponent удаляет компонент
// @Summary Удаление компонента
// @Description Удаляет кастомный UI компонент
// @Tags CustomComponents
// @Param id path int true "ID компонента"
// @Success 204 "Компонент удален"
// @Failure 400 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components/{id} [delete]
func (h *CustomComponentHandler) DeleteComponent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID")
	}

	if err := h.storage.DeleteComponent(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListComponents возвращает список компонентов
// @Summary Список компонентов
// @Description Возвращает список всех кастомных UI компонентов
// @Tags CustomComponents
// @Produce json
// @Param component_type query string false "Фильтр по типу компонента"
// @Param active query bool false "Фильтр по активности"
// @Success 200 {array} models.CustomUIComponent
// @Security BearerAuth
// @Router /api/admin/custom-components [get]
func (h *CustomComponentHandler) ListComponents(c *fiber.Ctx) error {
	log.Printf("ListComponents called")
	
	filters := map[string]interface{}{
		"component_type": c.Query("component_type"),
		"active":         c.Query("active"),
	}

	components, err := h.storage.ListComponents(c.Context(), filters)
	if err != nil {
		// Добавляем логирование ошибки
		log.Printf("Error listing components: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Логируем успешный результат
	log.Printf("Listed %d components", len(components))
	
	return c.JSON(components)
}

// AddComponentUsage добавляет использование компонента для категории
// @Summary Добавление использования компонента
// @Description Добавляет использование компонента для конкретной категории
// @Tags CustomComponents
// @Accept json
// @Produce json
// @Param usage body models.CreateComponentUsageRequest true "Данные использования"
// @Success 201 {object} models.CustomUIComponentUsage
// @Failure 400 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components/usage [post]
func (h *CustomComponentHandler) AddComponentUsage(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var req models.CreateComponentUsageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	usage.ID = id
	return c.Status(fiber.StatusCreated).JSON(usage)
}

// RemoveComponentUsage удаляет использование компонента
// @Summary Удаление использования компонента
// @Description Удаляет использование компонента для категории
// @Tags CustomComponents
// @Param id path int true "ID использования"
// @Success 204 "Использование удалено"
// @Failure 400 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components/usage/{id} [delete]
func (h *CustomComponentHandler) RemoveComponentUsage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID")
	}

	if err := h.storage.RemoveComponentUsage(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetCategoryComponents возвращает компоненты для категории
// @Summary Компоненты категории
// @Description Возвращает кастомные компоненты для конкретной категории
// @Tags CustomComponents
// @Produce json
// @Param category_id path int true "ID категории"
// @Param context query string false "Контекст использования"
// @Success 200 {array} models.CustomUIComponentUsage
// @Security BearerAuth
// @Router /api/admin/categories/{category_id}/components [get]
func (h *CustomComponentHandler) GetCategoryComponents(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID категории")
	}

	context := c.Query("context")
	components, err := h.storage.GetCategoryComponents(c.Context(), categoryID, context)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(components)
}

// CreateTemplate создает шаблон компонента
// @Summary Создание шаблона
// @Description Создает шаблон для кастомного компонента
// @Tags CustomComponents
// @Accept json
// @Produce json
// @Param template body models.CreateTemplateRequest true "Данные шаблона"
// @Success 201 {object} models.CustomUITemplate
// @Failure 400 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/admin/custom-components/templates [post]
func (h *CustomComponentHandler) CreateTemplate(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(int)

	var req models.CreateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	template := &models.CustomUITemplate{
		Name:         req.Name,
		Description:  req.Description,
		TemplateCode: req.TemplateCode,
		Variables:    req.Variables,
		CreatedBy:    &userID,
		UpdatedBy:    &userID,
	}

	id, err := h.storage.CreateTemplate(c.Context(), template)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	template.ID = id
	return c.Status(fiber.StatusCreated).JSON(template)
}

// ListTemplates возвращает список шаблонов
// @Summary Список шаблонов
// @Description Возвращает список всех шаблонов компонентов
// @Tags CustomComponents
// @Produce json
// @Success 200 {array} models.CustomUITemplate
// @Security BearerAuth
// @Router /api/admin/custom-components/templates [get]
func (h *CustomComponentHandler) ListTemplates(c *fiber.Ctx) error {
	log.Printf("ListTemplates called")
	
	templates, err := h.storage.ListTemplates(c.Context())
	if err != nil {
		// Добавляем логирование ошибки
		log.Printf("Error listing templates: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Логируем успешный результат
	log.Printf("Listed %d templates", len(templates))
	
	return c.JSON(templates)
}