// backend/internal/proj/users/storage/postgres/admin_users.go
package postgres

import (
	"context"
	"log"

	"backend/internal/domain/models"
)

// IsUserAdmin проверяет, является ли пользователь администратором по email
func (s *Storage) IsUserAdmin(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM admin_users WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking admin status for email %s: %v", email, err)
		return false, err
	}

	return exists, nil
}

// GetAllAdmins возвращает список всех администраторов
func (s *Storage) GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, email, created_at, created_by, notes
		FROM admin_users
		ORDER BY id
	`)
	if err != nil {
		log.Printf("Error getting admin users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var admins []*models.AdminUser

	for rows.Next() {
		admin := &models.AdminUser{}
		err := rows.Scan(&admin.ID, &admin.Email, &admin.CreatedAt, &admin.CreatedBy, &admin.Notes)
		if err != nil {
			log.Printf("Error scanning admin user: %v", err)
			return nil, err
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating admin users: %v", err)
		return nil, err
	}

	return admins, nil
}

// AddAdmin добавляет нового администратора
func (s *Storage) AddAdmin(ctx context.Context, admin *models.AdminUser) error {
	err := s.pool.QueryRow(ctx, `
		INSERT INTO admin_users (email, created_by, notes)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
		RETURNING id
	`, admin.Email, admin.CreatedBy, admin.Notes).Scan(&admin.ID)
	if err != nil {
		log.Printf("Error adding admin user %s: %v", admin.Email, err)
		return err
	}

	return nil
}

// RemoveAdmin удаляет администратора по email
func (s *Storage) RemoveAdmin(ctx context.Context, email string) error {
	_, err := s.pool.Exec(ctx, `
		DELETE FROM admin_users WHERE email = $1
	`, email)
	if err != nil {
		log.Printf("Error removing admin user %s: %v", email, err)
		return err
	}

	return nil
}
