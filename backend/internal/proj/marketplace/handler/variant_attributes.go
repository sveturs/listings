package handler

import (
	"net/http"

	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	_ "backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// VariantAttributesHandler обрабатывает запросы связанные с вариативными атрибутами
type VariantAttributesHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewVariantAttributesHandler создает новый обработчик вариативных атрибутов
func NewVariantAttributesHandler(services globalService.ServicesInterface) *VariantAttributesHandler {
	return &VariantAttributesHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// GetProductVariantAttributes возвращает список доступных вариативных атрибутов
// @Summary Get product variant attributes
// @Description Returns list of all product variant attributes
// @Tags marketplace-variant-attributes
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.ProductVariantAttribute} "List of variant attributes"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/product-variant-attributes [get]
func (h *VariantAttributesHandler) GetProductVariantAttributes(c *fiber.Ctx) error {
	attributes, err := h.marketplaceService.GetProductVariantAttributes(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "marketplace.getVariantAttributesError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    attributes,
	})
}

// GetCategoryVariantAttributes возвращает вариативные атрибуты для конкретной категории
// @Summary Get variant attributes for category
// @Description Returns variant attributes suitable for specific category
// @Tags marketplace-variant-attributes
// @Accept json
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.ProductVariantAttribute} "Variant attributes for category"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Category not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/categories/{slug}/variant-attributes [get]
func (h *VariantAttributesHandler) GetCategoryVariantAttributes(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "marketplace.categorySlugRequired",
		})
	}

	attributes, err := h.marketplaceService.GetCategoryVariantAttributes(c.Context(), slug)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "marketplace.getCategoryVariantAttributesError",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    attributes,
	})
}
