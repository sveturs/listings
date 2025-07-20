package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const localeKey contextKey = "locale"

// LocaleMiddleware извлекает язык из запроса и добавляет его в контекст
func (m *Middleware) LocaleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Приоритет определения языка:
		// 1. Query параметр "lang" или "locale"
		// 2. Заголовок "Accept-Language"
		// 3. Cookie "locale"
		// 4. Язык по умолчанию "sr"

		var locale string

		// 1. Проверяем query параметры
		locale = c.Query("lang")
		if locale == "" {
			locale = c.Query("locale")
		}

		// 2. Проверяем заголовок Accept-Language
		if locale == "" {
			acceptLang := c.Get("Accept-Language")
			if acceptLang != "" {
				// Извлекаем первый язык из заголовка
				parts := strings.Split(acceptLang, ",")
				if len(parts) > 0 {
					lang := strings.TrimSpace(parts[0])
					// Убираем суффиксы типа en-US -> en
					if idx := strings.Index(lang, "-"); idx != -1 {
						lang = lang[:idx]
					}
					locale = lang
				}
			}
		}

		// 3. Проверяем cookie
		if locale == "" {
			locale = c.Cookies("locale")
		}

		// 4. Валидация и установка языка по умолчанию
		validLocales := map[string]bool{
			"sr": true,
			"en": true,
			"ru": true,
		}

		if !validLocales[locale] {
			locale = "sr" // По умолчанию сербский
		}

		// Добавляем язык в контекст
		ctx := context.WithValue(c.Context(), localeKey, locale)
		c.SetUserContext(ctx)

		// Добавляем заголовок Content-Language в ответ
		c.Set("Content-Language", locale)

		return c.Next()
	}
}
