-- backend/migrations/000003_add_user_balance.up.sql

-- Создание таблиц для баланса
CREATE TABLE IF NOT EXISTS user_balances (
    user_id INT PRIMARY KEY REFERENCES users(id),
    balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS balance_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL, -- deposit, withdrawal, transfer, service_payment
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, completed, failed
    payment_method VARCHAR(50), -- bank_transfer, payment_slip, crypto, etc.
    payment_details JSONB,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL, -- bank, payment_system, crypto
    is_active BOOLEAN DEFAULT true,
    minimum_amount DECIMAL(12,2),
    maximum_amount DECIMAL(12,2),
    fee_percentage DECIMAL(5,2),
    fixed_fee DECIMAL(12,2),
    credentials JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем индексы
CREATE INDEX IF NOT EXISTS idx_transactions_user ON balance_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON balance_transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created ON balance_transactions(created_at);

-- Добавляем методы оплаты по умолчанию
INSERT INTO payment_methods (
    name, code, type, is_active, minimum_amount, maximum_amount, fee_percentage, fixed_fee
) VALUES 
    ('Bank transfer', 'bank_transfer', 'bank', true, 1000, 10000000, 0, 100),
    ('Post office', 'post_office', 'cash', true, 500, 500000, 1.5, 50),
    ('IPS QR code', 'ips_qr', 'digital', true, 100, 1000000, 0.8, 0);