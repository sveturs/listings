-- Migration for table: search_synonyms_config

CREATE SEQUENCE public.search_synonyms_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_synonyms_config (
    id bigint NOT NULL,
    term character varying(255) NOT NULL,
    synonyms text[] NOT NULL,
    language character varying(10) DEFAULT 'ru'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.search_synonyms_config_id_seq OWNED BY public.search_synonyms_config.id;

CREATE INDEX idx_search_synonyms_config_term ON public.search_synonyms_config USING btree (term);

ALTER TABLE ONLY public.search_synonyms_config
    ADD CONSTRAINT search_synonyms_config_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.search_synonyms_config
    ADD CONSTRAINT search_synonyms_config_term_language_key UNIQUE (term, language);