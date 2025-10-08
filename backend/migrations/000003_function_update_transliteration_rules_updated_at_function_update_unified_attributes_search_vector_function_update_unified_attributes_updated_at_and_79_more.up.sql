CREATE FUNCTION public.update_transliteration_rules_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_unified_attributes_search_vector() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', coalesce(NEW.name, '') || ' ' || coalesce(NEW.code, ''));
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.update_unified_attributes_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
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
CREATE FUNCTION public.update_vin_decode_cache_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
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
CREATE FUNCTION public.check_return_visit() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Проверяем, был ли этот пользователь раньше
    IF NEW.user_id IS NOT NULL THEN
        NEW.is_return_visit := EXISTS (
            SELECT 1 FROM user_view_history
            WHERE user_id = NEW.user_id
            AND listing_id = NEW.listing_id
            AND id < NEW.id
            LIMIT 1
        );
    END IF;
    RETURN NEW;
END;
$$;
CREATE FUNCTION public.check_subscription_limits(p_user_id integer, p_resource_type character varying, p_count integer DEFAULT 1) RETURNS boolean
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_subscription RECORD;
    v_plan RECORD;
    v_current_usage INTEGER;
BEGIN
    -- Получаем активную подписку пользователя
    SELECT us.*, sp.*
    INTO v_subscription
    FROM user_subscriptions us
    JOIN subscription_plans sp ON us.plan_id = sp.id
    WHERE us.user_id = p_user_id
    AND us.status IN ('active', 'trial')
    LIMIT 1;
    IF NOT FOUND THEN
        -- Если нет подписки, используем бесплатный план
        SELECT * INTO v_plan
        FROM subscription_plans
        WHERE code = 'starter'
        LIMIT 1;
        v_subscription.id = NULL;
    ELSE
        v_plan = v_subscription;
    END IF;
    -- Проверяем лимиты в зависимости от типа ресурса
    CASE p_resource_type
        WHEN 'storefront' THEN
            -- Для unlimited возвращаем true
            IF v_plan.max_storefronts = -1 THEN
                RETURN TRUE;
            END IF;
            -- Считаем текущие витрины
            SELECT COUNT(*) INTO v_current_usage
            FROM storefronts
            WHERE user_id = p_user_id
            AND is_active = TRUE;
            RETURN (v_current_usage + p_count) <= v_plan.max_storefronts;
        WHEN 'product' THEN
            IF v_plan.max_products_per_storefront = -1 THEN
                RETURN TRUE;
            END IF;
            -- Логика для продуктов будет добавлена позже
            RETURN TRUE;
        WHEN 'staff' THEN
            IF v_plan.max_staff_per_storefront = -1 THEN
                RETURN TRUE;
            END IF;
            -- Логика для персонала будет добавлена позже
            RETURN TRUE;
        ELSE
            RETURN FALSE;
    END CASE;
END;
$$;
CREATE FUNCTION public.cleanup_detection_cache() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM category_detection_cache WHERE expires_at < NOW();
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
CREATE FUNCTION public.close_expired_viber_sessions() RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE viber_sessions
    SET active = false
    WHERE active = true
    AND expires_at < CURRENT_TIMESTAMP;
