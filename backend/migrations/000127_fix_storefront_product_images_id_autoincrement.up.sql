-- Создаем последовательность для storefront_product_images
CREATE SEQUENCE IF NOT EXISTS storefront_product_images_id_seq;

-- Устанавливаем владельца последовательности
ALTER SEQUENCE storefront_product_images_id_seq OWNED BY storefront_product_images.id;

-- Устанавливаем значение по умолчанию
ALTER TABLE storefront_product_images 
    ALTER COLUMN id SET DEFAULT nextval('storefront_product_images_id_seq');

-- Устанавливаем текущее значение последовательности
SELECT setval('storefront_product_images_id_seq', COALESCE((SELECT MAX(id) FROM storefront_product_images), 0) + 1, false);

-- Добавляем комментарий для документации
COMMENT ON SEQUENCE storefront_product_images_id_seq IS 'Автоинкремент для ID изображений товаров витрин';