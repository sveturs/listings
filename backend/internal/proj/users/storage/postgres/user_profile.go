
// backend/internal/proj/users/storage/postgres/user_profile.go
package postgres

import (
    "context"
    "backend/internal/domain/models"
)

func (s *UserStorage) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
    profile := &models.UserProfile{}
    err := s.pool.QueryRow(ctx, `
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

func (s *UserStorage) UpdateUserProfile(ctx context.Context, userID int, update *models.UserProfileUpdate) error {
    // Создаем слайс для параметров и строки обновления
    setFields := make([]string, 0)
    params := make([]interface{}, 0)
    paramCount := 1

    // Добавляем только не-nil поля в обновление
    if update.Phone != nil {
        setFields = append(setFields, fmt.Sprintf("phone = $%d", paramCount))
        params = append(params, *update.Phone)
        paramCount++
    }
    if update.Bio != nil {
        setFields = append(setFields, fmt.Sprintf("bio = $%d", paramCount))
        params = append(params, *update.Bio)
        paramCount++
    }
    if update.NotificationEmail != nil {
        setFields = append(setFields, fmt.Sprintf("notification_email = $%d", paramCount))
        params = append(params, *update.NotificationEmail)
        paramCount++
    }
    if update.NotificationPush != nil {
        setFields = append(setFields, fmt.Sprintf("notification_push = $%d", paramCount))
        params = append(params, *update.NotificationPush)
        paramCount++
    }
    if update.Timezone != nil {
        setFields = append(setFields, fmt.Sprintf("timezone = $%d", paramCount))
        params = append(params, *update.Timezone)
        paramCount++
    }
    if update.Settings != nil {
        setFields = append(setFields, fmt.Sprintf("settings = $%d", paramCount))
        params = append(params, update.Settings)
        paramCount++
    }

    if len(setFields) == 0 {
        return nil
    }

    params = append(params, userID)
    query := fmt.Sprintf(`
        UPDATE users 
        SET %s, updated_at = CURRENT_TIMESTAMP
        WHERE id = $%d
    `, strings.Join(setFields, ", "), paramCount)

    _, err := s.pool.Exec(ctx, query, params...)
    return err
}

func (s *UserStorage) UpdateLastSeen(ctx context.Context, userID int) error {
    _, err := s.pool.Exec(ctx, `
        UPDATE users 
        SET last_seen = CURRENT_TIMESTAMP
        WHERE id = $1
    `, userID)
    return err
}