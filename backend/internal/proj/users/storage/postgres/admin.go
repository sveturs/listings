// backend/internal/proj/users/storage/postgres/admin.go
package postgres

import (
	"context"
	"database/sql"
	"log"

	"backend/internal/domain/models"
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
	// Логируем начало процесса удаления
	log.Printf("Starting deletion process for user ID %d", id)

	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	// 1. Сначала проверим все объявления и удалим зависимые данные
	log.Printf("Deleting data for user's marketplace listings")

	// Удаляем все изображения для объявлений пользователя одним запросом
	result, err := tx.Exec(ctx, `
		DELETE FROM marketplace_images 
		WHERE listing_id IN (
			SELECT id FROM marketplace_listings WHERE user_id = $1
		)`, id)
	if err != nil {
		log.Printf("Error deleting marketplace_images: %v", err)
		return err
	}
	imagesDeleted := result.RowsAffected()
	log.Printf("Deleted %d images for user %d listings", imagesDeleted, id)

	// 2. Удаляем избранное (маркетплейс)
	log.Printf("Deleting marketplace_favorites")
	_, err = tx.Exec(ctx, `DELETE FROM marketplace_favorites WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting marketplace_favorites: %v", err)
		return err
	}

	// 3. Теперь удаляем сообщения в чатах
	log.Printf("Deleting marketplace_messages")
	_, err = tx.Exec(ctx, `DELETE FROM marketplace_messages WHERE sender_id = $1 OR receiver_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting marketplace_messages: %v", err)
		return err
	}

	// 4. Удаляем чаты
	log.Printf("Deleting marketplace_chats")
	_, err = tx.Exec(ctx, `DELETE FROM marketplace_chats WHERE buyer_id = $1 OR seller_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting marketplace_chats: %v", err)
		return err
	}

	// 5. Теперь можно удалить сами объявления
	log.Printf("Deleting marketplace_listings")
	_, err = tx.Exec(ctx, `DELETE FROM marketplace_listings WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting marketplace_listings: %v", err)
		return err
	}

	// 6. Для витрин и с ними связанных данных
	log.Printf("Processing user_storefronts")

	// Удаляем историю импортов для всех витрин пользователя одним запросом
	log.Printf("Deleting import_history for all user's storefronts")
	_, err = tx.Exec(ctx, `
		DELETE FROM import_history 
		WHERE source_id IN (
			SELECT is.id 
			FROM import_sources is
			JOIN user_storefronts us ON is.storefront_id = us.id
			WHERE us.user_id = $1
		)`, id)
	if err != nil {
		log.Printf("Error deleting import_history: %v", err)
		return err
	}

	// Удаляем источники импорта
	log.Printf("Deleting import_sources")
	_, err = tx.Exec(ctx, `
		DELETE FROM import_sources 
		WHERE storefront_id IN (
			SELECT id FROM user_storefronts WHERE user_id = $1
		)`, id)
	if err != nil {
		log.Printf("Error deleting import_sources: %v", err)
		return err
	}

	// Удаляем витрины
	log.Printf("Deleting user_storefronts")
	_, err = tx.Exec(ctx, `DELETE FROM user_storefronts WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting user_storefronts: %v", err)
		return err
	}

	// 7. Удаляем данные баланса
	log.Printf("Deleting balance_transactions")
	_, err = tx.Exec(ctx, `DELETE FROM balance_transactions WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting balance_transactions: %v", err)
		return err
	}

	log.Printf("Deleting user_balances")
	_, err = tx.Exec(ctx, `DELETE FROM user_balances WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting user_balances: %v", err)
		return err
	}

	// 8. Удаляем уведомления и настройки уведомлений
	log.Printf("Deleting notifications")
	_, err = tx.Exec(ctx, `DELETE FROM notifications WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting notifications: %v", err)
		return err
	}

	log.Printf("Deleting notification_settings")
	_, err = tx.Exec(ctx, `DELETE FROM notification_settings WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting notification_settings: %v", err)
		return err
	}

	// 9. Удаляем телеграм-соединения
	log.Printf("Deleting user_telegram_connections")
	_, err = tx.Exec(ctx, `DELETE FROM user_telegram_connections WHERE user_id = $1`, id)
	if err != nil {
		log.Printf("Error deleting user_telegram_connections: %v", err)
		return err
	}

	// 10. Отзывы имеют каскадное удаление (ON DELETE CASCADE), но проверим на всякий случай
	log.Printf("Checking reviews (should be handled by ON DELETE CASCADE)")
	var reviewCount int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM reviews WHERE user_id = $1`, id).Scan(&reviewCount)
	if err != nil {
		log.Printf("Error checking reviews count: %v", err)
	} else if reviewCount > 0 {
		log.Printf("Found %d reviews that should be deleted by ON DELETE CASCADE", reviewCount)
	}

	// 11. Наконец, удаляем самого пользователя
	log.Printf("Deleting user")
	_, err = tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}

	// Фиксируем транзакцию
	log.Printf("Committing transaction for user deletion")
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("Successfully deleted user ID %d", id)
	return nil
}
