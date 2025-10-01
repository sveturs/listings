package handler

import (
	"fmt"
	"net/url"

	"backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
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
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid callback parameters"
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
			if err.Error() == "failed to exchange code" {
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
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
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

// GetUserRoles returns roles for a specific user
// @Summary Get user roles
// @Description Returns all roles assigned to a specific user
// @Tags roles
// @Security BearerAuth
// @Param userId path int true "User ID"
// @Produce json
// @Success 200 {object} map[string]interface{} "User roles"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "User not found"
// @Router /api/v1/users/{userId}/roles [get]
func (h *AuthHandler) GetUserRoles(c *fiber.Ctx) error {
	userID := c.Params("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID required",
		})
	}

	// For now, return empty roles
	// TODO: Implement through auth service
	return c.JSON(fiber.Map{
		"roles": []string{},
	})
}

// GetRoles returns all available roles
// @Summary Get all roles
// @Description Returns all available roles in the system
// @Tags roles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} []map[string]interface{} "List of roles"
// @Router /api/v1/roles [get]
func (h *AuthHandler) GetRoles(c *fiber.Ctx) error {
	// TODO: Implement through auth service
	return c.JSON(fiber.Map{
		"roles": []fiber.Map{
			{"id": 1, "name": "admin", "description": "Administrator"},
			{"id": 2, "name": "user", "description": "Regular user"},
		},
	})
}

// AssignRole assigns a role to a user
// @Summary Assign role to user
// @Description Assigns a specific role to a user
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Role assigned successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Router /api/v1/roles/assign [post]
func (h *AuthHandler) AssignRole(c *fiber.Ctx) error {
	// TODO: Implement through auth service when EntityAssignRoleRequest is available
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Role assignment not implemented",
	})
}

// RevokeRole revokes a role from a user
// @Summary Revoke role from user
// @Description Revokes a specific role from a user
// @Tags roles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag "Role revoked successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Invalid request"
// @Router /api/v1/roles/revoke [post]
func (h *AuthHandler) RevokeRole(c *fiber.Ctx) error {
	// TODO: Implement through auth service when EntityRevokeRoleRequest is available
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Role revocation not implemented",
	})
}
