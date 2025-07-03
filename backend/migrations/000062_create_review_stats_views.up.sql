-- Migration: Create materialized views for review statistics
-- This migration creates the user_ratings and storefront_ratings materialized views
-- that were needed to fix the 500 error in /api/v1/entity/user/{id}/stats endpoint

-- Create user_ratings materialized view
CREATE MATERIALIZED VIEW IF NOT EXISTS user_ratings AS
SELECT 
    u.id AS user_id,
    COALESCE(AVG(r.rating), 0) AS average_rating,
    COUNT(r.id) AS total_reviews,
    COUNT(CASE WHEN r.verified THEN 1 END) AS verified_reviews,
    COUNT(CASE WHEN array_length(r.photo_urls, 1) > 0 THEN 1 END) AS photo_reviews,
    -- Rating distribution as JSON
    jsonb_object_agg(
        r.rating::text, 
        COALESCE(rating_counts.count, 0)
    ) FILTER (WHERE r.rating IS NOT NULL) AS rating_distribution
FROM users u
LEFT JOIN reviews r ON r.entity_type = 'user' AND r.entity_id = u.id
LEFT JOIN (
    -- Subquery to count ratings distribution
    SELECT 
        entity_id,
        rating,
        COUNT(*) as count
    FROM reviews 
    WHERE entity_type = 'user'
    GROUP BY entity_id, rating
) rating_counts ON rating_counts.entity_id = u.id AND rating_counts.rating = r.rating
GROUP BY u.id;

-- Create storefront_ratings materialized view with owner_id
CREATE MATERIALIZED VIEW IF NOT EXISTS storefront_ratings AS
SELECT 
    s.id AS storefront_id,
    s.user_id AS owner_id,  -- Added owner_id field that was missing
    COALESCE(AVG(r.rating), 0) AS average_rating,
    COUNT(r.id) AS total_reviews,
    COUNT(CASE WHEN r.verified THEN 1 END) AS verified_reviews,
    COUNT(CASE WHEN array_length(r.photo_urls, 1) > 0 THEN 1 END) AS photo_reviews,
    -- Rating distribution as JSON
    jsonb_object_agg(
        r.rating::text, 
        COALESCE(rating_counts.count, 0)
    ) FILTER (WHERE r.rating IS NOT NULL) AS rating_distribution
FROM storefronts s
LEFT JOIN reviews r ON r.entity_type = 'storefront' AND r.entity_id = s.id
LEFT JOIN (
    -- Subquery to count ratings distribution
    SELECT 
        entity_id,
        rating,
        COUNT(*) as count
    FROM reviews 
    WHERE entity_type = 'storefront'
    GROUP BY entity_id, rating
) rating_counts ON rating_counts.entity_id = s.id AND rating_counts.rating = r.rating
GROUP BY s.id, s.user_id;

-- Create listing_ratings materialized view for completeness
CREATE MATERIALIZED VIEW IF NOT EXISTS listing_ratings AS
SELECT 
    ml.id AS listing_id,
    ml.user_id AS owner_id,
    COALESCE(AVG(r.rating), 0) AS average_rating,
    COUNT(r.id) AS total_reviews,
    COUNT(CASE WHEN r.verified THEN 1 END) AS verified_reviews,
    COUNT(CASE WHEN array_length(r.photo_urls, 1) > 0 THEN 1 END) AS photo_reviews,
    -- Rating distribution as JSON
    jsonb_object_agg(
        r.rating::text, 
        COALESCE(rating_counts.count, 0)
    ) FILTER (WHERE r.rating IS NOT NULL) AS rating_distribution
FROM marketplace_listings ml
LEFT JOIN reviews r ON r.entity_type = 'listing' AND r.entity_id = ml.id
LEFT JOIN (
    -- Subquery to count ratings distribution
    SELECT 
        entity_id,
        rating,
        COUNT(*) as count
    FROM reviews 
    WHERE entity_type = 'listing'
    GROUP BY entity_id, rating
) rating_counts ON rating_counts.entity_id = ml.id AND rating_counts.rating = r.rating
GROUP BY ml.id, ml.user_id;

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_ratings_user_id ON user_ratings(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_storefront_ratings_storefront_id ON storefront_ratings(storefront_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_storefront_ratings_owner_id ON storefront_ratings(owner_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_listing_ratings_listing_id ON listing_ratings(listing_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_listing_ratings_owner_id ON listing_ratings(owner_id);

-- Function to refresh all rating views
CREATE OR REPLACE FUNCTION refresh_rating_views() RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW user_ratings;
    REFRESH MATERIALIZED VIEW storefront_ratings;
    REFRESH MATERIALIZED VIEW listing_ratings;
END;
$$ LANGUAGE plpgsql;

-- Initial refresh
SELECT refresh_rating_views();