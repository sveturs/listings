// backend/internal/proj/users/handler/auth.go
package handler

import (
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
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("GoogleAuth handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "OAuth authentication should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
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
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("GoogleCallback handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "OAuth callback should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
	})
}

// GetSession - DEPRECATED: This endpoint should be handled by Auth Service
// @Summary Get current session
// @Description This endpoint is deprecated and should be proxied to Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=SessionResponse} "Session information"
// @Router /auth/session [get]
func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
	// This handler should never be called - all auth requests are proxied to Auth Service
	logger.Error().Msg("GetSession handler called in monolith - should be proxied to Auth Service")
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error":   "Session endpoint should be handled by Auth Service",
		"message": "Auth Service proxy is not configured correctly",
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
