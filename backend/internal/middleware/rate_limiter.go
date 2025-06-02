package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"backend/pkg/utils"
)

// AuthRateLimit создает rate limiter для аутентификации
// Ограничивает количество попыток входа/регистрации для защиты от brute force атак
func (m *Middleware) AuthRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 5 попыток входа за 15 минут
		Max:        5,
		Expiration: 15 * time.Minute,
		
		// Генерация ключа по IP адресу для контроля попыток
		KeyGenerator: func(c *fiber.Ctx) string {
			return "auth_" + c.IP()
		},
		
		// Обработка превышения лимита
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "users.errors.tooManyAttempts")
		},
		
		// Пропускаем localhost в режиме разработки
		Next: func(c *fiber.Ctx) bool {
			if m.config.Environment == "development" || m.config.Environment == "dev" {
				return c.IP() == "127.0.0.1" || c.IP() == "::1"
			}
			return false
		},
		
		// Не учитываем неуспешные запросы (например, invalid JSON)
		SkipFailedRequests: true,
		
		// Учитываем только успешные запросы для более точного контроля
		SkipSuccessfulRequests: false,
	})
}

// StrictAuthRateLimit создает более строгий rate limiter для повторных нарушений
// Применяется после первого превышения лимита для усиления защиты
func (m *Middleware) StrictAuthRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 2 попытки за 1 час после первого нарушения
		Max:        2,
		Expiration: 1 * time.Hour,
		
		// Ключ для строгого лимита
		KeyGenerator: func(c *fiber.Ctx) string {
			return "strict_auth_" + c.IP()
		},
		
		// Более строгий ответ при превышении
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "users.errors.accountTemporarilyLocked")
		},
		
		// Не пропускаем никого в строгом режиме
		Next: nil,
		
		SkipFailedRequests: true,
		SkipSuccessfulRequests: false,
	})
}

// RegistrationRateLimit создает отдельный rate limiter для регистрации
// Более мягкие ограничения, так как регистрация происходит реже
func (m *Middleware) RegistrationRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 3 регистрации за 1 час с одного IP
		Max:        3,
		Expiration: 1 * time.Hour,
		
		KeyGenerator: func(c *fiber.Ctx) string {
			return "register_" + c.IP()
		},
		
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "users.register.error.tooManyAttempts")
		},
		
		// Пропускаем localhost в режиме разработки
		Next: func(c *fiber.Ctx) bool {
			if m.config.Environment == "development" || m.config.Environment == "dev" {
				return c.IP() == "127.0.0.1" || c.IP() == "::1"
			}
			return false
		},
		
		SkipFailedRequests: true,
		SkipSuccessfulRequests: false,
	})
}