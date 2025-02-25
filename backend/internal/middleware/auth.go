// backend/internal/middleware/auth.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "backend/pkg/utils"
    //"context"
    //"log"
)

func (m *Middleware) AuthRequired(c *fiber.Ctx) error {
    sessionToken := c.Cookies("session_token")
    if sessionToken == "" {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
    }

    session, err := m.services.Auth().GetSession(c.Context(), sessionToken)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
    }

    if session == nil || session.UserID == 0 {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
    }

    c.Locals("user_id", session.UserID)
    c.Locals("session_token", sessionToken)
    
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