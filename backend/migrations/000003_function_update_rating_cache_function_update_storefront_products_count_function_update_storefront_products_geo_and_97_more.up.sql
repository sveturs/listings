CREATE FUNCTION public.update_rating_cache() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Обновляем rating_cache для измененного объявления
    INSERT INTO rating_cache (entity_type, entity_id, average_rating, total_reviews, calculated_at)
    SELECT
        'listing' as entity_type,
        COALESCE(NEW.entity_id, OLD.entity_id) as entity_id,
        ROUND(AVG(rating)::numeric, 2) as average_rating,
        COUNT(*) as total_reviews,
        NOW() as calculated_at
    FROM reviews
    WHERE entity_type = 'listing'
      AND entity_id = COALESCE(NEW.entity_id, OLD.entity_id)
      AND status = 'published'
    GROUP BY entity_id
    ON CONFLICT (entity_type, entity_id)
    DO UPDATE SET
        average_rating = EXCLUDED.average_rating,
        total_reviews = EXCLUDED.total_reviews,
        calculated_at = EXCLUDED.calculated_at;
    -- Если нет отзывов, удаляем запись из rating_cache
    IF NOT EXISTS (
        SELECT 1 FROM reviews
        WHERE entity_type = 'listing'
          AND entity_id = COALESCE(NEW.entity_id, OLD.entity_id)
          AND status = 'published'
    ) THEN
        DELETE FROM rating_cache
        WHERE entity_type = 'listing'
          AND entity_id = COALESCE(NEW.entity_id, OLD.entity_id);
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$;
CREATE FUNCTION public.update_storefront_products_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- При вставке или удалении товара обновляем счётчик
    IF TG_OP = 'INSERT' THEN
        UPDATE storefronts
        SET products_count = (
            SELECT COUNT(*)
            FROM storefront_products
            WHERE storefront_id = NEW.storefront_id
            AND is_active = true
        )
        WHERE id = NEW.storefront_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE storefronts
        SET products_count = (
            SELECT COUNT(*)
            FROM storefront_products
            WHERE storefront_id = OLD.storefront_id
            AND is_active = true
        )
        WHERE id = OLD.storefront_id;
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        -- При обновлении проверяем изменение is_active или storefront_id
        IF OLD.is_active != NEW.is_active OR OLD.storefront_id != NEW.storefront_id THEN
            -- Обновляем счётчик для старой витрины
            IF OLD.storefront_id != NEW.storefront_id THEN
                UPDATE storefronts
                SET products_count = (
                    SELECT COUNT(*)
                    FROM storefront_products
                    WHERE storefront_id = OLD.storefront_id
                    AND is_active = true
                )
                WHERE id = OLD.storefront_id;
            END IF;
            -- Обновляем счётчик для новой витрины
            UPDATE storefronts
            SET products_count = (
                SELECT COUNT(*)
                FROM storefront_products
                WHERE storefront_id = NEW.storefront_id
                AND is_active = true
            )
            WHERE id = NEW.storefront_id;
        END IF;
        RETURN NEW;
    END IF;
