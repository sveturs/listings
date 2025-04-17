-- Объединение миграций add_user_balance и create_storefronts
-- Сначала создаем все таблицы для баланса пользователей

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
    ('IPS QR code', 'ips_qr', 'digital', true, 100, 1000000, 0.8, 0)
ON CONFLICT (code) DO NOTHING;  


-- Добавляем начальный баланс для пользователей
INSERT INTO user_balances (user_id, balance, frozen_balance, currency, updated_at) 
VALUES (2, 15000000.00, 0.00, 'RSD', CURRENT_TIMESTAMP)
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_balances (user_id, balance, frozen_balance, currency, updated_at) 
VALUES (3, 1500000.00, 0.00, 'RSD', CURRENT_TIMESTAMP)
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_balances (user_id, balance, frozen_balance, currency, updated_at) 
VALUES (4, 150000.00, 0.00, 'RSD', CURRENT_TIMESTAMP)
ON CONFLICT (user_id) DO NOTHING;

-- Теперь создаем таблицы для витрин магазинов
CREATE TABLE IF NOT EXISTS user_storefronts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo_path VARCHAR(255),
    slug VARCHAR(100) UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    creation_transaction_id INT REFERENCES balance_transactions(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для настроек импорта
CREATE TABLE IF NOT EXISTS import_sources (
    id SERIAL PRIMARY KEY,
    storefront_id INT NOT NULL REFERENCES user_storefronts(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- csv, xml, json
    url VARCHAR(512),
    auth_data JSONB,
    schedule VARCHAR(50), -- cron-like schedule
    mapping JSONB, -- маппинг полей
    last_import_at TIMESTAMP,
    last_import_status VARCHAR(20),
    last_import_log TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения истории импорта
CREATE TABLE IF NOT EXISTS import_history (
    id SERIAL PRIMARY KEY,
    source_id INT NOT NULL REFERENCES import_sources(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    items_total INT DEFAULT 0,
    items_imported INT DEFAULT 0,
    items_failed INT DEFAULT 0,
    log TEXT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_user_storefronts_user ON user_storefronts(user_id);
CREATE INDEX IF NOT EXISTS idx_user_storefronts_status ON user_storefronts(status);
CREATE INDEX IF NOT EXISTS idx_import_sources_storefront ON import_sources(storefront_id);
CREATE INDEX IF NOT EXISTS idx_import_history_source ON import_history(source_id);

-- Добавляем поле для ссылки на витрину в таблицу объявлений
ALTER TABLE marketplace_listings ADD COLUMN IF NOT EXISTS storefront_id INT REFERENCES user_storefronts(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_marketplace_listings_storefront ON marketplace_listings(storefront_id);