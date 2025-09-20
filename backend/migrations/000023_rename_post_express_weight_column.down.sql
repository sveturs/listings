-- Откат: возвращение колонки weight обратно к weight_kg
ALTER TABLE post_express_shipments RENAME COLUMN weight TO weight_kg;

-- Удаление комментария
COMMENT ON COLUMN post_express_shipments.weight_kg IS NULL;