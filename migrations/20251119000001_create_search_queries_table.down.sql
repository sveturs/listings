-- Migration Rollback: Drop Search Queries Analytics Table
-- Phase: Phase 28 - Search Analytics Infrastructure
-- Date: 2025-11-19

-- Drop all indexes first (for clean rollback)
DROP INDEX IF EXISTS idx_search_queries_ctr_analysis;
DROP INDEX IF EXISTS idx_search_queries_query_text_fts;
DROP INDEX IF EXISTS idx_search_queries_category_agg;
DROP INDEX IF EXISTS idx_search_queries_created_at;
DROP INDEX IF EXISTS idx_search_queries_session_history;
DROP INDEX IF EXISTS idx_search_queries_user_history;
DROP INDEX IF EXISTS idx_search_queries_trending;

-- Drop the table
DROP TABLE IF EXISTS search_queries;

-- Note: This migration is safe to rollback at any time
-- No dependent tables or foreign keys reference search_queries
-- Data loss is acceptable (analytics/logging table, can be rebuilt)
