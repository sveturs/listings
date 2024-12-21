package middleware

import (
    "github.com/gofiber/fiber/v2"
    "backend/pkg/utils"
    "log"
)

func (m *Middleware) AuthRequired(c *fiber.Ctx) error {
    sessionToken := c.Cookies("session_token")
    if sessionToken == "" {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required")
    }

    sessionData, ok := m.services.Auth().GetSession(sessionToken)
    if !ok {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired session")
    }

    log.Printf("AuthRequired: setting userID=%d for session token %s", sessionData.UserID, sessionToken)

    // Устанавливаем user_id в контекст
    c.Locals("user_id", sessionData.UserID)
    c.Locals("user", sessionData)

    return c.Next()
}