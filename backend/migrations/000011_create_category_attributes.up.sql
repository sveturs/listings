-- backend/migrations/000011_create_category_attributes.up.sql
-- Таблица для определения атрибутов категорий
CREATE TABLE category_attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    attribute_type VARCHAR(50) NOT NULL, -- text, number, select, boolean, etc.
    options JSONB, -- Для хранения вариантов выбора, диапазонов и т.д.
    validation_rules JSONB, -- Правила валидации
    is_searchable BOOLEAN DEFAULT true,
    is_filterable BOOLEAN DEFAULT true,
    is_required BOOLEAN DEFAULT false,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для связи атрибутов с категориями
CREATE TABLE category_attribute_mapping (
    category_id INT NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES category_attributes(id) ON DELETE CASCADE,
    is_enabled BOOLEAN DEFAULT true,
    is_required BOOLEAN DEFAULT false,
    PRIMARY KEY (category_id, attribute_id)
);

-- Таблица для хранения значений атрибутов объявлений
-- ИЗМЕНЕНО: добавлен serial id как PRIMARY KEY вместо композитного ключа
CREATE TABLE listing_attribute_values (
    id SERIAL PRIMARY KEY,  -- Добавлен уникальный идентификатор
    listing_id INT NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES category_attributes(id) ON DELETE CASCADE,
    text_value TEXT,
    numeric_value DECIMAL(15,5),
    boolean_value BOOLEAN,
    json_value JSONB
);

-- Создаем индекс для быстрого поиска атрибутов конкретного объявления
CREATE INDEX idx_listing_attr_listing_id ON listing_attribute_values(listing_id);

-- Индексы для ускорения поиска по атрибутам
CREATE INDEX idx_listing_attr_text ON listing_attribute_values(attribute_id, text_value) WHERE text_value IS NOT NULL;
CREATE INDEX idx_listing_attr_numeric ON listing_attribute_values(attribute_id, numeric_value) WHERE numeric_value IS NOT NULL;
CREATE INDEX idx_listing_attr_boolean ON listing_attribute_values(attribute_id, boolean_value) WHERE boolean_value IS NOT NULL;