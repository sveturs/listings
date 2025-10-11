// Package handler
// backend/internal/proj/c2c/handler/admin_variant_attributes.go
package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
)

// AdminVariantAttributesHandler обрабатывает запросы админки для управления вариативными атрибутами
type AdminVariantAttributesHandler struct {
	*CategoriesHandler
}

// NewAdminVariantAttributesHandler создает новый обработчик админки для вариативных атрибутов
func NewAdminVariantAttributesHandler(services globalService.ServicesInterface) *AdminVariantAttributesHandler {
	return &AdminVariantAttributesHandler{
		CategoriesHandler: NewCategoriesHandler(services),
	}
}

// GetVariantAttributes получает список всех вариативных атрибутов
// @Summary Get variant attributes
// @Description Gets list of all variant attributes from product_variant_attributes table
// @Tags marketplace-admin-variant-attributes
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param search query string false "Search by name or display_name"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.ProductVariantAttribute} "Variant attributes list"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getVariantAttributesError"
// @Security BearerAuth
// @Router /api/admin/variant-attributes [get]
func (h *AdminVariantAttributesHandler) GetVariantAttributes(c *fiber.Ctx) error {
	ctx := context.Background()

	// Получаем параметры пагинации
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	search := c.Query("search", "")

	// Базовый запрос
	query := `
		SELECT 
			id, name, display_name, type, is_required, sort_order, 
			affects_stock, created_at, updated_at
		FROM product_variant_attributes
	`
	args := []interface{}{}
	argIndex := 0

	// Добавляем поиск если указан
	if search != "" {
		query += " WHERE (name ILIKE $" + strconv.Itoa(argIndex+1) + " OR display_name ILIKE $" + strconv.Itoa(argIndex+2) + ")"
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	// Добавляем сортировку и лимиты
	query += " ORDER BY sort_order, name LIMIT $" + strconv.Itoa(argIndex+1) + " OFFSET $" + strconv.Itoa(argIndex+2)
	args = append(args, limit, offset)

	rows, err := h.services.Storage().Query(ctx, query, args...)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get variant attributes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getVariantAttributesError")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var attributes []models.ProductVariantAttribute
	for rows.Next() {
		var attr models.ProductVariantAttribute
		err := rows.Scan(
			&attr.ID,
			&attr.Name,
			&attr.DisplayName,
			&attr.Type,
			&attr.IsRequired,
			&attr.SortOrder,
			&attr.AffectsStock,
			&attr.CreatedAt,
			&attr.UpdatedAt,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan variant attribute")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getVariantAttributesError")
		}
		attributes = append(attributes, attr)
	}

	return utils.SuccessResponse(c, attributes)
}

