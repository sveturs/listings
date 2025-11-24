-- Rollback: Drop trending analytics materialized view and related objects

-- Drop function
DROP FUNCTION IF EXISTS refresh_analytics_trending_cache();

-- Drop index
DROP INDEX IF EXISTS idx_analytics_trending_cache_singleton;

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS analytics_trending_cache;
