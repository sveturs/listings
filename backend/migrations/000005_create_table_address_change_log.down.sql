-- Drop table: address_change_log
DROP SEQUENCE IF EXISTS public.address_change_log_id_seq;
DROP TABLE IF EXISTS public.address_change_log;
DROP INDEX IF EXISTS public.idx_address_log_change_reason;
DROP INDEX IF EXISTS public.idx_address_log_confidence_after;
DROP INDEX IF EXISTS public.idx_address_log_created_at;
DROP INDEX IF EXISTS public.idx_address_log_listing_id;
DROP INDEX IF EXISTS public.idx_address_log_new_location;
DROP INDEX IF EXISTS public.idx_address_log_old_location;
DROP INDEX IF EXISTS public.idx_address_log_user_id;