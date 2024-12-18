// backend/internal/proj/users/storage/postgres/storage.go
package postgres

import (
    "github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
    pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
    return &UserStorage{
        pool: pool,
    }
}

