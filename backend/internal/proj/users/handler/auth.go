// backend/internal/proj/users/handler/auth.go
package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	"backend/internal/types"
	"backend/pkg/jwt"
	"backend/pkg/utils"
)

type AuthHandler struct {
	services    globalService.ServicesInterface
	authService service.AuthServiceInterface
}

func NewAuthHandler(services globalService.ServicesInterface) *AuthHandler {
	return &AuthHandler{
		services:    services,
		authService: services.Auth(),
	}
}

// GoogleAuth - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Initiate Google OAuth authentication
// @Description This endpoint is deprecated and should be proxied to Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Param returnTo query string false "URL to return after authentication"
// @Success 302 {string} string "Redirect to Google OAuth"
// @Router /auth/google [get]
func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	// This handler should never be called if Auth Service proxy is configured correctly
	logger.Error().Msg("GoogleAuth handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "OAuth authentication should be handled by Auth Service",
		"message": "Please ensure USE_AUTH_SERVICE=true and Auth Service is running",
	})
}

// GoogleCallback - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Handle Google OAuth callback
// @Description This endpoint is deprecated and should be proxied to Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 302 {string} string "Redirect to frontend with session"
// @Failure 500 {object} utils.ErrorResponseSwag "auth.google_callback.error.authentication_failed"
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// This handler should never be called if Auth Service proxy is configured correctly
	logger.Error().Msg("GoogleCallback handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "OAuth callback should be handled by Auth Service",
		"message": "Please ensure USE_AUTH_SERVICE=true and Auth Service is running",
	})
}

// GetSession returns current user session information
// @Summary Get current session
// @Description Returns information about the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=SessionResponse} "Session information"
// @Router /auth/session [get]
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
		logger.Info().Str("auth_header", authHeader).Msg("GetSession: Authorization header")

		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			jwtToken := authHeader[7:]
			if len(jwtToken) > 20 {
				logger.Info().
					Str("token_prefix", jwtToken[:20]+"...").
					Msg("GetSession: Extracted JWT token")
			}

			// Проверяем JWT токен
			claims, validateErr := utils.ValidateJWTToken(jwtToken, h.services.Config().JWTSecret)
			if validateErr != nil {
				logger.Info().
					Err(validateErr).
					Str("ip", c.IP()).
					Str("user_agent", c.Get("User-Agent")).
					Msg("JWT token validation failed")
			} else if claims != nil {
				logger.Info().
					Int("user_id", claims.UserID).
					Str("email", claims.Email).
					Msg("JWT claims validated")

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
					logger.Info().
						Str("email", claims.Email).
						Int("user_id", claims.UserID).
						Msg("JWT session restored for user")
				} else {
					logger.Error().
						Err(err).
						Str("email", claims.Email).
						Msg("Failed to get user data for JWT claims")
				}
			}
		} else {
			logger.Info().Msg("GetSession: No valid Authorization header found")
		}
	}

	// Если ни session cookie, ни JWT не дали результата
	if sessionData == nil {
		response := SessionResponse{
			Authenticated: false,
		}
		return utils.SuccessResponse(c, response)
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

	response := SessionResponse{
		Authenticated: true,
		User: &SessionUserResponse{
			ID:         sessionData.UserID,
			Name:       sessionData.Name,
			Email:      sessionData.Email,
			Provider:   sessionData.Provider,
			PictureURL: sessionData.PictureURL,
			IsAdmin:    isAdmin,
			City:       city,
			Country:    country,
			Phone:      phone,
		},
	}
	return utils.SuccessResponse(c, response)
}

