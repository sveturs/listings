-- Добавление поля affects_stock в таблицу category_attributes
-- Это поле указывает, влияет ли атрибут на учет остатков товара
ALTER TABLE category_attributes 
ADD COLUMN affects_stock BOOLEAN DEFAULT FALSE;

-- Комментарий для поля
COMMENT ON COLUMN category_attributes.affects_stock IS 'Указывает, влияет ли данный атрибут на учет остатков товара (например, размер или цвет могут влиять на остатки)';

-- Обновляем существующие атрибуты, которые обычно влияют на остатки
UPDATE category_attributes 
SET affects_stock = TRUE 
WHERE LOWER(name) IN ('size', 'color', 'memory', 'storage', 'ram', 'capacity');