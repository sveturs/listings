-- Удаление индексов для materialized views

DROP INDEX IF EXISTS storefront_ratings_storefront_id_idx;
DROP INDEX IF EXISTS user_ratings_user_id_idx;