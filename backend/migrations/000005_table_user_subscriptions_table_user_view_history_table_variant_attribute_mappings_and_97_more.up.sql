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
CREATE TABLE public.user_view_history (
    id integer NOT NULL,
    user_id integer,
    listing_id integer,
    category_id integer,
    listing_type character varying(50) DEFAULT 'marketplace'::character varying,
    session_id character varying(100),
    viewed_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    view_duration_seconds integer,
    interaction_type character varying(50) DEFAULT 'view'::character varying,
    device_type character varying(50),
    referrer character varying(255),
    ip_address inet,
    user_agent text,
    viewport_width integer,
    viewport_height integer,
    page_depth integer,
    is_return_visit boolean DEFAULT false,
    source character varying(100),
    medium character varying(100),
    campaign character varying(200),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
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
    CONSTRAINT viber_messages_direction_check CHECK (((direction)::text = ANY (ARRAY[('incoming'::character varying)::text, ('outgoing'::character varying)::text])))
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
CREATE TABLE public.view_statistics (
    id integer NOT NULL,
    listing_id integer,
    category_id integer,
    date date NOT NULL,
    views_count integer DEFAULT 0,
    unique_users_count integer DEFAULT 0,
    unique_sessions_count integer DEFAULT 0,
    avg_view_duration numeric(10,2),
    mobile_views integer DEFAULT 0,
    desktop_views integer DEFAULT 0,
    tablet_views integer DEFAULT 0,
    contact_clicks integer DEFAULT 0,
    favorite_adds integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.vin_accident_history (
    id bigint NOT NULL,
    vin character varying(17) NOT NULL,
    accident_date date,
    accident_type character varying(100),
    severity character varying(50),
    damage_areas text[],
    airbag_deployed boolean,
    structural_damage boolean,
    repair_cost numeric(10,2),
    repair_date date,
    data_source character varying(100),
    report_number character varying(100),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.vin_check_history (
    id bigint NOT NULL,
    user_id bigint,
    vin character varying(17) NOT NULL,
    listing_id bigint,
    decode_success boolean DEFAULT true,
    decode_cache_id bigint,
    check_type character varying(50) DEFAULT 'manual'::character varying,
    ip_address inet,
    user_agent text,
    checked_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.vin_decode_cache (
    id bigint NOT NULL,
    vin character varying(17) NOT NULL,
    make character varying(100),
    model character varying(100),
    year integer,
    engine_type character varying(100),
    engine_displacement character varying(50),
    transmission_type character varying(100),
    drivetrain character varying(50),
    body_type character varying(100),
    fuel_type character varying(50),
    doors integer,
    seats integer,
    color_exterior character varying(100),
    color_interior character varying(100),
    manufacturer character varying(200),
    country_of_origin character varying(100),
    assembly_plant character varying(200),
    vehicle_class character varying(100),
    vehicle_type character varying(100),
    gross_vehicle_weight character varying(50),
    decode_status character varying(50) DEFAULT 'success'::character varying,
    error_message text,
    raw_response jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT vin_length CHECK ((length((vin)::text) = 17))
);
CREATE TABLE public.vin_ownership_history (
    id bigint NOT NULL,
    vin character varying(17) NOT NULL,
    owner_number integer,
    ownership_type character varying(50),
    purchase_date date,
    sale_date date,
    state character varying(100),
    city character varying(100),
    mileage_at_purchase integer,
    mileage_at_sale integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.vin_recalls (
    id bigint NOT NULL,
    vin character varying(17) NOT NULL,
    recall_id character varying(100),
    campaign_number character varying(100),
    component character varying(500),
    summary text,
    consequence text,
    remedy text,
    report_date date,
    remedy_date date,
    status character varying(50),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE ONLY public.address_change_log ALTER COLUMN id SET DEFAULT nextval('public.address_change_log_id_seq'::regclass);
ALTER TABLE ONLY public.ai_category_decisions ALTER COLUMN id SET DEFAULT nextval('public.ai_category_decisions_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_group_items ALTER COLUMN id SET DEFAULT nextval('public.attribute_group_items_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_groups ALTER COLUMN id SET DEFAULT nextval('public.attribute_groups_id_seq'::regclass);
ALTER TABLE ONLY public.attribute_option_translations ALTER COLUMN id SET DEFAULT nextval('public.attribute_option_translations_id_seq'::regclass);
ALTER TABLE ONLY public.balance_transactions ALTER COLUMN id SET DEFAULT nextval('public.balance_transactions_id_seq'::regclass);
ALTER TABLE ONLY public.bex_configuration ALTER COLUMN id SET DEFAULT nextval('public.bex_configuration_id_seq'::regclass);
ALTER TABLE ONLY public.bex_shipments ALTER COLUMN id SET DEFAULT nextval('public.bex_shipments_id_seq'::regclass);
ALTER TABLE ONLY public.bex_tracking_events ALTER COLUMN id SET DEFAULT nextval('public.bex_tracking_events_id_seq'::regclass);
ALTER TABLE ONLY public.car_generations ALTER COLUMN id SET DEFAULT nextval('public.car_generations_id_seq'::regclass);
ALTER TABLE ONLY public.car_makes ALTER COLUMN id SET DEFAULT nextval('public.car_makes_id_seq'::regclass);
ALTER TABLE ONLY public.car_market_analysis ALTER COLUMN id SET DEFAULT nextval('public.car_market_analysis_id_seq'::regclass);
ALTER TABLE ONLY public.car_models ALTER COLUMN id SET DEFAULT nextval('public.car_models_id_seq'::regclass);
ALTER TABLE ONLY public.category_ai_mappings ALTER COLUMN id SET DEFAULT nextval('public.category_ai_mappings_id_seq'::regclass);
ALTER TABLE ONLY public.category_attribute_groups ALTER COLUMN id SET DEFAULT nextval('public.category_attribute_groups_id_seq'::regclass);
ALTER TABLE ONLY public.category_detection_cache ALTER COLUMN id SET DEFAULT nextval('public.category_detection_cache_id_seq'::regclass);
ALTER TABLE ONLY public.category_detection_feedback ALTER COLUMN id SET DEFAULT nextval('public.category_detection_feedback_id_seq'::regclass);
ALTER TABLE ONLY public.category_keyword_weights ALTER COLUMN id SET DEFAULT nextval('public.category_keyword_weights_id_seq'::regclass);
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
ALTER TABLE ONLY public.delivery_category_defaults ALTER COLUMN id SET DEFAULT nextval('public.delivery_category_defaults_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_notifications ALTER COLUMN id SET DEFAULT nextval('public.delivery_notifications_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_pricing_rules ALTER COLUMN id SET DEFAULT nextval('public.delivery_pricing_rules_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_providers ALTER COLUMN id SET DEFAULT nextval('public.delivery_providers_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_shipments ALTER COLUMN id SET DEFAULT nextval('public.delivery_shipments_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_tracking_events ALTER COLUMN id SET DEFAULT nextval('public.delivery_tracking_events_id_seq'::regclass);
ALTER TABLE ONLY public.delivery_zones ALTER COLUMN id SET DEFAULT nextval('public.delivery_zones_id_seq'::regclass);
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
ALTER TABLE ONLY public.notification_templates ALTER COLUMN id SET DEFAULT nextval('public.notification_templates_id_seq'::regclass);
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
ALTER TABLE ONLY public.reviews ALTER COLUMN id SET DEFAULT nextval('public.reviews_id_seq'::regclass);
ALTER TABLE ONLY public.role_audit_log ALTER COLUMN id SET DEFAULT nextval('public.role_audit_log_id_seq'::regclass);
ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);
ALTER TABLE ONLY public.saved_search_notifications ALTER COLUMN id SET DEFAULT nextval('public.saved_search_notifications_id_seq'::regclass);
ALTER TABLE ONLY public.saved_searches ALTER COLUMN id SET DEFAULT nextval('public.saved_searches_id_seq'::regclass);
ALTER TABLE ONLY public.search_behavior_metrics ALTER COLUMN id SET DEFAULT nextval('public.search_behavior_metrics_id_seq'::regclass);
ALTER TABLE ONLY public.search_config ALTER COLUMN id SET DEFAULT nextval('public.search_config_id_seq'::regclass);
ALTER TABLE ONLY public.search_optimization_sessions ALTER COLUMN id SET DEFAULT nextval('public.search_optimization_sessions_id_seq'::regclass);
ALTER TABLE ONLY public.search_queries ALTER COLUMN id SET DEFAULT nextval('public.search_queries_id_seq'::regclass);
ALTER TABLE ONLY public.search_statistics ALTER COLUMN id SET DEFAULT nextval('public.search_statistics_id_seq'::regclass);
ALTER TABLE ONLY public.search_synonyms ALTER COLUMN id SET DEFAULT nextval('public.search_synonyms_id_seq'::regclass);
ALTER TABLE ONLY public.search_synonyms_config ALTER COLUMN id SET DEFAULT nextval('public.search_synonyms_config_id_seq'::regclass);
