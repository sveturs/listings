-- Откат миграции: удаление таблиц import_jobs и import_errors
DROP TABLE IF EXISTS import_errors CASCADE;
DROP TABLE IF EXISTS import_jobs CASCADE;
DROP FUNCTION IF EXISTS update_import_jobs_updated_at() CASCADE;
