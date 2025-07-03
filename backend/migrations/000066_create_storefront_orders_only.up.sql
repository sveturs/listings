-- Migration: Create Storefront Orders System (without shopping carts)
-- Description: Создает систему заказов для витрин (корзины уже созданы в миграции 000065)

-- Таблица заказов
CREATE TABLE IF NOT EXISTS storefront_orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(32) UNIQUE NOT NULL, -- человекочитаемый номер заказа
    storefront_id INTEGER REFERENCES storefronts(id) ON DELETE RESTRICT,
    customer_id INTEGER REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Связь с платежной системой
    payment_transaction_id BIGINT REFERENCES payment_transactions(id) ON DELETE SET NULL,
    
    -- Финансовые данные (фиксируются на момент создания заказа)
    subtotal_amount DECIMAL(12,2) NOT NULL, -- сумма товаров
    shipping_amount DECIMAL(12,2) DEFAULT 0, -- стоимость доставки
    tax_amount DECIMAL(12,2) DEFAULT 0, -- налоги (если применимо)
    total_amount DECIMAL(12,2) NOT NULL, -- итоговая сумма
    commission_amount DECIMAL(12,2) NOT NULL, -- комиссия платформы
    seller_amount DECIMAL(12,2) NOT NULL, -- сумма к выплате продавцу
    
    currency CHAR(3) DEFAULT 'RSD',
    
    -- Статус заказа
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- Возможные статусы: pending, confirmed, processing, shipped, delivered, cancelled, refunded
    
    -- Escrow параметры
    escrow_release_date DATE, -- дата автоматического освобождения средств
    escrow_days INTEGER DEFAULT 3, -- количество дней удержания (3-7 для витрин)
    
    -- Информация о доставке
    shipping_address JSONB, -- адрес доставки
    billing_address JSONB, -- адрес оплаты
    shipping_method VARCHAR(100), -- способ доставки
    shipping_provider VARCHAR(50), -- провайдер доставки
    tracking_number VARCHAR(100), -- номер отслеживания
    
    -- Контактная информация
    customer_notes TEXT, -- комментарии покупателя
    seller_notes TEXT, -- заметки продавца
    
    -- Временные метки
    confirmed_at TIMESTAMP, -- когда заказ подтвержден
    shipped_at TIMESTAMP, -- когда отправлен
    delivered_at TIMESTAMP, -- когда доставлен
    cancelled_at TIMESTAMP, -- когда отменен
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Позиции заказа
CREATE TABLE IF NOT EXISTS storefront_order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES storefront_orders(id) ON DELETE CASCADE,
    product_id BIGINT REFERENCES storefront_products(id) ON DELETE RESTRICT,
    variant_id BIGINT REFERENCES storefront_product_variants(id) ON DELETE SET NULL,
    
    -- Фиксированная информация о товаре на момент заказа
    product_name VARCHAR(255) NOT NULL, -- название товара
    product_sku VARCHAR(100), -- артикул
    variant_name VARCHAR(255), -- название варианта
    
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_per_unit DECIMAL(12,2) NOT NULL, -- цена за единицу
    total_price DECIMAL(12,2) NOT NULL, -- общая стоимость позиции
    
    -- Snapshot атрибутов товара на момент заказа
    product_attributes JSONB, -- сохраняем все атрибуты товара
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Резервирование товаров (для управления остатками)
CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES storefront_orders(id) ON DELETE CASCADE,
    product_id BIGINT REFERENCES storefront_products(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES storefront_product_variants(id) ON DELETE SET NULL,
    
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    
    -- Статус резерва
    status VARCHAR(20) DEFAULT 'active', -- active, committed, released, expired
    
    -- Автоматическое освобождение через X часов если заказ не оплачен
    expires_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '2 hours'),
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    released_at TIMESTAMP -- когда резерв освобожден
);

