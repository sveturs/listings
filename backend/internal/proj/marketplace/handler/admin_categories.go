// Package handler
// backend/internal/proj/marketplace/handler/admin_categories.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"
)

// AdminCategoriesHandler обрабатывает запросы админки для управления категориями
type AdminCategoriesHandler struct {
	*CategoriesHandler
}

// NewAdminCategoriesHandler создает новый обработчик админки для категорий
func NewAdminCategoriesHandler(categoriesHandler *CategoriesHandler) *AdminCategoriesHandler {
	return &AdminCategoriesHandler{
		CategoriesHandler: categoriesHandler,
	}
}

// CreateCategory создает новую категорию
// @Summary Create category
// @Description Creates a new marketplace category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param body body object{name=string,slug=string,icon=string,parent_id=int} true "Category data"
// @Success 200 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Category created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData or marketplace.categoryNameRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories [post]
func (h *AdminCategoriesHandler) CreateCategory(c *fiber.Ctx) error {
	// Парсим JSON из запроса в map для гибкой обработки типов
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Создаем структуру категории
	var category models.MarketplaceCategory

	// Обрабатываем основные поля
	if name, ok := requestData["name"].(string); ok {
		category.Name = name
	}
	if slug, ok := requestData["slug"].(string); ok {
		category.Slug = slug
	}
	if icon, ok := requestData["icon"].(string); ok {
		category.Icon = icon
	}

	// Обрабатываем parent_id - может прийти как строка или число
	if parentIDRaw, ok := requestData["parent_id"]; ok && parentIDRaw != nil {
		switch v := parentIDRaw.(type) {
		case string:
			if v != "" && v != "0" {
				if parentID, err := strconv.Atoi(v); err == nil && parentID > 0 {
					category.ParentID = &parentID
				}
			}
		case float64:
			if v > 0 {
				parentID := int(v)
				category.ParentID = &parentID
			}
		case int:
			if v > 0 {
				category.ParentID = &v
			}
		}
	}

	// Проверяем обязательные поля
	if category.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.categoryNameRequired")
	}

	// Если slug не указан, генерируем его из названия
	if category.Slug == "" {
		category.Slug = utils.GenerateSlug(category.Name)
	}

	// Создаем категорию
	id, err := h.marketplaceService.CreateCategory(c.Context(), &category)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createCategoryError")
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	// Возвращаем ID созданной категории
	return utils.SuccessResponse(c, IDMessageResponse{
		ID:      id,
		Message: "marketplace.categoryCreated",
	})
}

// GetCategoryByID получает информацию о категории по ID
// @Summary Get category by ID
// @Description Returns detailed information about a specific category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceCategory} "Category information"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.categoryNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getCategoriesError"
// @Security BearerAuth
// @Router /api/admin/categories/{id} [get]
func (h *AdminCategoriesHandler) GetCategoryByID(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем все категории
	categories, err := h.marketplaceService.GetCategories(c.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get categories")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getCategoriesError")
	}

	// Ищем нужную категорию
	var category *models.MarketplaceCategory
	for i := range categories {
		if categories[i].ID == categoryID {
			category = &categories[i]
			break
		}
	}

	// Если категория не найдена, возвращаем ошибку
	if category == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.categoryNotFound")
	}

	// Возвращаем информацию о категории
	return utils.SuccessResponse(c, category)
}

// UpdateCategory обновляет существующую категорию
// @Summary Update category
// @Description Updates an existing marketplace category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body object{name=string,slug=string,icon=string,parent_id=int} true "Updated category data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Category updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.categoryNameRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories/{id} [put]
func (h *AdminCategoriesHandler) UpdateCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Парсим JSON из запроса в map для гибкой обработки типов
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Создаем структуру категории
	var category models.MarketplaceCategory
	category.ID = categoryID

	// Обрабатываем основные поля
	if name, ok := requestData["name"].(string); ok {
		category.Name = name
	}
	if slug, ok := requestData["slug"].(string); ok {
		category.Slug = slug
	}
	if icon, ok := requestData["icon"].(string); ok {
		category.Icon = icon
	}

	// Обрабатываем parent_id - может прийти как строка или число
	if parentIDRaw, ok := requestData["parent_id"]; ok && parentIDRaw != nil {
		switch v := parentIDRaw.(type) {
		case string:
			if v != "" && v != "0" {
				if parentID, err := strconv.Atoi(v); err == nil && parentID > 0 {
					category.ParentID = &parentID
				}
			}
		case float64:
			if v > 0 {
				parentID := int(v)
				category.ParentID = &parentID
			}
		case int:
			if v > 0 {
				category.ParentID = &v
			}
		}
	}

	// Проверяем обязательные поля
	if category.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.categoryNameRequired")
	}

	// Если slug не указан, генерируем его из названия
	if category.Slug == "" {
		category.Slug = utils.GenerateSlug(category.Name)
	}

	// Обновляем категорию
	err = h.marketplaceService.UpdateCategory(c.Context(), &category)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateCategoryError")
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.categoryUpdated",
	})
}

