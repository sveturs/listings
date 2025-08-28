-- Откат изменений - возвращаем полные URL
UPDATE marketplace_images 
SET public_url = 'http://localhost:9000' || public_url
WHERE public_url LIKE '/listings/%' 
  AND storage_type = 'minio';