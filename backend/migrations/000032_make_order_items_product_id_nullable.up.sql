-- Изменяем storefront_order_items.product_id на nullable с ON DELETE SET NULL
-- Это позволяет сохранять историю заказов даже после удаления товара

-- 1. Удаляем старый foreign key constraint
ALTER TABLE storefront_order_items
    DROP CONSTRAINT IF EXISTS storefront_order_items_product_id_fkey;

-- 2. Делаем product_id nullable
ALTER TABLE storefront_order_items
    ALTER COLUMN product_id DROP NOT NULL;

-- 3. Создаём новый foreign key с ON DELETE SET NULL
ALTER TABLE storefront_order_items
    ADD CONSTRAINT storefront_order_items_product_id_fkey
    FOREIGN KEY (product_id)
    REFERENCES storefront_products(id)
    ON DELETE SET NULL;

-- Комментарий для документации
COMMENT ON COLUMN storefront_order_items.product_id IS 'Reference to product. Can be NULL if product was deleted. Product details are preserved in product_name, variant_name, product_sku columns.';
