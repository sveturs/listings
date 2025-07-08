-- Создаем таблицу для групп атрибутов
CREATE TABLE IF NOT EXISTS attribute_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(100),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_system BOOLEAN DEFAULT false, -- системные группы нельзя удалить
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу для связи атрибутов с группами
CREATE TABLE IF NOT EXISTS attribute_group_items (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES attribute_groups(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES category_attributes(id) ON DELETE CASCADE,
    icon VARCHAR(100), -- иконка для конкретного атрибута в группе
    sort_order INT DEFAULT 0,
    custom_display_name VARCHAR(255), -- переопределенное имя для отображения в группе
    visibility_condition JSONB, -- условия видимости атрибута в группе
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, attribute_id)
);

-- Создаем таблицу для привязки групп к категориям
CREATE TABLE IF NOT EXISTS category_attribute_groups (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES marketplace_categories(id) ON DELETE CASCADE,
    group_id INT NOT NULL REFERENCES attribute_groups(id) ON DELETE CASCADE,
    component_id INT REFERENCES custom_ui_components(id) ON DELETE SET NULL, -- связь с кастомным компонентом
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    display_mode VARCHAR(50) DEFAULT 'list', -- list, grid, accordion, tabs, custom
    collapsed_by_default BOOLEAN DEFAULT false,
    configuration JSONB DEFAULT '{}', -- дополнительная конфигурация для отображения группы
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category_id, group_id)
);

-- Создаем индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_attribute_groups_name ON attribute_groups(name);
CREATE INDEX IF NOT EXISTS idx_attribute_groups_active ON attribute_groups(is_active);
CREATE INDEX IF NOT EXISTS idx_attribute_group_items_group ON attribute_group_items(group_id);
CREATE INDEX IF NOT EXISTS idx_attribute_group_items_attribute ON attribute_group_items(attribute_id);
CREATE INDEX IF NOT EXISTS idx_category_attribute_groups_category ON category_attribute_groups(category_id);
CREATE INDEX IF NOT EXISTS idx_category_attribute_groups_group ON category_attribute_groups(group_id);
CREATE INDEX IF NOT EXISTS idx_category_attribute_groups_component ON category_attribute_groups(component_id);

-- Создаем функцию для обновления updated_at
CREATE OR REPLACE FUNCTION update_attribute_groups_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создаем триггер для автоматического обновления updated_at
DROP TRIGGER IF EXISTS update_attribute_groups_updated_at ON attribute_groups;
DROP TRIGGER IF EXISTS update_attribute_groups_updated_at ON update_attribute_groups_updated_at;
CREATE TRIGGER update_attribute_groups_updated_at BEFORE UPDATE ON attribute_groups
    FOR EACH ROW EXECUTE FUNCTION update_attribute_groups_updated_at();

-- Добавляем примеры групп атрибутов
INSERT INTO attribute_groups (name, display_name, description, icon, sort_order, is_system) VALUES
('basic_info', 'Основная информация', 'Базовые характеристики товара', 'info', 1, true),
('technical_specs', 'Технические характеристики', 'Технические параметры товара', 'settings', 2, true),
('dimensions', 'Размеры и вес', 'Физические параметры товара', 'straighten', 3, true),
('appearance', 'Внешний вид', 'Цвет, материал и другие визуальные характеристики', 'palette', 4, true),
('condition_details', 'Состояние и комплектация', 'Детали о состоянии товара и его комплектации', 'inventory_2', 5, true),
('additional_info', 'Дополнительная информация', 'Прочие характеристики и особенности', 'more_horiz', 6, true);

-- Добавляем представление для удобного получения групп с атрибутами
CREATE OR REPLACE VIEW v_attribute_groups_with_items AS
SELECT 
    ag.id AS group_id,
    ag.name AS group_name,
    ag.display_name AS group_display_name,
    ag.icon AS group_icon,
    ag.sort_order AS group_sort_order,
    agi.id AS item_id,
    agi.attribute_id,
    ca.name AS attribute_name,
    ca.display_name AS attribute_display_name,
    agi.icon AS attribute_icon,
    agi.custom_display_name,
    agi.sort_order AS attribute_sort_order
FROM attribute_groups ag
LEFT JOIN attribute_group_items agi ON ag.id = agi.group_id
LEFT JOIN category_attributes ca ON agi.attribute_id = ca.id
WHERE ag.is_active = true
ORDER BY ag.sort_order, agi.sort_order;