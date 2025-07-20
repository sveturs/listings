-- Migration for table: storefront_staff

CREATE SEQUENCE public.storefront_staff_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_staff (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    user_id integer NOT NULL,
    role character varying(50) DEFAULT 'staff'::character varying NOT NULL,
    permissions jsonb DEFAULT '{}'::jsonb,
    last_active_at timestamp without time zone,
    actions_count integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.storefront_staff_id_seq OWNED BY public.storefront_staff.id;

CREATE INDEX idx_staff_storefront_id ON public.storefront_staff USING btree (storefront_id);

CREATE INDEX idx_staff_user_id ON public.storefront_staff USING btree (user_id);

ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_storefront_id_user_id_key UNIQUE (storefront_id, user_id);

ALTER TABLE ONLY public.storefront_staff
    ADD CONSTRAINT storefront_staff_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;