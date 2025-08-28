-- Сначала проверяем, не существует ли уже таблица translation_versions
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'translation_versions') THEN
        -- Создание таблицы версий переводов для отслеживания истории изменений
        CREATE TABLE translation_versions (
            id BIGSERIAL PRIMARY KEY,
            translation_id INTEGER NOT NULL REFERENCES translations(id) ON DELETE CASCADE,
            entity_type VARCHAR(50) NOT NULL,
            entity_id INTEGER NOT NULL,
            field_name VARCHAR(50) NOT NULL,
            language VARCHAR(10) NOT NULL,
            translated_text TEXT NOT NULL,
            previous_text TEXT,
            version INTEGER NOT NULL DEFAULT 1,
            change_type VARCHAR(20) NOT NULL CHECK (change_type IN ('created', 'updated', 'deleted', 'restored')),
            changed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
            changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
            change_reason TEXT,
            metadata JSONB DEFAULT '{}',
            
            -- Индексы для быстрого поиска
            CONSTRAINT translation_versions_unique_version UNIQUE (translation_id, version)
        );

        -- Индексы для оптимизации запросов
        CREATE INDEX idx_translation_versions_translation_id ON translation_versions(translation_id);
        CREATE INDEX idx_translation_versions_entity ON translation_versions(entity_type, entity_id);
        CREATE INDEX idx_translation_versions_changed_at ON translation_versions(changed_at DESC);
        CREATE INDEX idx_translation_versions_changed_by ON translation_versions(changed_by);
        CREATE INDEX idx_translation_versions_change_type ON translation_versions(change_type);
    END IF;
END $$;

-- Добавление полей версионирования в основную таблицу переводов
ALTER TABLE translations ADD COLUMN IF NOT EXISTS current_version INTEGER DEFAULT 1;
ALTER TABLE translations ADD COLUMN IF NOT EXISTS last_modified_by BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE translations ADD COLUMN IF NOT EXISTS last_modified_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Удаляем старую функцию и триггер, если они существуют
DROP TRIGGER IF EXISTS translations_versioning ON translations;
DROP FUNCTION IF EXISTS create_translation_version();

-- Функция для автоматического создания версии при изменении перевода
CREATE OR REPLACE FUNCTION create_translation_version() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO translation_versions (
            translation_id, entity_type, entity_id, field_name, language, 
            translated_text, previous_text, version, change_type, changed_by, change_reason
        ) VALUES (
            NEW.id, NEW.entity_type, NEW.entity_id, NEW.field_name, NEW.language,
            NEW.translated_text, NULL, 1, 'created', NEW.last_modified_by, 'Initial creation'
        );
        NEW.current_version := 1;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Только если значение действительно изменилось
        IF OLD.translated_text IS DISTINCT FROM NEW.translated_text THEN
            NEW.current_version := COALESCE(OLD.current_version, 1) + 1;
            NEW.last_modified_at := CURRENT_TIMESTAMP;
            
            INSERT INTO translation_versions (
                translation_id, entity_type, entity_id, field_name, language,
                translated_text, previous_text, version, change_type, changed_by
            ) VALUES (
                NEW.id, NEW.entity_type, NEW.entity_id, NEW.field_name, NEW.language,
                NEW.translated_text, OLD.translated_text, NEW.current_version, 'updated', NEW.last_modified_by
            );
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматического версионирования (только для UPDATE, так как на INSERT пока не работает из-за отсутствия last_modified_by)
CREATE TRIGGER translations_versioning
    BEFORE UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION create_translation_version();

-- Таблица для хранения снимков состояния всех переводов (для массового отката)
CREATE TABLE IF NOT EXISTS translation_snapshots (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    translations_count INTEGER NOT NULL DEFAULT 0,
    snapshot_data JSONB NOT NULL,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_translation_snapshots_created_at ON translation_snapshots(created_at DESC);
CREATE INDEX idx_translation_snapshots_created_by ON translation_snapshots(created_by);

-- Комментарии к таблицам
COMMENT ON TABLE translation_versions IS 'История всех изменений переводов с возможностью отката';
COMMENT ON TABLE translation_snapshots IS 'Снимки состояния всех переводов для массового восстановления';
COMMENT ON COLUMN translation_versions.change_type IS 'Тип изменения: created, updated, deleted, restored';
COMMENT ON COLUMN translation_versions.metadata IS 'Дополнительная информация: источник перевода (manual, ai), провайдер AI и т.д.';