// Package handler
// backend/internal/proj/marketplace/handler/admin_attributes.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
)

const (
	queryOptionTranslations = `
			SELECT language, field_name, translated_text
			FROM translations
			WHERE entity_type = 'attribute_option' AND entity_id = $1
		`
	attributeTypeSelect = "select"
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

// CreateAttribute создает новый атрибут
// @Summary Create attribute
// @Description Creates a new attribute for categories
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param body body backend_internal_domain_models.CategoryAttribute true "Attribute data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=AttributeCreateResponse} "marketplace.attributeCreated"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.invalidData or marketplace.requiredFieldsMissing"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.createAttributeError"
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

			// Обрабатываем переводы display_name
			if translations, ok := requestData["translations"].(map[string]interface{}); ok {
				attribute.Translations = make(map[string]string)
				for lang, trans := range translations {
					if transStr, ok := trans.(string); ok {
						attribute.Translations[lang] = transStr
					}
				}
			}

			// Обрабатываем переводы опций
			if optionTranslations, ok := requestData["option_translations"].(map[string]interface{}); ok {
				attribute.OptionTranslations = make(map[string]map[string]string)
				for lang, options := range optionTranslations {
					if optionsMap, ok := options.(map[string]interface{}); ok {
						attribute.OptionTranslations[lang] = make(map[string]string)
						for optKey, optValue := range optionsMap {
							if optStr, ok := optValue.(string); ok {
								attribute.OptionTranslations[lang][optKey] = optStr
							}
						}
					}
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
		logger.Error().Err(err).Msg("Failed to create attribute")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createAttributeError")
	}

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, AttributeCreateResponse{
		ID:      id,
		Message: "marketplace.attributeCreated",
	})
}

// GetAttributes returns paginated list of attributes with search and filter
// @Summary Get all attributes with pagination, search and filter
// @Description Returns paginated list of all category attributes sorted by sort_order and ID with optional search and filter
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Param search query string false "Search term for name or display_name"
// @Param type query string false "Filter by attribute type"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_domain_models.PaginatedResponse} "Paginated list of attributes"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid pagination parameters"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes [get]
func (h *AdminAttributesHandler) GetAttributes(c *fiber.Ctx) error {
	// Получаем параметры пагинации
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 20)
	searchTerm := c.Query("search", "")
	filterType := c.Query("type", "")

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Вычисляем offset
	offset := (page - 1) * pageSize

	// Строим запросы с учетом фильтров
	whereConditions := []string{}
	queryParams := []interface{}{}
	paramCounter := 1

	// Добавляем условие поиска
	if searchTerm != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(LOWER(name) LIKE LOWER($%d) OR LOWER(display_name) LIKE LOWER($%d))", paramCounter, paramCounter+1))
		searchPattern := "%" + searchTerm + "%"
		queryParams = append(queryParams, searchPattern, searchPattern)
		paramCounter += 2
	}

	// Добавляем фильтр по типу
	if filterType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("attribute_type = $%d", paramCounter))
		queryParams = append(queryParams, filterType)
		paramCounter++
	}

	// Формируем WHERE часть запроса
	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Получаем общее количество атрибутов с учетом фильтров
	var total int
	countQuery := `SELECT COUNT(*) FROM category_attributes` + whereClause
	err := h.marketplaceService.Storage().QueryRow(c.Context(), countQuery, queryParams...).Scan(&total)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to count attributes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getAttributesError")
	}

	// Добавляем параметры для LIMIT и OFFSET
	queryParams = append(queryParams, pageSize, offset)

	// Получаем атрибуты с пагинацией и фильтрами
	query := fmt.Sprintf(`
		SELECT id, name, display_name, attribute_type, COALESCE(icon, '') as icon, options, validation_rules,
		is_searchable, is_filterable, is_required, sort_order, created_at, COALESCE(custom_component, '') as custom_component,
		is_variant_compatible, affects_stock
		FROM category_attributes
		%s
		ORDER BY sort_order, id
		LIMIT $%d OFFSET $%d
	`, whereClause, paramCounter, paramCounter+1)

	rows, err := h.marketplaceService.Storage().Query(c.Context(), query, queryParams...)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get attributes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getAttributesError")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

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
			&attribute.CustomComponent,
			&attribute.IsVariantCompatible,
			&attribute.AffectsStock,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan attribute")
			continue
		}

		attribute.Options = optionsJSON
		attribute.ValidRules = validRulesJSON

		// Загружаем переводы для атрибута
		translationsQuery := `
			SELECT language, field_name, translated_text
			FROM translations
			WHERE entity_type = 'attribute' AND entity_id = $1 AND field_name = 'display_name'
		`
		tRows, err := h.marketplaceService.Storage().Query(c.Context(), translationsQuery, attribute.ID)
		if err == nil {
			attribute.Translations = make(map[string]string)
			for tRows.Next() {
				var lang, field, text string
				if err := tRows.Scan(&lang, &field, &text); err == nil {
					attribute.Translations[lang] = text
				}
			}
			if err := tRows.Close(); err != nil {
				logger.Error().Err(err).Msg("Failed to close translation rows")
			}
		}

		// Загружаем переводы для опций атрибута
		optionTranslationsQuery := queryOptionTranslations
		oRows, err := h.marketplaceService.Storage().Query(c.Context(), optionTranslationsQuery, attribute.ID)
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
			if err := oRows.Close(); err != nil {
				logger.Error().Err(err).Msg("Failed to close option translation rows")
			}
		}

		attributes = append(attributes, attribute)
	}

	// Вычисляем общее количество страниц
	totalPages := (total + pageSize - 1) / pageSize

	// Формируем пагинированный ответ
	response := models.PaginatedResponse{
		Data:       attributes,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return utils.SuccessResponse(c, response)
}