// Logout terminates user session and revokes tokens
// @Summary Logout user
// @Description Logs out the user by revoking all refresh tokens and clearing cookies
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Logout successful"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Получаем userID из контекста (может быть из JWT или session)
	userID := 0

	// Пробуем получить из JWT
	if user, ok := c.Locals("user").(*jwt.Claims); ok && user != nil {
		userID = user.UserID
		logger.Info().Int("user_id", userID).Msg("Logout: Got userID from JWT")
	}

	// Если не получили из JWT, пробуем из session
	if userID == 0 {
		sessionToken := c.Cookies("session_token")
		if sessionToken != "" {
			sessionData, _ := h.services.Auth().GetSession(c.Context(), sessionToken)
			if sessionData != nil {
				userID = sessionData.UserID
				logger.Info().Int("user_id", userID).Msg("Logout: Got userID from session")
			}
		}
	}

	// Отзываем ВСЕ refresh токены пользователя
	if userID > 0 {
		logger.Info().Int("user_id", userID).Msg("Logout: Revoking ALL refresh tokens for userID")
		if err := h.services.Auth().RevokeUserRefreshTokens(c.Context(), userID); err != nil {
			logger.Error().Err(err).Msg("Failed to revoke user refresh tokens")
		} else {
			logger.Info().Msg("All user refresh tokens revoked successfully")
		}
	} else {
		// Fallback: отзываем только текущий токен
		refreshToken := c.Cookies("refresh_token")
		if refreshToken != "" {
			if len(refreshToken) > 20 {
				logger.Info().
					Str("token_prefix", refreshToken[:20]+"...").
					Msg("Logout: No userID, revoking single token")
			}
			if err := h.services.Auth().RevokeRefreshToken(c.Context(), refreshToken); err != nil {
				logger.Error().Err(err).Msg("Failed to revoke refresh token")
			}
		}
	}

	// Удаляем refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   h.services.Config().GetCookieDomain(),
		MaxAge:   -1,
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	// Удаляем session token (для обратной совместимости)
	sessionToken := c.Cookies("session_token")
	if sessionToken != "" {
		h.services.Auth().DeleteSession(sessionToken)
		c.Cookie(&fiber.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			Domain:   h.services.Config().GetCookieDomain(),
			MaxAge:   -1,
			Secure:   h.services.Config().GetCookieSecure(),
			HTTPOnly: true,
			SameSite: h.services.Config().GetCookieSameSite(),
		})
	}

	// Удаляем JWT token cookie (для обратной совместимости)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		Domain:   h.services.Config().GetCookieDomain(),
		MaxAge:   -1,
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	response := MessageResponse{
		Message: "auth.logout.success",
	}
	return utils.SuccessResponse(c, response)
}

// Login authenticates user with email and password
// @Summary Login with email and password
// @Description Authenticates user and returns JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} utils.SuccessResponseSwag{data=AuthResponse} "Authentication successful"
// @Failure 400 {object} utils.ErrorResponseSwag "auth.login.error.invalid_request_body or auth.login.error.email_password_required"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.login.error.invalid_credentials"
// @Failure 500 {object} utils.ErrorResponseSwag "auth.login.error.failed"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginData LoginRequest

	if err := c.BodyParser(&loginData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "auth.login.error.invalid_request_body")
	}

	// Валидация данных
	if loginData.Email == "" || loginData.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "auth.login.error.email_password_required")
	}

	// Аутентификация с получением access и refresh токенов
	accessToken, refreshToken, user, err := h.services.Auth().LoginWithRefreshToken(c.Context(), loginData.Email, loginData.Password, c.IP(), c.Get("User-Agent"))
	if err != nil {
		if errors.Is(err, types.ErrInvalidCredentials) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.login.error.invalid_credentials")
		}
		logger.Error().Err(err).Msg("Login failed")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "auth.login.error.failed")
	}

	// Устанавливаем refresh token в httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/", // Используем корневой путь для доступности на всех роутах
		Domain:   h.services.Config().GetCookieDomain(),
		MaxAge:   30 * 24 * 3600, // 30 дней
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})
	logger.Info().
		Str("email", user.Email).
		Str("domain", h.services.Config().GetCookieDomain()).
		Bool("secure", h.services.Config().GetCookieSecure()).
		Str("same_site", h.services.Config().GetCookieSameSite()).
		Msg("Setting refresh_token cookie for user")

	// Логируем успешную выдачу токенов
	tokenPrefix := accessToken
	if len(accessToken) > 20 {
		tokenPrefix = accessToken[:20] + "..."
	}
	logger.Info().
		Int("user_id", user.ID).
		Str("email", user.Email).
		Str("access_token_prefix", tokenPrefix).
		Int("expires_in", h.services.Config().JWTExpirationHours*3600).
		Msg("Tokens issued for login")

	// Возвращаем access токен и данные пользователя в формате, ожидаемом фронтендом
	response := AuthResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   h.services.Config().JWTExpirationHours * 3600,
		User: UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			PictureURL: user.PictureURL,
		},
	}
	return utils.SuccessResponse(c, response)
}

