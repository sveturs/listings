-- Удаляем триггер
DROP TRIGGER IF EXISTS update_ratings_after_review_change ON reviews;

-- Удаляем функции
DROP FUNCTION IF EXISTS refresh_rating_views();
DROP FUNCTION IF EXISTS rebuild_all_ratings();

-- Удаляем таблицу кеша
DROP TABLE IF EXISTS rating_cache;

-- Удаляем материализованные представления
DROP MATERIALIZED VIEW IF EXISTS storefront_ratings;
DROP MATERIALIZED VIEW IF EXISTS user_ratings;