// backend/internal/middleware/auth.go
package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"backend/pkg/utils"
)

func (m *Middleware) AuthRequired(c *fiber.Ctx) error {
	log.Printf("AuthRequired middleware called for path: %s", c.Path())

	// Сначала проверяем JWT токен в заголовке Authorization
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		// Извлекаем токен
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString != "" {
			// Валидируем JWT токен
			claims, err := utils.ValidateJWTToken(tokenString, m.config.JWTSecret)
			if err == nil {
				log.Printf("AuthRequired: JWT authentication successful for user ID: %d", claims.UserID)

				// Сохраняем данные пользователя в контексте
				c.Locals("user_id", claims.UserID)
				c.Locals("user_email", claims.Email)
				c.Locals("auth_method", "jwt")

				// Обновляем last_seen асинхронно
				go func() {
					ctx := context.Background()
					_ = m.services.User().UpdateLastSeen(ctx, claims.UserID)
				}()

				return c.Next()
			}
			log.Printf("AuthRequired: JWT validation failed: %v", err)
		}
	}

	// Если JWT не удался или отсутствует, пробуем сессионную аутентификацию
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		log.Printf("AuthRequired: No authentication method found")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.token_required")
	}

	log.Printf("AuthRequired: Trying session authentication with token: %s", sessionToken[:5])

	session, err := m.services.Auth().GetSession(c.Context(), sessionToken)
	if err != nil {
		log.Printf("AuthRequired: Error getting session: %v", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_session")
	}

	if session == nil || session.UserID == 0 {
		log.Printf("AuthRequired: Session is nil or UserID is 0")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_session")
	}

	log.Printf("AuthRequired: Session authentication successful for user ID: %d", session.UserID)

	c.Locals("user_id", session.UserID)
	c.Locals("session_token", sessionToken)
	c.Locals("auth_method", "session")

	go func() {
		ctx := context.Background()
		_ = m.services.User().UpdateLastSeen(ctx, session.UserID)
	}()

	// Обновляем cookie при каждом запросе, чтобы сессия не истекала
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   3600 * 24, // 1 день
		Secure:   true,
		HTTPOnly: false, // Разрешаем доступ из JavaScript
		SameSite: "Lax",
	})

	return c.Next()
}
