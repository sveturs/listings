-- Drop table: user_privacy_settings
DROP TABLE IF EXISTS public.user_privacy_settings;
DROP INDEX IF EXISTS public.idx_user_privacy_settings_user_id;
DROP TRIGGER IF EXISTS update_user_privacy_settings_updated_at ON public.user_privacy_settings;