CREATE EXTENSION IF NOT EXISTS cube WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS earthdistance WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE TYPE public.storefront_invitation_status AS ENUM (
    'pending',
    'accepted',
    'declined',
    'expired',
    'revoked'
);
CREATE SEQUENCE public.b2c_product_variants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TYPE public.storefront_invitation_type AS ENUM (
    'email',
    'link'
);
CREATE SEQUENCE public.attribute_options_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.attribute_search_cache_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.attribute_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.c2c_chats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.c2c_messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.cart_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.category_attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.category_variant_attributes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.chat_attachments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.chats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.indexing_queue_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.inventory_movements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.inventory_reservations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listing_attribute_values_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listing_attributes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listing_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listing_locations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listing_tags_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listing_variants_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.listings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.search_queries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.shopping_carts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_delivery_options_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_hours_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_invitations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_payment_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_staff_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefronts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE FUNCTION public.archive_old_analytics_events() RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_deleted_count INTEGER;
BEGIN
    -- Delete events older than 90 days (after materialized views are refreshed)
    DELETE FROM analytics_events
    WHERE created_at < NOW() - INTERVAL '90 days'
    RETURNING COUNT(*) INTO v_deleted_count;
    RETURN v_deleted_count;
END;
$$;
CREATE FUNCTION public.auto_expire_reservations() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Mark reservation as expired if past expiration time
    IF NEW.status = 'active' AND NEW.expires_at < CURRENT_TIMESTAMP THEN
        NEW.status = 'expired';
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.auto_update_variant_status() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Set status to out_of_stock when stock reaches zero
    IF NEW.stock_quantity = 0 AND OLD.stock_quantity > 0 THEN
        NEW.status = 'out_of_stock';
    END IF;
    -- Set status to active when stock becomes available again
    IF NEW.stock_quantity > 0 AND OLD.status = 'out_of_stock' THEN
        NEW.status = 'active';
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.is_valid_file_size(file_type text, file_size bigint) RETURNS boolean
    LANGUAGE plpgsql IMMUTABLE
    AS $$
BEGIN
    RETURN (
        (file_type = 'image' AND file_size > 0 AND file_size <= 10485760) OR
        (file_type = 'video' AND file_size > 0 AND file_size <= 52428800) OR
        (file_type = 'document' AND file_size > 0 AND file_size <= 20971520)
    );
END;
$$;
CREATE FUNCTION public.log_analytics_event(p_event_type character varying, p_entity_type character varying, p_entity_id bigint, p_user_id bigint DEFAULT NULL::bigint, p_session_id character varying DEFAULT NULL::character varying, p_metadata jsonb DEFAULT '{}'::jsonb) RETURNS uuid
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_uuid UUID;
BEGIN
    INSERT INTO analytics_events (
        event_type, entity_type, entity_id,
        user_id, session_id, metadata
    )
    VALUES (
        p_event_type, p_entity_type, p_entity_id,
        p_user_id, p_session_id, p_metadata
    )
    RETURNING uuid INTO v_uuid;
    RETURN v_uuid;
END;
$$;
CREATE FUNCTION public.refresh_analytics_trending_cache() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_trending_cache;
END;
$$;
CREATE FUNCTION public.refresh_analytics_views() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_overview_daily;
    REFRESH MATERIALIZED VIEW CONCURRENTLY analytics_listing_stats;
