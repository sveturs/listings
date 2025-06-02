-- Добавление поля password в таблицу users
ALTER TABLE users ADD COLUMN password VARCHAR(255);

-- Создание индекса для быстрого поиска по email (если еще не существует)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Удаление уникального ограничения с колонки google_id таблицы users
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_google_id_key;