-- Migration for table: import_history

CREATE SEQUENCE public.import_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.import_history (
    id integer NOT NULL,
    source_id integer NOT NULL,
    status character varying(20) NOT NULL,
    items_total integer DEFAULT 0,
    items_imported integer DEFAULT 0,
    items_failed integer DEFAULT 0,
    log text,
    started_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    finished_at timestamp without time zone
);

ALTER SEQUENCE public.import_history_id_seq OWNED BY public.import_history.id;

CREATE INDEX idx_import_history_source ON public.import_history USING btree (source_id);

ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.import_history
    ADD CONSTRAINT import_history_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;