// backend/internal/proj/users/storage/postgres/user_profile.go
package postgres

import (
    "context"
    "fmt"
    "strings"
    "backend/internal/domain/models"
    "database/sql"
)

// backend/internal/proj/users/storage/postgres/user_profile.go

// backend/internal/proj/users/storage/postgres/user_profile.go

func (s *Storage) GetUserProfile(ctx context.Context, userID int) (*models.UserProfile, error) {
    profile := &models.UserProfile{}
    
    // Используем нативные nullable типы для полей, которые могут быть NULL
    var city, country sql.NullString
    
    err := s.pool.QueryRow(ctx, `
        SELECT 
            id, name, email, google_id, picture_url, created_at,
            phone, bio, notification_email, 
            timezone, last_seen, account_status, settings,
            city, country
        FROM users 
        WHERE id = $1
    `, userID).Scan(
        &profile.ID, &profile.Name, &profile.Email, &profile.GoogleID, &profile.PictureURL,
        &profile.CreatedAt, &profile.Phone, &profile.Bio, &profile.NotificationEmail,
        &profile.Timezone, &profile.LastSeen,
        &profile.AccountStatus, &profile.Settings,
        &city, &country,
    )
    if err != nil {
        return nil, err
    }
    
    // Преобразуем nullable типы в обычные строки
    if city.Valid {
        profile.City = city.String
    }
    if country.Valid {
        profile.Country = country.String
    }
    
    return profile, nil
}

// backend/internal/proj/users/storage/postgres/user_profile.go

func (s *Storage) UpdateUserProfile(ctx context.Context, userID int, update *models.UserProfileUpdate) error {
    var setFields []string
    var params []interface{}
    paramCount := 1

    if update.Phone != nil {
        setFields = append(setFields, fmt.Sprintf("phone = $%d", paramCount))
        params = append(params, update.Phone)
        paramCount++
    }
    if update.Bio != nil {
        setFields = append(setFields, fmt.Sprintf("bio = $%d", paramCount))
        params = append(params, update.Bio)
        paramCount++
    }
    if update.NotificationEmail != nil {
        setFields = append(setFields, fmt.Sprintf("notification_email = $%d", paramCount))
        params = append(params, update.NotificationEmail)
        paramCount++
    }
    if update.Timezone != nil {
        setFields = append(setFields, fmt.Sprintf("timezone = $%d", paramCount))
        params = append(params, update.Timezone)
        paramCount++
    }
    if update.Settings != nil {
        setFields = append(setFields, fmt.Sprintf("settings = $%d", paramCount))
        params = append(params, update.Settings)
        paramCount++
    }
    // Добавляем новые поля для города и страны
    if update.City != nil {
        setFields = append(setFields, fmt.Sprintf("city = $%d", paramCount))
        params = append(params, update.City)
        paramCount++
    }
    if update.Country != nil {
        setFields = append(setFields, fmt.Sprintf("country = $%d", paramCount))
        params = append(params, update.Country)
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

func (s *Storage) UpdateLastSeen(ctx context.Context, userID int) error {
    _, err := s.pool.Exec(ctx, `
        UPDATE users 
        SET last_seen = CURRENT_TIMESTAMP
        WHERE id = $1
    `, userID)
    return err
}