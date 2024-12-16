//backend/internal/services/user.go
package services

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/storage"
)

type UserService struct {
    storage storage.Storage
}

func NewUserService(storage storage.Storage) *UserService {
    return &UserService{
        storage: storage,
    }
}
func (s *UserService) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
    return s.storage.GetUserProfile(ctx, id)
}

func (s *UserService) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
    // Валидация данных
    if err := update.Validate(); err != nil {
        return err
    }
    
    return s.storage.UpdateUserProfile(ctx, id, update)
}

func (s *UserService) UpdateLastSeen(ctx context.Context, id int) error {
    return s.storage.UpdateLastSeen(ctx, id)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.storage.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.storage.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.storage.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.storage.UpdateUser(ctx, user)
}
