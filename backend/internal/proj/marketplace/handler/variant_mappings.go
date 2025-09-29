package handler

import (
	"strconv"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/pkg/utils"
	globalService "backend/internal/proj/global/service"
	"backend/internal/services/attributes"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

// VariantMappingsHandler обрабатывает запросы для управления связями вариативных атрибутов
type VariantMappingsHandler struct {
	services         globalService.ServicesInterface
	attributeService *attributes.UnifiedAttributeService
}

// NewVariantMappingsHandler создает новый обработчик для вариативных связей
func NewVariantMappingsHandler(
	services globalService.ServicesInterface,
	storage postgres.UnifiedAttributeStorage,
	featureFlags *config.FeatureFlags,
) *VariantMappingsHandler {
	// Создаем сервис с поддержкой feature flags
	attributeService := attributes.NewUnifiedAttributeService(
		storage,
		featureFlags.UnifiedAttributesFallback,
		featureFlags.DualWriteAttributes,
	)

	return &VariantMappingsHandler{
		services:         services,
		attributeService: attributeService,
	}
}

// GetVariantCompatibleAttributes возвращает все атрибуты с is_variant_compatible = true
// @Summary Get variant compatible attributes
// @Description Returns all attributes that can be used as product variants
// @Tags admin-variant-attributes
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_domain_models.UnifiedAttribute} "List of variant compatible attributes"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/attributes/variant-compatible [get]
func (h *VariantMappingsHandler) GetVariantCompatibleAttributes(c *fiber.Ctx) error {
	attributes, err := h.attributeService.GetVariantAttributes(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "attributes.getVariantAttributesError")
	}

	return utils.SuccessResponse(c, attributes)
}

// GetCategoryVariantMappings возвращает вариативные атрибуты для категории
// @Summary Get category variant mappings
// @Description Returns variant attribute mappings for a specific category
// @Tags admin-variant-attributes
// @Accept json
// @Produce json
// @Param category_id query int true "Category ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_domain_models.VariantAttributeMapping} "Category variant mappings"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid category ID"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/variant-attributes/mappings [get]
func (h *VariantMappingsHandler) GetCategoryVariantMappings(c *fiber.Ctx) error {
	categoryIDStr := c.Query("category_id")
	if categoryIDStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.categoryIdRequired")
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidCategoryId")
	}

	mappings, err := h.attributeService.GetCategoryVariantAttributes(c.Context(), categoryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "attributes.getMappingsError")
	}

	return utils.SuccessResponse(c, mappings)
}

// CreateVariantMapping создает новую связь между вариативным атрибутом и категорией
// @Summary Create variant mapping
// @Description Creates a new mapping between variant attribute and category
// @Tags admin-variant-attributes
// @Accept json
// @Produce json
// @Param body body backend_internal_domain_models.VariantAttributeMappingCreateRequest true "Mapping data"
// @Success 201 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_domain_models.VariantAttributeMapping} "Created mapping"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid data"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/variant-attributes/mappings [post]
func (h *VariantMappingsHandler) CreateVariantMapping(c *fiber.Ctx) error {
	var request models.VariantAttributeMappingCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidData")
	}

	// Валидация
	if request.VariantAttributeID <= 0 || request.CategoryID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidMappingData")
	}

	mapping, err := h.attributeService.CreateVariantAttributeMapping(c.Context(), &request)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "attributes.createMappingError")
	}

	return utils.SuccessResponse(c, mapping)
}

// UpdateVariantMapping обновляет связь между вариативным атрибутом и категорией
// @Summary Update variant mapping
// @Description Updates an existing variant attribute mapping
// @Tags admin-variant-attributes
// @Accept json
// @Produce json
// @Param id path int true "Mapping ID"
// @Param body body backend_internal_domain_models.VariantAttributeMappingUpdateRequest true "Update data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Mapping updated"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid data"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Mapping not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/variant-attributes/mappings/{id} [patch]
func (h *VariantMappingsHandler) UpdateVariantMapping(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidMappingId")
	}

	var request models.VariantAttributeMappingUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidData")
	}

	err = h.attributeService.UpdateVariantAttributeMapping(c.Context(), id, &request)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "attributes.updateMappingError")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "attributes.mappingUpdated"})
}

// DeleteVariantMapping удаляет связь между вариативным атрибутом и категорией
// @Summary Delete variant mapping
// @Description Deletes a variant attribute mapping
// @Tags admin-variant-attributes
// @Accept json
// @Produce json
// @Param id path int true "Mapping ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Mapping deleted"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid ID"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Mapping not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/variant-attributes/mappings/{id} [delete]
func (h *VariantMappingsHandler) DeleteVariantMapping(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidMappingId")
	}

	err = h.attributeService.DeleteVariantAttributeMapping(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "attributes.deleteMappingError")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "attributes.mappingDeleted"})
}

// UpdateCategoryVariantAttributes обновляет все вариативные атрибуты для категории
// @Summary Update category variant attributes
// @Description Updates all variant attributes for a category (replaces existing)
// @Tags admin-variant-attributes
// @Accept json
// @Produce json
// @Param body body backend_internal_domain_models.CategoryVariantAttributesUpdateRequest true "Update request"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Attributes updated"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid data"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/categories/variant-attributes [put]
func (h *VariantMappingsHandler) UpdateCategoryVariantAttributes(c *fiber.Ctx) error {
	var request models.CategoryVariantAttributesUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidData")
	}

	// Валидация
	if request.CategoryID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "attributes.invalidCategoryId")
	}

	err := h.attributeService.UpdateCategoryVariantAttributes(c.Context(), &request)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "attributes.updateCategoryAttributesError")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "attributes.categoryAttributesUpdated"})
}
