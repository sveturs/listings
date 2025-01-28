// backend/internal/proj/users/handler/auth.go
package handler

import (
 	"backend/pkg/utils"
    globalService "backend/internal/proj/global/service"
    "backend/internal/proj/users/service" 
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
    services globalService.ServicesInterface
	authService service.AuthServiceInterface
}

func NewAuthHandler(services globalService.ServicesInterface) *AuthHandler {
	return &AuthHandler{
		services:    services,
		authService: services.Auth(),
	}
}

func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	// Получаем returnTo из query параметров
	returnTo := c.Query("returnTo")
	if returnTo != "" {
		// Сохраняем в cookie
		c.Cookie(&fiber.Cookie{
			Name:     "returnTo",
			Value:    returnTo,
			Path:     "/",
			MaxAge:   300, // 5 минут
			Secure:   true,
			HTTPOnly: true,
		})
	}
	url := h.services.Auth().GetGoogleAuthURL()
	return c.Redirect(url)
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

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
		MaxAge:   3600 * 24,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})
	returnTo := h.services.Config().FrontendURL // значение по умолчанию
	if saved := c.Cookies("returnTo"); saved != "" {
		returnTo = h.services.Config().FrontendURL + saved
		// Удаляем cookie
		c.Cookie(&fiber.Cookie{
			Name:   "returnTo",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	// Устанавливаем cookie с токеном
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   3600 * 24,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Redirect(returnTo)
}

func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
    sessionToken := c.Cookies("session_token")
    if sessionToken == "" {
        return c.JSON(fiber.Map{
            "authenticated": false,
        })
    }

    // Исправляем вызов метода GetSession, добавляя контекст
    sessionData, err := h.services.Auth().GetSession(c.Context(), sessionToken)
    if err != nil {
        return c.JSON(fiber.Map{
            "authenticated": false,
        })
    }

    // Проверяем на nil вместо проверки ошибки
    if sessionData == nil {
        return c.JSON(fiber.Map{
            "authenticated": false,
        })
    }

    return c.JSON(fiber.Map{
        "authenticated": true,
        "user": fiber.Map{
            "id":          sessionData.UserID,
            "name":        sessionData.Name,
            "email":       sessionData.Email,
            "provider":    sessionData.Provider,
            "picture_url": sessionData.PictureURL,
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
