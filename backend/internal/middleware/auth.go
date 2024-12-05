package middleware

import (
    "github.com/gofiber/fiber/v2"
    "backend/pkg/utils"
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

    c.Locals("user", sessionData)
    return c.Next()
}