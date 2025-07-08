-- Migration 000082 Down: Удаление таблиц для заказов маркетплейса

-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS marketplace_orders_updated_at_trigger ON marketplace_orders;
DROP FUNCTION IF EXISTS update_marketplace_orders_updated_at();

-- Удаляем таблицы в правильном порядке (с учетом внешних ключей)
DROP TABLE IF EXISTS order_messages;
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS marketplace_orders;