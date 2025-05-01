// backend/internal/proj/users/storage/postgres/admin.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
)

// GetAllUsers возвращает список всех пользователей с пагинацией
func (s *Storage) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	// Сначала получаем общее количество пользователей
	var total int
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Затем получаем пользователей с пагинацией
	query := `
        SELECT
            id, name, email, google_id, picture_url, created_at,
            phone, bio, notification_email,
            timezone, last_seen, account_status, settings,
            city, country
        FROM users
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	rows, err := s.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.UserProfile

	for rows.Next() {
		profile := &models.UserProfile{}

		// Используем sql.Null* типы для всех полей, которые могут быть NULL
		var googleID, pictureURL, phone, bio, city, country sql.NullString
		var lastSeen sql.NullTime
		var notificationEmail sql.NullBool
		var settings sql.NullString

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Email,
			&googleID,
			&pictureURL,
			&profile.CreatedAt,
			&phone,
			&bio,
			&notificationEmail,
			&profile.Timezone,
			&lastSeen,
			&profile.AccountStatus,
			&settings,
			&city,
			&country,
		)
		if err != nil {
			return nil, 0, err
		}

		// Преобразуем nullable типы в обычные
		if googleID.Valid {
			profile.GoogleID = googleID.String
		}
		if pictureURL.Valid {
			profile.PictureURL = pictureURL.String
		}
		if phone.Valid {
			phoneStr := phone.String
			profile.Phone = &phoneStr
		}
		if bio.Valid {
			bioStr := bio.String
			profile.Bio = &bioStr
		}
		if notificationEmail.Valid {
			profile.NotificationEmail = notificationEmail.Bool
		}
		if lastSeen.Valid {
			profile.LastSeen = &lastSeen.Time
		}
		if settings.Valid {
			profile.Settings = []byte(settings.String)
		}
		if city.Valid {
			profile.City = city.String
		}
		if country.Valid {
			profile.Country = country.String
		}

		users = append(users, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUserStatus обновляет статус пользователя
func (s *Storage) UpdateUserStatus(ctx context.Context, id int, status string) error {
	_, err := s.pool.Exec(ctx, `
		UPDATE users
		SET account_status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, id)

	return err
}

// DeleteUser удаляет пользователя
func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Удаляем связанные данные (можно добавить больше таблиц по мере необходимости)

	// Удаляем уведомления пользователя
	_, err = tx.Exec(ctx, `DELETE FROM notifications WHERE user_id = $1`, id)
	if err != nil {
		return err
	}

	// Удаляем телеграм-соединения пользователя
	_, err = tx.Exec(ctx, `DELETE FROM user_telegram_connections WHERE user_id = $1`, id)
	if err != nil {
		return err
	}

	// Удаляем избранное пользователя
	_, err = tx.Exec(ctx, `DELETE FROM user_favorites WHERE user_id = $1`, id)
	if err != nil {
		return err
	}

	// Удаляем самого пользователя
	_, err = tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	// Фиксируем транзакцию
	return tx.Commit(ctx)
}
