-- Migration for table: address_change_log

CREATE SEQUENCE public.address_change_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.address_change_log (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    user_id bigint NOT NULL,
    old_address text,
    new_address text,
    old_location public.geography(Point,4326),
    new_location public.geography(Point,4326),
    change_reason character varying(100),
    confidence_before numeric(3,2),
    confidence_after numeric(3,2),
    ip_address inet,
    user_agent text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER SEQUENCE public.address_change_log_id_seq OWNED BY public.address_change_log.id;

CREATE INDEX idx_address_log_change_reason ON public.address_change_log USING btree (change_reason);

CREATE INDEX idx_address_log_confidence_after ON public.address_change_log USING btree (confidence_after);

CREATE INDEX idx_address_log_created_at ON public.address_change_log USING btree (created_at);

CREATE INDEX idx_address_log_listing_id ON public.address_change_log USING btree (listing_id);

CREATE INDEX idx_address_log_new_location ON public.address_change_log USING gist (new_location);

CREATE INDEX idx_address_log_old_location ON public.address_change_log USING gist (old_location);

CREATE INDEX idx_address_log_user_id ON public.address_change_log USING btree (user_id);

ALTER TABLE ONLY public.address_change_log
    ADD CONSTRAINT address_change_log_pkey PRIMARY KEY (id);