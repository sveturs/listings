-- Migration for table: review_votes

CREATE TABLE public.review_votes (
    review_id integer NOT NULL,
    user_id integer NOT NULL,
    vote_type character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT review_votes_vote_type_check CHECK (((vote_type)::text = ANY (ARRAY[('helpful'::character varying)::text, ('not_helpful'::character varying)::text])))
);

ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_pkey PRIMARY KEY (review_id, user_id);

ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_review_id_fkey FOREIGN KEY (review_id) REFERENCES public.reviews(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.review_votes
    ADD CONSTRAINT review_votes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;