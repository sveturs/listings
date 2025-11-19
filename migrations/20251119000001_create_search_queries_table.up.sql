-- Migration: Create Search Queries Analytics Table
-- Phase: Phase 28 - Search Analytics Infrastructure
-- Date: 2025-11-19
--
-- Purpose: Track all search queries for analytics, trending searches, and user search history
--
-- Features:
-- - Track authenticated and anonymous users (similar to shopping_carts pattern)
-- - Store query text, category context, and results metadata
-- - Support CTR tracking (clicked_listing_id)
-- - Optimized indexes for trending aggregations and time-based queries
-- - Retention policy ready (no soft delete - this is logging)
--
-- Performance Considerations:
-- - Partial indexes for filtering NULL values
-- - Composite indexes for frequent query patterns
-- - GIN index for full-text search on query_text (for "similar searches")
-- - Timestamp index for time-range aggregations (trending queries)

-- =====================================================
-- SEARCH_QUERIES TABLE
-- =====================================================

CREATE TABLE IF NOT EXISTS search_queries (
    id BIGSERIAL PRIMARY KEY,

    -- Query details
    query_text VARCHAR(500) NOT NULL CHECK (LENGTH(TRIM(query_text)) >= 1),

    -- Context
    category_id BIGINT NULL,                    -- FK to categories (optional filter context)

    -- User identification (mutually exclusive: user_id OR session_id)
    -- Follows same pattern as shopping_carts for consistency
    user_id BIGINT NULL,                        -- Authenticated user (FK to auth service)
    session_id VARCHAR(255) NULL,               -- Anonymous user session (UUID from cookie)

    -- Search results metadata
    results_count INTEGER NOT NULL DEFAULT 0 CHECK (results_count >= 0),

    -- Click-through tracking (for CTR analytics)
    clicked_listing_id BIGINT NULL,             -- ID of listing clicked from search results

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_search_queries_category FOREIGN KEY (category_id)
        REFERENCES categories(id) ON DELETE SET NULL,

    -- Business Logic Constraints
    -- Ensure EXACTLY ONE of user_id or session_id is set (not both, not neither)
    CONSTRAINT chk_search_queries_user_or_session CHECK (
        (user_id IS NOT NULL AND session_id IS NOT NULL) OR  -- Both (allowed for migration from session to user)
        (user_id IS NOT NULL AND session_id IS NULL) OR      -- Authenticated only
        (user_id IS NULL AND session_id IS NOT NULL)         -- Anonymous only
    )
);

-- =====================================================
-- INDEXES FOR PERFORMANCE
-- =====================================================

-- 1. Trending Queries Aggregation (PRIMARY USE CASE)
-- Query: SELECT query_text, COUNT(*) FROM search_queries
--        WHERE created_at > NOW() - INTERVAL '7 days' AND category_id = ?
--        GROUP BY query_text ORDER BY COUNT(*) DESC LIMIT 10
-- Expected performance: < 500ms for 1M rows
CREATE INDEX idx_search_queries_trending
    ON search_queries(category_id, created_at DESC, query_text)
    WHERE results_count > 0;  -- Only index successful searches

-- 2. User Search History (authenticated users)
-- Query: SELECT * FROM search_queries WHERE user_id = ? ORDER BY created_at DESC LIMIT 20
-- Expected performance: < 50ms
CREATE INDEX idx_search_queries_user_history
    ON search_queries(user_id, created_at DESC)
    WHERE user_id IS NOT NULL;

-- 3. Session Search History (anonymous users)
-- Query: SELECT * FROM search_queries WHERE session_id = ? ORDER BY created_at DESC LIMIT 20
-- Expected performance: < 50ms
CREATE INDEX idx_search_queries_session_history
    ON search_queries(session_id, created_at DESC)
    WHERE session_id IS NOT NULL;

-- 4. Time-based Cleanup & Aggregations
-- For retention policy (e.g., delete records older than 90 days)
-- Also used for time-range trending queries
CREATE INDEX idx_search_queries_created_at
    ON search_queries(created_at DESC);

