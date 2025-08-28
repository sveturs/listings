-- Drop table: marketplace_messages
DROP SEQUENCE IF EXISTS public.marketplace_messages_id_seq;
DROP TABLE IF EXISTS public.marketplace_messages;
DROP INDEX IF EXISTS public.idx_marketplace_messages_chat;
DROP INDEX IF EXISTS public.idx_marketplace_messages_chat_last;
DROP INDEX IF EXISTS public.idx_marketplace_messages_chat_ordered;
DROP INDEX IF EXISTS public.idx_marketplace_messages_chat_unread;
DROP INDEX IF EXISTS public.idx_marketplace_messages_created;
DROP INDEX IF EXISTS public.idx_marketplace_messages_listing;
DROP INDEX IF EXISTS public.idx_marketplace_messages_receiver;
DROP INDEX IF EXISTS public.idx_marketplace_messages_receiver_unread_count;
DROP INDEX IF EXISTS public.idx_marketplace_messages_sender;
DROP INDEX IF EXISTS public.idx_marketplace_messages_unread;
DROP TRIGGER IF EXISTS update_marketplace_messages_timestamp ON public.marketplace_messages;