package service

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/types"
)

type AuthServiceInterface interface {
    GetGoogleAuthURL() string
    HandleGoogleCallback(ctx context.Context, code string) (*types.SessionData, error)
    SaveSession(token string, data *types.SessionData)
    GetSession(token string) (*types.SessionData, bool)
    DeleteSession(token string)
}

type UserServiceInterface interface {
    GetUserByID(ctx context.Context, id int) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    UpdateUser(ctx context.Context, user *models.User) error
    GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error)
    UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error
    UpdateLastSeen(ctx context.Context, id int) error
}