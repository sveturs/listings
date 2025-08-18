package postgres

import (
	"backend/internal/proj/marketplace/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool               *pgxpool.Pool
	translationService service.TranslationServiceInterface
	AttributeGroups    AttributeGroupStorage
}

func NewStorage(pool *pgxpool.Pool, translationService service.TranslationServiceInterface) *Storage {
	return &Storage{
		pool:               pool,
		translationService: translationService,
		AttributeGroups:    NewAttributeGroupStorage(pool),
	}
}

// GetPool returns the database pool
func (s *Storage) GetPool() *pgxpool.Pool {
	return s.pool
}
