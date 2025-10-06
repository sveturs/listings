-- Migration: 000028_sync_storefront_to_marketplace
-- Description: Создает триггер для автоматической синхронизации товаров из storefront_products в marketplace_listings
-- Author: Claude Code
-- Date: 2025-10-06

-- Функция для синхронизации товара витрины с маркетплейсом
CREATE OR REPLACE FUNCTION sync_storefront_product_to_marketplace()
RETURNS TRIGGER AS $$
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
$$ LANGUAGE plpgsql;

-- Создать триггер на storefront_products
CREATE TRIGGER trigger_sync_storefront_to_marketplace
    AFTER INSERT OR UPDATE OR DELETE
    ON storefront_products
    FOR EACH ROW
    EXECUTE FUNCTION sync_storefront_product_to_marketplace();

-- Комментарии
COMMENT ON FUNCTION sync_storefront_product_to_marketplace() IS
'Автоматически синхронизирует товары из storefront_products в marketplace_listings.
При создании/обновлении товара создается/обновляется соответствующий listing.
При удалении товара listing деактивируется (не удаляется) для сохранения связей.';

COMMENT ON TRIGGER trigger_sync_storefront_to_marketplace ON storefront_products IS
'Триггер для автоматической синхронизации товаров витрин с маркетплейсом';
