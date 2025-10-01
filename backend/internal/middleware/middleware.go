// backend/internal/middleware/middleware.go
package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

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
	config            *config.Config
	services          globalService.ServicesInterface
	metrics           *monitoring.MetricsCollector
	authServicePubKey *rsa.PublicKey
	authService       *authService.AuthService
	jwtParserMW       fiber.Handler
}

// JWTParser возвращает middleware для парсинга JWT токенов из auth service
func (m *Middleware) JWTParser() fiber.Handler {
	return m.jwtParserMW
}

func NewMiddleware(cfg *config.Config, services globalService.ServicesInterface, authSvc *authService.AuthService, jwtParser fiber.Handler) *Middleware {
	m := &Middleware{
		config:      cfg,
		services:    services,
		metrics:     monitoring.NewMetricsCollector(pkglogger.GetLogger()),
		authService: authSvc,
		jwtParserMW: jwtParser,
	}

	// Загружаем публичный ключ auth service
	if err := m.loadAuthServicePublicKey(); err != nil {
		logger.Error().Err(err).Msg("Failed to load auth service public key, RS256 tokens will not be validated")
	}

	return m
}

func (m *Middleware) loadAuthServicePublicKey() error {
	// Загружаем из пути, указанного в конфигурации
	pubKeyPath := m.config.AuthServicePubKeyPath
	if pubKeyPath == "" {
		// Используем дефолтный путь если не указан в конфиге
		pubKeyPath = "keys/auth_service_public.pem"
	}
	logger.Info().Str("path", pubKeyPath).Msg("Loading auth service public key from path")

	// Проверяем существование файла
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		logger.Error().Str("path", pubKeyPath).Msg("Public key file does not exist")
		return fmt.Errorf("public key file not found: %s", pubKeyPath)
	}

	pubKeyData, err := os.ReadFile(pubKeyPath) // #nosec G304 - path is from config or hardcoded default
	if err != nil {
		logger.Error().Err(err).Str("path", pubKeyPath).Msg("Failed to read public key file")
		return err
	}

	// Временное логирование для отладки
	logger.Info().
		Str("path", pubKeyPath).
		Int("size", len(pubKeyData)).
		Str("first_line", strings.Split(string(pubKeyData), "\n")[0]).
		Msg("Public key file loaded, details")

	block, _ := pem.Decode(pubKeyData)
	if block == nil {
		logger.Error().Str("content_preview", string(pubKeyData[:100])).Msg("Failed to parse PEM block")
		return errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse public key")
		return err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		logger.Error().Msg("Not an RSA public key")
		return errors.New("not an RSA public key")
	}

	m.authServicePubKey = rsaPub

	// Временное детальное логирование ключа для отладки
	modulusBytes := rsaPub.N.Bytes()
	previewLen := 32
	if len(modulusBytes) < previewLen {
		previewLen = len(modulusBytes)
	}

	logger.Info().
		Int("key_size", rsaPub.Size()).
		Int("exponent", rsaPub.E).
		Str("modulus_first_32_bytes", fmt.Sprintf("%x", modulusBytes[:previewLen])).
		Msg("Auth service RSA public key loaded successfully with details")

	return nil
}

// Setup - deprecated, middleware регистрируется в server.go
// func (m *Middleware) Setup(app *fiber.App) {
// 	app.Use(m.Logger())
// 	app.Use(m.CORS())
// }

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

	// ПЕРВЫМ ДЕЛОМ проверяем роли из JWT токена (установлены в JWTParser из sveturs/auth)
	if roles, ok := c.Locals("roles").([]string); ok {
		for _, role := range roles {
			if role == "admin" {
				logger.Info().
					Int("user_id", userID).
					Str("source", "jwt_token").
					Strs("roles", roles).
					Msg("AdminRequired: Access granted - user has admin role in JWT")
				c.Locals("admin_id", userID)
				// ИСПРАВЛЕНИЕ: не возвращаем результат c.Next(), а просто вызываем его
				return nil
			}
		}
	}

	// Затем проверяем ID пользователя (для обратной совместимости)
	if userID == 1 || userID == 2 || userID == 3 || userID == 6 || userID == 11 {
		logger.Info().
			Int("user_id", userID).
			Msg("AdminRequired: Access granted for hardcoded user ID")
		// Устанавливаем admin_id для использования в handlers
		c.Locals("admin_id", userID)
		// ИСПРАВЛЕНИЕ: не возвращаем результат c.Next(), а просто вызываем его
		return nil
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
