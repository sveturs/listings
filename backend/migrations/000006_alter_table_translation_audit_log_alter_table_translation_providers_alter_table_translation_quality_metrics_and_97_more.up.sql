ALTER TABLE ONLY public.translation_audit_log ALTER COLUMN id SET DEFAULT nextval('public.translation_audit_log_id_seq'::regclass);
ALTER TABLE ONLY public.translation_providers ALTER COLUMN id SET DEFAULT nextval('public.translation_providers_id_seq'::regclass);
ALTER TABLE ONLY public.translation_quality_metrics ALTER COLUMN id SET DEFAULT nextval('public.translation_quality_metrics_id_seq'::regclass);
ALTER TABLE ONLY public.translation_sync_conflicts ALTER COLUMN id SET DEFAULT nextval('public.translation_sync_conflicts_id_seq'::regclass);
ALTER TABLE ONLY public.translation_tasks ALTER COLUMN id SET DEFAULT nextval('public.translation_tasks_id_seq'::regclass);
ALTER TABLE ONLY public.translations ALTER COLUMN id SET DEFAULT nextval('public.translations_id_seq'::regclass);
ALTER TABLE ONLY public.transliteration_rules ALTER COLUMN id SET DEFAULT nextval('public.transliteration_rules_id_seq'::regclass);
ALTER TABLE ONLY public.unified_attribute_stats ALTER COLUMN id SET DEFAULT nextval('public.unified_attribute_stats_id_seq'::regclass);
ALTER TABLE ONLY public.unified_attribute_values ALTER COLUMN id SET DEFAULT nextval('public.unified_attribute_values_id_seq'::regclass);
ALTER TABLE ONLY public.unified_attributes ALTER COLUMN id SET DEFAULT nextval('public.unified_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.unified_category_attributes ALTER COLUMN id SET DEFAULT nextval('public.unified_category_attributes_id_seq'::regclass);
ALTER TABLE ONLY public.unified_geo ALTER COLUMN id SET DEFAULT nextval('public.unified_geo_id_seq'::regclass);
ALTER TABLE ONLY public.user_behavior_events ALTER COLUMN id SET DEFAULT nextval('public.user_behavior_events_id_seq'::regclass);
ALTER TABLE ONLY public.user_car_view_history ALTER COLUMN id SET DEFAULT nextval('public.user_car_view_history_id_seq'::regclass);
ALTER TABLE ONLY public.user_contacts ALTER COLUMN id SET DEFAULT nextval('public.user_contacts_id_seq'::regclass);
ALTER TABLE ONLY public.user_notification_contacts ALTER COLUMN id SET DEFAULT nextval('public.user_notification_contacts_id_seq'::regclass);
ALTER TABLE ONLY public.user_notification_preferences ALTER COLUMN id SET DEFAULT nextval('public.user_notification_preferences_id_seq'::regclass);
ALTER TABLE ONLY public.user_subscriptions ALTER COLUMN id SET DEFAULT nextval('public.user_subscriptions_id_seq'::regclass);
ALTER TABLE ONLY public.user_view_history ALTER COLUMN id SET DEFAULT nextval('public.user_view_history_id_seq'::regclass);
ALTER TABLE ONLY public.variant_attribute_mappings ALTER COLUMN id SET DEFAULT nextval('public.variant_attribute_mappings_id_seq'::regclass);
ALTER TABLE ONLY public.viber_messages ALTER COLUMN id SET DEFAULT nextval('public.viber_messages_id_seq'::regclass);
ALTER TABLE ONLY public.viber_sessions ALTER COLUMN id SET DEFAULT nextval('public.viber_sessions_id_seq'::regclass);
ALTER TABLE ONLY public.viber_tracking_sessions ALTER COLUMN id SET DEFAULT nextval('public.viber_tracking_sessions_id_seq'::regclass);
ALTER TABLE ONLY public.viber_users ALTER COLUMN id SET DEFAULT nextval('public.viber_users_id_seq'::regclass);
ALTER TABLE ONLY public.view_statistics ALTER COLUMN id SET DEFAULT nextval('public.view_statistics_id_seq'::regclass);
ALTER TABLE ONLY public.vin_accident_history ALTER COLUMN id SET DEFAULT nextval('public.vin_accident_history_id_seq'::regclass);
ALTER TABLE ONLY public.vin_check_history ALTER COLUMN id SET DEFAULT nextval('public.vin_check_history_id_seq'::regclass);
ALTER TABLE ONLY public.vin_decode_cache ALTER COLUMN id SET DEFAULT nextval('public.vin_decode_cache_id_seq'::regclass);
ALTER TABLE ONLY public.vin_ownership_history ALTER COLUMN id SET DEFAULT nextval('public.vin_ownership_history_id_seq'::regclass);
ALTER TABLE ONLY public.vin_recalls ALTER COLUMN id SET DEFAULT nextval('public.vin_recalls_id_seq'::regclass);
CREATE VIEW public.category_attributes AS
 SELECT unified_attributes.id,
    unified_attributes.code AS name,
    unified_attributes.display_name,
    unified_attributes.attribute_type,
    COALESCE(unified_attributes.icon, ''::character varying) AS icon,
    unified_attributes.options,
    unified_attributes.validation_rules,
    unified_attributes.is_searchable,
    unified_attributes.is_filterable,
    unified_attributes.is_required,
    unified_attributes.is_required AS is_mandatory,
    unified_attributes.is_active,
    unified_attributes.sort_order,
    unified_attributes.created_at,
    unified_attributes.updated_at,
    ''::text AS custom_component,
    unified_attributes.is_variant_compatible,
    unified_attributes.affects_stock,
    unified_attributes.show_in_card
   FROM public.unified_attributes
  WHERE (unified_attributes.is_active = true);
CREATE VIEW public.category_detection_accuracy AS
 SELECT date(category_detection_feedback.created_at) AS date,
    category_detection_feedback.algorithm_version,
    count(*) AS total_detections,
    sum(
        CASE
            WHEN category_detection_feedback.user_confirmed THEN 1
            ELSE 0
        END) AS confirmed,
    round(((100.0 * (sum(
        CASE
            WHEN category_detection_feedback.user_confirmed THEN 1
            ELSE 0
        END))::numeric) / (NULLIF(count(*), 0))::numeric), 2) AS accuracy_percent,
    avg(category_detection_feedback.confidence_score) AS avg_confidence,
    percentile_cont((0.5)::double precision) WITHIN GROUP (ORDER BY ((category_detection_feedback.processing_time_ms)::double precision)) AS median_time_ms,
    percentile_cont((0.95)::double precision) WITHIN GROUP (ORDER BY ((category_detection_feedback.processing_time_ms)::double precision)) AS p95_time_ms,
    percentile_cont((0.99)::double precision) WITHIN GROUP (ORDER BY ((category_detection_feedback.processing_time_ms)::double precision)) AS p99_time_ms
   FROM public.category_detection_feedback
  WHERE (category_detection_feedback.created_at > (now() - '30 days'::interval))
  GROUP BY (date(category_detection_feedback.created_at)), category_detection_feedback.algorithm_version
  ORDER BY (date(category_detection_feedback.created_at)) DESC, category_detection_feedback.algorithm_version;
CREATE MATERIALIZED VIEW public.category_listing_counts AS
 SELECT c.id AS category_id,
    count(l.id) AS listing_count
   FROM (public.c2c_categories c
     LEFT JOIN public.c2c_listings l ON (((l.category_id = c.id) AND ((l.status)::text = 'active'::text))))
  GROUP BY c.id
  WITH NO DATA;
CREATE VIEW public.unified_listings AS
 SELECT l.id,
    'c2c'::text AS source_type,
    l.user_id,
    l.category_id,
    l.title,
    l.description,
    l.price,
    l.condition,
    l.status,
    l.location,
    l.latitude,
    l.longitude,
    l.address_city,
    l.address_country,
    l.views_count,
    l.show_on_map,
    l.original_language,
    l.created_at,
    l.updated_at,
    NULL::integer AS storefront_id,
    l.external_id,
    l.metadata,
    l.needs_reindex,
    l.address_multilingual,
    COALESCE(( SELECT jsonb_agg(jsonb_build_object('id', img.id, 'file_path', img.file_path, 'file_name', img.file_name, 'public_url', img.public_url, 'is_main', img.is_main, 'storage_type', img.storage_type) ORDER BY img.is_main DESC, img.id) AS jsonb_agg
           FROM public.c2c_images img
          WHERE (img.listing_id = l.id)), '[]'::jsonb) AS images
   FROM public.c2c_listings l
UNION ALL
 SELECT p.id,
    'b2c'::text AS source_type,
    s.user_id,
    p.category_id,
    p.name AS title,
    p.description,
    p.price,
    'new'::character varying AS condition,
        CASE
            WHEN p.is_active THEN 'active'::character varying
            ELSE 'inactive'::character varying
        END AS status,
    COALESCE(p.individual_address, s.address) AS location,
    COALESCE(p.individual_latitude, s.latitude) AS latitude,
    COALESCE(p.individual_longitude, s.longitude) AS longitude,
    s.city AS address_city,
    s.country AS address_country,
    p.view_count AS views_count,
    COALESCE(p.show_on_map, true) AS show_on_map,
    'sr'::character varying AS original_language,
    (p.created_at)::timestamp without time zone AS created_at,
    (p.updated_at)::timestamp without time zone AS updated_at,
    p.storefront_id,
    p.sku AS external_id,
    jsonb_build_object('source', 'storefront', 'storefront_id', p.storefront_id, 'stock_quantity', p.stock_quantity, 'stock_status', p.stock_status, 'currency', p.currency, 'barcode', p.barcode, 'attributes', p.attributes, 'has_variants', p.has_variants) AS metadata,
    false AS needs_reindex,
    NULL::jsonb AS address_multilingual,
    COALESCE(( SELECT jsonb_agg(jsonb_build_object('id', img.id, 'image_url', img.image_url, 'thumbnail_url', img.thumbnail_url, 'is_default', img.is_default, 'display_order', img.display_order) ORDER BY img.is_default DESC, img.display_order, img.id) AS jsonb_agg
           FROM public.b2c_product_images img
          WHERE (img.storefront_product_id = p.id)), '[]'::jsonb) AS images
   FROM (public.b2c_products p
     JOIN public.b2c_stores s ON ((s.id = p.storefront_id)));
CREATE MATERIALIZED VIEW public.user_rating_summary AS
 WITH review_stats AS (
         SELECT r.entity_id AS user_id,
            count(*) AS total_reviews,
            avg(r.rating) AS average_rating,
            count(*) FILTER (WHERE (r.rating = 1)) AS rating_1,
            count(*) FILTER (WHERE (r.rating = 2)) AS rating_2,
            count(*) FILTER (WHERE (r.rating = 3)) AS rating_3,
            count(*) FILTER (WHERE (r.rating = 4)) AS rating_4,
            count(*) FILTER (WHERE (r.rating = 5)) AS rating_5
           FROM public.reviews r
          WHERE (((r.entity_type)::text = 'user'::text) AND ((r.status)::text = 'published'::text))
          GROUP BY r.entity_id
        )
 SELECT review_stats.user_id,
    review_stats.total_reviews,
    review_stats.average_rating,
    review_stats.rating_1,
    review_stats.rating_2,
    review_stats.rating_3,
    review_stats.rating_4,
    review_stats.rating_5
   FROM review_stats
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.user_ratings AS
 SELECT users.user_id,
    count(DISTINCT r.id) AS total_reviews,
    COALESCE(avg(r.rating), (0)::numeric) AS average_rating,
    count(DISTINCT
        CASE
            WHEN ((r.entity_type)::text = 'user'::text) THEN r.id
            ELSE NULL::integer
        END) AS direct_reviews,
    count(DISTINCT
        CASE
            WHEN ((r.entity_type)::text = 'marketplace_listing'::text) THEN r.id
            ELSE NULL::integer
        END) AS listing_reviews,
    count(DISTINCT
        CASE
            WHEN ((r.entity_type)::text = 'storefront'::text) THEN r.id
            ELSE NULL::integer
        END) AS storefront_reviews,
    count(DISTINCT
        CASE
            WHEN r.is_verified_purchase THEN r.id
            ELSE NULL::integer
        END) AS verified_reviews,
    count(
        CASE
            WHEN (r.rating = 1) THEN 1
            ELSE NULL::integer
        END) AS rating_1,
    count(
        CASE
            WHEN (r.rating = 2) THEN 1
            ELSE NULL::integer
        END) AS rating_2,
    count(
        CASE
            WHEN (r.rating = 3) THEN 1
            ELSE NULL::integer
        END) AS rating_3,
    count(
        CASE
            WHEN (r.rating = 4) THEN 1
            ELSE NULL::integer
        END) AS rating_4,
    count(
        CASE
            WHEN (r.rating = 5) THEN 1
            ELSE NULL::integer
        END) AS rating_5,
    avg(
        CASE
            WHEN (r.created_at > (now() - '30 days'::interval)) THEN r.rating
            ELSE NULL::integer
        END) AS recent_rating,
    count(
        CASE
            WHEN (r.created_at > (now() - '30 days'::interval)) THEN 1
            ELSE NULL::integer
        END) AS recent_reviews,
    max(r.created_at) AS last_review_at
   FROM (( SELECT DISTINCT COALESCE(reviews.entity_origin_id, reviews.entity_id) AS user_id
           FROM public.reviews
          WHERE ((((reviews.entity_type)::text = 'user'::text) OR ((reviews.entity_origin_type)::text = 'user'::text)) AND ((reviews.status)::text = 'published'::text))) users
     LEFT JOIN public.reviews r ON ((((((r.entity_type)::text = 'user'::text) AND (r.entity_id = users.user_id)) OR (((r.entity_origin_type)::text = 'user'::text) AND (r.entity_origin_id = users.user_id))) AND ((r.status)::text = 'published'::text))))
  GROUP BY users.user_id
  WITH NO DATA;
CREATE INDEX idx_c2c_listings_location ON public.c2c_listings USING btree (latitude, longitude) WHERE (((status)::text = 'active'::text) AND (latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX idx_couriers_location ON public.couriers USING btree (current_latitude, current_longitude) WHERE (is_online = true);
CREATE INDEX idx_geocoding_cache_location ON public.geocoding_cache USING gist (location);
CREATE INDEX idx_map_items_cache_location ON public.map_items_cache USING btree (latitude, longitude);
CREATE INDEX idx_search_weights_history_reason ON public.search_weights_history USING btree (change_reason);
CREATE INDEX idx_search_weights_version ON public.search_weights USING btree (version);
CREATE INDEX idx_unified_geo_location ON public.unified_geo USING gist (location);
CREATE INDEX idx_viber_messages_user_session ON public.viber_messages USING btree (viber_user_id, session_id, created_at DESC);
CREATE INDEX b2c_delivery_options_is_active_idx ON public.b2c_delivery_options USING btree (is_active);
CREATE INDEX b2c_delivery_options_storefront_id_idx ON public.b2c_delivery_options USING btree (storefront_id);
CREATE INDEX b2c_favorites_created_at_idx ON public.b2c_favorites USING btree (created_at DESC);
CREATE INDEX b2c_favorites_product_id_idx ON public.b2c_favorites USING btree (product_id);
CREATE INDEX b2c_favorites_user_id_idx ON public.b2c_favorites USING btree (user_id);
CREATE INDEX b2c_inventory_movements_created_at_idx ON public.b2c_inventory_movements USING btree (created_at);
CREATE INDEX b2c_inventory_movements_storefront_product_id_idx ON public.b2c_inventory_movements USING btree (storefront_product_id);
CREATE INDEX b2c_inventory_movements_type_idx ON public.b2c_inventory_movements USING btree (type);
CREATE INDEX b2c_inventory_movements_variant_id_idx ON public.b2c_inventory_movements USING btree (variant_id) WHERE (variant_id IS NOT NULL);
CREATE INDEX b2c_order_items_order_id_idx ON public.b2c_order_items USING btree (order_id);
CREATE INDEX b2c_order_items_product_id_idx ON public.b2c_order_items USING btree (product_id);
CREATE INDEX b2c_order_items_variant_id_idx ON public.b2c_order_items USING btree (variant_id);
CREATE INDEX b2c_orders_customer_id_created_at_idx ON public.b2c_orders USING btree (customer_id, created_at DESC);
CREATE INDEX b2c_orders_escrow_release_date_idx ON public.b2c_orders USING btree (escrow_release_date) WHERE (escrow_release_date IS NOT NULL);
CREATE INDEX b2c_orders_status_idx ON public.b2c_orders USING btree (status);
CREATE INDEX b2c_orders_storefront_id_created_at_idx ON public.b2c_orders USING btree (storefront_id, created_at DESC);
CREATE INDEX b2c_payment_methods_is_enabled_idx ON public.b2c_payment_methods USING btree (is_enabled);
CREATE INDEX b2c_payment_methods_storefront_id_idx ON public.b2c_payment_methods USING btree (storefront_id);
CREATE INDEX b2c_product_attributes_attribute_id_idx ON public.b2c_product_attributes USING btree (attribute_id);
CREATE INDEX b2c_product_attributes_is_enabled_idx ON public.b2c_product_attributes USING btree (is_enabled);
CREATE INDEX b2c_product_attributes_product_id_idx ON public.b2c_product_attributes USING btree (product_id);
CREATE INDEX b2c_product_images_display_order_idx ON public.b2c_product_images USING btree (display_order);
CREATE INDEX b2c_product_images_storefront_product_id_idx ON public.b2c_product_images USING btree (storefront_product_id);
CREATE INDEX b2c_product_variant_images_is_main_idx ON public.b2c_product_variant_images USING btree (is_main);
CREATE INDEX b2c_product_variant_images_variant_id_idx ON public.b2c_product_variant_images USING btree (variant_id);
CREATE INDEX b2c_product_variants_barcode_idx ON public.b2c_product_variants USING btree (barcode) WHERE (barcode IS NOT NULL);
CREATE INDEX b2c_product_variants_is_active_idx ON public.b2c_product_variants USING btree (is_active);
CREATE INDEX b2c_product_variants_is_default_idx ON public.b2c_product_variants USING btree (is_default);
CREATE UNIQUE INDEX b2c_product_variants_product_id_idx ON public.b2c_product_variants USING btree (product_id) WHERE (is_default = true);
CREATE INDEX b2c_product_variants_product_id_idx1 ON public.b2c_product_variants USING btree (product_id);
CREATE INDEX b2c_product_variants_sku_idx ON public.b2c_product_variants USING btree (sku) WHERE (sku IS NOT NULL);
CREATE INDEX b2c_product_variants_variant_attributes_idx ON public.b2c_product_variants USING gin (variant_attributes);
CREATE INDEX b2c_products_barcode_idx ON public.b2c_products USING btree (barcode) WHERE (barcode IS NOT NULL);
CREATE INDEX b2c_products_category_id_idx ON public.b2c_products USING btree (category_id);
CREATE INDEX b2c_products_has_individual_location_idx ON public.b2c_products USING btree (has_individual_location);
CREATE INDEX b2c_products_has_variants_idx ON public.b2c_products USING btree (has_variants);
CREATE INDEX b2c_products_is_active_idx ON public.b2c_products USING btree (is_active);
CREATE INDEX b2c_products_location_privacy_idx ON public.b2c_products USING btree (location_privacy);
CREATE INDEX b2c_products_show_on_map_idx ON public.b2c_products USING btree (show_on_map);
CREATE INDEX b2c_products_sku_idx ON public.b2c_products USING btree (sku) WHERE (sku IS NOT NULL);
CREATE INDEX b2c_products_stock_status_idx ON public.b2c_products USING btree (stock_status);
CREATE UNIQUE INDEX b2c_products_storefront_id_barcode_idx ON public.b2c_products USING btree (storefront_id, barcode) WHERE (barcode IS NOT NULL);
CREATE INDEX b2c_products_storefront_id_idx ON public.b2c_products USING btree (storefront_id);
CREATE UNIQUE INDEX b2c_products_storefront_id_sku_idx ON public.b2c_products USING btree (storefront_id, sku) WHERE (sku IS NOT NULL);
CREATE INDEX b2c_products_storefront_id_view_count_idx ON public.b2c_products USING btree (storefront_id, view_count);
CREATE INDEX b2c_store_hours_storefront_id_idx ON public.b2c_store_hours USING btree (storefront_id);
CREATE INDEX b2c_store_staff_storefront_id_idx ON public.b2c_store_staff USING btree (storefront_id);
CREATE INDEX b2c_store_staff_user_id_idx ON public.b2c_store_staff USING btree (user_id);
CREATE INDEX b2c_stores_city_idx ON public.b2c_stores USING btree (city);
CREATE INDEX b2c_stores_geo_strategy_idx ON public.b2c_stores USING btree (geo_strategy);
CREATE INDEX b2c_stores_is_active_idx ON public.b2c_stores USING btree (is_active);
CREATE INDEX b2c_stores_latitude_longitude_idx ON public.b2c_stores USING btree (latitude, longitude) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX b2c_stores_rating_idx ON public.b2c_stores USING btree (rating DESC);
CREATE INDEX b2c_stores_slug_idx ON public.b2c_stores USING btree (slug);
CREATE INDEX b2c_stores_user_id_idx ON public.b2c_stores USING btree (user_id);
CREATE INDEX c2c_categories_external_id_idx ON public.c2c_categories USING btree (external_id);
