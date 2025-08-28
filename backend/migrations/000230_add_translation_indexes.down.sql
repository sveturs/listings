-- Откат добавления индексов для системы переводов

-- Удаление индексов таблицы translations
DROP INDEX IF EXISTS idx_translations_entity;
DROP INDEX IF EXISTS idx_translations_lang;
DROP INDEX IF EXISTS idx_translations_field;
DROP INDEX IF EXISTS idx_translations_composite;
DROP INDEX IF EXISTS idx_translations_machine;
DROP INDEX IF EXISTS idx_translations_verified;
DROP INDEX IF EXISTS idx_translations_pending;
DROP INDEX IF EXISTS idx_translations_missing_ru;
DROP INDEX IF EXISTS idx_translations_missing_en;
DROP INDEX IF EXISTS idx_translations_missing_sr;
DROP INDEX IF EXISTS idx_translations_text_gin;

-- Удаление индексов таблицы attribute_option_translations
DROP INDEX IF EXISTS idx_attr_opt_trans_option;

-- Удаление индексов таблицы unit_translations
DROP INDEX IF EXISTS idx_unit_trans_composite;

-- Удаление индексов таблицы translation_versions
DROP INDEX IF EXISTS idx_trans_vers_trans_id;
DROP INDEX IF EXISTS idx_trans_vers_entity;