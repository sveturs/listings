-- Drop table: price_history
DROP SEQUENCE IF EXISTS public.price_history_id_seq;
DROP TABLE IF EXISTS public.price_history;
DROP INDEX IF EXISTS public.idx_price_history_current;
DROP INDEX IF EXISTS public.idx_price_history_effective;
DROP INDEX IF EXISTS public.idx_price_history_listing_id;
DROP TRIGGER IF EXISTS trig_update_metadata_after_price_change ON public.price_history;