-- Migration for table: marketplace_categories

CREATE SEQUENCE public.marketplace_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.marketplace_categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    parent_id integer,
    icon character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_custom_ui boolean DEFAULT false,
    custom_ui_component character varying(255),
    sort_order integer DEFAULT 0,
    level integer DEFAULT 0,
    count integer DEFAULT 0,
    external_id character varying(255),
    description text,
    is_active boolean DEFAULT true,
    seo_title character varying(255),
    seo_description text,
    seo_keywords text
);

ALTER SEQUENCE public.marketplace_categories_id_seq OWNED BY public.marketplace_categories.id;

CREATE INDEX idx_marketplace_categories_external_id ON public.marketplace_categories USING btree (external_id);

CREATE INDEX idx_marketplace_categories_parent ON public.marketplace_categories USING btree (parent_id);

CREATE INDEX idx_marketplace_categories_slug ON public.marketplace_categories USING btree (slug);

ALTER TABLE ONLY public.marketplace_categories
    ADD CONSTRAINT marketplace_categories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.marketplace_categories
    ADD CONSTRAINT marketplace_categories_slug_key UNIQUE (slug);

ALTER TABLE ONLY public.marketplace_categories
    ADD CONSTRAINT marketplace_categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES public.marketplace_categories(id);