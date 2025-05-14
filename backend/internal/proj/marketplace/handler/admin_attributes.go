// backend/internal/proj/marketplace/handler/admin_attributes.go
package handler

import (
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// AdminAttributesHandler обрабатывает запросы админки для управления атрибутами
type AdminAttributesHandler struct {
	*CategoriesHandler
}

// NewAdminAttributesHandler создает новый обработчик админки для атрибутов
func NewAdminAttributesHandler(categoriesHandler *CategoriesHandler) *AdminAttributesHandler {
	return &AdminAttributesHandler{
		CategoriesHandler: categoriesHandler,
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
		is_searchable, is_filterable, is_required, sort_order, created_at, custom_component
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
			&attribute.CustomComponent,
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
