-- Удаление таблиц аналитики
DROP TRIGGER IF EXISTS update_storefront_analytics_updated_at ON storefront_analytics;
DROP FUNCTION IF EXISTS update_storefront_analytics_updated_at();
DROP TABLE IF EXISTS storefront_events;
DROP TABLE IF EXISTS storefront_analytics;