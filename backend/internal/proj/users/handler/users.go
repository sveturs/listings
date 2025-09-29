package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
)

type UserHandler struct {
	services    globalService.ServicesInterface
	userService service.UserServiceInterface
}

func NewUserHandler(services globalService.ServicesInterface) *UserHandler {
	return &UserHandler{
		services:    services,
		userService: services.User(),
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
	Message string `json:"message" example:"Операция выполнена успешно"`
}

// RegisterResponse представляет ответ после регистрации
type RegisterResponse struct {
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
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=backend_internal_domain_models.UserProfile} "User profile"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "users.profile.error.fetch"
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
		logger.Error().Err(err).Int("user_id", userID).Msg("Error checking admin status")
		isAdmin = false
	}

	// Добавляем информацию об админе в профиль
	profile.IsAdmin = isAdmin

	return utils.SuccessResponse(c, profile)
}

// UpdateProfile updates current user profile
// @Summary Update user profile
// @Description Updates profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param profile body backend_internal_domain_models.UserProfileUpdate true "Profile update data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Profile updated successfully"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "users.profile.error.invalid_data or users.profile.error.validation"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "users.profile.error.update"
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

	return utils.SuccessResponse(c, &MessageResponse{
		Message: "users.profile.success.updated",
	})
}

// GetProfileByID returns public user profile
// @Summary Get public user profile
// @Description Returns public information about user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=PublicUserResponse} "Public user profile"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "users.profile.error.invalid_id"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "users.profile.error.not_found"
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

	return utils.SuccessResponse(c, publicUser)
}

// IsAdminSimple checks if user is administrator (simple implementation)
// @Summary Check admin status (simple)
// @Description Checks if user with specified email is administrator (simplified check by ID)
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=AdminCheckResponse} "Admin status"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "users.admin_check.error.email_required"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "users.admin_check.error.user_not_found"
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

	return utils.SuccessResponse(c, &AdminCheckResponse{
		IsAdmin: isAdmin,
	})
}

func (h *UserHandler) GetPrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	settings, err := h.userService.GetPrivacySettings(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("user_id", userID).Msg("Error fetching privacy settings")
		// Возвращаем дефолтные настройки если их нет
		settings = &models.UserPrivacySettings{
			UserID:                        userID,
			AllowContactRequests:          true,
			AllowMessagesFromContactsOnly: false,
		}
	}

	return utils.SuccessResponse(c, settings)
}

// UpdatePrivacySettings обновляет настройки приватности пользователя
// @Summary Update privacy settings
// @Description Updates privacy settings for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param settings body backend_internal_domain_models.UpdatePrivacySettingsRequest true "Privacy settings"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=MessageResponse} "Settings updated"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "users.privacy.error.invalid_data"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "users.privacy.error.update"
// @Security BearerAuth
// @Router /api/v1/users/privacy-settings [put]
func (h *UserHandler) UpdatePrivacySettings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var settings models.UpdatePrivacySettingsRequest
	if err := c.BodyParser(&settings); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.privacy.error.invalid_data")
	}

	err := h.userService.UpdatePrivacySettings(c.Context(), userID, &settings)
	if err != nil {
		logger.Error().Err(err).Int("user_id", userID).Msg("Error updating privacy settings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.privacy.error.update")
	}

	return utils.SuccessResponse(c, &MessageResponse{
		Message: "users.privacy.success.updated",
	})
}
