-- Исправление автоинкремента для таблицы storefront_staff
-- Устанавливаем DEFAULT значение для id с использованием существующей sequence
ALTER TABLE storefront_staff 
ALTER COLUMN id SET DEFAULT nextval('storefront_staff_id_seq');

-- Убедимся, что sequence принадлежит колонке
ALTER SEQUENCE storefront_staff_id_seq OWNED BY storefront_staff.id;

-- Синхронизируем sequence с максимальным значением id
SELECT setval('storefront_staff_id_seq', COALESCE((SELECT MAX(id) FROM storefront_staff), 1), true);