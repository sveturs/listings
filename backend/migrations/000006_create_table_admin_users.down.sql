-- Drop table: admin_users
DROP SEQUENCE IF EXISTS public.admin_users_id_seq;
DROP TABLE IF EXISTS public.admin_users;
DROP INDEX IF EXISTS public.admin_users_email_idx;