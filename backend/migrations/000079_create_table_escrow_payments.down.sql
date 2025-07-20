-- Drop table: escrow_payments
DROP SEQUENCE IF EXISTS public.escrow_payments_id_seq;
DROP TABLE IF EXISTS public.escrow_payments;
DROP INDEX IF EXISTS public.idx_escrow_payments_buyer_id;
DROP INDEX IF EXISTS public.idx_escrow_payments_payment_transaction_id;
DROP INDEX IF EXISTS public.idx_escrow_payments_seller_id;
DROP INDEX IF EXISTS public.idx_escrow_payments_status;