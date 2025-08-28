-- Устанавливаем значение по умолчанию для NULL storage_bucket
-- Используем 'listings' как значение по умолчанию для совместимости с MinIO
UPDATE marketplace_images 
SET storage_bucket = 'listings'
WHERE storage_bucket IS NULL AND storage_type = 'minio';

-- Для локальных файлов оставляем NULL или устанавливаем пустую строку
UPDATE marketplace_images 
SET storage_bucket = ''
WHERE storage_bucket IS NULL AND storage_type = 'local';

-- storefront_product_images не имеет колонки storage_bucket, поэтому пропускаем