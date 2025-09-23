-- Удаляем индекс
DROP INDEX IF EXISTS idx_marketplace_listings_address_multilingual;

-- Удаляем колонку с мультиязычными адресами
ALTER TABLE marketplace_listings
DROP COLUMN IF EXISTS address_multilingual;