END;
$$;
CREATE FUNCTION public.update_storefront_products_geo() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- If storefront coordinates changed, update all products that use storefront location
    IF (OLD.latitude IS DISTINCT FROM NEW.latitude OR OLD.longitude IS DISTINCT FROM NEW.longitude) THEN
        -- Trigger re-geocoding for all products that don't have individual locations
        UPDATE storefront_products
        SET updated_at = NOW()
        WHERE storefront_id = NEW.id AND has_individual_location = false;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_storefront_views_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Обновляем views_count для витрины при изменении view_count товара
    IF TG_OP = 'UPDATE' THEN
        -- Обновляем для старой витрины (если товар переместили)
        IF OLD.storefront_id IS DISTINCT FROM NEW.storefront_id AND OLD.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = OLD.storefront_id
            ), 0)
            WHERE id = OLD.storefront_id;
        END IF;
        -- Обновляем для новой/текущей витрины
        IF NEW.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = NEW.storefront_id
            ), 0)
            WHERE id = NEW.storefront_id;
        END IF;
    ELSIF TG_OP = 'INSERT' THEN
        -- При добавлении нового товара
        IF NEW.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = NEW.storefront_id
            ), 0)
            WHERE id = NEW.storefront_id;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении товара
        IF OLD.storefront_id IS NOT NULL THEN
            UPDATE storefronts
            SET views_count = COALESCE((
                SELECT SUM(view_count)
                FROM storefront_products
                WHERE storefront_id = OLD.storefront_id
            ), 0)
            WHERE id = OLD.storefront_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_subscription_usage() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- При создании витрины увеличиваем счетчик
        UPDATE user_subscriptions
        SET used_storefronts = used_storefronts + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.user_id;
        -- Записываем в историю использования
        INSERT INTO subscription_usage (subscription_id, storefront_id, resource_type, action)
        SELECT id, NEW.id, 'storefront', 'created'
        FROM user_subscriptions
        WHERE user_id = NEW.user_id
        LIMIT 1;
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении витрины уменьшаем счетчик
        UPDATE user_subscriptions
        SET used_storefronts = GREATEST(0, used_storefronts - 1),
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = OLD.user_id;
        -- Записываем в историю использования
        INSERT INTO subscription_usage (subscription_id, storefront_id, resource_type, action)
        SELECT id, OLD.id, 'storefront', 'deleted'
        FROM user_subscriptions
        WHERE user_id = OLD.user_id
        LIMIT 1;
    END IF;
    RETURN NEW;
