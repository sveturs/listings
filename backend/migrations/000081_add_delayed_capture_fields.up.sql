-- Добавляем поля для управления delayed capture
ALTER TABLE payment_transactions 
ADD COLUMN IF NOT EXISTS capture_mode VARCHAR(20) DEFAULT 'manual' CHECK (capture_mode IN ('auto', 'manual')),
ADD COLUMN IF NOT EXISTS auto_capture_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS capture_deadline_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS capture_attempted_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS capture_attempts INT DEFAULT 0;

-- Индекс для поиска транзакций требующих auto-capture
CREATE INDEX IF NOT EXISTS idx_payment_transactions_auto_capture 
ON payment_transactions(auto_capture_at, status) 
WHERE status = 'authorized' AND capture_mode = 'auto';

-- Добавляем настройки в payment_gateways (если не JSONB)
ALTER TABLE payment_gateways 
ALTER COLUMN config TYPE JSONB USING config::jsonb;

-- Создаем таблицу orders для связи платежей с покупками
CREATE TABLE IF NOT EXISTS marketplace_orders (
    id SERIAL PRIMARY KEY,
    
    -- Участники сделки
    buyer_id INT NOT NULL REFERENCES users(id),
    seller_id INT NOT NULL REFERENCES users(id),
    listing_id INT NOT NULL REFERENCES marketplace_listings(id),
    
    -- Финансы
    item_price DECIMAL(10,2) NOT NULL,
    platform_fee_rate DECIMAL(5,2) DEFAULT 5.00,
    platform_fee_amount DECIMAL(10,2) NOT NULL,
    seller_payout_amount DECIMAL(10,2) NOT NULL,
    
    -- Связь с платежом
    payment_transaction_id INT REFERENCES payment_transactions(id),
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN (
        'pending',           -- Ожидает оплаты
        'paid',              -- Оплачен, ожидает отправки
        'shipped',           -- Отправлен
        'delivered',         -- Доставлен, ожидает подтверждения
        'completed',         -- Завершен, средства переведены продавцу
        'disputed',          -- Открыт спор
        'cancelled',         -- Отменен
        'refunded'           -- Возвращен
    )),
    
    -- Защитный период
    protection_period_days INT DEFAULT 7,
    protection_expires_at TIMESTAMP WITH TIME ZONE,
    
    -- Доставка
    shipping_method VARCHAR(100),
    tracking_number VARCHAR(255),
    shipped_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    
    -- Метаданные
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для orders
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_buyer ON marketplace_orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_seller ON marketplace_orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_listing ON marketplace_orders(listing_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_status ON marketplace_orders(status);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_protection ON marketplace_orders(protection_expires_at) 
WHERE status IN ('delivered', 'shipped');

-- Таблица для истории изменения статусов
CREATE TABLE IF NOT EXISTS order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES marketplace_orders(id),
    old_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    reason TEXT,
    created_by INT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица для сообщений в заказе (не чат, а структурированные сообщения)
CREATE TABLE IF NOT EXISTS order_messages (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES marketplace_orders(id),
    sender_id INT NOT NULL REFERENCES users(id),
    message_type VARCHAR(50) DEFAULT 'text' CHECK (message_type IN (
        'text',              -- Обычное сообщение
        'shipping_update',   -- Обновление доставки
        'dispute_opened',    -- Открыт спор
        'dispute_message',   -- Сообщение в споре
        'system'             -- Системное уведомление
    )),
    content TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS update_marketplace_orders_updated_at ON marketplace_orders;
CREATE TRIGGER update_marketplace_orders_updated_at BEFORE UPDATE
ON marketplace_orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();