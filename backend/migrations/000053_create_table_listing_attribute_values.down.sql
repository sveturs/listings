-- Drop table: listing_attribute_values
DROP SEQUENCE IF EXISTS public.listing_attribute_values_id_seq;
DROP TABLE IF EXISTS public.listing_attribute_values;
DROP INDEX IF EXISTS public.idx_attr_name_num_val;
DROP INDEX IF EXISTS public.idx_attr_name_text_val;
DROP INDEX IF EXISTS public.idx_attr_unit;
DROP INDEX IF EXISTS public.idx_listing_attr_boolean;
DROP INDEX IF EXISTS public.idx_listing_attr_listing_id;
DROP INDEX IF EXISTS public.idx_listing_attr_numeric;
DROP INDEX IF EXISTS public.idx_listing_attr_text;
DROP INDEX IF EXISTS public.idx_listing_attr_unique;