-- Удаление таблицы search_behavior_metrics и связанных объектов

DROP TRIGGER IF EXISTS trigger_update_search_behavior_metrics_updated_at ON search_behavior_metrics;
DROP FUNCTION IF EXISTS update_search_behavior_metrics_updated_at();
DROP TABLE IF EXISTS search_behavior_metrics CASCADE;