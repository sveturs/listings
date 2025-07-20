-- Drop table: attribute_option_translations
DROP SEQUENCE IF EXISTS public.attribute_option_translations_id_seq;
DROP TABLE IF EXISTS public.attribute_option_translations;
DROP INDEX IF EXISTS public.idx_attribute_option_translations;
DROP TRIGGER IF EXISTS update_listings_on_attribute_translation_change ON public.attribute_option_translations;