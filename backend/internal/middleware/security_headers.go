package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders добавляет заголовки безопасности к ответам
func (m *Middleware) SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Защита от Clickjacking
		c.Set("X-Frame-Options", "DENY")

		// Предотвращение MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Включение встроенной защиты XSS в браузерах
		c.Set("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		// Разрешаем загрузку ресурсов только с нашего домена
		// unsafe-inline для стилей нужен для DaisyUI
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://accounts.google.com https://apis.google.com; "+
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data: https: blob:; "+
				"connect-src 'self' wss: https://accounts.google.com https://www.googleapis.com; "+
				"frame-src https://accounts.google.com;")

		// Referrer Policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy (ранее Feature Policy)
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// HSTS будет добавлен только для production через конфигурацию
		if m.services.Config().IsProduction() {
			// Strict Transport Security
			// max-age=31536000 (1 год)
			// includeSubDomains - применяется ко всем поддоменам
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		return c.Next()
	}
}
