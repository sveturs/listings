// backend/internal/proj/users/handler/routes.go
package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует маршруты с использованием auth-service
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные auth эндпоинты (без middleware аутентификации)
	// Все handlers проксируют к auth-service
	app.Post("/api/v1/auth/register", mw.RegistrationRateLimit(), h.Auth.Register)
	app.Post("/api/v1/auth/login", mw.AuthRateLimit(), h.Auth.Login)
	app.Post("/api/v1/auth/logout", mw.RateLimitByIP(10, time.Minute), h.Auth.Logout)
	// Временно отключаем rate limit для refresh в development из-за проблемы с частыми вызовами
	// TODO: исправить логику refresh на frontend чтобы не было избыточных вызовов
	app.Post("/api/v1/auth/refresh", h.Auth.RefreshToken)

	// Защищенные auth эндпоинты (требуют токен)
	app.Post("/api/v1/auth/logout", h.jwtParserMW, authMiddleware.RequireAuth(), mw.RateLimitByIP(10, time.Minute), h.Auth.Logout)
	app.Get("/api/v1/auth/logout", h.jwtParserMW, authMiddleware.RequireAuth(), h.Auth.Logout)
	app.Get("/api/v1/auth/session", h.jwtParserMW, authMiddleware.RequireAuth(), h.Auth.GetSession)
	app.Get("/api/v1/auth/validate", h.jwtParserMW, authMiddleware.RequireAuth(), h.Auth.Validate)
	app.Get("/api/v1/auth/me", h.jwtParserMW, authMiddleware.RequireAuth(), h.Auth.GetCurrentUser)

	// OAuth эндпоинты
	app.Get("/api/v1/auth/google", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleAuth)
	app.Get("/api/v1/auth/google/callback", mw.RateLimitByIP(10, time.Minute), h.Auth.GoogleCallback)

	// Admin check остается локальным (использует нашу БД)
	app.Get("/api/v1/admin-check/:email", h.User.IsAdminPublic)

	// User profile endpoints (БЕЗ CSRF - используем BFF proxy архитектуру)
	users := app.Group("/api/v1/users", h.jwtParserMW, authMiddleware.RequireAuthString())
	users.Get("/profile", h.User.GetProfile)
	users.Put("/profile", h.User.UpdateProfile)
	users.Get("/:id/profile", h.User.GetProfileByID)
	users.Get("/privacy-settings", h.User.GetPrivacySettings)
	users.Put("/privacy-settings", h.User.UpdatePrivacySettings)
	users.Get("/chat-settings", h.User.GetChatSettings)
	users.Put("/chat-settings", h.User.UpdateChatSettings)

	// Roles endpoints for specific users
	users.Get("/:userId/roles", h.Auth.GetUserRoles)

	// Admin endpoints (БЕЗ CSRF - используем BFF proxy архитектуру)
	adminUsersRoutes := app.Group("/api/v1/admin/users", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))
	adminUsersRoutes.Get("/", h.User.GetAllUsers)
	adminUsersRoutes.Get("/:id", h.User.GetUserByIDAdmin)
	adminUsersRoutes.Put("/:id", h.User.UpdateUserAdmin)
	adminUsersRoutes.Put("/:id/status", h.User.UpdateUserStatus)
	adminUsersRoutes.Put("/:id/role", h.User.UpdateUserRole)
	adminUsersRoutes.Delete("/:id", h.User.DeleteUser)
	adminUsersRoutes.Get("/:id/balance", h.User.GetUserBalance)
	adminUsersRoutes.Get("/:id/transactions", h.User.GetUserTransactions)

	// Roles management (general endpoints)
	rolesRoutes := app.Group("/api/v1/roles", h.jwtParserMW, authMiddleware.RequireAuthString())
	rolesRoutes.Get("/", h.Auth.GetRoles)
	rolesRoutes.Post("/assign", h.Auth.AssignRole)
	rolesRoutes.Post("/revoke", h.Auth.RevokeRole)

	// Admin roles management
	adminRolesRoutes := app.Group("/api/v1/admin/roles", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))
	adminRolesRoutes.Get("/", h.User.GetAllRoles)

	// Admin management (БЕЗ CSRF - используем BFF proxy архитектуру)
	adminAdminsRoutes := app.Group("/api/v1/admin/admins", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))
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
