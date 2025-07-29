-- Удаление триггера и функции
DROP TRIGGER IF EXISTS update_inventory_reservations_updated_at_trigger ON inventory_reservations;
DROP FUNCTION IF EXISTS update_inventory_reservations_updated_at();

-- Удаление индексов
DROP INDEX IF EXISTS idx_inventory_reservations_expires_at;
DROP INDEX IF EXISTS idx_inventory_reservations_status;
DROP INDEX IF EXISTS idx_inventory_reservations_variant_id;
DROP INDEX IF EXISTS idx_inventory_reservations_product_id;
DROP INDEX IF EXISTS idx_inventory_reservations_order_id;

-- Удаление таблицы
DROP TABLE IF EXISTS inventory_reservations;