-- 000045_create_marketplace_listing_variants.up.sql

-- Создание таблицы для вариантов объявлений маркетплейса
CREATE TABLE IF NOT EXISTS marketplace_listing_variants (
    id SERIAL PRIMARY KEY,
    listing_id INTEGER NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2),
    stock INTEGER DEFAULT 0,
    attributes JSONB NOT NULL DEFAULT '{}', -- Атрибуты варианта (размер, цвет и т.д.)
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индексов для оптимизации поиска
CREATE INDEX IF NOT EXISTS idx_marketplace_listing_variants_listing_id ON marketplace_listing_variants(listing_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_listing_variants_sku ON marketplace_listing_variants(sku);
CREATE INDEX IF NOT EXISTS idx_marketplace_listing_variants_attributes ON marketplace_listing_variants USING GIN(attributes);
CREATE INDEX IF NOT EXISTS idx_marketplace_listing_variants_is_active ON marketplace_listing_variants(is_active);

-- Добавление уникального ограничения на SKU в пределах объявления
ALTER TABLE marketplace_listing_variants ADD CONSTRAINT uk_marketplace_listing_variants_sku_per_listing 
UNIQUE (listing_id, sku);

-- Создание триггера для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_marketplace_listing_variants_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER marketplace_listing_variants_updated_at
    BEFORE UPDATE ON marketplace_listing_variants
    FOR EACH ROW
    EXECUTE FUNCTION update_marketplace_listing_variants_updated_at();