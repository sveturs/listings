-- Создание таблицы для элементов заказа витрины
CREATE TABLE IF NOT EXISTS storefront_order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES storefront_orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES storefront_products(id),
    variant_id BIGINT REFERENCES storefront_product_variants(id),
    
    -- Информация о товаре на момент заказа
    product_name VARCHAR(255) NOT NULL,
    variant_name VARCHAR(255),
    product_sku VARCHAR(100),
    
    -- Цены и количество
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price NUMERIC(12,2) NOT NULL,
    total_price NUMERIC(12,2) NOT NULL,
    
    -- Статус и примечания
    status VARCHAR(50) DEFAULT 'pending',
    notes TEXT,
    
    -- Временные метки
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для производительности
CREATE INDEX idx_storefront_order_items_order_id ON storefront_order_items(order_id);
CREATE INDEX idx_storefront_order_items_product_id ON storefront_order_items(product_id);
CREATE INDEX idx_storefront_order_items_variant_id ON storefront_order_items(variant_id);

-- Триггер для обновления updated_at
CREATE TRIGGER update_storefront_order_items_updated_at
    BEFORE UPDATE ON storefront_order_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();