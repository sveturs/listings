-- Drop table: search_behavior_metrics
DROP SEQUENCE IF EXISTS public.search_behavior_metrics_id_seq;
DROP TABLE IF EXISTS public.search_behavior_metrics;
DROP INDEX IF EXISTS public.idx_search_behavior_metrics_conversions;
DROP INDEX IF EXISTS public.idx_search_behavior_metrics_ctr;
DROP INDEX IF EXISTS public.idx_search_behavior_metrics_period;
DROP INDEX IF EXISTS public.idx_search_behavior_metrics_query;
DROP INDEX IF EXISTS public.idx_search_behavior_metrics_unique;
DROP TRIGGER IF EXISTS trigger_update_search_behavior_metrics_updated_at ON public.search_behavior_metrics;