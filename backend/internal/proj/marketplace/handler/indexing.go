// backend/internal/proj/marketplace/handler/indexing.go
package handler

import (
	"context"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
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
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Reindexing started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-listings [post]
func (h *IndexingHandler) ReindexAll(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("user_id", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userID", userID).Msg("Failed to get user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		logger.Error().Err(err).Int("userID", userID).Msg("User is not admin")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию в отдельной горутине, чтобы не блокировать запрос
	go func() {
		err := h.marketplaceService.ReindexAllListings(context.Background())
		if err != nil {
			logger.Error().Err(err).Msg("Reindex error")
		} else {
			logger.Info().Msg("Reindex completed successfully")
		}
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.reindexStarted",
	})
}

// ReindexAllWithTranslations reindexes all listings with translations
// @Summary Reindex all listings with translations
// @Description Reindexes all marketplace listings with their translations in the search index
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=ReindexStartedResponse} "Reindexing with translations started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-listings-with-translations [post]
func (h *IndexingHandler) ReindexAllWithTranslations(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("user_id", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userID", userID).Msg("Failed to get user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		logger.Error().Err(err).Int("userID", userID).Msg("User is not admin")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию с переводами в отдельной горутине
	go func() {
		startTime := time.Now()
		logger.Info().Time("startTime", startTime).Msg("Starting full reindex with translations")

		// Получаем все объявления
		filters := make(map[string]string)
		offset := 0
		limit := 100
		total := 0

		for {
			listings, count, err := h.services.Storage().GetListings(context.Background(), filters, limit, offset)
			if err != nil {
				logger.Error().Err(err).Msg("Error fetching listings")
				break
			}

			if len(listings) == 0 {
				break
			}

			total += len(listings)
			logger.Info().Int("count", len(listings)).Int("offset", offset).Msg("Processing listings")

			// Обрабатываем каждое объявление
			for _, listing := range listings {
				// Индексируем объявление
				err = h.services.Storage().IndexListing(context.Background(), &listing)
				if err != nil {
					logger.Error().Err(err).Int("listingID", listing.ID).Msg("Error indexing listing")
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
		logger.Info().
			Time("endTime", endTime).
			Dur("duration", duration).
			Int("processedListings", total).
			Msg("Reindex with translations completed")
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, ReindexStartedResponse{
		Success: true,
		Message: "marketplace.reindexWithTranslationsStarted",
	})
}

// ReindexAllListings reindexes all listings (alias method)
// @Summary Reindex all listings (alias)
// @Description Alternative endpoint to reindex all marketplace listings
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=ReindexStartedResponse} "Reindexing started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-all-listings [post]
func (h *IndexingHandler) ReindexAllListings(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("user_id", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userID", userID).Msg("Failed to get user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		logger.Error().Err(err).Int("userID", userID).Msg("User is not admin")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию
	go func() {
		err := h.marketplaceService.ReindexAllListings(context.Background())
		if err != nil {
			logger.Error().Err(err).Msg("Error during reindex")
		}
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, ReindexStartedResponse{
		Success: true,
		Message: "marketplace.reindexStarted",
	})
}

// ReindexRatings reindexes ratings for all listings
// @Summary Reindex listing ratings
// @Description Reindexes ratings and review counts for all marketplace listings
// @Tags marketplace-admin-indexing
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=ReindexStartedResponse} "Rating reindexing started"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.adminCheckError"
// @Security BearerAuth
// @Router /api/v1/admin/reindex-ratings [post]
func (h *IndexingHandler) ReindexRatings(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logger.Error().Interface("user_id", c.Locals("user_id")).Msg("Failed to get user_id from context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userID", userID).Msg("Failed to get user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		logger.Error().Err(err).Int("userID", userID).Msg("User is not admin")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Запускаем переиндексацию рейтингов в отдельной горутине
	go func() {
		logger.Info().Msg("Starting ratings reindex")
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
				logger.Error().Err(err).Msg("Error fetching listings")
				lastError = err
				break
			}

			if len(listings) == 0 {
				break
			}

			total += len(listings)
			logger.Info().Int("count", len(listings)).Int("offset", offset).Msg("Processing ratings for listings")

			// Обрабатываем каждое объявление
			for _, listing := range listings {
				// Получаем рейтинг объявления
				avgRating, err := h.services.Storage().GetEntityRating(context.Background(), "listing", listing.ID)
				if err != nil {
					logger.Error().Err(err).Int("listingID", listing.ID).Msg("Error fetching rating for listing")
					continue
				}

				// Получаем количество отзывов
				reviews, _, err := h.services.Review().GetReviews(context.Background(), models.ReviewsFilter{
					EntityType: "listing",
					EntityID:   listing.ID,
				})
				if err != nil {
					logger.Error().Err(err).Int("listingID", listing.ID).Msg("Error fetching reviews for listing")
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
					logger.Error().Err(err).Int("listingID", listing.ID).Msg("Error updating rating for listing")
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
		logger.Info().
			Dur("duration", duration).
			Int("processedListings", total).
			Int("reindexedListings", reindexed).
			Err(lastError).
			Msg("Ratings reindex completed")
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, ReindexStartedResponse{
		Success: true,
		Message: "marketplace.ratingsReindexStarted",
	})
}
