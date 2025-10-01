package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes registers subscription routes
func (h *SubscriptionHandler) RegisterRoutes(app *fiber.App, authMiddleware *middleware.Middleware) {
	// Public routes
	public := app.Group("/api/v1/subscriptions")
	public.Get("/plans", h.GetPlans)

	// Protected routes - require authentication
	protected := app.Group("/api/v1/subscriptions", authMiddleware.RequireAuth())
	protected.Get("/current", h.GetCurrentSubscription)
	protected.Post("/", h.CreateSubscription)
	protected.Post("/upgrade", h.UpgradeSubscription)
	protected.Post("/cancel", h.CancelSubscription)
	protected.Post("/check-limits", h.CheckLimits)
	protected.Post("/initiate-payment", h.InitiatePayment)
	protected.Post("/complete-payment", h.CompletePayment)

	// Admin routes
	admin := app.Group("/api/v1/admin/subscriptions", authMiddleware.RequireAuth(), authMiddleware.RequireAdmin())
	admin.Get("/users/:user_id/subscription", h.AdminGetUserSubscription)
}
