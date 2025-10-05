// backend/internal/proj/users/service/interface.go
package service

import (
	"context"

	"backend/internal/domain/models"
)

type UserServiceInterface interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error
	UpdateLastSeen(ctx context.Context, id int) error

	// Административные методы
	GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error)
	GetAllUsersWithSort(ctx context.Context, limit, offset int, sortBy, sortOrder, statusFilter string) ([]*models.UserProfile, int, error)
	UpdateUserStatus(ctx context.Context, id int, status string) error
	UpdateUserRole(ctx context.Context, id int, roleID int) error
	GetAllRoles(ctx context.Context) ([]*models.Role, error)
	DeleteUser(ctx context.Context, id int) error

	// Методы для управления администраторами
	IsUserAdmin(ctx context.Context, email string) (bool, error)
	GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error)
	AddAdmin(ctx context.Context, admin *models.AdminUser) error
	RemoveAdmin(ctx context.Context, email string) error

	// Методы для настроек приватности
	GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error)
	UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error

	// Методы для настроек чата
	GetChatSettings(ctx context.Context, userID int) (*models.ChatUserSettings, error)
	UpdateChatSettings(ctx context.Context, userID int, settings *models.ChatUserSettings) error
}
