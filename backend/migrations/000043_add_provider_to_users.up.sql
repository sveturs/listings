-- Добавление поля provider в таблицу users
ALTER TABLE users ADD COLUMN IF NOT EXISTS provider VARCHAR(50) DEFAULT 'email';

-- Обновляем provider для существующих Google пользователей
UPDATE users SET provider = 'google' WHERE google_id IS NOT NULL AND google_id != '';

-- Создаем индекс для быстрого поиска по provider
CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider);