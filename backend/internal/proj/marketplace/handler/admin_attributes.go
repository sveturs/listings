// backend/internal/proj/marketplace/handler/admin_attributes.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// AdminAttributesHandler обрабатывает запросы админки для управления атрибутами
type AdminAttributesHandler struct {
	*CategoriesHandler
}

// NewAdminAttributesHandler создает новый обработчик админки для атрибутов
func NewAdminAttributesHandler(services globalService.ServicesInterface) *AdminAttributesHandler {
	return &AdminAttributesHandler{
		CategoriesHandler: NewCategoriesHandler(services),
	}
}

// RegisterRoutes регистрирует маршруты для админки атрибутов
func (h *AdminAttributesHandler) RegisterRoutes(app *fiber.App, adminMiddleware fiber.Handler) {
	adminGroup := app.Group("/api/admin", adminMiddleware)

	// Маршруты для управления атрибутами
	adminGroup.Post("/attributes", h.CreateAttribute)
	adminGroup.Get("/attributes", h.GetAttributes)
	adminGroup.Get("/attributes/:id", h.GetAttributeByID)
	adminGroup.Put("/attributes/:id", h.UpdateAttribute)
	adminGroup.Delete("/attributes/:id", h.DeleteAttribute)
	adminGroup.Post("/attributes/bulk-update", h.BulkUpdateAttributes)

	// Маршруты для экспорта/импорта настроек атрибутов - регистрируем до других маршрутов категорий и атрибутов
	// чтобы избежать конфликта маршрутизации
	adminGroup.Get("/categories/:categoryId/attributes/export", h.ExportCategoryAttributes)
	adminGroup.Post("/categories/:categoryId/attributes/import", h.ImportCategoryAttributes)
	
	// Маршрут для копирования настроек между категориями
	adminGroup.Post("/categories/:targetCategoryId/attributes/copy", h.CopyAttributesSettings)

	// Маршруты для управления связями атрибутов с категориями
	adminGroup.Post("/categories/:categoryId/attributes/:attributeId", h.AddAttributeToCategory)
	adminGroup.Delete("/categories/:categoryId/attributes/:attributeId", h.RemoveAttributeFromCategory)
	adminGroup.Put("/categories/:categoryId/attributes/:attributeId", h.UpdateAttributeCategory)
}

