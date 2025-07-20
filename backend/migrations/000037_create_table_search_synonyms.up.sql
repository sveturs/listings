-- Migration for table: search_synonyms

CREATE SEQUENCE public.search_synonyms_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_synonyms (
    id integer NOT NULL,
    term character varying(255) NOT NULL,
    synonym character varying(255) NOT NULL,
    language character varying(10) DEFAULT 'ru'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.search_synonyms_id_seq OWNED BY public.search_synonyms.id;

CREATE INDEX idx_search_synonyms_active ON public.search_synonyms USING btree (is_active) WHERE (is_active = true);

CREATE INDEX idx_search_synonyms_language ON public.search_synonyms USING btree (language);

CREATE INDEX idx_search_synonyms_synonym ON public.search_synonyms USING btree (synonym);

CREATE INDEX idx_search_synonyms_term ON public.search_synonyms USING btree (term);

CREATE UNIQUE INDEX idx_search_synonyms_unique ON public.search_synonyms USING btree (term, synonym, language) WHERE (is_active = true);

ALTER TABLE ONLY public.search_synonyms
    ADD CONSTRAINT search_synonyms_pkey PRIMARY KEY (id);

CREATE TRIGGER trigger_update_search_synonyms_updated_at BEFORE UPDATE ON public.search_synonyms FOR EACH ROW EXECUTE FUNCTION public.update_search_synonyms_updated_at();