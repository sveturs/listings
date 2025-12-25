CREATE TABLE public.chats (
    id bigint NOT NULL,
    buyer_id bigint NOT NULL,
    seller_id bigint NOT NULL,
    listing_id bigint,
    storefront_product_id bigint,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    is_archived boolean DEFAULT false NOT NULL,
    last_message_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT check_chat_context CHECK ((((listing_id IS NOT NULL) AND (storefront_product_id IS NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NOT NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NULL)))),
    CONSTRAINT check_chat_status CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'archived'::character varying, 'blocked'::character varying])::text[]))),
    CONSTRAINT check_participants CHECK ((buyer_id <> seller_id))
);
CREATE TABLE public.indexing_queue (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    operation character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    retry_count integer DEFAULT 0 NOT NULL,
    max_retries integer DEFAULT 3 NOT NULL,
    error_message text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    processed_at timestamp with time zone,
    CONSTRAINT indexing_queue_operation_check CHECK (((operation)::text = ANY ((ARRAY['index'::character varying, 'update'::character varying, 'delete'::character varying])::text[]))),
    CONSTRAINT indexing_queue_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'processing'::character varying, 'completed'::character varying, 'failed'::character varying])::text[])))
);
CREATE TABLE public.inventory_movements (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    variant_id bigint,
    movement_type character varying(50) NOT NULL,
    quantity integer NOT NULL,
    reason character varying(255),
    notes text,
    user_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    metadata jsonb,
    CONSTRAINT inventory_movements_quantity_check CHECK ((quantity <> 0))
);
CREATE TABLE public.inventory_reservations (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    variant_id bigint,
    reference_id bigint,
    quantity integer NOT NULL,
    status character varying(20) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    committed_at timestamp with time zone,
    released_at timestamp with time zone,
    reference_type character varying(20) DEFAULT 'order'::character varying NOT NULL,
    CONSTRAINT chk_inventory_reservations_quantity_positive CHECK ((quantity > 0)),
    CONSTRAINT chk_inventory_reservations_reference_type CHECK (((reference_type)::text = ANY ((ARRAY['order'::character varying, 'transfer'::character varying])::text[]))),
    CONSTRAINT chk_inventory_reservations_status CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'committed'::character varying, 'released'::character varying, 'expired'::character varying])::text[])))
);
CREATE TABLE public.listing_attribute_values (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    attribute_id integer NOT NULL,
    value_text text,
    value_number numeric(20,4),
    value_boolean boolean,
    value_date date,
    value_json jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.listing_attributes (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    attribute_key character varying(100) NOT NULL,
    attribute_value text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.listing_images (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    url text NOT NULL,
    storage_path text,
    thumbnail_url text,
    display_order integer DEFAULT 0 NOT NULL,
    is_primary boolean DEFAULT false NOT NULL,
    width integer,
    height integer,
    file_size bigint,
    mime_type character varying(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.listing_locations (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    country character varying(100),
    city character varying(100),
    postal_code character varying(20),
    address_line1 text,
    address_line2 text,
    latitude numeric(10,8),
    longitude numeric(11,8),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.listing_tags (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    tag character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.listing_variants (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    sku character varying(100),
    attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    price numeric(12,2),
    stock integer DEFAULT 0 NOT NULL,
    image_url character varying(500),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.messages (
    id bigint NOT NULL,
    chat_id bigint NOT NULL,
    sender_id bigint NOT NULL,
    receiver_id bigint NOT NULL,
    content text NOT NULL,
    original_language character varying(2) DEFAULT 'en'::character varying NOT NULL,
    listing_id bigint,
    storefront_product_id bigint,
    status character varying(20) DEFAULT 'sent'::character varying NOT NULL,
    is_read boolean DEFAULT false NOT NULL,
    has_attachments boolean DEFAULT false NOT NULL,
    attachments_count integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    read_at timestamp without time zone,
    is_system boolean DEFAULT false NOT NULL,
    CONSTRAINT check_attachments_consistency CHECK ((((has_attachments = true) AND (attachments_count > 0)) OR ((has_attachments = false) AND (attachments_count = 0)))),
    CONSTRAINT check_content_length CHECK (((length(content) >= 1) AND (length(content) <= 10000))),
    CONSTRAINT check_message_context CHECK ((((listing_id IS NOT NULL) AND (storefront_product_id IS NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NOT NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NULL)))),
    CONSTRAINT check_message_status CHECK (((status)::text = ANY ((ARRAY['sent'::character varying, 'delivered'::character varying, 'read'::character varying, 'failed'::character varying])::text[]))),
    CONSTRAINT check_original_language CHECK (((original_language)::text ~ '^[a-z]{2}$'::text))
);
CREATE TABLE public.order_items (
    id bigint NOT NULL,
    order_id bigint NOT NULL,
    listing_id bigint NOT NULL,
    variant_id bigint,
    listing_name character varying(255) NOT NULL,
    sku character varying(100),
    variant_data jsonb,
    attributes jsonb,
    quantity integer NOT NULL,
    price numeric(10,2) NOT NULL,
    total numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    subtotal numeric(10,2) DEFAULT 0 NOT NULL,
    discount numeric(10,2) DEFAULT 0 NOT NULL,
    image_url character varying(500),
    stock_reservation_id text,
    variant_uuid uuid,
    CONSTRAINT chk_order_items_discount_non_negative CHECK ((discount >= (0)::numeric)),
    CONSTRAINT chk_order_items_price_non_negative CHECK ((price >= (0)::numeric)),
    CONSTRAINT chk_order_items_quantity_positive CHECK ((quantity > 0)),
    CONSTRAINT chk_order_items_subtotal_non_negative CHECK ((subtotal >= (0)::numeric)),
    CONSTRAINT chk_order_items_total_non_negative CHECK ((total >= (0)::numeric))
);
CREATE TABLE public.search_queries (
    id bigint NOT NULL,
    query_text character varying(500) NOT NULL,
    category_id bigint,
    user_id bigint,
    session_id character varying(255),
    results_count integer DEFAULT 0 NOT NULL,
    clicked_listing_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT chk_search_queries_user_or_session CHECK ((((user_id IS NOT NULL) AND (session_id IS NOT NULL)) OR ((user_id IS NOT NULL) AND (session_id IS NULL)) OR ((user_id IS NULL) AND (session_id IS NOT NULL)))),
    CONSTRAINT search_queries_query_text_check CHECK ((length(TRIM(BOTH FROM query_text)) >= 1)),
    CONSTRAINT search_queries_results_count_check CHECK ((results_count >= 0))
);
CREATE TABLE public.shopping_carts (
    id bigint NOT NULL,
    user_id bigint,
    session_id character varying(255),
    storefront_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_shopping_carts_user_or_session CHECK ((((user_id IS NOT NULL) AND (session_id IS NULL)) OR ((user_id IS NULL) AND (session_id IS NOT NULL))))
);
CREATE TABLE public.storefront_delivery_options (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    base_price numeric(10,2) DEFAULT 0.00,
    price_per_km numeric(10,2) DEFAULT 0.00,
    price_per_kg numeric(10,2) DEFAULT 0.00,
    free_above_amount numeric(10,2),
    min_order_amount numeric(10,2),
    max_weight_kg numeric(10,2),
    max_distance_km numeric(10,2),
    estimated_days_min integer DEFAULT 1,
    estimated_days_max integer DEFAULT 3,
    zones jsonb DEFAULT '[]'::jsonb,
    available_days jsonb DEFAULT '[1, 2, 3, 4, 5]'::jsonb,
    cutoff_time time without time zone,
    provider character varying(50),
    provider_config jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true,
    display_order integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.storefront_events (
    id bigint NOT NULL,
    storefront_id integer NOT NULL,
    event_type character varying(50) NOT NULL,
    event_data jsonb DEFAULT '{}'::jsonb,
    user_id integer,
    session_id character varying(100) NOT NULL,
    ip_address character varying(45),
    user_agent text,
    referrer text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.storefront_hours (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    day_of_week integer NOT NULL,
    open_time time without time zone,
    close_time time without time zone,
    is_closed boolean DEFAULT false,
    special_date date,
    special_note character varying(255),
    CONSTRAINT storefront_hours_day_of_week_check CHECK (((day_of_week >= 0) AND (day_of_week <= 6)))
);
CREATE TABLE public.storefront_invitations (
    id bigint NOT NULL,
    storefront_id bigint NOT NULL,
    role character varying(20) DEFAULT 'staff'::character varying NOT NULL,
    type public.storefront_invitation_type NOT NULL,
    invited_email character varying(255),
    invited_user_id bigint,
    invite_code character varying(32),
    expires_at timestamp with time zone,
    max_uses integer,
    current_uses integer DEFAULT 0,
    invited_by_id bigint NOT NULL,
    status public.storefront_invitation_status DEFAULT 'pending'::public.storefront_invitation_status NOT NULL,
    comment text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    accepted_at timestamp with time zone,
    declined_at timestamp with time zone,
    CONSTRAINT valid_sf_email_invitation CHECK (((type <> 'email'::public.storefront_invitation_type) OR (invited_email IS NOT NULL))),
    CONSTRAINT valid_sf_link_invitation CHECK (((type <> 'link'::public.storefront_invitation_type) OR (invite_code IS NOT NULL))),
    CONSTRAINT valid_sf_role CHECK (((role)::text = ANY ((ARRAY['owner'::character varying, 'manager'::character varying, 'staff'::character varying, 'cashier'::character varying])::text[])))
);
CREATE TABLE public.storefront_payment_methods (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    method_type character varying(50) NOT NULL,
    is_enabled boolean DEFAULT true,
    provider character varying(50),
    settings jsonb DEFAULT '{}'::jsonb,
    transaction_fee numeric(5,2) DEFAULT 0.00,
    min_amount numeric(10,2),
    max_amount numeric(10,2),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.storefront_staff (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    user_id integer NOT NULL,
    role character varying(50) DEFAULT 'staff'::character varying NOT NULL,
    permissions jsonb DEFAULT '{}'::jsonb,
    last_active_at timestamp without time zone,
    actions_count integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    invitation_id bigint
);
CREATE TABLE public.users (
    id bigint NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(255),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
ALTER TABLE ONLY public.attribute_options ALTER COLUMN id SET DEFAULT nextval('public.attribute_options_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_search_cache ALTER COLUMN id SET DEFAULT nextval('public.attribute_search_cache_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_values ALTER COLUMN id SET DEFAULT nextval('public.attribute_values_id_seq'::regclass);
ALTER TABLE ONLY public.attributes ALTER COLUMN id SET DEFAULT nextval('public.attributes_id_seq'::regclass);
ALTER TABLE ONLY public.c2c_chats ALTER COLUMN id SET DEFAULT nextval('public.c2c_chats_id_seq'::regclass);
ALTER TABLE ONLY public.c2c_messages ALTER COLUMN id SET DEFAULT nextval('public.c2c_messages_id_seq'::regclass);
ALTER TABLE ONLY public.cart_items ALTER COLUMN id SET DEFAULT nextval('public.cart_items_id_seq'::regclass);
ALTER TABLE ONLY public.category_attributes ALTER COLUMN id SET DEFAULT nextval('public.category_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.category_variant_attributes ALTER COLUMN id SET DEFAULT nextval('public.category_variant_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.chat_attachments ALTER COLUMN id SET DEFAULT nextval('public.chat_attachments_id_seq'::regclass);
ALTER TABLE ONLY public.chats ALTER COLUMN id SET DEFAULT nextval('public.chats_id_seq'::regclass);
ALTER TABLE ONLY public.indexing_queue ALTER COLUMN id SET DEFAULT nextval('public.indexing_queue_id_seq'::regclass);
ALTER TABLE ONLY public.inventory_movements ALTER COLUMN id SET DEFAULT nextval('public.inventory_movements_id_seq'::regclass);
ALTER TABLE ONLY public.inventory_reservations ALTER COLUMN id SET DEFAULT nextval('public.inventory_reservations_id_seq'::regclass);
ALTER TABLE ONLY public.listing_attribute_values ALTER COLUMN id SET DEFAULT nextval('public.listing_attribute_values_id_seq'::regclass);
ALTER TABLE ONLY public.listing_attributes ALTER COLUMN id SET DEFAULT nextval('public.listing_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.listing_images ALTER COLUMN id SET DEFAULT nextval('public.listing_images_id_seq'::regclass);
ALTER TABLE ONLY public.listing_locations ALTER COLUMN id SET DEFAULT nextval('public.listing_locations_id_seq'::regclass);
ALTER TABLE ONLY public.listing_tags ALTER COLUMN id SET DEFAULT nextval('public.listing_tags_id_seq'::regclass);
ALTER TABLE ONLY public.listing_variants ALTER COLUMN id SET DEFAULT nextval('public.listing_variants_id_seq'::regclass);
ALTER TABLE ONLY public.listings ALTER COLUMN id SET DEFAULT nextval('public.listings_id_seq'::regclass);
ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);
ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);
ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);
ALTER TABLE ONLY public.search_queries ALTER COLUMN id SET DEFAULT nextval('public.search_queries_id_seq'::regclass);
ALTER TABLE ONLY public.shopping_carts ALTER COLUMN id SET DEFAULT nextval('public.shopping_carts_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_delivery_options ALTER COLUMN id SET DEFAULT nextval('public.storefront_delivery_options_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_events ALTER COLUMN id SET DEFAULT nextval('public.storefront_events_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_hours ALTER COLUMN id SET DEFAULT nextval('public.storefront_hours_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_invitations ALTER COLUMN id SET DEFAULT nextval('public.storefront_invitations_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_payment_methods ALTER COLUMN id SET DEFAULT nextval('public.storefront_payment_methods_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_staff ALTER COLUMN id SET DEFAULT nextval('public.storefront_staff_id_seq'::regclass);
ALTER TABLE ONLY public.storefronts ALTER COLUMN id SET DEFAULT nextval('public.storefronts_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
CREATE MATERIALIZED VIEW public.analytics_storefront_stats AS
 SELECT s.id AS storefront_id,
    s.name AS storefront_name,
    s.user_id AS owner_id,
    count(DISTINCT l.id) FILTER (WHERE (((l.status)::text = 'active'::text) AND (l.is_deleted = false))) AS active_listings,
    count(DISTINCT l.id) FILTER (WHERE (l.is_deleted = false)) AS total_listings,
    count(DISTINCT o.id) FILTER (WHERE ((o.status)::text = 'delivered'::text)) AS total_sales,
    COALESCE(sum(o.total) FILTER (WHERE ((o.status)::text = 'delivered'::text)), (0)::numeric) AS total_revenue,
        CASE
            WHEN (count(DISTINCT o.id) FILTER (WHERE ((o.status)::text = 'delivered'::text)) > 0) THEN (COALESCE(sum(o.total) FILTER (WHERE ((o.status)::text = 'delivered'::text)), (0)::numeric) / (count(DISTINCT o.id) FILTER (WHERE ((o.status)::text = 'delivered'::text)))::numeric)
            ELSE (0)::numeric
        END AS average_order_value,
    COALESCE(sum(l.view_count), (0)::bigint) AS total_views,
    COALESCE(sum(l.favorites_count), (0)::bigint) AS total_favorites,
        CASE
            WHEN (sum(l.view_count) > 0) THEN (((count(DISTINCT o.id) FILTER (WHERE ((o.status)::text = 'delivered'::text)))::numeric / (sum(l.view_count))::numeric) * (100)::numeric)
            ELSE (0)::numeric
        END AS conversion_rate,
    now() AS last_updated_at
   FROM ((public.storefronts s
     LEFT JOIN public.listings l ON (((l.storefront_id = s.id) AND (l.is_deleted = false))))
     LEFT JOIN public.orders o ON ((o.storefront_id = s.id)))
  WHERE (s.deleted_at IS NULL)
  GROUP BY s.id, s.name, s.user_id
  WITH NO DATA;
CREATE INDEX idx_listing_attr_values_json ON public.listing_attribute_values USING gin (value_json);
CREATE INDEX idx_listings_location ON public.listings USING btree (individual_latitude, individual_longitude) WHERE ((individual_latitude IS NOT NULL) AND (individual_longitude IS NOT NULL) AND (has_individual_location = true));
CREATE INDEX idx_orders_payment_session ON public.orders USING btree (payment_session_id) WHERE (payment_session_id IS NOT NULL);
CREATE INDEX idx_orders_provider_session ON public.orders USING btree (payment_provider, payment_session_id) WHERE ((payment_provider IS NOT NULL) AND (payment_session_id IS NOT NULL));
CREATE INDEX idx_storefront_staff_invitation ON public.storefront_staff USING btree (invitation_id) WHERE (invitation_id IS NOT NULL);
CREATE INDEX idx_storefronts_location ON public.storefronts USING gist (public.ll_to_earth((latitude)::double precision, (longitude)::double precision)) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX idx_variant_attr_values_json ON public.variant_attribute_values USING gin (value_json) WHERE (value_json IS NOT NULL);
CREATE INDEX c2c_chats_buyer_id_idx ON public.c2c_chats USING btree (buyer_id);
CREATE INDEX c2c_chats_buyer_id_seller_id_last_message_at_idx ON public.c2c_chats USING btree (buyer_id, seller_id, last_message_at DESC);
CREATE INDEX c2c_chats_is_archived_idx ON public.c2c_chats USING btree (is_archived) WHERE (NOT is_archived);
CREATE INDEX c2c_chats_last_message_at_idx ON public.c2c_chats USING btree (last_message_at DESC) WHERE (NOT is_archived);
CREATE UNIQUE INDEX c2c_chats_least_greatest_idx1 ON public.c2c_chats USING btree (LEAST(buyer_id, seller_id), GREATEST(buyer_id, seller_id)) WHERE ((listing_id IS NULL) AND (storefront_product_id IS NULL));
CREATE INDEX c2c_chats_listing_id_idx ON public.c2c_chats USING btree (listing_id) WHERE (listing_id IS NOT NULL);
CREATE INDEX c2c_chats_seller_id_idx ON public.c2c_chats USING btree (seller_id);
CREATE INDEX c2c_chats_storefront_product_id_idx ON public.c2c_chats USING btree (storefront_product_id);
CREATE INDEX c2c_chats_updated_at_idx ON public.c2c_chats USING btree (updated_at);
CREATE INDEX c2c_messages_chat_id_created_at_idx ON public.c2c_messages USING btree (chat_id, created_at DESC);
CREATE INDEX c2c_messages_chat_id_id_idx ON public.c2c_messages USING btree (chat_id, id DESC);
CREATE INDEX c2c_messages_chat_id_idx ON public.c2c_messages USING btree (chat_id);
CREATE INDEX c2c_messages_chat_id_receiver_id_idx ON public.c2c_messages USING btree (chat_id, receiver_id) WHERE (NOT is_read);
CREATE INDEX c2c_messages_created_at_idx ON public.c2c_messages USING btree (created_at);
CREATE INDEX c2c_messages_listing_id_idx ON public.c2c_messages USING btree (listing_id);
CREATE INDEX c2c_messages_receiver_id_chat_id_idx ON public.c2c_messages USING btree (receiver_id, chat_id) WHERE (NOT is_read);
CREATE INDEX c2c_messages_receiver_id_idx ON public.c2c_messages USING btree (receiver_id);
CREATE INDEX c2c_messages_receiver_id_is_read_idx ON public.c2c_messages USING btree (receiver_id, is_read) WHERE (NOT is_read);
CREATE INDEX c2c_messages_sender_id_idx ON public.c2c_messages USING btree (sender_id);
CREATE INDEX c2c_messages_storefront_product_id_idx ON public.c2c_messages USING btree (storefront_product_id) WHERE (storefront_product_id IS NOT NULL);
CREATE UNIQUE INDEX idx_ai_mapping_name ON public.category_ai_mapping USING btree (ai_category_name);
CREATE INDEX idx_ai_mapping_priority ON public.category_ai_mapping USING btree (priority DESC);
CREATE INDEX idx_ai_mapping_target ON public.category_ai_mapping USING btree (target_category_id);
CREATE UNIQUE INDEX idx_analytics_storefront_stats_storefront_id ON public.analytics_storefront_stats USING btree (storefront_id);
CREATE INDEX idx_attachments_created_at ON public.chat_attachments USING btree (created_at);
CREATE INDEX idx_attachments_file_size ON public.chat_attachments USING btree (file_size DESC) WHERE (file_size > 1048576);
CREATE INDEX idx_attachments_file_type ON public.chat_attachments USING btree (file_type);
CREATE INDEX idx_attachments_message_id ON public.chat_attachments USING btree (message_id);
CREATE INDEX idx_attachments_metadata ON public.chat_attachments USING gin (metadata);
CREATE INDEX idx_attachments_storage ON public.chat_attachments USING btree (storage_type, storage_bucket);
CREATE INDEX idx_attr_search_cache_filterable ON public.attribute_search_cache USING gin (attributes_filterable);
CREATE INDEX idx_attr_search_cache_flat ON public.attribute_search_cache USING gin (attributes_flat);
CREATE INDEX idx_attr_search_cache_listing ON public.attribute_search_cache USING btree (listing_id);
CREATE INDEX idx_attr_search_cache_updated ON public.attribute_search_cache USING btree (last_updated);
CREATE INDEX idx_attribute_options_active ON public.attribute_options USING btree (is_active, attribute_id, sort_order) WHERE (is_active = true);
CREATE INDEX idx_attribute_options_attribute ON public.attribute_options USING btree (attribute_id);
CREATE INDEX idx_attribute_values_active ON public.attribute_values USING btree (is_active) WHERE (is_active = true);
