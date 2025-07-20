-- Migration for table: attribute_group_items

CREATE SEQUENCE public.attribute_group_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.attribute_group_items (
    id integer NOT NULL,
    group_id integer NOT NULL,
    attribute_id integer NOT NULL,
    icon character varying(100),
    sort_order integer DEFAULT 0,
    custom_display_name character varying(255),
    visibility_condition jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.attribute_group_items_id_seq OWNED BY public.attribute_group_items.id;

CREATE INDEX idx_attribute_group_items_attribute ON public.attribute_group_items USING btree (attribute_id);

CREATE INDEX idx_attribute_group_items_group ON public.attribute_group_items USING btree (group_id);

ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_group_id_attribute_id_key UNIQUE (group_id, attribute_id);

ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_attribute_id_fkey FOREIGN KEY (attribute_id) REFERENCES public.category_attributes(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.attribute_group_items
    ADD CONSTRAINT attribute_group_items_group_id_fkey FOREIGN KEY (group_id) REFERENCES public.attribute_groups(id) ON DELETE CASCADE;