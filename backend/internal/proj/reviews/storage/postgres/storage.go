package postgres

import (
    "github.com/jackc/pgx/v5/pgxpool"
    "backend/internal/proj/marketplace/service"
    "context"
    
    
    "github.com/jackc/pgx/v5"
)

type Storage struct {
    pool *pgxpool.Pool
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
        pool: pool,
        translationService: translationService,
    }
}


func (s *Storage) saveTranslation(ctx context.Context, tx pgx.Tx, entityType string, entityID int, 
    language string, fieldName string, text string, isMachineTranslated bool, isVerified bool) error {
    
    _, err := tx.Exec(ctx, `
        INSERT INTO translations (
            entity_type, entity_id, language, field_name,
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (entity_type, entity_id, language, field_name)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            updated_at = CURRENT_TIMESTAMP
    `, entityType, entityID, language, fieldName, text, isMachineTranslated, isVerified)
    
    return err
}