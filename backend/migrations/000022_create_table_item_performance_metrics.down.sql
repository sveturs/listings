-- Drop table: item_performance_metrics
DROP SEQUENCE IF EXISTS public.item_performance_metrics_id_seq;
DROP TABLE IF EXISTS public.item_performance_metrics;
DROP INDEX IF EXISTS public.idx_item_performance_metrics_conversions;
DROP INDEX IF EXISTS public.idx_item_performance_metrics_ctr;
DROP INDEX IF EXISTS public.idx_item_performance_metrics_impressions;
DROP INDEX IF EXISTS public.idx_item_performance_metrics_item;
DROP INDEX IF EXISTS public.idx_item_performance_metrics_period;
DROP INDEX IF EXISTS public.idx_item_performance_metrics_unique;
DROP TRIGGER IF EXISTS trigger_update_item_performance_metrics_updated_at ON public.item_performance_metrics;