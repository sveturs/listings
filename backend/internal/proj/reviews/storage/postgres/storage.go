package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	//   "context"
	//   "github.com/jackc/pgx/v5"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool, _ interface{}) *Storage {
	if pool == nil {
		panic("pool cannot be nil")
	}
	// translation service removed, keeping interface compatible
	return &Storage{
		pool: pool,
	}
}
