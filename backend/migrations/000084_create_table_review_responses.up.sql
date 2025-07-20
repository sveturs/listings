-- Migration for table: review_responses

CREATE SEQUENCE public.review_responses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.review_responses (
    id integer NOT NULL,
    review_id integer NOT NULL,
    user_id integer NOT NULL,
    response text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER SEQUENCE public.review_responses_id_seq OWNED BY public.review_responses.id;

CREATE INDEX idx_review_responses_review ON public.review_responses USING btree (review_id);

ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.review_responses
    ADD CONSTRAINT review_responses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

CREATE TRIGGER update_review_responses_updated_at BEFORE UPDATE ON public.review_responses FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();