-- Drop table: review_dispute_messages
DROP SEQUENCE IF EXISTS public.review_dispute_messages_id_seq;
DROP TABLE IF EXISTS public.review_dispute_messages;
DROP INDEX IF EXISTS public.idx_dispute_messages_dispute_id;