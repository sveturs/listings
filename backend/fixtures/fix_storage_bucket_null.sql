-- Исправляем NULL значения в storage_bucket для marketplace_images
-- Это исправит ошибку "can't scan into dest[9]: cannot scan NULL into *string"

-- Обновляем все записи где storage_bucket равен NULL
UPDATE marketplace_images 
SET storage_bucket = 'listings' 
WHERE storage_bucket IS NULL;

-- Проверяем результат
SELECT 
    COUNT(*) as total_images,
    COUNT(storage_bucket) as images_with_bucket,
    COUNT(*) FILTER (WHERE storage_bucket IS NULL) as images_without_bucket
FROM marketplace_images;

-- Показываем несколько примеров исправленных записей
SELECT id, listing_id, public_url, storage_bucket, file_path 
FROM marketplace_images 
WHERE listing_id IN (
    SELECT id FROM marketplace_listings WHERE title LIKE 'Test Phone%'
)
ORDER BY id DESC
LIMIT 10;