// CreateAttribute создает новый атрибут
// @Summary Create attribute
// @Description Creates a new attribute for categories
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param body body models.CategoryAttribute true "Attribute data"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{id=int,message=string}} "marketplace.attributeCreated"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData or marketplace.requiredFieldsMissing"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createAttributeError"
// @Security BearerAuth
// @Router /api/admin/attributes [post]
func (h *AdminAttributesHandler) CreateAttribute(c *fiber.Ctx) error {
	var attribute models.CategoryAttribute

	// Парсим JSON из запроса
	if err := c.BodyParser(&attribute); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if attribute.Name == "" || attribute.DisplayName == "" || attribute.AttributeType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.requiredFieldsMissing")
	}

	// Обрабатываем полные данные из JSON запроса
	if c.Get("Content-Type") == "application/json" {
		var requestData map[string]interface{}
		if err := json.Unmarshal([]byte(c.Body()), &requestData); err == nil {
			// Преобразуем Options, если они есть
			if options, ok := requestData["options"]; ok {
				optionsJSON, err := json.Marshal(options)
				if err == nil {
					attribute.Options = optionsJSON
				}
			}

			// Формируем validation_rules из отдельных полей
			validationRules := make(map[string]interface{})
			
			// Проверяем и добавляем правила валидации из запроса
			if minLength, ok := requestData["min_length"]; ok && minLength != nil {
				if val, ok := minLength.(float64); ok && val > 0 {
					validationRules["min_length"] = val
				}
			}
			if maxLength, ok := requestData["max_length"]; ok && maxLength != nil {
				if val, ok := maxLength.(float64); ok && val > 0 {
					validationRules["max_length"] = val
				}
			}
			if minValue, ok := requestData["min_value"]; ok && minValue != nil {
				if val, ok := minValue.(float64); ok {
					validationRules["min_value"] = val
				}
			}
			if maxValue, ok := requestData["max_value"]; ok && maxValue != nil {
				if val, ok := maxValue.(float64); ok {
					validationRules["max_value"] = val
				}
			}
			if pattern, ok := requestData["pattern"]; ok && pattern != nil {
				if val, ok := pattern.(string); ok && val != "" {
					validationRules["pattern"] = val
				}
			}
			if unit, ok := requestData["unit"]; ok && unit != nil {
				if val, ok := unit.(string); ok && val != "" {
					validationRules["unit"] = val
				}
			}
			if defaultValue, ok := requestData["default_value"]; ok && defaultValue != nil {
				validationRules["default_value"] = defaultValue
			}

			// Обрабатываем поле icon
			if icon, ok := requestData["icon"]; ok && icon != nil {
				if val, ok := icon.(string); ok {
					attribute.Icon = val
				}
			}

			// Преобразуем validation_rules в JSON, если есть правила
			if len(validationRules) > 0 {
				validRulesJSON, err := json.Marshal(validationRules)
				if err == nil {
					attribute.ValidRules = validRulesJSON
				}
			}
		}
	}

	// Создаем атрибут
	id, err := h.marketplaceService.CreateAttribute(c.Context(), &attribute)
	if err != nil {
		log.Printf("Failed to create attribute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createAttributeError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"message": "marketplace.attributeCreated",
	})
}

// GetAttributes returns list of all attributes
// @Summary Get all attributes
// @Description Returns list of all category attributes sorted by sort_order and ID
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Success 200 {array} models.CategoryAttribute "List of attributes"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes [get]
func (h *AdminAttributesHandler) GetAttributes(c *fiber.Ctx) error {
	// В текущей реализации нет метода для получения всех атрибутов,
	// поэтому можно создать запрос к базе данных напрямую
	query := `
		SELECT id, name, display_name, attribute_type, icon, options, validation_rules, 
		is_searchable, is_filterable, is_required, sort_order, created_at
		FROM category_attributes
		ORDER BY sort_order, id
	`

	rows, err := h.marketplaceService.Storage().Query(c.Context(), query)
	if err != nil {
		log.Printf("Failed to get attributes: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getAttributesError")
	}
	defer rows.Close()

	attributes := make([]models.CategoryAttribute, 0)
	for rows.Next() {
		var attribute models.CategoryAttribute
		var optionsJSON, validRulesJSON []byte

		err := rows.Scan(
			&attribute.ID,
			&attribute.Name,
			&attribute.DisplayName,
			&attribute.AttributeType,
			&attribute.Icon,
			&optionsJSON,
			&validRulesJSON,
			&attribute.IsSearchable,
			&attribute.IsFilterable,
			&attribute.IsRequired,
			&attribute.SortOrder,
			&attribute.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan attribute: %v", err)
			continue
		}

		attribute.Options = optionsJSON
		attribute.ValidRules = validRulesJSON
		attributes = append(attributes, attribute)
	}

	return utils.SuccessResponse(c, attributes)
}

// GetAttributeByID returns attribute information by ID
// @Summary Get attribute by ID
// @Description Returns detailed information about a specific attribute
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} models.CategoryAttribute "Attribute information"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid attribute ID"
// @Failure 404 {object} utils.ErrorResponseSwag "Attribute not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes/{id} [get]
func (h *AdminAttributesHandler) GetAttributeByID(c *fiber.Ctx) error {
	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Получаем атрибут по ID
	attribute, err := h.marketplaceService.GetAttributeByID(c.Context(), attributeID)
	if err != nil {
		log.Printf("Failed to get attribute by ID: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getAttributeError")
	}

	// Если атрибут не найден, возвращаем ошибку
	if attribute == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.attributeNotFound")
	}

	return utils.SuccessResponse(c, attribute)
}

