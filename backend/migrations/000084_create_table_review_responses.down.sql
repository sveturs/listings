-- Drop table: review_responses
DROP SEQUENCE IF EXISTS public.review_responses_id_seq;
DROP TABLE IF EXISTS public.review_responses;
DROP INDEX IF EXISTS public.idx_review_responses_review;
DROP TRIGGER IF EXISTS update_review_responses_updated_at ON public.review_responses;