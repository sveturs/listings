-- Создание таблицы для отслеживания истории цен
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    listing_id INT NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    price DECIMAL(12,2) NOT NULL,
    effective_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    effective_to TIMESTAMP,
    change_source VARCHAR(50) NOT NULL, -- manual, import, system, etc.
    change_percentage DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска истории цен конкретного объявления
CREATE INDEX IF NOT EXISTS idx_price_history_listing_id ON price_history(listing_id);

-- Индекс для быстрого поиска актуальных цен (где effective_to IS NULL)
CREATE INDEX IF NOT EXISTS idx_price_history_effective ON price_history(listing_id, effective_to);

-- Добавление поля metadata к таблице marketplace_listings если его еще нет
ALTER TABLE marketplace_listings ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Создаем индекс для быстрого поиска товаров со скидками
CREATE INDEX IF NOT EXISTS idx_listings_metadata_discount ON marketplace_listings USING GIN ((metadata -> 'discount'));

-- Удаляем функции и триггеры, если они существуют
DROP TRIGGER IF EXISTS trig_update_metadata_after_price_change ON price_history;
DROP TRIGGER IF EXISTS trg_new_listing_price_history ON marketplace_listings;
DROP TRIGGER IF EXISTS trg_update_listing_price_history ON marketplace_listings;
DROP FUNCTION IF EXISTS update_listing_metadata_after_price_change();
DROP FUNCTION IF EXISTS update_price_history();
DROP FUNCTION IF EXISTS check_price_manipulation(INT);
DROP FUNCTION IF EXISTS refresh_discount_metadata();

-- Функция для проверки манипуляций с ценами
CREATE OR REPLACE FUNCTION check_price_manipulation(p_listing_id INT)
RETURNS BOOLEAN AS $$
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
$$ LANGUAGE plpgsql;

-- Создаем триггерную функцию для автоматического обновления истории цен
CREATE OR REPLACE FUNCTION update_price_history()
RETURNS TRIGGER AS $$
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
$$ LANGUAGE plpgsql;

-- Функция для обновления метаданных о скидках после изменения цены
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
$$ LANGUAGE plpgsql;

-- Функция для ручного обновления метаданных скидок
CREATE OR REPLACE FUNCTION refresh_discount_metadata()
RETURNS void AS $$
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
$$ LANGUAGE plpgsql;

-- Создаем триггеры
DROP TRIGGER IF EXISTS trg_new_listing_price_history ON marketplace_listings;
CREATE TRIGGER trg_new_listing_price_history
AFTER INSERT ON marketplace_listings
FOR EACH ROW
EXECUTE FUNCTION update_price_history('create');

DROP TRIGGER IF EXISTS trg_update_listing_price_history ON marketplace_listings;
CREATE TRIGGER trg_update_listing_price_history
AFTER UPDATE OF price ON marketplace_listings
FOR EACH ROW
WHEN (OLD.price IS DISTINCT FROM NEW.price)
EXECUTE FUNCTION update_price_history('update');

DROP TRIGGER IF EXISTS trig_update_metadata_after_price_change ON price_history;
CREATE TRIGGER trig_update_metadata_after_price_change
AFTER INSERT ON price_history
FOR EACH ROW
EXECUTE FUNCTION update_listing_metadata_after_price_change();

-- Вставляем запись в schema_migrations
INSERT INTO schema_migrations (version, dirty) VALUES (13, false) ON CONFLICT (version) DO NOTHING;