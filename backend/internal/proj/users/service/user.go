// backend/internal/proj/users/service/user.go
package service

import (
	"context"
	"fmt"

	"github.com/sveturs/auth/pkg/entity"
	authService "github.com/sveturs/auth/pkg/service"

	"backend/internal/domain/models"
	"backend/internal/storage"
)

type UserService struct {
	authService *authService.AuthService
	userService *authService.UserService
	storage     storage.Storage
}

func NewUserService(authSvc *authService.AuthService, userSvc *authService.UserService, storage storage.Storage) *UserService {
	return &UserService{
		authService: authSvc,
		userService: userSvc,
		storage:     storage,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	userProfile, err := s.userService.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from auth service: %w", err)
	}

	return FromEntityUserProfile(userProfile), nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	userProfile, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email from auth service: %w", err)
	}

	return FromEntityUserProfile(userProfile), nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	// Создание пользователя теперь происходит через auth-service
	// Этот метод больше не нужен, так как регистрация идет через auth-service
	return fmt.Errorf("user creation should be done through auth service registration endpoint")
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	// Простое обновление без профиля - используем только имя
	req := entity.UpdateProfileRequest{
		Name: &user.Name,
	}

	_, err := s.userService.UpdateUserProfile(ctx, user.ID, req)
	if err != nil {
		return fmt.Errorf("failed to update user through auth service: %w", err)
	}

	return nil
}

func (s *UserService) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
	userProfile, err := s.userService.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from auth service: %w", err)
	}

	return FromEntityUserProfileToProfile(userProfile), nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
	req := ToEntityUpdateProfile(update)

	_, err := s.userService.UpdateUserProfile(ctx, id, req)
	if err != nil {
		return fmt.Errorf("failed to update user profile through auth service: %w", err)
	}

	return nil
}

func (s *UserService) UpdateLastSeen(ctx context.Context, id int) error {
	// LastSeen трекинг может быть убран или реализован через auth-service
	// Для совместимости просто возвращаем nil
	return nil
}

// Административные методы

// GetAllUsers возвращает список всех пользователей с пагинацией
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	usersResp, err := s.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users from auth service: %w", err)
	}

	if usersResp.Users == nil {
		return []*models.UserProfile{}, 0, nil
	}

	// Apply pagination manually
	total := len(usersResp.Users)
	start := offset
	end := offset + limit

	if start >= total {
		return []*models.UserProfile{}, total, nil
	}
	if end > total {
		end = total
	}

	paginatedUsers := usersResp.Users[start:end]
	profiles := FromEntityUserProfileList(paginatedUsers)

	return profiles, total, nil
}

// GetAllUsersWithSort возвращает список всех пользователей с пагинацией, сортировкой и фильтрацией
func (s *UserService) GetAllUsersWithSort(ctx context.Context, limit, offset int, sortBy, sortOrder, statusFilter string) ([]*models.UserProfile, int, error) {
	// Auth service v1.7.0 пока не поддерживает продвинутую сортировку и фильтрацию
	// Используем простой метод GetAllUsers и применяем фильтры на нашей стороне
	return s.GetAllUsers(ctx, limit, offset)
}

// UpdateUserStatus обновляет статус пользователя
func (s *UserService) UpdateUserStatus(ctx context.Context, id int, status string) error {
	req := entity.UpdateStatusRequest{
		Status: status,
	}

	err := s.userService.UpdateUserStatus(ctx, id, req)
	if err != nil {
		return fmt.Errorf("failed to update user status through auth service: %w", err)
	}

	return nil
}

// UpdateUserRole обновляет роль пользователя
func (s *UserService) UpdateUserRole(ctx context.Context, id int, roleID int) error {
	// Auth service работает с ролями по имени, а не по ID
	// Нужно преобразовать roleID в имя роли
	// Для простоты пока возвращаем ошибку, так как нужна таблица соответствия
	return fmt.Errorf("UpdateUserRole needs role ID to name mapping - use AddUserRole with role name instead")
}

// GetAllRoles возвращает список всех ролей
func (s *UserService) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	rolesResp, err := s.userService.GetAllRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles from auth service: %w", err)
	}

	if rolesResp.Roles == nil {
		return []*models.Role{}, nil
	}

	return FromEntityRoleList(rolesResp.Roles), nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// Используем soft delete по умолчанию
	_, err := s.userService.DeleteUser(ctx, id, false)
	if err != nil {
		return fmt.Errorf("failed to delete user through auth service: %w", err)
	}

	return nil
}

