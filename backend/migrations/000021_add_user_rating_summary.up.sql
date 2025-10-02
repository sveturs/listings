-- Create user_rating_summary materialized view
-- This view aggregates reviews for users from auth service
CREATE MATERIALIZED VIEW IF NOT EXISTS user_rating_summary AS
WITH review_stats AS (
    SELECT
        r.entity_id AS user_id,
        COUNT(*) AS total_reviews,
        AVG(r.rating) AS average_rating,
        COUNT(*) FILTER (WHERE r.rating = 1) AS rating_1,
        COUNT(*) FILTER (WHERE r.rating = 2) AS rating_2,
        COUNT(*) FILTER (WHERE r.rating = 3) AS rating_3,
        COUNT(*) FILTER (WHERE r.rating = 4) AS rating_4,
        COUNT(*) FILTER (WHERE r.rating = 5) AS rating_5
    FROM reviews r
    WHERE r.entity_type = 'user'
        AND r.status = 'published'
    GROUP BY r.entity_id
)
SELECT
    CAST(user_id AS INTEGER) AS user_id,
    total_reviews,
    average_rating,
    rating_1,
    rating_2,
    rating_3,
    rating_4,
    rating_5
FROM review_stats;

-- Create unique index for CONCURRENTLY refresh support
CREATE UNIQUE INDEX IF NOT EXISTS user_rating_summary_user_id_idx
    ON user_rating_summary (user_id);

-- Fix refresh_rating_summaries function to handle missing views gracefully
CREATE OR REPLACE FUNCTION public.refresh_rating_summaries()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
    -- Refresh storefront_rating_summary if it exists
    IF EXISTS (
        SELECT 1 FROM pg_matviews
        WHERE schemaname = 'public' AND matviewname = 'storefront_rating_summary'
    ) THEN
        REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_summary;
    END IF;

    -- Refresh user_rating_summary if it exists
    IF EXISTS (
        SELECT 1 FROM pg_matviews
        WHERE schemaname = 'public' AND matviewname = 'user_rating_summary'
    ) THEN
        REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_summary;
    END IF;

    RETURN NULL;
END;
$function$;

-- Add comment
COMMENT ON MATERIALIZED VIEW user_rating_summary IS 'Aggregated user ratings from reviews';
