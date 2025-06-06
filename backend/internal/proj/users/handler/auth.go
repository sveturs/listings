// backend/internal/proj/users/handler/auth.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	"backend/internal/types"
	"backend/pkg/jwt"
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

// GoogleAuth redirects to Google OAuth2 authorization page
// @Summary Initiate Google OAuth authentication
// @Description Redirects user to Google OAuth2 authorization page
// @Tags auth
// @Accept json
// @Produce json
// @Param returnTo query string false "URL to return after authentication"
// @Success 302 {string} string "Redirect to Google OAuth"
// @Router /auth/google [get]
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

// GoogleCallback handles Google OAuth2 callback
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth2 callback from Google and creates user session
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Success 302 {string} string "Redirect to frontend with session"
// @Failure 500 {object} ErrorResponse "Authentication failed"
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	sessionData, err := h.services.Auth().HandleGoogleCallback(c.Context(), code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "auth.google_callback.error.authentication_failed")
	}

	// Генерация токена сессии
	sessionToken := utils.GenerateSessionToken()
	h.services.Auth().SaveSession(sessionToken, sessionData)

	// Генерация JWT и Refresh токенов
	accessToken, refreshToken, err := h.services.Auth().GenerateTokensForOAuth(c.Context(), sessionData.UserID, sessionData.Email, c.IP(), c.Get("User-Agent"))
	if err != nil {
		h.logger.Error("Failed to generate tokens: %v", err)
		// Продолжаем с session token как fallback
	}

	// Установка session cookie для обратной совместимости
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   3600 * 24,
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	// Установка refresh token в httpOnly cookie
	if refreshToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			MaxAge:   30 * 24 * 3600, // 30 дней
			Secure:   h.services.Config().GetCookieSecure(),
			HTTPOnly: true,
			SameSite: h.services.Config().GetCookieSameSite(),
		})
		h.logger.Info("OAuth: Set refresh_token cookie for user %s, access token generated: %v", sessionData.Email, accessToken != "")
	}
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

	return c.Redirect(returnTo)
}

// GetSession returns current user session information
// @Summary Get current session
// @Description Returns information about the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SessionResponse "Session information"
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
		h.logger.Info("GetSession: Authorization header = '%s'", authHeader)
		
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			jwtToken := authHeader[7:]
			h.logger.Info("GetSession: Extracted JWT token (first 20 chars): %s...", jwtToken[:20])

			// Проверяем JWT токен
			claims, err := utils.ValidateJWTToken(jwtToken, h.services.Config().JWTSecret)
			if err != nil {
				h.logger.Info("JWT token validation failed: %v (IP: %s, UserAgent: %s)",
					err, c.IP(), c.Get("User-Agent"))
			} else if claims != nil {
				h.logger.Info("JWT claims validated: UserID=%d, Email=%s", claims.UserID, claims.Email)
				
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
		} else {
			h.logger.Info("GetSession: No valid Authorization header found")
		}
	}

	// Если ни session cookie, ни JWT не дали результата
	if sessionData == nil {
		response := SessionResponse{
			Authenticated: false,
		}
		return c.JSON(response)
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
	return c.JSON(response)
}

// Logout terminates user session and revokes tokens
// @Summary Logout user
// @Description Logs out the user by revoking all refresh tokens and clearing cookies
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {string} string "OK"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Получаем userID из контекста (может быть из JWT или session)
	userID := 0
	
	// Пробуем получить из JWT
	if user, ok := c.Locals("user").(*jwt.Claims); ok && user != nil {
		userID = user.UserID
		h.logger.Info("Logout: Got userID from JWT: %d", userID)
	}
	
	// Если не получили из JWT, пробуем из session
	if userID == 0 {
		sessionToken := c.Cookies("session_token")
		if sessionToken != "" {
			sessionData, _ := h.services.Auth().GetSession(c.Context(), sessionToken)
			if sessionData != nil {
				userID = sessionData.UserID
				h.logger.Info("Logout: Got userID from session: %d", userID)
			}
		}
	}
	
	// Отзываем ВСЕ refresh токены пользователя
	if userID > 0 {
		h.logger.Info("Logout: Revoking ALL refresh tokens for userID: %d", userID)
		if err := h.services.Auth().RevokeUserRefreshTokens(c.Context(), userID); err != nil {
			h.logger.Error("Failed to revoke user refresh tokens: %v", err)
		} else {
			h.logger.Info("All user refresh tokens revoked successfully")
		}
	} else {
		// Fallback: отзываем только текущий токен
		refreshToken := c.Cookies("refresh_token")
		if refreshToken != "" {
			h.logger.Info("Logout: No userID, revoking single token: %s...", refreshToken[:20])
			if err := h.services.Auth().RevokeRefreshToken(c.Context(), refreshToken); err != nil {
				h.logger.Error("Failed to revoke refresh token: %v", err)
			}
		}
	}

	// Удаляем refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
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
		MaxAge:   -1,
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})
	
	return c.SendStatus(fiber.StatusOK)
}

