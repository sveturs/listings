-- Drop table: balance_transactions
DROP SEQUENCE IF EXISTS public.balance_transactions_id_seq;
DROP TABLE IF EXISTS public.balance_transactions;
DROP INDEX IF EXISTS public.idx_transactions_created;
DROP INDEX IF EXISTS public.idx_transactions_status;
DROP INDEX IF EXISTS public.idx_transactions_user;