-- Откат миграции 000263_full_images_sync_from_minio
-- ВНИМАНИЕ: Полный откат невозможен без резервной копии

-- Возврат URL с IP адресами для тестового окружения (если требуется)
UPDATE storefront_product_images
SET 
    image_url = 'http://100.88.44.15:9000' || image_url,
    thumbnail_url = 'http://100.88.44.15:9000' || thumbnail_url
WHERE 
    image_url NOT LIKE 'http://%' 
    AND image_url NOT LIKE 'https://%'
    AND (image_url LIKE '/storefront-products/%' OR image_url LIKE '/listings/%');

-- Примечание: удаленные записи с паттерном listing_XXX.jpg не могут быть восстановлены
-- Добавленные новые записи останутся в БД (их можно удалить вручную при необходимости)