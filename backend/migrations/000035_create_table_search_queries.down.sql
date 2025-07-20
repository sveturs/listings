-- Drop table: search_queries
DROP SEQUENCE IF EXISTS public.search_queries_id_seq;
DROP TABLE IF EXISTS public.search_queries;
DROP INDEX IF EXISTS public.idx_search_queries_language;
DROP INDEX IF EXISTS public.idx_search_queries_normalized_query;
DROP INDEX IF EXISTS public.idx_search_queries_search_count;
DROP TRIGGER IF EXISTS update_search_queries_updated_at_trigger ON public.search_queries;