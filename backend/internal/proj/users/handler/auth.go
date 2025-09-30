package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/sveturs/auth/pkg/http/entity"
	"github.com/sveturs/auth/pkg/http/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	oauthService *service.OAuthService
	backendURL   string
	frontendURL  string
	log          zerolog.Logger
}

func NewAuthHandler(
	authService *service.AuthService,
	oauthService *service.OAuthService,
	backendURL string,
	frontendURL string,
	log zerolog.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
		backendURL:   backendURL,
		frontendURL:  frontendURL,
		log:          log,
	}
}

// Register handles user registration via auth-service
// @Summary Register a new user
// @Description Registers a new user via auth-service
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Registration request"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 409 {object} utils.ErrorResponseSwag "User already exists"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Log that we reached the handler
	h.log.Debug().Msg("Register handler called")

	var req entity.UserRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to parse registration request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	h.log.Debug().Interface("parsed_request", req).Msg("Registration request parsed")

	resp, err := h.authService.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to auth service",
		})
	}

	c.Set("Content-Type", "application/json")
	switch resp.StatusCode() {
	case 201:
		return c.Status(201).JSON(resp.JSON201)
	case 400:
		return c.Status(400).JSON(resp.JSON400)
	case 409:
		return c.Status(409).JSON(resp.JSON409)
	case 500:
		return c.Status(500).JSON(resp.JSON500)
	default:
		return c.Status(resp.StatusCode()).Send(resp.Body)
	}
}

// Login handles user login via auth-service
// @Summary User login
// @Description Authenticates user via auth-service
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Login request"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 401 {object} utils.ErrorResponseSwag "Invalid credentials"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req entity.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Error().Err(err).Str("body", string(c.Body())).Msg("Failed to parse login request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Add default device info if not provided
	if req.DeviceID == "" {
		req.DeviceID = "web-browser"
	}
	if req.DeviceName == "" {
		req.DeviceName = "Web Browser"
	}

	resp, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to auth service",
		})
	}

	c.Set("Content-Type", "application/json")
	switch resp.StatusCode() {
	case 200:
		return c.Status(200).JSON(resp.JSON200)
	case 400:
		return c.Status(400).JSON(resp.JSON400)
	case 401:
		return c.Status(401).JSON(resp.JSON401)
	case 403:
		return c.Status(403).JSON(resp.JSON403)
	case 429:
		return c.Status(429).JSON(resp.JSON429)
	default:
		return c.Status(resp.StatusCode()).Send(resp.Body)
	}
}

// Logout handles user logout via auth-service
// @Summary User logout
// @Description Logs out the current user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Logout successful"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	resp, err := h.authService.Logout(c.Context(), authHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to auth service",
		})
	}

	c.Set("Content-Type", "application/json")
	if resp.JSON200 != nil {
		return c.Status(200).JSON(resp.JSON200)
	}
	return c.Status(resp.StatusCode()).Send(resp.Body)
}

// RefreshToken handles token refresh via auth-service
// @Summary Refresh access token
// @Description Refreshes the access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "Refresh token request"
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 401 {object} utils.ErrorResponseSwag "Invalid refresh token"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req entity.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	resp, err := h.authService.RefreshToken(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to auth service",
		})
	}

	c.Set("Content-Type", "application/json")
	switch resp.StatusCode() {
	case 200:
		return c.Status(200).JSON(resp.JSON200)
	case 400:
		return c.Status(400).JSON(resp.JSON400)
	case 401:
		return c.Status(401).JSON(resp.JSON401)
	default:
		return c.Status(resp.StatusCode()).Send(resp.Body)
	}
}

// Validate validates the current token
// @Summary Validate token
// @Description Validates the current access token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/auth/validate [get]
func (h *AuthHandler) Validate(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	resp, err := h.authService.ValidateToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to connect to auth service",
		})
	}

	if resp.Valid {
		return c.JSON(fiber.Map{
			"valid":   resp.Valid,
			"user_id": resp.UserID,
			"email":   resp.Email,
			"roles":   resp.Roles,
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "invalid token",
	})
}

// GetCurrentUser returns current authenticated user info
// @Summary Get current user info
// @Description Returns info about the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "User info"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	email := c.Locals("email")
	roles := c.Locals("roles")

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    userID,
			"email": email,
			"roles": roles,
		},
	})
}