// GetAttributeByID returns attribute information by ID
// @Summary Get attribute by ID
// @Description Returns detailed information about a specific attribute
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_domain_models.CategoryAttribute} "Attribute information"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid attribute ID"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Attribute not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
		logger.Error().Err(err).Msg("Failed to get attribute by ID")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getAttributeError")
	}

	// Если атрибут не найден, возвращаем ошибку
	if attribute == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.attributeNotFound")
	}

	// Загружаем переводы для атрибута
	translationsQuery := `
		SELECT language, field_name, translated_text
		FROM translations
		WHERE entity_type = 'attribute' AND entity_id = $1 AND field_name = 'display_name'
	`
	tRows, err := h.marketplaceService.Storage().Query(c.Context(), translationsQuery, attribute.ID)
	if err == nil {
		attribute.Translations = make(map[string]string)
		for tRows.Next() {
			var lang, field, text string
			if err := tRows.Scan(&lang, &field, &text); err == nil {
				attribute.Translations[lang] = text
			}
		}
		if err := tRows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close translation rows")
		}
	}

	// Загружаем переводы опций, если атрибут типа select
	if attribute.AttributeType == attributeTypeSelect && attribute.Options != nil {
		optionTranslationsQuery := queryOptionTranslations
		oRows, err := h.marketplaceService.Storage().Query(c.Context(), optionTranslationsQuery, attribute.ID)
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
			if err := oRows.Close(); err != nil {
				logger.Error().Err(err).Msg("Failed to close option translation rows")
			}
		}
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
// @Param attribute body backend_internal_domain_models.CategoryAttribute true "Attribute data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Success message"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid data"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes/{id} [put]
func (h *AdminAttributesHandler) UpdateAttribute(c *fiber.Ctx) error {
	logger.Info().Msg("UpdateAttribute method called")

	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		logger.Error().Err(err).Msg("Invalid attribute ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	logger.Info().Int("attributeID", attributeID).Msg("Updating attribute")

	// Парсим JSON из запроса
	var attribute models.CategoryAttribute
	if err := c.BodyParser(&attribute); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	logger.Info().Interface("attribute", attribute).Msg("Parsed attribute data")

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

			// Обрабатываем переводы display_name
			if translations, ok := requestData["translations"].(map[string]interface{}); ok {
				attribute.Translations = make(map[string]string)
				for lang, trans := range translations {
					if transStr, ok := trans.(string); ok {
						attribute.Translations[lang] = transStr
					}
				}
			}

			// Обрабатываем переводы опций
			if optionTranslations, ok := requestData["option_translations"].(map[string]interface{}); ok {
				attribute.OptionTranslations = make(map[string]map[string]string)
				for lang, options := range optionTranslations {
					if optionsMap, ok := options.(map[string]interface{}); ok {
						attribute.OptionTranslations[lang] = make(map[string]string)
						for optKey, optValue := range optionsMap {
							if optStr, ok := optValue.(string); ok {
								attribute.OptionTranslations[lang][optKey] = optStr
							}
						}
					}
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
		logger.Error().Err(err).Msg("Failed to update attribute")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateAttributeError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeUpdated",
	})
}

// DeleteAttribute deletes an attribute
// @Summary Delete attribute
// @Description Deletes a category attribute by ID
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "marketplace.attributeDeleted"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid attribute ID"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
		logger.Error().Err(err).Msg("Failed to delete attribute")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteAttributeError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeDeleted",
	})
}

