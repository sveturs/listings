-- Drop table: user_contacts
DROP SEQUENCE IF EXISTS public.user_contacts_id_seq;
DROP TABLE IF EXISTS public.user_contacts;
DROP INDEX IF EXISTS public.idx_user_contacts_contact_user_id;
DROP INDEX IF EXISTS public.idx_user_contacts_created_at;
DROP INDEX IF EXISTS public.idx_user_contacts_status;
DROP INDEX IF EXISTS public.idx_user_contacts_user_id;
DROP TRIGGER IF EXISTS update_user_contacts_updated_at ON public.user_contacts;