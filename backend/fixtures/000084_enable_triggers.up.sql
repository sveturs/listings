-- Включаем триггеры обратно
SET session_replication_role = 'origin';

-- Обновляем материализованные представления после загрузки данных
REFRESH MATERIALIZED VIEW category_listing_counts;
REFRESH MATERIALIZED VIEW gis_listing_density_grid;
REFRESH MATERIALIZED VIEW mv_category_statistics;
REFRESH MATERIALIZED VIEW mv_popular_category_attributes;
REFRESH MATERIALIZED VIEW storefront_rating_distribution;
REFRESH MATERIALIZED VIEW storefront_rating_summary;
REFRESH MATERIALIZED VIEW storefront_ratings;
-- DEPRECATED: user_* materialized views removed after migration to auth-service
-- REFRESH MATERIALIZED VIEW user_rating_distribution;
-- REFRESH MATERIALIZED VIEW user_rating_summary;
-- REFRESH MATERIALIZED VIEW user_ratings;