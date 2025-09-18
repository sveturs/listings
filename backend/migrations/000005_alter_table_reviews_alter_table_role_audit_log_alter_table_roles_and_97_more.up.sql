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
ALTER TABLE ONLY public.tracking_websocket_connections ALTER COLUMN id SET DEFAULT nextval('public.tracking_websocket_connections_id_seq'::regclass);
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
ALTER TABLE ONLY public.user_contacts ALTER COLUMN id SET DEFAULT nextval('public.user_contacts_id_seq'::regclass);
ALTER TABLE ONLY public.user_storefronts ALTER COLUMN id SET DEFAULT nextval('public.user_storefronts_id_seq'::regclass);
ALTER TABLE ONLY public.user_subscriptions ALTER COLUMN id SET DEFAULT nextval('public.user_subscriptions_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
ALTER TABLE ONLY public.variant_attribute_mappings ALTER COLUMN id SET DEFAULT nextval('public.variant_attribute_mappings_id_seq'::regclass);
ALTER TABLE ONLY public.viber_messages ALTER COLUMN id SET DEFAULT nextval('public.viber_messages_id_seq'::regclass);
ALTER TABLE ONLY public.viber_sessions ALTER COLUMN id SET DEFAULT nextval('public.viber_sessions_id_seq'::regclass);
ALTER TABLE ONLY public.viber_tracking_sessions ALTER COLUMN id SET DEFAULT nextval('public.viber_tracking_sessions_id_seq'::regclass);
ALTER TABLE ONLY public.viber_users ALTER COLUMN id SET DEFAULT nextval('public.viber_users_id_seq'::regclass);
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
    unified_attributes.affects_stock
   FROM public.unified_attributes
  WHERE (unified_attributes.is_active = true);
CREATE MATERIALIZED VIEW public.category_listing_counts AS
 WITH RECURSIVE category_tree AS (
         SELECT marketplace_categories.id,
            ARRAY[marketplace_categories.id] AS category_path,
            marketplace_categories.name,
            1 AS depth,
            ( SELECT count(*) AS count
                   FROM public.marketplace_listings ml
                  WHERE ((ml.category_id = marketplace_categories.id) AND ((ml.status)::text = 'active'::text))) AS direct_count
           FROM public.marketplace_categories
          WHERE (marketplace_categories.parent_id IS NULL)
        UNION ALL
         SELECT c.id,
            (ct_1.category_path || c.id),
            c.name,
            (ct_1.depth + 1),
            ( SELECT count(*) AS count
                   FROM public.marketplace_listings ml
                  WHERE ((ml.category_id = c.id) AND ((ml.status)::text = 'active'::text))) AS direct_count
           FROM (public.marketplace_categories c
             JOIN category_tree ct_1 ON ((c.parent_id = ct_1.id)))
          WHERE (ct_1.depth < 10)
        )
 SELECT ct.id AS category_id,
    ((ct.direct_count)::numeric + COALESCE(( SELECT sum(ch.direct_count) AS sum
           FROM category_tree ch
          WHERE ((ch.category_path[1:array_length(ct.category_path, 1)] = ct.category_path) AND (ch.id <> ct.id))), (0)::numeric)) AS listing_count,
    max(ct.depth) AS category_depth
   FROM category_tree ct
  GROUP BY ct.id, ct.direct_count, ct.category_path
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.gis_listing_density_grid AS
 WITH grid_cells AS (
         SELECT x.x AS grid_x,
            y.y AS grid_y,
            (public.st_makeenvelope((x.x)::double precision, (y.y)::double precision, ((x.x + 0.005))::double precision, ((y.y + 0.005))::double precision, 4326))::public.geography AS cell
           FROM generate_series(18.0, 23.0, 0.005) x(x),
            generate_series(42.0, 46.5, 0.005) y(y)
        )
 SELECT row_number() OVER () AS id,
    g.grid_x,
    g.grid_y,
    g.cell,
    count(l.id) AS listing_count,
    (public.st_area(g.cell) / (1000000.0)::double precision) AS area_km2,
        CASE
            WHEN (public.st_area(g.cell) > (0)::double precision) THEN ((count(l.id))::double precision / (public.st_area(g.cell) / (1000000.0)::double precision))
            ELSE (0)::double precision
        END AS density
   FROM ((grid_cells g
     LEFT JOIN public.listings_geo lg ON (((lg.location IS NOT NULL) AND public.st_within((lg.location)::public.geometry, (g.cell)::public.geometry))))
     LEFT JOIN public.marketplace_listings l ON (((l.id = lg.listing_id) AND ((l.status)::text = 'active'::text))))
  GROUP BY g.grid_x, g.grid_y, g.cell
  WITH NO DATA;
CREATE VIEW public.map_items_view AS
 SELECT ml.id,
    ml.title,
    ml.description,
    ml.price,
    ml.condition,
    ml.location,
    ml.latitude,
    ml.longitude,
    ml.address_city AS city,
    ml.address_country AS country,
    ml.status,
    ml.created_at,
    ml.updated_at,
    ml.user_id,
    ml.category_id,
    mc.name AS category_name,
    mc.slug AS category_slug,
    ( SELECT mi.public_url
           FROM public.marketplace_images mi
          WHERE ((mi.listing_id = ml.id) AND (mi.is_main = true))
         LIMIT 1) AS main_image_url,
    u.name AS user_name,
    ml.show_on_map
   FROM ((public.marketplace_listings ml
     LEFT JOIN public.marketplace_categories mc ON ((ml.category_id = mc.id)))
     LEFT JOIN public.users u ON ((ml.user_id = u.id)))
  WHERE (((ml.status)::text = 'active'::text) AND (ml.show_on_map = true) AND (ml.latitude IS NOT NULL) AND (ml.longitude IS NOT NULL));
CREATE MATERIALIZED VIEW public.mv_category_statistics AS
 SELECT c.id AS category_id,
    c.name,
    c.slug,
    count(DISTINCT l.id) AS total_listings,
    count(DISTINCT l.id) FILTER (WHERE ((l.status)::text = 'active'::text)) AS active_listings,
    count(DISTINCT l.user_id) AS unique_sellers,
    avg(l.price) AS avg_price,
    min(l.price) AS min_price,
    max(l.price) AS max_price,
    count(DISTINCT ua.id) AS total_attributes,
    now() AS last_updated
   FROM (((public.marketplace_categories c
     LEFT JOIN public.marketplace_listings l ON ((c.id = l.category_id)))
     LEFT JOIN public.unified_category_attributes uca ON ((c.id = uca.category_id)))
     LEFT JOIN public.unified_attributes ua ON (((uca.attribute_id = ua.id) AND (ua.is_active = true))))
  GROUP BY c.id, c.name, c.slug
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.mv_popular_category_attributes AS
 SELECT c.id AS category_id,
    c.name AS category_name,
    ua.id AS attribute_id,
    ua.code AS attribute_code,
    ua.name AS attribute_name,
    ua.attribute_type,
    uca.is_required,
    uca.sort_order,
    count(uav.id) AS usage_count
   FROM (((public.marketplace_categories c
     JOIN public.unified_category_attributes uca ON ((c.id = uca.category_id)))
     JOIN public.unified_attributes ua ON ((uca.attribute_id = ua.id)))
     LEFT JOIN public.unified_attribute_values uav ON ((ua.id = uav.attribute_id)))
  WHERE ((uca.is_enabled = true) AND (ua.is_active = true))
  GROUP BY c.id, c.name, ua.id, ua.code, ua.name, ua.attribute_type, uca.is_required, uca.sort_order
 HAVING (count(uav.id) > 5)
  ORDER BY c.id, uca.sort_order, ua.name
  WITH NO DATA;
CREATE VIEW public.storefront_orders_view AS
 SELECT storefront_orders.id,
    storefront_orders.order_number,
    storefront_orders.storefront_id,
    storefront_orders.customer_id,
    storefront_orders.payment_transaction_id,
    storefront_orders.subtotal_amount AS subtotal,
    storefront_orders.shipping_amount AS shipping,
    storefront_orders.tax_amount AS tax,
    storefront_orders.discount,
    storefront_orders.total_amount AS total,
    storefront_orders.commission_amount,
    storefront_orders.seller_amount,
    storefront_orders.currency,
    storefront_orders.status,
    storefront_orders.escrow_release_date,
    storefront_orders.escrow_days,
    storefront_orders.shipping_address,
    storefront_orders.billing_address,
    storefront_orders.shipping_method,
    storefront_orders.shipping_provider,
    storefront_orders.tracking_number,
    storefront_orders.customer_notes,
    storefront_orders.seller_notes,
    storefront_orders.payment_method,
    storefront_orders.payment_status,
    storefront_orders.notes,
    storefront_orders.metadata,
    storefront_orders.confirmed_at,
    storefront_orders.shipped_at,
    storefront_orders.delivered_at,
    storefront_orders.cancelled_at,
    storefront_orders.created_at,
    storefront_orders.updated_at
   FROM public.storefront_orders;
CREATE MATERIALIZED VIEW public.storefront_rating_distribution AS
 SELECT reviews.entity_origin_id AS storefront_id,
    reviews.rating,
    count(*) AS count
   FROM public.reviews
  WHERE (((reviews.entity_origin_type)::text = 'storefront'::text) AND ((reviews.status)::text = 'published'::text))
  GROUP BY reviews.entity_origin_id, reviews.rating
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.storefront_rating_summary AS
 WITH review_stats AS (
         SELECT COALESCE(r.entity_origin_id, ml.storefront_id) AS storefront_id,
            count(*) AS total_reviews,
            avg(r.rating) AS average_rating,
            count(*) FILTER (WHERE (r.rating = 1)) AS rating_1,
            count(*) FILTER (WHERE (r.rating = 2)) AS rating_2,
            count(*) FILTER (WHERE (r.rating = 3)) AS rating_3,
            count(*) FILTER (WHERE (r.rating = 4)) AS rating_4,
            count(*) FILTER (WHERE (r.rating = 5)) AS rating_5
           FROM (public.reviews r
             JOIN public.marketplace_listings ml ON ((r.entity_id = ml.id)))
          WHERE ((((r.entity_type)::text = 'listing'::text) AND (ml.storefront_id IS NOT NULL) AND (r.entity_origin_type IS NULL)) OR ((r.entity_origin_type)::text = 'storefront'::text))
          GROUP BY COALESCE(r.entity_origin_id, ml.storefront_id)
        )
 SELECT s.id AS storefront_id,
    s.name,
    rs.total_reviews,
    rs.average_rating,
    rs.rating_1,
    rs.rating_2,
    rs.rating_3,
    rs.rating_4,
    rs.rating_5
   FROM (public.user_storefronts s
     LEFT JOIN review_stats rs ON ((s.id = rs.storefront_id)))
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.storefront_ratings AS
 SELECT reviews.entity_origin_id AS storefront_id,
    count(*) AS total_reviews,
    avg(reviews.rating) AS average_rating,
    count(*) FILTER (WHERE ((reviews.entity_type)::text = 'storefront'::text)) AS direct_reviews,
    count(*) FILTER (WHERE ((reviews.entity_type)::text = 'listing'::text)) AS listing_reviews,
    count(*) FILTER (WHERE (reviews.is_verified_purchase = true)) AS verified_reviews,
    count(*) FILTER (WHERE (array_length(reviews.photos, 1) > 0)) AS photo_reviews,
    count(*) FILTER (WHERE (reviews.rating = 1)) AS rating_1,
    count(*) FILTER (WHERE (reviews.rating = 2)) AS rating_2,
    count(*) FILTER (WHERE (reviews.rating = 3)) AS rating_3,
    count(*) FILTER (WHERE (reviews.rating = 4)) AS rating_4,
    count(*) FILTER (WHERE (reviews.rating = 5)) AS rating_5,
    avg(reviews.rating) FILTER (WHERE (reviews.created_at >= (now() - '30 days'::interval))) AS recent_rating,
    count(*) FILTER (WHERE (reviews.created_at >= (now() - '30 days'::interval))) AS recent_reviews,
    max(reviews.created_at) AS last_review_at
   FROM public.reviews
  WHERE (((reviews.entity_origin_type)::text = 'storefront'::text) AND ((reviews.status)::text = 'published'::text))
  GROUP BY reviews.entity_origin_id
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.user_rating_distribution AS
 SELECT reviews.entity_origin_id AS user_id,
    reviews.rating,
    count(*) AS count
   FROM public.reviews
  WHERE (((reviews.entity_origin_type)::text = 'user'::text) AND ((reviews.status)::text = 'published'::text))
  GROUP BY reviews.entity_origin_id, reviews.rating
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.user_rating_summary AS
 WITH review_stats AS (
         SELECT COALESCE(reviews.entity_origin_id, reviews.user_id) AS user_id,
            count(*) AS total_reviews,
            avg(reviews.rating) AS average_rating,
            count(*) FILTER (WHERE (reviews.rating = 1)) AS rating_1,
            count(*) FILTER (WHERE (reviews.rating = 2)) AS rating_2,
            count(*) FILTER (WHERE (reviews.rating = 3)) AS rating_3,
            count(*) FILTER (WHERE (reviews.rating = 4)) AS rating_4,
            count(*) FILTER (WHERE (reviews.rating = 5)) AS rating_5
           FROM public.reviews
          WHERE ((((reviews.entity_type)::text = 'listing'::text) AND (reviews.entity_origin_type IS NULL)) OR ((reviews.entity_origin_type)::text = 'user'::text))
          GROUP BY COALESCE(reviews.entity_origin_id, reviews.user_id)
        )
 SELECT u.id AS user_id,
    u.name,
    rs.total_reviews,
    rs.average_rating,
    rs.rating_1,
    rs.rating_2,
    rs.rating_3,
    rs.rating_4,
    rs.rating_5
   FROM (public.users u
     LEFT JOIN review_stats rs ON ((u.id = rs.user_id)))
  WITH NO DATA;
CREATE MATERIALIZED VIEW public.user_ratings AS
 SELECT reviews.entity_origin_id AS user_id,
    count(*) AS total_reviews,
    avg(reviews.rating) AS average_rating,
    count(*) FILTER (WHERE ((reviews.entity_type)::text = 'user'::text)) AS direct_reviews,
    count(*) FILTER (WHERE ((reviews.entity_type)::text = 'listing'::text)) AS listing_reviews,
    count(*) FILTER (WHERE ((reviews.entity_type)::text = 'storefront'::text)) AS storefront_reviews,
    count(*) FILTER (WHERE (reviews.is_verified_purchase = true)) AS verified_reviews,
    count(*) FILTER (WHERE (array_length(reviews.photos, 1) > 0)) AS photo_reviews,
    count(*) FILTER (WHERE (reviews.rating = 1)) AS rating_1,
    count(*) FILTER (WHERE (reviews.rating = 2)) AS rating_2,
    count(*) FILTER (WHERE (reviews.rating = 3)) AS rating_3,
    count(*) FILTER (WHERE (reviews.rating = 4)) AS rating_4,
    count(*) FILTER (WHERE (reviews.rating = 5)) AS rating_5,
    avg(reviews.rating) FILTER (WHERE (reviews.created_at >= (now() - '30 days'::interval))) AS recent_rating,
    count(*) FILTER (WHERE (reviews.created_at >= (now() - '30 days'::interval))) AS recent_reviews,
    max(reviews.created_at) AS last_review_at
   FROM public.reviews
  WHERE (((reviews.entity_origin_type)::text = 'user'::text) AND ((reviews.status)::text = 'published'::text))
  GROUP BY reviews.entity_origin_id
  WITH NO DATA;
CREATE VIEW public.user_role_permissions AS
 SELECT u.id AS user_id,
    u.email,
    u.name AS user_name,
    r.id AS role_id,
    r.name AS role_name,
    r.display_name AS role_display_name,
    p.id AS permission_id,
    p.name AS permission_name,
    p.resource,
    p.action
   FROM (((public.users u
     LEFT JOIN public.roles r ON ((u.role_id = r.id)))
     LEFT JOIN public.role_permissions rp ON ((r.id = rp.role_id)))
     LEFT JOIN public.permissions p ON ((rp.permission_id = p.id)))
  WHERE ((u.account_status)::text = 'active'::text);
CREATE INDEX idx_address_log_change_reason ON public.address_change_log USING btree (change_reason);
CREATE INDEX idx_address_log_new_location ON public.address_change_log USING gist (new_location);
CREATE INDEX idx_address_log_old_location ON public.address_change_log USING gist (old_location);
CREATE INDEX idx_couriers_location ON public.couriers USING btree (current_latitude, current_longitude) WHERE (is_online = true);
CREATE INDEX idx_filter_analytics_session ON public.gis_filter_analytics USING btree (session_id);
CREATE INDEX idx_geocoding_cache_location ON public.geocoding_cache USING gist (location);
CREATE INDEX idx_listings_geo_blurred_location ON public.listings_geo USING gist (blurred_location);
CREATE INDEX idx_listings_geo_location ON public.listings_geo USING gist (location);
CREATE INDEX idx_map_items_cache_location ON public.map_items_cache USING btree (latitude, longitude);
CREATE INDEX idx_marketplace_listings_location ON public.marketplace_listings USING btree (latitude, longitude) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX idx_marketplace_orders_protection ON public.marketplace_orders USING btree (protection_expires_at) WHERE ((status)::text = ANY (ARRAY[('delivered'::character varying)::text, ('shipped'::character varying)::text]));
CREATE INDEX idx_poi_cache_location ON public.gis_poi_cache USING gist (location);
CREATE INDEX idx_search_weights_history_reason ON public.search_weights_history USING btree (change_reason);
CREATE INDEX idx_search_weights_version ON public.search_weights USING btree (version);
CREATE INDEX idx_storefront_products_individual_location ON public.storefront_products USING btree (has_individual_location);
CREATE INDEX idx_unified_geo_location ON public.unified_geo USING gist (location);
CREATE INDEX idx_viber_messages_user_session ON public.viber_messages USING btree (viber_user_id, session_id, created_at DESC);
CREATE INDEX admin_users_email_idx ON public.admin_users USING btree (email);
CREATE INDEX idx_address_log_confidence_after ON public.address_change_log USING btree (confidence_after);
CREATE INDEX idx_address_log_created_at ON public.address_change_log USING btree (created_at);
CREATE INDEX idx_address_log_listing_id ON public.address_change_log USING btree (listing_id);
CREATE INDEX idx_address_log_user_id ON public.address_change_log USING btree (user_id);
CREATE INDEX idx_attribute_group_items_attribute ON public.attribute_group_items USING btree (attribute_id);
CREATE INDEX idx_attribute_group_items_group ON public.attribute_group_items USING btree (group_id);
CREATE INDEX idx_attribute_option_translations ON public.attribute_option_translations USING btree (attribute_name, option_value);
CREATE INDEX idx_audit_log_entity ON public.translation_audit_log USING btree (entity_type, entity_id);
CREATE INDEX idx_audit_log_user_date ON public.translation_audit_log USING btree (user_id, created_at DESC);
CREATE INDEX idx_car_generations_active ON public.car_generations USING btree (is_active) WHERE (is_active = true);
CREATE INDEX idx_car_generations_model_id ON public.car_generations USING btree (model_id);
CREATE INDEX idx_car_generations_slug ON public.car_generations USING btree (slug);
CREATE INDEX idx_car_generations_years ON public.car_generations USING btree (year_start, year_end);
CREATE INDEX idx_car_makes_external_id ON public.car_makes USING btree (external_id);
CREATE INDEX idx_car_makes_is_domestic ON public.car_makes USING btree (is_domestic);
