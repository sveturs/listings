-- Drop table: storefront_inventory_movements
DROP SEQUENCE IF EXISTS public.storefront_inventory_movements_id_seq;
DROP TABLE IF EXISTS public.storefront_inventory_movements;
DROP INDEX IF EXISTS public.idx_storefront_inventory_movements_created_at;
DROP INDEX IF EXISTS public.idx_storefront_inventory_movements_product_id;
DROP INDEX IF EXISTS public.idx_storefront_inventory_movements_type;
DROP INDEX IF EXISTS public.idx_storefront_inventory_movements_variant_id;