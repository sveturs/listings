-- Исправляем категории для объявлений о животных, которые ошибочно находятся в категории "Industrijske mašine"
UPDATE marketplace_listings 
SET category_id = 1011, -- Pets category
    updated_at = NOW()
WHERE id IN (262, 263);

-- Добавим комментарий для отслеживания изменений
COMMENT ON COLUMN marketplace_listings.category_id IS 'Category ID for the listing. Fixed listings 262, 263 from incorrect category 1701 to 1011 (Pets)';