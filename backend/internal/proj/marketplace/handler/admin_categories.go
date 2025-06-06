// backend/internal/proj/marketplace/handler/admin_categories.go
package handler

import (
	"backend/internal/domain/models"
	//"backend/internal/middleware"
	"backend/pkg/utils"
	//"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
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

// RegisterRoutes регистрирует маршруты для админки категорий
func (h *AdminCategoriesHandler) RegisterRoutes(app *fiber.App, adminMiddleware fiber.Handler) {
	adminGroup := app.Group("/api/admin", adminMiddleware)

	// Маршруты для управления категориями
	adminGroup.Post("/categories", h.CreateCategory)
	adminGroup.Get("/categories", h.GetCategories)
	adminGroup.Get("/categories/:id", h.GetCategoryByID)
	adminGroup.Put("/categories/:id", h.UpdateCategory)
	adminGroup.Delete("/categories/:id", h.DeleteCategory)
	adminGroup.Post("/categories/:id/reorder", h.ReorderCategories)
	adminGroup.Put("/categories/:id/move", h.MoveCategory)

	// Маршруты для управления связями категорий и атрибутов
	adminGroup.Post("/categories/:id/attributes", h.AddAttributeToCategory)
	adminGroup.Delete("/categories/:id/attributes/:attr_id", h.RemoveAttributeFromCategory)
	adminGroup.Put("/categories/:id/attributes/:attr_id", h.UpdateAttributeCategory)
}

// CreateCategory создает новую категорию
func (h *AdminCategoriesHandler) CreateCategory(c *fiber.Ctx) error {
	// Парсим JSON из запроса в map для гибкой обработки типов
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Название категории не может быть пустым")
	}

	// Если slug не указан, генерируем его из названия
	if category.Slug == "" {
		category.Slug = utils.GenerateSlug(category.Name)
	}

	// Создаем категорию
	id, err := h.marketplaceService.CreateCategory(c.Context(), &category)
	if err != nil {
		log.Printf("Failed to create category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось создать категорию: "+err.Error())
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	// Возвращаем ID созданной категории
	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"message": "Категория успешно создана",
	})
}

// GetCategoryByID получает информацию о категории по ID
func (h *AdminCategoriesHandler) GetCategoryByID(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Получаем все категории
	categories, err := h.marketplaceService.GetCategories(c.Context())
	if err != nil {
		log.Printf("Failed to get categories: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить категории")
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
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Категория не найдена")
	}

	// Возвращаем информацию о категории
	return utils.SuccessResponse(c, category)
}

// UpdateCategory обновляет существующую категорию
func (h *AdminCategoriesHandler) UpdateCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Парсим JSON из запроса в map для гибкой обработки типов
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Название категории не может быть пустым")
	}

	// Если slug не указан, генерируем его из названия
	if category.Slug == "" {
		category.Slug = utils.GenerateSlug(category.Name)
	}

	// Обновляем категорию
	err = h.marketplaceService.UpdateCategory(c.Context(), &category)
	if err != nil {
		log.Printf("Failed to update category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось обновить категорию: "+err.Error())
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Категория успешно обновлена",
	})
}

// DeleteCategory удаляет категорию
func (h *AdminCategoriesHandler) DeleteCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Удаляем категорию
	err = h.marketplaceService.DeleteCategory(c.Context(), categoryID)
	if err != nil {
		log.Printf("Failed to delete category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось удалить категорию: "+err.Error())
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Категория успешно удалена",
	})
}

// ReorderCategories изменяет порядок категорий
func (h *AdminCategoriesHandler) ReorderCategories(c *fiber.Ctx) error {
	// Получаем входные данные
	var input struct {
		OrderedIDs []int `json:"ordered_ids"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	if len(input.OrderedIDs) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Список идентификаторов не может быть пустым")
	}

	// Изменяем порядок категорий
	err := h.marketplaceService.ReorderCategories(c.Context(), input.OrderedIDs)
	if err != nil {
		log.Printf("Failed to reorder categories: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось изменить порядок категорий")
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Порядок категорий успешно изменен",
	})
}

// MoveCategory перемещает категорию в иерархии
func (h *AdminCategoriesHandler) MoveCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Получаем входные данные
	var input struct {
		NewParentID int `json:"new_parent_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Перемещаем категорию
	err = h.marketplaceService.MoveCategory(c.Context(), categoryID, input.NewParentID)
	if err != nil {
		log.Printf("Failed to move category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось переместить категорию: "+err.Error())
	}

	// Инвалидируем кеш категорий
	h.InvalidateCategoryCache()

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Категория успешно перемещена",
	})
}

// AddAttributeToCategory привязывает атрибут к категории
func (h *AdminCategoriesHandler) AddAttributeToCategory(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Получаем входные данные
	var input struct {
		AttributeID int  `json:"attribute_id"`
		IsRequired  bool `json:"is_required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Привязываем атрибут к категории
	err = h.marketplaceService.AddAttributeToCategory(c.Context(), categoryID, input.AttributeID, input.IsRequired)
	if err != nil {
		log.Printf("Failed to add attribute to category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось привязать атрибут к категории: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Атрибут успешно привязан к категории",
	})
}

// RemoveAttributeFromCategory отвязывает атрибут от категории
func (h *AdminCategoriesHandler) RemoveAttributeFromCategory(c *fiber.Ctx) error {
	// Получаем ID категории и ID атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	attributeID, err := strconv.Atoi(c.Params("attr_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Отвязываем атрибут от категории
	err = h.marketplaceService.RemoveAttributeFromCategory(c.Context(), categoryID, attributeID)
	if err != nil {
		log.Printf("Failed to remove attribute from category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось отвязать атрибут от категории: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Атрибут успешно отвязан от категории",
	})
}

// UpdateAttributeCategory обновляет настройки связи атрибута с категорией
func (h *AdminCategoriesHandler) UpdateAttributeCategory(c *fiber.Ctx) error {
	// Получаем ID категории и ID атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	attributeID, err := strconv.Atoi(c.Params("attr_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Получаем входные данные
	var input struct {
		IsRequired bool `json:"is_required"`
		IsEnabled  bool `json:"is_enabled"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Обновляем настройки связи
	err = h.marketplaceService.UpdateAttributeCategory(c.Context(), categoryID, attributeID, input.IsRequired, input.IsEnabled)
	if err != nil {
		log.Printf("Failed to update attribute category settings: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось обновить настройки связи атрибута с категорией: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Настройки связи успешно обновлены",
	})
}
