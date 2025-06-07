// backend/internal/proj/marketplace/handler/indexing.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

// IndexingHandler handles requests related to listing indexing
type IndexingHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewIndexingHandler creates a new indexing handler
func NewIndexingHandler(services globalService.ServicesInterface) *IndexingHandler {
	return &IndexingHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// ReindexAll reindexes all listings
// @Summary Reindex all listings
// @Description Reindexes all marketplace listings in the search index
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string} "Reindexing started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-listings [post]
func (h *IndexingHandler) ReindexAll(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию в отдельной горутине, чтобы не блокировать запрос
	go func() {
		err := h.marketplaceService.ReindexAllListings(context.Background())
		if err != nil {
			log.Printf("Reindex error: %v", err)
		} else {
			log.Println("Reindex completed successfully")
		}
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.reindexStarted",
	})
}

// ReindexAllWithTranslations reindexes all listings with translations
// @Summary Reindex all listings with translations
// @Description Reindexes all marketplace listings with their translations in the search index
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string} "Reindexing with translations started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-listings-with-translations [post]
func (h *IndexingHandler) ReindexAllWithTranslations(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию с переводами в отдельной горутине
	go func() {
		startTime := time.Now()
		log.Printf("Starting full reindex with translations at %s", startTime.Format(time.RFC3339))

		// Получаем все объявления
		filters := make(map[string]string)
		offset := 0
		limit := 100
		total := 0

		for {
			listings, count, err := h.services.Storage().GetListings(context.Background(), filters, limit, offset)
			if err != nil {
				log.Printf("Error fetching listings: %v", err)
				break
			}

			if len(listings) == 0 {
				break
			}

			total += len(listings)
			log.Printf("Processing %d listings (offset %d)", len(listings), offset)

			// Обрабатываем каждое объявление
			for _, listing := range listings {
				// Индексируем объявление
				err = h.services.Storage().IndexListing(context.Background(), &listing)
				if err != nil {
					log.Printf("Error indexing listing %d: %v", listing.ID, err)
				}
			}

			offset += limit
			if offset >= int(count) {
				break
			}

			// Небольшая пауза, чтобы не перегружать сервер
			time.Sleep(100 * time.Millisecond)
		}

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		log.Printf("Reindex with translations completed at %s, duration: %s, processed %d listings",
			endTime.Format(time.RFC3339), duration, total)
	}()

	// Возвращаем успешный результат
	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.reindexWithTranslationsStarted",
	})
}

// ReindexAllListings reindexes all listings (alias method)
// @Summary Reindex all listings (alias)
// @Description Alternative endpoint to reindex all marketplace listings
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string} "Reindexing started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-all-listings [post]
func (h *IndexingHandler) ReindexAllListings(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию
	go func() {
		err := h.marketplaceService.ReindexAllListings(context.Background())
		if err != nil {
			log.Printf("Error during reindex: %v", err)
		}
	}()

	// Возвращаем успешный результат
	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.reindexStarted",
	})
}

// ReindexRatings reindexes ratings for all listings
// @Summary Reindex listing ratings
// @Description Reindexes ratings and review counts for all marketplace listings
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} object{success=bool,message=string} "Rating reindexing started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-ratings [post]
func (h *IndexingHandler) ReindexRatings(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию рейтингов в отдельной горутине
	go func() {
		log.Println("Starting ratings reindex")
		startTime := time.Now()

		// Получаем все объявления
		filters := make(map[string]string)
		offset := 0
		limit := 100
		total := 0
		reindexed := 0
		var lastError error

		for {
			listings, count, err := h.services.Storage().GetListings(context.Background(), filters, limit, offset)
			if err != nil {
				log.Printf("Error fetching listings: %v", err)
				lastError = err
				break
			}

			if len(listings) == 0 {
				break
			}

			total += len(listings)
			log.Printf("Processing ratings for %d listings (offset %d)", len(listings), offset)

			// Обрабатываем каждое объявление
			for _, listing := range listings {
				// Получаем рейтинг объявления
				avgRating, err := h.services.Storage().GetEntityRating(context.Background(), "listing", listing.ID)
				if err != nil {
					log.Printf("Error fetching rating for listing %d: %v", listing.ID, err)
					continue
				}

				// Получаем количество отзывов
				reviews, _, err := h.services.Review().GetReviews(context.Background(), models.ReviewsFilter{
					EntityType: "listing",
					EntityID:   listing.ID,
				})
				if err != nil {
					log.Printf("Error fetching reviews for listing %d: %v", listing.ID, err)
					continue
				}

				reviewCount := len(reviews)

				// Для обновления рейтинга в индексе можно использовать IndexListing
				// Установим рейтинг в объект listing
				listing.AverageRating = avgRating
				listing.ReviewCount = reviewCount

				// Обновляем индекс
				err = h.services.Storage().IndexListing(context.Background(), &listing)
				if err != nil {
					log.Printf("Error updating rating for listing %d: %v", listing.ID, err)
					lastError = err
					continue
				}

				reindexed++
			}

			offset += limit
			if offset >= int(count) {
				break
			}

			// Небольшая пауза, чтобы не перегружать сервер
			time.Sleep(100 * time.Millisecond)
		}

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		log.Printf("Ratings reindex completed in %s, processed %d listings, reindexed %d, last error: %v",
			duration, total, reindexed, lastError)
	}()

	// Возвращаем успешный результат
	return c.JSON(fiber.Map{
		"success": true,
		"message": "marketplace.ratingsReindexStarted",
	})
}
