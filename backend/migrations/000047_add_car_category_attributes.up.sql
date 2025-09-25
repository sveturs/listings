-- Добавление атрибутов для категории "Lični automobili" (id: 1301)
-- и родительской категории "Automobili" (id: 1003)

-- Удаляем существующие маппинги если они есть
DELETE FROM category_attribute_mapping WHERE category_id IN (1003, 1301);

-- Основные атрибуты для автомобилей (категория 1301 - Lični automobili)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order) VALUES
-- Основные атрибуты
(1301, 113, true, 1),   -- brand (Brend) - обязательное
(1301, 91, true, 2),    -- car_model (Model) - обязательное
(1301, 86, true, 3),    -- year (Godište) - обязательное
(1301, 87, false, 4),   -- mileage (Kilometraža)
(1301, 140, false, 5),  -- body_type (Tip karoserije)
(1301, 149, false, 6),  -- fuel_type (Gorivo)
(1301, 148, false, 7);  -- transmission (Menjač)

-- Те же атрибуты для родительской категории (1003 - Automobili)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order) VALUES
(1003, 113, true, 1),   -- brand (Brend) - обязательное
(1003, 91, true, 2),    -- car_model (Model) - обязательное
(1003, 86, true, 3),    -- year (Godište) - обязательное
(1003, 87, false, 4),   -- mileage (Kilometraža)
(1003, 140, false, 5),  -- body_type (Tip karoserije)
(1003, 149, false, 6),  -- fuel_type (Gorivo)
(1003, 148, false, 7);  -- transmission (Menjač)

-- Добавляем атрибуты для всех подкатегорий автомобилей
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_required, sort_order)
SELECT DISTINCT mc.id, cam.attribute_id, cam.is_required, cam.sort_order
FROM marketplace_categories mc
CROSS JOIN category_attribute_mapping cam
WHERE mc.parent_id = 1301  -- все подкатегории Lični automobili
  AND cam.category_id = 1301
  AND NOT EXISTS (
    SELECT 1 FROM category_attribute_mapping cam2
    WHERE cam2.category_id = mc.id AND cam2.attribute_id = cam.attribute_id
  );