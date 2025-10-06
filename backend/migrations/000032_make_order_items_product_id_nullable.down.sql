-- Откат изменений storefront_order_items.product_id обратно на NOT NULL

-- 1. Удаляем новый foreign key constraint
ALTER TABLE storefront_order_items
    DROP CONSTRAINT IF EXISTS storefront_order_items_product_id_fkey;

-- 2. Удаляем записи где product_id = NULL (если есть)
-- ВНИМАНИЕ: это удалит исторические данные о заказах с удаленными товарами!
DELETE FROM storefront_order_items WHERE product_id IS NULL;

-- 3. Возвращаем NOT NULL constraint
ALTER TABLE storefront_order_items
    ALTER COLUMN product_id SET NOT NULL;

-- 4. Восстанавливаем оригинальный foreign key (без ON DELETE)
ALTER TABLE storefront_order_items
    ADD CONSTRAINT storefront_order_items_product_id_fkey
    FOREIGN KEY (product_id)
    REFERENCES storefront_products(id);

-- Удаляем комментарий
COMMENT ON COLUMN storefront_order_items.product_id IS NULL;