END;
$$;
CREATE FUNCTION public.update_ai_mapping_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_attribute_values_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_attributes_search_vector() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.name->>'en', '')), 'A') ||
        setweight(to_tsvector('russian', COALESCE(NEW.name->>'ru', '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.code, '')), 'B');
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_attributes_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_b2c_product_variants_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    -- Automatically update parent product's has_variants flag
    UPDATE b2c_products
    SET has_variants = EXISTS(
        SELECT 1 FROM b2c_product_variants
        WHERE product_id = NEW.product_id
    )
    WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_b2c_products_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_cart_items_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_chats_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_inventory_reservations_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_listing_variants_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_messages_read_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.is_read = true AND OLD.is_read = false THEN
        NEW.read_at = CURRENT_TIMESTAMP;
        NEW.status = 'read';
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_messages_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_orders_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_shopping_carts_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_storefront_delivery_options_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_storefront_invitations_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_storefront_staff_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_storefronts_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.check_category_level_constraint() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- L1 categories (level=1) must have parent_id = NULL
    IF NEW.level = 1 AND NEW.parent_id IS NOT NULL THEN
        RAISE EXCEPTION 'Level 1 categories must have parent_id = NULL';
    END IF;
    -- L2/L3 categories must have parent_id
    IF NEW.level > 1 AND NEW.parent_id IS NULL THEN
        RAISE EXCEPTION 'Level % categories must have a parent_id', NEW.level;
    END IF;
    -- Verify parent's level is exactly (current level - 1)
    IF NEW.parent_id IS NOT NULL THEN
        DECLARE
            parent_level INTEGER;
        BEGIN
            SELECT level INTO parent_level FROM categories WHERE id = NEW.parent_id;
            IF parent_level IS NULL THEN
                RAISE EXCEPTION 'Parent category not found';
            END IF;
            IF parent_level != NEW.level - 1 THEN
                RAISE EXCEPTION 'Parent level must be % for level % category', NEW.level - 1, NEW.level;
            END IF;
        END;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.cleanup_expired_reservations() RETURNS TABLE(cleaned_count bigint)
    LANGUAGE plpgsql
    AS $$
DECLARE
    result BIGINT;
BEGIN
    -- Update expired reservations
    WITH expired AS (
        UPDATE stock_reservations
        SET status = 'expired'
        WHERE status = 'active' AND expires_at < CURRENT_TIMESTAMP
        RETURNING id
    )
    SELECT COUNT(*) INTO result FROM expired;
    RETURN QUERY SELECT result;
END;
$$;
CREATE FUNCTION public.enforce_single_default_variant() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.is_default = true THEN
        -- Unset is_default for all other variants of the same product
        UPDATE product_variants
        SET is_default = false
        WHERE product_id = NEW.product_id AND id != NEW.id;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.generate_slug_from_title() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.slug IS NULL OR NEW.slug = '' THEN
        NEW.slug := LOWER(
            TRIM(
                BOTH '-' FROM
                REGEXP_REPLACE(
                    REGEXP_REPLACE(
                        REGEXP_REPLACE(NEW.title, '[^a-zA-Z0-9\s-]', '', 'g'),
                        '\s+', '-', 'g'
                    ),
                    '-+', '-', 'g'
                )
            )
        );
        -- Handle potential duplicates by appending random suffix
        IF EXISTS (SELECT 1 FROM listings WHERE slug = NEW.slug AND id != COALESCE(NEW.id, 0) AND is_deleted = false) THEN
            NEW.slug := NEW.slug || '-' || COALESCE(NEW.id, (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT);
        END IF;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.get_category_attributes_with_inheritance(p_category_id uuid) RETURNS TABLE(attribute_id integer, category_id uuid, is_enabled boolean, is_required boolean, is_searchable boolean, is_filterable boolean, sort_order integer)
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE category_path AS (
        SELECT
            c.id,
            c.parent_id,
            c.level,
            1 as distance
        FROM categories c
        WHERE c.id = p_category_id
        UNION ALL
        SELECT
            c.id,
            c.parent_id,
            c.level,
            cp.distance + 1
        FROM categories c
        INNER JOIN category_path cp ON c.id = cp.parent_id
    )
    SELECT DISTINCT ON (a.id)
        a.id as attribute_id,
        ca.category_id,
        ca.is_enabled,
        COALESCE(ca.is_required, false) as is_required,
        COALESCE(ca.is_searchable, a.is_searchable) as is_searchable,
        COALESCE(ca.is_filterable, a.is_filterable) as is_filterable,
        COALESCE(ca.sort_order, a.sort_order) as sort_order
    FROM category_path cp
    INNER JOIN category_attributes ca ON cp.id = ca.category_id
    INNER JOIN attributes a ON ca.attribute_id = a.id
    WHERE ca.is_enabled = true
      AND a.is_active = true
    ORDER BY a.id, cp.distance ASC;
END;
$$;
CREATE FUNCTION public.get_category_attributes_with_inheritance(p_category_id integer, p_locale character varying DEFAULT 'sr'::character varying) RETURNS TABLE(attribute_id integer, attribute_code character varying, attribute_name text, attribute_type character varying, purpose character varying, is_required boolean, is_filterable boolean, is_searchable boolean, sort_order integer, options jsonb, validation_rules jsonb, ui_settings jsonb, source_category_id integer, is_inherited boolean)
    LANGUAGE plpgsql STABLE
    AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE category_path AS (
        -- Start with the given category
        SELECT
            c.id,
            c.parent_id,
            c.level,
            1 as distance
        FROM categories c
        WHERE c.id = p_category_id
        UNION ALL
        -- Recursively get parent categories
        SELECT
            c.id,
            c.parent_id,
            c.level,
            cp.distance + 1
        FROM categories c
        INNER JOIN category_path cp ON c.id = cp.parent_id
    ),
    all_attributes AS (
        -- Get category-specific attributes (including inherited)
        SELECT DISTINCT ON (a.id)
            a.id as attribute_id,
            a.code as attribute_code,
            a.name->>p_locale as attribute_name,
            a.attribute_type,
            a.purpose,
            COALESCE(ca.is_required, a.is_required) as is_required,
            COALESCE(ca.is_filterable, a.is_filterable) as is_filterable,
            COALESCE(ca.is_searchable, a.is_searchable) as is_searchable,
            COALESCE(ca.sort_order, a.sort_order) as sort_order,
            COALESCE(ca.category_options, a.options) as options,
            COALESCE(ca.custom_validation, a.validation_rules) as validation_rules,
            COALESCE(ca.custom_ui_settings, a.ui_settings) as ui_settings,
            ca.category_id as source_category_id,
            (ca.category_id != p_category_id) as is_inherited,
            cp.distance as priority
        FROM category_path cp
        INNER JOIN category_attributes ca ON cp.id = ca.category_id
        INNER JOIN attributes a ON ca.attribute_id = a.id
        WHERE ca.is_enabled = true
          AND a.is_active = true
        ORDER BY a.id, priority ASC
    )
    SELECT
        aa.attribute_id,
        aa.attribute_code,
        aa.attribute_name,
        aa.attribute_type,
        aa.purpose,
        aa.is_required,
        aa.is_filterable,
        aa.is_searchable,
        aa.sort_order,
        aa.options,
        aa.validation_rules,
        aa.ui_settings,
        aa.source_category_id,
        aa.is_inherited
    FROM all_attributes aa
    ORDER BY aa.sort_order, aa.attribute_code;
END;
$$;
CREATE FUNCTION public.get_file_type_from_content_type(content_type text) RETURNS character varying
    LANGUAGE plpgsql IMMUTABLE
    AS $$
BEGIN
    CASE
        WHEN content_type LIKE 'image/%' THEN RETURN 'image';
        WHEN content_type LIKE 'video/%' THEN RETURN 'video';
        ELSE RETURN 'document';
    END CASE;
END;
$$;
CREATE FUNCTION public.sync_variant_reserved_quantity() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    old_quantity INTEGER := 0;
    new_quantity INTEGER := 0;
BEGIN
    -- Handle INSERT
    IF (TG_OP = 'INSERT') THEN
        IF NEW.status = 'active' THEN
            UPDATE product_variants
            SET reserved_quantity = reserved_quantity + NEW.quantity
            WHERE id = NEW.variant_id;
        END IF;
        RETURN NEW;
    END IF;
    -- Handle UPDATE
    IF (TG_OP = 'UPDATE') THEN
        -- Calculate quantity delta
        IF OLD.status = 'active' THEN
            old_quantity := OLD.quantity;
        END IF;
        IF NEW.status = 'active' THEN
            new_quantity := NEW.quantity;
        END IF;
        -- Update variant reserved_quantity
        IF old_quantity != new_quantity THEN
            UPDATE product_variants
            SET reserved_quantity = reserved_quantity - old_quantity + new_quantity
            WHERE id = NEW.variant_id;
        END IF;
        RETURN NEW;
    END IF;
    -- Handle DELETE
    IF (TG_OP = 'DELETE') THEN
        IF OLD.status = 'active' THEN
            UPDATE product_variants
            SET reserved_quantity = reserved_quantity - OLD.quantity
            WHERE id = OLD.variant_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.update_chat_last_message_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE chats
    SET last_message_at = NEW.created_at,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.chat_id;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_message_attachments_count() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    attachment_count INT;
BEGIN
    -- Count attachments for the message
    SELECT COUNT(*) INTO attachment_count
    FROM chat_attachments
    WHERE message_id = COALESCE(NEW.message_id, OLD.message_id);
    -- Update message
    UPDATE messages
    SET
        attachments_count = attachment_count,
        has_attachments = (attachment_count > 0),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = COALESCE(NEW.message_id, OLD.message_id);
    RETURN COALESCE(NEW, OLD);
END;
$$;
CREATE FUNCTION public.validate_variant_attribute_value() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    attr_type VARCHAR(50);
BEGIN
    -- Get attribute type
    SELECT attribute_type INTO attr_type
    FROM attributes
    WHERE id = NEW.attribute_id;
    -- Validate correct value field is used
    CASE attr_type
        WHEN 'text', 'textarea', 'select', 'color', 'size' THEN
            IF NEW.value_text IS NULL THEN
                RAISE EXCEPTION 'value_text must be set for attribute type %', attr_type;
            END IF;
        WHEN 'number' THEN
            IF NEW.value_number IS NULL THEN
                RAISE EXCEPTION 'value_number must be set for attribute type %', attr_type;
            END IF;
        WHEN 'boolean' THEN
            IF NEW.value_boolean IS NULL THEN
                RAISE EXCEPTION 'value_boolean must be set for attribute type %', attr_type;
            END IF;
        WHEN 'date' THEN
            IF NEW.value_date IS NULL THEN
                RAISE EXCEPTION 'value_date must be set for attribute type %', attr_type;
            END IF;
        WHEN 'multiselect' THEN
            IF NEW.value_json IS NULL THEN
                RAISE EXCEPTION 'value_json must be set for attribute type %', attr_type;
            END IF;
    END CASE;
    RETURN NEW;
END;
$$;
CREATE TABLE public.brand_category_mapping (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    brand_name character varying(100) NOT NULL,
    brand_aliases text[] DEFAULT '{}'::text[],
    category_slug character varying(100) NOT NULL,
    confidence double precision DEFAULT 0.95 NOT NULL,
    is_verified boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
CREATE TABLE public.product_variants (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    product_id uuid NOT NULL,
    sku character varying(100) NOT NULL,
    price numeric(12,2),
    compare_at_price numeric(12,2),
    stock_quantity integer DEFAULT 0 NOT NULL,
    reserved_quantity integer DEFAULT 0 NOT NULL,
    low_stock_alert integer DEFAULT 5,
    weight_grams numeric(8,3),
    barcode character varying(50),
    is_default boolean DEFAULT false NOT NULL,
    "position" integer DEFAULT 0 NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT positive_reserved CHECK ((reserved_quantity >= 0)),
    CONSTRAINT positive_stock CHECK ((stock_quantity >= 0)),
    CONSTRAINT product_variants_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'out_of_stock'::character varying, 'discontinued'::character varying])::text[]))),
    CONSTRAINT reserved_not_exceed_stock CHECK ((reserved_quantity <= stock_quantity)),
    CONSTRAINT valid_compare_price CHECK (((compare_at_price IS NULL) OR (compare_at_price >= (0)::numeric))),
    CONSTRAINT valid_price CHECK (((price IS NULL) OR (price >= (0)::numeric)))
);
CREATE TABLE public.stock_reservations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    variant_id uuid NOT NULL,
    order_id uuid NOT NULL,
    quantity integer NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT positive_quantity CHECK ((quantity > 0)),
    CONSTRAINT stock_reservations_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT stock_reservations_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'confirmed'::character varying, 'cancelled'::character varying, 'expired'::character varying])::text[])))
);
CREATE TABLE public.listings (
    id bigint NOT NULL,
    uuid uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id bigint NOT NULL,
    storefront_id bigint,
    title character varying(255) NOT NULL,
    description text,
    price numeric(15,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying NOT NULL,
    status character varying(50) DEFAULT 'draft'::character varying NOT NULL,
    visibility character varying(50) DEFAULT 'public'::character varying NOT NULL,
    quantity integer DEFAULT 1 NOT NULL,
    sku character varying(100),
    view_count integer DEFAULT 0 NOT NULL,
    favorites_count integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    published_at timestamp with time zone,
    deleted_at timestamp with time zone,
    is_deleted boolean DEFAULT false NOT NULL,
    source_type character varying(10) DEFAULT 'c2c'::character varying NOT NULL,
    attributes jsonb DEFAULT '{}'::jsonb,
    stock_status character varying(50) DEFAULT 'in_stock'::character varying,
    sold_count integer DEFAULT 0 NOT NULL,
    has_individual_location boolean DEFAULT false,
    individual_address text,
    individual_latitude numeric(10,8),
    individual_longitude numeric(11,8),
    location_privacy character varying(20) DEFAULT 'exact'::character varying,
    show_on_map boolean DEFAULT true,
    has_variants boolean DEFAULT false,
    slug character varying(255),
    expires_at timestamp with time zone,
    title_translations jsonb DEFAULT '{}'::jsonb,
    description_translations jsonb DEFAULT '{}'::jsonb,
    location_translations jsonb DEFAULT '{}'::jsonb,
    city_translations jsonb DEFAULT '{}'::jsonb,
    country_translations jsonb DEFAULT '{}'::jsonb,
    original_language character varying(10) DEFAULT 'sr'::character varying,
    category_id uuid,
    CONSTRAINT chk_original_language CHECK (((original_language)::text = ANY ((ARRAY['sr'::character varying, 'en'::character varying, 'ru'::character varying])::text[]))),
    CONSTRAINT listings_location_privacy_check CHECK (((location_privacy)::text = ANY ((ARRAY['exact'::character varying, 'approximate'::character varying, 'hidden'::character varying])::text[]))),
    CONSTRAINT listings_price_check CHECK ((price > (0)::numeric)),
    CONSTRAINT listings_quantity_check CHECK ((quantity >= 0)),
    CONSTRAINT listings_source_type_check CHECK (((source_type)::text = ANY ((ARRAY['c2c'::character varying, 'b2c'::character varying])::text[]))),
    CONSTRAINT listings_status_check CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'active'::character varying, 'inactive'::character varying, 'sold'::character varying, 'archived'::character varying])::text[]))),
    CONSTRAINT listings_stock_status_check CHECK (((stock_status)::text = ANY ((ARRAY['in_stock'::character varying, 'out_of_stock'::character varying, 'low_stock'::character varying, 'discontinued'::character varying])::text[]))),
    CONSTRAINT listings_title_check CHECK ((length(TRIM(BOTH FROM title)) >= 3)),
    CONSTRAINT listings_user_id_check CHECK ((user_id > 0)),
    CONSTRAINT listings_visibility_check CHECK (((visibility)::text = ANY ((ARRAY['public'::character varying, 'private'::character varying, 'unlisted'::character varying])::text[])))
);
CREATE TABLE public.listing_favorites (
    user_id integer NOT NULL,
    listing_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.listing_stats (
    listing_id bigint NOT NULL,
    views_count integer DEFAULT 0 NOT NULL,
    favorites_count integer DEFAULT 0 NOT NULL,
    inquiries_count integer DEFAULT 0 NOT NULL,
    last_viewed_at timestamp with time zone,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.orders (
    id bigint NOT NULL,
    order_number character varying(50) NOT NULL,
    user_id bigint,
    storefront_id bigint NOT NULL,
    status character varying(20) NOT NULL,
    payment_status character varying(20) NOT NULL,
    payment_method character varying(50),
    payment_transaction_id character varying(255),
    subtotal numeric(10,2) NOT NULL,
    tax numeric(10,2) DEFAULT 0.00,
    shipping numeric(10,2) DEFAULT 0.00,
    discount numeric(10,2) DEFAULT 0.00,
    total numeric(10,2) NOT NULL,
    commission numeric(10,2) DEFAULT 0.00,
    seller_amount numeric(10,2) NOT NULL,
    currency character varying(3) DEFAULT 'RSD'::character varying,
    shipping_address jsonb,
    billing_address jsonb,
    shipping_provider character varying(100),
    tracking_number character varying(255),
    escrow_release_date timestamp with time zone,
    escrow_days integer DEFAULT 3,
    shipment_id bigint,
    cancelled_at timestamp with time zone,
    cancellation_reason text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    payment_completed_at timestamp with time zone,
    customer_name character varying(255),
    customer_email character varying(255),
    customer_phone character varying(50),
    customer_notes text,
    admin_notes text,
    confirmed_at timestamp with time zone,
    shipped_at timestamp with time zone,
    delivered_at timestamp with time zone,
    seller_notes text,
    accepted_at timestamp with time zone,
    label_url text,
    shipping_cost_cents integer DEFAULT 0,
    shipping_method_id character varying(100),
    shipping_details jsonb DEFAULT '{}'::jsonb,
    delivery_address_id integer,
    delivery_address_snapshot jsonb,
    payment_provider character varying(50),
    payment_session_id character varying(255),
    payment_intent_id character varying(255),
    payment_idempotency_key character varying(255),
    shipping_method character varying(100),
    notes text,
    CONSTRAINT chk_orders_escrow_days_positive CHECK ((escrow_days >= 0)),
    CONSTRAINT chk_orders_payment_status CHECK (((payment_status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying, 'completed'::character varying, 'cod_pending'::character varying, 'failed'::character varying, 'refunded'::character varying, 'partially_refunded'::character varying])::text[]))),
    CONSTRAINT chk_orders_seller_amount_positive CHECK ((seller_amount >= (0)::numeric)),
    CONSTRAINT chk_orders_status CHECK (((status)::text = ANY ((ARRAY['unspecified'::character varying, 'pending'::character varying, 'confirmed'::character varying, 'accepted'::character varying, 'processing'::character varying, 'shipped'::character varying, 'delivered'::character varying, 'cancelled'::character varying, 'refunded'::character varying, 'failed'::character varying])::text[]))),
    CONSTRAINT chk_orders_subtotal_positive CHECK ((subtotal >= (0)::numeric)),
    CONSTRAINT chk_orders_total_positive CHECK ((total >= (0)::numeric)),
    CONSTRAINT valid_shipping_cost_cents CHECK (((shipping_cost_cents IS NULL) OR (shipping_cost_cents >= 0)))
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
    formatted_address text,
    geo_strategy character varying(50) DEFAULT 'storefront_location'::character varying,
    default_privacy_level character varying(20) DEFAULT 'exact'::character varying,
    address_verified boolean DEFAULT false,
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
    followers_count integer DEFAULT 0,
    subscription_plan character varying(50) DEFAULT 'starter'::character varying,
    subscription_expires_at timestamp without time zone,
    subscription_id integer,
    is_subscription_active boolean DEFAULT true,
    commission_rate numeric(5,2) DEFAULT 3.00,
    ai_agent_enabled boolean DEFAULT false,
    ai_agent_config jsonb DEFAULT '{}'::jsonb,
    live_shopping_enabled boolean DEFAULT false,
    group_buying_enabled boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);
