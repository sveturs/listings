package middleware

import (
	"fmt"
	"sync"
	"time"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"backend/internal/logger"
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
			if m.config.Environment == envDevelopment || m.config.Environment == envDev {
				return c.IP() == localhostIPv4 || c.IP() == localhostIPv6
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

		SkipFailedRequests:     true,
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
			if m.config.Environment == envDevelopment || m.config.Environment == envDev {
				return c.IP() == localhostIPv4 || c.IP() == localhostIPv6
			}
			return false
		},

		SkipFailedRequests:     true,
		SkipSuccessfulRequests: false,
	})
}

// RateLimiter структура для хранения информации о rate limiting
type RateLimiter struct {
	mu              sync.RWMutex
	requests        map[string]*userRequests
	cleanupInterval time.Duration
}

// userRequests хранит информацию о запросах пользователя
type userRequests struct {
	timestamps []time.Time
	mu         sync.Mutex
}

// NewRateLimiter создает новый экземпляр RateLimiter
func NewRateLimiter(cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:        make(map[string]*userRequests),
		cleanupInterval: cleanupInterval,
	}

	// Запускаем горутину для периодической очистки старых записей
	go rl.cleanup()

	return rl
}

// cleanup периодически очищает старые записи
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, ur := range rl.requests {
			ur.mu.Lock()
			// Удаляем записи старше часа
			cutoff := now.Add(-time.Hour)
			newTimestamps := make([]time.Time, 0)
			for _, ts := range ur.timestamps {
				if ts.After(cutoff) {
					newTimestamps = append(newTimestamps, ts)
				}
			}
			ur.timestamps = newTimestamps

			// Если нет активных запросов, удаляем пользователя
			if len(ur.timestamps) == 0 {
				delete(rl.requests, key)
			}
			ur.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// getUserRequests получает или создает запись для пользователя
func (rl *RateLimiter) getUserRequests(key string) *userRequests {
	rl.mu.RLock()
	ur, exists := rl.requests[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		ur, exists = rl.requests[key]
		if !exists {
			ur = &userRequests{
				timestamps: make([]time.Time, 0),
			}
			rl.requests[key] = ur
		}
		rl.mu.Unlock()
	}

	return ur
}

// isAllowed проверяет, разрешен ли запрос
func (rl *RateLimiter) isAllowed(key string, limit int, window time.Duration) bool {
	ur := rl.getUserRequests(key)

	ur.mu.Lock()
	defer ur.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// Удаляем старые timestamp
	newTimestamps := make([]time.Time, 0)
	for _, ts := range ur.timestamps {
		if ts.After(cutoff) {
			newTimestamps = append(newTimestamps, ts)
		}
	}
	ur.timestamps = newTimestamps

	// Проверяем лимит
	if len(ur.timestamps) >= limit {
		return false
	}

	// Добавляем новый timestamp
	ur.timestamps = append(ur.timestamps, now)
	return true
}

// RateLimitByUser создает middleware для rate limiting по пользователю
func (m *Middleware) RateLimitByUser(limit int, window time.Duration) fiber.Handler {
	rl := NewRateLimiter(5 * time.Minute)

	return func(c *fiber.Ctx) error {
		// Получаем userID из контекста
		userID, ok := authMiddleware.GetUserID(c)
		if !ok || userID == 0 {
			// Если нет userID, используем IP
			key := fmt.Sprintf("ip:%s", c.IP())
			if !rl.isAllowed(key, limit/2, window) { // Меньший лимит для неавторизованных
				logger.Info().
					Str("ip", c.IP()).
					Str("path", c.Path()).
					Msg("Rate limit exceeded for IP")
				return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many requests")
			}
			return c.Next()
		}

		// Для авторизованных пользователей
		key := fmt.Sprintf("user:%d", userID)
		if !rl.isAllowed(key, limit, window) {
			logger.Info().
				Int("user_id", userID).
				Str("path", c.Path()).
				Msg("Rate limit exceeded for user")
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many requests")
		}

		return c.Next()
	}
}

// RateLimitByIP создает middleware для rate limiting по IP
func (m *Middleware) RateLimitByIP(limit int, window time.Duration) fiber.Handler {
	rl := NewRateLimiter(5 * time.Minute)

	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("ip:%s:%s", c.IP(), c.Path())

		if !rl.isAllowed(key, limit, window) {
			logger.Info().
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Msg("Rate limit exceeded for IP")
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Too many requests")
		}

		return c.Next()
	}
}

// RateLimitMessages специальный rate limiter для сообщений
func (m *Middleware) RateLimitMessages() fiber.Handler {
	// Разные лимиты для разных типов действий
	messageLimiter := NewRateLimiter(5 * time.Minute)
	fileLimiter := NewRateLimiter(5 * time.Minute)

	return func(c *fiber.Ctx) error {
		userID, ok := authMiddleware.GetUserID(c)
		if !ok || userID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		path := c.Path()
		key := fmt.Sprintf("user:%d", userID)

		// Разные лимиты для разных эндпоинтов
		switch path {
		case "/api/v1/c2c/chat/messages":
			// 30 сообщений в минуту
			if !messageLimiter.isAllowed(key, 30, time.Minute) {
				logger.Info().
					Int("user_id", userID).
					Str("action", "send_message").
					Msg("Message rate limit exceeded")
				return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Слишком много сообщений. Подождите немного.")
			}

		case "/api/v1/c2c/chat/messages/:id/attachments":
			// 10 загрузок файлов в минуту
			if !fileLimiter.isAllowed(key, 10, time.Minute) {
				logger.Info().
					Int("user_id", userID).
					Str("action", "upload_file").
					Msg("File upload rate limit exceeded")
				return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "Слишком много загрузок файлов. Подождите немного.")
			}
		}

		return c.Next()
	}
}

