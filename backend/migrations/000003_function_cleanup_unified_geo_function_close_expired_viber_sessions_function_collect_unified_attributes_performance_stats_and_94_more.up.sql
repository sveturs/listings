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
    CONSTRAINT ai_category_decisions_confidence_check CHECK (((confidence >= (0)::numeric) AND (confidence <= (1)::numeric)))
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
CREATE TABLE public.category_detection_cache (
    id bigint NOT NULL,
    cache_key character varying(255) NOT NULL,
    category_id integer,
    confidence_score numeric(5,4),
    ai_hints jsonb,
    keywords text[],
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now()
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
CREATE TABLE public.category_keyword_weights (
    id bigint NOT NULL,
    keyword character varying(255) NOT NULL,
    category_id integer,
    weight numeric(5,4) DEFAULT 1.0000,
    occurrence_count integer DEFAULT 1,
    success_rate numeric(5,4) DEFAULT 0.5000,
    language character varying(10) DEFAULT 'ru'::character varying,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
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
    needs_reindex boolean DEFAULT false,
    address_multilingual jsonb
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
CREATE TABLE public.delivery_category_defaults (
    id integer NOT NULL,
    category_id integer,
    default_weight_kg numeric(10,3),
    default_length_cm numeric(10,2),
    default_width_cm numeric(10,2),
    default_height_cm numeric(10,2),
    default_packaging_type character varying(50),
    is_typically_fragile boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_notifications (
    id integer NOT NULL,
    shipment_id integer,
    user_id integer,
    channel character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    template character varying(50),
    data jsonb,
    sent_at timestamp with time zone,
    error_message text,
    retry_count integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_pricing_rules (
    id integer NOT NULL,
    provider_id integer,
    rule_type character varying(50) NOT NULL,
    weight_ranges jsonb,
    volume_ranges jsonb,
    zone_multipliers jsonb,
    fragile_surcharge numeric(10,2) DEFAULT 0,
    oversized_surcharge numeric(10,2) DEFAULT 0,
    special_handling_surcharge numeric(10,2) DEFAULT 0,
    min_price numeric(10,2),
    max_price numeric(10,2),
    custom_formula text,
    priority integer DEFAULT 0,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_providers (
    id integer NOT NULL,
    code character varying(50) NOT NULL,
    name character varying(255) NOT NULL,
    logo_url character varying(500),
    is_active boolean DEFAULT false,
    supports_cod boolean DEFAULT false,
    supports_insurance boolean DEFAULT false,
    supports_tracking boolean DEFAULT true,
    api_config jsonb,
    capabilities jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_shipments (
    id integer NOT NULL,
    provider_id integer,
    order_id integer,
    external_id character varying(255),
    tracking_number character varying(255),
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL,
    sender_info jsonb NOT NULL,
    recipient_info jsonb NOT NULL,
    package_info jsonb NOT NULL,
    delivery_cost numeric(10,2),
    insurance_cost numeric(10,2),
    cod_amount numeric(10,2),
    cost_breakdown jsonb,
    pickup_date date,
    estimated_delivery date,
    actual_delivery_date timestamp with time zone,
    provider_response jsonb,
    labels jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_tracking_events (
    id integer NOT NULL,
    shipment_id integer,
    provider_id integer,
    event_time timestamp with time zone NOT NULL,
    status character varying(100) NOT NULL,
    location character varying(500),
    description text,
    raw_data jsonb,
    created_at timestamp with time zone DEFAULT now()
);
CREATE TABLE public.delivery_zones (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    type character varying(50) NOT NULL,
    countries text[],
    regions text[],
    cities text[],
    postal_codes text[],
    boundary public.geometry(Polygon,4326),
    center_point public.geometry(Point,4326),
    radius_km numeric(10,2),
    created_at timestamp with time zone DEFAULT now()
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
    unit character varying(50),
    value_type character varying(50) NOT NULL
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
