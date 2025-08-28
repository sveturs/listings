-- Откат миграции Post Express интеграции

-- Удаление триггеров
DROP TRIGGER IF EXISTS log_post_express_shipment_status ON post_express_shipments;
DROP TRIGGER IF EXISTS update_post_express_shipments_updated_at ON post_express_shipments;
DROP TRIGGER IF EXISTS update_post_express_rates_updated_at ON post_express_rates;
DROP TRIGGER IF EXISTS update_post_express_offices_updated_at ON post_express_offices;
DROP TRIGGER IF EXISTS update_post_express_locations_updated_at ON post_express_locations;
DROP TRIGGER IF EXISTS update_post_express_settings_updated_at ON post_express_settings;

-- Удаление функций
DROP FUNCTION IF EXISTS log_shipment_status_change();
DROP FUNCTION IF EXISTS update_post_express_updated_at();

-- Удаление полей из существующих таблиц
ALTER TABLE marketplace_orders 
DROP COLUMN IF EXISTS delivery_method,
DROP COLUMN IF EXISTS delivery_cost,
DROP COLUMN IF EXISTS delivery_address,
DROP COLUMN IF EXISTS delivery_notes;

ALTER TABLE storefront_orders 
DROP COLUMN IF EXISTS delivery_method,
DROP COLUMN IF EXISTS delivery_cost,
DROP COLUMN IF EXISTS delivery_address,
DROP COLUMN IF EXISTS delivery_notes;

-- Удаление таблиц в обратном порядке (из-за foreign keys)
DROP TABLE IF EXISTS post_express_tracking_events;
DROP TABLE IF EXISTS post_express_api_logs;
DROP TABLE IF EXISTS post_express_shipments;
DROP TABLE IF EXISTS post_express_rates;
DROP TABLE IF EXISTS post_express_offices;
DROP TABLE IF EXISTS post_express_locations;
DROP TABLE IF EXISTS post_express_settings;