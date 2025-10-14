// backend/internal/storage/postgres/db_users.go
package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// User methods - DEPRECATED: moved to auth-service
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, fmt.Errorf("GetUserByEmail: moved to auth-service")
}

func (db *Database) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return nil, fmt.Errorf("GetUserByID: moved to auth-service")
}

func (db *Database) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, fmt.Errorf("CreateUser: moved to auth-service")
}

func (db *Database) UpdateUser(ctx context.Context, user *models.User) error {
	return fmt.Errorf("UpdateUser: moved to auth-service")
}

func (db *Database) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, fmt.Errorf("GetOrCreateGoogleUser: moved to auth-service")
}

func (db *Database) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
	return nil, fmt.Errorf("GetUserProfile: moved to auth-service")
}

func (db *Database) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
	return fmt.Errorf("UpdateUserProfile: moved to auth-service")
}

func (db *Database) UpdateLastSeen(ctx context.Context, id int) error {
	return fmt.Errorf("UpdateLastSeen: moved to auth-service")
}

// Административные методы для управления пользователями - DEPRECATED: moved to auth-service
func (db *Database) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	return nil, 0, fmt.Errorf("GetAllUsers: moved to auth-service")
}

func (db *Database) GetAllUsersWithSort(ctx context.Context, limit, offset int, sortBy, sortOrder, statusFilter string) ([]*models.UserProfile, int, error) {
	return nil, 0, fmt.Errorf("GetAllUsersWithSort: moved to auth-service")
}

func (db *Database) UpdateUserStatus(ctx context.Context, id int, status string) error {
	return fmt.Errorf("UpdateUserStatus: moved to auth-service")
}

func (db *Database) UpdateUserRole(ctx context.Context, id int, roleID int) error {
	return fmt.Errorf("UpdateUserRole: moved to auth-service")
}

func (db *Database) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	return nil, fmt.Errorf("GetAllRoles: moved to auth-service")
}

func (db *Database) DeleteUser(ctx context.Context, id int) error {
	return fmt.Errorf("DeleteUser: moved to auth-service")
}
