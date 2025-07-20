-- Migration for table: category_attributes

CREATE SEQUENCE public.category_attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.category_attributes (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(255) NOT NULL,
    attribute_type character varying(50) NOT NULL,
    options jsonb DEFAULT '{}'::jsonb,
    validation_rules jsonb,
    is_searchable boolean DEFAULT true,
    is_filterable boolean DEFAULT true,
    is_required boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    custom_component character varying(255),
    show_in_card boolean DEFAULT true,
    show_in_list boolean DEFAULT false,
    icon character varying(10) DEFAULT ''::character varying
);

ALTER SEQUENCE public.category_attributes_id_seq OWNED BY public.category_attributes.id;

CREATE INDEX idx_category_attributes_name ON public.category_attributes USING btree (name);

ALTER TABLE ONLY public.category_attributes
    ADD CONSTRAINT category_attributes_pkey PRIMARY KEY (id);