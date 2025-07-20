-- Drop table: payment_transactions
DROP SEQUENCE IF EXISTS public.payment_transactions_id_seq;
DROP TABLE IF EXISTS public.payment_transactions;
DROP INDEX IF EXISTS public.idx_payment_transactions_auto_capture;
DROP INDEX IF EXISTS public.idx_payment_transactions_created_at;
DROP INDEX IF EXISTS public.idx_payment_transactions_gateway_transaction_id;
DROP INDEX IF EXISTS public.idx_payment_transactions_listing_id;
DROP INDEX IF EXISTS public.idx_payment_transactions_order_reference;
DROP INDEX IF EXISTS public.idx_payment_transactions_source;
DROP INDEX IF EXISTS public.idx_payment_transactions_status;
DROP INDEX IF EXISTS public.idx_payment_transactions_user_id;