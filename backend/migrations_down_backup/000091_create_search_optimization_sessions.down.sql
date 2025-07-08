-- Откат миграции search_optimization_sessions

-- Удаление триггера и функции
DROP TRIGGER IF EXISTS trigger_update_search_optimization_sessions_updated_at ON search_optimization_sessions;
DROP FUNCTION IF EXISTS update_search_optimization_sessions_updated_at();

-- Удаление таблицы
DROP TABLE IF EXISTS search_optimization_sessions;