// UpdateAttribute updates an existing attribute
// @Summary Update attribute
// @Description Updates an existing category attribute with new data
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param attribute body models.CategoryAttribute true "Attribute data"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid data"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes/{id} [put]
func (h *AdminAttributesHandler) UpdateAttribute(c *fiber.Ctx) error {
	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Парсим JSON из запроса
	var attribute models.CategoryAttribute
	if err := c.BodyParser(&attribute); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Устанавливаем ID атрибута
	attribute.ID = attributeID

	// Проверяем обязательные поля
	if attribute.Name == "" || attribute.DisplayName == "" || attribute.AttributeType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.requiredFieldsMissing")
	}

	// Обрабатываем полные данные из JSON запроса
	if c.Get("Content-Type") == "application/json" {
		var requestData map[string]interface{}
		if err := json.Unmarshal([]byte(c.Body()), &requestData); err == nil {
			// Преобразуем Options, если они есть
			if options, ok := requestData["options"]; ok {
				optionsJSON, err := json.Marshal(options)
				if err == nil {
					attribute.Options = optionsJSON
				}
			}

			// Формируем validation_rules из отдельных полей
			validationRules := make(map[string]interface{})
			
			// Проверяем и добавляем правила валидации из запроса
			if minLength, ok := requestData["min_length"]; ok && minLength != nil {
				if val, ok := minLength.(float64); ok && val > 0 {
					validationRules["min_length"] = val
				}
			}
			if maxLength, ok := requestData["max_length"]; ok && maxLength != nil {
				if val, ok := maxLength.(float64); ok && val > 0 {
					validationRules["max_length"] = val
				}
			}
			if minValue, ok := requestData["min_value"]; ok && minValue != nil {
				if val, ok := minValue.(float64); ok {
					validationRules["min_value"] = val
				}
			}
			if maxValue, ok := requestData["max_value"]; ok && maxValue != nil {
				if val, ok := maxValue.(float64); ok {
					validationRules["max_value"] = val
				}
			}
			if pattern, ok := requestData["pattern"]; ok && pattern != nil {
				if val, ok := pattern.(string); ok && val != "" {
					validationRules["pattern"] = val
				}
			}
			if unit, ok := requestData["unit"]; ok && unit != nil {
				if val, ok := unit.(string); ok && val != "" {
					validationRules["unit"] = val
				}
			}
			if defaultValue, ok := requestData["default_value"]; ok && defaultValue != nil {
				validationRules["default_value"] = defaultValue
			}

			// Обрабатываем поле icon
			if icon, ok := requestData["icon"]; ok && icon != nil {
				if val, ok := icon.(string); ok {
					attribute.Icon = val
				}
			}

			// Преобразуем validation_rules в JSON, если есть правила
			if len(validationRules) > 0 {
				validRulesJSON, err := json.Marshal(validationRules)
				if err == nil {
					attribute.ValidRules = validRulesJSON
				}
			}
		}
	}

	// Обновляем атрибут
	err = h.marketplaceService.UpdateAttribute(c.Context(), &attribute)
	if err != nil {
		log.Printf("Failed to update attribute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateAttributeError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attributeUpdated",
	})
}

// DeleteAttribute deletes an attribute
// @Summary Delete attribute
// @Description Deletes a category attribute by ID
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid attribute ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes/{id} [delete]
func (h *AdminAttributesHandler) DeleteAttribute(c *fiber.Ctx) error {
	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Удаляем атрибут
	err = h.marketplaceService.DeleteAttribute(c.Context(), attributeID)
	if err != nil {
		log.Printf("Failed to delete attribute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteAttributeError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attributeDeleted",
	})
}

