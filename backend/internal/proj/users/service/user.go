// backend/internal/proj/users/service/user.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
)

type UserService struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) *UserService {
	return &UserService{
		storage: storage,
	}
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

func (s *UserService) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
	return s.storage.GetUserProfile(ctx, id)
}

func (s *UserService) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
	return s.storage.UpdateUserProfile(ctx, id, update)
}

func (s *UserService) UpdateLastSeen(ctx context.Context, id int) error {
	return s.storage.UpdateLastSeen(ctx, id)
}

// Административные методы

// GetAllUsers возвращает список всех пользователей с пагинацией
func (s *UserService) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	return s.storage.GetAllUsers(ctx, limit, offset)
}

// UpdateUserStatus обновляет статус пользователя
func (s *UserService) UpdateUserStatus(ctx context.Context, id int, status string) error {
	return s.storage.UpdateUserStatus(ctx, id, status)
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.storage.DeleteUser(ctx, id)
}