// Login authenticates user with email and password
// @Summary Login with email and password
// @Description Authenticates user and returns JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse "Authentication successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
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
		if err == types.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.login.error.invalid_credentials")
		}
		h.logger.Error("Login failed: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "auth.login.error.failed")
	}

	// Устанавливаем refresh token в httpOnly cookie
	cookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/", // Используем корневой путь для доступности на всех роутах
		MaxAge:   30 * 24 * 3600, // 30 дней
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	}
	
	// В development не устанавливаем Domain, чтобы cookie работала на localhost
	if !h.services.Config().IsDevelopment() {
		cookie.Domain = ".svetu.rs" // Для production
	}
	
	c.Cookie(cookie)
	h.logger.Info("Setting refresh_token cookie for user %s, cookie config: Path=%s, Secure=%v, SameSite=%s", 
		user.Email, cookie.Path, cookie.Secure, cookie.SameSite)

	// Логируем успешную выдачу токенов
	tokenPrefix := accessToken
	if len(accessToken) > 20 {
		tokenPrefix = accessToken[:20] + "..."
	}
	h.logger.Info("Tokens issued for login - UserID: %d, Email: %s, AccessToken: %s, ExpiresIn: %d",
		user.ID, user.Email, tokenPrefix, h.services.Config().JWTExpirationHours * 3600)

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
	return c.JSON(response)
}

// Register creates a new user account
// @Summary Register new user
// @Description Creates a new user account and returns JWT access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Registration data"
// @Success 200 {object} AuthResponse "Registration successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 409 {object} ErrorResponse "User already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
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
		if err == types.ErrUserAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "auth.register.error.email_exists")
		}
		h.logger.Error("Registration failed: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "auth.register.error.failed")
	}

	// Устанавливаем refresh token в httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/", // Используем корневой путь для доступности на всех роутах
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
	h.logger.Info("New user registered and tokens issued - UserID: %d, Email: %s, AccessToken: %s, ExpiresIn: %d",
		user.ID, user.Email, tokenPrefix, h.services.Config().JWTExpirationHours * 3600)

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
	return c.JSON(response)
}

// RefreshToken refreshes JWT access token using refresh token
// @Summary Refresh access token
// @Description Refreshes expired access token using valid refresh token from cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} TokenResponse "New access token"
// @Failure 401 {object} ErrorResponse "Invalid or missing refresh token"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Получаем refresh токен из cookie
	refreshToken := c.Cookies("refresh_token")
	h.logger.Info("RefreshToken called. Token present: %v, Token length: %d, Cookie header: %s", 
		refreshToken != "", len(refreshToken), c.Get("Cookie"))
	
	if refreshToken == "" {
		h.logger.Error("Refresh token not found in cookie")
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
			MaxAge:   -1,
			Secure:   h.services.Config().GetCookieSecure(),
			HTTPOnly: true,
			SameSite: h.services.Config().GetCookieSameSite(),
		})
		
		h.logger.Error("Token refresh failed: %v", err)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.refresh_token.error.invalid_token")
	}

	// Устанавливаем новый refresh токен в cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 дней
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	// Логируем успешное обновление токенов
	h.logger.Info("Tokens refreshed successfully for refresh_token")
	
	// Возвращаем новый access токен в формате, ожидаемом фронтендом
	response := TokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   h.services.Config().JWTExpirationHours * 3600,
	}
	
	h.logger.Info("Refresh endpoint returning: %+v", response)
	return c.JSON(response)
}
