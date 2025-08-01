-- Создание таблицы для резервирования товаров
CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    variant_id BIGINT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    order_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'reserved',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Внешние ключи
    CONSTRAINT fk_reservation_product FOREIGN KEY (product_id) 
        REFERENCES storefront_products(id) ON DELETE CASCADE,
    CONSTRAINT fk_reservation_variant FOREIGN KEY (variant_id) 
        REFERENCES storefront_product_variants(id) ON DELETE CASCADE,
    CONSTRAINT fk_reservation_order FOREIGN KEY (order_id) 
        REFERENCES storefront_orders(id) ON DELETE CASCADE
);

-- Индексы для быстрого поиска
CREATE INDEX idx_inventory_reservations_order_id ON inventory_reservations(order_id);
CREATE INDEX idx_inventory_reservations_product_id ON inventory_reservations(product_id);
CREATE INDEX idx_inventory_reservations_variant_id ON inventory_reservations(variant_id) WHERE variant_id IS NOT NULL;
CREATE INDEX idx_inventory_reservations_status ON inventory_reservations(status);
CREATE INDEX idx_inventory_reservations_expires_at ON inventory_reservations(expires_at) WHERE status = 'reserved';

-- Комментарии к таблице и колонкам
COMMENT ON TABLE inventory_reservations IS 'Резервирования товаров для заказов';
COMMENT ON COLUMN inventory_reservations.id IS 'Уникальный идентификатор резервирования';
COMMENT ON COLUMN inventory_reservations.product_id IS 'ID товара';
COMMENT ON COLUMN inventory_reservations.variant_id IS 'ID варианта товара (если применимо)';
COMMENT ON COLUMN inventory_reservations.quantity IS 'Зарезервированное количество';
COMMENT ON COLUMN inventory_reservations.order_id IS 'ID заказа';
COMMENT ON COLUMN inventory_reservations.status IS 'Статус резервирования: reserved, committed, released';
COMMENT ON COLUMN inventory_reservations.expires_at IS 'Время истечения резервирования';

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_inventory_reservations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для обновления updated_at
CREATE TRIGGER update_inventory_reservations_updated_at_trigger
    BEFORE UPDATE ON inventory_reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_inventory_reservations_updated_at();