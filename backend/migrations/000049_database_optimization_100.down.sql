-- Откат оптимизации базы данных
-- Migration 000049 DOWN: Remove database optimizations

-- Удаление таблиц производительности
DROP TABLE IF EXISTS performance_metrics;
DROP TABLE IF EXISTS query_cache;

-- Удаление материализованных представлений
DROP MATERIALIZED VIEW IF EXISTS mv_category_statistics;

-- Удаление функций
DROP FUNCTION IF EXISTS refresh_category_statistics();
DROP FUNCTION IF EXISTS create_monthly_partitions(text, date, date);

-- Удаление индексов (в обратном порядке)
DROP INDEX IF EXISTS idx_performance_metrics_name_time;
DROP INDEX IF EXISTS idx_query_cache_expires;
DROP INDEX IF EXISTS idx_query_cache_hash;
DROP INDEX IF EXISTS idx_translations_composite;
DROP INDEX IF EXISTS idx_unified_cat_attrs_composite;
DROP INDEX IF EXISTS idx_unified_attributes_name_trgm;
DROP INDEX IF EXISTS idx_marketplace_listings_description_trgm;
DROP INDEX IF EXISTS idx_marketplace_listings_title_trgm;
DROP INDEX IF EXISTS idx_marketplace_listings_created_at_desc;
DROP INDEX IF EXISTS idx_marketplace_listings_user_status;
DROP INDEX IF EXISTS idx_marketplace_listings_status_category;
DROP INDEX IF EXISTS idx_translations_entity_lang_field;
DROP INDEX IF EXISTS idx_unified_attributes_active_filterable;
DROP INDEX IF EXISTS idx_unified_attributes_active_searchable;

-- Сброс настроек автовакуума
ALTER TABLE marketplace_listings RESET (
    autovacuum_vacuum_threshold,
    autovacuum_analyze_threshold,
    autovacuum_vacuum_scale_factor,
    autovacuum_analyze_scale_factor
);

ALTER TABLE unified_attributes RESET (
    autovacuum_vacuum_threshold,
    autovacuum_analyze_threshold
);