// BulkUpdateAttributes updates multiple attributes in batch
// @Summary Bulk update attributes
// @Description Updates multiple category attributes in a single request
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param attributes body object{attributes=[]models.CategoryAttribute} true "List of attributes to update"
// @Success 200 {object} map[string]interface{} "Update results with success count"
// @Success 206 {object} map[string]interface{} "Partial update with errors"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid data"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes/bulk [put]
func (h *AdminAttributesHandler) BulkUpdateAttributes(c *fiber.Ctx) error {
	// Получаем входные данные
	var input struct {
		Attributes []models.CategoryAttribute `json:"attributes"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	if len(input.Attributes) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.emptyAttributesList")
	}

	// Обновляем каждый атрибут
	errors := make([]string, 0)
	success := 0

	for _, attribute := range input.Attributes {
		// Проверяем обязательные поля
		if attribute.ID == 0 || attribute.Name == "" || attribute.DisplayName == "" || attribute.AttributeType == "" {
			errors = append(errors, "marketplace.incompleteAttributeData")
			continue
		}

		// Обновляем атрибут
		err := h.marketplaceService.UpdateAttribute(c.Context(), &attribute)
		if err != nil {
			errors = append(errors, "marketplace.updateAttributeError")
			continue
		}

		success++
	}

	result := fiber.Map{
		"success_count": success,
		"total_count":   len(input.Attributes),
	}

	if len(errors) > 0 {
		result["errors"] = errors
		return c.Status(fiber.StatusPartialContent).JSON(result)
	}

	return utils.SuccessResponse(c, result)
}

// AddAttributeToCategory links an attribute to a category
// @Summary Add attribute to category
// @Description Links a category attribute to a specific category with optional settings
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Param attributeId path int true "Attribute ID"
// @Param settings body object{is_required=bool,sort_order=int} false "Attribute settings for category"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/categories/{categoryId}/attributes/{attributeId} [post]
func (h *AdminAttributesHandler) AddAttributeToCategory(c *fiber.Ctx) error {
	// Получаем ID категории и атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	attributeID, err := strconv.Atoi(c.Params("attributeId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Парсим параметры из запроса
	var input struct {
		IsRequired bool `json:"is_required"`
		SortOrder  int  `json:"sort_order"`
	}

	if err := c.BodyParser(&input); err != nil {
		// Если JSON не передан, устанавливаем значение по умолчанию
		input.IsRequired = false
		input.SortOrder = 0
	}

	// Если указан порядок сортировки, используем расширенный метод
	var addErr error
	if input.SortOrder > 0 {
		addErr = h.marketplaceService.AddAttributeToCategoryWithOrder(c.Context(), categoryID, attributeID, input.IsRequired, input.SortOrder)
	} else {
		// Иначе используем обычный метод (порядок будет взят из атрибута)
		addErr = h.marketplaceService.AddAttributeToCategory(c.Context(), categoryID, attributeID, input.IsRequired)
	}

	if addErr != nil {
		log.Printf("Failed to add attribute to category: %v", addErr)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.addAttributeToCategoryError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attributeAddedToCategory",
	})
}

// RemoveAttributeFromCategory unlinks an attribute from a category
// @Summary Remove attribute from category
// @Description Removes the link between a category attribute and a category
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Param attributeId path int true "Attribute ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/categories/{categoryId}/attributes/{attributeId} [delete]
func (h *AdminAttributesHandler) RemoveAttributeFromCategory(c *fiber.Ctx) error {
	// Получаем ID категории и атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	attributeID, err := strconv.Atoi(c.Params("attributeId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Отвязываем атрибут от категории
	err = h.marketplaceService.RemoveAttributeFromCategory(c.Context(), categoryID, attributeID)
	if err != nil {
		log.Printf("Failed to remove attribute from category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.removeAttributeFromCategoryError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attributeRemovedFromCategory",
	})
}

// UpdateAttributeCategory updates attribute settings in a category
// @Summary Update attribute category settings
// @Description Updates settings for an attribute within a specific category
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Param attributeId path int true "Attribute ID"
// @Param settings body object{is_required=bool,is_enabled=bool,sort_order=int,custom_component=string} true "Attribute category settings"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/categories/{categoryId}/attributes/{attributeId} [put]
func (h *AdminAttributesHandler) UpdateAttributeCategory(c *fiber.Ctx) error {
	// Получаем ID категории и атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	attributeID, err := strconv.Atoi(c.Params("attributeId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Парсим параметры из запроса
	var input struct {
		IsRequired      bool   `json:"is_required"`
		IsEnabled       bool   `json:"is_enabled"`
		SortOrder       int    `json:"sort_order"`
		CustomComponent string `json:"custom_component"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Обновляем настройки атрибута в категории
	err = h.marketplaceService.UpdateAttributeCategoryExtended(
		c.Context(),
		categoryID,
		attributeID,
		input.IsRequired,
		input.IsEnabled,
		input.SortOrder,
		input.CustomComponent,
	)

	if err != nil {
		log.Printf("Failed to update attribute category settings: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateAttributeCategoryError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attributeCategoryUpdated",
	})
}

// ExportCategoryAttributes exports category attribute settings
// @Summary Export category attributes
// @Description Exports all attribute settings for a specific category as JSON
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Success 200 {array} models.CategoryAttributeMapping "Category attributes with settings"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid category ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/categories/{categoryId}/attributes/export [get]
func (h *AdminAttributesHandler) ExportCategoryAttributes(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем атрибуты категории с их настройками
	categoryAttributes, err := h.getCategoryAttributesWithSettings(c.Context(), categoryID)
	if err != nil {
		log.Printf("Failed to export category attributes: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.exportAttributesError")
	}

	// Устанавливаем заголовок для скачивания файла
	categoryName := "category_" + strconv.Itoa(categoryID)
	// Получаем информацию о категории для более точного имени файла
	category, err := h.marketplaceService.GetCategoryByID(c.Context(), categoryID)
	if err == nil && category != nil {
		categoryName = category.Name
	}

	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_attributes.json", categoryName))
	c.Set("Content-Type", "application/json")

	return c.JSON(categoryAttributes)
}

// Вспомогательный метод для получения атрибутов категории с настройками
func (h *AdminAttributesHandler) getCategoryAttributesWithSettings(ctx context.Context, categoryID int) ([]models.CategoryAttributeMapping, error) {
	// Запрос для получения атрибутов категории с их настройками
	query := `
		SELECT cam.category_id, cam.attribute_id, cam.is_enabled, cam.is_required, cam.sort_order, 
			   COALESCE(cam.custom_component, '') as mapping_custom_component,
			   ca.id, ca.name, ca.display_name, ca.attribute_type, ca.options, ca.validation_rules,
			   ca.is_searchable, ca.is_filterable, ca.is_required, ca.sort_order, ca.created_at, 
			   COALESCE(ca.custom_component, '') as attribute_custom_component
		FROM category_attribute_mapping cam
		JOIN category_attributes ca ON cam.attribute_id = ca.id
		WHERE cam.category_id = $1
		ORDER BY cam.sort_order, ca.sort_order, ca.id
	`

	rows, err := h.marketplaceService.Storage().Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("marketplace.getCategoryAttributesError: %w", err)
	}
	defer rows.Close()

	result := make([]models.CategoryAttributeMapping, 0)
	for rows.Next() {
		var mapping models.CategoryAttributeMapping
		var attribute models.CategoryAttribute
		var optionsJSON, validRulesJSON []byte
		var mappingCustomComponent, attributeCustomComponent string

		err := rows.Scan(
			&mapping.CategoryID,
			&mapping.AttributeID,
			&mapping.IsEnabled,
			&mapping.IsRequired,
			&mapping.SortOrder,
			&mappingCustomComponent,
			&attribute.ID,
			&attribute.Name,
			&attribute.DisplayName,
			&attribute.AttributeType,
			&optionsJSON,
			&validRulesJSON,
			&attribute.IsSearchable,
			&attribute.IsFilterable,
			&attribute.IsRequired,
			&attribute.SortOrder,
			&attribute.CreatedAt,
			&attributeCustomComponent,
		)
		if err != nil {
			return nil, fmt.Errorf("marketplace.readAttributeError: %w", err)
		}

		// Используем пользовательский компонент из маппинга, если он есть, иначе из атрибута
		if mappingCustomComponent != "" {
			attribute.CustomComponent = mappingCustomComponent
		} else {
			attribute.CustomComponent = attributeCustomComponent
		}

		// Устанавливаем Options и ValidRules как json.RawMessage
		attribute.Options = optionsJSON
		attribute.ValidRules = validRulesJSON

		// Получаем переводы для атрибута
		translationsQuery := `
			SELECT language, field_name, translated_text
			FROM translations
			WHERE entity_type = 'attribute' AND entity_id = $1 AND field_name = 'display_name'
		`
		tRows, err := h.marketplaceService.Storage().Query(ctx, translationsQuery, attribute.ID)
		if err == nil {
			attribute.Translations = make(map[string]string)
			for tRows.Next() {
				var lang, field, text string
				if err := tRows.Scan(&lang, &field, &text); err == nil {
					attribute.Translations[lang] = text
				}
			}
			tRows.Close()
		}

		// Получаем переводы для опций атрибута
		optionTranslationsQuery := `
			SELECT language, field_name, translated_text
			FROM translations
			WHERE entity_type = 'attribute_option' AND entity_id = $1
		`
		oRows, err := h.marketplaceService.Storage().Query(ctx, optionTranslationsQuery, attribute.ID)
		if err == nil {
			attribute.OptionTranslations = make(map[string]map[string]string)
			for oRows.Next() {
				var lang, option, text string
				if err := oRows.Scan(&lang, &option, &text); err == nil {
					if attribute.OptionTranslations[lang] == nil {
						attribute.OptionTranslations[lang] = make(map[string]string)
					}
					attribute.OptionTranslations[lang][option] = text
				}
			}
			oRows.Close()
		}

		mapping.Attribute = &attribute
		result = append(result, mapping)
	}

	return result, nil
}

