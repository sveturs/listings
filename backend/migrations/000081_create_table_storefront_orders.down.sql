-- Drop table: storefront_orders
DROP SEQUENCE IF EXISTS public.storefront_orders_id_seq;
DROP TABLE IF EXISTS public.storefront_orders;
DROP INDEX IF EXISTS public.idx_storefront_orders_customer;
DROP INDEX IF EXISTS public.idx_storefront_orders_escrow_date;
DROP INDEX IF EXISTS public.idx_storefront_orders_status;
DROP INDEX IF EXISTS public.idx_storefront_orders_storefront;
DROP TRIGGER IF EXISTS calculate_escrow_release_date_trigger ON public.storefront_orders;
DROP TRIGGER IF EXISTS set_order_number_trigger ON public.storefront_orders;