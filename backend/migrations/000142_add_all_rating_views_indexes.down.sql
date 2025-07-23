-- Удаление индексов для materialized views

DROP INDEX IF EXISTS storefront_rating_distribution_storefront_id_idx;
DROP INDEX IF EXISTS storefront_rating_summary_storefront_id_idx;
DROP INDEX IF EXISTS user_rating_distribution_user_id_idx;