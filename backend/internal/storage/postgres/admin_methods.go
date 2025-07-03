// backend/internal/storage/postgres/admin_methods.go
package postgres

import (
	"context"

	"backend/internal/domain/models"
)

// IsUserAdmin проверяет, является ли пользователь администратором по email
func (db *Database) IsUserAdmin(ctx context.Context, email string) (bool, error) {
	return db.usersDB.IsUserAdmin(ctx, email)
}

// GetAllAdmins возвращает список всех администраторов
func (db *Database) GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error) {
	return db.usersDB.GetAllAdmins(ctx)
}

// AddAdmin добавляет нового администратора
func (db *Database) AddAdmin(ctx context.Context, admin *models.AdminUser) error {
	return db.usersDB.AddAdmin(ctx, admin)
}

// RemoveAdmin удаляет администратора по email
func (db *Database) RemoveAdmin(ctx context.Context, email string) error {
	return db.usersDB.RemoveAdmin(ctx, email)
}
