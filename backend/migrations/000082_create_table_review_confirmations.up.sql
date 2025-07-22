-- Migration for table: review_confirmations

CREATE SEQUENCE public.review_confirmations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.review_confirmations (
    id integer NOT NULL,
    review_id integer NOT NULL,
    confirmed_by integer NOT NULL,
    confirmation_status character varying(50) NOT NULL,
    confirmed_at timestamp without time zone DEFAULT now() NOT NULL,
    notes text,
    CONSTRAINT review_confirmations_confirmation_status_check CHECK (((confirmation_status)::text = ANY (ARRAY[('confirmed'::character varying)::text, ('disputed'::character varying)::text])))
);

ALTER SEQUENCE public.review_confirmations_id_seq OWNED BY public.review_confirmations.id;

ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_review_id_key UNIQUE (review_id);

ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_confirmed_by_fkey FOREIGN KEY (confirmed_by) REFERENCES public.users(id);

ALTER TABLE ONLY public.review_confirmations
    ADD CONSTRAINT review_confirmations_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;