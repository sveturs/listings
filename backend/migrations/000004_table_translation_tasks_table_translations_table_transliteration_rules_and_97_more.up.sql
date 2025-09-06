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
ALTER TABLE ONLY public.address_change_log ALTER COLUMN id SET DEFAULT nextval('public.address_change_log_id_seq'::regclass);
ALTER TABLE ONLY public.admin_users ALTER COLUMN id SET DEFAULT nextval('public.admin_users_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_group_items ALTER COLUMN id SET DEFAULT nextval('public.attribute_group_items_id_seq'::regclass);
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
ALTER TABLE ONLY public.custom_ui_component_usage ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_component_usage_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_components ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_components_id_seq'::regclass);
ALTER TABLE ONLY public.custom_ui_templates ALTER COLUMN id SET DEFAULT nextval('public.custom_ui_templates_id_seq'::regclass);
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
ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);
ALTER TABLE ONLY public.review_confirmations ALTER COLUMN id SET DEFAULT nextval('public.review_confirmations_id_seq'::regclass);
ALTER TABLE ONLY public.review_disputes ALTER COLUMN id SET DEFAULT nextval('public.review_disputes_id_seq'::regclass);
ALTER TABLE ONLY public.review_responses ALTER COLUMN id SET DEFAULT nextval('public.review_responses_id_seq'::regclass);
ALTER TABLE ONLY public.reviews ALTER COLUMN id SET DEFAULT nextval('public.reviews_id_seq'::regclass);
ALTER TABLE ONLY public.role_audit_log ALTER COLUMN id SET DEFAULT nextval('public.role_audit_log_id_seq'::regclass);
ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);
ALTER TABLE ONLY public.search_behavior_metrics ALTER COLUMN id SET DEFAULT nextval('public.search_behavior_metrics_id_seq'::regclass);
ALTER TABLE ONLY public.search_config ALTER COLUMN id SET DEFAULT nextval('public.search_config_id_seq'::regclass);
ALTER TABLE ONLY public.search_optimization_sessions ALTER COLUMN id SET DEFAULT nextval('public.search_optimization_sessions_id_seq'::regclass);
ALTER TABLE ONLY public.search_queries ALTER COLUMN id SET DEFAULT nextval('public.search_queries_id_seq'::regclass);
ALTER TABLE ONLY public.search_statistics ALTER COLUMN id SET DEFAULT nextval('public.search_statistics_id_seq'::regclass);
ALTER TABLE ONLY public.search_synonyms ALTER COLUMN id SET DEFAULT nextval('public.search_synonyms_id_seq'::regclass);
ALTER TABLE ONLY public.search_synonyms_config ALTER COLUMN id SET DEFAULT nextval('public.search_synonyms_config_id_seq'::regclass);
ALTER TABLE ONLY public.search_weights ALTER COLUMN id SET DEFAULT nextval('public.search_weights_id_seq'::regclass);
ALTER TABLE ONLY public.search_weights_history ALTER COLUMN id SET DEFAULT nextval('public.search_weights_history_id_seq'::regclass);
ALTER TABLE ONLY public.shopping_cart_items ALTER COLUMN id SET DEFAULT nextval('public.shopping_cart_items_id_seq'::regclass);
ALTER TABLE ONLY public.shopping_carts ALTER COLUMN id SET DEFAULT nextval('public.shopping_carts_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_delivery_options ALTER COLUMN id SET DEFAULT nextval('public.storefront_delivery_options_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_hours ALTER COLUMN id SET DEFAULT nextval('public.storefront_hours_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_inventory_movements ALTER COLUMN id SET DEFAULT nextval('public.storefront_inventory_movements_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_order_items ALTER COLUMN id SET DEFAULT nextval('public.storefront_order_items_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_orders ALTER COLUMN id SET DEFAULT nextval('public.storefront_orders_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_payment_methods ALTER COLUMN id SET DEFAULT nextval('public.storefront_payment_methods_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_product_attributes ALTER COLUMN id SET DEFAULT nextval('public.storefront_product_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_product_images ALTER COLUMN id SET DEFAULT nextval('public.storefront_product_images_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_product_variant_images ALTER COLUMN id SET DEFAULT nextval('public.storefront_product_variant_images_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_product_variants ALTER COLUMN id SET DEFAULT nextval('public.storefront_product_variants_id_seq'::regclass);
ALTER TABLE ONLY public.storefront_staff ALTER COLUMN id SET DEFAULT nextval('public.storefront_staff_id_seq'::regclass);
ALTER TABLE ONLY public.storefronts ALTER COLUMN id SET DEFAULT nextval('public.storefronts_id_seq'::regclass);
ALTER TABLE ONLY public.subscription_history ALTER COLUMN id SET DEFAULT nextval('public.subscription_history_id_seq'::regclass);
ALTER TABLE ONLY public.subscription_payments ALTER COLUMN id SET DEFAULT nextval('public.subscription_payments_id_seq'::regclass);
ALTER TABLE ONLY public.subscription_plans ALTER COLUMN id SET DEFAULT nextval('public.subscription_plans_id_seq'::regclass);
ALTER TABLE ONLY public.subscription_usage ALTER COLUMN id SET DEFAULT nextval('public.subscription_usage_id_seq'::regclass);
ALTER TABLE ONLY public.translation_audit_log ALTER COLUMN id SET DEFAULT nextval('public.translation_audit_log_id_seq'::regclass);
ALTER TABLE ONLY public.translation_providers ALTER COLUMN id SET DEFAULT nextval('public.translation_providers_id_seq'::regclass);
ALTER TABLE ONLY public.translation_quality_metrics ALTER COLUMN id SET DEFAULT nextval('public.translation_quality_metrics_id_seq'::regclass);
