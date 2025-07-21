// Package handler
// backend/internal/proj/marketplace/handler/admin_categories.go
package handler

import (
	"context"
	"fmt"
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

// GetAllCategories возвращает все категории включая неактивные (для админки)
// @Summary Get all categories including inactive
// @Description Returns all marketplace categories including inactive ones for admin panel
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceCategory} "Categories list"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getCategoriesError"
// @Security BearerAuth
// @Router /api/admin/categories/all [get]
func (h *AdminCategoriesHandler) GetAllCategories(c *fiber.Ctx) error {
	logger.Info().Str("method", c.Method()).Str("path", c.Path()).Msg("GetAllCategories handler called")

	// Получаем язык из query параметра
	lang := c.Query("lang", "en")
	
	// Создаем контекст с языком
	ctx := context.WithValue(c.UserContext(), "locale", lang)
	
	categories, err := h.marketplaceService.GetAllCategories(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get all categories")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getCategoriesError")
	}

	logger.Info().Int("count", len(categories)).Msg("Successfully retrieved all categories")
	return utils.SuccessResponse(c, categories)
}

// CreateCategory создает новую категорию
// @Summary Create category
// @Description Creates a new marketplace category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param body body object{name=string,slug=string,icon=string,parent_id=int,description=string,is_active=bool,seo_title=string,seo_description=string,seo_keywords=string} true "Category data"
// @Success 200 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Category created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData or marketplace.categoryNameRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createCategoryError"
// @Security BearerAuth
// @Router /api/admin/categories [post]
func (h *AdminCategoriesHandler) CreateCategory(c *fiber.Ctx) error {
	logger.Info().Str("method", c.Method()).Str("path", c.Path()).Msg("CreateCategory handler called - START")

	// Парсим JSON из запроса в map для гибкой обработки типов
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	logger.Info().Interface("requestData", requestData).Msg("Parsed request data")

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
	if description, ok := requestData["description"].(string); ok {
		category.Description = description
	}

	// Обрабатываем is_active (по умолчанию true)
	category.IsActive = true
	if isActive, ok := requestData["is_active"].(bool); ok {
		category.IsActive = isActive
	}

	// Обрабатываем SEO поля
	if seoTitle, ok := requestData["seo_title"].(string); ok {
		category.SEOTitle = seoTitle
	}
	if seoDescription, ok := requestData["seo_description"].(string); ok {
		category.SEODescription = seoDescription
	}
	if seoKeywords, ok := requestData["seo_keywords"].(string); ok {
		category.SEOKeywords = seoKeywords
	}

	// Обрабатываем переводы
	if translations, ok := requestData["translations"].(map[string]interface{}); ok {
		category.Translations = make(map[string]string)
		for lang, trans := range translations {
			if transStr, ok := trans.(string); ok {
				category.Translations[lang] = transStr
			}
		}
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
	logger.Info().Interface("category", category).Msg("About to create category via service")
	id, err := h.marketplaceService.CreateCategory(c.Context(), &category)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createCategoryError")
	}

	logger.Info().Int("categoryId", id).Msg("Category created successfully")

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
// @Param body body object{name=string,slug=string,icon=string,parent_id=int,description=string,is_active=bool,seo_title=string,seo_description=string,seo_keywords=string} true "Updated category data"
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
	if description, ok := requestData["description"].(string); ok {
		category.Description = description
	}
	if isActive, ok := requestData["is_active"].(bool); ok {
		category.IsActive = isActive
	}

	// Обрабатываем SEO поля
	if seoTitle, ok := requestData["seo_title"].(string); ok {
		category.SEOTitle = seoTitle
	}
	if seoDescription, ok := requestData["seo_description"].(string); ok {
		category.SEODescription = seoDescription
	}
	if seoKeywords, ok := requestData["seo_keywords"].(string); ok {
		category.SEOKeywords = seoKeywords
	}

	// Обрабатываем переводы
	if translations, ok := requestData["translations"].(map[string]interface{}); ok {
		category.Translations = make(map[string]string)
		for lang, trans := range translations {
			if transStr, ok := trans.(string); ok {
				category.Translations[lang] = transStr
			}
		}
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

// GetCategoryAttributeGroups получает группы атрибутов, привязанные к категории
// @Summary Get category attribute groups
// @Description Returns attribute groups attached to a category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.AttributeGroup} "Category groups"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getCategoryGroupsError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/groups [get]
func (h *AdminCategoriesHandler) GetCategoryAttributeGroups(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем группы через MarketplaceHandler
	// Поскольку AdminCategoriesHandler включает CategoriesHandler, но не имеет прямого доступа к storage.AttributeGroups,
	// мы можем добавить метод в CategoriesHandler или использовать прямой вызов к сервису
	groups, err := h.marketplaceService.GetCategoryAttributeGroups(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get category attribute groups")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getCategoryGroupsError")
	}

	return utils.SuccessResponse(c, groups)
}

// AttachAttributeGroupToCategory привязывает группу атрибутов к категории
// @Summary Attach attribute group to category
// @Description Attaches an attribute group to a category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body object{group_id=int,sort_order=int} true "Group attachment data"
// @Success 201 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Group attached successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.attachGroupError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/groups [post]
func (h *AdminCategoriesHandler) AttachAttributeGroupToCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем входные данные
	var input struct {
		GroupID   int `json:"group_id"`
		SortOrder int `json:"sort_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Привязываем группу к категории
	id, err := h.marketplaceService.AttachAttributeGroupToCategory(c.Context(), categoryID, input.GroupID, input.SortOrder)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to attach attribute group to category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.attachGroupError")
	}

	return utils.SuccessResponse(c, IDMessageResponse{
		ID:      id,
		Message: "marketplace.groupAttachedToCategory",
	})
}

// DetachAttributeGroupFromCategory отвязывает группу атрибутов от категории
// @Summary Detach attribute group from category
// @Description Detaches an attribute group from a category
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param group_id path int true "Group ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Group detached successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId or marketplace.invalidGroupId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.detachGroupError"
// @Security BearerAuth
// @Router /api/admin/categories/{id}/groups/{group_id} [delete]
func (h *AdminCategoriesHandler) DetachAttributeGroupFromCategory(c *fiber.Ctx) error {
	// Получаем ID категории и ID группы из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	groupID, err := strconv.Atoi(c.Params("group_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidGroupId")
	}

	// Отвязываем группу от категории
	err = h.marketplaceService.DetachAttributeGroupFromCategory(c.Context(), categoryID, groupID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to detach attribute group from category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.detachGroupError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.groupDetachedFromCategory",
	})
}

// TranslateCategory automatically translates category name, description and SEO fields
// @Summary Auto-translate category
// @Description Automatically translates category name, description, seo_title and seo_description to all supported languages using Google Translate
// @Tags marketplace-admin-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param languages body object{source_language=string,target_languages=[]string} false "Translation settings"
// @Success 200 {object} utils.SuccessResponseSwag{data=CategoryTranslationResult} "Translation results"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid category ID"
// @Failure 404 {object} utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Translation error"
// @Security BearerAuth
// @Router /api/v1/admin/categories/{id}/translate [post]
func (h *AdminCategoriesHandler) TranslateCategory(c *fiber.Ctx) error {
	logger.Info().Msg("TranslateCategory method called")

	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем категорию по ID
	category, err := h.marketplaceService.GetCategoryByID(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get category by ID")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getCategoryError")
	}

	if category == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.categoryNotFound")
	}

	// Парсим настройки перевода
	var input struct {
		SourceLanguage  string   `json:"source_language"`
		TargetLanguages []string `json:"target_languages"`
	}

	// Значения по умолчанию
	input.SourceLanguage = "en"
	input.TargetLanguages = []string{"ru", "sr"}

	// Если есть тело запроса, парсим его
	if err := c.BodyParser(&input); err == nil {
		// Проверяем валидность языков
		if input.SourceLanguage == "" {
			input.SourceLanguage = "en"
		}
		if len(input.TargetLanguages) == 0 {
			input.TargetLanguages = []string{"ru", "sr"}
		}
	}

	// Результаты перевода
	translationResults := make(map[string]any)
	errors := make([]string, 0)

	// Поля для перевода
	fieldsToTranslate := map[string]string{
		"name":            category.Name,
		"description":     category.Description,
		"seo_title":       category.SEOTitle,
		"seo_description": category.SEODescription,
	}

	// Переводим каждое поле
	for fieldName, fieldValue := range fieldsToTranslate {
		if fieldValue == "" {
			continue
		}

		fieldTranslations := make(map[string]string)
		for _, targetLang := range input.TargetLanguages {
			if targetLang == input.SourceLanguage {
				continue
			}

			translatedText, err := h.marketplaceService.TranslateText(c.Context(), fieldValue, input.SourceLanguage, targetLang)
			if err != nil {
				logger.Error().Err(err).Str("field", fieldName).Str("text", fieldValue).Str("target_lang", targetLang).Msg("Failed to translate field")
				errors = append(errors, fmt.Sprintf("Failed to translate %s to %s", fieldName, targetLang))
				continue
			}
			fieldTranslations[targetLang] = translatedText
		}

		// Сохраняем переводы для поля
		for lang, text := range fieldTranslations {
			err := h.marketplaceService.SaveTranslation(c.Context(), "category", categoryID, lang, fieldName, text, nil)
			if err != nil {
				logger.Error().Err(err).Str("field", fieldName).Msg("Failed to save translation")
				errors = append(errors, fmt.Sprintf("Failed to save %s translation for %s", fieldName, lang))
			}
		}

		if len(fieldTranslations) > 0 {
			translationResults[fieldName] = fieldTranslations
		}
	}

	// Формируем результат
	result := CategoryTranslationResult{
		CategoryID:   categoryID,
		Translations: translationResults,
		Errors:       errors,
	}

	if len(errors) > 0 {
		c.Status(fiber.StatusPartialContent)
	}

	return utils.SuccessResponse(c, result)
}
