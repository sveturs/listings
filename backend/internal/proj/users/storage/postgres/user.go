// backend/internal/proj/users/storage/postgres/user.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"errors"
	"strings"
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
		user.Provider = "google"
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
            SET name = $1, google_id = $2, picture_url = $3, provider = 'google'
            WHERE id = $4
        `, user.Name, user.GoogleID, user.PictureURL, userID)
		if err != nil {
			return nil, err
		}
		user.ID = userID
		user.Provider = "google"
		return user, nil
	}

	// Если пользователь не найден ни по google_id, ни по email - создаем нового
	err = s.pool.QueryRow(ctx, `
        INSERT INTO users (name, email, google_id, picture_url, provider)
        VALUES ($1, $2, $3, $4, 'google')
        RETURNING id, created_at
    `, user.Name, user.Email, user.GoogleID, user.PictureURL).Scan(&userID, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	user.ID = userID
	user.Provider = "google"
	return user, nil
}
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, `
        SELECT id, name, email, google_id, picture_url, phone, password, provider, created_at
        FROM users WHERE email = $1
    `, email).Scan(&user.ID, &user.Name, &user.Email, &user.GoogleID, &user.PictureURL, &user.Phone, &user.Password, &user.Provider, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := s.pool.QueryRow(ctx, `
        SELECT id, name, email, google_id, picture_url, phone, password, provider, created_at
        FROM users WHERE id = $1
    `, id).Scan(&user.ID, &user.Name, &user.Email, &user.GoogleID, &user.PictureURL, &user.Phone, &user.Password, &user.Provider, &user.CreatedAt)
	if err != nil {
		s.logger.Error("GetUserByID failed for id=%d: %v", id, err)
		return nil, err
	}
	s.logger.Info("GetUserByID successful for id=%d, provider=%s", id, user.Provider)
	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Установить provider по умолчанию, если не указан
	if user.Provider == "" {
		if user.GoogleID != "" {
			user.Provider = "google"
		} else {
			user.Provider = "email"
		}
	}

	// Для SQL запроса преобразуем *string в интерфейс, который может быть nil
	var passwordValue interface{}
	if user.Password != nil {
		passwordValue = *user.Password
	}

	err := s.pool.QueryRow(ctx, `
        INSERT INTO users (name, email, google_id, picture_url, phone, password, provider)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `, user.Name, user.Email, user.GoogleID, user.PictureURL, user.Phone, passwordValue, user.Provider).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}
	return user, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := s.pool.Exec(ctx, `
        UPDATE users 
        SET name = $1, email = $2, picture_url = $3
        WHERE id = $4
    `, user.Name, user.Email, user.PictureURL, user.ID)
	return err
}

// Refresh Token methods

func (s *Storage) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id, token, expires_at, created_at, user_agent, ip, device_name
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id`

	err := s.pool.QueryRow(
		ctx,
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
		token.UserAgent,
		token.IP,
		token.DeviceName,
	).Scan(&token.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetRefreshToken(ctx context.Context, tokenValue string) (*models.RefreshToken, error) {
	token := &models.RefreshToken{}
	query := `
		SELECT 
			id, user_id, token, expires_at, created_at, 
			user_agent, ip, device_name, is_revoked, revoked_at
		FROM refresh_tokens
		WHERE token = $1 AND NOT is_revoked`

	err := s.pool.QueryRow(ctx, query, tokenValue).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UserAgent,
		&token.IP,
		&token.DeviceName,
		&token.IsRevoked,
		&token.RevokedAt,
	)

	if err != nil {
		s.logger.Info("GetRefreshToken: token=%s..., error=%v", tokenValue[:20], err)
		return nil, err
	}

	s.logger.Info("GetRefreshToken: found token for userID=%d, isRevoked=%v",
		token.UserID, token.IsRevoked)
	return token, nil
}

func (s *Storage) GetRefreshTokenByID(ctx context.Context, id int) (*models.RefreshToken, error) {
	token := &models.RefreshToken{}
	query := `
		SELECT 
			id, user_id, token, expires_at, created_at, 
			user_agent, ip, device_name, is_revoked, revoked_at
		FROM refresh_tokens
		WHERE id = $1`

	err := s.pool.QueryRow(ctx, query, id).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UserAgent,
		&token.IP,
		&token.DeviceName,
		&token.IsRevoked,
		&token.RevokedAt,
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Storage) GetUserRefreshTokens(ctx context.Context, userID int) ([]*models.RefreshToken, error) {
	query := `
		SELECT 
			id, user_id, token, expires_at, created_at, 
			user_agent, ip, device_name, is_revoked, revoked_at
		FROM refresh_tokens
		WHERE user_id = $1 AND NOT is_revoked AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*models.RefreshToken
	for rows.Next() {
		token := &models.RefreshToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.UserAgent,
			&token.IP,
			&token.DeviceName,
			&token.IsRevoked,
			&token.RevokedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (s *Storage) UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	query := `
		UPDATE refresh_tokens 
		SET user_agent = $2, ip = $3, device_name = $4
		WHERE id = $1`

	_, err := s.pool.Exec(
		ctx,
		query,
		token.ID,
		token.UserAgent,
		token.IP,
		token.DeviceName,
	)

	return err
}

func (s *Storage) RevokeRefreshToken(ctx context.Context, tokenID int) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, revoked_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND NOT is_revoked`

	_, err := s.pool.Exec(ctx, query, tokenID)
	return err
}

func (s *Storage) RevokeRefreshTokenByValue(ctx context.Context, tokenValue string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, revoked_at = CURRENT_TIMESTAMP
		WHERE token = $1 AND NOT is_revoked`

	result, err := s.pool.Exec(ctx, query, tokenValue)
	if err != nil {
		s.logger.Error("Failed to revoke refresh token: %v", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	s.logger.Info("RevokeRefreshTokenByValue: token=%s..., rowsAffected=%d",
		tokenValue[:20], rowsAffected)

	if rowsAffected == 0 {
		s.logger.Info("Warning: No rows affected when revoking token, token might already be revoked or not exist")
	}

	return nil
}

func (s *Storage) RevokeUserRefreshTokens(ctx context.Context, userID int) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = true, revoked_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND NOT is_revoked`

	_, err := s.pool.Exec(ctx, query, userID)
	return err
}

func (s *Storage) DeleteExpiredRefreshTokens(ctx context.Context) (int64, error) {
	query := `
		DELETE FROM refresh_tokens 
		WHERE expires_at < CURRENT_TIMESTAMP 
		OR (is_revoked = TRUE AND revoked_at < CURRENT_TIMESTAMP - INTERVAL '30 days')`

	result, err := s.pool.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (s *Storage) CountActiveUserTokens(ctx context.Context, userID int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM refresh_tokens 
		WHERE user_id = $1 AND NOT is_revoked AND expires_at > CURRENT_TIMESTAMP`

	err := s.pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
