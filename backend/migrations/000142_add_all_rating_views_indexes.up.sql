-- Добавление индексов для всех materialized views, связанных с ratings
-- Необходимо для поддержки REFRESH MATERIALIZED VIEW CONCURRENTLY

-- storefront_rating_distribution
CREATE UNIQUE INDEX IF NOT EXISTS storefront_rating_distribution_storefront_id_idx 
ON storefront_rating_distribution (storefront_id);

-- storefront_rating_summary
CREATE UNIQUE INDEX IF NOT EXISTS storefront_rating_summary_storefront_id_idx 
ON storefront_rating_summary (storefront_id);

-- user_rating_distribution
CREATE UNIQUE INDEX IF NOT EXISTS user_rating_distribution_user_id_idx 
ON user_rating_distribution (user_id);