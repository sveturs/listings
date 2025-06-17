-- Создаем материализованное представление для рейтингов пользователей
CREATE MATERIALIZED VIEW IF NOT EXISTS user_ratings AS
SELECT 
    u.id as user_id,
    COUNT(DISTINCT r.id) as total_reviews,
    COALESCE(AVG(r.rating), 0) as average_rating,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'user' THEN r.id END) as direct_reviews,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'listing' THEN r.id END) as listing_reviews,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'storefront' THEN r.id END) as storefront_reviews,
    COUNT(DISTINCT CASE WHEN r.is_verified_purchase THEN r.id END) as verified_reviews,
    -- Распределение оценок
    COUNT(CASE WHEN r.rating = 1 THEN 1 END) as rating_1,
    COUNT(CASE WHEN r.rating = 2 THEN 1 END) as rating_2,
    COUNT(CASE WHEN r.rating = 3 THEN 1 END) as rating_3,
    COUNT(CASE WHEN r.rating = 4 THEN 1 END) as rating_4,
    COUNT(CASE WHEN r.rating = 5 THEN 1 END) as rating_5,
    -- Тренд за последние 30 дней
    AVG(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN r.rating END) as recent_rating,
    COUNT(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN 1 END) as recent_reviews,
    MAX(r.created_at) as last_review_at
FROM users u
LEFT JOIN reviews r ON (
    -- Прямые отзывы на пользователя
    (r.entity_type = 'user' AND r.entity_id = u.id) OR
    -- Отзывы через origin после удаления
    (r.entity_origin_type = 'user' AND r.entity_origin_id = u.id)
) AND r.status = 'published'
GROUP BY u.id;

-- Создаем индексы для материализованного представления
CREATE UNIQUE INDEX idx_user_ratings_user_id ON user_ratings(user_id);
CREATE INDEX idx_user_ratings_average ON user_ratings(average_rating DESC);
CREATE INDEX idx_user_ratings_total ON user_ratings(total_reviews DESC);

-- Создаем материализованное представление для рейтингов магазинов
CREATE MATERIALIZED VIEW IF NOT EXISTS storefront_ratings AS
SELECT 
    s.id as storefront_id,
    COUNT(DISTINCT r.id) as total_reviews,
    COALESCE(AVG(r.rating), 0) as average_rating,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'storefront' THEN r.id END) as direct_reviews,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'listing' THEN r.id END) as listing_reviews,
    COUNT(DISTINCT CASE WHEN r.is_verified_purchase THEN r.id END) as verified_reviews,
    -- Распределение оценок
    COUNT(CASE WHEN r.rating = 1 THEN 1 END) as rating_1,
    COUNT(CASE WHEN r.rating = 2 THEN 1 END) as rating_2,
    COUNT(CASE WHEN r.rating = 3 THEN 1 END) as rating_3,
    COUNT(CASE WHEN r.rating = 4 THEN 1 END) as rating_4,
    COUNT(CASE WHEN r.rating = 5 THEN 1 END) as rating_5,
    -- Тренд за последние 30 дней
    AVG(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN r.rating END) as recent_rating,
    COUNT(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN 1 END) as recent_reviews,
    MAX(r.created_at) as last_review_at,
    -- Информация о владельце
    s.user_id as owner_id
FROM user_storefronts s
LEFT JOIN reviews r ON (
    -- Прямые отзывы на магазин
    (r.entity_type = 'storefront' AND r.entity_id = s.id) OR
    -- Отзывы через origin после удаления
    (r.entity_origin_type = 'storefront' AND r.entity_origin_id = s.id)
) AND r.status = 'published'
GROUP BY s.id, s.user_id;

-- Создаем индексы для материализованного представления магазинов
CREATE UNIQUE INDEX idx_storefront_ratings_id ON storefront_ratings(storefront_id);
CREATE INDEX idx_storefront_ratings_average ON storefront_ratings(average_rating DESC);
CREATE INDEX idx_storefront_ratings_owner ON storefront_ratings(owner_id);

-- Создаем функцию для обновления материализованных представлений
CREATE OR REPLACE FUNCTION refresh_rating_views() RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем только затронутые строки, а не всё представление
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- Для пользователей
        IF NEW.entity_origin_type = 'user' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;
        
        -- Для магазинов
        IF NEW.entity_origin_type = 'storefront' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении также обновляем
        IF OLD.entity_origin_type = 'user' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;
        
        IF OLD.entity_origin_type = 'storefront' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для автоматического обновления представлений
CREATE TRIGGER update_ratings_after_review_change
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION refresh_rating_views();

-- Создаем функцию для полного пересчета рейтингов (для cron job)
CREATE OR REPLACE FUNCTION rebuild_all_ratings() RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
END;
$$ LANGUAGE plpgsql;

-- Создаем таблицу для кеширования агрегированных рейтингов
CREATE TABLE IF NOT EXISTS rating_cache (
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    average_rating DECIMAL(3,2),
    total_reviews INTEGER DEFAULT 0,
    distribution JSONB,
    breakdown JSONB,
    verified_percentage INTEGER DEFAULT 0,
    recent_trend VARCHAR(10) CHECK (recent_trend IN ('up', 'down', 'stable')),
    calculated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (entity_type, entity_id)
);

-- Создаем индексы для быстрого доступа
CREATE INDEX idx_rating_cache_type_rating ON rating_cache(entity_type, average_rating DESC);
CREATE INDEX idx_rating_cache_updated ON rating_cache(calculated_at);

-- Комментарии для документации
COMMENT ON MATERIALIZED VIEW user_ratings IS 'Агрегированные рейтинги пользователей с учетом всех типов отзывов';
COMMENT ON MATERIALIZED VIEW storefront_ratings IS 'Агрегированные рейтинги магазинов с учетом отзывов на товары';
COMMENT ON TABLE rating_cache IS 'Кеш для быстрого доступа к агрегированным рейтингам';
COMMENT ON FUNCTION refresh_rating_views IS 'Обновляет материализованные представления рейтингов';
COMMENT ON FUNCTION rebuild_all_ratings IS 'Полный пересчет всех рейтингов (для планировщика)';