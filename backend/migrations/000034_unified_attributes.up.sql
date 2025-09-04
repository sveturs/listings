-- Миграция для создания унифицированной системы атрибутов
-- ВАЖНО: Создаем новые таблицы параллельно существующим, без удаления старых!
-- Дата: 02.09.2025

-- =====================================================
-- 1. УНИФИЦИРОВАННАЯ ТАБЛИЦА АТРИБУТОВ
-- =====================================================
CREATE TABLE IF NOT EXISTS unified_attributes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) UNIQUE NOT NULL, -- Уникальный код атрибута
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    attribute_type VARCHAR(50) NOT NULL CHECK (attribute_type IN (
        'text', 'textarea', 'number', 'boolean', 
        'select', 'multiselect', 'date', 'color', 'size'
    )),
    purpose VARCHAR(20) NOT NULL DEFAULT 'regular' CHECK (purpose IN (
        'regular',    -- Обычный атрибут для фильтрации/поиска
        'variant',    -- Вариативный атрибут (влияет на SKU)
        'both'        -- Может использоваться в обоих случаях
    )),
    
    -- Настройки атрибута
    options JSONB DEFAULT '{}',           -- Опции для select/multiselect
    validation_rules JSONB DEFAULT '{}',  -- Правила валидации
    ui_settings JSONB DEFAULT '{}',       -- Настройки отображения
    
    -- Флаги использования
    is_searchable BOOLEAN DEFAULT false,
    is_filterable BOOLEAN DEFAULT false,
    is_required BOOLEAN DEFAULT false,
    affects_stock BOOLEAN DEFAULT false,  -- Для вариативных атрибутов
    affects_price BOOLEAN DEFAULT false,  -- Для вариативных атрибутов
    
    -- Метаданные
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Связь со старой системой (временно, для миграции)
    legacy_category_attribute_id INTEGER,
    legacy_product_variant_attribute_id INTEGER
);

-- =====================================================
-- 2. СВЯЗЬ АТРИБУТОВ С КАТЕГОРИЯМИ
-- =====================================================
CREATE TABLE IF NOT EXISTS unified_category_attributes (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
    
    -- Настройки для конкретной категории
    is_enabled BOOLEAN DEFAULT true,
    is_required BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    
    -- Переопределение настроек атрибута для категории
    category_specific_options JSONB,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(category_id, attribute_id)
);

-- =====================================================
-- 3. ЗНАЧЕНИЯ АТРИБУТОВ (универсальная таблица)
-- =====================================================
CREATE TABLE IF NOT EXISTS unified_attribute_values (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN (
        'listing',           -- Объявление маркетплейса
        'product',           -- Товар витрины
        'product_variant'    -- Вариант товара
    )),
    entity_id INTEGER NOT NULL,
    attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
    
    -- Значения разных типов
    text_value TEXT,
    numeric_value NUMERIC,
    boolean_value BOOLEAN,
    date_value DATE,
    json_value JSONB,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Индекс для быстрого поиска
    UNIQUE(entity_type, entity_id, attribute_id)
);

-- =====================================================
-- 4. ИНДЕКСЫ ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ
-- =====================================================
CREATE INDEX idx_unified_attributes_code ON unified_attributes(code);
CREATE INDEX idx_unified_attributes_purpose ON unified_attributes(purpose);
CREATE INDEX idx_unified_attributes_active ON unified_attributes(is_active);

CREATE INDEX idx_unified_category_attributes_category ON unified_category_attributes(category_id);
CREATE INDEX idx_unified_category_attributes_enabled ON unified_category_attributes(is_enabled);

CREATE INDEX idx_unified_attribute_values_entity ON unified_attribute_values(entity_type, entity_id);
CREATE INDEX idx_unified_attribute_values_attribute ON unified_attribute_values(attribute_id);
CREATE INDEX idx_unified_attribute_values_text ON unified_attribute_values(text_value) WHERE text_value IS NOT NULL;
CREATE INDEX idx_unified_attribute_values_numeric ON unified_attribute_values(numeric_value) WHERE numeric_value IS NOT NULL;

-- =====================================================
-- 5. ТРИГГЕР ДЛЯ ОБНОВЛЕНИЯ updated_at
-- =====================================================
CREATE OR REPLACE FUNCTION update_unified_attributes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_unified_attributes_updated_at
    BEFORE UPDATE ON unified_attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_unified_attributes_updated_at();

CREATE TRIGGER update_unified_category_attributes_updated_at
    BEFORE UPDATE ON unified_category_attributes
    FOR EACH ROW
    EXECUTE FUNCTION update_unified_attributes_updated_at();

CREATE TRIGGER update_unified_attribute_values_updated_at
    BEFORE UPDATE ON unified_attribute_values
    FOR EACH ROW
    EXECUTE FUNCTION update_unified_attributes_updated_at();

-- =====================================================
-- 6. КОММЕНТАРИИ К ТАБЛИЦАМ
-- =====================================================
COMMENT ON TABLE unified_attributes IS 'Унифицированная таблица атрибутов для всей системы';
COMMENT ON TABLE unified_category_attributes IS 'Связь атрибутов с категориями и их настройки';
COMMENT ON TABLE unified_attribute_values IS 'Значения атрибутов для всех сущностей системы';

COMMENT ON COLUMN unified_attributes.purpose IS 'regular - обычный атрибут, variant - вариативный, both - универсальный';
COMMENT ON COLUMN unified_attribute_values.entity_type IS 'Тип сущности: listing, product, product_variant';