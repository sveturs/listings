-- Добавляем колонку для хранения типа объекта, чтобы различать товары, продавцов и витрины
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS entity_origin_type VARCHAR(50);
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS entity_origin_id INT;

-- Создаем индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_reviews_entity_origin ON reviews(entity_origin_type, entity_origin_id);

-- Создаем представление для агрегации рейтингов пользователей
CREATE MATERIALIZED VIEW user_rating_summary AS
WITH review_stats AS (
    SELECT 
        COALESCE(entity_origin_id, user_id) as user_id,
        COUNT(*) as total_reviews,
        AVG(rating) as average_rating,
        COUNT(*) FILTER (WHERE rating = 1) as rating_1,
        COUNT(*) FILTER (WHERE rating = 2) as rating_2,
        COUNT(*) FILTER (WHERE rating = 3) as rating_3,
        COUNT(*) FILTER (WHERE rating = 4) as rating_4,
        COUNT(*) FILTER (WHERE rating = 5) as rating_5
    FROM reviews
    WHERE (entity_type = 'listing' AND entity_origin_type IS NULL) 
       OR (entity_origin_type = 'user')
    GROUP BY COALESCE(entity_origin_id, user_id)
)
SELECT 
    u.id as user_id,
    u.name,
    rs.total_reviews,
    rs.average_rating,
    rs.rating_1,
    rs.rating_2,
    rs.rating_3,
    rs.rating_4,
    rs.rating_5
FROM users u
LEFT JOIN review_stats rs ON u.id = rs.user_id;

-- Создаем уникальный индекс для обновления
CREATE UNIQUE INDEX IF NOT EXISTS user_rating_summary_idx ON user_rating_summary(user_id);

-- Создаем представление для агрегации рейтингов витрин
CREATE MATERIALIZED VIEW storefront_rating_summary AS
WITH review_stats AS (
    SELECT 
        COALESCE(entity_origin_id, storefront_id) as storefront_id,
        COUNT(*) as total_reviews,
        AVG(rating) as average_rating,
        COUNT(*) FILTER (WHERE rating = 1) as rating_1,
        COUNT(*) FILTER (WHERE rating = 2) as rating_2,
        COUNT(*) FILTER (WHERE rating = 3) as rating_3,
        COUNT(*) FILTER (WHERE rating = 4) as rating_4,
        COUNT(*) FILTER (WHERE rating = 5) as rating_5
    FROM reviews r
    JOIN marketplace_listings ml ON r.entity_id = ml.id
    WHERE (r.entity_type = 'listing' AND ml.storefront_id IS NOT NULL AND r.entity_origin_type IS NULL) 
       OR (r.entity_origin_type = 'storefront')
    GROUP BY COALESCE(r.entity_origin_id, ml.storefront_id)
)
SELECT 
    s.id as storefront_id,
    s.name,
    rs.total_reviews,
    rs.average_rating,
    rs.rating_1,
    rs.rating_2,
    rs.rating_3,
    rs.rating_4,
    rs.rating_5
FROM user_storefronts s
LEFT JOIN review_stats rs ON s.id = rs.storefront_id;

-- Создаем уникальный индекс для обновления
CREATE UNIQUE INDEX IF NOT EXISTS storefront_rating_summary_idx ON storefront_rating_summary(storefront_id);

-- Создаем функцию для обновления представлений
CREATE OR REPLACE FUNCTION refresh_rating_summaries()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_summary;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_summary;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для обновления представлений
CREATE TRIGGER refresh_rating_summaries_trigger
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_rating_summaries();

-- Функция для сохранения информации о происхождении отзыва при удалении объявления
CREATE OR REPLACE FUNCTION preserve_review_origin()
RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем отзывы удаляемого товара, сохраняя информацию о продавце
    UPDATE reviews
    SET entity_origin_type = 'user', 
        entity_origin_id = OLD.user_id
    WHERE entity_type = 'listing' 
      AND entity_id = OLD.id 
      AND entity_origin_type IS NULL;
    
    -- Если товар из витрины, добавляем и эту информацию
    IF OLD.storefront_id IS NOT NULL THEN
        UPDATE reviews
        SET entity_origin_type = 'storefront',
            entity_origin_id = OLD.storefront_id
        WHERE entity_type = 'listing' 
          AND entity_id = OLD.id
          AND entity_origin_type = 'user';
    END IF;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для сохранения происхождения отзыва
CREATE TRIGGER preserve_review_origin_trigger
BEFORE DELETE ON marketplace_listings
FOR EACH ROW
EXECUTE FUNCTION preserve_review_origin();