// ImportCategoryAttributes imports attribute settings into a category
// @Summary Import category attributes
// @Description Imports attribute settings into a category, replacing existing settings
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Param attributes body []models.CategoryAttributeMapping true "List of attribute mappings to import"
// @Success 200 {object} map[string]interface{} "Import results with success count"
// @Success 206 {object} map[string]interface{} "Partial import with errors"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid data"
// @Failure 404 {object} utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/categories/{categoryId}/attributes/import [post]
func (h *AdminAttributesHandler) ImportCategoryAttributes(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Проверяем существование категории
	var categoryExists bool
	err = h.marketplaceService.Storage().QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", categoryID).Scan(&categoryExists)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.checkCategoryExistenceError")
	}
	if !categoryExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.categoryNotFound")
	}

	// Парсим данные из запроса
	var attributeMappings []models.CategoryAttributeMapping
	if err := c.BodyParser(&attributeMappings); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	if len(attributeMappings) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.emptyAttributesList")
	}

	// Начинаем транзакцию
	tx, err := h.marketplaceService.Storage().BeginTx(c.Context(), nil)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.beginTransactionError")
	}
	defer tx.Rollback()

	// Удаляем существующие связи категории с атрибутами
	_, err = tx.Exec(c.Context(), "DELETE FROM category_attribute_mapping WHERE category_id = $1", categoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.clearAttributeMappingsError")
	}

	// Добавляем новые связи
	successCount := 0
	errors := make([]string, 0)

	for _, mapping := range attributeMappings {
		// Проверяем существование атрибута
		var attributeExists bool
		err = tx.QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM category_attributes WHERE id = $1)", mapping.AttributeID).Scan(&attributeExists)
		if err != nil {
			errors = append(errors, "marketplace.checkAttributeError")
			continue
		}
		if !attributeExists {
			errors = append(errors, "marketplace.attributeNotExists")
			continue
		}

		// Добавляем связь с учетом sort_order и custom_component
		_, err = tx.Exec(c.Context(), `
			INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order, custom_component)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, categoryID, mapping.AttributeID, mapping.IsEnabled, mapping.IsRequired, mapping.SortOrder, mapping.CustomComponent)

		if err != nil {
			errors = append(errors, "marketplace.addAttributeError")
			continue
		}

		successCount++
	}

	// Если были успешные добавления, фиксируем транзакцию
	if successCount > 0 {
		if err := tx.Commit(); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.commitTransactionError")
		}

		// Инвалидируем кеш атрибутов для категории
		if err := h.marketplaceService.InvalidateAttributeCache(c.Context(), categoryID); err != nil {
			log.Printf("Failed to invalidate attribute cache: %v", err)
		}
	}

	result := fiber.Map{
		"success_count": successCount,
		"total_count":   len(attributeMappings),
	}

	if len(errors) > 0 {
		result["errors"] = errors
		if successCount == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(result)
		}
		return c.Status(fiber.StatusPartialContent).JSON(result)
	}

	return utils.SuccessResponse(c, result)
}

// CopyAttributesSettings copies attribute settings from one category to another
// @Summary Copy attribute settings between categories
// @Description Copies all attribute settings from a source category to a target category
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param targetCategoryId path int true "Target category ID"
// @Param source body object{source_category_id=int} true "Source category ID"
// @Success 200 {object} map[string]interface{} "Success message"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 404 {object} utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/categories/{targetCategoryId}/attributes/copy [post]
func (h *AdminAttributesHandler) CopyAttributesSettings(c *fiber.Ctx) error {
	// Получаем ID целевой категории из параметров URL
	targetCategoryID, err := strconv.Atoi(c.Params("targetCategoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidTargetCategoryId")
	}

	// Получаем ID исходной категории из запроса
	var input struct {
		SourceCategoryID int `json:"source_category_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	if input.SourceCategoryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.sourceCategoryIdRequired")
	}

	// Проверяем существование обеих категорий
	var targetExists, sourceExists bool
	err = h.marketplaceService.Storage().QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", targetCategoryID).Scan(&targetExists)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.checkTargetCategoryExistenceError")
	}
	if !targetExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.targetCategoryNotFound")
	}

	err = h.marketplaceService.Storage().QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", input.SourceCategoryID).Scan(&sourceExists)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.checkSourceCategoryExistenceError")
	}
	if !sourceExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.sourceCategoryNotFound")
	}

	// Начинаем транзакцию
	tx, err := h.marketplaceService.Storage().BeginTx(c.Context(), nil)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.beginTransactionError")
	}
	defer tx.Rollback()

	// Удаляем существующие связи в целевой категории
	_, err = tx.Exec(c.Context(), "DELETE FROM category_attribute_mapping WHERE category_id = $1", targetCategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.clearAttributeMappingsError")
	}

	// Копируем связи из исходной категории в целевую с учетом sort_order и custom_component
	query := `
		INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order, custom_component)
		SELECT $1, attribute_id, is_enabled, is_required, sort_order, custom_component
		FROM category_attribute_mapping
		WHERE category_id = $2
	`
	_, err = tx.Exec(c.Context(), query, targetCategoryID, input.SourceCategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.copyAttributeSettingsError")
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.commitTransactionError")
	}

	// Инвалидируем кеш атрибутов для целевой категории
	if err := h.marketplaceService.InvalidateAttributeCache(c.Context(), targetCategoryID); err != nil {
		log.Printf("Failed to invalidate attribute cache: %v", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attributeSettingsCopied",
	})
}
