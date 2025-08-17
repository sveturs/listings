ALTER TABLE ONLY public.user_behavior_events ALTER COLUMN id SET DEFAULT nextval('public.user_behavior_events_id_seq'::regclass);
ALTER TABLE ONLY public.user_contacts ALTER COLUMN id SET DEFAULT nextval('public.user_contacts_id_seq'::regclass);
ALTER TABLE ONLY public.user_storefronts ALTER COLUMN id SET DEFAULT nextval('public.user_storefronts_id_seq'::regclass);
ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
ALTER TABLE ONLY public.variant_attribute_mappings ALTER COLUMN id SET DEFAULT nextval('public.variant_attribute_mappings_id_seq'::regclass);
ALTER TABLE ONLY public.warehouse_inventory ALTER COLUMN id SET DEFAULT nextval('public.warehouse_inventory_id_seq'::regclass);
ALTER TABLE ONLY public.warehouse_invoices ALTER COLUMN id SET DEFAULT nextval('public.warehouse_invoices_id_seq'::regclass);
ALTER TABLE ONLY public.warehouse_movements ALTER COLUMN id SET DEFAULT nextval('public.warehouse_movements_id_seq'::regclass);
ALTER TABLE ONLY public.warehouse_pickup_orders ALTER COLUMN id SET DEFAULT nextval('public.warehouse_pickup_orders_id_seq'::regclass);
ALTER TABLE ONLY public.warehouses ALTER COLUMN id SET DEFAULT nextval('public.warehouses_id_seq'::regclass);
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
CREATE VIEW public.v_attribute_groups_with_items AS
 SELECT ag.id AS group_id,
    ag.name AS group_name,
    ag.display_name AS group_display_name,
    ag.icon AS group_icon,
    ag.sort_order AS group_sort_order,
    agi.id AS item_id,
    agi.attribute_id,
    ca.name AS attribute_name,
    ca.display_name AS attribute_display_name,
    agi.icon AS attribute_icon,
    agi.custom_display_name,
    agi.sort_order AS attribute_sort_order
   FROM ((public.attribute_groups ag
     LEFT JOIN public.attribute_group_items agi ON ((ag.id = agi.group_id)))
     LEFT JOIN public.category_attributes ca ON ((agi.attribute_id = ca.id)))
  WHERE (ag.is_active = true)
  ORDER BY ag.sort_order, agi.sort_order;
CREATE VIEW public.v_category_attributes AS
 SELECT cam.category_id,
    cam.attribute_id,
    cam.is_enabled,
    cam.is_required,
    cam.sort_order,
    cam.custom_component,
    ca.name,
    ca.display_name,
    ca.attribute_type,
    ca.options,
    ca.validation_rules,
    ca.is_searchable,
    ca.is_filterable,
    ca.custom_component AS default_custom_component,
    mc.name AS category_name,
    mc.slug AS category_slug
   FROM ((public.category_attribute_mapping cam
     JOIN public.category_attributes ca ON ((cam.attribute_id = ca.id)))
     JOIN public.marketplace_categories mc ON ((cam.category_id = mc.id)))
  ORDER BY cam.category_id, cam.sort_order, ca.sort_order;
