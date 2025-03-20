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

-- Создаем триггер для новых объявлений
CREATE TRIGGER trg_new_listing_price_history
AFTER INSERT ON marketplace_listings
FOR EACH ROW
EXECUTE FUNCTION update_price_history('create');

-- Создаем триггер для обновления объявлений
CREATE TRIGGER trg_update_listing_price_history
AFTER UPDATE OF price ON marketplace_listings
FOR EACH ROW
WHEN (OLD.price IS DISTINCT FROM NEW.price)
EXECUTE FUNCTION update_price_history('update');

-- Функция для проверки манипуляций с ценами
CREATE OR REPLACE FUNCTION check_price_manipulation(p_listing_id INT)
RETURNS BOOLEAN AS $$
DECLARE
    price_increases INT := 0;
    price_decreases INT := 0;
    last_action VARCHAR(10) := NULL;
    last_price DECIMAL(12,2) := NULL;
    current_price DECIMAL(12,2);
    days_since_last_increase INT := 0;
    manipulation_detected BOOLEAN := FALSE;
    r RECORD;
BEGIN
    -- Получаем историю цен за последние 30 дней, отсортированную по времени
    FOR r IN (
        SELECT 
            price,
            effective_from,
            CASE 
                WHEN LAG(price) OVER (ORDER BY effective_from) IS NULL THEN 'none'
                WHEN price > LAG(price) OVER (ORDER BY effective_from) THEN 'increase'
                WHEN price < LAG(price) OVER (ORDER BY effective_from) THEN 'decrease'
                ELSE 'same'
            END as action
        FROM price_history
        WHERE listing_id = p_listing_id
        AND effective_from > CURRENT_TIMESTAMP - INTERVAL '30 days'
        ORDER BY effective_from
    ) LOOP
        -- Подозрительные паттерны:
        -- 1. Увеличение цены с последующим быстрым снижением (менее 7 дней)
        IF r.action = 'increase' THEN
            price_increases := price_increases + 1;
            last_action := 'increase';
            last_price := r.price;
        ELSIF r.action = 'decrease' AND last_action = 'increase' THEN
            price_decreases := price_decreases + 1;
            days_since_last_increase := EXTRACT(EPOCH FROM (r.effective_from - last_price)) / 86400;
            
            -- Если цена была увеличена и снижена в течение 7 дней, это подозрительно
            IF days_since_last_increase < 7 THEN
                manipulation_detected := TRUE;
            END IF;
            
            last_action := 'decrease';
        END IF;
        
        current_price := r.price;
    END LOOP;
    
    -- Если обнаружено слишком много циклов увеличение-снижение в течение месяца
    IF price_increases > 2 AND price_decreases > 2 THEN
        manipulation_detected := TRUE;
    END IF;
    
    RETURN manipulation_detected;
END;
$$ LANGUAGE plpgsql;