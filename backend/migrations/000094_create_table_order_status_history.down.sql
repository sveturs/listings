-- Drop table: order_status_history
DROP SEQUENCE IF EXISTS public.order_status_history_id_seq;
DROP TABLE IF EXISTS public.order_status_history;
DROP INDEX IF EXISTS public.idx_order_status_history_created_at;
DROP INDEX IF EXISTS public.idx_order_status_history_order_id;