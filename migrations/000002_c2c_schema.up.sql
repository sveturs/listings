-- Phase 7.4: C2C Schema Migration
-- Migrating c2c_* tables from main svetu database to listings microservice
-- This migration creates 8 tables: categories, chats, favorites, images, listing_variants, listings, messages, orders

-- =============================================================================
-- c2c_categories
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    parent_id integer,
    icon character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_custom_ui boolean DEFAULT false,
    custom_ui_component character varying(255),
    sort_order integer DEFAULT 0,
    level integer DEFAULT 0,
    count integer DEFAULT 0,
    external_id character varying(255),
    description text,
    is_active boolean DEFAULT true,
    seo_title character varying(255),
    seo_description text,
    seo_keywords text,
    CONSTRAINT check_root_categories_level CHECK ((((parent_id IS NULL) AND (level = 0)) OR ((parent_id IS NOT NULL) AND (level > 0))))
);

CREATE SEQUENCE IF NOT EXISTS c2c_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_categories_id_seq OWNED BY c2c_categories.id;
ALTER TABLE ONLY c2c_categories ALTER COLUMN id SET DEFAULT nextval('c2c_categories_id_seq'::regclass);

ALTER TABLE ONLY c2c_categories
    ADD CONSTRAINT c2c_categories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY c2c_categories
    ADD CONSTRAINT c2c_categories_slug_key UNIQUE (slug);

CREATE INDEX IF NOT EXISTS c2c_categories_external_id_idx ON c2c_categories USING btree (external_id);
CREATE INDEX IF NOT EXISTS c2c_categories_parent_id_idx ON c2c_categories USING btree (parent_id);
CREATE INDEX IF NOT EXISTS c2c_categories_slug_idx ON c2c_categories USING btree (slug);

-- =============================================================================
-- c2c_listings
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_listings (
    id integer NOT NULL,
    user_id integer,
    category_id integer,
    title character varying(255) NOT NULL,
    description text,
    price numeric(12,2),
    condition character varying(50),
    status character varying(20) DEFAULT 'active'::character varying,
    location character varying(255),
    latitude numeric(10,8),
    longitude numeric(11,8),
    address_city character varying(100),
    address_country character varying(100),
    views_count integer DEFAULT 0,
    show_on_map boolean DEFAULT true NOT NULL,
    original_language character varying(10) DEFAULT 'sr'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    storefront_id integer,
    external_id character varying(255),
    metadata jsonb,
    needs_reindex boolean DEFAULT false,
    address_multilingual jsonb
);

-- Note: Using a separate sequence for listings as global_product_id_seq is not in listings DB
CREATE SEQUENCE IF NOT EXISTS c2c_listings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_listings_id_seq OWNED BY c2c_listings.id;
ALTER TABLE ONLY c2c_listings ALTER COLUMN id SET DEFAULT nextval('c2c_listings_id_seq'::regclass);

ALTER TABLE ONLY c2c_listings
    ADD CONSTRAINT c2c_listings_pkey PRIMARY KEY (id);

