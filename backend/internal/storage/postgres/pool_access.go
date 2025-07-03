package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetPool возвращает пул подключений к базе данных
func (db *Database) GetPool() *pgxpool.Pool {
	return db.pool
}
