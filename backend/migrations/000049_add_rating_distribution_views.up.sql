-- Создаем представление для распределения рейтингов пользователей
CREATE MATERIALIZED VIEW IF NOT EXISTS user_rating_distribution AS
SELECT 
    entity_origin_id as user_id,
    rating,
    COUNT(*) as count
FROM reviews
WHERE entity_origin_type = 'user' 
    AND status = 'published'
GROUP BY entity_origin_id, rating;

-- Создаем индекс для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_user_rating_distribution_user_id ON user_rating_distribution(user_id);

-- Создаем представление для распределения рейтингов витрин
CREATE MATERIALIZED VIEW IF NOT EXISTS storefront_rating_distribution AS
SELECT 
    entity_origin_id as storefront_id,
    rating,
    COUNT(*) as count
FROM reviews
WHERE entity_origin_type = 'storefront' 
    AND status = 'published'
GROUP BY entity_origin_id, rating;

-- Создаем индекс для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_storefront_rating_distribution_storefront_id ON storefront_rating_distribution(storefront_id);

-- Добавляем автоматическое обновление при изменении отзывов
CREATE OR REPLACE FUNCTION refresh_rating_distributions() RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_distribution;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_distribution;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Триггер для обновления при изменении отзывов
DROP TRIGGER IF EXISTS trigger_refresh_rating_distributions ON reviews;
CREATE TRIGGER trigger_refresh_rating_distributions
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_rating_distributions();