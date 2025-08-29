-- Откат миграции 000262_fix_images_sync_minio
-- ВНИМАНИЕ: Этот откат не восстанавливает удаленные записи с паттерном listing_XXX.jpg

-- Возврат URL с IP адресами (для тестового окружения)
UPDATE storefront_product_images
SET 
    image_url = 'http://100.88.44.15:9000' || image_url,
    thumbnail_url = 'http://100.88.44.15:9000' || thumbnail_url
WHERE 
    image_url NOT LIKE 'http://%' AND 
    image_url NOT LIKE 'https://%' AND
    image_url LIKE '/storefront-products/%';

-- Возврат старого формата путей в marketplace_images
UPDATE marketplace_images
SET 
    file_path = 'listing_' || listing_id || '_' || 
                CASE 
                    WHEN file_name = 'main.jpg' THEN '1.jpg'
                    WHEN file_name = 'image2.jpg' THEN '2.jpg'
                    WHEN file_name = 'image3.jpg' THEN '3.jpg'
                    WHEN file_name = 'image4.jpg' THEN '4.jpg'
                    WHEN file_name = 'image5.jpg' THEN '5.jpg'
                    ELSE file_name
                END,
    public_url = 'http://localhost:9000/listings/listing_' || listing_id || '_' ||
                CASE 
                    WHEN file_name = 'main.jpg' THEN '1.jpg'
                    WHEN file_name = 'image2.jpg' THEN '2.jpg'
                    WHEN file_name = 'image3.jpg' THEN '3.jpg'
                    WHEN file_name = 'image4.jpg' THEN '4.jpg'
                    WHEN file_name = 'image5.jpg' THEN '5.jpg'
                    ELSE file_name
                END
WHERE listing_id BETWEEN 250 AND 267;

-- Удаление добавленных записей (только те, что были добавлены миграцией)
DELETE FROM marketplace_images 
WHERE listing_id = 252 AND file_name = 'image5.jpg' 
AND created_at >= (CURRENT_TIMESTAMP - INTERVAL '1 hour');

DELETE FROM marketplace_images 
WHERE listing_id = 183 AND file_name IN ('main.jpg', 'image2.jpg')
AND created_at >= (CURRENT_TIMESTAMP - INTERVAL '1 hour');

-- Восстановление is_main для объявления 252
UPDATE marketplace_images SET is_main = false WHERE listing_id = 252;
UPDATE marketplace_images SET is_main = true WHERE listing_id = 252 AND file_name = 'image2.jpg';

-- Примечание: полный откат невозможен без резервной копии,
-- так как некоторые данные были удалены в процессе миграции