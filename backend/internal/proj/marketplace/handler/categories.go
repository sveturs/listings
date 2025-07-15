// Package handler
// backend/internal/proj/marketplace/handler/categories.go
package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
)

// Используем переменные кеша из marketplace.go
// var (
// 	categoryTreeCache      []models.CategoryTreeNode
// 	categoryTreeLastUpdate time.Time
// 	categoryTreeMutex      sync.RWMutex
// )

// CategoriesHandler обрабатывает запросы, связанные с категориями
type CategoriesHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewCategoriesHandler создает новый обработчик категорий
func NewCategoriesHandler(services globalService.ServicesInterface) *CategoriesHandler {
	return &CategoriesHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// GetCategories получает список категорий
// @Summary Get categories list
// @Description Returns all marketplace categories
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.MarketplaceCategory}
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.categoriesError"
// @Router /api/v1/marketplace/categories [get]
func (h *CategoriesHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.marketplaceService.GetCategories(c.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get categories")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.categoriesError")
	}

	return utils.SuccessResponse(c, categories)
}

// GetCategoryTree получает дерево категорий
// @Summary Get category tree
// @Description Returns hierarchical tree of all categories with caching
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.CategoryTreeNode}
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.categoryTreeError"
// @Router /api/v1/marketplace/category-tree [get]
func (h *CategoriesHandler) GetCategoryTree(c *fiber.Ctx) error {
	// Оптимизация: используем кеш, если он актуален (не старше 5 минут)
	categoryTreeMutex.RLock()
	cacheValid := len(categoryTreeCache) > 0 && time.Since(categoryTreeLastUpdate) < 5*time.Minute
	cachedTree := categoryTreeCache
	categoryTreeMutex.RUnlock()

	if cacheValid {
		return utils.SuccessResponse(c, cachedTree)
	}

	// Если кеш устарел или пуст, загружаем дерево категорий из хранилища
	categoryTree, err := h.marketplaceService.GetCategoryTree(c.Context())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get category tree")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.categoryTreeError")
	}

	// Обновляем кеш
	categoryTreeMutex.Lock()
	categoryTreeCache = categoryTree
	categoryTreeLastUpdate = time.Now()
	categoryTreeMutex.Unlock()

	return utils.SuccessResponse(c, categoryTree)
}

// GetCategoryAttributes получает атрибуты для категории
// @Summary Get category attributes
// @Description Returns all attributes available for a specific category
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param lang query string false "Language code (e.g., 'sr', 'en', 'ru')"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_models.CategoryAttribute}
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.attributesError"
// @Router /api/v1/marketplace/categories/{id}/attributes [get]
func (h *CategoriesHandler) GetCategoryAttributes(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем язык из query параметров
	lang := c.Query("lang", "en") // По умолчанию английский

	// Получаем атрибуты категории с указанным языком
	attributes, err := h.marketplaceService.GetCategoryAttributesWithLang(c.Context(), categoryID, lang)
	if err != nil {
		logger.Error().Err(err).Int("categoryId", categoryID).Str("lang", lang).Msg("Failed to get attributes for category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.attributesError")
	}

	return utils.SuccessResponse(c, attributes)
}

// GetAttributeRanges получает диапазоны значений для числовых атрибутов
// @Summary Get attribute value ranges
// @Description Returns min/max values for numeric attributes in a category
// @Tags marketplace-categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=AttributeRangesResponse}
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidCategoryId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.rangesError"
// @Router /api/v1/marketplace/categories/{id}/attribute-ranges [get]
func (h *CategoriesHandler) GetAttributeRanges(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidCategoryId")
	}

	// Получаем диапазоны значений атрибутов
	ranges, err := h.marketplaceService.GetAttributeRanges(c.Context(), categoryID)
	if err != nil {
		logger.Error().Err(err).Int("categoryId", categoryID).Msg("Failed to get attribute ranges for category")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.rangesError")
	}

	return utils.SuccessResponse(c, ranges)
}

// InvalidateCategoryCache инвалидирует кеш категорий
func (h *CategoriesHandler) InvalidateCategoryCache() {
	categoryTreeMutex.Lock()
	categoryTreeCache = nil
	categoryTreeLastUpdate = time.Time{}
	categoryTreeMutex.Unlock()
}