// CreateVariantAttribute создает новый вариативный атрибут
// @Summary Create variant attribute
// @Description Creates a new variant attribute in product_variant_attributes table
// @Tags marketplace-admin-variant-attributes
// @Accept json
// @Produce json
// @Param body body models.ProductVariantAttribute true "Variant attribute data"
// @Success 201 {object} utils.SuccessResponseSwag{data=models.ProductVariantAttribute} "marketplace.variantAttributeCreated"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData or marketplace.requiredFieldsMissing"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createVariantAttributeError"
// @Security BearerAuth
// @Router /api/admin/variant-attributes [post]
func (h *AdminVariantAttributesHandler) CreateVariantAttribute(c *fiber.Ctx) error {
	ctx := context.Background()

	var attr models.ProductVariantAttribute
	if err := c.BodyParser(&attr); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if attr.Name == "" || attr.DisplayName == "" || attr.Type == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.requiredFieldsMissing")
	}

	// Создаем атрибут в БД
	query := `
		INSERT INTO product_variant_attributes (
			name, display_name, type, is_required, sort_order, affects_stock
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := h.services.Storage().QueryRow(
		ctx, query,
		attr.Name,
		attr.DisplayName,
		attr.Type,
		attr.IsRequired,
		attr.SortOrder,
		attr.AffectsStock,
	).Scan(&attr.ID, &attr.CreatedAt, &attr.UpdatedAt)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create variant attribute")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createVariantAttributeError")
	}

	return utils.SuccessResponse(c, attr)
}

// UpdateVariantAttribute обновляет существующий вариативный атрибут
// @Summary Update variant attribute
// @Description Updates an existing variant attribute
// @Tags marketplace-admin-variant-attributes
// @Accept json
// @Produce json
// @Param id path int true "Variant attribute ID"
// @Param body body models.ProductVariantAttribute true "Updated variant attribute data"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ProductVariantAttribute} "marketplace.variantAttributeUpdated"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.variantAttributeNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateVariantAttributeError"
// @Security BearerAuth
// @Router /api/admin/variant-attributes/{id} [put]
func (h *AdminVariantAttributesHandler) UpdateVariantAttribute(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidID")
	}

	var attr models.ProductVariantAttribute
	if err := c.BodyParser(&attr); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем обязательные поля
	if attr.Name == "" || attr.DisplayName == "" || attr.Type == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.requiredFieldsMissing")
	}

	// Обновляем атрибут в БД
	query := `
		UPDATE product_variant_attributes
		SET 
			name = $1, 
			display_name = $2, 
			type = $3, 
			is_required = $4, 
			sort_order = $5,
			affects_stock = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING id, created_at, updated_at
	`

	err = h.services.Storage().QueryRow(
		ctx, query,
		attr.Name,
		attr.DisplayName,
		attr.Type,
		attr.IsRequired,
		attr.SortOrder,
		attr.AffectsStock,
		id,
	).Scan(&attr.ID, &attr.CreatedAt, &attr.UpdatedAt)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update variant attribute")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateVariantAttributeError")
	}

	return utils.SuccessResponse(c, attr)
}

// DeleteVariantAttribute удаляет вариативный атрибут
// @Summary Delete variant attribute
// @Description Deletes a variant attribute (only if not used by any products)
// @Tags marketplace-admin-variant-attributes
// @Produce json
// @Param id path int true "Variant attribute ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "marketplace.variantAttributeDeleted"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidID"
// @Failure 409 {object} utils.ErrorResponseSwag "marketplace.variantAttributeInUse"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteVariantAttributeError"
// @Security BearerAuth
// @Router /api/admin/variant-attributes/{id} [delete]
func (h *AdminVariantAttributesHandler) DeleteVariantAttribute(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidID")
	}

	// Проверяем, используется ли атрибут в товарах
	var count int
	err = h.services.Storage().QueryRow(ctx, `
		SELECT COUNT(*) FROM b2c_product_variants spv
		JOIN b2c_products sp ON spv.product_id = sp.id
		WHERE spv.variant_attributes ? (SELECT name FROM product_variant_attributes WHERE id = $1)
	`, id).Scan(&count)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to check variant attribute usage")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteVariantAttributeError")
	}

	if count > 0 {
		return utils.ErrorResponse(c, fiber.StatusConflict, "marketplace.variantAttributeInUse")
	}

	// Удаляем атрибут
	_, err = h.services.Storage().Exec(ctx,
		"DELETE FROM product_variant_attributes WHERE id = $1", id)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete variant attribute")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteVariantAttributeError")
	}

	return utils.SuccessResponse(c, "marketplace.variantAttributeDeleted")
}

// GetVariantAttributeByID получает вариативный атрибут по ID
// @Summary Get variant attribute by ID
// @Description Gets a single variant attribute by its ID
// @Tags marketplace-admin-variant-attributes
// @Produce json
// @Param id path int true "Variant attribute ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ProductVariantAttribute} "Variant attribute details"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidID"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.variantAttributeNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getVariantAttributeError"
// @Security BearerAuth
// @Router /api/admin/variant-attributes/{id} [get]
func (h *AdminVariantAttributesHandler) GetVariantAttributeByID(c *fiber.Ctx) error {
	ctx := context.Background()

	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidID")
	}

	var attr models.ProductVariantAttribute
	query := `
		SELECT 
			id, name, display_name, type, is_required, sort_order, 
			affects_stock, created_at, updated_at
		FROM product_variant_attributes
		WHERE id = $1
	`

	err = h.services.Storage().QueryRow(ctx, query, id).Scan(
		&attr.ID,
		&attr.Name,
		&attr.DisplayName,
		&attr.Type,
		&attr.IsRequired,
		&attr.SortOrder,
		&attr.AffectsStock,
		&attr.CreatedAt,
		&attr.UpdatedAt,
	)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get variant attribute by ID")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.variantAttributeNotFound")
	}

	return utils.SuccessResponse(c, attr)
}

// GetVariantAttributeMappings получает связи вариативного атрибута с атрибутами категорий
// @Summary Get variant attribute mappings
// @Description Gets all category attribute mappings for a variant attribute
// @Tags marketplace-admin-variant-attributes
// @Produce json
// @Param id path int true "Variant attribute ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.CategoryAttribute} "Category attributes linked to this variant attribute"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidID"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getMappingsError"
// @Security BearerAuth
// @Router /api/v1/admin/variant-attributes/{id}/mappings [get]
func (h *AdminVariantAttributesHandler) GetVariantAttributeMappings(c *fiber.Ctx) error {
	ctx := context.Background()

	variantAttrID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidID")
	}

	// Получаем все атрибуты категорий, связанные с этим вариативным атрибутом
	query := `
		SELECT 
			ca.id,
			ca.name,
			ca.display_name,
			ca.attribute_type,
			ca.options,
			ca.is_searchable,
			ca.is_filterable,
			ca.is_required,
			ca.is_variant_compatible,
			ca.affects_stock,
			ca.sort_order
		FROM category_attributes ca
		INNER JOIN variant_attribute_mappings vam ON ca.id = vam.category_attribute_id
		WHERE vam.variant_attribute_id = $1
		ORDER BY ca.name
	`

	rows, err := h.services.Storage().Query(ctx, query, variantAttrID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get variant attribute mappings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getMappingsError")
	}
	defer func() { _ = rows.Close() }()

	var mappings []models.CategoryAttribute
	for rows.Next() {
		var attr models.CategoryAttribute
		err := rows.Scan(
			&attr.ID,
			&attr.Name,
			&attr.DisplayName,
			&attr.AttributeType,
			&attr.Options,
			&attr.IsSearchable,
			&attr.IsFilterable,
			&attr.IsRequired,
			&attr.IsVariantCompatible,
			&attr.AffectsStock,
			&attr.SortOrder,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan category attribute")
			continue
		}
		mappings = append(mappings, attr)
	}

	return utils.SuccessResponse(c, mappings)
}

// UpdateVariantAttributeMappings обновляет связи вариативного атрибута с атрибутами категорий
// @Summary Update variant attribute mappings
// @Description Updates category attribute mappings for a variant attribute
// @Tags marketplace-admin-variant-attributes
// @Accept json
// @Produce json
// @Param id path int true "Variant attribute ID"
// @Param body body object{category_attribute_ids=[]int} true "Array of category attribute IDs to link"
// @Success 200 {object} utils.SuccessResponseSwag{data=string} "marketplace.mappingsUpdated"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateMappingsError"
// @Security BearerAuth
// @Router /api/v1/admin/variant-attributes/{id}/mappings [put]
func (h *AdminVariantAttributesHandler) UpdateVariantAttributeMappings(c *fiber.Ctx) error {
	ctx := context.Background()

	variantAttrID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidID")
	}

	// Парсим тело запроса
	var body struct {
		CategoryAttributeIDs []int `json:"category_attribute_ids"`
	}
	if err := c.BodyParser(&body); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Начинаем транзакцию
	tx, err := h.services.Storage().BeginTx(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to begin transaction")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateMappingsError")
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Игнорируем ошибку, если транзакция уже завершена
			logger.Debug().Err(err).Msg("Transaction rollback")
		}
	}()

	// Удаляем старые связи
	_, err = tx.Exec(ctx, "DELETE FROM variant_attribute_mappings WHERE variant_attribute_id = $1", variantAttrID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete old mappings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateMappingsError")
	}

	// Создаем новые связи
	for _, catAttrID := range body.CategoryAttributeIDs {
		// Проверяем, что атрибут категории существует и имеет is_variant_compatible = true
		var isCompatible bool
		err = tx.QueryRow(ctx,
			"SELECT is_variant_compatible FROM category_attributes WHERE id = $1",
			catAttrID,
		).Scan(&isCompatible)
		if err != nil {
			logger.Error().Err(err).Int("category_attribute_id", catAttrID).Msg("Category attribute not found")
			continue
		}

		if !isCompatible {
			logger.Warn().Int("category_attribute_id", catAttrID).Msg("Category attribute is not variant compatible")
			continue
		}

		// Создаем связь
		_, err = tx.Exec(ctx, `
			INSERT INTO variant_attribute_mappings (variant_attribute_id, category_attribute_id)
			VALUES ($1, $2)
			ON CONFLICT (variant_attribute_id, category_attribute_id) DO NOTHING
		`, variantAttrID, catAttrID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to insert mapping")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateMappingsError")
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		logger.Error().Err(err).Msg("Failed to commit transaction")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateMappingsError")
	}

	return utils.SuccessResponse(c, "marketplace.mappingsUpdated")
}