// Register creates a new user account
// @Summary Register new user
// @Description Creates a new user account and returns JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Registration data"
// @Success 200 {object} utils.SuccessResponseSwag{data=AuthResponse} "Registration successful"
// @Failure 400 {object} utils.ErrorResponseSwag "auth.register.error.invalid_request_body or auth.register.error.fields_required"
// @Failure 409 {object} utils.ErrorResponseSwag "auth.register.error.email_exists"
// @Failure 500 {object} utils.ErrorResponseSwag "auth.register.error.failed"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var registerData RegisterRequest

	if err := c.BodyParser(&registerData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "auth.register.error.invalid_request_body")
	}

	// Валидация данных
	if registerData.Name == "" || registerData.Email == "" || registerData.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "auth.register.error.fields_required")
	}

	// Регистрация с получением access и refresh токенов
	accessToken, refreshToken, user, err := h.services.Auth().RegisterWithRefreshToken(
		c.Context(),
		registerData.Name,
		registerData.Email,
		registerData.Password,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		if errors.Is(err, types.ErrUserAlreadyExists) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "auth.register.error.email_exists")
		}
		logger.Error().Err(err).Msg("Registration failed")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "auth.register.error.failed")
	}

	// Устанавливаем refresh token в httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/", // Используем корневой путь для доступности на всех роутах
		Domain:   h.services.Config().GetCookieDomain(),
		MaxAge:   30 * 24 * 3600, // 30 дней
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	// Логируем успешную регистрацию и выдачу токенов
	tokenPrefix := accessToken
	if len(accessToken) > 20 {
		tokenPrefix = accessToken[:20] + "..."
	}
	logger.Info().
		Int("user_id", user.ID).
		Str("email", user.Email).
		Str("access_token_prefix", tokenPrefix).
		Int("expires_in", h.services.Config().JWTExpirationHours*3600).
		Msg("New user registered and tokens issued")

	// Возвращаем access токен и данные пользователя в формате, ожидаемом фронтендом
	response := AuthResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   h.services.Config().JWTExpirationHours * 3600,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}
	return utils.SuccessResponse(c, response)
}

// RefreshToken refreshes JWT access token using refresh token
// @Summary Refresh access token
// @Description Refreshes expired access token using valid refresh token from cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=TokenResponse} "New access token"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.refresh_token.error.token_not_found or auth.refresh_token.error.invalid_token"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Получаем refresh токен из cookie
	refreshToken := c.Cookies("refresh_token")
	logger.Info().
		Bool("token_present", refreshToken != "").
		Int("token_length", len(refreshToken)).
		Str("cookie_header", c.Get("Cookie")).
		Str("host", c.Hostname()).
		Str("origin", c.Get("Origin")).
		Str("referer", c.Get("Referer")).
		Msg("RefreshToken called")

	if refreshToken == "" {
		logger.Error().
			Str("all_cookies", c.Get("Cookie")).
			Msg("Refresh token not found in cookie")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.refresh_token.error.token_not_found")
	}

	// Обновляем токены
	newAccessToken, newRefreshToken, err := h.services.Auth().RefreshTokens(
		c.Context(),
		refreshToken,
		c.IP(),
		c.Get("User-Agent"),
	)
	if err != nil {
		// Удаляем невалидный refresh токен из cookie
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			Domain:   h.services.Config().GetCookieDomain(),
			MaxAge:   -1,
			Secure:   h.services.Config().GetCookieSecure(),
			HTTPOnly: true,
			SameSite: h.services.Config().GetCookieSameSite(),
		})

		logger.Error().Err(err).Msg("Token refresh failed")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.refresh_token.error.invalid_token")
	}

	// Устанавливаем новый refresh токен в cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		Domain:   h.services.Config().GetCookieDomain(),
		MaxAge:   30 * 24 * 3600, // 30 дней
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	// Логируем успешное обновление токенов
	logger.Info().Msg("Tokens refreshed successfully for refresh_token")

	// Возвращаем новый access токен в формате, ожидаемом фронтендом
	response := TokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   h.services.Config().JWTExpirationHours * 3600,
	}

	logger.Info().
		Interface("response", response).
		Msg("Refresh endpoint returning")
	return utils.SuccessResponse(c, response)
}
