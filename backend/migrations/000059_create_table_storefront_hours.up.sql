-- Migration for table: storefront_hours

CREATE SEQUENCE public.storefront_hours_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.storefront_hours (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    day_of_week integer NOT NULL,
    open_time time without time zone,
    close_time time without time zone,
    is_closed boolean DEFAULT false,
    special_date date,
    special_note character varying(255),
    CONSTRAINT storefront_hours_day_of_week_check CHECK (((day_of_week >= 0) AND (day_of_week <= 6)))
);

ALTER SEQUENCE public.storefront_hours_id_seq OWNED BY public.storefront_hours.id;

CREATE INDEX idx_hours_storefront_id ON public.storefront_hours USING btree (storefront_id);

ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_storefront_id_day_of_week_special_date_key UNIQUE (storefront_id, day_of_week, special_date);

ALTER TABLE ONLY public.storefront_hours
    ADD CONSTRAINT storefront_hours_storefront_id_fkey FOREIGN KEY (storefront_id) REFERENCES public.storefronts(id) ON DELETE CASCADE;