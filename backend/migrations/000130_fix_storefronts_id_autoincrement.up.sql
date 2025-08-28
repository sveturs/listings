-- Исправление автоинкремента для таблицы storefronts
-- Устанавливаем DEFAULT значение для id с использованием существующей sequence
ALTER TABLE storefronts 
ALTER COLUMN id SET DEFAULT nextval('storefronts_id_seq');

-- Убедимся, что sequence принадлежит колонке
ALTER SEQUENCE storefronts_id_seq OWNED BY storefronts.id;

-- Синхронизируем sequence с максимальным значением id
SELECT setval('storefronts_id_seq', COALESCE((SELECT MAX(id) FROM storefronts), 1), true);