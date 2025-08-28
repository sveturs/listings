-- Откат миграции складской системы и самовывоза

-- Удаление триггеров
DROP TRIGGER IF EXISTS update_warehouse_invoices_updated_at ON warehouse_invoices;
DROP TRIGGER IF EXISTS update_storefront_fbs_settings_updated_at ON storefront_fbs_settings;
DROP TRIGGER IF EXISTS update_warehouse_pickup_orders_updated_at ON warehouse_pickup_orders;
DROP TRIGGER IF EXISTS update_warehouse_inventory_updated_at ON warehouse_inventory;
DROP TRIGGER IF EXISTS update_warehouses_updated_at ON warehouses;
DROP TRIGGER IF EXISTS track_warehouse_inventory_changes ON warehouse_inventory;

-- Удаление функций
DROP FUNCTION IF EXISTS track_inventory_changes();
DROP FUNCTION IF EXISTS reserve_inventory_on_order();
DROP FUNCTION IF EXISTS generate_pickup_code();
DROP FUNCTION IF EXISTS update_warehouse_updated_at();

-- Удаление таблиц в обратном порядке
DROP TABLE IF EXISTS warehouse_invoices;
DROP TABLE IF EXISTS storefront_fbs_settings;
DROP TABLE IF EXISTS warehouse_pickup_orders;
DROP TABLE IF EXISTS warehouse_movements;
DROP TABLE IF EXISTS warehouse_inventory;
DROP TABLE IF EXISTS warehouses;