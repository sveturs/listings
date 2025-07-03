-- Создание таблиц для интеграции с AllSecure Payment Gateway

-- Таблица конфигурации платежных шлюзов
CREATE TABLE payment_gateways (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, -- 'allsecure'
    is_active BOOLEAN DEFAULT true,
    config JSONB NOT NULL, -- API credentials, endpoints
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица платежных транзакций
CREATE TABLE payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    gateway_id INT REFERENCES payment_gateways(id),
    user_id INT REFERENCES users(id),
    
    -- Ссылка на покупку
    listing_id INT REFERENCES marketplace_listings(id),
    order_reference VARCHAR(255) UNIQUE NOT NULL,
    
    -- AllSecure данные
    gateway_transaction_id VARCHAR(255),
    gateway_reference_id VARCHAR(255),
    
    -- Финансовые данные
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    marketplace_commission DECIMAL(12,2),
    seller_amount DECIMAL(12,2),
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'pending', -- pending, authorized, captured, failed, refunded, voided
    gateway_status VARCHAR(50),
    
    -- Дополнительная информация
    payment_method VARCHAR(50), -- card, bank_transfer, etc.
    customer_email VARCHAR(255),
    description TEXT,
    
    -- Metadata
    gateway_response JSONB,
    error_details JSONB,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    authorized_at TIMESTAMP WITH TIME ZONE,
    captured_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT payment_transactions_amount_positive CHECK (amount > 0),
    CONSTRAINT payment_transactions_status_valid CHECK (
        status IN ('pending', 'authorized', 'captured', 'failed', 'refunded', 'voided')
    )
);

-- Таблица escrow платежей для marketplace
CREATE TABLE escrow_payments (
    id BIGSERIAL PRIMARY KEY,
    payment_transaction_id BIGINT REFERENCES payment_transactions(id) ON DELETE CASCADE,
    seller_id INT REFERENCES users(id),
    buyer_id INT REFERENCES users(id),
    listing_id INT REFERENCES marketplace_listings(id),
    
    amount DECIMAL(12,2) NOT NULL,
    marketplace_commission DECIMAL(12,2) NOT NULL,
    seller_amount DECIMAL(12,2) NOT NULL,
    
    status VARCHAR(50) DEFAULT 'held', -- held, released, refunded
    release_date TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT escrow_payments_amount_positive CHECK (amount > 0),
    CONSTRAINT escrow_payments_status_valid CHECK (
        status IN ('held', 'released', 'refunded')
    ),
    CONSTRAINT escrow_payments_amounts_sum CHECK (
        marketplace_commission + seller_amount = amount
    )
);

-- Таблица выплат продавцам
CREATE TABLE merchant_payouts (
    id BIGSERIAL PRIMARY KEY,
    seller_id INT REFERENCES users(id),
    gateway_id INT REFERENCES payment_gateways(id),
    
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    
    -- AllSecure payout данные
    gateway_payout_id VARCHAR(255),
    gateway_reference_id VARCHAR(255),
    
    status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
    
    -- Банковские данные получателя
    bank_account_info JSONB,
    
    -- Metadata
    gateway_response JSONB,
    error_details JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT merchant_payouts_amount_positive CHECK (amount > 0),
    CONSTRAINT merchant_payouts_status_valid CHECK (
        status IN ('pending', 'processing', 'completed', 'failed')
    )
);

-- Индексы для производительности
CREATE INDEX idx_payment_transactions_user_id ON payment_transactions(user_id);
CREATE INDEX idx_payment_transactions_listing_id ON payment_transactions(listing_id);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);
CREATE INDEX idx_payment_transactions_gateway_transaction_id ON payment_transactions(gateway_transaction_id);
CREATE INDEX idx_payment_transactions_order_reference ON payment_transactions(order_reference);
CREATE INDEX idx_payment_transactions_created_at ON payment_transactions(created_at);

CREATE INDEX idx_escrow_payments_seller_id ON escrow_payments(seller_id);
CREATE INDEX idx_escrow_payments_buyer_id ON escrow_payments(buyer_id);
CREATE INDEX idx_escrow_payments_status ON escrow_payments(status);
CREATE INDEX idx_escrow_payments_payment_transaction_id ON escrow_payments(payment_transaction_id);

CREATE INDEX idx_merchant_payouts_seller_id ON merchant_payouts(seller_id);
CREATE INDEX idx_merchant_payouts_status ON merchant_payouts(status);
CREATE INDEX idx_merchant_payouts_gateway_payout_id ON merchant_payouts(gateway_payout_id);

-- Проверяем существование таблицы user_transactions и добавляем поддержку AllSecure если она существует
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_transactions') THEN
        -- Добавляем колонки если таблица существует
        ALTER TABLE user_transactions ADD COLUMN IF NOT EXISTS payment_gateway VARCHAR(50);
        ALTER TABLE user_transactions ADD COLUMN IF NOT EXISTS gateway_transaction_id VARCHAR(255);
        ALTER TABLE user_transactions ADD COLUMN IF NOT EXISTS gateway_reference VARCHAR(255);
        
        -- Создаем индексы
        CREATE INDEX IF NOT EXISTS idx_user_transactions_gateway_transaction_id ON user_transactions(gateway_transaction_id);
        CREATE INDEX IF NOT EXISTS idx_user_transactions_payment_gateway ON user_transactions(payment_gateway);
    END IF;
END $$;

-- Добавляем базовую конфигурацию AllSecure
INSERT INTO payment_gateways (name, is_active, config) VALUES (
    'allsecure',
    false, -- изначально неактивен, будет активирован после настройки
    '{
        "base_url": "https://asxgw.com",
        "sandbox_mode": true,
        "supported_currencies": ["RSD", "EUR", "USD"],
        "marketplace_commission_rate": 0.05,
        "escrow_release_days": 7
    }'::jsonb
);

-- Добавляем комментарии к таблицам
COMMENT ON TABLE payment_gateways IS 'Конфигурация платежных шлюзов';
COMMENT ON TABLE payment_transactions IS 'Платежные транзакции через внешние шлюзы';
COMMENT ON TABLE escrow_payments IS 'Эскроу платежи для marketplace (удержание средств)';
COMMENT ON TABLE merchant_payouts IS 'Выплаты продавцам через платежные шлюзы';

-- Комментарии к ключевым полям
COMMENT ON COLUMN payment_transactions.order_reference IS 'Уникальная ссылка на заказ для отслеживания';
COMMENT ON COLUMN payment_transactions.gateway_transaction_id IS 'UUID транзакции в AllSecure';
COMMENT ON COLUMN payment_transactions.gateway_reference_id IS 'Purchase ID в AllSecure';
COMMENT ON COLUMN escrow_payments.release_date IS 'Дата автоматического освобождения средств';
COMMENT ON COLUMN merchant_payouts.bank_account_info IS 'IBAN и другие банковские данные для выплаты';