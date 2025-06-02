-- Удаление поля password из таблицы users
ALTER TABLE users DROP COLUMN IF EXISTS password;

-- Восстановление уникального ограничения на колонку google_id таблицы users
ALTER TABLE users ADD CONSTRAINT users_google_id_key UNIQUE (google_id);