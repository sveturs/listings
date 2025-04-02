// backend/internal/storage/postgres/translations.go
package postgres

import (
    "context"
    "fmt"
    "backend/internal/domain/models"
)

// Методы для работы с переводами
func (db *Database) SaveTranslation(ctx context.Context, translation *models.Translation) error {
    query := `
        INSERT INTO translations (
            entity_type, entity_id, field_name, language,
            translated_text, is_machine_translated, is_verified
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (entity_type, entity_id, field_name, language)
        DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            is_machine_translated = EXCLUDED.is_machine_translated,
            is_verified = EXCLUDED.is_verified,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id, created_at, updated_at
    `
    
    err := db.pool.QueryRow(ctx, query,
        translation.EntityType,
        translation.EntityID,
        translation.FieldName,
        translation.Language,
        translation.TranslatedText,
        translation.IsMachineTranslated,
        translation.IsVerified,
    ).Scan(&translation.ID, &translation.CreatedAt, &translation.UpdatedAt)

    if err != nil {
        return fmt.Errorf("failed to save translation: %w", err)
    }

    return nil
}

func (db *Database) GetTranslationsForEntity(ctx context.Context, entityType string, entityID int) ([]models.Translation, error) {
    query := `
        SELECT 
            id, entity_type, entity_id, field_name, language,
            translated_text, is_machine_translated, is_verified,
            created_at, updated_at
        FROM translations
        WHERE entity_type = $1 AND entity_id = $2
    `
    
    rows, err := db.pool.Query(ctx, query, entityType, entityID)
    if err != nil {
        return nil, fmt.Errorf("failed to query translations: %w", err)
    }
    defer rows.Close()

    var translations []models.Translation
    for rows.Next() {
        var t models.Translation
        err := rows.Scan(
            &t.ID, &t.EntityType, &t.EntityID, &t.FieldName, &t.Language,
            &t.TranslatedText, &t.IsMachineTranslated, &t.IsVerified,
            &t.CreatedAt, &t.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan translation: %w", err)
        }
        translations = append(translations, t)
    }

    return translations, nil
}

func (db *Database) UpdateTranslation(ctx context.Context, translation *models.Translation) error {
    query := `
        UPDATE translations
        SET 
            translated_text = $1,
            is_machine_translated = $2,
            is_verified = $3,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $4
        RETURNING updated_at
    `
    
    err := db.pool.QueryRow(ctx, query,
        translation.TranslatedText,
        translation.IsMachineTranslated,
        translation.IsVerified,
        translation.ID,
    ).Scan(&translation.UpdatedAt)

    if err != nil {
        return fmt.Errorf("failed to update translation: %w", err)
    }

    return nil
}

func (db *Database) DeleteTranslations(ctx context.Context, entityType string, entityID int) error {
    _, err := db.pool.Exec(ctx, `
        DELETE FROM translations
        WHERE entity_type = $1 AND entity_id = $2
    `, entityType, entityID)

    if err != nil {
        return fmt.Errorf("failed to delete translations: %w", err)
    }

    return nil
}