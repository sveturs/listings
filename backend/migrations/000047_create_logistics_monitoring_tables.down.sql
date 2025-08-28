-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_logistics_metrics_updated_at ON logistics_metrics;
DROP TRIGGER IF EXISTS update_problem_shipments_updated_at ON problem_shipments;
DROP TRIGGER IF EXISTS update_logistics_monitoring_settings_updated_at ON logistics_monitoring_settings;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_logistics_updated_at();

-- Удаляем таблицы в обратном порядке зависимостей
DROP TABLE IF EXISTS logistics_dashboard_cache;
DROP TABLE IF EXISTS logistics_monitoring_settings;
DROP TABLE IF EXISTS logistics_admin_logs;
DROP TABLE IF EXISTS problem_shipments;
DROP TABLE IF EXISTS logistics_metrics;