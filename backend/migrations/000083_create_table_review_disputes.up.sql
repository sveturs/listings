-- Migration for table: review_disputes

CREATE SEQUENCE public.review_disputes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.review_disputes (
    id integer NOT NULL,
    review_id integer NOT NULL,
    disputed_by integer NOT NULL,
    dispute_reason character varying(100) NOT NULL,
    dispute_description text NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    admin_id integer,
    admin_notes text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    resolved_at timestamp without time zone,
    CONSTRAINT review_disputes_dispute_reason_check CHECK (((dispute_reason)::text = ANY ((ARRAY['not_a_customer'::character varying, 'false_information'::character varying, 'deal_cancelled'::character varying, 'spam'::character varying, 'other'::character varying])::text[]))),
    CONSTRAINT review_disputes_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'in_review'::character varying, 'resolved_keep_review'::character varying, 'resolved_remove_review'::character varying, 'resolved_remove_verification'::character varying, 'cancelled'::character varying])::text[])))
);

ALTER SEQUENCE public.review_disputes_id_seq OWNED BY public.review_disputes.id;

CREATE INDEX idx_disputes_review_id ON public.review_disputes USING btree (review_id);

CREATE INDEX idx_disputes_status ON public.review_disputes USING btree (status);

ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_admin_id_fkey FOREIGN KEY (admin_id) REFERENCES public.users(id);

ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_disputed_by_fkey FOREIGN KEY (disputed_by) REFERENCES public.users(id);

ALTER TABLE ONLY public.review_disputes
    ADD CONSTRAINT review_disputes_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;