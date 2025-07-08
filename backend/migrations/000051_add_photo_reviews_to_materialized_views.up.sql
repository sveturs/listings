-- Удаляем старые представления
DROP MATERIALIZED VIEW IF EXISTS user_rating_distribution CASCADE;
DROP MATERIALIZED VIEW IF EXISTS storefront_rating_distribution CASCADE;
DROP MATERIALIZED VIEW IF EXISTS user_ratings CASCADE;
DROP MATERIALIZED VIEW IF EXISTS storefront_ratings CASCADE;

-- Пересоздаем user_ratings с полем photo_reviews
CREATE MATERIALIZED VIEW user_ratings AS
SELECT 
    entity_origin_id as user_id,
    COUNT(*) as total_reviews,
    AVG(rating) as average_rating,
    COUNT(*) FILTER (WHERE entity_type = 'user') as direct_reviews,
    COUNT(*) FILTER (WHERE entity_type = 'listing') as listing_reviews,
    COUNT(*) FILTER (WHERE entity_type = 'storefront') as storefront_reviews,
    COUNT(*) FILTER (WHERE is_verified_purchase = true) as verified_reviews,
    COUNT(*) FILTER (WHERE array_length(photos, 1) > 0) as photo_reviews,
    COUNT(*) FILTER (WHERE rating = 1) as rating_1,
    COUNT(*) FILTER (WHERE rating = 2) as rating_2,
    COUNT(*) FILTER (WHERE rating = 3) as rating_3,
    COUNT(*) FILTER (WHERE rating = 4) as rating_4,
    COUNT(*) FILTER (WHERE rating = 5) as rating_5,
    AVG(rating) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as recent_rating,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as recent_reviews,
    MAX(created_at) as last_review_at
FROM reviews
WHERE entity_origin_type = 'user' 
    AND status = 'published'
GROUP BY entity_origin_id;

-- Пересоздаем storefront_ratings с полем photo_reviews
CREATE MATERIALIZED VIEW storefront_ratings AS
SELECT 
    entity_origin_id as storefront_id,
    COUNT(*) as total_reviews,
    AVG(rating) as average_rating,
    COUNT(*) FILTER (WHERE entity_type = 'storefront') as direct_reviews,
    COUNT(*) FILTER (WHERE entity_type = 'listing') as listing_reviews,
    COUNT(*) FILTER (WHERE is_verified_purchase = true) as verified_reviews,
    COUNT(*) FILTER (WHERE array_length(photos, 1) > 0) as photo_reviews,
    COUNT(*) FILTER (WHERE rating = 1) as rating_1,
    COUNT(*) FILTER (WHERE rating = 2) as rating_2,
    COUNT(*) FILTER (WHERE rating = 3) as rating_3,
    COUNT(*) FILTER (WHERE rating = 4) as rating_4,
    COUNT(*) FILTER (WHERE rating = 5) as rating_5,
    AVG(rating) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as recent_rating,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as recent_reviews,
    MAX(created_at) as last_review_at
FROM reviews
WHERE entity_origin_type = 'storefront' 
    AND status = 'published'
GROUP BY entity_origin_id;

-- Пересоздаем индексы для user_ratings
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_ratings_user_id ON user_ratings(user_id);
CREATE INDEX IF NOT EXISTS idx_user_ratings_average ON user_ratings(average_rating DESC);
CREATE INDEX IF NOT EXISTS idx_user_ratings_total ON user_ratings(total_reviews DESC);

-- Пересоздаем индексы для storefront_ratings
CREATE UNIQUE INDEX IF NOT EXISTS idx_storefront_ratings_storefront_id ON storefront_ratings(storefront_id);
CREATE INDEX IF NOT EXISTS idx_storefront_ratings_average ON storefront_ratings(average_rating DESC);
CREATE INDEX IF NOT EXISTS idx_storefront_ratings_total ON storefront_ratings(total_reviews DESC);

-- Пересоздаем представления для распределения рейтингов
CREATE MATERIALIZED VIEW user_rating_distribution AS
SELECT 
    entity_origin_id as user_id,
    rating,
    COUNT(*) as count
FROM reviews
WHERE entity_origin_type = 'user' 
    AND status = 'published'
GROUP BY entity_origin_id, rating;

CREATE MATERIALIZED VIEW storefront_rating_distribution AS
SELECT 
    entity_origin_id as storefront_id,
    rating,
    COUNT(*) as count
FROM reviews
WHERE entity_origin_type = 'storefront' 
    AND status = 'published'
GROUP BY entity_origin_id, rating;

-- Пересоздаем индексы для распределений
CREATE INDEX IF NOT EXISTS idx_user_rating_distribution_user_id ON user_rating_distribution(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_rating_distribution_unique ON user_rating_distribution(user_id, rating);

CREATE INDEX IF NOT EXISTS idx_storefront_rating_distribution_storefront_id ON storefront_rating_distribution(storefront_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storefront_rating_distribution_unique ON storefront_rating_distribution(storefront_id, rating);

-- Удаляем старый триггер если существует
DROP TRIGGER IF EXISTS trigger_refresh_rating_distributions ON reviews;
DROP FUNCTION IF EXISTS refresh_rating_distributions();

-- Пересоздаем функцию для обновления всех представлений
CREATE OR REPLACE FUNCTION refresh_rating_distributions() RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем все материализованные представления
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_distribution;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_distribution;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Пересоздаем триггер
CREATE TRIGGER trigger_refresh_rating_distributions
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_rating_distributions();