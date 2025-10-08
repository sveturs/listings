-- Migration: 000170_add_storefront_images_reindex_trigger (DOWN)
-- Description: Откат триггера для установки needs_reindex при изменении изображений
-- Author: Claude Code
-- Date: 2025-10-08

-- Удалить триггер
DROP TRIGGER IF EXISTS trigger_storefront_images_reindex ON storefront_product_images;

-- Удалить функцию
DROP FUNCTION IF EXISTS set_storefront_listing_needs_reindex();
