-- Исправление функции update_listing_metadata_after_price_change
CREATE OR REPLACE FUNCTION update_listing_metadata_after_price_change()
RETURNS TRIGGER AS $$
DECLARE
    last_price DECIMAL(12,2);
    max_price DECIMAL(12,2);
    max_price_date TIMESTAMP;
    price_diff DECIMAL(10,2);
    percentage DECIMAL(10,2);
    current_timestamp_var TIMESTAMP := CURRENT_TIMESTAMP AT TIME ZONE 'UTC';
    metadata_json JSONB;
    listing_data RECORD;
    min_price_duration INT := 1; -- Минимальная длительность цены в днях для учета в расчете скидки
BEGIN
    -- Получаем текущее состояние объявления и его метаданные
    SELECT price, metadata INTO listing_data
    FROM marketplace_listings
    WHERE id = NEW.listing_id;
    
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
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Создаем новый триггер, который будет срабатывать ПОСЛЕ вставки записи в price_history
CREATE TRIGGER trig_update_metadata_after_price_change
AFTER INSERT ON price_history
FOR EACH ROW
EXECUTE FUNCTION update_listing_metadata_after_price_change();

-- Исправляем существующие метаданные о скидках для всех объявлений
DO $$
DECLARE
    r RECORD;
    max_price DECIMAL(12,2);
    current_price DECIMAL(12,2);
    percentage DECIMAL(10,2);
    metadata_json JSONB;
    min_price_duration INT := 1; -- Минимальная длительность цены в днях
BEGIN
    -- Сначала исправляем объявления, у которых уже есть метаданные о скидке
    FOR r IN SELECT id FROM marketplace_listings WHERE metadata->'discount' IS NOT NULL LOOP
        -- Получаем текущую цену
        SELECT price INTO current_price FROM marketplace_listings WHERE id = r.id;
        
        -- Находим максимальную цену из истории
        WITH price_data AS (
            SELECT 
                price,
                EXTRACT(EPOCH FROM (COALESCE(effective_to, CURRENT_TIMESTAMP) - effective_from))/86400 as duration_days
            FROM price_history
            WHERE listing_id = r.id
            ORDER BY price DESC -- Сортируем по цене для поиска максимума
        )
        SELECT MAX(price) INTO max_price
        FROM price_data
        WHERE duration_days >= min_price_duration;
        
        -- Проверяем, нужна ли скидка
        IF max_price IS NOT NULL AND current_price < max_price THEN
            -- Вычисляем процент скидки
            percentage := ((current_price - max_price) / max_price) * 100;
            
            IF ABS(percentage) >= 5 THEN
                -- Обновляем метаданные
                SELECT metadata INTO metadata_json FROM marketplace_listings WHERE id = r.id;
                
                metadata_json := jsonb_set(
                    COALESCE(metadata_json, '{}'::jsonb),
                    '{discount}',
                    jsonb_build_object(
                        'discount_percent', ROUND(ABS(percentage)),
                        'previous_price', max_price,
                        'effective_from', CURRENT_TIMESTAMP,
                        'has_price_history', true
                    )
                );
                
                UPDATE marketplace_listings
                SET metadata = metadata_json
                WHERE id = r.id;
                
                RAISE NOTICE 'Исправлены метаданные о скидке для объявления %d: %.2f -> %.2f (скидка %.0f%%)',
                    r.id, max_price, current_price, ABS(percentage);
            ELSE
                -- Удаляем информацию о скидке, так как она меньше 5%
                UPDATE marketplace_listings
                SET metadata = metadata - 'discount'
                WHERE id = r.id;
                
                RAISE NOTICE 'Удалена малая скидка (%.1f%%) для объявления %d', ABS(percentage), r.id;
            END IF;
        ELSE
            -- Удаляем информацию о скидке, если она больше не актуальна
            UPDATE marketplace_listings
            SET metadata = metadata - 'discount'
            WHERE id = r.id;
            
            RAISE NOTICE 'Удалена неактуальная информация о скидке для объявления %d', r.id;
        END IF;
    END LOOP;
    
    -- Затем проверяем все объявления с историей цен на наличие неучтенных скидок
    FOR r IN 
        SELECT DISTINCT l.id
        FROM marketplace_listings l
        JOIN price_history ph ON l.id = ph.listing_id
        WHERE (l.metadata IS NULL OR l.metadata->'discount' IS NULL)
        AND EXISTS (
            SELECT 1 FROM price_history ph2 
            WHERE ph2.listing_id = l.id 
            GROUP BY ph2.listing_id 
            HAVING COUNT(*) > 1
        )
    LOOP
        -- Получаем текущую цену
        SELECT price INTO current_price FROM marketplace_listings WHERE id = r.id;
        
        -- Находим максимальную цену из истории
        WITH price_data AS (
            SELECT 
                price,
                EXTRACT(EPOCH FROM (COALESCE(effective_to, CURRENT_TIMESTAMP) - effective_from))/86400 as duration_days
            FROM price_history
            WHERE listing_id = r.id
            ORDER BY price DESC
        )
        SELECT MAX(price) INTO max_price
        FROM price_data
        WHERE duration_days >= min_price_duration;
        
        -- Если нашли максимальную цену и текущая цена ниже
        IF max_price IS NOT NULL AND current_price < max_price THEN
            -- Вычисляем процент скидки
            percentage := ((current_price - max_price) / max_price) * 100;
            
            -- Если скидка значительная (>= 5%)
            IF ABS(percentage) >= 5 THEN
                -- Обновляем метаданные
                SELECT COALESCE(metadata, '{}'::jsonb) INTO metadata_json FROM marketplace_listings WHERE id = r.id;
                
                metadata_json := jsonb_set(
                    metadata_json,
                    '{discount}',
                    jsonb_build_object(
                        'discount_percent', ROUND(ABS(percentage)),
                        'previous_price', max_price,
                        'effective_from', CURRENT_TIMESTAMP,
                        'has_price_history', true
                    )
                );
                
                UPDATE marketplace_listings
                SET metadata = metadata_json
                WHERE id = r.id;
                
                RAISE NOTICE 'Добавлена информация о скидке для объявления %d: %.2f -> %.2f (скидка %.0f%%)',
                    r.id, max_price, current_price, ABS(percentage);
            END IF;
        END IF;
    END LOOP;
END$$;

-- Создаем или обновляем функцию проверки манипуляций с ценой
CREATE OR REPLACE FUNCTION check_price_manipulation(p_listing_id INT)
RETURNS BOOLEAN AS $$
DECLARE
    price_changes RECORD;
    manipulation_detected BOOLEAN := FALSE;
BEGIN
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
    
    RETURN manipulation_detected;
END;
$$ LANGUAGE plpgsql;