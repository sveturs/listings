-- =====================================================
-- Migration: 20251119000002_create_analytics_tables.up.sql
-- Description: Analytics events tracking with time-series optimization
-- Author: Phase 29 - Advanced Analytics Implementation
-- Date: 2025-11-19
-- =====================================================

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- Core Analytics Events Table
-- =====================================================
CREATE TABLE IF NOT EXISTS analytics_events (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,

    -- Event classification
    event_type VARCHAR(50) NOT NULL, -- 'view', 'favorite', 'search', 'order_created', 'order_completed'
    entity_type VARCHAR(50) NOT NULL, -- 'listing', 'order', 'search'
    entity_id BIGINT NOT NULL,

    -- User tracking
    user_id BIGINT,
    session_id VARCHAR(100),
    ip_address INET,

    -- Additional context
    metadata JSONB DEFAULT '{}'::jsonb,

    -- Time-series optimization
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
    -- Note: date_partition removed - use date_trunc in queries instead for partitioning
);

-- =====================================================
-- Time-Series Optimized Indexes
-- =====================================================

-- BRIN index for time-series range queries (very space-efficient)
CREATE INDEX idx_analytics_events_created_at_brin
    ON analytics_events USING BRIN (created_at)
    WITH (pages_per_range = 128);

-- Composite indexes for common query patterns
CREATE INDEX idx_analytics_events_entity_time
    ON analytics_events (entity_type, entity_id, created_at DESC);

CREATE INDEX idx_analytics_events_type_time
    ON analytics_events (event_type, created_at DESC);

CREATE INDEX idx_analytics_events_user_time
    ON analytics_events (user_id, created_at DESC)
    WHERE user_id IS NOT NULL;

-- Partial indexes for hot queries
CREATE INDEX idx_analytics_events_listing_views
    ON analytics_events (entity_id, created_at DESC)
    WHERE event_type = 'view' AND entity_type = 'listing';

CREATE INDEX idx_analytics_events_orders
    ON analytics_events (entity_id, created_at DESC, metadata)
    WHERE event_type IN ('order_created', 'order_completed');

-- GIN index for JSONB metadata queries
CREATE INDEX idx_analytics_events_metadata
    ON analytics_events USING GIN (metadata jsonb_path_ops);

-- =====================================================
-- Materialized View: Daily Overview Stats
-- =====================================================
CREATE MATERIALIZED VIEW analytics_overview_daily AS
SELECT
    date_trunc('day', created_at)::DATE AS stat_date,

    -- Listings metrics
    COUNT(*) FILTER (WHERE event_type = 'view' AND entity_type = 'listing') AS total_views,
    COUNT(DISTINCT entity_id) FILTER (WHERE event_type = 'view' AND entity_type = 'listing') AS unique_listings_viewed,

    -- Favorites metrics
    COUNT(*) FILTER (WHERE event_type = 'favorite' AND entity_type = 'listing') AS total_favorites,

    -- Search metrics
    COUNT(*) FILTER (WHERE event_type = 'search') AS total_searches,

    -- Orders metrics
    COUNT(*) FILTER (WHERE event_type = 'order_created') AS orders_created,
    COUNT(*) FILTER (WHERE event_type = 'order_completed') AS orders_completed,

    -- Revenue metrics (from metadata.total_amount)
    COALESCE(SUM((metadata->>'total_amount')::DECIMAL) FILTER (
        WHERE event_type = 'order_completed'
    ), 0) AS total_revenue,

    -- User engagement
    COUNT(DISTINCT user_id) FILTER (WHERE user_id IS NOT NULL) AS active_users,
    COUNT(DISTINCT session_id) AS total_sessions

FROM analytics_events
GROUP BY stat_date
ORDER BY stat_date DESC;

-- Indexes for materialized view
CREATE UNIQUE INDEX idx_analytics_overview_daily_date
    ON analytics_overview_daily (stat_date DESC);

