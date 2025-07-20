-- Drop table: search_weights_history
DROP SEQUENCE IF EXISTS public.search_weights_history_id_seq;
DROP TABLE IF EXISTS public.search_weights_history;
DROP INDEX IF EXISTS public.idx_search_weights_history_changed_at;
DROP INDEX IF EXISTS public.idx_search_weights_history_changed_by;
DROP INDEX IF EXISTS public.idx_search_weights_history_reason;
DROP INDEX IF EXISTS public.idx_search_weights_history_weight_id;