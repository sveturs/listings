-- Откат изменений для автоинкремента id в таблице storefront_products

-- Удаляем значение по умолчанию
ALTER TABLE storefront_products 
    ALTER COLUMN id DROP DEFAULT;

-- Удаляем последовательность
DROP SEQUENCE IF EXISTS storefront_products_id_seq;