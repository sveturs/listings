-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS trigger_refresh_rating_distributions ON reviews;
DROP FUNCTION IF EXISTS refresh_rating_distributions();

-- Удаляем все материализованные представления
DROP MATERIALIZED VIEW IF EXISTS user_rating_distribution CASCADE;
DROP MATERIALIZED VIEW IF EXISTS storefront_rating_distribution CASCADE;
DROP MATERIALIZED VIEW IF EXISTS user_ratings CASCADE;
DROP MATERIALIZED VIEW IF EXISTS storefront_ratings CASCADE;