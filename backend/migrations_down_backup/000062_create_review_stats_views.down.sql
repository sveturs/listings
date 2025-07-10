-- Rollback migration: Drop review statistics views
-- This migration removes the materialized views created for review statistics

-- Drop function
DROP FUNCTION IF EXISTS refresh_rating_views();

-- Drop indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_listing_ratings_owner_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_listing_ratings_listing_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_storefront_ratings_owner_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_storefront_ratings_storefront_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_ratings_user_id;

-- Drop materialized views
DROP MATERIALIZED VIEW IF EXISTS listing_ratings;
DROP MATERIALIZED VIEW IF EXISTS storefront_ratings;
DROP MATERIALIZED VIEW IF EXISTS user_ratings;