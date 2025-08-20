CREATE SEQUENCE public.storefront_product_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_product_variant_images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.storefront_product_variants_id_seq
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
CREATE SEQUENCE public.translation_audit_log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.translation_providers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.translation_quality_metrics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.translation_sync_conflicts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.translation_tasks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.translations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.transliteration_rules_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.unified_geo_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.user_behavior_events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.user_contacts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.user_storefronts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.variant_attribute_mappings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.warehouse_inventory_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.warehouse_invoices_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.warehouse_movements_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.warehouse_pickup_orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.warehouses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE SEQUENCE public.webhook_audit_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE FUNCTION public.calculate_escrow_release_date() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Рассчитываем дату освобождения на основе escrow_days
    IF NEW.escrow_release_date IS NULL AND NEW.escrow_days IS NOT NULL THEN
        NEW.escrow_release_date := CURRENT_DATE + INTERVAL '1 day' * NEW.escrow_days;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.generate_pickup_code() RETURNS character varying
    LANGUAGE plpgsql
    AS $$
DECLARE
    chars VARCHAR := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR := '';
    i INTEGER;
BEGIN
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$;
CREATE FUNCTION public.log_shipment_status_change() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        NEW.status_history = NEW.status_history || jsonb_build_object(
            'status', NEW.status,
            'timestamp', NOW(),
            'old_status', OLD.status
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.rebuild_all_ratings() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
END;
$$;
CREATE FUNCTION public.refresh_category_listing_counts() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    current_ts TIMESTAMP;
    last_refresh TIMESTAMP;
BEGIN
    -- Проверяем, не обновлялось ли представление в последние N секунд
    SELECT INTO last_refresh COALESCE(
        (SELECT obj_description('category_listing_counts'::regclass)::timestamp),
        '1970-01-01'::timestamp
    );
    current_ts := CURRENT_TIMESTAMP;
    IF current_ts - last_refresh > interval '5 seconds' THEN
        -- Обновляем представление
        REFRESH MATERIALIZED VIEW category_listing_counts;
        -- Сохраняем время последнего обновления
        EXECUTE format(
            'COMMENT ON MATERIALIZED VIEW category_listing_counts IS %L',
            current_ts::text
        );
    END IF;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.refresh_density_grid() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY gis_listing_density_grid;
END;
$$;
CREATE FUNCTION public.refresh_map_items_cache() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY map_items_cache;
END;
$$;
CREATE FUNCTION public.refresh_rating_distributions() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Обновляем все материализованные представления
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_distribution;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_distribution;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.refresh_rating_summaries() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY user_rating_summary;
    REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_rating_summary;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.refresh_rating_views() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Обновляем только затронутые строки, а не всё представление
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- Для пользователей
        IF NEW.entity_origin_type = 'user' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;
        -- Для магазинов
        IF NEW.entity_origin_type = 'storefront' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении также обновляем
        IF OLD.entity_origin_type = 'user' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;
        IF OLD.entity_origin_type = 'storefront' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    END IF;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.set_order_number() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.order_number IS NULL OR NEW.order_number = '' THEN
        NEW.order_number := generate_order_number();
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.trigger_cleanup_geocoding_cache() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Очищаем кэш только в 1% случаев для избежания частых очисток
    IF random() < 0.01 THEN
        PERFORM cleanup_expired_geocoding_cache();
    END IF;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.trigger_refresh_map_cache() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Use pg_notify to trigger async refresh
    PERFORM pg_notify('refresh_map_cache', '');
    RETURN COALESCE(NEW, OLD);
END;
$$;
CREATE FUNCTION public.update_attribute_groups_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_bex_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_car_trims_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_custom_ui_component_usage_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_failed_webhooks_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_geocoding_cache_updated_at() RETURNS trigger
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
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_item_performance_metrics_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_listing_search_vector() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('russian_unaccent', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('russian_unaccent', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english_unaccent', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english_unaccent', COALESCE(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_listings_geo_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_marketplace_chats_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_marketplace_orders_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_notification_settings_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_post_express_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_product_stock_status() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.stock_quantity = 0 THEN
        NEW.stock_status = 'out_of_stock';
    ELSIF NEW.stock_quantity <= 5 THEN
        NEW.stock_status = 'low_stock';
    ELSE
        NEW.stock_status = 'in_stock';
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_product_variants_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_search_behavior_metrics_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_search_optimization_sessions_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_search_queries_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_search_synonyms_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_search_weights_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_shopping_cart_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_translation_providers_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_translations_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_transliteration_rules_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_unified_geo_updated_at() RETURNS trigger
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
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_user_contacts_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_user_privacy_settings_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_user_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_warehouse_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.assign_district_municipality() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Find district containing the point
    SELECT id INTO NEW.district_id
    FROM districts
    WHERE ST_Contains(boundary, NEW.location::geometry)
    LIMIT 1;
    -- Find municipality containing the point
    SELECT id INTO NEW.municipality_id
    FROM municipalities
    WHERE ST_Contains(boundary, NEW.location::geometry)
    LIMIT 1;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.auto_geocode_storefront_product() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    storefront_rec RECORD;
    product_location geography(Point, 4326);
    original_location geography(Point, 4326);
    calculated_geohash VARCHAR(12);
    privacy_level location_privacy_level;
BEGIN
    -- Get storefront data
    SELECT s.geo_strategy, s.address, s.latitude, s.longitude, s.default_privacy_level
    INTO storefront_rec
    FROM storefronts s
    WHERE s.id = NEW.storefront_id;
    -- Determine location and privacy level
    IF NEW.has_individual_location AND NEW.individual_latitude IS NOT NULL AND NEW.individual_longitude IS NOT NULL THEN
        -- Product has individual location
        original_location := ST_SetSRID(ST_MakePoint(NEW.individual_longitude, NEW.individual_latitude), 4326)::geography;
        privacy_level := COALESCE(NEW.location_privacy, storefront_rec.default_privacy_level);
    ELSIF storefront_rec.latitude IS NOT NULL AND storefront_rec.longitude IS NOT NULL THEN
        -- Use storefront location
        original_location := ST_SetSRID(ST_MakePoint(storefront_rec.longitude, storefront_rec.latitude), 4326)::geography;
        privacy_level := 'exact'; -- Storefront location is always exact
    ELSE
        -- No location available
        RETURN NEW;
    END IF;
    -- Calculate display location based on privacy level
    -- Исправлено: передаем geography напрямую вместо numeric параметров
    product_location := calculate_blurred_location(
        original_location,
        privacy_level
    );
    -- Skip if location should be hidden
    IF product_location IS NULL OR NOT NEW.show_on_map THEN
        -- Remove from geo table if exists
        DELETE FROM unified_geo
        WHERE source_type = 'storefront_product' AND source_id = NEW.id;
        RETURN NEW;
    END IF;
    -- Calculate geohash
    calculated_geohash := substring(md5(ST_Y(product_location::geometry)::text || ST_X(product_location::geometry)::text), 1, 12);
    -- Insert or update geo data
    INSERT INTO unified_geo (
        source_type, source_id, location, original_location, geohash,
        privacy_level, blur_radius_meters,
        formatted_address, created_at, updated_at
    ) VALUES (
        'storefront_product',
        NEW.id,
        product_location,
        original_location,
        calculated_geohash,
        privacy_level,
        CASE privacy_level
            WHEN 'street' THEN 200  -- Обновлено с 'approximate'
            WHEN 'district' THEN 1000  -- Обновлено с 'city_only'
            WHEN 'city' THEN 5000  -- Новое значение
            ELSE 0
        END,
        COALESCE(NEW.individual_address, storefront_rec.address),
        NOW(),
        NOW()
    ) ON CONFLICT (source_type, source_id) DO UPDATE SET
        location = EXCLUDED.location,
        original_location = EXCLUDED.original_location,
        geohash = EXCLUDED.geohash,
        privacy_level = EXCLUDED.privacy_level,
        blur_radius_meters = EXCLUDED.blur_radius_meters,
        formatted_address = EXCLUDED.formatted_address,
        updated_at = EXCLUDED.updated_at;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.calculate_blurred_location(exact_location public.geography, privacy_level public.location_privacy_level) RETURNS public.geography
    LANGUAGE plpgsql IMMUTABLE
    AS $$
BEGIN
    CASE privacy_level
        WHEN 'exact' THEN
            RETURN exact_location;
        WHEN 'street' THEN
            -- Смещаем точку на случайное расстояние в пределах 200м
            RETURN ST_Project(
                exact_location,
                200 + (random() * 100 - 50), -- 200м ± 50м
                radians(random() * 360)       -- случайное направление
            )::geography;
        WHEN 'district' THEN
            -- Смещаем точку на случайное расстояние в пределах 1км
            RETURN ST_Project(
                exact_location,
                1000 + (random() * 200 - 100), -- 1км ± 100м
                radians(random() * 360)        -- случайное направление
            )::geography;
        WHEN 'city' THEN
            -- Смещаем точку на случайное расстояние в пределах 5км
            RETURN ST_Project(
                exact_location,
                5000 + (random() * 1000 - 500), -- 5км ± 500м
                radians(random() * 360)         -- случайное направление
            )::geography;
        ELSE
            RETURN NULL;
    END CASE;
END;
$$;
CREATE FUNCTION public.calculate_entity_rating(p_entity_type character varying, p_entity_id integer) RETURNS numeric
    LANGUAGE plpgsql
    AS $$
DECLARE
    avg_rating NUMERIC;
BEGIN
    SELECT COALESCE(AVG(rating)::NUMERIC(3,2), 0)
    INTO avg_rating
    FROM reviews
    WHERE entity_type = p_entity_type
    AND entity_id = p_entity_id
    AND status = 'published';
    RETURN avg_rating;
END;
$$;
CREATE FUNCTION public.check_price_manipulation(p_listing_id integer) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    manipulation_detected BOOLEAN := FALSE;
    manipulation_date TIMESTAMP;
    rehabilitation_period INTERVAL := INTERVAL '30 days'; -- Период "реабилитации"
BEGIN
    -- Проверяем наличие записи о манипуляции в метаданных
    SELECT metadata->>'manipulation_detected_at' INTO manipulation_date
    FROM marketplace_listings
    WHERE id = p_listing_id;
    -- Если запись есть и дата манипуляции меньше 30 дней назад
    IF manipulation_date IS NOT NULL AND
       (manipulation_date::TIMESTAMP + rehabilitation_period) > CURRENT_TIMESTAMP THEN
        RETURN TRUE;
    END IF;
    -- Ищем паттерн: резкое повышение цены > 30% с последующим быстрим снижением в течение недели
    WITH price_history_ordered AS (
        SELECT
            price,
            effective_from,
            effective_to,
            EXTRACT(EPOCH FROM (COALESCE(effective_to, CURRENT_TIMESTAMP) - effective_from))/86400 as duration_days,
            LAG(price) OVER (ORDER BY effective_from) as prev_price,
            LEAD(price) OVER (ORDER BY effective_from) as next_price
        FROM price_history
        WHERE listing_id = p_listing_id
        AND effective_from > CURRENT_TIMESTAMP - INTERVAL '30 days'
        ORDER BY effective_from
    )
    SELECT
        COUNT(*) > 0 INTO manipulation_detected
    FROM price_history_ordered pho
    WHERE pho.prev_price IS NOT NULL
      AND pho.next_price IS NOT NULL
      AND pho.price > pho.prev_price * 1.3  -- повышение более чем на 30%
      AND pho.duration_days < 7             -- действовало менее 7 дней
      AND pho.next_price < pho.price * 0.9  -- быстрое снижение более чем на 10%
      AND pho.next_price > pho.prev_price * 0.9; -- но не слишком низкое по отношению к начальной цене
    -- Если обнаружена манипуляция, сохраняем дату в метаданных
    IF manipulation_detected THEN
        UPDATE marketplace_listings
        SET metadata = COALESCE(metadata, '{}'::jsonb) ||
                      jsonb_build_object('manipulation_detected_at', CURRENT_TIMESTAMP)
        WHERE id = p_listing_id;
    ELSE
        -- Если манипуляция не обнаружена, но была ранее - очищаем метку
        IF manipulation_date IS NOT NULL THEN
            UPDATE marketplace_listings
            SET metadata = metadata - 'manipulation_detected_at'
            WHERE id = p_listing_id;
        END IF;
    END IF;
    RETURN manipulation_detected;
END;
$$;
CREATE FUNCTION public.check_user_permission(p_user_id integer, p_permission_name character varying) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    has_permission BOOLEAN;
BEGIN
    -- Check if user has permission through their role
    SELECT EXISTS (
        SELECT 1
        FROM users u
        JOIN roles r ON u.role_id = r.id
        JOIN role_permissions rp ON r.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = p_user_id
        AND p.name = p_permission_name
    ) INTO has_permission;
    -- Also check user_roles table for multiple roles
    IF NOT has_permission THEN
        SELECT EXISTS (
            SELECT 1
            FROM user_roles ur
            JOIN role_permissions rp ON ur.role_id = rp.role_id
            JOIN permissions p ON rp.permission_id = p.id
            WHERE ur.user_id = p_user_id
            AND p.name = p_permission_name
        ) INTO has_permission;
    END IF;
    RETURN has_permission;
END;
$$;
CREATE FUNCTION public.cleanup_expired_geocoding_cache() RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM geocoding_cache WHERE expires_at < CURRENT_TIMESTAMP;
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$;
CREATE FUNCTION public.cleanup_expired_refresh_tokens() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM refresh_tokens
    WHERE expires_at < CURRENT_TIMESTAMP
    OR (is_revoked = TRUE AND revoked_at < CURRENT_TIMESTAMP - INTERVAL '30 days');
END;
$$;
CREATE FUNCTION public.cleanup_old_address_logs() RETURNS integer
    LANGUAGE plpgsql
    AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM address_change_log
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '2 years';
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    -- Логируем результат очистки
    INSERT INTO address_change_log (
        listing_id, user_id, change_reason, created_at
    ) VALUES (
        0, 0, 'cleanup_old_logs_' || deleted_count, CURRENT_TIMESTAMP
    );
    RETURN deleted_count;
END;
$$;
CREATE FUNCTION public.cleanup_unified_geo() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM unified_geo
    WHERE source_type = 'storefront_product' AND source_id = OLD.id;
    RETURN OLD;
END;
$$;
CREATE FUNCTION public.expand_search_query(query_text text, query_language character varying DEFAULT 'ru'::character varying) RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
    word TEXT;
    synonym_text TEXT;
    expanded_query TEXT := query_text;
    synonyms_array TEXT[];
BEGIN
    -- Split query into words
    FOREACH word IN ARRAY string_to_array(lower(query_text), ' ')
    LOOP
        -- Find all active synonyms for this word
        SELECT array_agg(DISTINCT synonym) INTO synonyms_array
        FROM search_synonyms
        WHERE is_active = true
          AND language = query_language
          AND (term = word OR synonym = word);
        -- If synonyms found, add them to the query
        IF synonyms_array IS NOT NULL THEN
            synonym_text := array_to_string(synonyms_array, ' | ');
            expanded_query := expanded_query || ' | ' || synonym_text;
        END IF;
    END LOOP;
    RETURN expanded_query;
END;
$$;
CREATE FUNCTION public.generate_order_number() RETURNS character varying
    LANGUAGE plpgsql
    AS $$
DECLARE
    order_num VARCHAR(32);
    counter INTEGER := 0;
BEGIN
    LOOP
        -- Генерируем номер в формате: STO-YYYYMMDD-XXXXX
        order_num := 'STO-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' ||
                    LPAD((EXTRACT(epoch FROM CURRENT_TIMESTAMP)::INTEGER % 100000)::TEXT, 5, '0');
        -- Проверяем уникальность
        IF NOT EXISTS (SELECT 1 FROM storefront_orders WHERE order_number = order_num) THEN
            RETURN order_num;
        END IF;
        -- Защита от бесконечного цикла
        counter := counter + 1;
        IF counter > 1000 THEN
            RAISE EXCEPTION 'Unable to generate unique order number after 1000 attempts';
        END IF;
        -- Небольшая задержка перед следующей попыткой
        PERFORM pg_sleep(0.001);
    END LOOP;
END;
$$;
CREATE FUNCTION public.generate_unique_slug(base_name character varying, table_name character varying) RETURNS character varying
    LANGUAGE plpgsql
    AS $$
DECLARE
    slug VARCHAR;
    counter INT := 0;
BEGIN
    -- Convert to lowercase, replace spaces and special chars
    slug := LOWER(REGEXP_REPLACE(base_name, '[^a-zA-Z0-9]+', '-', 'g'));
    slug := TRIM(BOTH '-' FROM slug);
    -- Check if slug exists
    LOOP
        IF counter = 0 THEN
            -- First try without number
            EXIT WHEN NOT EXISTS (
                SELECT 1 FROM storefronts WHERE slug = slug
            );
        ELSE
            -- Add counter
            EXIT WHEN NOT EXISTS (
                SELECT 1 FROM storefronts WHERE slug = slug || '-' || counter
            );
            slug := slug || '-' || counter;
        END IF;
        counter := counter + 1;
    END LOOP;
    RETURN slug;
END;
$$;
CREATE FUNCTION public.get_car_generations_with_translations(p_model_id integer, p_language character varying DEFAULT 'ru'::character varying) RETURNS TABLE(id integer, name character varying, year_from integer, year_to integer, facelift_year integer, body_types jsonb, engine_types jsonb, specs jsonb, image_url character varying)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        g.id,
        COALESCE(
            (SELECT translated_text FROM translations
             WHERE entity_type = 'car_generation'
             AND entity_id = g.id
             AND field_name = 'name'
             AND language = p_language),
            g.name
        ) as name,
        g.year_start as year_from,
        g.year_end as year_to,
        g.facelift_year,
        g.body_types,
        g.engine_types,
        g.specs,
        g.image_url
    FROM car_generations g
    WHERE g.model_id = p_model_id
    AND g.is_active = true
    ORDER BY g.year_start DESC;
END;
$$;
CREATE FUNCTION public.get_geocoding_cache_stats() RETURNS TABLE(total_entries bigint, active_entries bigint, expired_entries bigint, total_cache_hits bigint, avg_confidence numeric, top_providers text[])
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*) as total_entries,
        COUNT(*) FILTER (WHERE expires_at > CURRENT_TIMESTAMP) as active_entries,
        COUNT(*) FILTER (WHERE expires_at <= CURRENT_TIMESTAMP) as expired_entries,
        COALESCE(SUM(cache_hits), 0) as total_cache_hits,
        ROUND(AVG(confidence), 2) as avg_confidence,
        ARRAY_AGG(DISTINCT provider) as top_providers
    FROM geocoding_cache;
END;
$$;
CREATE FUNCTION public.get_telegram_translation(p_key character varying, p_language character varying DEFAULT 'ru'::character varying) RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_translation TEXT;
BEGIN
    SELECT translated_text INTO v_translation
    FROM translations
    WHERE entity_type = 'telegram_bot'
    AND field_name = p_key
    AND language = p_language
    LIMIT 1;
    -- Если перевод не найден, пробуем английский
    IF v_translation IS NULL AND p_language != 'en' THEN
        SELECT translated_text INTO v_translation
        FROM translations
        WHERE entity_type = 'telegram_bot'
        AND field_name = p_key
        AND language = 'en'
        LIMIT 1;
    END IF;
    -- Если и английский не найден, возвращаем ключ
    IF v_translation IS NULL THEN
        RETURN p_key;
    END IF;
    RETURN v_translation;
END;
$$;
CREATE FUNCTION public.increment_keyword_usage(p_category_id integer, p_keywords text[], p_language character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE category_keywords
    SET usage_count = usage_count + 1,
        updated_at = CURRENT_TIMESTAMP
    WHERE category_id = p_category_id
        AND keyword = ANY(p_keywords)
        AND (language = p_language OR language = '*');
END;
$$;
CREATE FUNCTION public.log_role_change() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.role_id IS DISTINCT FROM NEW.role_id THEN
        INSERT INTO role_audit_log (
            target_user_id,
            action,
            old_role_id,
            new_role_id,
            details
        ) VALUES (
            NEW.id,
            'role_changed',
            OLD.role_id,
            NEW.role_id,
            jsonb_build_object(
                'old_role', (SELECT name FROM roles WHERE id = OLD.role_id),
                'new_role', (SELECT name FROM roles WHERE id = NEW.role_id),
                'timestamp', CURRENT_TIMESTAMP
            )
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.log_search_weight_changes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Логируем только изменения веса
    IF OLD.weight <> NEW.weight THEN
        INSERT INTO search_weights_history (
            weight_id,
            old_weight,
            new_weight,
            change_reason,
            changed_by
        ) VALUES (
            NEW.id,
            OLD.weight,
            NEW.weight,
            'manual',  -- По умолчанию, может быть переопределено в коде
            NEW.updated_by
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.preserve_review_origin() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Обновляем отзывы удаляемого товара, сохраняя информацию о продавце
    UPDATE reviews
    SET entity_origin_type = 'user',
        entity_origin_id = OLD.user_id
    WHERE entity_type = 'listing'
      AND entity_id = OLD.id
      AND entity_origin_type IS NULL;
    -- Если товар из витрины, добавляем и эту информацию
    IF OLD.storefront_id IS NOT NULL THEN
        UPDATE reviews
        SET entity_origin_type = 'storefront',
            entity_origin_id = OLD.storefront_id
        WHERE entity_type = 'listing'
          AND entity_id = OLD.id
          AND entity_origin_type = 'user';
    END IF;
    RETURN OLD;
END;
$$;
CREATE FUNCTION public.refresh_discount_metadata() RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    r RECORD;
    max_price DECIMAL(12,2);
    max_price_date TIMESTAMP;
    current_price DECIMAL(12,2);
    percentage DECIMAL(10,2);
    metadata_json JSONB;
    min_price_duration INT := 3; -- Минимальная длительность цены в днях
    is_manipulation BOOLEAN := FALSE;
BEGIN
    -- Обработка всех объявлений с историей цен
    FOR r IN
        SELECT DISTINCT ml.id, ml.price, ml.metadata
        FROM marketplace_listings ml
        JOIN price_history ph ON ml.id = ph.listing_id
        WHERE ml.status = 'active'
        GROUP BY ml.id, ml.price, ml.metadata
        HAVING COUNT(ph.*) > 1
    LOOP
        -- Проверяем на манипуляции с ценами
        SELECT check_price_manipulation(r.id) INTO is_manipulation;
        -- Если обнаружены манипуляции, удаляем метку скидки и переходим к следующему объявлению
        IF is_manipulation THEN
            -- Если есть метаданные о скидке, удаляем их
            IF r.metadata IS NOT NULL AND r.metadata ? 'discount' THEN
                UPDATE marketplace_listings
                SET metadata = metadata - 'discount'
                WHERE id = r.id;
                RAISE NOTICE 'Удалена информация о скидке из-за обнаружения манипуляций с ценой для объявления %d', r.id;
            END IF;
            CONTINUE; -- Переходим к следующему объявлению
        END IF;
        -- Получаем текущую цену
        current_price := r.price;
        -- Находим максимальную цену из истории с учетом длительности
        SELECT price, effective_from INTO max_price, max_price_date
        FROM price_history
        WHERE listing_id = r.id
          AND EXTRACT(EPOCH FROM (COALESCE(effective_to, CURRENT_TIMESTAMP) - effective_from))/86400 >= min_price_duration
        ORDER BY price DESC
        LIMIT 1;
        -- Если максимальная цена не найдена, ищем просто максимальную выше текущей
        IF max_price IS NULL THEN
            SELECT price, effective_from INTO max_price, max_price_date
            FROM price_history
            WHERE listing_id = r.id
              AND price > current_price -- Только если выше текущей
            ORDER BY price DESC
            LIMIT 1;
        END IF;
        -- Если максимальная цена найдена и текущая цена ниже
        IF max_price IS NOT NULL AND current_price < max_price THEN
            -- Вычисляем процент скидки
            percentage := ((current_price - max_price) / max_price) * 100;
            -- Если процент скидки значительный (>= 5%)
            IF ABS(percentage) >= 5 THEN
                -- Подготавливаем метаданные
                metadata_json := COALESCE(r.metadata, '{}'::jsonb);
                -- Создаем информацию о скидке
                metadata_json := jsonb_set(
                    metadata_json,
                    '{discount}',
                    jsonb_build_object(
                        'discount_percent', ROUND(ABS(percentage)),
                        'previous_price', max_price,
                        'effective_from', max_price_date,
                        'has_price_history', true
                    )
                );
                -- Обновляем метаданные
                UPDATE marketplace_listings
                SET metadata = metadata_json
                WHERE id = r.id;
                RAISE NOTICE 'Обновлена информация о скидке для объявления %: %.2f -> %.2f (скидка %.0f%%)',
                    r.id, max_price, current_price, ABS(percentage);
            ELSIF r.metadata IS NOT NULL AND r.metadata ? 'discount' THEN
                -- Если скидка меньше 5%, но были метаданные о скидке - удаляем их
                UPDATE marketplace_listings
                SET metadata = metadata - 'discount'
                WHERE id = r.id;
                RAISE NOTICE 'Удалена неактуальная информация о скидке для объявления %d (скидка меньше 5%%)', r.id;
            END IF;
        ELSIF r.metadata IS NOT NULL AND r.metadata ? 'discount' THEN
            -- Если нет условий для скидки, но были метаданные о скидке - удаляем их
            UPDATE marketplace_listings
            SET metadata = metadata - 'discount'
            WHERE id = r.id;
            RAISE NOTICE 'Удалена неактуальная информация о скидке для объявления %d', r.id;
        END IF;
    END LOOP;
END;
$$;
CREATE FUNCTION public.reserve_inventory_on_order() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    inv_record RECORD;
BEGIN
    -- Находим товар на складе
    SELECT * INTO inv_record
    FROM warehouse_inventory
    WHERE (storefront_product_id = NEW.product_id OR marketplace_listing_id = NEW.listing_id)
    AND quantity_available >= NEW.quantity
    LIMIT 1;
    IF FOUND THEN
        -- Резервируем товар
        UPDATE warehouse_inventory
        SET
            quantity_available = quantity_available - NEW.quantity,
            quantity_reserved = quantity_reserved + NEW.quantity,
            updated_at = NOW()
        WHERE id = inv_record.id;
        -- Создаем запись о движении
        INSERT INTO warehouse_movements (
            warehouse_id, inventory_id, movement_type, movement_reason,
            quantity, quantity_before, quantity_after,
            order_id, storefront_order_id
        ) VALUES (
            inv_record.warehouse_id, inv_record.id, 'outbound', 'order_fulfillment',
            -NEW.quantity, inv_record.quantity_available, inv_record.quantity_available - NEW.quantity,
            NEW.marketplace_order_id, NEW.storefront_order_id
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.sync_attribute_option_translations() RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    attr RECORD;
    opt RECORD;
    opt_index INTEGER;
    options_json JSONB;
    attr_id INTEGER;
BEGIN
    FOR attr IN SELECT DISTINCT attribute_name FROM attribute_option_translations LOOP
        -- Получение ID атрибута
        SELECT id, options INTO attr_id, options_json FROM category_attributes WHERE name = attr.attribute_name;
        IF attr_id IS NOT NULL AND options_json IS NOT NULL THEN
            -- Получение массива значений из JSON
            IF options_json ? 'values' THEN
                FOR opt_index IN 0..jsonb_array_length(options_json->'values')-1 LOOP
                    -- Получение оригинального значения
                    DECLARE
                        option_value TEXT := options_json->'values'->opt_index;
                        option_text TEXT;
                    BEGIN
                        -- Удаление кавычек из JSON-значения
                        option_text := trim(both '"' from option_value::text);
                        -- Добавление перевода на русский
                        INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at)
                        SELECT 'attribute_option', attr_id, 'ru', 'option_' || option_text, aot.ru_translation, false, true, NOW(), NOW()
                        FROM attribute_option_translations aot
                        WHERE aot.attribute_name = attr.attribute_name AND aot.option_value = option_text
                        ON CONFLICT DO NOTHING;
                        -- Добавление перевода на сербский
                        INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at)
                        SELECT 'attribute_option', attr_id, 'sr', 'option_' || option_text, aot.sr_translation, false, true, NOW(), NOW()
                        FROM attribute_option_translations aot
                        WHERE aot.attribute_name = attr.attribute_name AND aot.option_value = option_text
                        ON CONFLICT DO NOTHING;
                    END;
                END LOOP;
            END IF;
        END IF;
    END LOOP;
END;
$$;
CREATE FUNCTION public.track_inventory_changes() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF OLD.quantity_total != NEW.quantity_total OR
       OLD.quantity_available != NEW.quantity_available OR
       OLD.quantity_reserved != NEW.quantity_reserved THEN
        INSERT INTO warehouse_movements (
            warehouse_id, inventory_id, movement_type, movement_reason,
            quantity, quantity_before, quantity_after
        ) VALUES (
            NEW.warehouse_id, NEW.id, 'adjustment', 'manual_adjustment',
            NEW.quantity_total - OLD.quantity_total,
            OLD.quantity_total, NEW.quantity_total
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.trigger_update_listings_on_attribute_translation_change() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Отмечаем все объявления с этим атрибутом как требующие переиндексации
    UPDATE marketplace_listings
    SET needs_reindex = true
    WHERE id IN (
        SELECT DISTINCT lav.listing_id
        FROM listing_attribute_values lav
        JOIN category_attributes ca ON lav.attribute_id = ca.id
        WHERE ca.name = NEW.attribute_name
    );
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_category_attribute_sort_order() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Если sort_order не указан, берем его из атрибута
    IF NEW.sort_order = 0 THEN
        SELECT sort_order INTO NEW.sort_order
        FROM category_attributes
        WHERE id = NEW.attribute_id;
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_listing_metadata_after_price_change() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    last_price DECIMAL(12,2);
    max_price DECIMAL(12,2);
    max_price_date TIMESTAMP;
    price_diff DECIMAL(10,2);
    percentage DECIMAL(10,2);
    current_timestamp_var TIMESTAMP := CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
    metadata_json JSONB;
    listing_data RECORD;
    min_price_duration INT := 3; -- Минимальная длительность цены в днях для учета в расчете скидки
    is_manipulation BOOLEAN := FALSE;
BEGIN
    -- Получаем текущее состояние объявления и его метаданные
    SELECT price, metadata INTO listing_data
    FROM marketplace_listings
    WHERE id = NEW.listing_id;
    -- Проверяем на манипуляции с ценами
    SELECT check_price_manipulation(NEW.listing_id) INTO is_manipulation;
    -- Если обнаружены манипуляции с ценой, не добавляем метку скидки
    IF is_manipulation THEN
        RAISE NOTICE 'Объявление % помечено как манипуляция с ценой, скидка не будет применена', NEW.listing_id;
        -- Если уже есть метаданные о скидке, удаляем их
        IF listing_data.metadata IS NOT NULL AND listing_data.metadata ? 'discount' THEN
            metadata_json := listing_data.metadata - 'discount';
            -- Обновляем метаданные, удаляя информацию о скидке
            UPDATE marketplace_listings
            SET metadata = metadata_json
            WHERE id = NEW.listing_id;
            RAISE NOTICE 'Удалена информация о скидке из-за обнаружения манипуляций с ценой для объявления %d', NEW.listing_id;
        END IF;
        RETURN NULL;
    END IF;
    -- Получаем максимальную цену из истории, которая существовала достаточно долго
    SELECT price, effective_from INTO max_price, max_price_date
    FROM price_history
    WHERE listing_id = NEW.listing_id
      AND EXTRACT(EPOCH FROM (COALESCE(effective_to, CURRENT_TIMESTAMP) - effective_from))/86400 >= min_price_duration
    ORDER BY price DESC
    LIMIT 1;
    -- Если максимальная цена не найдена, ищем просто максимальную
    IF max_price IS NULL THEN
        SELECT price, effective_from INTO max_price, max_price_date
        FROM price_history
        WHERE listing_id = NEW.listing_id
          AND price > NEW.price -- Только если выше текущей
        ORDER BY price DESC
        LIMIT 1;
    END IF;
    -- Получаем предыдущую цену (непосредственно перед текущим изменением)
    SELECT price INTO last_price
    FROM price_history
    WHERE listing_id = NEW.listing_id
    AND effective_to IS NOT NULL
    ORDER BY effective_to DESC
    LIMIT 1;
    -- Ключевая логика обработки скидок - если есть скидка от максимальной цены
    IF max_price IS NOT NULL AND NEW.price < max_price THEN
        -- Вычисляем процент скидки
        percentage := ((NEW.price - max_price) / max_price) * 100;
        -- Если процент скидки значительный (>= 5%)
        IF ABS(percentage) >= 5 THEN
            -- Подготавливаем или обновляем метаданные
            metadata_json := COALESCE(listing_data.metadata, '{}'::jsonb);
            -- Создаем информацию о скидке
            metadata_json := jsonb_set(
                metadata_json,
                '{discount}',
                jsonb_build_object(
                    'discount_percent', ROUND(ABS(percentage)),
                    'previous_price', max_price,
                    'effective_from', max_price_date,
                    'has_price_history', true
                )
            );
            -- Обновляем метаданные
            UPDATE marketplace_listings
            SET metadata = metadata_json
            WHERE id = NEW.listing_id;
            RAISE NOTICE 'Обновлена информация о скидке для объявления %: %.2f -> %.2f (скидка %.0f%%)',
                NEW.listing_id, max_price, NEW.price, ABS(percentage);
        END IF;
    ELSIF listing_data.metadata IS NOT NULL AND listing_data.metadata ? 'discount' THEN
        -- Если скидка больше не актуальна, удаляем информацию о ней
        -- (например, если цена выросла выше максимальной)
        metadata_json := listing_data.metadata - 'discount';
        -- Обновляем метаданные, удаляя информацию о скидке
        UPDATE marketplace_listings
        SET metadata = metadata_json
        WHERE id = NEW.listing_id;
        RAISE NOTICE 'Удалена неактуальная информация о скидке для объявления %d', NEW.listing_id;
    END IF;
    RETURN NULL;
END;
$$;
CREATE FUNCTION public.update_price_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    last_price DECIMAL(12,2);
    price_diff DECIMAL(10,2);
    percentage DECIMAL(10,2);
    current_timestamp_var TIMESTAMP := CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
BEGIN
    -- Получаем последнюю цену
    SELECT price INTO last_price
    FROM price_history
    WHERE listing_id = NEW.id
    AND effective_to IS NULL
    ORDER BY effective_from DESC
    LIMIT 1;
    -- Если цена изменилась или это новое объявление
    IF last_price IS NULL OR last_price != NEW.price THEN
        -- Если это не новое объявление, закрываем старую запись
        IF last_price IS NOT NULL THEN
            UPDATE price_history
            SET effective_to = current_timestamp_var
            WHERE listing_id = NEW.id
            AND effective_to IS NULL;
            -- Вычисляем процент изменения цены
            price_diff := NEW.price - last_price;
            IF last_price != 0 THEN  -- Избегаем деления на ноль
                percentage := (price_diff / last_price) * 100;
            ELSE
                percentage := NULL;
            END IF;
        END IF;
        -- Создаем новую запись с текущей ценой
        INSERT INTO price_history (
            listing_id,
            price,
            effective_from,
            change_source,
            change_percentage
        ) VALUES (
            NEW.id,
            NEW.price,
            current_timestamp_var,
            TG_ARGV[0],  -- Источник изменения передается как аргумент триггера
            percentage
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_product_has_variants() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE storefront_products
        SET has_variants = true
        WHERE id = NEW.product_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE storefront_products
        SET has_variants = EXISTS (
            SELECT 1
            FROM storefront_product_variants
            WHERE product_id = OLD.product_id
        )
        WHERE id = OLD.product_id;
    END IF;
    RETURN NULL;
END;
$$;
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
