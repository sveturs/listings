package postgres

import (
	"backend/internal/proj/marketplace/service"
	"github.com/jackc/pgx/v5/pgxpool"
	//   "context"
	//   "github.com/jackc/pgx/v5"
)

type Storage struct {
	pool               *pgxpool.Pool
	translationService service.TranslationServiceInterface
}

func NewStorage(pool *pgxpool.Pool, translationService service.TranslationServiceInterface) *Storage {
	if pool == nil {
		panic("pool cannot be nil")
	}
	if translationService == nil {
		panic("translationService cannot be nil")
	}
	return &Storage{
		pool:               pool,
		translationService: translationService,
	}
}
