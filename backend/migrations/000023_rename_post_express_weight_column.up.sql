-- Переименование колонки weight_kg в weight для унификации с другими таблицами доставки
ALTER TABLE post_express_shipments RENAME COLUMN weight_kg TO weight;

-- Обновление комментария для ясности
COMMENT ON COLUMN post_express_shipments.weight IS 'Weight of the package in kilograms';