// Package handler содержит обработчики HTTP запросов для работы с отзывами
// backend/internal/proj/reviews/handler/handler.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
)

func NewHandler(services globalService.ServicesInterface) *Handler {

	return &Handler{
		Review: NewReviewHandler(services),
	}
}

type Handler struct {
	Review *ReviewHandler
}

func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	app.Get("/api/v1/public/storefronts/:id/reviews", h.Review.GetStorefrontReviews)
	app.Get("/api/v1/public/storefronts/:id/rating", h.Review.GetStorefrontRatingSummary)

	entityStats := app.Group("/api/v1/entity")
	entityStats.Get("/:type/:id/rating", h.Review.GetEntityRating)
	entityStats.Get("/:type/:id/stats", h.Review.GetEntityStats)
	
	// Новые endpoints для агрегированных рейтингов
	app.Get("/api/v1/users/:id/aggregated-rating", h.Review.GetUserAggregatedRating)
	app.Get("/api/v1/storefronts/:id/aggregated-rating", h.Review.GetStorefrontAggregatedRating)
	
	// Endpoint для проверки возможности оставить отзыв
	app.Get("/api/v1/reviews/can-review/:type/:id", mw.AuthRequiredJWT, h.Review.CanReview)

	review := app.Group("/api/v1/reviews")
	review.Get("/", h.Review.GetReviews)
	review.Get("/:id", h.Review.GetReviewByID)
	review.Get("/stats", h.Review.GetStats)

	app.Get("/api/v1/users/:id/reviews", mw.AuthRequiredJWT, mw.CSRFProtection(), h.Review.GetUserReviews)
	app.Get("/api/v1/users/:id/rating", mw.AuthRequiredJWT, mw.CSRFProtection(), h.Review.GetUserRatingSummary)
	app.Get("/api/v1/storefronts/:id/reviews", mw.AuthRequiredJWT, mw.CSRFProtection(), h.Review.GetStorefrontReviews)
	app.Get("/api/v1/storefronts/:id/rating", mw.AuthRequiredJWT, mw.CSRFProtection(), h.Review.GetStorefrontRatingSummary)

	protectedReviews := app.Group("/api/v1/reviews", mw.AuthRequiredJWT, mw.CSRFProtection())
	protectedReviews.Post("/", h.Review.CreateReview)
	protectedReviews.Put("/:id", h.Review.UpdateReview)
	protectedReviews.Delete("/:id", h.Review.DeleteReview)
	protectedReviews.Post("/:id/vote", h.Review.VoteForReview)
	protectedReviews.Post("/:id/response", h.Review.AddResponse)
	protectedReviews.Post("/:id/photos", h.Review.UploadPhotos)
	
	// Новые endpoints для подтверждений и споров
	protectedReviews.Post("/:id/confirm", h.Review.ConfirmReview)
	protectedReviews.Post("/:id/dispute", h.Review.DisputeReview)

	return nil
}

func (h *Handler) GetPrefix() string {
	return "/api/v1/review"
}