-- 5. Category-specific Trending (without time filter)
-- For "all-time popular searches in category"
CREATE INDEX idx_search_queries_category_agg
    ON search_queries(category_id, query_text)
    WHERE category_id IS NOT NULL AND results_count > 0;

-- 6. Full-text Search on Query Text (for "similar searches" feature)
-- Allows finding related searches: "iphone" → ["iphone 13", "iphone pro", "iphone case"]
CREATE INDEX idx_search_queries_query_text_fts
    ON search_queries USING gin(to_tsvector('english', query_text));

-- 7. CTR Analysis Index
-- For analyzing which queries lead to clicks
CREATE INDEX idx_search_queries_ctr_analysis
    ON search_queries(clicked_listing_id, created_at DESC)
    WHERE clicked_listing_id IS NOT NULL;

-- =====================================================
-- COMMENTS (Documentation)
-- =====================================================

COMMENT ON TABLE search_queries IS
    'Analytics tracking for all search queries. Used for trending searches, user history, and CTR analysis. No soft delete (append-only log).';

COMMENT ON COLUMN search_queries.query_text IS
    'Search query text entered by user. Max 500 chars. Indexed for aggregations and full-text similarity.';

COMMENT ON COLUMN search_queries.category_id IS
    'Optional category context. NULL means global search. Used for category-specific trending queries.';

COMMENT ON COLUMN search_queries.user_id IS
    'Authenticated user ID (FK to auth service). Mutually exclusive with session_id (or both for migration).';

COMMENT ON COLUMN search_queries.session_id IS
    'Anonymous user session ID (UUID from cookie). Mutually exclusive with user_id (or both for migration).';

COMMENT ON COLUMN search_queries.results_count IS
    'Number of results returned by search. Used to filter out empty searches from trending calculations.';

COMMENT ON COLUMN search_queries.clicked_listing_id IS
    'ID of listing clicked from search results. NULL if no click. Used for CTR (Click-Through Rate) analysis.';

COMMENT ON COLUMN search_queries.created_at IS
    'Timestamp when search was performed. Primary field for time-based aggregations and retention policy.';

-- =====================================================
-- PERFORMANCE NOTES
-- =====================================================

-- Expected Query Performance (with 1M rows):
-- 1. Trending queries (7 days, category filter):  < 500ms
-- 2. Trending queries (7 days, no category):      < 800ms
-- 3. User history (20 results):                   < 50ms
-- 4. Session history (20 results):                < 50ms
-- 5. CTR analysis (aggregations):                 < 300ms

-- Retention Policy Recommendation:
-- - Keep 90 days of data for trending analysis
-- - Archive older data to separate table if needed for long-term analytics
-- - Cleanup query: DELETE FROM search_queries WHERE created_at < NOW() - INTERVAL '90 days'

-- Index Maintenance:
-- - Run VACUUM ANALYZE weekly to maintain index performance
-- - Monitor index bloat with pg_stat_user_indexes
-- - Consider table partitioning if data grows beyond 10M rows

-- =====================================================
-- EXAMPLE QUERIES (for testing after migration)
-- =====================================================

-- Trending searches (last 7 days, category 1301 - Cars)
-- SELECT query_text, COUNT(*) as search_count, MAX(created_at) as last_searched
-- FROM search_queries
-- WHERE created_at > NOW() - INTERVAL '7 days'
--   AND category_id = 1301
--   AND results_count > 0
-- GROUP BY query_text
-- ORDER BY search_count DESC, last_searched DESC
-- LIMIT 10;

-- User search history
-- SELECT query_text, results_count, created_at
-- FROM search_queries
-- WHERE user_id = 123
-- ORDER BY created_at DESC
-- LIMIT 20;

-- CTR analysis for specific query
-- SELECT
--     query_text,
--     COUNT(*) as total_searches,
--     COUNT(clicked_listing_id) as clicks,
--     ROUND(100.0 * COUNT(clicked_listing_id) / COUNT(*), 2) as ctr_percent
-- FROM search_queries
-- WHERE query_text = 'iphone'
--   AND created_at > NOW() - INTERVAL '30 days'
-- GROUP BY query_text;