CREATE INDEX idx_address_log_change_reason ON public.address_change_log USING btree (change_reason);
CREATE INDEX idx_address_log_new_location ON public.address_change_log USING gist (new_location);
CREATE INDEX idx_address_log_old_location ON public.address_change_log USING gist (old_location);
CREATE INDEX idx_filter_analytics_session ON public.gis_filter_analytics USING btree (session_id);
CREATE INDEX idx_geocoding_cache_location ON public.geocoding_cache USING gist (location);
CREATE INDEX idx_listings_geo_blurred_location ON public.listings_geo USING gist (blurred_location);
CREATE INDEX idx_listings_geo_location ON public.listings_geo USING gist (location);
CREATE INDEX idx_map_items_cache_location ON public.map_items_cache USING btree (latitude, longitude);
CREATE INDEX idx_marketplace_listings_location ON public.marketplace_listings USING btree (latitude, longitude) WHERE ((latitude IS NOT NULL) AND (longitude IS NOT NULL));
CREATE INDEX idx_marketplace_orders_protection ON public.marketplace_orders USING btree (protection_expires_at) WHERE ((status)::text = ANY (ARRAY[('delivered'::character varying)::text, ('shipped'::character varying)::text]));
CREATE INDEX idx_poi_cache_location ON public.gis_poi_cache USING gist (location);
CREATE INDEX idx_post_express_api_logs_transaction ON public.post_express_api_logs USING btree (transaction_id);
CREATE INDEX idx_post_express_offices_location ON public.post_express_offices USING btree (location_id);
CREATE INDEX idx_search_weights_history_reason ON public.search_weights_history USING btree (change_reason);
CREATE INDEX idx_search_weights_version ON public.search_weights USING btree (version);
CREATE INDEX idx_storefront_products_individual_location ON public.storefront_products USING btree (has_individual_location);
CREATE INDEX idx_unified_geo_location ON public.unified_geo USING gist (location);
CREATE INDEX admin_users_email_idx ON public.admin_users USING btree (email);
CREATE INDEX idx_address_log_confidence_after ON public.address_change_log USING btree (confidence_after);
CREATE INDEX idx_address_log_created_at ON public.address_change_log USING btree (created_at);
CREATE INDEX idx_address_log_listing_id ON public.address_change_log USING btree (listing_id);
CREATE INDEX idx_address_log_user_id ON public.address_change_log USING btree (user_id);
CREATE INDEX idx_attr_name_num_val ON public.listing_attribute_values USING btree (attribute_id, numeric_value) WHERE (numeric_value IS NOT NULL);
CREATE INDEX idx_attr_name_text_val ON public.listing_attribute_values USING btree (attribute_id, text_value) WHERE (text_value IS NOT NULL);
CREATE INDEX idx_attr_unit ON public.listing_attribute_values USING btree (unit) WHERE (unit IS NOT NULL);
CREATE INDEX idx_attribute_group_items_attribute ON public.attribute_group_items USING btree (attribute_id);
CREATE INDEX idx_attribute_group_items_group ON public.attribute_group_items USING btree (group_id);
CREATE INDEX idx_attribute_groups_active ON public.attribute_groups USING btree (is_active);
CREATE INDEX idx_attribute_groups_name ON public.attribute_groups USING btree (name);
CREATE INDEX idx_attribute_option_translations ON public.attribute_option_translations USING btree (attribute_name, option_value);
CREATE INDEX idx_audit_log_entity ON public.translation_audit_log USING btree (entity_type, entity_id);
CREATE INDEX idx_audit_log_user_date ON public.translation_audit_log USING btree (user_id, created_at DESC);
CREATE INDEX idx_car_generations_active ON public.car_generations USING btree (is_active) WHERE (is_active = true);
CREATE INDEX idx_car_generations_model_id ON public.car_generations USING btree (model_id);
CREATE INDEX idx_car_generations_slug ON public.car_generations USING btree (slug);
CREATE INDEX idx_car_generations_years ON public.car_generations USING btree (year_start, year_end);
CREATE INDEX idx_car_makes_external_id ON public.car_makes USING btree (external_id);
CREATE INDEX idx_car_makes_is_domestic ON public.car_makes USING btree (is_domestic);
CREATE INDEX idx_car_makes_popularity ON public.car_makes USING btree (popularity_rs DESC);
CREATE INDEX idx_car_makes_slug ON public.car_makes USING btree (slug);
CREATE INDEX idx_car_models_battery_capacity ON public.car_models USING btree (battery_capacity_kwh) WHERE (battery_capacity_kwh IS NOT NULL);
CREATE INDEX idx_car_models_body_type ON public.car_models USING btree (body_type);
CREATE INDEX idx_car_models_drive_type ON public.car_models USING btree (drive_type);
CREATE INDEX idx_car_models_engine_type ON public.car_models USING btree (engine_type);
CREATE INDEX idx_car_models_external_id ON public.car_models USING btree (external_id);
CREATE INDEX idx_car_models_fuel_type ON public.car_models USING btree (fuel_type);
CREATE INDEX idx_car_models_is_electric ON public.car_models USING btree (is_electric) WHERE (is_electric = true);
CREATE INDEX idx_car_models_make_slug ON public.car_models USING btree (make_id, slug);
CREATE INDEX idx_car_models_search ON public.car_models USING gin (to_tsvector('simple'::regconfig, (name)::text));
CREATE INDEX idx_car_models_serbia_popularity ON public.car_models USING btree (serbia_popularity_score DESC);
CREATE INDEX idx_car_models_transmission_type ON public.car_models USING btree (transmission_type);
CREATE INDEX idx_category_attr_mapping ON public.category_attribute_mapping USING btree (category_id, is_enabled) WHERE (is_enabled = true);
CREATE INDEX idx_category_attribute_groups_category ON public.category_attribute_groups USING btree (category_id);
CREATE INDEX idx_category_attribute_groups_component ON public.category_attribute_groups USING btree (component_id);
CREATE INDEX idx_category_attribute_groups_group ON public.category_attribute_groups USING btree (group_id);
CREATE INDEX idx_category_attribute_map_attr_id ON public.category_attribute_mapping USING btree (attribute_id);
CREATE INDEX idx_category_attribute_map_cat_id ON public.category_attribute_mapping USING btree (category_id);
CREATE INDEX idx_category_attribute_mapping_custom_component ON public.category_attribute_mapping USING btree (custom_component);
CREATE INDEX idx_category_attributes_data_source ON public.category_attributes USING btree (data_source);
CREATE INDEX idx_category_attributes_name ON public.category_attributes USING btree (name);
CREATE INDEX idx_category_attributes_variant_compatible ON public.category_attributes USING btree (is_variant_compatible) WHERE (is_variant_compatible = true);
CREATE INDEX idx_category_variant_attributes_category_id ON public.category_variant_attributes USING btree (category_id);
CREATE INDEX idx_category_variant_attributes_variant_name ON public.category_variant_attributes USING btree (variant_attribute_name);
CREATE INDEX idx_category_weight ON public.category_keywords USING btree (category_id, weight DESC);
CREATE INDEX idx_chat_attachments_created_at ON public.chat_attachments USING btree (created_at);
CREATE INDEX idx_chat_attachments_file_type ON public.chat_attachments USING btree (file_type);
CREATE INDEX idx_chat_attachments_message ON public.chat_attachments USING btree (message_id);
CREATE INDEX idx_chat_attachments_message_id ON public.chat_attachments USING btree (message_id);
CREATE INDEX idx_cities_boundary ON public.cities USING gist (boundary);
CREATE INDEX idx_cities_center ON public.cities USING gist (center_point);
CREATE INDEX idx_cities_country ON public.cities USING btree (country_code);
CREATE INDEX idx_cities_has_districts ON public.cities USING btree (has_districts);
CREATE INDEX idx_cities_name ON public.cities USING btree (name);
CREATE INDEX idx_cities_priority ON public.cities USING btree (priority);
CREATE INDEX idx_cities_slug ON public.cities USING btree (slug);
CREATE INDEX idx_component_templates_cat ON public.component_templates USING btree (category_id);
CREATE INDEX idx_component_templates_comp ON public.component_templates USING btree (component_id);
CREATE INDEX idx_custom_ui_components_active ON public.custom_ui_components USING btree (is_active);
