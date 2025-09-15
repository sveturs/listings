-- Откат миграции: удаление синхронизированных данных

-- Удаление данных из listings_geo для объявлений из marketplace_listings
DELETE FROM listings_geo
WHERE listing_id IN (
    SELECT id FROM marketplace_listings
);

-- Удаление созданных индексов
DROP INDEX IF EXISTS idx_listings_geo_privacy_level;
DROP INDEX IF EXISTS idx_listings_geo_listing_id;