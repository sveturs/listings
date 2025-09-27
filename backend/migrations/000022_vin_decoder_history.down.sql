-- Откат миграции VIN декодера

-- Удаление триггера
DROP TRIGGER IF EXISTS update_vin_decode_cache_timestamp ON vin_decode_cache;
DROP FUNCTION IF EXISTS update_vin_decode_cache_updated_at();

-- Удаление индексов
DROP INDEX IF EXISTS idx_vin_accident_history_vin;
DROP INDEX IF EXISTS idx_vin_ownership_history_vin;
DROP INDEX IF EXISTS idx_vin_recalls_vin;
DROP INDEX IF EXISTS idx_vin_check_history_checked_at;
DROP INDEX IF EXISTS idx_vin_check_history_listing_id;
DROP INDEX IF EXISTS idx_vin_check_history_vin;
DROP INDEX IF EXISTS idx_vin_check_history_user_id;
DROP INDEX IF EXISTS idx_vin_decode_cache_year;
DROP INDEX IF EXISTS idx_vin_decode_cache_make_model;
DROP INDEX IF EXISTS idx_vin_decode_cache_vin;

-- Удаление таблиц (в обратном порядке из-за внешних ключей)
DROP TABLE IF EXISTS vin_accident_history;
DROP TABLE IF EXISTS vin_ownership_history;
DROP TABLE IF EXISTS vin_recalls;
DROP TABLE IF EXISTS vin_check_history;
DROP TABLE IF EXISTS vin_decode_cache;