-- Foreign key to categories
ALTER TABLE ONLY c2c_listings
    ADD CONSTRAINT fk_c2c_listings_category_id FOREIGN KEY (category_id) REFERENCES c2c_categories(id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS c2c_listings_address_city_idx ON c2c_listings USING btree (address_city) WHERE (address_city IS NOT NULL);
CREATE INDEX IF NOT EXISTS c2c_listings_address_city_status_idx ON c2c_listings USING btree (address_city, status) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_address_multilingual_idx ON c2c_listings USING gin (address_multilingual);
CREATE INDEX IF NOT EXISTS c2c_listings_category_id_status_idx ON c2c_listings USING btree (category_id, status) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_category_id_status_price_idx ON c2c_listings USING btree (category_id, status, price) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_created_at_idx ON c2c_listings USING btree (created_at DESC) WHERE (((status)::text = 'active'::text) AND ((metadata ->> 'is_featured'::text) = 'true'::text));
CREATE INDEX IF NOT EXISTS c2c_listings_created_at_idx1 ON c2c_listings USING btree (created_at DESC) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_expr_idx ON c2c_listings USING gin (((metadata -> 'discount'::text)));
CREATE INDEX IF NOT EXISTS c2c_listings_external_id_idx ON c2c_listings USING btree (external_id);
CREATE INDEX IF NOT EXISTS c2c_listings_external_id_storefront_id_idx ON c2c_listings USING btree (external_id, storefront_id);
CREATE INDEX IF NOT EXISTS c2c_listings_latitude_longitude_idx ON c2c_listings USING btree (latitude, longitude) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX IF NOT EXISTS c2c_listings_price_idx ON c2c_listings USING btree (price) WHERE ((price IS NOT NULL) AND ((status)::text = 'active'::text));
CREATE INDEX IF NOT EXISTS c2c_listings_status_category_id_idx ON c2c_listings USING btree (status, category_id) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_status_created_at_idx ON c2c_listings USING btree (status, created_at DESC) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_status_created_at_views_count_idx ON c2c_listings USING btree (status, created_at DESC, views_count DESC) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS c2c_listings_status_idx ON c2c_listings USING btree (status);
CREATE INDEX IF NOT EXISTS c2c_listings_storefront_id_idx ON c2c_listings USING btree (storefront_id);
CREATE INDEX IF NOT EXISTS c2c_listings_user_id_status_created_at_idx ON c2c_listings USING btree (user_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS c2c_listings_views_count_idx ON c2c_listings USING btree (views_count DESC) WHERE (((status)::text = 'active'::text) AND (views_count > 10));
CREATE INDEX IF NOT EXISTS idx_c2c_listings_active_created ON c2c_listings USING btree (status, created_at DESC) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS idx_c2c_listings_category_active ON c2c_listings USING btree (category_id, status) WHERE ((status)::text = 'active'::text);
CREATE INDEX IF NOT EXISTS idx_c2c_listings_location ON c2c_listings USING btree (latitude, longitude) WHERE (((status)::text = 'active'::text) AND (latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX IF NOT EXISTS idx_c2c_listings_price ON c2c_listings USING btree (price) WHERE (((status)::text = 'active'::text) AND (price IS NOT NULL));

-- =============================================================================
-- c2c_chats
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_chats (
    id integer NOT NULL,
    listing_id integer,
    buyer_id integer,
    seller_id integer,
    last_message_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_archived boolean DEFAULT false,
    storefront_product_id integer,
    CONSTRAINT check_chat_target CHECK ((NOT ((listing_id IS NOT NULL) AND (storefront_product_id IS NOT NULL))))
);

CREATE SEQUENCE IF NOT EXISTS c2c_chats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_chats_id_seq OWNED BY c2c_chats.id;
ALTER TABLE ONLY c2c_chats ALTER COLUMN id SET DEFAULT nextval('c2c_chats_id_seq'::regclass);

ALTER TABLE ONLY c2c_chats
    ADD CONSTRAINT c2c_chats_pkey PRIMARY KEY (id);

ALTER TABLE ONLY c2c_chats
    ADD CONSTRAINT c2c_chats_listing_id_buyer_id_seller_id_key UNIQUE (listing_id, buyer_id, seller_id);

ALTER TABLE ONLY c2c_chats
    ADD CONSTRAINT c2c_chats_storefront_product_id_buyer_id_seller_id_key UNIQUE (storefront_product_id, buyer_id, seller_id);

CREATE INDEX IF NOT EXISTS c2c_chats_buyer_id_idx ON c2c_chats USING btree (buyer_id);
CREATE INDEX IF NOT EXISTS c2c_chats_buyer_id_seller_id_last_message_at_idx ON c2c_chats USING btree (buyer_id, seller_id, last_message_at DESC);
CREATE INDEX IF NOT EXISTS c2c_chats_is_archived_idx ON c2c_chats USING btree (is_archived) WHERE (NOT is_archived);
CREATE INDEX IF NOT EXISTS c2c_chats_last_message_at_idx ON c2c_chats USING btree (last_message_at DESC) WHERE (NOT is_archived);
CREATE INDEX IF NOT EXISTS c2c_chats_least_greatest_idx ON c2c_chats USING btree (LEAST(buyer_id, seller_id), GREATEST(buyer_id, seller_id));
CREATE UNIQUE INDEX IF NOT EXISTS c2c_chats_least_greatest_idx1 ON c2c_chats USING btree (LEAST(buyer_id, seller_id), GREATEST(buyer_id, seller_id)) WHERE ((listing_id IS NULL) AND (storefront_product_id IS NULL));
CREATE INDEX IF NOT EXISTS c2c_chats_listing_id_buyer_id_seller_id_idx ON c2c_chats USING btree (listing_id, buyer_id, seller_id) WHERE (listing_id IS NOT NULL);
CREATE INDEX IF NOT EXISTS c2c_chats_listing_id_idx ON c2c_chats USING btree (listing_id) WHERE (listing_id IS NOT NULL);
CREATE INDEX IF NOT EXISTS c2c_chats_seller_id_idx ON c2c_chats USING btree (seller_id);
CREATE INDEX IF NOT EXISTS c2c_chats_storefront_product_id_idx ON c2c_chats USING btree (storefront_product_id);
CREATE INDEX IF NOT EXISTS c2c_chats_updated_at_idx ON c2c_chats USING btree (updated_at);

-- =============================================================================
-- c2c_favorites
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_favorites (
    user_id integer NOT NULL,
    listing_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE ONLY c2c_favorites
    ADD CONSTRAINT c2c_favorites_pkey PRIMARY KEY (user_id, listing_id);

ALTER TABLE ONLY c2c_favorites
    ADD CONSTRAINT fk_c2c_favorites_listing_id FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON UPDATE CASCADE ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS c2c_favorites_listing_id_idx ON c2c_favorites USING btree (listing_id);
CREATE INDEX IF NOT EXISTS c2c_favorites_user_id_created_at_idx ON c2c_favorites USING btree (user_id, created_at DESC);

-- =============================================================================
-- c2c_images
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_images (
    id integer NOT NULL,
    listing_id integer,
    file_path character varying(255) NOT NULL,
    file_name character varying(255) NOT NULL,
    file_size integer NOT NULL,
    content_type character varying(100) NOT NULL,
    is_main boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    storage_type character varying(20) DEFAULT 'local'::character varying,
    storage_bucket character varying(100),
    public_url text
);

CREATE SEQUENCE IF NOT EXISTS c2c_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_images_id_seq OWNED BY c2c_images.id;
ALTER TABLE ONLY c2c_images ALTER COLUMN id SET DEFAULT nextval('c2c_images_id_seq'::regclass);

ALTER TABLE ONLY c2c_images
    ADD CONSTRAINT c2c_images_pkey PRIMARY KEY (id);

ALTER TABLE ONLY c2c_images
    ADD CONSTRAINT fk_c2c_images_listing_id FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON UPDATE CASCADE ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS c2c_images_listing_id_is_main_idx ON c2c_images USING btree (listing_id, is_main) WHERE (is_main = true);
CREATE INDEX IF NOT EXISTS idx_c2c_images_listing_main ON c2c_images USING btree (listing_id, is_main DESC);

-- =============================================================================
-- c2c_listing_variants
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_listing_variants (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    sku character varying(100) NOT NULL,
    price numeric(10,2),
    stock integer DEFAULT 0,
    attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    image_url text,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE IF NOT EXISTS c2c_listing_variants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_listing_variants_id_seq OWNED BY c2c_listing_variants.id;
ALTER TABLE ONLY c2c_listing_variants ALTER COLUMN id SET DEFAULT nextval('c2c_listing_variants_id_seq'::regclass);

ALTER TABLE ONLY c2c_listing_variants
    ADD CONSTRAINT c2c_listing_variants_pkey PRIMARY KEY (id);

ALTER TABLE ONLY c2c_listing_variants
    ADD CONSTRAINT c2c_listing_variants_listing_id_sku_key UNIQUE (listing_id, sku);

ALTER TABLE ONLY c2c_listing_variants
    ADD CONSTRAINT fk_c2c_listing_variants_listing_id FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON UPDATE CASCADE ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS c2c_listing_variants_attributes_idx ON c2c_listing_variants USING gin (attributes);
CREATE INDEX IF NOT EXISTS c2c_listing_variants_is_active_idx ON c2c_listing_variants USING btree (is_active);
CREATE INDEX IF NOT EXISTS c2c_listing_variants_listing_id_idx ON c2c_listing_variants USING btree (listing_id);
CREATE INDEX IF NOT EXISTS c2c_listing_variants_sku_idx ON c2c_listing_variants USING btree (sku);

-- =============================================================================
-- c2c_messages
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_messages (
    id integer NOT NULL,
    chat_id integer,
    listing_id integer,
    sender_id integer,
    receiver_id integer,
    content text NOT NULL,
    is_read boolean DEFAULT false,
    original_language character varying(2) DEFAULT 'en'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_attachments boolean DEFAULT false,
    attachments_count integer DEFAULT 0,
    storefront_product_id integer,
    CONSTRAINT check_message_target CHECK ((((listing_id IS NOT NULL) AND (storefront_product_id IS NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NOT NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NULL))))
);

CREATE SEQUENCE IF NOT EXISTS c2c_messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_messages_id_seq OWNED BY c2c_messages.id;
ALTER TABLE ONLY c2c_messages ALTER COLUMN id SET DEFAULT nextval('c2c_messages_id_seq'::regclass);

ALTER TABLE ONLY c2c_messages
    ADD CONSTRAINT c2c_messages_pkey PRIMARY KEY (id);

CREATE INDEX IF NOT EXISTS c2c_messages_chat_id_created_at_idx ON c2c_messages USING btree (chat_id, created_at DESC);
CREATE INDEX IF NOT EXISTS c2c_messages_chat_id_id_idx ON c2c_messages USING btree (chat_id, id DESC);
CREATE INDEX IF NOT EXISTS c2c_messages_chat_id_idx ON c2c_messages USING btree (chat_id);
CREATE INDEX IF NOT EXISTS c2c_messages_chat_id_receiver_id_idx ON c2c_messages USING btree (chat_id, receiver_id) WHERE (NOT is_read);
CREATE INDEX IF NOT EXISTS c2c_messages_created_at_idx ON c2c_messages USING btree (created_at);
CREATE INDEX IF NOT EXISTS c2c_messages_listing_id_idx ON c2c_messages USING btree (listing_id);
CREATE INDEX IF NOT EXISTS c2c_messages_receiver_id_chat_id_idx ON c2c_messages USING btree (receiver_id, chat_id) WHERE (NOT is_read);
CREATE INDEX IF NOT EXISTS c2c_messages_receiver_id_idx ON c2c_messages USING btree (receiver_id);
CREATE INDEX IF NOT EXISTS c2c_messages_receiver_id_is_read_idx ON c2c_messages USING btree (receiver_id, is_read) WHERE (NOT is_read);
CREATE INDEX IF NOT EXISTS c2c_messages_sender_id_idx ON c2c_messages USING btree (sender_id);
CREATE INDEX IF NOT EXISTS c2c_messages_storefront_product_id_idx ON c2c_messages USING btree (storefront_product_id) WHERE (storefront_product_id IS NOT NULL);

-- =============================================================================
-- c2c_orders
-- =============================================================================
CREATE TABLE IF NOT EXISTS c2c_orders (
    id integer NOT NULL,
    buyer_id integer NOT NULL,
    seller_id integer NOT NULL,
    listing_id integer NOT NULL,
    item_price numeric(10,2) NOT NULL,
    platform_fee_rate numeric(5,2) DEFAULT 5.00,
    platform_fee_amount numeric(10,2) NOT NULL,
    seller_payout_amount numeric(10,2) NOT NULL,
    payment_transaction_id integer,
    status character varying(50) DEFAULT 'pending'::character varying,
    protection_period_days integer DEFAULT 7,
    protection_expires_at timestamp with time zone,
    shipping_method character varying(100),
    tracking_number character varying(255),
    shipped_at timestamp with time zone,
    delivered_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    shipment_id bigint,
    shipping_provider character varying(50),
    CONSTRAINT marketplace_orders_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('paid'::character varying)::text, ('shipped'::character varying)::text, ('delivered'::character varying)::text, ('completed'::character varying)::text, ('disputed'::character varying)::text, ('cancelled'::character varying)::text, ('refunded'::character varying)::text])))
);

COMMENT ON COLUMN c2c_orders.tracking_number IS 'Tracking number from delivery microservice (UI only, single source of truth in microservice DB)';
COMMENT ON COLUMN c2c_orders.shipment_id IS 'Shipment ID in delivery microservice DB (links to microservice shipments table)';
COMMENT ON COLUMN c2c_orders.shipping_provider IS 'Delivery provider code: post_express, bex, aks, d_express, city_express (for UI icons/filters)';

CREATE SEQUENCE IF NOT EXISTS c2c_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE c2c_orders_id_seq OWNED BY c2c_orders.id;
ALTER TABLE ONLY c2c_orders ALTER COLUMN id SET DEFAULT nextval('c2c_orders_id_seq'::regclass);

ALTER TABLE ONLY c2c_orders
    ADD CONSTRAINT c2c_orders_pkey PRIMARY KEY (id);

ALTER TABLE ONLY c2c_orders
    ADD CONSTRAINT fk_c2c_orders_listing_id FOREIGN KEY (listing_id) REFERENCES c2c_listings(id) ON UPDATE CASCADE ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS c2c_orders_buyer_id_idx ON c2c_orders USING btree (buyer_id);
CREATE INDEX IF NOT EXISTS c2c_orders_buyer_id_idx1 ON c2c_orders USING btree (buyer_id);
CREATE INDEX IF NOT EXISTS c2c_orders_created_at_idx ON c2c_orders USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS c2c_orders_listing_id_idx ON c2c_orders USING btree (listing_id);
CREATE INDEX IF NOT EXISTS c2c_orders_listing_id_idx1 ON c2c_orders USING btree (listing_id);
CREATE INDEX IF NOT EXISTS c2c_orders_payment_transaction_id_idx ON c2c_orders USING btree (payment_transaction_id);
CREATE INDEX IF NOT EXISTS c2c_orders_protection_expires_at_idx ON c2c_orders USING btree (protection_expires_at) WHERE ((status)::text = ANY (ARRAY[('delivered'::character varying)::text, ('shipped'::character varying)::text]));
CREATE INDEX IF NOT EXISTS c2c_orders_protection_expires_at_idx1 ON c2c_orders USING btree (protection_expires_at) WHERE (protection_expires_at IS NOT NULL);
CREATE INDEX IF NOT EXISTS c2c_orders_seller_id_idx ON c2c_orders USING btree (seller_id);
CREATE INDEX IF NOT EXISTS c2c_orders_seller_id_idx1 ON c2c_orders USING btree (seller_id);
CREATE INDEX IF NOT EXISTS c2c_orders_status_idx ON c2c_orders USING btree (status);
CREATE INDEX IF NOT EXISTS idx_c2c_orders_shipment_id ON c2c_orders USING btree (shipment_id) WHERE (shipment_id IS NOT NULL);
CREATE INDEX IF NOT EXISTS idx_c2c_orders_tracking_number ON c2c_orders USING btree (tracking_number) WHERE (tracking_number IS NOT NULL);
