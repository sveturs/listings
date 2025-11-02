package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool            *pgxpool.Pool
	AttributeGroups AttributeGroupStorage
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool:            pool,
		AttributeGroups: NewAttributeGroupStorage(pool),
	}
}

// GetPool returns the database pool
func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}
