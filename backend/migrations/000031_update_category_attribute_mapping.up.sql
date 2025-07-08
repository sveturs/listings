-- Миграция для обновления таблицы связей атрибутов категорий
-- Добавляем поле sort_order в category_attribute_mapping, если оно еще не существует
ALTER TABLE category_attribute_mapping 
    ADD COLUMN IF NOT EXISTS sort_order INT DEFAULT 0;

-- Добавляем индекс для ускорения поиска атрибутов по категории
CREATE INDEX IF NOT EXISTS idx_category_attribute_map_cat_id ON category_attribute_mapping(category_id);

-- Добавляем индекс для ускорения поиска атрибутов по ID атрибута
CREATE INDEX IF NOT EXISTS idx_category_attribute_map_attr_id ON category_attribute_mapping(attribute_id);

-- Добавляем поле под кастомный компонент ввода для конкретной категории
ALTER TABLE category_attribute_mapping
    ADD COLUMN IF NOT EXISTS custom_component VARCHAR(255);

-- Обновляем существующие записи для установки порядка сортировки
-- Это важно, чтобы атрибуты отображались в правильном порядке в интерфейсе
UPDATE category_attribute_mapping cam
SET sort_order = ca.sort_order
FROM category_attributes ca
WHERE cam.attribute_id = ca.id AND cam.sort_order = 0;

-- Создаем функцию для автоматического обновления порядка сортировки при добавлении новых атрибутов
CREATE OR REPLACE FUNCTION update_category_attribute_sort_order() RETURNS TRIGGER AS $$
BEGIN
    -- Если sort_order не указан, берем его из атрибута
    IF NEW.sort_order = 0 THEN
        SELECT sort_order INTO NEW.sort_order 
        FROM category_attributes 
        WHERE id = NEW.attribute_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер для автоматического обновления порядка сортировки
DROP TRIGGER IF EXISTS tr_update_category_attribute_sort_order ON category_attribute_mapping;
DROP TRIGGER IF EXISTS tr_update_category_attribute_sort_order ON tr_update_category_attribute_sort_order;
CREATE TRIGGER tr_update_category_attribute_sort_order
BEFORE INSERT ON category_attribute_mapping
FOR EACH ROW
EXECUTE FUNCTION update_category_attribute_sort_order();

-- Создаем или обновляем представление для получения полной информации об атрибутах категорий
CREATE OR REPLACE VIEW v_category_attributes AS
SELECT 
    cam.category_id,
    cam.attribute_id,
    cam.is_enabled,
    cam.is_required,
    cam.sort_order,
    cam.custom_component,
    ca.name,
    ca.display_name,
    ca.attribute_type,
    ca.options,
    ca.validation_rules,
    ca.is_searchable,
    ca.is_filterable,
    ca.custom_component as default_custom_component,
    mc.name as category_name,
    mc.slug as category_slug
FROM 
    category_attribute_mapping cam
JOIN 
    category_attributes ca ON cam.attribute_id = ca.id
JOIN 
    marketplace_categories mc ON cam.category_id = mc.id
ORDER BY 
    cam.category_id, cam.sort_order, ca.sort_order;

-- Комментарии к представлению
COMMENT ON VIEW v_category_attributes IS 'Представление для получения полной информации об атрибутах категорий';