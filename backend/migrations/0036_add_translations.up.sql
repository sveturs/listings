-- /backend/migrations/0036_add_translations.up.sql
CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL, -- 'listing', 'review', 'message'
    entity_id INT NOT NULL,
    field_name VARCHAR(50) NOT NULL, -- 'title', 'description', 'content'
    language VARCHAR(2) NOT NULL, -- 'en', 'sr', 'ru'
    translated_text TEXT NOT NULL,
    is_machine_translated BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_translations_entity ON translations(entity_type, entity_id);
CREATE INDEX idx_translations_language ON translations(language);
CREATE INDEX idx_translations_lookup ON translations(entity_type, entity_id, field_name, language);

-- Добавляем язык оригинала в существующие таблицы
ALTER TABLE marketplace_listings ADD COLUMN original_language VARCHAR(2) DEFAULT 'en';
ALTER TABLE reviews ADD COLUMN original_language VARCHAR(2) DEFAULT 'en';
ALTER TABLE marketplace_messages ADD COLUMN original_language VARCHAR(2) DEFAULT 'en';

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