END;
$$;
CREATE TABLE public.category_attribute_mapping (
    category_id integer NOT NULL,
    attribute_id integer NOT NULL,
    is_enabled boolean DEFAULT true,
    is_required boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    custom_component character varying(255),
    show_in_card boolean,
    show_in_list boolean
);
CREATE TABLE public.cities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    slug character varying(255) NOT NULL,
    country_code character varying(2) DEFAULT 'RS'::character varying NOT NULL,
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326) NOT NULL,
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    has_districts boolean DEFAULT false NOT NULL,
    priority integer DEFAULT 100 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.districts_leskovac_backup (
    id uuid,
    name character varying(255),
    city_id uuid,
    country_code character varying(2),
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    city_name character varying(255)
);
CREATE TABLE public.districts_novi_sad_backup_20250715 (
    id uuid,
    name character varying(255),
    city_id uuid,
    country_code character varying(2),
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);
CREATE TABLE public.map_items_cache (
    id integer NOT NULL,
    item_type character varying(50) NOT NULL,
    item_id integer NOT NULL,
    latitude numeric(10,7) NOT NULL,
    longitude numeric(10,7) NOT NULL,
    title text NOT NULL,
    price numeric(20,5),
    category_id integer,
    category_name character varying(255),
    image_url text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at timestamp with time zone
);
CREATE TABLE public.rating_cache (
    entity_type character varying(50) NOT NULL,
    entity_id integer NOT NULL,
    average_rating numeric(3,2),
    total_reviews integer DEFAULT 0,
    distribution jsonb,
    breakdown jsonb,
    verified_percentage integer DEFAULT 0,
    recent_trend character varying(10),
    calculated_at timestamp without time zone DEFAULT now(),
    CONSTRAINT rating_cache_recent_trend_check CHECK (((recent_trend)::text = ANY (ARRAY[('up'::character varying)::text, ('down'::character varying)::text, ('stable'::character varying)::text])))
);
CREATE TABLE public.unit_translations (
    unit character varying(20) NOT NULL,
    language character varying(10) NOT NULL,
    translated_unit character varying(20) NOT NULL,
    display_format character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.user_balances (
    user_id integer NOT NULL,
    balance numeric(12,2) DEFAULT 0 NOT NULL,
    frozen_balance numeric(12,2) DEFAULT 0 NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.user_privacy_settings (
    user_id integer NOT NULL,
    allow_contact_requests boolean DEFAULT true,
    allow_messages_from_contacts_only boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.districts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    city_id uuid,
    country_code character varying(2) DEFAULT 'RS'::character varying NOT NULL,
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.municipalities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    district_id uuid,
    country_code character varying(2) DEFAULT 'RS'::character varying NOT NULL,
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    population integer,
    area_km2 numeric(10,2),
    postal_codes text[],
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.address_change_log (
    id bigint NOT NULL,
    listing_id bigint NOT NULL,
    user_id bigint NOT NULL,
    old_address text,
    new_address text,
    old_location public.geography(Point,4326),
    new_location public.geography(Point,4326),
    change_reason character varying(100),
    confidence_before numeric(3,2),
    confidence_after numeric(3,2),
    ip_address inet,
    user_agent text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.admin_users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    created_by integer,
    notes text
);
CREATE TABLE public.attribute_group_items (
    id integer NOT NULL,
    group_id integer NOT NULL,
    attribute_id integer NOT NULL,
    icon character varying(100),
    sort_order integer DEFAULT 0,
    custom_display_name character varying(255),
    visibility_condition jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.attribute_groups (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    display_name character varying(255),
    icon character varying(100),
    is_system boolean DEFAULT false
);
CREATE TABLE public.attribute_option_translations (
    id integer NOT NULL,
    attribute_name character varying(100) NOT NULL,
    option_value text NOT NULL,
    ru_translation text NOT NULL,
    sr_translation text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.balance_transactions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    type character varying(20) NOT NULL,
    amount numeric(12,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    payment_method character varying(50),
    payment_details jsonb,
    description text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    completed_at timestamp without time zone
);
CREATE TABLE public.car_generations (
    id integer NOT NULL,
    model_id integer,
    name character varying(100),
    year_start integer,
    year_end integer,
    facelift_year integer,
    specs jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    body_types jsonb DEFAULT '[]'::jsonb,
    engine_types jsonb DEFAULT '[]'::jsonb,
    image_url character varying(500),
    is_active boolean DEFAULT true,
    slug character varying(100) NOT NULL,
    sort_order integer DEFAULT 0,
    external_id character varying(100),
    platform character varying(100),
    production_country character varying(100),
    metadata jsonb DEFAULT '{}'::jsonb,
    last_sync_at timestamp without time zone
);
CREATE TABLE public.car_makes (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    logo_url character varying(500),
    country character varying(50),
    is_active boolean DEFAULT true,
    sort_order integer DEFAULT 0,
    is_domestic boolean DEFAULT false,
    popularity_rs integer DEFAULT 0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    external_id character varying(100),
    manufacturer_id character varying(100),
    last_sync_at timestamp without time zone,
    metadata jsonb DEFAULT '{}'::jsonb
);
CREATE TABLE public.car_market_analysis (
    id integer NOT NULL,
    brand character varying(100),
    model character varying(100),
    total_listings integer,
    avg_price numeric(10,2),
    min_year integer,
    max_year integer,
    popularity_score numeric(5,2),
    analyzed_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.car_models (
    id integer NOT NULL,
    make_id integer,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    generation character varying(50),
    production_start integer,
    production_end integer,
    body_types jsonb,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    sort_order integer DEFAULT 0,
    external_id character varying(100),
    body_type character varying(50),
    segment character varying(50),
    years_range int4range,
    metadata jsonb DEFAULT '{}'::jsonb,
    last_sync_at timestamp without time zone,
    engine_type character varying(50),
    engine_power_kw integer,
    engine_power_hp integer,
    engine_torque_nm integer,
    fuel_type character varying(50),
    fuel_consumption_city numeric(4,2),
    fuel_consumption_highway numeric(4,2),
    fuel_consumption_combined numeric(4,2),
    co2_emissions integer,
    euro_standard character varying(10),
    transmission_type character varying(50),
    transmission_gears integer,
    drive_type character varying(50),
    length_mm integer,
    width_mm integer,
    height_mm integer,
    wheelbase_mm integer,
    trunk_volume_l integer,
    fuel_tank_l integer,
    weight_kg integer,
    max_speed_kmh integer,
    acceleration_0_100 numeric(3,1),
    seats integer,
    doors integer,
    serbia_popularity_score integer DEFAULT 0,
    serbia_average_price_eur numeric(10,2),
    serbia_listings_count integer DEFAULT 0,
    is_electric boolean DEFAULT false,
    battery_capacity_kwh numeric(10,2),
    battery_capacity_net_kwh numeric(10,2),
    electric_range_km integer,
    electric_range_wltp_km integer,
    electric_range_standard character varying(50),
    charging_time_0_100 numeric(10,2),
    charging_time_10_80 numeric(10,2),
    fast_charging_power_kw numeric(10,2),
    onboard_charger_kw numeric(10,2),
    popularity_score numeric(5,2) DEFAULT 0
);
CREATE TABLE public.category_attribute_groups (
    id integer NOT NULL,
    category_id integer NOT NULL,
    group_id integer NOT NULL,
    component_id integer,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    display_mode character varying(50) DEFAULT 'list'::character varying,
    collapsed_by_default boolean DEFAULT false,
    configuration jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.unified_attributes (
    id integer NOT NULL,
    code character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    display_name character varying(200) NOT NULL,
    attribute_type character varying(50) NOT NULL,
    purpose character varying(20) DEFAULT 'regular'::character varying NOT NULL,
    options jsonb DEFAULT '{}'::jsonb,
    validation_rules jsonb DEFAULT '{}'::jsonb,
    ui_settings jsonb DEFAULT '{}'::jsonb,
    is_searchable boolean DEFAULT false,
    is_filterable boolean DEFAULT false,
    is_required boolean DEFAULT false,
    affects_stock boolean DEFAULT false,
    affects_price boolean DEFAULT false,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    legacy_category_attribute_id integer,
    legacy_product_variant_attribute_id integer,
    search_vector tsvector,
    is_variant_compatible boolean DEFAULT false,
    icon character varying(255) DEFAULT ''::character varying,
    CONSTRAINT unified_attributes_attribute_type_check CHECK (((attribute_type)::text = ANY (ARRAY[('text'::character varying)::text, ('textarea'::character varying)::text, ('number'::character varying)::text, ('boolean'::character varying)::text, ('select'::character varying)::text, ('multiselect'::character varying)::text, ('date'::character varying)::text, ('color'::character varying)::text, ('size'::character varying)::text]))),
    CONSTRAINT unified_attributes_purpose_check CHECK (((purpose)::text = ANY (ARRAY[('regular'::character varying)::text, ('variant'::character varying)::text, ('both'::character varying)::text])))
)
WITH (autovacuum_vacuum_scale_factor='0.2', autovacuum_analyze_scale_factor='0.1', autovacuum_vacuum_threshold='50', autovacuum_analyze_threshold='50');
CREATE TABLE public.category_keywords (
    id integer NOT NULL,
    category_id integer NOT NULL,
    keyword character varying(100) NOT NULL,
    language character varying(2) DEFAULT 'en'::character varying,
    weight double precision DEFAULT 1.0,
    keyword_type character varying(20) DEFAULT 'general'::character varying,
    is_negative boolean DEFAULT false,
    source character varying(50) DEFAULT 'manual'::character varying,
    usage_count integer DEFAULT 0,
    success_rate double precision DEFAULT 0.0,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT category_keywords_keyword_type_check CHECK (((keyword_type)::text = ANY (ARRAY[('main'::character varying)::text, ('synonym'::character varying)::text, ('brand'::character varying)::text, ('attribute'::character varying)::text, ('context'::character varying)::text, ('pattern'::character varying)::text]))),
    CONSTRAINT category_keywords_source_check CHECK (((source)::text = ANY (ARRAY[('manual'::character varying)::text, ('ai_extracted'::character varying)::text, ('user_confirmed'::character varying)::text, ('auto_learned'::character varying)::text]))),
    CONSTRAINT category_keywords_success_rate_check CHECK (((success_rate >= (0.0)::double precision) AND (success_rate <= (1.0)::double precision))),
    CONSTRAINT category_keywords_weight_check CHECK (((weight >= (0.0)::double precision) AND (weight <= (10.0)::double precision)))
);
CREATE TABLE public.marketplace_categories (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(100) NOT NULL,
    parent_id integer,
    icon character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_custom_ui boolean DEFAULT false,
    custom_ui_component character varying(255),
    sort_order integer DEFAULT 0,
    level integer DEFAULT 0,
    count integer DEFAULT 0,
    external_id character varying(255),
    description text,
    is_active boolean DEFAULT true,
    seo_title character varying(255),
    seo_description text,
    seo_keywords text,
    CONSTRAINT check_root_categories_level CHECK ((((parent_id IS NULL) AND (level = 0)) OR ((parent_id IS NOT NULL) AND (level > 0))))
);
CREATE TABLE public.marketplace_listings (
    id integer DEFAULT nextval('public.global_product_id_seq'::regclass) NOT NULL,
    user_id integer,
    category_id integer,
    title character varying(255) NOT NULL,
    description text,
    price numeric(12,2),
    condition character varying(50),
    status character varying(20) DEFAULT 'active'::character varying,
    location character varying(255),
    latitude numeric(10,8),
    longitude numeric(11,8),
    address_city character varying(100),
    address_country character varying(100),
    views_count integer DEFAULT 0,
    show_on_map boolean DEFAULT true NOT NULL,
    original_language character varying(10) DEFAULT 'sr'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    storefront_id integer,
    external_id character varying(255),
    metadata jsonb,
    needs_reindex boolean DEFAULT false
)
WITH (autovacuum_vacuum_threshold='100', autovacuum_analyze_threshold='100', autovacuum_vacuum_scale_factor='0.1', autovacuum_analyze_scale_factor='0.05');
CREATE TABLE public.category_variant_attributes (
    id integer NOT NULL,
    category_id integer NOT NULL,
    variant_attribute_name character varying(100) NOT NULL,
    sort_order integer DEFAULT 0,
    is_required boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.chat_attachments (
    id integer NOT NULL,
    message_id integer NOT NULL,
    file_type character varying(20) NOT NULL,
    file_path character varying(500) NOT NULL,
    file_name character varying(255) NOT NULL,
    file_size bigint NOT NULL,
    content_type character varying(100) NOT NULL,
    storage_type character varying(20) DEFAULT 'minio'::character varying,
    storage_bucket character varying(100) DEFAULT 'chat-files'::character varying,
    public_url text,
    thumbnail_url text,
    metadata jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chat_attachments_file_type_check CHECK (((file_type)::text = ANY (ARRAY[('image'::character varying)::text, ('video'::character varying)::text, ('document'::character varying)::text])))
);
CREATE TABLE public.component_templates (
    id integer NOT NULL,
    component_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    template_config jsonb DEFAULT '{}'::jsonb,
    preview_image text,
    category_id integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer
);
CREATE TABLE public.courier_location_history (
    id integer NOT NULL,
    delivery_id integer,
    courier_id integer,
    latitude numeric(10,8) NOT NULL,
    longitude numeric(11,8) NOT NULL,
    altitude_meters numeric(7,2),
    speed_kmh numeric(5,2),
    heading integer,
    accuracy_meters numeric(6,2),
    battery_level integer,
    is_mock_location boolean DEFAULT false,
    recorded_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_key_point boolean DEFAULT false,
    CONSTRAINT courier_location_history_heading_check CHECK (((heading >= 0) AND (heading <= 360)))
);
CREATE TABLE public.courier_zones (
    id integer NOT NULL,
    courier_id integer,
    zone_name character varying(100),
    polygon jsonb NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.couriers (
    id integer NOT NULL,
    user_id integer,
    name character varying(255) NOT NULL,
    phone character varying(50),
    photo_url text,
    vehicle_type character varying(50),
    vehicle_number character varying(50),
    is_online boolean DEFAULT false,
    is_available boolean DEFAULT true,
    current_latitude numeric(10,8),
    current_longitude numeric(11,8),
    current_heading integer,
    current_speed numeric(5,2),
    last_location_update timestamp with time zone,
    rating numeric(3,2) DEFAULT 5.0,
    total_deliveries integer DEFAULT 0,
    total_distance_km numeric(10,2) DEFAULT 0,
    working_hours jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT couriers_vehicle_type_check CHECK (((vehicle_type)::text = ANY ((ARRAY['bike'::character varying, 'car'::character varying, 'scooter'::character varying, 'on_foot'::character varying, 'van'::character varying])::text[])))
);
CREATE TABLE public.custom_ui_component_usage (
    id integer NOT NULL,
    component_id integer NOT NULL,
    category_id integer,
    usage_context character varying(50) DEFAULT 'listing'::character varying NOT NULL,
    placement character varying(50) DEFAULT 'default'::character varying,
    priority integer DEFAULT 0,
    configuration jsonb DEFAULT '{}'::jsonb,
    conditions_logic jsonb,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer,
    updated_by integer
);
CREATE TABLE public.custom_ui_components (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    component_type character varying(50) NOT NULL,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by integer,
    updated_by integer,
    template_code text DEFAULT ''::text NOT NULL,
    styles text DEFAULT ''::text,
    props_schema jsonb DEFAULT '{}'::jsonb,
    CONSTRAINT custom_ui_components_component_type_check CHECK (((component_type)::text = ANY (ARRAY[('category'::character varying)::text, ('attribute'::character varying)::text, ('filter'::character varying)::text])))
);
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
    CONSTRAINT deliveries_package_size_check CHECK (((package_size)::text = ANY ((ARRAY['small'::character varying, 'medium'::character varying, 'large'::character varying, 'xl'::character varying])::text[]))),
    CONSTRAINT delivery_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'assigned'::character varying, 'accepted'::character varying, 'picked_up'::character varying, 'in_transit'::character varying, 'delivered'::character varying, 'cancelled'::character varying, 'failed'::character varying])::text[])))
);
CREATE TABLE public.delivery_notifications (
    id integer NOT NULL,
    delivery_id integer,
    viber_user_id integer,
    notification_type character varying(50),
    sent_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    is_read boolean DEFAULT false,
    read_at timestamp with time zone
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
    unit character varying(50)
);
CREATE TABLE public.listing_views (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    user_id integer,
    ip_hash character varying(255),
    view_time timestamp without time zone DEFAULT now(),
    CONSTRAINT at_least_one_identifier CHECK (((user_id IS NOT NULL) OR (ip_hash IS NOT NULL)))
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
CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(150) NOT NULL,
    google_id character varying(255),
    picture_url text,
    phone character varying(20),
    bio text,
    notification_email boolean DEFAULT true,
    timezone character varying(50) DEFAULT 'UTC'::character varying,
    last_seen timestamp without time zone,
    account_status character varying(20) DEFAULT 'active'::character varying,
    settings jsonb DEFAULT '{}'::jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    city character varying(100),
    country character varying(100),
    password character varying(255),
    provider character varying(50) DEFAULT 'email'::character varying,
    preferred_language character varying(10) DEFAULT 'ru'::character varying,
    role_id integer,
    old_email text,
    CONSTRAINT users_account_status_check CHECK (((account_status)::text = ANY (ARRAY[('active'::character varying)::text, ('inactive'::character varying)::text, ('suspended'::character varying)::text]))),
    CONSTRAINT users_preferred_language_check CHECK (((preferred_language)::text = ANY (ARRAY[('ru'::character varying)::text, ('sr'::character varying)::text, ('en'::character varying)::text])))
);
CREATE TABLE public.marketplace_favorites (
    user_id integer NOT NULL,
    listing_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.notification_settings (
    user_id integer NOT NULL,
    notification_type character varying(50) NOT NULL,
    telegram_enabled boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    email_enabled boolean DEFAULT false
);
CREATE TABLE public.user_telegram_connections (
    user_id integer NOT NULL,
    telegram_chat_id character varying(100) NOT NULL,
    telegram_username character varying(100),
    connected_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
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
    original_language character varying(2) DEFAULT 'en'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    has_attachments boolean DEFAULT false,
    attachments_count integer DEFAULT 0,
    storefront_product_id integer,
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
    weight_kg numeric(10,3) NOT NULL,
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
CREATE TABLE public.user_roles (
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    assigned_by integer
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
