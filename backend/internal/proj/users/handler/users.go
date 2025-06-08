package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	"backend/internal/types"
	"backend/pkg/logger"
	"backend/pkg/utils"
)

type UserHandler struct {
	services    globalService.ServicesInterface
	userService service.UserServiceInterface
	logger      *logger.Logger
}

func NewUserHandler(services globalService.ServicesInterface) *UserHandler {
	return &UserHandler{
		services:    services,
		userService: services.User(),
		logger:      logger.New(),
	}
}

// Response structures for Swagger documentation

// UserProfileResponse представляет ответ с профилем пользователя
type UserProfileResponse struct {
	Success bool                `json:"success" example:"true"`
	Data    *models.UserProfile `json:"data"`
}

// MessageResponse представляет ответ с сообщением
type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Операция выполнена успешно"`
}

// RegisterResponse представляет ответ после регистрации
type RegisterResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Пользователь успешно зарегистрирован"`
	User    *models.User `json:"user"`
}

// PublicUserResponse представляет публичные данные пользователя
type PublicUserResponse struct {
	ID         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"Иван Иванов"`
	Email      string `json:"email" example:"user@example.com"`
	PictureURL string `json:"picture_url" example:"https://example.com/avatar.jpg"`
	CreatedAt  string `json:"created_at" example:"2023-01-01T12:00:00Z"`
}

// PublicUserResponseWrapper обертка для публичного профиля
type PublicUserResponseWrapper struct {
	Success bool               `json:"success" example:"true"`
	Data    PublicUserResponse `json:"data"`
}

// AdminCheckResponse представляет ответ проверки администратора
type AdminCheckResponse struct {
	Success bool `json:"success" example:"true"`
	IsAdmin bool `json:"is_admin" example:"false"`
}

// AdminCheckResponseWrapper обертка для проверки администратора
type AdminCheckResponseWrapper struct {
	Success bool               `json:"success" example:"true"`
	Data    AdminCheckResponse `json:"data"`
}

// GetProfile returns current user profile
// @Summary Get current user profile
// @Description Returns full profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserProfileResponse "User profile"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "users.profile.error.fetch"
// @Security BearerAuth
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	profile, err := h.services.User().GetUserProfile(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.profile.error.fetch")
	}

	// Проверяем, является ли пользователь администратором
	isAdmin, err := h.userService.IsUserAdmin(c.Context(), profile.Email)
	if err != nil {
		// Если ошибка при проверке админа, логируем но не прерываем запрос
		log.Printf("Error checking admin status for user %d: %v", userID, err)
		isAdmin = false
	}
	
	// Добавляем информацию об админе в профиль
	profile.IsAdmin = isAdmin

	return c.JSON(UserProfileResponse{
		Success: true,
		Data:    profile,
	})
}

// UpdateProfile updates current user profile
// @Summary Update user profile
// @Description Updates profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param profile body models.UserProfileUpdate true "Profile update data"
// @Success 200 {object} MessageResponse "Profile updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "users.profile.error.invalid_data or users.profile.error.validation"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "users.profile.error.update"
// @Security BearerAuth
// @Router /api/v1/users/me [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var update models.UserProfileUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.profile.error.invalid_data")
	}

	// Валидация данных
	if err := update.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.profile.error.validation")
	}

	err := h.services.User().UpdateUserProfile(c.Context(), userID, &update)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.profile.error.update")
	}

	return c.JSON(MessageResponse{
		Success: true,
		Message: "users.profile.success.updated",
	})
}

// RegisterOld регистрирует нового пользователя
// @Summary Register user (deprecated)
// @Description Creates new user in the system. DEPRECATED: Use /api/v1/auth/register instead
// @Tags users
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Registration data (name, email, password required, phone optional)"
// @Success 201 {object} RegisterResponse "User created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "users.register.error.invalid_data or validation errors"
// @Failure 500 {object} utils.ErrorResponseSwag "users.register.error.password_hash_failed or users.register.error.create_failed"
// @Router /api/v1/users/register [post]
// @Deprecated
func (h *UserHandler) RegisterOld(c *fiber.Ctx) error {
	var registerData RegisterRequest

	if err := c.BodyParser(&registerData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.invalid_data")
	}

	// Базовая валидация
	if registerData.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.name_required")
	}
	if registerData.Email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.email_required")
	}
	if registerData.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.password_required")
	}
	if len(registerData.Password) < 6 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.password_too_short")
	}

	// Хеширование пароля
	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.register.error.password_hash_failed")
	}

	user := &models.User{
		Name:     registerData.Name,
		Email:    registerData.Email,
		Password: &hashedPassword,
		Phone:    registerData.Phone,
	}

	err = h.services.User().CreateUser(c.Context(), user)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.email_exists")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.register.error.create_failed")
	}

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		Success: true,
		Message: "users.register.success.created",
		User:    user,
	})
}

// LoginOld авторизует пользователя по email и паролю
// @Summary Login user (deprecated)
// @Description Authenticates user by email and password, creates session and sets session cookie. DEPRECATED: Use /api/v1/auth/login instead
// @Tags users
// @Accept json
// @Produce json
// @Param user body LoginRequest true "Login credentials (email and password required)"
// @Success 200 {object} LoginResponse "Authentication successful"
// @Failure 400 {object} utils.ErrorResponseSwag "users.login.error.invalid_data or validation errors"
// @Failure 401 {object} utils.ErrorResponseSwag "users.login.error.invalid_credentials"
// @Router /api/v1/users/login [post]
// @Deprecated
func (h *UserHandler) LoginOld(c *fiber.Ctx) error {
	var loginData LoginRequest

	if err := c.BodyParser(&loginData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.login.error.invalid_data")
	}

	// Базовая валидация
	if loginData.Email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.login.error.email_required")
	}
	if loginData.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.login.error.password_required")
	}

	// Получаем пользователя по email
	user, err := h.services.User().GetUserByEmail(c.Context(), loginData.Email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.login.error.invalid_credentials")
	}

	// Проверяем пароль
	if user.Password == nil || !utils.CheckPasswordHash(loginData.Password, *user.Password) {
		h.logger.Info("Failed login attempt for user: %s (IP: %s, UserAgent: %s)",
			loginData.Email, c.IP(), c.Get("User-Agent"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.login.error.invalid_credentials")
	}

	// Создаем сессию для пользователя (как и при Google auth)
	sessionToken := utils.GenerateSessionToken()
	sessionData := &types.SessionData{
		UserID:     user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Provider:   "password",
		PictureURL: user.PictureURL,
	}

	// Сохраняем сессию
	h.services.Auth().SaveSession(sessionToken, sessionData)

	// Устанавливаем session cookie (как и при Google auth)
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   3600 * 24, // 24 часа
		Secure:   h.services.Config().GetCookieSecure(),
		HTTPOnly: true,
		SameSite: h.services.Config().GetCookieSameSite(),
	})

	// Логируем успешный логин
	h.logger.Info("Successful login for user: %s (UserID: %d, IP: %s, Provider: password)",
		user.Email, user.ID, c.IP())

	// Возвращаем ответ без JWT токена (используем только сессию)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "users.login.success.authenticated",
		"user":    user,
	})
}

// GetProfileByID returns public user profile
// @Summary Get public user profile
// @Description Returns public information about user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} PublicUserResponseWrapper "Public user profile"
// @Failure 400 {object} utils.ErrorResponseSwag "users.profile.error.invalid_id"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "users.profile.error.not_found"
// @Security BearerAuth
// @Router /api/v1/users/{id}/profile [get]
func (h *UserHandler) GetProfileByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.profile.error.invalid_id")
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "users.profile.error.not_found")
	}

	publicUser := PublicUserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		PictureURL: user.PictureURL,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return c.JSON(PublicUserResponseWrapper{
		Success: true,
		Data:    publicUser,
	})
}

// IsAdminSimple checks if user is administrator (simple implementation)
// @Summary Check admin status (simple)
// @Description Checks if user with specified email is administrator (simplified check by ID)
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} AdminCheckResponseWrapper "Admin status"
// @Failure 400 {object} utils.ErrorResponseSwag "users.admin_check.error.email_required"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "users.admin_check.error.user_not_found"
// @Security BearerAuth
// @Router /api/v1/users/admin-check/{email} [get]
func (h *UserHandler) IsAdminSimple(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.admin_check.error.email_required")
	}

	// Получаем пользователя по email
	user, err := h.userService.GetUserByEmail(c.Context(), email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "users.admin_check.error.user_not_found")
	}

	// Проверяем, является ли пользователь администратором
	// В этой простой версии считаем администраторами только пользователей с ID 1, 2, 3
	isAdmin := user.ID == 1 || user.ID == 2 || user.ID == 3

	return c.JSON(AdminCheckResponseWrapper{
		Success: true,
		Data: AdminCheckResponse{
			IsAdmin: isAdmin,
		},
	})
}
