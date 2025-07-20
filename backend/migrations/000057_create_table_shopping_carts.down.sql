-- Drop table: shopping_carts
DROP SEQUENCE IF EXISTS public.shopping_carts_id_seq;
DROP TABLE IF EXISTS public.shopping_carts;
DROP INDEX IF EXISTS public.idx_shopping_carts_session_id;
DROP INDEX IF EXISTS public.idx_shopping_carts_storefront_id;
DROP INDEX IF EXISTS public.idx_shopping_carts_user_id;
DROP TRIGGER IF EXISTS trigger_shopping_carts_updated_at ON public.shopping_carts;