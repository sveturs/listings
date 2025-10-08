-- Migration: 000169_fix_storefront_sync_trigger
-- Description: Исправляет триггер синхронизации: правильный storefront_id и создание translations
-- Author: Claude Code
-- Date: 2025-10-08

-- Удалить старый триггер и функцию
DROP TRIGGER IF EXISTS trigger_sync_storefront_to_marketplace ON storefront_products;
DROP FUNCTION IF EXISTS sync_storefront_product_to_marketplace();

-- Создать улучшенную функцию синхронизации
CREATE OR REPLACE FUNCTION sync_storefront_product_to_marketplace()
RETURNS TRIGGER AS $$
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
    -- Если не найден, берём первый для пользователя (backward compatibility)
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
                v_user_storefront_id,  -- user_storefronts.id для обратной совместимости
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

            -- Создать translations
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
        END IF;

    -- Обработка DELETE: деактивация listing
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE marketplace_listings
        SET
            status = 'inactive',
            updated_at = CURRENT_TIMESTAMP,
            needs_reindex = true
        WHERE id = OLD.id;
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
Версия 2: исправлен поиск storefront_id, добавлено создание translations для Serbian.
При создании/обновлении товара создается/обновляется соответствующий listing и translations.
При удалении товара listing деактивируется (не удаляется) для сохранения связей.';

COMMENT ON TRIGGER trigger_sync_storefront_to_marketplace ON storefront_products IS
'Триггер для автоматической синхронизации товаров витрин с маркетплейсом (v2 с translations)';
