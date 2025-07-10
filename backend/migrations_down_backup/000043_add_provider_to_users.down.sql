-- Удаление индекса
DROP INDEX IF EXISTS idx_users_provider;

-- Удаление поля provider
ALTER TABLE users DROP COLUMN IF EXISTS provider;