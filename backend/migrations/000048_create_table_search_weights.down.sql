-- Drop table: search_weights
DROP SEQUENCE IF EXISTS public.search_weights_id_seq;
DROP TABLE IF EXISTS public.search_weights;
DROP INDEX IF EXISTS public.idx_search_weights_category_id;
DROP INDEX IF EXISTS public.idx_search_weights_field_name;
DROP INDEX IF EXISTS public.idx_search_weights_is_active;
DROP INDEX IF EXISTS public.idx_search_weights_item_type;
DROP INDEX IF EXISTS public.idx_search_weights_unique_category;
DROP INDEX IF EXISTS public.idx_search_weights_unique_global;
DROP INDEX IF EXISTS public.idx_search_weights_version;
DROP TRIGGER IF EXISTS trigger_log_search_weight_changes ON public.search_weights;
DROP TRIGGER IF EXISTS trigger_update_search_weights_updated_at ON public.search_weights;