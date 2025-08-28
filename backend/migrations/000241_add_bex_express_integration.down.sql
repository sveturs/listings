-- Удаление интеграции BEX Express

-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_bex_settings_updated_at ON bex_settings;
DROP TRIGGER IF EXISTS update_bex_municipalities_updated_at ON bex_municipalities;
DROP TRIGGER IF EXISTS update_bex_places_updated_at ON bex_places;
DROP TRIGGER IF EXISTS update_bex_streets_updated_at ON bex_streets;
DROP TRIGGER IF EXISTS update_bex_parcel_shops_updated_at ON bex_parcel_shops;
DROP TRIGGER IF EXISTS update_bex_shipments_updated_at ON bex_shipments;
DROP TRIGGER IF EXISTS update_bex_rates_updated_at ON bex_rates;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_bex_updated_at();

-- Удаляем колонки из таблиц заказов
ALTER TABLE marketplace_orders 
DROP COLUMN IF EXISTS delivery_provider,
DROP COLUMN IF EXISTS delivery_tracking_number,
DROP COLUMN IF EXISTS delivery_status,
DROP COLUMN IF EXISTS delivery_metadata;

ALTER TABLE storefront_orders 
DROP COLUMN IF EXISTS delivery_provider,
DROP COLUMN IF EXISTS delivery_tracking_number,
DROP COLUMN IF EXISTS delivery_status,
DROP COLUMN IF EXISTS delivery_metadata;

-- Удаляем таблицы
DROP TABLE IF EXISTS bex_tracking_events;
DROP TABLE IF EXISTS bex_shipments;
DROP TABLE IF EXISTS bex_parcel_shops;
DROP TABLE IF EXISTS bex_streets;
DROP TABLE IF EXISTS bex_places;
DROP TABLE IF EXISTS bex_municipalities;
DROP TABLE IF EXISTS bex_rates;
DROP TABLE IF EXISTS bex_settings;