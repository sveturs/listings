-- Создаем последовательность для marketplace_categories
CREATE SEQUENCE IF NOT EXISTS marketplace_categories_id_seq;

-- Устанавливаем владельца последовательности
ALTER SEQUENCE marketplace_categories_id_seq OWNED BY marketplace_categories.id;

-- Устанавливаем значение по умолчанию
ALTER TABLE marketplace_categories 
    ALTER COLUMN id SET DEFAULT nextval('marketplace_categories_id_seq');

-- Устанавливаем текущее значение последовательности
-- Используем COALESCE чтобы обработать случай пустой таблицы
SELECT setval('marketplace_categories_id_seq', COALESCE((SELECT MAX(id) FROM marketplace_categories), 0) + 1, false);

-- Добавляем комментарий для документации
COMMENT ON SEQUENCE marketplace_categories_id_seq IS 'Автоинкремент для ID категорий маркетплейса';