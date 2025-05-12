// backend/internal/middleware/auth.go
package middleware

import (
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	//    "log"
	"context"
)

func (m *Middleware) AuthRequired(c *fiber.Ctx) error {
	//    log.Printf("AuthRequired middleware called for path: %s", c.Path())

	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		//        log.Printf("AuthRequired: No session_token cookie found")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	//    log.Printf("AuthRequired: Found session_token: %s (first 5 chars)", sessionToken[:5])

	session, err := m.services.Auth().GetSession(c.Context(), sessionToken)
	if err != nil {
		//        log.Printf("AuthRequired: Error getting session: %v", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	if session == nil || session.UserID == 0 {
		//        log.Printf("AuthRequired: Session is nil or UserID is 0")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	//    log.Printf("AuthRequired: Valid session found for user ID: %d", session.UserID)

	c.Locals("user_id", session.UserID)
	c.Locals("session_token", sessionToken)
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

	//    log.Printf("AuthRequired: Authentication successful for user ID: %d", session.UserID)
	return c.Next()
}