// BulkUpdateAttributes updates multiple attributes in batch
// @Summary Bulk update attributes
// @Description Updates multiple category attributes in a single request
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param attributes body object{attributes=[]backend_internal_domain_models.CategoryAttribute} true "List of attributes to update"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=BulkUpdateResult} "Update results with success count"
// @Success 206 {object} PartialOperationResponse{data=BulkUpdateResult} "Partial update with errors"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid data"
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

	result := BulkUpdateResult{
		SuccessCount: success,
		TotalCount:   len(input.Attributes),
	}

	if len(errors) > 0 {
		result.Errors = errors
		return c.Status(fiber.StatusPartialContent).JSON(PartialOperationResponse{
			Success: false,
			Error:   "marketplace.partialUpdateCompleted",
			Data:    result,
		})
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Success message"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
		logger.Error().Err(addErr).Msg("Failed to add attribute to category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.addAttributeToCategoryError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeAddedToCategory",
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Success message"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
		logger.Error().Err(err).Msg("Failed to remove attribute from category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.removeAttributeFromCategoryError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeRemovedFromCategory",
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Success message"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
		logger.Error().Err(err).Msg("Failed to update attribute category settings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateAttributeCategoryError")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeCategoryUpdated",
	})
}

// ExportCategoryAttributes exports category attribute settings
// @Summary Export category attributes
// @Description Exports all attribute settings for a specific category as JSON
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param categoryId path int true "Category ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_domain_models.CategoryAttributeMapping} "Category attributes with settings"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid category ID"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
		logger.Error().Err(err).Msg("Failed to export category attributes")
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
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

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
			if err := tRows.Close(); err != nil {
				logger.Error().Err(err).Msg("Failed to close translation rows")
			}
		}

		// Получаем переводы для опций атрибута
		optionTranslationsQuery := queryOptionTranslations
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
			if err := oRows.Close(); err != nil {
				logger.Error().Err(err).Msg("Failed to close option translation rows")
			}
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
// @Param attributes body []backend_internal_domain_models.CategoryAttributeMapping true "List of attribute mappings to import"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=ImportAttributesResult} "Import results with success count"
// @Success 206 {object} PartialOperationResponse{data=ImportAttributesResult} "Partial import with errors"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid data"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Error().Err(err).Msg("Failed to rollback transaction")
		}
	}()

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
			logger.Error().Err(err).Msg("Failed to invalidate attribute cache")
		}
	}

	result := ImportAttributesResult{
		SuccessCount: successCount,
		TotalCount:   len(attributeMappings),
	}

	if len(errors) > 0 {
		result.Errors = errors
		if successCount == 0 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.importAttributesFailed")
		}
		return c.Status(fiber.StatusPartialContent).JSON(PartialOperationResponse{
			Success: false,
			Error:   "marketplace.partialImportCompleted",
			Data:    result,
		})
	}

	return utils.SuccessResponse(c, result)
}

