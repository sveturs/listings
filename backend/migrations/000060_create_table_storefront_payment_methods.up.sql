-- Migration for table: storefront_payment_methods

CREATE SEQUENCE public.storefront_payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_payment_methods (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    method_type character varying(50) NOT NULL,
    is_enabled boolean DEFAULT true,
    provider character varying(50),
    settings jsonb DEFAULT '{}'::jsonb,
    transaction_fee numeric(5,2) DEFAULT 0.00,
    min_amount numeric(10,2),
    max_amount numeric(10,2),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.storefront_payment_methods_id_seq OWNED BY public.storefront_payment_methods.id;

CREATE INDEX idx_payment_is_enabled ON public.storefront_payment_methods USING btree (is_enabled);

CREATE INDEX idx_payment_storefront_id ON public.storefront_payment_methods USING btree (storefront_id);

ALTER TABLE ONLY public.storefront_payment_methods
    ADD CONSTRAINT storefront_payment_methods_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_payment_methods
    ADD CONSTRAINT storefront_payment_methods_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;