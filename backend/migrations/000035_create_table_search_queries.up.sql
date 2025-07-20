-- Migration for table: search_queries

CREATE SEQUENCE public.search_queries_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.search_queries (
    id integer NOT NULL,
    query text NOT NULL,
    normalized_query text NOT NULL,
    search_count integer DEFAULT 1 NOT NULL,
    last_searched timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    language character varying(10) DEFAULT 'ru'::character varying NOT NULL,
    results_count integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.search_queries_id_seq OWNED BY public.search_queries.id;

CREATE INDEX idx_search_queries_language ON public.search_queries USING btree (language);

CREATE INDEX idx_search_queries_normalized_query ON public.search_queries USING btree (normalized_query);

CREATE INDEX idx_search_queries_search_count ON public.search_queries USING btree (search_count DESC);

ALTER TABLE ONLY public.search_queries
    ADD CONSTRAINT search_queries_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.search_queries
    ADD CONSTRAINT unique_normalized_query_language UNIQUE (normalized_query, language);

CREATE TRIGGER update_search_queries_updated_at_trigger BEFORE UPDATE ON public.search_queries FOR EACH ROW EXECUTE FUNCTION public.update_search_queries_updated_at();