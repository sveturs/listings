-- Drop table: order_messages
DROP SEQUENCE IF EXISTS public.order_messages_id_seq;
DROP TABLE IF EXISTS public.order_messages;
DROP INDEX IF EXISTS public.idx_order_messages_created_at;
DROP INDEX IF EXISTS public.idx_order_messages_order_id;
DROP INDEX IF EXISTS public.idx_order_messages_sender_id;