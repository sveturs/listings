-- =====================================================
-- Migration: 20251120000001_analytics_storefront_stats.down.sql
-- Description: Rollback storefront analytics stats
-- Author: Phase 30.1 - Storefront Performance Analytics
-- Date: 2025-11-20
-- =====================================================

-- Drop refresh function for storefront stats
DROP FUNCTION IF EXISTS refresh_analytics_storefront_stats();

-- Restore original refresh_analytics_views function (without storefront stats)
CREATE OR REPLACE FUNCTION refresh_analytics_views()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_overview_daily;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_listing_stats;
END;
$$ LANGUAGE plpgsql;

-- Drop materialized view
DROP MATERIALIZED VIEW IF EXISTS analytics_storefront_stats CASCADE;
