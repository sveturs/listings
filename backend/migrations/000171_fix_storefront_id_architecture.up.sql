-- Migration: Fix storefront_id architecture
-- Проблема 1: marketplace_listings.storefront_id ссылается на user_storefronts.id (неправильно!)
-- Проблема 2: Функция sync_storefront_product_to_marketplace использовала user_storefronts.id
-- Решение: storefront_id должен ссылаться на storefronts.id

-- Шаг 1: Удалить неправильный foreign key constraint
ALTER TABLE marketplace_listings
DROP CONSTRAINT IF EXISTS marketplace_listings_storefront_id_fkey;

-- Шаг 2: Исправить storefront_id в существующих marketplace_listings
-- Берем правильный storefront_id из metadata.storefront_id
UPDATE marketplace_listings ml
SET storefront_id = (ml.metadata->>'storefront_id')::INTEGER
WHERE ml.metadata->>'source' = 'storefront'
  AND ml.metadata->>'storefront_id' IS NOT NULL
  AND ml.storefront_id != (ml.metadata->>'storefront_id')::INTEGER;

-- Шаг 3: Создать правильный foreign key constraint на storefronts
ALTER TABLE marketplace_listings
ADD CONSTRAINT marketplace_listings_storefront_id_fkey
FOREIGN KEY (storefront_id) REFERENCES storefronts(id) ON DELETE SET NULL;

-- Шаг 4: Пересоздать функцию sync_storefront_product_to_marketplace с исправлением
CREATE OR REPLACE FUNCTION public.sync_storefront_product_to_marketplace()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_user_id INTEGER;
    v_existing_listing_id INTEGER;
    v_new_status VARCHAR(20);
    v_storefront_slug VARCHAR(100);
BEGIN
    -- Получить user_id и slug владельца витрины
    SELECT user_id, slug INTO v_user_id, v_storefront_slug
    FROM storefronts
    WHERE id = COALESCE(NEW.storefront_id, OLD.storefront_id);

    -- Определить статус для marketplace_listing
    IF TG_OP = 'DELETE' THEN
        v_new_status := 'inactive';
    ELSE
        v_new_status := CASE WHEN NEW.is_active THEN 'active' ELSE 'inactive' END;
    END IF;

    -- Обработка INSERT: создание нового listing
    IF TG_OP = 'INSERT' THEN
        -- Проверить, не существует ли уже listing с таким же ID
        SELECT id INTO v_existing_listing_id
        FROM marketplace_listings
        WHERE id = NEW.id
        LIMIT 1;

        -- Если listing не существует, создать новый
        IF v_existing_listing_id IS NULL THEN
            -- Создать marketplace_listing
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
                NEW.storefront_id,  -- ИСПРАВЛЕНО: используем storefronts.id напрямую
                NEW.sku,  -- SKU как external_id для отслеживания
                jsonb_build_object(
                    'source', 'storefront',
                    'storefront_id', NEW.storefront_id,  -- storefronts.id в metadata
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

            -- Создать translations для Serbian (по умолчанию)
            INSERT INTO translations (
                entity_type,
                entity_id,
                language,
                field_name,
                translated_text,
                is_machine_translated,
                is_verified,
                created_at,
                updated_at
            ) VALUES
                ('listing', NEW.id, 'sr', 'title', NEW.name, false, true, NEW.created_at, NEW.updated_at),
                ('listing', NEW.id, 'sr', 'description', NEW.description, false, true, NEW.created_at, NEW.updated_at)
            ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE SET
                translated_text = EXCLUDED.translated_text,
                updated_at = EXCLUDED.updated_at;

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
                storefront_id = NEW.storefront_id,  -- ИСПРАВЛЕНО: обновляем правильный storefront_id
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

            -- Обновить translations
            INSERT INTO translations (
                entity_type,
                entity_id,
                language,
                field_name,
                translated_text,
                is_machine_translated,
                is_verified,
                updated_at
            ) VALUES
                ('listing', NEW.id, 'sr', 'title', NEW.name, false, true, NEW.updated_at),
                ('listing', NEW.id, 'sr', 'description', NEW.description, false, true, NEW.updated_at)
            ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE SET
                translated_text = EXCLUDED.translated_text,
                updated_at = EXCLUDED.updated_at;
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
            storefront_id = NEW.storefront_id,  -- ИСПРАВЛЕНО: обновляем правильный storefront_id
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
        WHERE id = NEW.id;

        -- Обновить translations (или создать если не существуют)
        INSERT INTO translations (
            entity_type,
            entity_id,
            language,
            field_name,
            translated_text,
            is_machine_translated,
            is_verified,
            updated_at
        ) VALUES
            ('listing', NEW.id, 'sr', 'title', NEW.name, false, true, NEW.updated_at),
            ('listing', NEW.id, 'sr', 'description', NEW.description, false, true, NEW.updated_at)
        ON CONFLICT (entity_type, entity_id, language, field_name) DO UPDATE SET
            translated_text = EXCLUDED.translated_text,
            updated_at = EXCLUDED.updated_at;

    -- Обработка DELETE: деактивация listing
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE marketplace_listings
        SET
            status = 'inactive',
            updated_at = NOW(),
            needs_reindex = true
        WHERE id = OLD.id;
    END IF;

    RETURN NEW;
END;
$function$;

COMMENT ON FUNCTION sync_storefront_product_to_marketplace() IS
'Автоматически синхронизирует товары из storefront_products в marketplace_listings.
Версия 3: ИСПРАВЛЕНА архитектура - storefront_id теперь ссылается на storefronts.id (а не user_storefronts.id).
При создании/обновлении товара создается/обновляется соответствующий listing и translations.
При удалении товара listing деактивируется (не удаляется) для сохранения связей.';
