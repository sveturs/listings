-- Migration: 000177_remove_storefront_sync_trigger (ROLLBACK)
-- Description: Восстанавливает триггер sync_storefront_to_marketplace (если нужен rollback)
-- Date: 2025-10-11

BEGIN;

-- ВНИМАНИЕ: Этот rollback восстановит триггер, но НЕ восстановит удаленные дубликаты!
-- Для полного восстановления нужно:
-- 1. Применить эту миграцию (восстановит триггер)
-- 2. Вручную запустить синхронизацию всех B2C товаров

-- Восстановить функцию триггера из миграции 000171
CREATE OR REPLACE FUNCTION sync_storefront_product_to_marketplace()
RETURNS TRIGGER AS $$
DECLARE
    v_user_id INTEGER;
    v_existing_listing_id INTEGER;
    v_new_status VARCHAR(20);
    v_storefront_slug VARCHAR(100);
BEGIN
    -- Получить user_id и slug владельца витрины
    SELECT user_id, slug INTO v_user_id, v_storefront_slug
    FROM b2c_stores
    WHERE id = COALESCE(NEW.storefront_id, OLD.storefront_id);

    -- Определить статус для c2c_listing
    IF TG_OP = 'DELETE' THEN
        v_new_status := 'inactive';
    ELSE
        v_new_status := CASE WHEN NEW.is_active THEN 'active' ELSE 'inactive' END;
    END IF;

    -- Обработка INSERT: создание нового listing
    IF TG_OP = 'INSERT' THEN
        SELECT id INTO v_existing_listing_id
        FROM c2c_listings
        WHERE id = NEW.id
        LIMIT 1;

        IF v_existing_listing_id IS NULL THEN
            INSERT INTO c2c_listings (
                id, user_id, category_id, title, description, price,
                condition, status, location, latitude, longitude, show_on_map,
                storefront_id, external_id, metadata, created_at, updated_at, needs_reindex
            ) VALUES (
                NEW.id, v_user_id, NEW.category_id, NEW.name, NEW.description, NEW.price,
                'new', v_new_status, COALESCE(NEW.individual_address, ''),
                NEW.individual_latitude, NEW.individual_longitude, COALESCE(NEW.show_on_map, true),
                NEW.storefront_id, NEW.sku,
                jsonb_build_object(
                    'source', 'storefront',
                    'storefront_id', NEW.storefront_id,
                    'stock_quantity', NEW.stock_quantity,
                    'stock_status', NEW.stock_status,
                    'currency', NEW.currency,
                    'barcode', NEW.barcode,
                    'attributes', NEW.attributes
                ),
                NEW.created_at, NEW.updated_at, true
            );
        END IF;

    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE c2c_listings SET
            title = NEW.name,
            description = NEW.description,
            price = NEW.price,
            status = v_new_status,
            category_id = NEW.category_id,
            updated_at = NEW.updated_at,
            needs_reindex = true
        WHERE id = NEW.id;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE c2c_listings SET
            status = 'inactive',
            updated_at = NOW(),
            needs_reindex = true
        WHERE id = OLD.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Восстановить триггер
CREATE TRIGGER trigger_sync_storefront_to_marketplace
    AFTER INSERT OR UPDATE OR DELETE
    ON b2c_products
    FOR EACH ROW
    EXECUTE FUNCTION sync_storefront_product_to_marketplace();

COMMENT ON FUNCTION sync_storefront_product_to_marketplace() IS
'Автоматически синхронизирует товары из b2c_products в c2c_listings.
ВНИМАНИЕ: Это костыль! Используется только для rollback миграции 000177.';

COMMIT;
