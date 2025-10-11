package handler

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/middleware"
)

// RegisterRoutes registers subscription routes
func (h *SubscriptionHandler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) {
	// Public routes
	public := app.Group("/api/v1/subscriptions")
	public.Get("/plans", h.GetPlans)

	// Protected routes - require authentication (using library middleware)
	protected := app.Group("/api/v1/subscriptions", h.jwtParserMW, authMiddleware.RequireAuthString())
	protected.Get("/current", h.GetCurrentSubscription)
	protected.Post("/", h.CreateSubscription)
	protected.Post("/upgrade", h.UpgradeSubscription)
	protected.Post("/cancel", h.CancelSubscription)
	protected.Post("/check-limits", h.CheckLimits)
	protected.Post("/initiate-payment", h.InitiatePayment)
	protected.Post("/complete-payment", h.CompletePayment)

	// Admin routes (using library middleware with admin role)
	admin := app.Group("/api/v1/admin/subscriptions", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))
	admin.Get("/users/:user_id/subscription", h.AdminGetUserSubscription)
}
