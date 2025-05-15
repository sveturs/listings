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

	// Маршруты для управления связями атрибутов с категориями
	adminGroup.Post("/categories/:categoryId/attributes/:attributeId", h.AddAttributeToCategory)
	adminGroup.Delete("/categories/:categoryId/attributes/:attributeId", h.RemoveAttributeFromCategory)
	adminGroup.Put("/categories/:categoryId/attributes/:attributeId", h.UpdateAttributeCategory)

	// Маршруты для экспорта/импорта настроек атрибутов
	adminGroup.Get("/categories/:categoryId/attributes/export", h.ExportCategoryAttributes)
	adminGroup.Post("/categories/:categoryId/attributes/import", h.ImportCategoryAttributes)

	// Маршрут для копирования настроек между категориями
	adminGroup.Post("/categories/:targetCategoryId/attributes/copy", h.CopyAttributesSettings)
}

// CreateAttribute создает новый атрибут
func (h *AdminAttributesHandler) CreateAttribute(c *fiber.Ctx) error {
	var attribute models.CategoryAttribute

	// Парсим JSON из запроса
	if err := c.BodyParser(&attribute); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Проверяем обязательные поля
	if attribute.Name == "" || attribute.DisplayName == "" || attribute.AttributeType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не все обязательные поля заполнены")
	}

	// Преобразуем Options из JSON, если они есть
	if c.Get("Content-Type") == "application/json" {
		var optionsData map[string]interface{}
		if err := json.Unmarshal([]byte(c.Body()), &optionsData); err == nil {
			if options, ok := optionsData["options"]; ok {
				optionsJSON, err := json.Marshal(options)
				if err == nil {
					attribute.Options = optionsJSON
				}
			}
		}
	}

	// Создаем атрибут
	id, err := h.marketplaceService.CreateAttribute(c.Context(), &attribute)
	if err != nil {
		log.Printf("Failed to create attribute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось создать атрибут: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"id":      id,
		"message": "Атрибут успешно создан",
	})
}

// GetAttributes получает список всех атрибутов
func (h *AdminAttributesHandler) GetAttributes(c *fiber.Ctx) error {
	// В текущей реализации нет метода для получения всех атрибутов,
	// поэтому можно создать запрос к базе данных напрямую
	query := `
		SELECT id, name, display_name, attribute_type, options, validation_rules, 
		is_searchable, is_filterable, is_required, sort_order, created_at
		FROM category_attributes
		ORDER BY sort_order, id
	`

	rows, err := h.marketplaceService.Storage().Query(c.Context(), query)
	if err != nil {
		log.Printf("Failed to get attributes: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить атрибуты")
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

// GetAttributeByID получает информацию об атрибуте по ID
func (h *AdminAttributesHandler) GetAttributeByID(c *fiber.Ctx) error {
	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Получаем атрибут по ID
	attribute, err := h.marketplaceService.GetAttributeByID(c.Context(), attributeID)
	if err != nil {
		log.Printf("Failed to get attribute by ID: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить атрибут")
	}

	// Если атрибут не найден, возвращаем ошибку
	if attribute == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Атрибут не найден")
	}

	return utils.SuccessResponse(c, attribute)
}

// UpdateAttribute обновляет существующий атрибут
func (h *AdminAttributesHandler) UpdateAttribute(c *fiber.Ctx) error {
	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Парсим JSON из запроса
	var attribute models.CategoryAttribute
	if err := c.BodyParser(&attribute); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Устанавливаем ID атрибута
	attribute.ID = attributeID

	// Проверяем обязательные поля
	if attribute.Name == "" || attribute.DisplayName == "" || attribute.AttributeType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не все обязательные поля заполнены")
	}

	// Преобразуем Options из JSON, если они есть
	if c.Get("Content-Type") == "application/json" {
		var optionsData map[string]interface{}
		if err := json.Unmarshal([]byte(c.Body()), &optionsData); err == nil {
			if options, ok := optionsData["options"]; ok {
				optionsJSON, err := json.Marshal(options)
				if err == nil {
					attribute.Options = optionsJSON
				}
			}
		}
	}

	// Обновляем атрибут
	err = h.marketplaceService.UpdateAttribute(c.Context(), &attribute)
	if err != nil {
		log.Printf("Failed to update attribute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось обновить атрибут: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Атрибут успешно обновлен",
	})
}

// DeleteAttribute удаляет атрибут
func (h *AdminAttributesHandler) DeleteAttribute(c *fiber.Ctx) error {
	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Удаляем атрибут
	err = h.marketplaceService.DeleteAttribute(c.Context(), attributeID)
	if err != nil {
		log.Printf("Failed to delete attribute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось удалить атрибут: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Атрибут успешно удален",
	})
}

