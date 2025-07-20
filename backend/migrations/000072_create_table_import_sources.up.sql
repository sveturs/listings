-- Migration for table: import_sources

CREATE SEQUENCE public.import_sources_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.import_sources (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    type character varying(20) NOT NULL,
    url character varying(512),
    auth_data jsonb,
    schedule character varying(50),
    mapping jsonb,
    last_import_at timestamp without time zone,
    last_import_status character varying(20),
    last_import_log text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.import_sources_id_seq OWNED BY public.import_sources.id;

CREATE INDEX idx_import_sources_storefront ON public.import_sources USING btree (storefront_id);

ALTER TABLE ONLY public.import_sources
    ADD CONSTRAINT import_sources_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.import_sources
    ADD CONSTRAINT import_sources_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.user_storefronts(id) ON DELETE CASCADE;