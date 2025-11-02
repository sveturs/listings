// backend/internal/proj/reviews/handler/reviews.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
)

// Create - TODO: temporarily disabled during refactoring
func (h *Handler) Create(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetByEntityID - TODO: temporarily disabled during refactoring
func (h *Handler) GetByEntityID(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// GetByID - TODO: temporarily disabled during refactoring
func (h *Handler) GetByID(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// Update - TODO: temporarily disabled during refactoring
func (h *Handler) Update(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// Delete - TODO: temporarily disabled during refactoring
func (h *Handler) Delete(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

// ReviewHandler - обработчик отзывов
type ReviewHandler struct {
	services globalService.ServicesInterface
}

// NewReviewHandler создаёт новый обработчик
func NewReviewHandler(services globalService.ServicesInterface) *ReviewHandler {
	return &ReviewHandler{
		services: services,
	}
}

// Все методы ReviewHandler - заглушки
func (h *ReviewHandler) GetStorefrontReviews(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetStorefrontRatingSummary(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetEntityRating(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetEntityStats(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetUserAggregatedRating(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetStorefrontAggregatedRating(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetUserReviews(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetUserRatingSummary(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetReviews(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetStats(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) CanReviewListing(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) CanReviewUser(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) CanReviewStorefront(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) GetReviewByID(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) CreateDraftReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) UploadPhotos(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) PublishReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) VoteForReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) AddResponse(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) UploadPhotosForNewReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) ConfirmReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}

func (h *ReviewHandler) DisputeReview(c *fiber.Ctx) error {
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
}
