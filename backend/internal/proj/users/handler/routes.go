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
	// OAuth routes are now handled by Auth Service proxy middleware
	// These routes are commented out to prevent conflicts
	// app.Get("/api/v1/auth/google", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleAuth)
	// app.Get("/api/v1/auth/google/callback", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleCallback)
	// app.Get("/auth/google", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleAuth)
	// app.Get("/auth/google/callback", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleCallback)
	app.Get("/api/v1/auth/logout", h.Auth.Logout)
	app.Post("/api/v1/auth/logout", mw.CSRFProtection(), h.Auth.Logout) // Поддержка POST для logout
	app.Get("/api/v1/admin-check/:email", h.User.IsAdminPublic)

	users := app.Group("/api/v1/users", mw.AuthRequiredJWT, mw.CSRFProtection())
	users.Get("/me", h.User.GetProfile)    // TODO: remove
	users.Put("/me", h.User.UpdateProfile) // TODO: remove
	users.Get("/profile", h.User.GetProfile)
	users.Put("/profile", h.User.UpdateProfile)
	users.Get("/:id/profile", h.User.GetProfileByID)
	users.Get("/privacy-settings", h.User.GetPrivacySettings)
	users.Put("/privacy-settings", h.User.UpdatePrivacySettings)

	// Use specific route groups to avoid conflicts with marketplace admin routes
	adminUsersRoutes := app.Group("/api/v1/admin/users", mw.AuthRequiredJWT, mw.AdminRequired, mw.CSRFProtection())
	adminUsersRoutes.Get("/", h.User.GetAllUsers)
	adminUsersRoutes.Get("/:id", h.User.GetUserByIDAdmin)
	adminUsersRoutes.Put("/:id", h.User.UpdateUserAdmin)
	adminUsersRoutes.Put("/:id/status", h.User.UpdateUserStatus)
	adminUsersRoutes.Put("/:id/role", h.User.UpdateUserRole)
	adminUsersRoutes.Delete("/:id", h.User.DeleteUser)
	adminUsersRoutes.Get("/:id/balance", h.User.GetUserBalance)
	adminUsersRoutes.Get("/:id/transactions", h.User.GetUserTransactions)

	// Routes for roles management
	adminRolesRoutes := app.Group("/api/v1/admin/roles", mw.AuthRequiredJWT, mw.AdminRequired)
	adminRolesRoutes.Get("/", h.User.GetAllRoles)

	adminAdminsRoutes := app.Group("/api/v1/admin/admins", mw.AuthRequiredJWT, mw.AdminRequired, mw.CSRFProtection())
	adminAdminsRoutes.Get("/", h.User.GetAllAdmins)
	adminAdminsRoutes.Post("/", h.User.AddAdmin)
	adminAdminsRoutes.Delete("/:email", h.User.RemoveAdmin)
	adminAdminsRoutes.Get("/check/:email", h.User.IsAdmin)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/users"
}
