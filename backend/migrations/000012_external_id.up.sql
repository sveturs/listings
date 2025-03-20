-- Добавление поля external_id в таблицу marketplace_listings
ALTER TABLE marketplace_listings ADD COLUMN IF NOT EXISTS external_id VARCHAR(255);

-- Добавление индекса для быстрого поиска по external_id
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_external_id ON marketplace_listings(external_id);

-- Добавление составного индекса для быстрого поиска по external_id и storefront_id
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_external_id_storefront_id 
ON marketplace_listings(external_id, storefront_id);