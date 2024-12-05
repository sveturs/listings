package postgres

import (
	"backend/internal/domain/models"
	"context"
)

// GetOrCreateGoogleUser получает или создает пользователя через Google OAuth
func (db *Database) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
	var userID int
	// Пробуем найти пользователя по google_id
	err := db.pool.QueryRow(ctx, `
		SELECT id FROM users WHERE google_id = $1
	`, user.GoogleID).Scan(&userID)

	if err == nil {
		// Пользователь найден, обновляем информацию
		_, err = db.pool.Exec(ctx, `
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

	// Пользователь не найден, создаём нового
	err = db.pool.QueryRow(ctx, `
		INSERT INTO users (name, email, google_id, picture_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO UPDATE 
			SET google_id = $3, 
				picture_url = $4,
				name = $1
		RETURNING id
	`, user.Name, user.Email, user.GoogleID, user.PictureURL).Scan(&userID)

	if err != nil {
		return nil, err
	}

	user.ID = userID
	return user, nil
}

// GetUserByEmail получает пользователя по email
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := db.pool.QueryRow(ctx, `
		SELECT id, name, email, google_id, picture_url, created_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.GoogleID, &user.PictureURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID получает пользователя по ID
func (db *Database) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := db.pool.QueryRow(ctx, `
		SELECT id, name, email, google_id, picture_url, created_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.GoogleID, &user.PictureURL, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser создает нового пользователя
func (db *Database) CreateUser(ctx context.Context, user *models.User) error {
	return db.pool.QueryRow(ctx, `
		INSERT INTO users (name, email, google_id, picture_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, user.Name, user.Email, user.GoogleID, user.PictureURL).Scan(&user.ID)
}

// UpdateUser обновляет информацию о пользователе
func (db *Database) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := db.pool.Exec(ctx, `
		UPDATE users 
		SET name = $1, email = $2, picture_url = $3
		WHERE id = $4
	`, user.Name, user.Email, user.PictureURL, user.ID)
	return err
}
