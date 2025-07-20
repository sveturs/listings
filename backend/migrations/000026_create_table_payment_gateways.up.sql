-- Migration for table: payment_gateways

CREATE SEQUENCE public.payment_gateways_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.payment_gateways (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_active boolean DEFAULT true,
    config jsonb NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER SEQUENCE public.payment_gateways_id_seq OWNED BY public.payment_gateways.id;

ALTER TABLE ONLY public.payment_gateways
    ADD CONSTRAINT payment_gateways_pkey PRIMARY KEY (id);