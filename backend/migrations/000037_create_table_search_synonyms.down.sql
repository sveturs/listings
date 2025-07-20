-- Drop table: search_synonyms
DROP SEQUENCE IF EXISTS public.search_synonyms_id_seq;
DROP TABLE IF EXISTS public.search_synonyms;
DROP INDEX IF EXISTS public.idx_search_synonyms_active;
DROP INDEX IF EXISTS public.idx_search_synonyms_language;
DROP INDEX IF EXISTS public.idx_search_synonyms_synonym;
DROP INDEX IF EXISTS public.idx_search_synonyms_term;
DROP INDEX IF EXISTS public.idx_search_synonyms_unique;
DROP TRIGGER IF EXISTS trigger_update_search_synonyms_updated_at ON public.search_synonyms;