ALTER TABLE marketplace_images
ADD COLUMN IF NOT EXISTS storage_type VARCHAR(20) DEFAULT 'local',
ADD COLUMN IF NOT EXISTS storage_bucket VARCHAR(100),
ADD COLUMN IF NOT EXISTS public_url TEXT;

-- Обновляем существующие записи
UPDATE marketplace_images
SET storage_type = 'local',
    public_url = CONCAT('/uploads/', file_path)
WHERE storage_type = 'local' OR storage_type IS NULL;