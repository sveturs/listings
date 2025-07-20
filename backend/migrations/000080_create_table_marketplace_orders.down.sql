-- Drop table: marketplace_orders
DROP SEQUENCE IF EXISTS public.marketplace_orders_id_seq;
DROP TABLE IF EXISTS public.marketplace_orders;
DROP INDEX IF EXISTS public.idx_marketplace_orders_buyer;
DROP INDEX IF EXISTS public.idx_marketplace_orders_buyer_id;
DROP INDEX IF EXISTS public.idx_marketplace_orders_created_at;
DROP INDEX IF EXISTS public.idx_marketplace_orders_listing;
DROP INDEX IF EXISTS public.idx_marketplace_orders_listing_id;
DROP INDEX IF EXISTS public.idx_marketplace_orders_payment_transaction_id;
DROP INDEX IF EXISTS public.idx_marketplace_orders_protection;
DROP INDEX IF EXISTS public.idx_marketplace_orders_protection_expires_at;
DROP INDEX IF EXISTS public.idx_marketplace_orders_seller;
DROP INDEX IF EXISTS public.idx_marketplace_orders_seller_id;
DROP INDEX IF EXISTS public.idx_marketplace_orders_status;
DROP TRIGGER IF EXISTS marketplace_orders_updated_at_trigger ON public.marketplace_orders;