END;
$$;
CREATE FUNCTION public.collect_unified_attributes_performance_stats() RETURNS TABLE(metric_name text, metric_value numeric, metric_unit text, category text)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        'total_attributes'::TEXT,
        COUNT(*)::NUMERIC,
        'count'::TEXT,
        'system'::TEXT
    FROM unified_attributes WHERE is_active = true
    UNION ALL
    SELECT
        'total_attribute_values'::TEXT,
        COUNT(*)::NUMERIC,
        'count'::TEXT,
        'system'::TEXT
    FROM unified_attribute_values
    UNION ALL
    SELECT
        'avg_attributes_per_category'::TEXT,
        AVG(attr_count)::NUMERIC,
        'count'::TEXT,
        'performance'::TEXT
    FROM (
        SELECT category_id, COUNT(*) as attr_count
        FROM unified_category_attributes
        WHERE is_enabled = true
        GROUP BY category_id
    ) t
    UNION ALL
    SELECT
        'categories_with_attributes'::TEXT,
        COUNT(DISTINCT category_id)::NUMERIC,
        'count'::TEXT,
        'system'::TEXT
    FROM unified_category_attributes WHERE is_enabled = true
    UNION ALL
    SELECT
        'unused_attributes'::TEXT,
        COUNT(*)::NUMERIC,
        'count'::TEXT,
        'optimization'::TEXT
    FROM unified_attributes ua
    LEFT JOIN unified_attribute_values uav ON ua.id = uav.attribute_id
    WHERE ua.is_active = true AND uav.id IS NULL;
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
CREATE FUNCTION public.extract_city_from_location(full_address text) RETURNS text
    LANGUAGE plpgsql IMMUTABLE
    AS $_$
DECLARE
    parts TEXT[];
    city_part TEXT;
BEGIN
    IF full_address IS NULL OR full_address = '' THEN
        RETURN NULL;
    END IF;
    -- Разбиваем адрес по запятой
    parts := string_to_array(full_address, ',');
    -- Если только одна часть - возвращаем как есть (скорее всего уже город)
    IF array_length(parts, 1) = 1 THEN
        RETURN trim(parts[1]);
    END IF;
    -- Если больше одной части - берем вторую (index 2 в PostgreSQL, т.к. массивы с 1)
    IF array_length(parts, 1) > 1 THEN
        city_part := trim(parts[2]);
        -- Убираем почтовый код (5 цифр в конце): "Нови Сад 21101" → "Нови Сад"
        city_part := regexp_replace(city_part, '\s+\d{5}.*$', '');
        RETURN trim(city_part);
    END IF;
    RETURN NULL;
END;
$_$;
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
CREATE FUNCTION public.get_category_accuracy_report(days_back integer DEFAULT 7) RETURNS TABLE(category_id integer, category_name character varying, total_detections bigint, correct_detections bigint, accuracy_percent numeric)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.id,
        c.name,
        COUNT(f.id) as total_detections,
        SUM(CASE WHEN f.detected_category_id = f.correct_category_id THEN 1 ELSE 0 END) as correct_detections,
        ROUND(100.0 * SUM(CASE WHEN f.detected_category_id = f.correct_category_id THEN 1 ELSE 0 END) /
              NULLIF(COUNT(f.id), 0), 2) as accuracy_percent
    FROM marketplace_categories c
    LEFT JOIN category_detection_feedback f ON f.detected_category_id = c.id
    WHERE f.created_at > NOW() - (days_back || ' days')::INTERVAL
        AND f.user_confirmed = TRUE
    GROUP BY c.id, c.name
    HAVING COUNT(f.id) > 0
    ORDER BY accuracy_percent DESC, total_detections DESC;
END;
$$;
CREATE FUNCTION public.get_category_by_ai_hints(p_domain character varying, p_product_type character varying) RETURNS TABLE(category_id integer, confidence numeric)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        m.category_id,
        m.weight * (1.0 + (m.success_count::DECIMAL / GREATEST(m.success_count + m.failure_count, 1)) * 0.2) as confidence
    FROM category_ai_mappings m
    WHERE m.ai_domain = p_domain
      AND m.product_type = p_product_type
      AND m.is_active = TRUE
    ORDER BY confidence DESC
    LIMIT 1;
END;
$$;
CREATE FUNCTION public.get_delivery_attributes(p_product_id integer, p_product_type character varying DEFAULT 'listing'::character varying) RETURNS jsonb
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_attributes JSONB;
    v_category_id INTEGER;
    v_defaults RECORD;
