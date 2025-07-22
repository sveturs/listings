-- Migration for table: reviews

CREATE SEQUENCE public.reviews_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.reviews (
    id integer NOT NULL,
    user_id integer NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id integer NOT NULL,
    rating integer NOT NULL,
    comment text,
    pros text,
    cons text,
    photos text[],
    likes_count integer DEFAULT 0,
    helpful_votes integer DEFAULT 0,
    not_helpful_votes integer DEFAULT 0,
    is_verified_purchase boolean DEFAULT false,
    status character varying(20) DEFAULT 'published'::character varying,
    original_language character varying(2) DEFAULT 'en'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    entity_origin_type character varying(50),
    entity_origin_id integer,
    seller_confirmed boolean DEFAULT false,
    has_active_dispute boolean DEFAULT false,
    CONSTRAINT reviews_rating_check CHECK (((rating >= 1) AND (rating <= 5))),
    CONSTRAINT reviews_status_check CHECK (((status)::text = ANY (ARRAY[('draft'::character varying)::text, ('published'::character varying)::text, ('hidden'::character varying)::text])))
);

ALTER SEQUENCE public.reviews_id_seq OWNED BY public.reviews.id;

CREATE INDEX idx_reviews_entity ON public.reviews USING btree (entity_type, entity_id);

CREATE INDEX idx_reviews_entity_origin ON public.reviews USING btree (entity_origin_type, entity_origin_id);

CREATE INDEX idx_reviews_has_dispute ON public.reviews USING btree (has_active_dispute) WHERE (has_active_dispute = true);

CREATE INDEX idx_reviews_rating ON public.reviews USING btree (rating);

CREATE INDEX idx_reviews_seller_confirmed ON public.reviews USING btree (seller_confirmed) WHERE (seller_confirmed = true);

CREATE INDEX idx_reviews_status ON public.reviews USING btree (status);

CREATE INDEX idx_reviews_user ON public.reviews USING btree (user_id);

CREATE UNIQUE INDEX idx_reviews_user_entity_unique ON public.reviews USING btree (user_id, entity_type, entity_id) WHERE ((status)::text <> 'deleted'::text);

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.reviews
    ADD CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;

CREATE TRIGGER refresh_rating_summaries_trigger AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH STATEMENT EXECUTE FUNCTION public.refresh_rating_summaries();

CREATE TRIGGER trigger_refresh_rating_distributions AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH STATEMENT EXECUTE FUNCTION public.refresh_rating_distributions();

CREATE TRIGGER update_ratings_after_review_change AFTER INSERT OR DELETE OR UPDATE ON public.reviews FOR EACH ROW EXECUTE FUNCTION public.refresh_rating_views();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON public.reviews FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();