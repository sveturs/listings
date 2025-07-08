-- Удаляем триггер
DROP TRIGGER IF EXISTS trigger_refresh_rating_distributions ON reviews;

-- Удаляем функцию
DROP FUNCTION IF EXISTS refresh_rating_distributions();

-- Удаляем материализованные представления
DROP MATERIALIZED VIEW IF EXISTS storefront_rating_distribution;
DROP MATERIALIZED VIEW IF EXISTS user_rating_distribution;