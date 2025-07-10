-- Удаление устаревших таблиц search_logs и связанных с ними структур
-- Эти таблицы заменены на user_behavior_events для единообразного трекинга

-- Удаляем триггеры
DROP TRIGGER IF EXISTS update_search_analytics_updated_at_trigger ON search_analytics;
DROP TRIGGER IF EXISTS update_search_trending_queries_updated_at_trigger ON search_trending_queries;

-- Удаляем функцию
DROP FUNCTION IF EXISTS update_search_analytics_updated_at();

-- Удаляем таблицы с зависимостями (в правильном порядке)
DROP TABLE IF EXISTS search_result_clicks CASCADE;
DROP TABLE IF EXISTS search_trending_queries CASCADE; 
DROP TABLE IF EXISTS search_analytics CASCADE;
DROP TABLE IF EXISTS search_logs CASCADE;

-- Комментарий для логов
-- search_logs заменена на user_behavior_events для единого трекинга всех событий
-- search_analytics заменена на search_behavior_metrics с агрегацией из user_behavior_events
-- search_result_clicks - клики теперь отслеживаются через result_clicked в user_behavior_events
-- search_trending_queries - трендовые запросы рассчитываются динамически из user_behavior_events