-- Удаление таблиц и связанных объектов
DROP TRIGGER IF EXISTS update_search_trending_queries_updated_at_trigger ON search_trending_queries;
DROP TRIGGER IF EXISTS update_search_analytics_updated_at_trigger ON search_analytics;
DROP FUNCTION IF EXISTS update_search_analytics_updated_at();

DROP TABLE IF EXISTS search_trending_queries;
DROP TABLE IF EXISTS search_result_clicks;
DROP TABLE IF EXISTS search_analytics;
DROP TABLE IF EXISTS search_logs;