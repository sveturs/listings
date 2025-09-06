-- Откат изменений - добавляем обратно /listings/ к путям
-- (используем /listings/ как дефолтное значение для обратной совместимости)

UPDATE marketplace_images
SET public_url = '/listings/' || public_url
WHERE public_url NOT LIKE '/%';

-- Откат URL в storefront_product_images если таблица существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'storefront_product_images') THEN
        UPDATE storefront_product_images
        SET 
            image_url = CASE 
                WHEN image_url NOT LIKE '/%' THEN '/listings/' || image_url
                ELSE image_url
            END,
            thumbnail_url = CASE
                WHEN thumbnail_url NOT LIKE '/%' THEN '/listings/' || thumbnail_url
                ELSE thumbnail_url
            END
        WHERE image_url NOT LIKE '/%' OR thumbnail_url NOT LIKE '/%';
    END IF;
END $$;

-- Откат URL в marketplace_listing_variants если колонка существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'marketplace_listing_variants' AND column_name = 'image_url') THEN
        UPDATE marketplace_listing_variants
        SET image_url = '/listings/' || image_url
        WHERE image_url NOT LIKE '/%';
    END IF;
END $$;