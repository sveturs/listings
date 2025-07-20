-- Drop table: chat_attachments
DROP SEQUENCE IF EXISTS public.chat_attachments_id_seq;
DROP TABLE IF EXISTS public.chat_attachments;
DROP INDEX IF EXISTS public.idx_chat_attachments_created_at;
DROP INDEX IF EXISTS public.idx_chat_attachments_file_type;
DROP INDEX IF EXISTS public.idx_chat_attachments_message;
DROP INDEX IF EXISTS public.idx_chat_attachments_message_id;