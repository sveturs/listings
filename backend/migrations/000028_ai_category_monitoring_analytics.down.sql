-- Откат миграции 000028_ai_category_monitoring_analytics

-- Удаление триггера
DROP TRIGGER IF EXISTS trigger_update_detection_stats ON category_detection_feedback;

-- Удаление функций
DROP FUNCTION IF EXISTS update_detection_stats();
DROP FUNCTION IF EXISTS check_algorithm_performance();
DROP FUNCTION IF EXISTS get_realtime_accuracy(INTEGER);

-- Удаление views
DROP VIEW IF EXISTS category_ai_mapping_performance;
DROP VIEW IF EXISTS category_detection_top_errors;
DROP VIEW IF EXISTS category_detection_by_category;
DROP VIEW IF EXISTS category_detection_daily_accuracy;

-- Удаление индексов
DROP INDEX IF EXISTS idx_feedback_algorithm;
DROP INDEX IF EXISTS idx_feedback_category;
DROP INDEX IF EXISTS idx_feedback_confirmed;
DROP INDEX IF EXISTS idx_feedback_created;

-- Удаление таблиц
DROP TABLE IF EXISTS category_detection_stats;
DROP TABLE IF EXISTS category_detection_experiments;