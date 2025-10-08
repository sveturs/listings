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
    CONSTRAINT tracking_websocket_connections_client_type_check CHECK (((client_type)::text = ANY (ARRAY[('customer'::character varying)::text, ('courier'::character varying)::text, ('merchant'::character varying)::text, ('admin'::character varying)::text])))
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
CREATE TABLE public.user_car_view_history (
    id integer NOT NULL,
    user_id integer,
    listing_id integer NOT NULL,
    session_id character varying(100),
    viewed_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    view_duration_seconds integer,
    referrer character varying(255),
    device_type character varying(50),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
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
CREATE TABLE public.user_notification_contacts (
    id integer NOT NULL,
    user_id integer,
    channel character varying(20) NOT NULL,
    contact_value character varying(255) NOT NULL,
    is_verified boolean DEFAULT false,
    is_primary boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.user_notification_preferences (
    id integer NOT NULL,
    user_id integer,
    channel character varying(20) NOT NULL,
    is_enabled boolean DEFAULT true,
    notify_on_confirmed boolean DEFAULT true,
    notify_on_picked_up boolean DEFAULT true,
    notify_on_in_transit boolean DEFAULT false,
    notify_on_out_for_delivery boolean DEFAULT true,
    notify_on_delivered boolean DEFAULT true,
    notify_on_failed boolean DEFAULT true,
    notify_on_returned boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
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
ALTER TABLE ONLY public.category_proposals ALTER COLUMN id SET DEFAULT nextval('public.category_proposals_id_seq'::regclass);
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
ALTER TABLE ONLY public.import_errors ALTER COLUMN id SET DEFAULT nextval('public.import_errors_id_seq'::regclass);
ALTER TABLE ONLY public.import_history ALTER COLUMN id SET DEFAULT nextval('public.import_history_id_seq'::regclass);
ALTER TABLE ONLY public.import_jobs ALTER COLUMN id SET DEFAULT nextval('public.import_jobs_id_seq'::regclass);
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
