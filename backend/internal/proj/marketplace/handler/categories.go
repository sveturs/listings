// backend/internal/proj/marketplace/handler/categories.go
package handler

import (
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"
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
func (h *CategoriesHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.marketplaceService.GetCategories(c.Context())
	if err != nil {
		log.Printf("Failed to get categories: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить категории")
	}

	return utils.SuccessResponse(c, categories)
}

// GetCategoryTree получает дерево категорий
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
		log.Printf("Failed to get category tree: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить дерево категорий")
	}

	// Обновляем кеш
	categoryTreeMutex.Lock()
	categoryTreeCache = categoryTree
	categoryTreeLastUpdate = time.Now()
	categoryTreeMutex.Unlock()

	return utils.SuccessResponse(c, categoryTree)
}

// GetCategoryAttributes получает атрибуты для категории
func (h *CategoriesHandler) GetCategoryAttributes(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Получаем атрибуты категории
	attributes, err := h.marketplaceService.GetCategoryAttributes(c.Context(), categoryID)
	if err != nil {
		log.Printf("Failed to get attributes for category %d: %v", categoryID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить атрибуты категории")
	}

	return utils.SuccessResponse(c, attributes)
}

// GetAttributeRanges получает диапазоны значений для числовых атрибутов
func (h *CategoriesHandler) GetAttributeRanges(c *fiber.Ctx) error {
	// Получаем ID категории из параметров URL
	categoryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID категории")
	}

	// Получаем диапазоны значений атрибутов
	ranges, err := h.marketplaceService.GetAttributeRanges(c.Context(), categoryID)
	if err != nil {
		log.Printf("Failed to get attribute ranges for category %d: %v", categoryID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить диапазоны значений атрибутов")
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
