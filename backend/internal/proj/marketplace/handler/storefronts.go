// TEMPORARY: Will be moved to microservice
package handler

import (
	"strconv"

	"backend/internal/proj/marketplace/storage"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// GetStorefronts возвращает список витрин (B2C stores)
// @Summary Get storefronts list
// @Description Get list of active B2C storefronts with filtering and pagination
// @Tags marketplace
// @Accept json
// @Produce json
// @Param is_active query boolean false "Filter by active status"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Param sort_by query string false "Sort by field: products_count, rating, created_at, views_count (default: products_count)"
// @Param sort_order query string false "Sort order: asc, desc (default: desc)"
// @Success 200 {object} utils.SuccessResponse{data=object{storefronts=[]models.Storefront,total=int,page=int,limit=int}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/b2c [get]
func (h *Handler) GetStorefronts(c *fiber.Ctx) error {
	// Парсим query параметры
	filters := storage.StorefrontFilters{
		Page:      c.QueryInt("page", 1),
		Limit:     c.QueryInt("limit", 10),
		SortBy:    c.Query("sort_by", "products_count"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	// Валидация лимита
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}

	// Парсим is_active фильтр
	if c.Query("is_active") != "" {
		isActiveStr := c.Query("is_active")
		isActive := isActiveStr == "true" || isActiveStr == "1"
		filters.IsActive = &isActive
	}

	// Получаем витрины из БД
	storefronts, total, err := h.storage.GetStorefronts(c.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get storefronts")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.fetch_failed")
	}

	return c.JSON(fiber.Map{
		"storefronts": storefronts,
		"total":       total,
		"page":        filters.Page,
		"limit":       filters.Limit,
	})
}

// GetStorefrontBySlug возвращает витрину по slug
// @Summary Get storefront by slug
// @Description Get single storefront by slug
// @Tags marketplace
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Success 200 {object} utils.SuccessResponse{data=models.Storefront}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/b2c/{slug} [get]
func (h *Handler) GetStorefrontBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.slug_required")
	}

	// Получаем одну витрину (limit=1, фильтр по slug можно добавить позже)
	// Пока используем простой запрос через GetStorefronts
	// TODO: Добавить отдельный метод GetStorefrontBySlug в storage
	filters := storage.StorefrontFilters{
		Page:  1,
		Limit: 1,
	}

	storefronts, _, err := h.storage.GetStorefronts(c.Context(), filters)
	if err != nil {
		h.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get storefront by slug")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.fetch_failed")
	}

	// Фильтруем по slug на уровне приложения (временное решение)
	for _, store := range storefronts {
		if store.Slug == slug {
			return c.JSON(store)
		}
	}

	return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.not_found")
}

// GetStorefrontProducts возвращает товары витрины
// @Summary Get storefront products
// @Description Get products list for specific storefront
// @Tags marketplace
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Success 200 {object} utils.SuccessResponse{data=object{products=[]models.Product,total=int}}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/b2c/{id}/products [get]
func (h *Handler) GetStorefrontProducts(c *fiber.Ctx) error {
	idStr := c.Params("id")
	storefrontID, err := strconv.Atoi(idStr)
	if err != nil || storefrontID <= 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	// TODO: Реализовать получение продуктов витрины
	// Пока возвращаем заглушку
	h.logger.Warn().Int("storefront_id", storefrontID).Msg("GetStorefrontProducts not yet implemented")

	return c.JSON(fiber.Map{
		"products": []interface{}{},
		"total":    0,
		"page":     c.QueryInt("page", 1),
		"limit":    c.QueryInt("limit", 20),
	})
}
