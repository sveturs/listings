-- Привязка атрибутов к категориям

-- Общие атрибуты для всех категорий
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    c.id as category_id,
    a.id as attribute_id,
    true as is_enabled,
    CASE 
        WHEN a.name = 'price' THEN true
        ELSE false
    END as is_required,
    a.sort_order
FROM 
    marketplace_categories c,
    category_attributes a
WHERE 
    a.name IN ('price', 'condition', 'brand', 'color')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping cam 
        WHERE cam.category_id = c.id AND cam.attribute_id = a.id
    );

-- Атрибуты для электроники (category_id = 1001 и подкатегории)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    c.id as category_id,
    a.id as attribute_id,
    true as is_enabled,
    false as is_required,
    a.sort_order + 10
FROM 
    marketplace_categories c,
    category_attributes a
WHERE 
    (c.id = 1001 OR c.parent_id = 1001)
    AND a.name IN ('storage', 'operating_system', 'processor', 'ram', 'storage_type')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping cam 
        WHERE cam.category_id = c.id AND cam.attribute_id = a.id
    );

-- Атрибуты для моды (category_id = 1002 и подкатегории)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    c.id as category_id,
    a.id as attribute_id,
    true as is_enabled,
    false as is_required,
    a.sort_order + 10
FROM 
    marketplace_categories c,
    category_attributes a
WHERE 
    (c.id = 1002 OR c.parent_id = 1002)
    AND a.name IN ('gender', 'clothing_size', 'shoe_size')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping cam 
        WHERE cam.category_id = c.id AND cam.attribute_id = a.id
    );

-- Атрибуты для автомобилей (category_id = 1003 и подкатегории)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required, sort_order)
SELECT 
    c.id as category_id,
    a.id as attribute_id,
    true as is_enabled,
    false as is_required,
    a.sort_order + 10
FROM 
    marketplace_categories c,
    category_attributes a
WHERE 
    (c.id = 1003 OR c.parent_id = 1003)
    AND a.name IN ('fuel_type', 'transmission')
    AND NOT EXISTS (
        SELECT 1 FROM category_attribute_mapping cam 
        WHERE cam.category_id = c.id AND cam.attribute_id = a.id
    );