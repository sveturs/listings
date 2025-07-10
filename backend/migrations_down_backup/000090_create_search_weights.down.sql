-- Откат миграции search_weights

-- Удаление триггеров
DROP TRIGGER IF EXISTS trigger_log_search_weight_changes ON search_weights;
DROP TRIGGER IF EXISTS trigger_update_search_weights_updated_at ON search_weights;

-- Удаление функций
DROP FUNCTION IF EXISTS log_search_weight_changes();
DROP FUNCTION IF EXISTS update_search_weights_updated_at();

-- Удаление таблиц
DROP TABLE IF EXISTS search_weights_history;
DROP TABLE IF EXISTS search_weights;