BEGIN
    IF p_product_type = 'listing' THEN
        SELECT
            ml.metadata->'delivery_attributes',
            ml.category_id
        INTO v_attributes, v_category_id
        FROM marketplace_listings ml
        WHERE ml.id = p_product_id;
    ELSE
        SELECT
            sp.attributes->'delivery_attributes',
            sp.category_id
        INTO v_attributes, v_category_id
        FROM storefront_products sp
        WHERE sp.id = p_product_id;
    END IF;
    -- Если атрибуты пустые, используем дефолтные из категории
    IF v_attributes IS NULL OR v_attributes = '{}'::jsonb THEN
        SELECT * INTO v_defaults
        FROM delivery_category_defaults
        WHERE category_id = v_category_id;
        IF FOUND THEN
            v_attributes = jsonb_build_object(
                'weight_kg', v_defaults.default_weight_kg,
                'dimensions', jsonb_build_object(
                    'length_cm', v_defaults.default_length_cm,
                    'width_cm', v_defaults.default_width_cm,
                    'height_cm', v_defaults.default_height_cm
                ),
                'packaging_type', v_defaults.default_packaging_type,
                'is_fragile', v_defaults.is_typically_fragile
            );
        END IF;
    END IF;
    RETURN COALESCE(v_attributes, '{}'::jsonb);
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
CREATE FUNCTION public.get_user_subscription(p_user_id integer) RETURNS TABLE(subscription_id integer, plan_code character varying, plan_name character varying, status character varying, expires_at timestamp without time zone, max_storefronts integer, used_storefronts integer, max_products integer, max_staff integer, max_images integer, has_ai boolean, has_live boolean, has_export boolean, has_custom_domain boolean)
    LANGUAGE plpgsql
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        us.id as subscription_id,
        sp.code as plan_code,
        sp.name as plan_name,
        us.status,
        us.expires_at,
        sp.max_storefronts,
        us.used_storefronts,
        sp.max_products_per_storefront as max_products,
        sp.max_staff_per_storefront as max_staff,
        sp.max_images_total as max_images,
        sp.has_ai_assistant as has_ai,
        sp.has_live_shopping as has_live,
        sp.has_export_data as has_export,
        sp.has_custom_domain
    FROM user_subscriptions us
    JOIN subscription_plans sp ON us.plan_id = sp.id
    WHERE us.user_id = p_user_id
    AND us.status IN ('active', 'trial')
    LIMIT 1;
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
CREATE FUNCTION public.sync_storefront_product_to_marketplace() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    v_user_id INTEGER;
    v_user_storefront_id INTEGER;
    v_existing_listing_id INTEGER;
    v_new_status VARCHAR(20);
