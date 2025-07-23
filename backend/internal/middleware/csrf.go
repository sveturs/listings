package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	CSRFTokenLength = 32
	CSRFCookieName  = "csrf_token"
	CSRFHeaderName  = "X-CSRF-Token"
)

// generateCSRFToken генерирует криптографически стойкий CSRF токен
func generateCSRFToken() (string, error) {
	bytes := make([]byte, CSRFTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CSRFProtection возвращает middleware для защиты от CSRF атак
func (m *Middleware) CSRFProtection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Пропускаем GET, HEAD, OPTIONS запросы
		method := c.Method()
		if method == httpMethodGet || method == httpMethodHead || method == httpMethodOptions {
			return c.Next()
		}

		// Получаем токен из заголовка
		headerToken := c.Get(CSRFHeaderName)
		if headerToken == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token required",
			})
		}

		// Получаем токен из cookie
		cookieToken := c.Cookies(CSRFCookieName)
		if cookieToken == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token not found in cookies",
			})
		}

		// Сравниваем токены
		if headerToken != cookieToken {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token mismatch",
			})
		}

		return c.Next()
	}
}

// GetCSRFToken возвращает обработчик для получения CSRF токена
func (m *Middleware) GetCSRFToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Генерируем новый токен
		token, err := generateCSRFToken()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate CSRF token",
			})
		}

		// Устанавливаем cookie с токеном
		c.Cookie(&fiber.Cookie{
			Name:     CSRFCookieName,
			Value:    token,
			Path:     "/",
			Domain:   m.config.GetCookieDomain(),
			HTTPOnly: true,
			Secure:   m.config.Environment == "production",
			SameSite: "Lax",
			MaxAge:   int(24 * time.Hour.Seconds()), // 24 часа
		})

		// Возвращаем токен в JSON ответе
		return c.JSON(fiber.Map{
			"csrf_token": token,
		})
	}
}
