-- Исправление проблемы с дубликатами версий в translation_versions

-- Удаляем старые триггеры
DROP TRIGGER IF EXISTS translation_version_trigger ON translations;
DROP TRIGGER IF EXISTS translations_versioning ON translations;

-- Удаляем старую функцию
DROP FUNCTION IF EXISTS create_translation_version();

-- Создаём исправленную функцию для версионирования
CREATE OR REPLACE FUNCTION create_translation_version() RETURNS TRIGGER AS $$
DECLARE
    next_version INTEGER;
BEGIN
    -- Для INSERT
    IF TG_OP = 'INSERT' THEN
        -- Получаем следующую версию для этого translation_id
        SELECT COALESCE(MAX(version), 0) + 1 INTO next_version
        FROM translation_versions
        WHERE translation_id = NEW.id;
        
        -- Устанавливаем версию в основной таблице
        NEW.current_version := next_version;
        
        -- Создаём запись версии только если это происходит после триггера
        IF TG_WHEN = 'AFTER' THEN
            INSERT INTO translation_versions (
                translation_id, entity_type, entity_id, field_name, language,
                translated_text, previous_text, version, change_type, changed_by
            ) VALUES (
                NEW.id, NEW.entity_type, NEW.entity_id, NEW.field_name, NEW.language,
                NEW.translated_text, NULL, next_version, 'created', NEW.last_modified_by
            )
            ON CONFLICT (translation_id, version) DO NOTHING;
        END IF;
        
        RETURN NEW;
        
    -- Для UPDATE
    ELSIF TG_OP = 'UPDATE' THEN
        -- Проверяем, действительно ли изменился текст
        IF OLD.translated_text IS DISTINCT FROM NEW.translated_text THEN
            -- Получаем следующую версию для этого translation_id
            SELECT COALESCE(MAX(version), 0) + 1 INTO next_version
            FROM translation_versions
            WHERE translation_id = NEW.id;
            
            -- Устанавливаем версию в основной таблице
            NEW.current_version := next_version;
            NEW.last_modified_at := CURRENT_TIMESTAMP;
            
            -- Создаём запись версии только если это происходит после триггера
            IF TG_WHEN = 'AFTER' THEN
                INSERT INTO translation_versions (
                    translation_id, entity_type, entity_id, field_name, language,
                    translated_text, previous_text, version, change_type, changed_by
                ) VALUES (
                    NEW.id, NEW.entity_type, NEW.entity_id, NEW.field_name, NEW.language,
                    NEW.translated_text, OLD.translated_text, next_version, 'updated', NEW.last_modified_by
                )
                ON CONFLICT (translation_id, version) DO NOTHING;
            END IF;
        END IF;
        
        RETURN NEW;
        
    -- Для DELETE
    ELSIF TG_OP = 'DELETE' THEN
        IF TG_WHEN = 'AFTER' THEN
            -- Получаем следующую версию для этого translation_id
            SELECT COALESCE(MAX(version), 0) + 1 INTO next_version
            FROM translation_versions
            WHERE translation_id = OLD.id;
            
            INSERT INTO translation_versions (
                translation_id, entity_type, entity_id, field_name, language,
                translated_text, previous_text, version, change_type, changed_by
            ) VALUES (
                OLD.id, OLD.entity_type, OLD.entity_id, OLD.field_name, OLD.language,
                NULL, OLD.translated_text, next_version, 'deleted', OLD.last_modified_by
            )
            ON CONFLICT (translation_id, version) DO NOTHING;
        END IF;
        
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Создаём один BEFORE триггер для установки версий
CREATE TRIGGER translations_versioning_before
    BEFORE INSERT OR UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION create_translation_version();

-- Создаём один AFTER триггер для записи в таблицу версий
CREATE TRIGGER translations_versioning_after
    AFTER INSERT OR UPDATE OR DELETE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION create_translation_version();

-- Очищаем существующие дубликаты в translation_versions если они есть
DELETE FROM translation_versions a
WHERE EXISTS (
    SELECT 1
    FROM translation_versions b
    WHERE a.translation_id = b.translation_id
      AND a.version = b.version
      AND a.id > b.id
);

-- Исправляем последовательность версий для существующих записей
WITH version_fix AS (
    SELECT 
        id,
        translation_id,
        ROW_NUMBER() OVER (PARTITION BY translation_id ORDER BY changed_at, id) as new_version
    FROM translation_versions
)
UPDATE translation_versions tv
SET version = vf.new_version
FROM version_fix vf
WHERE tv.id = vf.id;

-- Обновляем current_version в таблице translations
UPDATE translations t
SET current_version = (
    SELECT COALESCE(MAX(version), 1)
    FROM translation_versions tv
    WHERE tv.translation_id = t.id
)
WHERE EXISTS (
    SELECT 1 FROM translation_versions tv
    WHERE tv.translation_id = t.id
);

-- Устанавливаем current_version = 1 для записей без версий
UPDATE translations
SET current_version = 1
WHERE current_version IS NULL OR current_version = 0;