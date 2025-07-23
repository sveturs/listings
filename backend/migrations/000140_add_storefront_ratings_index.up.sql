-- Добавление уникального индекса для materialized view storefront_ratings
-- Необходимо для поддержки REFRESH MATERIALIZED VIEW CONCURRENTLY

-- Создаем уникальный индекс на storefront_id
CREATE UNIQUE INDEX IF NOT EXISTS storefront_ratings_storefront_id_idx 
ON storefront_ratings (storefront_id);

-- Аналогично для user_ratings, если он тоже используется
CREATE UNIQUE INDEX IF NOT EXISTS user_ratings_user_id_idx 
ON user_ratings (user_id);