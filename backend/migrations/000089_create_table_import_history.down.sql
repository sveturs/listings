-- Drop table: import_history
DROP SEQUENCE IF EXISTS public.import_history_id_seq;
DROP TABLE IF EXISTS public.import_history;
DROP INDEX IF EXISTS public.idx_import_history_source;