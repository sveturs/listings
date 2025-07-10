-- Удаление таблицы item_performance_metrics и связанных объектов

DROP TRIGGER IF EXISTS trigger_update_item_performance_metrics_updated_at ON item_performance_metrics;
DROP FUNCTION IF EXISTS update_item_performance_metrics_updated_at();
DROP TABLE IF EXISTS item_performance_metrics CASCADE;