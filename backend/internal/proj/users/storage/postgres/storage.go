// backend/internal/proj/users/storage/postgres/storage.go
package postgres

import (
    "github.com/jackc/pgx/v5/pgxpool"
    "backend/internal/proj/users/storage"
    "backend/pkg/logger"
)

type Storage struct {
    pool   *pgxpool.Pool
    logger *logger.Logger
}

var _ storage.Repository = (*Storage)(nil)

func NewStorage(pool *pgxpool.Pool) *Storage {
    return &Storage{
        pool:   pool,
        logger: logger.New(),
    }
}