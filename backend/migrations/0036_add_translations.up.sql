-- backend/migrations/0036_add_translations.up.sql
CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    language VARCHAR(10) NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    translated_text TEXT NOT NULL,
    is_machine_translated BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(entity_type, entity_id, language, field_name)
);

CREATE INDEX idx_translations_lookup ON translations(entity_type, entity_id, language);

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_translations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_translations_timestamp
    BEFORE UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION update_translations_updated_at();

-- Добавляем поле original_language в таблицу marketplace_listings, если его ещё нет
ALTER TABLE marketplace_listings 
ADD COLUMN IF NOT EXISTS original_language VARCHAR(10) DEFAULT 'sr';