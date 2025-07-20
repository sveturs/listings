-- Drop table: users
DROP SEQUENCE IF EXISTS public.users_id_seq;
DROP TABLE IF EXISTS public.users;
DROP INDEX IF EXISTS public.idx_users_active;
DROP INDEX IF EXISTS public.idx_users_email;
DROP INDEX IF EXISTS public.idx_users_email_lower;
DROP INDEX IF EXISTS public.idx_users_phone;
DROP INDEX IF EXISTS public.idx_users_provider;
DROP INDEX IF EXISTS public.idx_users_status;
DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;