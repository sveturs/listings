// backend/internal/proj/balance/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта balance
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Защищенные маршруты с аутентификацией
	balanceRoutes := app.Group("/api/v1/balance", mw.AuthRequiredJWT)

	balanceRoutes.Get("/", h.Balance.GetBalance)
	balanceRoutes.Get("/transactions", h.Balance.GetTransactions)
	balanceRoutes.Get("/payment-methods", h.Balance.GetPaymentMethods)
	balanceRoutes.Post("/deposit", h.Balance.CreateDeposit)

	// Mock payment completion endpoint (только для разработки)
	balanceRoutes.Post("/mock/complete", h.Balance.CompleteMockPayment)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/balance"
}
