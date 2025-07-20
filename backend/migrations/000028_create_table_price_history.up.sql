-- Migration for table: price_history

CREATE SEQUENCE public.price_history_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.price_history (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    price numeric(12,2) NOT NULL,
    effective_from timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    effective_to timestamp without time zone,
    change_source character varying(50) NOT NULL,
    change_percentage numeric(10,2),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.price_history_id_seq OWNED BY public.price_history.id;

CREATE INDEX idx_price_history_current ON public.price_history USING btree (listing_id, effective_from DESC) WHERE (effective_to IS NULL);

CREATE INDEX idx_price_history_effective ON public.price_history USING btree (listing_id, effective_to);

CREATE INDEX idx_price_history_listing_id ON public.price_history USING btree (listing_id);

ALTER TABLE ONLY public.price_history
    ADD CONSTRAINT price_history_pkey PRIMARY KEY (id);

CREATE TRIGGER trig_update_metadata_after_price_change AFTER INSERT ON public.price_history FOR EACH ROW EXECUTE FUNCTION public.update_listing_metadata_after_price_change();