BEGIN
    -- Получить user_id владельца витрины
    SELECT user_id INTO v_user_id
    FROM storefronts
    WHERE id = COALESCE(NEW.storefront_id, OLD.storefront_id);
    -- Получить user_storefront_id из user_storefronts
    -- marketplace_listings.storefront_id ссылается на user_storefronts.id, а НЕ на storefronts.id
    SELECT id INTO v_user_storefront_id
    FROM user_storefronts
    WHERE user_id = v_user_id
    LIMIT 1;
    -- Определить статус для marketplace_listing
    IF TG_OP = 'DELETE' THEN
        v_new_status := 'inactive';
    ELSE
        v_new_status := CASE WHEN NEW.is_active THEN 'active' ELSE 'inactive' END;
    END IF;
    -- Обработка INSERT: создание нового listing
    IF TG_OP = 'INSERT' THEN
        -- Проверить, не существует ли уже listing с таким же SKU для этой витрины
        IF NEW.sku IS NOT NULL THEN
            SELECT id INTO v_existing_listing_id
            FROM marketplace_listings
            WHERE storefront_id = NEW.storefront_id
              AND external_id = NEW.sku
            LIMIT 1;
        END IF;
        -- Если listing не существует, создать новый
        IF v_existing_listing_id IS NULL THEN
            INSERT INTO marketplace_listings (
                id,
                user_id,
                category_id,
                title,
                description,
                price,
                condition,
                status,
                location,
                latitude,
                longitude,
                show_on_map,
                storefront_id,
                external_id,
                metadata,
                created_at,
                updated_at,
                needs_reindex
            ) VALUES (
                NEW.id,  -- Используем тот же ID (shared sequence)
                v_user_id,
                NEW.category_id,
                NEW.name,
                NEW.description,
                NEW.price,
                'new',  -- Товары витрин всегда новые
                v_new_status,
                COALESCE(NEW.individual_address, ''),
                NEW.individual_latitude,
                NEW.individual_longitude,
                COALESCE(NEW.show_on_map, true),
                v_user_storefront_id,  -- Используем user_storefronts.id, а не storefronts.id
                NEW.sku,  -- SKU как external_id для отслеживания
                jsonb_build_object(
                    'source', 'storefront',
                    'storefront_id', NEW.storefront_id,  -- Сохраняем оригинальный storefront_id в metadata
                    'stock_quantity', NEW.stock_quantity,
                    'stock_status', NEW.stock_status,
                    'currency', NEW.currency,
                    'barcode', NEW.barcode,
                    'attributes', NEW.attributes
                ),
                NEW.created_at,
                NEW.updated_at,
                true  -- Требуется переиндексация в OpenSearch
            );
        ELSE
            -- Если listing существует, обновить его
            UPDATE marketplace_listings
            SET
                title = NEW.name,
                description = NEW.description,
                price = NEW.price,
                status = v_new_status,
                category_id = NEW.category_id,
                location = COALESCE(NEW.individual_address, ''),
                latitude = NEW.individual_latitude,
                longitude = NEW.individual_longitude,
                show_on_map = COALESCE(NEW.show_on_map, true),
                metadata = jsonb_build_object(
                    'source', 'storefront',
                    'storefront_id', NEW.storefront_id,
                    'stock_quantity', NEW.stock_quantity,
                    'stock_status', NEW.stock_status,
                    'currency', NEW.currency,
                    'barcode', NEW.barcode,
                    'attributes', NEW.attributes
                ),
                updated_at = NEW.updated_at,
                needs_reindex = true
            WHERE id = v_existing_listing_id;
        END IF;
    -- Обработка UPDATE: обновление существующего listing
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE marketplace_listings
        SET
            title = NEW.name,
            description = NEW.description,
            price = NEW.price,
            status = v_new_status,
            category_id = NEW.category_id,
            location = COALESCE(NEW.individual_address, ''),
            latitude = NEW.individual_latitude,
            longitude = NEW.individual_longitude,
            show_on_map = COALESCE(NEW.show_on_map, true),
            external_id = NEW.sku,
            metadata = jsonb_build_object(
                'source', 'storefront',
                'storefront_id', NEW.storefront_id,
                'stock_quantity', NEW.stock_quantity,
                'stock_status', NEW.stock_status,
                'currency', NEW.currency,
                'barcode', NEW.barcode,
                'attributes', NEW.attributes
            ),
            updated_at = NEW.updated_at,
            needs_reindex = true
        WHERE id = NEW.id
          AND storefront_id = v_user_storefront_id;
        -- Если listing не найден, создать его (edge case)
        IF NOT FOUND THEN
            INSERT INTO marketplace_listings (
                id,
                user_id,
                category_id,
                title,
                description,
                price,
                condition,
                status,
                location,
                latitude,
                longitude,
                show_on_map,
                storefront_id,
                external_id,
                metadata,
                created_at,
                updated_at,
                needs_reindex
            ) VALUES (
                NEW.id,
                v_user_id,
                NEW.category_id,
                NEW.name,
                NEW.description,
                NEW.price,
                'new',
                v_new_status,
                COALESCE(NEW.individual_address, ''),
                NEW.individual_latitude,
                NEW.individual_longitude,
                COALESCE(NEW.show_on_map, true),
                v_user_storefront_id,
                NEW.sku,
                jsonb_build_object(
                    'source', 'storefront',
                    'storefront_id', NEW.storefront_id,
                    'stock_quantity', NEW.stock_quantity,
                    'stock_status', NEW.stock_status,
                    'currency', NEW.currency,
                    'barcode', NEW.barcode,
                    'attributes', NEW.attributes
                ),
                NEW.created_at,
                NEW.updated_at,
                true
            );
        END IF;
    -- Обработка DELETE: деактивация listing (не удаляем, т.к. могут быть связи)
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE marketplace_listings
        SET
            status = 'inactive',
            updated_at = CURRENT_TIMESTAMP,
            needs_reindex = true
        WHERE id = OLD.id
          AND storefront_id = v_user_storefront_id;
    END IF;
    RETURN COALESCE(NEW, OLD);
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
CREATE FUNCTION public.update_mapping_stats(p_domain character varying, p_product_type character varying, p_category_id integer, p_success boolean) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF p_success THEN
        UPDATE category_ai_mappings
        SET success_count = success_count + 1,
            weight = LEAST(weight * 1.01, 1.0)
        WHERE ai_domain = p_domain
          AND product_type = p_product_type
          AND category_id = p_category_id;
    ELSE
        UPDATE category_ai_mappings
        SET failure_count = failure_count + 1,
            weight = GREATEST(weight * 0.99, 0.1)
        WHERE ai_domain = p_domain
          AND product_type = p_product_type
          AND category_id = p_category_id;
    END IF;
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
CREATE FUNCTION public.update_view_statistics(target_date date DEFAULT CURRENT_DATE) RETURNS void
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Удаляем старую статистику за эту дату
    DELETE FROM view_statistics WHERE date = target_date;
    -- Вставляем новую агрегированную статистику
    INSERT INTO view_statistics (
        listing_id,
        category_id,
        date,
        views_count,
        unique_users_count,
        unique_sessions_count,
        avg_view_duration,
        mobile_views,
        desktop_views,
        tablet_views,
        contact_clicks,
        favorite_adds
    )
    SELECT
        uvh.listing_id,
        ml.category_id,
        target_date,
        COUNT(*) as views_count,
        COUNT(DISTINCT uvh.user_id) as unique_users_count,
        COUNT(DISTINCT uvh.session_id) as unique_sessions_count,
        AVG(uvh.view_duration_seconds) as avg_view_duration,
        COUNT(*) FILTER (WHERE uvh.device_type = 'mobile') as mobile_views,
        COUNT(*) FILTER (WHERE uvh.device_type = 'desktop') as desktop_views,
        COUNT(*) FILTER (WHERE uvh.device_type = 'tablet') as tablet_views,
        COUNT(*) FILTER (WHERE uvh.interaction_type = 'click_phone') as contact_clicks,
        COUNT(*) FILTER (WHERE uvh.interaction_type = 'add_favorite') as favorite_adds
    FROM user_view_history uvh
    JOIN marketplace_listings ml ON ml.id = uvh.listing_id
    WHERE DATE(uvh.viewed_at) = target_date
    GROUP BY uvh.listing_id, ml.category_id;
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
CREATE TABLE public.notification_settings (
    user_id integer NOT NULL,
    notification_type character varying(50) NOT NULL,
    telegram_enabled boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    email_enabled boolean DEFAULT false
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
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    settings jsonb DEFAULT '{}'::jsonb
);
CREATE TABLE public.user_telegram_connections (
    user_id integer NOT NULL,
    telegram_chat_id character varying(100) NOT NULL,
    telegram_username character varying(100),
    connected_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
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
CREATE TABLE public.ai_category_decisions (
    id integer NOT NULL,
    title_hash character varying(64) NOT NULL,
    title text NOT NULL,
    description text,
    category_id integer NOT NULL,
    confidence numeric(5,4) NOT NULL,
    reasoning text,
    alternative_category_ids integer[],
    ai_model character varying(100) DEFAULT 'claude-3-haiku-20240307'::character varying,
    processing_time_ms integer,
    ai_domain character varying(100),
    ai_product_type character varying(100),
    ai_keywords text[],
    user_confirmed boolean DEFAULT false,
    user_corrected_category_id integer,
    user_feedback_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    entity_type character varying(20) DEFAULT 'listing'::character varying NOT NULL,
    CONSTRAINT ai_category_decisions_confidence_check CHECK (((confidence >= (0)::numeric) AND (confidence <= (1)::numeric))),
    CONSTRAINT ai_category_decisions_entity_type_check CHECK (((entity_type)::text = ANY (ARRAY[('listing'::character varying)::text, ('product'::character varying)::text])))
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
CREATE TABLE public.bex_configuration (
    id integer NOT NULL,
    key character varying(100) NOT NULL,
    value text,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE public.bex_shipments (
    id integer NOT NULL,
    tracking_number character varying(100) NOT NULL,
    order_id integer,
    provider_order_id character varying(100),
    sender_name character varying(255) NOT NULL,
    sender_phone character varying(50),
    sender_email character varying(255),
    sender_address text,
    sender_city character varying(100),
    sender_postal_code character varying(20),
    recipient_name character varying(255) NOT NULL,
    recipient_phone character varying(50),
    recipient_email character varying(255),
    recipient_address text,
    recipient_city character varying(100),
    recipient_postal_code character varying(20),
    weight numeric(10,3),
    width numeric(10,2),
    height numeric(10,2),
    length numeric(10,2),
    package_type character varying(50),
    content_description text,
    content_value numeric(10,2),
    status_description text,
    last_status_update timestamp with time zone,
    delivery_type character varying(50),
    delivery_date date,
    estimated_delivery timestamp with time zone,
    actual_delivery timestamp with time zone,
    delivery_cost numeric(10,2),
    cod_amount numeric(10,2),
    insurance_amount numeric(10,2),
    total_cost numeric(10,2),
    bex_shipment_id character varying(100),
    bex_barcode character varying(100),
    label_url text,
    raw_response jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    marketplace_order_id integer,
    storefront_order_id bigint,
    status_text character varying(255),
    registered_at timestamp with time zone,
    delivered_at timestamp with time zone,
    failed_reason text,
    status_history jsonb,
    status integer DEFAULT 1
);
CREATE TABLE public.bex_tracking_events (
    id integer NOT NULL,
    shipment_id integer,
    event_date timestamp with time zone NOT NULL,
    event_type character varying(100),
    event_code character varying(50),
    description text,
    location character varying(255),
    raw_data jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
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
CREATE TABLE public.category_ai_mappings (
    id integer NOT NULL,
    ai_domain character varying(100) NOT NULL,
    product_type character varying(100) NOT NULL,
    category_id integer,
    weight numeric(3,2) DEFAULT 1.00,
    success_count integer DEFAULT 0,
    failure_count integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
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
    show_in_card boolean DEFAULT false,
    CONSTRAINT unified_attributes_attribute_type_check CHECK (((attribute_type)::text = ANY (ARRAY[('text'::character varying)::text, ('textarea'::character varying)::text, ('number'::character varying)::text, ('boolean'::character varying)::text, ('select'::character varying)::text, ('multiselect'::character varying)::text, ('date'::character varying)::text, ('color'::character varying)::text, ('size'::character varying)::text]))),
    CONSTRAINT unified_attributes_purpose_check CHECK (((purpose)::text = ANY (ARRAY[('regular'::character varying)::text, ('variant'::character varying)::text, ('both'::character varying)::text])))
)
WITH (autovacuum_vacuum_scale_factor='0.2', autovacuum_analyze_scale_factor='0.1', autovacuum_vacuum_threshold='50', autovacuum_analyze_threshold='50');
CREATE TABLE public.category_detection_feedback (
    id bigint NOT NULL,
    listing_id integer,
    detected_category_id integer,
    correct_category_id integer,
    ai_hints jsonb DEFAULT '{}'::jsonb NOT NULL,
    keywords text[] DEFAULT '{}'::text[],
    confidence_score numeric(5,4) DEFAULT 0.0000,
    user_confirmed boolean DEFAULT false,
    algorithm_version character varying(50) DEFAULT 'stable_v1'::character varying,
    processing_time_ms integer,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
