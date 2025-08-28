-- Migration for table: map_items_cache
-- Cache table for map items to improve performance

CREATE SEQUENCE public.map_items_cache_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.map_items_cache (
    id integer NOT NULL,
    item_type VARCHAR(50) NOT NULL,
    item_id integer NOT NULL,
    latitude numeric(10,7) NOT NULL,
    longitude numeric(10,7) NOT NULL,
    title TEXT NOT NULL,
    price numeric(20,5),
    category_id integer,
    category_name VARCHAR(255),
    image_url TEXT,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone
);

ALTER SEQUENCE public.map_items_cache_id_seq OWNED BY public.map_items_cache.id;

-- Primary key
ALTER TABLE ONLY public.map_items_cache
    ADD CONSTRAINT map_items_cache_pkey PRIMARY KEY (id);

-- Indexes for performance
CREATE INDEX idx_map_items_cache_type_id ON public.map_items_cache USING btree (item_type, item_id);
CREATE INDEX idx_map_items_cache_location ON public.map_items_cache USING btree (latitude, longitude);
CREATE INDEX idx_map_items_cache_expires ON public.map_items_cache USING btree (expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_map_items_cache_category ON public.map_items_cache USING btree (category_id) WHERE category_id IS NOT NULL;

-- Unique constraint to prevent duplicates
CREATE UNIQUE INDEX idx_map_items_cache_unique ON public.map_items_cache USING btree (item_type, item_id);