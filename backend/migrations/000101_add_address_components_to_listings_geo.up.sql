-- Миграция 101: Добавление полей для умного ввода адресов (Phase 2)
-- Добавляем новые поля в таблицу listings_geo для поддержки Phase 2 функциональности

-- Добавляем новые колонки
ALTER TABLE listings_geo 
ADD COLUMN address_components JSONB,
ADD COLUMN geocoding_confidence NUMERIC(3, 2) DEFAULT 0.0,
ADD COLUMN address_verified BOOLEAN DEFAULT false,
ADD COLUMN input_method VARCHAR(50) DEFAULT 'manual',
ADD COLUMN location_privacy VARCHAR(20) DEFAULT 'exact',
ADD COLUMN blurred_location geography(Point, 4326),
ADD COLUMN formatted_address TEXT;

-- Создаем индексы для новых полей
CREATE INDEX IF NOT EXISTS idx_listings_geo_confidence ON listings_geo (geocoding_confidence);
CREATE INDEX IF NOT EXISTS idx_listings_geo_verified ON listings_geo (address_verified);
CREATE INDEX IF NOT EXISTS idx_listings_geo_privacy ON listings_geo (location_privacy);
CREATE INDEX IF NOT EXISTS idx_listings_geo_input_method ON listings_geo (input_method);
CREATE INDEX IF NOT EXISTS idx_listings_geo_blurred_location ON listings_geo USING GIST (blurred_location);
CREATE INDEX IF NOT EXISTS idx_listings_geo_address_components ON listings_geo USING GIN (address_components);

-- Добавляем constraint для location_privacy
ALTER TABLE listings_geo 
ADD CONSTRAINT chk_location_privacy 
CHECK (location_privacy IN ('exact', 'street', 'district', 'city'));

-- Добавляем constraint для input_method
ALTER TABLE listings_geo 
ADD CONSTRAINT chk_input_method 
CHECK (input_method IN ('manual', 'geocoded', 'map_click', 'current_location'));

-- Добавляем constraint для geocoding_confidence (от 0.0 до 1.0)
ALTER TABLE listings_geo 
ADD CONSTRAINT chk_geocoding_confidence 
CHECK (geocoding_confidence >= 0.0 AND geocoding_confidence <= 1.0);

-- Комментарии к новым полям
COMMENT ON COLUMN listings_geo.address_components IS 'JSON структура с компонентами адреса (страна, город, улица и т.д.)';
COMMENT ON COLUMN listings_geo.geocoding_confidence IS 'Уровень доверия геокодирования от 0.0 до 1.0';
COMMENT ON COLUMN listings_geo.address_verified IS 'Подтвержден ли адрес пользователем';
COMMENT ON COLUMN listings_geo.input_method IS 'Способ ввода адреса: manual, geocoded, map_click, current_location';
COMMENT ON COLUMN listings_geo.location_privacy IS 'Уровень приватности: exact, street, district, city';
COMMENT ON COLUMN listings_geo.blurred_location IS 'Размытая локация для отображения в зависимости от настроек приватности';
COMMENT ON COLUMN listings_geo.formatted_address IS 'Отформатированный адрес от провайдера геокодирования';

-- Обновляем триггер updated_at если он существует
DROP TRIGGER IF EXISTS trigger_listings_geo_updated_at ON listings_geo;
CREATE TRIGGER trigger_listings_geo_updated_at
    BEFORE UPDATE ON listings_geo
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();