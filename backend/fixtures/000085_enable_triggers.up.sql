-- Включаем триггеры обратно
SET session_replication_role = 'origin';

-- Обновляем материализованные представления после загрузки данных
REFRESH MATERIALIZED VIEW category_listing_counts;
REFRESH MATERIALIZED VIEW user_rating_summary;
REFRESH MATERIALIZED VIEW user_ratings;
