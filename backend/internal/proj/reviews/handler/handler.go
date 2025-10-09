// Package handler содержит обработчики HTTP запросов для работы с отзывами
// backend/internal/proj/reviews/handler/handler.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
)

func NewHandler(services globalService.ServicesInterface, jwtParserMW fiber.Handler) *Handler {
	return &Handler{
		Review:      NewReviewHandler(services),
		jwtParserMW: jwtParserMW,
	}
}

type Handler struct {
	Review      *ReviewHandler
	jwtParserMW fiber.Handler
}

func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	app.Get("/api/v1/public/b2c_stores/:id/reviews", h.Review.GetStorefrontReviews)
	app.Get("/api/v1/public/b2c_stores/:id/rating", h.Review.GetStorefrontRatingSummary)

	entityStats := app.Group("/api/v1/entity")
	entityStats.Get("/:type/:id/rating", h.Review.GetEntityRating)
	entityStats.Get("/:type/:id/stats", h.Review.GetEntityStats)

	// Новые endpoints для агрегированных рейтингов
	app.Get("/api/v1/users/:id/aggregated-rating", h.Review.GetUserAggregatedRating)
	app.Get("/api/v1/b2c_stores/:id/aggregated-rating", h.Review.GetStorefrontAggregatedRating)

	// Protected GET endpoints (требуют только Auth, БЕЗ CSRF так как это GET)
	app.Get("/api/v1/users/:id/reviews", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.GetUserReviews)
	app.Get("/api/v1/users/:id/rating", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.GetUserRatingSummary)
	app.Get("/api/v1/b2c_stores/:id/reviews", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.GetStorefrontReviews)
	app.Get("/api/v1/b2c_stores/:id/rating", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.GetStorefrontRatingSummary)

	// Публичные GET endpoints (без middleware)
	app.Get("/api/v1/reviews", h.Review.GetReviews)
	app.Get("/api/v1/reviews/stats", h.Review.GetStats)

	// can-review endpoint - отдельные роуты для каждого типа чтобы избежать Fiber routing багов
	app.Get("/api/v1/review-permission/listing/:id", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.CanReviewListing)
	app.Get("/api/v1/review-permission/user/:id", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.CanReviewUser)
	app.Get("/api/v1/review-permission/storefront/:id", h.jwtParserMW, authMiddleware.RequireAuth(), h.Review.CanReviewStorefront)

	// ⚠️ ВАЖНО: Параметризованный роут /reviews/:id должен быть ПОСЛЕДНИМ среди GET
	// Fiber может ошибочно матчить /review-permission/listing/351 как /reviews/:id
	// если этот роут зарегистрирован раньше
	app.Get("/api/v1/reviews/:id", h.Review.GetReviewByID)

	// Protected POST/PUT/DELETE endpoints (требуют и Auth и CSRF)
	app.Post("/api/v1/reviews/draft", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.CreateDraftReview)
	app.Post("/api/v1/reviews/:id/photos", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.UploadPhotos)
	app.Post("/api/v1/reviews/:id/publish", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.PublishReview)
	app.Put("/api/v1/reviews/:id", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.UpdateReview)
	app.Delete("/api/v1/reviews/:id", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.DeleteReview)
	app.Post("/api/v1/reviews/:id/vote", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.VoteForReview)
	app.Post("/api/v1/reviews/:id/response", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.AddResponse)
	app.Post("/api/v1/reviews/upload-photos", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.UploadPhotosForNewReview)
	app.Post("/api/v1/reviews/:id/confirm", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.ConfirmReview)
	app.Post("/api/v1/reviews/:id/dispute", h.jwtParserMW, authMiddleware.RequireAuth(), mw.CSRFProtection(), h.Review.DisputeReview)

	return nil
}

func (h *Handler) GetPrefix() string {
	return "/api/v1/review"
}
