-- Drop table: inventory_reservations
DROP SEQUENCE IF EXISTS public.inventory_reservations_id_seq;
DROP TABLE IF EXISTS public.inventory_reservations;
DROP INDEX IF EXISTS public.idx_inventory_reservations_expires;
DROP INDEX IF EXISTS public.idx_inventory_reservations_product;