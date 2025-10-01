-- Create user_ratings materialized view (comprehensive version)
-- This view includes direct reviews and reviews through entity_origin
CREATE MATERIALIZED VIEW IF NOT EXISTS user_ratings AS
SELECT
    CAST(users.user_id AS INTEGER) as user_id,
    COUNT(DISTINCT r.id) as total_reviews,
    COALESCE(AVG(r.rating), 0) as average_rating,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'user' THEN r.id END) as direct_reviews,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'marketplace_listing' THEN r.id END) as listing_reviews,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'storefront' THEN r.id END) as storefront_reviews,
    COUNT(DISTINCT CASE WHEN r.is_verified_purchase THEN r.id END) as verified_reviews,
    -- Rating distribution
    COUNT(CASE WHEN r.rating = 1 THEN 1 END) as rating_1,
    COUNT(CASE WHEN r.rating = 2 THEN 1 END) as rating_2,
    COUNT(CASE WHEN r.rating = 3 THEN 1 END) as rating_3,
    COUNT(CASE WHEN r.rating = 4 THEN 1 END) as rating_4,
    COUNT(CASE WHEN r.rating = 5 THEN 1 END) as rating_5,
    -- Recent trend (last 30 days)
    AVG(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN r.rating END) as recent_rating,
    COUNT(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN 1 END) as recent_reviews,
    MAX(r.created_at) as last_review_at
FROM (
    -- Get all unique user_ids from reviews (both as entity and origin)
    SELECT DISTINCT COALESCE(entity_origin_id, entity_id) as user_id
    FROM reviews
    WHERE (entity_type = 'user' OR entity_origin_type = 'user')
      AND status = 'published'
) users
LEFT JOIN reviews r ON (
    -- Direct reviews on user
    (r.entity_type = 'user' AND r.entity_id = users.user_id) OR
    -- Reviews through origin after deletion
    (r.entity_origin_type = 'user' AND r.entity_origin_id = users.user_id)
) AND r.status = 'published'
GROUP BY users.user_id;

-- Create unique index for CONCURRENTLY refresh support
CREATE UNIQUE INDEX IF NOT EXISTS user_ratings_user_id_idx
    ON user_ratings (user_id);

-- Create additional indexes
CREATE INDEX IF NOT EXISTS idx_user_ratings_average ON user_ratings(average_rating DESC);
CREATE INDEX IF NOT EXISTS idx_user_ratings_total ON user_ratings(total_reviews DESC);

-- Create storefront_ratings materialized view
CREATE MATERIALIZED VIEW IF NOT EXISTS storefront_ratings AS
SELECT
    s.id as storefront_id,
    COUNT(DISTINCT r.id) as total_reviews,
    COALESCE(AVG(r.rating), 0) as average_rating,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'storefront' THEN r.id END) as direct_reviews,
    COUNT(DISTINCT CASE WHEN r.entity_type = 'marketplace_listing' THEN r.id END) as listing_reviews,
    COUNT(DISTINCT CASE WHEN r.is_verified_purchase THEN r.id END) as verified_reviews,
    -- Rating distribution
    COUNT(CASE WHEN r.rating = 1 THEN 1 END) as rating_1,
    COUNT(CASE WHEN r.rating = 2 THEN 1 END) as rating_2,
    COUNT(CASE WHEN r.rating = 3 THEN 1 END) as rating_3,
    COUNT(CASE WHEN r.rating = 4 THEN 1 END) as rating_4,
    COUNT(CASE WHEN r.rating = 5 THEN 1 END) as rating_5,
    -- Recent trend (last 30 days)
    AVG(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN r.rating END) as recent_rating,
    COUNT(CASE WHEN r.created_at > NOW() - INTERVAL '30 days' THEN 1 END) as recent_reviews,
    MAX(r.created_at) as last_review_at,
    -- Owner info
    s.user_id as owner_id
FROM user_storefronts s
LEFT JOIN reviews r ON (
    -- Direct reviews on storefront
    (r.entity_type = 'storefront' AND r.entity_id = s.id) OR
    -- Reviews through origin after deletion
    (r.entity_origin_type = 'storefront' AND r.entity_origin_id = s.id)
) AND r.status = 'published'
GROUP BY s.id;

-- Create unique index for CONCURRENTLY refresh support
CREATE UNIQUE INDEX IF NOT EXISTS storefront_ratings_storefront_id_idx
    ON storefront_ratings (storefront_id);

-- Create additional indexes
CREATE INDEX IF NOT EXISTS idx_storefront_ratings_average ON storefront_ratings(average_rating DESC);
CREATE INDEX IF NOT EXISTS idx_storefront_ratings_owner ON storefront_ratings(owner_id);

-- Fix refresh_rating_views function to handle missing views gracefully
CREATE OR REPLACE FUNCTION public.refresh_rating_views()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
    -- Update only affected rows, not entire view
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- For users
        IF NEW.entity_origin_type = 'user' AND EXISTS (
            SELECT 1 FROM pg_matviews
            WHERE schemaname = 'public' AND matviewname = 'user_ratings'
        ) THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;

        -- For storefronts
        IF NEW.entity_origin_type = 'storefront' AND EXISTS (
            SELECT 1 FROM pg_matviews
            WHERE schemaname = 'public' AND matviewname = 'storefront_ratings'
        ) THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- Also update on deletion
        IF OLD.entity_origin_type = 'user' AND EXISTS (
            SELECT 1 FROM pg_matviews
            WHERE schemaname = 'public' AND matviewname = 'user_ratings'
        ) THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;

        IF OLD.entity_origin_type = 'storefront' AND EXISTS (
            SELECT 1 FROM pg_matviews
            WHERE schemaname = 'public' AND matviewname = 'storefront_ratings'
        ) THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    END IF;

    RETURN NULL;
END;
$function$;

-- Add comments
COMMENT ON MATERIALIZED VIEW user_ratings IS 'Comprehensive user ratings aggregation with all review types';
COMMENT ON MATERIALIZED VIEW storefront_ratings IS 'Comprehensive storefront ratings aggregation';
COMMENT ON FUNCTION refresh_rating_views IS 'Updates rating materialized views with graceful error handling';
