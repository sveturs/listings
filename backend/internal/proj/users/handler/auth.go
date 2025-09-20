// backend/internal/proj/users/handler/auth.go
package handler

import (
	"crypto/rand"
	"math/big"

	"github.com/gofiber/fiber/v2"

	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	_ "backend/pkg/utils" // For swagger documentation
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
	// Check for explicit returnTo parameter first
	returnTo := c.Query("returnTo")

	if returnTo == "" {
		// Fallback to Origin header
		returnTo = c.Get("Origin")
	}
	if returnTo == "" {
		// Fallback to Referer header
		returnTo = c.Get("Referer")
	}
	if returnTo == "" {
		// Use configured frontend URL as last resort
		returnTo = h.cfg.FrontendURL
	}

	authURL := h.authService.GetGoogleAuthURL(returnTo)
	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
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
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "auth.google_callback.error.missing_code",
		})
	}

	state := c.Query("state")
	sessionData, err := h.authService.HandleGoogleCallback(c.Context(), code)
	if err != nil {
		logger.Error().
			Err(err).
			Str("code", code[:10]+"...").
			Msg("GoogleCallback: Failed to handle callback")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "auth.google_callback.error.authentication_failed",
		})
	}

	// Generate JWT token for OAuth user
	jwtToken, err := h.authService.GenerateJWT(sessionData.UserID, sessionData.Email)
	if err != nil {
		logger.Error().
			Err(err).
			Int("user_id", sessionData.UserID).
			Msg("GoogleCallback: Failed to generate JWT token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "auth.google_callback.error.token_generation_failed",
		})
	}

	// Generate session token and save session
	sessionToken := generateSessionToken()
	h.authService.SaveSession(sessionToken, sessionData)

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		MaxAge:   24 * 60 * 60, // 24 hours
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // Set to true in production
		SameSite: "Lax",
	})

	// Redirect to frontend with JWT token
	redirectURL := state
	if redirectURL == "" || redirectURL == "default" {
		redirectURL = h.cfg.FrontendURL // Use configured frontend URL
	}

	// Добавляем JWT токен в URL для передачи на frontend
	redirectURL = redirectURL + "?token=" + jwtToken

	return c.Redirect(redirectURL, fiber.StatusTemporaryRedirect)
}

// GetSession - Local implementation when Auth Service is disabled
// @Summary Get current session
// @Description Get current user session information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=SessionResponse} "Session information"
// @Router /auth/session [get]
func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
	// Check if user is authenticated via JWT middleware
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "auth.session.error.not_authenticated",
		})
	}

	// Get user profile from database
	userProfile, err := h.services.User().GetUserProfile(c.Context(), userID.(int))
	if err != nil {
		logger.Error().
			Err(err).
			Int("user_id", userID.(int)).
			Msg("GetSession: Failed to get user profile")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "auth.session.error.failed",
		})
	}

	// Return session information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"authenticated": true,
			"user": fiber.Map{
				"id":       userProfile.ID,
				"email":    userProfile.Email,
				"name":     userProfile.Name,
				"is_admin": userProfile.IsAdmin,
			},
		},
	})
}

// Logout - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Logout user
// @Description This endpoint is deprecated and should be proxied to Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Logout successful"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("Logout handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "Logout endpoint should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
	})
}

// Login - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Login with email and password
// @Description This endpoint is deprecated and should be proxied to Auth Service
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
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("Login handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "Login endpoint should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
	})
}

// Register - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Register new user
// @Description This endpoint is deprecated and should be proxied to Auth Service
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
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("Register handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "Register endpoint should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
	})
}

// RefreshToken - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Refresh access token
// @Description This endpoint is deprecated and should be proxied to Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=TokenResponse} "New access token"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.refresh_token.error.token_not_found or auth.refresh_token.error.invalid_token"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("RefreshToken handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "RefreshToken endpoint should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
	})
}

// generateSessionToken generates a secure random session token
func generateSessionToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 32)
	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a less random but still functional approach
			// This should never happen with crypto/rand
			continue
		}
		token[i] = charset[n.Int64()]
	}
	return string(token)
}
