-- Drop table: transliteration_rules
DROP SEQUENCE IF EXISTS public.transliteration_rules_id_seq;
DROP TABLE IF EXISTS public.transliteration_rules;
DROP INDEX IF EXISTS public.idx_transliteration_rules_active;
DROP INDEX IF EXISTS public.idx_transliteration_rules_enabled;
DROP INDEX IF EXISTS public.idx_transliteration_rules_language;
DROP INDEX IF EXISTS public.idx_transliteration_rules_priority;
DROP INDEX IF EXISTS public.idx_transliteration_rules_type;
DROP TRIGGER IF EXISTS trigger_update_transliteration_rules_updated_at ON public.transliteration_rules;