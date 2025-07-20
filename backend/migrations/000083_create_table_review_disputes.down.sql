-- Drop table: review_disputes
DROP SEQUENCE IF EXISTS public.review_disputes_id_seq;
DROP TABLE IF EXISTS public.review_disputes;
DROP INDEX IF EXISTS public.idx_disputes_review_id;
DROP INDEX IF EXISTS public.idx_disputes_status;