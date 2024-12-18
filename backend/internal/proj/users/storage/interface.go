// backend/internal/proj/users/storage/interface.go

package storage

import (
    "context"
    "backend/internal/domain/models"
)

type UserStorage interface {
    // User methods
    GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    GetUserByID(ctx context.Context, id int) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    UpdateUser(ctx context.Context, user *models.User) error
    
    // Profile methods
    GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error)
    UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error
    UpdateLastSeen(ctx context.Context, id int) error
}

