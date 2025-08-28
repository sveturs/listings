-- Исправление автоинкремента для таблицы users
-- Устанавливаем DEFAULT значение для id с использованием существующей sequence
ALTER TABLE users 
ALTER COLUMN id SET DEFAULT nextval('users_id_seq');

-- Убедимся, что sequence принадлежит колонке
ALTER SEQUENCE users_id_seq OWNED BY users.id;

-- Синхронизируем sequence с максимальным значением id
SELECT setval('users_id_seq', COALESCE((SELECT MAX(id) FROM users), 1), true);