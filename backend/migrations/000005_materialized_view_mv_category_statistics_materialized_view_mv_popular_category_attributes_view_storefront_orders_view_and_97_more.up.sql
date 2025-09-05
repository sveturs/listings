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
CREATE INDEX idx_custom_ui_components_name ON public.custom_ui_components USING btree (name);
CREATE INDEX idx_custom_ui_components_type ON public.custom_ui_components USING btree (component_type);
CREATE INDEX idx_custom_ui_templates_name ON public.custom_ui_templates USING btree (name);
CREATE INDEX idx_delivery_is_active ON public.storefront_delivery_options USING btree (is_active);
CREATE INDEX idx_delivery_storefront_id ON public.storefront_delivery_options USING btree (storefront_id);
CREATE INDEX idx_density_grid_cell ON public.gis_listing_density_grid USING gist (cell);
CREATE INDEX idx_density_grid_density ON public.gis_listing_density_grid USING btree (density);
CREATE INDEX idx_disputes_review_id ON public.review_disputes USING btree (review_id);
CREATE INDEX idx_disputes_status ON public.review_disputes USING btree (status);
CREATE INDEX idx_districts_boundary ON public.districts USING gist (boundary);
CREATE INDEX idx_districts_center ON public.districts USING gist (center_point);
CREATE INDEX idx_districts_city_id ON public.districts USING btree (city_id);
CREATE INDEX idx_districts_country ON public.districts USING btree (country_code);
CREATE INDEX idx_districts_name ON public.districts USING btree (name);
CREATE INDEX idx_escrow_payments_buyer_id ON public.escrow_payments USING btree (buyer_id);
CREATE INDEX idx_escrow_payments_payment_transaction_id ON public.escrow_payments USING btree (payment_transaction_id);
CREATE INDEX idx_escrow_payments_seller_id ON public.escrow_payments USING btree (seller_id);
CREATE INDEX idx_escrow_payments_status ON public.escrow_payments USING btree (status);
CREATE INDEX idx_filter_analytics_created ON public.gis_filter_analytics USING btree (created_at);
CREATE INDEX idx_filter_analytics_type ON public.gis_filter_analytics USING btree (filter_type);
CREATE INDEX idx_filter_analytics_user ON public.gis_filter_analytics USING btree (user_id);
CREATE INDEX idx_geocoding_cache_address_components ON public.geocoding_cache USING gin (address_components);
