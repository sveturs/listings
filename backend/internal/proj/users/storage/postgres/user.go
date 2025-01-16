// backend/internal/proj/users/storage/postgres/user.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
)

func (s *Storage) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
    var userID int

    // Сначала пробуем найти по google_id
    err := s.pool.QueryRow(ctx, `
        SELECT id FROM users WHERE google_id = $1
    `, user.GoogleID).Scan(&userID)

    if err == nil {
        // Пользователь найден, обновляем информацию
        _, err = s.pool.Exec(ctx, `
            UPDATE users 
            SET name = $1, email = $2, picture_url = $3
            WHERE id = $4
        `, user.Name, user.Email, user.PictureURL, userID)
        if err != nil {
            return nil, err
        }
        user.ID = userID
        return user, nil
    }

    // Если не нашли по google_id, пробуем найти по email и обновить
    err = s.pool.QueryRow(ctx, `
        SELECT id FROM users WHERE email = $1
    `, user.Email).Scan(&userID)

    if err == nil {
        // Нашли пользователя по email, обновляем его данные
        _, err = s.pool.Exec(ctx, `
            UPDATE users 
            SET name = $1, google_id = $2, picture_url = $3
            WHERE id = $4
        `, user.Name, user.GoogleID, user.PictureURL, userID)
        if err != nil {
            return nil, err
        }
        user.ID = userID
        return user, nil
    }

    // Если пользователь не найден ни по google_id, ни по email - создаем нового
    err = s.pool.QueryRow(ctx, `
        INSERT INTO users (name, email, google_id, picture_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, user.Name, user.Email, user.GoogleID, user.PictureURL).Scan(&userID)

    if err != nil {
        return nil, err
    }

    user.ID = userID
    return user, nil
}
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, `
        SELECT id, name, email, google_id, picture_url, created_at
        FROM users WHERE email = $1
    `, email).Scan(&user.ID, &user.Name, &user.Email, &user.GoogleID, &user.PictureURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, `
        SELECT id, name, email, google_id, picture_url, created_at
        FROM users WHERE id = $1
    `, id).Scan(&user.ID, &user.Name, &user.Email, &user.GoogleID, &user.PictureURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	return s.pool.QueryRow(ctx, `
        INSERT INTO users (name, email, google_id, picture_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, user.Name, user.Email, user.GoogleID, user.PictureURL).Scan(&user.ID)
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE users 
        SET name = $1, email = $2, picture_url = $3
        WHERE id = $4
    `, user.Name, user.Email, user.PictureURL, user.ID)
	return err
}
