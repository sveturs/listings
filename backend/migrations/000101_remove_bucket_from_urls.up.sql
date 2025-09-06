-- Удаляем имя bucket из URL изображений, оставляя только путь к файлу
-- Теперь bucket будет подставляться динамически из конфигурации

-- Обновляем URL изображений в таблице marketplace_images
-- Убираем /listings/ и /development-listings/ из начала пути
UPDATE marketplace_images
SET public_url = SUBSTRING(public_url FROM '/[^/]+/(.*)$')
WHERE public_url LIKE '/listings/%' OR public_url LIKE '/development-listings/%';

-- Обновляем URL изображений в таблице storefront_product_images если она существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'storefront_product_images') THEN
        UPDATE storefront_product_images
        SET 
            image_url = CASE 
                WHEN image_url LIKE '/listings/%' OR image_url LIKE '/development-listings/%'
                THEN SUBSTRING(image_url FROM '/[^/]+/(.*)$')
                ELSE image_url
            END,
            thumbnail_url = CASE
                WHEN thumbnail_url LIKE '/listings/%' OR thumbnail_url LIKE '/development-listings/%'
                THEN SUBSTRING(thumbnail_url FROM '/[^/]+/(.*)$')
                ELSE thumbnail_url
            END
        WHERE image_url LIKE '/listings/%' OR image_url LIKE '/development-listings/%' 
           OR thumbnail_url LIKE '/listings/%' OR thumbnail_url LIKE '/development-listings/%';
    END IF;
END $$;

-- Обновляем URL в marketplace_listing_variants если колонка существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'marketplace_listing_variants' AND column_name = 'image_url') THEN
        UPDATE marketplace_listing_variants
        SET image_url = SUBSTRING(image_url FROM '/[^/]+/(.*)$')
        WHERE image_url LIKE '/listings/%' OR image_url LIKE '/development-listings/%';
    END IF;
END $$;