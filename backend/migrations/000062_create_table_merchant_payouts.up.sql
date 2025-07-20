-- Migration for table: merchant_payouts

CREATE SEQUENCE public.merchant_payouts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.merchant_payouts (
    id bigint NOT NULL,
    seller_id integer,
    gateway_id integer,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying,
    gateway_payout_id character varying(255),
    gateway_reference_id character varying(255),
    status character varying(50) DEFAULT 'pending'::character varying,
    bank_account_info jsonb,
    gateway_response jsonb,
    error_details jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    processed_at timestamp with time zone,
    CONSTRAINT merchant_payouts_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT merchant_payouts_status_valid CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'processing'::character varying, 'completed'::character varying, 'failed'::character varying])::text[])))
);

ALTER SEQUENCE public.merchant_payouts_id_seq OWNED BY public.merchant_payouts.id;

CREATE INDEX idx_merchant_payouts_gateway_payout_id ON public.merchant_payouts USING btree (gateway_payout_id);

CREATE INDEX idx_merchant_payouts_seller_id ON public.merchant_payouts USING btree (seller_id);

CREATE INDEX idx_merchant_payouts_status ON public.merchant_payouts USING btree (status);

ALTER TABLE ONLY public.merchant_payouts
    ADD CONSTRAINT merchant_payouts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.merchant_payouts
    ADD CONSTRAINT merchant_payouts_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.payment_gateways(id);