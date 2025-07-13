-- Откат миграции 102: Удаление таблицы логов изменений адресов

-- Удаляем функцию
DROP FUNCTION IF EXISTS cleanup_old_address_logs();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_address_log_listing_id;
DROP INDEX IF EXISTS idx_address_log_user_id;
DROP INDEX IF EXISTS idx_address_log_created_at;
DROP INDEX IF EXISTS idx_address_log_change_reason;
DROP INDEX IF EXISTS idx_address_log_confidence_after;
DROP INDEX IF EXISTS idx_address_log_old_location;
DROP INDEX IF EXISTS idx_address_log_new_location;

-- Удаляем таблицу
DROP TABLE IF EXISTS address_change_log;