-- Migration for table: imported_categories

CREATE SEQUENCE public.imported_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.imported_categories (
    id integer NOT NULL,
    source_id integer NOT NULL,
    source_category character varying(255) NOT NULL,
    category_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.imported_categories_id_seq OWNED BY public.imported_categories.id;

CREATE INDEX idx_imported_categories_source_id ON public.imported_categories USING btree (source_id);

ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_source_id_source_category_key UNIQUE (source_id, source_category);

ALTER TABLE ONLY public.imported_categories
    ADD CONSTRAINT imported_categories_source_id_fkey FOREIGN KEY (source_id) REFERENCES public.import_sources(id) ON DELETE CASCADE;