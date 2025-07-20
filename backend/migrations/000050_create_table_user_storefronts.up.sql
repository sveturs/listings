-- Migration for table: user_storefronts

CREATE SEQUENCE public.user_storefronts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.user_storefronts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    logo_path character varying(255),
    slug character varying(100),
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    creation_transaction_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    phone character varying(50),
    email character varying(255),
    website character varying(255),
    address character varying(255),
    city character varying(100),
    country character varying(100),
    latitude numeric(10,8),
    longitude numeric(11,8)
);

ALTER SEQUENCE public.user_storefronts_id_seq OWNED BY public.user_storefronts.id;

CREATE INDEX idx_user_storefronts_status ON public.user_storefronts USING btree (status);

CREATE INDEX idx_user_storefronts_user ON public.user_storefronts USING btree (user_id);

ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_slug_key UNIQUE (slug);

ALTER TABLE ONLY public.user_storefronts
    ADD CONSTRAINT user_storefronts_creation_transaction_id_fkey FOREIGN KEY (creation_transaction_id) REFERENCES public.balance_transactions(id);