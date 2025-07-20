-- Migration for table: payment_methods

CREATE SEQUENCE public.payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.payment_methods (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    type character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    minimum_amount numeric(12,2),
    maximum_amount numeric(12,2),
    fee_percentage numeric(5,2),
    fixed_fee numeric(12,2),
    credentials jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.payment_methods_id_seq OWNED BY public.payment_methods.id;

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_code_key UNIQUE (code);

ALTER TABLE ONLY public.payment_methods
    ADD CONSTRAINT payment_methods_pkey PRIMARY KEY (id);