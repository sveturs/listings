-- Добавление поля metadata к таблице marketplace_listings
ALTER TABLE marketplace_listings ADD COLUMN IF NOT EXISTS metadata JSONB;

-- Создаем индекс для быстрого поиска товаров со скидками
CREATE INDEX IF NOT EXISTS idx_listings_metadata_discount ON marketplace_listings USING GIN ((metadata -> 'discount'));

-- Обновляем триггерную функцию для обработки изменений цены
CREATE OR REPLACE FUNCTION update_price_history()
RETURNS TRIGGER AS $$
DECLARE
    last_price DECIMAL(12,2);
    price_diff DECIMAL(10,2);
    percentage DECIMAL(10,2);
        current_timestamp_var TIMESTAMP := CURRENT_TIMESTAMP AT TIME ZONE 'UTC';

    metadata_json JSONB;
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
            
            -- Если цена снизилась существенно (10% и более), добавляем информацию о скидке в метаданные
            IF price_diff < 0 AND percentage <= -10 THEN
                metadata_json := COALESCE(NEW.metadata, '{}'::jsonb);
                
                -- Проверяем на подозрительные манипуляции с ценой
                IF NOT check_price_manipulation(NEW.id) THEN
                    -- Добавляем информацию о скидке в метаданные
                    metadata_json := jsonb_set(
                        metadata_json, 
                        '{discount}', 
                        jsonb_build_object(
                            'discount_percent', -1 * ROUND(percentage),
                            'previous_price', last_price,
                            'effective_from', current_timestamp_var,
                            'has_price_history', true
                        )
                    );
                    
                    -- Обновляем метаданные объявления
                    NEW.metadata := metadata_json;
                END IF;
            ELSIF price_diff >= 0 THEN
                -- Если цена выросла или не изменилась, удаляем информацию о скидке из метаданных
                metadata_json := COALESCE(NEW.metadata, '{}'::jsonb);
                
                IF metadata_json ? 'discount' THEN
                    metadata_json := metadata_json - 'discount';
                    NEW.metadata := metadata_json;
                END IF;
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