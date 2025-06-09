// backend/internal/proj/users/storage/interface.go
package storage

import (
	"backend/internal/domain/models"
	"context"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	// User methods
	GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

// ProfileRepository определяет интерфейс для работы с профилями пользователей
type ProfileRepository interface {
	GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error
	UpdateLastSeen(ctx context.Context, id int) error
}

// AdminRepository определяет интерфейс для работы с администраторами
type AdminRepository interface {
	IsUserAdmin(ctx context.Context, email string) (bool, error)
	GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error)
	AddAdmin(ctx context.Context, admin *models.AdminUser) error
	RemoveAdmin(ctx context.Context, email string) error
}

// Repository объединяет все репозитории для домена users
type Repository interface {
	UserRepository
	ProfileRepository
	AdminRepository
}
