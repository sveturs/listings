-- Migration for table: marketplace_favorites

CREATE TABLE public.marketplace_favorites (
    user_id integer NOT NULL,
    listing_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_marketplace_favorites_listing ON public.marketplace_favorites USING btree (listing_id);

CREATE INDEX idx_marketplace_favorites_user_count ON public.marketplace_favorites USING btree (user_id, created_at DESC);

ALTER TABLE ONLY public.marketplace_favorites
    ADD CONSTRAINT marketplace_favorites_pkey PRIMARY KEY (user_id, listing_id);

ALTER TABLE ONLY public.marketplace_favorites
    ADD CONSTRAINT marketplace_favorites_listing_id_fkey FOREIGN KEY (listing_id) REFERENCES public.marketplace_listings(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.marketplace_favorites
    ADD CONSTRAINT marketplace_favorites_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);