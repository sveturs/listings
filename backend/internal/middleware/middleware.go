// backend/internal/middleware/middleware.go
package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
)

type Middleware struct {
	config   *config.Config
	services globalService.ServicesInterface
}

func NewMiddleware(cfg *config.Config, services globalService.ServicesInterface) *Middleware {
	return &Middleware{
		config:   cfg,
		services: services,
	}
}

func (m *Middleware) Setup(app *fiber.App) {
	app.Use(m.Logger())
	app.Use(m.CORS())
}

// ErrorHandler обрабатывает все ошибки приложения
func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	logger.Error().Err(err).Msg("Error handling request")

	return utils.ErrorResponse(c, code, message)
}

func (m *Middleware) AdminRequired(c *fiber.Ctx) error {
	logger.Info().Str("path", c.Path()).Msg("AdminRequired middleware called")

	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	logger.Info().
		Int("user_id", userID).
		Bool("user_id_ok", ok).
		Msg("AdminRequired: checking user ID")

	if !ok || userID == 0 {
		logger.Info().Msg("AdminRequired: User ID not found or invalid")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Сначала проверяем ID пользователя (для обратной совместимости)
	if userID == 1 || userID == 2 || userID == 3 {
		logger.Info().
			Int("user_id", userID).
			Msg("AdminRequired: Access granted for hardcoded user ID")
		// Устанавливаем admin_id для использования в handlers
		c.Locals("admin_id", userID)
		return c.Next()
	}

	logger.Info().
		Int("user_id", userID).
		Msg("AdminRequired: User ID not in hardcoded list, checking email")

	// Если ID не подходит, проверяем email пользователя
	ctx := c.Context()

	if m.services == nil {
		logger.Error().Msg("AdminRequired: services is nil")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	userService := m.services.User()
	if userService == nil {
		logger.Error().Msg("AdminRequired: User service is nil")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	user, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error().
			Err(err).
			Int("user_id", userID).
			Msg("AdminRequired: Error getting user")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
	}

	if user == nil {
		logger.Info().
			Int("user_id", userID).
			Msg("AdminRequired: User not found")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
	}

	logger.Info().
		Int("user_id", userID).
		Str("email", user.Email).
		Msg("AdminRequired: Found user")

	// Проверка администратора по БД (нужно учитывать, что таблица может еще не существовать)
	isAdmin, err := userService.IsUserAdmin(ctx, user.Email)
	if err != nil {
		// Если таблица еще не создана, проверяем по жесткому списку
		logger.Error().Err(err).Msg("AdminRequired: Error checking admin status in DB, falling back to hardcoded list")
	} else if isAdmin {
		logger.Info().
			Str("email", user.Email).
			Str("source", "database").
			Msg("AdminRequired: Access granted for admin")
		// Устанавливаем admin_id для использования в handlers
		c.Locals("admin_id", userID)
		return c.Next()
	}

	// Резервный список email-адресов администраторов (на случай если миграция еще не применена)
	adminEmails := []string{
		"bevzenko.sergey@gmail.com",
		"voroshilovdo@gmail.com",
		"admin@svetu.rs",
		// Здесь можно добавить дополнительные email-адреса администраторов
	}

	logger.Info().
		Str("email", user.Email).
		Strs("admin_emails", adminEmails).
		Msg("AdminRequired: Checking user email against hardcoded list")

	// Проверяем, является ли email пользователя админским
	for _, adminEmail := range adminEmails {
		if user.Email == adminEmail {
			logger.Info().
				Str("email", user.Email).
				Str("source", "hardcoded").
				Msg("AdminRequired: Access granted for admin")
			// Устанавливаем admin_id для использования в handlers
			c.Locals("admin_id", userID)
			return c.Next()
		}
	}

	logger.Info().
		Int("user_id", userID).
		Str("email", user.Email).
		Msg("AdminRequired: Access denied - user is not an admin")
	return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
}

// RequireAdmin является алиасом для AdminRequired для удобства
func (m *Middleware) RequireAdmin() fiber.Handler {
	return m.AdminRequired
}

// RequireAuth требует обязательную аутентификацию
func (m *Middleware) RequireAuth() fiber.Handler {
	return m.AuthRequiredJWT
}

// OptionalAuth опциональная аутентификация
func (m *Middleware) OptionalAuth() fiber.Handler {
	return m.OptionalAuthJWT
}
