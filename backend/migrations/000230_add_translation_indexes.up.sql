-- Добавление оптимизационных индексов для системы переводов
-- Дата: 2025-08-12

-- Основные индексы для таблицы translations
CREATE INDEX IF NOT EXISTS idx_translations_entity ON translations(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_translations_lang ON translations(language);
CREATE INDEX IF NOT EXISTS idx_translations_field ON translations(field_name);
CREATE INDEX IF NOT EXISTS idx_translations_composite ON translations(entity_type, entity_id, language, field_name);

-- Частичные индексы для оптимизации частых запросов
CREATE INDEX IF NOT EXISTS idx_translations_machine ON translations(is_machine_translated) WHERE is_machine_translated = true;
CREATE INDEX IF NOT EXISTS idx_translations_verified ON translations(is_verified) WHERE is_verified = false;
CREATE INDEX IF NOT EXISTS idx_translations_pending ON translations(entity_type, entity_id) 
    WHERE is_machine_translated = true AND is_verified = false;

-- Индексы для поиска недостающих переводов
CREATE INDEX IF NOT EXISTS idx_translations_missing_ru ON translations(entity_type, entity_id) 
    WHERE language = 'ru' AND translated_text IS NULL;
CREATE INDEX IF NOT EXISTS idx_translations_missing_en ON translations(entity_type, entity_id) 
    WHERE language = 'en' AND translated_text IS NULL;
CREATE INDEX IF NOT EXISTS idx_translations_missing_sr ON translations(entity_type, entity_id) 
    WHERE language = 'sr' AND translated_text IS NULL;

-- Дополнительные индексы для attribute_option_translations
CREATE INDEX IF NOT EXISTS idx_attr_opt_trans_option ON attribute_option_translations(option_value);

-- Индексы для unit_translations
CREATE INDEX IF NOT EXISTS idx_unit_trans_composite ON unit_translations(unit, language);

-- Индексы для translation_versions
CREATE INDEX IF NOT EXISTS idx_trans_vers_trans_id ON translation_versions(translation_id);
CREATE INDEX IF NOT EXISTS idx_trans_vers_entity ON translation_versions(entity_type, entity_id);

-- GIN индекс для полнотекстового поиска (опционально)
CREATE INDEX IF NOT EXISTS idx_translations_text_gin ON translations USING gin(to_tsvector('simple', translated_text));

-- Комментарии к основным индексам
COMMENT ON INDEX idx_translations_entity IS 'Оптимизация поиска переводов по типу и ID сущности';
COMMENT ON INDEX idx_translations_composite IS 'Композитный индекс для основных запросов переводов';
COMMENT ON INDEX idx_translations_pending IS 'Частичный индекс для непроверенных машинных переводов';