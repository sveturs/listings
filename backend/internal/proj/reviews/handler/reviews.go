// Package handler
// backend/internal/proj/reviews/handler/reviews.go
package handler

import (
	"log"
	"strconv"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/reviews/service"
	"backend/pkg/utils"
)

// ReviewHandler - обработчик отзывов
type ReviewHandler struct {
	services      globalService.ServicesInterface
	reviewService service.ReviewServiceInterface
}

// NewReviewHandler создаёт новый обработчик
func NewReviewHandler(services globalService.ServicesInterface) *ReviewHandler {
	if services == nil {
		log.Fatal("services cannot be nil")
	}
	if services.Review() == nil {
		log.Fatal("review service cannot be nil")
	}

	return &ReviewHandler{
		services:      services,
		reviewService: services.Review(),
	}
}

// GetReviews returns filtered list of reviews
// @Summary Get reviews list
// @Description Returns paginated list of reviews with filters
// @Tags reviews
// @Accept json
// @Produce json
// @Param entity_type query string false "Entity type filter"
// @Param entity_id query int false "Entity ID filter"
// @Param user_id query int false "User ID filter"
// @Param min_rating query int false "Minimum rating filter"
// @Param max_rating query int false "Maximum rating filter"
// @Param status query string false "Status filter" default(published)
// @Param sort_by query string false "Sort by field" default(date)
// @Param sort_order query string false "Sort order" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewsListResponse} "List of reviews"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews [get]
func (h *ReviewHandler) GetReviews(c *fiber.Ctx) error {
	filter := models.ReviewsFilter{
		EntityType: c.Query("entity_type"),
		EntityID:   utils.StringToInt(c.Query("entity_id"), 0),
		UserID:     utils.StringToInt(c.Query("user_id"), 0),
		MinRating:  utils.StringToInt(c.Query("min_rating"), 0),
		MaxRating:  utils.StringToInt(c.Query("max_rating"), 5),
		Status:     c.Query("status", "published"),
		SortBy:     c.Query("sort_by", "date"),
		SortOrder:  c.Query("sort_order", "desc"),
		Page:       utils.StringToInt(c.Query("page"), 1),
		Limit:      utils.StringToInt(c.Query("limit"), 20),
	}

	reviews, total, err := h.services.Review().GetReviews(c.Context(), filter)
	if err != nil {
		log.Printf("Error getting reviews: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.list.error.fetch_failed")
	}

	return utils.SuccessResponse(c, ReviewsListResponse{
		Success: true,
		Data:    reviews,
		Meta: ReviewsMeta{
			Total: int(total),
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	})
}

// GetEntityStats returns entity review statistics
// @Summary Get entity review statistics
// @Description Returns detailed review statistics for a specific entity
// @Tags reviews
// @Accept json
// @Produce json
// @Param type path string true "Entity type"
// @Param id path int true "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ReviewStats} "Entity review statistics"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid entity ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/stats/{type}/{id} [get]
func (h *ReviewHandler) GetEntityStats(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	log.Printf("Getting stats for %s ID=%d", entityType, entityId)

	stats, err := h.services.Review().GetReviewStats(c.Context(), entityType, entityId)
	if err != nil {
		log.Printf("Error getting stats for %s ID=%d: %v", entityType, entityId, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.stats.error.fetch_failed")
	}

	log.Printf("Got stats for %s ID=%d: %+v", entityType, entityId, stats)

	return utils.SuccessResponse(c, stats)
}

// GetEntityRating returns entity rating
// @Summary Get entity rating
// @Description Returns average rating for a specific entity
// @Tags reviews
// @Accept json
// @Produce json
// @Param type path string true "Entity type"
// @Param id path int true "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=RatingResponse} "Entity rating"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid entity ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/rating/{type}/{id} [get]
func (h *ReviewHandler) GetEntityRating(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	rating, err := h.services.Review().GetEntityRating(c.Context(), entityType, entityId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.rating.error.fetch_failed")
	}

	return utils.SuccessResponse(c, RatingResponse{
		Success: true,
		Rating:  rating,
	})
}

// CanReviewListing проверяет может ли пользователь оставить отзыв на listing
// @Summary Check if user can review listing
// @Description Checks if the current user can leave a review for the specified listing
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.CanReviewResponse} "Permission check result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/can-review/listing/{id} [get]
func (h *ReviewHandler) CanReviewListing(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, "listing", entityID)
	if err != nil {
		log.Printf("Error checking review permission for listing: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// CanReviewUser проверяет может ли пользователь оставить отзыв на user
// @Summary Check if user can review user
// @Description Checks if the current user can leave a review for the specified user
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.CanReviewResponse} "Permission check result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/can-review/user/{id} [get]
func (h *ReviewHandler) CanReviewUser(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, "user", entityID)
	if err != nil {
		log.Printf("Error checking review permission for user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// CanReviewStorefront проверяет может ли пользователь оставить отзыв на storefront
// @Summary Check if user can review storefront
// @Description Checks if the current user can leave a review for the specified storefront
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.CanReviewResponse} "Permission check result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/can-review/storefront/{id} [get]
func (h *ReviewHandler) CanReviewStorefront(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, "storefront", entityID)
	if err != nil {
		log.Printf("Error checking review permission for storefront: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// ===============================================
// Заглушки для POST/PUT/DELETE методов
// ===============================================

// GetStorefrontReviews - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetStorefrontReviews(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetStorefrontRatingSummary - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetStorefrontRatingSummary(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetUserAggregatedRating - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetUserAggregatedRating(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetStorefrontAggregatedRating - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetStorefrontAggregatedRating(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetUserReviews - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetUserReviews(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetUserRatingSummary - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetUserRatingSummary(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetStats - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetStats(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetReviewByID - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) GetReviewByID(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// CreateDraftReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) CreateDraftReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// UploadPhotos - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) UploadPhotos(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// PublishReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) PublishReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// UpdateReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// DeleteReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// VoteForReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) VoteForReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// AddResponse - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) AddResponse(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// UploadPhotosForNewReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) UploadPhotosForNewReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// ConfirmReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) ConfirmReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// DisputeReview - TODO: temporarily disabled during refactoring
func (h *ReviewHandler) DisputeReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}
