// TEMPORARY: Will be moved to microservice
package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/domain/models"
	"backend/pkg/utils"
)

// GetCategories godoc
// @Summary Получить список категорий
// @Description Получить список всех активных категорий
// @Tags marketplace
// @Accept json
// @Produce json
// @Param lang query string false "Язык (ru, en, sr)" default(ru)
// @Success 200 {object} utils.SuccessResponse{data=[]models.MarketplaceCategory}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/categories [get]
func (h *Handler) GetCategories(c *fiber.Ctx) error {
	lang := c.Query("lang", "ru")

	// Phase 7.4: Route to microservice if feature flag is enabled
	if h.useListingsMicroservice && h.listingsClient != nil {
		h.logger.Info().
			Bool("use_microservice", true).
			Str("lang", lang).
			Msg("Routing GetCategories to listings microservice")

		// Call microservice via gRPC
		grpcResp, err := h.listingsClient.GetAllCategories(c.Context())
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to get categories from microservice, falling back to monolith")
			// Fallback to monolith
			categories, err := h.storage.GetCategories(c.Context(), lang)
			if err != nil {
				h.logger.Error().Err(err).Msg("Failed to get categories from monolith")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.get_categories_failed")
			}
			return utils.SuccessResponse(c, categories)
		}

		// Convert proto.Category to models.MarketplaceCategory
		categories := make([]models.MarketplaceCategory, 0, len(grpcResp.Categories))
		for _, protoCategory := range grpcResp.Categories {
			category := models.MarketplaceCategory{
				ID:          int(protoCategory.Id),
				Name:        protoCategory.Name,
				Slug:        protoCategory.Slug,
				IsActive:    protoCategory.IsActive,
				Translations: protoCategory.Translations,
				ListingCount: int(protoCategory.ListingCount),
				HasCustomUI: protoCategory.HasCustomUi,
				SortOrder:   int(protoCategory.SortOrder),
				Level:       int(protoCategory.Level),
			}

			// Optional fields
			if protoCategory.ParentId != nil {
				parentID := int(*protoCategory.ParentId)
				category.ParentID = &parentID
			}
			if protoCategory.Icon != nil {
				category.Icon = protoCategory.Icon
			}
			if protoCategory.Description != nil {
				category.Description = protoCategory.Description
			}
			if protoCategory.CustomUiComponent != nil {
				category.CustomUIComponent = protoCategory.CustomUiComponent
			}
			// Parse created_at (proto uses string, model uses time.Time)
			if createdAt, err := time.Parse(time.RFC3339, protoCategory.CreatedAt); err == nil {
				category.CreatedAt = createdAt
			}

			categories = append(categories, category)
		}

		h.logger.Info().
			Int("count", len(categories)).
			Bool("served_by_microservice", true).
			Msg("Successfully retrieved categories from microservice")

		// Add header to indicate microservice was used
		c.Set("X-Served-By", "microservice")
		return utils.SuccessResponse(c, categories)
	}

	// Default: use monolith storage
	h.logger.Debug().
		Bool("use_microservice", false).
		Str("lang", lang).
		Msg("Routing GetCategories to monolith")

	categories, err := h.storage.GetCategories(c.Context(), lang)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get categories")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.get_categories_failed")
	}

	c.Set("X-Served-By", "monolith")
	return utils.SuccessResponse(c, categories)
}

// GetPopularCategories godoc
// @Summary Получить популярные категории
// @Description Получить список популярных категорий (с наибольшим количеством объявлений)
// @Tags marketplace
// @Accept json
// @Produce json
// @Param lang query string false "Язык (ru, en, sr)" default(ru)
// @Param limit query int false "Количество категорий" default(8)
// @Success 200 {object} utils.SuccessResponse{data=[]models.MarketplaceCategory}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/popular-categories [get]
func (h *Handler) GetPopularCategories(c *fiber.Ctx) error {
	lang := c.Query("lang", "ru")
	limit := c.QueryInt("limit", 8)

	categories, err := h.storage.GetPopularCategories(c.Context(), lang, limit)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get popular categories")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.get_popular_categories_failed")
	}

	return utils.SuccessResponse(c, categories)
}

// GetCategoryAttributes godoc
// @Summary Получить атрибуты категории
// @Description Получить список атрибутов для указанной категории
// @Tags marketplace
// @Accept json
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} utils.SuccessResponse{data=[]models.CategoryAttribute}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/categories/{id}/attributes [get]
func (h *Handler) GetCategoryAttributes(c *fiber.Ctx) error {
	categoryID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_category_id")
	}

	attributes, err := h.storage.GetCategoryAttributes(c.Context(), categoryID)
	if err != nil {
		h.logger.Error().Err(err).Int("category_id", categoryID).Msg("Failed to get category attributes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.get_attributes_failed")
	}

	return utils.SuccessResponse(c, attributes)
}

// GetVariantAttributes godoc
// @Summary Получить вариативные атрибуты категории
// @Description Получить список вариативных атрибутов для категории по slug
// @Tags marketplace
// @Accept json
// @Produce json
// @Param slug path string true "Slug категории"
// @Success 200 {object} utils.SuccessResponse{data=[]models.CategoryVariantAttribute}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/categories/{slug}/variant-attributes [get]
func (h *Handler) GetVariantAttributes(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_category_slug")
	}

	attributes, err := h.storage.GetVariantAttributes(c.Context(), slug)
	if err != nil {
		h.logger.Error().Err(err).Str("slug", slug).Msg("Failed to get variant attributes")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.get_variant_attributes_failed")
	}

	return utils.SuccessResponse(c, attributes)
}

// GetFavorites godoc
// @Summary Получить избранное пользователя
// @Description Получить список ID избранных объявлений текущего пользователя
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.SuccessResponse{data=[]int}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/favorites [get]
func (h *Handler) GetFavorites(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	favorites, err := h.storage.GetUserFavorites(c.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Int("user_id", userID).Msg("Failed to get favorites")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.get_favorites_failed")
	}

	return utils.SuccessResponse(c, favorites)
}

// AddToFavorites godoc
// @Summary Добавить в избранное
// @Description Добавить объявление в избранное текущего пользователя
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "ID объявления"
// @Success 200 {object} utils.SuccessResponse{data=map[string]bool}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/favorites/{id} [post]
func (h *Handler) AddToFavorites(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	listingID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_listing_id")
	}

	if err := h.storage.AddToFavorites(c.Context(), userID, listingID); err != nil {
		h.logger.Error().Err(err).Int("user_id", userID).Int("listing_id", listingID).Msg("Failed to add to favorites")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.add_favorite_failed")
	}

	return utils.SuccessResponse(c, map[string]bool{"added": true})
}

// RemoveFromFavorites godoc
// @Summary Удалить из избранного
// @Description Удалить объявление из избранного текущего пользователя
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "ID объявления"
// @Success 200 {object} utils.SuccessResponse{data=map[string]bool}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/marketplace/favorites/{id} [delete]
func (h *Handler) RemoveFromFavorites(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	listingID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_listing_id")
	}

	if err := h.storage.RemoveFromFavorites(c.Context(), userID, listingID); err != nil {
		h.logger.Error().Err(err).Int("user_id", userID).Int("listing_id", listingID).Msg("Failed to remove from favorites")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.remove_favorite_failed")
	}

	return utils.SuccessResponse(c, map[string]bool{"removed": true})
}