// BulkUpdateAttributes массово обновляет атрибуты
func (h *AdminAttributesHandler) BulkUpdateAttributes(c *fiber.Ctx) error {
	// Получаем входные данные
	var input struct {
		Attributes []models.CategoryAttribute `json:"attributes"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	if len(input.Attributes) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Список атрибутов не может быть пустым")
	}

	// Обновляем каждый атрибут
	errors := make([]string, 0)
	success := 0

	for _, attribute := range input.Attributes {
		// Проверяем обязательные поля
		if attribute.ID == 0 || attribute.Name == "" || attribute.DisplayName == "" || attribute.AttributeType == "" {
			errors = append(errors, "Атрибут с неполными данными, ID: "+strconv.Itoa(attribute.ID))
			continue
		}

		// Обновляем атрибут
		err := h.marketplaceService.UpdateAttribute(c.Context(), &attribute)
		if err != nil {
			errors = append(errors, "Не удалось обновить атрибут с ID "+strconv.Itoa(attribute.ID)+": "+err.Error())
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

// AddAttributeToCategory привязывает атрибут к категории
func (h *AdminAttributesHandler) AddAttributeToCategory(c *fiber.Ctx) error {
	// Получаем ID категории и атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	attributeID, err := strconv.Atoi(c.Params("attributeId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Парсим параметры из запроса
	var input struct {
		IsRequired bool `json:"is_required"`
	}

	if err := c.BodyParser(&input); err != nil {
		// Если JSON не передан, устанавливаем значение по умолчанию
		input.IsRequired = false
	}

	// Привязываем атрибут к категории
	err = h.marketplaceService.AddAttributeToCategory(c.Context(), categoryID, attributeID, input.IsRequired)
	if err != nil {
		log.Printf("Failed to add attribute to category: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось привязать атрибут к категории: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Атрибут успешно привязан к категории",
	})
}

// RemoveAttributeFromCategory отвязывает атрибут от категории
func (h *AdminAttributesHandler) RemoveAttributeFromCategory(c *fiber.Ctx) error {
	// Получаем ID категории и атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	attributeID, err := strconv.Atoi(c.Params("attributeId"))
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

// UpdateAttributeCategory обновляет настройки атрибута в категории
func (h *AdminAttributesHandler) UpdateAttributeCategory(c *fiber.Ctx) error {
	// Получаем ID категории и атрибута из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	attributeID, err := strconv.Atoi(c.Params("attributeId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID атрибута")
	}

	// Парсим параметры из запроса
	var input struct {
		IsRequired bool `json:"is_required"`
		IsEnabled  bool `json:"is_enabled"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	// Обновляем настройки атрибута в категории
	err = h.marketplaceService.UpdateAttributeCategory(c.Context(), categoryID, attributeID, input.IsRequired, input.IsEnabled)
	if err != nil {
		log.Printf("Failed to update attribute category settings: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось обновить настройки атрибута в категории: "+err.Error())
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Настройки атрибута в категории успешно обновлены",
	})
}