// TranslateAttribute automatically translates attribute display name and options
// @Summary Auto-translate attribute
// @Description Automatically translates attribute display name and options to all supported languages using Google Translate
// @Tags marketplace-admin-attributes
// @Accept json
// @Produce json
// @Param id path int true "Attribute ID"
// @Param languages body object{source_language=string,target_languages=[]string} false "Translation settings"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=TranslationResult} "Translation results"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid attribute ID"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Attribute not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Translation error"
// @Security BearerAuth
// @Router /api/v1/admin/marketplace/attributes/{id}/translate [post]
func (h *AdminAttributesHandler) TranslateAttribute(c *fiber.Ctx) error {
	logger.Info().Msg("TranslateAttribute method called")

	// Получаем ID атрибута из параметров URL
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttributeId")
	}

	// Получаем атрибут по ID
	attribute, err := h.marketplaceService.GetAttributeByID(c.Context(), attributeID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get attribute by ID")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getAttributeError")
	}

	if attribute == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.attributeNotFound")
	}

	// Парсим настройки перевода
	var input struct {
		SourceLanguage  string   `json:"source_language"`
		TargetLanguages []string `json:"target_languages"`
	}

	// Значения по умолчанию
	input.SourceLanguage = "en" // TODO: сделать enums для языков
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

	// Переводим display_name
	displayNameTranslations := make(map[string]string)
	for _, targetLang := range input.TargetLanguages {
		if targetLang == input.SourceLanguage {
			continue
		}

		translatedText, err := h.marketplaceService.TranslateText(c.Context(), attribute.DisplayName, input.SourceLanguage, targetLang)
		if err != nil {
			logger.Error().Err(err).Str("text", attribute.DisplayName).Str("target_lang", targetLang).Msg("Failed to translate display name")
			errors = append(errors, fmt.Sprintf("Failed to translate display_name to %s", targetLang))
			continue
		}
		displayNameTranslations[targetLang] = translatedText
	}

	// Сохраняем переводы display_name
	for lang, text := range displayNameTranslations {
		err := h.marketplaceService.SaveTranslation(c.Context(), "attribute", attributeID, lang, "display_name", text, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to save display_name translation")
			errors = append(errors, fmt.Sprintf("Failed to save display_name translation for %s", lang))
		}
	}

	translationResults["display_name"] = displayNameTranslations

	// Переводим опции для select/multiselect атрибутов
	if attribute.AttributeType == "select" || attribute.AttributeType == "multiselect" {
		var options []string
		if err := json.Unmarshal(attribute.Options, &options); err == nil && len(options) > 0 {
			optionTranslations := make(map[string]map[string]string)

			for _, targetLang := range input.TargetLanguages {
				if targetLang == input.SourceLanguage {
					continue
				}
				optionTranslations[targetLang] = make(map[string]string)

				for _, option := range options {
					translatedOption, err := h.marketplaceService.TranslateText(c.Context(), option, input.SourceLanguage, targetLang)
					if err != nil {
						logger.Error().Err(err).Str("option", option).Str("target_lang", targetLang).Msg("Failed to translate option")
						errors = append(errors, fmt.Sprintf("Failed to translate option '%s' to %s", option, targetLang))
						continue
					}
					optionTranslations[targetLang][option] = translatedOption
				}
			}

			// Сохраняем переводы опций
			for lang, options := range optionTranslations {
				for option, translatedOption := range options {
					err := h.marketplaceService.SaveTranslation(c.Context(), "attribute_option", attributeID, lang, option, translatedOption, nil)
					if err != nil {
						logger.Error().Err(err).Msg("Failed to save option translation")
						errors = append(errors, fmt.Sprintf("Failed to save option '%s' translation for %s", option, lang))
					}
				}
			}

			translationResults["options"] = optionTranslations
		}
	}

	// Формируем результат
	result := TranslationResult{
		AttributeID:  attributeID,
		Translations: translationResults,
		Errors:       errors,
	}

	if len(errors) > 0 {
		c.Status(fiber.StatusPartialContent)
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Success message"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid parameters"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
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
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Error().Err(err).Msg("Failed to rollback transaction")
		}
	}()

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
		logger.Error().Err(err).Msg("Failed to invalidate attribute cache")
	}

	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.attributeSettingsCopied",
	})
}
