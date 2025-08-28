-- Rollback migration: Remove translation admin tables

-- Удаляем триггеры
DROP TRIGGER IF EXISTS translation_versioning ON translations;
DROP TRIGGER IF EXISTS update_translation_providers_timestamp ON translation_providers;

-- Удаляем функции
DROP FUNCTION IF EXISTS create_translation_version();
DROP FUNCTION IF EXISTS update_translation_providers_updated_at();

-- Удаляем таблицы в правильном порядке (учитывая foreign keys)
DROP TABLE IF EXISTS translation_tasks CASCADE;
DROP TABLE IF EXISTS translation_quality_metrics CASCADE;
DROP TABLE IF EXISTS translation_providers CASCADE;
DROP TABLE IF EXISTS translation_audit_log CASCADE;
DROP TABLE IF EXISTS translation_sync_conflicts CASCADE;
DROP TABLE IF EXISTS translation_versions CASCADE;

-- Удаляем добавленную колонку version из translations
ALTER TABLE translations DROP COLUMN IF EXISTS version;