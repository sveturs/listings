-- Migration rollback: возврат к старой архитектуре (с багом)
-- ВНИМАНИЕ: Этот откат восстанавливает неправильную архитектуру!

-- Шаг 1: Удалить правильный foreign key constraint на storefronts
ALTER TABLE marketplace_listings
DROP CONSTRAINT IF EXISTS marketplace_listings_storefront_id_fkey;

-- Шаг 2: Восстановить связь через user_storefronts
-- Конвертировать storefront_id обратно в user_storefront_id
UPDATE marketplace_listings ml
SET storefront_id = us.id
FROM storefronts s
JOIN user_storefronts us ON us.user_id = s.user_id AND us.slug LIKE '%' || s.slug || '%'
WHERE ml.storefront_id = s.id
  AND ml.metadata->>'source' = 'storefront';

-- Шаг 3: Создать старый foreign key constraint на user_storefronts
ALTER TABLE marketplace_listings
ADD CONSTRAINT marketplace_listings_storefront_id_fkey
FOREIGN KEY (storefront_id) REFERENCES user_storefronts(id) ON DELETE SET NULL;

-- Шаг 4: Восстановить старую версию функции sync_storefront_product_to_marketplace
CREATE OR REPLACE FUNCTION public.sync_storefront_product_to_marketplace()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
    v_user_id INTEGER;
    v_user_storefront_id INTEGER;
    v_existing_listing_id INTEGER;
    v_new_status VARCHAR(20);
    v_storefront_slug VARCHAR(100);
BEGIN
    -- Получить user_id и slug владельца витрины
    SELECT user_id, slug INTO v_user_id, v_storefront_slug
    FROM storefronts
    WHERE id = COALESCE(NEW.storefront_id, OLD.storefront_id);

    -- Получить user_storefront_id из user_storefronts по user_id и slug
    SELECT id INTO v_user_storefront_id
    FROM user_storefronts
    WHERE user_id = v_user_id
      AND slug LIKE '%' || v_storefront_slug || '%'
    ORDER BY id DESC
    LIMIT 1;

    -- Fallback: берём первый user_storefront для пользователя
    IF v_user_storefront_id IS NULL THEN
        SELECT id INTO v_user_storefront_id
        FROM user_storefronts
        WHERE user_id = v_user_id
        LIMIT 1;
    END IF;

    -- Определить статус для marketplace_listing
    IF TG_OP = 'DELETE' THEN
        v_new_status := 'inactive';
    ELSE
        v_new_status := CASE WHEN NEW.is_active THEN 'active' ELSE 'inactive' END;
    END IF;

    -- Обработка INSERT: создание нового listing
    IF TG_OP = 'INSERT' THEN
        SELECT id INTO v_existing_listing_id
        FROM marketplace_listings
        WHERE id = NEW.id
        LIMIT 1;

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
                v_user_storefront_id,  -- СТАРАЯ ВЕРСИЯ: используем user_storefronts.id
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
        WHERE id = NEW.id;

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
Версия 2: исправлен поиск storefront_id, добавлено создание translations для Serbian.
При создании/обновлении товара создается/обновляется соответствующий listing и translations.
При удалении товара listing деактивируется (не удаляется) для сохранения связей.';
