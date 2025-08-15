-- Включаем триггеры обратно
SET session_replication_role = 'origin';

-- Обновляем материализованные представления после загрузки данных
REFRESH MATERIALIZED VIEW user_ratings;
REFRESH MATERIALIZED VIEW storefront_ratings;
