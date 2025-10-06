CREATE TABLE public.custom_ui_templates (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    display_name character varying(255) NOT NULL,
    description text,
    template_code text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    variables jsonb DEFAULT '{}'::jsonb,
    is_shared boolean DEFAULT false,
    created_by integer,
    updated_by integer
);
CREATE TABLE public.deliveries (
    id integer NOT NULL,
    order_id integer,
    courier_id integer,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    pickup_address text NOT NULL,
    pickup_latitude numeric(10,8),
    pickup_longitude numeric(11,8),
    pickup_contact_name character varying(255),
    pickup_contact_phone character varying(50),
    delivery_address text NOT NULL,
    delivery_latitude numeric(10,8),
    delivery_longitude numeric(11,8),
    delivery_contact_name character varying(255),
    delivery_contact_phone character varying(50),
    assigned_at timestamp with time zone,
    accepted_at timestamp with time zone,
    picked_up_at timestamp with time zone,
    delivered_at timestamp with time zone,
    cancelled_at timestamp with time zone,
    estimated_pickup_time timestamp with time zone,
    estimated_delivery_time timestamp with time zone,
    actual_delivery_time timestamp with time zone,
    tracking_token character varying(100) DEFAULT (gen_random_uuid())::text NOT NULL,
    tracking_url text,
    share_location_enabled boolean DEFAULT true,
    distance_meters integer,
    duration_seconds integer,
    route_polyline text,
    courier_fee numeric(10,2),
    courier_tip numeric(10,2) DEFAULT 0,
    notes text,
    package_size character varying(20),
    package_weight_kg numeric(6,2),
    requires_signature boolean DEFAULT false,
    photo_proof_url text,
    customer_rating integer,
    customer_feedback text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT deliveries_customer_rating_check CHECK (((customer_rating >= 1) AND (customer_rating <= 5))),
    CONSTRAINT deliveries_package_size_check CHECK (((package_size)::text = ANY (ARRAY[('small'::character varying)::text, ('medium'::character varying)::text, ('large'::character varying)::text, ('xl'::character varying)::text]))),
    CONSTRAINT delivery_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('assigned'::character varying)::text, ('accepted'::character varying)::text, ('picked_up'::character varying)::text, ('in_transit'::character varying)::text, ('delivered'::character varying)::text, ('cancelled'::character varying)::text, ('failed'::character varying)::text])))
);
CREATE TABLE public.delivery_category_defaults (
    id integer NOT NULL,
    category_id integer,
    default_weight_kg numeric(10,3),
    default_length_cm numeric(10,2),
    default_width_cm numeric(10,2),
    default_height_cm numeric(10,2),
    default_packaging_type character varying(50),
    is_typically_fragile boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_notifications (
    id integer NOT NULL,
    shipment_id integer,
    user_id integer,
    channel character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    template character varying(50),
    data jsonb,
    sent_at timestamp with time zone,
    error_message text,
    retry_count integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_pricing_rules (
    id integer NOT NULL,
    provider_id integer,
    rule_type character varying(50) NOT NULL,
    weight_ranges jsonb,
    volume_ranges jsonb,
    zone_multipliers jsonb,
    fragile_surcharge numeric(10,2) DEFAULT 0,
    oversized_surcharge numeric(10,2) DEFAULT 0,
    special_handling_surcharge numeric(10,2) DEFAULT 0,
    min_price numeric(10,2),
    max_price numeric(10,2),
    custom_formula text,
    priority integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_providers (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255) NOT NULL,
    logo_url character varying(500),
    is_active boolean DEFAULT false,
    supports_cod boolean DEFAULT false,
    supports_insurance boolean DEFAULT false,
    supports_tracking boolean DEFAULT true,
    api_config jsonb,
    capabilities jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_shipments (
    id integer NOT NULL,
    provider_id integer,
    order_id integer,
    external_id character varying(255),
    tracking_number character varying(255),
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    sender_info jsonb NOT NULL,
    recipient_info jsonb NOT NULL,
    package_info jsonb NOT NULL,
    delivery_cost numeric(10,2),
    insurance_cost numeric(10,2),
    cod_amount numeric(10,2),
    cost_breakdown jsonb,
    pickup_date date,
    estimated_delivery date,
    actual_delivery_date timestamp with time zone,
    provider_response jsonb,
    labels jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_tracking_events (
    id integer NOT NULL,
    shipment_id integer,
    provider_id integer,
    event_time timestamp with time zone NOT NULL,
    status character varying(100) NOT NULL,
    location character varying(500),
    description text,
    raw_data jsonb,
    created_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_zones (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    type character varying(50) NOT NULL,
    countries text[],
    regions text[],
    cities text[],
    postal_codes text[],
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    radius_km numeric(10,2),
    created_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.escrow_payments (
    id bigint NOT NULL,
    payment_transaction_id bigint,
    seller_id integer,
    buyer_id integer,
    listing_id integer,
    amount numeric(12,2) NOT NULL,
    marketplace_commission numeric(12,2) NOT NULL,
    seller_amount numeric(12,2) NOT NULL,
    status character varying(50) DEFAULT 'held'::character varying,
    release_date timestamp with time zone,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT escrow_payments_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT escrow_payments_amounts_sum CHECK (((marketplace_commission + seller_amount) = amount)),
    CONSTRAINT escrow_payments_status_valid CHECK (((status)::text = ANY (ARRAY[('held'::character varying)::text, ('released'::character varying)::text, ('refunded'::character varying)::text])))
);
CREATE TABLE public.geocoding_cache (
    id bigint NOT NULL,
    input_address text NOT NULL,
    normalized_address text NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    address_components jsonb NOT NULL,
    formatted_address text NOT NULL,
    confidence numeric(3,2) NOT NULL,
    provider character varying(50) DEFAULT 'mapbox'::character varying NOT NULL,
    language character varying(5) DEFAULT 'en'::character varying,
    country_code character varying(2),
    cache_hits bigint DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '30 days'::interval) NOT NULL
);
CREATE TABLE public.gis_filter_analytics (
    id integer NOT NULL,
    user_id integer,
    session_id character varying(255) NOT NULL,
    filter_type character varying(50) NOT NULL,
    filter_params jsonb NOT NULL,
    result_count integer NOT NULL,
    response_time_ms integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.gis_isochrone_cache (
    id integer NOT NULL,
    center_point public.geography(Point,4326) NOT NULL,
    transport_mode character varying(20) NOT NULL,
    max_minutes integer NOT NULL,
    polygon public.geography(Polygon,4326) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL
);
CREATE TABLE public.listings_geo (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    geohash character varying(12) NOT NULL,
    is_precise boolean DEFAULT true NOT NULL,
    blur_radius numeric(10,2) DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    address_components jsonb,
    geocoding_confidence numeric(3,2) DEFAULT 0.0,
    address_verified boolean DEFAULT false,
    input_method character varying(50) DEFAULT 'manual'::character varying,
    location_privacy character varying(20) DEFAULT 'exact'::character varying,
    blurred_location public.geography(Point,4326),
    formatted_address text,
    district_id uuid,
    municipality_id uuid,
    CONSTRAINT chk_geocoding_confidence CHECK (((geocoding_confidence >= 0.0) AND (geocoding_confidence <= 1.0))),
    CONSTRAINT chk_input_method CHECK (((input_method)::text = ANY (ARRAY[('manual'::character varying)::text, ('geocoded'::character varying)::text, ('map_click'::character varying)::text, ('current_location'::character varying)::text]))),
    CONSTRAINT chk_location_privacy CHECK (((location_privacy)::text = ANY (ARRAY[('exact'::character varying)::text, ('street'::character varying)::text, ('district'::character varying)::text, ('city'::character varying)::text])))
);
CREATE TABLE public.gis_poi_cache (
    id integer NOT NULL,
    external_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    poi_type character varying(50) NOT NULL,
    location public.geography(Point,4326) NOT NULL,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone NOT NULL
);
CREATE TABLE public.import_history (
    id integer NOT NULL,
    source_id integer NOT NULL,
    status character varying(20) NOT NULL,
    items_total integer DEFAULT 0,
    items_imported integer DEFAULT 0,
    items_failed integer DEFAULT 0,
    log text,
    started_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    finished_at timestamp without time zone
);
CREATE TABLE public.import_sources (
    id integer NOT NULL,
    storefront_id integer NOT NULL,
    type character varying(20) NOT NULL,
    url character varying(512),
    auth_data jsonb,
    schedule character varying(50),
    mapping jsonb,
    last_import_at timestamp without time zone,
    last_import_status character varying(20),
    last_import_log text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.imported_categories (
    id integer NOT NULL,
    source_id integer NOT NULL,
    source_category character varying(255) NOT NULL,
    category_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.inventory_reservations (
    id bigint NOT NULL,
    product_id bigint NOT NULL,
    variant_id bigint,
    order_id bigint NOT NULL,
    quantity integer NOT NULL,
    status public.reservation_status DEFAULT 'active'::public.reservation_status NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT inventory_reservations_quantity_check CHECK ((quantity > 0))
);
CREATE TABLE public.item_performance_metrics (
    id bigint NOT NULL,
    item_id character varying(50) NOT NULL,
    item_type character varying(20) NOT NULL,
    impressions integer DEFAULT 0,
    clicks integer DEFAULT 0,
    ctr double precision DEFAULT 0,
    conversions integer DEFAULT 0,
    avg_position double precision DEFAULT 0,
    period_start timestamp with time zone NOT NULL,
    period_end timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT item_performance_metrics_item_type_check CHECK (((item_type)::text = ANY (ARRAY[('marketplace'::character varying)::text, ('storefront'::character varying)::text])))
);
CREATE TABLE public.listing_attribute_values (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    attribute_id integer NOT NULL,
    text_value text,
    numeric_value numeric(15,2),
    boolean_value boolean,
    date_value date,
    json_value jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    unit character varying(50),
    value_type character varying(50) NOT NULL
);
CREATE TABLE public.listing_views (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    user_id integer,
    ip_hash character varying(255),
    view_time timestamp without time zone DEFAULT now(),
    CONSTRAINT at_least_one_identifier CHECK (((user_id IS NOT NULL) OR (ip_hash IS NOT NULL)))
);
CREATE TABLE public.marketplace_chats (
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
CREATE TABLE public.marketplace_images (
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
CREATE TABLE public.marketplace_listing_variants (
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
CREATE TABLE public.marketplace_messages (
    id integer NOT NULL,
    chat_id integer,
    listing_id integer,
    sender_id integer,
    receiver_id integer,
    content text NOT NULL,
    is_read boolean DEFAULT false,
    original_language character varying(10) DEFAULT 'en'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_attachments boolean DEFAULT false,
    attachments_count integer DEFAULT 0,
    storefront_product_id integer,
    translations jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT check_message_target CHECK ((((listing_id IS NOT NULL) AND (storefront_product_id IS NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NOT NULL)) OR ((listing_id IS NULL) AND (storefront_product_id IS NULL))))
);
CREATE TABLE public.marketplace_orders (
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
    delivery_shipment_id integer,
    CONSTRAINT marketplace_orders_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('paid'::character varying)::text, ('shipped'::character varying)::text, ('delivered'::character varying)::text, ('completed'::character varying)::text, ('disputed'::character varying)::text, ('cancelled'::character varying)::text, ('refunded'::character varying)::text])))
);
CREATE TABLE public.merchant_payouts (
    id bigint NOT NULL,
    seller_id integer,
    gateway_id integer,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying,
    gateway_payout_id character varying(255),
    gateway_reference_id character varying(255),
    status character varying(50) DEFAULT 'pending'::character varying,
    bank_account_info jsonb,
    gateway_response jsonb,
    error_details jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    processed_at timestamp with time zone,
    CONSTRAINT merchant_payouts_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT merchant_payouts_status_valid CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('processing'::character varying)::text, ('completed'::character varying)::text, ('failed'::character varying)::text])))
);
CREATE TABLE public.unified_category_attributes (
    id integer NOT NULL,
    category_id integer NOT NULL,
    attribute_id integer NOT NULL,
    is_enabled boolean DEFAULT true,
    is_required boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    category_specific_options jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.unified_attribute_values (
    id integer NOT NULL,
    entity_type character varying(50) NOT NULL,
    entity_id integer NOT NULL,
    attribute_id integer NOT NULL,
    text_value text,
    numeric_value numeric,
    boolean_value boolean,
    date_value date,
    json_value jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unified_attribute_values_entity_type_check CHECK (((entity_type)::text = ANY (ARRAY[('listing'::character varying)::text, ('product'::character varying)::text, ('product_variant'::character varying)::text])))
)
WITH (autovacuum_vacuum_scale_factor='0.1', autovacuum_analyze_scale_factor='0.05', autovacuum_vacuum_cost_delay='10');
CREATE TABLE public.notification_templates (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    channel character varying(20) NOT NULL,
    name character varying(255) NOT NULL,
    subject text,
    body_template text NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.notifications (
    id integer NOT NULL,
    user_id integer NOT NULL,
    type character varying(50) NOT NULL,
    title text NOT NULL,
    message text NOT NULL,
    data jsonb,
    is_read boolean DEFAULT false,
    delivered_to jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.payment_gateways (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    is_active boolean DEFAULT true,
    config jsonb NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.payment_methods (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    code character varying(50) NOT NULL,
    type character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    minimum_amount numeric(12,2),
    maximum_amount numeric(12,2),
    fee_percentage numeric(5,2),
    fixed_fee numeric(12,2),
    credentials jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.payment_transactions (
    id bigint NOT NULL,
    gateway_id integer,
    user_id integer,
    listing_id integer,
    order_reference character varying(255) NOT NULL,
    gateway_transaction_id character varying(255),
    gateway_reference_id character varying(255),
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying,
    marketplace_commission numeric(12,2),
    seller_amount numeric(12,2),
    status character varying(50) DEFAULT 'pending'::character varying,
    gateway_status character varying(50),
    payment_method character varying(50),
    customer_email character varying(255),
    description text,
    gateway_response jsonb,
    error_details jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    authorized_at timestamp with time zone,
    captured_at timestamp with time zone,
    failed_at timestamp with time zone,
    source_type character varying(20) DEFAULT 'marketplace_listing'::character varying,
    source_id bigint,
    storefront_id integer,
    capture_mode character varying(20) DEFAULT 'manual'::character varying,
    auto_capture_at timestamp with time zone,
    capture_deadline_at timestamp with time zone,
    capture_attempted_at timestamp with time zone,
    capture_attempts integer DEFAULT 0,
    CONSTRAINT payment_transactions_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT payment_transactions_capture_mode_check CHECK (((capture_mode)::text = ANY (ARRAY[('auto'::character varying)::text, ('manual'::character varying)::text]))),
    CONSTRAINT payment_transactions_status_valid CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('authorized'::character varying)::text, ('captured'::character varying)::text, ('failed'::character varying)::text, ('refunded'::character varying)::text, ('voided'::character varying)::text])))
);
CREATE TABLE public.permissions (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    resource character varying(50) NOT NULL,
    action character varying(50) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.post_express_locations (
    id integer NOT NULL,
    post_express_id integer,
    name character varying(255) NOT NULL,
    name_cyrillic character varying(255),
    postal_code character varying(20),
    municipality character varying(255),
    latitude double precision,
    longitude double precision,
    region character varying(255),
    district character varying(255),
    delivery_zone character varying(50),
    is_active boolean DEFAULT true,
    supports_cod boolean DEFAULT true,
    supports_express boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.post_express_offices (
    id integer NOT NULL,
    office_code character varying(50) NOT NULL,
    location_id integer,
    name character varying(255) NOT NULL,
    address character varying(500) NOT NULL,
    phone character varying(50),
    email character varying(255),
    working_hours jsonb,
    latitude double precision,
    longitude double precision,
    accepts_packages boolean DEFAULT true,
    issues_packages boolean DEFAULT true,
    has_atm boolean DEFAULT false,
    has_parking boolean DEFAULT false,
    wheelchair_accessible boolean DEFAULT false,
    is_active boolean DEFAULT true,
    temporary_closed boolean DEFAULT false,
    closed_until timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.post_express_rates (
    id integer NOT NULL,
    weight_from numeric(10,3) NOT NULL,
    weight_to numeric(10,3) NOT NULL,
    base_price numeric(10,2) NOT NULL,
    insurance_included_up_to numeric(10,2) DEFAULT 15000,
    insurance_rate_percent numeric(5,2) DEFAULT 1.0,
    cod_fee numeric(10,2) DEFAULT 45,
    max_length_cm integer DEFAULT 60,
    max_width_cm integer DEFAULT 60,
    max_height_cm integer DEFAULT 60,
    max_dimensions_sum_cm integer DEFAULT 180,
    delivery_days_min integer DEFAULT 1,
    delivery_days_max integer DEFAULT 3,
    is_active boolean DEFAULT true,
    is_special_offer boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.post_express_settings (
    id integer NOT NULL,
    api_username character varying(255) NOT NULL,
    api_password character varying(255) NOT NULL,
    api_endpoint character varying(500) DEFAULT 'https://wsp.postexpress.rs/api'::character varying NOT NULL,
    sender_name character varying(255) NOT NULL,
    sender_address character varying(500) NOT NULL,
    sender_city character varying(255) NOT NULL,
    sender_postal_code character varying(20) NOT NULL,
    sender_phone character varying(50) NOT NULL,
    sender_email character varying(255),
    enabled boolean DEFAULT true,
    test_mode boolean DEFAULT true,
    auto_print_labels boolean DEFAULT false,
    auto_track_shipments boolean DEFAULT false,
    notify_on_pickup boolean DEFAULT false,
    notify_on_delivery boolean DEFAULT true,
    notify_on_failed_delivery boolean DEFAULT true,
    total_shipments integer DEFAULT 0,
    successful_deliveries integer DEFAULT 0,
    failed_deliveries integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.post_express_shipments (
    id integer NOT NULL,
    marketplace_order_id integer,
    storefront_order_id bigint,
    tracking_number character varying(255),
    barcode character varying(255),
    post_express_id character varying(255),
    sender_name character varying(255) NOT NULL,
    sender_address character varying(500) NOT NULL,
    sender_city character varying(255) NOT NULL,
    sender_postal_code character varying(20) NOT NULL,
    sender_phone character varying(50) NOT NULL,
    sender_email character varying(255),
    recipient_name character varying(255) NOT NULL,
    recipient_address character varying(500) NOT NULL,
    recipient_city character varying(255) NOT NULL,
    recipient_postal_code character varying(20) NOT NULL,
    recipient_phone character varying(50) NOT NULL,
    recipient_email character varying(255),
    weight numeric(10,3) NOT NULL,
    length_cm integer,
    width_cm integer,
    height_cm integer,
    package_contents text,
    declared_value numeric(10,2),
    service_type character varying(50) DEFAULT 'standard'::character varying,
    cod_amount numeric(10,2),
    insurance_amount numeric(10,2),
    express_delivery boolean DEFAULT false,
    office_pickup boolean DEFAULT false,
    office_code character varying(50),
    status character varying(50) DEFAULT 'created'::character varying,
    status_description text,
    last_tracking_update timestamp with time zone,
    pickup_date timestamp with time zone,
    delivery_date timestamp with time zone,
    label_url text,
    label_printed_at timestamp with time zone,
    receipt_url text,
    status_history jsonb DEFAULT '[]'::jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    sender_location_id integer,
    recipient_location_id integer,
    cod_reference character varying(255),
    base_price numeric(10,2),
    insurance_fee numeric(10,2),
    cod_fee numeric(10,2),
    total_price numeric(10,2),
    delivery_status character varying(100),
    delivery_instructions text,
    notes text,
    invoice_url text,
    invoice_number character varying(255),
    invoice_date timestamp with time zone,
    pod_url text,
    registered_at timestamp with time zone,
    picked_up_at timestamp with time zone,
    delivered_at timestamp with time zone,
    failed_at timestamp with time zone,
    returned_at timestamp with time zone,
    internal_notes text,
    failed_reason text
);
CREATE TABLE public.post_express_tracking_events (
    id integer NOT NULL,
    shipment_id integer,
    event_code character varying(50),
    event_description text,
    event_location character varying(255),
    event_timestamp timestamp with time zone,
    additional_info jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.price_history (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    price numeric(12,2) NOT NULL,
    effective_from timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    effective_to timestamp without time zone,
    change_source character varying(50) NOT NULL,
    change_percentage numeric(10,2),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.product_variant_attribute_values (
    id integer NOT NULL,
    attribute_id integer NOT NULL,
    value character varying(255) NOT NULL,
    display_name character varying(255) NOT NULL,
    color_hex character varying(7),
    image_url text,
    sort_order integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.product_variant_attributes (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(255) NOT NULL,
    type character varying(50) DEFAULT 'text'::character varying NOT NULL,
    is_required boolean DEFAULT false NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    affects_stock boolean DEFAULT false
);
CREATE TABLE public.query_cache (
    id integer NOT NULL,
    query_hash character varying(64) NOT NULL,
    query_text text NOT NULL,
    result_data jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    expires_at timestamp without time zone NOT NULL,
    hit_count integer DEFAULT 0
);
CREATE TABLE public.review_confirmations (
    id integer NOT NULL,
    review_id integer NOT NULL,
    confirmed_by integer NOT NULL,
    confirmation_status character varying(50) NOT NULL,
    confirmed_at timestamp without time zone DEFAULT now() NOT NULL,
    notes text,
    CONSTRAINT review_confirmations_confirmation_status_check CHECK (((confirmation_status)::text = ANY (ARRAY[('confirmed'::character varying)::text, ('disputed'::character varying)::text])))
);
CREATE TABLE public.review_disputes (
    id integer NOT NULL,
    review_id integer NOT NULL,
    disputed_by integer NOT NULL,
    dispute_reason character varying(100) NOT NULL,
    dispute_description text NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    admin_id integer,
    admin_notes text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    resolved_at timestamp without time zone,
    CONSTRAINT review_disputes_dispute_reason_check CHECK (((dispute_reason)::text = ANY (ARRAY[('not_a_customer'::character varying)::text, ('false_information'::character varying)::text, ('deal_cancelled'::character varying)::text, ('spam'::character varying)::text, ('other'::character varying)::text]))),
    CONSTRAINT review_disputes_status_check CHECK (((status)::text = ANY (ARRAY[('pending'::character varying)::text, ('in_review'::character varying)::text, ('resolved_keep_review'::character varying)::text, ('resolved_remove_review'::character varying)::text, ('resolved_remove_verification'::character varying)::text, ('cancelled'::character varying)::text])))
);
CREATE TABLE public.review_responses (
    id integer NOT NULL,
    review_id integer NOT NULL,
    user_id integer NOT NULL,
    response text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.review_votes (
    review_id integer NOT NULL,
    user_id integer NOT NULL,
    vote_type character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT review_votes_vote_type_check CHECK (((vote_type)::text = ANY (ARRAY[('helpful'::character varying)::text, ('not_helpful'::character varying)::text])))
);
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
CREATE TABLE public.role_audit_log (
    id integer NOT NULL,
    user_id integer,
    target_user_id integer,
    action character varying(50) NOT NULL,
    old_role_id integer,
    new_role_id integer,
    details jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.role_permissions (
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.roles (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    display_name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_system boolean DEFAULT false,
    is_assignable boolean DEFAULT true,
    priority integer DEFAULT 100
);
CREATE TABLE public.saved_search_notifications (
    id integer NOT NULL,
    saved_search_id integer NOT NULL,
    new_listings_count integer DEFAULT 0 NOT NULL,
    notification_sent boolean DEFAULT false,
    sent_at timestamp with time zone,
    error_message text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.saved_searches (
    id integer NOT NULL,
    user_id integer NOT NULL,
    name character varying(100) NOT NULL,
    filters jsonb DEFAULT '{}'::jsonb NOT NULL,
    search_type character varying(50) DEFAULT 'cars'::character varying,
    notify_enabled boolean DEFAULT false,
    notify_frequency character varying(20) DEFAULT 'daily'::character varying,
    last_notified_at timestamp with time zone,
    results_count integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.search_behavior_metrics (
    id bigint NOT NULL,
    search_query text NOT NULL,
    total_searches integer DEFAULT 0,
    total_clicks integer DEFAULT 0,
    ctr double precision DEFAULT 0,
    avg_click_position double precision DEFAULT 0,
    conversions integer DEFAULT 0,
    conversion_rate double precision DEFAULT 0,
    period_start timestamp with time zone NOT NULL,
    period_end timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.search_config (
    id bigint NOT NULL,
    min_search_length integer DEFAULT 2 NOT NULL,
    max_suggestions integer DEFAULT 10 NOT NULL,
    fuzzy_enabled boolean DEFAULT true NOT NULL,
    fuzzy_max_edits integer DEFAULT 2 NOT NULL,
    synonyms_enabled boolean DEFAULT true NOT NULL,
    transliteration_enabled boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.search_optimization_sessions (
    id bigint NOT NULL,
    status character varying(20) DEFAULT 'running'::character varying NOT NULL,
    start_time timestamp with time zone DEFAULT now() NOT NULL,
    end_time timestamp with time zone,
    total_fields integer DEFAULT 0 NOT NULL,
    processed_fields integer DEFAULT 0 NOT NULL,
    results jsonb,
    error_message text,
    created_by integer NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT search_optimization_sessions_status_check CHECK (((status)::text = ANY (ARRAY[('running'::character varying)::text, ('completed'::character varying)::text, ('failed'::character varying)::text, ('cancelled'::character varying)::text])))
);
CREATE TABLE public.search_queries (
    id integer NOT NULL,
    query text NOT NULL,
    normalized_query text NOT NULL,
    search_count integer DEFAULT 1 NOT NULL,
    last_searched timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    language character varying(10) DEFAULT 'ru'::character varying NOT NULL,
    results_count integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.search_statistics (
    id bigint NOT NULL,
    query text NOT NULL,
    results_count integer DEFAULT 0 NOT NULL,
    search_duration_ms bigint NOT NULL,
    user_id bigint,
    search_filters jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.search_synonyms (
    id integer NOT NULL,
    term character varying(255) NOT NULL,
    synonym character varying(255) NOT NULL,
    language character varying(10) DEFAULT 'ru'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.search_synonyms_config (
    id bigint NOT NULL,
    term character varying(255) NOT NULL,
    synonyms text[] NOT NULL,
    language character varying(10) DEFAULT 'ru'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.search_weights (
    id bigint NOT NULL,
    field_name character varying(100) NOT NULL,
    weight double precision NOT NULL,
    search_type character varying(20) DEFAULT 'fulltext'::character varying NOT NULL,
    item_type character varying(20) DEFAULT 'global'::character varying NOT NULL,
    category_id integer,
    description text,
    is_active boolean DEFAULT true,
    version integer DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    created_by integer,
    updated_by integer,
    CONSTRAINT search_weights_item_type_check CHECK (((item_type)::text = ANY (ARRAY[('marketplace'::character varying)::text, ('storefront'::character varying)::text, ('global'::character varying)::text]))),
    CONSTRAINT search_weights_search_type_check CHECK (((search_type)::text = ANY (ARRAY[('fulltext'::character varying)::text, ('fuzzy'::character varying)::text, ('exact'::character varying)::text]))),
    CONSTRAINT search_weights_weight_check CHECK (((weight >= (0.0)::double precision) AND (weight <= (1.0)::double precision)))
);
CREATE TABLE public.search_weights_history (
    id bigint NOT NULL,
    weight_id bigint NOT NULL,
    old_weight double precision NOT NULL,
    new_weight double precision NOT NULL,
    change_reason character varying(50) DEFAULT 'manual'::character varying NOT NULL,
    change_metadata jsonb DEFAULT '{}'::jsonb,
    changed_by integer,
    changed_at timestamp with time zone DEFAULT now(),
    CONSTRAINT search_weights_history_change_reason_check CHECK (((change_reason)::text = ANY (ARRAY[('manual'::character varying)::text, ('optimization'::character varying)::text, ('rollback'::character varying)::text, ('initialization'::character varying)::text])))
);
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
