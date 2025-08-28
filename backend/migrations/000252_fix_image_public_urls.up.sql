-- Исправляем public_url для изображений, чтобы все имели относительный путь
-- Убираем полные URL с http://localhost:9000 и оставляем только относительные пути

UPDATE marketplace_images 
SET public_url = REPLACE(public_url, 'http://localhost:9000', '')
WHERE public_url LIKE 'http://localhost:9000%';

UPDATE marketplace_images 
SET public_url = REPLACE(public_url, 'https://localhost:9000', '')
WHERE public_url LIKE 'https://localhost:9000%';

-- Убеждаемся, что все пути начинаются с /listings/
UPDATE marketplace_images 
SET public_url = '/' || file_path
WHERE public_url NOT LIKE '/%' AND storage_type = 'minio';

-- Исправляем пути для основных изображений
UPDATE marketplace_images 
SET public_url = '/listings/' || file_path
WHERE file_path LIKE 'listings/%' 
  AND public_url NOT LIKE '/listings/%'
  AND storage_type = 'minio';

-- Исправляем пути без префикса listings
UPDATE marketplace_images 
SET public_url = '/listings/' || file_path
WHERE file_path NOT LIKE 'listings/%' 
  AND public_url NOT LIKE '/listings/%'
  AND storage_type = 'minio';