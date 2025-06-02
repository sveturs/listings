// backend/internal/proj/users/handler/auth.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	"backend/internal/types"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type AuthHandler struct {
	services    globalService.ServicesInterface
	authService service.AuthServiceInterface
	logger      *logger.Logger
}

func NewAuthHandler(services globalService.ServicesInterface) *AuthHandler {
	return &AuthHandler{
		services:    services,
		authService: services.Auth(),
		logger:      logger.New(),
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

// backend/internal/proj/users/handler/auth.go
func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
	var sessionData *types.SessionData
	var err error

	// Попробуем получить сессию через session cookie
	sessionToken := c.Cookies("session_token")
	if sessionToken != "" {
		sessionData, err = h.services.Auth().GetSession(c.Context(), sessionToken)
		if err != nil || sessionData == nil {
			sessionData = nil // Очищаем если ошибка
		}
	}

	// Если сессия через cookie не найдена, проверяем JWT токен
	if sessionData == nil {
		authHeader := c.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			jwtToken := authHeader[7:]

			// Проверяем JWT токен
			claims, err := utils.ValidateJWTToken(jwtToken, h.services.Config().JWTSecret)
			if err != nil {
				h.logger.Info("JWT token validation failed: %v (IP: %s, UserAgent: %s)", 
					err, c.IP(), c.Get("User-Agent"))
			} else if claims != nil {
				// Получаем данные пользователя из базы
				user, err := h.services.User().GetUserByEmail(c.Context(), claims.Email)
				if err == nil && user != nil {
					// Создаем временную sessionData для JWT пользователя
					sessionData = &types.SessionData{
						UserID:     claims.UserID,
						Email:      claims.Email,
						Name:       user.Name,
						Provider:   "password", // Предполагаем что JWT = password login
						PictureURL: user.PictureURL,
					}
					h.logger.Info("JWT session restored for user: %s (UserID: %d)", 
						claims.Email, claims.UserID)
				} else {
					h.logger.Error("Failed to get user data for JWT claims: %v (Email: %s)", 
						err, claims.Email)
				}
			}
		}
	}

	// Если ни session cookie, ни JWT не дали результата
	if sessionData == nil {
		return c.JSON(fiber.Map{
			"authenticated": false,
		})
	}

	// Проверяем, является ли пользователь администратором
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), sessionData.Email)
	if err != nil {
		isAdmin = false // По умолчанию не администратор в случае ошибки
	}

	// Проверяем, есть ли у пользователя дополнительная информация
	userProfile, err := h.services.User().GetUserProfile(c.Context(), sessionData.UserID)
	var city, country, phone string
	if err == nil && userProfile != nil {
		city = userProfile.City
		country = userProfile.Country
		if userProfile.Phone != nil {
			phone = *userProfile.Phone
		}
	}

	return c.JSON(fiber.Map{
		"authenticated": true,
		"user": fiber.Map{
			"id":          sessionData.UserID,
			"name":        sessionData.Name,
			"email":       sessionData.Email,
			"provider":    sessionData.Provider,
			"picture_url": sessionData.PictureURL,
			"is_admin":    isAdmin, // Добавляем поле is_admin!
			"city":        city,    // Добавляем поля профиля
			"country":     country,
			"phone":       phone,
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
