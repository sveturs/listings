package repository

import (
	"context"

	"backend/internal/domain/models"
)

// UserRepositoryInterface определяет методы для работы с пользователями
type UserRepositoryInterface interface {
	GetByID(ctx context.Context, userID int) (*models.User, error)
}
