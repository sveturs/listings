-- Migration: 000170_add_storefront_images_reindex_trigger
-- Description: Триггер для установки needs_reindex при изменении изображений товаров витрин
-- Author: Claude Code
-- Date: 2025-10-08

-- Создать функцию для установки needs_reindex при изменении изображений
CREATE OR REPLACE FUNCTION set_storefront_listing_needs_reindex()
RETURNS TRIGGER AS $$
DECLARE
    v_product_id INTEGER;
BEGIN
    -- Получить ID товара (из NEW или OLD в зависимости от операции)
    v_product_id := COALESCE(NEW.storefront_product_id, OLD.storefront_product_id);

    -- Установить needs_reindex=true для соответствующего listing
    UPDATE marketplace_listings
    SET needs_reindex = true,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = v_product_id
      AND metadata->>'source' = 'storefront';

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Создать триггер на storefront_product_images
CREATE TRIGGER trigger_storefront_images_reindex
    AFTER INSERT OR UPDATE OR DELETE
    ON storefront_product_images
    FOR EACH ROW
    EXECUTE FUNCTION set_storefront_listing_needs_reindex();

-- Комментарии
COMMENT ON FUNCTION set_storefront_listing_needs_reindex() IS
'Устанавливает needs_reindex=true для товара витрины при изменении его изображений.
Это гарантирует, что изменения в изображениях будут отражены в OpenSearch после переиндексации.';

COMMENT ON TRIGGER trigger_storefront_images_reindex ON storefront_product_images IS
'Триггер для установки needs_reindex при добавлении/изменении/удалении изображений товаров витрин';
