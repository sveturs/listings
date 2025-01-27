DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 
                  FROM information_schema.columns 
                  WHERE table_name='marketplace_listings' 
                  AND column_name='status') THEN
        ALTER TABLE marketplace_listings 
        ADD COLUMN status VARCHAR(20) DEFAULT 'active';
    END IF;
END $$;

-- Создаем индекс для ускорения поиска по статусу
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_status 
ON marketplace_listings(status);