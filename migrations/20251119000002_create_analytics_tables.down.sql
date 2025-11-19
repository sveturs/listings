-- =====================================================
-- Migration: 20251119000002_create_analytics_tables.down.sql
-- Description: Rollback analytics infrastructure
-- =====================================================

-- Drop functions
DROP FUNCTION IF EXISTS archive_old_analytics_events();
DROP FUNCTION IF EXISTS refresh_analytics_views();
DROP FUNCTION IF EXISTS log_analytics_event(VARCHAR, VARCHAR, BIGINT, BIGINT, VARCHAR, JSONB);

-- Drop materialized views
DROP MATERIALIZED VIEW IF EXISTS analytics_listing_stats;
DROP MATERIALIZED VIEW IF EXISTS analytics_overview_daily;

-- Drop main table (cascade will remove all indexes)
DROP TABLE IF EXISTS analytics_events CASCADE;
