-- Migration for table: attribute_groups

CREATE SEQUENCE public.attribute_groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.attribute_groups (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(255) NOT NULL,
    description text,
    icon character varying(100),
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    is_system boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.attribute_groups_id_seq OWNED BY public.attribute_groups.id;

CREATE INDEX idx_attribute_groups_active ON public.attribute_groups USING btree (is_active);

CREATE INDEX idx_attribute_groups_name ON public.attribute_groups USING btree (name);

ALTER TABLE ONLY public.attribute_groups
    ADD CONSTRAINT attribute_groups_name_key UNIQUE (name);

ALTER TABLE ONLY public.attribute_groups
    ADD CONSTRAINT attribute_groups_pkey PRIMARY KEY (id);

CREATE TRIGGER update_attribute_groups_updated_at BEFORE UPDATE ON public.attribute_groups FOR EACH ROW EXECUTE FUNCTION public.update_attribute_groups_updated_at();