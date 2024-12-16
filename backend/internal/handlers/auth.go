package handlers

import (
    "github.com/gofiber/fiber/v2"
    "backend/internal/services"
    "backend/pkg/utils"
)

type AuthHandler struct {
    services services.ServicesInterface
}

func NewAuthHandler(services services.ServicesInterface) *AuthHandler {
    return &AuthHandler{
        services: services,
    }
}

func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
    url := h.services.Auth().GetGoogleAuthURL()
    return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing code")
	}

	sessionData, err := h.services.Auth().HandleGoogleCallback(c.Context(), code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Authentication failed")
	}

	// Генерация токена сессии
	sessionToken := utils.GenerateSessionToken()
	h.services.Auth().SaveSession(sessionToken, sessionData)

	// Установка cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   3600 * 24, // 24 часа
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

    return c.Redirect(h.services.Config().FrontendURL)
}

func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
    sessionToken := c.Cookies("session_token")
    if sessionToken == "" {
        return c.JSON(fiber.Map{
            "authenticated": false,
        })
    }

    sessionData, ok := h.services.Auth().GetSession(sessionToken)
    if !ok {
        return c.JSON(fiber.Map{
            "authenticated": false,
        })
    }

    return c.JSON(fiber.Map{
        "authenticated": true,
        "user": fiber.Map{
            "id": sessionData.UserID,      // Добавляем ID пользователя
            "name": sessionData.Name,
            "email": sessionData.Email,
            "provider": sessionData.Provider,
            "picture_url": sessionData.PictureURL, // Добавляем URL изображения
        },
    })
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken != "" {
		h.services.Auth().DeleteSession(sessionToken)
		c.Cookie(&fiber.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Lax",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
