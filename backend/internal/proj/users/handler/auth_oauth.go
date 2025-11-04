package handler

import (
	"errors"
	"fmt"
	"net/url"

	"backend/internal/domain"
	"backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	_ "backend/pkg/utils" // For Swagger documentation
)

// GoogleAuth redirects to Google OAuth
// @Summary Start Google OAuth authentication
// @Description Redirects user to Google OAuth consent page
// @Tags auth
// @Param locale query string false "User locale (en, ru, sr)"
// @Param return_url query string false "URL to return after auth"
// @Produce json
// @Success 302 "Redirect to Google OAuth"
// @Router /api/v1/auth/google [get]
func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	// Get user context from query parameters
	locale := c.Query("locale", "en")
	returnPath := c.Query("return_url", "/")

	logger.Info().
		Str("locale", locale).
		Str("return_path", returnPath).
		Msg("Starting OAuth with user context")

	redirectURI := fmt.Sprintf("%s/api/v1/auth/google/callback", h.backendURL)

	logger.Info().
		Str("backend_url", h.backendURL).
		Str("redirect_uri", redirectURI).
		Msg("OAuth redirectURI configuration")

	// Use the new method with locale and returnPath
	oauthURL, err := h.oauthService.StartGoogleOAuth(c.Context(), redirectURI, locale, returnPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to start Google OAuth")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to start OAuth flow",
		})
	}

	logger.Info().Str("url", oauthURL).Msg("Redirecting to Google OAuth")
	return c.Redirect(oauthURL)
}

// GoogleCallback handles Google OAuth callback
// @Summary Handle Google OAuth callback
// @Description Processes OAuth callback from Google and authenticates user
// @Tags auth
// @Param code query string true "OAuth authorization code"
// @Param state query string true "OAuth state parameter"
// @Produce json
// @Success 302 "Redirect to frontend with auth tokens"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid callback parameters"
// @Router /api/v1/auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	logger.Info().
		Bool("has_code", code != "").
		Str("state", state).
		Msg("OAuth callback received")

	// Complete OAuth flow - now result contains Locale and ReturnPath
	result, err := h.oauthService.CompleteGoogleOAuth(c.Context(), code, state)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to complete OAuth flow")

		// Determine error type for redirect
		errorType := "auth_failed"
		switch err.Error() {
		case "missing parameters":
			errorType = "missing_parameters"
		case "invalid state: invalid state":
			errorType = "invalid_state"
		default:
			if errors.Is(err, domain.ErrOAuthCodeExchangeFailed) {
				errorType = "exchange_failed"
			}
		}

		// Use default locale if not available
		locale := "en"
		if result != nil && result.Locale != "" {
			locale = result.Locale
		}

		// Redirect with locale to correct error page
		return c.Redirect(fmt.Sprintf("%s/%s/login?error=%s", h.frontendURL, locale, errorType))
	}

	logger.Info().
		Str("email", result.Email).
		Str("locale", result.Locale).
		Str("return_path", result.ReturnPath).
		Msg("OAuth successful with context")

	// Get locale and return path from result
	locale := result.Locale
	if locale == "" {
		locale = "en"
	}
	returnPath := result.ReturnPath
	if returnPath == "" {
		returnPath = "/"
	}

	// URL encode all parameters for redirect
	encodedReturnPath := url.QueryEscape(returnPath)
	encodedAccessToken := url.QueryEscape(result.AccessToken)
	encodedRefreshToken := url.QueryEscape(result.RefreshToken)

	// Redirect to Next.js API route which will set cookies
	// Next.js API route will then redirect to frontend callback page
	redirectURL := fmt.Sprintf("%s/api/auth/google/callback?access_token=%s&refresh_token=%s&locale=%s&return_url=%s",
		h.frontendURL, encodedAccessToken, encodedRefreshToken, locale, encodedReturnPath)

	logger.Info().
		Str("redirect_url_without_tokens", fmt.Sprintf("%s/api/auth/google/callback?locale=%s&return_url=%s", h.frontendURL, locale, encodedReturnPath)).
		Str("frontend_url", h.frontendURL).
		Str("locale", locale).
		Msg("Redirecting to Next.js API route to set cookies")

	return c.Redirect(redirectURL)
}

// GetSession returns current session info
// @Summary Get current session
// @Description Returns information about the current authenticated session
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Session info"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/auth/session [get]
func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	email, _ := authMiddleware.GetEmail(c)
	roles, _ := authMiddleware.GetRoles(c)

	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    userID,
			"email": email,
			"roles": roles,
		},
		"authenticated": true,
	})
}
