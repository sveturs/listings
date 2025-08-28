-- Удаляем значение по умолчанию
ALTER TABLE storefront_product_images 
    ALTER COLUMN id DROP DEFAULT;

-- Удаляем последовательность
DROP SEQUENCE IF EXISTS storefront_product_images_id_seq;