-- Индексы для inventory_reservations
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_product ON inventory_reservations(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_reservations_expires ON inventory_reservations(expires_at) WHERE status = 'active';

-- Расширяем payment_transactions для поддержки заказов
ALTER TABLE payment_transactions 
ADD COLUMN IF NOT EXISTS source_type VARCHAR(20) DEFAULT 'marketplace_listing', -- 'marketplace_listing' | 'storefront_order'
ADD COLUMN IF NOT EXISTS source_id BIGINT, -- listing_id или order_id
ADD COLUMN IF NOT EXISTS storefront_id INTEGER REFERENCES storefronts(id);

-- Обновляем существующие записи
UPDATE payment_transactions 
SET source_type = 'marketplace_listing', source_id = listing_id 
WHERE listing_id IS NOT NULL AND source_type IS NULL;

-- Создаем индексы для производительности
CREATE INDEX IF NOT EXISTS idx_payment_transactions_source ON payment_transactions(source_type, source_id);
CREATE INDEX IF NOT EXISTS idx_storefront_orders_customer ON storefront_orders(customer_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storefront_orders_storefront ON storefront_orders(storefront_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storefront_orders_status ON storefront_orders(status);
CREATE INDEX IF NOT EXISTS idx_storefront_orders_escrow_date ON storefront_orders(escrow_release_date) WHERE escrow_release_date IS NOT NULL;

-- Триггеры для автоматического обновления timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_storefront_orders_updated_at BEFORE UPDATE ON storefront_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Функция для генерации номеров заказов
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS VARCHAR(32) AS $$
DECLARE
    order_num VARCHAR(32);
    counter INTEGER := 0;
BEGIN
    LOOP
        -- Генерируем номер в формате: STO-YYYYMMDD-XXXXX
        order_num := 'STO-' || TO_CHAR(CURRENT_DATE, 'YYYYMMDD') || '-' || 
                    LPAD((EXTRACT(epoch FROM CURRENT_TIMESTAMP)::INTEGER % 100000)::TEXT, 5, '0');
        
        -- Проверяем уникальность
        IF NOT EXISTS (SELECT 1 FROM storefront_orders WHERE order_number = order_num) THEN
            RETURN order_num;
        END IF;
        
        -- Защита от бесконечного цикла
        counter := counter + 1;
        IF counter > 1000 THEN
            RAISE EXCEPTION 'Unable to generate unique order number after 1000 attempts';
        END IF;
        
        -- Небольшая задержка перед следующей попыткой
        PERFORM pg_sleep(0.001);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Триггер для автоматической генерации номера заказа
CREATE OR REPLACE FUNCTION set_order_number()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.order_number IS NULL OR NEW.order_number = '' THEN
        NEW.order_number := generate_order_number();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_order_number_trigger 
    BEFORE INSERT ON storefront_orders 
    FOR EACH ROW EXECUTE FUNCTION set_order_number();

-- Функция для автоматического расчета escrow_release_date
CREATE OR REPLACE FUNCTION calculate_escrow_release_date()
RETURNS TRIGGER AS $$
BEGIN
    -- Рассчитываем дату освобождения на основе escrow_days
    IF NEW.escrow_release_date IS NULL AND NEW.escrow_days IS NOT NULL THEN
        NEW.escrow_release_date := CURRENT_DATE + INTERVAL '1 day' * NEW.escrow_days;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER calculate_escrow_release_date_trigger 
    BEFORE INSERT OR UPDATE ON storefront_orders 
    FOR EACH ROW EXECUTE FUNCTION calculate_escrow_release_date();

-- Добавляем комментарии к таблицам
COMMENT ON TABLE storefront_orders IS 'Заказы в витринах с полной информацией о покупке';
COMMENT ON TABLE storefront_order_items IS 'Позиции заказов с фиксированной информацией о товарах';
COMMENT ON TABLE inventory_reservations IS 'Резервирование товаров на время оформления заказа';

COMMENT ON COLUMN storefront_orders.escrow_days IS 'Количество дней удержания средств (3-7 для витрин)';
COMMENT ON COLUMN storefront_orders.commission_amount IS 'Комиссия платформы рассчитанная по тарифному плану витрины';
COMMENT ON COLUMN payment_transactions.source_type IS 'Тип источника платежа: marketplace_listing или storefront_order';