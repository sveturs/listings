-- Migration for table: category_attribute_groups

CREATE SEQUENCE public.category_attribute_groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.category_attribute_groups (
    id integer NOT NULL,
    category_id integer NOT NULL,
    group_id integer NOT NULL,
    component_id integer,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    display_mode character varying(50) DEFAULT 'list'::character varying,
    collapsed_by_default boolean DEFAULT false,
    configuration jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.category_attribute_groups_id_seq OWNED BY public.category_attribute_groups.id;

CREATE INDEX idx_category_attribute_groups_category ON public.category_attribute_groups USING btree (category_id);

CREATE INDEX idx_category_attribute_groups_component ON public.category_attribute_groups USING btree (component_id);

CREATE INDEX idx_category_attribute_groups_group ON public.category_attribute_groups USING btree (group_id);

ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_category_id_group_id_key UNIQUE (category_id, group_id);

ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.category_attribute_groups
    ADD CONSTRAINT category_attribute_groups_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.attribute_groups(id) ON DELETE CASCADE;