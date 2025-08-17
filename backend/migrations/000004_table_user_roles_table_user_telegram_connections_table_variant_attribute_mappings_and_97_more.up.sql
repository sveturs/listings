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
ALTER TABLE ONLY public.category_attributes ALTER COLUMN id SET DEFAULT nextval('public.category_attributes_id_seq'::regclass);
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
ALTER TABLE ONLY public.storefront_fbs_settings ALTER COLUMN id SET DEFAULT nextval('public.storefront_fbs_settings_id_seq'::regclass);
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
ALTER TABLE ONLY public.translation_audit_log ALTER COLUMN id SET DEFAULT nextval('public.translation_audit_log_id_seq'::regclass);
ALTER TABLE ONLY public.translation_providers ALTER COLUMN id SET DEFAULT nextval('public.translation_providers_id_seq'::regclass);
ALTER TABLE ONLY public.translation_quality_metrics ALTER COLUMN id SET DEFAULT nextval('public.translation_quality_metrics_id_seq'::regclass);
ALTER TABLE ONLY public.translation_sync_conflicts ALTER COLUMN id SET DEFAULT nextval('public.translation_sync_conflicts_id_seq'::regclass);
ALTER TABLE ONLY public.translation_tasks ALTER COLUMN id SET DEFAULT nextval('public.translation_tasks_id_seq'::regclass);
ALTER TABLE ONLY public.translations ALTER COLUMN id SET DEFAULT nextval('public.translations_id_seq'::regclass);
ALTER TABLE ONLY public.transliteration_rules ALTER COLUMN id SET DEFAULT nextval('public.transliteration_rules_id_seq'::regclass);
ALTER TABLE ONLY public.unified_geo ALTER COLUMN id SET DEFAULT nextval('public.unified_geo_id_seq'::regclass);
