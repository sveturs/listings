-- Drop table: merchant_payouts
DROP SEQUENCE IF EXISTS public.merchant_payouts_id_seq;
DROP TABLE IF EXISTS public.merchant_payouts;
DROP INDEX IF EXISTS public.idx_merchant_payouts_gateway_payout_id;
DROP INDEX IF EXISTS public.idx_merchant_payouts_seller_id;
DROP INDEX IF EXISTS public.idx_merchant_payouts_status;