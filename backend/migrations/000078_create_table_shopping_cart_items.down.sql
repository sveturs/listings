-- Drop table: shopping_cart_items
DROP SEQUENCE IF EXISTS public.shopping_cart_items_id_seq;
DROP TABLE IF EXISTS public.shopping_cart_items;
DROP INDEX IF EXISTS public.idx_shopping_cart_items_cart_id;
DROP INDEX IF EXISTS public.idx_shopping_cart_items_product_id;
DROP TRIGGER IF EXISTS trigger_shopping_cart_items_updated_at ON public.shopping_cart_items;