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
    pickup_address jsonb,
    delivery_method character varying(50) DEFAULT 'pickup'::character varying,
    delivery_cost numeric(12,2) DEFAULT 0,
    delivery_address jsonb,
    delivery_notes text,
    delivery_provider character varying(50) DEFAULT 'standard'::character varying,
    delivery_tracking_number character varying(100),
    delivery_status character varying(50),
    delivery_metadata jsonb
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
    is_active boolean DEFAULT false,
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
    address_verified boolean DEFAULT false
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
    version integer DEFAULT 1
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
CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    google_id character varying(255),
    picture_url text,
    phone character varying(20),
    bio text,
    notification_email boolean DEFAULT true,
    timezone character varying(50) DEFAULT 'UTC'::character varying,
    last_seen timestamp without time zone,
    account_status character varying(20) DEFAULT 'active'::character varying,
    settings jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    city character varying(100),
    country character varying(100),
    password character varying(255),
    provider character varying(50) DEFAULT 'email'::character varying,
    preferred_language character varying(10) DEFAULT 'ru'::character varying,
    role_id integer,
    CONSTRAINT users_account_status_check CHECK (((account_status)::text = ANY (ARRAY[('active'::character varying)::text, ('inactive'::character varying)::text, ('suspended'::character varying)::text]))),
    CONSTRAINT users_preferred_language_check CHECK (((preferred_language)::text = ANY ((ARRAY['ru'::character varying, 'sr'::character varying, 'en'::character varying])::text[])))
);
CREATE TABLE public.user_roles (
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    assigned_by integer
);
CREATE TABLE public.user_telegram_connections (
    user_id integer NOT NULL,
    telegram_chat_id character varying(100) NOT NULL,
    telegram_username character varying(100),
    connected_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.variant_attribute_mappings (
    id integer NOT NULL,
    variant_attribute_id integer NOT NULL,
    category_attribute_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.warehouse_inventory (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    storefront_product_id bigint,
    marketplace_listing_id integer,
    sku character varying(100) NOT NULL,
    barcode character varying(100),
    external_id character varying(100),
    product_name character varying(500) NOT NULL,
    product_description text,
    quantity_total integer DEFAULT 0 NOT NULL,
    quantity_available integer DEFAULT 0 NOT NULL,
    quantity_reserved integer DEFAULT 0 NOT NULL,
    quantity_damaged integer DEFAULT 0,
    unit_weight_kg numeric(10,3),
    unit_length_cm numeric(10,2),
    unit_width_cm numeric(10,2),
    unit_height_cm numeric(10,2),
    unit_volume_m3 numeric(10,4),
    location_zone character varying(50),
    location_rack character varying(50),
    location_shelf character varying(50),
    location_bin character varying(50),
    received_at timestamp with time zone,
    expiry_date date,
    storage_fee_daily numeric(10,2) DEFAULT 0,
    is_fragile boolean DEFAULT false,
    requires_refrigeration boolean DEFAULT false,
    is_hazardous boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.warehouse_invoices (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    warehouse_id integer NOT NULL,
    period_start date NOT NULL,
    period_end date NOT NULL,
    invoice_number character varying(100) NOT NULL,
    invoice_date date NOT NULL,
    due_date date NOT NULL,
    line_items jsonb NOT NULL,
    subtotal numeric(12,2) NOT NULL,
    tax_amount numeric(12,2) DEFAULT 0,
    total_amount numeric(12,2) NOT NULL,
    status character varying(50) DEFAULT 'draft'::character varying,
    paid_at timestamp with time zone,
    payment_method character varying(50),
    payment_reference character varying(100),
    pdf_url character varying(500),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.warehouse_movements (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    inventory_id integer,
    movement_type character varying(50) NOT NULL,
    movement_reason character varying(100),
    quantity integer NOT NULL,
    quantity_before integer NOT NULL,
    quantity_after integer NOT NULL,
    order_id integer,
    storefront_order_id bigint,
    shipment_id integer,
    document_number character varying(100),
    document_type character varying(50),
    performed_by integer,
    notes text,
    created_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.warehouse_pickup_orders (
    id integer NOT NULL,
    warehouse_id integer NOT NULL,
    marketplace_order_id integer,
    storefront_order_id bigint,
    pickup_code character varying(10) NOT NULL,
    qr_code_url character varying(500),
    status character varying(50) DEFAULT 'pending'::character varying,
    ready_at timestamp with time zone,
    picked_up_at timestamp with time zone,
    expires_at timestamp with time zone,
    customer_name character varying(200) NOT NULL,
    customer_phone character varying(50) NOT NULL,
    customer_email character varying(200),
    pickup_confirmed_by character varying(200),
    id_document_type character varying(50),
    id_document_number character varying(100),
    signature_url character varying(500),
    notification_sent_at timestamp with time zone,
    reminder_sent_at timestamp with time zone,
    notes text,
    pickup_photo_url character varying(500),
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.warehouses (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(200) NOT NULL,
    type character varying(50) DEFAULT 'main'::character varying,
    address character varying(500) NOT NULL,
    city character varying(100) NOT NULL,
    postal_code character varying(20) NOT NULL,
    country character varying(2) DEFAULT 'RS'::character varying,
    phone character varying(50),
    email character varying(200),
    manager_name character varying(200),
    manager_phone character varying(50),
    latitude numeric(10,8),
    longitude numeric(11,8),
    working_hours jsonb,
    total_area_m2 numeric(10,2),
    storage_area_m2 numeric(10,2),
    max_capacity_m3 numeric(10,2),
    current_occupancy_m3 numeric(10,2) DEFAULT 0,
    supports_fbs boolean DEFAULT true,
    supports_pickup boolean DEFAULT true,
    has_refrigeration boolean DEFAULT false,
    has_loading_dock boolean DEFAULT true,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.webhook_audit_log (
    id bigint NOT NULL,
    webhook_id character varying(255),
    webhook_type character varying(100) NOT NULL,
    action character varying(50) NOT NULL,
    status character varying(50) NOT NULL,
    details jsonb,
    error_message text,
    processing_time_ms integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE ONLY public.address_change_log ALTER COLUMN id SET DEFAULT nextval('public.address_change_log_id_seq'::regclass);
ALTER TABLE ONLY public.admin_users ALTER COLUMN id SET DEFAULT nextval('public.admin_users_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_group_items ALTER COLUMN id SET DEFAULT nextval('public.attribute_group_items_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_groups ALTER COLUMN id SET DEFAULT nextval('public.attribute_groups_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_option_translations ALTER COLUMN id SET DEFAULT nextval('public.attribute_option_translations_id_seq'::regclass);
ALTER TABLE ONLY public.balance_transactions ALTER COLUMN id SET DEFAULT nextval('public.balance_transactions_id_seq'::regclass);
ALTER TABLE ONLY public.bex_municipalities ALTER COLUMN id SET DEFAULT nextval('public.bex_municipalities_id_seq'::regclass);
ALTER TABLE ONLY public.bex_parcel_shops ALTER COLUMN id SET DEFAULT nextval('public.bex_parcel_shops_id_seq'::regclass);
ALTER TABLE ONLY public.bex_places ALTER COLUMN id SET DEFAULT nextval('public.bex_places_id_seq'::regclass);
ALTER TABLE ONLY public.bex_rates ALTER COLUMN id SET DEFAULT nextval('public.bex_rates_id_seq'::regclass);
ALTER TABLE ONLY public.bex_settings ALTER COLUMN id SET DEFAULT nextval('public.bex_settings_id_seq'::regclass);
ALTER TABLE ONLY public.bex_shipments ALTER COLUMN id SET DEFAULT nextval('public.bex_shipments_id_seq'::regclass);
ALTER TABLE ONLY public.bex_streets ALTER COLUMN id SET DEFAULT nextval('public.bex_streets_id_seq'::regclass);
ALTER TABLE ONLY public.bex_tracking_events ALTER COLUMN id SET DEFAULT nextval('public.bex_tracking_events_id_seq'::regclass);
ALTER TABLE ONLY public.car_generations ALTER COLUMN id SET DEFAULT nextval('public.car_generations_id_seq'::regclass);
ALTER TABLE ONLY public.car_makes ALTER COLUMN id SET DEFAULT nextval('public.car_makes_id_seq'::regclass);
ALTER TABLE ONLY public.car_market_analysis ALTER COLUMN id SET DEFAULT nextval('public.car_market_analysis_id_seq'::regclass);
ALTER TABLE ONLY public.car_models ALTER COLUMN id SET DEFAULT nextval('public.car_models_id_seq'::regclass);
ALTER TABLE ONLY public.category_attribute_groups ALTER COLUMN id SET DEFAULT nextval('public.category_attribute_groups_id_seq'::regclass);
ALTER TABLE ONLY public.category_attributes ALTER COLUMN id SET DEFAULT nextval('public.category_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.category_keywords ALTER COLUMN id SET DEFAULT nextval('public.category_keywords_id_seq'::regclass);
ALTER TABLE ONLY public.category_variant_attributes ALTER COLUMN id SET DEFAULT nextval('public.category_variant_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.chat_attachments ALTER COLUMN id SET DEFAULT nextval('public.chat_attachments_id_seq'::regclass);
ALTER TABLE ONLY public.component_templates ALTER COLUMN id SET DEFAULT nextval('public.component_templates_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_component_usage ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_component_usage_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_components ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_components_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_templates ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_templates_id_seq'::regclass);
ALTER TABLE ONLY public.escrow_payments ALTER COLUMN id SET DEFAULT nextval('public.escrow_payments_id_seq'::regclass);
ALTER TABLE ONLY public.failed_webhooks ALTER COLUMN id SET DEFAULT nextval('public.failed_webhooks_id_seq'::regclass);
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
ALTER TABLE ONLY public.marketplace_messages ALTER COLUMN id SET DEFAULT nextval('public.marketplace_messages_id_seq'::regclass);
ALTER TABLE ONLY public.marketplace_orders ALTER COLUMN id SET DEFAULT nextval('public.marketplace_orders_id_seq'::regclass);
ALTER TABLE ONLY public.merchant_payouts ALTER COLUMN id SET DEFAULT nextval('public.merchant_payouts_id_seq'::regclass);
ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);
ALTER TABLE ONLY public.payment_gateways ALTER COLUMN id SET DEFAULT nextval('public.payment_gateways_id_seq'::regclass);
ALTER TABLE ONLY public.payment_methods ALTER COLUMN id SET DEFAULT nextval('public.payment_methods_id_seq'::regclass);
ALTER TABLE ONLY public.payment_transactions ALTER COLUMN id SET DEFAULT nextval('public.payment_transactions_id_seq'::regclass);
ALTER TABLE ONLY public.permissions ALTER COLUMN id SET DEFAULT nextval('public.permissions_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_api_logs ALTER COLUMN id SET DEFAULT nextval('public.post_express_api_logs_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_locations ALTER COLUMN id SET DEFAULT nextval('public.post_express_locations_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_offices ALTER COLUMN id SET DEFAULT nextval('public.post_express_offices_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_rates ALTER COLUMN id SET DEFAULT nextval('public.post_express_rates_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_settings ALTER COLUMN id SET DEFAULT nextval('public.post_express_settings_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_shipments ALTER COLUMN id SET DEFAULT nextval('public.post_express_shipments_id_seq'::regclass);
ALTER TABLE ONLY public.post_express_tracking_events ALTER COLUMN id SET DEFAULT nextval('public.post_express_tracking_events_id_seq'::regclass);
ALTER TABLE ONLY public.price_history ALTER COLUMN id SET DEFAULT nextval('public.price_history_id_seq'::regclass);
ALTER TABLE ONLY public.product_variant_attribute_values ALTER COLUMN id SET DEFAULT nextval('public.product_variant_attribute_values_id_seq'::regclass);
ALTER TABLE ONLY public.product_variant_attributes ALTER COLUMN id SET DEFAULT nextval('public.product_variant_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);
ALTER TABLE ONLY public.review_confirmations ALTER COLUMN id SET DEFAULT nextval('public.review_confirmations_id_seq'::regclass);
ALTER TABLE ONLY public.review_disputes ALTER COLUMN id SET DEFAULT nextval('public.review_disputes_id_seq'::regclass);
ALTER TABLE ONLY public.review_responses ALTER COLUMN id SET DEFAULT nextval('public.review_responses_id_seq'::regclass);
ALTER TABLE ONLY public.reviews ALTER COLUMN id SET DEFAULT nextval('public.reviews_id_seq'::regclass);
ALTER TABLE ONLY public.role_audit_log ALTER COLUMN id SET DEFAULT nextval('public.role_audit_log_id_seq'::regclass);
ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);
ALTER TABLE ONLY public.search_behavior_metrics ALTER COLUMN id SET DEFAULT nextval('public.search_behavior_metrics_id_seq'::regclass);
