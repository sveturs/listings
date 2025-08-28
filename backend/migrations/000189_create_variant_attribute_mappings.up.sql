-- Создание таблицы для связи вариативных атрибутов с атрибутами категорий
CREATE TABLE IF NOT EXISTS variant_attribute_mappings (
    id SERIAL PRIMARY KEY,
    variant_attribute_id INTEGER NOT NULL,
    category_attribute_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Внешние ключи
    CONSTRAINT fk_variant_attribute 
        FOREIGN KEY (variant_attribute_id) 
        REFERENCES product_variant_attributes(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_category_attribute 
        FOREIGN KEY (category_attribute_id) 
        REFERENCES category_attributes(id) 
        ON DELETE CASCADE,
    
    -- Уникальный индекс для предотвращения дублирования связей
    CONSTRAINT unique_variant_category_mapping 
        UNIQUE (variant_attribute_id, category_attribute_id)
);

-- Индексы для быстрого поиска
CREATE INDEX idx_variant_attribute_mappings_variant_id 
    ON variant_attribute_mappings(variant_attribute_id);
CREATE INDEX idx_variant_attribute_mappings_category_id 
    ON variant_attribute_mappings(category_attribute_id);

-- Триггер для обновления updated_at
CREATE TRIGGER update_variant_attribute_mappings_updated_at
    BEFORE UPDATE ON variant_attribute_mappings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии для документации
COMMENT ON TABLE variant_attribute_mappings IS 'Связи между вариативными атрибутами и атрибутами категорий';
COMMENT ON COLUMN variant_attribute_mappings.variant_attribute_id IS 'ID вариативного атрибута из product_variant_attributes';
COMMENT ON COLUMN variant_attribute_mappings.category_attribute_id IS 'ID атрибута категории из category_attributes';