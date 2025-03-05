-- Создание таблицы для витрин магазинов
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