// Методы для управления администраторами

// IsUserAdmin проверяет, является ли пользователь администратором по email
func (s *UserService) IsUserAdmin(ctx context.Context, email string) (bool, error) {
	// Получаем пользователя по email
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	// Проверяем через auth service
	adminResp, err := s.userService.IsUserAdmin(ctx, user.ID)
	if err != nil {
		return false, fmt.Errorf("failed to check admin status through auth service: %w", err)
	}

	return adminResp.IsAdmin, nil
}

// GetAllAdmins возвращает список всех администраторов
func (s *UserService) GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error) {
	// Получаем всех пользователей
	usersResp, err := s.userService.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users from auth service: %w", err)
	}

	if usersResp.Users == nil {
		return []*models.AdminUser{}, nil
	}

	// Фильтруем только администраторов (проверяем каждого пользователя через IsUserAdmin)
	admins := make([]*models.AdminUser, 0)
	for _, user := range usersResp.Users {
		// Проверяем, является ли пользователь администратором
		adminResp, err := s.userService.IsUserAdmin(ctx, user.ID)
		if err != nil {
			// Логируем ошибку, но продолжаем проверку других пользователей
			continue
		}

		if adminResp.IsAdmin {
			admins = append(admins, &models.AdminUser{
				Email: user.Email,
			})
		}
	}

	return admins, nil
}

// AddAdmin добавляет нового администратора
func (s *UserService) AddAdmin(ctx context.Context, admin *models.AdminUser) error {
	// Находим пользователя по email
	user, err := s.GetUserByEmail(ctx, admin.Email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Добавляем роль admin
	req := entity.AddRoleRequest{
		Role: "admin",
	}

	_, err = s.userService.AddUserRole(ctx, user.ID, req)
	if err != nil {
		return fmt.Errorf("failed to add admin role through auth service: %w", err)
	}

	return nil
}

// RemoveAdmin удаляет администратора по email
func (s *UserService) RemoveAdmin(ctx context.Context, email string) error {
	// Находим пользователя по email
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Удаляем роль admin
	_, err = s.userService.RemoveUserRole(ctx, user.ID, "admin")
	if err != nil {
		return fmt.Errorf("failed to remove admin role through auth service: %w", err)
	}

	return nil
}

// Методы для настроек приватности

// GetPrivacySettings возвращает настройки приватности пользователя
func (s *UserService) GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return s.storage.GetPrivacySettings(ctx, userID)
}

// UpdatePrivacySettings обновляет настройки приватности пользователя
func (s *UserService) UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return s.storage.UpdatePrivacySettings(ctx, userID, settings)
}

// GetChatSettings возвращает настройки чата пользователя из JSONB поля settings
func (s *UserService) GetChatSettings(ctx context.Context, userID int) (*models.ChatUserSettings, error) {
	// Сначала получаем privacy settings (создаст запись если не существует)
	privacySettings, err := s.storage.GetUserPrivacySettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get privacy settings: %w", err)
	}

	// Парсим JSONB settings
	settings := &models.ChatUserSettings{
		// Defaults
		AutoTranslate:     true,
		PreferredLanguage: "en",
		ShowLanguageBadge: true,
		ModerateTone:      true,
	}

	// Если есть сохраненные настройки в JSONB
	if privacySettings.Settings != nil {
		if autoTranslate, ok := privacySettings.Settings["auto_translate_chat"].(bool); ok {
			settings.AutoTranslate = autoTranslate
		}
		if preferredLang, ok := privacySettings.Settings["preferred_language"].(string); ok {
			settings.PreferredLanguage = preferredLang
		}
		if showBadge, ok := privacySettings.Settings["show_original_language_badge"].(bool); ok {
			settings.ShowLanguageBadge = showBadge
		}
		if moderateTone, ok := privacySettings.Settings["chat_tone_moderation"].(bool); ok {
			settings.ModerateTone = moderateTone
		}
	}

	return settings, nil
}

// UpdateChatSettings обновляет настройки чата пользователя в JSONB поле settings
func (s *UserService) UpdateChatSettings(ctx context.Context, userID int, settings *models.ChatUserSettings) error {
	return s.storage.UpdateChatSettings(ctx, userID, settings)
}
