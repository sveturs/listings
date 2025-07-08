-- Добавляем уникальные индексы для CONCURRENTLY обновления
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_rating_distribution_unique ON user_rating_distribution(user_id, rating);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storefront_rating_distribution_unique ON storefront_rating_distribution(storefront_id, rating);