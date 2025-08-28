-- Добавляем автоинкремент для поля id в таблице storefront_products

-- Создаем последовательность
CREATE SEQUENCE IF NOT EXISTS storefront_products_id_seq;

-- Устанавливаем владельца последовательности
ALTER SEQUENCE storefront_products_id_seq OWNED BY storefront_products.id;

-- Устанавливаем значение по умолчанию для колонки id
ALTER TABLE storefront_products 
    ALTER COLUMN id SET DEFAULT nextval('storefront_products_id_seq');

-- Устанавливаем текущее значение последовательности на максимальный существующий id + 1
SELECT setval('storefront_products_id_seq', COALESCE((SELECT MAX(id) FROM storefront_products), 0) + 1, false);