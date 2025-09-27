// backend/internal/proj/marketplace/handler/favorites.go
package handler

import (
	"strconv"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// FavoritesHandler обрабатывает запросы, связанные с избранными объявлениями
type FavoritesHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewFavoritesHandler создает новый обработчик избранного
func NewFavoritesHandler(services globalService.ServicesInterface) *FavoritesHandler {
	return &FavoritesHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// AddToFavorites добавляет объявление в избранное
// @Summary Add listing to favorites
// @Description Adds a listing to user's favorites
// @Tags marketplace-favorites
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Added to favorites"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.addToFavoritesError"
// @Security BearerAuth
// @Router /api/v1/marketplace/favorites/{id} [post]
func (h *FavoritesHandler) AddToFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Проверяем тип - обычное объявление или товар витрины
	isStorefront := c.Query("type") == "storefront"

	if isStorefront {
		// Для товаров витрин используем отдельную таблицу storefront_favorites
		logger.Info().Int("productId", id).Int("userId", userID).Msg("Adding storefront product to favorites")
		err = h.marketplaceService.AddStorefrontToFavorites(c.Context(), userID, id)
	} else {
		// Проверяем, существует ли обычное объявление
		_, err = h.marketplaceService.GetListingByID(c.Context(), id)
		if err != nil {
			logger.Error().Err(err).Int("listingId", id).Msg("Listing not found")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
		}

		// Добавляем объявление в избранное
		err = h.marketplaceService.AddToFavorites(c.Context(), userID, id)
	}
	if err != nil {
		logger.Error().Err(err).Int("id", id).Int("userId", userID).Bool("isStorefront", isStorefront).Msg("Failed to add to favorites")
		// Проверяем, было ли объявление уже в избранном
		if strings.Contains(err.Error(), "already in favorites") {
			return utils.SuccessResponse(c, MessageResponse{
				Message: "marketplace.alreadyInFavorites",
			})
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.addToFavoritesError")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.addedToFavorites",
	})
}

// RemoveFromFavorites удаляет объявление из избранного
// @Summary Remove listing from favorites
// @Description Removes a listing from user's favorites
// @Tags marketplace-favorites
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Removed from favorites"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.removeFromFavoritesError"
// @Security BearerAuth
// @Router /api/v1/marketplace/favorites/{id} [delete]
func (h *FavoritesHandler) RemoveFromFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Проверяем тип - обычное объявление или товар витрины
	isStorefront := c.Query("type") == "storefront"

	// Удаляем из избранного в зависимости от типа
	if isStorefront {
		err = h.marketplaceService.RemoveStorefrontFromFavorites(c.Context(), userID, id)
	} else {
		err = h.marketplaceService.RemoveFromFavorites(c.Context(), userID, id)
	}
	if err != nil {
		logger.Error().Err(err).Int("id", id).Int("userId", userID).Bool("isStorefront", isStorefront).Msg("Failed to remove from favorites")
		// Проверяем, было ли объявление в избранном
		if strings.Contains(err.Error(), "not in favorites") {
			return utils.SuccessResponse(c, MessageResponse{
				Message: "marketplace.notInFavorites",
			})
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.removeFromFavoritesError")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.removedFromFavorites",
	})
}

// GetFavorites получает список избранных объявлений пользователя
// @Summary Get user's favorite listings
// @Description Returns all listings marked as favorites by the user
// @Tags marketplace-favorites
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceListing} "List of favorite listings"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getFavoritesError"
// @Security BearerAuth
// @Router /api/v1/marketplace/favorites [get]
func (h *FavoritesHandler) GetFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем список избранных объявлений
	listings, err := h.marketplaceService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get favorites for user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getFavoritesError")
	}

	// Проверяем на nil и возвращаем
	if listings == nil {
		listings = []models.MarketplaceListing{}
	}

	// Возвращаем список избранных объявлений
	return utils.SuccessResponse(c, listings)
}

// IsInFavorites проверяет, находится ли объявление в избранном пользователя
// @Summary Check if listing is in favorites
// @Description Checks if a specific listing is in user's favorites
// @Tags marketplace-favorites
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=FavoriteStatusData} "Favorite status"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.checkFavoritesError"
// @Security BearerAuth
// @Router /api/v1/marketplace/favorites/{id}/check [get]
func (h *FavoritesHandler) IsInFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем список избранных объявлений
	favorites, err := h.marketplaceService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get favorites for user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.checkFavoritesError")
	}

	// Проверяем, есть ли объявление в избранном
	isInFavorites := false
	for _, fav := range favorites {
		if fav.ID == listingID {
			isInFavorites = true
			break
		}
	}

	// Возвращаем результат проверки
	return utils.SuccessResponse(c, FavoriteStatusData{
		IsInFavorites: isInFavorites,
	})
}

// GetFavoritesCount получает количество избранных объявлений пользователя
// @Summary Get favorites count
// @Description Returns the total number of listings in user's favorites
// @Tags marketplace-favorites
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=FavoritesCountData} "Favorites count"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getFavoritesCountError"
// @Security BearerAuth
// @Router /api/v1/marketplace/favorites/count [get]
func (h *FavoritesHandler) GetFavoritesCount(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("userId", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем список избранных объявлений и считаем их количество
	favorites, err := h.marketplaceService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get favorites count for user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getFavoritesCountError")
	}

	count := len(favorites)

	// Возвращаем количество
	return utils.SuccessResponse(c, FavoritesCountData{
		Count: count,
	})
}
