-- Откат миграции: удаление поля affects_stock
ALTER TABLE category_attributes 
DROP COLUMN IF EXISTS affects_stock;