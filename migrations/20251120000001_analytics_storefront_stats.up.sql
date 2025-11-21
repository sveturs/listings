-- =====================================================
-- Migration: 20251120000001_analytics_storefront_stats.up.sql
-- Description: Create materialized view for storefront analytics stats
-- Author: Phase 30.1 - Storefront Performance Analytics
-- Date: 2025-11-20
-- =====================================================

-- =====================================================
-- Materialized View: Storefront Performance Stats
-- =====================================================
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_storefront_stats AS
SELECT
    s.id AS storefront_id,
    s.name AS storefront_name,
    s.user_id AS owner_id,

    -- Listings metrics (all time)
    COUNT(DISTINCT l.id) FILTER (WHERE l.status = 'active' AND l.is_deleted = false) AS active_listings,
    COUNT(DISTINCT l.id) FILTER (WHERE l.is_deleted = false) AS total_listings,

    -- Sales metrics (completed orders only)
    COUNT(DISTINCT o.id) FILTER (WHERE o.status = 'delivered') AS total_sales,
    COALESCE(SUM(o.total) FILTER (WHERE o.status = 'delivered'), 0) AS total_revenue,
    CASE
        WHEN COUNT(DISTINCT o.id) FILTER (WHERE o.status = 'delivered') > 0
        THEN COALESCE(SUM(o.total) FILTER (WHERE o.status = 'delivered'), 0) /
             COUNT(DISTINCT o.id) FILTER (WHERE o.status = 'delivered')
        ELSE 0
    END AS average_order_value,

    -- Engagement metrics (from analytics_events)
    COALESCE(COUNT(ae.id) FILTER (
        WHERE ae.event_type = 'view'
        AND ae.entity_type = 'listing'
        AND ae.entity_id IN (SELECT id FROM listings WHERE storefront_id = s.id AND is_deleted = false)
    ), 0) AS total_views,

    -- Favorites count (from listings)
    COALESCE(SUM(l.favorites_count), 0) AS total_favorites,

    -- Conversion rate calculation
    CASE
        WHEN COUNT(ae.id) FILTER (
            WHERE ae.event_type = 'view'
            AND ae.entity_type = 'listing'
            AND ae.entity_id IN (SELECT id FROM listings WHERE storefront_id = s.id AND is_deleted = false)
        ) > 0
        THEN (COUNT(DISTINCT o.id) FILTER (WHERE o.status = 'delivered')::DECIMAL /
              COUNT(ae.id) FILTER (
                  WHERE ae.event_type = 'view'
                  AND ae.entity_type = 'listing'
                  AND ae.entity_id IN (SELECT id FROM listings WHERE storefront_id = s.id AND is_deleted = false)
              )) * 100
        ELSE 0
    END AS conversion_rate,

    -- Timestamps for cache invalidation
    GREATEST(
        MAX(s.updated_at),
        MAX(l.updated_at),
        MAX(o.updated_at),
        MAX(ae.created_at)
    ) AS last_updated_at

FROM storefronts s
LEFT JOIN listings l ON l.storefront_id = s.id AND l.is_deleted = false
LEFT JOIN orders o ON o.storefront_id = s.id
LEFT JOIN analytics_events ae ON ae.entity_id = l.id
    AND ae.entity_type = 'listing'
    AND ae.event_type = 'view'

WHERE s.deleted_at IS NULL

GROUP BY s.id, s.name, s.user_id;

-- =====================================================
-- Indexes for Materialized View
-- =====================================================

-- Unique index on storefront_id for CONCURRENTLY refresh
CREATE UNIQUE INDEX idx_analytics_storefront_stats_storefront_id
    ON analytics_storefront_stats(storefront_id);

-- Index for owner lookups
CREATE INDEX idx_analytics_storefront_stats_owner
    ON analytics_storefront_stats(owner_id);

-- Index for sorting by revenue
CREATE INDEX idx_analytics_storefront_stats_revenue
    ON analytics_storefront_stats(total_revenue DESC);

-- Index for sorting by sales count
CREATE INDEX idx_analytics_storefront_stats_sales
    ON analytics_storefront_stats(total_sales DESC);

-- =====================================================
-- Refresh Function
-- =====================================================

-- Function to refresh storefront stats view
CREATE OR REPLACE FUNCTION refresh_analytics_storefront_stats()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_storefront_stats;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Update existing refresh function
-- =====================================================

-- Update the main analytics refresh function to include storefront stats
CREATE OR REPLACE FUNCTION refresh_analytics_views()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_overview_daily;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_listing_stats;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_storefront_stats;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Comments for Documentation
-- =====================================================

COMMENT ON MATERIALIZED VIEW analytics_storefront_stats IS
'Pre-aggregated storefront performance metrics (refresh every 15 minutes)';

COMMENT ON FUNCTION refresh_analytics_storefront_stats IS
'Refresh storefront stats materialized view concurrently';
