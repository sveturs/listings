-- Удаление триггера и функции версионирования
DROP TRIGGER IF EXISTS translations_versioning ON translations;
DROP FUNCTION IF EXISTS create_translation_version() CASCADE;

-- Удаление таблиц версионирования
DROP TABLE IF EXISTS translation_snapshots;
DROP TABLE IF EXISTS translation_versions;

-- Удаление полей версионирования из основной таблицы
ALTER TABLE translations DROP COLUMN IF EXISTS current_version;
ALTER TABLE translations DROP COLUMN IF EXISTS last_modified_by;
ALTER TABLE translations DROP COLUMN IF EXISTS last_modified_at;