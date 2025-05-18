package postgres

import (
    "github.com/jackc/pgx/v5/pgxpool"
    "backend/internal/proj/marketplace/service"
)

type Storage struct {
    pool *pgxpool.Pool
    translationService service.TranslationServiceInterface
    AttributeGroups AttributeGroupStorage
}

func NewStorage(pool *pgxpool.Pool, translationService service.TranslationServiceInterface) *Storage {
    return &Storage{
        pool: pool,
        translationService: translationService,
        AttributeGroups: NewAttributeGroupStorage(pool),
    }
}