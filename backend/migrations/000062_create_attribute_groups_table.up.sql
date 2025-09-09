-- Создание таблицы attribute_groups для группировки атрибутов

CREATE TABLE IF NOT EXISTS attribute_groups (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_attribute_groups_code ON attribute_groups(code);
CREATE INDEX IF NOT EXISTS idx_attribute_groups_is_active ON attribute_groups(is_active);
CREATE INDEX IF NOT EXISTS idx_attribute_groups_sort_order ON attribute_groups(sort_order);

-- Добавление поля group_id в таблицу attributes если его нет
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'attributes' 
        AND column_name = 'group_id'
    ) THEN
        ALTER TABLE attributes ADD COLUMN group_id INTEGER REFERENCES attribute_groups(id);
        CREATE INDEX idx_attributes_group_id ON attributes(group_id);
    END IF;
END $$;

-- Вставка базовых групп атрибутов
INSERT INTO attribute_groups (code, name, description, sort_order) VALUES
    ('basic', 'Основные характеристики', 'Основные характеристики товара', 1),
    ('technical', 'Технические характеристики', 'Технические параметры и спецификации', 2),
    ('dimensions', 'Размеры и вес', 'Физические размеры и вес товара', 3),
    ('appearance', 'Внешний вид', 'Цвет, материал и другие визуальные характеристики', 4),
    ('features', 'Особенности', 'Дополнительные функции и возможности', 5),
    ('package', 'Комплектация', 'Комплект поставки', 6)
ON CONFLICT (code) DO NOTHING;