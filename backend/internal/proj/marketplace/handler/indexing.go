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

// IndexingHandler обрабатывает запросы, связанные с индексацией объявлений
type IndexingHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewIndexingHandler создает новый обработчик индексации
func NewIndexingHandler(services globalService.ServicesInterface) *IndexingHandler {
	return &IndexingHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// ReindexAll переиндексирует все объявления
func (h *IndexingHandler) ReindexAll(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить права администратора")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Требуются права администратора")
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
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Переиндексация запущена",
	})
}

// ReindexAllWithTranslations переиндексирует все объявления с переводами
func (h *IndexingHandler) ReindexAllWithTranslations(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить права администратора")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Требуются права администратора")
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
		"message": "Переиндексация с переводами запущена",
	})
}

// ReindexAllListings переиндексирует все объявления
func (h *IndexingHandler) ReindexAllListings(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить права администратора")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Требуются права администратора")
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
		"message": "Переиндексация запущена",
	})
}

// ReindexRatings переиндексирует рейтинги всех объявлений
func (h *IndexingHandler) ReindexRatings(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить права администратора")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Требуются права администратора")
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
		"message": "Переиндексация рейтингов запущена",
	})
}
