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
ALTER TABLE ONLY public.user_notification_contacts ALTER COLUMN id SET DEFAULT nextval('public.user_notification_contacts_id_seq'::regclass);
ALTER TABLE ONLY public.user_notification_preferences ALTER COLUMN id SET DEFAULT nextval('public.user_notification_preferences_id_seq'::regclass);
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
