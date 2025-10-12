// backend/internal/middleware/middleware.go
package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/monitoring"
	globalService "backend/internal/proj/global/service"
	pkglogger "backend/pkg/logger"
	"backend/pkg/utils"

	authService "github.com/sveturs/auth/pkg/http/service"
)

type Middleware struct {
	config      *config.Config
	services    globalService.ServicesInterface
	metrics     *monitoring.MetricsCollector
	authService *authService.AuthService
	jwtParserMW fiber.Handler
}

// JWTParser возвращает middleware для парсинга JWT токенов из auth service
func (m *Middleware) JWTParser() fiber.Handler {
	return m.jwtParserMW
}

func NewMiddleware(cfg *config.Config, services globalService.ServicesInterface, authSvc *authService.AuthService, jwtParser fiber.Handler) *Middleware {
	return &Middleware{
		config:      cfg,
		services:    services,
		metrics:     monitoring.NewMetricsCollector(pkglogger.GetLogger()),
		authService: authSvc,
		jwtParserMW: jwtParser,
	}
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

	// Проверяем user_id из JWT (установлен в JWTParser из sveturs/auth)
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		logger.Info().Msg("AdminRequired: User ID not found or invalid - требуется авторизация")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Проверяем роли из JWT токена (установлены в JWTParser из sveturs/auth)
	// Auth-service - единственный источник истины для ролей!
	if roles, ok := c.Locals("roles").([]string); ok {
		for _, role := range roles {
			if role == "admin" {
				logger.Info().
					Int("user_id", userID).
					Strs("roles", roles).
					Msg("AdminRequired: Access granted - user has admin role")
				c.Locals("admin_id", userID)
				return c.Next()
			}
		}
	}

	// Роль admin не найдена - доступ запрещён
	logger.Info().
		Int("user_id", userID).
		Msg("AdminRequired: Access denied - user does not have admin role")
	return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
}

// RequireAdmin является алиасом для AdminRequired для удобства
func (m *Middleware) RequireAdmin() fiber.Handler {
	return m.AdminRequired
}

// GetJWTParser возвращает JWT parser middleware из sveturs/auth
func (m *Middleware) GetJWTParser() fiber.Handler {
	return m.jwtParserMW
}

// GetAuthService возвращает auth service из sveturs/auth
func (m *Middleware) GetAuthService() *authService.AuthService {
	return m.authService
}

// AuthRequiredJWT это алиас для совместимости
// Использует комбинацию JWTParser и RequireAuth из sveturs/auth
func (m *Middleware) AuthRequiredJWT(c *fiber.Ctx) error {
	// Сначала парсим JWT
	if err := m.jwtParserMW(c); err != nil {
		return err
	}

	// Затем проверяем аутентификацию
	// Если userID есть в контексте - значит пользователь аутентифицирован
	if userID, ok := c.Locals("user_id").(int); !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	return c.Next()
}

// OptionalAuthJWT это алиас для совместимости
// Просто парсит JWT если он есть, но не требует его
func (m *Middleware) OptionalAuthJWT(c *fiber.Ctx) error {
	// Просто парсим JWT если он есть
	if m.jwtParserMW != nil {
		return m.jwtParserMW(c)
	}
	return c.Next()
}

// RequireAuth требует обязательную аутентификацию
// Возвращает middleware, которое требует аутентификацию
func (m *Middleware) RequireAuth() fiber.Handler {
	// Возвращаем функцию, которая использует AuthRequiredJWT
	return func(c *fiber.Ctx) error {
		return m.AuthRequiredJWT(c)
	}
}

// OptionalAuth опциональная аутентификация
// Возвращает middleware, которое опционально проверяет JWT
func (m *Middleware) OptionalAuth() fiber.Handler {
	// Возвращаем функцию, которая использует OptionalAuthJWT
	return func(c *fiber.Ctx) error {
		return m.OptionalAuthJWT(c)
	}
}
