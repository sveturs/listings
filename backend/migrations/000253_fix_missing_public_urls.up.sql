-- Исправляем пустые public_url для всех изображений
UPDATE marketplace_images 
SET public_url = CASE
    WHEN file_path LIKE 'listings/%' THEN '/' || file_path
    ELSE '/listings/' || file_path
END
WHERE public_url IS NULL OR public_url = '';

-- Устанавливаем NOT NULL constraint для public_url чтобы предотвратить эту проблему в будущем
-- (закомментировано, так как может сломать старый код)
-- ALTER TABLE marketplace_images ALTER COLUMN public_url SET NOT NULL;