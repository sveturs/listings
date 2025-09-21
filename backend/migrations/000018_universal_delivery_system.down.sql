-- ================================================================================
-- ОТКАТ МИГРАЦИИ 000018: УНИВЕРСАЛЬНАЯ СИСТЕМА ДОСТАВКИ
-- ================================================================================

-- 1. УДАЛЯЕМ ТРИГГЕРЫ
-- ================================================================================

DROP TRIGGER IF EXISTS update_delivery_providers_updated_at ON delivery_providers;
DROP TRIGGER IF EXISTS update_delivery_category_defaults_updated_at ON delivery_category_defaults;
DROP TRIGGER IF EXISTS update_delivery_pricing_rules_updated_at ON delivery_pricing_rules;
DROP TRIGGER IF EXISTS update_delivery_shipments_updated_at ON delivery_shipments;

-- 2. УДАЛЯЕМ ФУНКЦИИ
-- ================================================================================

DROP FUNCTION IF EXISTS get_delivery_attributes(INTEGER, VARCHAR);
DROP FUNCTION IF EXISTS calculate_volumetric_weight(NUMERIC, NUMERIC, NUMERIC, INTEGER);
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- 3. УДАЛЯЕМ СВЯЗИ С ЗАКАЗАМИ
-- ================================================================================

ALTER TABLE marketplace_orders
DROP COLUMN IF EXISTS delivery_shipment_id;

-- 4. УДАЛЯЕМ ИНДЕКСЫ
-- ================================================================================

DROP INDEX IF EXISTS idx_delivery_shipments_tracking;
DROP INDEX IF EXISTS idx_delivery_shipments_status;
DROP INDEX IF EXISTS idx_delivery_shipments_order;
DROP INDEX IF EXISTS idx_delivery_tracking_events_shipment;
DROP INDEX IF EXISTS idx_delivery_tracking_events_time;
DROP INDEX IF EXISTS idx_delivery_category_defaults_category;
DROP INDEX IF EXISTS idx_delivery_pricing_rules_provider;
DROP INDEX IF EXISTS idx_delivery_pricing_rules_active;
DROP INDEX IF EXISTS idx_delivery_zones_boundary;
DROP INDEX IF EXISTS idx_delivery_zones_center;

-- 5. УДАЛЯЕМ ТАБЛИЦЫ
-- ================================================================================

DROP TABLE IF EXISTS delivery_zones CASCADE;
DROP TABLE IF EXISTS delivery_pricing_rules CASCADE;
DROP TABLE IF EXISTS delivery_category_defaults CASCADE;
DROP TABLE IF EXISTS delivery_tracking_events CASCADE;
DROP TABLE IF EXISTS delivery_shipments CASCADE;
DROP TABLE IF EXISTS delivery_providers CASCADE;

-- 6. УДАЛЯЕМ АТРИБУТЫ ДОСТАВКИ ИЗ СУЩЕСТВУЮЩИХ ТАБЛИЦ
-- ================================================================================

-- Удаляем атрибуты доставки из marketplace_listings
UPDATE marketplace_listings
SET metadata = metadata - 'delivery_attributes'
WHERE metadata ? 'delivery_attributes';

-- Удаляем атрибуты доставки из storefront_products
UPDATE storefront_products
SET attributes = attributes - 'delivery_attributes'
WHERE attributes ? 'delivery_attributes';

-- ================================================================================
-- КОНЕЦ ОТКАТА
-- ================================================================================