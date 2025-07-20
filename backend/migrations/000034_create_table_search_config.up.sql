-- Migration for table: search_config

CREATE SEQUENCE public.search_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_config (
    id bigint NOT NULL,
    min_search_length integer DEFAULT 2 NOT NULL,
    max_suggestions integer DEFAULT 10 NOT NULL,
    fuzzy_enabled boolean DEFAULT true NOT NULL,
    fuzzy_max_edits integer DEFAULT 2 NOT NULL,
    synonyms_enabled boolean DEFAULT true NOT NULL,
    transliteration_enabled boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.search_config_id_seq OWNED BY public.search_config.id;

ALTER TABLE ONLY public.search_config
    ADD CONSTRAINT search_config_pkey PRIMARY KEY (id);