-- Drop table: translations
DROP SEQUENCE IF EXISTS public.translations_id_seq;
DROP TABLE IF EXISTS public.translations;
DROP INDEX IF EXISTS public.idx_translations_lookup;
DROP INDEX IF EXISTS public.idx_translations_metadata;
DROP INDEX IF EXISTS public.idx_translations_type_lang;
DROP TRIGGER IF EXISTS update_translations_timestamp ON public.translations;