// backend/internal/middleware/auth.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "backend/pkg/utils"
    "log"
)

func (m *Middleware) AuthRequired(c *fiber.Ctx) error {
    log.Printf("Starting AuthRequired middleware")
    sessionToken := c.Cookies("session_token")
    if sessionToken == "" {
        log.Printf("No session token found")
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
    }

    log.Printf("Found session token: %s", sessionToken)
    session, err := m.services.Auth().GetSession(c.Context(), sessionToken)
    if err != nil {
        log.Printf("Error getting session: %v", err)
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
    }

    if session == nil || session.UserID == 0 {
        log.Printf("Invalid session: %+v", session)
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
    }

    log.Printf("Setting userID=%d for session token %s", session.UserID, sessionToken)
    c.Locals("user_id", session.UserID)
    return c.Next()
}