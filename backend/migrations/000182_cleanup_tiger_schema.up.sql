-- Удаление Tiger schema (PostGIS US addresses - не нужен для проекта в Сербии)
-- Освобождает ~2.1 MB места

DROP SCHEMA IF EXISTS tiger CASCADE;

-- Удаляем также topology schema если она не используется
DROP SCHEMA IF EXISTS topology CASCADE;
