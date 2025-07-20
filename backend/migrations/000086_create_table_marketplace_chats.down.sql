-- Drop table: marketplace_chats
DROP SEQUENCE IF EXISTS public.marketplace_chats_id_seq;
DROP TABLE IF EXISTS public.marketplace_chats;
DROP INDEX IF EXISTS public.idx_marketplace_chats_active_sorted;
DROP INDEX IF EXISTS public.idx_marketplace_chats_archived;
DROP INDEX IF EXISTS public.idx_marketplace_chats_buyer;
DROP INDEX IF EXISTS public.idx_marketplace_chats_listing;
DROP INDEX IF EXISTS public.idx_marketplace_chats_listing_participants;
DROP INDEX IF EXISTS public.idx_marketplace_chats_participants;
DROP INDEX IF EXISTS public.idx_marketplace_chats_seller;
DROP INDEX IF EXISTS public.idx_marketplace_chats_updated;
DROP INDEX IF EXISTS public.idx_marketplace_chats_user_lookup;
DROP INDEX IF EXISTS public.idx_unique_direct_chat;
DROP TRIGGER IF EXISTS update_marketplace_chats_timestamp ON public.marketplace_chats;