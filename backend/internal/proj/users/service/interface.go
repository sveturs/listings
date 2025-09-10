// backend/internal/proj/users/service/interface.go
package service

import (
	"context"

	"backend/internal/domain/models"
	"backend/internal/types"
)

type AuthServiceInterface interface {
	GetGoogleAuthURL(origin string) string
	HandleGoogleCallback(ctx context.Context, code string) (*types.SessionData, error)
	SaveSession(token string, data *types.SessionData)
	DeleteSession(token string)
	GetSession(ctx context.Context, token string) (*types.SessionData, error)

	// JWT методы
	GenerateJWT(userID int, email string) (string, error)
	ValidateJWT(tokenString string) (interface{}, error)

	// Email/Password аутентификация
	LoginWithEmailPassword(ctx context.Context, email, password string) (string, *models.User, error)
	RegisterWithEmailPassword(ctx context.Context, name, email, password string) (string, *models.User, error)

	// Refresh Token методы
	LoginWithRefreshToken(ctx context.Context, email, password, ip, userAgent string) (accessToken, refreshToken string, user *models.User, err error)
	RegisterWithRefreshToken(ctx context.Context, name, email, password, ip, userAgent string) (accessToken, refreshToken string, user *models.User, err error)
	GenerateTokensForOAuth(ctx context.Context, userID int, email, ip, userAgent string) (accessToken, refreshToken string, err error)
	RefreshTokens(ctx context.Context, refreshToken, ip, userAgent string) (newAccessToken, newRefreshToken string, err error)
	RevokeRefreshToken(ctx context.Context, refreshToken string) error
	RevokeUserRefreshTokens(ctx context.Context, userID int) error
}

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
}
