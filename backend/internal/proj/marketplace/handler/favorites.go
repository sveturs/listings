// backend/internal/proj/marketplace/handler/favorites.go
package handler

import (
	"backend/internal/domain/models"
	//"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"strings"
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
func (h *FavoritesHandler) AddToFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Проверяем, существует ли объявление
	_, err = h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Listing %d not found: %v", listingID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Добавляем объявление в избранное
	err = h.marketplaceService.AddToFavorites(c.Context(), userID, listingID)
	if err != nil {
		log.Printf("Failed to add listing %d to favorites for user %d: %v", listingID, userID, err)
		// Проверяем, было ли объявление уже в избранном
		if strings.Contains(err.Error(), "already in favorites") {
			return utils.SuccessResponse(c, fiber.Map{
				"message": "Объявление уже в избранном",
			})
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось добавить объявление в избранное")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "Объявление добавлено в избранное",
	})
}

// RemoveFromFavorites удаляет объявление из избранного
func (h *FavoritesHandler) RemoveFromFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Удаляем объявление из избранного
	err = h.marketplaceService.RemoveFromFavorites(c.Context(), userID, listingID)
	if err != nil {
		log.Printf("Failed to remove listing %d from favorites for user %d: %v", listingID, userID, err)
		// Проверяем, было ли объявление в избранном
		if strings.Contains(err.Error(), "not in favorites") {
			return utils.SuccessResponse(c, fiber.Map{
				"message": "Объявление не было в избранном",
			})
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось удалить объявление из избранного")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "Объявление удалено из избранного",
	})
}

// GetFavorites получает список избранных объявлений пользователя
func (h *FavoritesHandler) GetFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем список избранных объявлений
	listings, err := h.marketplaceService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get favorites for user %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить список избранных объявлений")
	}

	// Проверяем на nil и возвращаем
	if listings == nil {
		listings = []models.MarketplaceListing{}
	}

	// Возвращаем список избранных объявлений
	return utils.SuccessResponse(c, listings)
}

// IsInFavorites проверяет, находится ли объявление в избранном пользователя
func (h *FavoritesHandler) IsInFavorites(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID объявления из параметров URL
	listingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Получаем список избранных объявлений
	favorites, err := h.marketplaceService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get favorites for user %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить избранное")
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
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"is_in_favorites": isInFavorites,
		},
	})
}

// GetFavoritesCount получает количество избранных объявлений пользователя
func (h *FavoritesHandler) GetFavoritesCount(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем список избранных объявлений и считаем их количество
	favorites, err := h.marketplaceService.GetUserFavorites(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get favorites count for user %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить количество избранных объявлений")
	}

	count := len(favorites)

	// Возвращаем количество
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"count": count,
		},
	})
}
