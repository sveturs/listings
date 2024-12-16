// backend/internal/storage/postgres/user_profile.go
package postgres

import (
    "context"
    "backend/internal/domain/models"
 //   "encoding/json"
	"fmt"
	"strings"
)

// GetUserProfile получает расширенный профиль пользователя
func (db *Database) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
    profile := &models.UserProfile{}
    err := db.pool.QueryRow(ctx, `
        SELECT 
            id, name, email, google_id, picture_url, created_at,
            phone, bio, notification_email, notification_push,
            timezone, last_seen, account_status, settings
        FROM users 
        WHERE id = $1
    `, userID).Scan(
        &profile.ID, &profile.Name, &profile.Email, &profile.GoogleID, &profile.PictureURL,
        &profile.CreatedAt, &profile.Phone, &profile.Bio, &profile.NotificationEmail,
        &profile.NotificationPush, &profile.Timezone, &profile.LastSeen,
        &profile.AccountStatus, &profile.Settings,
    )
    if err != nil {
        return nil, err
    }
    return profile, nil
}

// UpdateUserProfile обновляет профиль пользователя
func (db *Database) UpdateUserProfile(ctx context.Context, userID int, update *models.UserProfileUpdate) error {
    // Создаем слайс для параметров и строки обновления
    var params []interface{}
    var updates []string
    paramCount := 1

    // Добавляем только не-nil поля в обновление
    if update.Phone != nil {
        updates = append(updates, fmt.Sprintf("phone = $%d", paramCount))
        params = append(params, update.Phone)
        paramCount++
    }
    if update.Bio != nil {
        updates = append(updates, fmt.Sprintf("bio = $%d", paramCount))
        params = append(params, update.Bio)
        paramCount++
    }
    if update.NotificationEmail != nil {
        updates = append(updates, fmt.Sprintf("notification_email = $%d", paramCount))
        params = append(params, update.NotificationEmail)
        paramCount++
    }
    if update.NotificationPush != nil {
        updates = append(updates, fmt.Sprintf("notification_push = $%d", paramCount))
        params = append(params, update.NotificationPush)
        paramCount++
    }
    if update.Timezone != nil {
        updates = append(updates, fmt.Sprintf("timezone = $%d", paramCount))
        params = append(params, update.Timezone)
        paramCount++
    }
    if update.Settings != nil {
        updates = append(updates, fmt.Sprintf("settings = $%d", paramCount))
        params = append(params, update.Settings)
        paramCount++
    }

    // Если нет обновлений, возвращаем nil
    if len(updates) == 0 {
        return nil
    }

    // Добавляем ID пользователя в параметры
    params = append(params, userID)

    // Формируем и выполняем запрос
    query := fmt.Sprintf(`
        UPDATE users 
        SET %s, updated_at = CURRENT_TIMESTAMP
        WHERE id = $%d
    `, strings.Join(updates, ", "), paramCount)

    _, err := db.pool.Exec(ctx, query, params...)
    return err
}

// UpdateLastSeen обновляет время последнего посещения пользователя
func (db *Database) UpdateLastSeen(ctx context.Context, userID int) error {
    _, err := db.pool.Exec(ctx, `
        UPDATE users 
        SET last_seen = CURRENT_TIMESTAMP
        WHERE id = $1
    `, userID)
    return err
}