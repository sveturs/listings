// backend/internal/storage/postgres/admin_methods.go
package postgres

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
)

// IsUserAdmin проверяет, является ли пользователь администратором по email
// DEPRECATED: moved to auth-service
func (db *Database) IsUserAdmin(ctx context.Context, email string) (bool, error) {
	return false, fmt.Errorf("IsUserAdmin: moved to auth-service")
}

// GetAllAdmins возвращает список всех администраторов
// DEPRECATED: moved to auth-service
func (db *Database) GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error) {
	return nil, fmt.Errorf("GetAllAdmins: moved to auth-service")
}

// AddAdmin добавляет нового администратора
// DEPRECATED: moved to auth-service
func (db *Database) AddAdmin(ctx context.Context, admin *models.AdminUser) error {
	return fmt.Errorf("AddAdmin: moved to auth-service")
}

// RemoveAdmin удаляет администратора по email
// DEPRECATED: moved to auth-service
func (db *Database) RemoveAdmin(ctx context.Context, email string) error {
	return fmt.Errorf("RemoveAdmin: moved to auth-service")
}
