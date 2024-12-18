
// backend/internal/proj/users/storage/postgres/user.go
package postgres

import (
    "context"
    "backend/internal/domain/models"
)

func (s *Storage) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
    var userID int
    err := s.pool.QueryRow(ctx, `
        SELECT id FROM users WHERE google_id = $1
    `, user.GoogleID).Scan(&userID)

    if err == nil {
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

    err = s.pool.QueryRow(ctx, `
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