-- =====================================================
-- Materialized View: Listing Performance Stats
-- =====================================================
CREATE MATERIALIZED VIEW analytics_listing_stats AS
SELECT
    entity_id AS listing_id,

    -- Engagement metrics (last 30 days)
    COUNT(*) FILTER (
        WHERE event_type = 'view'
        AND created_at >= NOW() - INTERVAL '30 days'
    ) AS views_30d,

    COUNT(DISTINCT user_id) FILTER (
        WHERE event_type = 'view'
        AND created_at >= NOW() - INTERVAL '30 days'
        AND user_id IS NOT NULL
    ) AS unique_viewers_30d,

    COUNT(*) FILTER (
        WHERE event_type = 'favorite'
        AND created_at >= NOW() - INTERVAL '30 days'
    ) AS favorites_30d,

    COUNT(*) FILTER (
        WHERE event_type = 'search'
        AND created_at >= NOW() - INTERVAL '30 days'
    ) AS search_appearances_30d,

    -- Orders metrics (all time)
    COUNT(*) FILTER (WHERE event_type = 'order_created') AS total_orders,
    COUNT(*) FILTER (WHERE event_type = 'order_completed') AS completed_orders,

    -- Revenue metrics
    COALESCE(SUM((metadata->>'total_amount')::DECIMAL) FILTER (
        WHERE event_type = 'order_completed'
    ), 0) AS total_revenue,

    -- Conversion rates (last 30 days)
    CASE
        WHEN COUNT(*) FILTER (WHERE event_type = 'view' AND created_at >= NOW() - INTERVAL '30 days') > 0
        THEN (COUNT(*) FILTER (WHERE event_type = 'favorite' AND created_at >= NOW() - INTERVAL '30 days')::DECIMAL /
              COUNT(*) FILTER (WHERE event_type = 'view' AND created_at >= NOW() - INTERVAL '30 days')) * 100
        ELSE 0
    END AS view_to_favorite_rate,

    CASE
        WHEN COUNT(*) FILTER (WHERE event_type = 'view' AND created_at >= NOW() - INTERVAL '30 days') > 0
        THEN (COUNT(*) FILTER (WHERE event_type = 'order_created' AND created_at >= NOW() - INTERVAL '30 days')::DECIMAL /
              COUNT(*) FILTER (WHERE event_type = 'view' AND created_at >= NOW() - INTERVAL '30 days')) * 100
        ELSE 0
    END AS view_to_order_rate,

    -- Timestamps
    MIN(created_at) AS first_event_at,
    MAX(created_at) AS last_event_at

FROM analytics_events
WHERE entity_type = 'listing'
GROUP BY entity_id;

-- Indexes for listing stats
CREATE UNIQUE INDEX idx_analytics_listing_stats_id
    ON analytics_listing_stats (listing_id);

CREATE INDEX idx_analytics_listing_stats_views
    ON analytics_listing_stats (views_30d DESC);

CREATE INDEX idx_analytics_listing_stats_revenue
    ON analytics_listing_stats (total_revenue DESC);

-- =====================================================
-- Automatic Refresh Functions
-- =====================================================

-- Function to refresh materialized views
CREATE OR REPLACE FUNCTION refresh_analytics_views()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_overview_daily;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_listing_stats;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Helper Functions for Analytics
-- =====================================================

-- Function to log analytics event
CREATE OR REPLACE FUNCTION log_analytics_event(
    p_event_type VARCHAR,
    p_entity_type VARCHAR,
    p_entity_id BIGINT,
    p_user_id BIGINT DEFAULT NULL,
    p_session_id VARCHAR DEFAULT NULL,
    p_metadata JSONB DEFAULT '{}'::jsonb
)
RETURNS UUID AS $$
DECLARE
    v_uuid UUID;
BEGIN
    INSERT INTO analytics_events (
        event_type, entity_type, entity_id,
        user_id, session_id, metadata
    )
    VALUES (
        p_event_type, p_entity_type, p_entity_id,
        p_user_id, p_session_id, p_metadata
    )
    RETURNING uuid INTO v_uuid;

    RETURN v_uuid;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Data Retention Policy (Optional - for production)
-- =====================================================

-- Function to archive old events (keep raw data for 90 days)
CREATE OR REPLACE FUNCTION archive_old_analytics_events()
RETURNS INTEGER AS $$
DECLARE
    v_deleted_count INTEGER;
BEGIN
    -- Delete events older than 90 days (after materialized views are refreshed)
    DELETE FROM analytics_events
    WHERE created_at < NOW() - INTERVAL '90 days'
    RETURNING COUNT(*) INTO v_deleted_count;

    RETURN v_deleted_count;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- Comments for documentation
-- =====================================================

COMMENT ON TABLE analytics_events IS 'Time-series event tracking for analytics with BRIN indexing optimization';
COMMENT ON COLUMN analytics_events.event_type IS 'Type of event: view, favorite, search, order_created, order_completed';
COMMENT ON COLUMN analytics_events.entity_type IS 'Entity being tracked: listing, order, search';
COMMENT ON COLUMN analytics_events.metadata IS 'Additional event context (e.g., order amount, search query)';

COMMENT ON MATERIALIZED VIEW analytics_overview_daily IS 'Pre-aggregated daily platform statistics (refresh every 1 hour)';
COMMENT ON MATERIALIZED VIEW analytics_listing_stats IS 'Pre-aggregated listing performance metrics (refresh every 15 minutes)';

COMMENT ON FUNCTION log_analytics_event IS 'Helper function to log analytics events with validation';
COMMENT ON FUNCTION refresh_analytics_views IS 'Refresh all materialized views concurrently';
COMMENT ON FUNCTION archive_old_analytics_events IS 'Archive events older than 90 days to manage table size';