// DeleteCategory удаляет категорию
// @Summary Delete category
// @Description Deletes a marketplace category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Category deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories/{id} [delete]
func (h *AdminCategoriesHandler) DeleteCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Удаляем категорию
	err = h.marketplaceService.DeleteCategory(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteCategoryError")
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.categoryDeleted",
	})
}

// ReorderCategories изменяет порядок категорий
// @Summary Reorder categories
// @Description Changes the order of categories based on provided IDs list
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Parent category ID"
// @Param body body object{ordered_ids=[]int} true "List of category IDs in new order"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Categories reordered successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData or marketplace.emptyIdList"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.reorderCategoriesError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/reorder [post]
func (h *AdminCategoriesHandler) ReorderCategories(c *fiber.Ctx) error {
	// Получаем входные данные
	var input struct {
		OrderedIDs []int `json:"ordered_ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	if len(input.OrderedIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.emptyIdList")
	}

	// Изменяем порядок категорий
	err := h.marketplaceService.ReorderCategories(c.Context(), input.OrderedIDs)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to reorder categories")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.reorderCategoriesError")
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.categoriesReordered",
	})
}

// MoveCategory перемещает категорию в иерархии
// @Summary Move category in hierarchy
// @Description Moves a category to a different parent in the hierarchy
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID to move"
// @Param body body object{new_parent_id=int} true "New parent category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Category moved successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.moveCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/move [put]
func (h *AdminCategoriesHandler) MoveCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем входные данные
	var input struct {
		NewParentID int `json:"new_parent_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Перемещаем категорию
	err = h.marketplaceService.MoveCategory(c.Context(), categoryID, input.NewParentID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to move category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.moveCategoryError")
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.categoryMoved",
	})
}

// AddAttributeToCategory привязывает атрибут к категории
// @Summary Add attribute to category
// @Description Links an attribute to a category with optional required setting
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body object{attribute_id=int,is_required=bool} true "Attribute data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Attribute added to category successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.addAttributeToCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/attributes [post]
func (h *AdminCategoriesHandler) AddAttributeToCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем входные данные
	var input struct {
		AttributeID int  `json:"attribute_id"`
		IsRequired  bool `json:"is_required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Привязываем атрибут к категории
	err = h.marketplaceService.AddAttributeToCategory(c.Context(), categoryID, input.AttributeID, input.IsRequired)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to add attribute to category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.addAttributeToCategoryError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeAddedToCategory",
	})
}

// RemoveAttributeFromCategory отвязывает атрибут от категории
// @Summary Remove attribute from category
// @Description Unlinks an attribute from a category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param attr_id path int true "Attribute ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Attribute removed from category successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidAttributeId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.removeAttributeFromCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/attributes/{attr_id} [delete]
func (h *AdminCategoriesHandler) RemoveAttributeFromCategory(c *fiber.Ctx) error {
	// Получаем ID категории и ID атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	attributeID, err := strconv.Atoi(c.Params("attr_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Отвязываем атрибут от категории
	err = h.marketplaceService.RemoveAttributeFromCategory(c.Context(), categoryID, attributeID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to remove attribute from category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.removeAttributeFromCategoryError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeRemovedFromCategory",
	})
}

// UpdateAttributeCategory обновляет настройки связи атрибута с категорией
// @Summary Update attribute category settings
// @Description Updates settings for an attribute-category relationship
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param attr_id path int true "Attribute ID"
// @Param body body object{is_required=bool,is_enabled=bool} true "Attribute settings"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Attribute category updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidAttributeId or marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateAttributeCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/attributes/{attr_id} [put]
func (h *AdminCategoriesHandler) UpdateAttributeCategory(c *fiber.Ctx) error {
	// Получаем ID категории и ID атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	attributeID, err := strconv.Atoi(c.Params("attr_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Получаем входные данные
	var input struct {
		IsRequired bool `json:"is_required"`
		IsEnabled  bool `json:"is_enabled"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Обновляем настройки связи
	err = h.marketplaceService.UpdateAttributeCategory(c.Context(), categoryID, attributeID, input.IsRequired, input.IsEnabled)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update attribute category settings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateAttributeCategoryError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeCategoryUpdated",
	})
}
