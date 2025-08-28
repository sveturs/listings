-- Удаляем индекс
DROP INDEX IF EXISTS idx_users_preferred_language;

-- Удаляем ограничение
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_preferred_language_check;

-- Удаляем колонку
ALTER TABLE users DROP COLUMN IF EXISTS preferred_language;