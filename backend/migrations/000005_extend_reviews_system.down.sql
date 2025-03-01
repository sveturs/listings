-- Удаляем триггеры
DROP TRIGGER IF EXISTS preserve_review_origin_trigger ON marketplace_listings;
DROP TRIGGER IF EXISTS refresh_rating_summaries_trigger ON reviews;

-- Удаляем функции
DROP FUNCTION IF EXISTS preserve_review_origin();
DROP FUNCTION IF EXISTS refresh_rating_summaries();

-- Удаляем материализованные представления
DROP MATERIALIZED VIEW IF EXISTS storefront_rating_summary;
DROP MATERIALIZED VIEW IF EXISTS user_rating_summary;

-- Удаляем добавленные поля и индексы
DROP INDEX IF EXISTS idx_reviews_entity_origin;
ALTER TABLE reviews DROP COLUMN IF EXISTS entity_origin_type;
ALTER TABLE reviews DROP COLUMN IF EXISTS entity_origin_id;