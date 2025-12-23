-- Brand category mapping table
-- Связывает бренды с категориями для быстрой детекции

CREATE TABLE IF NOT EXISTS brand_category_mapping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_name VARCHAR(100) NOT NULL,
    brand_aliases TEXT[] DEFAULT '{}',
    category_slug VARCHAR(100) NOT NULL,
    confidence FLOAT NOT NULL DEFAULT 0.95,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_brand_category UNIQUE (brand_name, category_slug)
);

-- Индекс для быстрого поиска по бренду (case-insensitive)
CREATE INDEX idx_brand_mapping_name ON brand_category_mapping(LOWER(brand_name));

-- GIN индекс для поиска по алиасам
CREATE INDEX idx_brand_mapping_aliases ON brand_category_mapping USING GIN(brand_aliases);

-- Индекс для фильтрации верифицированных брендов
CREATE INDEX idx_brand_mapping_verified ON brand_category_mapping(is_verified) WHERE is_verified = true;

-- Комментарии
COMMENT ON TABLE brand_category_mapping IS 'Маппинг брендов на категории для автоматической детекции';
COMMENT ON COLUMN brand_category_mapping.brand_name IS 'Основное название бренда (нормализованное)';
COMMENT ON COLUMN brand_category_mapping.brand_aliases IS 'Альтернативные написания бренда (массив)';
COMMENT ON COLUMN brand_category_mapping.category_slug IS 'Slug категории из таблицы categories';
COMMENT ON COLUMN brand_category_mapping.confidence IS 'Уверенность маппинга (0.0-1.0, default 0.95)';
COMMENT ON COLUMN brand_category_mapping.is_verified IS 'Подтверждён ли маппинг вручную';