CREATE TABLE public.attribute_options (
    id integer NOT NULL,
    attribute_id integer NOT NULL,
    option_value character varying(255) NOT NULL,
    option_label jsonb DEFAULT '{}'::jsonb NOT NULL,
    color_hex character varying(7),
    image_url text,
    icon character varying(100),
    is_default boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.attribute_search_cache (
    id integer NOT NULL,
    listing_id integer NOT NULL,
    attributes_flat jsonb DEFAULT '{}'::jsonb NOT NULL,
    attributes_searchable text,
    attributes_filterable jsonb DEFAULT '{}'::jsonb NOT NULL,
    last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    cache_version integer DEFAULT 1 NOT NULL
);
CREATE TABLE public.attribute_values (
    id integer NOT NULL,
    attribute_id integer NOT NULL,
    value character varying(255) NOT NULL,
    label jsonb DEFAULT '{}'::jsonb NOT NULL,
    metadata jsonb DEFAULT '{}'::jsonb,
    sort_order integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.attributes (
    id integer NOT NULL,
    code character varying(100) NOT NULL,
    name jsonb DEFAULT '{}'::jsonb NOT NULL,
    display_name jsonb DEFAULT '{}'::jsonb NOT NULL,
    attribute_type character varying(50) NOT NULL,
    purpose character varying(20) DEFAULT 'regular'::character varying NOT NULL,
    options jsonb DEFAULT '{}'::jsonb,
    validation_rules jsonb DEFAULT '{}'::jsonb,
    ui_settings jsonb DEFAULT '{}'::jsonb,
    is_searchable boolean DEFAULT false NOT NULL,
    is_filterable boolean DEFAULT false NOT NULL,
    is_required boolean DEFAULT false NOT NULL,
    is_variant_compatible boolean DEFAULT false NOT NULL,
    affects_stock boolean DEFAULT false NOT NULL,
    affects_price boolean DEFAULT false NOT NULL,
    show_in_card boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    legacy_category_attribute_id integer,
    icon character varying(255) DEFAULT ''::character varying,
    search_vector tsvector,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT attributes_attribute_type_check CHECK (((attribute_type)::text = ANY ((ARRAY['text'::character varying, 'textarea'::character varying, 'number'::character varying, 'boolean'::character varying, 'select'::character varying, 'multiselect'::character varying, 'date'::character varying, 'color'::character varying, 'size'::character varying])::text[]))),
    CONSTRAINT attributes_purpose_check CHECK (((purpose)::text = ANY ((ARRAY['regular'::character varying, 'variant'::character varying, 'both'::character varying])::text[])))
);
CREATE TABLE public.variant_attribute_values (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    variant_id uuid NOT NULL,
    attribute_id integer NOT NULL,
    value_text text,
    value_number numeric(20,4),
    value_boolean boolean,
    value_date date,
    value_json jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE public.c2c_chats (
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
CREATE TABLE public.c2c_messages (
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
CREATE TABLE public.cart_items (
    id bigint NOT NULL,
    cart_id bigint NOT NULL,
    listing_id bigint NOT NULL,
    variant_id bigint,
    quantity integer NOT NULL,
    price_snapshot numeric(10,2) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_cart_items_quantity_positive CHECK ((quantity > 0))
);
CREATE TABLE public.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    slug character varying(255) NOT NULL,
    parent_id uuid,
    level integer NOT NULL,
    path character varying(1000) NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    name jsonb DEFAULT '{}'::jsonb NOT NULL,
    description jsonb DEFAULT '{}'::jsonb,
    meta_title jsonb DEFAULT '{}'::jsonb,
    meta_description jsonb DEFAULT '{}'::jsonb,
    meta_keywords jsonb DEFAULT '{}'::jsonb,
    icon character varying(50),
    image_url character varying(500),
    is_active boolean DEFAULT true NOT NULL,
    google_category_id integer,
    facebook_category_id character varying(100),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    has_custom_ui boolean DEFAULT false,
    custom_ui_component character varying(100),
    CONSTRAINT categories_level_check CHECK (((level >= 1) AND (level <= 3)))
);
CREATE TABLE public.category_ai_mapping (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ai_category_name character varying(100) NOT NULL,
    target_category_id uuid NOT NULL,
    confidence_boost numeric(3,2) DEFAULT 0.15,
    priority integer DEFAULT 100,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    notes text,
    CONSTRAINT category_ai_mapping_confidence_boost_check CHECK (((confidence_boost >= (0)::numeric) AND (confidence_boost <= 1.0))),
    CONSTRAINT category_ai_mapping_priority_check CHECK ((priority > 0))
);
CREATE TABLE public.category_detections (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    input_title text NOT NULL,
    input_description text,
    input_language character varying(5) DEFAULT 'sr'::character varying NOT NULL,
    detected_category_id uuid,
    confidence_score numeric(4,3),
    detection_method character varying(50) NOT NULL,
    matched_keywords text[],
    alternatives jsonb DEFAULT '[]'::jsonb,
    user_confirmed boolean,
    user_selected_category_id uuid,
    processing_time_ms integer,
    created_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.category_attributes (
    id integer NOT NULL,
    category_id uuid NOT NULL,
    attribute_id integer NOT NULL,
    is_enabled boolean DEFAULT true,
    is_required boolean,
    sort_order integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_searchable boolean DEFAULT false,
    is_filterable boolean DEFAULT true,
    category_options jsonb,
    custom_validation jsonb,
    custom_ui_settings jsonb
);
CREATE TABLE public.category_variant_attributes (
    id integer NOT NULL,
    category_id character varying(36) NOT NULL,
    attribute_id integer NOT NULL,
    is_required boolean DEFAULT false NOT NULL,
    affects_price boolean DEFAULT false NOT NULL,
    affects_stock boolean DEFAULT true NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    display_as character varying(50) DEFAULT 'dropdown'::character varying,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT category_variant_attributes_display_as_check CHECK (((display_as)::text = ANY ((ARRAY['dropdown'::character varying, 'buttons'::character varying, 'swatches'::character varying, 'radio'::character varying])::text[])))
);
CREATE TABLE public.chat_attachments (
    id bigint NOT NULL,
    message_id bigint NOT NULL,
    file_type character varying(20) NOT NULL,
    file_name character varying(255) NOT NULL,
    file_size bigint NOT NULL,
    content_type character varying(100) NOT NULL,
    storage_type character varying(20) DEFAULT 'minio'::character varying NOT NULL,
    storage_bucket character varying(100) DEFAULT 'chat-files'::character varying NOT NULL,
    file_path character varying(500) NOT NULL,
    public_url text,
    thumbnail_url text,
    metadata jsonb,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT check_content_type CHECK (((content_type)::text ~ '^[a-z]+/[a-z0-9\+\-\.]+$'::text)),
    CONSTRAINT check_file_name CHECK (((length(TRIM(BOTH FROM file_name)) > 0) AND (length((file_name)::text) <= 255))),
    CONSTRAINT check_file_path CHECK (((length(TRIM(BOTH FROM file_path)) > 0) AND (length((file_path)::text) <= 500))),
    CONSTRAINT check_file_size CHECK (((((file_type)::text = 'image'::text) AND (file_size > 0) AND (file_size <= 10485760)) OR (((file_type)::text = 'video'::text) AND (file_size > 0) AND (file_size <= 52428800)) OR (((file_type)::text = 'document'::text) AND (file_size > 0) AND (file_size <= 20971520)))),
    CONSTRAINT check_file_type CHECK (((file_type)::text = ANY ((ARRAY['image'::character varying, 'video'::character varying, 'document'::character varying])::text[]))),
    CONSTRAINT check_storage_type CHECK (((storage_type)::text = ANY ((ARRAY['minio'::character varying, 's3'::character varying, 'local'::character varying])::text[])))
);
