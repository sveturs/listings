-- Drop table: storefront_products
DROP SEQUENCE IF EXISTS public.storefront_products_id_seq;
DROP TABLE IF EXISTS public.storefront_products;
DROP INDEX IF EXISTS public.idx_storefront_products_barcode;
DROP INDEX IF EXISTS public.idx_storefront_products_category_id;
DROP INDEX IF EXISTS public.idx_storefront_products_individual_location;
DROP INDEX IF EXISTS public.idx_storefront_products_is_active;
DROP INDEX IF EXISTS public.idx_storefront_products_name_gin;
DROP INDEX IF EXISTS public.idx_storefront_products_privacy;
DROP INDEX IF EXISTS public.idx_storefront_products_show_on_map;
DROP INDEX IF EXISTS public.idx_storefront_products_sku;
DROP INDEX IF EXISTS public.idx_storefront_products_stock_status;
DROP INDEX IF EXISTS public.idx_storefront_products_storefront_id;
DROP INDEX IF EXISTS public.unique_storefront_product_barcode;
DROP INDEX IF EXISTS public.unique_storefront_product_sku;
DROP TRIGGER IF EXISTS trigger_auto_geocode_storefront_product ON public.storefront_products;
DROP TRIGGER IF EXISTS trigger_cleanup_storefront_product_geo ON public.storefront_products;
DROP TRIGGER IF EXISTS trigger_storefront_products_cache_refresh ON public.storefront_products;
DROP TRIGGER IF EXISTS update_stock_status_trigger ON public.storefront_products;