-- Materialized views migration

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

REFRESH MATERIALIZED VIEW public.category_listing_counts;

CREATE MATERIALIZED VIEW public.map_items_cache AS
 WITH combined_items AS (
         SELECT ml.id,
            ml.title AS name,
            ml.description,
            ml.price,
            c.name AS category_name,
            ug.location,
            public.st_y((ug.location)::public.geometry) AS latitude,
            public.st_x((ug.location)::public.geometry) AS longitude,
            ug.formatted_address,
            ml.user_id,
            ml.storefront_id,
            ml.status,
            ml.created_at,
            ml.updated_at,
            ml.views_count,
            0 AS rating,
            'marketplace_listing'::text AS item_type,
            COALESCE(ug.privacy_level, 'exact'::public.location_privacy_level) AS privacy_level,
            COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
            'individual'::text AS display_strategy
           FROM ((public.marketplace_listings ml
             JOIN public.unified_geo ug ON (((ug.source_type = 'marketplace_listing'::public.geo_source_type) AND (ug.source_id = ml.id))))
             LEFT JOIN public.marketplace_categories c ON ((ml.category_id = c.id)))
          WHERE (((ml.status)::text = 'active'::text) AND (ml.show_on_map = true))
        UNION ALL
         SELECT s.id,
            s.name,
            s.description,
            0 AS price,
            'Витрина'::character varying AS category_name,
            ug.location,
            public.st_y((ug.location)::public.geometry) AS latitude,
            public.st_x((ug.location)::public.geometry) AS longitude,
            ug.formatted_address,
            s.user_id,
            NULL::integer AS storefront_id,
                CASE
                    WHEN s.is_active THEN 'active'::text
                    ELSE 'inactive'::text
                END AS status,
            s.created_at,
            s.updated_at,
            s.views_count,
            0 AS rating,
            'storefront'::text AS item_type,
            COALESCE(ug.privacy_level, 'exact'::public.location_privacy_level) AS privacy_level,
            COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
            'grouped'::text AS display_strategy
           FROM (public.storefronts s
             JOIN public.unified_geo ug ON (((ug.source_type = 'storefront'::public.geo_source_type) AND (ug.source_id = s.id))))
          WHERE (s.is_active = true)
        UNION ALL
         SELECT sp.id,
            sp.name,
            sp.description,
            COALESCE(spv.price, sp.price) AS price,
            spc.name AS category_name,
            ug.location,
            public.st_y((ug.location)::public.geometry) AS latitude,
            public.st_x((ug.location)::public.geometry) AS longitude,
            ug.formatted_address,
            s.user_id,
            sp.storefront_id,
                CASE
                    WHEN (sp.is_active AND s.is_active) THEN 'active'::text
                    ELSE 'inactive'::text
                END AS status,
            sp.created_at,
            sp.updated_at,
            sp.view_count AS views_count,
            0 AS rating,
            'storefront_product'::text AS item_type,
            COALESCE(ug.privacy_level, 'exact'::public.location_privacy_level) AS privacy_level,
            COALESCE(ug.blur_radius_meters, 0) AS blur_radius_meters,
            'grouped'::text AS display_strategy
           FROM ((((public.storefront_products sp
             JOIN public.storefronts s ON ((sp.storefront_id = s.id)))
             JOIN public.unified_geo ug ON (((ug.source_type = 'storefront'::public.geo_source_type) AND (ug.source_id = s.id))))
             LEFT JOIN public.storefront_product_variants spv ON (((sp.id = spv.product_id) AND (spv.is_default = true))))
             LEFT JOIN public.marketplace_categories spc ON ((sp.category_id = spc.id)))
          WHERE ((sp.is_active = true) AND (s.is_active = true))
        )
 SELECT combined_items.id,
    combined_items.name,
    combined_items.description,
    combined_items.price,
    combined_items.category_name,
    combined_items.location,
    combined_items.latitude,
    combined_items.longitude,
    combined_items.formatted_address,
    combined_items.user_id,
    combined_items.storefront_id,
    combined_items.status,
    combined_items.created_at,
    combined_items.updated_at,
    combined_items.views_count,
    combined_items.rating,
    combined_items.item_type,
    combined_items.privacy_level,
    combined_items.blur_radius_meters,
    combined_items.display_strategy
   FROM combined_items
  WITH NO DATA;

REFRESH MATERIALIZED VIEW public.map_items_cache;

CREATE MATERIALIZED VIEW public.storefront_rating_distribution AS
 SELECT reviews.entity_origin_id AS storefront_id,
    reviews.rating,
    count(*) AS count
   FROM public.reviews
  WHERE (((reviews.entity_origin_type)::text = 'storefront'::text) AND ((reviews.status)::text = 'published'::text))
  GROUP BY reviews.entity_origin_id, reviews.rating
  WITH NO DATA;

REFRESH MATERIALIZED VIEW public.storefront_rating_distribution;

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

REFRESH MATERIALIZED VIEW public.storefront_rating_summary;

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

REFRESH MATERIALIZED VIEW public.storefront_ratings;

CREATE MATERIALIZED VIEW public.user_rating_distribution AS
 SELECT reviews.entity_origin_id AS user_id,
    reviews.rating,
    count(*) AS count
   FROM public.reviews
  WHERE (((reviews.entity_origin_type)::text = 'user'::text) AND ((reviews.status)::text = 'published'::text))
  GROUP BY reviews.entity_origin_id, reviews.rating
  WITH NO DATA;

REFRESH MATERIALIZED VIEW public.user_rating_distribution;

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

REFRESH MATERIALIZED VIEW public.user_rating_summary;

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

REFRESH MATERIALIZED VIEW public.user_ratings;