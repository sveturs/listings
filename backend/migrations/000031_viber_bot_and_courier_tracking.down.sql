-- Удаляем в обратном порядке

-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_viber_users_updated_at ON viber_users;
DROP TRIGGER IF EXISTS update_couriers_updated_at ON couriers;
DROP TRIGGER IF EXISTS update_deliveries_updated_at ON deliveries;

-- Удаляем функции
DROP FUNCTION IF EXISTS close_expired_viber_sessions();
DROP FUNCTION IF EXISTS calculate_distance(NUMERIC, NUMERIC, NUMERIC, NUMERIC);

-- Удаляем таблицы
DROP TABLE IF EXISTS delivery_notifications CASCADE;
DROP TABLE IF EXISTS tracking_websocket_connections CASCADE;
DROP TABLE IF EXISTS viber_tracking_sessions CASCADE;
DROP TABLE IF EXISTS courier_location_history CASCADE;
DROP TABLE IF EXISTS deliveries CASCADE;
DROP TABLE IF EXISTS courier_zones CASCADE;
DROP TABLE IF EXISTS couriers CASCADE;
DROP TABLE IF EXISTS viber_messages CASCADE;
DROP TABLE IF EXISTS viber_sessions CASCADE;
DROP TABLE IF EXISTS viber_users CASCADE;