// RefreshTokenRateLimit создает специальный rate limiter для refresh токенов
// Более строгие ограничения для защиты от атак на refresh endpoint
func (m *Middleware) RefreshTokenRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 5 запросов на обновление токена за 15 минут
		Max:        5,
		Expiration: 15 * time.Minute,

		// Генерация ключа по IP адресу
		KeyGenerator: func(c *fiber.Ctx) string {
			return "refresh_" + c.IP()
		},

		// Обработка превышения лимита
		LimitReached: func(c *fiber.Ctx) error {
			logger.Info().
				Str("ip", c.IP()).
				Str("action", "refresh_token").
				Msg("Refresh token rate limit exceeded")
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "users.auth.error.too_many_refresh_attempts")
		},

		// Не пропускаем localhost в продакшене
		Next: func(c *fiber.Ctx) bool {
			if m.config.Environment == envDevelopment || m.config.Environment == envDev {
				return c.IP() == localhostIPv4 || c.IP() == localhostIPv6
			}
			return false
		},

		SkipFailedRequests:     false, // Учитываем все попытки
		SkipSuccessfulRequests: false,
	})
}

// PaymentAPIRateLimit создает rate limiter для платежных API
// Строгие ограничения для защиты от злоупотреблений и мошенничества
func (m *Middleware) PaymentAPIRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 10 запросов на создание платежа за 1 минуту
		Max:        10,
		Expiration: 1 * time.Minute,

		// Генерация ключа по userID если авторизован, иначе по IP
		KeyGenerator: func(c *fiber.Ctx) string {
			userID, ok := authMiddleware.GetUserID(c)
			if ok && userID > 0 {
				return fmt.Sprintf("payment_user_%d", userID)
			}
			return fmt.Sprintf("payment_ip_%s", c.IP())
		},

		// Обработка превышения лимита
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn().
				Str("ip", c.IP()).
				Interface("user_id", c.Locals("user_id")).
				Str("path", c.Path()).
				Msg("Payment API rate limit exceeded")
			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "payments.errors.tooManyRequests")
		},

		// Пропускаем localhost только в режиме разработки
		Next: func(c *fiber.Ctx) bool {
			if m.config.Environment == envDevelopment || m.config.Environment == envDev {
				return c.IP() == localhostIPv4 || c.IP() == localhostIPv6
			}
			return false
		},

		SkipFailedRequests:     false, // Учитываем все попытки для безопасности
		SkipSuccessfulRequests: false,
	})
}

// WebhookRateLimit создает rate limiter для webhook endpoints
// Более высокие лимиты для внешних сервисов
func (m *Middleware) WebhookRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 100 webhook запросов за 1 минуту
		Max:        100,
		Expiration: 1 * time.Minute,

		// Генерация ключа по IP и endpoint
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("webhook_%s_%s", c.IP(), c.Path())
		},

		// Обработка превышения лимита
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn().
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Str("webhook_signature", c.Get("X-Webhook-Signature")).
				Msg("Webhook rate limit exceeded")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},

		// Webhooks не пропускаем даже в dev режиме для тестирования
		Next: nil,

		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

// StrictPaymentRateLimit создает очень строгий rate limiter для критических платежных операций
// Используется для операций типа capture, refund, etc.
func (m *Middleware) StrictPaymentRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		// Максимум 3 критические операции за 5 минут
		Max:        3,
		Expiration: 5 * time.Minute,

		// Генерация ключа по userID
		KeyGenerator: func(c *fiber.Ctx) string {
			userID, ok := authMiddleware.GetUserID(c)
			if ok && userID > 0 {
				return fmt.Sprintf("strict_payment_user_%d", userID)
			}
			// Для неавторизованных запрещаем полностью
			return fmt.Sprintf("strict_payment_ip_%s", c.IP())
		},

		// Обработка превышения лимита с уведомлением
		LimitReached: func(c *fiber.Ctx) error {
			userID, _ := c.Locals("user_id").(int)

			logger.Error().
				Str("ip", c.IP()).
				Interface("user_id", userID).
				Str("path", c.Path()).
				Msg("CRITICAL: Strict payment rate limit exceeded - possible fraud attempt")

			// Record metrics for monitoring
			if m.metrics != nil {
				m.metrics.RecordRateLimitExceeded(c.Path(), userID, c.IP())
			}

			return utils.ErrorResponse(c, fiber.StatusTooManyRequests, "payments.errors.suspiciousActivity")
		},

		Next: nil, // Не пропускаем никого, даже в dev режиме

		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}
