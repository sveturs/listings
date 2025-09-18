CREATE TABLE public.shopping_cart_items (
    id bigint NOT NULL,
    cart_id bigint NOT NULL,
    product_id bigint NOT NULL,
    variant_id bigint,
    quantity integer NOT NULL,
    price_per_unit numeric(10,2) NOT NULL,
    total_price numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT shopping_cart_items_price_per_unit_check CHECK ((price_per_unit >= (0)::numeric)),
    CONSTRAINT shopping_cart_items_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT shopping_cart_items_total_price_check CHECK ((total_price >= (0)::numeric))
);
CREATE TABLE public.shopping_carts (
    id bigint NOT NULL,
    user_id integer,
    storefront_id integer NOT NULL,
    session_id character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT check_cart_owner CHECK ((((user_id IS NOT NULL) AND (session_id IS NULL)) OR ((user_id IS NULL) AND (session_id IS NOT NULL))))
);
CREATE TABLE public.storefront_delivery_options (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    base_price numeric(10,2) DEFAULT 0.00 NOT NULL,
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
CREATE TABLE public.storefront_favorites (
    user_id integer NOT NULL,
    product_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
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
CREATE TABLE public.storefront_inventory_movements (
    id integer NOT NULL,
    storefront_product_id integer NOT NULL,
    variant_id integer,
    type character varying(20) NOT NULL,
    quantity integer NOT NULL,
    reason character varying(50) NOT NULL,
    order_id integer,
    notes text,
    user_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT storefront_inventory_movements_type_check CHECK (((type)::text = ANY (ARRAY[('in'::character varying)::text, ('out'::character varying)::text, ('adjustment'::character varying)::text])))
);
CREATE TABLE public.storefront_order_items (
    id bigint NOT NULL,
    order_id bigint NOT NULL,
    product_id bigint NOT NULL,
    variant_id bigint,
    product_name character varying(255) NOT NULL,
    variant_name character varying(255),
    product_sku character varying(100),
    quantity integer DEFAULT 1 NOT NULL,
    unit_price numeric(12,2) NOT NULL,
    total_price numeric(12,2) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    notes text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT storefront_order_items_quantity_check CHECK ((quantity > 0))
);
CREATE TABLE public.storefront_orders (
    id bigint NOT NULL,
    order_number character varying(32) NOT NULL,
    storefront_id integer,
    customer_id integer,
    payment_transaction_id bigint,
    subtotal_amount numeric(12,2) NOT NULL,
    shipping_amount numeric(12,2) DEFAULT 0,
    tax_amount numeric(12,2) DEFAULT 0,
    total_amount numeric(12,2) NOT NULL,
    commission_amount numeric(12,2) NOT NULL,
    seller_amount numeric(12,2) NOT NULL,
    currency character(3) DEFAULT 'RSD'::bpchar,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    escrow_release_date date,
    escrow_days integer DEFAULT 3,
    shipping_address jsonb,
    billing_address jsonb,
    shipping_method character varying(100),
    shipping_provider character varying(50),
    tracking_number character varying(100),
    customer_notes text,
    seller_notes text,
    confirmed_at timestamp without time zone,
    shipped_at timestamp without time zone,
    delivered_at timestamp without time zone,
    cancelled_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    payment_method character varying(50) DEFAULT 'allsecure'::character varying NOT NULL,
    payment_status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    notes text,
    metadata jsonb DEFAULT '{}'::jsonb,
    discount numeric(12,2) DEFAULT 0,
    pickup_address jsonb
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
CREATE TABLE public.storefront_product_attributes (
    id integer NOT NULL,
    product_id integer NOT NULL,
    attribute_id integer NOT NULL,
    is_enabled boolean DEFAULT true NOT NULL,
    is_required boolean DEFAULT false NOT NULL,
    custom_values jsonb DEFAULT '[]'::jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.storefront_product_images (
    id integer NOT NULL,
    storefront_product_id integer NOT NULL,
    image_url text NOT NULL,
    thumbnail_url text NOT NULL,
    display_order integer DEFAULT 0 NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.storefront_product_variant_images (
    id integer NOT NULL,
    variant_id integer NOT NULL,
    image_url text NOT NULL,
    thumbnail_url text,
    alt_text character varying(255),
    display_order integer DEFAULT 0 NOT NULL,
    is_main boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.storefront_product_variants (
    id integer NOT NULL,
    product_id integer NOT NULL,
    sku character varying(100),
    barcode character varying(100),
    price numeric(15,2),
    compare_at_price numeric(15,2),
    cost_price numeric(15,2),
    stock_quantity integer DEFAULT 0 NOT NULL,
    stock_status character varying(20) DEFAULT 'in_stock'::character varying NOT NULL,
    low_stock_threshold integer DEFAULT 5,
    variant_attributes jsonb DEFAULT '{}'::jsonb NOT NULL,
    weight numeric(10,3),
    dimensions jsonb,
    is_active boolean DEFAULT true NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    view_count integer DEFAULT 0 NOT NULL,
    sold_count integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT storefront_product_variants_stock_quantity_check CHECK ((stock_quantity >= 0)),
    CONSTRAINT storefront_product_variants_stock_status_check CHECK (((stock_status)::text = ANY (ARRAY[('in_stock'::character varying)::text, ('low_stock'::character varying)::text, ('out_of_stock'::character varying)::text])))
);
CREATE TABLE public.storefront_products (
    id integer DEFAULT nextval('public.global_product_id_seq'::regclass) NOT NULL,
    storefront_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text NOT NULL,
    price numeric(15,2) NOT NULL,
    currency character(3) DEFAULT 'USD'::bpchar NOT NULL,
    category_id integer NOT NULL,
    sku character varying(100),
    barcode character varying(100),
    stock_quantity integer DEFAULT 0 NOT NULL,
    stock_status character varying(20) DEFAULT 'in_stock'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    attributes jsonb DEFAULT '{}'::jsonb,
    view_count integer DEFAULT 0 NOT NULL,
    sold_count integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    has_individual_location boolean DEFAULT false,
    individual_address text,
    individual_latitude numeric(10,8),
    individual_longitude numeric(11,8),
    location_privacy public.location_privacy_level DEFAULT 'exact'::public.location_privacy_level,
    show_on_map boolean DEFAULT true,
    has_variants boolean DEFAULT false,
    CONSTRAINT storefront_products_price_check CHECK ((price >= (0)::numeric)),
    CONSTRAINT storefront_products_stock_quantity_check CHECK ((stock_quantity >= 0)),
    CONSTRAINT storefront_products_stock_status_check CHECK (((stock_status)::text = ANY (ARRAY[('in_stock'::character varying)::text, ('low_stock'::character varying)::text, ('out_of_stock'::character varying)::text])))
);
CREATE TABLE public.user_storefronts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    logo_path character varying(255),
    slug character varying(100),
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    creation_transaction_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    phone character varying(50),
    email character varying(255),
    website character varying(255),
    address character varying(255),
    city character varying(100),
    country character varying(100),
    latitude numeric(10,8),
    longitude numeric(11,8)
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
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.storefronts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    slug character varying(100) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    logo_url character varying(500),
    banner_url character varying(500),
    theme jsonb DEFAULT '{"layout": "grid", "primaryColor": "#1976d2"}'::jsonb,
    phone character varying(50),
    email character varying(255),
    website character varying(255),
    address text,
    city character varying(100),
    postal_code character varying(20),
    country character varying(2) DEFAULT 'RS'::character varying,
    latitude numeric(10,8),
    longitude numeric(11,8),
    settings jsonb DEFAULT '{}'::jsonb,
    seo_meta jsonb DEFAULT '{}'::jsonb,
    is_active boolean DEFAULT true,
    is_verified boolean DEFAULT false,
    verification_date timestamp without time zone,
    rating numeric(3,2) DEFAULT 0.00,
    reviews_count integer DEFAULT 0,
    products_count integer DEFAULT 0,
    sales_count integer DEFAULT 0,
    views_count integer DEFAULT 0,
    subscription_plan character varying(50) DEFAULT 'starter'::character varying,
    subscription_expires_at timestamp without time zone,
    commission_rate numeric(5,2) DEFAULT 3.00,
    ai_agent_enabled boolean DEFAULT false,
    ai_agent_config jsonb DEFAULT '{}'::jsonb,
    live_shopping_enabled boolean DEFAULT false,
    group_buying_enabled boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    formatted_address text,
    geo_strategy public.storefront_geo_strategy DEFAULT 'storefront_location'::public.storefront_geo_strategy,
    default_privacy_level public.location_privacy_level DEFAULT 'exact'::public.location_privacy_level,
    address_verified boolean DEFAULT false,
    subscription_id integer,
    is_subscription_active boolean DEFAULT true,
    followers_count integer DEFAULT 0
);
CREATE TABLE public.subscription_history (
    id integer NOT NULL,
    subscription_id integer NOT NULL,
    user_id integer NOT NULL,
    action character varying(50) NOT NULL,
    from_plan_id integer,
    to_plan_id integer,
    reason text,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer
);
CREATE TABLE public.subscription_payments (
    id integer NOT NULL,
    subscription_id integer NOT NULL,
    user_id integer NOT NULL,
    payment_id integer,
    amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'EUR'::character varying,
    period_start date NOT NULL,
    period_end date NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    payment_method character varying(50),
    transaction_data jsonb DEFAULT '{}'::jsonb,
    paid_at timestamp without time zone,
    failed_at timestamp without time zone,
    refunded_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.subscription_plans (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    price_monthly numeric(10,2) DEFAULT 0,
    price_yearly numeric(10,2) DEFAULT 0,
    max_storefronts integer DEFAULT 1,
    max_products_per_storefront integer DEFAULT 50,
    max_staff_per_storefront integer DEFAULT 1,
    max_images_total integer DEFAULT 100,
    has_ai_assistant boolean DEFAULT false,
    has_live_shopping boolean DEFAULT false,
    has_export_data boolean DEFAULT false,
    has_custom_domain boolean DEFAULT false,
    has_analytics boolean DEFAULT true,
    has_priority_support boolean DEFAULT false,
    commission_rate numeric(5,2) DEFAULT 10.00,
    free_trial_days integer DEFAULT 0,
    sort_order integer DEFAULT 1,
    is_active boolean DEFAULT true,
    is_recommended boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.subscription_usage (
    id integer NOT NULL,
    subscription_id integer NOT NULL,
    storefront_id integer,
    resource_type character varying(50) NOT NULL,
    resource_id integer,
    resource_count integer DEFAULT 1,
    action character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.tracking_websocket_connections (
    id integer NOT NULL,
    connection_id character varying(100) NOT NULL,
    delivery_id integer,
    client_type character varying(20),
    user_id integer,
    viber_user_id integer,
    connected_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    disconnected_at timestamp with time zone,
    last_ping_at timestamp with time zone,
    ip_address inet,
    user_agent text,
    CONSTRAINT tracking_websocket_connections_client_type_check CHECK (((client_type)::text = ANY ((ARRAY['customer'::character varying, 'courier'::character varying, 'merchant'::character varying, 'admin'::character varying])::text[])))
);
CREATE TABLE public.translation_audit_log (
    id integer NOT NULL,
    user_id integer,
    action character varying(100) NOT NULL,
    entity_type character varying(50),
    entity_id integer,
    old_value text,
    new_value text,
    ip_address inet,
    user_agent text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.translation_providers (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    provider_type character varying(50) NOT NULL,
    api_key text,
    settings jsonb DEFAULT '{}'::jsonb,
    usage_limit integer,
    usage_current integer DEFAULT 0,
    is_active boolean DEFAULT true,
    priority integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.translation_quality_metrics (
    id integer NOT NULL,
    translation_id integer,
    quality_score numeric(3,2),
    character_count integer,
    word_count integer,
    has_placeholders boolean DEFAULT false,
    has_html_tags boolean DEFAULT false,
    checked_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    checked_by character varying(50),
    issues jsonb DEFAULT '[]'::jsonb
);
CREATE TABLE public.translation_sync_conflicts (
    id integer NOT NULL,
    source_type character varying(50) NOT NULL,
    target_type character varying(50) NOT NULL,
    entity_identifier text NOT NULL,
    source_value text,
    target_value text,
    conflict_type character varying(50),
    resolved boolean DEFAULT false,
    resolved_by integer,
    resolved_at timestamp without time zone,
    resolution_type character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.translation_tasks (
    id integer NOT NULL,
    task_type character varying(50) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    source_language character varying(10),
    target_languages text[],
    entity_references jsonb DEFAULT '[]'::jsonb,
    provider_id integer,
    created_by integer,
    assigned_to integer,
    started_at timestamp without time zone,
    completed_at timestamp without time zone,
    error_message text,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.translations (
    id integer NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id integer NOT NULL,
    language character varying(10) NOT NULL,
    field_name character varying(50) NOT NULL,
    translated_text text NOT NULL,
    is_machine_translated boolean DEFAULT true,
    is_verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    metadata jsonb DEFAULT '{}'::jsonb,
    version integer DEFAULT 1,
    last_modified_by integer
);
CREATE TABLE public.transliteration_rules (
    id integer NOT NULL,
    source_char character varying(10) NOT NULL,
    target_char character varying(20) NOT NULL,
    language character varying(2) NOT NULL,
    enabled boolean DEFAULT true,
    priority integer DEFAULT 0,
    description text,
    rule_type character varying(20) DEFAULT 'custom'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.unified_attribute_stats (
    id integer NOT NULL,
    attribute_id integer NOT NULL,
    category_id integer,
    usage_count bigint DEFAULT 0,
    avg_query_time_ms numeric(8,3) DEFAULT 0,
    last_updated timestamp with time zone DEFAULT now()
);
CREATE TABLE public.unified_geo (
    id bigint NOT NULL,
    source_type public.geo_source_type NOT NULL,
    source_id bigint NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    geohash character varying(12) NOT NULL,
    formatted_address text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    privacy_level public.location_privacy_level DEFAULT 'exact'::public.location_privacy_level,
    original_location public.geography(Point,4326),
    blur_radius_meters integer DEFAULT 0
);
CREATE TABLE public.user_behavior_events (
    id bigint NOT NULL,
    event_type character varying(50) NOT NULL,
    user_id integer,
    session_id character varying(100) NOT NULL,
    search_query text,
    item_id character varying(50),
    item_type character varying(20),
    "position" integer,
    metadata jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT user_behavior_events_item_type_check CHECK (((item_type)::text = ANY (ARRAY[('marketplace'::character varying)::text, ('storefront'::character varying)::text, (NULL::character varying)::text])))
);
CREATE TABLE public.user_contacts (
    id bigint NOT NULL,
    user_id integer NOT NULL,
    contact_user_id integer NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    added_from_chat_id integer,
    notes text,
    CONSTRAINT user_contacts_check CHECK ((user_id <> contact_user_id)),
    CONSTRAINT user_contacts_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('accepted'::character varying)::text, ('blocked'::character varying)::text])))
);
CREATE TABLE public.user_subscriptions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    plan_id integer NOT NULL,
    status character varying(50) DEFAULT 'active'::character varying,
    billing_cycle character varying(20) DEFAULT 'monthly'::character varying,
    started_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    trial_ends_at timestamp without time zone,
    current_period_start timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    current_period_end timestamp without time zone NOT NULL,
    cancelled_at timestamp without time zone,
    expires_at timestamp without time zone,
    last_payment_id integer,
    last_payment_at timestamp without time zone,
    next_payment_at timestamp without time zone,
    payment_method character varying(50),
    auto_renew boolean DEFAULT true,
    used_storefronts integer DEFAULT 0,
    metadata jsonb DEFAULT '{}'::jsonb,
    notes text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.variant_attribute_mappings (
    id integer NOT NULL,
    variant_attribute_id integer NOT NULL,
    category_id integer NOT NULL,
    sort_order integer DEFAULT 0,
    is_required boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
CREATE TABLE public.viber_messages (
    id integer NOT NULL,
    viber_user_id integer,
    session_id integer,
    message_token character varying(100),
    direction character varying(20),
    message_type character varying(50),
    content text,
    rich_media jsonb,
    is_billable boolean DEFAULT false,
    status character varying(20),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT viber_messages_direction_check CHECK (((direction)::text = ANY ((ARRAY['incoming'::character varying, 'outgoing'::character varying])::text[])))
);
CREATE TABLE public.viber_sessions (
    id integer NOT NULL,
    viber_user_id integer,
    started_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_message_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp with time zone DEFAULT (CURRENT_TIMESTAMP + '24:00:00'::interval),
    message_count integer DEFAULT 0,
    context jsonb DEFAULT '{}'::jsonb,
    active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.viber_tracking_sessions (
    id integer NOT NULL,
    viber_user_id integer,
    delivery_id integer,
    tracking_token character varying(100) NOT NULL,
    started_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_viewed_at timestamp with time zone,
    page_views integer DEFAULT 1,
    is_active boolean DEFAULT true,
    device_info jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.viber_users (
    id integer NOT NULL,
    viber_id character varying(100) NOT NULL,
    user_id integer,
    name character varying(255) DEFAULT 'Unknown User'::character varying NOT NULL,
    avatar_url text,
    language character varying(10) DEFAULT 'sr'::character varying,
    country_code character varying(5),
    api_version integer DEFAULT 1,
    subscribed boolean DEFAULT true,
    subscribed_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    last_session_at timestamp with time zone,
    conversation_started_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE ONLY public.address_change_log ALTER COLUMN id SET DEFAULT nextval('public.address_change_log_id_seq'::regclass);
ALTER TABLE ONLY public.admin_users ALTER COLUMN id SET DEFAULT nextval('public.admin_users_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_group_items ALTER COLUMN id SET DEFAULT nextval('public.attribute_group_items_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_groups ALTER COLUMN id SET DEFAULT nextval('public.attribute_groups_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_option_translations ALTER COLUMN id SET DEFAULT nextval('public.attribute_option_translations_id_seq'::regclass);
ALTER TABLE ONLY public.balance_transactions ALTER COLUMN id SET DEFAULT nextval('public.balance_transactions_id_seq'::regclass);
ALTER TABLE ONLY public.car_generations ALTER COLUMN id SET DEFAULT nextval('public.car_generations_id_seq'::regclass);
ALTER TABLE ONLY public.car_makes ALTER COLUMN id SET DEFAULT nextval('public.car_makes_id_seq'::regclass);
ALTER TABLE ONLY public.car_market_analysis ALTER COLUMN id SET DEFAULT nextval('public.car_market_analysis_id_seq'::regclass);
ALTER TABLE ONLY public.car_models ALTER COLUMN id SET DEFAULT nextval('public.car_models_id_seq'::regclass);
ALTER TABLE ONLY public.category_attribute_groups ALTER COLUMN id SET DEFAULT nextval('public.category_attribute_groups_id_seq'::regclass);
ALTER TABLE ONLY public.category_keywords ALTER COLUMN id SET DEFAULT nextval('public.category_keywords_id_seq'::regclass);
ALTER TABLE ONLY public.category_variant_attributes ALTER COLUMN id SET DEFAULT nextval('public.category_variant_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.chat_attachments ALTER COLUMN id SET DEFAULT nextval('public.chat_attachments_id_seq'::regclass);
ALTER TABLE ONLY public.component_templates ALTER COLUMN id SET DEFAULT nextval('public.component_templates_id_seq'::regclass);
ALTER TABLE ONLY public.courier_location_history ALTER COLUMN id SET DEFAULT nextval('public.courier_location_history_id_seq'::regclass);
ALTER TABLE ONLY public.courier_zones ALTER COLUMN id SET DEFAULT nextval('public.courier_zones_id_seq'::regclass);
ALTER TABLE ONLY public.couriers ALTER COLUMN id SET DEFAULT nextval('public.couriers_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_component_usage ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_component_usage_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_components ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_components_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_templates ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_templates_id_seq'::regclass);
ALTER TABLE ONLY public.deliveries ALTER COLUMN id SET DEFAULT nextval('public.deliveries_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_notifications ALTER COLUMN id SET DEFAULT nextval('public.delivery_notifications_id_seq'::regclass);
ALTER TABLE ONLY public.escrow_payments ALTER COLUMN id SET DEFAULT nextval('public.escrow_payments_id_seq'::regclass);
ALTER TABLE ONLY public.geocoding_cache ALTER COLUMN id SET DEFAULT nextval('public.geocoding_cache_id_seq'::regclass);
ALTER TABLE ONLY public.gis_filter_analytics ALTER COLUMN id SET DEFAULT nextval('public.gis_filter_analytics_id_seq'::regclass);
ALTER TABLE ONLY public.gis_isochrone_cache ALTER COLUMN id SET DEFAULT nextval('public.gis_isochrone_cache_id_seq'::regclass);
ALTER TABLE ONLY public.gis_poi_cache ALTER COLUMN id SET DEFAULT nextval('public.gis_poi_cache_id_seq'::regclass);
ALTER TABLE ONLY public.import_history ALTER COLUMN id SET DEFAULT nextval('public.import_history_id_seq'::regclass);
ALTER TABLE ONLY public.import_sources ALTER COLUMN id SET DEFAULT nextval('public.import_sources_id_seq'::regclass);
ALTER TABLE ONLY public.imported_categories ALTER COLUMN id SET DEFAULT nextval('public.imported_categories_id_seq'::regclass);
ALTER TABLE ONLY public.inventory_reservations ALTER COLUMN id SET DEFAULT nextval('public.inventory_reservations_id_seq'::regclass);
ALTER TABLE ONLY public.item_performance_metrics ALTER COLUMN id SET DEFAULT nextval('public.item_performance_metrics_id_seq'::regclass);
ALTER TABLE ONLY public.listing_attribute_values ALTER COLUMN id SET DEFAULT nextval('public.listing_attribute_values_id_seq'::regclass);
ALTER TABLE ONLY public.listing_views ALTER COLUMN id SET DEFAULT nextval('public.listing_views_id_seq'::regclass);
ALTER TABLE ONLY public.listings_geo ALTER COLUMN id SET DEFAULT nextval('public.listings_geo_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_categories ALTER COLUMN id SET DEFAULT nextval('public.marketplace_categories_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_chats ALTER COLUMN id SET DEFAULT nextval('public.marketplace_chats_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_images ALTER COLUMN id SET DEFAULT nextval('public.marketplace_images_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_listing_variants ALTER COLUMN id SET DEFAULT nextval('public.marketplace_listing_variants_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_messages ALTER COLUMN id SET DEFAULT nextval('public.marketplace_messages_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_orders ALTER COLUMN id SET DEFAULT nextval('public.marketplace_orders_id_seq'::regclass);
ALTER TABLE ONLY public.merchant_payouts ALTER COLUMN id SET DEFAULT nextval('public.merchant_payouts_id_seq'::regclass);
ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);
ALTER TABLE ONLY public.payment_gateways ALTER COLUMN id SET DEFAULT nextval('public.payment_gateways_id_seq'::regclass);
ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);
ALTER TABLE ONLY public.payment_transactions ALTER COLUMN id SET DEFAULT nextval('public.payment_transactions_id_seq'::regclass);
ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_locations ALTER COLUMN id SET DEFAULT nextval('public.post_express_locations_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_offices ALTER COLUMN id SET DEFAULT nextval('public.post_express_offices_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_rates ALTER COLUMN id SET DEFAULT nextval('public.post_express_rates_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_settings ALTER COLUMN id SET DEFAULT nextval('public.post_express_settings_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_shipments ALTER COLUMN id SET DEFAULT nextval('public.post_express_shipments_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_tracking_events ALTER COLUMN id SET DEFAULT nextval('public.post_express_tracking_events_id_seq'::regclass);
ALTER TABLE ONLY public.price_history ALTER COLUMN id SET DEFAULT nextval('public.price_history_id_seq'::regclass);
ALTER TABLE ONLY public.product_variant_attribute_values ALTER COLUMN id SET DEFAULT nextval('public.product_variant_attribute_values_id_seq'::regclass);
ALTER TABLE ONLY public.product_variant_attributes ALTER COLUMN id SET DEFAULT nextval('public.product_variant_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.query_cache ALTER COLUMN id SET DEFAULT nextval('public.query_cache_id_seq'::regclass);
ALTER TABLE ONLY public.review_confirmations ALTER COLUMN id SET DEFAULT nextval('public.review_confirmations_id_seq'::regclass);
ALTER TABLE ONLY public.review_disputes ALTER COLUMN id SET DEFAULT nextval('public.review_disputes_id_seq'::regclass);
ALTER TABLE ONLY public.review_responses ALTER COLUMN id SET DEFAULT nextval('public.review_responses_id_seq'::regclass);
