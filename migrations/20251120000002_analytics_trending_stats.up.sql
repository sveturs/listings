-- Migration: Create materialized view for trending analytics
-- This provides cached trending data that refreshes periodically for optimal performance

-- Create materialized view for trending stats
CREATE MATERIALIZED VIEW IF NOT EXISTS analytics_trending_cache AS
WITH category_trends AS (
    -- Trending Categories: rate of change in orders
    WITH category_orders_30d AS (
        SELECT
            l.category_id,
            c.name as category_name,
            COUNT(DISTINCT o.id) as order_count
        FROM orders o
        JOIN order_items oi ON oi.order_id = o.id
        JOIN listings l ON l.id = oi.listing_id
        JOIN categories c ON c.id = l.category_id
        WHERE o.created_at >= NOW() - INTERVAL '30 days'
            AND o.status = 'completed'
        GROUP BY l.category_id, c.name
    ),
    category_orders_7d AS (
        SELECT
            l.category_id,
            COUNT(DISTINCT o.id) as order_count
        FROM orders o
        JOIN order_items oi ON oi.order_id = o.id
        JOIN listings l ON l.id = oi.listing_id
        WHERE o.created_at >= NOW() - INTERVAL '7 days'
            AND o.status = 'completed'
        GROUP BY l.category_id
    )
    SELECT
        c30.category_id,
        c30.category_name,
        c30.order_count as order_count_30d,
        COALESCE(c7.order_count, 0) as order_count_7d,
        -- Growth rate: (7d_rate - 30d_rate) / 30d_rate * 100
        CASE
            WHEN c30.order_count > 0 THEN
                ((COALESCE(c7.order_count, 0)::float / 7) - (c30.order_count::float / 30))
                / (c30.order_count::float / 30) * 100
            ELSE 0
        END as growth_rate,
        -- Trend score: combines volume and growth (weights recent activity higher)
        (COALESCE(c7.order_count, 0) * 10 + c30.order_count * 1) as trend_score
    FROM category_orders_30d c30
    LEFT JOIN category_orders_7d c7 ON c7.category_id = c30.category_id
    ORDER BY trend_score DESC
    LIMIT 10
),
hot_listings AS (
    -- Hot Listings: based on recent order activity (views data not available in listings microservice)
    -- Using orders and order_items as proxy for popularity
    WITH listing_orders_24h AS (
        SELECT
            oi.listing_id,
            COUNT(DISTINCT o.id) as order_count,
            SUM(oi.quantity) as total_quantity
        FROM orders o
        JOIN order_items oi ON oi.order_id = o.id
        WHERE o.created_at >= NOW() - INTERVAL '24 hours'
            AND o.status IN ('confirmed', 'processing', 'shipped', 'delivered', 'completed')
        GROUP BY oi.listing_id
    ),
    listing_orders_7d AS (
        SELECT
            oi.listing_id,
            COUNT(DISTINCT o.id) as order_count,
            SUM(oi.quantity) as total_quantity
        FROM orders o
        JOIN order_items oi ON oi.order_id = o.id
        WHERE o.created_at >= NOW() - INTERVAL '7 days'
            AND o.status IN ('confirmed', 'processing', 'shipped', 'delivered', 'completed')
        GROUP BY oi.listing_id
    )
    SELECT
        l.id as listing_id,
        l.title,
        COALESCE(lo24.order_count, 0) as orders_24h,
        COALESCE(lo7.order_count, 0) as orders_7d,
        -- Growth: ratio of 24h orders to 7d average
        CASE
            WHEN COALESCE(lo7.order_count, 0) > 0 THEN
                (COALESCE(lo24.order_count, 0)::float / (COALESCE(lo7.order_count, 0)::float / 7))
            ELSE 0
        END as orders_growth,
        COALESCE(lo24.total_quantity, 0) as quantity_sold_24h,
        l.price
    FROM listings l
    LEFT JOIN listing_orders_24h lo24 ON lo24.listing_id = l.id
    LEFT JOIN listing_orders_7d lo7 ON lo7.listing_id = l.id
    WHERE l.status = 'active'
        AND COALESCE(lo24.order_count, 0) > 0  -- Must have at least 1 order in 24h
    ORDER BY orders_growth DESC, orders_24h DESC
    LIMIT 20
),
popular_searches AS (
    -- Popular Searches: reuse Phase 28 implementation
    SELECT
        query_text as query,
        COUNT(*) as search_count
    FROM search_queries
    WHERE created_at >= NOW() - INTERVAL '7 days'
        AND results_count > 0  -- Only queries that returned results
    GROUP BY query_text
    ORDER BY search_count DESC
    LIMIT 10
)
SELECT
    json_build_object(
        'trending_categories', (SELECT json_agg(row_to_json(category_trends.*)) FROM category_trends),
        'hot_listings', (SELECT json_agg(row_to_json(hot_listings.*)) FROM hot_listings),
        'popular_searches', (SELECT json_agg(row_to_json(popular_searches.*)) FROM popular_searches),
        'generated_at', NOW()
    ) as trending_data;

-- Create index on materialized view for fast lookup
CREATE UNIQUE INDEX IF NOT EXISTS idx_analytics_trending_cache_singleton
ON analytics_trending_cache ((trending_data IS NOT NULL));

-- Create function to refresh the materialized view
CREATE OR REPLACE FUNCTION refresh_analytics_trending_cache()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_trending_cache;
END;
$$ LANGUAGE plpgsql;

-- Add comment for documentation
COMMENT ON MATERIALIZED VIEW analytics_trending_cache IS
'Cached trending analytics data: trending categories, hot listings, popular searches. Refresh hourly via cron or trigger.';

COMMENT ON FUNCTION refresh_analytics_trending_cache() IS
'Refresh the trending analytics cache. Run this hourly via cron job or trigger.';
