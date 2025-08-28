-- Откат изменений версионирования переводов

-- Удаляем новые триггеры
DROP TRIGGER IF EXISTS translations_versioning_before ON translations;
DROP TRIGGER IF EXISTS translations_versioning_after ON translations;

-- Удаляем функцию
DROP FUNCTION IF EXISTS create_translation_version();

-- Восстанавливаем старую функцию
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

-- Восстанавливаем старый триггер
CREATE TRIGGER translations_versioning
    BEFORE UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION create_translation_version();