// ExportCategoryAttributes экспортирует настройки атрибутов категории
func (h *AdminAttributesHandler) ExportCategoryAttributes(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Получаем атрибуты категории с их настройками
	categoryAttributes, err := h.getCategoryAttributesWithSettings(c.Context(), categoryID)
	if err != nil {
		log.Printf("Failed to export category attributes: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось экспортировать настройки атрибутов: "+err.Error())
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
		SELECT cam.category_id, cam.attribute_id, cam.is_enabled, cam.is_required,
			   ca.id, ca.name, ca.display_name, ca.attribute_type, ca.options, ca.validation_rules,
			   ca.is_searchable, ca.is_filterable, ca.is_required, ca.sort_order, ca.created_at, ca.custom_component
		FROM category_attribute_mapping cam
		JOIN category_attributes ca ON cam.attribute_id = ca.id
		WHERE cam.category_id = $1
		ORDER BY ca.sort_order, ca.id
	`

	rows, err := h.marketplaceService.Storage().Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить атрибуты категории: %w", err)
	}
	defer rows.Close()

	result := make([]models.CategoryAttributeMapping, 0)
	for rows.Next() {
		var mapping models.CategoryAttributeMapping
		var attribute models.CategoryAttribute
		var optionsJSON, validRulesJSON []byte

		err := rows.Scan(
			&mapping.CategoryID,
			&mapping.AttributeID,
			&mapping.IsEnabled,
			&mapping.IsRequired,
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
			&attribute.CustomComponent,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать атрибут: %w", err)
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

// ImportCategoryAttributes импортирует настройки атрибутов в категорию
func (h *AdminAttributesHandler) ImportCategoryAttributes(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Проверяем существование категории
	var categoryExists bool
	err = h.marketplaceService.Storage().QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", categoryID).Scan(&categoryExists)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить существование категории")
	}
	if !categoryExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Категория с указанным ID не найдена")
	}

	// Парсим данные из запроса
	var attributeMappings []models.CategoryAttributeMapping
	if err := c.BodyParser(&attributeMappings); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	if len(attributeMappings) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Список атрибутов не может быть пустым")
	}

	// Начинаем транзакцию
	tx, err := h.marketplaceService.Storage().BeginTx(c.Context(), nil)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось начать транзакцию")
	}
	defer tx.Rollback()

	// Удаляем существующие связи категории с атрибутами
	_, err = tx.Exec(c.Context(), "DELETE FROM category_attribute_mapping WHERE category_id = $1", categoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось очистить существующие связи атрибутов")
	}

	// Добавляем новые связи
	successCount := 0
	errors := make([]string, 0)

	for _, mapping := range attributeMappings {
		// Проверяем существование атрибута
		var attributeExists bool
		err = tx.QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM category_attributes WHERE id = $1)", mapping.AttributeID).Scan(&attributeExists)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Не удалось проверить атрибут %d: %s", mapping.AttributeID, err.Error()))
			continue
		}
		if !attributeExists {
			errors = append(errors, fmt.Sprintf("Атрибут с ID %d не существует", mapping.AttributeID))
			continue
		}

		// Добавляем связь
		_, err = tx.Exec(c.Context(), `
			INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
			VALUES ($1, $2, $3, $4)
		`, categoryID, mapping.AttributeID, mapping.IsEnabled, mapping.IsRequired)

		if err != nil {
			errors = append(errors, fmt.Sprintf("Не удалось добавить атрибут %d: %s", mapping.AttributeID, err.Error()))
			continue
		}

		successCount++
	}

	// Если были успешные добавления, фиксируем транзакцию
	if successCount > 0 {
		if err := tx.Commit(); err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось завершить транзакцию: "+err.Error())
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

// CopyAttributesSettings копирует настройки атрибутов из одной категории в другую
func (h *AdminAttributesHandler) CopyAttributesSettings(c *fiber.Ctx) error {
	// Получаем ID целевой категории из параметров URL
	targetCategoryID, err := strconv.Atoi(c.Params("targetCategoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID целевой категории")
	}

	// Получаем ID исходной категории из запроса
	var input struct {
		SourceCategoryID int `json:"source_category_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный формат данных")
	}

	if input.SourceCategoryID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID исходной категории не указан")
	}

	// Проверяем существование обеих категорий
	var targetExists, sourceExists bool
	err = h.marketplaceService.Storage().QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", targetCategoryID).Scan(&targetExists)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить существование целевой категории")
	}
	if !targetExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Целевая категория не найдена")
	}

	err = h.marketplaceService.Storage().QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM marketplace_categories WHERE id = $1)", input.SourceCategoryID).Scan(&sourceExists)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить существование исходной категории")
	}
	if !sourceExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Исходная категория не найдена")
	}

	// Начинаем транзакцию
	tx, err := h.marketplaceService.Storage().BeginTx(c.Context(), nil)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось начать транзакцию")
	}
	defer tx.Rollback()

	// Удаляем существующие связи в целевой категории
	_, err = tx.Exec(c.Context(), "DELETE FROM category_attribute_mapping WHERE category_id = $1", targetCategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось очистить существующие связи атрибутов")
	}

	// Копируем связи из исходной категории в целевую
	query := `
		INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
		SELECT $1, attribute_id, is_enabled, is_required
		FROM category_attribute_mapping
		WHERE category_id = $2
	`
	_, err = tx.Exec(c.Context(), query, targetCategoryID, input.SourceCategoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось скопировать настройки атрибутов: "+err.Error())
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось завершить транзакцию: "+err.Error())
	}

	// Инвалидируем кеш атрибутов для целевой категории
	if err := h.marketplaceService.InvalidateAttributeCache(c.Context(), targetCategoryID); err != nil {
		log.Printf("Failed to invalidate attribute cache: %v", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Настройки атрибутов успешно скопированы",
	})
}
