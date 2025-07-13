-- Откат миграции 101: Удаление полей Phase 2 из listings_geo

-- Удаляем триггер
DROP TRIGGER IF EXISTS trigger_listings_geo_updated_at ON listings_geo;

-- Удаляем constraints
ALTER TABLE listings_geo DROP CONSTRAINT IF EXISTS chk_location_privacy;
ALTER TABLE listings_geo DROP CONSTRAINT IF EXISTS chk_input_method;
ALTER TABLE listings_geo DROP CONSTRAINT IF EXISTS chk_geocoding_confidence;

-- Удаляем индексы
DROP INDEX IF EXISTS idx_listings_geo_confidence;
DROP INDEX IF EXISTS idx_listings_geo_verified;
DROP INDEX IF EXISTS idx_listings_geo_privacy;
DROP INDEX IF EXISTS idx_listings_geo_input_method;
DROP INDEX IF EXISTS idx_listings_geo_blurred_location;
DROP INDEX IF EXISTS idx_listings_geo_address_components;

-- Удаляем колонки
ALTER TABLE listings_geo 
DROP COLUMN IF EXISTS address_components,
DROP COLUMN IF EXISTS geocoding_confidence,
DROP COLUMN IF EXISTS address_verified,
DROP COLUMN IF EXISTS input_method,
DROP COLUMN IF EXISTS location_privacy,
DROP COLUMN IF EXISTS blurred_location,
DROP COLUMN IF EXISTS formatted_address;

-- Восстанавливаем оригинальный триггер updated_at
CREATE TRIGGER trigger_listings_geo_updated_at
    BEFORE UPDATE ON listings_geo
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();