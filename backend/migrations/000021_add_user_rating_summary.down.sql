-- Restore original refresh_rating_summaries function
CREATE OR REPLACE FUNCTION public.refresh_rating_summaries()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_summary;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_summary;
    RETURN NULL;
END;
$function$;

-- Drop user_rating_summary view
DROP MATERIALIZED VIEW IF EXISTS user_rating_summary CASCADE;
