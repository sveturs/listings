// backend/internal/proj/notifications/routes.go
package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта users
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {

	app.Post("/api/v1/auth/register", mw.RegistrationRateLimit(), mw.CSRFProtection(), h.Auth.Register)
	app.Post("/api/v1/auth/login", mw.AuthRateLimit(), mw.CSRFProtection(), h.Auth.Login)
	app.Post("/api/v1/auth/logout", mw.RateLimitByIP(10, time.Minute), h.Auth.Logout)
	// Временно отключаем rate limit для refresh в development из-за проблемы с частыми вызовами
	// TODO: исправить логику refresh на frontend чтобы не было избыточных вызовов
	app.Post("/api/v1/auth/refresh", h.Auth.RefreshToken)
	app.Get("/api/v1/auth/session", h.Auth.GetSession)
	app.Get("/api/v1/auth/google", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleAuth)
	app.Get("/api/v1/auth/google/callback", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleCallback)
	app.Get("/auth/google", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleAuth)              // TODO: изменить в google console и потом удалить тут
	app.Get("/auth/google/callback", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleCallback) // TODO: изменить в google console и потом удалить тут
	app.Get("/api/v1/auth/logout", h.Auth.Logout)
	app.Post("/api/v1/auth/logout", mw.CSRFProtection(), h.Auth.Logout) // Поддержка POST для logout
	app.Get("/api/v1/admin-check/:email", h.User.IsAdminPublic)

	users := app.Group("/api/v1/users", mw.AuthRequiredJWT, mw.CSRFProtection())
	users.Get("/me", h.User.GetProfile)    // TODO: remove
	users.Put("/me", h.User.UpdateProfile) // TODO: remove
	users.Get("/profile", h.User.GetProfile)
	users.Put("/profile", h.User.UpdateProfile)
	users.Get("/:id/profile", h.User.GetProfileByID)

	adminRoutes := app.Group("/api/v1/admin", mw.AuthRequiredJWT, mw.AdminRequired, mw.CSRFProtection())
	adminRoutes.Get("/users", h.User.GetAllUsers)
	adminRoutes.Get("/users/:id", h.User.GetUserByIDAdmin)
	adminRoutes.Put("/users/:id", h.User.UpdateUserAdmin)
	adminRoutes.Put("/users/:id/status", h.User.UpdateUserStatus)
	adminRoutes.Delete("/users/:id", h.User.DeleteUser)
	adminRoutes.Get("/users/:id/balance", h.User.GetUserBalance)
	adminRoutes.Get("/users/:id/transactions", h.User.GetUserTransactions)
	adminRoutes.Get("/admins", h.User.GetAllAdmins)
	adminRoutes.Post("/admins", h.User.AddAdmin)
	adminRoutes.Delete("/admins/:email", h.User.RemoveAdmin)
	adminRoutes.Get("/admins/check/:email", h.User.IsAdmin)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/users"
}
