-- Migration 000082: Создание таблиц для заказов маркетплейса
-- Эта миграция создает систему заказов специально для маркетплейса с поддержкой delayed capture

-- Таблица заказов маркетплейса
CREATE TABLE IF NOT EXISTS marketplace_orders (
    id                    BIGSERIAL PRIMARY KEY,
    buyer_id              BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    seller_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id            BIGINT NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    
    -- Финансовые данные
    item_price            DECIMAL(10,2) NOT NULL CHECK (item_price > 0),
    platform_fee_rate     DECIMAL(5,2) NOT NULL DEFAULT 5.0 CHECK (platform_fee_rate >= 0 AND platform_fee_rate <= 100),
    platform_fee_amount   DECIMAL(10,2) NOT NULL CHECK (platform_fee_amount >= 0),
    seller_payout_amount  DECIMAL(10,2) NOT NULL CHECK (seller_payout_amount >= 0),
    
    -- Связь с платежной системой
    payment_transaction_id BIGINT REFERENCES payment_transactions(id),
    
    -- Статус заказа
    status                TEXT NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'paid', 'shipped', 'delivered', 'completed', 'disputed', 'cancelled', 'refunded')),
    
    -- Защитный период
    protection_period_days INTEGER NOT NULL DEFAULT 7 CHECK (protection_period_days >= 0),
    protection_expires_at  TIMESTAMP WITH TIME ZONE,
    
    -- Доставка
    shipping_method       TEXT,
    tracking_number       TEXT,
    shipped_at           TIMESTAMP WITH TIME ZONE,
    delivered_at         TIMESTAMP WITH TIME ZONE,
    
    -- Аудит
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at           TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ограничения
    CONSTRAINT marketplace_orders_different_users CHECK (buyer_id != seller_id),
    CONSTRAINT marketplace_orders_consistent_fees CHECK (platform_fee_amount + seller_payout_amount = item_price)
);

-- Таблица истории изменения статусов заказов
CREATE TABLE IF NOT EXISTS order_status_history (
    id          BIGSERIAL PRIMARY KEY,
    order_id    BIGINT NOT NULL REFERENCES marketplace_orders(id) ON DELETE CASCADE,
    old_status  TEXT,
    new_status  TEXT NOT NULL,
    reason      TEXT,
    created_by  BIGINT REFERENCES users(id),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица сообщений в заказах
CREATE TABLE IF NOT EXISTS order_messages (
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT NOT NULL REFERENCES marketplace_orders(id) ON DELETE CASCADE,
    sender_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_type TEXT NOT NULL DEFAULT 'text' 
        CHECK (message_type IN ('text', 'shipping_update', 'dispute_opened', 'dispute_message', 'system')),
    content      TEXT NOT NULL,
    metadata     JSONB DEFAULT '{}',
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для производительности
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_buyer_id ON marketplace_orders(buyer_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_seller_id ON marketplace_orders(seller_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_listing_id ON marketplace_orders(listing_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_status ON marketplace_orders(status);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_payment_transaction_id ON marketplace_orders(payment_transaction_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_created_at ON marketplace_orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_marketplace_orders_protection_expires_at ON marketplace_orders(protection_expires_at) WHERE protection_expires_at IS NOT NULL;

-- Индексы для истории статусов
CREATE INDEX IF NOT EXISTS idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX IF NOT EXISTS idx_order_status_history_created_at ON order_status_history(created_at DESC);

-- Индексы для сообщений
CREATE INDEX IF NOT EXISTS idx_order_messages_order_id ON order_messages(order_id);
CREATE INDEX IF NOT EXISTS idx_order_messages_sender_id ON order_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_order_messages_created_at ON order_messages(created_at DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_marketplace_orders_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS marketplace_orders_updated_at_trigger ON marketplace_orders;
CREATE TRIGGER marketplace_orders_updated_at_trigger
    BEFORE UPDATE ON marketplace_orders
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_orders_updated_at();

-- Комментарии для документации
COMMENT ON TABLE marketplace_orders IS 'Заказы в маркетплейсе с поддержкой delayed capture';
COMMENT ON COLUMN marketplace_orders.protection_period_days IS 'Количество дней защитного периода после доставки';
COMMENT ON COLUMN marketplace_orders.protection_expires_at IS 'Дата окончания защитного периода, после которой можно захватить платеж';
COMMENT ON TABLE order_status_history IS 'История изменений статусов заказов для аудита';
COMMENT ON TABLE order_messages IS 'Сообщения между покупателем и продавцом в рамках заказа';